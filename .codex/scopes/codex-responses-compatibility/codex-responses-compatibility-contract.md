---
description: Contract for Codex Responses API compatibility repair.
---

# Codex Responses Compatibility Contract

## Context

- Current worktree: `main` has local deployment and wiki changes outside this scope; preserve them.
- Relevant source paths: `internal/web/multimodal.go`, `internal/web/protocol_compat.go`, `internal/web/codex_catalog.go`, `internal/web/codex_responses.go`, `internal/web/codex_usage.go`, `internal/web/protocol_handlers.go`, `internal/web/server.go`, and focused `internal/web/*_test.go` files.
- Relevant durable knowledge: `.codex/wiki/concepts/gateway-architecture.md` records `internal/web/` as the public protocol boundary.

## Findings

- Codex v0.144.5 sends Responses input content as `{"type":"input_text","text":"hello"}`. `parseContent` recognizes only `text`, so the adapter creates an empty OpenAI chat prompt.
- `streamResponsesAdapter` emits `response.created` before the inner request completes, then silently returns when the inner SSE produces no text or tool calls. A client observes a closed stream without `response.completed`.
- `/v1/models` returns standard `object` and `data` fields only. Codex's model manager requires an additional top-level `models` list.
- Evidence: a local authenticated request with `input_text` produced only `response.created`; Codex reported missing `models` and retried after the stream closed.

## Outcome

- Done when: Codex v0.144.5 can refresh the m365-native model catalog and complete a streamed `gpt1` Responses request without reconnecting.
- User-visible/runtime state: normal OpenAI-compatible callers retain `data`; Codex receives its `models` compatibility field and a terminal SSE event for both success and translated failures.
- Durable knowledge to preserve: Responses content blocks and terminal SSE events are a compatibility boundary owned by `internal/web/`.

## Goals / Non-goals

Goals:

- Parse Responses `input_text` and compatible output-text blocks without changing attachment behavior.
- Return the catalog in both OpenAI `data` and Codex-required `models` forms.
- Include the Codex-required stable `slug` on every catalog model.
- Include the Codex-required `display_name` on every catalog model.
- Include the Codex-required `supported_reasoning_levels` on every catalog model.
- Serialize `reasoning_efforts` and `supported_reasoning_levels` as Codex `ReasoningEffortPreset` objects, not strings.
- Mirror the remaining Codex execution metadata that is local-client-only: shell type, visibility, API availability, priority, and tier lists.
- Return the Codex-required `base_instructions` and matching `model_messages` template with an explicitly maintained gateway-safe instruction string; do not copy a user-local bundled prompt into the repository.
- Add administrator-managed public-model mappings that drive both Codex catalog entries and the corresponding upstream ChatHub tone. Seed the current Codex GPT-5.6 `Sol`, `Terra`, and `Luna` aliases.
- Ensure an inner streaming failure terminates with `response.failed` instead of a silent close.
- Return explicit local-estimate `usage` values in Responses completion objects so Codex can display context progress when ChatHub has no provider token counts. Include visible message framing, tool schemas, tool choice, tool calls, and completion framing. Use the embedded-vocabulary `tiktoken-go/tokenizer` BPE estimator for GPT routes and a clearly identified heuristic only when no compatible tokenizer is known.
- Add focused regression tests, build the image, and validate `codex exec --profile gpt1` against the local container.

Non-goals:

- Change API-key storage, account OAuth state, model routing, Nginx, Docker networking, or Codex user configuration.
- Add unsupported Codex tools or alter ChatHub's upstream protocol.

## Target files / modules

- `internal/web/multimodal.go`
- `internal/web/codex_catalog.go`, `internal/web/codex_usage.go`, `internal/web/codex_responses.go`
- `internal/web/protocol_handlers.go`
- `internal/web/server.go`
- `internal/web/codex_catalog_test.go`, `internal/web/codex_responses_compat_test.go`
- `internal/web/protocol_response.go`
- `go.mod`, `go.sum`, and `Dockerfile`
- focused tests in `internal/web/`
- `.codex/wiki/` and this contract for durable evidence
- `internal/web/settings.go`, focused settings tests, and `web/index.html`

