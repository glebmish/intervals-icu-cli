# intervals CLI design

This document describes the design of the `intervals` command-line interface for
the intervals.icu API. The CLI's primary user is an AI agent, so the design
optimizes for predictability and defense-in-depth: a stable command surface, raw
JSON payloads, machine-readable errors, and self-describing schema discovery. The
CLI itself — via `intervals schema` — is the canonical source of truth about the
API surface, not external docs.

## Command surface
- Shape: `intervals <resource> <action>` (single API surface).
- Action verbs map 1:1 to the OpenAPI `operationId`s via `operationIDToCommand`
  in `internal/cmd/schema.go`. The bijection is enforced at build time by
  `TestAllOperationIdsMapped` and `TestEveryMappedCLINameExists` in
  `schema_test.go`, so the CLI cannot drift from the spec.
- One file per resource group in `internal/cmd/`.

## Build identity
- CLI short name: `intervals`
- Module path: `github.com/glebmish/intervals-icu-cli`

## Authentication
- Scheme: HTTP Basic auth — username `API_KEY`, password = the user's API key.
- Auth flag: `--api-key`
- Env var: `INTERVALS_API_KEY`
- Config field: `api_key`
- Resolution precedence: flag → env var → config file.
- Config file: `~/.config/intervals/config.yaml`, mode `0600`, overridable via
  `INTERVALS_CONFIG`.
- Token acquisition: https://intervals.icu/settings
- Headless bootstrap: set `INTERVALS_API_KEY`, or use
  `intervals config show --unmasked` to print credentials for piping into a
  config file, dotenv, or secret manager.

## Tenant scope
Most paths are scoped to an athlete (`/api/v1/athlete/{id}/...`), so the athlete
ID is a top-level concern stored in config.
- Flag: `--athlete-id`
- Env var: `INTERVALS_ATHLETE_ID`
- Config field: `athlete_id`
- The client substitutes both `{id}` and `{athleteId}` placeholders with the
  configured athlete ID. Any other `{key}` placeholder is filled from the params
  map.

## Source
- HTTP API with an OpenAPI spec.
- The canonical spec is checked in at the repo root as `openapi-spec.json`; a
  byte-identical mirror at `internal/cmd/openapi-spec.json` is what `//go:embed`
  in `schema.go` consumes (keep both in sync). 148 operations across 13 resource
  groups. The `schema` command is derived from it and runs offline, without
  credentials.

## Resource groups
activities, activity, athlete, chats, custom-items, events, folders, gear, misc,
shared-events, sport-settings, wellness, workouts.

## Global flags
Inherited by every command from the root:
- `--format json|ndjson|text` (default `json`)
- `--fields a.b,c.d` — dotted-path field mask, comma-separated
- `--dry-run` — print the request that would be sent and exit without executing
- `--yes` — skip the interactive confirmation on destructive ops
- `--json '<body>'` — raw request body for write ops (canonical, not a fallback)
- `--params '<obj>'` — raw query/path params as a JSON object
- `--api-key`, `--athlete-id`, `--base-url`

Per-command convenience flags (`--oldest`, `--query`, `--limit`, …) are sugar that
build the params map; `--json` and `--params` remain the canonical raw path.

## Schema discovery
- `intervals schema --list` — every operation in `<resource>.<action>` form.
- `intervals schema <op>` — method, path, parameters, request and response shapes.
- `intervals schema <TypeName>` — type definitions; `--resolve-refs` inlines `$ref`s.
- Runs offline and without credentials.

## Output contract
- Success output is JSON by default on stdout. `--format ndjson` streams one object
  per line; `--format text` passes non-JSON bodies through unchanged.
- `--fields` filters by dotted path, descending into arrays implicitly.
- Errors are emitted on stderr (via cobra); success and error output never mix.
- Envelope: none. intervals.icu returns raw arrays/objects with no constant
  wrapper key, so `--fields` paths are unprefixed (e.g. `--fields id,name`).

## Pagination
The API does not paginate: the spec exposes no page/offset/cursor parameters.
Listing endpoints window by date range (`oldest`/`newest`) plus an optional
`limit`. No page-walking flags are provided because there is nothing to walk.

## Exit codes
- 0 success
- 1 API error (4xx/5xx response)
- 2 auth error (missing/invalid credentials)
- 3 validation error (input rejected before the request)
- 4 discovery/schema error
- 5 internal error (panic, recovered in `main.go`)

Codes 1–4 are mapped in `main.go::exitCode`; code 5 is emitted by the
panic-recover block in `main.go::main`.

## Error messages
The API client's error formatter (`internal/api/client.go::Error()`) always
includes the server's response body inline, followed by an actionable hint keyed
by status code: 401, 403, 404, 409, 422, 429, 500, 503.

## Input hardening
All flag values are treated as adversarial and validated before any request:
- Path params: reject control chars, path traversal (`..`), query-injection
  characters (`?`, `#`, `&`), and pre-encoded `%`.
- Dates: strict `YYYY-MM-DD`.
- JSON bodies (`--json`, `--params`): valid JSON, no rogue control chars.
- Upload paths: reject control chars and NUL, clean the path.

Helpers live in `internal/validate/`; each returns a validation-typed error that
maps to exit code 3.

## File uploads
- Upload commands: `activities upload`, `custom-items upload-image`,
  `shared-events upload-image`.
- Canonical flag: `--upload <path>` (a peer of `--json`/`--params`). `--file` is
  retained as a deprecated back-compat alias.
- Hardening: paths are validated by `validate.FilePath` via the shared
  `requireUploadFile` helper before reading. Policy is pragmatic — absolute and
  `../` paths are allowed (no CWD sandbox) so uploads from e.g. `~/Downloads`
  work; only injection shapes are blocked.

## Response sanitization
String values in API responses are sanitized before stdout
(`internal/format/output.go::sanitize`): control characters are stripped and
`<system>`/`<assistant>`/`<tool_use>`/`<tool_result>` tag wrappers are removed, to
defend against prompt injection embedded in user-controlled fields (activity and
event names/descriptions, chat messages, gear and folder names).

## Bundled skills
Agent-facing knowledge ships inside the binary under `internal/cmd/skills/`:
- `intervals-shared` — the front door: global flags, exit codes, schema discovery.
- `intervals-training` — activities, activity analysis, events, workouts, wellness.
- `intervals-settings` — athlete profile, sport settings, gear, folders, chats.
- `recipe-plan-training`, `recipe-training-review` — multi-step workflows.

Runtime access without installing: `intervals skills list` and
`intervals skills get <name>` (both offline, no token). `intervals skills install`
writes them to disk.

## Open questions / deferred
- Resource skills are coalesced into two (`intervals-training`,
  `intervals-settings`) rather than one per resource group.
- No MCP surface (`intervals mcp`) yet.
