package vault

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func makeRenameServer(t *testing.T) *httptest.Server {
	t.Helper()
	store := map[string]map[string]string{
		"secret/data/src": {"API_KEY": "abc123"},
	}
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := strings.TrimPrefix(r.URL.Path, "/v1/")
		switch r.Method {
		case http.MethodGet:
			data, ok := store[path]
			if !ok {
				w.WriteHeader(404)
				return
			}
			j, _ := json.Marshal(map[string]any{"data": map[string]any{"data": data}})
			w.Write(j)
		case http.MethodPut, http.MethodPost:
			var body map[string]any
			json.NewDecoder(r.Body).Decode(&body)
			if d, ok := body["data"]; ok {
				m := map[string]string{}
				for k, v := range d.(map[string]any) {
					m[k] = v.(string)
				}
				store[path] = m
			}
			w.WriteHeader(200)
		case http.MethodDelete:
			if _, ok := store[path]; !ok {
				w.WriteHeader(404)
				return
			}
			delete(store, path)
			w.WriteHeader(204)
		}
	}))
}

func TestRenameSecret_Success(t *testing.T) {
	srv := makeRenameServer(t)
	defer srv.Close()
	c, _ := NewClient(srv.URL, "token")
	if err := c.RenameSecret("secret", "src", "dst", false); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestDeleteSecret_NotFound(t *testing.T) {
	srv := makeRenameServer(t)
	defer srv.Close()
	c, _ := NewClient(srv.URL, "token")
	if err := c.DeleteSecret("secret", "nonexistent"); err == nil {
		t.Fatal("expected error for missing path")
	}
}

func TestRenameSecret_OverwriteBlocked(t *testing.T) {
	srv := makeRenameServer(t)
	defer srv.Close()
	c, _ := NewClient(srv.URL, "token")
	// pre-create dst
	c.CopySecret("secret", "src", "dst", true)
	if err := c.RenameSecret("secret", "src", "dst", false); err == nil {
		t.Fatal("expected overwrite error")
	}
}
