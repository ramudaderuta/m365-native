---
title: Maintenance Log
type: maintenance-log
status: current
updated: 2026-07-18T08:33:14Z
---

# Maintenance Log

Append-only history for wiki updates caused by scope work, implementation closeout, or knowledge refresh.

## 2026-07-18T08:35:31Z [repository-operating-docs]

- Summary: Initialized current-code-backed repository knowledge and aligned README and AGENTS.md with the wiki.
- Pages: concepts/gateway-architecture.md, how-to/container-operation.md, reference/credential-and-state-handling.md, decisions/local-compose-overrides.md, debugging/docker-web-assets.md
- Verification: ok-skill run wiki-note rebuild; ok-skill run wiki-note lint; ok-skill run wiki-note doctor --json; ok-skill run wiki-note surface-check --json
- Residual risk: Wiki records local deployment assumptions; recheck them after upstream Docker or Compose changes.

## 2026-07-18T10:29:08Z [state-persistence-write-errors]

- Summary: Recorded the bind-mounted state-directory permission failure and non-destructive probe.
- Pages: debugging/docker-state-directory-permissions.md
- Verification: stat ownership and container write probe
- Residual risk: Previously missing state files require operator re-creation or restoration from an external backup.

## 2026-07-18T12:25:51Z [codex-responses-compatibility]

- Summary: Merged upstream/main c8f40d4 on an integration branch; resolved Codex compatibility overlaps and validated the rebuilt local container.
- Pages: none
- Verification: go test ./...; go vet ./...; go build ./...; docker compose up -d --build
- Residual risk: Authenticated Codex and Settings UI checks still require user-authorized credentials.
