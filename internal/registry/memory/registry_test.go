package memory

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"testing"

	"github.com/toolvault/toolvault/internal/registry"
	"github.com/toolvault/toolvault/internal/registry/spec"
)

var _ registry.Registry = (*Registry)(nil)

func TestRegisterSuccess(t *testing.T) {
	store := New()
	input := validToolSpec("tool-search")

	got, err := store.Register(context.Background(), registry.RegisterRequest{Spec: input})
	if err != nil {
		t.Fatalf("Register() error = %v, want nil", err)
	}
	if got.ID != input.ID {
		t.Fatalf("Register() ID = %q, want %q", got.ID, input.ID)
	}

	input.Tags[0] = "changed"
	stored, err := store.Get(context.Background(), registry.GetRequest{ID: got.ID})
	if err != nil {
		t.Fatalf("Get() error = %v, want nil", err)
	}
	if stored.Tags[0] == "changed" {
		t.Fatalf("Register() stored caller-owned Tags slice")
	}
}

func TestRegisterZeroValueRegistry(t *testing.T) {
	var store Registry
	input := validToolSpec("tool-search")

	got, err := store.Register(context.Background(), registry.RegisterRequest{Spec: input})
	if err != nil {
		t.Fatalf("Register() error = %v, want nil", err)
	}
	if got.ID != input.ID {
		t.Fatalf("Register() ID = %q, want %q", got.ID, input.ID)
	}
}

func TestRegisterDuplicate(t *testing.T) {
	store := New()
	tool := validToolSpec("tool-search")
	mustRegister(t, store, tool)

	_, err := store.Register(context.Background(), registry.RegisterRequest{Spec: tool})
	if !errors.Is(err, registry.AlreadyExists) {
		t.Fatalf("Register() error = %v, want %v", err, registry.AlreadyExists)
	}
}

func TestRegisterInvalidSpec(t *testing.T) {
	store := New()
	tool := validToolSpec("tool-search")
	tool.Name = ""

	_, err := store.Register(context.Background(), registry.RegisterRequest{Spec: tool})
	if !errors.Is(err, registry.InvalidSpec) {
		t.Fatalf("Register() error = %v, want %v", err, registry.InvalidSpec)
	}
}

func TestRegisterInvalidInitialStatus(t *testing.T) {
	store := New()
	tool := validToolSpec("tool-search")
	tool.Status = spec.ToolStatusDisabled

	_, err := store.Register(context.Background(), registry.RegisterRequest{Spec: tool})
	if !errors.Is(err, registry.InvalidStateTransition) {
		t.Fatalf("Register() error = %v, want %v", err, registry.InvalidStateTransition)
	}
}

func TestGetSuccess(t *testing.T) {
	store := New()
	tool := mustRegister(t, store, validToolSpec("tool-search"))

	got, err := store.Get(context.Background(), registry.GetRequest{ID: tool.ID})
	if err != nil {
		t.Fatalf("Get() error = %v, want nil", err)
	}
	if got.ID != tool.ID {
		t.Fatalf("Get() ID = %q, want %q", got.ID, tool.ID)
	}

	got.Tags[0] = "changed"
	got.Actions[0].Name = "changed"
	got.Metadata["owner"] = "changed"

	again, err := store.Get(context.Background(), registry.GetRequest{ID: tool.ID})
	if err != nil {
		t.Fatalf("Get() second error = %v, want nil", err)
	}
	if again.Tags[0] == "changed" || again.Actions[0].Name == "changed" || again.Metadata["owner"] == "changed" {
		t.Fatalf("Get() returned aliased mutable fields")
	}
}

func TestGetNotFound(t *testing.T) {
	store := New()

	_, err := store.Get(context.Background(), registry.GetRequest{ID: "missing"})
	if !errors.Is(err, registry.NotFound) {
		t.Fatalf("Get() error = %v, want %v", err, registry.NotFound)
	}
}

func TestListEmpty(t *testing.T) {
	store := New()

	got, err := store.List(context.Background(), registry.ListRequest{})
	if err != nil {
		t.Fatalf("List() error = %v, want nil", err)
	}
	if len(got) != 0 {
		t.Fatalf("List() len = %d, want 0", len(got))
	}
}

