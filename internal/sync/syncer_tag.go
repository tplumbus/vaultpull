package sync

import (
	"fmt"

	"github.com/user/vaultpull/internal/audit"
)

// TagFetcher can retrieve tags for a vault secret path.
type TagFetcher interface {
	GetTags(path string) (map[string]string, error)
}

// RunTags fetches and logs tags for each of the given secret paths.
func RunTags(client TagFetcher, paths []string, logger *audit.TagLogger) error {
	if len(paths) == 0 {
		return fmt.Errorf("no paths provided")
	}
	for _, p := range paths {
		tags, err := client.GetTags(p)
		if err != nil {
			return fmt.Errorf("failed to get tags for %s: %w", p, err)
		}
		logger.Log(p, tags)
	}
	return nil
}
