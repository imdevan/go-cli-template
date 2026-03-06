package editor

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Adapter launches the configured editor.
type Adapter struct {
	Command string
}

// New returns an editor adapter using the given command.
func New(command string) *Adapter {
	return &Adapter{Command: command}
}

// Open launches the editor with the provided file path.
func (a Adapter) Open(path string) error {
	command := ResolveCommand(a.Command)
	if command == "" {
		return errors.New("editor command is required")
	}
	return runEditorCommand(command, []string{path})
}

// OpenAtLine opens a file at a specific line number.
func (a Adapter) OpenAtLine(path string, line int) error {
	command := ResolveCommand(a.Command)
	if command == "" {
		return errors.New("editor command is required")
	}

	if IsVim(command) {
		return OpenVimAtLine(command, path, line)
	}
	if IsNano(command) {
		return OpenNanoAtLine(command, path, line)
	}
	if IsVSCode(command) {
		return OpenVSCodeAtLine(command, path, line)
	}
	if IsEmacs(command) {
		return OpenEmacsAtLine(command, path, line)
	}

	// For unknown editors, just open normally
	return a.Open(path)
}

// OpenAtEnd opens a file and positions the cursor at the end when supported.
func (a Adapter) OpenAtEnd(path string) error {
	command := ResolveCommand(a.Command)
	if command == "" {
		return errors.New("editor command is required")
	}
	if IsVim(command) {
		return OpenVimAtEnd(command, path)
	}
	return a.Open(path)
}

// ResolveCommand resolves the editor command from config or environment.
func ResolveCommand(command string) string {
	resolved := strings.TrimSpace(command)
	if resolved == "" {
		resolved = strings.TrimSpace(os.Getenv("VISUAL"))
	}
	if resolved == "" {
		resolved = strings.TrimSpace(os.Getenv("EDITOR"))
	}
	return resolved
}

// runEditorCommand executes an editor command with the given arguments.
func runEditorCommand(command string, args []string) error {
	cmd := exec.Command(command, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd.Run()
}

// getEditorBase extracts the base command name from a command string.
func getEditorBase(command string) string {
	fields := strings.Fields(command)
	if len(fields) == 0 {
		return ""
	}
	return strings.ToLower(filepath.Base(fields[0]))
}

// IsVim checks if the command is a vim-based editor.
func IsVim(command string) bool {
	base := getEditorBase(command)
	if base == "" {
		return false
	}
	return strings.Contains(base, "nvim") || strings.Contains(base, "vim") || base == "vi"
}

// IsNano checks if the command is nano.
func IsNano(command string) bool {
	base := getEditorBase(command)
	if base == "" {
		return false
	}
	return strings.Contains(base, "nano")
}

// IsVSCode checks if the command is a VSCode-based editor.
func IsVSCode(command string) bool {
	base := getEditorBase(command)
	if base == "" {
		return false
	}
	return base == "code" || base == "code-insiders" || base == "codium" || base == "cursor"
}

// IsEmacs checks if the command is emacs.
func IsEmacs(command string) bool {
	base := getEditorBase(command)
	if base == "" {
		return false
	}
	return strings.Contains(base, "emacs")
}

// OpenVimInsert opens vim at a specific line in insert mode.
func OpenVimInsert(command, path string, line int) error {
	fields := strings.Fields(command)
	if len(fields) == 0 {
		return errors.New("editor command is required")
	}
	args := append(fields[1:], fmt.Sprintf("+call cursor(%d,1)", line), "+startinsert", path)
	return runEditorCommand(fields[0], args)
}

// OpenVimAtLine opens vim at a specific line.
func OpenVimAtLine(command, path string, line int) error {
	fields := strings.Fields(command)
	if len(fields) == 0 {
		return errors.New("editor command is required")
	}
	args := append(fields[1:], fmt.Sprintf("+%d", line), path)
	return runEditorCommand(fields[0], args)
}

// OpenVimAtEnd opens vim at the end of the file.
func OpenVimAtEnd(command, path string) error {
	fields := strings.Fields(command)
	if len(fields) == 0 {
		return errors.New("editor command is required")
	}
	args := append(fields[1:], "+normal G$", path)
	return runEditorCommand(fields[0], args)
}

// OpenNanoAtLine opens nano at a specific line.
func OpenNanoAtLine(command, path string, line int) error {
	fields := strings.Fields(command)
	if len(fields) == 0 {
		return errors.New("editor command is required")
	}
	args := append(fields[1:], fmt.Sprintf("+%d", line), path)
	return runEditorCommand(fields[0], args)
}

// OpenVSCodeAtLine opens VSCode at a specific line.
func OpenVSCodeAtLine(command, path string, line int) error {
	fields := strings.Fields(command)
	if len(fields) == 0 {
		return errors.New("editor command is required")
	}
	args := append(fields[1:], "-g", fmt.Sprintf("%s:%d", path, line))
	return runEditorCommand(fields[0], args)
}

// OpenEmacsAtLine opens emacs at a specific line.
func OpenEmacsAtLine(command, path string, line int) error {
	fields := strings.Fields(command)
	if len(fields) == 0 {
		return errors.New("editor command is required")
	}
	args := append(fields[1:], fmt.Sprintf("+%d", line), path)
	return runEditorCommand(fields[0], args)
}
