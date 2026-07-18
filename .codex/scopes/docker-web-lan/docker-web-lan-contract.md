---
description: Contract for docker-web-lan.
---

# docker-web-lan Contract

## Context

- Current repo/worktree: `<workspace-root>/work/m365-native` on `main`.
- Relevant source paths: `Dockerfile`, `docker-compose.yml`, `web/login.html`, and `web/index.html`.
- Relevant archived scope references: none.

## Findings

- The runtime image copied only `/app/m365-native`; `internal/web/security_http.go` opens `web/login.html` and `web/index.html` relative to `/app`.
- `GET /` returned HTTP 500 with `web interface unavailable`; the container was otherwise listening on port 4141.
- The Compose file was locally changed to bind `0.0.0.0:4141`, which would create an upstream-merge-sensitive deployment preference.

## Outcome

- Done when: the runtime image contains `/app/web`, Compose preserves its upstream localhost default, and a local override exposes port 4141 to the LAN.
- User-visible/runtime state: `GET /` returns the login page; port 4141 listens on all IPv4 interfaces when the local override is active.
- Durable knowledge to preserve: web files are runtime assets rather than embedded binary resources; LAN exposure is a local deployment decision.

## Goals / Non-goals

Goals:
- Copy `web/` into the final runtime image.
- Move LAN port publishing into an ignored `docker-compose.override.yml`.
- Rebuild and smoke-test the running container.

Non-goals:
- Do not alter application routes, authentication, secrets, or public API behavior.
- Do not commit or push changes without explicit user authorization.

## Target files / modules

- `Dockerfile`
- `docker-compose.yml`
- `.gitignore`
- `docker-compose.override.yml` (ignored local deployment file)

## Constraints

- Keep the administrator password file private and mounted read-only.
- Preserve the upstream default bind address in `docker-compose.yml`.
- Retain the existing local-only data and secret volumes.

## Boundaries

Allowed changes:
- Docker image asset copying and local Compose port configuration.

Forbidden changes:
- Credential content, firewall rules, application authorization behavior, repository history, commits, and pushes.

## Decision Summary

| Decision | Evidence Source | Evidence Strength | Conflict | Result | Confidence Reason |
| --- | --- | --- | --- | --- | --- |
| Include `web/` in the final stage | `Dockerfile`, `internal/web/security_http.go`, HTTP 500 observation | High | None | Copy `/src/web` to `/app/web` | Runtime path requirement is direct and reproducible |
| Keep LAN publishing local | User requested LAN access; `docker-compose.yml` is upstream-managed | High | None | Ignored Compose override | Avoids an upstream sync conflict for a host-specific setting |

## Verification surface

- `docker compose config --quiet`
- `docker compose up -d --build`
- `docker compose ps`
- `ss -ltn 'sport = :4141'`
- `curl --fail --max-time 5 http://127.0.0.1:4141/`
- `docker exec` assertion that `/app/web/login.html` and `/app/web/index.html` exist

## Escalation triggers

- Escalate only when code/runtime evidence, authoritative wiki, and scope docs materially conflict and the conflict cannot be resolved from local evidence.
- Escalate for data deletion, permission semantics, production access model, or public API compatibility decisions outside the stated boundaries.
- Escalate when user-specified boundaries cannot be satisfied together.

## Rollback

- Remove the `COPY --from=build /src/web /app/web` line and recreate the container to return to the prior image behavior.
- Remove `docker-compose.override.yml` to restore localhost-only publishing without changing tracked Compose configuration.

## Open questions

- None.

## Execution log / evidence updates

- 2026-07-18: Scope created after reproducing the missing-web-assets failure and confirming the current Compose port binding.
- 2026-07-18: Added `COPY --from=build /src/web /app/web` to `Dockerfile`.
- 2026-07-18: Restored the tracked Compose file to localhost publishing and added ignored `docker-compose.override.yml` with a `!override` port list for `0.0.0.0:4141`.
- 2026-07-18: `docker compose config` confirmed one published port (`0.0.0.0:4141`). `docker compose up -d --build` passed. The container contains both required HTML files, `GET /` returned HTTP 200, and `ss` confirmed the LAN listener.
- 2026-07-18: User authorized committing and pushing the validated tracked changes to the `ramudaderuta/m365-native` fork.
