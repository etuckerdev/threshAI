# Basic Usage

## Overview
This guide provides examples of basic usage scenarios for threshAI, covering common tasks such as running prompts, monitoring system health, and managing plugins.

## Running a Prompt

### Using the CLI
To run a prompt using the CLI, use the following command:
```bash
threshAI prompt run --input "Your prompt input here" --config '{"optimizationLevel": "high", "plugins": ["plugin1", "plugin2"]}'
```

### Using the Web Interface
1. Open the threshAI web interface.
2. Navigate to the "Run Prompt" section.
3. Enter your prompt input and select the desired configuration options.
4. Click "Run" to process the prompt.

## Monitoring System Health

### Using the CLI
To monitor system health using the CLI, use the following command:
```bash
threshAI system monitor
```

### Using the Web Interface
1. Open the threshAI web interface.
2. Navigate to the "System Health" section.
3. View real-time monitoring metrics and system status.

## Managing Plugins

### Adding a Plugin
To add a plugin, follow these steps:

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

### Removing a Plugin
To remove a plugin, update the `THRESHAI_PLUGINS` environment variable to exclude the plugin you want to remove:
```bash
THRESHAI_PLUGINS=plugin1
```

## Example Scenarios

### Scenario 1: Simple Prompt Processing
```bash
threshAI prompt run --input "Translate 'Hello, world!' to French"
```

### Scenario 2: Prompt Processing with Plugins
```bash
threshAI prompt run --input "Analyze sentiment of 'I love threshAI!'" --config '{"plugins": ["sentiment_analysis"]}'
```

### Scenario 3: Monitoring System Health
```bash
threshAI system monitor
```

## Conclusion
These examples demonstrate basic usage scenarios for threshAI. Use the CLI and web interface to run prompts, monitor system health, and manage plugins effectively.