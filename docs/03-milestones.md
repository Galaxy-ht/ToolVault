# ToolVault Milestones

## M0 Bootstrap

Goal: create the project operating system for multi-agent work.

Deliverables:

- Repository structure
- README
- Agent operating rules
- Vision, scope, architecture, milestones, and rules docs
- Codex role skills
- Makefile
- Verification scripts
- CI bootstrap checks
- M0 acceptance checklist
- M1 Registry task breakdown

Exit criteria:

- `make bootstrap-check` passes.
- No core business logic exists.
- M1 tasks are decomposed with goal, allowed directories, forbidden
  directories, acceptance criteria, verification commands, and risks.
- Human approval items are documented.

## M1 Registry

Goal: define and implement the first Tool Registry slice.

M1 is limited to Tool Registry. It must not implement Gateway, Runtime Manager,
Policy Engine, Credential Vault, or Observability business logic.

M1 work starts from `docs/06-m1-registry-tasks.md`.

## M2 Gateway

Goal: introduce controlled tool invocation entrypoints after Registry boundaries
are stable.

Requires a human-approved architecture decision for protocol surface and
module dependency direction.

## M3 Runtime Governance

Goal: introduce runtime lifecycle decisions and execution governance.

Requires approved interfaces between Gateway, Runtime Manager, Policy Engine,
Credential Vault, and Observability.

## M4 Policy And Credentials

Goal: add policy enforcement and credential injection boundaries.

Requires explicit secret-handling rules and test strategy approval.

## M5 Observability Baseline

Goal: add audit, metrics, and tracing baseline for tool control-plane flows.

Requires approved event model and data retention assumptions.
