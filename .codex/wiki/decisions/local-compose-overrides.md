---
title: Local Compose overrides own LAN exposure
type: decision
status: accepted
scope: docker-web-lan
related_scopes:
  - repository-operating-docs
related_files:
  - docker-compose.yml
  - Dockerfile
  - .gitignore
source_docs:
  - .codex/scopes/docker-web-lan/docker-web-lan-contract.md
tags:
  - docker
  - lan
  - deployment
last_checked: 2026-07-18
updated: 2026-07-18T08:34:37Z
decision_date: 2026-07-18
---

# Local Compose overrides own LAN exposure

## Context

The tracked Compose file uses a localhost port binding as the safe upstream default. This workstation needs optional LAN access without turning that host-specific choice into a merge conflict during upstream synchronization.

## Decision

Keep `docker-compose.yml` loopback-only. Add `docker-compose.override.yml` to `.gitignore`; when LAN access is required, an untracked local override replaces the port list using the Compose `!override` tag.

## Rationale

A normal Compose override merges port entries, producing duplicate host bindings and a startup failure. Replacing the list gives one effective published port while isolating the deployment choice from the tracked upstream configuration.

## Consequences

The override is intentionally local and is not committed. Operators must verify the effective configuration and listener after changing it. LAN exposure remains protected only by the application authentication and host network controls, so TLS and edge access control are required before any broader exposure.
