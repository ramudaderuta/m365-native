---
title: Upstream synchronization and push safety
type: reference
status: current
scope: repository-operating-docs
related_scopes: []
related_files:
  - AGENTS.md
source_docs:
  - AGENTS.md
tags:
  - git
  - upstream
  - synchronization
last_checked: 2026-07-18
updated: 2026-07-18T08:45:00Z
---

# Upstream synchronization and push safety

## Scope

This repository normally synchronizes the current local branch with the explicitly selected branch on `origin`. Synchronization and pushes require user authorization; a configured remote or authenticated GitHub account does not grant that authority by itself.

## Before integrating remote updates

Check `git status --short --branch`, the selected remote URL, and the current branch. Run `git fetch --prune <remote>`, inspect divergence, and review the incoming commits and diff. Do not automatically stash, discard, reset, or rebase local work to make an update apply.

Integrate only the intended remote branch. Prefer a fast-forward update. When a local modification or conflict blocks the update, retain the worktree and stop for review instead of forcing an integration.

## Before pushing combined work

Run the validation required for the complete combined change, then inspect the resulting status and diff. Confirm the target remote and branch, active GitHub account, and write access just before using a normal push. Do not use force-push unless separately and explicitly authorized after verifying the target ref.

## Verification

For documentation-only changes, run `git diff --check`, `ok-skill run wiki-note rebuild`, `ok-skill run wiki-note lint`, `ok-skill run wiki-note doctor --stale-refs --json`, and `ok-skill run wiki-note surface-check --json`. Confirm the branch is synchronized after the push with `git status --short --branch`.
