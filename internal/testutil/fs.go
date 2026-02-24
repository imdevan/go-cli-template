package testutil

import (
	"os"
	"path/filepath"
	"testing"
)

// WithTempWorkspace sets CWD to a temporary directory.
func WithTempWorkspace(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	if err := os.Chdir(root); err != nil {
		t.Fatalf("chdir: %v", err)
	}
	return root
}

// FixturesPath resolves a path under tests/fixtures.
func FixturesPath(parts ...string) string {
	segments := append([]string{"tests", "fixtures"}, parts...)
	return filepath.Join(segments...)
}

// CopyFixture copies a fixture file into the destination path.
func CopyFixture(t *testing.T, fixturePath, destPath string) {
	t.Helper()
	data, err := os.ReadFile(fixturePath)
	if err != nil {
		t.Fatalf("read fixture %s: %v", fixturePath, err)
	}
	if err := os.MkdirAll(filepath.Dir(destPath), 0o755); err != nil {
		t.Fatalf("mkdir %s: %v", destPath, err)
	}
	if err := os.WriteFile(destPath, data, 0o644); err != nil {
		t.Fatalf("write fixture %s: %v", destPath, err)
	}
}
