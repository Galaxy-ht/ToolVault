# M0 Bootstrap Acceptance Checklist

This checklist mirrors `scripts/bootstrap-check.sh`. The final M0 gate is
`make bootstrap-check`; if this document and the script ever disagree, update
the document or script before accepting M0.

## Required Files

- [ ] `README.md`
- [ ] `AGENTS.md`
- [ ] `PROJECT_CONSTITUTION.md`
- [ ] `docs/00-vision.md`
- [ ] `docs/01-scope.md`
- [ ] `docs/02-architecture.md`
- [ ] `docs/03-milestones.md`
- [ ] `docs/04-rules.md`
- [ ] `docs/05-m0-bootstrap-checklist.md`
- [ ] `docs/06-m1-registry-tasks.md`
- [ ] `Makefile`
- [ ] `scripts/bootstrap-check.sh`
- [ ] `.github/workflows/bootstrap.yml`
- [ ] `go.mod`

## Required Directories

- [ ] `cmd/toolvault/`
- [ ] `internal/registry/`
- [ ] `internal/gateway/`
- [ ] `internal/runtime/`
- [ ] `internal/policy/`
- [ ] `internal/credential/`
- [ ] `internal/observability/`
- [ ] `pkg/`
- [ ] `docs/proposals/`
- [ ] `.codex/skills/toolvault-bootstrap-orchestrator/`
- [ ] `.codex/skills/toolvault-module-builder/`
- [ ] `.codex/skills/toolvault-reviewer/`
- [ ] `.codex/skills/toolvault-registry-builder/`
- [ ] `.codex/skills/toolvault-architect/`

## Required Constraints

- [ ] No `.go`, `.sql`, or `.proto` files under `internal/`, `cmd/`, or `pkg/`.
- [ ] No Registry business logic.
- [ ] No Gateway business logic.
- [ ] No Runtime Manager business logic.
- [ ] No Policy Engine business logic.
- [ ] No Credential Vault business logic.
- [ ] No Observability business logic.
- [ ] No `package.json`.
- [ ] No `go.sum`.
- [ ] No unapproved third-party dependency.
- [ ] Every `.codex/skills/*/SKILL.md` contains `name:` metadata.
- [ ] Every `.codex/skills/*/SKILL.md` contains `description:` metadata.
- [ ] M1 task breakdown includes `Goal`.
- [ ] M1 task breakdown includes `Allowed directories`.
- [ ] M1 task breakdown includes `Forbidden directories`.
- [ ] M1 task breakdown includes `Acceptance criteria`.
- [ ] M1 task breakdown includes `Verification commands`.
- [ ] M1 task breakdown includes `Risks`.

## Verification

Run:

```sh
make bootstrap-check
```

M0 is complete only when `make bootstrap-check` passes and the checklist is
kept aligned with `scripts/bootstrap-check.sh`.
