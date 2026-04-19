package vault

import (
	"fmt"
	"strings"
)

// PromoteSecret copies a secret from one environment path prefix to another.
// e.g. promote "secret/dev/app" -> "secret/prod/app"
func (c *Client) PromoteSecret(srcPath, dstPath string, overwrite bool) error {
	src := strings.TrimRight(srcPath, "/")
	dst := strings.TrimRight(dstPath, "/")

	data, err := c.GetSecret(src)
	if err != nil {
		return fmt.Errorf("promote: read src %q: %w", src, err)
	}

	if !overwrite {
		existing, err := c.GetSecret(dst)
		if err == nil && len(existing) > 0 {
			return fmt.Errorf("promote: destination %q already exists; use overwrite=true to force", dst)
		}
	}

	if err := c.WriteSecret(dst, data); err != nil {
		return fmt.Errorf("promote: write dst %q: %w", dst, err)
	}

	return nil
}

// ReplaceEnvInPath swaps the environment segment in a vault path.
// e.g. ReplaceEnvInPath("secret/dev/app/db", "dev", "prod") -> "secret/prod/app/db"
func ReplaceEnvInPath(path, fromEnv, toEnv string) string {
	return strings.Replace(path, "/"+fromEnv+"/", "/"+toEnv+"/", 1)
}
