package config

import (
	"os"
	"path/filepath"
	"testing"

	"go-cli-template/internal/domain"
	"go-cli-template/internal/utils"
)

func TestManagerLoadsDefaults(t *testing.T) {
	root := t.TempDir()
	cwd := filepath.Join(root, "project")
	if err := os.MkdirAll(cwd, 0o755); err != nil {
		t.Fatalf("mkdir cwd: %v", err)
	}
	t.Setenv("XDG_CONFIG_HOME", filepath.Join(root, "config"))
	t.Setenv("XDG_DATA_HOME", filepath.Join(root, "data"))

	manager := NewManager(cwd)
	cfg, err := manager.Load()
	if err != nil {
		t.Fatalf("load: %v", err)
	}

	if cfg.Editor == "" {
		t.Fatal("expected editor to have default value")
	}
}

func TestManagerLoadsFromFile(t *testing.T) {
	root := t.TempDir()
	cwd := filepath.Join(root, "project")
	if err := os.MkdirAll(cwd, 0o755); err != nil {
		t.Fatalf("mkdir cwd: %v", err)
	}
	t.Setenv("XDG_CONFIG_HOME", filepath.Join(root, "config"))

	configPath := utils.ConfigPathGlobal()
	if err := os.MkdirAll(filepath.Dir(configPath), 0o755); err != nil {
		t.Fatalf("mkdir config dir: %v", err)
	}

	data := []byte("editor = \"vim\"\n")
	if err := os.WriteFile(configPath, data, 0o644); err != nil {
		t.Fatalf("write config: %v", err)
	}

	manager := NewManager(cwd)
	cfg, err := manager.Load()
	if err != nil {
		t.Fatalf("load: %v", err)
	}

	if cfg.Editor != "vim" {
		t.Fatalf("expected editor from config, got %q", cfg.Editor)
	}
}

func TestManagerSavesConfig(t *testing.T) {
	root := t.TempDir()
	cwd := filepath.Join(root, "project")
	if err := os.MkdirAll(cwd, 0o755); err != nil {
		t.Fatalf("mkdir cwd: %v", err)
	}
	t.Setenv("XDG_CONFIG_HOME", filepath.Join(root, "config"))

	manager := NewManager(cwd)
	cfg := domain.DefaultConfig()
	cfg.Editor = "emacs"

	if err := manager.Save(cfg); err != nil {
		t.Fatalf("save: %v", err)
	}

	configPath := utils.ConfigPathGlobal()
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Fatal("expected config file to exist")
	}
}

func TestManagerExists(t *testing.T) {
	t.Run("returns false when no config exists", func(t *testing.T) {
		root := t.TempDir()
		cwd := filepath.Join(root, "project")
		if err := os.MkdirAll(cwd, 0o755); err != nil {
			t.Fatalf("mkdir cwd: %v", err)
		}
		t.Setenv("XDG_CONFIG_HOME", filepath.Join(root, "config"))

		manager := NewManager(cwd)
		exists, err := manager.Exists()
		if err != nil {
			t.Fatalf("exists: %v", err)
		}
		if exists {
			t.Error("expected config to not exist")
		}
	})

	t.Run("returns true when global config exists", func(t *testing.T) {
		root := t.TempDir()
		cwd := filepath.Join(root, "project")
		if err := os.MkdirAll(cwd, 0o755); err != nil {
			t.Fatalf("mkdir cwd: %v", err)
		}
		t.Setenv("XDG_CONFIG_HOME", filepath.Join(root, "config"))

		configPath := utils.ConfigPathGlobal()
		if err := os.MkdirAll(filepath.Dir(configPath), 0o755); err != nil {
			t.Fatalf("mkdir config dir: %v", err)
		}
		if err := os.WriteFile(configPath, []byte("editor = \"vim\"\n"), 0o644); err != nil {
			t.Fatalf("write config: %v", err)
		}

		manager := NewManager(cwd)
		exists, err := manager.Exists()
		if err != nil {
			t.Fatalf("exists: %v", err)
		}
		if !exists {
			t.Error("expected config to exist")
		}
	})

	t.Run("returns true when local config exists", func(t *testing.T) {
		root := t.TempDir()
		cwd := filepath.Join(root, "project")
		if err := os.MkdirAll(cwd, 0o755); err != nil {
			t.Fatalf("mkdir cwd: %v", err)
		}
		t.Setenv("XDG_CONFIG_HOME", filepath.Join(root, "config"))

		localConfigPath := utils.ConfigPathLocal(cwd)
		if err := os.MkdirAll(filepath.Dir(localConfigPath), 0o755); err != nil {
			t.Fatalf("mkdir local config dir: %v", err)
		}
		if err := os.WriteFile(localConfigPath, []byte("editor = \"code\"\n"), 0o644); err != nil {
			t.Fatalf("write local config: %v", err)
		}

		manager := NewManager(cwd)
		exists, err := manager.Exists()
		if err != nil {
			t.Fatalf("exists: %v", err)
		}
		if !exists {
			t.Error("expected config to exist")
		}
	})
}

