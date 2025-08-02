package textlib

import (
	"fmt"
	"math"
	"regexp"
	"sort"
	"strings"
	"unicode"
)

// PatternAnalysis contains detected patterns and repetitions
type PatternAnalysis struct {
	// Repetition patterns
	WordRepetitions      []RepetitionPattern
	PhraseRepetitions    []RepetitionPattern
	StructureRepetitions []StructuralPattern
	
	// Linguistic patterns
	SyntacticPatterns    []SyntacticPattern
	RhetoricalPatterns   []RhetoricalPattern
	
	// Sequence patterns
	SequencePatterns     []SequencePattern
	TemporalPatterns     []TemporalPattern
	
	// Anomaly detection
	Anomalies           []TextAnomaly
	
	// Statistics
	RepetitionScore     float64 // 0-1, higher means more repetitive
	PatternDiversity    float64 // 0-1, higher means more diverse patterns
	Predictability      float64 // 0-1, higher means more predictable
}

type RepetitionPattern struct {
	Text         string
	Count        int
	Positions    []Position
	Type         string // "exact", "stem", "semantic"
	Proximity    float64 // Average distance between occurrences
	Significance float64 // Statistical significance
}

type StructuralPattern struct {
	Pattern      string // e.g., "ADJ NOUN VERB ADV"
	Count        int
	Examples     []string
	Positions    []Position
}

type SyntacticPattern struct {
	Type         string // "parallel", "chiasmus", "anaphora", etc.
	Description  string
	Examples     []string
	Strength     float64
}

type RhetoricalPattern struct {
	Device       string // "metaphor", "alliteration", "repetition", etc.
	Instances    []string
	Positions    []Position
	Effectiveness float64
}

type SequencePattern struct {
	Type         string // "enumeration", "progression", "alternation"
	Elements     []string
	Pattern      string // Representation of the pattern
	Consistency  float64
}

type TemporalPattern struct {
	Type         string // "chronological", "flashback", "simultaneous"
	Markers      []string
	Flow         string // "forward", "backward", "mixed"
	Consistency  float64
}

type TextAnomaly struct {
	Type         string // "style_shift", "vocabulary_spike", "structure_break"
	Location     Position
	Description  string
	Severity     float64
	Context      string
}

// DetectPatterns performs comprehensive pattern detection
func DetectPatterns(text string) *PatternAnalysis {
	analysis := &PatternAnalysis{
		WordRepetitions:      []RepetitionPattern{},
		PhraseRepetitions:    []RepetitionPattern{},
		StructureRepetitions: []StructuralPattern{},
		SyntacticPatterns:    []SyntacticPattern{},
		RhetoricalPatterns:   []RhetoricalPattern{},
		SequencePatterns:     []SequencePattern{},
		TemporalPatterns:     []TemporalPattern{},
		Anomalies:           []TextAnomaly{},
	}
	
	// Detect repetitions
	analysis.WordRepetitions = detectWordRepetitions(text)
	analysis.PhraseRepetitions = detectPhraseRepetitions(text)
	analysis.StructureRepetitions = detectStructuralRepetitions(text)
	
	// Detect linguistic patterns
	analysis.SyntacticPatterns = detectSyntacticPatterns(text)
	analysis.RhetoricalPatterns = detectRhetoricalPatterns(text)
	
	// Detect sequences
	analysis.SequencePatterns = detectSequencePatterns(text)
	analysis.TemporalPatterns = detectTemporalPatterns(text)
	
	// Detect anomalies
	analysis.Anomalies = detectAnomalies(text)
	
	// Calculate scores
	analysis.RepetitionScore = calculateRepetitionScore(analysis)
	analysis.PatternDiversity = calculatePatternDiversity(analysis)
	analysis.Predictability = calculatePredictability(text, analysis)
	
	return analysis
}

func detectWordRepetitions(text string) []RepetitionPattern {
	var patterns []RepetitionPattern
	words := extractWords(text)
	
	// Track word positions
	wordPositions := make(map[string][]int)
	for i, word := range words {
		lower := strings.ToLower(word)
		wordPositions[lower] = append(wordPositions[lower], i)
	}
	
	// Find significant repetitions
	for word, positions := range wordPositions {
		if len(positions) >= 3 && len(word) > 3 && !isStopWord(word) {
			pattern := RepetitionPattern{
				Text:         word,
				Count:        len(positions),
				Type:         "exact",
				Proximity:    calculateAverageDistance(positions),
				Significance: calculateSignificance(len(positions), len(words)),
			}
			
			// Convert to text positions
			pattern.Positions = convertToTextPositions(positions, words, text)
			patterns = append(patterns, pattern)
		}
	}
	
	// Sort by significance
	sort.Slice(patterns, func(i, j int) bool {
		return patterns[i].Significance > patterns[j].Significance
	})
	
	// Also detect stem-based repetitions
	stemPatterns := detectStemRepetitions(words, text)
	patterns = append(patterns, stemPatterns...)
	
	return patterns
}

