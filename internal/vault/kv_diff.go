package vault

import "fmt"

// DiffResult holds the comparison between remote Vault secrets and local env values.
type DiffResult struct {
	Added   map[string]string // keys present in Vault but not locally
	Changed map[string]string // keys present in both but with different values
	Removed []string          // keys present locally but not in Vault
}

// Diff compares remote secrets from Vault against local env entries.
// remote is the map fetched from Vault; local is the map read from the .env file.
func Diff(remote, local map[string]string) DiffResult {
	result := DiffResult{
		Added:   make(map[string]string),
		Changed: make(map[string]string),
	}

	for k, rv := range remote {
		lv, exists := local[k]
		if !exists {
			result.Added[k] = rv
		} else if lv != rv {
			result.Changed[k] = rv
		}
	}

	for k := range local {
		if _, exists := remote[k]; !exists {
			result.Removed = append(result.Removed, k)
		}
	}

	return result
}

// HasChanges returns true if there is any difference between remote and local.
func (d DiffResult) HasChanges() bool {
	return len(d.Added) > 0 || len(d.Changed) > 0 || len(d.Removed) > 0
}

// Summary returns a human-readable one-line summary of the diff.
func (d DiffResult) Summary() string {
	return fmt.Sprintf("added=%d changed=%d removed=%d",
		len(d.Added), len(d.Changed), len(d.Removed))
}
