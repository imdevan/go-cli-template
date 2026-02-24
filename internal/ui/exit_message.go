package ui

import "github.com/charmbracelet/lipgloss"

// ClipboardConfirm renders the standard clipboard confirmation message.
func ClipboardConfirm(theme Theme) string {
	return ExitMessage(theme, "  Copied to clipboard  󱁖", false)
}

// ExitMessage renders a standard framed exit message.
func ExitMessage(theme Theme, message string, mutedText bool) string {
	style := lipgloss.NewStyle().
		Foreground(lipgloss.Color(theme.Text)).
		Margin(0, 2, 0, 2)

	if mutedText {
		return style.Render(message)
	}

	return style.Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(theme.Muted)).
		Foreground(lipgloss.Color(theme.Secondary)).
		Bold(true).
		Margin(1, 1).
		Padding(1, 2).Render(message)
}
