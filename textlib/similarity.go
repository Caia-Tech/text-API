package textlib

import (
	"math"
	"sort"
	"strings"
	"unicode"
)

// SimilarityResult contains various similarity metrics between two texts
type SimilarityResult struct {
	// Character-level metrics
	CharacterOverlap     float64
	LevenshteinDistance  int
	NormalizedLevenshtein float64
	
	// Word-level metrics
	WordOverlap          float64
	JaccardIndex         float64
	DiceCoefficient      float64
	
	// Semantic metrics
	CosineSimilarity     float64
	TFIDFSimilarity      float64
	
	// Structural metrics
	StructuralSimilarity float64
	SentenceAlignment    []SentenceAlignment
	
	// N-gram metrics
	BigramSimilarity     float64
	TrigramSimilarity    float64
	
	// Longest common subsequence
	LCSLength            int
	LCSRatio             float64
}

type SentenceAlignment struct {
	Text1Index    int
	Text2Index    int
	Similarity    float64
	AlignmentType string // "exact", "similar", "different"
}

// TextDiff represents differences between two texts
type TextDiff struct {
	Operations     []DiffOperation
	AddedWords     []string
	RemovedWords   []string
	ChangedWords   []WordChange
	
	// Structural changes
	AddedSentences   []string
	RemovedSentences []string
	ReorderedSentences []SentenceReorder
	
	// Statistical changes
	WordCountDelta   int
	SentenceCountDelta int
	
	// Character-level stats
	CharsAdded       int
	CharsRemoved     int
	CharsUnchanged   int
}

type DiffOperation struct {
	Type     string // "insert", "delete", "equal", "replace"
	Text     string
	Position int
	Length   int
}

type WordChange struct {
	Original string
	Changed  string
	Position int
}

type SentenceReorder struct {
	Sentence      string
	OldPosition   int
	NewPosition   int
}

// CalculateSimilarity computes comprehensive similarity metrics between two texts
func CalculateSimilarity(text1, text2 string) *SimilarityResult {
	result := &SimilarityResult{}
	
	// Character-level similarity
	result.LevenshteinDistance = calculateLevenshtein(text1, text2)
	maxLen := max(len(text1), len(text2))
	if maxLen > 0 {
		result.NormalizedLevenshtein = 1.0 - float64(result.LevenshteinDistance)/float64(maxLen)
	}
	result.CharacterOverlap = calculateCharacterOverlap(text1, text2)
	
	// Word-level similarity
	words1 := extractWords(text1)
	words2 := extractWords(text2)
	result.WordOverlap = calculateWordSetOverlap(words1, words2)
	result.JaccardIndex = calculateJaccardIndex(words1, words2)
	result.DiceCoefficient = calculateDiceCoefficient(words1, words2)
	
	// Semantic similarity
	result.CosineSimilarity = calculateCosineSimilarity(words1, words2)
	result.TFIDFSimilarity = calculateTFIDFSimilarity(text1, text2)
	
	// Structural similarity
	sentences1 := SplitIntoSentences(text1)
	sentences2 := SplitIntoSentences(text2)
	result.StructuralSimilarity = calculateStructuralSimilarity(sentences1, sentences2)
	result.SentenceAlignment = alignSentences(sentences1, sentences2)
	
	// N-gram similarity
	result.BigramSimilarity = calculateNGramSimilarity(text1, text2, 2)
	result.TrigramSimilarity = calculateNGramSimilarity(text1, text2, 3)
	
	// Longest common subsequence
	result.LCSLength = calculateLCS(text1, text2)
	if maxLen > 0 {
		result.LCSRatio = float64(result.LCSLength) / float64(maxLen)
	}
	
	return result
}

