package registry

import (
	"errors"
	"testing"

	"github.com/toolvault/toolvault/internal/registry/spec"
)

func TestRegisterStatusContract(t *testing.T) {
	tests := []struct {
		name    string
		status  spec.ToolStatus
		allowed bool
	}{
		{name: "draft", status: spec.ToolStatusDraft, allowed: true},
		{name: "active", status: spec.ToolStatusActive, allowed: true},
		{name: "deprecated", status: spec.ToolStatusDeprecated, allowed: false},
		{name: "disabled", status: spec.ToolStatusDisabled, allowed: false},
		{name: "unknown", status: "archived", allowed: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CanRegisterStatus(tt.status); got != tt.allowed {
				t.Fatalf("CanRegisterStatus(%q) = %v, want %v", tt.status, got, tt.allowed)
			}

			err := ValidateRegisterStatus(tt.status)
			if tt.allowed && err != nil {
				t.Fatalf("ValidateRegisterStatus(%q) error = %v, want nil", tt.status, err)
			}
			wantKind := InvalidStateTransition
			if !IsValidStatus(tt.status) {
				wantKind = InvalidSpec
			}
			if !tt.allowed && !errors.Is(err, wantKind) {
				t.Fatalf("ValidateRegisterStatus(%q) error = %v, want %v", tt.status, err, wantKind)
			}
		})
	}
}

func TestStatusTransitionContract(t *testing.T) {
	tests := []struct {
		name    string
		from    spec.ToolStatus
		to      spec.ToolStatus
		allowed bool
	}{
		{name: "same draft", from: spec.ToolStatusDraft, to: spec.ToolStatusDraft, allowed: true},
		{name: "draft to active", from: spec.ToolStatusDraft, to: spec.ToolStatusActive, allowed: true},
		{name: "draft to disabled", from: spec.ToolStatusDraft, to: spec.ToolStatusDisabled, allowed: true},
		{name: "active to deprecated", from: spec.ToolStatusActive, to: spec.ToolStatusDeprecated, allowed: true},
		{name: "active to disabled", from: spec.ToolStatusActive, to: spec.ToolStatusDisabled, allowed: true},
		{name: "deprecated to disabled", from: spec.ToolStatusDeprecated, to: spec.ToolStatusDisabled, allowed: true},
		{name: "disabled to active", from: spec.ToolStatusDisabled, to: spec.ToolStatusActive, allowed: false},
		{name: "deprecated to active", from: spec.ToolStatusDeprecated, to: spec.ToolStatusActive, allowed: false},
		{name: "draft to deprecated", from: spec.ToolStatusDraft, to: spec.ToolStatusDeprecated, allowed: false},
		{name: "unknown from", from: "archived", to: spec.ToolStatusActive, allowed: false},
		{name: "unknown to", from: spec.ToolStatusActive, to: "archived", allowed: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CanTransitionStatus(tt.from, tt.to); got != tt.allowed {
				t.Fatalf("CanTransitionStatus(%q, %q) = %v, want %v", tt.from, tt.to, got, tt.allowed)
			}

			err := ValidateStatusTransition(tt.from, tt.to)
			if tt.allowed && err != nil {
				t.Fatalf("ValidateStatusTransition(%q, %q) error = %v, want nil", tt.from, tt.to, err)
			}
			wantKind := InvalidStateTransition
			if !IsValidStatus(tt.from) || !IsValidStatus(tt.to) {
				wantKind = InvalidSpec
			}
			if !tt.allowed && !errors.Is(err, wantKind) {
				t.Fatalf("ValidateStatusTransition(%q, %q) error = %v, want %v", tt.from, tt.to, err, wantKind)
			}
		})
	}
}
