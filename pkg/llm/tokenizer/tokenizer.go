package tokenizer

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

// Token represents a subword token and its ID
type Token struct {
	Text  string
	ID    int
	Count int
}

// Tokenizer implements Byte-Pair Encoding (BPE) tokenization
type Tokenizer struct {
	vocab     map[string]int // token text -> id
	merges    map[string]int // merge rule -> priority
	decoder   map[int]string // id -> token text
	unkToken  int
	padToken  int
	eosToken  int
	maxLength int
}

// NewTokenizer creates a new BPE tokenizer
func NewTokenizer() *Tokenizer {
	return &Tokenizer{
		vocab:     make(map[string]int),
		merges:    make(map[string]int),
		decoder:   make(map[int]string),
		unkToken:  0,
		padToken:  1,
		eosToken:  2,
		maxLength: 512,
	}
}

// LoadGPT2Tokenizer loads the pre-trained GPT-2 tokenizer
func LoadGPT2Tokenizer(vocabPath, mergePath string) (*Tokenizer, error) {
	t := NewTokenizer()

	// Load vocabulary
	vocabBytes, err := os.ReadFile(vocabPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read vocab file: %v", err)
	}
	if err := json.Unmarshal(vocabBytes, &t.vocab); err != nil {
		return nil, fmt.Errorf("failed to parse vocab: %v", err)
	}

	// Build decoder
	for token, id := range t.vocab {
		t.decoder[id] = token
	}

	// Load merges
	mergeBytes, err := os.ReadFile(mergePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read merges file: %v", err)
	}
	merges := strings.Split(string(mergeBytes), "\n")
	for i, merge := range merges {
		if merge == "" {
			continue
		}
		t.merges[merge] = i
	}

	return t, nil
}

// Encode converts text into token IDs using BPE
func (t *Tokenizer) Encode(text string, maxLength int) ([]int, error) {
	if maxLength == 0 {
		maxLength = t.maxLength
	}

	// Pre-tokenize into words
	words := strings.Fields(text)
	tokens := make([]int, 0)

	// Apply BPE to each word
	for _, word := range words {
		wordTokens := t.tokenizeWord(word)
		tokens = append(tokens, wordTokens...)

		if len(tokens) >= maxLength {
			tokens = tokens[:maxLength]
			break
		}
	}

	// Add EOS token if there's room
	if len(tokens) < maxLength {
		tokens = append(tokens, t.eosToken)
	}

	return tokens, nil
}

// Decode converts token IDs back into text
func (t *Tokenizer) Decode(tokens []int) string {
	parts := make([]string, len(tokens))
	for i, id := range tokens {
		if text, ok := t.decoder[id]; ok {
			parts[i] = text
		} else {
			parts[i] = "[UNK]"
		}
	}
	return strings.Join(parts, "")
}

func (t *Tokenizer) tokenizeWord(word string) []int {
	// Start with character-level tokens
	parts := strings.Split(word, "")

	for {
		bestPair := ""
		bestScore := -1

		// Find the best merge
		for i := 0; i < len(parts)-1; i++ {
			pair := parts[i] + parts[i+1]
			if score, ok := t.merges[pair]; ok {
				if score > bestScore {
					bestScore = score
					bestPair = pair
				}
			}
		}

		// No more merges possible
		if bestPair == "" {
			break
		}

		// Apply the merge
		newParts := make([]string, 0)
		i := 0
		for i < len(parts)-1 {
			if parts[i]+parts[i+1] == bestPair {
				newParts = append(newParts, bestPair)
				i += 2
			} else {
				newParts = append(newParts, parts[i])
				i++
			}
		}
		if i < len(parts) {
			newParts = append(newParts, parts[i])
		}
		parts = newParts
	}

	// Convert subwords to token IDs
	tokens := make([]int, len(parts))
	for i, part := range parts {
		if id, ok := t.vocab[part]; ok {
			tokens[i] = id
		} else {
			tokens[i] = t.unkToken
		}
	}

	return tokens
}

// Save saves the tokenizer state to files
func (t *Tokenizer) Save(vocabPath, mergePath string) error {
	// Save vocabulary
	vocabBytes, err := json.Marshal(t.vocab)
	if err != nil {
		return fmt.Errorf("failed to marshal vocab: %v", err)
	}
	if err := os.WriteFile(vocabPath, vocabBytes, 0644); err != nil {
		return fmt.Errorf("failed to write vocab file: %v", err)
	}

	// Save merges
	merges := make([]string, len(t.merges))
	for merge, priority := range t.merges {
		merges[priority] = merge
	}
	mergeStr := strings.Join(merges, "\n")
	if err := os.WriteFile(mergePath, []byte(mergeStr), 0644); err != nil {
		return fmt.Errorf("failed to write merges file: %v", err)
	}

	return nil
}
