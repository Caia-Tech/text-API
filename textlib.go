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
	"fmt"
	"regexp"
	"sort"
	"strings"
	"unicode"
)

// Position represents a position in text
type Position struct {
	Start int
	End   int
}

// SplitIntoSentences splits text into individual sentences
func SplitIntoSentences(text string) []string {
	if text == "" {
		return []string{}
	}
	
	// Handle common abbreviations
	abbreviations := []string{
		"Mr.", "Mrs.", "Ms.", "Dr.", "Prof.", "Sr.", "Jr.",
		"vs.", "etc.", "i.e.", "e.g.", "U.S.", "U.K.", "U.N.",
		"Inc.", "Corp.", "Ltd.", "Co.", "a.m.", "p.m.",
	}
	
	// Temporarily replace abbreviations
	replacements := make(map[string]string)
	for i, abbr := range abbreviations {
		placeholder := fmt.Sprintf("\x00ABBR%d\x00", i)
		replacements[placeholder] = abbr
		text = strings.ReplaceAll(text, abbr, placeholder)
	}
	
	// Split on sentence-ending punctuation followed by whitespace and capital letter
	sentencePattern := `([.!?]+)\s+([A-Z])`
	re := regexp.MustCompile(sentencePattern)
	
	// Find all split points
	matches := re.FindAllStringSubmatchIndex(text, -1)
	
	sentences := []string{}
	lastEnd := 0
	
	for _, match := range matches {
		// Split after the punctuation, including the punctuation
		splitPoint := match[3] // End of punctuation group (after the punctuation)
		sentence := strings.TrimSpace(text[lastEnd:splitPoint])
		
		if sentence != "" {
			sentences = append(sentences, sentence)
		}
		
		lastEnd = match[4] // Start of next sentence (capital letter)
	}
	
	// Add the final sentence
	if lastEnd < len(text) {
		finalSentence := strings.TrimSpace(text[lastEnd:])
		if finalSentence != "" {
			sentences = append(sentences, finalSentence)
		}
	}
	
	// If no sentences were found, treat the whole text as one sentence
	if len(sentences) == 0 && strings.TrimSpace(text) != "" {
		sentences = append(sentences, strings.TrimSpace(text))
	}
	
	// Restore abbreviations
	for i := range sentences {
		for placeholder, abbr := range replacements {
			sentences[i] = strings.ReplaceAll(sentences[i], placeholder, abbr)
		}
	}
	
	return sentences
}

// SplitIntoParagraphs splits text into paragraphs
func SplitIntoParagraphs(text string) []string {
	if text == "" {
		return []string{}
	}
	
	// Split on double newlines (paragraph breaks)
	paragraphs := regexp.MustCompile(`\n\s*\n`).Split(text, -1)
	
	result := []string{}
	for _, para := range paragraphs {
		para = strings.TrimSpace(para)
		if para != "" {
			result = append(result, para)
		}
	}
	
	// If no paragraph breaks found, treat whole text as one paragraph
	if len(result) == 0 {
		cleaned := strings.TrimSpace(text)
		if cleaned != "" {
			result = append(result, cleaned)
		}
	}
	
	return result
}

// ExtractNamedEntities extracts basic named entities from text
func ExtractNamedEntities(text string) []Entity {
	entities := []Entity{}
	
	// Simple patterns for basic entity recognition
	patterns := map[string]string{
		"PERSON": `\b[A-Z][a-z]+\s+[A-Z][a-z]+\b`,
		"ORGANIZATION": `\b[A-Z][a-z]*\s+(Inc|Corp|LLC|Ltd|Company|Corporation)\b`,
		"LOCATION": `\b[A-Z][a-z]+,\s*[A-Z][A-Z]\b`, // City, State
		"DATE": `\b(January|February|March|April|May|June|July|August|September|October|November|December)\s+\d{1,2},?\s+\d{4}\b`,
		"TIME": `\b\d{1,2}:\d{2}\s*(AM|PM|am|pm)?\b`,
		"MONEY": `\$\d+(?:,\d{3})*(?:\.\d{2})?\b`,
		"PERCENT": `\d+(?:\.\d+)?%\b`,
	}
	
	for entityType, pattern := range patterns {
		re := regexp.MustCompile(pattern)
		matches := re.FindAllStringSubmatchIndex(text, -1)
		
		for _, match := range matches {
			entity := Entity{
				Type: entityType,
				Text: text[match[0]:match[1]],
				Position: Position{
					Start: match[0],
					End:   match[1],
				},
			}
			entities = append(entities, entity)
		}
	}
	
	// Sort entities by position
	sort.Slice(entities, func(i, j int) bool {
		return entities[i].Position.Start < entities[j].Position.Start
	})
	
	return entities
}

