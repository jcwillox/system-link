# Custom Buttons

Create buttons that run your own commands.

## Basic Custom Button

```yaml
buttons:
  - custom:
      name: "Restart Service"
      icon: mdi:restart
      unique_id: restart_nginx
      command: "systemctl restart nginx"
```

## Custom Button with Script

```yaml
buttons:
  - custom:
      name: "Run Backup"
      icon: mdi:backup-restore
      unique_id: run_backup
      command: "/scripts/backup.sh"
      shell: "bash"
```

## Custom Button Configuration

```yaml
buttons:
  - custom:
      name: "Custom Button"
      unique_id: "custom_button"
      icon: "mdi:gesture-tap-button"
      command: "echo 'Hello World'"
      shell: "bash"              # Shell to use
      detached: false            # Run in background
      show_output: true          # Log command output
      show_errors: true          # Log command errors
      env:                       # Environment variables
        MY_VAR: "value"
```

## Custom Button Options

| Option | Type | Description | Default |
|--------|------|-------------|---------|
| `name` | string | Display name | Required |
| `unique_id` | string | Unique identifier | Required |
| `icon` | string | MDI icon | `mdi:gesture-tap-button` |
| `command` | string | Command to execute | Required |
| `shell` | string | Shell (bash, powershell, cmd) | System default |
| `detached` | boolean | Run without waiting for completion | `false` |
| `show_output` | boolean | Log standard output | `false` |
| `show_errors` | boolean | Log error output | `false` |
| `env` | map | Environment variables | - |

## Template Variables

Custom button commands can be templated using Go templates:

```yaml
buttons:
  - custom:
      name: "Show Info"
      unique_id: show_info
      command: |
        echo "Device: {{ .DeviceName }}"
        echo "Host ID: {{ .HostID }}"
        echo "Directory: {{ .ExeDirectory }}"
```

See [Template Variables Reference](../../advanced/templating.md) for all available variables.

## Custom Button Examples

### Restart Application

=== "Linux"
    ```yaml
    buttons:
      - custom:
          name: "Restart Docker"
          icon: mdi:docker
          unique_id: restart_docker
          command: "systemctl restart docker"
    ```

=== "Windows"
    ```yaml
    buttons:
      - custom:
          name: "Restart Service"
          icon: mdi:application-cog
          unique_id: restart_service
          shell: "powershell"
          command: "Restart-Service -Name 'MyService'"
    ```

### Open Application

```yaml
buttons:
  - custom:
      name: "Open Notepad"
      icon: mdi:note-text
      unique_id: open_notepad
      detached: true  # Don't wait for app to close
      command: "C:\\Windows\\System32\\notepad.exe"
```

### Download File

```yaml
buttons:
  - custom:
      name: "Download Config"
      icon: mdi:cloud-download-outline
      unique_id: download_config
      shell: "powershell"
      command: |
        $output = "{{ .ExeDirectory }}\config.yaml"
        Invoke-WebRequest -Uri "https://example.com/config.yaml" -OutFile $output
```

### Clean Temp Files

```yaml
buttons:
  - custom:
      name: "Clean Temp"
      icon: mdi:delete-sweep
      unique_id: clean_temp
      shell: "bash"
      command: |
        rm -rf /tmp/*
        rm -rf ~/.cache/*
```

### Docker Operations

```yaml
buttons:
  - custom:
      name: "Restart Container"
      icon: mdi:docker
      unique_id: restart_nginx
      command: "docker restart nginx"

  - custom:
      name: "Update Container"
      icon: mdi:update
      unique_id: update_nginx
      shell: "bash"
      command: |
        docker pull nginx:latest
        docker stop nginx
        docker rm nginx
        docker run -d --name nginx nginx:latest
```

### System Maintenance

```yaml
buttons:
  - custom:
      name: "Clear Package Cache"
      icon: mdi:package-variant
      unique_id: clear_cache
      command: "sudo apt-get clean"

  - custom:
      name: "Update System"
      icon: mdi:update
      unique_id: update_system
      shell: "bash"
      command: |
        sudo apt-get update
        sudo apt-get upgrade -y
```
