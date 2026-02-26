# Cron Switches

Cron switches are scheduled tasks that run on a schedule and can be enabled/disabled.

Check out the [Scheduled Tasks](../../advanced/scheduled-tasks.md) section for more info.

## Basic Cron Switch

```yaml
switches:
  - cron:
      name: "Daily Backup"
      unique_id: daily_backup
      schedule: "0 2 * * *"  # Daily at 2 AM
      command: "/scripts/backup.sh"
```

## Cron Switch Configuration

```yaml
switches:
  - cron:
      name: "Scheduled Task"
      unique_id: "scheduled_task"
      icon: "mdi:clock-outline"
      schedule: "*/5 * * * *"    # Cron expression (required)
      command: "echo 'Hello'"    # Command to run (required)
      shell: "bash"              # Shell to use
      detached: false            # Run in background
      show_output: true          # Log command output
      show_errors: true          # Log command errors
      env:                       # Environment variables
        MY_VAR: "value"
      entities:                  # Optional sub-entities
        duration: {}             # Task execution time
        run: {}                  # Button to run immediately
        exit_code: {}            # Last exit code
        successful: {}           # Was last run successful
        output: {}               # Captured stdout/stderr from last run
        next_run: {}             # Time of next scheduled run
        last_run: {}             # Time of last run
```

## Cron Expression Format

System Link uses the standard cron format:

```
┌───────────── second (0-59)
│ ┌───────────── minute (0-59)
│ │ ┌───────────── hour (0-23)
│ │ │ ┌───────────── day of month (1-31)
│ │ │ │ ┌───────────── month (1-12)
│ │ │ │ │ ┌───────────── day of week (0-6) (Sunday=0)
│ │ │ │ │ │
* * * * * *
```

**Special Characters**:
- `*` - Any value
- `,` - Value list (1,3,5)
- `-` - Range (1-5)
- `/` - Step values (*/5)

## Cron Examples

```yaml
# Every minute
schedule: "* * * * *"

# Every 5 minutes
schedule: "*/5 * * * *"

# Every hour at minute 30
schedule: "30 * * * *"

# Daily at 2:30 AM
schedule: "30 2 * * *"

# Every Monday at 9 AM
schedule: "0 9 * * 1"

# First day of month at midnight
schedule: "0 0 1 * *"

# Every 30 seconds
schedule: "*/30 * * * * *"
```

## Cron Switch Sub-Entities

Enable additional entities to monitor task execution:

```yaml
switches:
  - cron:
      name: "Health Check"
      unique_id: health_check
      schedule: "*/5 * * * *"
      command: "/scripts/health-check.sh"
      entities:
        duration: {}      # Sensor: execution duration in seconds
        run: {}           # Button: run task immediately
        exit_code: {}     # Sensor: last exit code
        output: {}        # Sensor: captured output from last execution
        successful: {}    # Binary sensor: was last run successful
        next_run: {}      # Sensor: timestamp of next run
        last_run: {}      # Sensor: timestamp of last run
```

### Duration Sensor

Tracks how long the command takes to execute.

**Unit**: seconds
**Icon**: mdi:timer

### Run Button

Manually trigger the task outside the schedule.

**Icon**: mdi:play

### Exit Code Sensor

Reports the exit code of the last execution.

**Values**: 0 = success, non-zero = error
**Icon**: mdi:numeric

### Successful Binary Sensor

Indicates whether the last run completed successfully (exit code 0).

**Icon**: mdi:check-circle / mdi:alert-circle

### Next Run Sensor

Timestamp of when the task will run next.

**Device Class**: timestamp
**Icon**: mdi:clock-outline

### Last Run Sensor

Timestamp of when the task last ran.

**Device Class**: timestamp
**Icon**: mdi:clock-check-outline

### Output Sensor

Captures the combined and individual stdout/stderr output from the last command execution. Useful for quick inspection of what the command printed without needing to fetch logs.

**Type**: string
**Icon**: mdi:terminal

