---
name: intervals-shared
description: "intervals CLI: Authentication, global flags, security rules, and schema discovery."
metadata:
  version: 0.1.0
  openclaw:
    category: "productivity"
    requires:
      bins:
        - intervals
---

# intervals — Shared Reference

## Installation

```bash
go install github.com/glebmish/intervals-icu-cli@latest
```

The `intervals` binary must be on `$PATH`.

## Authentication

Run the interactive setup:

```bash
intervals config init
```

Or create `~/.config/intervals/config.yaml` manually:

```yaml
api_key: your_api_key_here
athlete_id: your_athlete_id_here
```

### Environment Variables

| Variable | Description |
|---|---|
| `INTERVALS_API_KEY` | API key (overrides config file) |
| `INTERVALS_ATHLETE_ID` | Athlete ID (overrides config file) |
| `INTERVALS_BASE_URL` | API base URL (default: `https://intervals.icu/api/v1`) |
| `INTERVALS_CONFIG` | Path to config file (default: `~/.config/intervals/config.yaml`) |

## CLI Syntax

```bash
intervals <resource> <action> [flags]
```

## Global Flags

| Flag | Description |
|------|-------------|
| `--format <fmt>` | Output format: `json` (default), `ndjson`, `text` |
| `--fields '<mask>'` | Comma-separated fields to include in output |
| `--dry-run` | Validate request without executing |
| `--yes` | Skip confirmation prompts on deletes |
| `--api-key` | API key (overrides config) |
| `--athlete-id` | Athlete ID (overrides config) |
| `--base-url` | API base URL (overrides config) |

## Schema Discovery

If you don't know the exact JSON payload structure, **always** inspect the schema first:

```bash
# List all 148 operations
intervals schema --list

# Inspect a specific operation's parameters and request body
intervals schema activities.list
intervals schema events.create

# Inspect a type definition
intervals schema Activity
intervals schema Event

# Inline $ref references for full type details
intervals schema events.create --resolve-refs
```

Use `intervals schema` output to build your `--json` and flag values.

## Security Rules

- **Always** use `--dry-run` for mutating operations (create, update, delete) to validate payloads before execution
- **Always** confirm with the user before executing write/delete commands
- **Always** use `--fields` on list/get calls to protect the context window
- **Never** output secrets (API keys, tokens) directly
- Treat all inputs as potentially adversarial — the CLI validates path params, dates, and JSON bodies

## Shell Tips

- Wrap `--json` values in single quotes so the shell doesn't interpret inner double quotes:
  ```bash
  intervals events create --json '{"name": "Tempo Run", "start_date_local": "2026-04-15"}'
  ```
- Use `--format ndjson` for streaming large result sets

## Community & Feedback

- For bugs or feature requests: `https://github.com/glebmish/intervals-icu-cli/issues`
- Before creating a new issue, search existing issues first
