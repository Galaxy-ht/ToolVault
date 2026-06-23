# Proposal 001: M1 Registry Implementation Plan

Status: proposed

Owner role: ToolVault Architect

## 1. Registry Goal

The Registry provides the authoritative catalog of tools known to ToolVault.
It owns tool metadata, lifecycle state, validation, lookup, listing, and update
semantics for registered tools.

M1 should produce a small, testable Registry module that other modules can later
depend on through explicit interfaces. M1 must not couple Registry to Gateway,
Runtime Manager, Policy Engine, Credential Vault, Observability, persistence, or
external protocols.

## 2. M1 Required Capabilities

M1 must complete:

- Define Registry-owned domain concepts and non-goals.
- Define the minimal `ToolSpec` model.
- Define Registry lifecycle states.
- Define a small public interface for registering, retrieving, listing, updating,
  and deleting tool specs.
- Define validation behavior and errors.
- Define an M1 Registry verification gate before adding Registry `.go` files.
- Implement an in-memory Registry only after this plan is approved.
- Add tests for all implemented behavior.
- Keep implementation standard-library only.
- Keep all Registry behavior inside `internal/registry/`.

## 3. M1 Explicit Non-Goals

M1 must not implement:

- Gateway request handling or protocol adaptation.
- Runtime execution or runtime selection.
- Policy decisions or authorization.
- Credential retrieval, storage, or injection.
- Observability pipelines, metrics exporters, tracing, or audit sinks.
- Database-backed persistence.
- Web UI, dashboard, or marketplace features.
- Streaming or tool composition.
- Public packages under `pkg/`.
- Network servers or external protocol endpoints.

## 4. Directory Boundaries

Allowed for M1 planning:

- `docs/`
- `docs/proposals/`

Allowed for approved Registry implementation tasks:

- `internal/registry/`
- `docs/`

Conditionally allowed:

- `cmd/toolvault/` only for an explicitly approved developer-only Registry probe.

Forbidden during M1:

- `internal/gateway/`
- `internal/runtime/`
- `internal/policy/`
- `internal/credential/`
- `internal/observability/`
- `pkg/`

## 5. Module Dependencies

Registry dependency rules:

- `internal/registry/` may depend on Go standard library packages.
- `internal/registry/` must not import other ToolVault core modules in M1.
- Other core modules must not be modified during M1.
- `pkg/` remains empty unless separately approved.
- Storage implementation must stay in-memory for M1 unless a persistence
  Proposal is approved.
- Validation and lifecycle rules belong inside Registry.

Expected future dependency direction:

- Gateway may later depend on Registry interfaces after an approved integration
  Proposal.
- Runtime Manager, Policy Engine, Credential Vault, and Observability must not
  depend on Registry internals.

## 6. Public Interface Draft

This is the M1 contract proposal, not an implementation.

```go
type Registry interface {
    Register(ctx context.Context, spec ToolSpec) (ToolSpec, error)
    Get(ctx context.Context, id ToolID) (ToolSpec, error)
    List(ctx context.Context, filter ListFilter) ([]ToolSpec, error)
    Update(ctx context.Context, id ToolID, update ToolSpecUpdate) (ToolSpec, error)
    Delete(ctx context.Context, id ToolID) error
}
```

Supporting types:

```go
type ToolID string

type ToolStatus string

const (
    ToolStatusDraft      ToolStatus = "draft"
    ToolStatusActive     ToolStatus = "active"
    ToolStatusDeprecated ToolStatus = "deprecated"
    ToolStatusDisabled   ToolStatus = "disabled"
)

type ListFilter struct {
    Status ToolStatus
    Tags   []string
}

type ToolSpecUpdate struct {
    DisplayName string
    Description string
    Status      ToolStatus
    Tags        []string
}
```

Error categories:

- invalid spec
- duplicate tool ID
- tool not found
- invalid lifecycle transition

M1 implementation decisions:

- `Register` accepts caller-provided IDs only.
- `ToolID` is required, stable, and non-empty.
- M1 does not generate IDs.
- M1 does not validate global uniqueness outside the in-memory Registry instance.
- `Version` is an opaque tool definition version string. It must be non-empty,
  but M1 does not require semantic versioning.
- `Delete` is a hard delete from the in-memory Registry. Lifecycle transitions
  are handled through `Update`, not `Delete`.

Lifecycle transition matrix:

| From | To | M1 Rule |
| ---- | -- | ------- |
| none | draft | allowed on register |
| none | active | allowed on register |
| draft | active | allowed |
| draft | disabled | allowed |
| active | deprecated | allowed |
| active | disabled | allowed |
| deprecated | disabled | allowed |
| disabled | active | rejected |
| deprecated | active | rejected |
| any | same status | allowed |

