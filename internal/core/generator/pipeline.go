package generator

import (
	"errors"

	"github.com/etuckerdev/threshAI/internal/core/config"
	"github.com/etuckerdev/threshAI/internal/render"
)

type Brutalizer interface {
	Brutalize(input string) string
}

type brutalizer struct{}
type quantumizer struct{}

func (b *brutalizer) brutalizePayload(input string) string {
	return render.Brutalize(input)
}

func (q *quantumizer) quantumizePayload(input string) string {
	return render.Quantumize(input)
}

func ProcessInput(mode config.GenerationMode, input string) (string, error) {
	// Nuclear option enforcement
	if mode == config.ModeBrutal && !config.AllowBrutal {
		return "", errors.New("nuclear option required: brutal mode needs --unsafe")
	}

	switch mode {
	case config.ModeBrutal:
		b := &brutalizer{}
		return b.brutalizePayload(input), nil
	case config.ModeQuantum:
		q := &quantumizer{}
		return q.quantumizePayload(input), nil
	default:
		return input, nil
	}
}
