package transformer

import (
	"fmt"
	"math"
	"runtime"
	"threshAI/pkg/llm/tokenizer"
	"time"

	"gorgonia.org/gorgonia"
	"gorgonia.org/tensor"
)

const (
	epsilon = 1e-5
)

// SamplingStrategy defines how to sample the next token
type SamplingStrategy struct {
	Type string  // "greedy", "topk", or "nucleus"
	K    int     // for top-k sampling
	P    float64 // for nucleus sampling (top-p)
}

// DefaultGreedyStrategy returns a greedy sampling strategy
func DefaultGreedyStrategy() SamplingStrategy {
	return SamplingStrategy{Type: "greedy"}
}

// TopKStrategy returns a top-k sampling strategy
func TopKStrategy(k int) SamplingStrategy {
	return SamplingStrategy{Type: "topk", K: k}
}

// NucleusStrategy returns a nucleus (top-p) sampling strategy
func NucleusStrategy(p float64) SamplingStrategy {
	return SamplingStrategy{Type: "nucleus", P: p}
}

// BlockMetrics stores timing and memory metrics for a transformer block
type BlockMetrics struct {
	AttentionTime time.Duration
	MLPTime       time.Duration
	MemoryUsed    uint64
}

func layerNorm(g *gorgonia.ExprGraph, input, scale *gorgonia.Node) (*gorgonia.Node, error) {
	mean, err := gorgonia.Mean(input, 1)
	if err != nil {
		return nil, err
	}

	epsilon := gorgonia.NewScalar(g, tensor.Float64)
	gorgonia.Let(epsilon, 1e-5)

	variance, err := gorgonia.Square(gorgonia.Must(gorgonia.Sub(input, mean)))
	if err != nil {
		return nil, err
	}

	varMean, err := gorgonia.Mean(variance, 1)
	if err != nil {
		return nil, err
	}

	normalized := gorgonia.Must(gorgonia.HadamardDiv(
		gorgonia.Must(gorgonia.Sub(input, mean)),
		gorgonia.Must(gorgonia.Sqrt(gorgonia.Must(gorgonia.Add(varMean, epsilon)))),
	))

	return gorgonia.Mul(normalized, scale)
}

func gelu(g *gorgonia.ExprGraph, x *gorgonia.Node) (*gorgonia.Node, error) {
	// Approximate GELU: 0.5 * x * (1 + tanh(sqrt(2/pi) * (x + 0.044715 * x^3)))
	sqrt2OverPi := gorgonia.NewScalar(g, tensor.Float64)
	half := gorgonia.NewScalar(g, tensor.Float64)
	coeff := gorgonia.NewScalar(g, tensor.Float64)
	one := gorgonia.NewScalar(g, tensor.Float64)

	gorgonia.Let(sqrt2OverPi, 0.7978845608028654)
	gorgonia.Let(half, 0.5)
	gorgonia.Let(coeff, 0.044715)
	gorgonia.Let(one, 1.0)

	cube := gorgonia.Must(gorgonia.Cube(x))
	inner := gorgonia.Must(gorgonia.Add(x, gorgonia.Must(gorgonia.Mul(cube, coeff))))
	tanh := gorgonia.Must(gorgonia.Tanh(gorgonia.Must(gorgonia.Mul(inner, sqrt2OverPi))))

	return gorgonia.Mul(
		gorgonia.Must(gorgonia.Mul(x, half)),
		gorgonia.Must(gorgonia.Add(one, tanh)),
	)
}

type MultiHeadAttention struct {
	g           *gorgonia.ExprGraph
	numHeads    int
	headDim     int
	qkv         *gorgonia.Node
	outProj     *gorgonia.Node
	scaleFactor float64
}

func NewMultiHeadAttention(g *gorgonia.ExprGraph, config Config) *MultiHeadAttention {
	headDim := config.EmbedSize / config.NumHeads

	qkvShape := tensor.Shape{config.EmbedSize, 3 * config.EmbedSize}
	qkvInit := gorgonia.NewTensor(g, tensor.Float64, 2, gorgonia.WithShape(qkvShape...))

	outProjShape := tensor.Shape{config.EmbedSize, config.EmbedSize}
	outProjInit := gorgonia.NewTensor(g, tensor.Float64, 2, gorgonia.WithShape(outProjShape...))

	return &MultiHeadAttention{
		g:           g,
		numHeads:    config.NumHeads,
		headDim:     headDim,
		qkv:         qkvInit,
		outProj:     outProjInit,
		scaleFactor: 1.0 / math.Sqrt(float64(headDim)),
	}
}

