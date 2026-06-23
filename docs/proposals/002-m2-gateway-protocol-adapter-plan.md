# Proposal 002: M2 Gateway And Protocol Adapter Implementation Plan

Status: approved with M2 constraints

Owner role: ToolVault M2 Proposal Architect

Date: 2026-06-23

Human decision date: 2026-06-23

Approved M2 positioning: Registry-backed read-only discovery plus
deterministic projection. M2 is not Gateway productization, not an MCP server,
and not an OpenAI integration.

## 1. M2 Background

M1 delivered the internal Tool Registry baseline under `internal/registry/`.
The Registry now owns `ToolSpec` metadata, lifecycle state, validation, lookup,
listing, update, deletion, an in-memory implementation, and the Registry error
contract.

M2 should build the next minimal layer above Registry: a read-only Gateway and
protocol projection surface that lets registered tools be discovered in a
consistent way and transformed into tool descriptions that agent/protocol
integrations can consume later.

M2 must not turn ToolVault into an Agent Framework or LLM Gateway. It should not
execute tools, route model traffic, manage credentials, evaluate policy, run
sandboxes, or implement a full MCP/OpenAI integration. The useful M2 slice is
discovery and description projection only.

Any M2 implementation attempt to add invocation, authentication, SDK
integration, streaming, admin/debug discovery, or production ingress behavior
must be rejected and redirected to a separate proposal.

## 2. M2 Goals

M2 should complete:

- Define a minimal Gateway interface for read-only tool discovery.
- Consume the approved M1 Registry contract through explicit dependency rules.
- Provide Registry-backed listing and lookup of discoverable tools.
- Add a REST read-only discovery endpoint using the Go standard library.
- Define `ToolSpec` to OpenAI Function tool-description projection.
- Define `ToolSpec` to MCP tool-schema projection as a draft contract.
- Define a small Protocol Adapter interface for projection only.
- Add Gateway and Protocol Adapter tests.
- Add an M2 guardrail that permits approved M2 packages while rejecting scope
  drift into Runtime, Policy, Credential, Observability, persistence, UI,
  execution, and third-party dependencies.

## 3. M2 Non-Goals

M2 must not implement:

- Tool execution or invocation.
- Runtime Manager behavior or runtime selection.
- Policy Engine behavior, authorization, or access control decisions.
- Credential Vault behavior, secret retrieval, or secret injection.
- Sandbox behavior.
- Dashboard, Web UI, marketplace, or admin console.
- Streaming.
- Tool composition or workflow behavior.
- Database persistence or durable Registry storage.
- Full MCP server implementation.
- Full OpenAI SDK integration.
- LLM Gateway behavior, model routing, prompt management, or chat completion
  proxying.
- Observability pipelines, metrics exporters, tracing, or audit sinks.

## 4. Dependency On M1 Registry

M2 depends on the accepted M1 Registry baseline:

- `internal/registry.Registry`
- `internal/registry.ListRequest`
- `internal/registry.GetRequest`
- `internal/registry.ListFilter`
- `internal/registry.Error` and `registry.KindOf`
- `internal/registry/spec.ToolSpec`
- `internal/registry/spec.ToolAction`
- `internal/registry/spec.ToolStatus`

M2 should treat Registry as the authoritative source of tool metadata and
lifecycle status. Gateway must not bypass Registry validation, inspect memory
store internals, or depend on `internal/registry/memory` except in tests or
developer-only wiring.

M2 accepts the M1 Registry contract deviations as the internal M2 input
contract:

- `ToolSpec.Actions` is accepted and required for projection.
- `ToolSpec.Metadata` is accepted but M2 must not depend on it.
- Duplicate `Name` rejection is accepted as useful for OpenAI/MCP name mapping.
- `ExpectedVersion`, `SetStatus`, and current lifecycle rules are accepted as
  part of the Registry baseline, but M2 discovery must not modify Registry
  state.

M2 must not change the Registry implementation. If projection needs fields not
present in M1, M2 must use deterministic placeholder schema behavior or create a
follow-up proposal/task to extend Registry outside M2.

