# Updaters

The updater component enables remote updating of System Link directly from Home Assistant. It will create an update entity that will show in Home Assistant's updates section.

## Available Updaters

### Update Component

- **[Update](update.md)** - Enable self-update capability for System Link

## Configuration Options

### Common Options

All updaters support these common options:

```yaml
updaters:
  - updater_type:
      name: "Custom Name"
      unique_id: "custom_id"
      availability: true
```

