# Images

Image components capture and stream images from your system to Home Assistant.

## Available Images

### Screen Capture

- **[Screen](screen.md)** - Capture screenshots from your system

## Configuration Options

### Common Options

All image components support these common options:

```yaml
images:
  - image_type:
      name: "Custom Name"
      unique_id: "custom_id"
      availability: false  # Don't track availability (recommended for images)
```

