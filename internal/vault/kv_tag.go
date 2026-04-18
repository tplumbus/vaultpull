package vault

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Tag represents a key-value label attached to a secret path.
type Tag struct {
	Key   string
	Value string
}

// GetTags retrieves custom metadata tags for a secret at the given path.
func (c *Client) GetTags(path string) (map[string]string, error) {
	metaPath := toMetadataPath(path, c.kvVersion)
	req, err := http.NewRequest(http.MethodGet, c.address+"/v1/"+metaPath, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-Vault-Token", c.token)

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
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
			CustomMetadata map[string]string `json:"custom_metadata"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return nil, err
	}
	if body.Data.CustomMetadata == nil {
		return map[string]string{}, nil
	}
	return body.Data.CustomMetadata, nil
}
