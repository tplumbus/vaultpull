// Package main is the entry point for the vaultpull CLI tool.
//
// vaultpull syncs secrets from HashiCorp Vault into a local .env file.
// Configuration is provided via environment variables:
//
//	VAULT_ADDR        - Vault server address (default: http://127.0.0.1:8200)
//	VAULT_TOKEN       - Vault authentication token (required)
//	VAULT_SECRET_PATH - Path to the secret in Vault (required)
//	OUTPUT_FILE       - Output .env file path (default: .env)
//
// Usage:
//
//	export VAULT_TOKEN=s.mytoken
//	export VAULT_SECRET_PATH=secret/data/myapp
//	vaultpull
package main
