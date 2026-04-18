package sync

import (
	"context"
	"fmt"
	"io"
	"os"

	"vaultpull/internal/env"
	"vaultpull/internal/vault"
)

// RunDiff prints a diff of Vault secrets vs the local .env file without writing any changes.
func (s *Syncer) RunDiff(ctx context.Context, w io.Writer) error {
	if w == nil {
		w = os.Stdout
	}

	remote, err := s.client.GetSecret(ctx, s.path)
	if err != nil {
		return fmt.Errorf("vault: %w", err)
	}

	local := make(map[string]string)
	if _, statErr := os.Stat(s.outputPath); statErr == nil {
		r := env.NewReader(s.outputPath)
		local, err = r.Read()
		if err != nil {
			return fmt.Errorf("reading local env: %w", err)
		}
	}

	diff := vault.Diff(remote, local)

	for k, v := range diff.Added {
		fmt.Fprintf(w, "+ %s=%s\n", k, v)
	}
	for k, v := range diff.Changed {
		fmt.Fprintf(w, "~ %s=%s\n", k, v)
	}
	for _, k := range diff.Removed {
		fmt.Fprintf(w, "- %s\n", k)
	}

	if !diff.HasChanges() {
		fmt.Fprintln(w, "no changes detected")
	} else {
		fmt.Fprintf(w, "summary: %s\n", diff.Summary())
	}

	return nil
}
