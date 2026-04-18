package vault

import (
	"context"
	"fmt"
)

// RestoreSecret writes a specific version's data back as the latest version,
// effectively restoring it. It reads the target version then writes it again.
func (c *Client) RestoreSecret(ctx context.Context, path string, version int) (map[string]string, error) {
	data, err := c.GetSecretVersion(ctx, path, version)
	if err != nil {
		return nil, fmt.Errorf("restore: read version %d: %w", version, err)
	}
	if len(data) == 0 {
		return nil, fmt.Errorf("restore: version %d is empty or deleted", version)
	}
	err = c.WriteSecret(ctx, path, data)
	if err != nil {
		return nil, fmt.Errorf("restore: write: %w", err)
	}
	return data, nil
}

// GetSecretVersion fetches a specific KV v2 version of a secret.
func (c *Client) GetSecretVersion(ctx context.Context, path string, version int) (map[string]string, error) {
	mountPath, secretPath := splitMount(path)
	url := fmt.Sprintf("%s/v1/%s/data/%s?version=%d", c.addr, mountPath, secretPath, version)
	return c.getSecretByURL(ctx, url)
}

// WriteSecret writes key-value pairs to a KV v2 path.
func (c *Client) WriteSecret(ctx context.Context, path string, data map[string]string) error {
	mountPath, secretPath := splitMount(path)
	url := fmt.Sprintf("%s/v1/%s/data/%s", c.addr, mountPath, secretPath)
	return c.writeSecretByURL(ctx, url, data)
}
