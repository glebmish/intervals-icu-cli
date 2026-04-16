---
name: recipe-training-review
description: "Review recent training: activities, wellness trends, and training load."
metadata:
  version: 0.1.0
  openclaw:
    category: "recipe"
    domain: "training"
    requires:
      bins:
        - intervals
      skills:
        - intervals-training
---

# Review Recent Training

> **PREREQUISITE:** Load the following skills to execute this recipe: `intervals-training`

Pull recent activities, wellness data, and training load to summarize the athlete's current training state. Adapt the date range to the user's request (default: last 7 days).

## Steps

### For short periods (up to ~30 days)

1. **Get recent activities:**
   ```bash
   intervals activities list --oldest <START> --newest <END> --fields id,name,start_date_local,type,icu_training_load,moving_time,distance
   ```

2. **Get wellness trends (includes CTL/ATL):**
   ```bash
   intervals wellness list --oldest <START> --newest <END> --fields id,restingHR,weight,sleepQuality,ctl,atl --format ndjson
   ```

3. **Get athlete summary** (note: uses `--start`/`--end`, not `--oldest`/`--newest`):
   ```bash
   intervals athlete summary --start <START> --end <END>
   ```

### For longer periods (months to years)

1. **Download all activities as CSV** (preferred over paginated list calls):
   ```bash
   intervals activities download-csv > /tmp/activities.csv
   ```
   The CSV contains all 88 fields for every activity. Use Python/scripts to analyze.
   Note: output has a UTF-8 BOM — use `encoding='utf-8-sig'` when parsing in Python.

2. **Download wellness data:**
   ```bash
   intervals wellness list --oldest <START> --newest <END> --fields id,restingHR,weight,sleepQuality,ctl,atl,sleepSecs --format ndjson > /tmp/wellness.ndjson
   ```

### Summarize for the user

- Total training load and volume (hours, distance)
- Activity breakdown by type (run, ride, swim, etc.)
- Wellness trends (resting HR, sleep, weight)
- Current fitness (CTL), fatigue (ATL), and form (TSB)
- Compare current CTL to previous peaks to contextualize fitness level
- Any notable patterns, gaps, or concerns

## Tips

- Adjust date range based on user request: "last week", "this month", "last 30 days"
- Use `--fields` to keep responses concise
- Compare current CTL/ATL to previous periods if the user asks about trends
- If the user asks to dig into a specific activity, use `intervals activity get` or analysis commands (power-curve, streams, etc.)
- For power/pace curve analysis, remember `--type` is required: `intervals athlete power-curves --type Ride --oldest ... --newest ...`
