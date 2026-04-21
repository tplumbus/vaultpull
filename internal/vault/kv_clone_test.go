package vault

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func makeCloneServer(t *testing.T) *httptest.Server {
	t.Helper()
	store := map[string]map[string]interface{}{
		"/v1/secret/data/src": {"foo": "bar", "baz": "qux"},
	}

	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			data, ok := store[r.URL.Path]
			if !ok {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			json.NewEncoder(w).Encode(map[string]interface{}{
				"data": map[string]interface{}{"data": data},
			})
		case http.MethodPost, http.MethodPut:
			var body map[string]interface{}
			json.NewDecoder(r.Body).Decode(&body)
			if d, ok := body["data"]; ok {
				store[r.URL.Path] = d.(map[string]interface{})
			}
			w.WriteHeader(http.StatusOK)
		}
	}))
}

func TestCloneSecret_Success(t *testing.T) {
	srv := makeCloneServer(t)
	defer srv.Close()

	c := newTestClient(t, srv.URL)
	err := c.CloneSecret("secret/data/src", []string{"secret/data/dst1", "secret/data/dst2"}, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCloneSecret_OverwriteBlocked(t *testing.T) {
	srv := makeCloneServer(t)
	defer srv.Close()

	c := newTestClient(t, srv.URL)
	// dst already exists in the store (same as src path for test simplicity)
	err := c.CloneSecret("secret/data/src", []string{"secret/data/src"}, false)
	if err == nil {
		t.Fatal("expected overwrite error, got nil")
	}
	if !strings.Contains(err.Error(), "already exists") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestCloneSecret_OverwriteAllowed(t *testing.T) {
	srv := makeCloneServer(t)
	defer srv.Close()

	c := newTestClient(t, srv.URL)
	err := c.CloneSecret("secret/data/src", []string{"secret/data/src"}, true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCloneResult_Summary(t *testing.T) {
	r := CloneResult{
		Source:       "secret/src",
		Destinations: []string{"secret/dst1", "secret/dst2"},
		Skipped:      []string{"secret/dst3"},
	}
	got := r.Summary()
	if !strings.Contains(got, "2 destination") {
		t.Errorf("expected destination count in summary, got: %s", got)
	}
	if !strings.Contains(got, "1 skipped") {
		t.Errorf("expected skipped count in summary, got: %s", got)
	}
}
