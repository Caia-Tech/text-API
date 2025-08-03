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
	"strings"
	"time"
)

// ExtractSentiment analyzes sentiment with adaptive algorithm selection
// accuracy: 0.7 = lexicon-based, 0.85 = rule-based, 0.95 = contextual analysis
func ExtractSentiment(text string, accuracy float64) SentimentResult {
	// Start metrics collection
	collector := StartMetricsCollection()
	startTime := time.Now()

	// Validate input
	if text == "" {
		return SentimentResult{
			OverallSentiment: Sentiment{
				Polarity:   0.0,
				Magnitude:  0.0,
				Label:      "neutral",
				Confidence: 0.0,
			},
			Method:         "none",
			ProcessingTime: time.Since(startTime),
			QualityMetrics: QualityMetrics{
				Accuracy:   0,
				Confidence: 0,
				Coverage:   0,
			},
		}
	}

	// Validate accuracy
	if accuracy < 0.7 {
		accuracy = 0.7
	} else if accuracy > 0.95 {
		accuracy = 0.95
	}

	var result SentimentResult

	// Choose analysis method based on accuracy requirement
	if accuracy <= 0.75 {
		// Fast lexicon-based analysis
		result = extractSentimentLexicon(text, collector)
		result.Method = "lexicon-based"
	} else if accuracy <= 0.85 {
		// Rule-based analysis with patterns
		result = extractSentimentRuleBased(text, collector)
		result.Method = "rule-based"
	} else {
		// Contextual analysis with advanced features
		result = extractSentimentContextual(text, collector)
		result.Method = "contextual-analysis"
	}

	result.ProcessingTime = time.Since(startTime)

	// Record metrics
	metrics := collector.GetMetrics()
	params := map[string]interface{}{
		"accuracy":    accuracy,
		"text_length": len(text),
		"method":      result.Method,
	}

	RecordFunctionCall("ExtractSentiment", params, metrics, &result.QualityMetrics)

	return result
}

// extractSentimentLexicon performs fast lexicon-based sentiment analysis
func extractSentimentLexicon(text string, collector *MetricsCollector) SentimentResult {
	collector.RecordProcessingTime("lexicon_analysis_start")

	// Analyze overall sentiment
	overallSentiment := analyzeSentimentLexicon(text)

	// Analyze by sentences
	sentences := SplitIntoSentences(text)
	sentenceSentiments := make([]SentenceSentiment, len(sentences))

	for i, sentence := range sentences {
		sentiment := analyzeSentimentLexicon(sentence)
		sentenceSentiments[i] = SentenceSentiment{
			Text:      sentence,
			Sentiment: sentiment,
			Position: RLPosition{
				Start: 0, // Simplified for this implementation
				End:   len(sentence),
			},
		}
	}

	// Basic emotion profile
	emotions := extractBasicEmotions(text)

	collector.IncrementAlgorithmSteps()

	return SentimentResult{
		OverallSentiment:   overallSentiment,
		SentenceSentiments: sentenceSentiments,
		EmotionProfile:     emotions,
		QualityMetrics: QualityMetrics{
			Accuracy:   0.75,
			Confidence: overallSentiment.Confidence,
			Coverage:   1.0,
		},
	}
}

