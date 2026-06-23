package memory

import (
	"context"
	"sort"
	"sync"

	"github.com/toolvault/toolvault/internal/registry"
	"github.com/toolvault/toolvault/internal/registry/spec"
)

const (
	opRegister  = "registry.memory.register"
	opGet       = "registry.memory.get"
	opList      = "registry.memory.list"
	opUpdate    = "registry.memory.update"
	opDelete    = "registry.memory.delete"
	opSetStatus = "registry.memory.set_status"
)

type Registry struct {
	mu    sync.RWMutex
	tools map[spec.ToolID]spec.ToolSpec
}

func New() *Registry {
	return &Registry{
		tools: make(map[spec.ToolID]spec.ToolSpec),
	}
}

func (r *Registry) Register(_ context.Context, req registry.RegisterRequest) (spec.ToolSpec, error) {
	tool := cloneToolSpec(req.Spec)
	if err := spec.Validate(tool); err != nil {
		return spec.ToolSpec{}, registry.NewError(registry.InvalidSpec, opRegister, tool.ID, err)
	}
	if err := registry.ValidateRegisterStatus(tool.Status); err != nil {
		return spec.ToolSpec{}, wrapRegistryError(err, opRegister, tool.ID)
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if r.tools == nil {
		r.tools = make(map[spec.ToolID]spec.ToolSpec)
	}
	if _, ok := r.tools[tool.ID]; ok {
		return spec.ToolSpec{}, registry.NewError(registry.AlreadyExists, opRegister, tool.ID, nil)
	}

	r.tools[tool.ID] = cloneToolSpec(tool)
	return cloneToolSpec(tool), nil
}

func (r *Registry) Get(_ context.Context, req registry.GetRequest) (spec.ToolSpec, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	tool, ok := r.tools[req.ID]
	if !ok {
		return spec.ToolSpec{}, registry.NewError(registry.NotFound, opGet, req.ID, nil)
	}

	return cloneToolSpec(tool), nil
}

func (r *Registry) List(_ context.Context, req registry.ListRequest) ([]spec.ToolSpec, error) {
	filter := req.Filter.Clone()
	if filter.Status != "" && !registry.IsValidStatus(filter.Status) {
		return nil, registry.NewError(registry.InvalidSpec, opList, "", registry.InvalidSpec)
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	tools := make([]spec.ToolSpec, 0, len(r.tools))
	for _, tool := range r.tools {
		if !matchesFilter(tool, filter) {
			continue
		}
		tools = append(tools, cloneToolSpec(tool))
	}

	sort.Slice(tools, func(i, j int) bool {
		return tools[i].ID < tools[j].ID
	})

	return tools, nil
}

func (r *Registry) Update(_ context.Context, req registry.UpdateRequest) (spec.ToolSpec, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	current, ok := r.tools[req.ID]
	if !ok {
		return spec.ToolSpec{}, registry.NewError(registry.NotFound, opUpdate, req.ID, nil)
	}

	updated := cloneToolSpec(current)
	updated.DisplayName = req.Update.DisplayName
	updated.Description = req.Update.Description
	updated.Version = req.Update.Version
	updated.Tags = cloneStrings(req.Update.Tags)
	updated.Actions = cloneActions(req.Update.Actions)
	updated.Metadata = cloneMetadata(req.Update.Metadata)

	if err := spec.Validate(updated); err != nil {
		return spec.ToolSpec{}, registry.NewError(registry.InvalidSpec, opUpdate, req.ID, err)
	}

	r.tools[req.ID] = cloneToolSpec(updated)
	return cloneToolSpec(updated), nil
}

func (r *Registry) Delete(_ context.Context, req registry.DeleteRequest) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.tools[req.ID]; !ok {
		return registry.NewError(registry.NotFound, opDelete, req.ID, nil)
	}

	delete(r.tools, req.ID)
	return nil
}

func (r *Registry) SetStatus(_ context.Context, req registry.SetStatusRequest) (spec.ToolSpec, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	current, ok := r.tools[req.ID]
	if !ok {
		return spec.ToolSpec{}, registry.NewError(registry.NotFound, opSetStatus, req.ID, nil)
	}
	if err := registry.ValidateStatusTransition(current.Status, req.Status); err != nil {
		return spec.ToolSpec{}, wrapRegistryError(err, opSetStatus, req.ID)
	}

	updated := cloneToolSpec(current)
	updated.Status = req.Status
	if err := spec.Validate(updated); err != nil {
		return spec.ToolSpec{}, registry.NewError(registry.InvalidSpec, opSetStatus, req.ID, err)
	}

	r.tools[req.ID] = cloneToolSpec(updated)
	return cloneToolSpec(updated), nil
}

func wrapRegistryError(err error, op string, id spec.ToolID) error {
	kind, ok := registry.KindOf(err)
	if !ok {
		return err
	}

	return registry.NewError(kind, op, id, err)
}

func matchesFilter(tool spec.ToolSpec, filter registry.ListFilter) bool {
	if filter.Status != "" && tool.Status != filter.Status {
		return false
	}
	for _, tag := range filter.Tags {
		if !hasTag(tool.Tags, tag) {
			return false
		}
	}

	return true
}

func hasTag(tags []string, want string) bool {
	for _, tag := range tags {
		if tag == want {
			return true
		}
	}

	return false
}

func cloneToolSpec(tool spec.ToolSpec) spec.ToolSpec {
	clone := tool
	clone.Tags = cloneStrings(tool.Tags)
	clone.Actions = cloneActions(tool.Actions)
	clone.Metadata = cloneMetadata(tool.Metadata)

	return clone
}

func cloneStrings(values []string) []string {
	if len(values) == 0 {
		return nil
	}

	return append([]string(nil), values...)
}

func cloneActions(actions []spec.ToolAction) []spec.ToolAction {
	if len(actions) == 0 {
		return nil
	}

	return append([]spec.ToolAction(nil), actions...)
}

func cloneMetadata(metadata map[string]string) map[string]string {
	if len(metadata) == 0 {
		return nil
	}

	clone := make(map[string]string, len(metadata))
	for key, value := range metadata {
		clone[key] = value
	}

	return clone
}
