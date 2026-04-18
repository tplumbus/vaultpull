package vault

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func makeRestoreServer(t *testing.T, version int, secret map[string]string) *httptest.Server {
	t.Helper()
	written := map[string]string{}
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			if len(secret) == 0 {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]any{
				"data": map[string]any{"data": secret},
			})
			return
		}
		if r.Method == http.MethodPut || r.Method == http.MethodPost {
			var body map[string]any
			json.NewDecoder(r.Body).Decode(&body)
			if d, ok := body["data"].(map[string]any); ok {
				for k, v := range d {
					written[k] = v.(string)
				}
			}
			w.WriteHeader(http.StatusOK)
			return
		}
		w.WriteHeader(http.StatusMethodNotAllowed)
	}))
}

func TestRestoreSecret_Success(t *testing.T) {
	secret := map[string]string{"KEY": "value123"}
	srv := makeRestoreServer(t, 2, secret)
	defer srv.Close()

	c, err := NewClient(srv.URL, "test-token", KVv2)
	if err != nil {
		t.Fatal(err)
	}
	result, err := c.RestoreSecret(context.Background(), "secret/myapp", 2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result["KEY"] != "value123" {
		t.Errorf("expected KEY=value123, got %q", result["KEY"])
	}
}

func TestRestoreSecret_EmptyVersion(t *testing.T) {
	srv := makeRestoreServer(t, 1, map[string]string{})
	defer srv.Close()

	c, err := NewClient(srv.URL, "test-token", KVv2)
	if err != nil {
		t.Fatal(err)
	}
	_, err = c.RestoreSecret(context.Background(), "secret/myapp", 1)
	if err == nil {
		t.Fatal("expected error for empty version")
	}
}
