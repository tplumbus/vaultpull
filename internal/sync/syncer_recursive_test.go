package sync

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"vaultpull/internal/vault"
)

func makeRecursiveVaultServer(t *testing.T) *httptest.Server {
	t.Helper()
	mux := http.NewServeMux()
	mux.HandleFunc("/v1/secret/", func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == "LIST":
			json.NewEncoder(w).Encode(map[string]interface{}{"data": map[string]interface{}{"keys": []string{"db"}}})
		case r.Method == http.MethodGet && strings.HasSuffix(r.URL.Path, "/data/db"):
			json.NewEncoder(w).Encode(map[string]interface{}{"data": map[string]interface{}{"data": map[string]string{"password": "s3cr3t"}}})
		default:
			http.NotFound(w, r)
		}
	})
	return httptest.NewServer(mux)
}

func TestRunRecursive_WritesSecrets(t *testing.T) {
	srv := makeRecursiveVaultServer(t)
	defer srv.Close()

	dir := t.TempDir()
	out := filepath.Join(dir, ".env")

	c, err := vault.NewClient(srv.URL, "tok", vault.KVv2)
	if err != nil {
		t.Fatal(err)
	}

	s := New(c, "secret", "", out, false)
	if err := s.RunRecursive(context.Background()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, err := os.ReadFile(out)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(data), "DB_PASSWORD") {
		t.Errorf("expected DB_PASSWORD in output, got:\n%s", data)
	}
}

func TestPathToEnvPrefix(t *testing.T) {
	cases := []struct{ in, want string }{
		{"app/config", "APP_CONFIG"},
		{"db", "DB"},
		{"/app/sub/key", "APP_SUB_KEY"},
	}
	for _, tc := range cases {
		got := pathToEnvPrefix(tc.in)
		if got != tc.want {
			t.Errorf("pathToEnvPrefix(%q) = %q, want %q", tc.in, got, tc.want)
		}
	}
}
