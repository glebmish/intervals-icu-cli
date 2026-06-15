# CLAUDE.md

## Project

Go CLI (`intervals`) for the intervals.icu API. 148 operations across 13 resource groups.

## Build & Test

```bash
go build -o intervals .
go test ./...
go vet ./...
```

**Always install after building:** `go install . && cp ~/go/bin/intervals-icu-cli ~/go/bin/intervals`

This must be done every time you make changes so the `intervals` binary on PATH stays current.

## Architecture

- `main.go` — entry point, calls `cmd.Execute()`
- `internal/cmd/` — Cobra commands. `root.go` wires config/client. `activities.go` defines shared helpers (`doGet`, `doMutate`, `doDelete`, `doDownload`). Each resource group is one file.
- `internal/config/` — YAML config loading, env/flag override, validation
- `internal/api/` — HTTP client with Basic Auth, path substitution, error handling
- `internal/validate/` — Input sanitization (path params, dates, JSON bodies)
- `internal/format/` — Output formatting (JSON, NDJSON, field filtering)
- `internal/cmd/openapi-spec.json` — Embedded OpenAPI spec for `schema` command

## Adding a New API Operation

1. Identify the resource group file in `internal/cmd/`
2. Add a `new<Resource><Action>Cmd()` function following the existing pattern
3. Register it in the `init()` function's `AddCommand` call
4. Update the command name mapping in `schema.go` if needed
5. `go build` to verify

## Conventions

- All commands use shared helpers: `doGet`, `doMutate`, `doDelete`, `doDownload`
- Required IDs validated with `validate.PathParam`, dates with `validate.DateParam`
- JSON bodies validated with `validate.JSONBody`
- API client retrieved from context: `api.FromContext(cmd.Context())`
- Path template variables: `{id}` and `{athleteId}` auto-substituted with athlete ID
- Worktree directory: `.worktrees/`
