# ZFS Pool Sensors

Monitor ZFS pool usage (FreeBSD/Linux with ZFS).

## Configuration

```yaml
sensors:
  - zpool: {}      # Pool usage percentage
  - zpool_used: {} # Pool used in GB
  - zpool_free: {} # Pool free in GB
```

## Options

```yaml
sensors:
  - zpool:
      pools:
        - cache # ZFS pool name
```

All [common sensor options](index.md#common-options) are also supported.
