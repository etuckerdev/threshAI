package examples

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"text/template"
	"threshAI/internal/core/plugin"
	"time"
)

// TemplatePlugin demonstrates a plugin that provides template processing capabilities
type TemplatePlugin struct {
	id      string
	config  *TemplateConfig
	mutex   sync.RWMutex
	tmpl    *template.Template
	running bool
	metrics map[string]string
}

// TemplateConfig holds plugin configuration
type TemplateConfig struct {
	TemplateName string            `json:"templateName"`
	Variables    map[string]string `json:"variables"`
}

// NewTemplatePlugin creates a new template plugin
func NewTemplatePlugin(id string) *TemplatePlugin {
	return &TemplatePlugin{
		id:      id,
		metrics: make(map[string]string),
	}
}

// ID implements Plugin interface
func (p *TemplatePlugin) ID() string {
	return p.id
}

// Initialize implements Plugin interface
func (p *TemplatePlugin) Initialize(ctx context.Context, config json.RawMessage) error {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	var cfg TemplateConfig
	if err := json.Unmarshal(config, &cfg); err != nil {
		return fmt.Errorf("failed to parse config: %v", err)
	}

	// Validate configuration
	if cfg.TemplateName == "" {
		return fmt.Errorf("template name is required")
	}

	// Create template
	tmpl, err := template.New(cfg.TemplateName).Parse("Template: {{ .Name }}")
	if err != nil {
		return fmt.Errorf("failed to parse template: %v", err)
	}

	p.config = &cfg
	p.tmpl = tmpl
	p.metrics["status"] = "initialized"

	return nil
}

// Start implements Plugin interface
func (p *TemplatePlugin) Start(ctx context.Context) error {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	if p.running {
		return fmt.Errorf("plugin already running")
	}

	p.running = true
	p.metrics["status"] = "running"

	// Start background processing if needed
	go func() {
		<-ctx.Done()
		p.mutex.Lock()
		p.running = false
		p.metrics["status"] = "stopped"
		p.mutex.Unlock()
	}()

	return nil
}

// Stop implements Plugin interface
func (p *TemplatePlugin) Stop(ctx context.Context) error {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	if !p.running {
		return nil
	}

	p.running = false
	p.metrics["status"] = "stopped"
	return nil
}

// Health implements Plugin interface
func (p *TemplatePlugin) Health() *plugin.HealthStatus {
	p.mutex.RLock()
	defer p.mutex.RUnlock()

	return &plugin.HealthStatus{
		Healthy: p.running,
		Status:  p.metrics["status"],
		Details: p.metrics,
	}
}

// ProcessTemplate applies the template with given variables
func (p *TemplatePlugin) ProcessTemplate(vars map[string]interface{}) (string, error) {
	p.mutex.RLock()
	defer p.mutex.RUnlock()

	if !p.running {
		return "", fmt.Errorf("plugin not running")
	}

	var output string
	builder := &strings.Builder{}

	// Merge configured variables with provided ones
	data := make(map[string]interface{})
	for k, v := range p.config.Variables {
		data[k] = v
	}
	for k, v := range vars {
		data[k] = v
	}

	if err := p.tmpl.Execute(builder, data); err != nil {
		return "", fmt.Errorf("template execution failed: %v", err)
	}

	output = builder.String()
	p.metrics["last_processed"] = time.Now().Format(time.RFC3339)
	p.metrics["total_processed"] = fmt.Sprintf("%d", p.getProcessedCount()+1)

	return output, nil
}

func (p *TemplatePlugin) getProcessedCount() int {
	count := 0
	if countStr, ok := p.metrics["total_processed"]; ok {
		if n, err := strconv.Atoi(countStr); err == nil {
			count = n
		}
	}
	return count
}
