package utils

import (
	"os"
	"path/filepath"
)

// XDGConfigHome returns the config home directory with sensible defaults.
func XDGConfigHome() string {
	return xdgHome("XDG_CONFIG_HOME", ".config")
}

// XDGDataHome returns the data home directory with sensible defaults.
func XDGDataHome() string {
	return xdgHome("XDG_DATA_HOME", filepath.Join(".local", "share"))
}

// XDGCacheHome returns the cache home directory with sensible defaults.
func XDGCacheHome() string {
	return xdgHome("XDG_CACHE_HOME", ".cache")
}

// ConfigPathGlobal returns the default global config file path.
func ConfigPathGlobal() string {
	return filepath.Join(XDGConfigHome(), "go-cli-template", "config.toml")
}

// ConfigPathLocal returns the local config file path for the given cwd.
func ConfigPathLocal(cwd string) string {
	return filepath.Join(cwd, ".go-cli-template", "config.toml")
}

func xdgHome(envKey, fallbackSuffix string) string {
	if value := os.Getenv(envKey); value != "" {
		return value
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(home, fallbackSuffix)
}
