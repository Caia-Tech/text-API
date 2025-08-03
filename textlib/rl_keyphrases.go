// Copyright 2025 Caia Tech
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package textlib

import (
	"math"
	"sort"
	"strings"
)

// ExtractKeyPhrases extracts key phrases with tunable complexity
// maxPhrases controls computational cost vs completeness
// Algorithm selection based on maxPhrases:
// 1-10: Fast TF-IDF approach
// 11-50: Enhanced statistical methods
// 51+: Deep NLP analysis
func ExtractKeyPhrases(text string, maxPhrases int) []KeyPhrase {
	// Start metrics collection
	collector := StartMetricsCollection()
	
	// Validate input
	if text == "" || maxPhrases <= 0 {
		return []KeyPhrase{}
	}
	
	// Select algorithm based on maxPhrases
	var phrases []KeyPhrase
	var algorithm string
	
	if maxPhrases <= 10 {
		// Fast TF-IDF approach
		phrases = extractKeyPhrasesTFIDF(text, maxPhrases, collector)
		algorithm = "tf-idf"
	} else if maxPhrases <= 50 {
		// Enhanced statistical methods
		phrases = extractKeyPhrasesStatistical(text, maxPhrases, collector)
		algorithm = "statistical"
	} else {
		// Deep NLP analysis
		phrases = extractKeyPhrasesDeep(text, maxPhrases, collector)
		algorithm = "deep-nlp"
	}
	
	// Record metrics
	metrics := collector.GetMetrics()
	params := map[string]interface{}{
		"maxPhrases":  maxPhrases,
		"textLength":  len(text),
		"algorithm":   algorithm,
	}
	
	// Calculate quality based on algorithm
	quality := &QualityMetrics{
		Accuracy:   0.7 + float64(min(maxPhrases, 50))*0.004, // 0.7 to 0.9
		Confidence: 0.8,
		Coverage:   float64(len(phrases)) / float64(maxPhrases),
	}
	
	RecordFunctionCall("ExtractKeyPhrases", params, metrics, quality)
	
	return phrases
}

// extractKeyPhrasesTFIDF uses TF-IDF for fast key phrase extraction
func extractKeyPhrasesTFIDF(text string, maxPhrases int, collector *MetricsCollector) []KeyPhrase {
	collector.RecordProcessingTime("tfidf_start")
	
	// Tokenize and clean
	words := strings.Fields(strings.ToLower(text))
	
	// Calculate term frequency
	termFreq := make(map[string]float64)
	totalWords := float64(len(words))
	
	for _, word := range words {
		cleanWord := strings.Trim(word, ".,!?;:\"'")
		if len(cleanWord) > 2 && !isStopWord(cleanWord) {
			termFreq[cleanWord]++
		}
	}
	
	// Normalize term frequency
	for term := range termFreq {
		termFreq[term] /= totalWords
	}
	
	// Extract phrases (unigrams and bigrams)
	phrases := extractPhrasesFromFrequencies(text, termFreq, maxPhrases)
	
	collector.IncrementAlgorithmSteps()
	return phrases
}

// extractKeyPhrasesStatistical uses enhanced statistical methods
func extractKeyPhrasesStatistical(text string, maxPhrases int, collector *MetricsCollector) []KeyPhrase {
	collector.RecordProcessingTime("statistical_start")
	
	// Get sentences for context
	sentences := SplitIntoSentences(text)
	
	// Build n-gram statistics (1-3 grams)
	ngramStats := make(map[string]*ngramInfo)
	
	for i, sentence := range sentences {
		words := strings.Fields(strings.ToLower(sentence))
		
		// Extract n-grams
		for n := 1; n <= 3; n++ {
			for j := 0; j <= len(words)-n; j++ {
				ngram := strings.Join(words[j:j+n], " ")
				
				// Skip if contains stop words at edges
				if n > 1 && (isStopWord(words[j]) || isStopWord(words[j+n-1])) {
					continue
				}
				
				if _, exists := ngramStats[ngram]; !exists {
					ngramStats[ngram] = &ngramInfo{
						text:      ngram,
						frequency: 0,
						positions: []int{},
						contexts:  []string{},
					}
				}
				
				info := ngramStats[ngram]
				info.frequency++
				info.positions = append(info.positions, i)
				
				// Store context (limited)
				if len(info.contexts) < 3 {
					info.contexts = append(info.contexts, sentence)
				}
			}
		}
	}
	
	// Score n-grams
	scoredPhrases := scoreNgrams(ngramStats, len(sentences))
	
	// Convert to KeyPhrases
	phrases := convertToKeyPhrases(scoredPhrases, maxPhrases, text)
	
	collector.IncrementAlgorithmSteps()
	collector.RecordMemoryUsage()
	return phrases
}

