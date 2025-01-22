package core

type Quantum interface {
	CUDACheck()
	TriggerQuantumRollback()
}

type Brutalizer interface {
	Brutalize(string) string
}

type Quantumizer interface {
	Quantize(string) string
}

const MAX_GPU_MEM = 4 * 1024 * 1024 * 1024 // 4GB Hard Limit

type Client struct {
	BaseURL string
}
