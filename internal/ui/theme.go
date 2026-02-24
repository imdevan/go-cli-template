package ui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"

	"go-cli-template/internal/domain"
)

// Theme holds configurable colors for UI output.
type Theme struct {
	Headings             lipgloss.Color
	Primary              lipgloss.Color
	Secondary            lipgloss.Color
	Text                 lipgloss.Color
	TextHighlight        lipgloss.Color
	DescriptionHighlight lipgloss.Color
	Tags                 lipgloss.Color
	Flags                lipgloss.Color
	Muted                lipgloss.Color
	Border               lipgloss.Color
}

// ThemeFromConfig builds a theme with safe fallbacks.
func ThemeFromConfig(cfg domain.Config) Theme {
	return Theme{
		Headings:             resolveColor(cfg.Headings, "15"),
		Primary:              resolveColor(cfg.Primary, "02"),
		Secondary:            resolveColor(cfg.Secondary, "06"),
		Text:                 resolveColor(cfg.Text, "07"),
		TextHighlight:        resolveColor(resolveFallback(cfg.TextHighlight, cfg.Secondary), "06"),
		DescriptionHighlight: resolveColor(resolveFallback(cfg.DescriptionHighlight, cfg.Secondary), "06"),
		Tags:                 resolveColor(resolveFallback(cfg.Tags, cfg.Accent), "13"),
		Flags:                resolveColor(cfg.Flags, "12"),
		Muted:                resolveColor(cfg.Muted, "08"),
		Border:               resolveColor(cfg.Border, "08"),
	}
}

func resolveColor(value, fallback string) lipgloss.Color {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		trimmed = fallback
	}
	return lipgloss.Color(trimmed)
}

func resolveFallback(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return value
		}
	}
	return ""
}
