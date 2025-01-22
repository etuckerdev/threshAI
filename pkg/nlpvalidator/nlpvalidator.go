package nlpvalidator

import (
	"fmt"
	"math"
	"regexp"
	"strings"

	"github.com/kljensen/snowball"
	"gonum.org/v1/gonum/mat"
)

// SplitSentences splits text into sentences using basic punctuation rules
func SplitSentences(text string) []string {
	// Split on sentence-ending punctuation while handling abbreviations
	re := regexp.MustCompile(`([.!?]+\\s+|…\\s+|$)`)
	sentences := re.Split(text, -1)

	// Filter out empty strings
	result := make([]string, 0, len(sentences))
	for _, s := range sentences {
		if len(strings.TrimSpace(s)) > 0 {
			result = append(result, s)
		}
	}
	return result
}

// CalculateTFIDFConsistency calculates semantic similarity between sentences
func CalculateTFIDFConsistency(sentences []string) float32 {
	if len(sentences) == 0 {
		return 0.0 // No sentences to compare
	}

	// Step 1: Tokenize sentences and filter out empty sentences
	var validTokens [][]string
	validSentences := []string{} // Keep valid sentences
	for _, s := range sentences {
		toks := strings.Fields(s)
		if len(toks) > 0 {
			validTokens = append(validTokens, toks)
			validSentences = append(validSentences, s)
		}
	}

	if len(validTokens) < 2 {
		return 1.0 // Not enough valid sentences to compare
	}

	// Step 2: Build common vocabulary
	vocabulary := make(map[string]int)
	for _, toks := range validTokens {
		for _, word := range toks {
			if _, ok := vocabulary[word]; !ok {
				vocabulary[word] = len(vocabulary)
			}
		}
	}

	// Step 3: Compute TF vectors based on common vocabulary
	tfVectors := make([]*mat.VecDense, len(validTokens))
	for i, toks := range validTokens {
		vector := mat.NewVecDense(len(vocabulary), nil)
		termCounts := make(map[string]int)
		for _, term := range toks {
			termCounts[term]++
		}
		for term, index := range vocabulary {
			tf := float64(termCounts[term]) / float64(len(toks))
			vector.SetVec(index, tf)
		}
		tfVectors[i] = vector
	}

	// Step 4: Compute IDF vector
	idfVector := mat.NewVecDense(len(vocabulary), nil)
	for term, index := range vocabulary {
		idf := inverseDocumentFrequency(term, validTokens)
		idfVector.SetVec(index, idf)
	}

	// Step 5: Compute TF-IDF vectors
	tfidfVectors := make([]*mat.VecDense, len(validTokens))
	for i := range validTokens {
		tfidfVector := &mat.VecDense{}
		tfidfVector.MulElemVec(tfVectors[i], idfVector)
		tfidfVectors[i] = tfidfVector
	}
	fmt.Printf("TFIDF Vectors: %v\\n", tfidfVectors)

	// Step 6: Calculate cosine similarity between all pairs
	var totalSimilarity float32
	pairs := 0
	for i := 0; i < len(tfidfVectors); i++ {
		for j := i + 1; j < len(tfidfVectors); j++ {
			similarity := cosineSimilarity(tfidfVectors[i], tfidfVectors[j])
			totalSimilarity += similarity
			pairs++
		}
	}

	if pairs == 0 {
		return 0.0 // No pairs to compare
	}
	return totalSimilarity / float32(pairs)
}

// cosineSimilarity calculates cosine similarity between two vectors
func cosineSimilarity(vecA, vecB *mat.VecDense) float32 {
	dot := mat.Dot(vecA, vecB)
	normA := mat.Norm(vecA, 2)
	normB := mat.Norm(vecB, 2)

	if normA == 0 || normB == 0 {
		return 0.0
	}
	if math.IsNaN(float64(dot / (normA * normB))) {
		return 0.0
	}
	return float32(dot / (normA * normB))
}

// DetectRepetition detects repeated phrases in text
func DetectRepetition(text string) float32 {
	words := tokenize(text)
	ngramSize := 3
	ngramCounts := make(map[string]int)

	// Count n-gram occurrences
	for i := 0; i <= len(words)-ngramSize; i++ {
		ngram := strings.Join(words[i:i+ngramSize], " ")
		ngramCounts[ngram]++
	}

	// Calculate repetition score
	var repeatedCount int
	for _, count := range ngramCounts {
		if count > 1 {
			repeatedCount++
		}
	}

	return float32(repeatedCount) / float32(len(ngramCounts))
}

// tokenize splits text into normalized terms
func tokenize(text string) []string {
	// Remove punctuation and convert to lowercase
	re := regexp.MustCompile(`[^\\p{L}\\p{N}]+`)
	text = re.ReplaceAllString(text, " ")
	text = strings.ToLower(text)

	// Split into words and stem
	words := strings.Fields(text)
	for i, word := range words {
		stemmed, err := snowball.Stem(word, "english", true)
		if err == nil {
			words[i] = stemmed
		}
	}
	if len(words) == 0 {
		return nil // Return nil for empty token list
	}
	if len(words) == 0 {
		return nil // Return nil for empty token list
	}
	return words
}

// tfidf calculates TF-IDF for a term in a document collection
func tfidf(term string, documents [][]string) float64 {
	tf := termFrequency(term, documents[0]) // Using first document for TF
	idf := inverseDocumentFrequency(term, documents)
	return tf * idf
}

// termFrequency calculates term frequency in a document
func termFrequency(term string, document []string) float64 {
	count := 0
	for _, t := range document {
		if t == term {
			count++
		}
	}
	return float64(count) / float64(len(document))
}

// inverseDocumentFrequency calculates inverse document frequency in document collection
func inverseDocumentFrequency(term string, documents [][]string) float64 {
	docCount := 0
	for _, doc := range documents {
		for _, t := range doc {
			if t == term {
				docCount++
				break // Term found in document, move to next document
			}
		}
	}

	if docCount == 0 {
		return 0 // Avoid division by zero
	}

	return math.Log(float64(len(documents)) / float64(docCount))
}