// CalculateDiff computes the differences between two texts
func CalculateDiff(text1, text2 string) *TextDiff {
	diff := &TextDiff{
		Operations:         []DiffOperation{},
		AddedWords:        []string{},
		RemovedWords:      []string{},
		ChangedWords:      []WordChange{},
		AddedSentences:    []string{},
		RemovedSentences:  []string{},
		ReorderedSentences: []SentenceReorder{},
	}
	
	// Word-level diff
	words1 := extractWords(text1)
	words2 := extractWords(text2)
	wordDiff := calculateWordDiff(words1, words2)
	diff.AddedWords = wordDiff.added
	diff.RemovedWords = wordDiff.removed
	diff.ChangedWords = detectWordChanges(words1, words2)
	
	// Sentence-level diff
	sentences1 := SplitIntoSentences(text1)
	sentences2 := SplitIntoSentences(text2)
	sentenceDiff := calculateSentenceDiff(sentences1, sentences2)
	diff.AddedSentences = sentenceDiff.added
	diff.RemovedSentences = sentenceDiff.removed
	diff.ReorderedSentences = detectSentenceReorders(sentences1, sentences2)
	
	// Character-level operations
	diff.Operations = calculateDiffOperations(text1, text2)
	
	// Statistics
	diff.WordCountDelta = len(words2) - len(words1)
	diff.SentenceCountDelta = len(sentences2) - len(sentences1)
	
	// Character stats
	for _, op := range diff.Operations {
		switch op.Type {
		case "insert":
			diff.CharsAdded += len(op.Text)
		case "delete":
			diff.CharsRemoved += len(op.Text)
		case "equal":
			diff.CharsUnchanged += len(op.Text)
		}
	}
	
	return diff
}

// Similarity calculation functions

func calculateLevenshtein(s1, s2 string) int {
	r1 := []rune(s1)
	r2 := []rune(s2)
	
	if len(r1) == 0 {
		return len(r2)
	}
	if len(r2) == 0 {
		return len(r1)
	}
	
	// Create matrix
	matrix := make([][]int, len(r1)+1)
	for i := range matrix {
		matrix[i] = make([]int, len(r2)+1)
	}
	
	// Initialize first row and column
	for i := 0; i <= len(r1); i++ {
		matrix[i][0] = i
	}
	for j := 0; j <= len(r2); j++ {
		matrix[0][j] = j
	}
	
	// Fill matrix
	for i := 1; i <= len(r1); i++ {
		for j := 1; j <= len(r2); j++ {
			cost := 0
			if r1[i-1] != r2[j-1] {
				cost = 1
			}
			
			matrix[i][j] = min(
				matrix[i-1][j]+1,      // deletion
				min(
					matrix[i][j-1]+1,  // insertion
					matrix[i-1][j-1]+cost, // substitution
				),
			)
		}
	}
	
	return matrix[len(r1)][len(r2)]
}

func calculateCharacterOverlap(text1, text2 string) float64 {
	chars1 := make(map[rune]int)
	chars2 := make(map[rune]int)
	
	for _, r := range text1 {
		if !unicode.IsSpace(r) {
			chars1[r]++
		}
	}
	
	for _, r := range text2 {
		if !unicode.IsSpace(r) {
			chars2[r]++
		}
	}
	
	overlap := 0
	total := 0
	
	for char, count1 := range chars1 {
		count2 := chars2[char]
		overlap += min(count1, count2)
		total += count1
	}
	
	for char, count2 := range chars2 {
		if _, exists := chars1[char]; !exists {
			total += count2
		}
	}
	
	if total == 0 {
		return 0
	}
	
	return float64(overlap) / float64(total)
}

func calculateWordSetOverlap(words1, words2 []string) float64 {
	set1 := make(map[string]bool)
	set2 := make(map[string]bool)
	
	for _, w := range words1 {
		set1[strings.ToLower(w)] = true
	}
	
	for _, w := range words2 {
		set2[strings.ToLower(w)] = true
	}
	
	intersection := 0
	for w := range set1 {
		if set2[w] {
			intersection++
		}
	}
	
	union := len(set1) + len(set2) - intersection
	
	if union == 0 {
		return 0
	}
	
	return float64(intersection) / float64(union)
}

func calculateJaccardIndex(words1, words2 []string) float64 {
	// Same as word set overlap
	return calculateWordSetOverlap(words1, words2)
}

func calculateDiceCoefficient(words1, words2 []string) float64 {
	set1 := make(map[string]bool)
	set2 := make(map[string]bool)
	
	for _, w := range words1 {
		set1[strings.ToLower(w)] = true
	}
	
	for _, w := range words2 {
		set2[strings.ToLower(w)] = true
	}
	
	intersection := 0
	for w := range set1 {
		if set2[w] {
			intersection++
		}
	}
	
	if len(set1)+len(set2) == 0 {
		return 0
	}
	
	return 2.0 * float64(intersection) / float64(len(set1)+len(set2))
}

