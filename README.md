# ToolVault

ToolVault is an enterprise Tool Control Plane for AI Agents.

It provides a unified control layer for tool registration, protocol adaptation,
runtime governance, policy enforcement, credential injection, and observability.

## v1 Core

- Tool Registry
- Gateway
- Runtime Manager
- Policy Engine
- Credential Vault
- Observability

## Future Plugins

- Streaming
- Tool Composition
- Sandbox
- Dashboard
- Marketplace

## Non-Goals

ToolVault is not:

- an Agent Framework
- an LLM Gateway
- a Workflow Engine replacement
- a Vector Database
- a Prompt Management Platform

## Technical Direction

- Go is the core implementation language.
- The default architecture is a modular monolith.
- M0 and M1 do not include a Web UI.
- M0 and M1 do not introduce heavy third-party dependencies without approval.
- Core modules must expose clear interfaces and tests when implemented.
- Module boundary changes must be proposed and approved before implementation.

## Repository Map

- `cmd/` - future command entrypoints.
- `internal/` - future internal modules for v1 core components.
- `pkg/` - future public packages, only when explicitly approved.
- `docs/` - vision, scope, architecture, milestones, rules, and task plans.
- `.codex/skills/` - Codex role skills for multi-agent collaboration.
- `.github/workflows/` - CI bootstrap checks.
- `scripts/` - local verification scripts.

## M0 Bootstrap

M0 establishes the project operating system. It creates documentation,
governance rules, collaboration roles, verification scripts, and CI checks.

M0 intentionally does not implement Registry, Gateway, Runtime, Policy,
Credential, or Observability business logic.

Run:

```sh
make bootstrap-check
```

## M1 Focus

M1 is limited to Tool Registry design and implementation planning. The task
breakdown lives in `docs/06-m1-registry-tasks.md`.