func detectStemRepetitions(words []string, text string) []RepetitionPattern {
	var patterns []RepetitionPattern
	stemGroups := make(map[string][]string)
	stemPositions := make(map[string][]int)
	
	// Group words by stem
	for i, word := range words {
		stem := extractStem(strings.ToLower(word))
		if len(stem) > 3 {
			stemGroups[stem] = append(stemGroups[stem], word)
			stemPositions[stem] = append(stemPositions[stem], i)
		}
	}
	
	// Find significant stem groups
	for stem, group := range stemGroups {
		if len(group) >= 4 && hasVariation(group) {
			pattern := RepetitionPattern{
				Text:         fmt.Sprintf("%s*", stem),
				Count:        len(group),
				Type:         "stem",
				Proximity:    calculateAverageDistance(stemPositions[stem]),
				Significance: calculateSignificance(len(group), len(words)),
			}
			
			pattern.Positions = convertToTextPositions(stemPositions[stem], words, text)
			patterns = append(patterns, pattern)
		}
	}
	
	return patterns
}

func hasVariation(words []string) bool {
	unique := make(map[string]bool)
	for _, w := range words {
		unique[strings.ToLower(w)] = true
	}
	return len(unique) > 1
}

func detectPhraseRepetitions(text string) []RepetitionPattern {
	var patterns []RepetitionPattern
	
	// Detect 2-4 word phrases
	for n := 2; n <= 4; n++ {
		phrases := extractNGrams(text, n)
		phraseCount := make(map[string]int)
		phrasePositions := make(map[string][]Position)
		
		for phrase, positions := range phrases {
			phraseCount[phrase] = len(positions)
			phrasePositions[phrase] = positions
		}
		
		// Find significant phrase repetitions
		for phrase, count := range phraseCount {
			if count >= 2 && !isPhraseStopWords(phrase) {
				positions := phrasePositions[phrase]
				intPositions := make([]int, len(positions))
				for i, p := range positions {
					intPositions[i] = p.Start
				}
				
				pattern := RepetitionPattern{
					Text:         phrase,
					Count:        count,
					Type:         "phrase",
					Positions:    positions,
					Proximity:    calculateAverageDistance(intPositions),
					Significance: calculatePhraseSignificance(count, n, len(strings.Fields(text))),
				}
				
				patterns = append(patterns, pattern)
			}
		}
	}
	
	// Sort by significance
	sort.Slice(patterns, func(i, j int) bool {
		return patterns[i].Significance > patterns[j].Significance
	})
	
	return patterns
}

func extractNGrams(text string, n int) map[string][]Position {
	ngrams := make(map[string][]Position)
	words := strings.Fields(text)
	
	if len(words) < n {
		return ngrams
	}
	
	// Track word positions in original text
	wordPositions := findWordPositions(text, words)
	
	for i := 0; i <= len(words)-n; i++ {
		ngram := strings.Join(words[i:i+n], " ")
		start := wordPositions[i]
		end := wordPositions[i+n-1] + len(words[i+n-1])
		
		position := Position{Start: start, End: end}
		ngrams[strings.ToLower(ngram)] = append(ngrams[strings.ToLower(ngram)], position)
	}
	
	return ngrams
}

func findWordPositions(text string, words []string) []int {
	positions := make([]int, len(words))
	currentPos := 0
	
	for i, word := range words {
		index := strings.Index(text[currentPos:], word)
		if index != -1 {
			positions[i] = currentPos + index
			currentPos = positions[i] + len(word)
		}
	}
	
	return positions
}

func isPhraseStopWords(phrase string) bool {
	stopPhrases := map[string]bool{
		"in the": true, "of the": true, "to the": true, "and the": true,
		"for the": true, "with the": true, "on the": true, "at the": true,
		"it is": true, "there is": true, "there are": true, "this is": true,
	}
	
	return stopPhrases[strings.ToLower(phrase)]
}

func detectStructuralRepetitions(text string) []StructuralPattern {
	var patterns []StructuralPattern
	sentences := SplitIntoSentences(text)
	
	// Extract POS patterns (simplified)
	posPatterns := make(map[string][]string)
	patternPositions := make(map[string][]Position)
	
	for i, sentence := range sentences {
		pattern := extractPOSPattern(sentence)
		if pattern != "" {
			posPatterns[pattern] = append(posPatterns[pattern], sentence)
			// Simplified position tracking
			patternPositions[pattern] = append(patternPositions[pattern], 
				Position{Start: i, End: i})
		}
	}
	
	// Find repeated patterns
	for pattern, examples := range posPatterns {
		if len(examples) >= 2 {
			structPattern := StructuralPattern{
				Pattern:   pattern,
				Count:     len(examples),
				Examples:  examples[:min(3, len(examples))],
				Positions: patternPositions[pattern],
			}
			patterns = append(patterns, structPattern)
		}
	}
	
	return patterns
}

