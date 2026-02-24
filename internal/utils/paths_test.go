package utils

import (
	"os"
	"path/filepath"
	"testing"
)

func TestXDGConfigHome(t *testing.T) {
	t.Run("uses XDG_CONFIG_HOME when set", func(t *testing.T) {
		expected := "/custom/config"
		t.Setenv("XDG_CONFIG_HOME", expected)
		
		got := XDGConfigHome()
		if got != expected {
			t.Errorf("XDGConfigHome() = %q, want %q", got, expected)
		}
	})

	t.Run("falls back to .config when not set", func(t *testing.T) {
		t.Setenv("XDG_CONFIG_HOME", "")
		home, err := os.UserHomeDir()
		if err != nil {
			t.Skip("cannot get home directory")
		}
		
		expected := filepath.Join(home, ".config")
		got := XDGConfigHome()
		if got != expected {
			t.Errorf("XDGConfigHome() = %q, want %q", got, expected)
		}
	})
}

func TestXDGDataHome(t *testing.T) {
	t.Run("uses XDG_DATA_HOME when set", func(t *testing.T) {
		expected := "/custom/data"
		t.Setenv("XDG_DATA_HOME", expected)
		
		got := XDGDataHome()
		if got != expected {
			t.Errorf("XDGDataHome() = %q, want %q", got, expected)
		}
	})

	t.Run("falls back to .local/share when not set", func(t *testing.T) {
		t.Setenv("XDG_DATA_HOME", "")
		home, err := os.UserHomeDir()
		if err != nil {
			t.Skip("cannot get home directory")
		}
		
		expected := filepath.Join(home, ".local", "share")
		got := XDGDataHome()
		if got != expected {
			t.Errorf("XDGDataHome() = %q, want %q", got, expected)
		}
	})
}

func TestXDGCacheHome(t *testing.T) {
	t.Run("uses XDG_CACHE_HOME when set", func(t *testing.T) {
		expected := "/custom/cache"
		t.Setenv("XDG_CACHE_HOME", expected)
		
		got := XDGCacheHome()
		if got != expected {
			t.Errorf("XDGCacheHome() = %q, want %q", got, expected)
		}
	})

	t.Run("falls back to .cache when not set", func(t *testing.T) {
		t.Setenv("XDG_CACHE_HOME", "")
		home, err := os.UserHomeDir()
		if err != nil {
			t.Skip("cannot get home directory")
		}
		
		expected := filepath.Join(home, ".cache")
		got := XDGCacheHome()
		if got != expected {
			t.Errorf("XDGCacheHome() = %q, want %q", got, expected)
		}
	})
}

func TestConfigPathGlobal(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", "/test/config")
	
	expected := filepath.Join("/test/config", "go-cli-template", "config.toml")
	got := ConfigPathGlobal()
	if got != expected {
		t.Errorf("ConfigPathGlobal() = %q, want %q", got, expected)
	}
}

func TestConfigPathLocal(t *testing.T) {
	cwd := "/project/dir"
	expected := filepath.Join(cwd, ".go-cli-template", "config.toml")
	
	got := ConfigPathLocal(cwd)
	if got != expected {
		t.Errorf("ConfigPathLocal(%q) = %q, want %q", cwd, got, expected)
	}
}