func calculateCosineSimilarity(words1, words2 []string) float64 {
	// Build frequency vectors
	freq1 := make(map[string]float64)
	freq2 := make(map[string]float64)
	
	for _, w := range words1 {
		freq1[strings.ToLower(w)]++
	}
	
	for _, w := range words2 {
		freq2[strings.ToLower(w)]++
	}
	
	// Get all unique words
	allWords := make(map[string]bool)
	for w := range freq1 {
		allWords[w] = true
	}
	for w := range freq2 {
		allWords[w] = true
	}
	
	// Calculate dot product and magnitudes
	dotProduct := 0.0
	magnitude1 := 0.0
	magnitude2 := 0.0
	
	for w := range allWords {
		f1 := freq1[w]
		f2 := freq2[w]
		
		dotProduct += f1 * f2
		magnitude1 += f1 * f1
		magnitude2 += f2 * f2
	}
	
	magnitude1 = math.Sqrt(magnitude1)
	magnitude2 = math.Sqrt(magnitude2)
	
	if magnitude1 == 0 || magnitude2 == 0 {
		return 0
	}
	
	return dotProduct / (magnitude1 * magnitude2)
}

func calculateTFIDFSimilarity(text1, text2 string) float64 {
	// Simple TF-IDF approximation
	docs := []string{text1, text2}
	
	// Calculate document frequencies
	df := make(map[string]int)
	for _, doc := range docs {
		words := extractWords(doc)
		seen := make(map[string]bool)
		
		for _, w := range words {
			lower := strings.ToLower(w)
			if !seen[lower] {
				df[lower]++
				seen[lower] = true
			}
		}
	}
	
	// Calculate TF-IDF vectors
	vectors := make([]map[string]float64, 2)
	
	for i, doc := range docs {
		words := extractWords(doc)
		tf := make(map[string]float64)
		
		// Calculate term frequencies
		for _, w := range words {
			tf[strings.ToLower(w)]++
		}
		
		// Normalize and apply IDF
		vectors[i] = make(map[string]float64)
		for term, freq := range tf {
			tf_norm := freq / float64(len(words))
			idf := math.Log(float64(len(docs)) / float64(df[term]))
			vectors[i][term] = tf_norm * idf
		}
	}
	
	// Calculate cosine similarity of TF-IDF vectors
	return calculateVectorCosineSimilarity(vectors[0], vectors[1])
}

func calculateVectorCosineSimilarity(v1, v2 map[string]float64) float64 {
	dotProduct := 0.0
	magnitude1 := 0.0
	magnitude2 := 0.0
	
	// Get all terms
	allTerms := make(map[string]bool)
	for t := range v1 {
		allTerms[t] = true
	}
	for t := range v2 {
		allTerms[t] = true
	}
	
	for term := range allTerms {
		val1 := v1[term]
		val2 := v2[term]
		
		dotProduct += val1 * val2
		magnitude1 += val1 * val1
		magnitude2 += val2 * val2
	}
	
	magnitude1 = math.Sqrt(magnitude1)
	magnitude2 = math.Sqrt(magnitude2)
	
	if magnitude1 == 0 || magnitude2 == 0 {
		return 0
	}
	
	return dotProduct / (magnitude1 * magnitude2)
}

func calculateStructuralSimilarity(sentences1, sentences2 []string) float64 {
	if len(sentences1) == 0 || len(sentences2) == 0 {
		return 0
	}
	
	// Compare sentence count
	countSim := 1.0 - math.Abs(float64(len(sentences1)-len(sentences2)))/
		float64(max(len(sentences1), len(sentences2)))
	
	// Compare average sentence lengths
	avgLen1 := calculateAverageSentenceLength(sentences1)
	avgLen2 := calculateAverageSentenceLength(sentences2)
	
	lengthSim := 1.0 - math.Abs(avgLen1-avgLen2)/math.Max(avgLen1, avgLen2)
	
	// Compare paragraph structure (if applicable)
	// This is simplified - could be enhanced with actual paragraph detection
	
	return (countSim + lengthSim) / 2
}