## 5. Module Boundaries

### Registry

Owns:

- Tool metadata and lifecycle state.
- `ToolSpec` validation.
- Tool lookup and listing.
- Registry error kinds and lifecycle rules.

M2 must not modify Registry implementation behavior.

### Gateway

Owns:

- Read-only discovery use cases over Registry.
- Gateway request/response contracts for discovery.
- Mapping Registry errors into Gateway errors and read-only HTTP responses.
- Discovery filters approved for M2.
- REST read-only discovery handlers.

Gateway does not own tool execution, runtime decisions, policy decisions,
credentials, observability sinks, protocol SDK integrations, or persistence.

### Protocol Projection / Adapter

M2 should keep protocol projection as Gateway-owned internal packages, not as a
new top-level v1 core module. The adapter layer owns:

- Projection of Registry `ToolSpec` into OpenAI Function-style descriptions.
- Projection of Registry `ToolSpec` into MCP tool-schema draft descriptions.
- Projection-specific validation of unsupported shapes.
- Stable deterministic output for tests.

The adapter layer must not call external SDKs, start protocol servers, execute
tools, or perform network I/O.

## 6. Allowed Modification Directories

Allowed for this proposal:

- `docs/proposals/`

Recommended for approved M2 implementation tasks:

- `internal/gateway/`
- `docs/`
- `scripts/`
- `Makefile`

Conditionally allowed with explicit task scope:

- `docs/acceptance/` for M2 acceptance checklist.
- `docs/releases/` for M2 release notes.
- `docs/retrospectives/` for M2 retrospective.

Deferred by default:

- `cmd/toolvault/` developer-only wiring. Add it only after explicit human
  request for a local manual probe, and label it developer-only, unauthenticated,
  and not a production entrypoint.

M2 should not require modifying `internal/registry/`. If tests need Registry
fixtures, they should use fakes or `internal/registry/memory` from test code
without changing Registry code.

## 7. Forbidden Modification Directories

Forbidden during M2 unless a separate human-approved proposal expands scope:

- `internal/registry/` implementation files.
- `internal/runtime/`
- `internal/policy/`
- `internal/credential/`
- `internal/observability/`
- `pkg/`
- Web UI directories or frontend application directories.
- Database migration directories or persistence packages.
- Vendor directories or dependency lockfiles created by new third-party
  dependencies.

## 8. Recommended Directory Structure

Recommended M2 implementation layout:

```text
internal/
  gateway/
    gateway.go
    errors.go
    discovery.go
    discovery_test.go
    http/
      handler.go
      handler_test.go
    protocol/
      adapter.go
      openai/
        projection.go
        projection_test.go
      mcp/
        projection.go
        projection_test.go
```

Rationale:

- `internal/gateway/` keeps M2 inside the existing Gateway core boundary.
- `internal/gateway/http/` isolates REST transport from Gateway use cases.
- `internal/gateway/protocol/` keeps protocol projection internal to Gateway
  for M2 and avoids adding a new top-level core module.
- `openai` and `mcp` subpackages are projection-only packages, not SDK
  integrations and not protocol servers.

If humans prefer Protocol Adapter to become a separate top-level module later,
that should be a separate boundary proposal.

Approved M2 decision: do not create `internal/protocol/`. Protocol projection
belongs under `internal/gateway/protocol/` for M2.

## 9. Public Interface Draft

This section is a contract draft for implementation review. It is not business
code and should be refined during M2 task approval.

### Gateway Discovery

```go
type Gateway interface {
    ListTools(ctx context.Context, req ListToolsRequest) ([]DiscoveredTool, error)
    GetTool(ctx context.Context, req GetToolRequest) (DiscoveredTool, error)
}

type ListToolsRequest struct {
    Filter DiscoveryFilter
}

type GetToolRequest struct {
    ID spec.ToolID
}

type DiscoveryFilter struct {
    Status spec.ToolStatus
    Tags   []string
}

type DiscoveredTool struct {
    Spec spec.ToolSpec
}
```

