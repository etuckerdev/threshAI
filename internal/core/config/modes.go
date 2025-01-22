package config

var (
	BrutalLevel int
	AllowBrutal bool
)

type GenerationMode int

const (
	ModeDefault GenerationMode = iota
	ModeBrutal
	ModeQuantum
)
