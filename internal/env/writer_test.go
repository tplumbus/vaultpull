package env

import (
	"bufio"
	"os"
	"strings"
	"testing"
)

func TestWriter_Write_CreatesFile(t *testing.T) {
	tmp, err := os.CreateTemp(t.TempDir(), ".env")
	if err != nil {
		t.Fatal(err)
	}
	tmp.Close()

	w := NewWriter(tmp.Name())
	secrets := map[string]string{
		"db_password": "s3cr3t",
		"api_key":     "abc123",
	}
	if err := w.Write(secrets); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	f, err := os.Open(tmp.Name())
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	lines := map[string]string{}
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		parts := strings.SplitN(scanner.Text(), "=", 2)
		if len(parts) == 2 {
			lines[parts[0]] = parts[1]
		}
	}

	if lines["DB_PASSWORD"] != "s3cr3t" {
		t.Errorf("expected DB_PASSWORD=s3cr3t, got %q", lines["DB_PASSWORD"])
	}
	if lines["API_KEY"] != "abc123" {
		t.Errorf("expected API_KEY=abc123, got %q", lines["API_KEY"])
	}
}

func TestWriter_Write_EscapesSpaces(t *testing.T) {
	tmp, err := os.CreateTemp(t.TempDir(), ".env")
	if err != nil {
		t.Fatal(err)
	}
	tmp.Close()

	w := NewWriter(tmp.Name())
	if err := w.Write(map[string]string{"note": "hello world"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, _ := os.ReadFile(tmp.Name())
	if !strings.Contains(string(data), `"hello world"`) {
		t.Errorf("expected quoted value, got: %s", string(data))
	}
}

func TestWriter_Write_InvalidPath(t *testing.T) {
	w := NewWriter("/nonexistent/path/.env")
	err := w.Write(map[string]string{"key": "val"})
	if err == nil {
		t.Error("expected error for invalid path, got nil")
	}
}
