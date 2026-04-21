package vault

import (
	"context"
	"fmt"
)

// ImportResult holds the outcome of a bulk import operation.
type ImportResult struct {
	Written []string
	Skipped []string
	Failed  []string
}

// Summary returns a human-readable summary of the import result.
func (r ImportResult) Summary() string {
	return fmt.Sprintf("import complete: %d written, %d skipped, %d failed",
		len(r.Written), len(r.Skipped), len(r.Failed))
}

// ImportSecrets writes multiple key/value pairs into Vault under the given
// path prefix. If overwrite is false, existing keys are skipped.
func ImportSecrets(
	ctx context.Context,
	client *Client,
	pathPrefix string,
	secrets map[string]map[string]string,
	overwrite bool,
) (ImportResult, error) {
	var result ImportResult

	for key, data := range secrets {
		fullPath := pathPrefix + "/" + key

		if !overwrite {
			existing, err := client.GetSecret(ctx, fullPath)
			if err == nil && len(existing) > 0 {
				result.Skipped = append(result.Skipped, fullPath)
				continue
			}
		}

		if err := client.PutSecret(ctx, fullPath, data); err != nil {
			result.Failed = append(result.Failed, fullPath)
			continue
		}

		result.Written = append(result.Written, fullPath)
	}

	return result, nil
}
