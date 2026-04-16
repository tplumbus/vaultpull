package sync

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/vaultpull/internal/config"
)

func makeVaultServer(t *testing.T, data map[string]interface{}, statusCode int) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(statusCode)
		if data != nil {
			_ = json.NewEncoder(w).Encode(map[string]interface{}{"data": data})
		}
	}))
}

func TestRun_WritesSecrets(t *testing.T) {
	server := makeVaultServer(t, map[string]interface{}{"KEY": "value", "FOO": "bar"}, http.StatusOK)
	defer server.Close()

	tmpDir := t.TempDir()
	outFile := filepath.Join(tmpDir, ".env")

	cfg := &config.Config{
		VaultAddr:  server.URL,
		VaultToken: "test-token",
		VaultPath:  "secret/data/app",
		OutputFile: outFile,
	}

	s, err := New(cfg)
	if err != nil {
		t.Fatalf("New() error: %v", err)
	}

	count, err := s.Run()
	if err != nil {
		t.Fatalf("Run() error: %v", err)
	}
	if count != 2 {
		t.Errorf("expected 2 secrets written, got %d", count)
	}

	if _, err := os.Stat(outFile); os.IsNotExist(err) {
		t.Errorf("expected output file to exist at %s", outFile)
	}
}

func TestRun_EmptySecrets(t *testing.T) {
	server := makeVaultServer(t, map[string]interface{}{}, http.StatusOK)
	defer server.Close()

	tmpDir := t.TempDir()
	cfg := &config.Config{
		VaultAddr:  server.URL,
		VaultToken: "test-token",
		VaultPath:  "secret/data/empty",
		OutputFile: filepath.Join(tmpDir, ".env"),
	}

	s, err := New(cfg)
	if err != nil {
		t.Fatalf("New() error: %v", err)
	}

	count, err := s.Run()
	if err != nil {
		t.Fatalf("Run() error: %v", err)
	}
	if count != 0 {
		t.Errorf("expected 0, got %d", count)
	}
}