func calculateAverageSentenceLength(sentences []string) float64 {
	if len(sentences) == 0 {
		return 0
	}
	
	totalWords := 0
	for _, s := range sentences {
		totalWords += len(strings.Fields(s))
	}
	
	return float64(totalWords) / float64(len(sentences))
}

func alignSentences(sentences1, sentences2 []string) []SentenceAlignment {
	alignments := []SentenceAlignment{}
	
	// Create similarity matrix
	simMatrix := make([][]float64, len(sentences1))
	for i := range simMatrix {
		simMatrix[i] = make([]float64, len(sentences2))
		for j := range simMatrix[i] {
			simMatrix[i][j] = calculateSentenceSimilarity(sentences1[i], sentences2[j])
		}
	}
	
	// Find best alignments (simplified greedy approach)
	used1 := make([]bool, len(sentences1))
	used2 := make([]bool, len(sentences2))
	
	// First pass: find exact and high similarity matches
	for i := 0; i < len(sentences1); i++ {
		for j := 0; j < len(sentences2); j++ {
			if !used1[i] && !used2[j] && simMatrix[i][j] > 0.8 {
				alignType := "similar"
				if simMatrix[i][j] > 0.95 {
					alignType = "exact"
				}
				
				alignments = append(alignments, SentenceAlignment{
					Text1Index:    i,
					Text2Index:    j,
					Similarity:    simMatrix[i][j],
					AlignmentType: alignType,
				})
				
				used1[i] = true
				used2[j] = true
			}
		}
	}
	
	// Second pass: find remaining alignments
	for i := 0; i < len(sentences1); i++ {
		if !used1[i] {
			bestJ := -1
			bestSim := 0.0
			
			for j := 0; j < len(sentences2); j++ {
				if !used2[j] && simMatrix[i][j] > bestSim {
					bestJ = j
					bestSim = simMatrix[i][j]
				}
			}
			
			if bestJ != -1 && bestSim > 0.3 {
				alignments = append(alignments, SentenceAlignment{
					Text1Index:    i,
					Text2Index:    bestJ,
					Similarity:    bestSim,
					AlignmentType: "different",
				})
				used2[bestJ] = true
			}
		}
	}
	
	// Sort by text1 index
	sort.Slice(alignments, func(i, j int) bool {
		return alignments[i].Text1Index < alignments[j].Text1Index
	})
	
	return alignments
}

func calculateSentenceSimilarity(sent1, sent2 string) float64 {
	words1 := extractWords(sent1)
	words2 := extractWords(sent2)
	
	// Combine word overlap and cosine similarity
	overlap := calculateWordSetOverlap(words1, words2)
	cosine := calculateCosineSimilarity(words1, words2)
	
	return (overlap + cosine) / 2
}

func calculateNGramSimilarity(text1, text2 string, n int) float64 {
	ngrams1 := extractCharNGrams(text1, n)
	ngrams2 := extractCharNGrams(text2, n)
	
	set1 := make(map[string]bool)
	set2 := make(map[string]bool)
	
	for _, ng := range ngrams1 {
		set1[ng] = true
	}
	
	for _, ng := range ngrams2 {
		set2[ng] = true
	}
	
	intersection := 0
	for ng := range set1 {
		if set2[ng] {
			intersection++
		}
	}
	
	union := len(set1) + len(set2) - intersection
	
	if union == 0 {
		return 0
	}
	
	return float64(intersection) / float64(union)
}

func extractCharNGrams(text string, n int) []string {
	ngrams := []string{}
	runes := []rune(text)
	
	if len(runes) < n {
		return ngrams
	}
	
	for i := 0; i <= len(runes)-n; i++ {
		ngram := string(runes[i : i+n])
		ngrams = append(ngrams, ngram)
	}
	
	return ngrams
}