M2 discovery defaults:

- Empty filter lists only `active` tools.
- Explicit `Status` may request `active` or `deprecated` only for ToolVault's
  own discovery JSON.
- `draft` and `disabled` are never exposed by M2 Gateway.
- M2 must not add debug/admin discovery mode.
- OpenAI and MCP projection endpoints project only `active` tools.
- `Tags` are passed through to Registry `ListFilter`.

### Gateway Error Categories

```go
type ErrorKind string

const (
    InvalidRequest ErrorKind = "invalid_request"
    ToolNotFound   ErrorKind = "tool_not_found"
    RegistryFailed ErrorKind = "registry_failed"
)
```

Gateway should preserve Registry error meaning without leaking storage details:

- Registry `NotFound` maps to `ToolNotFound`.
- Registry `InvalidSpec` from invalid filters maps to `InvalidRequest`.
- Other Registry errors map to `RegistryFailed`.

### Protocol Adapter

```go
type Adapter interface {
    Name() string
    ProjectTool(ctx context.Context, tool spec.ToolSpec) (ToolDescription, error)
    ProjectTools(ctx context.Context, tools []spec.ToolSpec) ([]ToolDescription, error)
}

type ToolDescription struct {
    Protocol string
    Name     string
    Payload  any
}
```

Adapter rules:

- Adapters accept `ToolSpec` only; they do not query Registry directly.
- Adapters produce descriptions only; they never execute or invoke tools.
- Output order is deterministic and follows Gateway/Registry list order.
- Unsupported fields should return projection errors, not trigger Registry
  changes during M2.

### OpenAI Function Projection Draft

M2 should map each `ToolAction` to one OpenAI-style tool description. Because
M1 `ToolAction` has no input parameter schema, M2 should use an empty object
schema until a future Registry schema proposal is approved.

```go
type OpenAITool struct {
    Type     string         `json:"type"`
    Function OpenAIFunction `json:"function"`
}

type OpenAIFunction struct {
    Name        string         `json:"name"`
    Description string         `json:"description"`
    Parameters  map[string]any `json:"parameters"`
}
```

Projection rules:

- Only `active` tools are projected.
- `Type` is always `function`.
- Function `Name` should be deterministic. Recommended format:
  `<tool.Name>__<action.Name>` when a tool has multiple actions, and
  `<tool.Name>` only when the tool has exactly one action whose name matches the
  tool name.
- `Description` uses action description, falling back to tool description only
  if the action description is empty. Current M1 validation requires action
  descriptions, so fallback should be defensive only.
- `Parameters` is:

```json
{"type":"object","properties":{},"additionalProperties":false}
```

M2 must not add OpenAI SDK imports.

M2 must not change Registry to add input schemas for this projection. The empty
object schema is an explicit M2 limitation.

### MCP Tool Schema Projection Draft

M2 should map each `ToolAction` to one MCP-style tool description draft. Because
M2 is not a full MCP server, this projection is a serializable shape compatible
with future MCP tool registration work, not an MCP runtime implementation.

```go
type MCPTool struct {
    Name        string         `json:"name"`
    Description string         `json:"description"`
    InputSchema map[string]any `json:"inputSchema"`
}
```

Projection rules:

- Only `active` tools are projected.
- `Name` uses the same deterministic naming rule as OpenAI projection.
- `Description` uses action description.
- `InputSchema` uses the same empty object schema as OpenAI projection.
- No MCP server lifecycle, request handling, session management, or SDK import
  is included in M2.

### REST Read-Only Discovery Endpoint

M2 should provide standard-library `net/http` handlers only.

Approved endpoints:

```text
GET /v1/tools
GET /v1/tools/{id}
GET /v1/tools:openai
GET /v1/tools:mcp
```

Endpoint behavior:

- `GET /v1/tools` returns Gateway discovery descriptions in ToolVault's own
  JSON shape. It may return `deprecated` tools only through an explicit
  `?status=deprecated` filter.
