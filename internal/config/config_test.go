package config

import (
	"os"
	"testing"
)

func setEnv(t *testing.T, key, value string) {
	t.Helper()
	t.Setenv(key, value)
}

func TestLoad_MissingToken(t *testing.T) {
	os.Unsetenv("VAULT_TOKEN")
	os.Unsetenv("VAULT_PATH")

	_, err := Load("")
	if err == nil {
		t.Fatal("expected error when VAULT_TOKEN is missing, got nil")
	}
}

func TestLoad_MissingPath(t *testing.T) {
	setEnv(t, "VAULT_TOKEN", "test-token")
	os.Unsetenv("VAULT_PATH")

	_, err := Load("")
	if err == nil {
		t.Fatal("expected error when VAULT_PATH is missing, got nil")
	}
}

func TestLoad_Defaults(t *testing.T) {
	setEnv(t, "VAULT_TOKEN", "test-token")
	setEnv(t, "VAULT_PATH", "secret/myapp")
	os.Unsetenv("VAULT_ADDR")
	os.Unsetenv("OUTPUT_FILE")

	cfg, err := Load("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.VaultAddr != "http://127.0.0.1:8200" {
		t.Errorf("expected default VAULT_ADDR, got %s", cfg.VaultAddr)
	}
	if cfg.OutputFile != ".env" {
		t.Errorf("expected default OUTPUT_FILE, got %s", cfg.OutputFile)
	}
}

func TestLoad_CustomValues(t *testing.T) {
	setEnv(t, "VAULT_TOKEN", "s.abc123")
	setEnv(t, "VAULT_PATH", "secret/data/prod")
	setEnv(t, "VAULT_ADDR", "https://vault.example.com")
	setEnv(t, "OUTPUT_FILE", ".env.prod")

	cfg, err := Load("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.VaultAddr != "https://vault.example.com" {
		t.Errorf("unexpected VaultAddr: %s", cfg.VaultAddr)
	}
	if cfg.OutputFile != ".env.prod" {
		t.Errorf("unexpected OutputFile: %s", cfg.OutputFile)
	}
}
