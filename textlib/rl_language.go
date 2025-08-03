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
	"unicode"
)

// Language detection common words and patterns
var languagePatterns = map[string]languageProfile{
	"en": {
		CommonWords: []string{"the", "is", "and", "of", "to", "in", "a", "that", "it", "for", "was", "with", "as", "on", "by", "at", "from", "have", "be", "are"},
		CharPatterns: map[string]float64{"th": 0.027, "he": 0.023, "in": 0.020, "er": 0.019, "an": 0.018},
		TypicalChars: "aeiourstlnmcpdhgbfwyvkxjqz",
	},
	"es": {
		CommonWords: []string{"el", "la", "de", "que", "y", "a", "en", "un", "ser", "se", "no", "haber", "por", "con", "su", "para", "como", "estar", "tener", "le"},
		CharPatterns: map[string]float64{"de": 0.025, "la": 0.022, "en": 0.020, "el": 0.019, "es": 0.018},
		TypicalChars: "aeiourstlnmcdpgbfvhñjqyzxkw",
	},
	"fr": {
		CommonWords: []string{"le", "de", "un", "être", "et", "à", "il", "avoir", "ne", "je", "son", "que", "se", "qui", "ce", "dans", "en", "du", "elle", "au"},
		CharPatterns: map[string]float64{"de": 0.024, "le": 0.021, "en": 0.019, "re": 0.018, "nt": 0.017},
		TypicalChars: "aeiourstlnmcdpgbfvhçèéêàâîôûùœæ",
	},
	"de": {
		CommonWords: []string{"der", "die", "und", "in", "den", "von", "zu", "das", "mit", "sich", "auf", "für", "ist", "im", "dem", "nicht", "ein", "eine", "als", "auch"},
		CharPatterns: map[string]float64{"en": 0.026, "er": 0.024, "ch": 0.022, "de": 0.020, "ei": 0.018},
		TypicalChars: "aeiourstlnmdhgbfwkzpvjßäöü",
	},
	"it": {
		CommonWords: []string{"di", "il", "che", "è", "e", "la", "a", "un", "in", "non", "si", "per", "con", "su", "sono", "da", "come", "una", "dei", "delle"},
		CharPatterns: map[string]float64{"di": 0.023, "in": 0.021, "la": 0.019, "re": 0.018, "to": 0.017},
		TypicalChars: "aeiourstlnmcdpgbfvhzàèéìòù",
	},
	"pt": {
		CommonWords: []string{"de", "o", "que", "a", "e", "do", "da", "em", "um", "para", "com", "não", "uma", "os", "no", "se", "na", "por", "mais", "as"},
		CharPatterns: map[string]float64{"de": 0.025, "ar": 0.020, "os": 0.019, "do": 0.018, "ra": 0.017},
		TypicalChars: "aeiourstlnmcdpgbfvhãõçàáâêéíóôúü",
	},
	"nl": {
		CommonWords: []string{"de", "het", "een", "van", "en", "in", "is", "dat", "op", "te", "zijn", "voor", "met", "die", "niet", "aan", "er", "om", "ook", "als"},
		CharPatterns: map[string]float64{"en": 0.028, "de": 0.024, "et": 0.020, "an": 0.019, "er": 0.018},
		TypicalChars: "aeiourstlnmdhgbfwkvpjzij",
	},
	"ru": {
		CommonWords: []string{"и", "в", "не", "на", "я", "что", "он", "с", "как", "это", "по", "но", "все", "она", "так", "его", "от", "за", "то", "мы"},
		CharPatterns: map[string]float64{"ст": 0.020, "на": 0.019, "ов": 0.018, "то": 0.017, "ен": 0.016},
		TypicalChars: "аеиоуыэюяёбвгджзйклмнпрстфхцчшщъь",
	},
	"ja": {
		CommonWords: []string{"の", "は", "に", "を", "が", "と", "で", "て", "た", "です", "も", "な", "い", "か", "ます", "から", "こと", "ある", "する", "れる"},
		CharPatterns: map[string]float64{},
		TypicalChars: "あいうえおかきくけこさしすせそたちつてとなにぬねのはひふへほまみむめもやゆよらりるれろわをん",
	},
	"zh": {
		CommonWords: []string{"的", "是", "在", "一", "有", "个", "我", "不", "这", "了", "他", "们", "人", "来", "到", "大", "和", "子", "说", "要"},
		CharPatterns: map[string]float64{},
		TypicalChars: "",
	},
}

