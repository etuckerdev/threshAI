package pipeline

import (
	"bytes"
	"fmt"
	"sync"
	"text/template"
)

// Chain represents a prompt execution pipeline
type Chain struct {
	steps     []Step
	variables map[string]interface{}
	mu        sync.RWMutex
}

// Step represents a single execution step in the chain
type Step struct {
	ID       string
	Template string
	Handler  StepHandler
}

// StepHandler defines the interface for step execution
type StepHandler interface {
	Execute(ctx *Context) error
}

// Context holds the execution context for a chain
type Context struct {
	Variables map[string]interface{}
	Input     interface{}
	Output    interface{}
}

// NewChain creates a new execution chain
func NewChain() *Chain {
	return &Chain{
		steps:     make([]Step, 0),
		variables: make(map[string]interface{}),
	}
}

// AddStep adds a new step to the chain
func (c *Chain) AddStep(id string, template string, handler StepHandler) {
	c.steps = append(c.steps, Step{
		ID:       id,
		Template: template,
		Handler:  handler,
	})
}

// SetVariable sets a variable in the chain context
func (c *Chain) SetVariable(key string, value interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.variables[key] = value
}

// GetVariable retrieves a variable from the chain context
func (c *Chain) GetVariable(key string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	val, ok := c.variables[key]
	return val, ok
}

// Execute runs the chain with the given input
func (c *Chain) Execute(input interface{}) (interface{}, error) {
	ctx := &Context{
		Variables: c.variables,
		Input:     input,
	}

	var lastOutput interface{}
	for _, step := range c.steps {
		// Process template variables
		processedTemplate, err := c.processTemplate(step.Template, ctx)
		if err != nil {
			return nil, fmt.Errorf("template processing error in step %s: %w", step.ID, err)
		}

		// Update context with processed template
		ctx.Input = processedTemplate

		// Execute the step
		if err := step.Handler.Execute(ctx); err != nil {
			return nil, fmt.Errorf("execution error in step %s: %w", step.ID, err)
		}

		lastOutput = ctx.Output
		ctx.Input = lastOutput // Pass output to next step's input
	}

	return lastOutput, nil
}

// processTemplate applies variable substitution to the template
func (c *Chain) processTemplate(tmpl string, ctx *Context) (string, error) {
	t, err := template.New("step").Parse(tmpl)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	err = t.Execute(&buf, ctx.Variables)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

// BatchExecutor handles parallel chain execution
type BatchExecutor struct {
	chains        []*Chain
	maxConcurrent int
}

// NewBatchExecutor creates a new batch executor
func NewBatchExecutor(maxConcurrent int) *BatchExecutor {
	if maxConcurrent <= 0 {
		maxConcurrent = 10 // Default concurrent chains
	}
	return &BatchExecutor{
		chains:        make([]*Chain, 0),
		maxConcurrent: maxConcurrent,
	}
}

// AddChain adds a chain to the batch
func (b *BatchExecutor) AddChain(chain *Chain) {
	b.chains = append(b.chains, chain)
}

// ExecuteAll runs all chains in parallel with rate limiting
func (b *BatchExecutor) ExecuteAll(input interface{}) []error {
	var wg sync.WaitGroup
	sem := make(chan struct{}, b.maxConcurrent)
	errors := make([]error, len(b.chains))

	for i, chain := range b.chains {
		wg.Add(1)
		go func(idx int, c *Chain) {
			defer wg.Done()
			sem <- struct{}{}        // Acquire semaphore
			defer func() { <-sem }() // Release semaphore

			_, err := c.Execute(input)
			if err != nil {
				errors[idx] = fmt.Errorf("chain %d error: %w", idx, err)
			}
		}(i, chain)
	}

	wg.Wait()
	return errors
}

// Predefined step handlers for common operations

// TemplateHandler processes a template with variables
type TemplateHandler struct {
	Template string
}

func (h *TemplateHandler) Execute(ctx *Context) error {
	tmpl, err := template.New("handler").Parse(h.Template)
	if err != nil {
		return err
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, ctx.Variables)
	if err != nil {
		return err
	}

	ctx.Output = buf.String()
	return nil
}

// TransformHandler applies a transformation function
type TransformHandler struct {
	Transform func(interface{}) (interface{}, error)
}

func (h *TransformHandler) Execute(ctx *Context) error {
	output, err := h.Transform(ctx.Input)
	if err != nil {
		return err
	}
	ctx.Output = output
	return nil
}