## Settings UI Extension

- Target surface: the existing authenticated `/` administration console, Settings page.
- Visual thesis: a compact operational mapping table integrated into the existing quiet settings form, with public model IDs on the left and the chosen ChatHub tone on the right.
- Content plan: show the configured mappings, permit adding a row, permit removing a row, and retain the existing single Save action for persistence.
- Interaction thesis: rows are immediately editable; add and remove preserve the table's stable layout; save reports the existing success/error state without a page reload.
- Design rules: use the existing form controls, spacing, and restrained admin styling; keep keyboard-operable labeled inputs and explicit remove controls; do not add decorative panels.
- Browser verification: authenticated desktop and narrow mobile screenshots; add/remove row interaction; settings save; no clipped table content or console errors.

## Constraints

- Preserve authenticated API behavior and avoid logging request bodies, keys, tokens, or account data.
- Preserve existing OpenAI `data` catalog compatibility.
- Run `gofmt`, `go test ./...`, `go vet ./...`, `go build ./...`, Docker rebuild, API smoke checks, and a Codex profile smoke test.

## Boundaries

Allowed changes:

- Compatibility parsing, catalog serialization, Responses SSE terminal-error handling, focused tests, scope/wiki evidence, and generated wiki indexes.

Forbidden changes:

- Secret files, OAuth/account caches, API-key values, host Nginx routing, Compose port publishing, and unrelated protocol refactors.

## Decision Summary

| Decision | Evidence Source | Evidence Strength | Conflict | Result | Confidence Reason |
| --- | --- | --- | --- | --- | --- |
| Parse `input_text` alongside `text` | Runtime reproduction and `parseContent` source | High | None | Accept | Codex's request becomes an empty prompt without it. |
| Add a `models` catalog alias | Codex v0.144.5 decode error and `/v1/models` response | High | None | Accept | Preserve `data` while supplying the required field. |
| Add `slug` equal to model `id` | Codex v0.144.5 model-manager decode error | High | None | Accept | Codex requires a stable per-model slug and the gateway has no separate public model identifier. |
| Add `display_name` equal to model `id` | Codex v0.144.5 model-manager decode error after `slug` was supplied | High | None | Accept | The gateway has no distinct model-display metadata, so the stable model identifier is the least-surprising display value. |
| Add `supported_reasoning_levels` matching `reasoning_efforts` | Codex v0.144.5 model-manager decode error after `display_name` was supplied | High | None | Accept | Codex requires a separate advertised list; the gateway already validates and routes these exact levels. |
| Serialize reasoning levels as preset objects | Codex v0.144.5 decode error: string `none` where `ReasoningEffortPreset` was required; `codex debug models --bundled` | High | Prior catalog used strings | Accept | The bundled catalog proves the required shape is `{effort, description}`. |
| Mirror Codex execution metadata | Codex v0.144.5 decode error for missing `shell_type`; `codex debug models --bundled` | High | Gateway has no upstream equivalent | Accept | These fields are consumed only by Codex's local model manager and do not cross the ChatHub boundary. |
| Return gateway-owned instruction metadata | Codex v0.144.5 runtime decode error for missing `base_instructions`; bundled record shape | High | `base_instructions` and `model_messages.instructions_template` are used by Codex to construct later requests | Accept | Supply the required matching fields with a concise maintained template, rather than copying the local bundled 16KB prompt. The catalog response itself is not forwarded by ChatHub; the client may include the template in later request instructions. |
| Emit `response.failed` after inner stream errors | Local SSE reproduction and adapter control flow | High | None | Accept | Clients receive a terminal protocol event instead of an ambiguous disconnect. |
| Estimate Responses usage locally | Codex session token-count events contain no provider usage and completion objects omitted `usage` | High | None | Accept | Enables progress UI without claiming the estimate is billing-grade provider usage. |
| Use embedded-vocabulary tiktoken-go for GPT usage estimates | ChatHub does not expose provider counts; dependency inspection showed the selected tokenizer embeds `o200k_base` rather than fetching it at runtime | High | Default tiktoken-go loader downloads and caches BPE data | Accept | Improves GPT text/code/JSON estimates without adding a runtime network dependency; unknown and Claude routes retain an explicit fallback. |
| Raise the Go build baseline to 1.23 | tiktoken-go/tokenizer v0.7.0 requires Go 1.23 and includes GPT-5 model support | High | Latest v0.8.1 requires Go 1.26 with no needed tokenizer capability | Accept | A one-minor build-image update keeps a current dependency without an unrelated major toolchain jump. |
| Include visible tools and protocol framing in usage | Codex exposes an automatic compaction token limit, while the prior estimate omitted request schema and wrapper cost | High | Exact hidden ChatHub prompt/tokenization remains unavailable | Accept | Reduces systematic undercounting for tool-heavy sessions and makes the residual boundary explicit in response metadata. |

