package registry

import (
	"context"
	"reflect"
	"testing"

	"github.com/toolvault/toolvault/internal/registry/spec"
)

var _ Registry = (*contractRegistry)(nil)

type contractRegistry struct{}

func (contractRegistry) Register(context.Context, RegisterRequest) (spec.ToolSpec, error) {
	return spec.ToolSpec{}, nil
}

func (contractRegistry) Get(context.Context, GetRequest) (spec.ToolSpec, error) {
	return spec.ToolSpec{}, nil
}

func (contractRegistry) List(context.Context, ListRequest) ([]spec.ToolSpec, error) {
	return nil, nil
}

func (contractRegistry) Update(context.Context, UpdateRequest) (spec.ToolSpec, error) {
	return spec.ToolSpec{}, nil
}

func (contractRegistry) Delete(context.Context, DeleteRequest) error {
	return nil
}

func (contractRegistry) SetStatus(context.Context, SetStatusRequest) (spec.ToolSpec, error) {
	return spec.ToolSpec{}, nil
}

func TestDefaultFilterAndOptions(t *testing.T) {
	if got := DefaultListFilter(); !got.IsZero() {
		t.Fatalf("DefaultListFilter().IsZero() = false, want true")
	}

	if got := DefaultRegisterOptions(); got != (RegisterOptions{}) {
		t.Fatalf("DefaultRegisterOptions() = %#v, want zero value", got)
	}
	if got := DefaultGetOptions(); got != (GetOptions{}) {
		t.Fatalf("DefaultGetOptions() = %#v, want zero value", got)
	}
	if got := DefaultListOptions(); got != (ListOptions{}) {
		t.Fatalf("DefaultListOptions() = %#v, want zero value", got)
	}
	if got := DefaultUpdateOptions(); got != (UpdateOptions{}) {
		t.Fatalf("DefaultUpdateOptions() = %#v, want zero value", got)
	}
	if got := DefaultDeleteOptions(); got != (DeleteOptions{}) {
		t.Fatalf("DefaultDeleteOptions() = %#v, want zero value", got)
	}
	if got := DefaultSetStatusOptions(); got != (SetStatusOptions{}) {
		t.Fatalf("DefaultSetStatusOptions() = %#v, want zero value", got)
	}
}

func TestListFilterCloneDoesNotAliasTags(t *testing.T) {
	filter := ListFilter{
		Status: spec.ToolStatusActive,
		Tags:   []string{"search", "internal"},
	}

	clone := filter.Clone()
	if !reflect.DeepEqual(clone, filter) {
		t.Fatalf("Clone() = %#v, want %#v", clone, filter)
	}

	clone.Tags[0] = "changed"
	if filter.Tags[0] == "changed" {
		t.Fatalf("Clone() aliased Tags slice")
	}
}

func TestRequestsUseSpecToolSpec(t *testing.T) {
	var registered spec.ToolSpec = RegisterRequest{}.Spec
	var actions []spec.ToolAction = ToolSpecUpdate{}.Actions
	var metadata map[string]string = ToolSpecUpdate{}.Metadata

	registered.ID = "tool-search"
	actions = append(actions, spec.ToolAction{Name: "search"})
	metadata = map[string]string{"owner": "platform"}

	if registered.ID == "" || actions[0].Name == "" || metadata["owner"] == "" {
		t.Fatalf("registry requests do not use spec package types as expected")
	}
}

func TestUpdateRequestExcludesIdentityAndStatusFields(t *testing.T) {
	updateType := reflect.TypeOf(ToolSpecUpdate{})
	for _, field := range []string{"ID", "Name", "Status"} {
		if _, ok := updateType.FieldByName(field); ok {
			t.Fatalf("ToolSpecUpdate contains %s, want identity and status fields excluded", field)
		}
	}
}
