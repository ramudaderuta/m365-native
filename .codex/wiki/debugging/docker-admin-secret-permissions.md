---
title: Docker administrator secret permission failure
type: debugging
status: current
scope: repository-operating-docs
related_scopes: []
related_files:
  - Dockerfile
  - docker-compose.yml
  - internal/web/admin_security.go
source_docs: []
tags:
  - docker
  - credentials
  - permissions
  - debugging
last_checked: 2026-07-18
updated: 2026-07-18T08:55:31Z
---

# Docker administrator secret permission failure

## Symptom

The administrator password stored in `secrets/m365_admin_password` cannot log in, while the initial password still works and the login page requests an immediate password change.

## Cause

The final image runs as the non-root `m365` user (UID/GID `100:101`), and Compose bind-mounts the host password file read-only at `/run/secrets/m365_admin_password`. Bind mounts preserve host ownership and mode. A host file owned by another user with mode 0600 is unreadable to `m365`; the password loader then falls back to its bootstrap password.

## Remediation

Keep the password file non-empty and secret. Set its owner to UID/GID `100:101` and mode 0600, then recreate the service so it reloads the file. This owner change means later manual rotation needs an authorized privileged edit or a deliberate ownership handoff; never loosen the file to world-readable permissions.

## Probe

Check the host file's numeric owner and mode with `stat` without reading its contents. In the container, test read access as `m365` using `test -r /run/secrets/m365_admin_password`; do not print, hash, copy, or log the password. Confirm the application is running and the root page returns HTTP 200 after recreation.
