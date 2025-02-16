# System-Link

[![GitHub Release](https://img.shields.io/github/v/release/jcwillox/system-link?style=flat-square)](https://github.com/jcwillox/system-link/releases/latest)
[![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/jcwillox/system-link?style=flat-square&label=go)](https://github.com/jcwillox/system-link/blob/main/go.mod)
[![License](https://img.shields.io/github/license/jcwillox/system-link?style=flat-square)](https://github.com/jcwillox/system-link/blob/main/LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/jcwillox/system-link?style=flat-square)](https://goreportcard.com/report/github.com/jcwillox/system-link)

## Installation

You can download the latest release from the [releases page](https://github.com/jcwillox/system-link/releases/latest). Or use one of the following methods:

### Script

```bash
sh -c "$(curl -fsSL jcwillox.com/l/syslink)"
```

Windows

```powershell
iwr -useb jcwillox.com/l/syslink-ps1 | iex
```

### Docker

You can pull the Docker image from GitHub Container Registry:

```sh
docker pull ghcr.io/jcwillox/system-link:latest
```

Full Config

```yaml
# override host name
# device_name: "My PC"

# unique id for the host must be globally unique, this is optional
# and only needed for running multiple instances on the same host,
# particularly useful for docker containers as they usually have
# the same machine id as their host.

# host_id: "my-pc-unique-id"
# host_id: "3c5fc61b-a47b-4437-b38c-5979a216e1d4"

mqtt:
  host: 127.0.0.1
  port: 8883
  tls: true
  username: admin
  password: password

# log_level: debug
# log_level_mqtt: debug
# log_timing: true

buttons:
  - lock: {}
  - shutdown: {}
  - force_shutdown: {}
  - sleep: {}
  - reload: {}
  - exit: {}
  - reset_topics: {}
  - custom:
      name: "Download config"
      icon: mdi:cloud-download-outline
      unique_id: download_config
      shell: "powershell"
      # language=powershell
      command: |
        $output = "{{ .ExeDirectory }}\config.yaml"
        Invoke-WebRequest -Uri "https://example.com/config.yaml" -OutFile $output

  - custom:
      name: "Download assets"
      icon: mdi:cloud-download-outline
      unique_id: download_assets
      shell: "powershell"
      # language=powershell
      command: |
        $zipPath = "{{ .ExeDirectory }}\assets.zip"
        $extractPath = "{{ .ExeDirectory }}\assets"

        Remove-Item -Path $zipPath -ErrorAction SilentlyContinue
        Remove-Item -Path $extractPath -Recurse -ErrorAction SilentlyContinue

        Invoke-WebRequest -Uri "https://example.com/assets.zip" -OutFile $zipPath
        Expand-Archive -Path $zipPath -DestinationPath $extractPath -Force

        Remove-Item -Path $zipPath -ErrorAction SilentlyContinue

images:
  - screen:
      availability: false

sensors:
  - cpu:
      update_interval: 10s
      filters:
        - or:
            - throttle: 5m
            - delta: 5

  - memory: {}
  - memory_used: {}
  - memory_free: {}

  - swap: {}
  - swap_used: {}
  - swap_free: {}

  - disk: {}
  - disk_used: {}
  - disk_free: {}

  - uptime: {}

  - battery: {}
  - battery_state: {}

binary_sensors:
  - status: {}

switches:
  - startup: {}
  - cron:
      name: "Open Notepad"
      icon: mdi:file-document
      unique_id: open_notepad
      schedule: "*/5 * * * * *"
      # detaches the spawned process from system link
      # and avoids waiting for it to finish
      detached: true
      command: "C:\\Windows\\System32\\notepad.exe"

locks:
  - custom:
      name: "Open and Kill Notepad"
      unique_id: open_and_kill_notepad
      optimistic: true
      lock:
        detached: true
        shell: "powershell"
        # language=powershell
        command: |
          taskkill /f /im notepad.exe
          C:\Windows\System32\notepad.exe
      unlock:
       command: taskkill /f /im notepad.exe

updaters:
  - update: {}
```
