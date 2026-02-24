package ui

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
)

type staticHelp struct {
	short []key.Binding
	full  [][]key.Binding
}

func (h staticHelp) ShortHelp() []key.Binding {
	return h.short
}

func (h staticHelp) FullHelp() [][]key.Binding {
	return h.full
}

func ListHelpView(model list.Model, short []key.Binding, full [][]key.Binding) string {
	return model.Styles.HelpStyle.Render(model.Help.View(staticHelp{
		short: short,
		full:  full,
	}))
}
