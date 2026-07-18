---
title: Gateway architecture and request flow
type: concept
status: current
scope: repository-operating-docs
related_scopes: []
related_files:
  - cmd/server/main.go
  - internal/web/server.go
  - internal/chathub/client.go
  - internal/auth/cache.go
source_docs:
  - README.md
tags:
  - architecture
  - gateway
  - chathub
last_checked: 2026-07-18
updated: 2026-07-18T08:33:43Z
---

# Gateway architecture and request flow

## Scope

`m365-native` is a local HTTP gateway for authorized Microsoft 365 Copilot sessions. It serves a web console and OpenAI- and Anthropic-compatible APIs; it is not an authentication bypass.

## Request flow

`cmd/server/main.go` loads persisted startup settings, constructs `web.Server`, and starts the HTTP listener. `internal/web/server.go` owns routing, request IDs, security headers, administrator-session checks, API-key validation, conversation handling, and protocol adapters. Requests that need model responses are translated into ChatHub requests by `internal/chathub/`, which maintains the upstream WebSocket interaction. `internal/auth/` stores and refreshes authorized account tokens.

## Boundaries

The web package is the public HTTP boundary. The ChatHub package is the upstream protocol boundary. The auth package owns account-token persistence and refresh. Keep protocol-specific transformations in `internal/web/` rather than leaking HTTP concerns into `internal/chathub/`.

## Verification

Run `go test ./...`, `go vet ./...`, and `go build ./...` after Go changes. Use `docker compose up -d --build` plus a local HTTP smoke test for container changes.

## Related knowledge

See the operations and security pages for container state, deployment boundaries, and credential handling.