// extractSentimentRuleBased performs rule-based sentiment analysis
func extractSentimentRuleBased(text string, collector *MetricsCollector) SentimentResult {
	collector.RecordProcessingTime("rule_based_analysis_start")

	// Start with lexicon analysis
	lexiconResult := extractSentimentLexicon(text, collector)

	// Apply rules to adjust sentiment
	adjustedSentiment := applySentimentRules(text, lexiconResult.OverallSentiment)

	// Enhanced sentence analysis
	sentences := SplitIntoSentences(text)
	sentenceSentiments := make([]SentenceSentiment, len(sentences))

	for i, sentence := range sentences {
		baseSentiment := analyzeSentimentLexicon(sentence)
		adjustedSentiment := applySentimentRules(sentence, baseSentiment)
		sentenceSentiments[i] = SentenceSentiment{
			Text:      sentence,
			Sentiment: adjustedSentiment,
			Position: RLPosition{
				Start: 0,
				End:   len(sentence),
			},
		}
	}

	// Enhanced emotion profile
	emotions := extractAdvancedEmotions(text)

	collector.IncrementAlgorithmSteps()
	collector.RecordMemoryUsage()

	return SentimentResult{
		OverallSentiment:   adjustedSentiment,
		SentenceSentiments: sentenceSentiments,
		EmotionProfile:     emotions,
		QualityMetrics: QualityMetrics{
			Accuracy:   0.82,
			Confidence: adjustedSentiment.Confidence,
			Coverage:   1.0,
		},
	}
}

// extractSentimentContextual performs contextual sentiment analysis
func extractSentimentContextual(text string, collector *MetricsCollector) SentimentResult {
	collector.RecordProcessingTime("contextual_analysis_start")

	// Start with rule-based analysis
	ruleBasedResult := extractSentimentRuleBased(text, collector)

	// Apply contextual adjustments
	contextualSentiment := applyContextualAnalysis(text, ruleBasedResult.OverallSentiment)

	// Advanced sentence analysis with context
	sentences := SplitIntoSentences(text)
	sentenceSentiments := make([]SentenceSentiment, len(sentences))

	for i, sentence := range sentences {
		// Consider surrounding context
		context := ""
		if i > 0 {
			context += sentences[i-1] + " "
		}
		context += sentence
		if i < len(sentences)-1 {
			context += " " + sentences[i+1]
		}

		baseSentiment := analyzeSentimentLexicon(sentence)
		ruleSentiment := applySentimentRules(sentence, baseSentiment)
		contextSentiment := applyContextualAnalysis(context, ruleSentiment)

		sentenceSentiments[i] = SentenceSentiment{
			Text:      sentence,
			Sentiment: contextSentiment,
			Position: RLPosition{
				Start: 0,
				End:   len(sentence),
			},
		}
	}

	// Comprehensive emotion profile
	emotions := extractComprehensiveEmotions(text)

	collector.IncrementAlgorithmSteps()
	collector.RecordMemoryUsage()

	return SentimentResult{
		OverallSentiment:   contextualSentiment,
		SentenceSentiments: sentenceSentiments,
		EmotionProfile:     emotions,
		QualityMetrics: QualityMetrics{
			Accuracy:   0.90,
			Confidence: contextualSentiment.Confidence,
			Coverage:   1.0,
		},
	}
}

// Core sentiment analysis functions

func analyzeSentimentLexicon(text string) Sentiment {
	words := strings.Fields(strings.ToLower(text))
	
	positiveScore := 0.0
	negativeScore := 0.0
	totalWords := 0

	for _, word := range words {
		word = strings.Trim(word, ".,!?;:\"'")
		if len(word) < 2 {
			continue
		}

		if score := getPositiveSentimentScore(word); score > 0 {
			positiveScore += score
			totalWords++
		} else if score := getNegativeSentimentScore(word); score > 0 {
			negativeScore += score
			totalWords++
		}
	}

	// Calculate final sentiment
	if totalWords == 0 {
		return Sentiment{
			Polarity:   0.0,
			Magnitude:  0.0,
			Label:      "neutral",
			Confidence: 0.5,
		}
	}

	netScore := (positiveScore - negativeScore) / float64(totalWords)
	magnitude := (positiveScore + negativeScore) / float64(totalWords)

	// Determine label
	label := "neutral"
	if netScore > 0.2 {
		label = "positive"
	} else if netScore < -0.2 {
		label = "negative"
	}

	// Calculate confidence based on magnitude
	confidence := magnitude
	if confidence > 0.9 {
		confidence = 0.9
	} else if confidence < 0.3 {
		confidence = 0.3
	}

	return Sentiment{
		Polarity:   netScore,
		Magnitude:  magnitude,
		Label:      label,
		Confidence: confidence,
	}
}

