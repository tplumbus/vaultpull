package vault

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func makeListServer(t *testing.T, status int, keys []string) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "LIST" {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		w.WriteHeader(status)
		if keys != nil {
			body, _ := json.Marshal(map[string]interface{}{
				"data": map[string]interface{}{"keys": keys},
			})
			w.Write(body)
		}
	}))
}

func TestListSecrets_Success(t *testing.T) {
	expected := []string{_PASSWORD", "API_KEY"}
	srv := makeListServer(t, http.StatusOK, expected)
	defer srv.Close()

	c, err := NewClient(srv.URL, "test-token")
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}

	keys, err := c.ListSecrets("secret/data/myapp")
	if err != nil {
		t.Fatalf("ListSecrets: %v", err)
	}
	if len(keys) != len(expected) {
		t.Fatalf("expected %d keys, got %d", len(expected), len(keys))
	}
	for i, k := range keys {
		if k != expected[i] {
			t.Errorf("key[%d]: want %q, got %q", i, expected[i], k)
		}
	}
}

func TestListSecrets_NotFound(t *testing.T) {
	srv := makeListServer(t, http.StatusNotFound, nil)
	defer srv.Close()

	c, err := NewClient(srv.URL, "test-token")
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}

	_, err = c.ListSecrets("secret/data/missing")
	if err == nil {
		t.Fatal("expected error for 404, got nil")
	}
}

func TestListSecrets_ServerError(t *testing.T) {
	srv := makeListServer(t, http.StatusInternalServerError, nil)
	defer srv.Close()

	c, err := NewClient(srv.URL, "test-token")
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}

	_, err = c.ListSecrets("secret/data/myapp")
	if err == nil {
		t.Fatal("expected error for 500, got nil")
	}
}

func TestToMetadataPath(t *testing.T) {
	cases := []struct{ in, out string }{
		{"secret/data/myapp", "secret/metadata/myapp"},
		{"kv/data/team/app", "kv/metadata/team/app"},
		{"secret/myapp", "secret/myapp"}, // no data segment, unchanged
	}
	for _, tc := range cases {
		got := toMetadataPath(tc.in)
		if got != tc.out {
			t.Errorf("toMetadataPath(%q) = %q, want %q", tc.in, got, tc.out)
		}
	}
}