type languageProfile struct {
	CommonWords  []string
	CharPatterns map[string]float64
	TypicalChars string
}

// DetectLanguage detects language with configurable confidence threshold
// confidence: 0.5 = fast heuristics, 0.95 = comprehensive analysis
func DetectLanguage(text string, confidence float64) LanguageResult {
	// Start metrics collection
	collector := StartMetricsCollection()
	startTime := time.Now()
	
	// Validate input
	if text == "" {
		return LanguageResult{
			Language:       "unknown",
			Confidence:     0.0,
			Alternatives:   []LanguageCandidate{},
			Method:         "empty",
			ProcessingTime: time.Since(startTime),
		}
	}
	
	// Validate confidence
	if confidence < 0.5 {
		confidence = 0.5
	} else if confidence > 0.95 {
		confidence = 0.95
	}
	
	var result LanguageResult
	
	// Choose detection method based on confidence requirement
	if confidence <= 0.6 {
		// Fast character-based detection
		result = detectLanguageFast(text, collector)
		result.Method = "character-frequency"
	} else if confidence <= 0.8 {
		// Statistical n-gram analysis
		result = detectLanguageStatistical(text, collector)
		result.Method = "statistical-ngram"
	} else {
		// Comprehensive analysis
		result = detectLanguageComprehensive(text, collector)
		result.Method = "comprehensive-analysis"
	}
	
	// Special handling for non-Latin scripts - they should have higher confidence
	if result.Language == "zh" || result.Language == "ja" || result.Language == "ru" {
		if result.Confidence < 0.8 {
			result.Confidence = 0.8
		}
	} else {
		// Adjust confidence based on text length for Latin scripts
		textLength := len([]rune(text))
		if textLength < 20 {
			result.Confidence *= 0.7
		} else if textLength < 50 {
			result.Confidence *= 0.85
		}
	}
	
	// Ensure confidence doesn't exceed requested threshold
	if result.Confidence > confidence {
		result.Confidence = confidence
	}
	
	result.ProcessingTime = time.Since(startTime)
	
	// Record metrics
	metrics := collector.GetMetrics()
	params := map[string]interface{}{
		"confidence":  confidence,
		"text_length": len(text),
		"method":      result.Method,
	}
	
	quality := &QualityMetrics{
		Accuracy:   result.Confidence,
		Confidence: result.Confidence,
		Coverage:   1.0,
	}
	
	RecordFunctionCall("DetectLanguage", params, metrics, quality)
	
	return result
}

