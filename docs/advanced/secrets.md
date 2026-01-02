---
icon: material/key
---
# Secrets

The `!secret` tag allows you to store sensitive information in a separate `secrets.yaml` file instead of your main configuration file. This is useful for keeping passwords, tokens, and other credentials out of version control.

## Usage

Create a `secrets.yaml` file in the same directory as your `config.yaml`:

```yaml
mqtt_server: mqtt.example.com
mqtt_password: my_secret_password
mqtt_username: my_username
```

Then reference these secrets in your `config.yaml` using the `!secret` tag:

```yaml
mqtt:
  host: !secret mqtt_server
  port: 1883
  username: !secret mqtt_username
  password: !secret mqtt_password
```

## Best Practices

- Add `secrets.yaml` to your `.gitignore` to prevent accidentally committing secrets
- Keep all sensitive information in `secrets.yaml`
- Use descriptive names for your secrets to make configuration easier to read
