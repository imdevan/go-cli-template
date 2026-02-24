package ui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

// ConfirmationModel is a reusable confirmation dialog component based on huh.Confirm.
// It provides a styled yes/no dialog with keyboard shortcuts.
//
// Example usage as a standalone dialog:
//
//	confirmed, err := ui.PromptConfirmation("Delete file?", "Are you sure?", theme)
//	if err != nil {
//	    return err
//	}
//	if confirmed {
//	    // perform action
//	}
//
// Example usage embedded in another Bubble Tea model:
//
//	// In your model struct:
//	type myModel struct {
//	    confirmMode bool
//	    confirm     *ui.ConfirmationModel
//	}
//
//	// To show the dialog:
//	confirmModel := ui.NewConfirmationModel("Delete?", "Are you sure?", theme)
//	m.confirmMode = true
//	m.confirm = &confirmModel
//	return m, confirmModel.Init()
//
//	// In Update method:
//	if m.confirmMode {
//	    updated, cmd := m.confirm.Update(msg)
//	    if updatedConfirm, ok := updated.(ui.ConfirmationModel); ok {
//	        m.confirm = &updatedConfirm
//	        if cmd != nil {
//	            if _, isQuit := cmd().(tea.QuitMsg); isQuit {
//	                confirmed := m.confirm.ChoiceValue()
//	                m.confirmMode = false
//	                // handle result
//	            }
//	        }
//	    }
//	    return m, cmd
//	}
//
//	// In View method:
//	if m.confirmMode && m.confirm != nil {
//	    return m.confirm.View()
//	}
type ConfirmationModel struct {
	title   string
	prompt  string
	choice  *bool
	confirm *huh.Confirm
	theme   Theme
}

// NewConfirmationModel creates a new confirmation dialog
func NewConfirmationModel(title, prompt string, theme Theme) ConfirmationModel {
	choice := false
	helpText := "Press y or n."
	width := confirmationDialogWidth(title, prompt, helpText)

	confirm := huh.NewConfirm().
		Title("").
		Description("").
		Value(&choice)
	confirm.WithKeyMap(huh.NewDefaultKeyMap())
	confirm.WithWidth(width)
	confirm.WithButtonAlignment(lipgloss.Center)
	confirm.WithTheme(confirmationHuhTheme(theme))
	confirm.Focus()

	return ConfirmationModel{
		title:   title,
		prompt:  prompt,
		choice:  &choice,
		confirm: confirm,
		theme:   theme,
	}
}

// Init initializes the confirmation dialog
func (m ConfirmationModel) Init() tea.Cmd {
	if m.confirm != nil {
		return m.confirm.Init()
	}
	return nil
}

// Update handles messages for the confirmation dialog
func (m ConfirmationModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			m.setChoice(false)
			return m, tea.Quit
		case tea.KeyEnter:
			return m, tea.Quit
		}
		switch msg.String() {
		case "y", "Y":
			m.setChoice(true)
			return m, tea.Quit
		case "n", "N":
			m.setChoice(false)
			return m, tea.Quit
		}
	}
	if m.confirm == nil {
		return m, nil
	}
	updated, cmd := m.confirm.Update(msg)
	if confirm, ok := updated.(*huh.Confirm); ok {
		m.confirm = confirm
	}
	return m, cmd
}

// View renders the confirmation dialog
func (m ConfirmationModel) View() string {
	titleText := m.title
	if strings.TrimSpace(titleText) == "" {
		titleText = "Confirm"
	}
	helpText := "Press y or n."

	title := lipgloss.NewStyle().Bold(true).Foreground(m.theme.Headings).Render(titleText)
	prompt := lipgloss.NewStyle().Foreground(m.theme.Text).Render(m.prompt)
	help := lipgloss.NewStyle().Foreground(m.theme.Muted).Render(helpText)

	confirmView := ""
	if m.confirm != nil {
		confirmView = m.confirm.View()
	}

	content := strings.Join([]string{
		title,
		"",
		confirmView,
		"",
		prompt,
		help,
	}, "\n")

	return lipgloss.NewStyle().
		Margin(1, 1).
		Padding(1, 2).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(m.theme.Border).
		Render(content)
}

// ChoiceValue returns the user's choice
func (m ConfirmationModel) ChoiceValue() bool {
	if m.choice == nil {
		return false
	}
	return *m.choice
}

func (m ConfirmationModel) setChoice(value bool) {
	if m.choice == nil {
		return
	}
	*m.choice = value
}

// PromptConfirmation runs a confirmation dialog and returns the user's choice
func PromptConfirmation(title, prompt string, theme Theme) (bool, error) {
	model := NewConfirmationModel(title, prompt, theme)
	program := tea.NewProgram(model, tea.WithoutSignalHandler())
	result, err := program.Run()
	if err != nil {
		return false, err
	}
	if m, ok := result.(ConfirmationModel); ok {
		return m.ChoiceValue(), nil
	}
	return false, fmt.Errorf("unexpected model result")
}

func confirmationDialogWidth(title, prompt, help string) int {
	width := lipgloss.Width(title)
	if promptWidth := lipgloss.Width(prompt); promptWidth > width {
		width = promptWidth
	}
	if helpWidth := lipgloss.Width(help); helpWidth > width {
		width = helpWidth
	}
	const minWidth = 32
	if width < minWidth {
		width = minWidth
	}
	return width
}

func confirmationHuhTheme(theme Theme) *huh.Theme {
	huhTheme := huh.ThemeBase()
	huhTheme.Focused.Base = lipgloss.NewStyle()
	huhTheme.Blurred.Base = lipgloss.NewStyle()
	huhTheme.Focused.FocusedButton = lipgloss.NewStyle().
		Padding(0, 2).
		MarginRight(1).
		Foreground(theme.Text).
		Background(theme.Border).
		Bold(true)
	huhTheme.Focused.BlurredButton = lipgloss.NewStyle().
		Padding(0, 2).
		MarginRight(1).
		Foreground(theme.Text)
	huhTheme.Blurred.FocusedButton = huhTheme.Focused.FocusedButton
	huhTheme.Blurred.BlurredButton = huhTheme.Focused.BlurredButton
	return huhTheme
}
