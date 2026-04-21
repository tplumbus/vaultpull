package vault

import (
	"fmt"
	"strings"
)

// SearchResult holds a matched secret path and the matching keys.
type SearchResult struct {
	Path string
	Keys []string
}

// SearchSummary returns a human-readable summary of search results.
func SearchSummary(results []SearchResult) string {
	if len(results) == 0 {
		return "no matches found"
	}
	total := 0
	for _, r := range results {
		total += len(r.Keys)
	}
	return fmt.Sprintf("%d match(es) across %d path(s)", total, len(results))
}

// SearchSecrets searches all secrets under basePath whose keys or values
// contain the given query string (case-insensitive).
func (c *Client) SearchSecrets(basePath, query string) ([]SearchResult, error) {
	paths, err := RecursiveList(c, basePath)
	if err != nil {
		return nil, fmt.Errorf("search: list failed: %w", err)
	}

	q := strings.ToLower(query)
	var results []SearchResult

	for _, p := range paths {
		data, err := c.GetSecret(p)
		if err != nil {
			continue
		}
		var matched []string
		for k, v := range data {
			if strings.Contains(strings.ToLower(k), q) ||
				strings.Contains(strings.ToLower(fmt.Sprintf("%v", v)), q) {
				matched = append(matched, k)
			}
		}
		if len(matched) > 0 {
			results = append(results, SearchResult{Path: p, Keys: matched})
		}
	}
	return results, nil
}
