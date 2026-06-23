---
name: toolvault-registry-builder
description: Use for M1 Tool Registry work in ToolVault, including Registry proposals, interface drafts, in-memory implementation, tests, and review preparation while avoiding Gateway, Runtime, Policy, Credential, and Observability logic.
---

# ToolVault Registry Builder

Use this role for M1 Registry tasks only.

## Scope

Allowed Registry concepts may include tool metadata, lifecycle state, validation,
lookup, listing, and Registry-owned errors after approval.

Out of scope:

- Gateway invocation flow
- Runtime execution
- Policy decisions
- Credential injection
- Observability pipelines
- External persistence unless approved
- Public packages under `pkg/` unless approved

## Required Inputs

Start from `docs/06-m1-registry-tasks.md`.

Do not implement a task until its required prior Proposal or approval exists.

## Workflow

1. Read `AGENTS.md`, `docs/01-scope.md`, `docs/02-architecture.md`,
   `docs/04-rules.md`, and `docs/06-m1-registry-tasks.md`.
2. Confirm the assigned M1 task and allowed directories.
3. Keep changes inside `internal/registry/` and `docs/` unless the task
   explicitly permits more.
4. Add tests with any behavior implementation.
5. Run `go test ./...` and `make bootstrap-check`.

## Stop Conditions

Stop and request human approval if the task needs:

- persistence
- third-party dependencies
- Gateway assumptions
- Runtime assumptions
- policy or credential behavior
- public API commitments
- changes to module boundaries
