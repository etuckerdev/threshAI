// pkg/core/generator.go
package core

type Generator interface {
	Generate(prompt string) (string, error)
}
