---
description: Return API key persistence failures to the administrator instead of silently retaining in-memory-only keys.
---

# state-persistence-write-errors Contract

## Context

- Current repo/worktree: `main`, with Compose state mounted at `/data`.
- Relevant source paths: `internal/web/keys.go`, `internal/web/server.go`, and `internal/web/keys_test.go`.
- Relevant durable knowledge: `.codex/wiki/debugging/docker-state-directory-permissions.md`.

## Findings

- `apiKeyStore.save` discarded directory-creation and file-write errors.
- `create` appended a record before saving, and `revoke` changed its state before saving, so an unwritable bind mount could produce session-only state.
- The Docker service runs as UID/GID `100:101`; the state bind mount must be writable by that identity.

## Outcome

- Done when: create and revoke requests fail with HTTP 500 when their state cannot be persisted, and their in-memory mutations are rolled back.
- User-visible/runtime state: API-key values and JSON format are unchanged; administrators receive a failure instead of a misleading success response.
- Durable knowledge to preserve: `/data` is the persistent state boundary and must be writable by the non-root service user.

## Goals / Non-goals

Goals:
- Propagate persistence errors from API-key create and revoke operations.
- Roll back in-memory mutations when those writes fail.
- Add regression coverage for both failure paths.

Non-goals:
- Do not change API-key generation, hashing, file schema, authentication semantics, or account-token storage.
- Do not change the Compose mount topology or run the application as root.

## Target files / modules

- `internal/web/keys.go`
- `internal/web/server.go`
- `internal/web/keys_test.go`

## Constraints

- Do not read, print, create, or modify real account, API-key, token, or password data during validation.
- Keep changes within the `internal/web` ownership boundary.

## Boundaries

Allowed changes:
- API-key persistence error propagation, request error handling, regression tests, and local scope/wiki evidence.

Forbidden changes:
- Data migration, deletion, credential rotation, Docker privilege escalation, or public API schema changes.

## Decision Summary

| Decision | Evidence Source | Evidence Strength | Conflict | Result | Confidence Reason |
| --- | --- | --- | --- | --- | --- |
| Return write errors and roll back mutations | `internal/web/keys.go`, runtime permission probe, user symptom | High | resolved | Selected | Preserves API semantics and prevents false success without weakening Docker isolation. |

## Verification surface

- `gofmt -w internal/web/keys.go internal/web/keys_test.go internal/web/server.go`
- `go test ./internal/web`
- `go test ./...`
- `go vet ./...`
- `go build ./...`
- `docker compose config --quiet`
- `docker compose up -d --build`
- Verify the service listener and `http://127.0.0.1:4141/` response.

## Escalation triggers

- Escalate only when code/runtime evidence, authoritative wiki, and scope docs materially conflict and the conflict cannot be resolved from local evidence.
- Escalate for data deletion, permission semantics, production access model, or public API compatibility decisions outside the stated boundaries.
- Escalate when user-specified boundaries cannot be satisfied together.

## Rollback

- Revert the source commit. No persistent data migration is involved, and the pre-existing state-file format remains unchanged.

## Open questions

- None.

## Execution log / evidence updates

- 2026-07-18: Confirmed the state bind mount was empty and unwritable by UID/GID `100:101`; corrected its local ownership and verified a write/remove probe without inspecting state contents.
- 2026-07-18: Implemented create/revoke error propagation and in-memory rollback. Targeted package tests passed.
- 2026-07-18: `go test ./...`, `go vet ./...`, `go build ./...`, Compose validation, rebuilt-container HTTP smoke test, and the service-user write probe all passed.
