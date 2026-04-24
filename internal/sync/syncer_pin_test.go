package sync_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	internalsync "github.com/yourusername/vaultpull/internal/sync"
	"github.com/yourusername/vaultpull/internal/vault"
)

func makePinVaultServer(t *testing.T, pinnedVersion string) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost, http.MethodPatch:
			w.WriteHeader(http.StatusNoContent)
		case http.MethodGet:
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{
				"data": map[string]interface{}{
					"custom_metadata": map[string]interface{}{
						"pinned_version": pinnedVersion,
					},
					"current_version": float64(3),
					"versions":        map[string]interface{}{},
				},
			})
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}))
}

func TestRunPin_Success(t *testing.T) {
	srv := makePinVaultServer(t, "")
	defer srv.Close()

	client, _ := vault.NewClient(vault.Config{Address: srv.URL, Token: "test"})
	syncer := internalsync.New(client, "")

	var buf bytes.Buffer
	err := syncer.RunPin(internalsync.PinOptions{
		Path:    "secret/myapp/db",
		Version: 2,
		Output:  &buf,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "pinned") {
		t.Errorf("expected 'pinned' in output, got: %s", buf.String())
	}
}

func TestRunPin_MissingPath(t *testing.T) {
	srv := makePinVaultServer(t, "")
	defer srv.Close()

	client, _ := vault.NewClient(vault.Config{Address: srv.URL, Token: "test"})
	syncer := internalsync.New(client, "")

	err := syncer.RunPin(internalsync.PinOptions{Version: 1})
	if err == nil {
		t.Fatal("expected error for missing path")
	}
}

func TestRunPin_InvalidVersion(t *testing.T) {
	srv := makePinVaultServer(t, "")
	defer srv.Close()

	client, _ := vault.NewClient(vault.Config{Address: srv.URL, Token: "test"})
	syncer := internalsync.New(client, "")

	err := syncer.RunPin(internalsync.PinOptions{Path: "secret/myapp/db", Version: 0})
	if err == nil {
		t.Fatal("expected error for invalid version")
	}
}

func TestRunPin_Unpin(t *testing.T) {
	srv := makePinVaultServer(t, "2")
	defer srv.Close()

	client, _ := vault.NewClient(vault.Config{Address: srv.URL, Token: "test"})
	syncer := internalsync.New(client, "")

	var buf bytes.Buffer
	err := syncer.RunPin(internalsync.PinOptions{
		Path:   "secret/myapp/db",
		Unpin:  true,
		Output: &buf,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "unpinned") {
		t.Errorf("expected 'unpinned' in output, got: %s", buf.String())
	}
}

func TestRunGetPin_NotPinned(t *testing.T) {
	srv := makePinVaultServer(t, "")
	defer srv.Close()

	client, _ := vault.NewClient(vault.Config{Address: srv.URL, Token: "test"})
	syncer := internalsync.New(client, "")

	var buf bytes.Buffer
	err := syncer.RunGetPin("secret/myapp/db", &buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "not pinned") {
		t.Errorf("expected 'not pinned' in output, got: %s", buf.String())
	}
}
