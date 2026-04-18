package vault

import (
	"context"
	"fmt"
	"net/http"
)

// RollbackSecret restores a secret to a previous version.
// It reads the data from the specified version and writes it back as a new version.
func (c *Client) RollbackSecret(ctx context.Context, path string, version int) error {
	metaPath := toMetadataPath(path)
	_ = metaPath

	// Read the specific version
	versionedPath := fmt.Sprintf("%s?version=%d", path, version)
	data, err := c.GetSecret(ctx, versionedPath)
	if err != nil {
		return fmt.Errorf("rollback: read version %d: %w", version, err)
	}
	if len(data) == 0 {
		return fmt.Errorf("rollback: version %d not found or empty at %s", version, path)
	}

	// Write the old data back as a new version
	if err := c.WriteSecret(ctx, path, data); err != nil {
		return fmt.Errorf("rollback: write restored data: %w", err)
	}
	return nil
}

// WriteSecret writes key-value pairs to a KV v2 secret path.
func (c *Client) WriteSecret(ctx context.Context, path string, data map[string]string) error {
	body := map[string]interface{}{
		"data": data,
	}

	req, err := c.newJSONRequest(ctx, http.MethodPost, path, body)
	if err != nil {
		return err
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("vault write %s: unexpected status %d", path, resp.StatusCode)
	}
	return nil
}
