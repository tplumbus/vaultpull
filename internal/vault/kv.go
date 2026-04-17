package vault

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// KVVersion represents the KV secrets engine version.
type KVVersion int

const (
	KVv1 KVVersion = 1
	KVv2 KVVersion = 2
)

// ListSecrets returns the keys available at the given path.
func (c *Client) ListSecrets(path string) ([]string, error) {
	listPath := path
	if c.kvVersion == KVv2 {
		listPath = toMetadataPath(path)
	}

	req, err := http.NewRequest("LIST", fmt.Sprintf("%s/v1/%s", c.addr, listPath), nil)
	if err != nil {
		return nil, fmt.Errorf("build list request: %w", err)
	}
	req.Header.Set("X-Vault-Token", c.token)

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("list request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("path not found: %s", path)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read body: %w", err)
	}

	var result struct {
		Data struct {
			Keys []string `json:"keys"`
		} `json:"data"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	return result.Data.Keys, nil
}

// toMetadataPath converts a KVv2 data path to its metadata equivalent.
func toMetadataPath(path string) string {
	// e.g. secret/data/myapp -> secret/metadata/myapp
	const dataSegment = "data/"
	for i := 0; i < len(path)-len(dataSegment); i++ {
		if path[i:i+len(dataSegment)] == dataSegment {
			return path[:i] + "metadata/" + path[i+len(dataSegment):]
		}
	}
	return path
}
