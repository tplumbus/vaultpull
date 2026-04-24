package sync

import (
	"fmt"
	"io"
	"os"

	"github.com/yourusername/vaultpull/internal/vault"
)

// RunArchive archives the current version of the secret at path and
// prints a summary to w (defaults to os.Stdout if nil).
func (s *Syncer) RunArchive(path string, w io.Writer) error {
	if w == nil {
		w = os.Stdout
	}

	res, err := vault.ArchiveSecret(s.client, path)
	if err != nil {
		return fmt.Errorf("RunArchive: %w", err)
	}

	fmt.Fprintln(w, res.Summary())
	return nil
}

// RunListArchives lists all archived versions for the secret at path
// and prints each entry to w (defaults to os.Stdout if nil).
func (s *Syncer) RunListArchives(path string, w io.Writer) error {
	if w == nil {
		w = os.Stdout
	}

	keys, err := vault.ListArchives(s.client, path)
	if err != nil {
		return fmt.Errorf("RunListArchives: %w", err)
	}

	if len(keys) == 0 {
		fmt.Fprintf(w, "no archives found for %s\n", path)
		return nil
	}

	fmt.Fprintf(w, "archives for %s:\n", path)
	for _, k := range keys {
		fmt.Fprintf(w, "  %s\n", k)
	}
	return nil
}
