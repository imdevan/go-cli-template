package ui

import "github.com/charmbracelet/lipgloss"

// Breakpoint defines responsive size categories.
type Breakpoint int

const (
	// BreakpointXS represents extra small screens (< 40 columns).
	BreakpointXS Breakpoint = iota
	// BreakpointSM represents small screens (40-60 columns).
	BreakpointSM
	// BreakpointMD represents medium screens (60-80 columns).
	BreakpointMD
	// BreakpointLG represents large screens (80-100 columns).
	BreakpointLG
	// BreakpointXL represents extra large screens (>= 100 columns).
	BreakpointXL
)

// ResponsiveManager tracks view width and provides responsive utilities.
type ResponsiveManager struct {
	width      int
	breakpoint Breakpoint
}

// NewResponsiveManager creates a new responsive manager with the given width.
func NewResponsiveManager(width int) *ResponsiveManager {
	rm := &ResponsiveManager{}
	rm.SetWidth(width)
	return rm
}

// SetWidth updates the tracked width and recalculates the breakpoint.
func (rm *ResponsiveManager) SetWidth(width int) {
	rm.width = width
	rm.breakpoint = rm.calculateBreakpoint()
}

// Width returns the current tracked width.
func (rm *ResponsiveManager) Width() int {
	return rm.width
}

// Breakpoint returns the current breakpoint category.
func (rm *ResponsiveManager) Breakpoint() Breakpoint {
	return rm.breakpoint
}

// calculateBreakpoint determines the breakpoint based on width.
func (rm *ResponsiveManager) calculateBreakpoint() Breakpoint {
	switch {
	case rm.width < 40:
		return BreakpointXS
	case rm.width < 60:
		return BreakpointSM
	case rm.width < 80:
		return BreakpointMD
	case rm.width < 100:
		return BreakpointLG
	default:
		return BreakpointXL
	}
}

// IsAtLeast checks if the current breakpoint is at least the given size.
func (rm *ResponsiveManager) IsAtLeast(bp Breakpoint) bool {
	return rm.breakpoint >= bp
}

// IsAtMost checks if the current breakpoint is at most the given size.
func (rm *ResponsiveManager) IsAtMost(bp Breakpoint) bool {
	return rm.breakpoint <= bp
}

// ShouldShowBorders returns true if borders should be rendered (MD or larger).
func (rm *ResponsiveManager) ShouldShowBorders() bool {
	return rm.IsAtLeast(BreakpointMD)
}

// AdaptiveFrameStyle returns a frame style that adapts based on breakpoint.
// XS: No padding, no margin, no border
// SM: Padding only, no margin, no border
// MD+: Full styling with padding, margin, and border
func (rm *ResponsiveManager) AdaptiveFrameStyle(theme Theme) lipgloss.Style {
	style := lipgloss.NewStyle()

	if rm.breakpoint == BreakpointXS {
		// Extra small: no padding, no margin, no border
		return style
	} else if rm.breakpoint == BreakpointSM {
		// Small: padding only, no margin, no border
		return style.Padding(1, 1)
	} else if rm.breakpoint == BreakpointMD || rm.breakpoint == BreakpointLG {
		// MD, LG
		return style.
			Padding(1, 2, 1, 1).
			Margin(1, 1).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(theme.Border)
	} else {
		// XL: full styling
		return style.
			Padding(1, 4, 1, 1).
			Margin(1, 1).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(theme.Border)
	}
}

// GetContentInsets returns the horizontal and vertical insets for content
// based on the current breakpoint. These values account for padding, margins,
// and borders applied by AdaptiveFrameStyle.
func (rm *ResponsiveManager) GetContentInsets() (horizontal, vertical int) {
	switch rm.breakpoint {
	case BreakpointXS:
		// No styling: no insets
		return 0, 0
	case BreakpointSM:
		// Padding only: 1 left + 1 right = 2, 1 top + 1 bottom = 2
		return 2, 2
	case BreakpointMD:
		// Padding + border + margin: (1+1+1)*2 = 6 horizontal, (1+1+1)*2 = 6 vertical
		return 6, 6
	case BreakpointLG:
		// Same as MD but with slightly more breathing room
		return 8, 6
	default:
		// BreakpointXL: maximum spacing
		return 9, 7
	}
}

// MaxContentWidth returns the maximum usable content width accounting for padding/borders.
func (rm *ResponsiveManager) MaxContentWidth() int {
	switch rm.breakpoint {
	case BreakpointXS:
		// No padding, margin, or border
		return max(rm.width, 20)
	case BreakpointSM:
		// Padding only (2)
		return max(rm.width-2, 20)
	default:
		// MD+: padding (2) + border (2) + margin (2)
		return max(rm.width-6, 20)
	}
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// GetListDimensions returns responsive width and height for a list component.
// It accounts for frame insets and provides sensible defaults based on breakpoint.
func (rm *ResponsiveManager) GetListDimensions(availableWidth, availableHeight int) (width, height int) {
	h, v := rm.GetContentInsets()

	// Calculate usable dimensions
	width = availableWidth - h
	height = availableHeight - v

	// Apply minimum constraints based on breakpoint
	switch rm.breakpoint {
	case BreakpointXS:
		// Extra small: maximize space, minimal constraints
		width = max(width, 20)
		height = max(height, 5)
	case BreakpointSM:
		// Small: slightly more breathing room
		width = max(width, 30)
		height = max(height, 8)
	case BreakpointMD:
		// Medium: comfortable viewing
		width = max(width, 40)
		height = max(height, 10)
	case BreakpointLG:
		// Large: spacious layout
		width = max(width, 60)
		height = max(height, 15)
	default:
		// Extra large: maximum comfort
		width = max(width, 80)
		height = max(height, 20)
	}

	return width, height
}

// ShouldShowFullHelp returns true if there's enough space to show full help.
func (rm *ResponsiveManager) ShouldShowFullHelp() bool {
	return rm.IsAtLeast(BreakpointMD)
}

// GetHelpHeight returns the recommended height for help text based on breakpoint.
func (rm *ResponsiveManager) GetHelpHeight() int {
	switch rm.breakpoint {
	case BreakpointXS:
		return 1 // Minimal help
	case BreakpointSM:
		return 2 // Short help
	default:
		return 3 // Full help with multiple lines
	}
}