func calculateLCS(text1, text2 string) int {
	r1 := []rune(text1)
	r2 := []rune(text2)
	
	if len(r1) == 0 || len(r2) == 0 {
		return 0
	}
	
	// Create DP table
	dp := make([][]int, len(r1)+1)
	for i := range dp {
		dp[i] = make([]int, len(r2)+1)
	}
	
	// Fill DP table
	for i := 1; i <= len(r1); i++ {
		for j := 1; j <= len(r2); j++ {
			if r1[i-1] == r2[j-1] {
				dp[i][j] = dp[i-1][j-1] + 1
			} else {
				dp[i][j] = max(dp[i-1][j], dp[i][j-1])
			}
		}
	}
	
	return dp[len(r1)][len(r2)]
}

// Diff calculation functions

type wordDiffResult struct {
	added   []string
	removed []string
}

func calculateWordDiff(words1, words2 []string) wordDiffResult {
	result := wordDiffResult{
		added:   []string{},
		removed: []string{},
	}
	
	// Convert to sets
	set1 := make(map[string]int)
	set2 := make(map[string]int)
	
	for _, w := range words1 {
		set1[w]++
	}
	
	for _, w := range words2 {
		set2[w]++
	}
	
	// Find removed words
	for w, count1 := range set1 {
		count2 := set2[w]
		if count2 < count1 {
			for i := 0; i < count1-count2; i++ {
				result.removed = append(result.removed, w)
			}
		}
	}
	
	// Find added words
	for w, count2 := range set2 {
		count1 := set1[w]
		if count1 < count2 {
			for i := 0; i < count2-count1; i++ {
				result.added = append(result.added, w)
			}
		}
	}
	
	return result
}

type sentenceDiffResult struct {
	added   []string
	removed []string
}

func calculateSentenceDiff(sentences1, sentences2 []string) sentenceDiffResult {
	result := sentenceDiffResult{
		added:   []string{},
		removed: []string{},
	}
	
	// Create sets for exact matching
	set1 := make(map[string]bool)
	set2 := make(map[string]bool)
	
	for _, s := range sentences1 {
		set1[s] = true
	}
	
	for _, s := range sentences2 {
		set2[s] = true
	}
	
	// Find removed sentences
	for s := range set1 {
		if !set2[s] {
			// Check if there's a similar sentence (not exact match)
			similar := false
			for s2 := range set2 {
				if calculateSentenceSimilarity(s, s2) > 0.9 {
					similar = true
					break
				}
			}
			
			if !similar {
				result.removed = append(result.removed, s)
			}
		}
	}
	
	// Find added sentences
	for s := range set2 {
		if !set1[s] {
			// Check if there's a similar sentence (not exact match)
			similar := false
			for s1 := range set1 {
				if calculateSentenceSimilarity(s, s1) > 0.9 {
					similar = true
					break
				}
			}
			
			if !similar {
				result.added = append(result.added, s)
			}
		}
	}
	
	return result
}

func detectWordChanges(words1, words2 []string) []WordChange {
	changes := []WordChange{}
	
	// Use LCS algorithm to find word-level changes
	lcs := calculateWordLCS(words1, words2)
	
	i, j := 0, 0
	for _, word := range lcs {
		// Skip to the LCS word in both arrays
		for i < len(words1) && words1[i] != word {
			i++
		}
		for j < len(words2) && words2[j] != word {
			j++
		}
		
		// Check if there's a substitution nearby
		if i > 0 && j > 0 && i < len(words1) && j < len(words2) {
			w1 := words1[i-1]
			w2 := words2[j-1]
			
			if w1 != w2 && calculateLevenshtein(w1, w2) <= 2 {
				changes = append(changes, WordChange{
					Original: w1,
					Changed:  w2,
					Position: i - 1,
				})
			}
		}
		
		i++
		j++
	}
	
	return changes
}

func calculateWordLCS(words1, words2 []string) []string {
	if len(words1) == 0 || len(words2) == 0 {
		return []string{}
	}
	
	// Create DP table
	dp := make([][]int, len(words1)+1)
	for i := range dp {
		dp[i] = make([]int, len(words2)+1)
	}
	
	// Fill DP table
	for i := 1; i <= len(words1); i++ {
		for j := 1; j <= len(words2); j++ {
			if words1[i-1] == words2[j-1] {
				dp[i][j] = dp[i-1][j-1] + 1
			} else {
				dp[i][j] = max(dp[i-1][j], dp[i][j-1])
			}
		}
	}
	
	// Reconstruct LCS
	lcs := []string{}
	i, j := len(words1), len(words2)
	
	for i > 0 && j > 0 {
		if words1[i-1] == words2[j-1] {
			lcs = append([]string{words1[i-1]}, lcs...)
			i--
			j--
		} else if dp[i-1][j] > dp[i][j-1] {
			i--
		} else {
			j--
		}
	}
	
	return lcs
}

