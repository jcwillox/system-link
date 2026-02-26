---
icon: material/code-json
---

# Template Variables

System Link use Go templates for templating client-side fields, e.g. `command`. Other fields such as `value_template` is executed server-side by Home Assistant and uses their [Jinja2 templating](https://www.home-assistant.io/docs/configuration/templating) engine.

| Variable        | Example                                                  | Description                                            |
|-----------------|----------------------------------------------------------|--------------------------------------------------------|
| `.ExePath`      | `~/.local/share/system-link/system-link`                 | Full path to the System Link executable                |
| `.ExeDirectory` | `~/.local/share/system-link`                             | Directory containing the executable                    |
| `.ExeName`      | `system-link`                                            | Executable filename                                    |
| `.Home`         | `~`                                                      | Current user's home directory                          |
| `.ShareDir`     | `~/.local/share/system-link`                             | Application data (share) directory                     |
| `.Portable`     | `false`                                                  | Whether System Link is running in portable mode (bool) |
| `.ConfigPath`   | `~/.config/system-link/config.yaml`                      | Full path to the active config file                    |
| `.ConfigDir`    | `~/.config/system-link`                                  | Directory containing the active config                 |
| `.LogsPath`     | `~/.local/share/system-link/system-link.8k2w9a1z4m.log`  | Path to the active log file                            |
| `.LockPath`     | `~/.local/share/system-link/system-link.8k2w9a1z4m.lock` | Lock file path used to detect running instance         |

## Example

```yaml
buttons:
  - custom:
      command: "{{ .ExeDirectory }}/scripts/backup.sh"
```
