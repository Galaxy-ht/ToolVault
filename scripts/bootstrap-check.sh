#!/usr/bin/env sh
set -eu

required_files="
README.md
AGENTS.md
PROJECT_CONSTITUTION.md
docs/00-vision.md
docs/01-scope.md
docs/02-architecture.md
docs/03-milestones.md
docs/04-rules.md
docs/05-m0-bootstrap-checklist.md
docs/06-m1-registry-tasks.md
Makefile
scripts/bootstrap-check.sh
.github/workflows/bootstrap.yml
go.mod
"

required_dirs="
cmd/toolvault
internal/registry
internal/gateway
internal/runtime
internal/policy
internal/credential
internal/observability
pkg
docs/proposals
.codex/skills/toolvault-bootstrap-orchestrator
.codex/skills/toolvault-module-builder
.codex/skills/toolvault-reviewer
.codex/skills/toolvault-registry-builder
.codex/skills/toolvault-architect
"

for file in $required_files; do
  test -f "$file" || {
    echo "missing required file: $file" >&2
    exit 1
  }
done

for dir in $required_dirs; do
  test -d "$dir" || {
    echo "missing required directory: $dir" >&2
    exit 1
  }
done

if find internal cmd pkg -type f \( -name '*.go' -o -name '*.sql' -o -name '*.proto' \) | grep -q .; then
  echo "M0 must not contain business implementation files under internal, cmd, or pkg" >&2
  exit 1
fi

if find . -path './.git' -prune -o -path './.codex' -prune -o -type f \( -name 'package.json' -o -name 'go.sum' \) -print | grep -q .; then
  echo "unexpected dependency manifest or lock file for M0" >&2
  exit 1
fi

for skill in .codex/skills/*/SKILL.md; do
  grep -q '^name:' "$skill" || {
    echo "skill missing name metadata: $skill" >&2
    exit 1
  }
  grep -q '^description:' "$skill" || {
    echo "skill missing description metadata: $skill" >&2
    exit 1
  }
done

grep -q "Goal" docs/06-m1-registry-tasks.md
grep -q "Allowed directories" docs/06-m1-registry-tasks.md
grep -q "Forbidden directories" docs/06-m1-registry-tasks.md
grep -q "Acceptance criteria" docs/06-m1-registry-tasks.md
grep -q "Verification commands" docs/06-m1-registry-tasks.md
grep -q "Risks" docs/06-m1-registry-tasks.md

echo "bootstrap check passed"
