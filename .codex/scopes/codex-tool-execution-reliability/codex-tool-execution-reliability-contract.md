---
description: Contract for codex-tool-execution-reliability.
---

# codex-tool-execution-reliability Contract

## Context

- Current worktree has unrelated local wiki/scope changes; preserve them.
- Codex sends local shell work through the Responses `custom: exec` contract.
- Relevant source paths: `internal/web/protocol_compat.go`,
  `internal/web/prompt.go`, `internal/web/protocol_handlers.go`,
  `internal/web/codex_responses.go`, and focused tests in `internal/web/`.
- Related active scopes: `codex-custom-tools` established the custom-exec
  boundary; `codex-responses-compatibility` owns the wider Responses adapter.

## Findings

- `responsesRequest` does not model the Responses top-level `instructions`
  field, so JSON decoding discards Codex's execution/workspace guidance before
  the request is adapted to ChatHub.
- The stream adapter recognizes a custom call but emits
  `response.function_call_arguments.delta` while assembling it. Codex expects
  `response.custom_tool_call_input.*` for `custom: exec`.
- A router can validly return an empty auto-mode plan, so the remaining model
  response must be constrained by the forwarded instructions rather than
  assumed to be an executable shell action.

## Outcome

- Done when Codex Responses instructions reach the ChatHub prompt, custom exec
  calls retain custom-only stream events, and focused regression tests cover
  both behavior.
- User-visible/runtime state: workspace tasks receive a structured custom exec
  call or an ordinary answer; a guessed absolute `/root/...` path is not added
  by gateway guidance and command text is not mislabelled as execution.
- Durable knowledge to preserve: model catalog instructions are client metadata;
  request `instructions` and executable-call event types are the enforcement
  path at the gateway boundary.

## Goals / Non-goals

Goals:
- Preserve Responses `instructions` as a system-level prompt component.
- Add a narrow custom-exec policy that uses the caller workspace and relative
  paths without granting execution authority.
- Emit only custom Responses argument events for custom calls in streaming mode.
- Add focused regression tests and run the Go validation matrix.

Non-goals:
- Alter Codex approval/sandbox policy, host permissions, or upstream ChatHub
  authentication.
- Force a tool call for every auto-mode user request or emulate unsupported tools.

## Target files / modules

- `internal/web/protocol_compat.go`
- `internal/web/protocol_handlers.go`
- focused `internal/web/*_test.go`
- `.codex/wiki/concepts/gateway-architecture.md`

## Constraints

- Preserve OpenAI/Anthropic compatibility and do not log raw instructions,
  command input, tool output, tokens, or account state.
- Keep the gateway tool-agnostic except for the existing Codex `custom: exec`
  compatibility path.

## Boundaries

Allowed changes:
- Responses request normalization, streaming event serialization, focused tests,
  local scope evidence, and aligned local wiki notes.

Forbidden changes:
- Broad execution-policy changes, Docker/Compose changes, permission changes,
  credential/config-file edits, or unsupported upstream tool emulation.

## Decision Summary

| Decision | Evidence Source | Evidence Strength | Conflict | Result | Confidence Reason |
| --- | --- | --- | --- | --- | --- |
| Preserve top-level Responses instructions | `responsesRequest` source and user reproduction | High | resolved | Prepend a system message during normalization | Without this, Codex workspace guidance is silently discarded. |
| Add custom-exec workspace policy | Existing custom-exec contract and user symptom | High | resolved | Add a narrow system policy only when `custom: exec` is present | Avoids guessing inaccessible absolute paths without altering other tools. |
| Emit custom-only streaming events | `streamResponsesAdapter` source and custom-tool contract | High | resolved | Suppress function-argument events and finalize the original custom item | Codex dispatch relies on custom input event names and stable item identity. |

## Verification surface

- `gofmt -w` on modified Go files
- focused `go test ./internal/web`
- `go test ./...`, `go vet ./...`, and `go build ./...`
- `git diff --check`, scope scans, wiki rebuild/lint/doctor

## Escalation triggers

- Escalate only when code/runtime evidence, authoritative wiki, and scope docs materially conflict and the conflict cannot be resolved from local evidence.
- Escalate for data deletion, permission semantics, production access model, or public API compatibility decisions outside the stated boundaries.
- Escalate when user-specified boundaries cannot be satisfied together.

## Rollback

- Revert the narrow source changes with a normal Git revert, then rebuild the
  service. No data migrations or persistent protocol state are introduced.

## Open questions

- None.

## Execution log / evidence updates

- 2026-07-19: scope created from a reproducible protocol-path diagnosis.
- 2026-07-19: added `responsesRequest.Instructions` and preserved it as a
  system message. The initial test exposed that string input replaced earlier
  system messages; the input normalization now appends instead.
- 2026-07-19: custom exec calls now suppress function-argument SSE deltas and
  finish the original custom output item using custom input events only.
- 2026-07-19: `gofmt`, focused web tests, full Go tests, vet, and build passed.
- 2026-07-19: wiki rebuild/lint/doctor and scope placeholder, text, and sync
  checks passed. Authenticated Codex end-to-end execution remains a caller-side
  validation because this scope does not inspect or change credentials.