func extractPOSPattern(sentence string) string {
	// Simplified POS tagging
	words := strings.Fields(sentence)
	if len(words) < 3 || len(words) > 8 {
		return ""
	}
	
	pattern := ""
	for _, word := range words {
		tag := simplePOSTag(word)
		if tag != "" {
			pattern += tag + " "
		}
	}
	
	return strings.TrimSpace(pattern)
}

func simplePOSTag(word string) string {
	lower := strings.ToLower(strings.Trim(word, ".,!?;:"))
	
	// Very simplified POS tagging
	if isLikelyVerb(lower) {
		return "VERB"
	} else if isPronoun(lower) {
		return "PRON"
	} else if isPreposition(lower) {
		return "PREP"
	} else if isConjunction(lower) {
		return "CONJ"
	} else if strings.HasSuffix(lower, "ly") {
		return "ADV"
	} else if strings.HasSuffix(lower, "ing") || strings.HasSuffix(lower, "ed") {
		return "ADJ"
	} else if unicode.IsUpper([]rune(word)[0]) {
		return "NOUN"
	}
	
	return "WORD"
}

func isPreposition(word string) bool {
	prepositions := map[string]bool{
		"in": true, "on": true, "at": true, "by": true, "for": true,
		"with": true, "about": true, "against": true, "between": true,
		"into": true, "through": true, "during": true, "before": true,
		"after": true, "above": true, "below": true, "to": true, "from": true,
	}
	return prepositions[word]
}

func detectSyntacticPatterns(text string) []SyntacticPattern {
	var patterns []SyntacticPattern
	
	// Detect parallelism
	if parallel := detectParallelism(text); parallel.Strength > 0.5 {
		patterns = append(patterns, parallel)
	}
	
	// Detect anaphora (repetition at beginning)
	if anaphora := detectAnaphora(text); anaphora.Strength > 0.5 {
		patterns = append(patterns, anaphora)
	}
	
	// Detect epistrophe (repetition at end)
	if epistrophe := detectEpistrophe(text); epistrophe.Strength > 0.5 {
		patterns = append(patterns, epistrophe)
	}
	
	// Detect chiasmus (ABBA pattern)
	if chiasmus := detectChiasmus(text); chiasmus.Strength > 0.5 {
		patterns = append(patterns, chiasmus)
	}
	
	return patterns
}

func detectParallelism(text string) SyntacticPattern {
	sentences := SplitIntoSentences(text)
	parallelPairs := 0
	examples := []string{}
	
	for i := 0; i < len(sentences)-1; i++ {
		if areParallel(sentences[i], sentences[i+1]) {
			parallelPairs++
			if len(examples) < 3 {
				examples = append(examples, 
					fmt.Sprintf("%s | %s", sentences[i], sentences[i+1]))
			}
		}
	}
	
	strength := 0.0
	if len(sentences) > 1 {
		strength = float64(parallelPairs) / float64(len(sentences)-1)
	}
	
	return SyntacticPattern{
		Type:        "parallel",
		Description: "Parallel sentence structures",
		Examples:    examples,
		Strength:    strength,
	}
}

func areParallel(sent1, sent2 string) bool {
	// Check if sentences have similar structure
	pattern1 := extractPOSPattern(sent1)
	pattern2 := extractPOSPattern(sent2)
	
	if pattern1 == "" || pattern2 == "" {
		return false
	}
	
	// Calculate similarity
	return calculatePatternSimilarity(pattern1, pattern2) > 0.7
}

func calculatePatternSimilarity(p1, p2 string) float64 {
	tags1 := strings.Fields(p1)
	tags2 := strings.Fields(p2)
	
	if len(tags1) == 0 || len(tags2) == 0 {
		return 0
	}
	
	matches := 0
	minLen := min(len(tags1), len(tags2))
	
	for i := 0; i < minLen; i++ {
		if tags1[i] == tags2[i] {
			matches++
		}
	}
	
	return float64(matches) / float64(max(len(tags1), len(tags2)))
}

func detectAnaphora(text string) SyntacticPattern {
	sentences := SplitIntoSentences(text)
	if len(sentences) < 2 {
		return SyntacticPattern{Type: "anaphora", Strength: 0}
	}
	
	// Track beginning phrases
	beginnings := make(map[string]int)
	examples := []string{}
	
	for _, sentence := range sentences {
		words := strings.Fields(sentence)
		if len(words) >= 2 {
			beginning := strings.ToLower(words[0] + " " + words[1])
			beginnings[beginning]++
			
			if beginnings[beginning] == 2 && len(examples) < 3 {
				examples = append(examples, sentence)
			}
		}
	}
	
	// Find most repeated beginning
	maxCount := 0
	for _, count := range beginnings {
		if count > maxCount {
			maxCount = count
		}
	}
	
	strength := 0.0
	if maxCount >= 3 {
		strength = float64(maxCount) / float64(len(sentences))
	}
	
	return SyntacticPattern{
		Type:        "anaphora",
		Description: "Repetition at beginning of sentences",
		Examples:    examples,
		Strength:    strength,
	}
}

