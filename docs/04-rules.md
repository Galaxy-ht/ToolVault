# ToolVault Engineering Rules

## Scope Rules

- Keep v1 limited to Registry, Gateway, Runtime Manager, Policy Engine,
  Credential Vault, and Observability.
- Keep Plugin/Future areas out of v1 unless approved.
- Keep M0 and M1 free of Web UI.
- Keep M0 free of business logic.
- Keep M1 limited to Registry unless approved.

## Dependency Rules

- Do not introduce heavyweight frameworks without approval.
- Do not introduce persistence without approval.
- Do not introduce network protocols without approval.
- Do not introduce distributed deployment without approval.
- Do not add packages under `pkg/` without approval.

## Module Boundary Rules

- Each core module owns its own concepts and tests.
- Cross-module calls require explicit interfaces.
- Boundary changes require a Proposal under `docs/proposals/`.
- Shared code must have clear ownership.

## Agent Task Rules

Every task must include:

- Goal
- Allowed directories
- Forbidden directories
- Acceptance criteria
- Verification commands
- Risks

Agents must not modify directories outside the task contract.

## Review Rules

Reviewers must check:

- scope drift
- boundary violations
- dependency creep
- missing tests
- weak acceptance criteria
- unapproved persistence
- unapproved UI
- unapproved distributed architecture
- secret leakage risk

## Proposal Rules

Create a Proposal when a task needs:

- module boundary changes
- new dependencies
- persistence
- external protocols
- Web UI
- distributed deployment
- future/plugin scope
- changes to non-goals

Proposal files belong under `docs/proposals/`.