func applySentimentRules(text string, baseSentiment Sentiment) Sentiment {
	adjusted := baseSentiment
	lowerText := strings.ToLower(text)

	// Negation rules
	negationWords := []string{"not", "no", "never", "neither", "nothing", "nowhere", "nobody", "none", "n't"}
	for _, negation := range negationWords {
		if strings.Contains(lowerText, negation) {
			// Flip and reduce sentiment when negation is found
			adjusted.Polarity = -adjusted.Polarity * 0.8
			break
		}
	}

	// Intensifier rules
	intensifiers := []string{"very", "extremely", "incredibly", "absolutely", "completely", "totally"}
	for _, intensifier := range intensifiers {
		if strings.Contains(lowerText, intensifier) {
			// Amplify sentiment
			adjusted.Polarity = adjusted.Polarity * 1.3
			adjusted.Magnitude = adjusted.Magnitude * 1.2
			break
		}
	}

	// Diminisher rules
	diminishers := []string{"slightly", "somewhat", "fairly", "rather", "quite", "a bit"}
	for _, diminisher := range diminishers {
		if strings.Contains(lowerText, diminisher) {
			// Reduce sentiment
			adjusted.Polarity = adjusted.Polarity * 0.7
			adjusted.Magnitude = adjusted.Magnitude * 0.8
			break
		}
	}

	// Question rules
	if strings.Contains(text, "?") {
		// Questions are generally more neutral
		adjusted.Polarity = adjusted.Polarity * 0.6
		adjusted.Magnitude = adjusted.Magnitude * 0.8
	}

	// Exclamation rules
	exclamationCount := strings.Count(text, "!")
	if exclamationCount > 0 {
		// Exclamations amplify sentiment
		multiplier := 1.0 + float64(exclamationCount)*0.1
		adjusted.Polarity = adjusted.Polarity * multiplier
		adjusted.Magnitude = adjusted.Magnitude * multiplier
	}

	// Update label and confidence
	if adjusted.Polarity > 0.2 {
		adjusted.Label = "positive"
	} else if adjusted.Polarity < -0.2 {
		adjusted.Label = "negative"
	} else {
		adjusted.Label = "neutral"
	}

	// Clamp values
	if adjusted.Polarity > 1.0 {
		adjusted.Polarity = 1.0
	} else if adjusted.Polarity < -1.0 {
		adjusted.Polarity = -1.0
	}
	
	if adjusted.Magnitude > 1.0 {
		adjusted.Magnitude = 1.0
	}

	adjusted.Confidence = adjusted.Magnitude
	if adjusted.Confidence > 0.9 {
		adjusted.Confidence = 0.9
	} else if adjusted.Confidence < 0.3 {
		adjusted.Confidence = 0.3
	}

	return adjusted
}

func applyContextualAnalysis(text string, baseSentiment Sentiment) Sentiment {
	adjusted := baseSentiment

	// Analyze comparative patterns
	if strings.Contains(strings.ToLower(text), "better than") ||
		strings.Contains(strings.ToLower(text), "worse than") {
		// Comparative statements are more confident
		adjusted.Confidence = math.Min(adjusted.Confidence*1.2, 0.95)
	}

	// Analyze conditional patterns
	conditionals := []string{"if", "unless", "provided", "assuming"}
	for _, conditional := range conditionals {
		if strings.Contains(strings.ToLower(text), conditional) {
			// Conditional statements are less certain
			adjusted.Polarity = adjusted.Polarity * 0.8
			adjusted.Confidence = adjusted.Confidence * 0.9
			break
		}
	}

	// Analyze temporal patterns
	past := []string{"was", "were", "had", "did"}
	future := []string{"will", "shall", "going to", "would"}
	
	hasPast := false
	hasFuture := false
	lowerText := strings.ToLower(text)
	
	for _, p := range past {
		if strings.Contains(lowerText, p) {
			hasPast = true
			break
		}
	}
	for _, f := range future {
		if strings.Contains(lowerText, f) {
			hasFuture = true
			break
		}
	}

	if hasPast {
		// Past events are more certain
		adjusted.Confidence = math.Min(adjusted.Confidence*1.1, 0.9)
	}
	if hasFuture {
		// Future events are less certain
		adjusted.Confidence = adjusted.Confidence * 0.9
	}

	return adjusted
}

