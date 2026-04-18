package audit

import (
	"fmt"
	"io"
	"os"
	"time"
)

// TagEntry records a tag retrieval event for a secret path.
type TagEntry struct {
	Timestamp time.Time
	Path      string
	Tags      map[string]string
}

// TagLogger logs tag fetch events.
type TagLogger struct {
	out     io.Writer
	entries []TagEntry
}

// NewTagLogger creates a TagLogger writing to w. If w is nil, os.Stdout is used.
func NewTagLogger(w io.Writer) *TagLogger {
	if w == nil {
		w = os.Stdout
	}
	return &TagLogger{out: w}
}

// Log records and prints a tag entry.
func (l *TagLogger) Log(path string, tags map[string]string) {
	e := TagEntry{Timestamp: time.Now(), Path: path, Tags: tags}
	l.entries = append(l.entries, e)
	fmt.Fprintf(l.out, "[tags] %s path=%s count=%d\n", e.Timestamp.Format(time.RFC3339), path, len(tags))
}

// Entries returns all recorded tag entries.
func (l *TagLogger) Entries() []TagEntry {
	return l.entries
}

// EntriesForPath returns all recorded tag entries matching the given path.
func (l *TagLogger) EntriesForPath(path string) []TagEntry {
	var result []TagEntry
	for _, e := range l.entries {
		if e.Path == path {
			result = append(result, e)
		}
	}
	return result
}