func TestListWithMultipleTools(t *testing.T) {
	store := New()
	mustRegister(t, store, validToolSpec("tool-zeta"))
	mustRegister(t, store, validToolSpec("tool-alpha"))

	got, err := store.List(context.Background(), registry.ListRequest{})
	if err != nil {
		t.Fatalf("List() error = %v, want nil", err)
	}
	if len(got) != 2 {
		t.Fatalf("List() len = %d, want 2", len(got))
	}
	if got[0].ID != "tool-alpha" || got[1].ID != "tool-zeta" {
		t.Fatalf("List() IDs = [%q, %q], want sorted IDs", got[0].ID, got[1].ID)
	}
}

func TestListFiltersByStatusAndTags(t *testing.T) {
	store := New()
	first := validToolSpec("tool-first")
	first.Tags = []string{"search", "internal"}
	mustRegister(t, store, first)

	second := validToolSpec("tool-second")
	second.Tags = []string{"search", "external"}
	mustRegister(t, store, second)
	mustSetStatus(t, store, second.ID, spec.ToolStatusDisabled)

	got, err := store.List(context.Background(), registry.ListRequest{
		Filter: registry.ListFilter{
			Status: spec.ToolStatusDraft,
			Tags:   []string{"search", "internal"},
		},
	})
	if err != nil {
		t.Fatalf("List() error = %v, want nil", err)
	}
	if len(got) != 1 || got[0].ID != first.ID {
		t.Fatalf("List() = %#v, want only %q", got, first.ID)
	}
}

func TestListRejectsInvalidStatusFilter(t *testing.T) {
	store := New()

	_, err := store.List(context.Background(), registry.ListRequest{
		Filter: registry.ListFilter{Status: "archived"},
	})
	if !errors.Is(err, registry.InvalidSpec) {
		t.Fatalf("List() error = %v, want %v", err, registry.InvalidSpec)
	}
}

func TestUpdateSuccess(t *testing.T) {
	store := New()
	registered := mustRegister(t, store, validToolSpec("tool-search"))

	got, err := store.Update(context.Background(), registry.UpdateRequest{
		ID: registered.ID,
		Update: registry.ToolSpecUpdate{
			DisplayName: "Search v2",
			Description: "Search known records with updated metadata.",
			Version:     "2.0.0",
			Tags:        []string{"records", "v2"},
			Actions: []spec.ToolAction{
				{Name: "search", DisplayName: "Search", Description: "Search updated records."},
			},
			Metadata: map[string]string{"owner": "platform", "tier": "dev"},
		},
	})
	if err != nil {
		t.Fatalf("Update() error = %v, want nil", err)
	}
	if got.ID != registered.ID || got.Name != registered.Name || got.Status != registered.Status {
		t.Fatalf("Update() changed immutable fields: %#v", got)
	}
	if got.DisplayName != "Search v2" || got.Version != "2.0.0" || got.Tags[1] != "v2" || got.Metadata["tier"] != "dev" {
		t.Fatalf("Update() = %#v, want updated mutable fields", got)
	}
}

func TestUpdateNotFound(t *testing.T) {
	store := New()

	_, err := store.Update(context.Background(), registry.UpdateRequest{
		ID:     "missing",
		Update: validUpdate(),
	})
	if !errors.Is(err, registry.NotFound) {
		t.Fatalf("Update() error = %v, want %v", err, registry.NotFound)
	}
}

func TestUpdateInvalidSpec(t *testing.T) {
	store := New()
	registered := mustRegister(t, store, validToolSpec("tool-search"))
	update := validUpdate()
	update.Version = "2.0 beta"

	_, err := store.Update(context.Background(), registry.UpdateRequest{
		ID:     registered.ID,
		Update: update,
	})
	if !errors.Is(err, registry.InvalidSpec) {
		t.Fatalf("Update() error = %v, want %v", err, registry.InvalidSpec)
	}
}

func TestDeleteSuccess(t *testing.T) {
	store := New()
	registered := mustRegister(t, store, validToolSpec("tool-search"))

	if err := store.Delete(context.Background(), registry.DeleteRequest{ID: registered.ID}); err != nil {
		t.Fatalf("Delete() error = %v, want nil", err)
	}

	_, err := store.Get(context.Background(), registry.GetRequest{ID: registered.ID})
	if !errors.Is(err, registry.NotFound) {
		t.Fatalf("Get() after Delete() error = %v, want %v", err, registry.NotFound)
	}
}

func TestDeleteNotFound(t *testing.T) {
	store := New()

	err := store.Delete(context.Background(), registry.DeleteRequest{ID: "missing"})
	if !errors.Is(err, registry.NotFound) {
		t.Fatalf("Delete() error = %v, want %v", err, registry.NotFound)
	}
}

