package sync

import (
	"context"
	"fmt"
	"io"
	"os"
)

// RollbackOptions configures a rollback operation.
type RollbackOptions struct {
	VaultPath string
	Version   int
	Out       io.Writer
}

// RunRollback restores a Vault secret at the given path to a prior version.
func (s *Syncer) RunRollback(ctx context.Context, opts RollbackOptions) error {
	if opts.VaultPath == "" {
		return fmt.Errorf("rollback: vault path is required")
	}
	if opts.Version < 1 {
		return fmt.Errorf("rollback: version must be >= 1, got %d", opts.Version)
	}

	out := opts.Out
	if out == nil {
		out = os.Stdout
	}

	fmt.Fprintf(out, "Rolling back %s to version %d...\n", opts.VaultPath, opts.Version)

	if err := s.vault.RollbackSecret(ctx, opts.VaultPath, opts.Version); err != nil {
		return fmt.Errorf("rollback failed: %w", err)
	}

	fmt.Fprintf(out, "Successfully restored %s to version %d\n", opts.VaultPath, opts.Version)
	return nil
}
