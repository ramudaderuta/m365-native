---
title: Container operation and local deployment
type: how-to
status: current
scope: repository-operating-docs
related_scopes: []
related_files:
  - Dockerfile
  - docker-compose.yml
  - deploy/nginx/nginx.conf
  - deploy/nginx/server.conf
  - web
  - README.md
source_docs:
  - README.md
tags:
  - docker
  - deployment
  - operations
last_checked: 2026-07-18
updated: 2026-07-18T09:15:52Z
---

# Container operation and local deployment

## Scope

The Compose service builds `m365-native:latest`, runs as the non-root `m365` user, and stores mutable state under `/data`. Its stable container name is `m365-native`, avoiding Compose's default repeated `<project>-<service>-<index>` name. The tracked Compose file publishes `127.0.0.1:4141` by default. The fixed container name and port mapping define this as a single-instance local deployment.

## First-run procedure

Create `data/` and `secrets/`, write a long administrator password to `secrets/m365_admin_password`, then run `docker compose up -d --build`. The container runs as the `m365` user (UID/GID `100:101`), while bind mounts preserve host ownership. Set the password file to owner `100:101` and mode 0600 so the service can read it without making it world-readable. Do not place real values in `.env` or tracked configuration.

## Runtime assets

The final image must contain both the binary at `/app/m365-native` and the runtime HTML assets at `/app/web`. `internal/web/security_http.go` opens `web/login.html` and `web/index.html` relative to the `/app` working directory. A container that omits `web/` starts but returns `web interface unavailable` from the root page.

## LAN deployment

For LAN access, keep the application published only at `localhost` and place Nginx in front of it on a designated LAN HTTPS endpoint such as `https://<lan-host>:4141`. Binding Nginx to a specific LAN address does not conflict with Docker's loopback-only listener because both use the same port number on distinct local addresses. Standard HTTP redirects to standard HTTPS, and the default 443 virtual host serves Nginx's `/var/www/html` page rather than m365-native. Other named Nginx virtual hosts can share ports 80 and 443. The m365-native HTTPS proxy forwards to the loopback application port. Do not publish the Docker port directly: administrator sessions use secure cookies and management data must not traverse LAN HTTP.

The tracked Nginx files are sanitized templates. Keep the effective Nginx configuration, certificate locations, internal CA material, and real server names outside the repository. The template centralizes TLS policy, common security headers, request buffers, proxy streaming behavior, and gzip defaults; the site template defines only routing and proxy behavior. Replace its documentation address and server name only in an untracked local deployment copy. Do not add source-address/User-Agent spoofing, broad CORS, or WebSocket rules unless the served application requires them. Distribute only a public CA certificate to authorized client trust stores. Validate a local deployment with `nginx -t`, a service reload, an HTTP-to-HTTPS redirect check, and an HTTPS request using the local CA certificate.

## Verification

After a container change, run `docker compose up -d --build`, `docker compose ps`, check that `/app/web/login.html` and `/app/web/index.html` exist, and request `http://127.0.0.1:4141/` expecting HTTP 200. For password-file changes, also verify the `m365` user can read `/run/secrets/m365_admin_password` without printing it.
