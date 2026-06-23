#!/usr/bin/env sh
set -eu

required_dirs="
internal/registry
"

for dir in $required_dirs; do
  test -d "$dir" || {
    echo "missing required Registry path: $dir" >&2
    exit 1
  }
done

for forbidden_dir in \
  cmd/toolvault \
  internal/gateway \
  internal/runtime \
  internal/policy \
  internal/credential \
  internal/protocol \
  internal/observability \
  pkg
do
  if test -d "$forbidden_dir"; then
    if find "$forbidden_dir" -type f \( -name '*.go' -o -name '*.sql' -o -name '*.proto' \) | grep -q .; then
      echo "M1 Registry must not add implementation files under $forbidden_dir" >&2
      exit 1
    fi
  fi
done

if find . \
  -path './.git' -prune -o \
  -path './.codex' -prune -o \
  -type f \( \
    -name 'go.sum' -o \
    -name 'package.json' -o \
    -name 'package-lock.json' -o \
    -name 'pnpm-lock.yaml' -o \
    -name 'yarn.lock' -o \
    -name 'bun.lockb' -o \
    -name 'vite.config.*' -o \
    -name 'next.config.*' -o \
    -name '*.sql' \
  \) -print | grep -q .; then
  echo "M1 Registry must not introduce dependency, Web UI, or database artifacts" >&2
  exit 1
fi

if grep -Eq '^[[:space:]]*require[[:space:]]' go.mod; then
  echo "M1 Registry must remain standard-library only; go.mod contains require directives" >&2
  exit 1
fi

if find internal/registry -type f -name '*.go' | grep -q .; then
  go test ./internal/registry/...
else
  echo "no Registry Go packages yet; reserved: go test ./internal/registry/..."
fi

echo "m1 registry check passed"
