package sync

import (
	"context"
	"fmt"
	"strings"
)

// RunRecursive fetches all secrets under the configured path recursively
// and merges them into the target .env file. Keys are derived from the
// secret path by replacing slashes with underscores and uppercasing.
func (s *Syncer) RunRecursive(ctx context.Context) error {
	paths, err := s.vault.RecursiveList(ctx, s.mount, s.prefix)
	if err != nil {
		return fmt.Errorf("recursive list: %w", err)
	}

	for _, p := range paths {
		secrets, err := s.vault.GetSecret(ctx, s.mount, p)
		if err != nil {
			return fmt.Errorf("get secret %q: %w", p, err)
		}

		prefixKey := pathToEnvPrefix(p)
		prefixed := make(map[string]string, len(secrets))
		for k, v := range secrets {
			prefixed[prefixKey+"_"+strings.ToUpper(k)] = v
		}

		if err := s.merge(ctx, prefixed); err != nil {
			return err
		}
	}
	return nil
}

// pathToEnvPrefix converts a vault path like "app/config" to "APP_CONFIG".
func pathToEnvPrefix(path string) string {
	return strings.ToUpper(strings.ReplaceAll(strings.Trim(path, "/"), "/", "_"))
}
