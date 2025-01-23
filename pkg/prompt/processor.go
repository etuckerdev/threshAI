package prompt

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"text/template"
)

// OutputProcessor handles the conversion of prompts into different output formats
type OutputProcessor struct {
	formatters map[string]OutputFormatter
}

// OutputFormatter defines the interface for different output formats
type OutputFormatter interface {
	Format(data interface{}) ([]byte, error)
}

// MarkdownFormatter implements markdown output
type MarkdownFormatter struct{}

// JSONFormatter implements JSON output
type JSONFormatter struct {
	PrettyPrint bool
}

// XMLFormatter implements XML output
type XMLFormatter struct {
	PrettyPrint bool
}

// NewOutputProcessor creates a new processor with default formatters
func NewOutputProcessor() *OutputProcessor {
	return &OutputProcessor{
		formatters: map[string]OutputFormatter{
			"md":   &MarkdownFormatter{},
			"json": &JSONFormatter{PrettyPrint: true},
			"xml":  &XMLFormatter{PrettyPrint: true},
		},
	}
}

// RegisterFormatter adds a new output formatter
func (p *OutputProcessor) RegisterFormatter(name string, formatter OutputFormatter) {
	p.formatters[name] = formatter
}

// Process converts input data to the specified format
func (p *OutputProcessor) Process(format string, data interface{}) ([]byte, error) {
	formatter, ok := p.formatters[format]
	if !ok {
		return nil, fmt.Errorf("unsupported format: %s", format)
	}
	return formatter.Format(data)
}

// Format implementations for each formatter type

func (m *MarkdownFormatter) Format(data interface{}) ([]byte, error) {
	tmpl := `# {{.Title}}

## Problem Statement
{{.ProblemStatement}}

## Code Targets
{{range .CodeTargets}}
- File: {{.File}}
  Function: {{.Function}}
{{end}}

## Success Metrics
{{range .Metrics}}
- {{.Category}}: {{.Value}}
{{end}}
`
	t, err := template.New("markdown").Parse(tmpl)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (j *JSONFormatter) Format(data interface{}) ([]byte, error) {
	if j.PrettyPrint {
		return json.MarshalIndent(data, "", "  ")
	}
	return json.Marshal(data)
}

func (x *XMLFormatter) Format(data interface{}) ([]byte, error) {
	if x.PrettyPrint {
		return xml.MarshalIndent(data, "", "  ")
	}
	return xml.Marshal(data)
}

// Helper types for template rendering
type CodeTarget struct {
	File     string
	Function string
}

type Metric struct {
	Category string
	Value    string
}

type TemplateData struct {
	Title            string
	ProblemStatement string
	CodeTargets      []CodeTarget
	Metrics          []Metric
}
