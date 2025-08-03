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
	"sort"
	"strings"
	"time"
	"unicode"
)

// SummarizeText creates a summary with adaptive algorithm selection based on max length
// maxLength: <100 = extractive, <300 = hybrid, 300+ = abstractive
func SummarizeText(text string, maxLength int) SummaryResult {
	// Start metrics collection
	collector := StartMetricsCollection()
	startTime := time.Now()

	// Validate input
	if text == "" || maxLength <= 0 {
		return SummaryResult{
			Summary:        "",
			Method:         "none",
			ProcessingTime: time.Since(startTime),
			QualityMetrics: QualityMetrics{
				Accuracy:   0,
				Confidence: 0,
				Coverage:   0,
			},
		}
	}

	var result SummaryResult

	// Choose summarization method based on maxLength parameter
	if maxLength < 100 {
		// Extractive: Select key sentences
		result = summarizeExtractive(text, maxLength, collector)
		result.Method = "extractive"
	} else if maxLength < 300 {
		// Hybrid: Combine extractive with compression
		result = summarizeHybrid(text, maxLength, collector)
		result.Method = "hybrid"
	} else {
		// Abstractive: Generate new text
		result = summarizeAbstractive(text, maxLength, collector)
		result.Method = "abstractive"
	}

	result.ProcessingTime = time.Since(startTime)

	// Calculate compression ratio
	result.CompressionRatio = float64(len(result.Summary)) / float64(len(text))

	// Record metrics
	metrics := collector.GetMetrics()
	params := map[string]interface{}{
		"max_length":  maxLength,
		"text_length": len(text),
		"method":      result.Method,
	}

	RecordFunctionCall("SummarizeText", params, metrics, &result.QualityMetrics)

	return result
}

// summarizeExtractive performs sentence extraction based on importance
func summarizeExtractive(text string, maxLength int, collector *MetricsCollector) SummaryResult {
	collector.RecordProcessingTime("extractive_start")

	sentences := SplitIntoSentences(text)
	if len(sentences) == 0 {
		return SummaryResult{
			Summary: "",
			QualityMetrics: QualityMetrics{
				Accuracy:   0,
				Confidence: 0,
				Coverage:   0,
			},
		}
	}

	// Calculate sentence scores
	sentenceScores := make([]sentenceScore, len(sentences))
	
	// Get word frequencies
	wordFreq := getWordFrequencies(text)
	
	// Score each sentence
	for i, sentence := range sentences {
		score := calculateSentenceScore(sentence, wordFreq, i, len(sentences))
		sentenceScores[i] = sentenceScore{
			index:    i,
			sentence: sentence,
			score:    score,
		}
	}

	// Sort by score
	sort.Slice(sentenceScores, func(i, j int) bool {
		return sentenceScores[i].score > sentenceScores[j].score
	})

	// Select top sentences while respecting order and length limit
	selectedIndices := []int{}
	currentLength := 0
	
	// Only select sentences with reasonable scores (>0.2) or at least one sentence
	for _, ss := range sentenceScores {
		sentenceLength := len(ss.sentence)
		if currentLength+sentenceLength <= maxLength {
			// Skip sentences with very low scores unless it's the only/first sentence
			if ss.score > 0.2 || len(selectedIndices) == 0 {
				selectedIndices = append(selectedIndices, ss.index)
				currentLength += sentenceLength
			}
		}
	}

	// Sort selected indices to maintain original order
	sort.Ints(selectedIndices)

	// Build summary
	selectedSentences := []string{}
	for _, idx := range selectedIndices {
		selectedSentences = append(selectedSentences, sentences[idx])
	}

	summary := strings.Join(selectedSentences, " ")

	// Calculate quality metrics
	coverage := float64(len(selectedIndices)) / float64(len(sentences))
	confidence := 0.7 + (coverage * 0.2) // Higher coverage = higher confidence
	
	if confidence > 0.9 {
		confidence = 0.9
	}

	collector.IncrementAlgorithmSteps()

	return SummaryResult{
		Summary: summary,
		OriginalSentences: len(sentences),
		SummarySentences: len(selectedIndices),
		QualityMetrics: QualityMetrics{
			Accuracy:   0.75,
			Confidence: confidence,
			Coverage:   coverage,
		},
	}
}

