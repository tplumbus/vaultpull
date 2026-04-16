// Package env provides utilities for writing secrets to .env files.
package env

import (
	"fmt"
	"os"\n	"strings"
)

// Writer handles writing key-value secrets to a .env file.
type Writer struct {
	filePath string
}

// NewWriter creates a new Writer targeting the given file path.
func NewWriter(filePath string) *Writer {
	return &Writer{filePath: filePath}
}

// Write writes the provided secrets map to the .env file.
// Existing file contents are replaced. Keys are uppercased.
func (w *Writer) Write(secrets map[string]string) error {
	f, err := os.OpenFile(w.filePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0600)
	if err != nil {
		return fmt.Errorf("env: open file %q: %w", w.filePath, err)
	}
	defer f.Close()

	for k, v := range secrets {
		key := strings.ToUpper(k)
		line := fmt.Sprintf("%s=%s\n", key, escapeValue(v))
		if _, err := f.WriteString(line); err != nil {
			return fmt.Errorf("env: write key %q: %w", key, err)
		}
	}
	return nil
}

// escapeValue wraps the value in double quotes if it contains spaces or special chars.
func escapeValue(v string) string {
	if strings.ContainsAny(v, " \t\n#$") {
		v = strings.ReplaceAll(v, `"`, `\"`)
		return `"` + v + `"`
	}
	return v
}
