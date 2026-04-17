package vault

// kvVersion stores the configured KV engine version for the client.
// It is set during client construction and used in path resolution.

// WithKVVersion sets the KV secrets engine version on the client.
// Defaults to KVv2 if not called.
func (c *Client) WithKVVersion(v KVVersion) *Client {
	c.kvVersion = v
	return c
}

// KVVersionFromString parses a version string ("1" or "2") into a KVVersion.
// Returns KVv2 and false if the string is unrecognised.
func KVVersionFromString(s string) (KVVersion, bool) {
	switch s {
	case "1":
		return KVv1, true
	case "2":
		return KVv2, true
	default:
		return KVv2, false
	}
}
