---
description: Bridge Codex custom exec tools through the M365 Responses adapter
---

# codex-custom-tools Contract

## Context

- Initial investigation baseline: local `main` at `35a4cb0`; the repair was
  subsequently integrated into `main`.
- Relevant source paths: `internal/web/protocol_compat.go`, `protocol_handlers.go`,
  `protocol_response.go`, `toolloop.go`, and `tool_response.go`.
- Relevant archived scope references: none.

## Findings

- Codex sends its local shell executor as `tools[].type = "custom"`,
  `name = "exec"`, with grammar input. The prior Responses adapter retained
  only `function` tools, silently dropping `exec` before ChatHub planning.
- A custom call must return a `custom_tool_call` output item and
  `response.custom_tool_call_input.*` stream events. The continuation also
  includes `custom_tool_call` plus `custom_tool_call_output`; both must map to
  an assistant call and a matching tool result for conversation validation.

## Outcome

- Done when: `custom: exec` is converted to an internal tool without losing
  its type; Responses output uses the custom-call shape; and a Codex CLI turn
  executes and completes a read-only command.
- User-visible/runtime state: `codex exec --profile gpt1` executed
  `/bin/bash -lc 'uname -s'` in the native worktree and completed with `Linux`.
- Durable knowledge to preserve: do not treat custom tools as JSON functions
  at the Responses boundary. Bridge their raw `input` separately.

## Goals / Non-goals

Goals:
- Support the Codex `custom: exec` request and continuation loop.
- Preserve existing standard function-call behavior.

Non-goals:
- Implementing `namespace` collaboration tools or hosted `web_search`.
- Changing Codex permissions, shell policy, credentials, or upstream auth.

## Target files / modules

- `internal/web/protocol_compat.go`
- `internal/web/protocol_handlers.go`
- `internal/web/protocol_response.go`
- tool-call conversion helpers and focused tests in `internal/web/`

## Constraints

- Keep the bridge tool-agnostic outside the explicitly supported `exec` custom
  tool and preserve the client-owned approval/sandbox model.
- Do not log raw prompts, command input, OAuth tokens, API keys, or tool output.

## Boundaries

Allowed changes:
- Responses adapter, internal tool-call representation, and focused tests.

Forbidden changes:
- Container privilege changes, broad execution policy changes, auth changes,
  or unsupported tool emulation.

## Decision Summary

| Decision | Evidence Source | Evidence Strength | Conflict | Result | Confidence Reason |
| --- | --- | --- | --- | --- | --- |
| Bridge `exec` as a custom tool with a JSON-string planning shim | Codex request metadata, current source, official Responses streaming schema | Direct | resolved | `custom_tool_call` output with raw `input` | Required for the client executor and verified end-to-end |

## Verification surface

- `go test ./...`
- `go vet ./...`
- `go build ./...`
- `docker compose config --quiet`
- `docker compose up -d --build`, `docker compose ps`, and local HTTP smoke
- `codex exec` read-only `uname -s` integration test via profile `gpt1`

## Escalation triggers

- Escalate only when code/runtime evidence, authoritative wiki, and scope docs materially conflict and the conflict cannot be resolved from local evidence.
- Escalate for data deletion, permission semantics, production access model, or public API compatibility decisions outside the stated boundaries.
- Escalate when user-specified boundaries cannot be satisfied together.

## Rollback

- Select a reviewed known-good commit, reset only with explicit authorization,
  and rebuild the Compose service.

## Open questions

- None.

## Execution log / evidence updates

- Implemented the custom-tool bridge and covered conversion, custom output
  serialization, streaming events, and continuation validation with unit tests.
- Integration evidence: first call emitted a custom tool call, Codex executed
  it locally, and the subsequent custom tool output completed without reconnect.
