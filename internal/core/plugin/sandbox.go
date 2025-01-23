package plugin

import (
	"context"
	"fmt"
	"os"
	"runtime"
	"sync"
	"time"
)

// NewPluginSandbox creates a new sandbox with the specified limits
func NewPluginSandbox(limits *PluginSandbox) *PluginSandbox {
	if limits == nil {
		limits = &PluginSandbox{
			ResourceLimits: struct {
				MaxMemoryMB   int `json:"maxMemoryMB"`
				MaxCPUPercent int `json:"maxCPUPercent"`
			}{
				MaxMemoryMB:   512, // Default 512MB
				MaxCPUPercent: 50,  // Default 50% CPU
			},
			Permissions: struct {
				AllowNetwork bool `json:"allowNetwork"`
				AllowFileIO  bool `json:"allowFileIO"`
				AllowExec    bool `json:"allowExec"`
			}{
				AllowNetwork: false,
				AllowFileIO:  false,
				AllowExec:    false,
			},
		}
	}
	return limits
}

// ResourceMonitor tracks resource usage of plugins
type ResourceMonitor struct {
	mu      sync.RWMutex
	plugins map[string]*PluginProcess
	limits  *PluginSandbox
}

// PluginProcess holds process information
type PluginProcess struct {
	pid       int
	startTime time.Time
	memStats  *runtime.MemStats
}

// NewResourceMonitor creates a new resource monitor
func NewResourceMonitor(limits *PluginSandbox) *ResourceMonitor {
	return &ResourceMonitor{
		plugins: make(map[string]*PluginProcess),
		limits:  limits,
	}
}

// RegisterPlugin adds a plugin to be monitored
func (rm *ResourceMonitor) RegisterPlugin(pluginID string) error {
	rm.mu.Lock()
	defer rm.mu.Unlock()

	proc := &PluginProcess{
		pid:       os.Getpid(),
		startTime: time.Now(),
		memStats:  &runtime.MemStats{},
	}

	rm.plugins[pluginID] = proc
	return nil
}

// UnregisterPlugin removes a plugin from monitoring
func (rm *ResourceMonitor) UnregisterPlugin(pluginID string) {
	rm.mu.Lock()
	defer rm.mu.Unlock()
	delete(rm.plugins, pluginID)
}

// StartMonitoring begins resource monitoring
func (rm *ResourceMonitor) StartMonitoring(ctx context.Context) {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			rm.checkResourceUsage()
		}
	}
}

// checkResourceUsage verifies resource limits aren't exceeded
func (rm *ResourceMonitor) checkResourceUsage() {
	rm.mu.RLock()
	defer rm.mu.RUnlock()

	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	for pluginID, proc := range rm.plugins {
		proc.memStats = &memStats

		// Check memory usage
		memoryMB := float64(memStats.Alloc) / 1024 / 1024
		if memoryMB > float64(rm.limits.ResourceLimits.MaxMemoryMB) {
			// Log memory limit exceeded
			rm.logResourceViolation(pluginID, fmt.Sprintf("Memory limit exceeded: %.2fMB", memoryMB))
		}

		// Check CPU usage via goroutine count as a rough approximation
		numGoroutines := runtime.NumGoroutine()
		if numGoroutines > rm.limits.ResourceLimits.MaxCPUPercent {
			rm.logResourceViolation(pluginID, fmt.Sprintf("CPU limit exceeded: %d goroutines", numGoroutines))
		}
	}
}

func (rm *ResourceMonitor) logResourceViolation(pluginID string, message string) {
	// In a real implementation, this would integrate with the logging system
	fmt.Printf("Resource violation for plugin %s: %s\n", pluginID, message)
}

// PermissionChecker validates plugin operations
type PermissionChecker struct {
	limits *PluginSandbox
}

// NewPermissionChecker creates a new permission checker
func NewPermissionChecker(limits *PluginSandbox) *PermissionChecker {
	return &PermissionChecker{limits: limits}
}

// CheckNetworkAccess verifies if network access is allowed
func (pc *PermissionChecker) CheckNetworkAccess() error {
	if !pc.limits.Permissions.AllowNetwork {
		return fmt.Errorf("network access denied")
	}
	return nil
}

// CheckFileAccess verifies if file I/O is allowed
func (pc *PermissionChecker) CheckFileAccess() error {
	if !pc.limits.Permissions.AllowFileIO {
		return fmt.Errorf("file I/O access denied")
	}
	return nil
}

// CheckExecAccess verifies if execution is allowed
func (pc *PermissionChecker) CheckExecAccess() error {
	if !pc.limits.Permissions.AllowExec {
		return fmt.Errorf("execution access denied")
	}
	return nil
}

// SandboxedContext wraps a context with resource monitoring
type SandboxedContext struct {
	context.Context
	rm   *ResourceMonitor
	pc   *PermissionChecker
	done chan struct{}
}

// NewSandboxedContext creates a new sandboxed context
func NewSandboxedContext(ctx context.Context, limits *PluginSandbox) *SandboxedContext {
	rm := NewResourceMonitor(limits)
	pc := NewPermissionChecker(limits)

	sandboxed := &SandboxedContext{
		Context: ctx,
		rm:      rm,
		pc:      pc,
		done:    make(chan struct{}),
	}

	go func() {
		sandboxed.rm.StartMonitoring(ctx)
	}()

	return sandboxed
}

// Done returns the done channel
func (sc *SandboxedContext) Done() <-chan struct{} {
	return sc.done
}

// CheckPermission verifies if an operation is allowed
func (sc *SandboxedContext) CheckPermission(op string) error {
	switch op {
	case "network":
		return sc.pc.CheckNetworkAccess()
	case "file":
		return sc.pc.CheckFileAccess()
	case "exec":
		return sc.pc.CheckExecAccess()
	default:
		return fmt.Errorf("unknown operation: %s", op)
	}
}
