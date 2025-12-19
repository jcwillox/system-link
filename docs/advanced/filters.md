---
icon: material/tune
---

# Filters

Filters allow you to control when sensor values are published to MQTT, reducing network traffic and database load in Home Assistant.

This filters implementation is based on ESPHome's [sensor filters](https://esphome.io/components/sensor/#sensor-filters), we support a subset of those filters.

Filters are applied in the order they are defined in your configuration, i.e. the output of one filter is the input to the next.

## Why Use Filters?

Without filters, sensors publish every time they update:

```yaml
sensors:
  - cpu:
      update_interval: 3s  # Publishes 1,200 times per hour!
```

With filters, you control when updates are sent:

```yaml
sensors:
  - cpu:
      update_interval: 3s
      filters:
        - throttle: 5m     # Publishes max 12 times per hour
```

**Benefits**:

- ⬇️ Reduced MQTT traffic

- ⬇️ Fewer database writes

- ⬇️ Extended SSD lifespan

- ⬇️ Lower network bandwidth usage

## Available Filters

### Throttle

Limit update frequency to a minimum interval.

```yaml
filters:
  - throttle: 5m  # Publish at most every 5 minutes
```

### Delta

Publish only when value changes by a specified amount.

```yaml
filters:
  - delta: 5  # Publish only if value changes by 5
```

### Throttle Average

Average readings over a time period before applying other filters.

```yaml
filters:
  - throttle_average: 30s  # Average values over 30 seconds
```

### OR Filter

Publish if ANY condition is met.

```yaml
filters:
  - or:
      - throttle: 5m   # Publish every 5 minutes
      - delta: 10      # OR if value changes by 10
```

## Examples

### CPU Usage

Frequent sampling with smart reporting:

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

**Result**:

1. Check every 3 seconds
2. Average over 30 seconds
3. Publish if 5 minutes passed OR value changed by 5%

**Use Case**: Regular updates for graphs + immediate notification of spikes.

### Memory Usage

Less frequent updates, change-based reporting:

```yaml
sensors:
  - memory:
      update_interval: 30s
      filters:
        - or:
            - throttle: 10m
            - delta: 5
```

- Checks: Every 30 seconds
- Publishes: Every 10 minutes OR when change > 5%

### Disk Usage

Infrequent checks, large change threshold:

```yaml
sensors:
  - disk:
      update_interval: 5m
      filters:
        - or:
            - throttle: 1h
            - delta: 2
```

- Checks: Every 5 minutes
- Publishes: Every hour OR when change > 2%

### Uptime

Very infrequent, no need for frequent updates:

```yaml
sensors:
  - uptime:
      update_interval: 1h
      filters:
        - throttle: 6h
```

- Checks: Every hour
- Publishes: Every 6 hours

### Battery

Moderate frequency with change detection:

```yaml
sensors:
  - battery:
      update_interval: 30s
      filters:
        - or:
            - throttle: 5m
            - delta: 1
```

- Checks: Every 30 seconds
- Publishes: Every 5 minutes OR when change > 1%

## Filter Behavior Examples

??? info "Filter Behavior Examples"

    ### Example 1: Throttle Only

    ```yaml
    filters:
      - throttle: 5m
    ```

    Timeline:
    ```
    00:00 - Value: 30% → Publish (first value)
    00:01 - Value: 35% → Skip (< 5m)
    00:02 - Value: 80% → Skip (< 5m)
    00:05 - Value: 75% → Publish (5m passed)
    00:10 - Value: 75% → Publish (5m passed, even if unchanged)
    ```

    ### Example 2: Delta Only

    ```yaml
    filters:
      - delta: 10
    ```

    Timeline:
    ```
    00:00 - Value: 30% → Publish (first value)
    00:01 - Value: 32% → Skip (change = 2)
    00:02 - Value: 35% → Skip (change = 5)
    00:03 - Value: 42% → Publish (change = 12)
    00:04 - Value: 43% → Skip (change = 1)
    ```

    ### Example 3: Throttle OR Delta

    ```yaml
    filters:
      - or:
          - throttle: 5m
          - delta: 10
    ```

    Timeline:
    ```
    00:00 - Value: 30% → Publish (first value)
    00:01 - Value: 32% → Skip (< 5m AND change < 10)
    00:02 - Value: 42% → Publish (change = 12)
    00:03 - Value: 43% → Skip (< 5m AND change < 10)
    00:07 - Value: 44% → Publish (5m passed since last publish)
    ```

    ### Example 4: Average + Throttle

    ```yaml
    filters:
      - throttle_average: 30s
      - throttle: 5m
    ```

    Values over 30s: 20%, 25%, 80%, 30%, 25% → Average: 36%

    Timeline:
    ```
    00:00 - Average: 36% → Publish (first value)
    00:30 - Average: 40% → Skip (< 5m)
    01:00 - Average: 35% → Skip (< 5m)
    05:00 - Average: 38% → Publish (5m passed)
    ```
