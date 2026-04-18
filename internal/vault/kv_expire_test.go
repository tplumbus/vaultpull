package vault

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func makeExpireServer(expiresAt string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost || r.Method == http.MethodPut {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		if strings.Contains(r.URL.Path, "metadata") {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{
				"data": map[string]interface{}{
					"custom_metadata": map[string]interface{}{
						"expires_at": expiresAt,
					},
				},
			})
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}))
}

func TestSetExpiry_Success(t *testing.T) {
	srv := makeExpireServer("")
	defer srv.Close()
	c := newTestClient(t, srv.URL)
	if err := c.SetExpiry("secret/data/myapp/db", 24*time.Hour); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestGetExpiry_Success(t *testing.T) {
	future := time.Now().UTC().Add(2 * time.Hour).Format(time.RFC3339)
	srv := makeExpireServer(future)
	defer srv.Close()
	c := newTestClient(t, srv.URL)
	exp, err := c.GetExpiry("secret/data/myapp/db")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if exp.Expired {
		t.Error("expected secret to not be expired")
	}
	if exp.TTL <= 0 {
		t.Error("expected positive TTL")
	}
}

func TestGetExpiry_Expired(t *testing.T) {
	past := time.Now().UTC().Add(-1 * time.Hour).Format(time.RFC3339)
	srv := makeExpireServer(past)
	defer srv.Close()
	c := newTestClient(t, srv.URL)
	exp, err := c.GetExpiry("secret/data/myapp/db")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !exp.Expired {
		t.Error("expected secret to be expired")
	}
}
