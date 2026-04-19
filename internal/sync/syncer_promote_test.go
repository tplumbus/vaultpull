package sync

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"vaultpull/internal/vault"
)

func makePromoteVaultServer(t *testing.T) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			// Source secret read
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]any{
				"data": map[string]any{
					"data": map[string]string{
						"API_KEY": "staging-key",
						"DB_PASS": "staging-pass",
					},
				},
			})
		case http.MethodPost, http.MethodPut:
			// Destination secret write
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]any{"data": map[string]any{}})
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}))
}

func TestRunPromote_Success(t *testing.T) {
	server := makePromoteVaultServer(t)
	defer server.Close()

	client, err := vault.NewClient(server.URL, "test-token")
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}

	tmpDir := t.TempDir()
	outPath := filepath.Join(tmpDir, "promoted.env")

	s := New(client, outPath)

	err = RunPromote(s, "secret/staging/app", "secret/production/app", false)
	if err != nil {
		t.Fatalf("RunPromote: %v", err)
	}
}

func TestRunPromote_WritesEnvFile(t *testing.T) {
	server := makePromoteVaultServer(t)
	defer server.Close()

	client, err := vault.NewClient(server.URL, "test-token")
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}

	tmpDir := t.TempDir()
	outPath := filepath.Join(tmpDir, "promoted.env")

	s := New(client, outPath)

	err = RunPromote(s, "secret/staging/app", "secret/production/app", false)
	if err != nil {
		t.Fatalf("RunPromote: %v", err)
	}

	data, err := os.ReadFile(outPath)
	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}

	content := string(data)
	if len(content) == 0 {
		t.Error("expected non-empty env file after promote")
	}
}

func TestRunPromote_OverwriteBlocked(t *testing.T) {
	// Server returns 403 on write when overwrite is blocked
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]any{
				"data": map[string]any{
					"data": map[string]string{"KEY": "val"},
				},
			})
			return
		}
		w.WriteHeader(http.StatusForbidden)
	}))
	defer server.Close()

	client, err := vault.NewClient(server.URL, "test-token")
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}

	tmpDir := t.TempDir()
	s := New(client, filepath.Join(tmpDir, "out.env"))

	err = RunPromote(s, "secret/staging/app", false)
	if err == nil {
		t.Error("expected error when promote write is blocked, got nil")
	}
}
