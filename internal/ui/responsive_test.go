package ui


import (
	"testing"

	"github.com/charmbracelet/lipgloss"
)

func TestResponsiveManager_Breakpoints(t *testing.T) {
	tests := []struct {
		name           string
		width          int
		wantBreakpoint Breakpoint
	}{
		{"extra small - 20", 20, BreakpointXS},
		{"extra small - 39", 39, BreakpointXS},
		{"small - 40", 40, BreakpointSM},
		{"small - 50", 50, BreakpointSM},
		{"small - 59", 59, BreakpointSM},
		{"medium - 60", 60, BreakpointMD},
		{"medium - 70", 70, BreakpointMD},
		{"medium - 79", 79, BreakpointMD},
		{"large - 80", 80, BreakpointLG},
		{"large - 90", 90, BreakpointLG},
		{"large - 99", 99, BreakpointLG},
		{"extra large - 100", 100, BreakpointXL},
		{"extra large - 120", 120, BreakpointXL},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rm := NewResponsiveManager(tt.width)
			if rm.Breakpoint() != tt.wantBreakpoint {
				t.Errorf("NewResponsiveManager(%d).Breakpoint() = %v, want %v",
					tt.width, rm.Breakpoint(), tt.wantBreakpoint)
			}
			if rm.Width() != tt.width {
				t.Errorf("NewResponsiveManager(%d).Width() = %d, want %d",
					tt.width, rm.Width(), tt.width)
			}
		})
	}
}

func TestResponsiveManager_SetWidth(t *testing.T) {
	rm := NewResponsiveManager(50)
	if rm.Width() != 50 {
		t.Errorf("initial width = %d, want 50", rm.Width())
	}
	if rm.Breakpoint() != BreakpointSM {
		t.Errorf("initial breakpoint = %v, want BreakpointSM", rm.Breakpoint())
	}

	rm.SetWidth(80)
	if rm.Width() != 80 {
		t.Errorf("updated width = %d, want 80", rm.Width())
	}
	if rm.Breakpoint() != BreakpointLG {
		t.Errorf("updated breakpoint = %v, want BreakpointLG", rm.Breakpoint())
	}
}

func TestResponsiveManager_IsAtLeast(t *testing.T) {
	tests := []struct {
		width int
		check Breakpoint
		want  bool
	}{
		{30, BreakpointXS, true},
		{30, BreakpointSM, false},
		{50, BreakpointSM, true},
		{50, BreakpointMD, false},
		{70, BreakpointMD, true},
		{70, BreakpointLG, false},
		{90, BreakpointLG, true},
		{90, BreakpointXL, false},
		{110, BreakpointXL, true},
	}

	for _, tt := range tests {
		rm := NewResponsiveManager(tt.width)
		if got := rm.IsAtLeast(tt.check); got != tt.want {
			t.Errorf("ResponsiveManager(%d).IsAtLeast(%v) = %v, want %v",
				tt.width, tt.check, got, tt.want)
		}
	}
}

func TestResponsiveManager_IsAtMost(t *testing.T) {
	tests := []struct {
		width int
		check Breakpoint
		want  bool
	}{
		{30, BreakpointXS, true},
		{30, BreakpointSM, true},
		{50, BreakpointSM, true},
		{50, BreakpointXS, false},
		{70, BreakpointMD, true},
		{70, BreakpointSM, false},
		{90, BreakpointLG, true},
		{90, BreakpointMD, false},
		{110, BreakpointXL, true},
		{110, BreakpointLG, false},
	}

	for _, tt := range tests {
		rm := NewResponsiveManager(tt.width)
		if got := rm.IsAtMost(tt.check); got != tt.want {
			t.Errorf("ResponsiveManager(%d).IsAtMost(%v) = %v, want %v",
				tt.width, tt.check, got, tt.want)
		}
	}
}

func TestResponsiveManager_ShouldShowBorders(t *testing.T) {
	tests := []struct {
		width int
		want  bool
	}{
		{30, false},  // XS
		{50, false},  // SM
		{60, true},   // MD
		{70, true},   // MD
		{80, true},   // LG
		{100, true},  // XL
		{120, true},  // XL
	}

	for _, tt := range tests {
		rm := NewResponsiveManager(tt.width)
		if got := rm.ShouldShowBorders(); got != tt.want {
			t.Errorf("ResponsiveManager(%d).ShouldShowBorders() = %v, want %v",
				tt.width, got, tt.want)
		}
	}
}

