package ui

import (
	"io"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// ItemWithMetadata is an optional interface that list items can implement
// to provide a third row of metadata under the description.
// Note: Currently items should embed metadata in Description() directly.
// This interface is reserved for future use when custom delegate rendering is implemented.
type ItemWithMetadata interface {
	list.Item
	Metadata() string
}

// ListDelegateOptions configures shared list presentation settings.
type ListDelegateOptions struct {
	Height              int
	PaddingLeft         int
	SelectedPaddingLeft int
	Spacing             string // "compact", "tight", or "space" (default)
	ShowMetadata        bool   // Enable metadata row support
	MetadataIndent      int    // Indentation for metadata row (default: 1)
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
	
	// Title style: primary color, bold, no background, left aligned
	model.Styles.Title = model.Styles.Title.
		Foreground(theme.Primary).
		Background(nil).
		Bold(true).
		Align(lipgloss.Left).
		Padding(0).
		Margin(0)
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
func NewListDelegate(theme Theme, opts ListDelegateOptions) list.ItemDelegate {
	// If metadata is enabled, use custom delegate
	if opts.ShowMetadata {
		return newMetadataDelegate(theme, opts)
	}
	
	// Otherwise use default delegate
	return newDefaultDelegate(theme, opts)
}

func newDefaultDelegate(theme Theme, opts ListDelegateOptions) list.DefaultDelegate {
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

// metadataDelegate wraps the default delegate and adds metadata row support.
type metadataDelegate struct {
	defaultDelegate list.DefaultDelegate
	theme           Theme
	metadataIndent  int
}

func newMetadataDelegate(theme Theme, opts ListDelegateOptions) *metadataDelegate {
	delegate := newDefaultDelegate(theme, opts)
	
	// Set height to 3 for title + description + metadata
	delegate.SetHeight(3)
	
	metadataIndent := opts.MetadataIndent
	if metadataIndent == 0 {
		metadataIndent = 1 // Default to 1 space to align with description
	}
	
	return &metadataDelegate{
		defaultDelegate: delegate,
		theme:           theme,
		metadataIndent:  metadataIndent,
	}
}

func (d *metadataDelegate) Height() int {
	return d.defaultDelegate.Height()
}

func (d *metadataDelegate) Spacing() int {
	return d.defaultDelegate.Spacing()
}

func (d *metadataDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd {
	return d.defaultDelegate.Update(msg, m)
}

func (d *metadataDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	// Check if item implements ItemWithMetadata
	itemWithMeta, hasMetadata := item.(ItemWithMetadata)
	
	if !hasMetadata {
		// Fall back to default rendering
		d.defaultDelegate.Render(w, m, index, item)
		return
	}
	
	metadata := itemWithMeta.Metadata()
	if metadata == "" {
		// No metadata, use default rendering
		d.defaultDelegate.Render(w, m, index, item)
		return
	}
	
	// Assert to DefaultItem for Title() and Description() access
	defaultItem, ok := item.(list.DefaultItem)
	if !ok {
		// Item doesn't implement DefaultItem, fall back
		d.defaultDelegate.Render(w, m, index, item)
		return
	}
	
	// Create a wrapper that includes metadata in description
	wrapper := &metadataItemWrapper{
		item:     defaultItem,
		metadata: metadata,
		indent:   d.metadataIndent,
	}
	
	d.defaultDelegate.Render(w, m, index, wrapper)
}

// metadataItemWrapper wraps a list item and appends metadata to its description.
type metadataItemWrapper struct {
	item     list.DefaultItem
	metadata string
	indent   int
}

func (w *metadataItemWrapper) Title() string {
	return w.item.Title()
}

func (w *metadataItemWrapper) Description() string {
	desc := w.item.Description()
	if w.metadata != "" {
		indentStr := ""
		for i := 0; i < w.indent; i++ {
			indentStr += " "
		}
		desc = desc + "\n" + indentStr + w.metadata
	}
	return desc
}

func (w *metadataItemWrapper) FilterValue() string {
	return w.item.FilterValue()
}
