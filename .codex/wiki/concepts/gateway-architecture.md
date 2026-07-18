---
title: Gateway architecture and request flow
type: concept
status: current
scope: repository-operating-docs
related_scopes:
  - codex-responses-compatibility
related_files:
  - cmd/server/main.go
  - internal/web/server.go
  - internal/web/codex_catalog.go
  - internal/web/codex_responses.go
  - internal/web/codex_usage.go
  - internal/web/protocol_handlers.go
  - internal/web/protocol_response.go
  - internal/chathub/client.go
  - internal/auth/cache.go
source_docs:
  - README.md
tags:
  - architecture
  - gateway
  - chathub
  - codex
  - tools
last_checked: 2026-07-18
updated: 2026-07-18T12:10:00Z
---

# Gateway architecture and request flow

## Scope

`m365-native` is a local HTTP gateway for authorized Microsoft 365 Copilot sessions. It serves a web console and OpenAI- and Anthropic-compatible APIs; it is not an authentication bypass.

## Request flow

`cmd/server/main.go` loads persisted startup settings, constructs `web.Server`, and starts the HTTP listener. `internal/web/server.go` owns routing, request IDs, security headers, administrator-session checks, API-key validation, conversation handling, and protocol adapters. Requests that need model responses are translated into ChatHub requests by `internal/chathub/`, which maintains the upstream WebSocket interaction. `internal/auth/` stores and refreshes authorized account tokens.

## Boundaries

The web package is the public HTTP boundary. The ChatHub package is the upstream protocol boundary. The auth package owns account-token persistence and refresh. Keep protocol-specific transformations in `internal/web/` rather than leaking HTTP concerns into `internal/chathub/`.

## Responses compatibility

The OpenAI Responses adapter accepts string input and typed content blocks. In
particular, Codex sends text as `input_text`, which is normalized alongside the
gateway's existing `text` and `output_text` blocks before ChatHub adaptation.
Responses image content may use either a nested `image_url.url` object or a
direct `image_url` string; both forms become ChatHub attachments without
changing the Codex text-block behavior.
The model catalog preserves OpenAI's `data` list and also exposes a `models`
alias required by Codex v0.144.5. Every catalog entry must provide stable
`id`, `slug`, `display_name`, and `supported_reasoning_levels` values; this
gateway uses the model ID for both derived identifier fields. Codex requires
both advertised reasoning fields to contain preset objects with `effort` and
`description`, not bare effort strings; the gateway mirrors that structure and
uses `medium` as its explicit default. The same local catalog also advertises
Codex execution metadata such as `shell_type`, visibility, API availability,
priority, tier lists, tool capabilities, and truncation settings. This metadata
is not forwarded to ChatHub. Codex also requires `base_instructions` and a
matching `model_messages.instructions_template`; the gateway supplies a concise
maintained template rather than copying a bundled local prompt. The catalog
response itself is not forwarded to ChatHub, although Codex may include this
template in later request instructions that use the normal upstream flow.
Administrator settings provide a single mapping contract for client-selected
model IDs: each mapping publishes a compatible catalog record and selects a
validated ChatHub tone for both Chat Completions and Responses calls. The UI
suggests the models currently bundled with Codex, preloads its GPT-5.6 Sol,
Terra, and Luna variants, and allows custom client aliases. It never advertises
the bundled model's local context size over the configured gateway limit.

Codex-specific code is organized inside `internal/web/` as
`codex_catalog.go`, `codex_usage.go`, and `codex_responses.go`. They remain in
the `web` package so they can use the HTTP layer's unexported request and
settings types without introducing exports or a package cycle. Generic OpenAI
and Anthropic handling remains in its protocol files, and ChatHub remains
unaware of Codex-specific metadata.
Streaming Responses calls must end in a terminal `response.completed` or
`response.failed` event; never close after only `response.created`, because
Codex treats that as a transport disconnect. ChatHub does not report provider
token counts, so Responses completion events carry clearly marked local usage
estimates for client context-progress displays rather than billing or quota
decisions. GPT route aliases use the embedded `o200k_base` vocabulary from
`tiktoken-go/tokenizer`; this avoids runtime BPE downloads and identifies the
result as `tiktoken_o200k_base_estimate`. Claude and unknown route aliases
retain a separately identified heuristic fallback. Usage estimates include
visible message framing, tool schemas, tool choice, tool calls, and completion
framing so client context and auto-compaction thresholds are not systematically
low for tool-heavy turns. Hidden upstream prompts, image tokenization, and
provider reasoning tokens remain unavailable.

## Custom tool compatibility

Codex can advertise local executors as Responses `custom` tools rather than
JSON `function` tools. The supported local shell path is `custom: exec`: its
grammar-constrained raw input is represented internally as an `input` string
solely for the ChatHub planner, while the Responses adapter returns a
`custom_tool_call` item and matching `response.custom_tool_call_input.*`
stream events. When Codex returns a `custom_tool_call_output`, the adapter must
also retain the preceding custom call as an assistant tool call so the gateway's
conversation validator can match the result. Do not coerce custom calls into
public `function_call` responses.

`namespace` and hosted tools such as `web_search` are not executed by this
bridge unless a separate compatibility path implements their semantics. Never
silently claim that an unsupported tool executed.

## Verification

Run `go test ./...`, `go vet ./...`, and `go build ./...` after Go changes. Use `docker compose up -d --build` plus a local HTTP smoke test for container changes.

## Related knowledge

See the operations and security pages for container state, deployment boundaries, and credential handling.
