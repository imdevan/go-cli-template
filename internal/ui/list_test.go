package ui

import (
	"testing"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
)

type testItem struct {
	title string
	desc  string
}

func (t testItem) Title() string       { return t.title }
func (t testItem) Description() string { return t.desc }
func (t testItem) FilterValue() string { return t.title }

func TestNewListModel(t *testing.T) {
	theme := Theme{
		Primary:   lipgloss.Color("2"),
		Secondary: lipgloss.Color("6"),
		Text:      lipgloss.Color("7"),
		Muted:     lipgloss.Color("8"),
	}

	items := []list.Item{
		testItem{title: "Item 1", desc: "Description 1"},
		testItem{title: "Item 2", desc: "Description 2"},
	}

	delegate := list.NewDefaultDelegate()
	model := NewListModel(items, delegate, 80, 20, theme)

	if len(model.Items()) != 2 {
		t.Errorf("NewListModel() items count = %d, want 2", len(model.Items()))
	}

	if model.Width() != 80 {
		t.Errorf("NewListModel() width = %d, want 80", model.Width())
	}

	if model.Height() != 20 {
		t.Errorf("NewListModel() height = %d, want 20", model.Height())
	}
}

func TestApplyListStyles(t *testing.T) {
	t.Run("applies styles to non-nil model", func(t *testing.T) {
		theme := Theme{
			Primary:              lipgloss.Color("2"),
			Secondary:            lipgloss.Color("6"),
			Text:                 lipgloss.Color("7"),
			TextHighlight:        lipgloss.Color("11"),
			DescriptionHighlight: lipgloss.Color("12"),
			Muted:                lipgloss.Color("8"),
		}

		model := list.New([]list.Item{}, list.NewDefaultDelegate(), 80, 20)
		
		// Should not panic
		ApplyListStyles(&model, theme)
		
		// Basic verification that function executed
		if model.Width() != 80 {
			t.Error("Model should maintain its properties after applying styles")
		}
	})

	t.Run("handles nil model gracefully", func(t *testing.T) {
		theme := Theme{}
		// Should not panic
		ApplyListStyles(nil, theme)
	})
}

func TestApplyListFilterStyles(t *testing.T) {
	t.Run("applies filter styles to non-nil model", func(t *testing.T) {
		theme := Theme{
			Secondary: lipgloss.Color("6"),
			Text:      lipgloss.Color("7"),
		}

		model := list.New([]list.Item{}, list.NewDefaultDelegate(), 80, 20)
		
		// Should not panic
		ApplyListFilterStyles(&model, theme)
		
		// Basic verification that function executed
		if model.Width() != 80 {
			t.Error("Model should maintain its properties after applying filter styles")
		}
	})

	t.Run("handles nil model gracefully", func(t *testing.T) {
		theme := Theme{}
		// Should not panic
		ApplyListFilterStyles(nil, theme)
	})
}

func TestNewListDelegate(t *testing.T) {
	theme := Theme{
		Primary:              lipgloss.Color("2"),
		TextHighlight:        lipgloss.Color("11"),
		DescriptionHighlight: lipgloss.Color("12"),
	}

	t.Run("default spacing", func(t *testing.T) {
		opts := ListDelegateOptions{}
		delegate := NewListDelegate(theme, opts)

		if !delegate.ShowDescription {
			t.Error("ShowDescription should be true for default spacing")
		}
		if delegate.Height() != 2 {
			t.Errorf("Height = %d, want 2 for default spacing", delegate.Height())
		}
		if delegate.Spacing() != 1 {
			t.Errorf("Spacing = %d, want 1 for default spacing", delegate.Spacing())
		}
	})

	t.Run("space spacing", func(t *testing.T) {
		opts := ListDelegateOptions{Spacing: "space"}
		delegate := NewListDelegate(theme, opts)

		if !delegate.ShowDescription {
			t.Error("ShowDescription should be true for space spacing")
		}
		if delegate.Height() != 2 {
			t.Errorf("Height = %d, want 2 for space spacing", delegate.Height())
		}
		if delegate.Spacing() != 1 {
			t.Errorf("Spacing = %d, want 1 for space spacing", delegate.Spacing())
		}
	})

	t.Run("compact spacing", func(t *testing.T) {
		opts := ListDelegateOptions{Spacing: "compact"}
		delegate := NewListDelegate(theme, opts)

		if delegate.ShowDescription {
			t.Error("ShowDescription should be false for compact spacing")
		}
		if delegate.Height() != 1 {
			t.Errorf("Height = %d, want 1 for compact spacing", delegate.Height())
		}
		if delegate.Spacing() != 0 {
			t.Errorf("Spacing = %d, want 0 for compact spacing", delegate.Spacing())
		}
	})

	t.Run("tight spacing", func(t *testing.T) {
		opts := ListDelegateOptions{Spacing: "tight"}
		delegate := NewListDelegate(theme, opts)

		if !delegate.ShowDescription {
			t.Error("ShowDescription should be true for tight spacing")
		}
		if delegate.Height() != 2 {
			t.Errorf("Height = %d, want 2 for tight spacing", delegate.Height())
		}
		if delegate.Spacing() != 0 {
			t.Errorf("Spacing = %d, want 0 for tight spacing", delegate.Spacing())
		}
	})

	t.Run("custom height overrides default", func(t *testing.T) {
		opts := ListDelegateOptions{
			Spacing: "space",
			Height:  3,
		}
		delegate := NewListDelegate(theme, opts)

		if delegate.Height() != 3 {
			t.Errorf("Height = %d, want 3 for custom height", delegate.Height())
		}
	})

	t.Run("invalid spacing falls back to space", func(t *testing.T) {
		opts := ListDelegateOptions{Spacing: "invalid"}
		delegate := NewListDelegate(theme, opts)

		if !delegate.ShowDescription {
			t.Error("ShowDescription should be true for invalid spacing (fallback)")
		}
		if delegate.Height() != 2 {
			t.Errorf("Height = %d, want 2 for invalid spacing (fallback)", delegate.Height())
		}
		if delegate.Spacing() != 1 {
			t.Errorf("Spacing = %d, want 1 for invalid spacing (fallback)", delegate.Spacing())
		}
	})

	t.Run("applies padding options", func(t *testing.T) {
		opts := ListDelegateOptions{
			PaddingLeft:         2,
			SelectedPaddingLeft: 3,
		}
		delegate := NewListDelegate(theme, opts)

		// Check that padding was applied (styles should have padding set)
		_, _, _, left := delegate.Styles.NormalTitle.GetPadding()
		if left != 2 {
			t.Errorf("NormalTitle left padding = %d, want 2", left)
		}

		_, _, _, left = delegate.Styles.SelectedTitle.GetPadding()
		if left != 3 {
			t.Errorf("SelectedTitle left padding = %d, want 3", left)
		}
	})
}