// Entity represents a named entity in text
type Entity struct {
	Type     string
	Text     string
	Position Position
}

// CalculateSyllableCount estimates the number of syllables in a word
func CalculateSyllableCount(word string) int {
	if len(word) == 0 {
		return 0
	}
	
	word = strings.ToLower(word)
	word = regexp.MustCompile(`[^a-z]`).ReplaceAllString(word, "")
	
	if len(word) == 0 {
		return 0
	}
	
	// Count vowel groups
	vowelPattern := regexp.MustCompile(`[aeiouy]+`)
	syllables := len(vowelPattern.FindAllString(word, -1))
	
	// Adjust for silent 'e'
	if strings.HasSuffix(word, "e") && syllables > 1 {
		syllables--
	}
	
	// Adjust for 'le' ending
	if strings.HasSuffix(word, "le") && len(word) > 2 && !isVowel(rune(word[len(word)-3])) {
		syllables++
	}
	
	// Every word has at least one syllable
	if syllables == 0 {
		syllables = 1
	}
	
	return syllables
}

func isVowel(r rune) bool {
	vowels := "aeiouy"
	return strings.ContainsRune(vowels, unicode.ToLower(r))
}

// Helper functions used across modules - using min/max from complexity.go

// IsCompleteSentence checks if text represents a complete sentence
func IsCompleteSentence(text string) bool {
	text = strings.TrimSpace(text)
	if len(text) < 3 {
		return false
	}
	
	// Check for sentence-ending punctuation
	lastChar := text[len(text)-1]
	if lastChar != '.' && lastChar != '!' && lastChar != '?' {
		return false
	}
	
	// Check for basic sentence structure (subject and predicate)
	words := strings.Fields(text)
	if len(words) < 2 {
		return false
	}
	
	// Simple heuristic: contains a verb
	hasVerb := false
	for _, word := range words {
		if isVerbSimple(word) {
			hasVerb = true
			break
		}
	}
	
	return hasVerb
}

func isVerbSimple(word string) bool {
	word = strings.ToLower(strings.Trim(word, ".,!?;:"))
	
	// Common verbs and patterns
	commonVerbs := []string{
		"is", "are", "was", "were", "be", "been", "being",
		"have", "has", "had", "do", "does", "did", "will",
		"would", "could", "should", "can", "may", "might",
	}
	
	for _, verb := range commonVerbs {
		if word == verb {
			return true
		}
	}
	
	// Simple patterns for regular verbs
	if strings.HasSuffix(word, "ed") || strings.HasSuffix(word, "ing") || strings.HasSuffix(word, "s") {
		return true
	}
	
	return false
}

// CountWords counts the number of words in text
func CountWords(text string) int {
	if text == "" {
		return 0
	}
	
	words := strings.Fields(text)
	return len(words)
}

// CountSentences counts the number of sentences in text
func CountSentences(text string) int {
	sentences := SplitIntoSentences(text)
	return len(sentences)
}

// CountSyllables counts total syllables in text
func CountSyllables(text string) int {
	words := strings.Fields(text)
	total := 0
	
	for _, word := range words {
		// Remove punctuation
		cleanWord := regexp.MustCompile(`[^a-zA-Z]`).ReplaceAllString(word, "")
		if cleanWord != "" {
			total += CalculateSyllableCount(cleanWord)
		}
	}
	
	return total
}

// CalculateFleschReadingEase calculates the Flesch Reading Ease score
func CalculateFleschReadingEase(text string) float64 {
	sentences := SplitIntoSentences(text)
	words := strings.Fields(text)
	syllables := CountSyllables(text)
	
	if len(sentences) == 0 || len(words) == 0 {
		return 0
	}
	
	avgSentenceLength := float64(len(words)) / float64(len(sentences))
	avgSyllablesPerWord := float64(syllables) / float64(len(words))
	
	score := 206.835 - (1.015 * avgSentenceLength) - (84.6 * avgSyllablesPerWord)
	
	return score
}