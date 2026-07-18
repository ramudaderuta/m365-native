---
title: Credential and persistent-state handling
type: reference
status: current
scope: repository-operating-docs
related_scopes: []
related_files:
  - internal/auth/cache.go
  - internal/web/keys.go
  - internal/web/sessions.go
  - internal/web/admin_security.go
  - docker-compose.yml
source_docs:
  - README.md
  - SECURITY.md
tags:
  - security
  - credentials
  - persistence
last_checked: 2026-07-18
updated: 2026-07-18T08:34:13Z
---

# Credential and persistent-state handling

## Scope

This gateway persists administrator credentials, OAuth account tokens, API-key hashes, and conversation mappings. These are sensitive operational data even when some records are hashed.

## Persistent paths

The container maps `./data` to `/data`. The Compose defaults place account metadata, token cache, session mapping, and API-key hashes there. The administrator password is read from the separate, read-only file mounted at `/run/secrets/m365_admin_password`. Local `data/` and `secrets/` are ignored by Git.

## Hard rules

Never commit, print, copy into documentation, or send OAuth tokens, refresh tokens, cookies, account caches, API keys, raw password files, or debug logs. Preserve file modes that restrict secret-bearing data to the service owner. Do not expose the service outside trusted networks without TLS and an additional access-control layer.

## Authentication behavior

Administrative management uses a session cookie after login. `/v1/` endpoints require a valid API key. The administrator password must be changed from the default bootstrap value before management operations proceed. Failed administrator logins are rate limited.

## Verification

Review changes for secret-bearing paths and run `git status --short` before committing. Validate only the presence, permissions, and mounting of secret files; do not print their contents.