// detectLanguageFast performs fast character-based language detection
func detectLanguageFast(text string, collector *MetricsCollector) LanguageResult {
	collector.RecordProcessingTime("fast_detection_start")
	
	// Count character frequencies
	charFreq := make(map[rune]int)
	totalChars := 0
	
	// Check for specific scripts first
	var hasLatin, hasChinese, hasJapanese, hasCyrillic bool
	
	for _, r := range text {
		if unicode.IsLetter(r) {
			charFreq[r]++
			totalChars++
			
			// Check script type
			if r >= 0x4E00 && r <= 0x9FFF {
				hasChinese = true
			} else if (r >= 0x3040 && r <= 0x309F) || (r >= 0x30A0 && r <= 0x30FF) {
				hasJapanese = true
			} else if r >= 0x0400 && r <= 0x04FF {
				hasCyrillic = true
			} else if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') {
				hasLatin = true
			}
		}
	}
	
	// Return early for non-Latin scripts
	if hasChinese && !hasJapanese {
		return LanguageResult{
			Language:     "zh",
			Confidence:   0.9,
			Alternatives: []LanguageCandidate{},
		}
	}
	
	if hasJapanese {
		return LanguageResult{
			Language:     "ja",
			Confidence:   0.9,
			Alternatives: []LanguageCandidate{},
		}
	}
	
	if hasCyrillic && !hasLatin {
		return LanguageResult{
			Language:     "ru",
			Confidence:   0.85,
			Alternatives: []LanguageCandidate{},
		}
	}
	
	// Score each language based on character usage
	scores := make(map[string]float64)
	
	// For Latin-based languages, use character frequency
	for lang, profile := range languagePatterns {
		if lang == "zh" || lang == "ja" || lang == "ru" {
			continue
		}
		
		score := 0.0
		charCount := 0
		
		// Check typical characters
		for _, char := range profile.TypicalChars {
			if count, exists := charFreq[char]; exists {
				score += float64(count)
				charCount++
			}
		}
		
		// Normalize by number of typical chars found
		if charCount > 0 && totalChars > 0 {
			score = score / float64(totalChars)
			// Boost for languages with special characters
			if lang == "es" && (charFreq['ñ'] > 0 || charFreq['¿'] > 0 || charFreq['¡'] > 0) {
				score *= 1.5
			} else if lang == "fr" && (charFreq['é'] > 0 || charFreq['è'] > 0 || charFreq['à'] > 0) {
				score *= 1.5
			} else if lang == "de" && (charFreq['ä'] > 0 || charFreq['ö'] > 0 || charFreq['ü'] > 0 || charFreq['ß'] > 0) {
				score *= 1.5
			} else if lang == "pt" && (charFreq['ã'] > 0 || charFreq['õ'] > 0 || charFreq['ç'] > 0) {
				score *= 1.5
			}
		}
		
		scores[lang] = score
	}
	
	// Find best match
	bestLang, bestScore := findBestLanguage(scores)
	alternatives := findAlternatives(scores, bestLang, 3)
	
	collector.IncrementAlgorithmSteps()
	
	// Calculate confidence based on score and text length
	confidence := bestScore * 2.0 // Scale up
	if len(text) < 20 {
		confidence *= 0.7
	} else if len(text) < 50 {
		confidence *= 0.85
	}
	
	if confidence > 0.8 {
		confidence = 0.8
	}
	if confidence < 0.3 {
		confidence = 0.3
	}
	
	return LanguageResult{
		Language:     bestLang,
		Confidence:   confidence,
		Alternatives: alternatives,
	}
}