## 7. Minimal ToolSpec Fields

M1 `ToolSpec` should contain only Registry-owned metadata:

```go
type ToolSpec struct {
    ID          ToolID
    Name        string
    DisplayName string
    Description string
    Version     string
    Status      ToolStatus
    Tags        []string
}
```

Field intent:

- `ID`: stable Registry identity.
- `Name`: stable machine-readable name.
- `DisplayName`: human-readable name.
- `Description`: short purpose statement.
- `Version`: opaque tool definition version, not runtime version.
- `Status`: Registry lifecycle state.
- `Tags`: optional metadata for filtering.

M1 required fields:

- `ID`
- `Name`
- `Version`
- `Status`

M1 optional fields:

- `DisplayName`
- `Description`
- `Tags`

Fields intentionally excluded from M1:

- executable command or runtime config
- protocol schema
- credentials
- policy rules
- observability config
- marketplace metadata
- owner/team model
- audit history
- persistence metadata

## 8. Subtask Breakdown

### M1.0 Registry Verification Gate

Goal: define the M1 verification command that permits approved Registry `.go`
files while continuing to reject non-Registry business logic.

Allowed modification directories:

- `docs/`
- `scripts/`
- `Makefile`

Forbidden directories:

- `cmd/toolvault/`
- `internal/registry/`
- `internal/gateway/`
- `internal/runtime/`
- `internal/policy/`
- `internal/credential/`
- `internal/observability/`
- `pkg/`

Acceptance criteria:

- M0 `make bootstrap-check` remains valid for M0 constraints.
- A separate M1 Registry verification command is documented or added.
- The M1 command allows approved files under `internal/registry/`.
- The M1 command rejects changes under forbidden core module directories.
- The M1 command rejects database, Web UI, and third-party dependency drift.

Verification commands:

```sh
make bootstrap-check
```

### M1.1 Registry Domain Proposal

Goal: finalize Registry-owned concepts, lifecycle states, validation scope, and
non-goals.

Allowed modification directories:

- `docs/`
- `docs/proposals/`

Forbidden directories:

- `cmd/toolvault/`
- `internal/registry/`
- `internal/gateway/`
- `internal/runtime/`
- `internal/policy/`
- `internal/credential/`
- `internal/observability/`
- `pkg/`

Acceptance criteria:

- Registry owns metadata, lifecycle, validation, lookup, listing, update, and
  deletion semantics.
- Proposal explicitly excludes Gateway, Runtime, Policy, Credential, and
  Observability behavior.
- ID strategy, version semantics, lifecycle transitions, delete semantics, and
  minimal `ToolSpec` fields are specified.
- Implementation waits for human approval of this Proposal.

Verification commands:

```sh
make bootstrap-check
```

### M1.2 Registry Interface And Error Contract

Goal: create approved Registry interfaces, domain types, and error categories.

Allowed modification directories:

- `internal/registry/`
- `docs/`

Forbidden directories:

- `cmd/toolvault/`
- `internal/gateway/`
- `internal/runtime/`
- `internal/policy/`
- `internal/credential/`
- `internal/observability/`
- `pkg/`

Acceptance criteria:

- Interfaces are small and do not expose storage details.
- `ToolSpec` includes only approved M1 fields.
- Error categories are testable and documented.
- No implementation storage is added in this subtask.
- No imports from other ToolVault core modules.

Verification commands:

```sh
go test ./...
make m1-registry-check
```

### M1.3 Registry Validation Tests

Goal: define contract tests for `ToolSpec` validation and lifecycle behavior.

Allowed modification directories:

- `internal/registry/`
- `docs/`

Forbidden directories:

- `cmd/toolvault/`
- `internal/gateway/`
- `internal/runtime/`
- `internal/policy/`
- `internal/credential/`
- `internal/observability/`
- `pkg/`

Acceptance criteria:

- Tests cover required fields.
- Tests cover duplicate identity behavior.
- Tests cover unknown ID behavior.
- Tests cover invalid lifecycle transitions.
- Tests cover delete as hard delete.
- Tests cover `Version` as required opaque metadata.
- Tests do not require database, network, or other core modules.

Verification commands:

```sh
go test ./...
make m1-registry-check
```

### M1.4 In-Memory Registry Implementation

Goal: implement the approved Registry contract using in-memory storage and the
Go standard library.

Allowed modification directories:

- `internal/registry/`
- `docs/`

Forbidden directories:

- `cmd/toolvault/`
- `internal/gateway/`
- `internal/runtime/`
- `internal/policy/`
- `internal/credential/`
- `internal/observability/`
- `pkg/`

Acceptance criteria:

- Implementation satisfies the approved interface.
- Register, get, list, update, and delete behavior is tested.
- Validation and lifecycle tests pass.
- `Register` requires caller-provided non-empty IDs.
- `Delete` performs hard delete only.
- No third-party dependency is introduced.
- No database or file persistence is introduced.
- No other core module is modified.

Verification commands:

```sh
go test ./...
make m1-registry-check
```

### M1.5 Optional Registry Developer Probe Proposal

Goal: decide whether M1 needs a developer-only CLI probe.

Allowed modification directories:

- `docs/`
- `docs/proposals/`

Forbidden directories:

- `cmd/toolvault/` unless this subtask receives explicit human approval to
  implement the probe in a later task.
- `internal/gateway/`
- `internal/runtime/`
- `internal/policy/`
- `internal/credential/`
- `internal/observability/`
- `pkg/`

Acceptance criteria:

- The probe is treated as optional and not a default M1 implementation item.
- Human explicitly approves or rejects the probe.
- If approved, the probe is documented as non-product surface.
- If rejected, no `cmd/toolvault/` implementation is added in M1.

Verification commands:

```sh
make bootstrap-check
```

### M1.6 Registry Review Gate

Goal: review M1 for scope, boundaries, dependencies, and test coverage before
any Gateway work begins.

Allowed modification directories:

- `docs/`
- `internal/registry/`

Forbidden directories:

- `cmd/toolvault/`
- `internal/gateway/`
- `internal/runtime/`
- `internal/policy/`
- `internal/credential/`
- `internal/observability/`
- `pkg/`

Acceptance criteria:

- Review confirms no non-Registry business logic.
- Review confirms no database, Web UI, or third-party dependency.
- Review confirms tests cover the approved Registry contract.
- Review lists required decisions before M2 Gateway.

Verification commands:

```sh
go test ./...
make m1-registry-check
```

## 9. Recommended Subagents

- `toolvault-architect`: finalize and revise this Proposal.
- `toolvault-registry-builder`: implement approved M1 Registry subtasks.
- `toolvault-reviewer`: review each M1 change for scope drift, dependency creep,
  and boundary violations.
- `toolvault-module-builder`: use only for narrowly assigned implementation
  subtasks with explicit directory boundaries.

## 10. Alternatives Considered

- Add database-backed Registry in M1: rejected because persistence is not
  approved and would expand M1 blast radius.
- Put shared types in `pkg/`: rejected because public package exposure is not
  approved.
- Add Gateway-facing protocol fields to `ToolSpec`: rejected because Gateway
  ownership is out of M1 scope.
- Add credentials or policy references to `ToolSpec`: rejected because Credential
  Vault and Policy Engine are separate modules.

## 11. Affected Modules

Affected in M1 after approval:

- `internal/registry/`
- `docs/`

Not affected in M1:

- `internal/gateway/`
- `internal/runtime/`
- `internal/policy/`
- `internal/credential/`
- `internal/observability/`
- `pkg/`

## 12. Dependency Impact

No new third-party dependency is proposed.

No database, network protocol, Web UI, or distributed deployment is proposed.

M1 requires an M1-specific verification gate because the M0 bootstrap gate
intentionally rejects all `.go` files under `internal/`, `cmd/`, and `pkg/`.

## 13. Test Impact

M1 implementation tasks must add Go tests under `internal/registry/`.

Minimum test coverage areas:

- valid registration
- invalid registration
- missing ID
- missing name
- missing version
- duplicate ID
- lookup by ID
- unknown ID
- list filtering
- update behavior
- hard delete behavior
- lifecycle transition validation
- forbidden lifecycle transition rejection
- contract tests that exercise behavior through the `Registry` interface

## 14. Rollback Plan

If M1 scope drifts, revert Registry implementation changes and keep only the
approved Proposal documents. Because M1 avoids persistence and external
dependencies, rollback should not require data migration or environment cleanup.

## 15. Required Human Decisions

Before implementation:

- Approve or revise this M1 Registry implementation plan.
- Confirm caller-provided IDs only.
- Confirm the minimal `ToolSpec` field set.
- Confirm lifecycle statuses and transition rules.
- Confirm hard delete semantics.

Before optional CLI probe:

- Approve whether `cmd/toolvault/` may contain a developer-only Registry probe.
