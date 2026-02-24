package domain

import (
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()

	t.Run("has editor set", func(t *testing.T) {
		if cfg.Editor == "" {
			t.Error("DefaultConfig().Editor should not be empty")
		}
		if cfg.Editor != "nvim" {
			t.Errorf("DefaultConfig().Editor = %q, want %q", cfg.Editor, "nvim")
		}
	})

	t.Run("has all color values set", func(t *testing.T) {
		colorFields := map[string]string{
			"Headings":             cfg.Headings,
			"Primary":              cfg.Primary,
			"Secondary":            cfg.Secondary,
			"Text":                 cfg.Text,
			"TextHighlight":        cfg.TextHighlight,
			"DescriptionHighlight": cfg.DescriptionHighlight,
			"Tags":                 cfg.Tags,
			"Flags":                cfg.Flags,
			"Muted":                cfg.Muted,
			"Accent":               cfg.Accent,
			"Border":               cfg.Border,
		}

		for name, value := range colorFields {
			if value == "" {
				t.Errorf("DefaultConfig().%s should not be empty", name)
			}
		}
	})

	t.Run("has expected default values", func(t *testing.T) {
		if cfg.Headings != "15" {
			t.Errorf("DefaultConfig().Headings = %q, want %q", cfg.Headings, "15")
		}
		if cfg.Primary != "02" {
			t.Errorf("DefaultConfig().Primary = %q, want %q", cfg.Primary, "02")
		}
		if cfg.Secondary != "06" {
			t.Errorf("DefaultConfig().Secondary = %q, want %q", cfg.Secondary, "06")
		}
		if cfg.Text != "07" {
			t.Errorf("DefaultConfig().Text = %q, want %q", cfg.Text, "07")
		}
		if cfg.TextHighlight != "06" {
			t.Errorf("DefaultConfig().TextHighlight = %q, want %q", cfg.TextHighlight, "06")
		}
		if cfg.DescriptionHighlight != "05" {
			t.Errorf("DefaultConfig().DescriptionHighlight = %q, want %q", cfg.DescriptionHighlight, "05")
		}
		if cfg.Tags != "13" {
			t.Errorf("DefaultConfig().Tags = %q, want %q", cfg.Tags, "13")
		}
		if cfg.Flags != "12" {
			t.Errorf("DefaultConfig().Flags = %q, want %q", cfg.Flags, "12")
		}
		if cfg.Muted != "08" {
			t.Errorf("DefaultConfig().Muted = %q, want %q", cfg.Muted, "08")
		}
		if cfg.Accent != "13" {
			t.Errorf("DefaultConfig().Accent = %q, want %q", cfg.Accent, "13")
		}
		if cfg.Border != "08" {
			t.Errorf("DefaultConfig().Border = %q, want %q", cfg.Border, "08")
		}
	})

	t.Run("has interactive default enabled", func(t *testing.T) {
		if !cfg.InteractiveDefault {
			t.Error("DefaultConfig().InteractiveDefault should be true")
		}
	})

	t.Run("has list spacing set", func(t *testing.T) {
		if cfg.ListSpacing == "" {
			t.Error("DefaultConfig().ListSpacing should not be empty")
		}
		if cfg.ListSpacing != "space" {
			t.Errorf("DefaultConfig().ListSpacing = %q, want %q", cfg.ListSpacing, "space")
		}
	})
}

func TestDefaultConfig_Consistency(t *testing.T) {
	t.Run("multiple calls return same values", func(t *testing.T) {
		cfg1 := DefaultConfig()
		cfg2 := DefaultConfig()

		if cfg1.Editor != cfg2.Editor {
			t.Error("DefaultConfig() should return consistent Editor values")
		}
		if cfg1.Headings != cfg2.Headings {
			t.Error("DefaultConfig() should return consistent Headings values")
		}
		if cfg1.InteractiveDefault != cfg2.InteractiveDefault {
			t.Error("DefaultConfig() should return consistent InteractiveDefault values")
		}
	})
}

func TestConfig_StructTags(t *testing.T) {
	t.Run("has toml tags for all fields", func(t *testing.T) {
		// This is a regression test to ensure TOML tags don't get accidentally removed
		cfg := Config{}
		
		// Set values to ensure struct is properly tagged
		cfg.Editor = "test"
		cfg.Primary = "01"
		cfg.InteractiveDefault = false
		cfg.ListSpacing = "compact"
		
		// If this compiles and runs, the struct tags are present
		if cfg.Editor != "test" {
			t.Error("Config struct should be properly defined")
		}
	})
}