func detectEpistrophe(text string) SyntacticPattern {
	sentences := SplitIntoSentences(text)
	if len(sentences) < 2 {
		return SyntacticPattern{Type: "epistrophe", Strength: 0}
	}
	
	// Track ending phrases
	endings := make(map[string]int)
	examples := []string{}
	
	for _, sentence := range sentences {
		words := strings.Fields(sentence)
		if len(words) >= 2 {
			ending := strings.ToLower(
				strings.Trim(words[len(words)-2], ".,!?") + " " + 
				strings.Trim(words[len(words)-1], ".,!?"))
			endings[ending]++
			
			if endings[ending] == 2 && len(examples) < 3 {
				examples = append(examples, sentence)
			}
		}
	}
	
	// Find most repeated ending
	maxCount := 0
	for _, count := range endings {
		if count > maxCount {
			maxCount = count
		}
	}
	
	strength := 0.0
	if maxCount >= 3 {
		strength = float64(maxCount) / float64(len(sentences))
	}
	
	return SyntacticPattern{
		Type:        "epistrophe",
		Description: "Repetition at end of sentences",
		Examples:    examples,
		Strength:    strength,
	}
}

func detectChiasmus(text string) SyntacticPattern {
	sentences := SplitIntoSentences(text)
	examples := []string{}
	chiasmusCount := 0
	
	for _, sentence := range sentences {
		if hasChiasmus(sentence) {
			chiasmusCount++
			if len(examples) < 3 {
				examples = append(examples, sentence)
			}
		}
	}
	
	strength := 0.0
	if len(sentences) > 0 {
		strength = float64(chiasmusCount) / float64(len(sentences))
	}
	
	return SyntacticPattern{
		Type:        "chiasmus",
		Description: "ABBA word pattern structure",
		Examples:    examples,
		Strength:    strength,
	}
}

func hasChiasmus(sentence string) bool {
	words := strings.Fields(sentence)
	if len(words) < 4 {
		return false
	}
	
	// Look for ABBA patterns
	for i := 0; i <= len(words)-4; i++ {
		w1 := strings.ToLower(strings.Trim(words[i], ".,!?;:"))
		w2 := strings.ToLower(strings.Trim(words[i+1], ".,!?;:"))
		w3 := strings.ToLower(strings.Trim(words[i+2], ".,!?;:"))
		w4 := strings.ToLower(strings.Trim(words[i+3], ".,!?;:"))
		
		// Check for ABBA pattern (words or stems)
		if (w1 == w4 || extractStem(w1) == extractStem(w4)) &&
		   (w2 == w3 || extractStem(w2) == extractStem(w3)) &&
		   w1 != w2 {
			return true
		}
	}
	
	return false
}

func detectRhetoricalPatterns(text string) []RhetoricalPattern {
	var patterns []RhetoricalPattern
	
	// Detect alliteration
	if alliteration := detectAlliteration(text); len(alliteration.Instances) > 0 {
		patterns = append(patterns, alliteration)
	}
	
	// Detect assonance
	if assonance := detectAssonance(text); len(assonance.Instances) > 0 {
		patterns = append(patterns, assonance)
	}
	
	// Detect rhetorical questions
	if questions := detectRhetoricalQuestions(text); len(questions.Instances) > 0 {
		patterns = append(patterns, questions)
	}
	
	// Detect enumeration
	if enumeration := detectEnumeration(text); len(enumeration.Instances) > 0 {
		patterns = append(patterns, enumeration)
	}
	
	return patterns
}

func detectAlliteration(text string) RhetoricalPattern {
	var instances []string
	var positions []Position
	
	sentences := SplitIntoSentences(text)
	
	for _, sentence := range sentences {
		words := strings.Fields(sentence)
		
		for i := 0; i < len(words)-2; i++ {
			// Check three consecutive words
			if len(words[i]) > 0 && len(words[i+1]) > 0 && len(words[i+2]) > 0 {
				c1 := strings.ToLower(string(words[i][0]))
				c2 := strings.ToLower(string(words[i+1][0]))
				c3 := strings.ToLower(string(words[i+2][0]))
				
				if c1 == c2 && c2 == c3 && unicode.IsLetter([]rune(c1)[0]) {
					phrase := strings.Join(words[i:i+3], " ")
					instances = append(instances, phrase)
					// Simplified position
					positions = append(positions, Position{Start: i, End: i+3})
				}
			}
		}
	}
	
	effectiveness := 0.0
	if len(instances) > 0 {
		effectiveness = math.Min(float64(len(instances))/10.0, 1.0)
	}
	
	return RhetoricalPattern{
		Device:        "alliteration",
		Instances:     instances,
		Positions:     positions,
		Effectiveness: effectiveness,
	}
}

