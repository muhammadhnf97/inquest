package internal

import (
	"os"
	"path/filepath"
	"testing"
)

func TestExecute(t *testing.T) {
	dir := t.TempDir()
	target := filepath.Join(dir, "node_modules")
	os.MkdirAll(target, 0755)
	os.WriteFile(filepath.Join(target, "pkg.js"), []byte("x"), 0644)

	if err := Execute(target); err != nil {
		t.Fatalf("Execute returned error: %v", err)
	}
	if _, err := os.Stat(target); !os.IsNotExist(err) {
		t.Errorf("expected %s to be deleted", target)
	}
}

func TestExecuteMissingPath(t *testing.T) {
	// os.RemoveAll is idempotent — deleting what doesn't exist is not an error
	missing := filepath.Join(t.TempDir(), "node_modules")
	if err := Execute(missing); err != nil {
		t.Errorf("Execute on missing path returned error: %v", err)
	}
}
