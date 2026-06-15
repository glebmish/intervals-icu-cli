---
name: intervals-settings
description: "intervals CLI: Athlete profile, sport settings, gear, folders, chats, and more."
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

# Athlete Settings & Configuration

> **PREREQUISITE:** Read `../intervals-shared/SKILL.md` for auth, global flags, and security rules.

## athlete

Athlete profile, performance curves, fitness model, routes.

```bash
intervals athlete <action> [flags]
```

| Action | Description |
|--------|-------------|
| `get` | Get athlete profile |
| `update` | Update athlete profile (`--json`) |
| `profile` | Get athlete profile details |
| `summary` | Get athlete summary (`--start`, `--end`, `--tags`, `--ext`) |
| `training-plan` | Get athlete training plan |
| `update-training-plan` | Update training plan (`--json`) |
| `apply-plan-changes` | Apply pending plan changes |
| `weather-config` | Get weather configuration |
| `update-weather-config` | Update weather configuration (`--json`) |
| `weather-forecast` | Get weather forecast |
| `routes` | List athlete routes |
| `route-get` | Get a specific route (`--route-id`) |
| `route-update` | Update a route (`--route-id`, `--json`) |
| `route-similarity` | Get similarity between routes (`--route-id`, `--other-route-id`) |
| `chats` | List athlete chats |
| `fitness-model-events` | Get fitness model events (may return empty — use `wellness list` with `--fields ctl,atl` for reliable CTL/ATL data) |
| `power-hr-curve` | Get power/HR curve (`--start`, `--end`) |
| `power-curves` | Get power curves (`--oldest`, `--newest`, **`--type` required**, `--fatigue`, `--filters`, `--secs`, `--ext`) |
| `pace-curves` | Get pace curves (`--oldest`, `--newest`, **`--type` required**, `--filters`, `--distances`, `--gap`, `--ext`) |
| `hr-curves` | Get HR curves (`--oldest`, `--newest`, **`--type` required**, `--filters`, `--secs`, `--ext`) |
| `mmp-model` | Get MMP model |
| `activity-power-curves` | Get activity power curves (`--oldest`, `--newest`, `--type`, `--fatigue`, `--filters`, `--secs`, `--ext`) |
| `activity-pace-curves` | Get activity pace curves (`--oldest`, `--newest`, `--type`, `--filters`, `--distances`, `--gap`, `--ext`) |
| `activity-hr-curves` | Get activity HR curves (`--oldest`, `--newest`, `--type`, `--filters`, `--secs`, `--ext`) |
| `duplicate-events` | Duplicate events (`--json`) |
| `download-workout` | Download a workout file (`--workout-id`) |
| `download-fit-files` | Download FIT files (`--ids`) |
| `tags` | List athlete activity tags |

### Examples

```bash
# Get athlete profile
intervals athlete get --fields id,name,email,ftp,weight

# Get athlete summary for a date range (note: uses --start/--end, not --oldest/--newest)
intervals athlete summary --start 2026-01-01 --end 2026-04-11

# Get power curves over last 90 days (--type is required)
intervals athlete power-curves --oldest 2026-01-11 --newest 2026-04-11 --type Ride

# Get pace curves for running
intervals athlete pace-curves --oldest 2026-01-11 --newest 2026-04-11 --type Run

# Get CTL/ATL/TSB — use wellness list (more reliable than fitness-model-events)
intervals wellness list --oldest 2026-01-01 --newest 2026-04-11 --fields id,ctl,atl --format ndjson
```

## sport-settings

Sport-specific settings and zones.

```bash
intervals sport-settings <action> [flags]
```

| Action | Description |
|--------|-------------|
| `list` | List sport settings |
| `get` | Get a sport setting (`--sport-id`) |
| `get-device` | Get device settings (`--sport-id`) |
| `create` | Create sport settings (`--json`) |
| `update` | Update sport settings (`--sport-id`, `--json`) |
| `update-multi` | Update multiple sport settings (`--json`) |
| `delete` | Delete sport settings (`--sport-id`) |
| `apply` | Apply sport settings (`--sport-id`) |
| `pace-distances` | Get pace distances (`--sport-id`) |
| `matching-activities` | Get activities matching settings (`--sport-id`) |

