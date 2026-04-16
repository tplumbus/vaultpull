package env

import (
	"testing"
)

func TestMerger_Merge_AddsNewKeys(t *testing.T) {
	m := NewMerger(false)
	existing := map[string]string{"FOO": "bar"}
	incoming := map[string]string{"BAZ": "qux"}

	out, report := m.Merge(existing, incoming)

	if out["FOO"] != "bar" {
		t.Errorf("expected FOO=bar, got %s", out["FOO"])
	}
	if out["BAZ"] != "qux" {
		t.Errorf("expected BAZ=qux, got %s", out["BAZ"])
	}
	if len(report.Added) != 1 || report.Added[0] != "BAZ" {
		t.Errorf("expected Added=[BAZ], got %v", report.Added)
	}
}

func TestMerger_Merge_NoOverwrite(t *testing.T) {
	m := NewMerger(false)
	existing := map[string]string{"FOO": "original"}
	incoming := map[string]string{"FOO": "new"}

	out, report := m.Merge(existing, incoming)

	if out["FOO"] != "original" {
		t.Errorf("expected FOO=original, got %s", out["FOO"])
	}
	if len(report.Unchanged) != 1 {
		t.Errorf("expected 1 unchanged, got %d", len(report.Unchanged))
	}
	if len(report.Updated) != 0 {
		t.Errorf("expected 0 updated, got %d", len(report.Updated))
	}
}

func TestMerger_Merge_WithOverwrite(t *testing.T) {
	m := NewMerger(true)
	existing := map[string]string{"FOO": "original"}
	incoming := map[string]string{"FOO": "new"}

	out, report := m.Merge(existing, incoming)

	if out["FOO"] != "new" {
		t.Errorf("expected FOO=new, got %s", out["FOO"])
	}
	if len(report.Updated) != 1 {
		t.Errorf("expected 1 updated, got %d", len(report.Updated))
	}
}

func TestMergeResult_Summary(t *testing.T) {
	r := MergeResult{
		Added:     []string{"A", "B"},
		Updated:   []string{"C"},
		Unchanged: []string{},
	}
	s := r.Summary()
	expected := "added=2 updated=1 unchanged=0"
	if s != expected {
		t.Errorf("expected %q, got %q", expected, s)
	}
}