// summarizeHybrid combines extraction with sentence compression
func summarizeHybrid(text string, maxLength int, collector *MetricsCollector) SummaryResult {
	collector.RecordProcessingTime("hybrid_start")

	// Start with extractive summary
	extractiveResult := summarizeExtractive(text, maxLength*2, collector) // Get more content initially
	
	if len(extractiveResult.Summary) <= maxLength {
		return extractiveResult
	}

	// Compress the extracted sentences
	sentences := SplitIntoSentences(extractiveResult.Summary)
	compressedSentences := []string{}
	
	for _, sentence := range sentences {
		compressed := compressSentence(sentence)
		compressedSentences = append(compressedSentences, compressed)
	}

	// Build compressed summary
	summary := strings.Join(compressedSentences, " ")
	
	// Trim to max length if still too long
	if len(summary) > maxLength {
		summary = truncateToWordBoundary(summary, maxLength)
	}

	// Calculate quality metrics
	originalSentences := len(SplitIntoSentences(text))
	summarySentences := len(compressedSentences)
	coverage := float64(summarySentences) / float64(originalSentences)
	
	confidence := 0.65 + (coverage * 0.15) // Slightly lower confidence due to compression
	if confidence > 0.85 {
		confidence = 0.85
	}

	collector.IncrementAlgorithmSteps()
	collector.RecordMemoryUsage()

	return SummaryResult{
		Summary:           summary,
		OriginalSentences: originalSentences,
		SummarySentences:  summarySentences,
		QualityMetrics: QualityMetrics{
			Accuracy:   0.70, // Lower accuracy due to compression
			Confidence: confidence,
			Coverage:   coverage,
		},
	}
}

// summarizeAbstractive generates new summary text
func summarizeAbstractive(text string, maxLength int, collector *MetricsCollector) SummaryResult {
	collector.RecordProcessingTime("abstractive_start")

	// For this implementation, we'll use a template-based approach
	// In a real system, this would use neural language models

	// Extract key information
	keyPhrases := ExtractKeyPhrases(text, 10)
	sentences := SplitIntoSentences(text)
	
	// Identify main topics
	mainTopics := []string{}
	for i, phrase := range keyPhrases {
		if i < 3 && phrase.Score > 0.5 {
			mainTopics = append(mainTopics, phrase.Text)
		}
	}

	// Generate summary based on templates
	summary := generateAbstractiveSummary(sentences, mainTopics, maxLength)

	// Calculate quality metrics
	originalSentences := len(sentences)
	summarySentences := len(SplitIntoSentences(summary))
	coverage := float64(len(mainTopics)) / float64(len(keyPhrases))
	
	confidence := 0.6 + (coverage * 0.2)
	if confidence > 0.8 {
		confidence = 0.8
	}

	collector.IncrementAlgorithmSteps()
	collector.RecordMemoryUsage()

	return SummaryResult{
		Summary:           summary,
		OriginalSentences: originalSentences,
		SummarySentences:  summarySentences,
		QualityMetrics: QualityMetrics{
			Accuracy:   0.65, // Lower accuracy for abstractive
			Confidence: confidence,
			Coverage:   coverage,
		},
	}
}

// Helper types and functions

type sentenceScore struct {
	index    int
	sentence string
	score    float64
}

func getWordFrequencies(text string) map[string]float64 {
	words := strings.Fields(strings.ToLower(text))
	freq := make(map[string]int)
	
	// Count frequencies
	for _, word := range words {
		word = strings.Trim(word, ".,!?;:\"'")
		if len(word) > 2 && !isStopWord(word) {
			freq[word]++
		}
	}

	// If no significant words found, return empty map
	if len(freq) == 0 {
		return make(map[string]float64)
	}

	// Normalize frequencies
	maxFreq := 0
	for _, count := range freq {
		if count > maxFreq {
			maxFreq = count
		}
	}

	normalized := make(map[string]float64)
	for word, count := range freq {
		normalized[word] = float64(count) / float64(maxFreq)
	}

	return normalized
}

func calculateSentenceScore(sentence string, wordFreq map[string]float64, position int, totalSentences int) float64 {
	words := strings.Fields(strings.ToLower(sentence))
	
	// If no significant words in the corpus, all sentences are equally unimportant
	if len(wordFreq) == 0 {
		return 0.1 // Very low score for all sentences
	}
	
	// Word frequency score
	freqScore := 0.0
	wordCount := 0
	for _, word := range words {
		word = strings.Trim(word, ".,!?;:\"'")
		if score, exists := wordFreq[word]; exists {
			freqScore += score
			wordCount++
		}
	}
	
	if wordCount > 0 {
		freqScore /= float64(wordCount)
	}

	// Position score (beginning and end sentences are often important)
	positionScore := 0.0
	if position == 0 {
		positionScore = 1.0
	} else if position == totalSentences-1 {
		positionScore = 0.7
	} else if float64(position) < float64(totalSentences)*0.2 {
		positionScore = 0.5
	}

	// Length score (prefer medium-length sentences)
	lengthScore := 0.0
	sentenceLength := len(words)
	if sentenceLength >= 10 && sentenceLength <= 30 {
		lengthScore = 1.0
	} else if sentenceLength >= 5 && sentenceLength <= 40 {
		lengthScore = 0.7
	} else {
		lengthScore = 0.3
	}

	// Combine scores
	return freqScore*0.5 + positionScore*0.3 + lengthScore*0.2
}

