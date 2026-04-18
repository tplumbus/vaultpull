package sync

import (
	"context"
	"fmt"
	"io"
	"os"
)

// RestoreOptions configures a secret restore operation.
type RestoreOptions struct {
	Path    string
	Version int
	Output  string
	Out     io.Writer
}

// RunRestore fetches a specific version of a Vault secret and writes it to the
// output .env file, useful for rolling back to a known-good state locally.
func (s *Syncer) RunRestore(ctx context.Context, opts RestoreOptions) error {
	if opts.Path == "" {
		return fmt.Errorf("restore: path is required")
	}
	if opts.Version <= 0 {
		return fmt.Errorf("restore: version must be >= 1")
	}
	out := opts.Out
	if out == nil {
		out = os.Stdout
	}

	data, err := s.vault.RestoreSecret(ctx, opts.Path, opts.Version)
	if err != nil {
		return fmt.Errorf("restore: %w", err)
	}

	outPath := opts.Output
	if outPath == "" {
		outPath = ".env"
	}

	w := s.writer
	if err := w.Write(outPath, data); err != nil {
		return fmt.Errorf("restore: write env: %w", err)
	}

	fmt.Fprintf(out, "restored %d keys from %s@v%d -> %s\n", len(data), opts.Path, opts.Version, outPath)
	return nil
}
