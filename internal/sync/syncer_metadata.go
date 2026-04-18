package sync

import (
	"fmt"

	"github.com/your-org/vaultpull/internal/audit"
)

// metadataFetcher is satisfied by vault.Client.
type metadataFetcher interface {
	GetSecretMetadata(path string) (interface{ GetVersion() int }, error)
}

// LogMetadataForPaths fetches and logs metadata for each secret path.
// It uses the audit MetadataLogger to record version info.
func (s *Syncer) LogMetadataForPaths(paths []string, logger *audit.MetadataLogger) error {
	for _, p := range paths {
		meta, err := s.vault.GetSecretMetadata(p)
		if err != nil {
			return fmt.Errorf("metadata for %s: %w", p, err)
		}
		logger.Log(audit.MetadataEntry{
			Path:           p,
			CurrentVersion: meta.CurrentVersion,
			UpdatedTime:    meta.UpdatedTime,
		})
	}
	return nil
}
