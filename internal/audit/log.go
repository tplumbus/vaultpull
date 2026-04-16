// Package audit provides a simple audit logger for vaultpull operations.
// It records which keys were added, updated, or skipped during a sync.
package audit

import (
	"fmt"
	"io"
	"os"
	"time"
)

// Entry represents a single audit log entry.
type Entry struct {
	Timestamp time.Time
	Action    string
	Key       string
	Message   string
}

// Logger writes audit entries to a given writer.
type Logger struct {
	out io.Writer
}

// NewLogger creates a new Logger writing to w.
// If w is nil, os.Stdout is used.
func NewLogger(w io.Writer) *Logger {
	if w == nil {
		w = os.Stdout
	}
	return &Logger{out: w}
}

// Log writes an audit entry for the given action and key.
func (l *Logger) Log(action, key, message string) {
	e := Entry{
		Timestamp: time.Now().UTC(),
		Action:    action,
		Key:       key,
		Message:   message,
	}
	fmt.Fprintf(l.out, "%s [%s] key=%q %s\n",
		e.Timestamp.Format(time.RFC3339), e.Action, e.Key, e.Message)
}

// LogAdded logs that a key was added.
func (l *Logger) LogAdded(key string) {
	l.Log("ADDED", key, "key added from Vault")
}

// LogUpdated logs that a key was updated.
func (l *Logger) LogUpdated(key string) {
	l.Log("UPDATED", key, "key updated from Vault")
}

// LogSkipped logs that a key was skipped (already exists, no overwrite).
func (l *Logger) LogSkipped(key string) {
	l.Log("SKIPPED", key, "key already exists, overwrite=false")
}
