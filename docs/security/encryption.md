---
icon: material/lock
---

# Encryption & TLS

System Link fully supports encrypted MQTT connections using TLS/SSL, ensuring all communication between System Link and your MQTT broker is secure and cannot be intercepted.

## Why Use TLS?

- 🔐 **Encryption** – All data transmitted over the network is encrypted
- ✅ **Authentication** – Verify that you're connecting to the correct MQTT broker
- 🚫 **MITM Protection** – Prevents man-in-the-middle attacks

## Basic TLS Configuration

The simplest way to enable TLS is to set the `tls` flag to `true`:

```yaml
mqtt:
  host: mqtt.home.local
  port: 8883  # Standard TLS port for MQTT
  tls: true
  username: system-link
  password: your_password
```

System Link will automatically verify the broker's certificate against the system's trusted CA certificates.

## Setting Up TLS with Mosquitto

Check the instructions for your MQTT broker to set up tls, and ideally also get a valid ssl certificate from a trusted CA like Let's Encrypt.
