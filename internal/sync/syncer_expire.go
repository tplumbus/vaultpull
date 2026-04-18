package sync

import (
	"fmt"
	"io"
	"os"
	"time"
)

// ExpiryResult holds the outcome of an expiry check for a single secret.
type ExpiryResult struct {
	Path    string
	Expired bool
	TTL     time.Duration
}

// ExpiryChecker defines what we need from the vault client.
type ExpiryChecker interface {
	GetExpiry(path string) (interface{ IsExpired() bool; GetTTL() time.Duration }, error)
}

// RunCheckExpiry checks expiry for each path and prints a summary.
func (s *Syncer) RunCheckExpiry(paths []string, out io.Writer) ([]ExpiryResult, error) {
	if out == nil {
		out = os.Stdout
	}
	var results []ExpiryResult
	for _, p := range paths {
		exp, err := s.vault.GetExpiry(p)
		if err != nil {
			fmt.Fprintf(out, "WARN: could not read expiry for %s: %v\n", p, err)
			continue
		}
		r := ExpiryResult{
			Path:    p,
			Expired: exp.Expired,
			TTL:     exp.TTL,
		}
		results = append(results, r)
		if exp.Expired {
			fmt.Fprintf(out, "EXPIRED  %s (expired %s ago)\n", p, (-exp.TTL).Round(time.Second))
		} else if exp.ExpiresAt.IsZero() {
			fmt.Fprintf(out, "NO_EXPIRY %s\n", p)
		} else {
			fmt.Fprintf(out, "OK       %s (expires in %s)\n", p, exp.TTL.Round(time.Second))
		}
	}
	return results, nil
}

// RunSetExpiry sets a TTL on each of the given paths.
func (s *Syncer) RunSetExpiry(paths []string, ttl time.Duration, out io.Writer) error {
	if out == nil {
		out = os.Stdout
	}
	for _, p := range paths {
		if err := s.vault.SetExpiry(p, ttl); err != nil {
			return fmt.Errorf("set expiry %s: %w", p, err)
		}
		fmt.Fprintf(out, "SET_EXPIRY %s ttl=%s\n", p, ttl)
	}
	return nil
}
