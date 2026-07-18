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

For LAN access, keep the application published only at `127.0.0.1:4141` and place Nginx in front of it only on the host LAN address at `192.0.2.10:4141`. Binding Nginx specifically to that LAN address does not conflict with Docker's loopback-only `127.0.0.1:4141`; both use the same port number on distinct local addresses. Standard HTTP redirects to standard HTTPS, and the default 443 virtual host serves Nginx's `/var/www/html` page rather than m365-native. Other named Nginx virtual hosts can share ports 80 and 443. The m365-native HTTPS proxy forwards to the loopback application port. Do not publish the Docker port directly: administrator sessions use secure cookies and management data must not traverse LAN HTTP.

The current local Nginx deployment uses `deploy/nginx/nginx.conf` as `/etc/nginx/nginx.conf` and loads the virtual host from `/etc/nginx/conf.d/server.conf`, sourced from `deploy/nginx/server.conf`. The main configuration centralizes the service certificate, TLS 1.2/1.3 policy, session policy, common security headers, HTTP-to-HTTPS redirect, default HTTPS rejection, request buffers, and the 25 MB request cap; the site enables HTTP/2 and defines only routing and proxy behavior. Proxy request/response buffering is globally disabled for streamed API responses. Gzip is configured for static response types but explicitly disabled for proxied authenticated responses. The `events` block allows 1024 connections per worker and enables batch acceptance, while leaving Linux to select its appropriate event mechanism rather than forcing `select`. TCP no-delay, MIME hash buckets, server-name hash buckets, and hidden Nginx version tokens are set in the main configuration. It does not apply unrelated source-address/User-Agent spoofing, broad CORS, or WebSocket rules. The deployment uses an internal CA and a certificate with the host's LAN IP and hostname as subject alternative names. The CA private key, service private key, and their certificates are stored in `/etc/nginx/conf.d/` with non-`.conf` extensions; both private keys remain root-only. Distribute only the public CA certificate to authorized client trust stores. Clients must trust that CA before accessing the HTTPS address. Validate with `nginx -t`, `systemctl reload nginx`, an HTTP-to-HTTPS redirect check, and an HTTPS request using the CA certificate.

## Verification

After a container change, run `docker compose up -d --build`, `docker compose ps`, check that `/app/web/login.html` and `/app/web/index.html` exist, and request `http://127.0.0.1:4141/` expecting HTTP 200. For password-file changes, also verify the `m365` user can read `/run/secrets/m365_admin_password` without printing it.
