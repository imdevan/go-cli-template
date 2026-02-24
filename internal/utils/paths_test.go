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
	
	got := ConfigPathGlobal()
	
	// Verify it uses the XDG_CONFIG_HOME
	if !filepath.HasPrefix(got, "/test/config") {
		t.Errorf("ConfigPathGlobal() = %q, should start with /test/config", got)
	}
	
	// Verify it ends with config.toml
	if filepath.Base(got) != "config.toml" {
		t.Errorf("ConfigPathGlobal() = %q, should end with config.toml", got)
	}
}

func TestConfigPathLocal(t *testing.T) {
	cwd := "/project/dir"
	got := ConfigPathLocal(cwd)
	
	// Verify it starts with the cwd
	if !filepath.HasPrefix(got, cwd) {
		t.Errorf("ConfigPathLocal(%q) = %q, should start with %q", cwd, got, cwd)
	}
	
	// Verify it ends with config.toml
	if filepath.Base(got) != "config.toml" {
		t.Errorf("ConfigPathLocal(%q) = %q, should end with config.toml", cwd, got)
	}
}