- `GET /v1/tools/{id}` returns one discoverable tool by Registry ID.
- `GET /v1/tools:openai` returns OpenAI projection for active tools only.
- `GET /v1/tools:mcp` returns MCP projection for active tools only.
- Endpoints are read-only.
- No execution endpoint is added.
- No authentication or policy behavior is implemented in M2. M2 handlers are
  library handlers and tests first, not a production ingress surface.
- No admin/debug discovery endpoint or mode is added.

HTTP status mapping:

- `200` for successful discovery.
- `400` for invalid filters or unsupported projection requests.
- `404` for unknown tool IDs or non-discoverable tool IDs.
- `405` for non-GET methods.
- `500` for unexpected Registry/Gateway failures.

## 10. Subtask Breakdown

### M2.0 Gateway Proposal Approval

Goal: approve M2 boundaries, Registry dependency direction, REST discovery
surface, and projection-only adapter scope.

Allowed modification directories:

- `docs/proposals/`

Forbidden directories:

- `cmd/`
- `internal/`
- `pkg/`
- `scripts/`
- `Makefile`

Acceptance criteria:

- M2 goals and non-goals are explicit.
- Registry dependency on M1 is documented.
- Protocol projection is scoped to description generation only.
- REST discovery is read-only and excludes execution.
- Human decisions are recorded before implementation starts.

Verification commands:

```sh
make bootstrap-check
```

### M2.1 M2 Guardrail

Goal: add a verification gate for approved M2 work before Gateway business code
is implemented.

Allowed modification directories:

- `scripts/`
- `Makefile`
- `docs/`

Forbidden directories:

- `cmd/toolvault/`
- `internal/registry/`
- `internal/gateway/`
- `internal/runtime/`
- `internal/policy/`
- `internal/credential/`
- `internal/observability/`
- `pkg/`

Acceptance criteria:

- Existing `make bootstrap-check`, `make m1-registry-check`, and `make verify`
  remain valid.
- A new M2 verification command is documented or added, for example
  `make m2-gateway-check`.
- The M2 gate permits approved `internal/gateway/` files after this subtask.
- The M2 gate rejects Runtime, Policy, Credential, Observability, Web UI,
  database persistence, dependency drift, SDK integration, streaming, debug/admin
  discovery, and tool execution endpoints.
- The gate rejects new third-party dependencies unless a human approval record
  exists.

Verification commands:

```sh
make bootstrap-check
make m1-registry-check
make verify
```

After the gate is added:

```sh
make m2-gateway-check
make verify
```

### M2.2 Gateway Interface And Error Contract

Goal: define the Gateway discovery interface, request/response types, and error
mapping.

Allowed modification directories:

- `internal/gateway/`
- `docs/`

Forbidden directories:

- `internal/registry/` implementation files.
- `internal/runtime/`
- `internal/policy/`
- `internal/credential/`
- `internal/observability/`
- `pkg/`
- `cmd/toolvault/`

Acceptance criteria:

- Gateway exposes read-only `ListTools` and `GetTool` behavior.
- Gateway depends on Registry through `internal/registry.Registry`, not memory
  internals.
- Default discovery scope is documented and tested.
- Gateway error kinds are testable.
- No REST handler, execution behavior, policy behavior, credential behavior, or
  protocol projection is implemented in this subtask.
- No debug/admin discovery mode is added.

Verification commands:

```sh
go test ./internal/gateway/...
make m2-gateway-check
```

### M2.3 Registry-Backed Discovery Implementation

Goal: implement Gateway discovery over the M1 Registry interface.

Allowed modification directories:

- `internal/gateway/`
- `docs/`

Forbidden directories:

- `internal/registry/` implementation files.
- `internal/runtime/`
- `internal/policy/`
- `internal/credential/`
- `internal/observability/`
- `pkg/`

Acceptance criteria:

- Listing uses Registry `List`.
- Lookup uses Registry `Get`.
- Default discovery exposes only active tools.
- Explicit filtering handles tags plus `active` or `deprecated` status for
  ToolVault-owned discovery JSON.