Notes:
- Output may be truncated in the frontend if very large.
- Use the `show_output` flag on the cron switch to enable logging of full output to logs in addition to this sensor.
- The output is also shown on the `successful` entity but only when the task fails, to help with debugging.

## Use Cases

### Scheduled Backups

```yaml
switches:
  - cron:
      name: "Nightly Backup"
      unique_id: nightly_backup
      icon: mdi:backup-restore
      schedule: "0 2 * * *"  # 2 AM daily
      shell: "bash"
      command: |
        rsync -av /data/ /backup/
        find /backup/ -mtime +7 -delete
      entities:
        duration: {}
        successful: {}
        last_run: {}
```

### Health Checks

```yaml
switches:
  - cron:
      name: "Service Health Check"
      unique_id: service_health
      schedule: "*/5 * * * *"  # Every 5 minutes
      command: |
        systemctl is-active nginx || systemctl restart nginx
      entities:
        successful: {}
        next_run: {}
```

### Log Rotation

```yaml
switches:
  - cron:
      name: "Rotate Logs"
      unique_id: rotate_logs
      schedule: "0 0 * * 0"  # Weekly on Sunday
      shell: "bash"
      command: |
        find /var/log/app -name "*.log" -mtime +30 -delete
        gzip /var/log/app/*.log.1
```

### Database Cleanup

```yaml
switches:
  - cron:
      name: "Database Cleanup"
      unique_id: db_cleanup
      schedule: "0 3 * * *"  # Daily at 3 AM
      command: "psql -c 'DELETE FROM logs WHERE created_at < NOW() - INTERVAL 90 days'"
      entities:
        duration: {}
        successful: {}
```

## Docker and Scheduled Scripts

System Link's Docker image includes Python and [UV](https://github.com/astral-sh/uv), making it perfect for running Python scripts with dependencies.

### Python Script with UV

```yaml
switches:
  - cron:
      name: "Python Health Check"
      unique_id: python_health
      schedule: "*/10 * * * *"
      command: "uv run /scripts/health_check.py"
      entities:
        successful: {}
        duration: {}
```

Example Python script with inline dependencies:

```python
# /// script
# dependencies = [
#   "requests",
#   "pydantic",
# ]
# ///

import requests
from pydantic import BaseModel

class HealthResponse(BaseModel):
    status: str
    uptime: int

def check_health():
    response = requests.get("http://localhost:8080/health")
    health = HealthResponse(**response.json())

    if health.status != "healthy":
        raise Exception(f"Service unhealthy: {health.status}")

    print(f"✓ Service healthy, uptime: {health.uptime}s")

if __name__ == "__main__":
    check_health()
```

See [Scheduled Tasks](../../advanced/scheduled-tasks.md) for more Python scripting examples.

## Template Variables

Cron switches have access to template variables:

```yaml
switches:
  - cron:
      name: "Backup with Timestamp"
      unique_id: timestamped_backup
      schedule: "0 2 * * *"
      shell: "bash"
      command: |
        BACKUP_DIR="{{ .ExeDirectory }}/backups"
        DATE=$(date +%Y%m%d)
        tar -czf "${BACKUP_DIR}/backup-${DATE}.tar.gz" /data/
```

See [Template Variables Reference](../../advanced/templating.md).

## Monitoring Cron Tasks

### Dashboard Card

```yaml
type: entities
title: Scheduled Tasks
entities:
  - entity: switch.nightly_backup
    name: Backup
  - entity: sensor.nightly_backup_last_run
    name: Last Run
  - entity: binary_sensor.nightly_backup_successful
    name: Status
  - entity: sensor.nightly_backup_duration
    name: Duration
  - entity: sensor.nightly_backup_output
    name: Output
  - entity: button.nightly_backup_run
    name: Run Now
```

### Automation: Alert on Failure

```yaml
automation:
  - alias: "Alert Backup Failed"
    trigger:
      - platform: state
        entity_id: binary_sensor.nightly_backup_successful
        to: "off"
    action:
      - service: notify.mobile_app
        data:
          message: "Nightly backup failed!"
          data:
            priority: high
```
