# Disk Sensors

Monitor disk space usage.

## Configuration

```yaml
sensors:
  - disk: {}      # Disk usage percentage
  - disk_used: {} # Disk used in GB
  - disk_free: {} # Disk free in GB
```

## Options

```yaml
sensors:
  - disk:
      mountpoints:
        - /boot # Mount point to monitor (Linux/macOS)
        - /mnt/cache
        - /mnt/user
        - "C:" # Drive to monitor (Windows)
```

All [common sensor options](index.md#common-options) are also supported.
