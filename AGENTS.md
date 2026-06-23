# ToolVault Agent Operating Rules

These instructions apply to the entire repository.

## Project Positioning

ToolVault is the Tool Control Plane for AI Agents. It manages tool
registration, protocol adaptation, runtime governance, policy enforcement,
credential injection, and observability.

## Hard Boundaries

- Do not implement Agent Framework behavior.
- Do not implement an LLM Gateway.
- Do not replace workflow engines.
- Do not add Vector DB functionality.
- Do not add Prompt Management Platform functionality.
- Do not add Web UI during M0 or M1.
- Do not add heavyweight dependencies without explicit human approval.
- Do not split into microservices in v1 without explicit human approval.

## Architecture Rules

- Go is the core implementation language.
- Default to a modular monolith.
- Keep module boundaries explicit and documented.
- Treat changes to module boundaries as Proposals requiring human approval.
- Prefer interface-first design for core modules.
- Every implemented core module must have tests.

## Agent Roles

- Bootstrap Orchestrator: creates and maintains the operating system,
  milestones, task graph, rules, and validation checks.
- Builder Agent: implements only assigned tasks and only in allowed directories.
- Reviewer Agent: checks boundary violations, dependency creep, missing tests,
  and architecture drift.
- Module Architect: drafts interfaces and proposals, but does not apply boundary
  changes without approval.

## Task Contract

Every future task must include:

- Goal
- Allowed directories
- Forbidden directories
- Acceptance criteria
- Verification commands
- Risks

## Change Control

If a task requires scope expansion, new dependencies, new persistence,
distributed deployment, Web UI, or cross-module boundary changes, create a
Proposal document under `docs/proposals/` and stop for human approval.

## Verification

Before claiming completion, run the narrowest relevant checks. For M0 changes,
run:

```sh
make bootstrap-check
```
