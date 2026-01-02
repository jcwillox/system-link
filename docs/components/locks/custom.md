# Custom Locks

Locks provide stateful control with separate lock and unlock commands. They're perfect for managing applications or services that need to be started and stopped.

Locks are always custom-defined with separate lock and unlock commands.

## Basic Lock Configuration

```yaml
locks:
  - custom:
      name: "Media Player"
      unique_id: media_player
      lock:
        command: "vlc /path/to/video.mp4"
      unlock:
        command: "pkill vlc"
```

## Lock Configuration Options

```yaml
locks:
  - custom:
      name: "Custom Lock"
      unique_id: "custom_lock"
      icon: "mdi:lock"
      optimistic: false          # Update state optimistically
      lock:                      # Lock command (required)
        command: "start_app.sh"
        shell: "bash"
        detached: true
        show_output: false
        show_errors: false
        env:
          VAR: "value"
      unlock:                    # Unlock command (required)
        command: "stop_app.sh"
        shell: "bash"
        show_output: false
        show_errors: false
```

## Optimistic Mode

When `optimistic: true`, the lock immediately updates its state in Home Assistant when commanded, without waiting for confirmation. This is useful for commands that don't provide feedback.

```yaml
locks:
  - custom:
      name: "Quick Lock"
      unique_id: quick_lock
      optimistic: true  # Update state immediately
      lock:
        command: "start_service.sh"
      unlock:
        command: "stop_service.sh"
```

## Advanced Examples

### Application with Cleanup

Start an application and ensure proper cleanup on stop:

```yaml
locks:
  - custom:
      name: "Notepad Manager"
      unique_id: notepad_manager
      optimistic: true
      lock:
        detached: true
        shell: "powershell"
        command: |
          taskkill /f /im notepad.exe /erroraction SilentlyContinue
          Start-Process notepad.exe
      unlock:
        shell: "powershell"
        command: |
          taskkill /f /im notepad.exe
```

### Docker Compose Stack

Control entire Docker Compose stacks:

```yaml
locks:
  - custom:
      name: "Application Stack"
      unique_id: app_stack
      lock:
        shell: "bash"
        command: |
          cd /opt/myapp
          docker-compose up -d
      unlock:
        shell: "bash"
        command: |
          cd /opt/myapp
          docker-compose down
```

### Process with Monitoring

Start a process and monitor it:

```yaml
locks:
  - custom:
      name: "Monitored Service"
      unique_id: monitored_service
      lock:
        shell: "bash"
        command: |
          /usr/local/bin/myservice &
          echo $! > /var/run/myservice.pid
      unlock:
        shell: "bash"
        command: |
          if [ -f /var/run/myservice.pid ]; then
            kill $(cat /var/run/myservice.pid)
            rm /var/run/myservice.pid
          fi
```

### Screen Lock with Application

Lock screen and start a screensaver app:

```yaml
locks:
  - custom:
      name: "Screensaver Lock"
      unique_id: screensaver_lock
      optimistic: true
      lock:
        detached: true
        shell: "bash"
        command: |
          gnome-screensaver-command -l
          /usr/local/bin/custom-screensaver &
      unlock:
        command: "pkill custom-screensaver"
```

## Template Variables

Lock commands have access to template variables:

```yaml
locks:
  - custom:
      name: "Custom App"
      unique_id: custom_app
      lock:
        shell: "bash"
        command: |
          cd "{{ .ExeDirectory }}/apps"
          ./myapp.sh start
      unlock:
        command: "{{ .ExeDirectory }}/apps/myapp.sh stop"
```

See [Template Variables Reference](../../advanced/templating.md).
