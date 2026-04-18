package vault

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func makeRollbackServer(t *testing.T) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodGet && strings.Contains(r.URL.RawQuery, "version=2"):
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{
				"data": map[string]interface{}{
					"data": map[string]string{"KEY": "old_value"},
				},
			})
		case r.Method == http.MethodGet && strings.Contains(r.URL.RawQuery, "version=99"):
			w.WriteHeader(http.StatusNotFound)
		case r.Method == http.MethodPost:
			w.WriteHeader(http.StatusOK)
		default:
			w.WriteHeader(http.StatusBadRequest)
		}
	}))
}

func TestRollbackSecret_Success(t *testing.T) {
	srv := makeRollbackServer(t)
	defer srv.Close()

	c, err := NewClient(srv.URL, "test-token", KVv2)
	if err != nil {
		t.Fatal(err)
	}

	err = c.RollbackSecret(context.Background(), "secret/data/myapp", 2)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestRollbackSecret_VersionNotFound(t *testing.T) {
	srv := makeRollbackServer(t)
	defer srv.Close()

	c, err := NewClient(srv.URL, "test-token", KVv2)
	if err != nil {
		t.Fatal(err)
	}

	err = c.RollbackSecret(context.Background(), "secret/data/myapp", 99)
	if err == nil {
		t.Fatal("expected error for missing version")
	}
	if !strings.Contains(err.Error(), "version 99") {
		t.Errorf("expected version info in error, got: %v", err)
	}
}
