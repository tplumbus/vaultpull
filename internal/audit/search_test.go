package audit

import (
	"bytes"
	"strings"
	"testing"
)

func TestSearchLogger_LogResult(t *testing.T) {
	var buf bytes.Buffer
	l := NewSearchLogger(&buf)
	l.LogResult("secret/app", []string{"DB_HOST", "DB_PORT"})

	out := buf.String()
	if !strings.Contains(out, "SEARCH") {
		t.Error("expected SEARCH in output")
	}
	if !strings.Contains(out, "DB_HOST") {
		t.Error("expected DB_HOST in output")
	}
	if !strings.Contains(out, "DB_PORT") {
		t.Error("expected DB_PORT in output")
	}
	if !strings.Contains(out, "secret/app") {
		t.Error("expected path in output")
	}
}

func TestSearchLogger_LogSummary(t *testing.T) {
	var buf bytes.Buffer
	l := NewSearchLogger(&buf)
	l.LogSummary("DB", 5)

	out := buf.String()
	if !strings.Contains(out, "query=\"DB\"") {
		t.Errorf("expected query in output, got: %s", out)
	}
	if !strings.Contains(out, "total_matches=5") {
		t.Errorf("expected total_matches in output, got: %s", out)
	}
}

func TestSearchLogger_NilUsesStdout(t *testing.T) {
	l := NewSearchLogger(nil)
	if l.out == nil {
		t.Error("expected non-nil writer")
	}
}
