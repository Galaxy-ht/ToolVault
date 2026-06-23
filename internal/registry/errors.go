package registry

import (
	"fmt"

	"github.com/toolvault/toolvault/internal/registry/spec"
)

type ErrorKind string

const (
	AlreadyExists          ErrorKind = "already_exists"
	NotFound               ErrorKind = "not_found"
	InvalidSpec            ErrorKind = "invalid_spec"
	InvalidStateTransition ErrorKind = "invalid_state_transition"
)

func (k ErrorKind) Error() string {
	return string(k)
}

type Error struct {
	Kind   ErrorKind
	Op     string
	ToolID spec.ToolID
	Err    error
}

func NewError(kind ErrorKind, op string, id spec.ToolID, err error) *Error {
	return &Error{
		Kind:   kind,
		Op:     op,
		ToolID: id,
		Err:    err,
	}
}

func (e *Error) Error() string {
	if e == nil {
		return "registry error"
	}

	message := string(e.Kind)
	if e.Op != "" {
		message = e.Op + ": " + message
	}
	if e.ToolID != "" {
		message = fmt.Sprintf("%s: tool %q", message, e.ToolID)
	}
	if e.Err != nil {
		message = message + ": " + e.Err.Error()
	}

	return message
}

func (e *Error) Unwrap() error {
	if e == nil {
		return nil
	}

	return e.Err
}

func (e *Error) Is(target error) bool {
	if e == nil {
		return false
	}

	kind, ok := target.(ErrorKind)
	return ok && e.Kind == kind
}

func KindOf(err error) (ErrorKind, bool) {
	for err != nil {
		if kind, ok := err.(ErrorKind); ok {
			return kind, true
		}
		if registryErr, ok := err.(*Error); ok {
			return registryErr.Kind, true
		}

		unwrap, ok := err.(interface{ Unwrap() error })
		if !ok {
			break
		}
		err = unwrap.Unwrap()
	}

	return "", false
}
