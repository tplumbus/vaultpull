package vault

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func makeWatchServer(t *testing.T, responses []map[string]interface{}) (*httptest.Server, int) {
	t.Helper()
	call := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		idx := call
		if idx >= len(responses) {
			idx = len(responses) - 1
		}
		call++
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{"data": responses[idx]})
	}))
	t.Cleanup(server.Close)
	return server, 0
}

func TestWatch_DetectsChange(t *testing.T) {
	responses := []map[string]interface{}{
		{"current_version": float64(1)},
		{"current_version": float64(2)},
	}
	server, _ := makeWatchServer(t, responses)

	c, err := NewClient(server.URL, "test-token")
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	ch := Watch(ctx, c, "secret/data/myapp", WatchOptions{
		Interval:  50 * time.Millisecond,
		KVVersion: KVv2,
	})

	var results []WatchResult
	for r := range ch {
		results = append(results, r)
		if len(results) >= 2 {
			cancel()
		}
	}

	if len(results) < 2 {
		t.Fatalf("expected at least 2 results, got %d", len(results))
	}
	if results[0].Version != 1 {
		t.Errorf("first version = %d, want 1", results[0].Version)
	}
	if results[1].Version != 2 {
		t.Errorf("second version = %d, want 2", results[1].Version)
	}
	if !results[1].Changed {
		t.Error("expected Changed=true on version bump")
	}
}

func TestWatch_NoChangeSkipped(t *testing.T) {
	responses := []map[string]interface{}{
		{"current_version": float64(3)},
		{"current_version": float64(3)},
	}
	server, _ := makeWatchServer(t, responses)

	c, err := NewClient(server.URL, "test-token")
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	ch := Watch(ctx, c, "secret/data/stable", WatchOptions{
		Interval:  40 * time.Millisecond,
		KVVersion: KVv2,
	})

	count := 0
	for range ch {
		count++
	}

	// Only the first poll should emit (Changed on initial detection).
	if count > 1 {
		t.Errorf("expected <=1 result for stable secret, got %d", count)
	}
}