- Draft and disabled tools are never exposed by M2 Gateway.
- Registry errors are mapped into Gateway errors.
- Tests use fakes or Registry memory from test code only.

Verification commands:

```sh
go test ./internal/gateway/...
make m2-gateway-check
make verify
```

### M2.4 Protocol Adapter Interface

Goal: define projection-only adapter interfaces and shared test helpers.

Allowed modification directories:

- `internal/gateway/protocol/`
- `docs/`

Forbidden directories:

- `internal/registry/` implementation files.
- `internal/runtime/`
- `internal/policy/`
- `internal/credential/`
- `internal/observability/`
- `pkg/`

Acceptance criteria:

- Adapter interface accepts `ToolSpec` inputs and returns serializable
  descriptions.
- Adapter interface does not import external SDKs.
- Adapter interface does not query Registry directly.
- Adapter interface has deterministic order semantics.
- Tests cover empty list, single tool, multiple actions, unsupported status, and
  projection errors.

Verification commands:

```sh
go test ./internal/gateway/...
make m2-gateway-check
```

### M2.5 OpenAI Function Projection

Goal: project discoverable `ToolSpec` actions into OpenAI Function-style tool
descriptions.

Allowed modification directories:

- `internal/gateway/protocol/openai/`
- `internal/gateway/protocol/`
- `docs/`

Forbidden directories:

- `internal/registry/` implementation files.
- `internal/runtime/`
- `internal/policy/`
- `internal/credential/`
- `internal/observability/`
- `pkg/`

Acceptance criteria:

- Each `ToolAction` maps to one OpenAI-style function tool description.
- Names are deterministic and valid for current M1 machine-name constraints.
- Descriptions are populated from action descriptions.
- Parameter schema is an empty JSON object schema until ToolSpec input schemas
  are approved.
- No OpenAI SDK dependency is introduced.
- Tests cover one-action tools, multi-action tools, deprecated/draft/disabled
  rejection, and deterministic JSON output.

Verification commands:

```sh
go test ./internal/gateway/protocol/...
make m2-gateway-check
```

### M2.6 MCP Tool Schema Projection Draft

Goal: project discoverable `ToolSpec` actions into an MCP-style tool-schema
draft without implementing an MCP server.

Allowed modification directories:

- `internal/gateway/protocol/mcp/`
- `internal/gateway/protocol/`
- `docs/`

Forbidden directories:

- `internal/registry/` implementation files.
- `internal/runtime/`
- `internal/policy/`
- `internal/credential/`
- `internal/observability/`
- `pkg/`

Acceptance criteria:

- Each `ToolAction` maps to one MCP-style tool schema.
- Output shape includes name, description, and input schema.
- Input schema is the same empty JSON object schema used by OpenAI projection.
- No MCP SDK dependency or MCP server is introduced.
- Tests cover one-action tools, multi-action tools, deprecated/draft/disabled
  rejection, and deterministic JSON output.

Verification commands:

```sh
go test ./internal/gateway/protocol/...
make m2-gateway-check
```

### M2.7 REST Read-Only Discovery Endpoint

Goal: expose Gateway discovery and protocol projections through read-only
standard-library HTTP handlers.

Allowed modification directories:

- `internal/gateway/http/`
- `internal/gateway/`
- `docs/`

Forbidden directories:

- `internal/registry/` implementation files.
- `internal/runtime/`
- `internal/policy/`
- `internal/credential/`
- `internal/observability/`
- `pkg/`

Acceptance criteria:

- Standard-library `net/http` is sufficient.
- Handlers expose only approved GET discovery endpoints.
- Non-GET methods return `405`.
- Unknown tools return `404`.
- Invalid filters return `400`.
- Successful responses are deterministic JSON.
- No execution endpoint is added.
- Deprecated tools are available only from ToolVault discovery JSON through
  explicit `?status=deprecated`.
- OpenAI/MCP endpoints project active tools only.
- No auth, policy, credential, observability, streaming, or SDK behavior is
  added.

