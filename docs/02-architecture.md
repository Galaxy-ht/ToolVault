# ToolVault Architecture

## Architecture Style

ToolVault starts as a modular monolith written in Go.

The repository should keep module boundaries explicit while avoiding premature
distribution. The system can later evolve into distributed services only after
human-approved architecture proposals.

## Initial Repository Structure

```text
cmd/
  toolvault/
internal/
  registry/
  gateway/
  runtime/
  policy/
  credential/
  observability/
pkg/
docs/
scripts/
.codex/skills/
```

## Module Responsibilities

### Registry

Owns tool definitions, metadata, lifecycle state, and lookup.

### Gateway

Owns controlled ingress for tool invocation and protocol-facing request flow.

### Runtime Manager

Owns runtime selection, execution lifecycle, and runtime governance decisions.

### Policy Engine

Owns authorization and policy decisions around tool access and execution.

### Credential Vault

Owns credential references and injection boundaries. Secrets must not be exposed
to agents or logs.

### Observability

Owns logs, audit events, metrics, and traces for control-plane operations.

## Dependency Direction

M1 must define concrete dependency rules before implementing business logic.
Until approved, use these provisional constraints:

- Core modules should not import each other directly without an approved
  interface boundary.
- Shared types must not become a dumping ground.
- `pkg/` must remain empty unless a public package is explicitly approved.
- New persistence, network protocols, and third-party dependencies require a
  Proposal.

## Interface-First Rule

Before business logic is added to a core module, the task must define:

- owned concepts
- public interface surface
- error model
- tests
- dependencies
- forbidden imports

## Proposals

All module boundary changes require a Proposal under `docs/proposals/`.

A Proposal must include:

- problem
- proposed change
- alternatives considered
- affected modules
- dependency impact
- test impact
- rollback plan
- required human decision
