package plugin

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"threshAI/internal/telemetry"
)

// NewPluginManager creates a new plugin manager
func NewPluginManager() *PluginManager {
	return &PluginManager{
		registry: NewRegistry(),
		sandbox:  NewPluginSandbox(nil),
	}
}

// LoadPlugin loads and initializes a plugin with the given configuration
func (pm *PluginManager) LoadPlugin(pluginID string, plugin Plugin, config json.RawMessage) error {
	// Create sandboxed context for plugin initialization
	ctx := NewSandboxedContext(context.Background(), pm.sandbox)

	// Initialize plugin
	if err := plugin.Initialize(ctx, config); err != nil {
		return fmt.Errorf("failed to initialize plugin %s: %v", pluginID, err)
	}

	// Register plugin
	if err := pm.registry.Register(plugin); err != nil {
		return fmt.Errorf("failed to register plugin %s: %v", pluginID, err)
	}

	return nil
}

// StartPlugin starts a specific plugin
func (pm *PluginManager) StartPlugin(pluginID string) error {
	plugin, err := pm.registry.Get(pluginID)
	if err != nil {
		return err
	}

	ctx := NewSandboxedContext(context.Background(), pm.sandbox)
	return plugin.Start(ctx)
}

// StopPlugin stops a specific plugin
func (pm *PluginManager) StopPlugin(pluginID string) error {
	plugin, err := pm.registry.Get(pluginID)
	if err != nil {
		return err
	}

	ctx := NewSandboxedContext(context.Background(), pm.sandbox)
	return plugin.Stop(ctx)
}

// UnloadPlugin stops and unregisters a plugin
func (pm *PluginManager) UnloadPlugin(pluginID string) error {
	if err := pm.StopPlugin(pluginID); err != nil {
		return err
	}
	return pm.registry.Unregister(pluginID)
}

// GetPluginStatus returns the current health status of a plugin
func (pm *PluginManager) GetPluginStatus(pluginID string) (*HealthStatus, error) {
	plugin, err := pm.registry.Get(pluginID)
	if err != nil {
		return nil, err
	}
	return plugin.Health(), nil
}

// UpdatePluginConfig updates plugin configuration
func (pm *PluginManager) UpdatePluginConfig(pluginID string, config json.RawMessage) error {
	plugin, err := pm.registry.Get(pluginID)
	if err != nil {
		return err
	}

	ctx := NewSandboxedContext(context.Background(), pm.sandbox)
	return plugin.Initialize(ctx, config)
}

// ListPlugins returns information about all registered plugins
func (pm *PluginManager) ListPlugins() []PluginInfo {
	plugins := pm.registry.List()
	info := make([]PluginInfo, len(plugins))

	for i, p := range plugins {
		status := p.Health()
		info[i] = PluginInfo{
			ID:      p.ID(),
			Health:  status.Healthy,
			Status:  status.Status,
			Details: status.Details,
		}
	}

	return info
}

// PluginInfo contains plugin status information
type PluginInfo struct {
	ID      string            `json:"id"`
	Health  bool              `json:"health"`
	Status  string            `json:"status"`
	Details map[string]string `json:"details,omitempty"`
}

// BatchOperation represents a bulk plugin operation
type BatchOperation struct {
	PluginIDs []string        `json:"pluginIds"`
	Config    json.RawMessage `json:"config,omitempty"`
}

// BatchStart starts multiple plugins
func (pm *PluginManager) BatchStart(batch BatchOperation) []error {
	var (
		wg     sync.WaitGroup
		mu     sync.Mutex
		errors []error
	)

	for _, id := range batch.PluginIDs {
		wg.Add(1)
		go func(pluginID string) {
			defer wg.Done()
			if err := pm.StartPlugin(pluginID); err != nil {
				mu.Lock()
				errors = append(errors, fmt.Errorf("plugin %s: %v", pluginID, err))
				mu.Unlock()
			}
		}(id)
	}

	wg.Wait()
	return errors
}

// BatchStop stops multiple plugins
func (pm *PluginManager) BatchStop(batch BatchOperation) []error {
	var (
		wg     sync.WaitGroup
		mu     sync.Mutex
		errors []error
	)

	for _, id := range batch.PluginIDs {
		wg.Add(1)
		go func(pluginID string) {
			defer wg.Done()
			if err := pm.StopPlugin(pluginID); err != nil {
				mu.Lock()
				errors = append(errors, fmt.Errorf("plugin %s: %v", pluginID, err))
				mu.Unlock()
			}
		}(id)
	}

	wg.Wait()
	return errors
}

// UpdateSandboxLimits updates resource limits for the plugin sandbox
func (pm *PluginManager) UpdateSandboxLimits(limits *PluginSandbox) {
	// Create new sandbox with updated limits
	newSandbox := NewPluginSandbox(limits)

	// Update sandbox reference
	pm.sandbox = newSandbox

	// Notify metrics of limit changes
	metrics := telemetry.GetMetrics()
	metrics.SetPluginMetric("sandbox", "max_memory_mb", float64(limits.ResourceLimits.MaxMemoryMB))
	metrics.SetPluginMetric("sandbox", "max_cpu_percent", float64(limits.ResourceLimits.MaxCPUPercent))
}
