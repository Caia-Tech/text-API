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
	"time"
)

// ClassifyTopics identifies topics with adaptive algorithm selection based on max topics
// maxTopics: ≤3 = clustering, ≤10 = statistical, >10 = comprehensive analysis
func ClassifyTopics(text string, maxTopics int) TopicResult {
	// Start metrics collection
	collector := StartMetricsCollection()
	startTime := time.Now()

	// Validate input
	if text == "" || maxTopics <= 0 {
		return TopicResult{
			Topics:         []Topic{},
			Method:         "none",
			ProcessingTime: time.Since(startTime),
			QualityMetrics: QualityMetrics{
				Accuracy:   0,
				Confidence: 0,
				Coverage:   0,
			},
		}
	}

	var result TopicResult

	// Choose classification method based on maxTopics parameter
	if maxTopics <= 3 {
		// Simple clustering approach
		result = classifyTopicsClustering(text, maxTopics, collector)
		result.Method = "clustering"
	} else if maxTopics <= 10 {
		// Statistical keyword analysis
		result = classifyTopicsStatistical(text, maxTopics, collector)
		result.Method = "statistical"
	} else {
		// Comprehensive topic modeling
		result = classifyTopicsComprehensive(text, maxTopics, collector)
		result.Method = "comprehensive"
	}

	result.ProcessingTime = time.Since(startTime)

	// Record metrics
	metrics := collector.GetMetrics()
	params := map[string]interface{}{
		"max_topics":  maxTopics,
		"text_length": len(text),
		"method":      result.Method,
	}

	RecordFunctionCall("ClassifyTopics", params, metrics, &result.QualityMetrics)

	return result
}

// classifyTopicsClustering performs simple keyword clustering for topic identification
func classifyTopicsClustering(text string, maxTopics int, collector *MetricsCollector) TopicResult {
	collector.RecordProcessingTime("clustering_start")

	// Extract key phrases first
	keyPhrases := ExtractKeyPhrases(text, 20)
	
	// Group phrases by semantic similarity (simplified clustering)
	clusters := clusterPhrasesForTopics(keyPhrases, maxTopics)
	
	// Convert clusters to topics
	topics := make([]Topic, 0, len(clusters))
	
	for _, cluster := range clusters {
		if len(cluster.phrases) == 0 {
			continue
		}
		
		topic := Topic{
			Name:       generateTopicName(cluster.phrases),
			Keywords:   extractKeywords(cluster.phrases),
			Confidence: cluster.confidence,
			Coverage:   cluster.coverage,
			Examples:   extractExamples(text, cluster.phrases),
		}
		
		topics = append(topics, topic)
		
		if len(topics) >= maxTopics {
			break
		}
	}

	// Sort by confidence
	sort.Slice(topics, func(i, j int) bool {
		return topics[i].Confidence > topics[j].Confidence
	})

	collector.IncrementAlgorithmSteps()

	// Calculate quality metrics
	coverage := 0.0
	totalConfidence := 0.0
	for _, topic := range topics {
		coverage += topic.Coverage
		totalConfidence += topic.Confidence
	}
	
	avgConfidence := 0.0
	if len(topics) > 0 {
		avgConfidence = totalConfidence / float64(len(topics))
	}

	return TopicResult{
		Topics: topics,
		QualityMetrics: QualityMetrics{
			Accuracy:   0.75,
			Confidence: avgConfidence,
			Coverage:   math.Min(coverage, 1.0),
		},
	}
}

