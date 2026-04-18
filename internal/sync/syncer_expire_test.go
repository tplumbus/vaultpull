package sync

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func makeExpireVaultServer(expiresAt string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost || r.Method == http.MethodPut {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		if strings.Contains(r.URL.Path, "metadata") {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{
				"data": map[string]interface{}{
					"custom_metadata": map[string]interface{}{
						"expires_at": expiresAt,
					},
				},
			})
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}))
}

func TestRunCheckExpiry_NotExpired(t *testing.T) {
	future := time.Now().UTC().Add(3 * time.Hour).Format(time.RFC3339)
	srv := makeExpireVaultServer(future)
	defer srv.Close()
	s := newTestSyncer(t, srv.URL)
	var buf bytes.Buffer
	results, err := s.RunCheckExpiry([]string{"secret/data/myapp/token"}, &buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Expired {
		t.Error("expected not expired")
	}
	if !strings.Contains(buf.String(), "OK") {
		t.Errorf("expected OK in output, got: %s", buf.String())
	}
}

func TestRunCheckExpiry_Expired(t *testing.T) {
	past := time.Now().UTC().Add(-2 * time.Hour).Format(time.RFC3339)
	srv := makeExpireVaultServer(past)
	defer srv.Close()
	s := newTestSyncer(t, srv.URL)
	var buf bytes.Buffer
	results, _ := s.RunCheckExpiry([]string{"secret/data/myapp/token"}, &buf)
	if len(results) == 0 || !results[0].Expired {
		t.Error("expected expired result")
	}
	if !strings.Contains(buf.String(), "EXPIRED") {
		t.Errorf("expected EXPIRED in output, got: %s", buf.String())
	}
}

func TestRunSetExpiry_Success(t *testing.T) {
	srv := makeExpireVaultServer("")
	defer srv.Close()
	s := newTestSyncer(t, srv.URL)
	var buf bytes.Buffer
	if err := s.RunSetExpiry([]string{"secret/data/myapp/token"}, 48*time.Hour, &buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "SET_EXPIRY") {
		t.Errorf("expected SET_EXPIRY in output, got: %s", buf.String())
	}
}