// Emotion analysis functions

func extractBasicEmotions(text string) EmotionProfile {
	words := strings.Fields(strings.ToLower(text))
	
	emotions := EmotionProfile{}
	totalWords := 0

	for _, word := range words {
		word = strings.Trim(word, ".,!?;:\"'")
		if len(word) < 2 {
			continue
		}

		// Simple emotion mapping
		if isJoyWord(word) {
			emotions.Joy += 1.0
			totalWords++
		}
		if isAngerWord(word) {
			emotions.Anger += 1.0
			totalWords++
		}
		if isFearWord(word) {
			emotions.Fear += 1.0
			totalWords++
		}
		if isSadnessWord(word) {
			emotions.Sadness += 1.0
			totalWords++
		}
		if isSurpriseWord(word) {
			emotions.Surprise += 1.0
			totalWords++
		}
		if isTrustWord(word) {
			emotions.Trust += 1.0
			totalWords++
		}
	}

	// Normalize emotions
	if totalWords > 0 {
		emotions.Joy /= float64(totalWords)
		emotions.Anger /= float64(totalWords)
		emotions.Fear /= float64(totalWords)
		emotions.Sadness /= float64(totalWords)
		emotions.Surprise /= float64(totalWords)
		emotions.Trust /= float64(totalWords)
	}

	return emotions
}

func extractAdvancedEmotions(text string) EmotionProfile {
	// Start with basic emotions
	emotions := extractBasicEmotions(text)

	// Apply rules for emotional intensifiers
	lowerText := strings.ToLower(text)
	
	// Boost emotions with intensifiers
	if strings.Contains(lowerText, "very") || strings.Contains(lowerText, "extremely") {
		emotions.Joy *= 1.2
		emotions.Anger *= 1.2
		emotions.Fear *= 1.2
		emotions.Sadness *= 1.2
	}

	// Check for emotional punctuation
	if strings.Contains(text, "!") {
		emotions.Joy *= 1.1
		emotions.Anger *= 1.1
		emotions.Surprise *= 1.3
	}

	if strings.Contains(text, "?") {
		emotions.Surprise *= 1.2
		emotions.Fear *= 1.1
	}

	// Clamp values
	emotions.Joy = math.Min(emotions.Joy, 1.0)
	emotions.Anger = math.Min(emotions.Anger, 1.0)
	emotions.Fear = math.Min(emotions.Fear, 1.0)
	emotions.Sadness = math.Min(emotions.Sadness, 1.0)
	emotions.Surprise = math.Min(emotions.Surprise, 1.0)
	emotions.Trust = math.Min(emotions.Trust, 1.0)

	return emotions
}

func extractComprehensiveEmotions(text string) EmotionProfile {
	// Start with advanced emotions
	emotions := extractAdvancedEmotions(text)

	// Contextual emotion adjustments
	sentences := SplitIntoSentences(text)
	
	// Look for emotional progressions
	if len(sentences) > 1 {
		// Analyze emotional arc
		firstHalf := strings.Join(sentences[:len(sentences)/2], " ")
		secondHalf := strings.Join(sentences[len(sentences)/2:], " ")
		
		firstEmotions := extractBasicEmotions(firstHalf)
		secondEmotions := extractBasicEmotions(secondHalf)
		
		// If emotions change significantly, boost surprise
		emotionalChange := math.Abs(firstEmotions.Joy-secondEmotions.Joy) +
			math.Abs(firstEmotions.Sadness-secondEmotions.Sadness)
		
		if emotionalChange > 0.3 {
			emotions.Surprise = math.Min(emotions.Surprise+0.2, 1.0)
		}
	}

	return emotions
}

