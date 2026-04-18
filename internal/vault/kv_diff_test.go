package vault

import (
	"testing"
)

func TestDiff_Added(t *testing.T) {
	remote := map[string]string{"FOO": "bar", "NEW": "value"}
	local := map[string]string{"FOO": "bar"}

	d := Diff(remote, local)

	if len(d.Added) != 1 || d.Added["NEW"] != "value" {
		t.Errorf("expected NEW to be added, got %v", d.Added)
	}
	if len(d.Changed) != 0 {
		t.Errorf("expected no changed keys, got %v", d.Changed)
	}
	if len(d.Removed) != 0 {
		t.Errorf("expected no removed keys, got %v", d.Removed)
	}
}

func TestDiff_Changed(t *testing.T) {
	remote := map[string]string{"FOO": "newval"}
	local := map[string]string{"FOO": "oldval"}

	d := Diff(remote, local)

	if len(d.Changed) != 1 || d.Changed["FOO"] != "newval" {
		t.Errorf("expected FOO to be changed, got %v", d.Changed)
	}
}

func TestDiff_Removed(t *testing.T) {
	remote := map[string]string{}
	local := map[string]string{"OLD": "gone"}

	d := Diff(remote, local)

	if len(d.Removed) != 1 || d.Removed[0] != "OLD" {
		t.Errorf("expected OLD to be removed, got %v", d.Removed)
	}
}

func TestDiff_NoChanges(t *testing.T) {
	m := map[string]string{"KEY": "val"}
	d := Diff(m, m)

	if d.HasChanges() {
		t.Error("expected no changes")
	}
}

func TestDiff_Summary(t *testing.T) {
	remote := map[string]string{"A": "1", "B": "new"}
	local := map[string]string{"B": "old", "C": "gone"}

	d := Diff(remote, local)
	summary := d.Summary()

	expected := "added=1 changed=1 removed=1"
	if summary != expected {
		t.Errorf("expected %q, got %q", expected, summary)
	}
}
