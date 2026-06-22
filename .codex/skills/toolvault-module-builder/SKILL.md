---
name: toolvault-module-builder
description: Use when implementing an approved ToolVault module task with explicit allowed directories, forbidden directories, acceptance criteria, verification commands, and risks.
---

# ToolVault Module Builder

Use this role only for an approved implementation task.

## Required Inputs

Do not start without a task contract containing:

- Goal
- Allowed directories
- Forbidden directories
- Acceptance criteria
- Verification commands
- Risks

## Rules

- Modify only allowed directories.
- Do not touch forbidden directories.
- Do not change module boundaries without an approved Proposal.
- Do not add dependencies without approval.
- Add focused tests for implemented behavior.
- Keep implementations small and aligned with existing docs.

## Workflow

1. Read `AGENTS.md`, `docs/01-scope.md`, and `docs/04-rules.md`.
2. Read the assigned task and any approved Proposal.
3. Inspect only relevant files.
4. Implement the smallest compliant change.
5. Run the task verification commands.
6. Report any residual risks or blocked decisions.
