package workflow

import "errors"

var (
	// ErrNotImplemented is returned when a workflow is not yet implemented
	ErrNotImplemented = errors.New("workflow not implemented")
)
