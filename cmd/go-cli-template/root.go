package main

import (
	"fmt"
	"os"
	"runtime/debug"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"

	"go-cli-template/internal/config"
	"go-cli-template/internal/domain"
	pkg "go-cli-template/internal/package"
	"go-cli-template/internal/ui"
	"go-cli-template/internal/utils"
)

// Metadata loaded from package.toml at build time
var (
	version = pkg.Version()
	name    = pkg.Name()
	short   = pkg.Short()
)

type rootOptions struct {
	configPath  string
	showVersion bool
}

var rootCmd = newRootCmd()

// Execute is the CLI entrypoint.
func Execute() error {
	return rootCmd.Execute()
}

func newRootCmd() *cobra.Command {
	opts := &rootOptions{}
	cmd := &cobra.Command{
		Use:   name,
		Short: short,
		Args:  cobra.ArbitraryArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if opts.showVersion {
				ver := resolvedVersion()
				cmd.Printf("%s\n", ver)
				return nil
			}
			// Interactive file manager will go here
			return runInteractive(cmd, opts, args)
		},
	}

	cmd.Flags().StringVarP(&opts.configPath, "config", "c", "", "config file path")
	cmd.Flags().BoolVarP(&opts.showVersion, "version", "V", false, "print version information")

	cmd.AddCommand(newConfigCmd())
	cmd.AddCommand(newCompletionCmd())

	return cmd
}

func resolvedVersion() string {
	ver := version
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return ver
	}
	if ver == "dev" && strings.TrimSpace(info.Main.Version) != "" && info.Main.Version != "(devel)" {
		ver = info.Main.Version
	}
	return ver
}

func runInteractive(cmd *cobra.Command, opts *rootOptions, args []string) error {
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get working directory: %w", err)
	}

	manager := config.NewManager(cwd)
	var cfg domain.Config
	if opts.configPath != "" {
		cfg, err = manager.LoadWithOverride(opts.configPath)
	} else {
		cfg, err = manager.Load()
	}
	if err != nil {
		cfg = domain.DefaultConfig()
	}

	return runDirectoryListing(cwd, cfg)
}

func runDirectoryListing(cwd string, cfg domain.Config) error {
	entries, err := os.ReadDir(cwd)
	if err != nil {
		return fmt.Errorf("failed to read directory: %w", err)
	}

	items := make([]list.Item, 0, len(entries))
	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			continue
		}
		items = append(items, fileItem{
			name:    entry.Name(),
			isDir:   entry.IsDir(),
			modTime: info.ModTime(),
		})
	}

	if len(items) == 0 {
		fmt.Println("Directory is empty")
		return nil
	}

	theme := ui.ThemeFromConfig(cfg)
	delegate := ui.NewListDelegate(theme, ui.ListDelegateOptions{
		Spacing: cfg.ListSpacing,
	})

	listModel := ui.NewListModel(items, delegate, 80, 20, theme)
	listModel.Title = fmt.Sprintf("Directory: %s", cwd)
	listModel.SetShowStatusBar(true)
	listModel.SetFilteringEnabled(true)

	model := directoryListModel{
		list:       listModel,
		theme:      theme,
		responsive: ui.NewResponsiveManager(80),
		cwd:        cwd,
	}

	// Set initial keybindings based on initial screen size
	model.list.AdditionalShortHelpKeys = model.getShortHelpKeys
	model.list.AdditionalFullHelpKeys = model.allHelpKeys

	p := tea.NewProgram(model, tea.WithoutSignalHandler())
	if _, err := p.Run(); err != nil {
		return fmt.Errorf("failed to run interactive list: %w", err)
	}

	return nil
}

type fileItem struct {
	name    string
	isDir   bool
	modTime time.Time
}

func (f fileItem) Title() string {
	if f.isDir {
		return f.name + "/"
	}
	return f.name
}

func (f fileItem) Description() string {
	return utils.TimeAgo(f.modTime)
}

func (f fileItem) FilterValue() string {
	return f.name
}

type directoryListModel struct {
	list          list.Model
	theme         ui.Theme
	responsive    *ui.ResponsiveManager
	cwd           string
	selected      string
	message       string
	confirmMode   bool
	confirmModel  *ui.ConfirmationModel
	pendingAction string
	pendingItem   fileItem
}

// allHelpKeys returns the complete list of keybindings in priority order
func (m directoryListModel) allHelpKeys() []key.Binding {
	return []key.Binding{
		key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "view file/directory"),
		),
		key.NewBinding(
			key.WithKeys("a"),
			key.WithHelp("a", "add new file"),
		),
		key.NewBinding(
			key.WithKeys("d"),
			key.WithHelp("d", "delete file"),
		),
		key.NewBinding(
			key.WithKeys("r"),
			key.WithHelp("r", "rename file"),
		),
		key.NewBinding(
			key.WithKeys("o"),
			key.WithHelp("o", "open in editor"),
		),
	}
}