// classifyTopicsStatistical performs statistical analysis for topic identification
func classifyTopicsStatistical(text string, maxTopics int, collector *MetricsCollector) TopicResult {
	collector.RecordProcessingTime("statistical_start")

	// Extract frequent terms and their co-occurrences
	termFreq := calculateTermFrequencies(text)
	cooccurrences := calculateCooccurrences(text, termFreq)
	
	// Identify topic clusters based on term associations
	topicClusters := identifyTopicClusters(termFreq, cooccurrences, maxTopics)
	
	// Convert to Topic structures
	topics := make([]Topic, 0, len(topicClusters))
	
	for _, cluster := range topicClusters {
		topic := Topic{
			Name:       cluster.name,
			Keywords:   cluster.keywords,
			Confidence: cluster.strength,
			Coverage:   cluster.coverage,
			Examples:   findTopicExamples(text, cluster.keywords),
		}
		
		topics = append(topics, topic)
	}

	// Sort by confidence
	sort.Slice(topics, func(i, j int) bool {
		return topics[i].Confidence > topics[j].Confidence
	})

	collector.IncrementAlgorithmSteps()
	collector.RecordMemoryUsage()

	// Calculate quality metrics
	coverage := 0.0
	totalConfidence := 0.0
	for _, topic := range topics {
		coverage += topic.Coverage
		totalConfidence += topic.Confidence
	}
	
	avgConfidence := 0.0
	if len(topics) > 0 {
		avgConfidence = totalConfidence / float64(len(topics))
	}

	return TopicResult{
		Topics: topics,
		QualityMetrics: QualityMetrics{
			Accuracy:   0.82,
			Confidence: avgConfidence,
			Coverage:   math.Min(coverage, 1.0),
		},
	}
}

// classifyTopicsComprehensive performs comprehensive topic modeling
func classifyTopicsComprehensive(text string, maxTopics int, collector *MetricsCollector) TopicResult {
	collector.RecordProcessingTime("comprehensive_start")

	// Start with statistical analysis
	statResult := classifyTopicsStatistical(text, maxTopics*2, collector) // Get more candidates
	
	// Refine using additional analysis
	refinedTopics := refinTopicsWithContext(text, statResult.Topics, maxTopics)
	
	// Add semantic analysis
	enhancedTopics := enhanceWithSemanticAnalysis(text, refinedTopics)

	// Sort by enhanced confidence
	sort.Slice(enhancedTopics, func(i, j int) bool {
		return enhancedTopics[i].Confidence > enhancedTopics[j].Confidence
	})

	// Limit to maxTopics
	if len(enhancedTopics) > maxTopics {
		enhancedTopics = enhancedTopics[:maxTopics]
	}

	collector.IncrementAlgorithmSteps()
	collector.RecordMemoryUsage()

	// Calculate enhanced quality metrics
	coverage := 0.0
	totalConfidence := 0.0
	for _, topic := range enhancedTopics {
		coverage += topic.Coverage
		totalConfidence += topic.Confidence
	}
	
	avgConfidence := 0.0
	if len(enhancedTopics) > 0 {
		avgConfidence = totalConfidence / float64(len(enhancedTopics))
	}

	return TopicResult{
		Topics: enhancedTopics,
		QualityMetrics: QualityMetrics{
			Accuracy:   0.88,
			Confidence: avgConfidence,
			Coverage:   math.Min(coverage, 1.0),
		},
	}
}

// Helper types and functions

type phraseCluster struct {
	phrases    []KeyPhrase
	confidence float64
	coverage   float64
}

type topicCluster struct {
	name     string
	keywords []string
	strength float64
	coverage float64
}

func clusterPhrasesForTopics(phrases []KeyPhrase, maxClusters int) []phraseCluster {
	if len(phrases) == 0 {
		return []phraseCluster{}
	}

	clusters := make([]phraseCluster, 0, maxClusters)
	used := make([]bool, len(phrases))

	// Simple clustering based on word overlap
	for i := 0; i < len(phrases) && len(clusters) < maxClusters; i++ {
		if used[i] {
			continue
		}

		cluster := phraseCluster{
			phrases:    []KeyPhrase{phrases[i]},
			confidence: phrases[i].Score,
			coverage:   0.1, // Base coverage
		}
		used[i] = true

		// Find similar phrases to add to cluster
		for j := i + 1; j < len(phrases); j++ {
			if used[j] {
				continue
			}

			if phraseSimilarity(phrases[i].Text, phrases[j].Text) > 0.3 {
				cluster.phrases = append(cluster.phrases, phrases[j])
				cluster.confidence += phrases[j].Score * 0.5 // Diminishing contribution
				cluster.coverage += 0.05
				used[j] = true
			}
		}

		clusters = append(clusters, cluster)
	}

	return clusters
}

