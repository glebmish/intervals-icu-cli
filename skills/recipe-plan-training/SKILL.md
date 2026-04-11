---
name: recipe-plan-training
description: "Plan and schedule training: create calendar events and apply workout plans."
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

# Plan Training

> **PREREQUISITE:** Load the following skills to execute this recipe: `intervals-training`

Create training events on the calendar. Can be used for planning a single session, a training week, or applying a structured plan.

## Steps

1. **Check existing schedule:**
   ```bash
   intervals events list --oldest <START> --newest <END> --fields id,name,start_date_local,category,icu_training_load
   ```

2. **Check current fitness state:**
   ```bash
   intervals athlete summary --oldest <30_DAYS_AGO> --newest <TODAY>
   ```

3. **Discover event schema:**
   ```bash
   intervals schema events.create --resolve-refs
   ```

4. **Create events (use --dry-run first):**
   ```bash
   # Single event
   intervals events create --dry-run --json '{"name": "Easy Run", "start_date_local": "2026-04-15", "category": "WORKOUT"}'

   # Multiple events at once
   intervals events create-bulk --dry-run --json '[{"name": "Intervals", "start_date_local": "2026-04-16"}, {"name": "Long Run", "start_date_local": "2026-04-17"}]'
   ```

5. **Confirm with user, then execute without --dry-run.**

6. **Verify the plan:**
   ```bash
   intervals events list --oldest <START> --newest <END> --fields id,name,start_date_local,category
   ```

## Tips

- **Always** dry-run first, show the user what will be created, and get confirmation
- Check existing events before creating to avoid double-booking
- Use `intervals workouts list` to browse the workout library for reuse
- Use `events apply-plan` to apply a structured training plan from the library
- For bulk creation, prefer `events create-bulk` over multiple individual creates
