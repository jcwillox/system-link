---
icon: material/rocket-launch
hide:
  - navigation
---

# Getting Started

This guide will help you install and configure System Link for your system.

## Installation

=== "Linux/macOS"
    ```bash
    sh -c "$(curl -fsSL jcwillox.com/l/syslink)"
    ```

=== "Windows"
    ```powershell
    iwr -useb jcwillox.com/l/syslink-ps1 | iex
    ```

The installation script will:

1. Download the latest System Link binary
2. Install it to `~/.local/share/system-link/`

**Configuration Location**: It will expect the config to be located at `~/.local/share/system-link/config.yaml`.

### Docker

For Docker deployments, pull the official image:

```bash
docker pull ghcr.io/jcwillox/system-link:python3-alpine
```

Run the container:

```bash
docker run -d \
  --name system-link \
  -v /path/to/config.yaml:/config/config.yaml \
  -v /path/to/logs:/logs \
  --restart unless-stopped \
  ghcr.io/jcwillox/system-link:python3-alpine
```

**Configuration Mounts**:

- `/config/config.yaml` - Mount your `config.yaml` file here
- `/logs` - Directory for System Link logs (optional)

## Basic Configuration

Create a `config.yaml` file in the System Link directory with the following:

```yaml
# device_name: "My PC"
# host_id: 5dc39758-0393-41ba-b019-f10172f5d373

mqtt:
  host: 127.0.0.1
  port: 1883
  # tls: true
  username: your_username
  password: your_password

# log_level: debug
# log_level_mqtt: debug
# log_timing: true

binary_sensors:
  - status: {}

buttons:
  - exit: {}
  - force_shutdown: {}
  - lock: {}
  - reload: {}
  - reset_topics: {}
  - shutdown: {}
  - sleep: {}

sensors:
  - cpu: {}

  - memory: {}
  - memory_used: {}
  - memory_free: {}

  - swap: {}
  - swap_used: {}
  - swap_free: {}

  - disk: {}
  - disk_used: {}
  - disk_free: {}

  - uptime: {}

  - battery: {}
  - battery_state: {}


switches:
  - startup: {}

updaters:
  - update: {}
```

### MQTT Configuration

System Link requires an MQTT broker to communicate. You have two options:

#### Option 1: Use Home Assistant's MQTT Addon

If you're using Home Assistant, the easiest way is to install the MQTT broker addon:

1. Open Home Assistant
2. Navigate to **Settings** → **Add-ons** → **Add-on Store**
3. Search for and install **Mosquitto Broker**
4. The addon provides MQTT credentials in its configuration
5. Add these credentials to your System Link `config.yaml`:

```yaml
mqtt:
  host: 192.168.1.100  # Your Home Assistant IP
  port: 1883
  username: mqtt_user
  password: mqtt_password
```

#### Option 2: Use External MQTT Broker

If using an external MQTT broker (e.g., mosquitto):

```yaml
mqtt:
  host: mqtt.example.com
  port: 1883  # or 8883 for TLS
  tls: false  # Set to true if using TLS
  username: system-link
  password: your_secure_password
```

### Home Assistant Integration

Once configured, System Link will automatically appear in Home Assistant via MQTT Discovery:

1. Open Home Assistant
2. Go to **Settings** → **Devices & Services** → **MQTT**
3. You should see your System Link device listed
4. All configured components (sensors, buttons, switches) will be automatically available

## Next Steps

- Check the [Security](../security/access-control.md) section to understand System Link's security features
- Explore [Components](../components/index.md) to see all available monitoring and control options
- Learn about [Advanced Features](../advanced/scheduled-tasks.md) like templating and scheduled tasks
- Configure custom [Buttons](../components/buttons/custom.md) and [Switches](../components/switches/cron.md) for your specific needs
