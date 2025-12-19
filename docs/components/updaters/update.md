# Update

The updater component enables remote updating of System Link directly from Home Assistant. It will create an update entity that will be shown in Home Assistant's updates section.

## Configuration

Enable self-update capability:

```yaml
updaters:
  - update: {}
```

You can omit this if you don't want to allow updates to be triggered via Home Assistant/MQTT, or if you are using the docker container.

## Options

All [common updater options](index.md#common-options) are supported.