## Verification Surface

- Focused `go test ./internal/web` cases for catalog aliases, Responses content conversion, successful completion, and translated failure completion.
- `gofmt -w` on modified Go files; `go test ./...`; `go vet ./...`; `go build ./...`.
- `docker compose config --quiet`; `docker compose up -d --build`; container/listener/root smoke checks.
- Authenticated `/v1/models` schema check and `codex exec "hello" --profile gpt1 --skip-git-repo-check` with no model-decoding or stream-completion errors.

## Escalation Triggers

- Stop for a Codex model schema whose required `models` shape cannot be resolved by local validation, or an upstream account failure that prevents all authenticated completion tests.

## Rollback

- Rebuild and recreate the prior image from the current checked-out source, then restore the prior container. Source edits are narrow and revertible with an ordinary Git revert.

## Open Questions

- None.

## Execution Log / Evidence Updates

- 2026-07-18: scope created after reproducing the Codex `input_text` empty-stream failure and model-catalog decode error.
- 2026-07-18: added typed text-block parsing, `models` catalog alias, and terminal `response.failed` translation. Focused tests, `go test ./...`, `go vet ./...`, `go build ./...`, Compose validation, image rebuild, container recreation, and root HTTP smoke test passed.
- 2026-07-18: upstream `https://github.com/uefi2333/m365-native` was merged before runtime validation. The pre-existing test API key now returns HTTP 401 after revocation, so authenticated catalog and Codex end-to-end checks require a new temporary key.
- 2026-07-18: a subsequent Codex catalog refresh reported `missing field slug`. Added `slug` equal to each catalog model `id`, with regression coverage for both `data` and `models`. Full Go checks, Compose validation, root smoke test, and unauthorized `/v1/models` and `/v1/responses` checks passed after rebuild. The agent process still inherits an older revoked key, while its Codex profile has a separate reused refresh token; authenticated end-to-end validation remains blocked by local Codex authentication state rather than a protocol failure.
- 2026-07-18: session evidence showed successful turns but `token_count.info` was null and Responses completion objects lacked `usage`, leaving Codex context progress at zero. Added clearly labeled local usage estimates to non-stream and streaming completion objects, with focused regression tests. Full Go checks and rebuilt-container root smoke test passed.
- 2026-07-18: a later Codex model refresh accepted `slug` and then reported missing `display_name`. Added `display_name` equal to model `id`, with coverage for both catalog arrays; validation pending rebuild and runtime retest.
- 2026-07-18: replaced the GPT Responses usage approximation with `github.com/tiktoken-go/tokenizer` v0.7.0 and its embedded `o200k_base` vocabulary. GPT route estimates now identify `tiktoken_o200k_base_estimate`; unknown/Claude routes identify `heuristic_character_estimate`. Raised the module and Docker build baseline to Go 1.23. Focused and full Go checks, image rebuild, container recreation, listener check, and root HTTP smoke test passed. Authenticated Codex end-to-end validation remains pending the user's active local profile.
- 2026-07-18: Codex accepted the prior catalog fields and then reported missing `supported_reasoning_levels`. Added it at both catalog locations with the already validated reasoning level set; validation pending rebuild and runtime retest.
- 2026-07-18: expanded Responses usage estimation to include visible message/reply framing, `ToolCallID`, stable JSON serialization for structured content and calls, tool schemas, tool choice, and completion framing. The `m365` metadata now describes this visible estimate scope. This reduces undercounting that could delay a client-side auto-compact threshold, but does not count hidden upstream prompts, images, or provider reasoning tokens. Validation pending rebuild and runtime retest.
- 2026-07-18: Codex then rejected `reasoning_efforts` because it expects `ReasoningEffortPreset` objects rather than strings. Local `codex debug models --bundled` established the exact `{effort, description}` shape. Updated both reasoning catalog fields and added default catalog metadata plus JSON-shape regression coverage. `gofmt`, focused and full Go tests, `go vet`, `go build`, Compose validation, Docker rebuild, listener check, root HTTP smoke test, unauthorized `/v1/models` check, wiki rebuild/lint/doctor, scope scans, and `git diff --check` passed. The user's authenticated Codex profile remains the required final runtime confirmation.
- 2026-07-18: authenticated Codex profile then advanced past reasoning-level decoding and reported missing `shell_type`. This proves the preset-object repair is accepted. Added the remaining bundled-catalog execution metadata (`shell_type`, visibility, API availability, priority, and tier arrays) with regression coverage; validation pending rebuild and runtime retest.
- 2026-07-18: compared all fields on the bundled model record. Added its safe static tool, context, truncation, and capability metadata. Deliberately did not mirror `base_instructions` or `model_messages`, because those are behavioral instruction payloads rather than endpoint capability metadata. Validation pending rebuild and runtime retest.
- 2026-07-18: refactored Codex-specific catalog, usage estimation, and Responses-result projection into `codex_catalog.go`, `codex_usage.go`, and `codex_responses.go` within package `web`. Generic protocol handlers and Anthropic serialization remain in their existing owners; no ChatHub interface changed. Focused and full Go tests, `go vet`, `go build`, Compose validation, Docker rebuild, listener and root HTTP smoke tests, unauthorized `/v1/models` check, wiki rebuild/lint/doctor, scope scans, and `git diff --check` passed.
- 2026-07-18: Codex then reported a required `base_instructions` field. `codex debug models --bundled` showed it is duplicated as `model_messages.instructions_template`, with three empty personality variables and null `approvals`/`auto_review`. The gateway now returns a concise maintained template in both required locations. The catalog document itself is local metadata; Codex may include this template in subsequent request instructions, which then travel through the normal gateway/upstream request flow. Focused and full Go tests, vet, build, Compose validation, Docker rebuild, listener/root smoke checks, and the unauthorized catalog check passed. The final profile check is blocked before catalog decoding because `gpt1` now has both an invalid local API key and a reused Codex refresh token (`401`); reauthenticate the profile before retrying.
- 2026-07-18: Codex subsequently selected `gpt-5.6-sol` and warned that gateway metadata was absent. The bundled catalog confirms Sol, Terra, and Luna are independently selectable Codex aliases. Added persisted administrator model mappings that drive both the catalog and `reasoningTone` route selection. The UI inserts the mapping table directly below the interface address and above tool scheduling, preloads Sol/Terra/Luna, suggests all current bundled Codex model IDs including GPT-5.5, accepts a custom public ID, and restricts target selection to known ChatHub tones. Existing public IDs are overridden rather than duplicated. Full Go checks, Compose validation, image rebuild, listener/root smoke checks, unauthenticated catalog rejection, JavaScript syntax check, and targeted settings/catalog tests passed. Authenticated screenshot and save interaction remain unverified because the local administrator password is no longer the documented default; no settings were changed during browser validation.
- 2026-07-18: merged `upstream/main` at `c8f40d4` on integration branch `integrate/upstream-main-20260718`. Resolved overlapping `internal/web/` changes by retaining Codex `output_text` normalization and model mappings while integrating upstream direct image-URL handling and `toolPlanningMode`. The Responses empty-result branch remains terminal and now uses the upstream `empty_upstream_response` identifier. `go test ./...`, `go vet ./...`, `go build ./...`, Compose config validation, image rebuild, container recreation, localhost listener/root smoke, and unauthenticated `/v1/models` 401 check passed. The merge is pending commit/push authorization.
