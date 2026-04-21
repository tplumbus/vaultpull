package sync

import (
	"fmt"
	"io"

	"github.com/yourusername/vaultpull/internal/audit"
	"github.com/yourusername/vaultpull/internal/vault"
)

// SearchOptions configures a secret search operation.
type SearchOptions struct {
	BasePath string
	Query    string
	LogOut   io.Writer
}

// RunSearch searches Vault secrets under BasePath for keys or values
// matching Query, logs results, and returns the matched results.
func (s *Syncer) RunSearch(opts SearchOptions) ([]vault.SearchResult, error) {
	if opts.Query == "" {
		return nil, fmt.Errorf("search query must not be empty")
	}
	if opts.BasePath == "" {
		return nil, fmt.Errorf("base path must not be empty")
	}

	results, err := s.client.SearchSecrets(opts.BasePath, opts.Query)
	if err != nil {
		return nil, fmt.Errorf("RunSearch: %w", err)
	}

	logger := audit.NewSearchLogger(opts.LogOut)
	total := 0
	for _, r := range results {
		logger.LogResult(r.Path, r.Keys)
		total += len(r.Keys)
	}
	logger.LogSummary(opts.Query, total)

	return results, nil
}
