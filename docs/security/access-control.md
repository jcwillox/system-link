---
icon: material/shield-lock
---

# MQTT Access Control Lists (ACLs)

System Link supports MQTT ACLs, allowing you to restrict what each device can control.

This is great for ensuring one compromised device can't control or compromise others on your network, especially when running system-link on an untrusted device, or if you just don't want to put your root mqtt creds everywhere.

## Per-Device ACLs

Restrict each System Link instance to only its own topics:

```toml
# /etc/mosquitto/acl

# Device 1 - Gaming PC
user gaming-pc-mqtt
topic system-link/3d8085e5-1fca-44e5-8d5a-0dd95b1a74ee/#
topic system-link/+/3d8085e5-1fca-44e5-8d5a-0dd95b1a74ee/#
topic homeassistant/+/3d8085e5-1fca-44e5-8d5a-0dd95b1a74ee/#

# Device 2 - Media Server
user media-server-mqtt
topic system-link/7f3a2b9c-8e4d-4a1f-9b2c-6d5e8f7a9b1c/#
topic system-link/+/7f3a2b9c-8e4d-4a1f-9b2c-6d5e8f7a9b1c/#
topic homeassistant/+/7f3a2b9c-8e4d-4a1f-9b2c-6d5e8f7a9b1c/#

# Home Assistant - Can control all devices
user homeassistant
topic #
```

## Finding Your Device UUID

The device UUID (`host_id`) is used in ACL rules. Find it by:

1. **Check System Link logs** on first run
2. **Home Assistant** device info page
3. **Explicitly set** in config.yaml:

```yaml
host_id: "3d8085e5-1fca-44e5-8d5a-0dd95b1a74ee"
```
