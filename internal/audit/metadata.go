package audit

import (
	"fmt"
	"io"
	"os"
	"time"
)

// MetadataEntry records metadata fetched for a secret path.
type MetadataEntry struct {
	Path           string
	CurrentVersion int
	UpdatedTime    time.Time
}

// MetadataLogger logs secret metadata information.
type MetadataLogger struct {
	out io.Writer
}

// NewMetadataLogger creates a MetadataLogger writing to out.
// If out is nil, os.Stdout is used.
func NewMetadataLogger(out io.Writer) *MetadataLogger {
	if out == nil {
		out = os.Stdout
	}
	return &MetadataLogger{out: out}
}

// Log writes a metadata entry to the output.
func (l *MetadataLogger) Log(e MetadataEntry) {
	fmt.Fprintf(l.out, "[metadata] path=%s version=%d updated=%s\n",
		e.Path, e.CurrentVersion, e.UpdatedTime.Format(time.RFC3339))
}

// LogBatch writes multiple metadata entries.
func (l *MetadataLogger) LogBatch(entries []MetadataEntry) {
	for _, e := range entries {
		l.Log(e)
	}
}
