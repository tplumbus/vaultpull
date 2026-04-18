package sync

import (
	"bytes"
	"errors"
	"strings"
	"testing"
)

type mockRenamer struct {
	err error
}

func (m *mockRenamer) RenameSecret(mount, src, dst string, overwrite bool) error {
	return m.err
}

func TestRunRename_Success(t *testing.T) {
	var buf bytes.Buffer
	err := RunRename(&mockRenamer{}, RenameOptions{
		Mount:   "secret",
		SrcPath: "old/key",
		DstPath: "new/key",
		Out:     &buf,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "old/key -> new/key") {
		t.Errorf("expected rename log, got: %s", buf.String())
	}
}

func TestRunRename_MissingSrc(t *testing.T) {
	err := RunRename(&mockRenamer{}, RenameOptions{Mount: "secret", DstPath: "dst"})
	if err == nil || !strings.Contains(err.Error(), "source path") {
		t.Fatalf("expected source path error, got: %v", err)
	}
}

func TestRunRename_MissingDst(t *testing.T) {
	err := RunRename(&mockRenamer{}, RenameOptions{Mount: "secret", SrcPath: "src"})
	if err == nil || !strings.Contains(err.Error(), "destination path") {
		t.Fatalf("expected destination path error, got: %v", err)
	}
}

func TestRunRename_VaultError(t *testing.T) {
	err := RunRename(&mockRenamer{err: errors.New("vault down")}, RenameOptions{
		Mount: "secret", SrcPath: "a", DstPath: "b",
	})
	if err == nil || !strings.Contains(err.Error(), "vault down") {
		t.Fatalf("expected vault error, got: %v", err)
	}
}
