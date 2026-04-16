// Package sync provides the core orchestration logic for vaultpull.
//
// It ties together the Vault client and the env file writer,
// coordinating the fetch-and-write lifecycle for a single sync run.
//
// Typical usage:
//
//	cfg, err := config.Load()
//	if err != nil { ... }
//
//	s, err := sync.New(cfg)
//	if err != nil { ... }
//
//	count, err := s.Run()
//	if err != nil { ... }
//	fmt.Printf("synced %d secrets\n", count)
package sync
