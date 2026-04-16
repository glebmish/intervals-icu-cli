---
name: intervals-training
description: "intervals CLI: Activities, activity analysis, events, workouts, and wellness."
metadata:
  version: 0.1.0
  openclaw:
    category: "productivity"
    requires:
      bins:
        - intervals
      skills:
        - intervals-shared
---

# Training Data & Calendar

> **PREREQUISITE:** Read `../intervals-shared/SKILL.md` for auth, global flags, and security rules.

## activities

Athlete activities — list, search, create, upload.

```bash
intervals activities <action> [flags]
```

| Action | Description |
|--------|-------------|
| `list` | List athlete activities (`--oldest`, `--newest`, `--format-ext`) |
| `search` | Search activities by text (`--query`) |
| `search-full` | Search activities with full details (`--query`) |
| `interval-search` | Search by interval criteria (`--json` for search body) |
| `create-manual` | Create a manual activity (`--json`) |
| `create-manual-bulk` | Create multiple manual activities (`--json`) |
| `upload` | Upload an activity file FIT/TCX/GPX (`--file`) |
| `download-csv` | Download all activities as CSV (prints to stdout, redirect with `> file.csv`) |
| `list-around` | List activities around a specific activity (`--activity-id`) |
| `get-multiple` | Get multiple activities by IDs (`--ids`) |

### Examples

```bash
# List recent activities (always use --fields)
intervals activities list --oldest 2026-01-01 --fields id,name,start_date_local,type,icu_training_load

# Search for specific activities
intervals activities search --query "tempo" --fields id,name,start_date_local

# Download all activities as CSV for bulk analysis (preferred over paginated list calls)
intervals activities download-csv > /tmp/activities.csv

# Create a manual activity (use --dry-run first)
intervals activities create-manual --dry-run --json '{"name": "Morning Run", "type": "Run", "start_date_local": "2026-04-11T07:00:00", "elapsed_time": 3600}'
```

> **Bulk analysis tip:** For analysis spanning many activities, prefer `download-csv` over paginated `list` calls. The CSV contains all 88 fields for the entire history. Use `--format ndjson` with `list` only when you need a filtered subset.

## activity

Single activity analysis — power curves, HR, pace, streams, intervals.

```bash
intervals activity <action> --activity-id <ID> [flags]
```

| Action | Description |
|--------|-------------|
| `get` | Get an activity by ID |
| `update` | Update an activity (`--json`) |
| `delete` | Delete an activity |
| `intervals` | Get intervals for an activity |
| `update-intervals` | Update intervals (`--json`) |
| `update-interval` | Update a single interval (`--interval-id`, `--json`) |
| `delete-intervals` | Delete specific intervals (`--json`) |
| `split-interval` | Split an interval (`--interval-id`, `--json`) |
| `streams` | Get activity streams (`--types`) |
| `update-streams` | Update activity streams (`--json`) |
| `upload-streams-csv` | Upload streams as CSV (`--file`) |
| `power-curve` | Get power curve |
| `power-curves` | Get power curves |
| `power-histogram` | Get power histogram |
| `power-spike-model` | Get power spike model |
| `power-vs-hr` | Get power vs heart rate data |
| `hr-curve` | Get heart rate curve |
| `hr-histogram` | Get heart rate histogram |
| `hr-load-model` | Get HR load model |
| `pace-curve` | Get pace curve |
| `pace-histogram` | Get pace histogram |
| `gap-histogram` | Get grade-adjusted pace histogram |
| `interval-stats` | Get interval statistics |
| `segments` | Get segments |
| `map` | Get map data |
| `weather-summary` | Get weather summary |
| `time-at-hr` | Get time at heart rate zones |
| `best-efforts` | Get best efforts |
| `messages` | Get messages |
| `send-message` | Send a message (`--json`) |
| `fit-file` | Download FIT file |
| `gpx-file` | Download GPX file |
| `file` | Download original file |

### Examples

