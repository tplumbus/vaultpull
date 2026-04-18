package vault

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func makeMetadataServer(t *testing.T, status int, payload interface{}) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(status)
		if payload != nil {
			json.NewEncoder(w).Encode(payload)
		}
	}))
}

func TestGetSecretMetadata_Success(t *testing.T) {
	now := time.Now().UTC().Truncate(time.Second)
	payload := map[string]interface{}{
		"data": map[string]interface{}{
			"created_time":    now.Format(time.RFC3339),
			"updated_time":    now.Format(time.RFC3339),
			"current_version": 3,
			"oldest_version":  1,
		},
	}
	srv := makeMetadataServer(t, http.StatusOK, payload)
	defer srv.Close()

	client, _ := NewClient(srv.URL, "test-token", KVv2)
	meta, err := client.GetSecretMetadata("secret/data/myapp")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if meta.CurrentVersion != 3 {
		t.Errorf("expected version 3, got %d", meta.CurrentVersion)
	}
	if meta.OldestVersion != 1 {
		t.Errorf("expected oldest 1, got %d", meta.OldestVersion)
	}
}

func TestGetSecretMetadata_NotFound(t *testing.T) {
	srv := makeMetadataServer(t, http.StatusNotFound, nil)
	defer srv.Close()

	client, _ := NewClient(srv.URL, "test-token", KVv2)
	_, err := client.GetSecretMetadata("secret/data/missing")
	if err == nil {
		t.Fatal("expected error for not found")
	}
}

func TestGetSecretMetadata_ServerError(t *testing.T) {
	srv := makeMetadataServer(t, http.StatusInternalServerError, nil)
	defer srv.Close()

	client, _ := NewClient(srv.URL, "test-token", KVv2)
	_, err := client.GetSecretMetadata("secret/data/myapp")
	if err == nil {
		t.Fatal("expected error for server error")
	}
}