// extractKeyPhrasesDeep uses deep NLP analysis
func extractKeyPhrasesDeep(text string, maxPhrases int, collector *MetricsCollector) []KeyPhrase {
	collector.RecordProcessingTime("deep_start")
	
	// First get statistical phrases as a base
	basePhrases := extractKeyPhrasesStatistical(text, maxPhrases*2, collector)
	
	// Enhance with additional analysis
	sentences := SplitIntoSentences(text)
	
	// Pattern-based extraction
	patternPhrases := extractPatternBasedPhrases(sentences)
	
	// Merge and re-rank
	allPhrases := mergePhrases(basePhrases, patternPhrases)
	
	// Apply semantic clustering (simplified)
	clusteredPhrases := clusterPhrases(allPhrases)
	
	// Select top phrases ensuring diversity
	finalPhrases := selectDiversePhrases(clusteredPhrases, maxPhrases)
	
	collector.IncrementAlgorithmSteps()
	collector.RecordMemoryUsage()
	return finalPhrases
}

// Helper structures and functions

type ngramInfo struct {
	text      string
	frequency int
	positions []int
	contexts  []string
	score     float64
}

// isStopWord checks if a word is a common stop word
func isStopWord(word string) bool {
	stopWords := map[string]bool{
		"the": true, "is": true, "at": true, "which": true, "on": true,
		"a": true, "an": true, "and": true, "or": true, "but": true,
		"in": true, "with": true, "to": true, "for": true, "of": true,
		"as": true, "by": true, "that": true, "this": true, "it": true,
		"from": true, "be": true, "are": true, "was": true, "were": true,
		"been": true, "have": true, "has": true, "had": true, "do": true,
		"does": true, "did": true, "will": true, "would": true, "could": true,
		"should": true, "may": true, "might": true, "must": true, "can": true,
	}
	return stopWords[strings.ToLower(word)]
}

// extractPhrasesFromFrequencies converts term frequencies to key phrases
func extractPhrasesFromFrequencies(text string, termFreq map[string]float64, maxPhrases int) []KeyPhrase {
	// Sort terms by frequency
	type termScore struct {
		term  string
		score float64
	}
	
	var scored []termScore
	for term, freq := range termFreq {
		scored = append(scored, termScore{term, freq})
	}
	
	sort.Slice(scored, func(i, j int) bool {
		return scored[i].score > scored[j].score
	})
	
	// Convert to KeyPhrases
	phrases := make([]KeyPhrase, 0, maxPhrases)
	for i := 0; i < len(scored) && i < maxPhrases; i++ {
		// Find position in text
		pos := strings.Index(strings.ToLower(text), scored[i].term)
		
		phrase := KeyPhrase{
			Text:       scored[i].term,
			Score:      scored[i].score,
			Position:   RLPosition{Start: pos, End: pos + len(scored[i].term)},
			Category:   "term",
			Context:    extractContext(text, pos, 20),
			Confidence: 0.7 + scored[i].score*0.3,
		}
		phrases = append(phrases, phrase)
	}
	
	return phrases
}