## gear

Equipment tracking and reminders.

```bash
intervals gear <action> [flags]
```

| Action | Description |
|--------|-------------|
| `list` | List athlete gear |
| `create` | Create gear item (`--json`) |
| `update` | Update gear item (`--gear-id`, `--json`) |
| `delete` | Delete gear item (`--gear-id`) |
| `replace` | Replace gear item (`--gear-id`, `--json`) |
| `calc` | Calculate gear statistics (`--gear-id`) |
| `create-reminder` | Create gear reminder (`--gear-id`, `--json`) |
| `update-reminder` | Update gear reminder (`--gear-id`, `--reminder-id`, `--json`) |
| `delete-reminder` | Delete gear reminder (`--gear-id`, `--reminder-id`) |

### Examples

```bash
# List gear with distance
intervals gear list --fields id,name,distance,time

# Check gear statistics
intervals gear calc --gear-id g123
```

## folders

Workout folder management.

```bash
intervals folders <action> [flags]
```

| Action | Description |
|--------|-------------|
| `list` | List workout folders |
| `create` | Create a folder (`--json`) |
| `update` | Update a folder (`--folder-id`, `--json`) |
| `delete` | Delete a folder (`--folder-id`) |
| `update-workouts` | Update workouts in a folder (`--folder-id`, `--json`) |
| `shared-with` | Get folder sharing settings (`--folder-id`) |
| `update-shared-with` | Update folder sharing (`--folder-id`, `--json`) |
| `import-workout` | Import a workout into a folder (`--folder-id`, `--json`) |

## chats

Coach/athlete messaging.

```bash
intervals chats <action> [flags]
```

| Action | Description |
|--------|-------------|
| `get` | Get a chat by ID (`--chat-id`) |
| `messages` | Get messages in a chat (`--chat-id`) |
| `send` | Send a chat message (`--json`) |
| `delete-message` | Delete a chat message (`--chat-id`, `--message-id`) |
| `mark-seen` | Mark a message as seen (`--chat-id`, `--message-id`) |

## shared-events

Public shared events.

```bash
intervals shared-events <action> [flags]
```

| Action | Description |
|--------|-------------|
| `get` | Get a shared event by ID (`--event-id`) |
| `get-by-slug` | Get a shared event by slug (`--slug`) |
| `create` | Create a shared event (`--json`) |
| `update` | Update a shared event (`--event-id`, `--json`) |
| `delete` | Delete a shared event (`--event-id`) |
| `upload-image` | Upload an image for a shared event (`--shared-event-id`, `--upload`) |

## custom-items

Custom dashboard items.

```bash
intervals custom-items <action> [flags]
```

| Action | Description |
|--------|-------------|
| `list` | List custom dashboard items |
| `get` | Get a custom item (`--item-id`) |
| `create` | Create a custom item (`--json`) |
| `update` | Update a custom item (`--item-id`, `--json`) |
| `delete` | Delete a custom item (`--item-id`) |
| `upload-image` | Upload image for a custom item (`--item-id`, `--upload`) |
| `update-indexes` | Update custom item indexes (`--json`) |

## misc

Miscellaneous operations.

```bash
intervals misc <action> [flags]
```

| Action | Description |
|--------|-------------|
| `pace-distances` | Get pace distances |
| `download-workout` | Download a workout file (`--workout-id`) |
| `update-athlete-plans` | Update athlete plans (`--json`) |
| `disconnect-app` | Disconnect an integrated app (`--app-id`) |

## Discovering Commands

```bash
intervals schema --list                    # All operations
intervals schema gear.create --resolve-refs  # Operation details
intervals schema Gear                       # Type definition
```

> [!CAUTION]
> **Write commands** — always use `--dry-run` first and confirm with the user before executing.
