// Package env provides utilities for reading and writing .env files.
package env

import "fmt"

// MergeResult holds the outcome of a merge operation.
type MergeResult struct {
	Added   []string
	Updated []string
	Unchanged []string
}

// Merger merges new secrets into an existing env map.
type Merger struct {
	overwrite bool
}

// NewMerger creates a Merger. If overwrite is true, existing keys are updated.
func NewMerger(overwrite bool) *Merger {
	return &Merger{overwrite: overwrite}
}

// Merge combines existing and incoming maps, returning merged result and a report.
func (m *Merger) Merge(existing, incoming map[string]string) (map[string]string, MergeResult) {
	result := make(map[string]string, len(existing))
	var report MergeResult

	for k, v := range existing {
		result[k] = v
	}

	for k, v := range incoming {
		old, exists := result[k]
		switch {
		case !exists:
			result[k] = v
			report.Added = append(report.Added, k)
		case m.overwrite && old != v:
			result[k] = v
			report.Updated = append(report.Updated, k)
		default:
			report.Unchanged = append(report.Unchanged, k)
		}
	}

	return result, report
}

// Summary returns a human-readable summary of the merge result.
func (r MergeResult) Summary() string {
	return fmt.Sprintf("added=%d updated=%d unchanged=%d",
		len(r.Added), len(r.Updated), len(r.Unchanged))
}
