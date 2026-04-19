package sync

import (
	"fmt"
	"io"
	"os"

	"github.com/user/vaultpull/internal/vault"
)

// PromoteOptions configures a secret promotion run.
type PromoteOptions struct {
	SrcPath   string
	DstPath   string
	Overwrite bool
	Out       io.Writer
}

// RunPromote copies secrets from a source vault path to a destination path,
// optionally replacing an environment segment in the path.
func RunPromote(client *vault.Client, opts PromoteOptions) error {
	out := opts.Out
	if out == nil {
		out = os.Stdout
	}

	if opts.SrcPath == "" {
		return fmt.Errorf("promote: src path is required")
	}
	if opts.DstPath == "" {
		return fmt.Errorf("promote: dst path is required")
	}

	fmt.Fprintf(out, "Promoting %q -> %q (overwrite=%v)\n", opts.SrcPath, opts.DstPath, opts.Overwrite)

	if err := client.PromoteSecret(opts.SrcPath, opts.DstPath, opts.Overwrite); err != nil {
		return fmt.Errorf("promote failed: %w", err)
	}

	fmt.Fprintf(out, "Promotion complete: %q\n", opts.DstPath)
	return nil
}
