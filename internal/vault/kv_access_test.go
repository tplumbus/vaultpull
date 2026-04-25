package vault

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func makeAccessServer(t *testing.T, status int, body interface{}) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		if body != nil {
			_ = json.NewEncoder(w).Encode(body)
		}
	}))
}

func TestGetAccessPolicies_Success(t *testing.T) {
	body := map[string]interface{}{
		"data": map[string]interface{}{
			"exact_rules": map[string]interface{}{
				"secret/data/myapp": map[string]interface{}{
					"capabilities": []string{"read", "list"},
				},
			},
		},
	}
	srv := makeAccessServer(t, http.StatusOK, body)
	defer srv.Close()

	client := &Client{addr: srv.URL, token: "tok", http: srv.Client()}
	result, err := client.GetAccessPolicies("secret/data/myapp")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.SecretPath != "secret/data/myapp" {
		t.Errorf("expected path secret/data/myapp, got %s", result.SecretPath)
	}
	if len(result.Policies) != 1 {
		t.Fatalf("expected 1 policy, got %d", len(result.Policies))
	}
	if result.Policies[0].Path != "secret/data/myapp" {
		t.Errorf("unexpected policy path: %s", result.Policies[0].Path)
	}
}

func TestGetAccessPolicies_NotFound(t *testing.T) {
	srv := makeAccessServer(t, http.StatusNotFound, nil)
	defer srv.Close()

	client := &Client{addr: srv.URL, token: "tok", http: srv.Client()}
	result, err := client.GetAccessPolicies("secret/data/missing")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Policies) != 0 {
		t.Errorf("expected empty policies, got %d", len(result.Policies))
	}
}

func TestGetAccessPolicies_ServerError(t *testing.T) {
	srv := makeAccessServer(t, http.StatusInternalServerError, nil)
	defer srv.Close()

	client := &Client{addr: srv.URL, token: "tok", http: srv.Client()}
	_, err := client.GetAccessPolicies("secret/data/myapp")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestAccessResult_Summary(t *testing.T) {
	result := AccessResult{
		SecretPath: "secret/data/myapp",
		Policies: []AccessPolicy{
			{Path: "secret/data/myapp", Capabilities: []string{"read", "list"}},
		},
	}
	summary := result.Summary()
	if summary == "" {
		t.Error("expected non-empty summary")
	}

	empty := AccessResult{SecretPath: "secret/data/none"}
	if empty.Summary() == "" {
		t.Error("expected fallback summary for empty policies")
	}
}
