# Advanced Scenarios

## Overview
This guide provides examples of advanced usage scenarios for threshAI, covering more complex tasks such as integrating with external systems, customizing prompt processing pipelines, and leveraging advanced features.

## Integrating with External Systems

### Scenario 1: Integrating with a Weather API
To integrate threshAI with a weather API, follow these steps:

1. **Set Up API Key**:
   ```bash
   export WEATHER_API_KEY=your_weather_api_key
   ```

2. **Create a Custom Plugin**:
   Create a new file in the `~/.threshai/plugins` directory. For example, `weather_plugin.go`:
   ```go
   package main

   import (
       "fmt"
       "net/http"
       "io/ioutil"
       "os"
   )

   func main() {
       apiKey := os.Getenv("WEATHER_API_KEY")
       if apiKey == "" {
           fmt.Println("WEATHER_API_KEY environment variable is not set")
           return
       }

       resp, err := http.Get(fmt.Sprintf("https://api.weather.com/data/2.5/weather?q=London&appid=%s", apiKey))
       if err != nil {
           fmt.Println("Error fetching weather data:", err)
           return
       }
       defer resp.Body.Close()

       body, err := ioutil.ReadAll(resp.Body)
       if err != nil {
           fmt.Println("Error reading response body:", err)
           return
       }

       fmt.Println("Weather data:", string(body))
   }
   ```

3. **Register the Plugin**:
   Update the `THRESHAI_PLUGINS` environment variable to include your custom plugin:
   ```bash
   THRESHAI_PLUGINS=weather_plugin
   ```

4. **Run a Prompt with the Plugin**:
   ```bash
   threshAI prompt run --input "Get the current weather in London" --config '{"plugins": ["weather_plugin"]}'
   ```

## Customizing Prompt Processing Pipelines

### Scenario 2: Customizing the Prompt Processing Pipeline
To customize the prompt processing pipeline, follow these steps:

1. **Create a Custom Pipeline**:
   Create a new file in the `~/.threshai/pipelines` directory. For example, `custom_pipeline.go`:
   ```go
   package main

   import (
       "fmt"
   )

   func main() {
       fmt.Println("Custom pipeline loaded")
       // Add custom processing logic here
   }
   ```

2. **Register the Pipeline**:
   Update the `THRESHAI_PIPELINE` environment variable to use your custom pipeline:
   ```bash
   THRESHAI_PIPELINE=custom_pipeline
   ```

3. **Run a Prompt with the Custom Pipeline**:
   ```bash
   threshAI prompt run --input "Your prompt input here" --config '{"pipeline": "custom_pipeline"}'
   ```

## Leveraging Advanced Features

### Scenario 3: Real-Time Monitoring with Telemetry
To enable real-time monitoring with telemetry, follow these steps:

1. **Enable Telemetry**:
   Update the `~/.threshai/config.yaml` file to enable telemetry:
   ```yaml
   telemetry:
     enabled: true
     interval: 60
   ```

2. **Run a Prompt with Telemetry**:
   ```bash
   threshAI prompt run --input "Your prompt input here" --config '{"telemetry": true}'
   ```

3. **View Telemetry Data**:
   Use the threshAI web interface to view real-time telemetry data.

## Example Scenarios

### Scenario 1: Integrating with a Weather API
```bash
threshAI prompt run --input "Get the current weather in London" --config '{"plugins": ["weather_plugin"]}'
```

### Scenario 2: Customizing the Prompt Processing Pipeline
```bash
threshAI prompt run --input "Your prompt input here" --config '{"pipeline": "custom_pipeline"}'
```

### Scenario 3: Real-Time Monitoring with Telemetry
```bash
threshAI prompt run --input "Your prompt input here" --config '{"telemetry": true}'
```

## Conclusion
These advanced scenarios demonstrate how to integrate threshAI with external systems, customize prompt processing pipelines, and leverage advanced features. Use these examples as a starting point for more complex use cases.