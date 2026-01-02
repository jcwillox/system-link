---
icon: material/code-json
---

# Template Variables

System Link use Go templates for templating client-side fields, e.g. `command`. Other fields such as `value_template` is executed server-side by Home Assistant and uses their [Jinja2 templating](https://www.home-assistant.io/docs/configuration/templating) engine.

| Variable        | Example                                  | Description                             |
|-----------------|------------------------------------------|-----------------------------------------|
| `.ExePath`      | `~/.local/share/system-link/system-link` | Full path to the System Link executable |
| `.ExeDirectory` | `~/.local/share/system-link`             | Directory containing the executable     |
| `.ExeName`      | `system-link`                            | Executable filename                     |

## Example

```yaml
buttons:
  - custom:
      command: "{{ .ExeDirectory }}/scripts/backup.sh"
```
