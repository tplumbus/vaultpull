package env

import (
	"os"
	"path/filepath"
	"testing"
)

func writeFile(t *testing.T, dir, name, content string) string {
	t.Helper()
	p := filepath.Join(dir, name)
	if err := os.WriteFile(p, []byte(content), 0600); err != nil {
		t.Fatalf("writeFile: %v", err)
	}
	return p
}

func TestReader_Read_ParsesKeyValues(t *testing.T) {
	dir := t.TempDir()
	p := writeFile(t, dir, ".env", "FOO=bar\nBAZ=\"hello world\"\n")

	r := NewReader(p)
	got, err := r.Read()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["FOO"] != "bar" {
		t.Errorf("FOO: got %q, want %q", got["FOO"], "bar")
	}
	if got["BAZ"] != "hello world" {
		t.Errorf("BAZ: got %q, want %q", got["BAZ"], "hello world")
	}
}

func TestReader_Read_IgnoresComments(t *testing.T) {
	dir := t.TempDir()
	p := writeFile(t, dir, ".env", "# comment\nKEY=value\n")

	got, err := NewReader(p).Read()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := got["# comment"]; ok {
		t.Error("comment line should not be parsed as key")
	}
	if got["KEY"] != "value" {
		t.Errorf("KEY: got %q, want %q", got["KEY"], "value")
	}
}

func TestReader_Read_MissingFile(t *testing.T) {
	r := NewReader("/nonexistent/.env")
	got, err := r.Read()
	if err != nil {
		t.Fatalf("expected no error for missing file, got: %v", err)
	}
	if len(got) != 0 {
		t.Errorf("expected empty map, got %v", got)
	}
}

func TestReader_Read_InvalidFormat(t *testing.T) {
	dir := t.TempDir()
	p := writeFile(t, dir, ".env", "BADLINE\n")

	_, err := NewReader(p).Read()
	if err == nil {
		t.Error("expected error for invalid format, got nil")
	}
}
