package sync

import (
	"fmt"

	"github.com/vaultpull/internal/config"
	"github.com/vaultpull/internal/env"
	"github.com/vaultpull/internal/vault"
)

// Syncer orchestrates fetching secrets from Vault and writing them to a .env file.
type Syncer struct {
	client *vault.Client
	writer *env.Writer
	cfg    *config.Config
}

// New creates a new Syncer from the provided config.
func New(cfg *config.Config) (*Syncer, error) {
	client, err := vault.NewClient(cfg)
	if err != nil {
		return nil, fmt.Errorf("syncer: failed to create vault client: %w", err)
	}

	writer, err := env.NewWriter(cfg.OutputFile)
	if err != nil {
		return nil, fmt.Errorf("syncer: failed to create env writer: %w", err)
	}

	return &Syncer{client: client, writer: writer, cfg: cfg}, nil
}

// Run fetches secrets from Vault and writes them to the configured output file.
func (s *Syncer) Run() (int, error) {
	secrets, err := s.client.GetSecret(s.cfg.VaultPath)
	if err != nil {
		return 0, fmt.Errorf("syncer: failed to get secrets: %w", err)
	}

	if len(secrets) == 0 {
		return 0, nil
	}

	if err := s.writer.Write(secrets); err != nil {
		return 0, fmt.Errorf("syncer: failed to write env file: %w", err)
	}

	return len(secrets), nil
}