func TestListFullHelpSections(t *testing.T) {
	model := list.New([]list.Item{}, list.NewDefaultDelegate(), 80, 20)

	t.Run("no options returns empty sections", func(t *testing.T) {
		opts := ListHelpOptions{}
		sections := ListFullHelpSections(model, opts)

		if len(sections) != 0 {
			t.Errorf("sections count = %d, want 0 for no options", len(sections))
		}
	})

	t.Run("paging option includes navigation keys", func(t *testing.T) {
		opts := ListHelpOptions{IncludePaging: true}
		sections := ListFullHelpSections(model, opts)

		if len(sections) != 1 {
			t.Errorf("sections count = %d, want 1 for paging only", len(sections))
		}
		if len(sections) > 0 && len(sections[0]) != 6 {
			t.Errorf("paging section keys count = %d, want 6", len(sections[0]))
		}
	})

	t.Run("filter option includes filter keys", func(t *testing.T) {
		opts := ListHelpOptions{IncludeFilter: true}
		sections := ListFullHelpSections(model, opts)

		if len(sections) != 1 {
			t.Errorf("sections count = %d, want 1 for filter only", len(sections))
		}
		if len(sections) > 0 && len(sections[0]) != 4 {
			t.Errorf("filter section keys count = %d, want 4", len(sections[0]))
		}
	})

	t.Run("quit option includes quit key", func(t *testing.T) {
		opts := ListHelpOptions{IncludeQuit: true}
		sections := ListFullHelpSections(model, opts)

		if len(sections) != 1 {
			t.Errorf("sections count = %d, want 1 for quit only", len(sections))
		}
		if len(sections) > 0 && len(sections[0]) != 1 {
			t.Errorf("quit section keys count = %d, want 1", len(sections[0]))
		}
	})

	t.Run("all options includes all sections", func(t *testing.T) {
		opts := ListHelpOptions{
			IncludePaging: true,
			IncludeFilter: true,
			IncludeQuit:   true,
		}
		sections := ListFullHelpSections(model, opts)

		if len(sections) != 2 {
			t.Errorf("sections count = %d, want 2 for all options", len(sections))
		}
		if len(sections) > 0 && len(sections[0]) != 6 {
			t.Errorf("first section (paging) keys count = %d, want 6", len(sections[0]))
		}
		if len(sections) > 1 && len(sections[1]) != 5 {
			t.Errorf("second section (filter+quit) keys count = %d, want 5", len(sections[1]))
		}
	})

	t.Run("filter and quit combined", func(t *testing.T) {
		opts := ListHelpOptions{
			IncludeFilter: true,
			IncludeQuit:   true,
		}
		sections := ListFullHelpSections(model, opts)

		if len(sections) != 1 {
			t.Errorf("sections count = %d, want 1 for filter+quit", len(sections))
		}
		if len(sections) > 0 && len(sections[0]) != 5 {
			t.Errorf("filter+quit section keys count = %d, want 5", len(sections[0]))
		}
	})
}

func TestListDelegateOptions_Defaults(t *testing.T) {
	t.Run("zero values work correctly", func(t *testing.T) {
		theme := Theme{
			Primary:              lipgloss.Color("2"),
			TextHighlight:        lipgloss.Color("11"),
			DescriptionHighlight: lipgloss.Color("12"),
		}

		opts := ListDelegateOptions{}
		delegate := NewListDelegate(theme, opts)

		// Should use defaults without panicking
		if delegate.Height() == 0 {
			t.Error("Height should have a default value")
		}
	})
}
