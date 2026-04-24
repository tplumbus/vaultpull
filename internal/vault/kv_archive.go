package vault

import (
	"fmt"
	"time"
)

// ArchiveEntry represents a single archived version of a secret.
type ArchiveEntry struct {
	Path      string
	Version   int
	ArchivedAt time.Time
	Data      map[string]string
}

// ArchiveResult holds the result of an archive operation.
type ArchiveResult struct {
	Path    string
	Version int
	OK      bool
}

// Summary returns a human-readable summary of the archive result.
func (r ArchiveResult) Summary() string {
	if r.OK {
		return fmt.Sprintf("archived %s at version %d", r.Path, r.Version)
	}
	return fmt.Sprintf("failed to archive %s", r.Path)
}

// ArchiveSecret reads the current secret at path and stores it as an
// archived entry by writing it to a versioned archive path under
// "archive/<path>/v<version>".
func ArchiveSecret(c *Client, path string) (ArchiveResult, error) {
	data, err := c.GetSecret(path)
	if err != nil {
		return ArchiveResult{Path: path}, fmt.Errorf("archive: read %s: %w", path, err)
	}

	meta, err := c.GetSecretMetadata(path)
	if err != nil {
		return ArchiveResult{Path: path}, fmt.Errorf("archive: metadata %s: %w", path, err)
	}

	archivePath := fmt.Sprintf("archive/%s/v%d", path, meta.CurrentVersion)

	if err := c.WriteSecret(archivePath, data); err != nil {
		return ArchiveResult{Path: path, Version: meta.CurrentVersion},
			fmt.Errorf("archive: write %s: %w", archivePath, err)
	}

	return ArchiveResult{
		Path:    path,
		Version: meta.CurrentVersion,
		OK:      true,
	}, nil
}

// ListArchives returns all archived versions stored under "archive/<path>/".
func ListArchives(c *Client, path string) ([]string, error) {
	archiveRoot := fmt.Sprintf("archive/%s", path)
	keys, err := c.ListSecrets(archiveRoot)
	if err != nil {
		return nil, fmt.Errorf("list archives %s: %w", path, err)
	}
	return keys, nil
}
