package vault

import (
	"context"
	"fmt"
)

// CopySecret copies all key-value pairs from one secret path to another.
// If overwrite is false, existing keys at the destination are preserved.
func (c *Client) CopySecret(ctx context.Context, src, dst string, overwrite bool) (int, error) {
	srcData, err := c.GetSecret(ctx, src)
	if err != nil {
		return 0, fmt.Errorf("read source %q: %w", src, err)
	}

	var dstData map[string]string
	if !overwrite {
		dstData, err = c.GetSecret(ctx, dst)
		if err != nil && !isNotFound(err) {
			return 0, fmt.Errorf("read destination %q: %w", dst, err)
		}
	}

	merged := make(map[string]string, len(srcData))
	for k, v := range srcData {
		merged[k] = v
	}
	if !overwrite {
		for k, v := range dstData {
			merged[k] = v
		}
	}

	if err := c.WriteSecret(ctx, dst, merged); err != nil {
		return 0, fmt.Errorf("write destination %q: %w", dst, err)
	}

	return len(srcData), nil
}

// isNotFound returns true if the error represents a 404 from Vault.
func isNotFound(err error) bool {
	if err == nil {
		return false
	}
	return err.Error() == ErrNotFound.Error()
}
