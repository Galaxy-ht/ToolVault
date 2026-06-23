# M1 Registry Acceptance Checklist

Status: accepted with documented residual risks

Date: 2026-06-23

## Scope Summary

M1 delivers the first Tool Registry slice inside `internal/registry/`.

Completed capabilities:

- `ToolSpec` model and validation in `internal/registry/spec`.
- Registry interface, request/option types, lifecycle helpers, and error
  contract in `internal/registry`.
- In-memory Registry implementation in `internal/registry/memory`.
- Registry-owned lifecycle states: `draft`, `active`, `deprecated`,
  `disabled`.
- Register, get, list, update, delete, and set-status behavior.
- Opaque version validation plus optimistic update checks through
  `ExpectedVersion`.
- Hard delete semantics for M1.
- Stable list ordering, status/tag filtering, duplicate identity checks, and
  defensive copies for mutable fields.

Out of scope and not implemented:

- Gateway/API server behavior.
- Runtime execution or runtime selection.
- Policy evaluation or authorization.
- Credential injection or secret storage.
- Protocol adapters such as MCP or OpenAI Function Calling.
- Database persistence.
- Observability pipelines.
- Dashboard, sandbox, streaming, composition, marketplace, or Web UI.

## Proposal Alignment

The implementation satisfies the M1 proposal direction: Registry behavior stays
inside `internal/registry/`, uses only the Go standard library, and does not
modify other core modules.

Documented deviations or extensions from Proposal 001:

- `ToolSpec` includes `Actions` and `Metadata`, while Proposal 001 listed a
  smaller minimum model.
- `Update` uses `ExpectedVersion`, requires a new opaque version, and may
  return `VersionConflict`; Proposal 001 did not include optimistic update
  semantics.
- Lifecycle status changes are handled by `SetStatus` rather than `Update`.
- The in-memory Registry rejects duplicate `Name` values in addition to
  duplicate `ID` values.

These deviations have tests and prior review coverage, but they should be
explicitly accepted or narrowed before M2 treats Registry as a stable dependency.

## Checklist

| Question | Answer |
| --- | --- |
| 1. ToolSpec model 是否完成？ | Yes. `ToolSpec`, `ToolID`, `ToolStatus`, `ToolAction`, and metadata fields exist under `internal/registry/spec`. Residual risk: model is wider than Proposal 001. |
| 2. ToolSpec validation 是否完成？ | Yes. Validation covers required ID/name/version/status, machine-readable names, opaque versions, valid statuses, actions, and metadata key limits. |
| 3. Registry interface 是否完成？ | Yes. `Registry` covers register, get, list, update, delete, and set-status. |
| 4. Error contract 是否完成？ | Yes. Error kinds cover already-exists, not-found, invalid-spec, invalid-state-transition, and version-conflict, with `errors.Is` support. |
| 5. In-memory Registry 是否完成？ | Yes. `memory.Registry` implements the interface with standard-library synchronization and defensive copies. |
| 6. Version / status lifecycle 是否完成？ | Yes. Version is opaque and required; metadata updates require expected-version and a new version; status transitions follow the M1 matrix through `SetStatus`. |
| 7. Tests 是否覆盖主要路径？ | Yes. Tests cover validation, errors, register/get/list/update/delete/status, filtering, sorting, duplicate checks, version conflict, hard delete, disabled/deprecated visibility, and concurrency basics. |
| 8. `make m1-registry-check` 是否通过？ | Yes. Last verified during this acceptance pass. |
| 9. `make verify` 是否通过？ | Yes. Last verified during this acceptance pass. |
| 10. M1 是否存在遗留风险？ | Yes. The main residual risks are proposal drift around `Actions`/`Metadata`, optimistic version update semantics, duplicate `Name` identity behavior, and the need to approve Registry consumption rules before Gateway work. |
| 11. 是否建议进入 M2？ | Yes, conditionally. M2 may begin after human confirmation that the documented Registry contract deviations are accepted or after a follow-up task narrows them to Proposal 001. |

## Acceptance Decision

M1 Registry is suitable to close as the internal Registry baseline for the next
planning step. The implementation is not yet a public API and should not be
consumed by Gateway, Runtime, Policy, Credential, Protocol, or Observability
modules until M2 has an approved integration proposal.