func TestManagerLoadWithOverride(t *testing.T) {
	t.Run("loads from override path", func(t *testing.T) {
		root := t.TempDir()
		cwd := filepath.Join(root, "project")
		if err := os.MkdirAll(cwd, 0o755); err != nil {
			t.Fatalf("mkdir cwd: %v", err)
		}

		overridePath := filepath.Join(root, "custom-config.toml")
		data := []byte("editor = \"emacs\"\nprimary = \"03\"\n")
		if err := os.WriteFile(overridePath, data, 0o644); err != nil {
			t.Fatalf("write override config: %v", err)
		}

		manager := NewManager(cwd)
		cfg, err := manager.LoadWithOverride(overridePath)
		if err != nil {
			t.Fatalf("load with override: %v", err)
		}

		if cfg.Editor != "emacs" {
			t.Errorf("expected editor from override, got %q", cfg.Editor)
		}
		if cfg.Primary != "03" {
			t.Errorf("expected primary from override, got %q", cfg.Primary)
		}
	})

	t.Run("falls back to Load when path is empty", func(t *testing.T) {
		root := t.TempDir()
		cwd := filepath.Join(root, "project")
		if err := os.MkdirAll(cwd, 0o755); err != nil {
			t.Fatalf("mkdir cwd: %v", err)
		}
		t.Setenv("XDG_CONFIG_HOME", filepath.Join(root, "config"))

		manager := NewManager(cwd)
		cfg, err := manager.LoadWithOverride("")
		if err != nil {
			t.Fatalf("load with empty override: %v", err)
		}

		if cfg.Editor == "" {
			t.Error("expected default config to be loaded")
		}
	})
}

func TestManagerLocalOverridesGlobal(t *testing.T) {
	root := t.TempDir()
	cwd := filepath.Join(root, "project")
	if err := os.MkdirAll(cwd, 0o755); err != nil {
		t.Fatalf("mkdir cwd: %v", err)
	}
	t.Setenv("XDG_CONFIG_HOME", filepath.Join(root, "config"))

	// Create global config
	globalPath := utils.ConfigPathGlobal()
	if err := os.MkdirAll(filepath.Dir(globalPath), 0o755); err != nil {
		t.Fatalf("mkdir global config dir: %v", err)
	}
	globalData := []byte("editor = \"vim\"\nprimary = \"01\"\n")
	if err := os.WriteFile(globalPath, globalData, 0o644); err != nil {
		t.Fatalf("write global config: %v", err)
	}

	// Create local config
	localPath := utils.ConfigPathLocal(cwd)
	if err := os.MkdirAll(filepath.Dir(localPath), 0o755); err != nil {
		t.Fatalf("mkdir local config dir: %v", err)
	}
	localData := []byte("editor = \"emacs\"\n")
	if err := os.WriteFile(localPath, localData, 0o644); err != nil {
		t.Fatalf("write local config: %v", err)
	}

	manager := NewManager(cwd)
	cfg, err := manager.Load()
	if err != nil {
		t.Fatalf("load: %v", err)
	}

	if cfg.Editor != "emacs" {
		t.Errorf("expected editor from local config, got %q", cfg.Editor)
	}
	if cfg.Primary != "01" {
		t.Errorf("expected primary from global config, got %q", cfg.Primary)
	}
}

