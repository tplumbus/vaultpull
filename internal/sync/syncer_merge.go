// Package sync orchestrates fetching secrets from Vault and writing them locally.
package sync

import (
	"fmt"
	"os"

	"github.com/user/vaultpull/internal/env"
)

// RunWithMerge fetches secrets from Vault, merges them with the existing .env
// file (if present), and writes the result. Returns a merge summary string.
func (s *Syncer) RunWithMerge(overwrite bool) (string, error) {
	secrets, err := s.vault.GetSecret(s.path)
	if err != nil {
		return "", fmt.Errorf("vault get secret: %w", err)
	}

	reader := env.NewReader(s.output)
	existing, readErr := reader.Read()
	if readErr != nil && !os.IsNotExist(readErr) {
		return "", fmt.Errorf("reading existing env file: %w", readErr)
	}
	if existing == nil {
		existing = map[string]string{}
	}

	merger := env.NewMerger(overwrite)
	merged, report := merger.Merge(existing, secrets)

	writer := env.NewWriter(s.output)
	if err := writer.Write(merged); err != nil {
		return "", fmt.Errorf("writing env file: %w", err)
	}

	return report.Summary(), nil
}
