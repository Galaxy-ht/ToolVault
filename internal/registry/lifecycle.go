package registry

import (
	"fmt"

	"github.com/toolvault/toolvault/internal/registry/spec"
)

func IsValidStatus(status spec.ToolStatus) bool {
	switch status {
	case spec.ToolStatusDraft, spec.ToolStatusActive, spec.ToolStatusDeprecated, spec.ToolStatusDisabled:
		return true
	default:
		return false
	}
}

func CanRegisterStatus(status spec.ToolStatus) bool {
	switch status {
	case spec.ToolStatusDraft, spec.ToolStatusActive:
		return true
	default:
		return false
	}
}

func ValidateRegisterStatus(status spec.ToolStatus) error {
	if !IsValidStatus(status) {
		return NewError(
			InvalidSpec,
			"registry.register",
			"",
			fmt.Errorf("status %q is not valid", status),
		)
	}
	if CanRegisterStatus(status) {
		return nil
	}

	return NewError(
		InvalidStateTransition,
		"registry.register",
		"",
		fmt.Errorf("initial status %q is not allowed", status),
	)
}

func CanTransitionStatus(from, to spec.ToolStatus) bool {
	if !IsValidStatus(from) || !IsValidStatus(to) {
		return false
	}
	if from == to {
		return true
	}

	switch from {
	case spec.ToolStatusDraft:
		return to == spec.ToolStatusActive || to == spec.ToolStatusDisabled
	case spec.ToolStatusActive:
		return to == spec.ToolStatusDeprecated || to == spec.ToolStatusDisabled
	case spec.ToolStatusDeprecated:
		return to == spec.ToolStatusDisabled
	default:
		return false
	}
}

func ValidateStatusTransition(from, to spec.ToolStatus) error {
	if !IsValidStatus(from) || !IsValidStatus(to) {
		return NewError(
			InvalidSpec,
			"registry.set_status",
			"",
			fmt.Errorf("status transition %q -> %q contains an invalid status", from, to),
		)
	}
	if CanTransitionStatus(from, to) {
		return nil
	}

	return NewError(
		InvalidStateTransition,
		"registry.set_status",
		"",
		fmt.Errorf("status transition %q -> %q is not allowed", from, to),
	)
}