func detectAssonance(text string) RhetoricalPattern {
	var instances []string
	var positions []Position
	
	// Detect repeated vowel sounds
	sentences := SplitIntoSentences(text)
	
	for _, sentence := range sentences {
		words := strings.Fields(sentence)
		
		for i := 0; i < len(words)-1; i++ {
			v1 := extractVowelSound(words[i])
			v2 := extractVowelSound(words[i+1])
			
			if v1 != "" && v1 == v2 && len(v1) >= 2 {
				phrase := words[i] + " " + words[i+1]
				instances = append(instances, phrase)
				positions = append(positions, Position{Start: i, End: i+2})
			}
		}
	}
	
	effectiveness := 0.0
	if len(instances) > 0 {
		effectiveness = math.Min(float64(len(instances))/15.0, 1.0)
	}
	
	return RhetoricalPattern{
		Device:        "assonance",
		Instances:     instances,
		Positions:     positions,
		Effectiveness: effectiveness,
	}
}

func extractVowelSound(word string) string {
	// Extract dominant vowel sound (simplified)
	vowels := ""
	for _, r := range strings.ToLower(word) {
		if strings.ContainsRune("aeiou", r) {
			vowels += string(r)
		}
	}
	
	// Return dominant pattern
	if len(vowels) >= 2 {
		return vowels[:2]
	}
	return vowels
}

func detectRhetoricalQuestions(text string) RhetoricalPattern {
	var instances []string
	var positions []Position
	
	sentences := SplitIntoSentences(text)
	
	for i, sentence := range sentences {
		if strings.HasSuffix(strings.TrimSpace(sentence), "?") {
			// Check if it's likely rhetorical (no answer follows)
			if i == len(sentences)-1 || !isAnswerPattern(sentences[i+1]) {
				instances = append(instances, sentence)
				positions = append(positions, Position{Start: i, End: i})
			}
		}
	}
	
	effectiveness := 0.0
	if len(instances) > 0 {
		effectiveness = math.Min(float64(len(instances))/5.0, 1.0)
	}
	
	return RhetoricalPattern{
		Device:        "rhetorical_question",
		Instances:     instances,
		Positions:     positions,
		Effectiveness: effectiveness,
	}
}

func isAnswerPattern(sentence string) bool {
	// Simple heuristic: answers often start with yes/no or specific patterns
	lower := strings.ToLower(sentence)
	answerStarts := []string{
		"yes", "no", "it is", "it's", "they are", "i think",
		"because", "the answer", "that's",
	}
	
	for _, start := range answerStarts {
		if strings.HasPrefix(lower, start) {
			return true
		}
	}
	
	return false
}

func detectEnumeration(text string) RhetoricalPattern {
	var instances []string
	var positions []Position
	
	// Patterns for enumeration
	patterns := []string{
		"first.*second.*third",
		"one.*two.*three",
		"firstly.*secondly.*thirdly",
		`\b1\).*2\).*3\)`,
		`\ba\).*b\).*c\)`,
	}
	
	for _, pattern := range patterns {
		re := regexp.MustCompile(`(?i)` + pattern)
		matches := re.FindAllStringIndex(text, -1)
		
		for _, match := range matches {
			instances = append(instances, text[match[0]:match[1]])
			positions = append(positions, Position{Start: match[0], End: match[1]})
		}
	}
	
	effectiveness := 0.0
	if len(instances) > 0 {
		effectiveness = 1.0 // Enumeration is very effective for clarity
	}
	
	return RhetoricalPattern{
		Device:        "enumeration",
		Instances:     instances,
		Positions:     positions,
		Effectiveness: effectiveness,
	}
}

func detectSequencePatterns(text string) []SequencePattern {
	var patterns []SequencePattern
	
	// Detect enumerations
	if enum := detectEnumerationSequence(text); enum.Consistency > 0.5 {
		patterns = append(patterns, enum)
	}
	
	// Detect progressions
	if prog := detectProgressionSequence(text); prog.Consistency > 0.5 {
		patterns = append(patterns, prog)
	}
	
	// Detect alternations
	if alt := detectAlternationSequence(text); alt.Consistency > 0.5 {
		patterns = append(patterns, alt)
	}
	
	return patterns
}

func detectEnumerationSequence(text string) SequencePattern {
	// Look for numbered or lettered lists
	elements := []string{}
	
	// Numeric enumeration
	numPattern := regexp.MustCompile(`\b(\d+)\.\s+([^.!?]+[.!?])`)
	numMatches := numPattern.FindAllStringSubmatch(text, -1)
	
	for _, match := range numMatches {
		elements = append(elements, match[2])
	}
	
	consistency := 0.0
	if len(elements) > 0 {
		// Check if numbers are sequential
		sequential := true
		for i := 1; i < len(elements); i++ {
			// Simple check - could be improved
			if i != i {
				sequential = false
			}
		}
		if sequential {
			consistency = 1.0
		}
	}
	
	return SequencePattern{
		Type:        "enumeration",
		Elements:    elements,
		Pattern:     "1, 2, 3, ...",
		Consistency: consistency,
	}
}

