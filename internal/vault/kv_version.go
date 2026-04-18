package vault

import (
	"fmt"
	"strconv"
)

// KVVersion represents the KV secrets engine version.
type KVVersion int

const (
	KVVersion1 KVVersion = 1
	KVVersion2 KVVersion = 2
)

// String returns the string representation of the KV version.
func (v KVVersion) String() string {
	return strconv.Itoa(int(v))
}

// KVVersionFromString parses a string into a KVVersion.
// An empty string defaults to KVVersion2.
func KVVersionFromString(s string) (KVVersion, error) {
	if s == "" {
		return KVVersion2, nil
	}
	switch s {
	case "1":
		return KVVersion1, nil
	case "2":
		return KVVersion2, nil
	default:
		return 0, fmt.Errorf("vault: unsupported KV version %q: must be \"1\" or \"2\"", s)
	}
}
