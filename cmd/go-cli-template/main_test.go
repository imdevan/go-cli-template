package main

import (
	"testing"
)

func TestCLIRuns(t *testing.T) {
	cmd := newRootCmd()
	if cmd == nil {
		t.Fatal("expected root command to be created")
	}
}
