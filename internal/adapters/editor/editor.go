package editor

import "errors"

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
