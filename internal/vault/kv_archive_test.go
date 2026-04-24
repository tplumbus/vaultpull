package vault

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func makeArchiveServer(t *testing.T) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodGet && strings.Contains(r.URL.Path, "/data/") && !strings.Contains(r.URL.Path, "archive"):
			json.NewEncoder(w).Encode(map[string]any{
				"data": map[string]any{
					"data": map[string]string{"KEY": "val"},
					"metadata": map[string]any{"version": 3},
				},
			})
		case r.Method == http.MethodGet && strings.Contains(r.URL.Path, "/metadata/"):
			json.NewEncoder(w).Encode(map[string]any{
				"data": map[string]any{
					"current_version": 3,
					"oldest_version":  1,
					"versions":        map[string]any{},
				},
			})
		case r.Method == http.MethodPost && strings.Contains(r.URL.Path, "archive"):
			w.WriteHeader(http.StatusOK)
		case r.Method == http.MethodGet && strings.Contains(r.URL.Path, "/metadata/archive"):
			json.NewEncoder(w).Encode(map[string]any{
				"data": map[string]any{
					"keys": []string{"v1", "v2", "v3"},
				},
			})
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
}

func TestArchiveSecret_Success(t *testing.T) {
	srv := makeArchiveServer(t)
	defer srv.Close()

	c, err := NewClient(srv.URL, "test-token", KVv2)
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}

	res, err := ArchiveSecret(c, "myapp/db")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !res.OK {
		t.Error("expected OK=true")
	}
	if res.Version != 3 {
		t.Errorf("expected version 3, got %d", res.Version)
	}
}

func TestArchiveResult_Summary(t *testing.T) {
	ok := ArchiveResult{Path: "myapp/db", Version: 3, OK: true}
	if s := ok.Summary(); s != "archived myapp/db at version 3" {
		t.Errorf("unexpected summary: %s", s)
	}

	fail := ArchiveResult{Path: "myapp/db", OK: false}
	if s := fail.Summary(); s != "failed to archive myapp/db" {
		t.Errorf("unexpected summary: %s", s)
	}
}

func TestListArchives_Success(t *testing.T) {
	srv := makeArchiveServer(t)
	defer srv.Close()

	c, err := NewClient(srv.URL, "test-token", KVv2)
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}

	keys, err := ListArchives(c, "myapp/db")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(keys) != 3 {
		t.Errorf("expected 3 keys, got %d", len(keys))
	}
}