// Sentiment lexicon functions

func getPositiveSentimentScore(word string) float64 {
	positiveWords := map[string]float64{
		"amazing": 0.9, "awesome": 0.8, "beautiful": 0.7, "best": 0.8, "brilliant": 0.8,
		"excellent": 0.9, "fantastic": 0.8, "good": 0.6, "great": 0.7, "happy": 0.7,
		"incredible": 0.8, "love": 0.8, "perfect": 0.9, "pleased": 0.6, "wonderful": 0.8,
		"outstanding": 0.8, "superb": 0.8, "magnificent": 0.8, "marvelous": 0.7, "terrific": 0.7,
		"delighted": 0.7, "thrilled": 0.8, "ecstatic": 0.9, "joyful": 0.7, "elated": 0.8,
		"satisfied": 0.6, "content": 0.5, "pleasant": 0.6, "nice": 0.5, "fine": 0.4,
	}
	
	if score, exists := positiveWords[word]; exists {
		return score
	}
	return 0.0
}

func getNegativeSentimentScore(word string) float64 {
	negativeWords := map[string]float64{
		"awful": 0.9, "bad": 0.6, "disgusting": 0.8, "hate": 0.8, "horrible": 0.8,
		"terrible": 0.8, "worst": 0.9, "angry": 0.7, "sad": 0.6, "disappointed": 0.7,
		"frustrated": 0.7, "annoyed": 0.6, "upset": 0.6, "furious": 0.9, "devastated": 0.9,
		"miserable": 0.8, "depressed": 0.8, "hopeless": 0.8, "pathetic": 0.7, "useless": 0.7,
		"ridiculous": 0.6, "stupid": 0.7, "crazy": 0.5, "insane": 0.6, "wrong": 0.5,
		"problem": 0.5, "issues": 0.5, "failed": 0.7, "broken": 0.6, "difficult": 0.4,
	}
	
	if score, exists := negativeWords[word]; exists {
		return score
	}
	return 0.0
}

// Emotion classification functions

func isJoyWord(word string) bool {
	joyWords := []string{"happy", "joy", "cheerful", "delighted", "ecstatic", "elated", "glad", "pleased", "merry", "blissful"}
	for _, jw := range joyWords {
		if word == jw {
			return true
		}
	}
	return false
}

func isAngerWord(word string) bool {
	angerWords := []string{"angry", "furious", "mad", "rage", "irritated", "annoyed", "frustrated", "outraged", "livid", "irate"}
	for _, aw := range angerWords {
		if word == aw {
			return true
		}
	}
	return false
}

func isFearWord(word string) bool {
	fearWords := []string{"afraid", "scared", "terrified", "frightened", "worried", "anxious", "nervous", "panic", "dread", "horror"}
	for _, fw := range fearWords {
		if word == fw {
			return true
		}
	}
	return false
}

func isSadnessWord(word string) bool {
	sadnessWords := []string{"sad", "depressed", "miserable", "disappointed", "heartbroken", "grief", "sorrow", "melancholy", "despair", "gloom"}
	for _, sw := range sadnessWords {
		if word == sw {
			return true
		}
	}
	return false
}

func isSurpriseWord(word string) bool {
	surpriseWords := []string{"surprised", "amazed", "astonished", "shocked", "stunned", "bewildered", "confused", "puzzled", "wonder", "awe"}
	for _, sw := range surpriseWords {
		if word == sw {
			return true
		}
	}
	return false
}

func isTrustWord(word string) bool {
	trustWords := []string{"trust", "confident", "secure", "reliable", "faithful", "loyal", "honest", "sincere", "genuine", "dependable"}
	for _, tw := range trustWords {
		if word == tw {
			return true
		}
	}
	return false
}