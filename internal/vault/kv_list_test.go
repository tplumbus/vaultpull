package vault

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"sort"
	"testing"
)

func makeRecursiveServer(t *testing.T) *httptest.Server {
	t.Helper()
	mux := http.NewServeMux()

	// root list
	mux.HandleFunc("/v1/secret/metadata/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "LIST" {
			http.NotFound(w, r)
			return
		}
		path := r.URL.Path
		switch path {
		case "/v1/secret/metadata/":
			json.NewEncoder(w).Encode(map[string]interface{}{"data": map[string]interface{}{"keys": []string{"app/", "db"}}})
		case "/v1/secret/metadata/app":
			json.NewEncoder(w).Encode(map[string]interface{}{"data": map[string]interface{}{"keys": []string{"config", "creds"}}})
		default:
			http.NotFound(w, r)
		}
	})
	return httptest.NewServer(mux)
}

func TestRecursiveList_Success(t *testing.T) {
	srv := makeRecursiveServer(t)
	defer srv.Close()

	c, err := NewClient(srv.URL, "test-token", KVv2)
	if err != nil {
		t.Fatal(err)
	}

	got, err := c.RecursiveList(context.Background(), "secret", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := []string{"app/config", "app/creds", "db"}
	sort.Strings(got)
	sort.Strings(want)
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}
