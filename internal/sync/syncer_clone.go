package sync

import (
	"fmt"
	"io"
	"os"

	"github.com/yourusername/vaultpull/internal/vault"
)

// RunClone reads secrets from srcPath in Vault and writes them to each path in
// dstPaths. When overwrite is false the operation aborts if any destination
// already exists. Progress is written to out (defaults to os.Stdout).
func RunClone(client *vault.Client, srcPath string, dstPaths []string, overwrite bool, out io.Writer) error {
	if out == nil {
		out = os.Stdout
	}

	if srcPath == "" {
		return fmt.Errorf("clone: source path must not be empty")
	}
	if len(dstPaths) == 0 {
		return fmt.Errorf("clone: at least one destination path is required")
	}

	fmt.Fprintf(out, "cloning %q to %d destination(s)...\n", srcPath, len(dstPaths))

	result := vault.CloneResult{Source: srcPath}

	for _, dst := range dstPaths {
		err := client.CloneSecret(srcPath, []string{dst}, overwrite)
		if err != nil {
			if !overwrite {
				result.Skipped = append(result.Skipped, dst)
				fmt.Fprintf(out, "  skipped %q: %v\n", dst, err)
				continue
			}
			return fmt.Errorf("clone: %w", err)
		}
		result.Destinations = append(result.Destinations, dst)
		fmt.Fprintf(out, "  cloned → %q\n", dst)
	}

	fmt.Fprintln(out, result.Summary())
	return nil
}
