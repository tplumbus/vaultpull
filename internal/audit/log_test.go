package audit

import (
	"bytes"
	"strings"
	"testing"
)

func TestLogger_LogAdded(t *testing.T) {
	var buf bytes.Buffer
	l := NewLogger(&buf)
	l.LogAdded("DB_PASSWORD")

	out := buf.String()
	if !strings.Contains(out, "ADDED") {
		t.Errorf("expected ADDED in output, got: %s", out)
	}
	if !strings.Contains(out, "DB_PASSWORD") {
		t.Errorf("expected key name in output, got: %s", out)
	}
}

func TestLogger_LogUpdated(t *testing.T) {
	var buf bytes.Buffer
	l := NewLogger(&buf)
	l.LogUpdated("API_KEY")

	out := buf.String()
	if !strings.Contains(out, "UPDATED") {
		t.Errorf("expected UPDATED in output, got: %s", out)
	}
	if !strings.Contains(out, "API_KEY") {
		t.Errorf("expected key name in output, got: %s", out)
	}
}

func TestLogger_LogSkipped(t *testing.T) {
	var buf bytes.Buffer
	l := NewLogger(&buf)
	l.LogSkipped("SECRET_TOKEN")

	out := buf.String()
	if !strings.Contains(out, "SKIPPED") {
		t.Errorf("expected SKIPPED in output, got: %s", out)
	}
	if !strings.Contains(out, "overwrite=false") {
		t.Errorf("expected overwrite=false message in output, got: %s", out)
	}
}

func TestNewLogger_NilUsesStdout(t *testing.T) {
	// Should not panic when nil is passed
	l := NewLogger(nil)
	if l == nil {
		t.Fatal("expected non-nil logger")
	}
	if l.out == nil {
		t.Fatal("expected non-nil writer")
	}
}

func TestLogger_MultipleEntries(t *testing.T) {
	var buf bytes.Buffer
	l := NewLogger(&buf)
	l.LogAdded("KEY_ONE")
	l.LogSkipped("KEY_TWO")
	l.LogUpdated("KEY_THREE")

	out := buf.String()
	lines := strings.Split(strings.TrimSpace(out), "\n")
	if len(lines) != 3 {
		t.Errorf("expected 3 log lines, got %d: %s", len(lines), out)
	}
}
