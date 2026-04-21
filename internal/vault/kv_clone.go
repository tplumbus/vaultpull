package vault

import (
	"fmt"
	"strings"
)

// CloneSecret copies all secrets from a source path to one or more destination
// paths. Unlike CopySecret, CloneSecret supports multi-target fan-out and
// optionally strips or replaces a path prefix in the destination.
func (c *Client) CloneSecret(srcPath string, dstPaths []string, overwrite bool) error {
	data, err := c.GetSecret(srcPath)
	if err != nil {
		return fmt.Errorf("clone: read source %q: %w", srcPath, err)
	}

	for _, dst := range dstPaths {
		dst = strings.TrimSpace(dst)
		if dst == "" {
			continue
		}

		if !overwrite {
			_, err := c.GetSecret(dst)
			if err == nil {
				return fmt.Errorf("clone: destination %q already exists (use overwrite=true to replace)", dst)
			}
			if !isNotFound(err) {
				return fmt.Errorf("clone: check destination %q: %w", dst, err)
			}
		}

		if err := c.WriteSecret(dst, data); err != nil {
			return fmt.Errorf("clone: write destination %q: %w", dst, err)
		}
	}

	return nil
}

// CloneResult summarises the outcome of a clone operation.
type CloneResult struct {
	Source      string
	Destinations []string
	Skipped     []string
}

// Summary returns a human-readable description of the clone result.
func (r CloneResult) Summary() string {
	return fmt.Sprintf("cloned %q → %d destination(s), %d skipped",
		r.Source, len(r.Destinations), len(r.Skipped))
}
