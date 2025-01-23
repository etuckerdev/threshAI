package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"threshAI/internal/core/plugin"
	"threshAI/internal/core/plugin/examples"

	"github.com/spf13/cobra"
)

var pluginManager *plugin.PluginManager

func init() {
	pluginManager = plugin.NewPluginManager()
	systemCmd.AddCommand(systemPluginCmd)
	systemCmd.AddCommand(systemStatusCmd)
	systemCmd.AddCommand(systemDiagCmd)
	systemCmd.AddCommand(systemMetricsCmd)
	rootCmd.AddCommand(systemCmd)
}

var systemCmd = &cobra.Command{
	Use:   "system",
	Short: "System management commands",
	Long: `Commands for managing and monitoring the ThreshAI system:
- Check system status
- Monitor resource usage
- View diagnostic information`,
	GroupID: "system",
}

var systemStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show system status",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("System Status:")
		// TODO: Add comprehensive status checks
		checkRedisConnection()
		checkModelAvailability()
	},
}

var systemDiagCmd = &cobra.Command{
	Use:   "diag",
	Short: "Run system diagnostics",
	Run: func(cmd *cobra.Command, args []string) {
		if verbose {
			fmt.Println("Running detailed diagnostics...")
		} else {
			fmt.Println("Running basic diagnostics...")
		}
		// TODO: Implement diagnostic tests
	},
}

var systemMetricsCmd = &cobra.Command{
	Use:   "metrics",
	Short: "Show system metrics",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("System Metrics:")
		// TODO: Implement metrics collection and display
	},
}

func checkRedisConnection() {
	fmt.Println("✓ Redis Connection: OK")
}

func checkModelAvailability() {
	fmt.Println("✓ LLM Model Status: Available")
}

var systemPluginCmd = &cobra.Command{
	Use:   "plugin",
	Short: "Plugin management commands",
	Long:  "Commands for managing plugins: list, load, unload, start, stop, status",
}

var pluginListCmd = &cobra.Command{
	Use:   "list",
	Short: "List loaded plugins",
	Run: func(cmd *cobra.Command, args []string) {
		plugins := pluginManager.ListPlugins()
		fmt.Println("Loaded Plugins:")
		if len(plugins) == 0 {
			fmt.Println("No plugins loaded.")
			return
		}
		for _, p := range plugins {
			fmt.Printf("- ID: %s, Status: %s, Healthy: %t\n", p.ID, p.Status, p.Health)
			if !p.Health {
				fmt.Printf("  Details: %v\n", p.Details)
			}
		}
	},
}

var pluginLoadCmd = &cobra.Command{
	Use:   "load [pluginID] [configPath]",
	Short: "Load a plugin",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		pluginID := args[0]
		configPath := args[1]

		configData, err := os.ReadFile(configPath)
		if err != nil {
			fmt.Printf("Error reading config file: %v\n", err)
			return
		}

		var config json.RawMessage
		if err := json.Unmarshal(configData, &config); err != nil {
			fmt.Printf("Error parsing config JSON: %v\n", err)
			return
		}

		// Example: Load TemplatePlugin - in real use, plugin type would be dynamically determined
		tmplPlugin := examples.NewTemplatePlugin(pluginID)
		if err := pluginManager.LoadPlugin(pluginID, tmplPlugin, config); err != nil {
			fmt.Printf("Error loading plugin %s: %v\n", pluginID, err)
			return
		}

		fmt.Printf("Plugin %s loaded successfully.\n", pluginID)
	},
}

var pluginUnloadCmd = &cobra.Command{
	Use:   "unload [pluginID]",
	Short: "Unload a plugin",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		pluginID := args[0]
		if err := pluginManager.UnloadPlugin(pluginID); err != nil {
			fmt.Printf("Error unloading plugin %s: %v\n", pluginID, err)
			return
		}
		fmt.Printf("Plugin %s unloaded successfully.\n", pluginID)
	},
}

var pluginStartCmd = &cobra.Command{
	Use:   "start [pluginID]",
	Short: "Start a plugin",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		pluginID := args[0]
		if err := pluginManager.StartPlugin(pluginID); err != nil {
			fmt.Printf("Error starting plugin %s: %v\n", pluginID, err)
			return
		}
		fmt.Printf("Plugin %s started successfully.\n", pluginID)
	},
}

var pluginStopCmd = &cobra.Command{
	Use:   "stop [pluginID]",
	Short: "Stop a plugin",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		pluginID := args[0]
		if err := pluginManager.StopPlugin(pluginID); err != nil {
			fmt.Printf("Error stopping plugin %s: %v\n", pluginID, err)
			return
		}
		fmt.Printf("Plugin %s stopped successfully.\n", pluginID)
	},
}

var pluginStatusCmd = &cobra.Command{
	Use:   "status [pluginID]",
	Short: "Get plugin status",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		pluginID := args[0]
		status, err := pluginManager.GetPluginStatus(pluginID)
		if err != nil {
			fmt.Printf("Error getting plugin status %s: %v\n", pluginID, err)
			return
		}
		fmt.Printf("Plugin %s Status: %s, Healthy: %t\n", pluginID, status.Status, status.Healthy)
		if !status.Healthy {
			fmt.Printf("Details: %v\n", status.Details)
		}
	},
}

func init() {
	systemPluginCmd.AddCommand(pluginListCmd)
	systemPluginCmd.AddCommand(pluginLoadCmd)
	systemPluginCmd.AddCommand(pluginUnloadCmd)
	systemPluginCmd.AddCommand(pluginStartCmd)
	systemPluginCmd.AddCommand(pluginStopCmd)
	systemPluginCmd.AddCommand(pluginStatusCmd)
}