func TestSetStatusSuccess(t *testing.T) {
	store := New()
	registered := mustRegister(t, store, validToolSpec("tool-search"))

	got := mustSetStatus(t, store, registered.ID, spec.ToolStatusActive)
	if got.Status != spec.ToolStatusActive {
		t.Fatalf("SetStatus() status = %q, want %q", got.Status, spec.ToolStatusActive)
	}
}

func TestSetStatusNotFound(t *testing.T) {
	store := New()

	_, err := store.SetStatus(context.Background(), registry.SetStatusRequest{
		ID:     "missing",
		Status: spec.ToolStatusActive,
	})
	if !errors.Is(err, registry.NotFound) {
		t.Fatalf("SetStatus() error = %v, want %v", err, registry.NotFound)
	}
}

func TestInvalidStatusTransition(t *testing.T) {
	store := New()
	registered := mustRegister(t, store, validToolSpec("tool-search"))
	mustSetStatus(t, store, registered.ID, spec.ToolStatusDisabled)

	_, err := store.SetStatus(context.Background(), registry.SetStatusRequest{
		ID:     registered.ID,
		Status: spec.ToolStatusActive,
	})
	if !errors.Is(err, registry.InvalidStateTransition) {
		t.Fatalf("SetStatus() error = %v, want %v", err, registry.InvalidStateTransition)
	}
}

func TestSetStatusRejectsInvalidStatus(t *testing.T) {
	store := New()
	registered := mustRegister(t, store, validToolSpec("tool-search"))

	_, err := store.SetStatus(context.Background(), registry.SetStatusRequest{
		ID:     registered.ID,
		Status: "archived",
	})
	if !errors.Is(err, registry.InvalidSpec) {
		t.Fatalf("SetStatus() error = %v, want %v", err, registry.InvalidSpec)
	}
}

func TestConcurrentRegisterGetList(t *testing.T) {
	store := New()
	const count = 64

	var wg sync.WaitGroup
	start := make(chan struct{})
	for i := 0; i < count; i++ {
		i := i
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-start

			id := spec.ToolID(fmt.Sprintf("tool-%02d", i))
			registered, err := store.Register(context.Background(), registry.RegisterRequest{Spec: validToolSpec(id)})
			if err != nil {
				t.Errorf("Register(%q) error = %v, want nil", id, err)
				return
			}
			if _, err := store.Get(context.Background(), registry.GetRequest{ID: registered.ID}); err != nil {
				t.Errorf("Get(%q) error = %v, want nil", id, err)
			}
			if _, err := store.List(context.Background(), registry.ListRequest{}); err != nil {
				t.Errorf("List() error = %v, want nil", err)
			}
		}()
	}
	close(start)
	wg.Wait()

	got, err := store.List(context.Background(), registry.ListRequest{})
	if err != nil {
		t.Fatalf("List() error = %v, want nil", err)
	}
	if len(got) != count {
		t.Fatalf("List() len = %d, want %d", len(got), count)
	}
}

func validToolSpec(id spec.ToolID) spec.ToolSpec {
	return spec.ToolSpec{
		ID:          id,
		Name:        string(id),
		DisplayName: "Search",
		Description: "Search known records.",
		Version:     "1.0.0",
		Status:      spec.ToolStatusDraft,
		Tags:        []string{"records"},
		Actions: []spec.ToolAction{
			{Name: "search", DisplayName: "Search", Description: "Search records."},
		},
		Metadata: map[string]string{"owner": "platform"},
	}
}

func validUpdate() registry.ToolSpecUpdate {
	return registry.ToolSpecUpdate{
		DisplayName: "Search v2",
		Description: "Search known records with updated metadata.",
		Version:     "2.0.0",
		Tags:        []string{"records", "v2"},
		Actions: []spec.ToolAction{
			{Name: "search", DisplayName: "Search", Description: "Search updated records."},
		},
		Metadata: map[string]string{"owner": "platform"},
	}
}

func mustRegister(t *testing.T, store *Registry, tool spec.ToolSpec) spec.ToolSpec {
	t.Helper()

	registered, err := store.Register(context.Background(), registry.RegisterRequest{Spec: tool})
	if err != nil {
		t.Fatalf("Register() error = %v, want nil", err)
	}

	return registered
}

func mustSetStatus(t *testing.T, store *Registry, id spec.ToolID, status spec.ToolStatus) spec.ToolSpec {
	t.Helper()

	updated, err := store.SetStatus(context.Background(), registry.SetStatusRequest{
		ID:     id,
		Status: status,
	})
	if err != nil {
		t.Fatalf("SetStatus() error = %v, want nil", err)
	}

	return updated
}