// detectLanguageStatistical performs statistical n-gram analysis
func detectLanguageStatistical(text string, collector *MetricsCollector) LanguageResult {
	collector.RecordProcessingTime("statistical_detection_start")
	
	// First check for non-Latin scripts - fast detection is more accurate for these
	fastResult := detectLanguageFast(text, collector)
	if fastResult.Language == "zh" || fastResult.Language == "ja" || fastResult.Language == "ru" {
		return fastResult
	}
	
	// Prepare text
	cleanText := strings.ToLower(text)
	words := strings.Fields(cleanText)
	
	// Count word occurrences and bigrams
	wordFreq := make(map[string]int)
	bigramFreq := make(map[string]int)
	totalBigrams := 0
	
	for _, word := range words {
		// Clean word
		word = strings.ToLower(strings.Trim(word, ".,!?;:\"'"))
		if len(word) > 0 {
			wordFreq[word]++
		}
		
		// Extract bigrams
		for i := 0; i < len(word)-1; i++ {
			bigram := word[i:i+2]
			bigramFreq[bigram]++
			totalBigrams++
		}
	}
	
	// Score languages based on common words and patterns
	scores := make(map[string]float64)
	
	for lang, profile := range languagePatterns {
		score := 0.0
		wordMatches := 0
		
		// Check common words (weighted by frequency)
		for i, commonWord := range profile.CommonWords {
			if count, exists := wordFreq[commonWord]; exists {
				wordMatches++
				// Give higher weight to more common words (earlier in list)
				weight := 1.0 / float64(i+1)
				score += float64(count) * weight * 0.2
			}
		}
		
		// Check character patterns (bigrams)
		if totalBigrams > 0 {
			for pattern, expectedFreq := range profile.CharPatterns {
				if count, exists := bigramFreq[pattern]; exists {
					actualFreq := float64(count) / float64(totalBigrams)
					// Score based on how close the frequency is to expected
					diff := math.Abs(actualFreq - expectedFreq)
					if diff < 0.01 {
						score += 0.1
					} else if diff < 0.02 {
						score += 0.05
					}
				}
			}
		}
		
		// Bonus for multiple word matches
		if wordMatches >= 3 {
			score *= (1.0 + float64(wordMatches) * 0.1)
		}
		
		scores[lang] = score
	}
	
	// Find best match
	bestLang, _ := findBestLanguage(scores)
	alternatives := findAlternatives(scores, bestLang, 3)
	
	// Calculate confidence based on score distribution
	confidence := calculateConfidence(scores, bestLang)
	
	collector.IncrementAlgorithmSteps()
	collector.RecordMemoryUsage()
	
	return LanguageResult{
		Language:     bestLang,
		Confidence:   confidence,
		Alternatives: alternatives,
	}
}

// detectLanguageComprehensive performs comprehensive language analysis
func detectLanguageComprehensive(text string, collector *MetricsCollector) LanguageResult {
	collector.RecordProcessingTime("comprehensive_detection_start")
	
	// Get results from both methods
	fastResult := detectLanguageFast(text, collector)
	statResult := detectLanguageStatistical(text, collector)
	
	// Additional analysis: sentence structure patterns
	sentences := SplitIntoSentences(text)
	structureScores := analyzeStructurePatterns(sentences)
	
	// Combine results
	combinedScores := make(map[string]float64)
	
	// Weight fast detection (character-based)
	combinedScores[fastResult.Language] += fastResult.Confidence * 0.3
	
	// Weight statistical detection (word/n-gram based)
	combinedScores[statResult.Language] += statResult.Confidence * 0.5
	
	// Weight structure analysis
	for lang, score := range structureScores {
		combinedScores[lang] += score * 0.2
	}
	
	// Also consider alternatives
	for _, alt := range fastResult.Alternatives {
		combinedScores[alt.Language] += alt.Confidence * 0.1
	}
	for _, alt := range statResult.Alternatives {
		combinedScores[alt.Language] += alt.Confidence * 0.15
	}
	
	// Find final best match
	bestLang, bestScore := findBestLanguage(combinedScores)
	alternatives := findAlternatives(combinedScores, bestLang, 4)
	
	// Enhanced confidence calculation
	confidence := bestScore
	if fastResult.Language == statResult.Language && fastResult.Language == bestLang {
		confidence = math.Min(confidence*1.2, 0.95)
	}
	
	collector.IncrementAlgorithmSteps()
	collector.RecordMemoryUsage()
	
	return LanguageResult{
		Language:     bestLang,
		Confidence:   confidence,
		Alternatives: alternatives,
	}
}

// Helper functions

func hasChineseCharacters(charFreq map[rune]int) bool {
	for r := range charFreq {
		if r >= 0x4E00 && r <= 0x9FFF {
			return true
		}
	}
	return false
}

func hasJapaneseCharacters(charFreq map[rune]int) bool {
	for r := range charFreq {
		// Hiragana or Katakana
		if (r >= 0x3040 && r <= 0x309F) || (r >= 0x30A0 && r <= 0x30FF) {
			return true
		}
	}
	return false
}

func hasCyrillicCharacters(charFreq map[rune]int) bool {
	for r := range charFreq {
		if r >= 0x0400 && r <= 0x04FF {
			return true
		}
	}
	return false
}

