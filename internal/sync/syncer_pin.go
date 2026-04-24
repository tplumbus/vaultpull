package sync

import (
	"fmt"
	"io"
	"os"
)

// PinOptions configures a pin or unpin operation.
type PinOptions struct {
	Path    string
	Version int
	Unpin   bool
	Output  io.Writer
}

// RunPin pins or unpins a specific version of a Vault secret.
func (s *Syncer) RunPin(opts PinOptions) error {
	out := opts.Output
	if out == nil {
		out = os.Stdout
	}

	if opts.Path == "" {
		return fmt.Errorf("pin: path is required")
	}

	if opts.Unpin {
		current, err := s.client.GetPinnedVersion(opts.Path)
		if err != nil {
			return fmt.Errorf("pin: failed to get current pin for %q: %w", opts.Path, err)
		}
		res, err := s.client.UnpinSecret(opts.Path, current)
		if err != nil {
			return fmt.Errorf("pin: unpin failed for %q: %w", opts.Path, err)
		}
		fmt.Fprintf(out, "✓ %s\n", res.Summary())
		return nil
	}

	if opts.Version <= 0 {
		return fmt.Errorf("pin: version must be a positive integer")
	}

	res, err := s.client.PinSecret(opts.Path, opts.Version)
	if err != nil {
		return fmt.Errorf("pin: failed for %q: %w", opts.Path, err)
	}
	fmt.Fprintf(out, "✓ %s\n", res.Summary())
	return nil
}

// RunGetPin prints the currently pinned version for a secret path.
func (s *Syncer) RunGetPin(path string, out io.Writer) error {
	if out == nil {
		out = os.Stdout
	}
	if path == "" {
		return fmt.Errorf("get-pin: path is required")
	}
	v, err := s.client.GetPinnedVersion(path)
	if err != nil {
		return fmt.Errorf("get-pin: %w", err)
	}
	if v == 0 {
		fmt.Fprintf(out, "%s is not pinned\n", path)
	} else {
		fmt.Fprintf(out, "%s is pinned at version %d\n", path, v)
	}
	return nil
}
