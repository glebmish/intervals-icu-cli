# intervals-icu-cli

[![CI](https://github.com/glebmish/intervals-icu-cli/actions/workflows/ci.yml/badge.svg)](https://github.com/glebmish/intervals-icu-cli/actions/workflows/ci.yml)
[![Release](https://img.shields.io/github/v/release/glebmish/intervals-icu-cli)](https://github.com/glebmish/intervals-icu-cli/releases/latest)
[![License: MIT](https://img.shields.io/github/license/glebmish/intervals-icu-cli)](LICENSE)

A command-line interface for the [intervals.icu](https://intervals.icu) training analytics API. Covers 148 operations across activities, events, workouts, wellness, and more. Designed for AI-agent use with schema discovery, field masking, and dry-run support.

> Not affiliated with or endorsed by intervals.icu. This is an independent client for its public API.

## Installation

### Homebrew

```bash
brew install glebmish/tap/intervals
```

Use the fully-qualified `glebmish/tap/intervals` form (the bare name `intervals` is generic).

### Prebuilt binary

Download the archive for your platform from the [latest release](https://github.com/glebmish/intervals-icu-cli/releases/latest), extract it, and put the `intervals` binary on your `PATH`.

### go install

```bash
go install github.com/glebmish/intervals-icu-cli@latest
```

This installs a binary named `intervals-icu-cli` (after the module path). Alias or rename it so the examples below work:

```bash
alias intervals=intervals-icu-cli
# or: mv "$(go env GOPATH)/bin/intervals-icu-cli" "$(go env GOPATH)/bin/intervals"
```

## Configuration

Run the interactive setup:

```bash
intervals config init
```

This writes `~/.config/intervals/config.yaml` with `0600` permissions. You can also create it manually:

```yaml
api_key: your_api_key_here
athlete_id: your_athlete_id_here
```

Get an API key at <https://intervals.icu/settings>.

### Environment Variables

| Variable | Description |
|---|---|
| `INTERVALS_API_KEY` | API key (overrides config file) |
| `INTERVALS_ATHLETE_ID` | Athlete ID (overrides config file) |
| `INTERVALS_BASE_URL` | API base URL (default: `https://intervals.icu`) |
| `INTERVALS_CONFIG` | Path to the config file (overrides default `~/.config/intervals/config.yaml`) |

## Quick Start

```bash
# List recent activities
intervals activities list --oldest 2026-01-01 --fields id,name,start_date_local

# Get a specific activity
intervals activity get --activity-id a1 --fields id,name,type,distance

# Create a calendar event
intervals events create --json '{"name": "Threshold Run", "start_date_local": "2026-04-15T08:00:00"}'

# Discover available operations
intervals schema --list

# Inspect a specific operation's parameters and payload shape
intervals schema activities.list
intervals schema Activity
```

`intervals schema --list` enumerates all 148 operations across 13 resource groups, offline and without credentials:

```
activities.create-manual                  POST /api/v1/athlete/{id}/activities/manual
activities.download-csv                   GET /api/v1/athlete/{id}/activities.csv
activities.get-multiple                   GET /api/v1/athlete/{athleteId}/activities/{ids}
activities.interval-search                GET /api/v1/athlete/{id}/activities/interval-search
activities.list                           GET /api/v1/athlete/{id}/activities
...
```

Every command supports `--dry-run` (print the request without sending it), `--fields` (filter output), and `--format json|ndjson|text`. Run `intervals --help` for the full surface and `intervals --version` for build info.

## Agent Reference

Agent-facing documentation ships inside the binary as skills. Browse or install them:

```bash
intervals skills list                 # list bundled skills
intervals skills get intervals-shared # print the shared reference (global flags, exit codes, schema discovery)
intervals skills install              # install skills into .claude/skills or .agents/skills
```

The `intervals-shared` skill is the canonical agent front door: global flags, exit codes, error hints, and schema discovery.

## How it works

See [docs/design.md](docs/design.md) for the design: the embedded OpenAPI spec, the shared request helpers, input validation, and output sanitization.

## License

MIT — see [LICENSE](LICENSE).
