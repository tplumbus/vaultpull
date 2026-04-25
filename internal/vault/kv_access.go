package vault

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strings"
)

// AccessPolicy represents a named set of capabilities on a Vault path.
type AccessPolicy struct {
	Path         string   `json:"path"`
	Capabilities []string `json:"capabilities"`
}

// AccessResult holds the resolved policies for a secret path.
type AccessResult struct {
	SecretPath string
	Policies   []AccessPolicy
}

// Summary returns a human-readable summary of the access result.
func (a AccessResult) Summary() string {
	if len(a.Policies) == 0 {
		return fmt.Sprintf("no policies found for %s", a.SecretPath)
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "access policies for %s:\n", a.SecretPath)
	for _, p := range a.Policies {
		sort.Strings(p.Capabilities)
		fmt.Fprintf(&sb, "  [%s] %s\n", p.Path, strings.Join(p.Capabilities, ", "))
	}
	return strings.TrimRight(sb.String(), "\n")
}

// GetAccessPolicies queries Vault's sys/internal/ui/resultant-acl endpoint
// to retrieve the effective capabilities for the given secret path.
func (c *Client) GetAccessPolicies(secretPath string) (AccessResult, error) {
	url := fmt.Sprintf("%s/v1/sys/internal/ui/resultant-acl", c.addr)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return AccessResult{}, fmt.Errorf("build request: %w", err)
	}
	req.Header.Set("X-Vault-Token", c.token)

	resp, err := c.http.Do(req)
	if err != nil {
		return AccessResult{}, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return AccessResult{SecretPath: secretPath}, nil
	}
	if resp.StatusCode != http.StatusOK {
		return AccessResult{}, fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	var body struct {
		Data struct {
			Exact map[string]struct {
				Capabilities []string `json:"capabilities"`
			} `json:"exact_rules"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return AccessResult{}, fmt.Errorf("decode response: %w", err)
	}

	var policies []AccessPolicy
	for path, rule := range body.Data.Exact {
		if strings.Contains(secretPath, strings.TrimSuffix(path, "*")) || path == secretPath {
			policies = append(policies, AccessPolicy{Path: path, Capabilities: rule.Capabilities})
		}
	}
	sort.Slice(policies, func(i, j int) bool { return policies[i].Path < policies[j].Path })

	return AccessResult{SecretPath: secretPath, Policies: policies}, nil
}
