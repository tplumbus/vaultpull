package audit

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func TestMetadataLogger_Log(t *testing.T) {
	var buf bytes.Buffer
	logger := NewMetadataLogger(&buf)

	now := time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)
	logger.Log(MetadataEntry{
		Path:           "secret/myapp",
		CurrentVersion: 5,
		UpdatedTime:    now,
	})

	out := buf.String()
	if !strings.Contains(out, "secret/myapp") {
		t.Errorf("expected path in output, got: %s", out)
	}
	if !strings.Contains(out, "version=5") {
		t.Errorf("expected version in output, got: %s", out)
	}
	if !strings.Contains(out, "2024-01-15") {
		t.Errorf("expected date in output, got: %s", out)
	}
}

func TestMetadataLogger_NilUsesStdout(t *testing.T) {
	logger := NewMetadataLogger(nil)
	if logger.out == nil {
		t.Error("expected non-nil writer")
	}
}

func TestMetadataLogger_LogBatch(t *testing.T) {
	var buf bytes.Buffer
	logger := NewMetadataLogger(&buf)

	entries := []MetadataEntry{
		{Path: "secret/a", CurrentVersion: 1, UpdatedTime: time.Now()},
		{Path: "secret/b", CurrentVersion: 2, UpdatedTime: time.Now()},
	}
	logger.LogBatch(entries)

	out := buf.String()
	if !strings.Contains(out, "secret/a") {
		t.Errorf("expected secret/a in output")
	}
	if !strings.Contains(out, "secret/b") {
		t.Errorf("expected secret/b in output")
	}
}
