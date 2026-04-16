package config

import (
	"errors"
	"os"

	"github.com/joho/godotenv"
)

// Config holds the configuration for vaultpull.
type Config struct {
	VaultAddr  string
	VaultToken string
	VaultPath  string
	OutputFile string
}

// Load reads configuration from environment variables and an optional .env file.
func Load(envFile string) (*Config, error) {
	if envFile != "" {
		_ = godotenv.Load(envFile)
	}

	cfg := &Config{
		VaultAddr:  getEnv("VAULT_ADDR", "http://127.0.0.1:8200"),
		VaultToken: os.Getenv("VAULT_TOKEN"),
		VaultPath:  os.Getenv("VAULT_PATH"),
		OutputFile: getEnv("OUTPUT_FILE", ".env"),
	}

	if err := cfg.validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

func (c *Config) validate() error {
	if c.VaultToken == "" {
		return errors.New("VAULT_TOKEN is required but not set")
	}
	if c.VaultPath == "" {
		return errors.New("VAULT_PATH is required but not set")
	}
	return nil
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