func detectSentenceReorders(sentences1, sentences2 []string) []SentenceReorder {
	reorders := []SentenceReorder{}
	
	// Create position maps
	pos1 := make(map[string]int)
	pos2 := make(map[string]int)
	
	for i, s := range sentences1 {
		pos1[s] = i
	}
	
	for i, s := range sentences2 {
		pos2[s] = i
	}
	
	// Find reordered sentences
	for s, newPos := range pos2 {
		if oldPos, exists := pos1[s]; exists && oldPos != newPos {
			reorders = append(reorders, SentenceReorder{
				Sentence:    s,
				OldPosition: oldPos,
				NewPosition: newPos,
			})
		}
	}
	
	// Sort by old position
	sort.Slice(reorders, func(i, j int) bool {
		return reorders[i].OldPosition < reorders[j].OldPosition
	})
	
	return reorders
}

func calculateDiffOperations(text1, text2 string) []DiffOperation {
	// Use Myers' diff algorithm (simplified)
	operations := []DiffOperation{}
	
	// For simplicity, use character-level diff
	r1 := []rune(text1)
	r2 := []rune(text2)
	
	// LCS-based diff
	lcs := calculateRuneLCS(r1, r2)
	
	i, j, lcsIdx := 0, 0, 0
	
	for i < len(r1) || j < len(r2) {
		if lcsIdx < len(lcs) && i < len(r1) && j < len(r2) && r1[i] == r2[j] && r1[i] == lcs[lcsIdx] {
			// Equal
			start := i
			for i < len(r1) && j < len(r2) && r1[i] == r2[j] {
				i++
				j++
				if lcsIdx < len(lcs) && i < len(r1) && r1[i] == lcs[lcsIdx] {
					lcsIdx++
				}
			}
			
			operations = append(operations, DiffOperation{
				Type:     "equal",
				Text:     string(r1[start:i]),
				Position: start,
				Length:   i - start,
			})
		} else if i < len(r1) && (lcsIdx >= len(lcs) || r1[i] != lcs[lcsIdx]) {
			// Delete
			start := i
			for i < len(r1) && (lcsIdx >= len(lcs) || j >= len(r2) || r1[i] != r2[j]) {
				i++
			}
			
			operations = append(operations, DiffOperation{
				Type:     "delete",
				Text:     string(r1[start:i]),
				Position: start,
				Length:   i - start,
			})
		} else if j < len(r2) {
			// Insert
			start := j
			for j < len(r2) && (i >= len(r1) || r1[i] != r2[j]) {
				j++
			}
			
			operations = append(operations, DiffOperation{
				Type:     "insert",
				Text:     string(r2[start:j]),
				Position: i,
				Length:   j - start,
			})
		}
	}
	
	return operations
}

func calculateRuneLCS(r1, r2 []rune) []rune {
	if len(r1) == 0 || len(r2) == 0 {
		return []rune{}
	}
	
	// Create DP table
	dp := make([][]int, len(r1)+1)
	for i := range dp {
		dp[i] = make([]int, len(r2)+1)
	}
	
	// Fill DP table
	for i := 1; i <= len(r1); i++ {
		for j := 1; j <= len(r2); j++ {
			if r1[i-1] == r2[j-1] {
				dp[i][j] = dp[i-1][j-1] + 1
			} else {
				dp[i][j] = max(dp[i-1][j], dp[i][j-1])
			}
		}
	}
	
	// Reconstruct LCS
	lcs := []rune{}
	i, j := len(r1), len(r2)
	
	for i > 0 && j > 0 {
		if r1[i-1] == r2[j-1] {
			lcs = append([]rune{r1[i-1]}, lcs...)
			i--
			j--
		} else if dp[i-1][j] > dp[i][j-1] {
			i--
		} else {
			j--
		}
	}
	
	return lcs
}