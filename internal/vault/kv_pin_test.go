package vault_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/yourusername/vaultpull/internal/vault"
)

func makePinServer(t *testing.T, metaResponse map[string]interface{}) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost, http.MethodPatch:
			w.WriteHeader(http.StatusNoContent)
		case http.MethodGet:
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{"data": metaResponse})
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}))
}

func TestPinSecret_Success(t *testing.T) {
	srv := makePinServer(t, nil)
	defer srv.Close()

	client := makeTestClient(t, srv.URL)
	res, err := client.PinSecret("secret/myapp/db", 3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !res.Pinned {
		t.Error("expected Pinned to be true")
	}
	if res.Version != 3 {
		t.Errorf("expected version 3, got %d", res.Version)
	}
	if got := res.Summary(); got == "" {
		t.Error("expected non-empty summary")
	}
}

func TestUnpinSecret_Success(t *testing.T) {
	srv := makePinServer(t, nil)
	defer srv.Close()

	client := makeTestClient(t, srv.URL)
	res, err := client.UnpinSecret("secret/myapp/db", 3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Pinned {
		t.Error("expected Pinned to be false")
	}
	if summary := res.Summary(); summary == "" {
		t.Error("expected non-empty summary")
	}
}

func TestGetPinnedVersion_Set(t *testing.T) {
	meta := map[string]interface{}{
		"custom_metadata": map[string]interface{}{"pinned_version": "2"},
		"current_version": float64(4),
		"versions":        map[string]interface{}{},
	}
	srv := makePinServer(t, meta)
	defer srv.Close()

	client := makeTestClient(t, srv.URL)
	v, err := client.GetPinnedVersion("secret/myapp/db")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v != 2 {
		t.Errorf("expected pinned version 2, got %d", v)
	}
}

func TestGetPinnedVersion_NotSet(t *testing.T) {
	meta := map[string]interface{}{
		"custom_metadata": map[string]interface{}{},
		"current_version": float64(1),
		"versions":        map[string]interface{}{},
	}
	srv := makePinServer(t, meta)
	defer srv.Close()

	client := makeTestClient(t, srv.URL)
	v, err := client.GetPinnedVersion("secret/myapp/db")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v != 0 {
		t.Errorf("expected 0, got %d", v)
	}
}

func TestPinResult_Summary(t *testing.T) {
	pinned := vault.PinResult{Path: "secret/app", Version: 5, Pinned: true}
	if s := pinned.Summary(); s == "" {
		t.Error("expected non-empty summary for pinned result")
	}
	unpinned := vault.PinResult{Path: "secret/app", Version: 5, Pinned: false}
	if s := unpinned.Summary(); s == "" {
		t.Error("expected non-empty summary for unpinned result")
	}
}
