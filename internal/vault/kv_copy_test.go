package vault

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func makeCopyServer(t *testing.T) *httptest.Server {
	t.Helper()
	secrets := map[string]map[string]string{
		"/v1/secret/data/src": {"FOO": "bar", "BAZ": "qux"},
		"/v1/secret/data/dst": {"EXISTING": "keep"},
	}
	written := map[string]map[string]string{}

	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			data, ok := secrets[r.URL.Path]
			if !ok {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			kv := map[string]interface{}{"data": map[string]interface{}{"data": data}}
			json.NewEncoder(w).Encode(kv)
			return
		}
		if r.Method == http.MethodPost || r.Method == http.MethodPut {
			var body map[string]interface{}
			json.NewDecoder(r.Body).Decode(&body)
			if d, ok := body["data"].(map[string]interface{}); ok {
				m := map[string]string{}
				for k, v := range d {
					m[k] = v.(string)
				}
				written[r.URL.Path] = m
			}
			w.WriteHeader(http.StatusOK)
			return
		}
		w.WriteHeader(http.StatusMethodNotAllowed)
	}))
}

func TestCopySecret_NoOverwrite(t *testing.T) {
	srv := makeCopyServer(t)
	defer srv.Close()

	client, err := NewClient(srv.URL, "test-token", KVv2)
	if err != nil {
		t.Fatal(err)
	}

	n, err := client.CopySecret(context.Background(), "secret/src", "secret/dst", false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n != 2 {
		t.Errorf("expected 2 keys copied, got %d", n)
	}
}

func TestCopySecret_WithOverwrite(t *testing.T) {
	srv := makeCopyServer(t)
	defer srv.Close()

	client, err := NewClient(srv.URL, "test-token", KVv2)
	if err != nil {
		t.Fatal(err)
	}

	n, err := client.CopySecret(context.Background(), "secret/src", "secret/dst", true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n != 2 {
		t.Errorf("expected 2 keys copied, got %d", n)
	}
}