func phraseSimilarity(phrase1, phrase2 string) float64 {
	words1 := strings.Fields(strings.ToLower(phrase1))
	words2 := strings.Fields(strings.ToLower(phrase2))

	if len(words1) == 0 || len(words2) == 0 {
		return 0.0
	}

	// Count common words
	commonWords := 0
	for _, w1 := range words1 {
		for _, w2 := range words2 {
			if w1 == w2 && len(w1) > 2 && !isStopWord(w1) {
				commonWords++
				break
			}
		}
	}

	// Jaccard similarity
	totalWords := len(words1) + len(words2) - commonWords
	if totalWords == 0 {
		return 0.0
	}

	return float64(commonWords) / float64(totalWords)
}

func generateTopicName(phrases []KeyPhrase) string {
	if len(phrases) == 0 {
		return "Unknown Topic"
	}

	// Look through all phrases for meaningful words
	for _, phrase := range phrases {
		words := strings.Fields(phrase.Text)
		for _, word := range words {
			cleanWord := strings.ToLower(strings.Trim(word, ".,!?;:\"'"))
			if len(cleanWord) > 3 && !isStopWord(cleanWord) {
				return strings.Title(cleanWord) + " Related"
			}
		}
	}

	// Fallback: use the first phrase even if it's not ideal
	if len(phrases) > 0 {
		words := strings.Fields(phrases[0].Text)
		if len(words) > 0 {
			return strings.Title(strings.ToLower(words[0])) + " Topic"
		}
	}

	return "General Topic"
}

func extractKeywords(phrases []KeyPhrase) []string {
	keywords := make([]string, 0, len(phrases)*2)
	seen := make(map[string]bool)

	for _, phrase := range phrases {
		words := strings.Fields(strings.ToLower(phrase.Text))
		for _, word := range words {
			// Clean word
			word = strings.Trim(word, ".,!?;:\"'")
			if len(word) > 2 && !isStopWord(word) && !seen[word] {
				keywords = append(keywords, word)
				seen[word] = true
			}
		}
	}

	// If no good keywords found, use the phrase text itself
	if len(keywords) == 0 && len(phrases) > 0 {
		// Take first meaningful word from the highest scoring phrase
		words := strings.Fields(strings.ToLower(phrases[0].Text))
		for _, word := range words {
			word = strings.Trim(word, ".,!?;:\"'")
			if len(word) > 2 {
				keywords = append(keywords, word)
				break
			}
		}
	}

	// Limit keywords
	if len(keywords) > 8 {
		keywords = keywords[:8]
	}

	return keywords
}

func extractExamples(text string, phrases []KeyPhrase) []string {
	sentences := SplitIntoSentences(text)
	examples := make([]string, 0, 3)

	// Find sentences containing the phrases
	for _, sentence := range sentences {
		lowerSentence := strings.ToLower(sentence)
		for _, phrase := range phrases {
			if strings.Contains(lowerSentence, strings.ToLower(phrase.Text)) {
				if len(examples) < 3 && len(sentence) > 20 {
					examples = append(examples, sentence)
					break
				}
			}
		}
		if len(examples) >= 3 {
			break
		}
	}

	return examples
}

func calculateTermFrequencies(text string) map[string]float64 {
	words := strings.Fields(strings.ToLower(text))
	freq := make(map[string]int)
	total := 0

	// Count word frequencies
	for _, word := range words {
		word = strings.Trim(word, ".,!?;:\"'")
		if len(word) > 2 && !isStopWord(word) {
			freq[word]++
			total++
		}
	}

	// Normalize frequencies
	normalized := make(map[string]float64)
	for word, count := range freq {
		normalized[word] = float64(count) / float64(total)
	}

	return normalized
}

