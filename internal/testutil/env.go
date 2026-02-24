package testutil

import (
	"path/filepath"
	"testing"
)

// WithTempXDG sets XDG paths to temporary directories for tests.
func WithTempXDG(t *testing.T) (configDir, dataDir, cacheDir string) {
	t.Helper()
	root := t.TempDir()
	configDir = filepath.Join(root, "config")
	dataDir = filepath.Join(root, "data")
	cacheDir = filepath.Join(root, "cache")
	t.Setenv("XDG_CONFIG_HOME", configDir)
	t.Setenv("XDG_DATA_HOME", dataDir)
	t.Setenv("XDG_CACHE_HOME", cacheDir)
	t.Setenv("NO_COLOR", "1")
	return configDir, dataDir, cacheDir
}
