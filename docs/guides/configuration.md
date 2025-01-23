# Configuration Guide

## Overview
This guide provides detailed instructions on how to configure threshAI to suit your specific needs. Configuration options include setting up environment variables, adjusting system settings, and customizing plugins.

## Environment Variables

### API Keys and Secrets
Set your API keys and secrets in the `.env` file:
```bash
# API Keys and Secrets
DEEPSEEK_API_KEY=your_key_here
```

### Optimization Level
Adjust the optimization level for prompt processing:
```bash
THRESHAI_OPTIMIZATION_LEVEL=high
```

### Plugins
Specify the plugins to use:
```bash
THRESHAI_PLUGINS=plugin1,plugin2
```

## System Settings

### Configuration File
You can also configure threshAI using a configuration file located at `~/.threshai/config.yaml`.

**Example Configuration File**:
```yaml
apiKey: YOUR_API_KEY
optimizationLevel: high
plugins:
  - plugin1
  - plugin2
```

## Customizing Plugins

### Adding Custom Plugins
To add custom plugins, follow these steps:

1. **Create a Plugin Directory**:
   ```bash
   mkdir -p ~/.threshai/plugins
   ```

2. **Write Your Plugin**:
   Create a new file in the `~/.threshai/plugins` directory. For example, `custom_plugin.go`:
   ```go
   package main

   import (
       "fmt"
   )

   func main() {
       fmt.Println("Custom plugin loaded")
   }
   ```

3. **Register the Plugin**:
   Update the `THRESHAI_PLUGINS` environment variable to include your custom plugin:
   ```bash
   THRESHAI_PLUGINS=plugin1,plugin2,custom_plugin
   ```

## Advanced Configuration

### Monitoring Settings
Configure monitoring settings in the `~/.threshai/monitoring.yaml` file:
```yaml
monitoring:
  interval: 60
  metrics:
    - cpu
    - memory
    - disk
```

### Security Settings
Configure security settings in the `~/.threshai/security.yaml` file:
```yaml
security:
  audit: true
  vulnerabilityScan: true
```

## Conclusion
Follow these steps to configure threshAI according to your needs. Ensure you set the appropriate environment variables and customize plugins as required.