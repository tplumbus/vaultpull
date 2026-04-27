package vault

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func makeSnapshotServer(t *testing.T, status int, payload interface{}) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(status)
		if payload != nil {
			_ = json.NewEncoder(w).Encode(payload)
		}
	}))
}

func TestTakeSnapshot_Success(t *testing.T) {
	payload := map[string]interface{}{
		"data": map[string]interface{}{
			"data": map[string]string{"API_KEY": "abc123", "DB_PASS": "secret"},
			"metadata": map[string]interface{}{"version": 3},
		},
	}
	srv := makeSnapshotServer(t, http.StatusOK, payload)
	defer srv.Close()

	result, err := TakeSnapshot(srv.Client(), srv.URL, "test-token", "secret/data/myapp")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Skipped {
		t.Fatal("expected snapshot to succeed, got skipped")
	}
	if result.Snapshot == nil {
		t.Fatal("expected non-nil snapshot")
	}
	if result.Snapshot.Version != 3 {
		t.Errorf("expected version 3, got %d", result.Snapshot.Version)
	}
	if len(result.Snapshot.Data) != 2 {
		t.Errorf("expected 2 keys, got %d", len(result.Snapshot.Data))
	}
	if result.Snapshot.Data["API_KEY"] != "abc123" {
		t.Errorf("expected API_KEY=abc123, got %s", result.Snapshot.Data["API_KEY"])
	}
}

func TestTakeSnapshot_NotFound(t *testing.T) {
	srv := makeSnapshotServer(t, http.StatusNotFound, nil)
	defer srv.Close()

	result, err := TakeSnapshot(srv.Client(), srv.URL, "test-token", "secret/data/missing")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.Skipped {
		t.Fatal("expected result to be skipped")
	}
	if !strings.Contains(result.Reason, "not found") {
		t.Errorf("expected 'not found' reason, got: %s", result.Reason)
	}
}

func TestTakeSnapshot_ServerError(t *testing.T) {
	srv := makeSnapshotServer(t, http.StatusInternalServerError, nil)
	defer srv.Close()

	_, err := TakeSnapshot(srv.Client(), srv.URL, "test-token", "secret/data/myapp")
	if err == nil {
		t.Fatal("expected error on server error, got nil")
	}
}

func TestSnapshotResult_Summary(t *testing.T) {
	result := SnapshotResult{
		Snapshot: &Snapshot{
			Path:    "secret/data/myapp",
			Data:    map[string]string{"A": "1", "B": "2"},
			Version: 5,
		},
	}
	summary := result.Summary()
	if !strings.Contains(summary, "2 keys") {
		t.Errorf("expected '2 keys' in summary, got: %s", summary)
	}
	if !strings.Contains(summary, "version 5") {
		t.Errorf("expected 'version 5' in summary, got: %s", summary)
	}

	skipped := SnapshotResult{Skipped: true, Reason: "path not found"}
	if !strings.Contains(skipped.Summary(), "skipped") {
		t.Errorf("expected 'skipped' in summary, got: %s", skipped.Summary())
	}
}