func TestManagerPartialConfig(t *testing.T) {
	t.Run("only overrides specified fields", func(t *testing.T) {
		root := t.TempDir()
		cwd := filepath.Join(root, "project")
		if err := os.MkdirAll(cwd, 0o755); err != nil {
			t.Fatalf("mkdir cwd: %v", err)
		}
		t.Setenv("XDG_CONFIG_HOME", filepath.Join(root, "config"))

		configPath := utils.ConfigPathGlobal()
		if err := os.MkdirAll(filepath.Dir(configPath), 0o755); err != nil {
			t.Fatalf("mkdir config dir: %v", err)
		}

		// Only set editor, leave other fields as defaults
		data := []byte("editor = \"code\"\n")
		if err := os.WriteFile(configPath, data, 0o644); err != nil {
			t.Fatalf("write config: %v", err)
		}

		manager := NewManager(cwd)
		cfg, err := manager.Load()
		if err != nil {
			t.Fatalf("load: %v", err)
		}

		if cfg.Editor != "code" {
			t.Errorf("expected editor from config, got %q", cfg.Editor)
		}
		// Check that defaults are still present
		if cfg.Primary != "02" {
			t.Errorf("expected default primary, got %q", cfg.Primary)
		}
		if cfg.Headings != "15" {
			t.Errorf("expected default headings, got %q", cfg.Headings)
		}
		if !cfg.InteractiveDefault {
			t.Error("expected default interactive_default to be true")
		}
	})
}

func TestManagerColorOverrides(t *testing.T) {
	root := t.TempDir()
	cwd := filepath.Join(root, "project")
	if err := os.MkdirAll(cwd, 0o755); err != nil {
		t.Fatalf("mkdir cwd: %v", err)
	}
	t.Setenv("XDG_CONFIG_HOME", filepath.Join(root, "config"))

	configPath := utils.ConfigPathGlobal()
	if err := os.MkdirAll(filepath.Dir(configPath), 0o755); err != nil {
		t.Fatalf("mkdir config dir: %v", err)
	}

	data := []byte(`
headings = "10"
primary = "04"
secondary = "05"
text = "09"
text_highlight = "11"
description_highlight = "12"
tags = "14"
flags = "15"
muted = "07"
accent = "16"
border = "06"
`)
	if err := os.WriteFile(configPath, data, 0o644); err != nil {
		t.Fatalf("write config: %v", err)
	}

	manager := NewManager(cwd)
	cfg, err := manager.Load()
	if err != nil {
		t.Fatalf("load: %v", err)
	}

	tests := []struct {
		name     string
		got      string
		expected string
	}{
		{"Headings", cfg.Headings, "10"},
		{"Primary", cfg.Primary, "04"},
		{"Secondary", cfg.Secondary, "05"},
		{"Text", cfg.Text, "09"},
		{"TextHighlight", cfg.TextHighlight, "11"},
		{"DescriptionHighlight", cfg.DescriptionHighlight, "12"},
		{"Tags", cfg.Tags, "14"},
		{"Flags", cfg.Flags, "15"},
		{"Muted", cfg.Muted, "07"},
		{"Accent", cfg.Accent, "16"},
		{"Border", cfg.Border, "06"},
	}

	for _, tt := range tests {
		if tt.got != tt.expected {
			t.Errorf("%s = %q, want %q", tt.name, tt.got, tt.expected)
		}
	}
}

func TestManagerBooleanOverrides(t *testing.T) {
	root := t.TempDir()
	cwd := filepath.Join(root, "project")
	if err := os.MkdirAll(cwd, 0o755); err != nil {
		t.Fatalf("mkdir cwd: %v", err)
	}
	t.Setenv("XDG_CONFIG_HOME", filepath.Join(root, "config"))

	configPath := utils.ConfigPathGlobal()
	if err := os.MkdirAll(filepath.Dir(configPath), 0o755); err != nil {
		t.Fatalf("mkdir config dir: %v", err)
	}

	data := []byte("interactive_default = false\n")
	if err := os.WriteFile(configPath, data, 0o644); err != nil {
		t.Fatalf("write config: %v", err)
	}

	manager := NewManager(cwd)
	cfg, err := manager.Load()
	if err != nil {
		t.Fatalf("load: %v", err)
	}

	if cfg.InteractiveDefault {
		t.Error("expected interactive_default to be false from config")
	}
}
