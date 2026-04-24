package vault

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// PinResult holds the outcome of a pin or unpin operation.
type PinResult struct {
	Path    string
	Version int
	Pinned  bool
}

func (r PinResult) Summary() string {
	if r.Pinned {
		return fmt.Sprintf("pinned %s at version %d", r.Path, r.Version)
	}
	return fmt.Sprintf("unpinned %s (was version %d)", r.Path, r.Version)
}

// PinSecret marks a specific version of a secret as pinned via custom metadata.
func (c *Client) PinSecret(path string, version int) (PinResult, error) {
	metaPath := toMetadataPath(path, c.cfg.KVVersion)
	payload := map[string]interface{}{
		"custom_metadata": map[string]string{
			"pinned_version": fmt.Sprintf("%d", version),
		},
	}
	body, _ := json.Marshal(payload)
	resp, err := c.postJSON(metaPath, body)
	if err != nil {
		return PinResult{}, fmt.Errorf("pin secret: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusOK {
		return PinResult{}, fmt.Errorf("pin secret: unexpected status %d", resp.StatusCode)
	}
	return PinResult{Path: path, Version: version, Pinned: true}, nil
}

// UnpinSecret removes the pinned_version custom metadata from a secret.
func (c *Client) UnpinSecret(path string, currentVersion int) (PinResult, error) {
	metaPath := toMetadataPath(path, c.cfg.KVVersion)
	payload := map[string]interface{}{
		"custom_metadata": map[string]string{
			"pinned_version": "",
		},
	}
	body, _ := json.Marshal(payload)
	resp, err := c.postJSON(metaPath, body)
	if err != nil {
		return PinResult{}, fmt.Errorf("unpin secret: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusOK {
		return PinResult{}, fmt.Errorf("unpin secret: unexpected status %d", resp.StatusCode)
	}
	return PinResult{Path: path, Version: currentVersion, Pinned: false}, nil
}

// GetPinnedVersion reads the pinned_version from a secret's custom metadata.
// Returns 0 if no pin is set.
func (c *Client) GetPinnedVersion(path string) (int, error) {
	meta, err := c.GetSecretMetadata(path)
	if err != nil {
		return 0, fmt.Errorf("get pinned version: %w", err)
	}
	raw, ok := meta.CustomMetadata["pinned_version"]
	if !ok || raw == "" {
		return 0, nil
	}
	var v int
	_, err = fmt.Sscanf(raw, "%d", &v)
	if err != nil {
		return 0, fmt.Errorf("parse pinned version: %w", err)
	}
	return v, nil
}
