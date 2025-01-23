package transformer

import (
	"encoding/gob"
	"fmt"
	"os"
	"threshAI/pkg/llm/tokenizer"

	"gorgonia.org/gorgonia"
	"gorgonia.org/tensor"
)

// TensorState holds tensor data and shape
type TensorState struct {
	Data  []float64
	Shape []int
}

// ModelState holds the weights and biases of the model
type ModelState struct {
	Config    Config
	Embedding TensorState
	Blocks    []BlockState
	LayerNorm TensorState
	Head      TensorState
}

type BlockState struct {
	QKV     TensorState
	OutProj TensorState
	MlpW1   TensorState
	MlpW2   TensorState
	Norm1   TensorState
	Norm2   TensorState
}

// getTensorState extracts data and shape from a Node's tensor value
func getTensorState(n *gorgonia.Node) (TensorState, error) {
	v := n.Value().(*tensor.Dense)
	data, ok := v.Data().([]float64)
	if !ok {
		return TensorState{}, fmt.Errorf("expected float64 tensor")
	}
	return TensorState{
		Data:  data,
		Shape: v.Shape(),
	}, nil
}

// SaveCheckpoint saves the model's state to a file
func (m *TransformerModel) SaveCheckpoint(path string) error {
	// Extract values
	embedding, err := getTensorState(m.embedding)
	if err != nil {
		return fmt.Errorf("failed to get embedding: %v", err)
	}

	layerNorm, err := getTensorState(m.lnf)
	if err != nil {
		return fmt.Errorf("failed to get layer norm: %v", err)
	}

	head, err := getTensorState(m.head)
	if err != nil {
		return fmt.Errorf("failed to get head: %v", err)
	}

	state := &ModelState{
		Config:    m.config,
		Embedding: embedding,
		Blocks:    make([]BlockState, len(m.blocks)),
		LayerNorm: layerNorm,
		Head:      head,
	}

	// Save each block's state
	for i, block := range m.blocks {
		qkv, err := getTensorState(block.attention.qkv)
		if err != nil {
			return fmt.Errorf("failed to get block %d qkv: %v", i, err)
		}

		outProj, err := getTensorState(block.attention.outProj)
		if err != nil {
			return fmt.Errorf("failed to get block %d outProj: %v", i, err)
		}

		mlpW1, err := getTensorState(block.mlpW1)
		if err != nil {
			return fmt.Errorf("failed to get block %d mlpW1: %v", i, err)
		}

		mlpW2, err := getTensorState(block.mlpW2)
		if err != nil {
			return fmt.Errorf("failed to get block %d mlpW2: %v", i, err)
		}

		norm1, err := getTensorState(block.norm1)
		if err != nil {
			return fmt.Errorf("failed to get block %d norm1: %v", i, err)
		}

		norm2, err := getTensorState(block.norm2)
		if err != nil {
			return fmt.Errorf("failed to get block %d norm2: %v", i, err)
		}

		state.Blocks[i] = BlockState{
			QKV:     qkv,
			OutProj: outProj,
			MlpW1:   mlpW1,
			MlpW2:   mlpW2,
			Norm1:   norm1,
			Norm2:   norm2,
		}
	}

	// Create the file
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create checkpoint file: %v", err)
	}
	defer f.Close()

	// Encode and save
	enc := gob.NewEncoder(f)
	if err := enc.Encode(state); err != nil {
		return fmt.Errorf("failed to encode model state: %v", err)
	}

	return nil
}

// createNode creates a new node with the given tensor state
func createNode(g *gorgonia.ExprGraph, state TensorState) *gorgonia.Node {
	t := tensor.New(tensor.WithShape(state.Shape...), tensor.WithBacking(state.Data))
	return gorgonia.NodeFromAny(g, t, gorgonia.WithName("loaded_tensor"))
}

// LoadCheckpoint loads the model's state from a file
func LoadCheckpoint(path string) (*TransformerModel, error) {
	// Open the file
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open checkpoint file: %v", err)
	}
	defer f.Close()

	// Decode the state
	var state ModelState
	dec := gob.NewDecoder(f)
	if err := dec.Decode(&state); err != nil {
		return nil, fmt.Errorf("failed to decode model state: %v", err)
	}

	// Create a new graph and model
	g := gorgonia.NewGraph()

	// Create nodes with loaded values
	embedding := createNode(g, state.Embedding)
	lnf := createNode(g, state.LayerNorm)
	head := createNode(g, state.Head)

	// Create blocks
	blocks := make([]*TransformerBlock, len(state.Blocks))
	for i, blockState := range state.Blocks {
		blocks[i] = &TransformerBlock{
			g: g,
			attention: &MultiHeadAttention{
				g:       g,
				qkv:     createNode(g, blockState.QKV),
				outProj: createNode(g, blockState.OutProj),
			},
			mlpW1: createNode(g, blockState.MlpW1),
			mlpW2: createNode(g, blockState.MlpW2),
			norm1: createNode(g, blockState.Norm1),
			norm2: createNode(g, blockState.Norm2),
		}
	}

	vm := gorgonia.NewTapeMachine(g)

	// Initialize tokenizer based on config
	var tok *tokenizer.Tokenizer
	if state.Config.TokenizerType == "bpe" {
		tok, err = tokenizer.LoadGPT2Tokenizer(state.Config.VocabPath, state.Config.MergePath)
		if err != nil {
			return nil, fmt.Errorf("failed to load tokenizer: %v", err)
		}
	} else {
		tok = tokenizer.NewTokenizer()
	}

	return &TransformerModel{
		g:         g,
		config:    state.Config,
		embedding: embedding,
		blocks:    blocks,
		lnf:       lnf,
		head:      head,
		vm:        vm,
		tokenizer: tok,
	}, nil
}
