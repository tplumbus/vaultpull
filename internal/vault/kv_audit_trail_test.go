package vault

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func makeAuditTrailServer(t *testing.T, status int, body interface{}) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		if body != nil {
			json.NewEncoder(w).Encode(body)
		}
	}))
}

func TestGetAuditTrail_Success(t *testing.T) {
	now := time.Now().UTC().Truncate(time.Second)
	body := map[string]interface{}{
		"data": map[string]interface{}{
			"versions": map[string]interface{}{
				"1": map[string]interface{}{
					"created_time":  now.Format(time.RFC3339),
					"deletion_time": "",
					"destroyed":     false,
				},
				"2": map[string]interface{}{
					"created_time":  now.Add(time.Minute).Format(time.RFC3339),
					"deletion_time": now.Add(2 * time.Minute).Format(time.RFC3339),
					"destroyed":     false,
				},
			},
		},
	}
	srv := makeAuditTrailServer(t, http.StatusOK, body)
	defer srv.Close()

	client, err := NewClient(srv.URL, "test-token")
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}

	result, err := GetAuditTrail(client, "secret", "myapp/db")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Path != "myapp/db" {
		t.Errorf("expected path myapp/db, got %s", result.Path)
	}
	if len(result.Entries) != 2 {
		t.Errorf("expected 2 entries, got %d", len(result.Entries))
	}
	summary := result.Summary()
	if summary == "" {
		t.Error("expected non-empty summary")
	}
}

func TestGetAuditTrail_NotFound(t *testing.T) {
	srv := makeAuditTrailServer(t, http.StatusNotFound, nil)
	defer srv.Close()

	client, err := NewClient(srv.URL, "test-token")
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}

	_, err = GetAuditTrail(client, "secret", "missing/path")
	if err == nil {
		t.Fatal("expected error for not found")
	}
}

func TestGetAuditTrail_ServerError(t *testing.T) {
	srv := makeAuditTrailServer(t, http.StatusInternalServerError, nil)
	defer srv.Close()

	client, err := NewClient(srv.URL, "test-token")
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}

	_, err = GetAuditTrail(client, "secret", "myapp/db")
	if err == nil {
		t.Fatal("expected error for server error")
	}
}

func TestAuditTrailResult_Summary(t *testing.T) {
	r := AuditTrailResult{Path: "myapp/db", Entries: []AuditEntry{{}, {}}}
	s := r.Summary()
	if s != `audit trail for "myapp/db": 2 entries` {
		t.Errorf("unexpected summary: %s", s)
	}
}
