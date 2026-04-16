// Package vault provides a lightweight client for reading secrets
// from a HashiCorp Vault instance using the KV v2 secrets engine.
//
// Usage:
//
//	client := vault.NewClient("https://vault.example.com", "s.mytoken")
//	data, err := client.GetSecret("secret/data/myapp/prod")
//	if err != nil {
//		log.Fatal(err)
//	}
//	// data is a map[string]string of key/value secret pairs
//
// The client expects secrets stored under the standard KV v2 path structure:
//
//	/v1/<mount>/data/<path>
//
// Authentication is performed via the X-Vault-Token header.
package vault
