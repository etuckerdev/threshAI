package plugin

import (
	"context"
	"encoding/json"
)

// Plugin represents a loadable extension that can modify or enhance system behavior
type Plugin interface {
	// ID returns the unique identifier for this plugin
	ID() string

	// Initialize sets up the plugin with its configuration
	Initialize(ctx context.Context, config json.RawMessage) error

	// Start begins plugin execution - should be non-blocking
	Start(ctx context.Context) error

	// Stop gracefully shuts down the plugin
	Stop(ctx context.Context) error

	// Health returns the current health status of the plugin
	Health() *HealthStatus
}

// HealthStatus represents the health of a plugin
type HealthStatus struct {
	Healthy bool              `json:"healthy"`
	Status  string            `json:"status"`
	Details map[string]string `json:"details,omitempty"`
}

// PluginMetadata contains information about a plugin
type PluginMetadata struct {
	ID          string   `json:"id"`
	Version     string   `json:"version"`
	Author      string   `json:"author"`
	Description string   `json:"description"`
	Permissions []string `json:"permissions,omitempty"`
}

// PluginRegistry manages plugin registration and lifecycle
type PluginRegistry interface {
	// Register adds a new plugin to the registry
	Register(plugin Plugin) error

	// Unregister removes a plugin from the registry
	Unregister(pluginID string) error

	// Get returns a plugin by ID
	Get(pluginID string) (Plugin, error)

	// List returns all registered plugins
	List() []Plugin

	// StartAll starts all registered plugins
	StartAll(ctx context.Context) error

	// StopAll stops all registered plugins
	StopAll(ctx context.Context) error
}

// PluginManager handles plugin operations and lifecycle
type PluginManager struct {
	registry PluginRegistry
	sandbox  *PluginSandbox
}

// PluginSandbox provides isolation for plugin execution
type PluginSandbox struct {
	// ResourceLimits defines memory and CPU constraints
	ResourceLimits struct {
		MaxMemoryMB   int `json:"maxMemoryMB"`
		MaxCPUPercent int `json:"maxCPUPercent"`
	}

	// Permissions defines allowed operations
	Permissions struct {
		AllowNetwork bool `json:"allowNetwork"`
		AllowFileIO  bool `json:"allowFileIO"`
		AllowExec    bool `json:"allowExec"`
	}
}
