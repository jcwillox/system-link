---
icon: material/home
hide:
  - navigation
---

# Overview

!!! note ""

    **Lightweight, secure system monitoring and control via MQTT for Home Assistant and more.**

System Link is a powerful tool for monitoring and controlling your systems, it integrates natively with Home Assistant via [MQTT Discovery](https://www.home-assistant.io/integrations/mqtt/#mqtt-discovery), communicates securely over MQTT, and is a single lightweight binary you can deploy anywhere.

## ✨ Features

### 🔌 Component Types

- **Sensors** – Monitor system metrics like CPU, memory, disk, uptime, battery, and ZFS pools
- **Binary Sensors** – Track system status and states
- **Buttons** – Execute actions like shutdown, sleep, lock, reload, and custom commands
- **Switches** – Toggle system behavior including startup programs and scheduled cron jobs
- **Locks** – Control complex state machines with lock/unlock commands
- **Images** – Capture and stream screenshots from your system
- **Updaters** – Keep System Link up to date remotely through Home Assistant

### 🚀 Key Capabilities

- **Filters** – Reduce unnecessary updates with averaging, throttling, and delta filters
- **Remote Updates** – Update System Link directly from Home Assistant
- **Custom Commands** – Define your own buttons, switches, and sensors with custom scripts
- **Secure by Design** – Preconfigured commands only, no arbitrary command execution, and opt-in to what components you want to expose.
- **Access Control** – Fine-grained MQTT ACL support for restricting device permissions
- **Encrypted Communication** – Full TLS support with hostname verification

### 🎯 Use Cases

- **System Monitoring** – Track CPU, memory, disk usage, and more from Home Assistant
- **Remote Control** – Safely shutdown, restart, or lock computers remotely
- **Automation Platform** – Run scheduled tasks and health checks with cron switches

## Quick Start

Get started in minutes:

=== "Linux/macOS"
    ```bash
    sh -c "$(curl -fsSL jcwillox.com/l/syslink)"
    ```

=== "Windows"
    ```powershell
    iwr -useb jcwillox.com/l/syslink-ps1 | iex
    ```

=== "Docker"
    ```bash
    docker pull ghcr.io/jcwillox/system-link:python3-alpine
    ```

Then create a `config.yaml` file with your MQTT broker details:

```yaml
mqtt:
  host: 127.0.0.1
  port: 1883
  username: your_username
  password: your_password

buttons:
  - shutdown: {}

sensors:
  - cpu: {}
  - memory: {}
```

See the [Getting Started Guide](./getting-started/index.md) for detailed instructions.