```bash
# Get activity overview (always use --fields)
intervals activity get --activity-id i12345 --fields id,name,type,icu_training_load,icu_ftp

# Get power curve for analysis
intervals activity power-curve --activity-id i12345

# Get activity streams
intervals activity streams --activity-id i12345 --types watts,heartrate,cadence

# Inspect schema before updating
intervals schema activity.update --resolve-refs
intervals activity update --activity-id i12345 --dry-run --json '{"name": "Updated Name"}'
```

## events

Training calendar events — CRUD, bulk operations, plans.

```bash
intervals events <action> [flags]
```

| Action | Description |
|--------|-------------|
| `list` | List events in date range (`--oldest`, `--newest`, `--category`, `--limit`) |
| `get` | Get a single event (`--event-id`) |
| `create` | Create a new event (`--json`) |
| `update` | Update an event (`--event-id`, `--json`) |
| `delete` | Delete an event (`--event-id`) |
| `update-bulk` | Bulk update events (`--json`) |
| `delete-range` | Delete events in date range (`--oldest`, `--newest`) |
| `delete-bulk` | Bulk delete events by IDs (`--json`) |
| `create-bulk` | Bulk create events (`--json`) |
| `mark-done` | Mark an event as done (`--event-id`) |
| `apply-plan` | Apply a training plan (`--json`) |
| `download-workout` | Download workout file for an event (`--event-id`) |
| `tags` | List event tags |

### Examples

```bash
# List this week's training plan
intervals events list --oldest 2026-04-07 --newest 2026-04-13 --fields id,name,start_date_local,category,icu_training_load

# Create a training event (use --dry-run first)
intervals events create --dry-run --json '{"name": "Threshold Intervals", "start_date_local": "2026-04-15T08:00:00", "category": "WORKOUT"}'

# Mark a workout as completed
intervals events mark-done --event-id 12345
```

## workouts

Workout library — CRUD, bulk, duplicate.

```bash
intervals workouts <action> [flags]
```

| Action | Description |
|--------|-------------|
| `list` | List workouts |
| `get` | Get a single workout (`--workout-id`) |
| `create` | Create a new workout (`--json`) |
| `update` | Update a workout (`--workout-id`, `--json`) |
| `delete` | Delete a workout (`--workout-id`) |
| `create-bulk` | Bulk create workouts (`--json`) |
| `duplicate` | Duplicate workouts (`--json`) |
| `download` | Download workouts as zip |
| `tags` | List workout tags |

### Examples

```bash
# List workout library
intervals workouts list --fields id,name,type,description

# Inspect workout schema before creating
intervals schema workouts.create --resolve-refs
```

## wellness

Daily wellness records.

```bash
intervals wellness <action> [flags]
```

| Action | Description |
|--------|-------------|
| `get` | Get wellness record for a date (`--date`) |
| `update` | Update wellness record (`--date`, `--json`) |
| `update-current` | Update today's wellness record (`--json`) |
| `update-bulk` | Bulk update wellness records (`--json`) |
| `upload` | Upload wellness data (`--file`) |
| `list` | List wellness records (`--oldest`, `--newest`) |

### Examples

```bash
# Get today's wellness
intervals wellness get --date 2026-04-11

# Log today's wellness (use --dry-run first)
intervals wellness update-current --dry-run --json '{"restingHR": 52, "weight": 72.5, "sleepQuality": 4}'

# List recent wellness trends
intervals wellness list --oldest 2026-04-01 --newest 2026-04-11 --fields id,restingHR,weight,sleepQuality,ctl,atl
```

## Discovering Commands

Before calling any operation, inspect it:

```bash
# Browse all operations
intervals schema --list

# Inspect a specific operation
intervals schema events.create --resolve-refs

# Inspect a type definition
intervals schema Activity
```

> [!CAUTION]
> **Write commands** (`create`, `update`, `delete`, `mark-done`, `apply-plan`) — always use `--dry-run` first and confirm with the user before executing.
