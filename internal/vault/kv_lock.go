package vault

import (
	"fmt"
	"net/http"
)

// LockSecret sets a custom metadata flag "locked" on a secret to prevent writes.
func (c *Client) LockSecret(path string) error {
	return c.setLockState(path, true)
}

// UnlockSecret removes the locked flag from a secret's custom metadata.
func (c *Client) UnlockSecret(path string) error {
	return c.setLockState(path, false)
}

// IsLocked returns true if the secret at path has the locked flag set.
func (c *Client) IsLocked(path string) (bool, error) {
	meta, err := c.GetSecretMetadata(path)
	if err != nil {
		return false, err
	}
	val, ok := meta.CustomMetadata["locked"]
	if !ok {
		return false, nil
	}
	return val == "true", nil
}

func (c *Client) setLockState(path string, locked bool) error {
	metaPath := toMetadataPath(path)
	value := "false"
	if locked {
		value = "true"
	}
	body := map[string]interface{}{
		"custom_metadata": map[string]string{
			"locked": value,
		},
	}
	resp, err := c.post(metaPath, body)
	if err != nil {
		return fmt.Errorf("lock state update failed: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("unexpected status %d setting lock on %s", resp.StatusCode, path)
	}
	return nil
}
