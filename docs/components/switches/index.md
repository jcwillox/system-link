# Switches

Switches toggle persistent states on and off. Unlike buttons, switches maintain their state and can be turned on or off.

## Available Switches

### Built-in Switches

- **[Startup](startup.md)** - Control whether System Link starts automatically when the system boots

### Custom Switches

- **[Cron](cron.md)** - Scheduled tasks that run on a schedule and can be enabled/disabled

## Configuration Options

### Common Options

All switches support these common options:

```yaml
switches:
  - switch_type:
      name: "Custom Name"
      unique_id: "custom_id"
      icon: "mdi:icon-name"
      availability: true
```

