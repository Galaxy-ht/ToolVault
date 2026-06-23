# Registry Basic Example

This example shows the smallest repository-local use of the M1 Registry:
create an in-memory Registry, register a valid tool spec, read it back, list it,
and move it through an allowed lifecycle transition.

The Registry is still an internal ToolVault module. This is not a public API,
network endpoint, Gateway flow, Runtime execution path, protocol adapter, or
persistence contract.

```go
package main

import (
	"context"
	"fmt"

	"github.com/toolvault/toolvault/internal/registry"
	"github.com/toolvault/toolvault/internal/registry/memory"
	"github.com/toolvault/toolvault/internal/registry/spec"
)

func main() {
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
}
```

Expected output:

```text
tool-search 1.0.0
1
active
```

For metadata updates, callers must provide `UpdateOptions.ExpectedVersion` and
a new opaque `ToolSpecUpdate.Version`. Status changes use `SetStatus` and do
not change the tool definition version.
