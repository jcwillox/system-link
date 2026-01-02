# Sensors

Sensors monitor numeric values and metrics from your system. They update periodically and can have filters applied to reduce MQTT traffic.

## Available Sensors

### Built-in Sensors

- **[CPU](cpu.md)** - Monitor CPU utilization as a percentage
- **[Memory](memory.md)** - Monitor RAM usage (percentage, used, free)
- **[Disk](disk.md)** - Monitor disk space usage (percentage, used, free)
- **[Swap](swap.md)** - Monitor swap/page file usage (percentage, used, free)
- **[ZFS Pool](zpool.md)** - Monitor ZFS pool usage (percentage, used, free)
- **[Uptime](uptime.md)** - Track system uptime
- **[Battery](battery.md)** - Monitor battery level and charging state

### Custom Sensors

- **[Custom Sensor](custom.md)** - Create your own sensors by running custom commands

## Configuration Options

### Common Options

All sensors support these common options:

```yaml
sensors:
  - sensor_type:
      name: "Custom Name"
      unique_id: "custom_id"
      icon: "mdi:icon-name"
      availability: true
      update_interval: 10s
      filters: []
```

### Update Interval

Control how often sensors check for updates:

```yaml
sensors:
  - cpu:
      update_interval: 5s   # Check every 5 seconds
  - disk:
      update_interval: 5m   # Check every 5 minutes
```

Valid time units: `s` (seconds), `m` (minutes), `h` (hours)

### Filters

Apply filters to reduce update frequency and control when values are published:

```yaml
sensors:
  - cpu:
      update_interval: 3s
      filters:
        - throttle_average: 30s  # Average over 30 seconds
        - or:
            - throttle: 5m       # Publish at least every 5 minutes
            - delta: 5           # Or when value changes by 5%
```

See [Filters Documentation](../../advanced/filters.md) for detailed information.

### Device Classes

Custom sensors can use Home Assistant device classes, see [available device classes](https://developers.home-assistant.io/docs/core/entity/sensor/#available-device-classes).

```yaml
sensors:
  - custom:
      device_class: "temperature"
      unit_of_measurement: "°C"
```

### State Classes

Enable long-term statistics in Home Assistant, see [available state classes](https://developers.home-assistant.io/docs/core/entity/sensor/#available-state-classes).

```yaml
sensors:
  - custom:
      state_class: "measurement"
```

## Examples

### Efficient CPU Monitoring

```yaml
sensors:
  - cpu:
      update_interval: 3s
      filters:
        - throttle_average: 30s
        - or:
            - throttle: 5m
            - delta: 5
```

This configuration:

- Checks CPU every 3 seconds
- Averages readings over 30 seconds
- Publishes if value changes by 5% OR every 5 minutes

### Complete System Monitoring

```yaml
sensors:
  - cpu:
      update_interval: 10s
      filters:
        - throttle_average: 30s
        - or:
            - throttle: 5m
            - delta: 5

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
```
