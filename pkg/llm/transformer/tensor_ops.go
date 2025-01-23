package transformer

import (
	"fmt"
	"math"
	"math/rand"
	"sort"

	"gorgonia.org/gorgonia"
	"gorgonia.org/tensor"
)

// TensorOps provides tensor operation utilities for the transformer model
type TensorOps struct {
	g *gorgonia.ExprGraph
}

// NewTensorOps creates a new TensorOps instance
func NewTensorOps(g *gorgonia.ExprGraph) *TensorOps {
	return &TensorOps{g: g}
}

// LayerNorm applies layer normalization to the input
func (ops *TensorOps) LayerNorm(input, scale *gorgonia.Node) (*gorgonia.Node, error) {
	mean, err := gorgonia.Mean(input, 1)
	if err != nil {
		return nil, fmt.Errorf("mean calculation failed: %v", err)
	}

	epsilon := gorgonia.NewScalar(ops.g, tensor.Float64, gorgonia.WithName("epsilon"))
	gorgonia.Let(epsilon, 1e-5)

	diff := gorgonia.Must(gorgonia.Sub(input, mean))
	variance := gorgonia.Must(gorgonia.Mean(gorgonia.Must(gorgonia.Square(diff)), 1))

	normalized := gorgonia.Must(gorgonia.HadamardDiv(
		diff,
		gorgonia.Must(gorgonia.Sqrt(gorgonia.Must(gorgonia.Add(variance, epsilon)))),
	))

	return gorgonia.Mul(normalized, scale)
}

// Gelu applies the Gaussian Error Linear Unit activation function
func (ops *TensorOps) Gelu(x *gorgonia.Node) (*gorgonia.Node, error) {
	// Initialize constants
	half := gorgonia.NewScalar(ops.g, tensor.Float64, gorgonia.WithName("half"))
	sqrt2OverPi := gorgonia.NewScalar(ops.g, tensor.Float64, gorgonia.WithName("sqrt2OverPi"))
	coeff := gorgonia.NewScalar(ops.g, tensor.Float64, gorgonia.WithName("coeff"))
	one := gorgonia.NewScalar(ops.g, tensor.Float64, gorgonia.WithName("one"))

	// Set values
	gorgonia.Let(half, 0.5)
	gorgonia.Let(sqrt2OverPi, 0.7978845608028654)
	gorgonia.Let(coeff, 0.044715)
	gorgonia.Let(one, 1.0)

	// Compute GELU
	cube := gorgonia.Must(gorgonia.Cube(x))
	inner := gorgonia.Must(gorgonia.Add(x, gorgonia.Must(gorgonia.Mul(cube, coeff))))
	tanh := gorgonia.Must(gorgonia.Tanh(gorgonia.Must(gorgonia.Mul(inner, sqrt2OverPi))))

	return gorgonia.Mul(
		gorgonia.Must(gorgonia.Mul(x, half)),
		gorgonia.Must(gorgonia.Add(one, tanh)),
	)
}

// ExtractLogits gets the logits for the last token
func (ops *TensorOps) ExtractLogits(t *tensor.Dense, vocabSize int) ([]float64, error) {
	data, ok := t.Data().([]float64)
	if !ok {
		return nil, fmt.Errorf("tensor data is not float64")
	}

	shape := t.Shape()
	if len(shape) < 2 || shape[1] != vocabSize {
		return nil, fmt.Errorf("invalid tensor shape: expected [_, %d], got %v", vocabSize, shape)
	}

	startIdx := len(data) - vocabSize
	if startIdx < 0 {
		return nil, fmt.Errorf("tensor too small for vocab size")
	}

	return data[startIdx:], nil
}

// Argmax returns the index of the maximum value in a slice
func (ops *TensorOps) Argmax(values []float64) (int, error) {
	if len(values) == 0 {
		return 0, fmt.Errorf("empty slice")
	}

	maxIdx := 0
	maxVal := values[0]
	for i, v := range values[1:] {
		if v > maxVal {
			maxVal = v
			maxIdx = i + 1
		}
	}

	return maxIdx, nil
}

// logitScore represents a token and its logit score
type logitScore struct {
	token int
	score float64
}

// computeSoftmax computes softmax probabilities for logits
func computeSoftmax(logits []float64) []float64 {
	maxLogit := logits[0]
	for _, l := range logits[1:] {
		if l > maxLogit {
			maxLogit = l
		}
	}

	expSum := 0.0
	probs := make([]float64, len(logits))
	for i, l := range logits {
		exp := math.Exp(l - maxLogit)
		probs[i] = exp
		expSum += exp
	}

	for i := range probs {
		probs[i] /= expSum
	}

	return probs
}

// SampleTopK performs top-k sampling on logits
func (ops *TensorOps) SampleTopK(logits []float64, k int) (int, error) {
	if k <= 0 || k > len(logits) {
		return 0, fmt.Errorf("invalid k value: %d", k)
	}

	// Create slice of token-score pairs
	scores := make([]logitScore, len(logits))
	for i, score := range logits {
		scores[i] = logitScore{token: i, score: score}
	}

	// Sort by score in descending order
	sort.Slice(scores, func(i, j int) bool {
		return scores[i].score > scores[j].score
	})

	// Take top k scores and compute softmax
	topK := scores[:k]
	topKLogits := make([]float64, k)
	for i, s := range topK {
		topKLogits[i] = s.score
	}

	// Compute probabilities with softmax
	probs := computeSoftmax(topKLogits)

	// Sample from the distribution
	r := rand.Float64()
	cumulativeProb := 0.0
	for i, prob := range probs {
		cumulativeProb += prob
		if r <= cumulativeProb {
			return topK[i].token, nil
		}
	}

	// Fallback to last token if no selection made
	return topK[len(topK)-1].token, nil
}

// SampleNucleus performs nucleus (top-p) sampling on logits
func (ops *TensorOps) SampleNucleus(logits []float64, p float64) (int, error) {
	if p <= 0 || p > 1 {
		return 0, fmt.Errorf("invalid p value: %f", p)
	}

	// Create slice of token-score pairs
	scores := make([]logitScore, len(logits))
	for i, score := range logits {
		scores[i] = logitScore{token: i, score: score}
	}

	// Sort by score in descending order
	sort.Slice(scores, func(i, j int) bool {
		return scores[i].score > scores[j].score
	})

	// Compute softmax probabilities
	probs := computeSoftmax(logits)

	// Find cutoff index where cumulative probability exceeds p
	cumulativeProb := 0.0
	var cutoff int
	for i, s := range scores {
		cumulativeProb += probs[s.token]
		if cumulativeProb >= p {
			cutoff = i + 1
			break
		}
	}

	// If no cutoff found, use all tokens
	if cutoff == 0 {
		cutoff = len(scores)
	}

	// Sample from the selected tokens
	selectedScores := scores[:cutoff]
	selectedLogits := make([]float64, cutoff)
	for i, s := range selectedScores {
		selectedLogits[i] = s.score
	}

	// Recompute probabilities for selected tokens
	selectedProbs := computeSoftmax(selectedLogits)

	// Sample from the distribution
	r := rand.Float64()
	cumulativeProb = 0.0
	for i, prob := range selectedProbs {
		cumulativeProb += prob
		if r <= cumulativeProb {
			return selectedScores[i].token, nil
		}
	}

	// Fallback to highest probability token
	return selectedScores[0].token, nil
}

// CreateInputTensor creates a tensor for model input
func (ops *TensorOps) CreateInputTensor(tokens []float64) *tensor.Dense {
	return tensor.New(
		tensor.WithBacking(tokens),
		tensor.WithShape(1, len(tokens)),
	)
}
