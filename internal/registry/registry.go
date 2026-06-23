package registry

import (
	"context"

	"github.com/toolvault/toolvault/internal/registry/spec"
)

// Registry defines the M1 contract for Tool Registry implementations.
// It intentionally contains no storage, protocol, or runtime assumptions.
type Registry interface {
	Register(context.Context, RegisterRequest) (spec.ToolSpec, error)
	Get(context.Context, GetRequest) (spec.ToolSpec, error)
	// List returns all non-deleted tools by default, including disabled and
	// deprecated tools. Implementations should return stable results.
	List(context.Context, ListRequest) ([]spec.ToolSpec, error)
	// Update replaces M1 mutable metadata fields after an expected-version
	// check. Successful updates must store a new opaque ToolSpec version.
	Update(context.Context, UpdateRequest) (spec.ToolSpec, error)
	// Delete is an M1 hard delete. Lifecycle changes belong in SetStatus.
	Delete(context.Context, DeleteRequest) error
	// SetStatus applies Registry-owned lifecycle transitions and preserves the
	// ToolSpec version because status is not the tool definition version.
	SetStatus(context.Context, SetStatusRequest) (spec.ToolSpec, error)
}

type RegisterRequest struct {
	Spec    spec.ToolSpec
	Options RegisterOptions
}

type RegisterOptions struct{}

func DefaultRegisterOptions() RegisterOptions {
	return RegisterOptions{}
}

type GetRequest struct {
	ID      spec.ToolID
	Options GetOptions
}

type GetOptions struct{}

func DefaultGetOptions() GetOptions {
	return GetOptions{}
}

type ListRequest struct {
	Filter  ListFilter
	Options ListOptions
}

type ListFilter struct {
	Status spec.ToolStatus
	Tags   []string
}

func DefaultListFilter() ListFilter {
	return ListFilter{}
}

func (f ListFilter) IsZero() bool {
	return f.Status == "" && len(f.Tags) == 0
}

func (f ListFilter) Clone() ListFilter {
	clone := ListFilter{
		Status: f.Status,
	}
	if len(f.Tags) > 0 {
		clone.Tags = append([]string(nil), f.Tags...)
	}

	return clone
}

type ListOptions struct{}

func DefaultListOptions() ListOptions {
	return ListOptions{}
}

type UpdateRequest struct {
	ID      spec.ToolID
	Update  ToolSpecUpdate
	Options UpdateOptions
}

// ToolSpecUpdate contains mutable ToolSpec fields for M1.
// ID, Name, and Status are intentionally excluded: ID and Name are stable
// identity fields, and lifecycle changes must go through SetStatus.
type ToolSpecUpdate struct {
	DisplayName string
	Description string
	// Version is the caller-provided next opaque tool definition version.
	// It must be non-empty and differ from the current stored version.
	Version  string
	Tags     []string
	Actions  []spec.ToolAction
	Metadata map[string]string
}

type UpdateOptions struct {
	// ExpectedVersion is required for M1 optimistic update checks.
	// It must match the currently stored opaque ToolSpec version.
	ExpectedVersion string
}

func DefaultUpdateOptions() UpdateOptions {
	return UpdateOptions{}
}

type DeleteRequest struct {
	ID      spec.ToolID
	Options DeleteOptions
}

type DeleteOptions struct{}

func DefaultDeleteOptions() DeleteOptions {
	return DeleteOptions{}
}

type SetStatusRequest struct {
	ID      spec.ToolID
	Status  spec.ToolStatus
	Options SetStatusOptions
}

type SetStatusOptions struct{}

func DefaultSetStatusOptions() SetStatusOptions {
	return SetStatusOptions{}
}
