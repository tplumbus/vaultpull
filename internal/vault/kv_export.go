package vault

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// ExportResult holds the exported secrets and metadata.
type ExportResult struct {
	Path    string
	Secrets map[string]string
	Version int
}

// Summary returns a human-readable summary of the export.
func (r ExportResult) Summary() string {
	return fmt.Sprintf("exported %d keys from %s at version %d", len(r.Secrets), r.Path, r.Version)
}

// ExportSecret reads all key-value pairs from a Vault KV path and returns them
// as an ExportResult suitable for serialisation or writing to disk.
func ExportSecret(client *Client, mount, path string) (ExportResult, error) {
	url := fmt.Sprintf("%s/v1/%s/data/%s", client.Address, mount, path)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return ExportResult{}, fmt.Errorf("export: build request: %w", err)
	}
	req.Header.Set("X-Vault-Token", client.Token)

	resp, err := client.HTTP.Do(req)
	if err != nil {
		return ExportResult{}, fmt.Errorf("export: request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return ExportResult{}, fmt.Errorf("export: path not found: %s", path)
	}
	if resp.StatusCode != http.StatusOK {
		return ExportResult{}, fmt.Errorf("export: unexpected status %d", resp.StatusCode)
	}

	var body struct {
		Data struct {
			Data     map[string]string `json:"data"`
			Metadata struct {
				Version int `json:"version"`
			} `json:"metadata"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return ExportResult{}, fmt.Errorf("export: decode response: %w", err)
	}

	return ExportResult{
		Path:    path,
		Secrets: body.Data.Data,
		Version: body.Data.Metadata.Version,
	}, nil
}
