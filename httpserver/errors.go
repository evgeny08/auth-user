package httpserver

import (
	"fmt"
)

// Error is a service error.
type Error struct {
	Kind    ErrorKind
	Message string
}

func (e *Error) Error() string {
	return e.Message
}

// ErrorKind is a kind of the service error.
type ErrorKind uint8

// Error kinds.
const (
	ErrBadParams ErrorKind = iota
	ErrNotFound
	ErrConflict
	ErrInternal
)

func errorf(kind ErrorKind, format string, v ...interface{}) error {
	return &Error{
		Kind:    kind,
		Message: fmt.Sprintf(format, v...),
	}
}