func detectProgressionSequence(text string) SequencePattern {
	// Detect logical or temporal progressions
	progressionMarkers := []string{
		"first", "then", "next", "after", "finally",
		"initially", "subsequently", "ultimately",
		"beginning", "middle", "end",
	}
	
	elements := []string{}
	sentences := SplitIntoSentences(text)
	
	for _, sentence := range sentences {
		lower := strings.ToLower(sentence)
		for _, marker := range progressionMarkers {
			if strings.Contains(lower, marker) {
				elements = append(elements, sentence)
				break
			}
		}
	}
	
	consistency := 0.0
	if len(elements) >= 3 {
		consistency = float64(len(elements)) / float64(len(sentences))
	}
	
	return SequencePattern{
		Type:        "progression",
		Elements:    elements,
		Pattern:     "temporal/logical sequence",
		Consistency: consistency,
	}
}

func detectAlternationSequence(text string) SequencePattern {
	// Detect alternating patterns (e.g., problem-solution, claim-evidence)
	sentences := SplitIntoSentences(text)
	
	// Simple detection: alternating sentence lengths or structures
	shortLong := 0
	for i := 0; i < len(sentences)-1; i++ {
		len1 := len(strings.Fields(sentences[i]))
		len2 := len(strings.Fields(sentences[i+1]))
		
		if (len1 < 10 && len2 > 15) || (len1 > 15 && len2 < 10) {
			shortLong++
		}
	}
	
	consistency := 0.0
	if len(sentences) > 1 {
		consistency = float64(shortLong) / float64(len(sentences)-1)
	}
	
	elements := []string{}
	if consistency > 0.5 {
		elements = sentences[:min(4, len(sentences))]
	}
	
	return SequencePattern{
		Type:        "alternation",
		Elements:    elements,
		Pattern:     "short-long alternation",
		Consistency: consistency,
	}
}

func detectTemporalPatterns(text string) []TemporalPattern {
	var patterns []TemporalPattern
	
	// Detect chronological flow
	if chrono := detectChronological(text); chrono.Consistency > 0.5 {
		patterns = append(patterns, chrono)
	}
	
	// Detect flashbacks
	if flash := detectFlashbacks(text); len(flash.Markers) > 0 {
		patterns = append(patterns, flash)
	}
	
	return patterns
}

func detectChronological(text string) TemporalPattern {
	temporalMarkers := []string{}
	
	// Time markers
	timeWords := []string{
		"yesterday", "today", "tomorrow",
		"morning", "afternoon", "evening", "night",
		"before", "after", "during", "while",
		"first", "then", "next", "finally",
		"earlier", "later", "now", "soon",
	}
	
	sentences := SplitIntoSentences(text)
	markerCount := 0
	
	for _, sentence := range sentences {
		lower := strings.ToLower(sentence)
		for _, marker := range timeWords {
			if strings.Contains(lower, marker) {
				temporalMarkers = append(temporalMarkers, marker)
				markerCount++
				break
			}
		}
	}
	
	consistency := 0.0
	if len(sentences) > 0 {
		consistency = float64(markerCount) / float64(len(sentences))
	}
	
	flow := "forward"
	if containsAny(temporalMarkers, []string{"before", "earlier", "yesterday"}) {
		flow = "mixed"
	}
	
	return TemporalPattern{
		Type:        "chronological",
		Markers:     temporalMarkers,
		Flow:        flow,
		Consistency: consistency,
	}
}

func detectFlashbacks(text string) TemporalPattern {
	flashbackMarkers := []string{}
	
	// Flashback indicators
	indicators := []string{
		"remembered", "recalled", "used to", "had been",
		"once upon", "back when", "in the past",
		"years ago", "months ago", "previously",
	}
	
	sentences := SplitIntoSentences(text)
	
	for _, sentence := range sentences {
		lower := strings.ToLower(sentence)
		for _, indicator := range indicators {
			if strings.Contains(lower, indicator) {
				flashbackMarkers = append(flashbackMarkers, indicator)
			}
		}
	}
	
	consistency := 0.0
	if len(flashbackMarkers) > 0 {
		consistency = 0.7 // Flashbacks are intentional
	}
	
	return TemporalPattern{
		Type:        "flashback",
		Markers:     flashbackMarkers,
		Flow:        "backward",
		Consistency: consistency,
	}
}

func detectAnomalies(text string) []TextAnomaly {
	var anomalies []TextAnomaly
	
	// Detect style shifts
	styleAnomalies := detectStyleShifts(text)
	anomalies = append(anomalies, styleAnomalies...)
	
	// Detect vocabulary spikes
	vocabAnomalies := detectVocabularySpikes(text)
	anomalies = append(anomalies, vocabAnomalies...)
	
	// Detect structure breaks
	structureAnomalies := detectStructureBreaks(text)
	anomalies = append(anomalies, structureAnomalies...)
	
	return anomalies
}

