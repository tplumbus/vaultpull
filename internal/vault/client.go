package vault

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Client is a minimal HashiCorp Vault HTTP client.
type Client struct {
	baseURL    string
	token      string
	httpClient *http.Client
}

// SecretData holds the key/value pairs returned from a Vault KV secret.
type SecretData map[string]string

// NewClient creates a new Vault client.
func NewClient(baseURL, token string) *Client {
	return &Client{
		baseURL: baseURL,
		token:   token,
		httpClient: &http.Client{Timeout: 10 * time.Second},
	}
}

// GetSecret fetches secrets at the given KV v2 path.
func (c *Client) GetSecret(path string) (SecretData, error) {
	url := fmt.Sprintf("%s/v1/%s", c.baseURL, path)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}
	req.Header.Set("X-Vault-Token", c.token)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("executing request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("secret not found at path: %s", path)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status %d for path: %s", resp.StatusCode, path)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response body: %w", err)
	}

	var result struct {
		Data struct {
			Data map[string]string `json:"data"`
		} `json:"data"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("parsing response: %w", err)
	}

	return result.Data.Data, nil
}
