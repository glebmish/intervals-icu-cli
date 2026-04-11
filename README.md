# intervals-icu-cli

A command-line interface for the [intervals.icu](https://intervals.icu) training analytics API. Covers 148 operations across activities, events, workouts, wellness, and more. Designed for AI-agent use with schema discovery, field masking, and dry-run support.

## Installation

```bash
go install github.com/glebmish/intervals-icu-cli@latest
```

## Configuration

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

## Quick Start

```bash
# List recent activities
intervals activities list --oldest 2026-01-01 --fields id,name,start_date_local

# Get a specific activity
intervals activity get --activity-id a1 --fields id,name,type,distance

# Create a calendar event
intervals events create --json '{"name": "Threshold Run", "start_date_local": "2026-04-15"}'

# Discover available operations
intervals schema --list

# Inspect a specific operation's parameters and payload shape
intervals schema activities.list
intervals schema Activity
```

## Agent Reference

See [CONTEXT.md](CONTEXT.md) for the full agent reference, including rules of engagement, usage patterns, and resource groups.
