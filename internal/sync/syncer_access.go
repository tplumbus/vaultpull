package sync

import (
	"fmt"
	"io"
	"os"
)

// accessChecker abstracts the vault client method used by RunAccess.
type accessChecker interface {
	GetAccessPolicies(secretPath string) (interface{ Summary() string }, error)
}

// RunAccessOptions configures the RunAccess operation.
type RunAccessOptions struct {
	SecretPath string
	Out        io.Writer
}

// RunAccess retrieves and prints the effective access policies for a Vault
// secret path. It writes a human-readable summary to opts.Out (defaults to
// os.Stdout when nil).
func (s *Syncer) RunAccess(opts RunAccessOptions) error {
	if opts.SecretPath == "" {
		return fmt.Errorf("secret path must not be empty")
	}
	out := opts.Out
	if out == nil {
		out = os.Stdout
	}

	result, err := s.vault.GetAccessPolicies(opts.SecretPath)
	if err != nil {
		return fmt.Errorf("get access policies: %w", err)
	}

	fmt.Fprintln(out, result.Summary())
	return nil
}
