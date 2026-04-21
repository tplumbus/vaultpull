package audit

import (
	"fmt"
	"io"
	"os"
	"time"
)

// SearchLogger logs the results of a vault secret search operation.
type SearchLogger struct {
	out io.Writer
}

// NewSearchLogger creates a SearchLogger writing to w.
// If w is nil, os.Stdout is used.
func NewSearchLogger(w io.Writer) *SearchLogger {
	if w == nil {
		w = os.Stdout
	}
	return &SearchLogger{out: w}
}

// LogResult logs a single search result entry.
func (l *SearchLogger) LogResult(path string, keys []string) {
	for _, k := range keys {
		fmt.Fprintf(l.out, "[%s] SEARCH match path=%s key=%s\n",
			time.Now().UTC().Format(time.RFC3339), path, k)
	}
}

// LogSummary logs the overall summary of a search operation.
func (l *SearchLogger) LogSummary(query string, total int) {
	fmt.Fprintf(l.out, "[%s] SEARCH query=%q total_matches=%d\n",
		time.Now().UTC().Format(time.RFC3339), query, total)
}
