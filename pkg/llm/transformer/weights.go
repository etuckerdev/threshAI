package transformer

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"

	"gorgonia.org/gorgonia"
	"gorgonia.org/tensor"
)

// WeightMapping defines how pre-trained weights map to our model
type WeightMapping struct {
	SourcePath string   // Path in pre-trained weights file
	TargetPath string   // Path in our model
	Transform  string   // Optional transformation to apply
	Args       []string // Optional arguments for transformation
}

// WeightConfig stores the mapping configuration for a pre-trained model
type WeightConfig struct {
	ModelType string // e.g. "gpt2", "mistral"
	Mappings  []WeightMapping
}

// DefaultGPT2Mapping returns the default weight mapping for GPT-2
func DefaultGPT2Mapping() WeightConfig {
	return WeightConfig{
		ModelType: "gpt2",
		Mappings: []WeightMapping{
			{SourcePath: "wte", TargetPath: "embedding"},
			{SourcePath: "h.{layer}.ln_1.weight", TargetPath: "blocks.{layer}.norm1"},
			{SourcePath: "h.{layer}.ln_2.weight", TargetPath: "blocks.{layer}.norm2"},
			{SourcePath: "h.{layer}.attn.c_attn", TargetPath: "blocks.{layer}.attention.qkv", Transform: "split_qkv"},
			{SourcePath: "h.{layer}.attn.c_proj", TargetPath: "blocks.{layer}.attention.outProj"},
			{SourcePath: "h.{layer}.mlp.c_fc", TargetPath: "blocks.{layer}.mlpW1"},
			{SourcePath: "h.{layer}.mlp.c_proj", TargetPath: "blocks.{layer}.mlpW2"},
			{SourcePath: "ln_f.weight", TargetPath: "lnf"},
		},
	}
}

// DefaultMistralMapping returns the default weight mapping for Mistral
func DefaultMistralMapping() WeightConfig {
	return WeightConfig{
		ModelType: "mistral",
		Mappings: []WeightMapping{
			{SourcePath: "embed_tokens", TargetPath: "embedding"},
			{SourcePath: "layers.{layer}.input_layernorm.weight", TargetPath: "blocks.{layer}.norm1"},
			{SourcePath: "layers.{layer}.post_attention_layernorm.weight", TargetPath: "blocks.{layer}.norm2"},
			{SourcePath: "layers.{layer}.self_attn", TargetPath: "blocks.{layer}.attention.qkv", Transform: "combine_qkv"},
			{SourcePath: "layers.{layer}.self_attn.o_proj", TargetPath: "blocks.{layer}.attention.outProj"},
			{SourcePath: "layers.{layer}.mlp.up_proj", TargetPath: "blocks.{layer}.mlpW1"},
			{SourcePath: "layers.{layer}.mlp.down_proj", TargetPath: "blocks.{layer}.mlpW2"},
			{SourcePath: "norm.weight", TargetPath: "lnf"},
		},
	}
}

// LoadPretrainedWeights loads weights from a pre-trained model into our architecture
func LoadPretrainedWeights(model *TransformerModel, weightsPath string, mappingConfig WeightConfig) error {
	// Read weights file
	data, err := ioutil.ReadFile(weightsPath)
	if err != nil {
		return fmt.Errorf("failed to read weights file: %v", err)
	}

	var pretrainedWeights map[string]interface{}
	if err := json.Unmarshal(data, &pretrainedWeights); err != nil {
		return fmt.Errorf("failed to parse weights file: %v", err)
	}

	// Apply mappings
	for _, mapping := range mappingConfig.Mappings {
		// Handle layer iteration if needed
		if strings.Contains(mapping.SourcePath, "{layer}") {
			for i := 0; i < len(model.blocks); i++ {
				sourcePath := strings.Replace(mapping.SourcePath, "{layer}", fmt.Sprintf("%d", i), -1)
				targetPath := strings.Replace(mapping.TargetPath, "{layer}", fmt.Sprintf("%d", i), -1)

				if err := applyMapping(model, pretrainedWeights, sourcePath, targetPath, mapping.Transform, mapping.Args); err != nil {
					return fmt.Errorf("failed to apply mapping for layer %d: %v", i, err)
				}
			}
		} else {
			if err := applyMapping(model, pretrainedWeights, mapping.SourcePath, mapping.TargetPath, mapping.Transform, mapping.Args); err != nil {
				return fmt.Errorf("failed to apply mapping: %v", err)
			}
		}
	}

	return nil
}