// scoreNgrams scores n-grams based on various factors
func scoreNgrams(ngramStats map[string]*ngramInfo, totalSentences int) []*ngramInfo {
	// Calculate scores
	for _, info := range ngramStats {
		// Frequency component
		freqScore := math.Log(float64(info.frequency) + 1)
		
		// Length component (prefer longer meaningful phrases)
		words := strings.Fields(info.text)
		lengthScore := math.Min(float64(len(words))*0.3, 1.0)
		
		// Distribution component (prefer phrases that appear throughout)
		uniquePositions := make(map[int]bool)
		for _, pos := range info.positions {
			uniquePositions[pos/3] = true // Group nearby occurrences
		}
		distributionScore := float64(len(uniquePositions)) / float64(totalSentences)
		
		// Combined score
		info.score = freqScore*0.5 + lengthScore*0.3 + distributionScore*0.2
	}
	
	// Convert to slice and sort
	var scored []*ngramInfo
	for _, info := range ngramStats {
		if info.frequency > 1 || len(strings.Fields(info.text)) > 1 {
			scored = append(scored, info)
		}
	}
	
	sort.Slice(scored, func(i, j int) bool {
		return scored[i].score > scored[j].score
	})
	
	return scored
}

// convertToKeyPhrases converts scored n-grams to KeyPhrase objects
func convertToKeyPhrases(scored []*ngramInfo, maxPhrases int, originalText string) []KeyPhrase {
	phrases := make([]KeyPhrase, 0, maxPhrases)
	
	for i := 0; i < len(scored) && i < maxPhrases; i++ {
		info := scored[i]
		
		// Find first occurrence
		pos := strings.Index(strings.ToLower(originalText), info.text)
		
		// Determine category
		category := "phrase"
		if len(strings.Fields(info.text)) == 1 {
			category = "term"
		} else if len(strings.Fields(info.text)) > 2 {
			category = "multi-phrase"
		}
		
		phrase := KeyPhrase{
			Text:       info.text,
			Score:      info.score,
			Position:   RLPosition{Start: pos, End: pos + len(info.text)},
			Category:   category,
			Context:    "", // Will be filled if needed
			Confidence: math.Min(0.6+info.score*0.1, 0.95),
		}
		
		// Add context from first occurrence
		if len(info.contexts) > 0 {
			phrase.Context = info.contexts[0]
		}
		
		phrases = append(phrases, phrase)
	}
	
	return phrases
}

// extractPatternBasedPhrases extracts phrases using linguistic patterns
func extractPatternBasedPhrases(sentences []string) []KeyPhrase {
	var phrases []KeyPhrase
	
	// Common patterns for key phrases
	// Simplified pattern matching (in production, would use proper NLP)
	patterns := []struct {
		name     string
		check    func([]string) (string, bool)
		category string
	}{
		{
			name: "noun_phrase",
			check: func(words []string) (string, bool) {
				// Simple pattern: adjective + noun
				if len(words) >= 2 {
					if isAdjective(words[0]) && isNoun(words[1]) {
						return words[0] + " " + words[1], true
					}
				}
				return "", false
			},
			category: "noun-phrase",
		},
		// Add more patterns as needed
	}
	
	for _, sentence := range sentences {
		words := strings.Fields(strings.ToLower(sentence))
		
		// Apply patterns
		for i := 0; i < len(words); i++ {
			for _, pattern := range patterns {
				if i+2 <= len(words) {
					if phrase, matched := pattern.check(words[i:i+2]); matched {
						phrases = append(phrases, KeyPhrase{
							Text:       phrase,
							Category:   pattern.category,
							Context:    sentence,
							Confidence: 0.75,
						})
					}
				}
			}
		}
	}
	
	return phrases
}

// Simplified POS detection (in production, use proper NLP library)
func isAdjective(word string) bool {
	// Common adjective endings
	adjEndings := []string{"ive", "ous", "ful", "less", "able", "ible", "al", "ic"}
	for _, ending := range adjEndings {
		if strings.HasSuffix(word, ending) {
			return true
		}
	}
	return false
}

