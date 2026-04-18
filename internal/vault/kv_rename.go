package vault

import (
	"fmt"
)

// RenameSecret copies a secret from src to dst, then deletes the source.
// If overwrite is false and dst already exists, it returns an error.
func (c *Client) RenameSecret(mount, src, dst string, overwrite bool) error {
	if err := c.CopySecret(mount, src, dst, overwrite); err != nil {
		return fmt.Errorf("rename: copy failed: %w", err)
	}
	if err := c.DeleteSecret(mount, src); err != nil {
		return fmt.Errorf("rename: delete source failed: %w", err)
	}
	return nil
}

// DeleteSecret removes a secret at the given path under mount.
func (c *Client) DeleteSecret(mount, path string) error {
	fullPath := fmt.Sprintf("/v1/%s/data/%s", mount, path)
	resp, err := c.http.Delete(c.addr + fullPath)
	if err != nil {
		return fmt.Errorf("delete secret: request failed: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode == 404 {
		return fmt.Errorf("delete secret: path not found: %s", path)
	}
	if resp.StatusCode >= 300 {
		return fmt.Errorf("delete secret: unexpected status %d", resp.StatusCode)
	}
	return nil
}
