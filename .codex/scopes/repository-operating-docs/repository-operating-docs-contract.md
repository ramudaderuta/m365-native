---
description: Contract for repository-operating-docs.
---

# repository-operating-docs Contract

## Context

- Current repo/worktree: `<workspace-root>/work/m365-native` on the `main` branch.
- Relevant source paths: `README.md`, `AGENTS.md`, `cmd/server/main.go`, `internal/auth/`, `internal/chathub/`, `internal/web/`, `Dockerfile`, and `docker-compose.yml`.
- Relevant archived scope references: none.

## Findings

- The repository had no agent contract or structured project wiki. Durable architecture, container-operation, security, and Docker debugging knowledge was otherwise split between source, README, and a previous scope record.
- Current behavior was verified by inspecting the startup path, route middleware, token cache, state stores, Dockerfile, Compose configuration, tests, and existing README/security guidance.

## Outcome

- Done when: `AGENTS.md` gives accurate source boundaries, security rules, and validation defaults, and `.codex/wiki/` provides typed, current pages with generated indexes.
- User-visible/runtime state: no runtime behavior changes.
- Durable knowledge to preserve: package ownership, Docker web-asset requirements, local LAN override policy, and secret-handling boundaries.

## Goals / Non-goals

Goals:
- Add a concise repository agent contract.
- Initialize a safe, structured local wiki backed by current code and docs.
- Link README, agent guidance, and wiki without duplicating long manuals.

Non-goals:
- Do not change application behavior, credentials, network/firewall policy, or upstream APIs.
- Do not commit or push without separate authorization.

## Target files / modules

- `AGENTS.md`
- `README.md`
- `.codex/wiki/`
- `.codex/scopes/repository-operating-docs/`

## Constraints

- Keep credentials and personal account data out of all new documents.
- Treat source and tests as authoritative over derived documentation.
- Use wiki tooling for wiki pages and generated indexes.

## Boundaries

Allowed changes:
- Documentation, repository agent instructions, wiki structure, and scope evidence.

Forbidden changes:
- Go behavior, Docker runtime behavior, local secrets, Compose port policy, Git history, commits, and pushes.

## Decision Summary

| Decision | Evidence Source | Evidence Strength | Conflict | Result | Confidence Reason |
| --- | --- | --- | --- | --- | --- |
| Agent contract scope | Current repository structure, README, CONTRIBUTING, and SECURITY | High | None | Concise boundaries and validation matrix in `AGENTS.md` | Future changes need source-aware constraints rather than duplicated manuals |
| Wiki contents | Current source and runtime observations | High | None | Four typed pages plus indexes | Covers architecture, operations, security, a deployment decision, and a known Docker failure |

## Verification surface

- `ok-skill run wiki-note rebuild`
- `ok-skill run wiki-note lint`
- `ok-skill run wiki-note doctor --json`
- `ok-skill run wiki-note surface-check --json`
- `ok-skill run repo-task-driven placeholder-scan .codex/scopes/repository-operating-docs`
- `ok-skill run repo-task-driven text-scan .codex/scopes/repository-operating-docs README.md AGENTS.md`
- `git diff --check`

## Escalation triggers

- Escalate only when code/runtime evidence, authoritative wiki, and scope docs materially conflict and the conflict cannot be resolved from local evidence.
- Escalate for data deletion, permission semantics, production access model, or public API compatibility decisions outside the stated boundaries.
- Escalate when user-specified boundaries cannot be satisfied together.

## Rollback

- Remove `AGENTS.md`, the README wiki link, and the `.codex/wiki/` tree together if the repository later adopts another documented knowledge system.

## Open questions

- None.

## Execution log / evidence updates

- 2026-07-18: Scope created after current-code and documentation discovery.
- 2026-07-18: Initialized the wiki and added pages for architecture, container operation, credential/state handling, local Compose policy, and Docker web-asset troubleshooting.
- 2026-07-18: Added `AGENTS.md` and a README link to keep the documentation surfaces aligned.
- 2026-07-18: Wiki rebuild, lint, doctor, navigation build/search, impact, and surface checks passed. The only surface-check note is informational: this repository has no separate `docs/` tree to archive or map.
- 2026-07-18: Scope placeholder and residual-text scans and `git diff --check` passed.
