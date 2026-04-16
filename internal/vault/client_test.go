package vault

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func makeServer(t *testing.T, status int, data map[string]string) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Vault-Token") == "" {
			w.WriteHeader(http.StatusForbidden)
			return
		}
		w.WriteHeader(status)
		if data != nil {
			body, _ := json.Marshal(map[string]interface{}{
				"data": map[string]interface{}{"data": data},
			})
			w.Write(body)
		}
	}))
}

func TestGetSecret_Success(t *testing.T) {
	expected := map[string]string{"API_KEY": "abc123", "DB_PASS": "secret"}
	srv := makeServer(t, http.StatusOK, expected)
	defer srv.Close()

	client := NewClient(srv.URL, "test-token")
	data, err := client.GetSecret("secret/data/myapp")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for k, v := range expected {
		if data[k] != v {
			t.Errorf("key %s: got %q, want %q", k, data[k], v)
		}
	}
}

func TestGetSecret_NotFound(t *testing.T) {
	srv := makeServer(t, http.StatusNotFound, nil)
	defer srv.Close()

	client := NewClient(srv.URL, "test-token")
	_, err := client.GetSecret("secret/data/missing")
	if err == nil {
		t.Fatal("expected error for 404, got nil")
	}
}

func TestGetSecret_ServerError(t *testing.T) {
	srv := makeServer(t, http.StatusInternalServerError, nil)
	defer srv.Close()

	client := NewClient(srv.URL, "test-token")
	_, err := client.GetSecret("secret/data/broken")
	if err == nil {
		t.Fatal("expected error for 500, got nil")
	}
}
