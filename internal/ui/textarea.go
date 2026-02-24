package ui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/lipgloss"
)

// ConfigureTextarea applies shared styles and keybindings for multi-line input.
func ConfigureTextarea(input *textarea.Model, theme Theme, altEnterSubmit bool) (key.Binding, key.Binding) {
	if input == nil {
		return key.Binding{}, key.Binding{}
	}
	input.ShowLineNumbers = false
	input.FocusedStyle.Base = lipgloss.NewStyle().Foreground(theme.Text)
	input.BlurredStyle.Base = lipgloss.NewStyle().Foreground(theme.Text)
	submitKey, newlineKey := TextareaKeyBindings(altEnterSubmit)
	input.KeyMap.InsertNewline = newlineKey
	return submitKey, newlineKey
}

// TextareaKeyBindings returns submit and newline bindings based on config.
func TextareaKeyBindings(altEnterSubmit bool) (key.Binding, key.Binding) {
	submitKey := key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "submit"))
	newlineKey := key.NewBinding(
		key.WithKeys("alt+enter"),
		key.WithHelp("alt+enter", "newline"),
	)
	if altEnterSubmit {
		submitKey = key.NewBinding(key.WithKeys("alt+enter"), key.WithHelp("alt+enter", "submit"))
		newlineKey = key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "newline"))
	}
	return submitKey, newlineKey
}

// TextareaSubmitHelp formats the key help for a multi-line submit action.
func TextareaSubmitHelp(submitKey key.Binding, action string) string {
	if submitKey.Help().Key == "alt+enter" {
		return fmt.Sprintf("Press Alt+Enter to %s. Enter for a new line.", action)
	}
	return fmt.Sprintf("Press Enter to %s. Alt+Enter for a new line.", action)
}
