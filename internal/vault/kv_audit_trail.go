package vault

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// AuditEntry represents a single audit trail record for a secret.
type AuditEntry struct {
	Timestamp time.Time         `json:"timestamp"`
	Action    string            `json:"action"`
	Path      string            `json:"path"`
	Version   int               `json:"version"`
	Metadata  map[string]string `json:"metadata,omitempty"`
}

// AuditTrailResult holds the full audit history for a secret path.
type AuditTrailResult struct {
	Path    string       `json:"path"`
	Entries []AuditEntry `json:"entries"`
}

// Summary returns a human-readable summary of the audit trail.
func (r AuditTrailResult) Summary() string {
	return fmt.Sprintf("audit trail for %q: %d entries", r.Path, len(r.Entries))
}

// GetAuditTrail fetches the audit trail for a secret by reading its version
// metadata history from the KV v2 metadata endpoint.
func GetAuditTrail(client *Client, mountPath, secretPath string) (*AuditTrailResult, error) {
	metaPath := toMetadataPath(mountPath, secretPath)
	resp, err := client.RawGet(metaPath)
	if err != nil {
		return nil, fmt.Errorf("audit trail request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("secret not found: %s", secretPath)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status %d fetching audit trail", resp.StatusCode)
	}

	var body struct {
		Data struct {
			Versions map[string]struct {
				CreatedTime  time.Time `json:"created_time"`
				DeletionTime string    `json:"deletion_time"`
				Destroyed    bool      `json:"destroyed"`
			} `json:"versions"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return nil, fmt.Errorf("failed to decode audit trail response: %w", err)
	}

	result := &AuditTrailResult{Path: secretPath}
	for vStr, v := range body.Data.Versions {
		var vNum int
		fmt.Sscanf(vStr, "%d", &vNum)
		action := "created"
		if v.Destroyed {
			action = "destroyed"
		} else if v.DeletionTime != "" {
			action = "deleted"
		}
		result.Entries = append(result.Entries, AuditEntry{
			Timestamp: v.CreatedTime,
			Action:    action,
			Path:      secretPath,
			Version:   vNum,
		})
	}
	return result, nil
}
