# Buttons

Buttons trigger one-time actions when pressed. They execute a command and return to their default state.

## Available Buttons

### Built-in Buttons

- **[Shutdown](shutdown.md)** – Gracefully shut down the system
- **[Force Shutdown](force-shutdown.md)** – Force immediate shutdown without saving
- **[Lock](lock.md)** – Lock the screen/session
- **[Sleep](sleep.md)** – Put the system to sleep/suspend
- **[Reload](reload.md)** - Reload System Link configuration (restarts system-link process)
- **[Exit](exit.md)** – Stops and exits System Link
- **[Reset Topics](reset-topics.md)** – Clear all MQTT topics published by System Link that are no longer used

### Custom Buttons

- **[Custom Button](custom.md)** – Create buttons that run your own commands

## Configuration Options

### Common Options

All buttons support these common options:

```yaml
buttons:
  - button_type:
      name: "Custom Name"
      unique_id: "custom_id"
      icon: "mdi:icon-name"
      availability: true
```


