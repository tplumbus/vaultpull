package sync

import (
	"fmt"
	"io"
	"os"
)

// RenameOptions holds parameters for a rename operation.
type RenameOptions struct {
	Mount     string
	SrcPath   string
	DstPath   string
	Overwrite bool
	Out       io.Writer
}

type secretRenamer interface {
	RenameSecret(mount, src, dst string, overwrite bool) error
}

// RunRename renames a secret in Vault from src to dst.
func RunRename(v secretRenamer, opts RenameOptions) error {
	out := opts.Out
	if out == nil {
		out = os.Stdout
	}
	if opts.SrcPath == "" {
		return fmt.Errorf("rename: source path is required")
	}
	if opts.DstPath == "" {
		return fmt.Errorf("rename: destination path is required")
	}
	if err := v.RenameSecret(opts.Mount, opts.SrcPath, opts.DstPath, opts.Overwrite); err != nil {
		return fmt.Errorf("rename failed: %w", err)
	}
	fmt.Fprintf(out, "renamed %s -> %s\n", opts.SrcPath, opts.DstPath)
	return nil
}
