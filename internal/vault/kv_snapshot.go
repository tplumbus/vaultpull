package vault

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// Snapshot represents a point-in-time capture of all key-value pairs at a path.
type Snapshot struct {
	Path      string            `json:"path"`
	Data      map[string]string `json:"data"`
	CapturedAt time.Time        `json:"captured_at"`
	Version   int               `json:"version"`
}

// SnapshotResult holds the outcome of a snapshot operation.
type SnapshotResult struct {
	Snapshot *Snapshot
	Skipped  bool
	Reason   string
}

// Summary returns a human-readable description of the snapshot result.
func (r SnapshotResult) Summary() string {
	if r.Skipped {
		return fmt.Sprintf("snapshot skipped: %s", r.Reason)
	}
	return fmt.Sprintf("snapshot captured %d keys from %s at version %d",
		len(r.Snapshot.Data), r.Snapshot.Path, r.Snapshot.Version)
}

// TakeSnapshot reads all secret data at the given path and returns a Snapshot.
func TakeSnapshot(client *http.Client, baseURL, token, path string) (SnapshotResult, error) {
	url := fmt.Sprintf("%s/v1/%s", baseURL, path)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return SnapshotResult{}, fmt.Errorf("snapshot: build request: %w", err)
	}
	req.Header.Set("X-Vault-Token", token)

	resp, err := client.Do(req)
	if err != nil {
		return SnapshotResult{}, fmt.Errorf("snapshot: request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return SnapshotResult{Skipped: true, Reason: "path not found"}, nil
	}
	if resp.StatusCode != http.StatusOK {
		return SnapshotResult{}, fmt.Errorf("snapshot: unexpected status %d", resp.StatusCode)
	}

	var body struct {
		Data struct {
			Data     map[string]string `json:"data"`
			Metadata struct {
				Version int `json:"version"`
			} `json:"metadata"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return SnapshotResult{}, fmt.Errorf("snapshot: decode response: %w", err)
	}

	snap := &Snapshot{
		Path:       path,
		Data:       body.Data.Data,
		CapturedAt: time.Now().UTC(),
		Version:    body.Data.Metadata.Version,
	}
	return SnapshotResult{Snapshot: snap}, nil
}