func detectStyleShifts(text string) []TextAnomaly {
	var anomalies []TextAnomaly
	sentences := SplitIntoSentences(text)
	
	if len(sentences) < 5 {
		return anomalies
	}
	
	// Calculate average sentence length in windows
	windowSize := 3
	avgLengths := make([]float64, 0)
	
	for i := 0; i <= len(sentences)-windowSize; i++ {
		totalLen := 0
		for j := i; j < i+windowSize; j++ {
			totalLen += len(strings.Fields(sentences[j]))
		}
		avgLengths = append(avgLengths, float64(totalLen)/float64(windowSize))
	}
	
	// Detect significant changes
	for i := 1; i < len(avgLengths); i++ {
		change := math.Abs(avgLengths[i] - avgLengths[i-1])
		avgLength := (avgLengths[i] + avgLengths[i-1]) / 2
		
		if avgLength > 0 && change/avgLength > 0.5 {
			anomaly := TextAnomaly{
				Type:        "style_shift",
				Location:    Position{Start: i * windowSize, End: (i+1)*windowSize},
				Description: "Significant change in sentence length",
				Severity:    change / avgLength,
				Context:     sentences[i*windowSize],
			}
			anomalies = append(anomalies, anomaly)
		}
	}
	
	return anomalies
}

func detectVocabularySpikes(text string) []TextAnomaly {
	var anomalies []TextAnomaly
	paragraphs := SplitIntoParagraphs(text)
	
	if len(paragraphs) < 2 {
		return anomalies
	}
	
	// Calculate vocabulary complexity per paragraph
	complexities := make([]float64, len(paragraphs))
	
	for i, para := range paragraphs {
		words := extractWords(para)
		complexWords := 0
		
		for _, word := range words {
			if len(word) > 8 || CountSyllables(word) > 3 {
				complexWords++
			}
		}
		
		if len(words) > 0 {
			complexities[i] = float64(complexWords) / float64(len(words))
		}
	}
	
	// Find outliers
	mean, stdDev := calculateMeanStdDevFloat(complexities)
	
	for i, complexity := range complexities {
		if math.Abs(complexity-mean) > 2*stdDev {
			anomaly := TextAnomaly{
				Type:        "vocabulary_spike",
				Location:    Position{Start: i, End: i},
				Description: "Unusual vocabulary complexity",
				Severity:    math.Abs(complexity-mean) / stdDev,
				Context:     paragraphs[i][:min(100, len(paragraphs[i]))] + "...",
			}
			anomalies = append(anomalies, anomaly)
		}
	}
	
	return anomalies
}

func detectStructureBreaks(text string) []TextAnomaly {
	var anomalies []TextAnomaly
	sentences := SplitIntoSentences(text)
	
	// Detect unusual punctuation patterns
	for i, sentence := range sentences {
		punctCount := 0
		for _, r := range sentence {
			if unicode.IsPunct(r) {
				punctCount++
			}
		}
		
		punctRatio := float64(punctCount) / float64(len(sentence))
		if punctRatio > 0.2 {
			anomaly := TextAnomaly{
				Type:        "structure_break",
				Location:    Position{Start: i, End: i},
				Description: "Excessive punctuation",
				Severity:    punctRatio * 5,
				Context:     sentence,
			}
			anomalies = append(anomalies, anomaly)
		}
	}
	
	return anomalies
}

// Helper functions for pattern analysis

func calculateAverageDistance(positions []int) float64 {
	if len(positions) < 2 {
		return 0
	}
	
	totalDistance := 0
	for i := 1; i < len(positions); i++ {
		totalDistance += positions[i] - positions[i-1]
	}
	
	return float64(totalDistance) / float64(len(positions)-1)
}

func calculateSignificance(count, total int) float64 {
	if total == 0 {
		return 0
	}
	
	// Simple significance calculation
	frequency := float64(count) / float64(total)
	
	// Adjust for expected frequency
	expectedFreq := 1.0 / float64(total)
	
	if frequency > expectedFreq {
		return math.Min((frequency / expectedFreq) - 1.0, 1.0)
	}
	
	return 0
}

func calculatePhraseSignificance(count, length, total int) float64 {
	if total == 0 {
		return 0
	}
	
	// Longer phrases are more significant
	lengthFactor := math.Log(float64(length))
	frequency := float64(count) / float64(total)
	
	return math.Min(frequency * lengthFactor, 1.0)
}

