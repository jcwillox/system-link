# Binary Sensors

Binary sensors track boolean (on/off) states from your system.

## Available Binary Sensors

### Built-in Binary Sensors

- **[Status](status.md)** - Reports whether System Link is online and connected

## Configuration Options

### Common Options

All binary sensors support these common options:

```yaml
binary_sensors:
  - sensor_type:
      name: "Custom Name"
      unique_id: "custom_id"
      icon: "mdi:icon-name"
      availability: true
```