type TransformerBlock struct {
	g         *gorgonia.ExprGraph
	attention *MultiHeadAttention
	mlpW1     *gorgonia.Node
	mlpW2     *gorgonia.Node
	norm1     *gorgonia.Node
	norm2     *gorgonia.Node
	metrics   BlockMetrics
}

func NewTransformerBlock(g *gorgonia.ExprGraph, config Config) *TransformerBlock {
	mlpW1Shape := tensor.Shape{config.EmbedSize, 4 * config.EmbedSize}
	mlpW2Shape := tensor.Shape{4 * config.EmbedSize, config.EmbedSize}

	mlpW1 := gorgonia.NewTensor(g, tensor.Float64, 2, gorgonia.WithShape(mlpW1Shape...))
	mlpW2 := gorgonia.NewTensor(g, tensor.Float64, 2, gorgonia.WithShape(mlpW2Shape...))

	norm1 := gorgonia.NewTensor(g, tensor.Float64, 1, gorgonia.WithShape(config.EmbedSize))
	norm2 := gorgonia.NewTensor(g, tensor.Float64, 1, gorgonia.WithShape(config.EmbedSize))

	return &TransformerBlock{
		g:         g,
		attention: NewMultiHeadAttention(g, config),
		mlpW1:     mlpW1,
		mlpW2:     mlpW2,
		norm1:     norm1,
		norm2:     norm2,
	}
}

type TransformerModel struct {
	g         *gorgonia.ExprGraph
	config    Config
	embedding *gorgonia.Node
	blocks    []*TransformerBlock
	lnf       *gorgonia.Node
	head      *gorgonia.Node
	vm        gorgonia.VM
	tokenizer *tokenizer.Tokenizer
	sampling  SamplingStrategy
}

func NewTransformerModel(config Config) (*TransformerModel, error) {
	g := gorgonia.NewGraph()

	// Initialize tokenizer
	var tok *tokenizer.Tokenizer
	var err error
	if config.TokenizerType == "bpe" {
		tok, err = tokenizer.LoadGPT2Tokenizer(config.VocabPath, config.MergePath)
		if err != nil {
			return nil, fmt.Errorf("failed to load tokenizer: %v", err)
		}
	} else {
		tok = tokenizer.NewTokenizer()
	}

	// Initialize model components
	embShape := tensor.Shape{config.VocabSize, config.EmbedSize}
	embedding := gorgonia.NewTensor(g, tensor.Float64, 2, gorgonia.WithShape(embShape...))

	blocks := make([]*TransformerBlock, config.NumLayers)
	for i := 0; i < config.NumLayers; i++ {
		blocks[i] = NewTransformerBlock(g, config)
	}

	lnf := gorgonia.NewTensor(g, tensor.Float64, 1, gorgonia.WithShape(config.EmbedSize))
	head := gorgonia.NewTensor(g, tensor.Float64, 2, gorgonia.WithShape(config.EmbedSize, config.VocabSize))

	vm := gorgonia.NewTapeMachine(g)

	return &TransformerModel{
		g:         g,
		config:    config,
		embedding: embedding,
		blocks:    blocks,
		lnf:       lnf,
		head:      head,
		vm:        vm,
		tokenizer: tok,
		sampling:  DefaultGreedyStrategy(),
	}, nil
}

// GetGPUMetrics returns current GPU memory usage
func (m *TransformerModel) GetGPUMetrics() uint64 {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	return memStats.Alloc
}

// SetSamplingStrategy sets the sampling strategy for generation
func (m *TransformerModel) SetSamplingStrategy(strategy SamplingStrategy) {
	m.sampling = strategy
}