func compressSentence(sentence string) string {
	// Simple compression: remove less important phrases and words
	words := strings.Fields(sentence)
	compressed := []string{}
	
	skipNext := false
	for i, word := range words {
		if skipNext {
			skipNext = false
			continue
		}

		// Skip certain phrases
		if i < len(words)-1 {
			phrase := strings.ToLower(word + " " + words[i+1])
			if isRemovablePhrase(phrase) {
				skipNext = true
				continue
			}
		}

		// Skip parenthetical expressions
		if strings.HasPrefix(word, "(") && strings.HasSuffix(word, ")") {
			continue
		}

		// Keep important words
		if !isRemovableWord(strings.ToLower(word)) {
			compressed = append(compressed, word)
		}
	}

	return strings.Join(compressed, " ")
}

func isRemovablePhrase(phrase string) bool {
	removable := []string{
		"for example", "for instance", "such as", "as well",
		"in addition", "on the other hand", "in fact",
	}
	
	for _, r := range removable {
		if phrase == r {
			return true
		}
	}
	return false
}

func isRemovableWord(word string) bool {
	// Remove certain adverbs and fillers
	removable := []string{
		"very", "really", "quite", "rather", "somewhat",
		"actually", "basically", "essentially", "generally",
	}
	
	for _, r := range removable {
		if word == r {
			return true
		}
	}
	return false
}

func truncateToWordBoundary(text string, maxLength int) string {
	if len(text) <= maxLength {
		return text
	}

	// Find last space before maxLength
	lastSpace := maxLength
	for i := maxLength - 1; i >= 0; i-- {
		if unicode.IsSpace(rune(text[i])) {
			lastSpace = i
			break
		}
	}

	return strings.TrimSpace(text[:lastSpace]) + "..."
}

func generateAbstractiveSummary(sentences []string, topics []string, maxLength int) string {
	if len(sentences) == 0 {
		return ""
	}

	// Template-based generation
	var summary strings.Builder

	// Opening statement
	if len(topics) > 0 {
		if len(topics) == 1 {
			summary.WriteString("This text discusses ")
			summary.WriteString(topics[0])
			summary.WriteString(". ")
		} else {
			summary.WriteString("This text covers ")
			for i, topic := range topics {
				if i > 0 {
					if i == len(topics)-1 {
						summary.WriteString(" and ")
					} else {
						summary.WriteString(", ")
					}
				}
				summary.WriteString(topic)
			}
			summary.WriteString(". ")
		}
	}

	// Extract key points from sentences
	keyPoints := extractKeyPoints(sentences, 3)
	
	// Add key points
	for i, point := range keyPoints {
		if summary.Len()+len(point) > maxLength {
			break
		}
		
		if i == 0 {
			summary.WriteString("The main point is that ")
		} else if i == 1 {
			summary.WriteString("Additionally, ")
		} else {
			summary.WriteString("Furthermore, ")
		}
		
		// Ensure point ends with period
		summary.WriteString(strings.ToLower(string(point[0])) + point[1:])
		if !strings.HasSuffix(point, ".") {
			summary.WriteString(".")
		}
		summary.WriteString(" ")
	}

	result := summary.String()
	if len(result) > maxLength {
		result = truncateToWordBoundary(result, maxLength)
	}

	return strings.TrimSpace(result)
}

func extractKeyPoints(sentences []string, maxPoints int) []string {
	points := []string{}
	
	for _, sentence := range sentences {
		// Skip very short sentences
		if len(strings.Fields(sentence)) < 5 {
			continue
		}
		
		// Look for sentences with key indicators
		lower := strings.ToLower(sentence)
		if strings.Contains(lower, "important") ||
			strings.Contains(lower, "significant") ||
			strings.Contains(lower, "key") ||
			strings.Contains(lower, "main") ||
			strings.Contains(lower, "conclude") ||
			strings.Contains(lower, "result") {
			points = append(points, sentence)
			if len(points) >= maxPoints {
				break
			}
		}
	}

	// If not enough key points found, use first few sentences
	if len(points) < maxPoints {
		for _, sentence := range sentences {
			if len(strings.Fields(sentence)) >= 5 {
				found := false
				for _, p := range points {
					if p == sentence {
						found = true
						break
					}
				}
				if !found {
					points = append(points, sentence)
					if len(points) >= maxPoints {
						break
					}
				}
			}
		}
	}

	return points
}

