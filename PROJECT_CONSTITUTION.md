# ToolVault Project Constitution

## Positioning

ToolVault is the Tool Control Plane for AI Agents.

It provides unified tool registration, protocol adaptation, runtime governance, policy enforcement, credential injection, and observability.

## v1 Core

- Tool Registry
- Gateway
- Runtime Manager
- Policy Engine
- Credential Vault
- Observability

## Plugins / Future

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

## Technical Constraints

- Go core
- Modular monolith first
- No Web UI in M0/M1
- No heavy dependencies in M0/M1 without explicit approval
- No microservices in v1 unless approved
- No database required in M0
- Prefer interface-first design

## Governance Rules

- Human owns project positioning and v1 scope.
- Codex may propose changes but must not apply scope changes without approval.
- Codex Main Agent may create docs, skills, milestones, task graph, CI, scripts, and bootstrap code.
- Builder Agents may only modify explicitly assigned directories.
- Reviewer Agents must check boundary violations, dependency creep, missing tests, and architecture drift.
- Every task must include acceptance criteria and verification commands.