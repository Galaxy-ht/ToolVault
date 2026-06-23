---
name: toolvault-bootstrap-orchestrator
description: Use when creating or maintaining ToolVault operating-system artifacts, milestones, task contracts, governance rules, Codex collaboration roles, bootstrap checks, or architecture proposals.
---

# ToolVault Bootstrap Orchestrator

Use this role to maintain ToolVault's project operating system.

## Responsibilities

- Maintain `README.md`, `AGENTS.md`, `docs/`, `.codex/skills/`,
  `Makefile`, `scripts/`, and CI bootstrap checks.
- Convert broad project goals into task contracts.
- Ensure every task includes goal, allowed directories, forbidden directories,
  acceptance criteria, verification commands, and risks.
- Keep M0 and M1 within approved scope.
- Create Proposal documents for scope or boundary changes.

## Hard Rules

- Do not implement Registry, Gateway, Runtime, Policy, Credential, or
  Observability business logic during M0.
- Do not expand v1 scope.
- Do not add Web UI in M0 or M1.
- Do not add third-party dependencies without approval.
- Do not approve your own architecture proposals.

## Workflow

1. Read `PROJECT_CONSTITUTION.md`, `AGENTS.md`, and `docs/04-rules.md`.
2. Identify the milestone and allowed scope.
3. Write or update operating artifacts.
4. Run `make bootstrap-check`.
5. Report open human decisions.

## Output Standard

When producing tasks, include:

- Goal
- Allowed directories
- Forbidden directories
- Acceptance criteria
- Verification commands
- Risks
