# CLI Reference

## Overview
The threshAI CLI provides a command-line interface for interacting with the platform, allowing users to manage prompts, monitor system health, and utilize various features.

## Installation
To install the threshAI CLI, use the following command:
```bash
go install github.com/threshai/threshai@latest
```

## Commands

### `threshAI prompt run`
Run a prompt with the specified input.

**Usage**:
```bash
threshAI prompt run --input "Your prompt input here" [--config CONFIG]
```

**Options**:
- `--input`: The input prompt to process.
- `--config`: Optional configuration for the prompt processing (e.g., optimization level, plugins).

**Example**:
```bash
threshAI prompt run --input "Hello, world!" --config '{"optimizationLevel": "high", "plugins": ["plugin1", "plugin2"]}'
```

### `threshAI prompt status`
Retrieve the status of a specific prompt.

**Usage**:
```bash
threshAI prompt status --prompt-id PROMPT_ID
```

**Options**:
- `--prompt-id`: The ID of the prompt.

**Example**:
```bash
threshAI prompt status --prompt-id 12345
```

### `threshAI prompt list`
List all prompts.

**Usage**:
```bash
threshAI prompt list
```

**Example**:
```bash
threshAI prompt list
```

### `threshAI prompt delete`
Delete a specific prompt.

**Usage**:
```bash
threshAI prompt delete --prompt-id PROMPT_ID
```

**Options**:
- `--prompt-id`: The ID of the prompt.

**Example**:
```bash
threshAI prompt delete --prompt-id 12345
```

### `threshAI system monitor`
Monitor the system health and performance.

**Usage**:
```bash
threshAI system monitor
```

**Example**:
```bash
threshAI system monitor
```

### `threshAI system info`
Retrieve system information.

**Usage**:
```bash
threshAI system info
```

**Example**:
```bash
threshAI system info
```

## Configuration
The CLI can be configured using a configuration file or environment variables. The configuration file should be located at `~/.threshai/config.yaml`.

**Example Configuration File**:
```yaml
apiKey: YOUR_API_KEY
optimizationLevel: high
plugins:
  - plugin1
  - plugin2
```

## Environment Variables
The following environment variables can be used to configure the CLI:

- `THRESHAI_API_KEY`: Your API key.
- `THRESHAI_OPTIMIZATION_LEVEL`: Optimization level for prompt processing.
- `THRESHAI_PLUGINS`: Comma-separated list of plugins to use.

## Conclusion
The threshAI CLI provides a powerful interface for managing prompts and monitoring system health. Ensure you follow the installation and configuration guidelines to ensure smooth CLI usage.