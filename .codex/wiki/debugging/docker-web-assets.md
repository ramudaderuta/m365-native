---
title: Docker image web asset failure
type: debugging
status: current
scope: docker-web-lan
related_scopes:
  - repository-operating-docs
related_files:
  - Dockerfile
  - internal/web/security_http.go
  - web
source_docs: []
tags:
  - docker
  - debugging
  - web
last_checked: 2026-07-18
updated: 2026-07-18T08:34:37Z
---

# Docker image web asset failure

## Symptom

The container is running and port 4141 is listening, but `GET /` returns HTTP 500 with `web interface unavailable`.

## Cause

The root handler opens `web/login.html` or `web/index.html` relative to the runtime working directory. The original final Docker stage copied only the compiled binary and omitted the `web/` directory.

## Fix

Copy `/src/web` from the build stage to `/app/web` in the final Docker stage, then rebuild and recreate the service with `docker compose up -d --build`.

## Probe

Confirm both `/app/web/login.html` and `/app/web/index.html` exist in the container, then request the root URL and expect HTTP 200. If the files are present but the response still fails, inspect the container working directory and the root handler path logic.
