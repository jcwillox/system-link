---
icon: material/file-tree
---
# Splitting Configuration

The `!include` tag allows you to split your configuration across multiple files, making it easier to organise and maintain large configurations.

## Usage

Instead of putting everything in `config.yaml`, you can include content from other YAML files:

```yaml title="config.yaml"
host_id: my-device
device_name: My Device

mqtt:
  host: mqtt.example.com
  port: 1883

buttons: !include buttons.yaml
sensors: !include sensors.yaml
```

The included files should contain valid YAML for that section:

```yaml title="buttons.yaml"
- name: Custom Action
  icon: mdi:play
  command: echo "Hello World"
```

## File Paths

- Relative paths are resolved relative to the directory containing `config.yaml`
- Absolute paths are also supported
