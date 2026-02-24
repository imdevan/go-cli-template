package testutil

import (
	"bytes"
	"io"
	"testing"
)

// RunCLI executes a Cobra command in-process and returns combined output.
func RunCLI(t *testing.T, cmd commandRunner, args ...string) (string, error) {
	t.Helper()
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetErr(&buf)
	cmd.SetArgs(args)
	err := cmd.Execute()
	return buf.String(), err
}

// commandRunner is the minimal cobra command surface for tests.
type commandRunner interface {
	SetArgs(args []string)
	SetOut(writer io.Writer)
	SetErr(writer io.Writer)
	Execute() error
}
