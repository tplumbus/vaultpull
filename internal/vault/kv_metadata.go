package vault

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// SecretMetadata holds version and timing info for a KV secret.
type SecretMetadata struct {
	CreatedTime    time.Time
	UpdatedTime    time.Time
	CurrentVersion int
	OldestVersion  int
}

// GetSecretMetadata fetches metadata for a KV v2 secret path.
func (c *Client) GetSecretMetadata(path string) (*SecretMetadata, error) {
	metaPath := toMetadataPath(path)
	req, err := http.NewRequest(http.MethodGet, c.addr+"/v1/"+metaPath, nil)
	if err != nil {
		return nil, fmt.Errorf("build request: %w", err)
	}
	req.Header.Set("X-Vault-Token", c.token)

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("secret not found: %s", path)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	var body struct {
		Data struct {
			CreatedTime    time.Time `json:"created_time"`
			UpdatedTime    time.Time `json:"updated_time"`
			CurrentVersion int       `json:"current_version"`
			OldestVersion  int       `json:"oldest_version"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	return &SecretMetadata{
		CreatedTime:    body.Data.CreatedTime,
		UpdatedTime:    body.Data.UpdatedTime,
		CurrentVersion: body.Data.CurrentVersion,
		OldestVersion:  body.Data.OldestVersion,
	}, nil
}