Verification commands:

```sh
go test ./internal/gateway/...
make m2-gateway-check
make verify
```

### M2.8 M2 Acceptance And Release Documentation

Goal: close M2 with acceptance, release notes, and residual-risk documentation.

Allowed modification directories:

- `docs/acceptance/`
- `docs/releases/`
- `docs/retrospectives/`
- `README.md`

Forbidden directories:

- `internal/registry/` implementation files.
- `internal/runtime/`
- `internal/policy/`
- `internal/credential/`
- `internal/observability/`
- `pkg/`

Acceptance criteria:

- Acceptance checklist records completed M2 capabilities and out-of-scope
  exclusions.
- Release notes explain that M2 is discovery/projection only.
- Retrospective records deviations, guardrail gaps, and review outcomes.
- README reflects M2 capabilities without implying execution or production
  ingress readiness.
- Residual risks are listed before entering M3 planning.

Verification commands:

```sh
make m2-gateway-check
make verify
```

## 11. Risk And Dependency Management

### Registry Contract Drift

Risk: M2 relies on M1 fields and behaviors that exceeded Proposal 001.

Management:

- Treat the accepted M1 deviations as M2's internal input contract.
- Do not modify Registry in M2 unless a separate proposal approves it.
- Keep Gateway tests focused on the accepted Registry interface.

### Parameter Schema Gap

Risk: M1 `ToolAction` has no input parameter schema, so OpenAI/MCP projections
cannot express real tool arguments.

Management:

- Use an empty object schema in M2.
- Document this as a known limitation.
- Do not modify Registry for schemas in M2.
- Defer richer input schema to a future Registry schema proposal.

### Protocol Scope Creep

Risk: Projection work may drift into full MCP server behavior or OpenAI SDK
integration.

Management:

- Keep adapter packages projection-only.
- Reject SDK imports in the M2 guardrail.
- Require tests that prove no execution endpoint exists.

### Gateway Scope Creep

Risk: Read-only Gateway may drift into invocation, runtime, policy, credential,
or observability behavior.

Management:

- Add M2 guardrail before Gateway implementation.
- Keep endpoints read-only.
- Map errors without adding auth/policy decisions.
- Reject admin/debug discovery modes in M2.
- Defer invocation to a later Gateway/Runtime proposal.

### Network Surface Risk

Risk: Even read-only HTTP handlers introduce a network protocol surface.

Management:

- Use `net/http` handlers without starting a production server by default.
- Defer `cmd/` server wiring unless explicitly requested for a developer-only
  manual probe.
- Document the lack of auth/policy as a non-production limitation.

### Dependency Drift

Risk: SDKs, routers, OpenAPI generators, or JSON schema packages may be added
prematurely.

Management:

- Use the Go standard library only.
- Require explicit human approval for any third-party dependency.
- Keep JSON schema payloads as simple typed structs/maps in M2.

## 12. Third-Party Dependencies

M2 should not add third-party dependencies.

Approved implementation should use:

- `context`
- `encoding/json`
- `errors`
- `fmt`
- `net/http`
- `sort`
- `strings`
- `testing`
- other Go standard-library packages as needed

Do not add:

- OpenAI SDKs.
- MCP SDKs.
- HTTP routers/frameworks.
- JSON schema libraries.
- Database drivers.
- Observability libraries.
- Web UI toolchains.

Any dependency addition requires a separate human approval decision before code
changes.

## 13. Resolved Human Decisions

Human decisions accepted for M2:

1. M1 Registry contract deviations are accepted as the M2 internal input
   contract. `Actions` is required for projection, duplicate `Name` rejection is
   useful for protocol naming, and `Metadata` may remain but M2 does not depend
   on it. M2 must not modify Registry.
2. Default discovery exposes only `active` tools.
3. `deprecated` tools may be returned only through ToolVault-owned discovery
   JSON with explicit `?status=deprecated`.
