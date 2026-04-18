package vault

import (
	"context"
	"fmt"
	"strings"
)

// RecursiveList walks a KV mount path and returns all secret leaf paths.
func (c *Client) RecursiveList(ctx context.Context, mount, prefix string) ([]string, error) {
	keys, err := c.ListSecrets(ctx, mount, prefix)
	if err != nil {
		return nil, err
	}

	var results []string
	for _, key := range keys {
		if strings.HasSuffix(key, "/") {
			sub := fmt.Sprintf("%s/%s", strings.TrimRight(prefix, "/"), strings.TrimRight(key, "/"))
			sub = strings.TrimPrefix(sub, "/")
			children, err := c.RecursiveList(ctx, mount, sub)
			if err != nil {
				return nil, err
			}
			results = append(results, children...)
		} else {
			leaf := key
			if prefix != "" {
				leaf = strings.TrimRight(prefix, "/") + "/" + key
			}
			results = append(results, leaf)
		}
	}
	return results, nil
}