func TestResponsiveManager_MaxContentWidth(t *testing.T) {
	tests := []struct {
		width int
		want  int
	}{
		{30, 30},   // XS: 30 - 0 = 30 (no padding/margin/border)
		{26, 26},   // XS: 26 - 0 = 26
		{22, 22},   // XS: 22 - 0 = 22
		{10, 20},   // XS: 10 - 0 = 10, clamped to 20 minimum
		{50, 48},   // SM: 50 - 2 = 48
		{40, 38},   // SM: 40 - 2 = 38
		{60, 54},   // MD: 60 - 6 (padding + border + margin) = 54
		{70, 64},   // MD: 70 - 6 = 64
		{80, 74},   // LG: 80 - 6 = 74
		{100, 94},  // XL: 100 - 6 = 94
	}

	for _, tt := range tests {
		rm := NewResponsiveManager(tt.width)
		if got := rm.MaxContentWidth(); got != tt.want {
			t.Errorf("ResponsiveManager(%d).MaxContentWidth() = %d, want %d",
				tt.width, got, tt.want)
		}
	}
}

func TestResponsiveManager_AdaptiveFrameStyle(t *testing.T) {
	theme := Theme{
		Border: lipgloss.Color("8"),
	}

	tests := []struct {
		name       string
		width      int
		wantPadding bool
		wantMargin bool
		wantBorder bool
	}{
		{"extra small - no styling", 30, false, false, false},
		{"small - padding only", 50, true, false, false},
		{"medium - full styling", 70, true, true, true},
		{"large - full styling", 90, true, true, true},
		{"extra large - full styling", 110, true, true, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
rm := NewResponsiveManager(tt.width)
style := rm.AdaptiveFrameStyle(theme)

// Check padding
top, right, bottom, left := style.GetPadding()
			hasPadding := top > 0 || right > 0 || bottom > 0 || left > 0
			if hasPadding != tt.wantPadding {
				t.Errorf("width %d: hasPadding = %v, want %v", tt.width, hasPadding, tt.wantPadding)
			}

			// Check margin
			top, right, bottom, left = style.GetMargin()
			hasMargin := top > 0 || right > 0 || bottom > 0 || left > 0
			if hasMargin != tt.wantMargin {
				t.Errorf("width %d: hasMargin = %v, want %v", tt.width, hasMargin, tt.wantMargin)
			}

			// Check border
			hasBorder := style.GetBorderTop() || style.GetBorderRight() || style.GetBorderBottom() || style.GetBorderLeft()
			if hasBorder != tt.wantBorder {
				t.Errorf("width %d: hasBorder = %v, want %v", tt.width, hasBorder, tt.wantBorder)
			}
		})
	}
}

func TestResponsiveManager_GetContentInsets(t *testing.T) {
	tests := []struct {
		name           string
		width          int
		wantHorizontal int
		wantVertical   int
	}{
		{"extra small - no insets", 30, 0, 0},
		{"extra small boundary", 39, 0, 0},
		{"small - padding only", 40, 2, 2},
		{"small mid", 50, 2, 2},
		{"small boundary", 59, 2, 2},
		{"medium - full frame", 60, 6, 6},
		{"medium mid", 70, 6, 6},
		{"medium boundary", 79, 6, 6},
		{"large - more spacing", 80, 8, 6},
		{"large mid", 90, 8, 6},
		{"large boundary", 99, 8, 6},
		{"extra large - max spacing", 100, 9, 7},
		{"extra large wide", 120, 9, 7},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
rm := NewResponsiveManager(tt.width)
h, v := rm.GetContentInsets()
			if h != tt.wantHorizontal {
				t.Errorf("GetContentInsets() horizontal = %d, want %d", h, tt.wantHorizontal)
			}
			if v != tt.wantVertical {
				t.Errorf("GetContentInsets() vertical = %d, want %d", v, tt.wantVertical)
			}
		})
	}
}
