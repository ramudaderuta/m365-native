---
title: Docker bind-mounted state write failure
type: debugging
status: current
scope: repository-operating-docs
related_scopes: []
related_files:
  - docker-compose.yml
  - internal/auth/cache.go
  - internal/web/keys.go
source_docs: []
tags:
  - docker
  - persistence
  - permissions
  - debugging
last_checked: 2026-07-18
updated: 2026-07-18T10:26:32Z
---

# Docker bind-mounted state write failure

# Docker bind-mounted state write failure

## Symptom

Accounts or API keys appear during an active session but are absent after a container restart or recreation.

## Cause

The final image runs as the non-root `m365` user (UID/GID `100:101`). Compose bind mounts preserve host ownership and mode. If the directory mounted at `/data` is not writable by that identity, the token and account stores cannot persist. The API-key store currently ignores write errors, so an in-memory key can appear valid until the process restarts.

## Remediation

Before first use, assign the state directory to UID/GID `100:101` while retaining restrictive permissions. Do not replace the bind mount with container-local storage; rebuilding or recreating a container discards such storage. Existing state files must be restored from a backup before recreating keys or accounts.

## Probe

Inspect only numeric ownership and mode with `stat`. From the running container, test that the `m365` user can create and remove a temporary probe file in `/data`; do not read or print state-file contents. Confirm the host directory remains mounted at `/data` after a rebuild.
