// Package env provides utilities for reading, writing, and merging
// .env files used by the vaultpull CLI tool.
//
// It supports:
//   - Reading existing .env files into key-value maps
//   - Writing key-value maps to .env files with proper escaping
//   - Merging Vault secrets into existing .env files with optional overwrite
package env
