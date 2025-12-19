# Screen

Capture screenshots from your system.

## Basic Configuration

```yaml
images:
  - screen: {}
```

## Configuration Options

```yaml
images:
  - screen:
      availability: false  # Don't track availability (recommended)
      entities:
        timing: {}         # Enable entity to track how long image capture takes
        interval:          # Enable entity to control screenshot interval
          min: 0.1
          max: 30
          step: 0.1
```

!!! note

    Capture interval should be greater than the time it takes to capture the screenshot, and higher capture rates will use more CPU and bandwidth, although on modern CPUs this should be minimal.

## Sub-Entities

### Timing

Creates a sensor to monitor how long each screenshot capture takes.

```yaml
images:
  - screen:
      entities:
        timing: {}
```

### Interval

This creates a number input in Home Assistant where you can control and override the default capture interval.

```yaml
images:
  - screen:
      entities:
        interval:
          min: 1      # Minimum 1 second
          max: 60     # Maximum 60 seconds
          step: 1     # Increment by 1 second
```

## Options

All [common image options](index.md#common-options) are supported.
