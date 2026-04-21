package vault

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func makeExportServer(t *testing.T, status int, payload interface{}) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(status)
		if payload != nil {
			_ = json.NewEncoder(w).Encode(payload)
		}
	}))
}

func TestExportSecret_Success(t *testing.T) {
	payload := map[string]interface{}{
		"data": map[string]interface{}{
			"data": map[string]string{"API_KEY": "abc123", "DB_PASS": "secret"},
			"metadata": map[string]interface{}{"version": 3},
		},
	}
	srv := makeExportServer(t, http.StatusOK, payload)
	defer srv.Close()

	client := &Client{Address: srv.URL, Token: "test-token", HTTP: srv.Client()}
	result, err := ExportSecret(client, "secret", "myapp/prod")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Path != "myapp/prod" {
		t.Errorf("expected path myapp/prod, got %s", result.Path)
	}
	if result.Version != 3 {
		t.Errorf("expected version 3, got %d", result.Version)
	}
	if result.Secrets["API_KEY"] != "abc123" {
		t.Errorf("expected API_KEY=abc123, got %s", result.Secrets["API_KEY"])
	}
	if len(result.Secrets) != 2 {
		t.Errorf("expected 2 secrets, got %d", len(result.Secrets))
	}
}

func TestExportSecret_NotFound(t *testing.T) {
	srv := makeExportServer(t, http.StatusNotFound, nil)
	defer srv.Close()

	client := &Client{Address: srv.URL, Token: "test-token", HTTP: srv.Client()}
	_, err := ExportSecret(client, "secret", "missing/path")
	if err == nil {
		t.Fatal("expected error for not found, got nil")
	}
}

func TestExportSecret_ServerError(t *testing.T) {
	srv := makeExportServer(t, http.StatusInternalServerError, nil)
	defer srv.Close()

	client := &Client{Address: srv.URL, Token: "test-token", HTTP: srv.Client()}
	_, err := ExportSecret(client, "secret", "myapp/prod")
	if err == nil {
		t.Fatal("expected error for server error, got nil")
	}
}

func TestExportResult_Summary(t *testing.T) {
	r := ExportResult{Path: "myapp/prod", Secrets: map[string]string{"A": "1", "B": "2"}, Version: 5}
	got := r.Summary()
	expected := "exported 2 keys from myapp/prod at version 5"
	if got != expected {
		t.Errorf("expected %q, got %q", expected, got)
	}
}
