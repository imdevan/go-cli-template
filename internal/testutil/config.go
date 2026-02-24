package testutil

import (
	"testing"

	"go-cli-template/internal/domain"
)

// ConfigOverrides lets tests override specific config fields.
type ConfigOverrides struct {
	Editor *string
}

// NewConfig returns a default config with optional overrides and temp XDG paths.
func NewConfig(t *testing.T, overrides ConfigOverrides) domain.Config {
	t.Helper()
	_, _, _ = WithTempXDG(t)
	cfg := domain.DefaultConfig()
	applyOverrides(&cfg, overrides)
	return cfg
}

func applyOverrides(cfg *domain.Config, overrides ConfigOverrides) {
	if overrides.Editor != nil {
		cfg.Editor = *overrides.Editor
	}
}