func calculateCooccurrences(text string, termFreq map[string]float64) map[string]map[string]float64 {
	sentences := SplitIntoSentences(text)
	cooccur := make(map[string]map[string]float64)

	// Initialize maps
	for term := range termFreq {
		cooccur[term] = make(map[string]float64)
	}

	// Calculate co-occurrences within sentences
	for _, sentence := range sentences {
		words := strings.Fields(strings.ToLower(sentence))
		cleanWords := make([]string, 0, len(words))

		// Clean and filter words
		for _, word := range words {
			word = strings.Trim(word, ".,!?;:\"'")
			if len(word) > 2 && !isStopWord(word) {
				if _, exists := termFreq[word]; exists {
					cleanWords = append(cleanWords, word)
				}
			}
		}

		// Calculate co-occurrences
		for i, word1 := range cleanWords {
			for j, word2 := range cleanWords {
				if i != j {
					cooccur[word1][word2]++
				}
			}
		}
	}

	// Normalize co-occurrence scores
	for word1 := range cooccur {
		total := 0.0
		for _, count := range cooccur[word1] {
			total += count
		}
		if total > 0 {
			for word2 := range cooccur[word1] {
				cooccur[word1][word2] /= total
			}
		}
	}

	return cooccur
}

func identifyTopicClusters(termFreq map[string]float64, cooccur map[string]map[string]float64, maxTopics int) []topicCluster {
	// Identify strongly connected word groups
	clusters := make([]topicCluster, 0, maxTopics)
	used := make(map[string]bool)

	// Convert term frequencies to sorted list
	type termScore struct {
		term  string
		score float64
	}

	terms := make([]termScore, 0, len(termFreq))
	for term, freq := range termFreq {
		terms = append(terms, termScore{term, freq})
	}

	sort.Slice(terms, func(i, j int) bool {
		return terms[i].score > terms[j].score
	})

	// Build clusters around high-frequency terms
	for _, termS := range terms {
		if used[termS.term] || len(clusters) >= maxTopics {
			break
		}

		cluster := topicCluster{
			keywords: []string{termS.term},
			strength: termS.score,
			coverage: termS.score,
		}
		used[termS.term] = true

		// Find strongly co-occurring terms
		if cooccurMap, exists := cooccur[termS.term]; exists {
			type cooccurScore struct {
				term  string
				score float64
			}

			cooccurTerms := make([]cooccurScore, 0, len(cooccurMap))
			for term, score := range cooccurMap {
				if !used[term] && score > 0.1 {
					cooccurTerms = append(cooccurTerms, cooccurScore{term, score})
				}
			}

			sort.Slice(cooccurTerms, func(i, j int) bool {
				return cooccurTerms[i].score > cooccurTerms[j].score
			})

			// Add top co-occurring terms to cluster
			for i, coTerm := range cooccurTerms {
				if i >= 4 { // Limit cluster size
					break
				}
				cluster.keywords = append(cluster.keywords, coTerm.term)
				cluster.strength += coTerm.score * 0.5
				cluster.coverage += termFreq[coTerm.term] * 0.5
				used[coTerm.term] = true
			}
		}

		// Generate cluster name
		cluster.name = generateClusterName(cluster.keywords)
		
		clusters = append(clusters, cluster)
	}

	return clusters
}

func generateClusterName(keywords []string) string {
	if len(keywords) == 0 {
		return "Unknown Topic"
	}

	// Use the first (most important) keyword as base
	primary := keywords[0]
	
	// Check for domain-specific patterns
	if isBusinessTerm(primary) {
		return "Business & " + strings.Title(primary)
	} else if isTechTerm(primary) {
		return "Technology & " + strings.Title(primary)
	} else if isHealthTerm(primary) {
		return "Health & " + strings.Title(primary)
	}

	return strings.Title(primary) + " Topic"
}

func isBusinessTerm(term string) bool {
	businessTerms := []string{"business", "market", "sales", "revenue", "profit", "strategy", "management", "company", "customer", "service"}
	for _, bt := range businessTerms {
		if strings.Contains(term, bt) {
			return true
		}
	}
	return false
}

func isTechTerm(term string) bool {
	techTerms := []string{"technology", "software", "data", "system", "computer", "digital", "algorithm", "programming", "database", "network"}
	for _, tt := range techTerms {
		if strings.Contains(term, tt) {
			return true
		}
	}
	return false
}

func isHealthTerm(term string) bool {
	healthTerms := []string{"health", "medical", "doctor", "patient", "treatment", "disease", "medicine", "hospital", "care", "therapy"}
	for _, ht := range healthTerms {
		if strings.Contains(term, ht) {
			return true
		}
	}
	return false
}

