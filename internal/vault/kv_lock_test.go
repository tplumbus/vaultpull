package vault

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func makeLockServer(t *testing.T, locked bool) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			val := "false"
			if locked {
				val = "true"
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{
				"data": map[string]interface{}{
					"custom_metadata": map[string]string{"locked": val},
					"versions":        map[string]interface{}{},
				},
			})
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}))
}

func TestIsLocked_True(t *testing.T) {
	srv := makeLockServer(t, true)
	defer srv.Close()
	c, _ := NewClient(srv.URL, "token", KVv2)
	locked, err := c.IsLocked("secret/data/myapp")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !locked {
		t.Error("expected secret to be locked")
	}
}

func TestIsLocked_False(t *testing.T) {
	srv := makeLockServer(t, false)
	defer srv.Close()
	c, _ := NewClient(srv.URL, "token", KVv2)
	locked, err := c.IsLocked("secret/data/myapp")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if locked {
		t.Error("expected secret to be unlocked")
	}
}

func TestLockSecret_Success(t *testing.T) {
	srv := makeLockServer(t, false)
	defer srv.Close()
	c, _ := NewClient(srv.URL, "token", KVv2)
	if err := c.LockSecret("secret/data/myapp"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestUnlockSecret_Success(t *testing.T) {
	srv := makeLockServer(t, true)
	defer srv.Close()
	c, _ := NewClient(srv.URL, "token", KVv2)
	if err := c.UnlockSecret("secret/data/myapp"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