// applyMapping applies a single weight mapping
func applyMapping(model *TransformerModel, weights map[string]interface{}, sourcePath, targetPath, transform string, args []string) error {
	// Get source weights
	sourceWeight, ok := weights[sourcePath]
	if !ok {
		return fmt.Errorf("source weight not found: %s", sourcePath)
	}

	// Convert source weight to tensor
	data, ok := sourceWeight.([]float64)
	if !ok {
		return fmt.Errorf("invalid weight data type for %s", sourcePath)
	}

	// Apply any required transformations
	var transformedData []float64
	var err error
	switch transform {
	case "split_qkv":
		transformedData, err = transformSplitQKV(data)
	case "combine_qkv":
		transformedData, err = transformCombineQKV(data)
	default:
		transformedData = data
	}
	if err != nil {
		return fmt.Errorf("failed to transform weights: %v", err)
	}

	// Create tensor
	t := tensor.New(tensor.WithBacking(transformedData))

	// Get target node
	targetNode, err := getNodeByPath(model, targetPath)
	if err != nil {
		return fmt.Errorf("failed to get target node: %v", err)
	}

	// Set values
	if err := gorgonia.Let(targetNode, t); err != nil {
		return fmt.Errorf("failed to set node values: %v", err)
	}

	return nil
}

// Helper functions for transforming weights between architectures
func transformSplitQKV(data []float64) ([]float64, error) {
	// Split concatenated QKV weights into separate Q, K, V
	if len(data)%3 != 0 {
		return nil, fmt.Errorf("invalid QKV weight size")
	}
	size := len(data) / 3
	result := make([]float64, len(data))

	// Rearrange from [Q1,Q2...,K1,K2...,V1,V2...] to [Q1,K1,V1,Q2,K2,V2,...]
	for i := 0; i < size; i++ {
		result[i*3] = data[i]
		result[i*3+1] = data[size+i]
		result[i*3+2] = data[2*size+i]
	}
	return result, nil
}

func transformCombineQKV(data []float64) ([]float64, error) {
	// Combine separate Q, K, V into concatenated weights
	if len(data)%3 != 0 {
		return nil, fmt.Errorf("invalid QKV weight size")
	}
	size := len(data) / 3
	result := make([]float64, len(data))

	// Rearrange from [Q1,K1,V1,Q2,K2,V2,...] to [Q1,Q2...,K1,K2...,V1,V2...]
	for i := 0; i < len(data)/3; i++ {
		result[i] = data[i*3]
		result[size+i] = data[i*3+1]
		result[2*size+i] = data[i*3+2]
	}
	return result, nil
}

// getNodeByPath gets a node from the model using a dot-separated path
func getNodeByPath(model *TransformerModel, path string) (*gorgonia.Node, error) {
	parts := strings.Split(path, ".")
	var curr interface{} = model

	for _, part := range parts {
		switch v := curr.(type) {
		case *TransformerModel:
			switch part {
			case "embedding":
				curr = v.embedding
			case "blocks":
				curr = v.blocks
			case "lnf":
				curr = v.lnf
			case "head":
				curr = v.head
			default:
				return nil, fmt.Errorf("invalid path component for model: %s", part)
			}
		case []*TransformerBlock:
			idx := -1
			_, err := fmt.Sscanf(part, "%d", &idx)
			if err != nil || idx < 0 || idx >= len(v) {
				return nil, fmt.Errorf("invalid block index: %s", part)
			}
			curr = v[idx]
		case *TransformerBlock:
			switch part {
			case "attention":
				curr = v.attention
			case "norm1":
				curr = v.norm1
			case "norm2":
				curr = v.norm2
			case "mlpW1":
				curr = v.mlpW1
			case "mlpW2":
				curr = v.mlpW2
			default:
				return nil, fmt.Errorf("invalid path component for block: %s", part)
			}
		case *MultiHeadAttention:
			switch part {
			case "qkv":
				curr = v.qkv
			case "outProj":
				curr = v.outProj
			default:
				return nil, fmt.Errorf("invalid path component for attention: %s", part)
			}
		case *gorgonia.Node:
			return v, nil
		default:
			return nil, fmt.Errorf("invalid path component type at %s", part)
		}
	}

	if node, ok := curr.(*gorgonia.Node); ok {
		return node, nil
	}
	return nil, fmt.Errorf("path does not resolve to a node")
}
