// Package env provides utilities for reading and writing .env files.
package env

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Reader reads key-value pairs from an existing .env file.
type Reader struct {
	path string
}

// NewReader creates a new Reader for the given file path.
func NewReader(path string) *Reader {
	return &Reader{path: path}
}

// Read parses the .env file and returns a map of key-value pairs.
// Lines beginning with '#' and empty lines are ignored.
// Values may optionally be wrapped in double quotes, which are stripped.
// Returns an empty map if the file does not exist.
func (r *Reader) Read() (map[string]string, error) {
	result := make(map[string]string)

	f, err := os.Open(r.path)
	if os.IsNotExist(err) {
		return result, nil
	}
	if err != nil {
		return nil, fmt.Errorf("env reader: open %q: %w", r.path, err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	lineNum := 0
	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("env reader: %q line %d: invalid format", r.path, lineNum)
		}
		key := strings.TrimSpace(parts[0])
		val := strings.Trim(strings.TrimSpace(parts[1]), `"`)
		result[key] = val
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("env reader: scan %q: %w", r.path, err)
	}
	return result, nil
}

// Keys returns the list of keys defined in the .env file, preserving
// the order in which they appear. Comments and blank lines are skipped.
func (r *Reader) Keys() ([]string, error) {
	f, err := os.Open(r.path)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("env reader: open %q: %w", r.path, err)
	}
	defer f.Close()

	var keys []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		if key, _, found := strings.Cut(line, "="); found {
			keys = append(keys, strings.TrimSpace(key))
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("env reader: scan %q: %w", r.path, err)
	}
	return keys, nil
}
