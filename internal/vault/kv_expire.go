package vault

import (
	"fmt"
	"time"
)

// SecretExpiry holds TTL metadata for a secret.
type SecretExpiry struct {
	Path      string
	ExpiresAt time.Time
	TTL       time.Duration
	Expired   bool
}

// SetExpiry writes a custom-metadata field "expires_at" to the given secret path.
func (c *Client) SetExpiry(path string, ttl time.Duration) error {
	expiresAt := time.Now().UTC().Add(ttl).Format(time.RFC3339)
	metaPath := toMetadataPath(path, c.kvVersion)
	body := map[string]interface{}{
		"custom_metadata": map[string]interface{}{
			"expires_at": expiresAt,
		},
	}
	_, err := c.client.Logical().Write(metaPath, body)
	if err != nil {
		return fmt.Errorf("set expiry %s: %w", path, err)
	}
	return nil
}

// GetExpiry reads the "expires_at" custom-metadata field for the given path.
func (c *Client) GetExpiry(path string) (*SecretExpiry, error) {
	metaPath := toMetadataPath(path, c.kvVersion)
	secret, err := c.client.Logical().Read(metaPath)
	if err != nil {
		return nil, fmt.Errorf("get expiry %s: %w", path, err)
	}
	if secret == nil {
		return nil, fmt.Errorf("get expiry %s: not found", path)
	}
	cm, ok := secret.Data["custom_metadata"].(map[string]interface{})
	if !ok {
		return &SecretExpiry{Path: path}, nil
	}
	raw, ok := cm["expires_at"].(string)
	if !ok {
		return &SecretExpiry{Path: path}, nil
	}
	t, err := time.Parse(time.RFC3339, raw)
	if err != nil {
		return nil, fmt.Errorf("parse expires_at %s: %w", path, err)
	}
	ttl := time.Until(t)
	return &SecretExpiry{
		Path:      path,
		ExpiresAt: t,
		TTL:       ttl,
		Expired:   ttl <= 0,
	}, nil
}
