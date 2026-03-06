package tty

import (
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"golang.org/x/term"
)

// IsTerminal checks if the file descriptor is a terminal
func IsTerminal(fd uintptr) bool {
	return term.IsTerminal(int(fd))
}

// GetProgramOptions returns Bubble Tea program options with TTY redirection
// when stdout is not a TTY (e.g., captured by command substitution).
// This allows the TUI to display properly while keeping stdout for result output.
func GetProgramOptions(baseOpts ...tea.ProgramOption) []tea.ProgramOption {
	opts := make([]tea.ProgramOption, 0, len(baseOpts)+2)
	opts = append(opts, baseOpts...)

	// When stdout is not a TTY (e.g., captured by $()), redirect TUI to /dev/tty
	if !IsTerminal(os.Stdout.Fd()) {
		tty, err := os.OpenFile("/dev/tty", os.O_RDWR, 0)
		if err == nil {
			// Note: caller should defer tty.Close() if they need to manage the lifecycle
			// For now, we let it be closed when the program exits
			opts = append(opts, tea.WithInput(tty), tea.WithOutput(tty))
		}
	}

	return opts
}
