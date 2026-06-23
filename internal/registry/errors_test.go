package registry

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/toolvault/toolvault/internal/registry/spec"
)

func TestErrorKindIsRecognizable(t *testing.T) {
	err := NewError(AlreadyExists, "registry.register", "tool-search", nil)

	if !errors.Is(err, AlreadyExists) {
		t.Fatalf("errors.Is(%v, %v) = false, want true", err, AlreadyExists)
	}

	kind, ok := KindOf(err)
	if !ok {
		t.Fatalf("KindOf(%v) ok = false, want true", err)
	}
	if kind != AlreadyExists {
		t.Fatalf("KindOf(%v) = %q, want %q", err, kind, AlreadyExists)
	}
}

func TestErrorWrappingAndMatching(t *testing.T) {
	cause := errors.New("validation failed")
	err := fmt.Errorf("outer: %w", NewError(InvalidSpec, "registry.register", "tool-search", cause))

	if !errors.Is(err, InvalidSpec) {
		t.Fatalf("errors.Is(%v, %v) = false, want true", err, InvalidSpec)
	}
	if !errors.Is(err, cause) {
		t.Fatalf("errors.Is(%v, cause) = false, want true", err)
	}

	var registryErr *Error
	if !errors.As(err, &registryErr) {
		t.Fatalf("errors.As(%v, *Error) = false, want true", err)
	}
	if registryErr.ToolID != spec.ToolID("tool-search") {
		t.Fatalf("registryErr.ToolID = %q, want %q", registryErr.ToolID, spec.ToolID("tool-search"))
	}
}

func TestErrorKindCanBeWrappedDirectly(t *testing.T) {
	err := fmt.Errorf("outer: %w", NotFound)

	if !errors.Is(err, NotFound) {
		t.Fatalf("errors.Is(%v, %v) = false, want true", err, NotFound)
	}

	kind, ok := KindOf(err)
	if !ok {
		t.Fatalf("KindOf(%v) ok = false, want true", err)
	}
	if kind != NotFound {
		t.Fatalf("KindOf(%v) = %q, want %q", err, kind, NotFound)
	}
}

func TestErrorMessageIncludesContractContext(t *testing.T) {
	err := NewError(InvalidStateTransition, "registry.set_status", "tool-search", errors.New("disabled to active"))
	message := err.Error()

	for _, want := range []string{"registry.set_status", string(InvalidStateTransition), "tool-search", "disabled to active"} {
		if !strings.Contains(message, want) {
			t.Fatalf("Error() = %q, want substring %q", message, want)
		}
	}
}
