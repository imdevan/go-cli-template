package ui

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
)

// FrameStyle defines a shared container style for bordered views.
func FrameStyle(theme Theme) lipgloss.Style {
	return lipgloss.NewStyle().
		Padding(1, 1).
		Margin(1, 1).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(theme.Border)
}

// FrameSizeOptions configures shared list sizing inside framed views.
type FrameSizeOptions struct {
	HorizontalInset int
	VerticalInset   int
	MinWidth        int
	MinHeight       int
}

// FrameSizeDefaults returns base inset values for framed views.
func FrameSizeDefaults() FrameSizeOptions {
	return FrameSizeOptions{
		HorizontalInset: 9,
		VerticalInset:   7,
		MinWidth:        40,
		MinHeight:       8,
	}
}

// ApplyFrameListSize clamps a list to fit a framed container.
func ApplyFrameListSize(model *list.Model, width, height int, opts FrameSizeOptions) {
	if model == nil {
		return
	}
	defaults := FrameSizeDefaults()
	if opts.HorizontalInset == 0 {
		opts.HorizontalInset = defaults.HorizontalInset
	}
	if opts.VerticalInset == 0 {
		opts.VerticalInset = defaults.VerticalInset
	}
	if opts.MinWidth == 0 {
		opts.MinWidth = defaults.MinWidth
	}
	if opts.MinHeight == 0 {
		opts.MinHeight = defaults.MinHeight
	}
	listWidth := width - opts.HorizontalInset
	listHeight := height - opts.VerticalInset
	if listWidth < opts.MinWidth {
		listWidth = opts.MinWidth
	}
	if listHeight < opts.MinHeight {
		listHeight = opts.MinHeight
	}
	model.SetSize(listWidth, listHeight)
}