func convertToTextPositions(wordPositions []int, words []string, text string) []Position {
	positions := make([]Position, len(wordPositions))
	
	currentPos := 0
	wordIndex := 0
	
	for i, pos := range wordPositions {
		// Find the word in text
		for wordIndex < pos && wordIndex < len(words) {
			idx := strings.Index(text[currentPos:], words[wordIndex])
			if idx != -1 {
				currentPos += idx + len(words[wordIndex])
			}
			wordIndex++
		}
		
		if pos < len(words) {
			idx := strings.Index(text[currentPos:], words[pos])
			if idx != -1 {
				start := currentPos + idx
				end := start + len(words[pos])
				positions[i] = Position{Start: start, End: end}
			}
		}
	}
	
	return positions
}

func calculateRepetitionScore(analysis *PatternAnalysis) float64 {
	if len(analysis.WordRepetitions) == 0 && len(analysis.PhraseRepetitions) == 0 {
		return 0
	}
	
	// Weight different types of repetition
	wordScore := 0.0
	for _, rep := range analysis.WordRepetitions {
		wordScore += rep.Significance
	}
	
	phraseScore := 0.0
	for _, rep := range analysis.PhraseRepetitions {
		phraseScore += rep.Significance * 2 // Phrases weighted more
	}
	
	totalScore := (wordScore + phraseScore) / 10.0
	return math.Min(totalScore, 1.0)
}

func calculatePatternDiversity(analysis *PatternAnalysis) float64 {
	// Count different pattern types found
	diversity := 0.0
	
	if len(analysis.WordRepetitions) > 0 {
		diversity += 0.15
	}
	if len(analysis.PhraseRepetitions) > 0 {
		diversity += 0.15
	}
	if len(analysis.StructureRepetitions) > 0 {
		diversity += 0.15
	}
	if len(analysis.SyntacticPatterns) > 0 {
		diversity += 0.15
	}
	if len(analysis.RhetoricalPatterns) > 0 {
		diversity += 0.15
	}
	if len(analysis.SequencePatterns) > 0 {
		diversity += 0.15
	}
	if len(analysis.TemporalPatterns) > 0 {
		diversity += 0.10
	}
	
	return diversity
}

func calculatePredictability(text string, analysis *PatternAnalysis) float64 {
	// Based on pattern regularity and repetition
	sentences := SplitIntoSentences(text)
	if len(sentences) == 0 {
		return 0
	}
	
	// Check structural regularity
	structureScore := 0.0
	if len(analysis.StructureRepetitions) > 0 {
		maxCount := 0
		for _, pattern := range analysis.StructureRepetitions {
			if pattern.Count > maxCount {
				maxCount = pattern.Count
			}
		}
		structureScore = float64(maxCount) / float64(len(sentences))
	}
	
	// Check repetition frequency
	repetitionScore := analysis.RepetitionScore
	
	// Combine scores
	predictability := (structureScore + repetitionScore) / 2
	return math.Min(predictability, 1.0)
}

func containsAny(list []string, items []string) bool {
	for _, item := range items {
		for _, l := range list {
			if l == item {
				return true
			}
		}
	}
	return false
}

func calculateMeanStdDevFloat(values []float64) (mean, stdDev float64) {
	if len(values) == 0 {
		return 0, 0
	}
	
	sum := 0.0
	for _, v := range values {
		sum += v
	}
	mean = sum / float64(len(values))
	
	variance := 0.0
	for _, v := range values {
		diff := v - mean
		variance += diff * diff
	}
	variance /= float64(len(values))
	stdDev = math.Sqrt(variance)
	
	return mean, stdDev
}

// Helper functions for pattern analysis

func isLikelyVerb(word string) bool {
	word = strings.ToLower(strings.Trim(word, ".,!?;:"))
	
	// Common verb patterns and endings
	verbPatterns := []string{
		"ed", "ing", "s", "es", "ies", "ied",
	}
	
	for _, pattern := range verbPatterns {
		if strings.HasSuffix(word, pattern) {
			return true
		}
	}
	
	// Common irregular verbs
	irregularVerbs := []string{
		"is", "are", "was", "were", "be", "been", "being",
		"have", "has", "had", "do", "does", "did", "will",
		"would", "could", "should", "can", "may", "might",
		"go", "went", "come", "came", "see", "saw", "get", "got",
		"make", "made", "take", "took", "give", "gave", "know", "knew",
	}
	
	for _, verb := range irregularVerbs {
		if word == verb {
			return true
		}
	}
	
	return false
}

func isConjunction(word string) bool {
	word = strings.ToLower(strings.Trim(word, ".,!?;:"))
	
	conjunctions := []string{
		// Coordinating conjunctions
		"and", "but", "or", "nor", "for", "yet", "so",
		// Subordinating conjunctions
		"because", "since", "although", "though", "while", "when",
		"where", "if", "unless", "until", "after", "before", "as",
		"that", "which", "who", "whom", "whose",
		// Correlative conjunctions
		"either", "neither", "both", "not", "whether",
	}
	
	for _, conj := range conjunctions {
		if word == conj {
			return true
		}
	}
	
	return false
}