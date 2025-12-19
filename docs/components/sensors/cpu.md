# CPU Sensor

Monitor CPU utilization as a percentage.

## Configuration

```yaml
sensors:
  - cpu:
      update_interval: 10s
```

## Options

All [common sensor options](index.md#common-options) are supported.

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
