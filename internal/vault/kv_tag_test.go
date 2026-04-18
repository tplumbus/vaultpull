package vault

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func makeTagServer(t *testing.T, tags map[string]string, status int) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(status)
		if status == http.StatusOK {
			body := map[string]interface{}{
				"data": map[string]interface{}{
					"custom_metadata": tags,
				},
			}
			_ = json.NewEncoder(w).Encode(body)
		}
	}))
}

func TestGetTags_Success(t *testing.T) {
	tags := map[string]string{"env": "production", "team": "platform"}
	srv := makeTagServer(t, tags, http.StatusOK)
	defer srv.Close()

	c, _ := NewClient(srv.URL, "token", KVVersion2)
	got, err := c.GetTags("secret/myapp")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["env"] != "production" {
		t.Errorf("expected production, got %s", got["env"])
	}
	if got["team"] != "platform" {
		t.Errorf("expected platform, got %s", got["team"])
	}
}

func TestGetTags_NotFound(t *testing.T) {
	srv := makeTagServer(t, nil, http.StatusNotFound)
	defer srv.Close()

	c, _ := NewClient(srv.URL, "token", KVVersion2)
	_, err := c.GetTags("secret/missing")
	if err == nil {
		t.Fatal("expected error for not found")
	}
}

func TestGetTags_ServerError(t *testing.T) {
	srv := makeTagServer(t, nil, http.StatusInternalServerError)
	defer srv.Close()

	c, _ := NewClient(srv.URL, "token", KVVersion2)
	_, err := c.GetTags("secret/myapp")
	if err == nil {
		t.Fatal("expected error for server error")
	}
}

func TestGetTags_Empty(t *testing.T) {
	srv := makeTagServer(t, map[string]string{}, http.StatusOK)
	defer srv.Close()

	c, _ := NewClient(srv.URL, "token", KVVersion2)
	got, err := c.GetTags("secret/myapp")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 0 {
		t.Errorf("expected empty tags, got %v", got)
	}
}
