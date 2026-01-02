---
icon: material/puzzle
---

# Components Overview

System Link provides a variety of component types that integrate seamlessly with Home Assistant through MQTT Discovery. Each component type serves a specific purpose in monitoring or controlling your system.

## Component Types

### [Sensors](sensors/index.md)
Monitor numeric system metrics that update periodically.

**Examples**: CPU usage, memory usage, disk space, uptime, battery level

```yaml
sensors:
  - cpu: {}
  - memory: {}
  - disk: {}
```

### [Binary Sensors](binary-sensors/index.md)
Track boolean (on/off) states.

**Examples**: System online/offline status, locked state

```yaml
binary_sensors:
  - status: {}
```

### [Buttons](buttons/index.md)
Trigger one-time actions when pressed.

**Examples**: Shutdown, restart, lock screen, custom commands

```yaml
buttons:
  - shutdown: {}
  - lock: {}
  - custom:
      name: "Run Backup"
      command: "/scripts/backup.sh"
```

### [Switches](switches/index.md)
Toggle persistent states on and off.

**Examples**: Startup programs, scheduled cron jobs

```yaml
switches:
  - startup: {}
  - cron:
      name: "Daily Backup"
      schedule: "0 2 * * *"
      command: "/scripts/backup.sh"
```

### [Locks](locks/index.md)
Control state machines with lock/unlock commands.

**Examples**: Start/stop applications, complex state management

```yaml
locks:
  - custom:
      name: "Media Player"
      lock:
        command: "vlc /path/to/video.mp4"
      unlock:
        command: "pkill vlc"
```

### [Images](images/index.md)
Capture and stream images from your system.

**Examples**: Screenshots, webcam feeds

```yaml
images:
  - screen: {}
```

### [Updaters](updaters/index.md)
Enable remote updating of System Link itself.

**Examples**: Self-update capability through Home Assistant

```yaml
updaters:
  - update: {}
```

## Common Configuration Options

All components share some common configuration options:

```yaml
component_type:
  - name: "Custom Display Name"
    unique_id: "custom_unique_id"
    icon: "mdi:icon-name"
    availability: true
```

| Option | Type | Description | Default |
|--------|------|-------------|---------|
| `name` | string | Display name in Home Assistant | Component-specific |
| `unique_id` | string | Unique identifier for the entity | Auto-generated |
| `icon` | string | Material Design Icon name | Component-specific |
| `availability` | boolean | Track and report availability | `true` |

## Component Categories

### Monitoring Components

Components that observe and report system state:

- **Sensors** - Numeric measurements
- **Binary Sensors** - Boolean states
- **Images** - Visual information

### Control Components

Components that trigger actions or change system state:

- **Buttons** - One-time actions
- **Switches** - Persistent toggles
- **Locks** - Stateful controls

### Management Components

Components for system management:

- **Updaters** - Self-update capability

## Selective Component Enablement

One of System Link's key features is the ability to enable only the components you need. This:

- Reduces MQTT traffic
- Minimizes attack surface
- Keeps Home Assistant entity list clean
- Improves performance

**Example**: Only monitoring, no control

```yaml
sensors:
  - cpu: {}
  - memory: {}
  - disk: {}

binary_sensors:
  - status: {}

# No buttons, switches, or locks
```

## Custom Components

Many component types support custom variants where you can define your own commands:

### Custom Button

```yaml
buttons:
  - custom:
      name: "Restart Service"
      icon: mdi:restart
      unique_id: restart_nginx
      command: "systemctl restart nginx"
```

### Custom Sensor

```yaml
sensors:
  - custom:
      name: "Service Status"
      unique_id: nginx_status
      command: "systemctl is-active nginx"
      value_template: "{{ value }}"
```

### Custom Switch (Cron)

```yaml
switches:
  - cron:
      name: "Hourly Health Check"
      schedule: "0 * * * *"
      command: "/scripts/health-check.sh"
```

## Availability Tracking

Most components report their availability status to Home Assistant. If System Link disconnects, entities will show as "unavailable".

Disable availability tracking for specific components:

```yaml
images:
  - screen:
      availability: false  # Don't track availability
```

