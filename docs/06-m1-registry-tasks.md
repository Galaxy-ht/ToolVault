# M1 Registry Task Breakdown

M1 is limited to Tool Registry. Do not implement Gateway, Runtime Manager,
Policy Engine, Credential Vault, or Observability business logic during M1.

Each task below is intentionally scoped for a separate Codex Builder Agent.

## M1.1 Registry Domain Proposal

Goal: define the Registry domain concepts, ownership boundary, lifecycle states,
and non-goals before implementation.

Allowed directories:

- `docs/`

Forbidden directories:

- `internal/gateway/`
- `internal/runtime/`
- `internal/policy/`
- `internal/credential/`
- `internal/observability/`
- `pkg/`

Acceptance criteria:

- A Proposal exists under `docs/proposals/`.
- Proposal defines Registry-owned concepts.
- Proposal defines what Registry does not own.
- Proposal identifies affected modules.
- Proposal lists open human decisions.

Verification commands:

```sh
make bootstrap-check
```

Risks:

- Registry absorbs Gateway or Runtime responsibilities.
- Lifecycle state becomes too broad.
- Shared types are introduced before ownership is clear.

## M1.2 Registry Interface Draft

Goal: draft Registry interfaces and errors after M1.1 is approved.

Allowed directories:

- `internal/registry/`
- `docs/`

Forbidden directories:

- `internal/gateway/`
- `internal/runtime/`
- `internal/policy/`
- `internal/credential/`
- `internal/observability/`
- `pkg/`

Acceptance criteria:

- Interfaces are small and focused.
- Error behavior is documented.
- No storage implementation is added unless separately approved.
- Tests describe expected contract behavior.

Verification commands:

```sh
go test ./...
make bootstrap-check
```

Risks:

- Interface exposes storage details.
- Interface assumes unapproved protocol or persistence.
- Tests assert implementation instead of contract.

## M1.3 In-Memory Registry Implementation

Goal: implement the first approved Registry behavior using only approved
interfaces and standard library dependencies.

Allowed directories:

- `internal/registry/`
- `docs/`

Forbidden directories:

- `internal/gateway/`
- `internal/runtime/`
- `internal/policy/`
- `internal/credential/`
- `internal/observability/`
- `pkg/`

Acceptance criteria:

- Implementation satisfies approved Registry interfaces.
- Tests cover create, update, lookup, list, validation, and error cases.
- No external persistence is introduced.
- No third-party dependency is introduced.

Verification commands:

```sh
go test ./...
make bootstrap-check
```

Risks:

- In-memory behavior leaks into public contract.
- Validation rules are under-specified.
- Concurrency assumptions are unclear.

## M1.4 Registry CLI Probe

Goal: add a minimal developer probe for Registry behavior only if approved.

Allowed directories:

- `cmd/toolvault/`
- `internal/registry/`
- `docs/`

Forbidden directories:

- `internal/gateway/`
- `internal/runtime/`
- `internal/policy/`
- `internal/credential/`
- `internal/observability/`
- `pkg/`

Acceptance criteria:

- CLI is clearly marked as developer-only.
- CLI does not become a public API commitment.
- CLI touches only Registry behavior.
- Tests or documented manual verification exist.

Verification commands:

```sh
go test ./...
make bootstrap-check
```

Risks:

- CLI becomes an accidental product surface.
- CLI pulls in unapproved dependencies.
- CLI creates pressure to add Gateway behavior early.

## M1.5 Registry Review Gate

Goal: review Registry implementation for scope, boundary, dependency, and test
quality before starting Gateway work.

Allowed directories:

- `docs/`
- `internal/registry/`

Forbidden directories:

- `internal/gateway/`
- `internal/runtime/`
- `internal/policy/`
- `internal/credential/`
- `internal/observability/`
- `pkg/`

Acceptance criteria:

- Review documents boundary compliance.
- Review confirms no unapproved dependencies.
- Review confirms tests cover the accepted Registry contract.
- Review lists required human decisions before M2.

Verification commands:

```sh
go test ./...
make bootstrap-check
```

Risks:

- Review misses architecture drift.
- Gateway assumptions appear inside Registry.
- Tests pass while contract remains ambiguous.
