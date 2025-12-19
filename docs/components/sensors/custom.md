# Custom Sensors

Create your own sensors by running custom commands.

## Basic Custom Sensor

```yaml
sensors:
  - custom:
      name: "Service Status"
      unique_id: nginx_status
      command: "systemctl is-active nginx"
```

## Custom Sensor with JSON Parsing

```yaml
sensors:
  - custom:
      name: "API Response Time"
      unique_id: api_response
      command: "curl -s https://api.example.com/health"
      value_template: "{{ value_json.response_time }}"
      json_attributes:
        - status
        - latency
        - version
      unit_of_measurement: "ms"
```

## Custom Sensor Configuration

```yaml
sensors:
  - custom:
      name: "Custom Sensor"
      unique_id: "custom_sensor"
      icon: "mdi:gauge"
      command: "echo 42"
      shell: "bash"              # Shell to use (bash, powershell, cmd, etc.)
      update_interval: 30s       # How often to run the command
      value_template: "{{ value }}"  # Template for extracting value
      unit_of_measurement: "units"   # Unit to display
      device_class: "temperature"    # Device class for Home Assistant
      state_class: "measurement"     # State class for statistics
      show_output: false         # Log command output
      show_errors: false         # Log command errors
      json_attributes:           # Extract JSON attributes
        - attr1
        - attr2
      json_attributes_path: "$.nested"  # JSON path to attributes
      env:                       # Environment variables
        MY_VAR: "value"
```

## Template Variables

Custom sensors commands can be templated using Go templates:

```yaml
sensors:
  - custom:
      command: |
        echo "{{ .ExeDirectory }}"
        echo "{{ .HostID }}"
        echo "{{ .DeviceName }}"
```

See [Template Variables Reference](../../advanced/templating.md) for all available variables.

## Custom Sensor Examples

### Monitor Service Status

```yaml
sensors:
  - custom:
      name: "Docker Status"
      unique_id: docker_status
      command: "systemctl is-active docker"
      update_interval: 60s
```

### Monitor GPU Temperature

=== "NVIDIA"
    ```yaml
    sensors:
      - custom:
          name: "GPU Temperature"
          unique_id: gpu_temp
          command: "nvidia-smi --query-gpu=temperature.gpu --format=csv,noheader"
          unit_of_measurement: "°C"
          device_class: "temperature"
          update_interval: 10s
    ```

=== "AMD"
    ```yaml
    sensors:
      - custom:
          name: "GPU Temperature"
          unique_id: gpu_temp
          command: "cat /sys/class/hwmon/hwmon0/temp1_input"
          value_template: "{{ value | int / 1000 }}"
          unit_of_measurement: "°C"
          device_class: "temperature"
    ```

### Monitor Container Stats

```yaml
sensors:
  - custom:
      name: "Container Memory"
      unique_id: nginx_memory
      shell: "bash"
      command: |
        docker stats --no-stream --format "{{.MemUsage}}" nginx | cut -d'/' -f1
      value_template: "{{ value | replace('MiB', '') | float }}"
      unit_of_measurement: "MiB"
      update_interval: 30s
```

## Device Classes

Custom sensors can use Home Assistant device classes, see [available device classes](https://developers.home-assistant.io/docs/core/entity/sensor/#available-device-classes).

```yaml
sensors:
  - custom:
      device_class: "temperature"
      unit_of_measurement: "°C"
```

## State Classes

Enable long-term statistics in Home Assistant, see [available state classes](https://developers.home-assistant.io/docs/core/entity/sensor/#available-state-classes).

```yaml
sensors:
  - custom:
      state_class: "measurement"
```
