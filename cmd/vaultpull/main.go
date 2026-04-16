package main

import (
	"fmt"
	"os"

	"github.com/vaultpull/internal/config"
	"github.com/vaultpull/internal/env"
	"github.com/vaultpull/internal/sync"
	"github.com/vaultpull/internal/vault"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}

	client, err := vault.NewClient(cfg.VaultAddr, cfg.VaultToken)
	if err != nil {
		return fmt.Errorf("creating vault client: %w", err)
	}

	writer, err := env.NewWriter(cfg.OutputFile)
	if err != nil {
		return fmt.Errorf("creating env writer: %w", err)
	}

	syncer := sync.New(client, writer)
	if err := syncer.Run(cfg.SecretPath); err != nil {
		return fmt.Errorf("syncing secrets: %w", err)
	}

	fmt.Printf("secrets synced to %s\n", cfg.OutputFile)
	return nil
}
