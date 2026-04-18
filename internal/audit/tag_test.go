package audit

import (
	"bytes"
	"strings"
	"testing"
)

func TestTagLogger_Log(t *testing.T) {
	var buf bytes.Buffer
	l := NewTagLogger(&buf)
	l.Log("secret/myapp", map[string]string{"env": "prod"})

	if len(l.Entries()) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(l.Entries()))
	}
	if l.Entries()[0].Path != "secret/myapp" {
		t.Errorf("unexpected path: %s", l.Entries()[0].Path)
	}
	if !strings.Contains(buf.String(), "secret/myapp") {
		t.Errorf("expected path in output, got: %s", buf.String())
	}
	if !strings.Contains(buf.String(), "count=1") {
		t.Errorf("expected count=1 in output, got: %s", buf.String())
	}
}

func TestTagLogger_NilUsesStdout(t *testing.T) {
	l := NewTagLogger(nil)
	if l.out == nil {
		t.Error("expected non-nil writer")
	}
}

func TestTagLogger_MultipleEntries(t *testing.T) {
	var buf bytes.Buffer
	l := NewTagLogger(&buf)
	l.Log("secret/a", map[string]string{"x": "1"})
	l.Log("secret/b", map[string]string{"y": "2", "z": "3"})

	if len(l.Entries()) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(l.Entries()))
	}
	if l.Entries()[1].Tags["z"] != "3" {
		t.Errorf("unexpected tag value")
	}
}
