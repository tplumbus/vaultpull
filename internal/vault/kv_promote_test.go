package vault

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func makePromoteServer(t *testing.T) *httptest.Server {
	t.Helper()
	store := map[string]map[string]interface{}{
		"/v1/secret/data/dev/app": {"API_KEY": "abc123"},
	}
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			d, ok := store[r.URL.Path]
			if !ok {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			json.NewEncoder(w).Encode(map[string]interface{}{"data": map[string]interface{}{"data": d}})
		case http.MethodPost, http.MethodPut:
			var body map[string]interface{}
			json.NewDecoder(r.Body).Decode(&body)
			if data, ok := body["data"]; ok {
				store[r.URL.Path] = data.(map[string]interface{})
			}
			w.WriteHeader(http.StatusOK)
		}
	}))
}

func TestPromoteSecret_Success(t *testing.T) {
	srv := makePromoteServer(t)
	defer srv.Close()
	c := newTestClient(t, srv.URL)

	err := c.PromoteSecret("secret/dev/app", "secret/prod/app", false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestPromoteSecret_OverwriteBlocked(t *testing.T) {
	srv := makePromoteServer(t)
	defer srv.Close()
	c := newTestClient(t, srv.URL)

	// First promote succeeds
	_ = c.PromoteSecret("secret/dev/app", "secret/prod/app", false)
	// Second should be blocked
	err := c.PromoteSecret("secret/dev/app", "secret/prod/app", false)
	if err == nil {
		t.Fatal("expected error when overwrite=false and dst exists")
	}
}

func TestReplaceEnvInPath(t *testing.T) {
	got := ReplaceEnvInPath("secret/dev/app/db", "dev", "prod")
	want := "secret/prod/app/db"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}
