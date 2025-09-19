package types

import "errors"

var (
	ErrInputNotFound     = errors.New("input not found")
	ErrOutputNotFound    = errors.New("output not found")
	ErrInvalidConnection = errors.New("invalid connection")
	ErrCyclicGraph       = errors.New("cyclic dependency detected")
	ErrNodeNotFound      = errors.New("node not found")
	ErrInvalidNodeType   = errors.New("invalid node type")
	ErrExecutionFailed   = errors.New("node execution failed")
)