# ToolVault Scope

## v1 Core

ToolVault v1 includes only the following core modules:

- Tool Registry
- Gateway
- Runtime Manager
- Policy Engine
- Credential Vault
- Observability

## Tool Registry

Owns tool metadata, lifecycle state, and lookup behavior.

Initial M1 work is limited to the Registry. See
`docs/06-m1-registry-tasks.md`.

## Gateway

Owns the controlled entrypoint for tool invocation.

M0 does not implement Gateway behavior.

## Runtime Manager

Owns execution lifecycle decisions and runtime governance.

M0 does not implement Runtime Manager behavior.

## Policy Engine

Owns policy evaluation before sensitive tool operations.

M0 does not implement Policy Engine behavior.

## Credential Vault

Owns credential reference, retrieval, and injection boundaries.

M0 does not implement Credential Vault behavior.

## Observability

Owns logging, audit events, metrics, and traces for tool operations.

M0 does not implement Observability behavior.

## Future Plugin Areas

These are explicitly out of v1 unless approved:

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

## Scope Control

Any task that adds a new module, expands v1 scope, introduces new persistence,
adds UI, introduces distributed deployment, or changes module ownership must be
captured as a Proposal under `docs/proposals/` and approved by a human.
