package transformer

type Config struct {
    VocabSize      int
    MaxContext     int
    EmbedSize      int
    NumLayers      int
    NumHeads       int
    BatchSize      int
    Device         string // "cuda" or "cpu"
    
    // Tokenizer settings
    TokenizerType  string // "char" or "bpe"
    VocabPath     string // Path to vocabulary file for BPE
    MergePath     string // Path to merges file for BPE
    CheckpointDir string // Directory for saving/loading model checkpoints
}

func DefaultConfig() Config {
    return Config{
        VocabSize:     50257, // Standard GPT-2 vocabulary size
        MaxContext:    512,   // Context window size
        EmbedSize:    768,   // Embedding dimension
        NumLayers:    6,     // Number of transformer layers
        NumHeads:     12,    // Number of attention heads
        BatchSize:    32,    // Default batch size
        Device:       "cuda",
        
        // Default to GPT-2 tokenizer
        TokenizerType: "bpe",
        VocabPath:    "models/gpt2-vocab.json",
        MergePath:    "models/gpt2-merges.txt",
        CheckpointDir: "checkpoints",
    }
}
