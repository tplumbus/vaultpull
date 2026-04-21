package vault

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func makeSearchServer(t *testing.T) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodGet && strings.Contains(r.URL.Path, "metadata"):
			// list
			keys := []string{"alpha", "beta"}
			body, _ := json.Marshal(map[string]interface{}{
				"data": map[string]interface{}{"keys": keys},
			})
			w.WriteHeader(http.StatusOK)
			w.Write(body)
		case r.Method == http.MethodGet && strings.HasSuffix(r.URL.Path, "alpha"):
			body, _ := json.Marshal(map[string]interface{}{
				"data": map[string]interface{}{"data": map[string]string{"DB_HOST": "localhost", "APP_ENV": "prod"}},
			})
			w.WriteHeader(http.StatusOK)
			w.Write(body)
		case r.Method == http.MethodGet && strings.HasSuffix(r.URL.Path, "beta"):
			body, _ := json.Marshal(map[string]interface{}{
				"data": map[string]interface{}{"data": map[string]string{"API_KEY": "secret123"}},
			})
			w.WriteHeader(http.StatusOK)
			w.Write(body)
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
}

func TestSearchSecrets_MatchesKey(t *testing.T) {
	srv := makeSearchServer(t)
	defer srv.Close()

	c, _ := NewClient(srv.URL, "test-token", KVv2)
	results, err := c.SearchSecrets("secret/myapp", "DB")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Keys[0] != "DB_HOST" {
		t.Errorf("expected DB_HOST, got %s", results[0].Keys[0])
	}
}

func TestSearchSecrets_NoMatch(t *testing.T) {
	srv := makeSearchServer(t)
	defer srv.Close()

	c, _ := NewClient(srv.URL, "test-token", KVv2)
	results, err := c.SearchSecrets("secret/myapp", "NOTEXIST")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 0 {
		t.Errorf("expected 0 results, got %d", len(results))
	}
}

func TestSearchSummary(t *testing.T) {
	results := []SearchResult{
		{Path: "secret/a", Keys: []string{"FOO", "BAR"}},
		{Path: "secret/b", Keys: []string{"BAZ"}},
	}
	got := SearchSummary(results)
	if got != "3 match(es) across 2 path(s)" {
		t.Errorf("unexpected summary: %s", got)
	}

	empty := SearchSummary(nil)
	if empty != "no matches found" {
		t.Errorf("unexpected empty summary: %s", empty)
	}
}
