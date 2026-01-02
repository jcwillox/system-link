# Reset Topics Button

Clear all MQTT topics published by System Link, and restarts System Link.

This is very useful for cleaning up entities which are no longer exposed, for example if you removed components from your configuration.

## Configuration

```yaml
buttons:
  - reset_topics: {}
```

## Options

All [common button options](index.md#common-options) are supported.
