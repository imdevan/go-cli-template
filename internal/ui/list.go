package ui

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
)

// ListDelegateOptions configures shared list presentation settings.
type ListDelegateOptions struct {
	Height              int
	PaddingLeft         int
	SelectedPaddingLeft int
	Spacing             string // "compact", "tight", or "space" (default)
}

// NewListModel creates a list with shared styles applied.
func NewListModel(items []list.Item, delegate list.ItemDelegate, width, height int, theme Theme) list.Model {
	model := list.New(items, delegate, width, height)
	ApplyListStyles(&model, theme)
	return model
}

// ApplyListStyles sets shared list styles.
func ApplyListStyles(model *list.Model, theme Theme) {
	if model == nil {
		return
	}
	ApplyListFilterStyles(model, theme)
	model.Styles.NoItems = model.Styles.NoItems.Foreground(theme.Muted)
	model.Styles.StatusBar = model.Styles.StatusBar.Foreground(theme.Muted)
	model.Styles.StatusEmpty = model.Styles.StatusEmpty.Foreground(theme.Muted)
	model.Styles.StatusBarActiveFilter = model.Styles.StatusBarActiveFilter.Foreground(theme.Secondary)
	model.Styles.StatusBarFilterCount = model.Styles.StatusBarFilterCount.Foreground(theme.Muted)
	model.Styles.HelpStyle = model.Styles.HelpStyle.Foreground(theme.Muted)
	model.Styles.PaginationStyle = model.Styles.PaginationStyle.Foreground(theme.Muted)
	model.Styles.ActivePaginationDot = model.Styles.ActivePaginationDot.Foreground(theme.Secondary)
	model.Styles.InactivePaginationDot = model.Styles.InactivePaginationDot.Foreground(theme.Muted)
	model.Styles.DividerDot = model.Styles.DividerDot.Foreground(theme.Muted)
}

// ApplyListFilterStyles sets shared filter styles for lists.
func ApplyListFilterStyles(model *list.Model, theme Theme) {
	if model == nil {
		return
	}
	model.Styles.FilterPrompt = model.Styles.FilterPrompt.Foreground(theme.Secondary)
	model.Styles.FilterCursor = model.Styles.FilterCursor.Foreground(theme.Secondary)
	model.FilterInput.PromptStyle = model.FilterInput.PromptStyle.Foreground(theme.Secondary)
	model.FilterInput.Cursor.Style = model.FilterInput.Cursor.Style.Foreground(theme.Secondary)
	model.FilterInput.TextStyle = model.FilterInput.TextStyle.Foreground(theme.Text)
	model.Styles.DefaultFilterCharacterMatch = model.Styles.DefaultFilterCharacterMatch.Foreground(theme.Secondary)
}

// NewListDelegate provides shared list focus styles.
func NewListDelegate(theme Theme, opts ListDelegateOptions) list.DefaultDelegate {
	delegate := list.NewDefaultDelegate()
	delegate.Styles.SelectedTitle = delegate.Styles.SelectedTitle.Foreground(theme.TextHighlight).BorderForeground(theme.Primary).Bold(true)
	delegate.Styles.SelectedDesc = delegate.Styles.SelectedDesc.Foreground(theme.DescriptionHighlight).BorderForeground(theme.Primary)
	
	// Apply spacing configuration
	spacing := opts.Spacing
	if spacing == "" {
		spacing = "space" // default
	}
	
	switch spacing {
	case "compact":
		// Only show title, no description
		delegate.ShowDescription = false
		delegate.SetHeight(1)
		delegate.SetSpacing(0)
	case "tight":
		// Show title and description with no margin
		delegate.ShowDescription = true
		delegate.SetHeight(2)
		delegate.SetSpacing(0)
	case "space":
		// Current default: title and description with spacing
		delegate.ShowDescription = true
		if opts.Height > 0 {
			delegate.SetHeight(opts.Height)
		} else {
			delegate.SetHeight(2)
		}
		delegate.SetSpacing(1)
	default:
		// Fallback to space
		delegate.ShowDescription = true
		if opts.Height > 0 {
			delegate.SetHeight(opts.Height)
		} else {
			delegate.SetHeight(2)
		}
		delegate.SetSpacing(1)
	}
	
	if opts.PaddingLeft > 0 {
		delegate.Styles.NormalTitle = delegate.Styles.NormalTitle.Padding(0, 0, 0, opts.PaddingLeft)
		delegate.Styles.NormalDesc = delegate.Styles.NormalDesc.Padding(0, 0, 0, opts.PaddingLeft)
	}
	if opts.SelectedPaddingLeft > 0 {
		delegate.Styles.SelectedTitle = delegate.Styles.SelectedTitle.Padding(0, 0, 0, opts.SelectedPaddingLeft)
		delegate.Styles.SelectedDesc = delegate.Styles.SelectedDesc.Padding(0, 0, 0, opts.SelectedPaddingLeft)
	}
	return delegate
}

type ListHelpOptions struct {
	IncludeFilter bool
	IncludePaging bool
	IncludeQuit   bool
}

func ListFullHelpSections(model list.Model, opts ListHelpOptions) [][]key.Binding {
	sections := make([][]key.Binding, 0, 2)
	if opts.IncludePaging {
		sections = append(sections, []key.Binding{
			model.KeyMap.CursorUp,
			model.KeyMap.CursorDown,
			model.KeyMap.NextPage,
			model.KeyMap.PrevPage,
			model.KeyMap.GoToStart,
			model.KeyMap.GoToEnd,
		})
	}
	if opts.IncludeFilter || opts.IncludeQuit {
		section := make([]key.Binding, 0, 5)
		if opts.IncludeFilter {
			section = append(section,
				model.KeyMap.Filter,
				model.KeyMap.ClearFilter,
				model.KeyMap.AcceptWhileFiltering,
				model.KeyMap.CancelWhileFiltering,
			)
		}
		if opts.IncludeQuit {
			section = append(section, model.KeyMap.Quit)
		}
		sections = append(sections, section)
	}
	return sections
}
