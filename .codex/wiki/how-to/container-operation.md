---
title: Container operation and local deployment
type: how-to
status: current
scope: repository-operating-docs
related_scopes: []
related_files:
  - Dockerfile
  - docker-compose.yml
  - web
  - README.md
source_docs:
  - README.md
tags:
  - docker
  - deployment
  - operations
last_checked: 2026-07-18
updated: 2026-07-18T08:34:13Z
---

# Container operation and local deployment

## Scope

The Compose service builds `m365-native:latest`, runs as the non-root `m365` user, and stores mutable state under `/data`. The tracked Compose file publishes `127.0.0.1:4141` by default.

## First-run procedure

Create `data/` and `secrets/`, write a long administrator password to `secrets/m365_admin_password`, restrict that file to mode 0600, then run `docker compose up -d --build`. Do not place real values in `.env` or tracked configuration.

## Runtime assets

The final image must contain both the binary at `/app/m365-native` and the runtime HTML assets at `/app/web`. `internal/web/security_http.go` opens `web/login.html` and `web/index.html` relative to the `/app` working directory. A container that omits `web/` starts but returns `web interface unavailable` from the root page.

## LAN deployment

LAN publishing is a host-specific choice. Keep the tracked default unchanged and use an ignored `docker-compose.override.yml` with `ports: !override` and `0.0.0.0:4141:4141` when LAN access is intended. Confirm the effective configuration with `docker compose config` and verify the listener using `ss -ltn sport = :4141`. Use TLS and an access-control layer before broader exposure.

## Verification

After a container change, run `docker compose up -d --build`, `docker compose ps`, check that `/app/web/login.html` and `/app/web/index.html` exist, and request `http://127.0.0.1:4141/` expecting HTTP 200.
