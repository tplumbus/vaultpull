package vault

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func makeImportServer(t *testing.T, existing map[string]bool) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := strings.TrimPrefix(r.URL.Path, "/v1/")
		switch r.Method {
		case http.MethodGet:
			if existing[path] {
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(map[string]interface{}{
					"data": map[string]interface{}{"data": map[string]string{"key": "val"}},
				})
				return
			}
			w.WriteHeader(http.StatusNotFound)
		case http.MethodPost, http.MethodPut:
			w.WriteHeader(http.StatusNoContent)
		}
	}))
}

func TestImportSecrets_WritesAll(t *testing.T) {
	srv := makeImportServer(t, map[string]bool{})
	defer srv.Close()

	client, _ := NewClient(srv.URL, "test-token", KVv2)
	secrets := map[string]map[string]string{
		"alpha": {"A": "1"},
		"beta":  {"B": "2"},
	}

	result, err := ImportSecrets(context.Background(), client, "secret/data/app", secrets, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Written) != 2 {
		t.Errorf("expected 2 written, got %d", len(result.Written))
	}
	if len(result.Skipped) != 0 {
		t.Errorf("expected 0 skipped, got %d", len(result.Skipped))
	}
}

func TestImportSecrets_SkipsExisting(t *testing.T) {
	srv := makeImportServer(t, map[string]bool{"secret/data/app/alpha": true})
	defer srv.Close()

	client, _ := NewClient(srv.URL, "test-token", KVv2)
	secrets := map[string]map[string]string{
		"alpha": {"A": "1"},
		"beta":  {"B": "2"},
	}

	result, err := ImportSecrets(context.Background(), client, "secret/data/app", secrets, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Skipped) != 1 {
		t.Errorf("expected 1 skipped, got %d", len(result.Skipped))
	}
	if len(result.Written) != 1 {
		t.Errorf("expected 1 written, got %d", len(result.Written))
	}
}

func TestImportResult_Summary(t *testing.T) {
	r := ImportResult{
		Written: []string{"a", "b"},
		Skipped: []string{"c"},
		Failed:  []string{},
	}
	got := r.Summary()
	want := "import complete: 2 written, 1 skipped, 0 failed"
	if got != want {
		t.Errorf("Summary() = %q, want %q", got, want)
	}
}
