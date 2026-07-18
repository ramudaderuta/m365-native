# m365-native Agent Guide

## Source of truth

- `README.md` is the operator-facing entry point.
- This file defines repository-work rules for coding agents.
- `.codex/wiki/index.md` is the detailed, durable project knowledge base. Read
  the relevant current wiki page before changing a covered subsystem, and keep
  it aligned when a change affects architecture, operations, security, or a
  reusable debugging finding.
- Current source and tests override stale wiki or scope notes. Do not treat
  archived scope records as implementation instructions.

## Architecture map

- `cmd/server/main.go` loads persisted startup settings, constructs the web
  server, and owns process startup and the `M365_LISTEN` bind address.
- `internal/web/` owns HTTP routes, web-console assets, admin sessions,
  API-key checks, protocol compatibility adapters, conversation persistence,
  and runtime settings.
- `internal/auth/` owns PKCE/device authorization flows, account-token cache,
  and token refresh. Its data is credential-bearing.
- `internal/chathub/` owns the upstream ChatHub WebSocket protocol, event
  parsing, attachments, images, and tool payloads. Keep it independent of HTTP
  request/response concerns.
- `web/` contains runtime HTML assets. The final Docker image must include it
  at `/app/web` because `internal/web/security_http.go` reads it at runtime.

## Configuration and security

- Treat account caches, OAuth access/refresh tokens, cookies, API keys,
  administrator passwords, HAR files, and debug logs as secrets. Never print,
  commit, copy, or place them in wiki or scope documents.
- `data/` and `secrets/` are local runtime state. Preserve restrictive file
  permissions and the read-only password-file mount in Compose.
- The tracked Compose configuration binds to localhost. LAN publishing is a
  deliberate local deployment choice using the ignored
  `docker-compose.override.yml`; do not change the tracked default for a
  workstation-specific preference.
- Keep external exposure behind TLS and an access-control layer. Do not weaken
  authentication, rate limiting, HTTP security headers, or API-key checks to
  make a workflow pass.

## Working rules

- Inspect `git status --short` before and after edits. Preserve unrelated user
  changes and do not commit, push, rebase, or reset unless explicitly asked.
- When an authorized upstream synchronization is needed, first inspect the
  current branch, remotes, working tree, and divergence after `git fetch
  --prune`. Review incoming commits and their diff before integrating them.
  Do not auto-stash user work or use reset, rebase, force-push, or a forced
  merge to make the sync succeed.
- Integrate only the explicitly selected remote branch, preferring a
  fast-forward update. If local work or a conflict prevents that update, stop
  and preserve the worktree until the conflict is reviewed and resolved.
  Re-run the validation appropriate to the combined change before an
  authorized normal push. Confirm the target remote, branch, active GitHub
  account, and write access immediately before pushing.
- Use `apply_patch` for deliberate text edits. Prefer narrow changes with
  targeted tests.
- Run `gofmt -w` on modified Go files. Keep Go code in its owning package;
  route-specific compatibility behavior belongs in `internal/web/`.
- When Docker behavior changes, build the image and verify both the service
  listener and an HTTP response. Do not assume a started container means the
  web assets are present.
- For durable knowledge changes, use the wiki workflow: search existing pages,
  update or add typed pages without secrets, then run `rebuild`, `lint`, and
  `doctor`. Run `surface-check` after changing README or this file.

## Validation matrix

| Change type | Required validation |
| --- | --- |
| Go source | `gofmt -w <changed files>`, `go test ./...`, `go vet ./...`, `go build ./...` |
| HTTP/auth/protocol behavior | Relevant package tests plus an authorized and an unauthorized-path check where safe |
| Docker or Compose | `docker compose config --quiet`, `docker compose up -d --build`, `docker compose ps`, listener check, and local HTTP smoke test |
| Documentation or wiki | `git diff --check`, wiki `rebuild`, `lint`, `doctor --json`, and `surface-check --json` when README/AGENTS changes |

## Useful wiki pages

- `.codex/wiki/concepts/gateway-architecture.md`
- `.codex/wiki/how-to/container-operation.md`
- `.codex/wiki/reference/credential-and-state-handling.md`
- `.codex/wiki/reference/upstream-synchronization.md`
- `.codex/wiki/debugging/docker-web-assets.md`
