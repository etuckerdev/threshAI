package config

var (
	BrutalLevel   int
	AllowBrutal   bool
	SecurityModel string
	Quantize      string
)

type GenerationMode int

const (
	ModeDefault GenerationMode = iota
	ModeBrutal
	ModeQuantum
)

// Current generation mode configuration
var (
	CurrentGenerationMode GenerationMode = ModeDefault
)