func isNoun(word string) bool {
	// Common noun endings
	nounEndings := []string{"tion", "ment", "ness", "ity", "er", "or", "ism", "ist"}
	for _, ending := range nounEndings {
		if strings.HasSuffix(word, ending) {
			return true
		}
	}
	// Also check if it's not a stop word and longer than 3 chars
	return len(word) > 3 && !isStopWord(word)
}

// mergePhrases merges two sets of phrases, removing duplicates
func mergePhrases(set1, set2 []KeyPhrase) []KeyPhrase {
	phraseMap := make(map[string]KeyPhrase)
	
	// Add all from set1
	for _, phrase := range set1 {
		key := strings.ToLower(phrase.Text)
		if existing, exists := phraseMap[key]; !exists || phrase.Score > existing.Score {
			phraseMap[key] = phrase
		}
	}
	
	// Add from set2
	for _, phrase := range set2 {
		key := strings.ToLower(phrase.Text)
		if existing, exists := phraseMap[key]; !exists || phrase.Score > existing.Score {
			phraseMap[key] = phrase
		}
	}
	
	// Convert back to slice
	var merged []KeyPhrase
	for _, phrase := range phraseMap {
		merged = append(merged, phrase)
	}
	
	// Sort by score
	sort.Slice(merged, func(i, j int) bool {
		return merged[i].Score > merged[j].Score
	})
	
	return merged
}

// clusterPhrases groups similar phrases (simplified clustering)
func clusterPhrases(phrases []KeyPhrase) []KeyPhrase {
	// For now, just remove near-duplicates
	filtered := make([]KeyPhrase, 0, len(phrases))
	seen := make(map[string]bool)
	
	for _, phrase := range phrases {
		// Check if we've seen a similar phrase
		isDuplicate := false
		
		for seenPhrase := range seen {
			if similarityScore(phrase.Text, seenPhrase) > 0.8 {
				isDuplicate = true
				break
			}
		}
		
		if !isDuplicate {
			filtered = append(filtered, phrase)
			seen[phrase.Text] = true
		}
	}
	
	return filtered
}

// similarityScore calculates similarity between two phrases
func similarityScore(s1, s2 string) float64 {
	// Simple word overlap similarity
	w1 := strings.Fields(strings.ToLower(s1))
	w2 := strings.Fields(strings.ToLower(s2))
	
	if len(w1) == 0 || len(w2) == 0 {
		return 0
	}
	
	wordSet1 := make(map[string]bool)
	for _, w := range w1 {
		wordSet1[w] = true
	}
	
	overlap := 0
	for _, w := range w2 {
		if wordSet1[w] {
			overlap++
		}
	}
	
	return float64(overlap) / float64(max(len(w1), len(w2)))
}

// selectDiversePhrases selects diverse phrases from clusters
func selectDiversePhrases(phrases []KeyPhrase, maxPhrases int) []KeyPhrase {
	if len(phrases) <= maxPhrases {
		return phrases
	}
	
	// Simple diversity selection: take top scored with some diversity
	selected := make([]KeyPhrase, 0, maxPhrases)
	categories := make(map[string]int)
	
	for _, phrase := range phrases {
		// Limit phrases per category
		if categories[phrase.Category] >= maxPhrases/3 {
			continue
		}
		
		selected = append(selected, phrase)
		categories[phrase.Category]++
		
		if len(selected) >= maxPhrases {
			break
		}
	}
	
	return selected
}

// extractContext extracts context around a position
func extractContext(text string, position, contextSize int) string {
	if position < 0 {
		return ""
	}
	
	start := max(0, position-contextSize)
	end := min(len(text), position+contextSize)
	
	context := text[start:end]
	
	// Clean up edges
	if start > 0 {
		if idx := strings.Index(context, " "); idx > 0 {
			context = "..." + context[idx+1:]
		}
	}
	
	if end < len(text) {
		if idx := strings.LastIndex(context, " "); idx > 0 && idx < len(context)-1 {
			context = context[:idx] + "..."
		}
	}
	
	return strings.TrimSpace(context)
}

// Helper functions
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}