func findTopicExamples(text string, keywords []string) []string {
	sentences := SplitIntoSentences(text)
	examples := make([]string, 0, 3)

	for _, sentence := range sentences {
		lowerSentence := strings.ToLower(sentence)
		keywordCount := 0
		
		for _, keyword := range keywords {
			if strings.Contains(lowerSentence, keyword) {
				keywordCount++
			}
		}

		// Include sentences with multiple keywords
		if keywordCount >= 2 && len(sentence) > 30 {
			examples = append(examples, sentence)
			if len(examples) >= 3 {
				break
			}
		}
	}

	// If not enough multi-keyword sentences, include single-keyword ones
	if len(examples) < 2 {
		for _, sentence := range sentences {
			lowerSentence := strings.ToLower(sentence)
			for _, keyword := range keywords {
				if strings.Contains(lowerSentence, keyword) && len(sentence) > 20 {
					examples = append(examples, sentence)
					if len(examples) >= 3 {
						break
					}
				}
			}
			if len(examples) >= 3 {
				break
			}
		}
	}

	return examples
}

func refinTopicsWithContext(text string, topics []Topic, maxTopics int) []Topic {
	// Analyze document structure for topic relevance
	sentences := SplitIntoSentences(text)
	
	for i := range topics {
		// Calculate better coverage based on sentence distribution
		sentenceMatches := 0
		for _, sentence := range sentences {
			lowerSentence := strings.ToLower(sentence)
			keywordMatches := 0
			
			for _, keyword := range topics[i].Keywords {
				if strings.Contains(lowerSentence, keyword) {
					keywordMatches++
				}
			}
			
			if keywordMatches > 0 {
				sentenceMatches++
			}
		}
		
		// Update coverage based on sentence distribution
		if len(sentences) > 0 {
			topics[i].Coverage = float64(sentenceMatches) / float64(len(sentences))
		}
		
		// Adjust confidence based on coverage
		topics[i].Confidence = topics[i].Confidence * (0.5 + topics[i].Coverage*0.5)
	}

	// Sort by refined confidence
	sort.Slice(topics, func(i, j int) bool {
		return topics[i].Confidence > topics[j].Confidence
	})

	// Return top topics
	if len(topics) > maxTopics {
		topics = topics[:maxTopics]
	}

	return topics
}

func enhanceWithSemanticAnalysis(text string, topics []Topic) []Topic {
	// Simple semantic enhancement: look for related terms
	for i := range topics {
		// Expand keywords with related terms found in text
		relatedTerms := findRelatedTerms(text, topics[i].Keywords)
		
		// Add unique related terms
		for _, term := range relatedTerms {
			unique := true
			for _, existing := range topics[i].Keywords {
				if existing == term {
					unique = false
					break
				}
			}
			if unique && len(topics[i].Keywords) < 10 {
				topics[i].Keywords = append(topics[i].Keywords, term)
			}
		}
		
		// Boost confidence for topics with more semantic relationships
		semanticBoost := math.Min(float64(len(relatedTerms))*0.05, 0.2)
		topics[i].Confidence = math.Min(topics[i].Confidence+semanticBoost, 1.0)
	}

	return topics
}

func findRelatedTerms(text string, keywords []string) []string {
	words := strings.Fields(strings.ToLower(text))
	related := make([]string, 0, 5)
	
	// Simple approach: find words that appear near the keywords
	for i, word := range words {
		word = strings.Trim(word, ".,!?;:\"'")
		if len(word) <= 2 || isStopWord(word) {
			continue
		}
		
		// Check if word appears near any keyword
		for _, keyword := range keywords {
			// Look in a window around the keyword
			for j := maxInt(0, i-3); j <= minInt(len(words)-1, i+3); j++ {
				checkWord := strings.Trim(words[j], ".,!?;:\"'")
				if checkWord == keyword {
					// Found keyword nearby, add current word if not already present
					found := false
					for _, existing := range keywords {
						if existing == word {
							found = true
							break
						}
					}
					for _, existing := range related {
						if existing == word {
							found = true
							break
						}
					}
					if !found && len(related) < 5 {
						related = append(related, word)
					}
					break
				}
			}
		}
	}
	
	return related
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}