// getShortHelpKeys returns keybindings for short help based on screen size
func (m directoryListModel) getShortHelpKeys() []key.Binding {
	allKeys := m.allHelpKeys()

	var splitAt int
	switch m.responsive.Breakpoint() {
	case ui.BreakpointXL:
		splitAt = 3
	case ui.BreakpointLG:
		splitAt = 1
	default:
		splitAt = 1
	}

	return allKeys[:splitAt]
}

// getFullHelpKeys returns keybindings for full help
func (m directoryListModel) getFullHelpKeys() []key.Binding {
	allKeys := m.allHelpKeys()

	// Determine split point based on breakpoint
	var splitAt int
	switch m.responsive.Breakpoint() {
	case ui.BreakpointXS:
		splitAt = 1 // Remaining keys after first
	case ui.BreakpointSM:
		splitAt = 2 // Remaining keys after first 2
	case ui.BreakpointMD:
		splitAt = 3 // Remaining keys after first 3
	default:
		// LG and XL: all shown in short, return empty for full
		return []key.Binding{}
	}

	return allKeys[splitAt:]
}

func (m directoryListModel) Init() tea.Cmd {
	return nil
}

func (m directoryListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Handle confirmation dialog if active
	if m.confirmMode && m.confirmModel != nil {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			updated, cmd := m.confirmModel.Update(msg)
			if updatedConfirm, ok := updated.(ui.ConfirmationModel); ok {
				m.confirmModel = &updatedConfirm
				if cmd != nil {
					if _, isQuit := cmd().(tea.QuitMsg); isQuit {
						confirmed := m.confirmModel.ChoiceValue()
						m.confirmMode = false
						if confirmed {
							return m.executeAction()
						} else {
							m.message = fmt.Sprintf("%s cancelled", m.pendingAction)
							m.pendingAction = ""
							return m, nil
						}
					}
				}
			}
			return m, cmd
		}
		return m, nil
	}

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		// Update responsive manager with new width
		m.responsive.SetWidth(msg.Width)

		// Get responsive list dimensions
		width, height := m.responsive.GetListDimensions(msg.Width, msg.Height)
		m.list.SetSize(width, height)

		// Update keybindings based on new screen size
		m.list.AdditionalShortHelpKeys = m.getShortHelpKeys
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			return m, tea.Quit
		case "enter":
			// View file
			if item, ok := m.list.SelectedItem().(fileItem); ok {
				m.selected = item.name
				if !item.isDir {
					m.message = fmt.Sprintf("Viewing: %s", item.name)
					// TODO: Implement file viewing
				} else {
					m.message = fmt.Sprintf("Directory: %s", item.name)
				}
			}
		case "d":
			// Delete file - show confirmation
			if item, ok := m.list.SelectedItem().(fileItem); ok {
				m.pendingAction = "Delete"
				m.pendingItem = item
				confirmModel := ui.NewConfirmationModel(
					"Delete File",
					fmt.Sprintf("Are you sure you want to delete '%s'?\nThis action cannot be undone.", item.name),
					m.theme,
				)
				m.confirmModel = &confirmModel
				m.confirmMode = true
				return m, confirmModel.Init()
			}
		case "r":
			// Rename file - show confirmation
			if item, ok := m.list.SelectedItem().(fileItem); ok {
				m.pendingAction = "Rename"
				m.pendingItem = item
				confirmModel := ui.NewConfirmationModel(
					"Rename File",
					fmt.Sprintf("Rename '%s'?", item.name),
					m.theme,
				)
				m.confirmModel = &confirmModel
				m.confirmMode = true
				return m, confirmModel.Init()
			}
		case "o":
			// Open file in editor
			if item, ok := m.list.SelectedItem().(fileItem); ok {
				m.message = fmt.Sprintf("Open: %s (not implemented)", item.name)
				// TODO: Implement file opening
			}
		case "a":
			// Add new file
			m.message = "Add file (not implemented)"
			// TODO: Implement file creation
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m directoryListModel) View() string {
	// Show confirmation dialog if active
	if m.confirmMode && m.confirmModel != nil {
		return m.confirmModel.View()
	}

	listView := m.list.View()

	// Add message if present
	if m.message != "" {
		listView = listView + "\n\n" + m.message
	}

	return m.responsive.AdaptiveFrameStyle(m.theme).Render(listView)
}

func (m directoryListModel) executeAction() (tea.Model, tea.Cmd) {
	switch m.pendingAction {
	case "Delete":
		// TODO: Implement actual file deletion
		m.message = fmt.Sprintf("✓ Deleted: %s (stub - not actually deleted)", m.pendingItem.name)
	case "Rename":
		// TODO: Implement actual file renaming
		m.message = fmt.Sprintf("✓ Renamed: %s (stub - not actually renamed)", m.pendingItem.name)
	}
	m.pendingAction = ""
	return m, nil
}
