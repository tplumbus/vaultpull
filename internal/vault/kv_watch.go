package vault

import (
	"context"
	"fmt"
	"time"
)

// WatchResult holds the outcome of a single watch poll cycle.
type WatchResult struct {
	Path    string
	Changed bool
	Version int
	Err     error
}

// WatchOptions configures the Watch behaviour.
type WatchOptions struct {
	// Interval between polls. Defaults to 30s if zero.
	Interval time.Duration
	// KVVersion of the target mount (KVv1 or KVv2).
	KVVersion KVVersion
}

// Watch polls a Vault KV secret at the given path and emits a WatchResult on
// the returned channel whenever the secret version changes or an error occurs.
// The caller must cancel ctx to stop watching.
func Watch(ctx context.Context, c *Client, path string, opts WatchOptions) <-chan WatchResult {
	ch := make(chan WatchResult, 1)

	interval := opts.Interval
	if interval <= 0 {
		interval = 30 * time.Second
	}

	go func() {
		defer close(ch)

		var lastVersion int

		for {
			result := poll(ctx, c, path, opts.KVVersion, &lastVersion)

			if result.Changed || result.Err != nil {
				select {
				case ch <- result:
				case <-ctx.Done():
					return
				}
			}

			select {
			case <-time.After(interval):
			case <-ctx.Done():
				return
			}
		}
	}()

	return ch
}

func poll(ctx context.Context, c *Client, path string, ver KVVersion, lastVersion *int) WatchResult {
	metaPath := toMetadataPath(path, ver)
	data, err := c.GetSecret(ctx, metaPath)
	if err != nil {
		return WatchResult{Path: path, Err: fmt.Errorf("watch poll: %w", err)}
	}

	current := 0
	if v, ok := data["current_version"]; ok {
		switch n := v.(type) {
		case float64:
			current = int(n)
		case int:
			current = n
		}
	}

	changed := current != *lastVersion
	if changed {
		*lastVersion = current
	}

	return WatchResult{Path: path, Changed: changed, Version: current}
}
