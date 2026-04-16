package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func makeServer(t *testing.T, body string, status int) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(status)
		fmt.Fprint(w, body)
	}))
}

func TestRun_Success(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"data":{"data":{"API_KEY":"abc123"}}}`))
	}))
	defer ts.Close()

	tmpDir := t.TempDir()
	outFile := filepath.Join(tmpDir, ".env")

	os.Setenv("VAULT_ADDR", ts.URL)
	os.Setenv("VAULT_TOKEN", "test-token")
	os.Setenv("VAULT_SECRET_PATH", "secret/data/myapp")
	os.Setenv("OUTPUT_FILE", outFile)
	defer func() {
		os.Unsetenv("VAULT_ADDR")
		os.Unsetenv("VAULT_TOKEN")
		os.Unsetenv("VAULT_SECRET_PATH")
		os.Unsetenv("OUTPUT_FILE")
	}()

	if err := run(); err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	data, err := os.ReadFile(outFile)
	if err != nil {
		t.Fatalf("could not read output file: %v", err)
	}
	if len(data) == 0 {
		t.Error("expected non-empty output file")
	}
}

func TestRun_MissingToken(t *testing.T) {
	os.Unsetenv("VAULT_TOKEN")
	os.Setenv("VAULT_SECRET_PATH", "secret/data/myapp")
	defer os.Unsetenv("VAULT_SECRET_PATH")

	if err := run(); err == nil {
		t.Fatal("expected error for missing token")
	}
}