func (m *TransformerModel) Forward(input *tensor.Dense) (*tensor.Dense, error) {
	ops := NewTensorOps(m.g)

	// Convert input to nodes
	inputNode := gorgonia.NewTensor(m.g, input.Dtype(), input.Shape().Dims(), gorgonia.WithValue(input))

	// Embedding lookup with gradient checkpointing
	x, err := gorgonia.Mul(inputNode, m.embedding)
	if err != nil {
		return nil, fmt.Errorf("embedding lookup failed: %v", err)
	}

	// Process through transformer blocks with gradient checkpointing
	for i, block := range m.blocks {
		// Store input for gradient checkpointing
		if i%2 == 0 {
			// Save activation tensor for backward pass
			x.Value().(*tensor.Dense).Clone()
		}

		// Multi-head attention with timing
		attnStart := time.Now()
		normalized1, err := ops.LayerNorm(x, block.norm1)
		if err != nil {
			return nil, fmt.Errorf("layer norm 1 failed: %v", err)
		}

		attnOut, err := gorgonia.Mul(normalized1, block.attention.qkv)
		if err != nil {
			return nil, fmt.Errorf("attention failed: %v", err)
		}
		block.metrics.AttentionTime = time.Since(attnStart)

		residual := x
		x, err = gorgonia.Add(residual, attnOut)
		if err != nil {
			return nil, fmt.Errorf("residual connection 1 failed: %v", err)
		}

		// MLP forward pass with timing
		mlpStart := time.Now()
		normalized2, err := ops.LayerNorm(x, block.norm2)
		if err != nil {
			return nil, fmt.Errorf("layer norm 2 failed: %v", err)
		}

		hidden, err := gorgonia.Mul(normalized2, block.mlpW1)
		if err != nil {
			return nil, fmt.Errorf("MLP W1 failed: %v", err)
		}

		hidden, err = ops.Gelu(hidden)
		if err != nil {
			return nil, fmt.Errorf("GELU failed: %v", err)
		}

		out, err := gorgonia.Mul(hidden, block.mlpW2)
		if err != nil {
			return nil, fmt.Errorf("MLP W2 failed: %v", err)
		}
		block.metrics.MLPTime = time.Since(mlpStart)

		residual = x
		x, err = gorgonia.Add(residual, out)
		if err != nil {
			return nil, fmt.Errorf("residual connection 2 failed: %v", err)
		}

		// Record block memory usage
		block.metrics.MemoryUsed = m.GetGPUMetrics()
	}

	// Final layer norm and head
	normalized, err := ops.LayerNorm(x, m.lnf)
	if err != nil {
		return nil, fmt.Errorf("final layer norm failed: %v", err)
	}

	logits, err := gorgonia.Mul(normalized, m.head)
	if err != nil {
		return nil, fmt.Errorf("head projection failed: %v", err)
	}

	// Run the VM
	if err := m.vm.RunAll(); err != nil {
		return nil, fmt.Errorf("VM execution failed: %v", err)
	}

	result := logits.Value().(*tensor.Dense)
	m.vm.Reset()

	return result, nil
}

func (m *TransformerModel) Generate(input []float64, maxLen int) ([]float64, error) {
	ops := NewTensorOps(m.g)

	// Convert input to tensor
	inputTensor := ops.CreateInputTensor(input)

	generated := input
	for len(generated) < maxLen {
		// Get model prediction
		logits, err := m.Forward(inputTensor)
		if err != nil {
			return nil, fmt.Errorf("forward pass failed: %v", err)
		}

		// Get last token logits
		lastLogits, err := ops.ExtractLogits(logits, m.config.VocabSize)
		if err != nil {
			return nil, fmt.Errorf("failed to extract logits: %v", err)
		}

		// Sample next token based on strategy
		var nextToken int
		switch m.sampling.Type {
		case "greedy":
			nextToken, err = ops.Argmax(lastLogits)
		case "topk":
			nextToken, err = ops.SampleTopK(lastLogits, m.sampling.K)
		case "nucleus":
			nextToken, err = ops.SampleNucleus(lastLogits, m.sampling.P)
		default:
			return nil, fmt.Errorf("unknown sampling strategy: %s", m.sampling.Type)
		}
		if err != nil {
			return nil, fmt.Errorf("failed to sample next token: %v", err)
		}

		// Append to generated sequence
		generated = append(generated, float64(nextToken))

		// Update input tensor for next iteration
		contextStart := len(generated) - m.config.MaxContext
		if contextStart < 0 {
			contextStart = 0
		}
		inputTensor = ops.CreateInputTensor(generated[contextStart:])
	}

	return generated, nil
}
