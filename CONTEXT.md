# intervals CLI Context

The `intervals` CLI provides access to the intervals.icu training analytics API. Designed for AI agents.

## Rules of Engagement for Agents

* **Schema Discovery:** If you don't know the exact JSON payload structure, run `intervals schema <resource>.<method>` first.
* **Context Window Protection:** Use `--fields` on every list/get call to limit response size.
* **Dry-Run Safety:** Always use `--dry-run` for mutating operations to validate payloads before execution.

## Core Syntax

```bash
intervals <resource> <action> [flags]
```

### Key Flags

- `--json '<JSON>'`: Request body for POST/PUT
- `--fields '<mask>'`: Comma-separated fields to include in output
- `--format <fmt>`: Output format (json, ndjson, text)
- `--dry-run`: Validate without executing
- `--yes`: Skip confirmation on deletes

## Usage Patterns

### Reading Data
```bash
intervals activities list --oldest 2026-01-01 --fields id,name,start_date_local
intervals activity get --activity-id a1 --fields id,name,type
intervals wellness get --date 2026-04-11
```

### Writing Data
```bash
intervals events create --json '{"name": "Threshold", "start_date_local": "2026-04-15"}'
intervals activity update --activity-id a1 --json '{"name": "Updated Name"}'
```

### Schema Discovery
```bash
intervals schema --list              # All operations
intervals schema activities.list     # Operation details
intervals schema Activity            # Type definition
```

## Resource Groups

- `activities` — List, search, create, upload athlete activities
- `activity` — Analyze a specific activity (power, HR, pace, streams, intervals)
- `athlete` — Profile, settings, curves, routes, fitness model
- `events` — Training calendar events (CRUD, bulk, plans)
- `workouts` — Workout library (CRUD, bulk, duplicate)
- `wellness` — Daily wellness records
- `gear` — Equipment tracking and reminders
- `folders` — Workout folder management
- `sport-settings` — Sport-specific settings and zones
- `chats` — Coach/athlete messaging
- `shared-events` — Public shared events
- `custom-items` — Custom dashboard items