4. `draft` and `disabled` tools are not exposed by M2 Gateway.
5. OpenAI/MCP projections include only `active` tools and never project
   `deprecated`, `draft`, or `disabled`.
6. Protocol adapter packages live under `internal/gateway/protocol/`; M2 must
   not create a top-level `internal/protocol/`.
7. Endpoint names are approved:
   `GET /v1/tools`, `GET /v1/tools/{id}`, `GET /v1/tools:openai`, and
   `GET /v1/tools:mcp`.
8. Empty object input schemas are accepted as an explicit M2 limitation:
   `{"type":"object","properties":{},"additionalProperties":false}`.
9. `cmd/toolvault/` developer-only wiring is excluded by default. It may be
   added only after explicit human request for a local manual probe.
10. M2 guardrail may modify `Makefile` and `scripts/` and should add
    `make m2-gateway-check`.

No open M2 direction decision remains in this proposal. Implementation tasks
must still go through review and must preserve the approved M2 scope.

## 14. Alternatives Considered

### Build Full MCP Server In M2

Rejected for M2. A full MCP server adds protocol lifecycle, session handling,
transport decisions, and compatibility obligations. M2 should only define a
draft schema projection.

### Integrate OpenAI SDK In M2

Rejected for M2. SDK integration would add third-party dependency and product
surface area before ToolVault has execution, policy, credential, or runtime
boundaries.

### Add Protocol Adapter As A Top-Level Core Module

Deferred. ToolVault positioning includes protocol adaptation, but the current
v1 core module list names Gateway, Runtime, Policy, Credential, Registry, and
Observability. Keeping projection under Gateway in M2 avoids expanding module
topology before there is a stronger reason.

### Extend Registry With Input Schemas Before M2

Deferred. Input schemas are useful, but changing `ToolSpec` is a Registry
contract change. M2 can still provide deterministic discovery/projection with
empty object schemas and record the limitation.

### Add Auth Or Policy To Discovery Endpoints

Rejected for M2. Policy Engine is a later v1 module and should not be
implemented implicitly inside Gateway.

## 15. Affected Modules

Affected by approved M2 implementation:

- `internal/gateway/`
- `docs/`
- `scripts/`
- `Makefile`

Consumed but not modified:

- `internal/registry/`

Explicitly unaffected:

- `internal/runtime/`
- `internal/policy/`
- `internal/credential/`
- `internal/observability/`
- `pkg/`

## 16. Test Impact

M2 should add tests for:

- Gateway interface contracts and error mapping.
- Registry-backed discovery list/get behavior.
- Discovery filtering and default active-only visibility.
- Non-discoverable status behavior.
- Protocol adapter deterministic output.
- OpenAI projection for single-action and multi-action tools.
- MCP projection for single-action and multi-action tools.
- HTTP handler response codes, methods, routes, and JSON bodies.
- Guardrail checks for forbidden directories, dependencies, execution endpoints,
  persistence, SDK integration, streaming, admin/debug discovery, and UI drift.

M2 should keep existing M1 tests passing.

## 17. Rollback Plan

If M2 implementation introduces boundary or dependency problems:

- Revert `internal/gateway/` implementation changes from the failing subtask.
- Keep this proposal and review notes for follow-up decisions.
- Preserve M1 Registry files and tests unchanged.
- Re-run `make m1-registry-check` and `make verify` to confirm the repository
  returns to the M1 baseline.
- If guardrail changes are wrong, revert the M2 guardrail task and keep M1
  verification commands intact.

## 18. Approval Record

This M2 implementation plan is approved with the constraints recorded above.

Approval means:

- M2 may create `internal/gateway/` implementation tasks.
- M2 may add a standard-library REST read-only discovery handler.
- M2 may add Gateway-owned protocol projection packages under
  `internal/gateway/protocol/`.
- M2 may add an M2 verification gate in `Makefile` and `scripts/`.
- M2 remains discovery/projection only and does not implement execution,
  Runtime, Policy, Credential, Observability, persistence, UI, SDK integration,
  streaming, admin/debug discovery, or full MCP server behavior.
