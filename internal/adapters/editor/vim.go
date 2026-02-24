package editor

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

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

func IsVim(command string) bool {
	fields := strings.Fields(command)
	if len(fields) == 0 {
		return false
	}
	base := strings.ToLower(filepath.Base(fields[0]))
	return strings.Contains(base, "nvim") || strings.Contains(base, "vim") || base == "vi"
}

func OpenVimInsert(command, path string, line int) error {
	fields := strings.Fields(command)
	if len(fields) == 0 {
		return errors.New("editor command is required")
	}
	args := append(fields[1:], fmt.Sprintf("+call cursor(%d,1)", line), "+startinsert", path)
	return runEditorCommand(fields[0], args)
}

func OpenVimAtEnd(command, path string) error {
	fields := strings.Fields(command)
	if len(fields) == 0 {
		return errors.New("editor command is required")
	}
	args := append(fields[1:], "+normal G$", path)
	return runEditorCommand(fields[0], args)
}

func runEditorCommand(command string, args []string) error {
	cmd := exec.Command(command, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd.Run()
}
