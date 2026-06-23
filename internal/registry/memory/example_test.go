package memory_test

import (
	"context"
	"fmt"

	"github.com/toolvault/toolvault/internal/registry"
	"github.com/toolvault/toolvault/internal/registry/memory"
	"github.com/toolvault/toolvault/internal/registry/spec"
)

func ExampleRegistry_basicUsage() {
	ctx := context.Background()
	store := memory.New()

	registered, err := store.Register(ctx, registry.RegisterRequest{
		Spec: spec.ToolSpec{
			ID:          "tool-search",
			Name:        "search",
			DisplayName: "Search",
			Description: "Search known records.",
			Version:     "1.0.0",
			Status:      spec.ToolStatusDraft,
			Tags:        []string{"records"},
			Actions: []spec.ToolAction{
				{
					Name:        "search",
					DisplayName: "Search",
					Description: "Search records.",
				},
			},
			Metadata: map[string]string{"owner": "platform"},
		},
	})
	if err != nil {
		panic(err)
	}

	lookup, err := store.Get(ctx, registry.GetRequest{ID: registered.ID})
	if err != nil {
		panic(err)
	}

	listed, err := store.List(ctx, registry.ListRequest{
		Filter: registry.ListFilter{Tags: []string{"records"}},
	})
	if err != nil {
		panic(err)
	}

	active, err := store.SetStatus(ctx, registry.SetStatusRequest{
		ID:     registered.ID,
		Status: spec.ToolStatusActive,
	})
	if err != nil {
		panic(err)
	}

	fmt.Println(lookup.ID, lookup.Version)
	fmt.Println(len(listed))
	fmt.Println(active.Status)

	// Output:
	// tool-search 1.0.0
	// 1
	// active
}