func findBestLanguage(scores map[string]float64) (string, float64) {
	bestLang := "en" // Default
	bestScore := 0.0
	
	for lang, score := range scores {
		if score > bestScore {
			bestScore = score
			bestLang = lang
		}
	}
	
	return bestLang, bestScore
}

func findAlternatives(scores map[string]float64, bestLang string, maxAlts int) []LanguageCandidate {
	type langScore struct {
		lang  string
		score float64
	}
	
	// Collect all scores except the best
	var alternatives []langScore
	for lang, score := range scores {
		if lang != bestLang && score > 0.1 {
			alternatives = append(alternatives, langScore{lang, score})
		}
	}
	
	// Sort by score
	sort.Slice(alternatives, func(i, j int) bool {
		return alternatives[i].score > alternatives[j].score
	})
	
	// Convert to LanguageCandidate
	candidates := []LanguageCandidate{}
	for i := 0; i < len(alternatives) && i < maxAlts; i++ {
		reason := "Statistical match"
		if alternatives[i].score > 0.5 {
			reason = "Strong statistical match"
		} else if alternatives[i].score > 0.3 {
			reason = "Moderate statistical match"
		} else {
			reason = "Weak statistical match"
		}
		
		candidates = append(candidates, LanguageCandidate{
			Language:   alternatives[i].lang,
			Confidence: alternatives[i].score,
			Reason:     reason,
		})
	}
	
	return candidates
}

func calculateConfidence(scores map[string]float64, bestLang string) float64 {
	if len(scores) == 0 {
		return 0.5
	}
	
	bestScore := scores[bestLang]
	
	// Calculate margin over second best
	secondBest := 0.0
	for lang, score := range scores {
		if lang != bestLang && score > secondBest {
			secondBest = score
		}
	}
	
	margin := bestScore - secondBest
	
	// Base confidence on absolute score and margin
	confidence := 0.3 // Base confidence
	
	// Add confidence based on absolute score
	if bestScore > 2.0 {
		confidence += 0.4
	} else if bestScore > 1.0 {
		confidence += 0.3
	} else if bestScore > 0.5 {
		confidence += 0.2
	} else if bestScore > 0.2 {
		confidence += 0.1
	}
	
	// Add confidence based on margin
	if margin > 1.0 {
		confidence += 0.3
	} else if margin > 0.5 {
		confidence += 0.2
	} else if margin > 0.2 {
		confidence += 0.1
	}
	
	// Cap confidence
	if confidence > 0.95 {
		confidence = 0.95
	} else if confidence < 0.3 {
		confidence = 0.3
	}
	
	return confidence
}

func analyzeStructurePatterns(sentences []string) map[string]float64 {
	scores := make(map[string]float64)
	
	// Simple heuristics based on sentence patterns
	for _, sentence := range sentences {
		// Spanish: inverted punctuation
		if strings.Contains(sentence, "¿") || strings.Contains(sentence, "¡") {
			scores["es"] += 0.2
		}
		
		// French: specific patterns
		if strings.Contains(sentence, "qu'") || strings.Contains(sentence, "c'") ||
			strings.Contains(sentence, "d'") || strings.Contains(sentence, "l'") {
			scores["fr"] += 0.1
		}
		
		// German: capitalized nouns, compound words
		words := strings.Fields(sentence)
		capitalizedWords := 0
		for _, word := range words {
			if len(word) > 0 && unicode.IsUpper(rune(word[0])) {
				capitalizedWords++
			}
		}
		if float64(capitalizedWords)/float64(len(words)) > 0.3 {
			scores["de"] += 0.1
		}
	}
	
	// Normalize scores
	maxScore := 0.0
	for _, score := range scores {
		if score > maxScore {
			maxScore = score
		}
	}
	
	if maxScore > 0 {
		for lang := range scores {
			scores[lang] /= maxScore
		}
	}
	
	return scores
}