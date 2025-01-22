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

const MODEL_MINISTRAL_8B = "cas/ministral-8b-instruct-2410_q4km"

type Client struct {
	BaseURL string
}
