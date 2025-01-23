package plugin

import (
	"context"
	"fmt"
	"sync"
	"threshAI/internal/telemetry"
	"time"
)

// defaultRegistry implements the PluginRegistry interface
type defaultRegistry struct {
	plugins  map[string]Plugin
	metrics  *telemetry.PipelineMetrics
	mu       sync.RWMutex
	statuses map[string]*HealthStatus
}

// NewRegistry creates a new plugin registry
func NewRegistry() PluginRegistry {
	return &defaultRegistry{
		plugins:  make(map[string]Plugin),
		metrics:  telemetry.GetMetrics(),
		statuses: make(map[string]*HealthStatus),
	}
}

func (r *defaultRegistry) Register(plugin Plugin) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	id := plugin.ID()
	if _, exists := r.plugins[id]; exists {
		return fmt.Errorf("plugin with ID %s already registered", id)
	}

	r.plugins[id] = plugin
	r.statuses[id] = &HealthStatus{
		Healthy: true,
		Status:  "Registered",
	}

	// Initialize metrics for the plugin
	r.metrics.SetPluginStatus(id, "registered", 1.0)
	return nil
}

func (r *defaultRegistry) Unregister(pluginID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	plugin, exists := r.plugins[pluginID]
	if !exists {
		return fmt.Errorf("plugin with ID %s not found", pluginID)
	}

	// Stop the plugin if it's running
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := plugin.Stop(ctx); err != nil {
		r.metrics.RecordPluginError(pluginID, "stop_error")
		return fmt.Errorf("failed to stop plugin %s: %v", pluginID, err)
	}

	delete(r.plugins, pluginID)
	delete(r.statuses, pluginID)
	r.metrics.SetPluginStatus(pluginID, "unregistered", 0.0)

	return nil
}

func (r *defaultRegistry) Get(pluginID string) (Plugin, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	plugin, exists := r.plugins[pluginID]
	if !exists {
		return nil, fmt.Errorf("plugin with ID %s not found", pluginID)
	}
	return plugin, nil
}

func (r *defaultRegistry) List() []Plugin {
	r.mu.RLock()
	defer r.mu.RUnlock()

	plugins := make([]Plugin, 0, len(r.plugins))
	for _, p := range r.plugins {
		plugins = append(plugins, p)
	}
	return plugins
}

func (r *defaultRegistry) StartAll(ctx context.Context) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	errChan := make(chan error, len(r.plugins))
	var wg sync.WaitGroup

	for id, p := range r.plugins {
		wg.Add(1)
		go func(pluginID string, plugin Plugin) {
			defer wg.Done()

			if err := plugin.Start(ctx); err != nil {
				r.metrics.RecordPluginError(pluginID, "start_error")
				errChan <- fmt.Errorf("failed to start plugin %s: %v", pluginID, err)
				return
			}

			r.metrics.SetPluginStatus(pluginID, "running", 1.0)
			r.statuses[pluginID] = &HealthStatus{
				Healthy: true,
				Status:  "Running",
			}
		}(id, p)
	}

	// Wait for all plugins to start
	wg.Wait()
	close(errChan)

	// Collect any errors
	var errors []error
	for err := range errChan {
		errors = append(errors, err)
	}

	if len(errors) > 0 {
		return fmt.Errorf("failed to start some plugins: %v", errors)
	}

	return nil
}

func (r *defaultRegistry) StopAll(ctx context.Context) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	errChan := make(chan error, len(r.plugins))
	var wg sync.WaitGroup

	for id, p := range r.plugins {
		wg.Add(1)
		go func(pluginID string, plugin Plugin) {
			defer wg.Done()

			if err := plugin.Stop(ctx); err != nil {
				r.metrics.RecordPluginError(pluginID, "stop_error")
				errChan <- fmt.Errorf("failed to stop plugin %s: %v", pluginID, err)
				return
			}

			r.metrics.SetPluginStatus(pluginID, "stopped", 0.0)
			r.statuses[pluginID] = &HealthStatus{
				Healthy: true,
				Status:  "Stopped",
			}
		}(id, p)
	}

	// Wait for all plugins to stop
	wg.Wait()
	close(errChan)

	// Collect any errors
	var errors []error
	for err := range errChan {
		errors = append(errors, err)
	}

	if len(errors) > 0 {
		return fmt.Errorf("failed to stop some plugins: %v", errors)
	}

	return nil
}

// startHealthCheck begins periodic health monitoring of plugins
func (r *defaultRegistry) startHealthCheck(ctx context.Context, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			r.checkPluginHealth()
		}
	}
}

func (r *defaultRegistry) checkPluginHealth() {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for id, plugin := range r.plugins {
		status := plugin.Health()
		r.statuses[id] = status

		if status.Healthy {
			r.metrics.SetPluginStatus(id, "healthy", 1.0)
		} else {
			r.metrics.SetPluginStatus(id, "unhealthy", 0.0)
			r.metrics.RecordPluginError(id, "health_check_failed")
		}

		// Record detailed metrics
		for key, value := range status.Details {
			r.metrics.SetPluginMetric(id, key, value)
		}
	}
}
