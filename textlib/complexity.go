package textlib

import (
	"math"
	"sort"
	"strings"
	"unicode"
)

// ComplexityAnalysis contains various complexity metrics
type ComplexityAnalysis struct {
	// Readability scores
	FleschReadingEase    float64
	FleschKincaidGrade   float64
	GunningFog           float64
	ColemanLiau          float64
	SMOG                 float64
	ARI                  float64 // Automated Readability Index
	LIX                  float64 // Läsbarhetsindex
	
	// Complexity components
	SyntacticComplexity  float64
	SemanticComplexity   float64
	StructuralComplexity float64
	
	// Detailed metrics
	ClausePerSentence    float64
	PassiveVoiceRatio    float64
	NominalizationRatio  float64
	SubordinationIndex   float64
	
	// Complexity classification
	OverallComplexity    string // "very_easy", "easy", "moderate", "difficult", "very_difficult"
	TargetAudience       string // "elementary", "middle_school", "high_school", "college", "graduate", "professional"
	EstimatedReadingTime int    // in seconds
}

// CoherenceAnalysis measures text coherence and flow
type CoherenceAnalysis struct {
	// Cohesion metrics
	LexicalCohesion      float64 // Word repetition and semantic relatedness
	ReferentialCohesion  float64 // Pronoun usage and entity consistency
	ConjunctiveCohesion  float64 // Connector usage
	
	// Flow metrics
	TopicContinuity      float64 // How well topics flow
	SentenceTransitions  float64 // Quality of sentence-to-sentence transitions
	ParagraphTransitions float64 // Quality of paragraph transitions
	
	// Structure metrics
	ArgumentStructure    float64 // Logical flow of arguments
	ThematicConsistency  float64 // Consistency of themes
	
	// Issues detected
	CoherenceIssues      []CoherenceIssue
	TransitionProblems   []TransitionProblem
	
	// Overall scores
	OverallCoherence     float64
	CoherenceRating      string // "excellent", "good", "fair", "poor"
}

type CoherenceIssue struct {
	Type        string // "topic_shift", "unclear_reference", "missing_transition"
	Location    Position
	Description string
	Severity    string // "minor", "moderate", "major"
}

type TransitionProblem struct {
	FromSentence int
	ToSentence   int
	Issue        string
	Suggestion   string
}

// AnalyzeComplexity performs comprehensive complexity analysis
func AnalyzeComplexity(text string) *ComplexityAnalysis {
	analysis := &ComplexityAnalysis{}
	
	// Get basic statistics
	sentences := SplitIntoSentences(text)
	words := extractWords(text)
	
	if len(sentences) == 0 || len(words) == 0 {
		return analysis
	}
	
	// Calculate readability scores
	analysis.FleschReadingEase = CalculateFleschReadingEase(text)
	analysis.FleschKincaidGrade = calculateFleschKincaidGrade(text, sentences, words)
	analysis.GunningFog = calculateGunningFog(text, sentences, words)
	analysis.ColemanLiau = calculateColemanLiau(text, sentences, words)
	analysis.SMOG = calculateSMOG(sentences, words)
	analysis.ARI = calculateARI(text, sentences, words)
	analysis.LIX = calculateLIX(text, sentences, words)
	
	// Calculate complexity components
	analysis.SyntacticComplexity = calculateSyntacticComplexity(sentences)
	analysis.SemanticComplexity = calculateSemanticComplexity(words)
	analysis.StructuralComplexity = calculateStructuralComplexity(text)
	
	// Detailed metrics
	analysis.ClausePerSentence = calculateClausePerSentence(sentences)
	analysis.PassiveVoiceRatio = calculatePassiveVoiceRatio(sentences)
	analysis.NominalizationRatio = calculateNominalizationRatio(words)
	analysis.SubordinationIndex = calculateSubordinationIndex(sentences)
	
	// Classify complexity
	analysis.OverallComplexity = classifyComplexity(analysis.FleschReadingEase)
	analysis.TargetAudience = determineTargetAudience(analysis.FleschKincaidGrade)
	analysis.EstimatedReadingTime = estimateReadingTime(words, analysis.FleschReadingEase)
	
	return analysis
}

// Readability formulas

func calculateFleschKincaidGrade(text string, sentences, words []string) float64 {
	if len(sentences) == 0 || len(words) == 0 {
		return 0
	}
	
	totalSyllables := 0
	for _, word := range words {
		totalSyllables += CountSyllables(word)
	}
	
	avgWordsPerSentence := float64(len(words)) / float64(len(sentences))
	avgSyllablesPerWord := float64(totalSyllables) / float64(len(words))
	
	grade := 0.39*avgWordsPerSentence + 11.8*avgSyllablesPerWord - 15.59
	return math.Max(0, grade)
}

func calculateGunningFog(text string, sentences, words []string) float64 {
	if len(sentences) == 0 || len(words) == 0 {
		return 0
	}
	
	complexWords := 0
	for _, word := range words {
		if CountSyllables(word) >= 3 && !isCommonSuffix(word) {
			complexWords++
		}
	}
	
	avgWordsPerSentence := float64(len(words)) / float64(len(sentences))
	complexWordRatio := float64(complexWords) / float64(len(words))
	
	return 0.4 * (avgWordsPerSentence + 100*complexWordRatio)
}

func calculateColemanLiau(text string, sentences, words []string) float64 {
	if len(words) == 0 {
		return 0
	}
	
	letters := 0
	for _, word := range words {
		for _, r := range word {
			if unicode.IsLetter(r) {
				letters++
			}
		}
	}
	
	L := float64(letters) / float64(len(words)) * 100      // avg letters per 100 words
	S := float64(len(sentences)) / float64(len(words)) * 100 // avg sentences per 100 words
	
	return 0.0588*L - 0.296*S - 15.8
}

func calculateSMOG(sentences, words []string) float64 {
	if len(sentences) < 3 {
		return 0
	}
	
	// Sample sentences (first, middle, last)
	sampleSize := min(len(sentences), 30)
	sampleSentences := sentences[:sampleSize]
	
	polysyllables := 0
	for _, sentence := range sampleSentences {
		sentWords := strings.Fields(sentence)
		for _, word := range sentWords {
			if CountSyllables(word) >= 3 {
				polysyllables++
			}
		}
	}
	
	return 1.0430*math.Sqrt(float64(polysyllables)*30/float64(sampleSize)) + 3.1291
}

func calculateARI(text string, sentences, words []string) float64 {
	if len(sentences) == 0 || len(words) == 0 {
		return 0
	}
	
	chars := 0
	for _, word := range words {
		chars += len(word)
	}
	
	avgCharsPerWord := float64(chars) / float64(len(words))
	avgWordsPerSentence := float64(len(words)) / float64(len(sentences))
	
	return 4.71*avgCharsPerWord + 0.5*avgWordsPerSentence - 21.43
}

func calculateLIX(text string, sentences, words []string) float64 {
	if len(sentences) == 0 || len(words) == 0 {
		return 0
	}
	
	longWords := 0
	for _, word := range words {
		if len(word) > 6 {
			longWords++
		}
	}
	
	avgWordsPerSentence := float64(len(words)) / float64(len(sentences))
	longWordRatio := float64(longWords) / float64(len(words)) * 100
	
	return avgWordsPerSentence + longWordRatio
}

// Complexity component calculations

func calculateSyntacticComplexity(sentences []string) float64 {
	if len(sentences) == 0 {
		return 0
	}
	
	totalComplexity := 0.0
	
	for _, sentence := range sentences {
		// Count subordinate clauses
		subordinates := countSubordinateClauses(sentence)
		
		// Count coordinate clauses
		coordinates := strings.Count(sentence, " and ") + strings.Count(sentence, " but ") + 
					   strings.Count(sentence, " or ")
		
		// Measure embedding depth
		embeddingDepth := calculateEmbeddingDepth(sentence)
		
		// Sentence complexity score
		complexity := float64(subordinates)*2 + float64(coordinates) + float64(embeddingDepth)*1.5
		totalComplexity += complexity
	}
	
	return totalComplexity / float64(len(sentences))
}

func calculateSemanticComplexity(words []string) float64 {
	if len(words) == 0 {
		return 0
	}
	
	// Factors: abstract words, technical terms, rare words
	abstractCount := 0
	technicalCount := 0
	rareCount := 0
	
	wordFreq := make(map[string]int)
	for _, word := range words {
		wordFreq[strings.ToLower(word)]++
	}
	
	for word, freq := range wordFreq {
		if isAbstractWord(word) {
			abstractCount++
		}
		if isTechnicalTerm(word) {
			technicalCount++
		}
		if freq == 1 && len(word) > 7 { // Rare long words
			rareCount++
		}
	}
	
	uniqueWords := len(wordFreq)
	abstractRatio := float64(abstractCount) / float64(uniqueWords)
	technicalRatio := float64(technicalCount) / float64(uniqueWords)
	rareRatio := float64(rareCount) / float64(uniqueWords)
	
	return abstractRatio*0.4 + technicalRatio*0.4 + rareRatio*0.2
}

func calculateStructuralComplexity(text string) float64 {
	paragraphs := SplitIntoParagraphs(text)
	if len(paragraphs) == 0 {
		return 0
	}
	
	// Measure paragraph length variation
	lengths := make([]int, len(paragraphs))
	for i, para := range paragraphs {
		lengths[i] = len(strings.Fields(para))
	}
	
	_, stdDev := calculateMeanStdDev(lengths)
	
	// Measure discourse markers
	discourseMarkers := countDiscourseMarkers(text)
	markerDensity := float64(discourseMarkers) / float64(len(strings.Fields(text)))
	
	// Combine metrics
	return stdDev/100 + markerDensity*10
}

// Helper functions for complexity

func countSubordinateClauses(sentence string) int {
	subordinators := []string{
		" because ", " although ", " though ", " while ", " when ",
		" if ", " unless ", " since ", " after ", " before ",
		" that ", " which ", " who ", " whom ", " whose ",
	}
	
	count := 0
	lower := strings.ToLower(sentence)
	for _, sub := range subordinators {
		count += strings.Count(lower, sub)
	}
	
	return count
}

func calculateEmbeddingDepth(sentence string) int {
	// Simple approximation using parentheses and commas
	maxDepth := 0
	currentDepth := 0
	
	for _, r := range sentence {
		switch r {
		case '(', '[':
			currentDepth++
			if currentDepth > maxDepth {
				maxDepth = currentDepth
			}
		case ')', ']':
			currentDepth--
		}
	}
	
	// Also consider comma-separated clauses
	commaDepth := strings.Count(sentence, ",") / 2
	
	return max(maxDepth, commaDepth)
}

func isCommonSuffix(word string) bool {
	commonSuffixes := []string{"ing", "ed", "es", "ly", "er", "est"}
	lower := strings.ToLower(word)
	
	for _, suffix := range commonSuffixes {
		if strings.HasSuffix(lower, suffix) {
			return true
		}
	}
	return false
}

func isAbstractWord(word string) bool {
	abstractWords := map[string]bool{
		"concept": true, "theory": true, "principle": true, "idea": true,
		"notion": true, "abstraction": true, "philosophy": true, "belief": true,
		"understanding": true, "knowledge": true, "wisdom": true, "truth": true,
		"justice": true, "freedom": true, "democracy": true, "ethics": true,
		"morality": true, "virtue": true, "consciousness": true, "awareness": true,
	}
	
	return abstractWords[strings.ToLower(word)]
}

func isTechnicalTerm(word string) bool {
	// Simple heuristic: long words with specific patterns
	if len(word) < 8 {
		return false
	}
	
	techPatterns := []string{
		"tion", "ization", "ology", "ometry", "graphy",
		"scopy", "metry", "lysis", "synthesis",
	}
	
	lower := strings.ToLower(word)
	for _, pattern := range techPatterns {
		if strings.Contains(lower, pattern) {
			return true
		}
	}
	
	return false
}

func countDiscourseMarkers(text string) int {
	markers := []string{
		"however", "therefore", "moreover", "furthermore",
		"nevertheless", "consequently", "additionally", "specifically",
		"in contrast", "on the other hand", "for example", "in conclusion",
		"first", "second", "finally", "in summary",
	}
	
	count := 0
	lower := strings.ToLower(text)
	for _, marker := range markers {
		count += strings.Count(lower, marker)
	}
	
	return count
}

// Analysis helper functions

func calculateClausePerSentence(sentences []string) float64 {
	if len(sentences) == 0 {
		return 0
	}
	
	totalClauses := 0
	for _, sentence := range sentences {
		// Count main clause (1) plus subordinate clauses
		clauses := 1 + countSubordinateClauses(sentence)
		totalClauses += clauses
	}
	
	return float64(totalClauses) / float64(len(sentences))
}

func calculatePassiveVoiceRatio(sentences []string) float64 {
	if len(sentences) == 0 {
		return 0
	}
	
	passiveCount := 0
	for _, sentence := range sentences {
		if hasPassiveVoice(sentence) {
			passiveCount++
		}
	}
	
	return float64(passiveCount) / float64(len(sentences))
}

func hasPassiveVoice(sentence string) bool {
	// Simple heuristic: look for "be" verbs + past participle
	beVerbs := []string{" is ", " are ", " was ", " were ", " been ", " being ", " be "}
	lower := strings.ToLower(sentence)
	
	for _, be := range beVerbs {
		if strings.Contains(lower, be) {
			// Check if followed by past participle (ends in -ed, -en, etc.)
			index := strings.Index(lower, be)
			after := lower[index+len(be):]
			words := strings.Fields(after)
			
			if len(words) > 0 {
				firstWord := words[0]
				if strings.HasSuffix(firstWord, "ed") || strings.HasSuffix(firstWord, "en") ||
				   isIrregularPastParticiple(firstWord) {
					return true
				}
			}
		}
	}
	
	return false
}

func isIrregularPastParticiple(word string) bool {
	irregular := map[string]bool{
		"done": true, "gone": true, "seen": true, "taken": true,
		"given": true, "written": true, "spoken": true, "broken": true,
		"chosen": true, "driven": true, "eaten": true, "fallen": true,
		"forgotten": true, "gotten": true, "known": true, "shown": true,
		"thrown": true, "worn": true, "torn": true, "born": true,
	}
	
	return irregular[word]
}

func calculateNominalizationRatio(words []string) float64 {
	if len(words) == 0 {
		return 0
	}
	
	nominalizationCount := 0
	for _, word := range words {
		if isNominalization(word) {
			nominalizationCount++
		}
	}
	
	return float64(nominalizationCount) / float64(len(words))
}

func isNominalization(word string) bool {
	// Words ending in -tion, -ment, -ance, -ence, etc.
	nominalSuffixes := []string{
		"tion", "sion", "ment", "ance", "ence", "ness",
		"ity", "ism", "ship", "hood", "dom",
	}
	
	lower := strings.ToLower(word)
	for _, suffix := range nominalSuffixes {
		if strings.HasSuffix(lower, suffix) && len(word) > len(suffix)+3 {
			return true
		}
	}
	
	return false
}

func calculateSubordinationIndex(sentences []string) float64 {
	if len(sentences) == 0 {
		return 0
	}
	
	totalIndex := 0.0
	for _, sentence := range sentences {
		subordinates := countSubordinateClauses(sentence)
		words := len(strings.Fields(sentence))
		
		if words > 0 {
			index := float64(subordinates) / float64(words) * 100
			totalIndex += index
		}
	}
	
	return totalIndex / float64(len(sentences))
}

func classifyComplexity(fleschScore float64) string {
	switch {
	case fleschScore >= 90:
		return "very_easy"
	case fleschScore >= 70:
		return "easy"
	case fleschScore >= 50:
		return "moderate"
	case fleschScore >= 30:
		return "difficult"
	default:
		return "very_difficult"
	}
}

func determineTargetAudience(gradeLevel float64) string {
	switch {
	case gradeLevel < 6:
		return "elementary"
	case gradeLevel < 9:
		return "middle_school"
	case gradeLevel < 13:
		return "high_school"
	case gradeLevel < 16:
		return "college"
	case gradeLevel < 18:
		return "graduate"
	default:
		return "professional"
	}
}

func estimateReadingTime(words []string, fleschScore float64) int {
	// Base reading speed (words per minute)
	baseSpeed := 250.0
	
	// Adjust based on complexity
	complexityFactor := (100 - fleschScore) / 100
	adjustedSpeed := baseSpeed * (1 - complexityFactor*0.5)
	
	// Calculate time in seconds
	minutes := float64(len(words)) / adjustedSpeed
	seconds := int(minutes * 60)
	
	return seconds
}

// AnalyzeCoherence performs comprehensive coherence analysis
func AnalyzeCoherence(text string) *CoherenceAnalysis {
	analysis := &CoherenceAnalysis{
		CoherenceIssues:    []CoherenceIssue{},
		TransitionProblems: []TransitionProblem{},
	}
	
	sentences := SplitIntoSentences(text)
	paragraphs := SplitIntoParagraphs(text)
	
	if len(sentences) == 0 {
		return analysis
	}
	
	// Calculate cohesion metrics
	analysis.LexicalCohesion = calculateLexicalCohesion(sentences)
	analysis.ReferentialCohesion = calculateReferentialCohesion(text)
	analysis.ConjunctiveCohesion = calculateConjunctiveCohesion(sentences)
	
	// Calculate flow metrics
	analysis.TopicContinuity = calculateTopicContinuity(sentences)
	analysis.SentenceTransitions = calculateSentenceTransitions(sentences)
	analysis.ParagraphTransitions = calculateParagraphTransitions(paragraphs)
	
	// Calculate structure metrics
	analysis.ArgumentStructure = calculateArgumentStructure(text)
	analysis.ThematicConsistency = calculateThematicConsistency(paragraphs)
	
	// Detect issues
	analysis.CoherenceIssues = detectCoherenceIssues(text, sentences)
	analysis.TransitionProblems = detectTransitionProblems(sentences)
	
	// Calculate overall coherence
	analysis.OverallCoherence = (
		analysis.LexicalCohesion*0.2 +
		analysis.ReferentialCohesion*0.2 +
		analysis.ConjunctiveCohesion*0.1 +
		analysis.TopicContinuity*0.2 +
		analysis.SentenceTransitions*0.15 +
		analysis.ParagraphTransitions*0.15)
	
	// Rate coherence
	analysis.CoherenceRating = rateCoherence(analysis.OverallCoherence)
	
	return analysis
}

func calculateLexicalCohesion(sentences []string) float64 {
	if len(sentences) < 2 {
		return 1.0
	}
	
	totalCohesion := 0.0
	
	for i := 0; i < len(sentences)-1; i++ {
		words1 := extractContentWords(sentences[i])
		words2 := extractContentWords(sentences[i+1])
		
		// Calculate word overlap
		overlap := calculateWordOverlap(words1, words2)
		
		// Calculate semantic similarity (simplified)
		similarity := calculateSimpleSemantic(words1, words2)
		
		cohesion := overlap*0.6 + similarity*0.4
		totalCohesion += cohesion
	}
	
	return totalCohesion / float64(len(sentences)-1)
}

func extractContentWords(sentence string) []string {
	words := strings.Fields(strings.ToLower(sentence))
	var content []string
	
	for _, word := range words {
		cleaned := strings.Trim(word, ".,!?;:")
		if len(cleaned) > 3 && !isStopWord(cleaned) {
			content = append(content, cleaned)
		}
	}
	
	return content
}

func isStopWord(word string) bool {
	stopWords := map[string]bool{
		"the": true, "a": true, "an": true, "and": true, "or": true,
		"but": true, "in": true, "on": true, "at": true, "to": true,
		"for": true, "of": true, "with": true, "by": true, "from": true,
		"is": true, "are": true, "was": true, "were": true, "be": true,
		"been": true, "being": true, "have": true, "has": true, "had": true,
	}
	
	return stopWords[word]
}

func calculateWordOverlap(words1, words2 []string) float64 {
	if len(words1) == 0 || len(words2) == 0 {
		return 0
	}
	
	set1 := make(map[string]bool)
	for _, w := range words1 {
		set1[w] = true
	}
	
	overlap := 0
	for _, w := range words2 {
		if set1[w] {
			overlap++
		}
	}
	
	return float64(overlap) / float64(max(len(words1), len(words2)))
}

func calculateSimpleSemantic(words1, words2 []string) float64 {
	// Simplified semantic similarity based on word stems and synonyms
	stemOverlap := 0
	
	for _, w1 := range words1 {
		stem1 := extractStem(w1)
		for _, w2 := range words2 {
			stem2 := extractStem(w2)
			if stem1 == stem2 || areSynonyms(w1, w2) {
				stemOverlap++
			}
		}
	}
	
	if len(words1) == 0 || len(words2) == 0 {
		return 0
	}
	
	return float64(stemOverlap) / float64(max(len(words1), len(words2)))
}

func extractStem(word string) string {
	// Very simple stemming
	suffixes := []string{"ing", "ed", "es", "s", "ly", "er", "est", "tion", "ment"}
	
	for _, suffix := range suffixes {
		if strings.HasSuffix(word, suffix) && len(word)-len(suffix) > 3 {
			return word[:len(word)-len(suffix)]
		}
	}
	
	return word
}

func areSynonyms(w1, w2 string) bool {
	// Simple synonym groups
	synonymGroups := [][]string{
		{"big", "large", "huge", "enormous"},
		{"small", "little", "tiny", "minute"},
		{"fast", "quick", "rapid", "speedy"},
		{"slow", "sluggish", "gradual"},
		{"good", "great", "excellent", "fine"},
		{"bad", "poor", "terrible", "awful"},
	}
	
	for _, group := range synonymGroups {
		inGroup1, inGroup2 := false, false
		for _, syn := range group {
			if syn == w1 {
				inGroup1 = true
			}
			if syn == w2 {
				inGroup2 = true
			}
		}
		if inGroup1 && inGroup2 {
			return true
		}
	}
	
	return false
}

func calculateReferentialCohesion(text string) float64 {
	// Simplified referential cohesion based on pronoun usage
	pronouns := extractPronouns(text)
	if len(pronouns) == 0 {
		return 1.0 // No pronouns means perfect cohesion
	}
	
	// Simple heuristic: fewer pronouns relative to nouns indicates better cohesion
	words := strings.Fields(text)
	nouns := 0
	for _, word := range words {
		if isNounBasic(word) {
			nouns++
		}
	}
	
	if nouns == 0 {
		return 0.5 // Neutral score if no nouns found
	}
	
	// Ratio of pronouns to nouns (lower is better for clarity)
	ratio := float64(len(pronouns)) / float64(nouns)
	cohesion := 1.0 - (ratio * 0.5) // Scale down the impact
	
	if cohesion < 0 {
		cohesion = 0
	}
	
	return cohesion
}

func isNounBasic(word string) bool {
	// Simple noun detection heuristic
	word = strings.ToLower(strings.Trim(word, ".,!?;:"))
	
	// Common noun endings
	nounSuffixes := []string{"tion", "sion", "ness", "ment", "ity", "er", "or"}
	for _, suffix := range nounSuffixes {
		if strings.HasSuffix(word, suffix) {
			return true
		}
	}
	
	// Capitalized words (proper nouns)
	if len(word) > 0 && unicode.IsUpper(rune(word[0])) {
		return true
	}
	
	return false
}

func extractPronouns(text string) []string {
	var pronouns []string
	pronounSet := map[string]bool{
		"he": true, "she": true, "it": true, "they": true,
		"him": true, "her": true, "them": true, "their": true,
		"his": true, "hers": true, "its": true, "theirs": true,
	}
	
	words := strings.Fields(text)
	for _, word := range words {
		cleaned := strings.ToLower(strings.Trim(word, ".,!?;:"))
		if pronounSet[cleaned] {
			pronouns = append(pronouns, word)
		}
	}
	
	return pronouns
}

func calculateConjunctiveCohesion(sentences []string) float64 {
	if len(sentences) < 2 {
		return 1.0
	}
	
	connectives := []string{
		"however", "therefore", "moreover", "furthermore",
		"nevertheless", "consequently", "additionally", "meanwhile",
		"subsequently", "alternatively", "similarly", "likewise",
		"in contrast", "on the other hand", "for example", "in addition",
	}
	
	connectiveCount := 0
	for _, sentence := range sentences {
		lower := strings.ToLower(sentence)
		for _, conn := range connectives {
			if strings.Contains(lower, conn) {
				connectiveCount++
				break
			}
		}
	}
	
	// Ideal ratio is about 0.3-0.5 connectives per sentence
	ratio := float64(connectiveCount) / float64(len(sentences))
	if ratio > 0.5 {
		ratio = 1.0 - (ratio - 0.5) // Penalize overuse
	}
	
	return minFloat(ratio*2, 1.0) // Scale to 0-1
}

func calculateTopicContinuity(sentences []string) float64 {
	if len(sentences) < 2 {
		return 1.0
	}
	
	// Track key topics through sentences
	topicChains := make(map[string][]int)
	
	for i, sentence := range sentences {
		topics := extractTopics(sentence)
		for _, topic := range topics {
			topicChains[topic] = append(topicChains[topic], i)
		}
	}
	
	// Calculate continuity score
	continuityScore := 0.0
	chainCount := 0
	
	for _, positions := range topicChains {
		if len(positions) > 1 {
			// Check if mentions are reasonably close
			maxGap := 0
			for i := 1; i < len(positions); i++ {
				gap := positions[i] - positions[i-1]
				if gap > maxGap {
					maxGap = gap
				}
			}
			
			// Score based on gap size (smaller gaps = better continuity)
			if maxGap <= 3 {
				continuityScore += 1.0
			} else if maxGap <= 5 {
				continuityScore += 0.7
			} else {
				continuityScore += 0.4
			}
			chainCount++
		}
	}
	
	if chainCount == 0 {
		return 0.5 // No topic chains
	}
	
	return continuityScore / float64(chainCount)
}

func extractTopics(sentence string) []string {
	// Extract noun phrases as potential topics
	var topics []string
	words := strings.Fields(sentence)
	
	for _, word := range words {
		cleaned := strings.Trim(word, ".,!?;:")
		// Simple heuristic: capitalized words (except sentence start) and long words
		if len(cleaned) > 5 && unicode.IsUpper([]rune(cleaned)[0]) {
			topics = append(topics, strings.ToLower(cleaned))
		}
	}
	
	return topics
}

func calculateSentenceTransitions(sentences []string) float64 {
	if len(sentences) < 2 {
		return 1.0
	}
	
	goodTransitions := 0
	
	for i := 0; i < len(sentences)-1; i++ {
		quality := evaluateTransition(sentences[i], sentences[i+1])
		if quality > 0.5 {
			goodTransitions++
		}
	}
	
	return float64(goodTransitions) / float64(len(sentences)-1)
}

func evaluateTransition(sent1, sent2 string) float64 {
	// Check for explicit connectives
	connectives := []string{
		"however", "therefore", "moreover", "furthermore",
		"additionally", "also", "similarly", "likewise",
	}
	
	sent2Lower := strings.ToLower(sent2)
	for _, conn := range connectives {
		if strings.HasPrefix(sent2Lower, conn) {
			return 0.9
		}
	}
	
	// Check for pronoun references
	if hasPronouns(sent2) && hasEntities(sent1) {
		return 0.8
	}
	
	// Check for lexical overlap
	overlap := calculateWordOverlap(
		extractContentWords(sent1),
		extractContentWords(sent2),
	)
	
	return overlap
}

func hasPronouns(sentence string) bool {
	pronouns := []string{"he", "she", "it", "they", "this", "that", "these", "those"}
	lower := strings.ToLower(sentence)
	
	for _, p := range pronouns {
		if strings.Contains(lower, " "+p+" ") {
			return true
		}
	}
	return false
}

func hasEntities(sentence string) bool {
	// Check for capitalized words (potential entities)
	words := strings.Fields(sentence)
	for i, word := range words {
		if i > 0 && len(word) > 0 && unicode.IsUpper([]rune(word)[0]) {
			return true
		}
	}
	return false
}

func calculateParagraphTransitions(paragraphs []string) float64 {
	if len(paragraphs) < 2 {
		return 1.0
	}
	
	goodTransitions := 0
	
	for i := 0; i < len(paragraphs)-1; i++ {
		// Get last sentence of paragraph and first of next
		sent1 := getLastSentence(paragraphs[i])
		sent2 := getFirstSentence(paragraphs[i+1])
		
		quality := evaluateTransition(sent1, sent2)
		if quality > 0.4 {
			goodTransitions++
		}
	}
	
	return float64(goodTransitions) / float64(len(paragraphs)-1)
}

func getLastSentence(paragraph string) string {
	sentences := SplitIntoSentences(paragraph)
	if len(sentences) > 0 {
		return sentences[len(sentences)-1]
	}
	return ""
}

func getFirstSentence(paragraph string) string {
	sentences := SplitIntoSentences(paragraph)
	if len(sentences) > 0 {
		return sentences[0]
	}
	return ""
}

func calculateArgumentStructure(text string) float64 {
	// Look for argument markers
	argumentMarkers := []string{
		"first", "second", "third", "finally",
		"in conclusion", "to summarize", "in summary",
		"for example", "for instance", "specifically",
		"on the contrary", "in contrast", "however",
		"therefore", "thus", "consequently", "as a result",
	}
	
	lower := strings.ToLower(text)
	markerCount := 0
	
	for _, marker := range argumentMarkers {
		markerCount += strings.Count(lower, marker)
	}
	
	words := strings.Fields(text)
	if len(words) == 0 {
		return 0
	}
	
	// Ideal density is about 1 marker per 100 words
	density := float64(markerCount) / float64(len(words)) * 100
	
	if density > 2 {
		density = 4 - density // Penalize overuse
	}
	
	return maxFloat(0, minFloat(1, density/2))
}

func calculateThematicConsistency(paragraphs []string) float64 {
	if len(paragraphs) == 0 {
		return 1.0
	}
	
	// Extract themes from each paragraph
	themes := make([][]string, len(paragraphs))
	for i, para := range paragraphs {
		themes[i] = extractThemes(para)
	}
	
	// Calculate consistency
	consistencyScore := 0.0
	comparisons := 0
	
	for i := 0; i < len(themes); i++ {
		for j := i + 1; j < len(themes); j++ {
			similarity := calculateThemeSimilarity(themes[i], themes[j])
			// Weight by distance (closer paragraphs should be more similar)
			weight := 1.0 / float64(j-i)
			consistencyScore += similarity * weight
			comparisons++
		}
	}
	
	if comparisons == 0 {
		return 1.0
	}
	
	return consistencyScore / float64(comparisons)
}

func extractThemes(paragraph string) []string {
	// Extract key content words as themes
	words := extractContentWords(paragraph)
	
	// Count frequencies
	freq := make(map[string]int)
	for _, word := range words {
		freq[word]++
	}
	
	// Get top words as themes
	type wordCount struct {
		word  string
		count int
	}
	
	var counts []wordCount
	for w, c := range freq {
		if c > 1 { // Only words that appear multiple times
			counts = append(counts, wordCount{w, c})
		}
	}
	
	sort.Slice(counts, func(i, j int) bool {
		return counts[i].count > counts[j].count
	})
	
	var themes []string
	for i := 0; i < len(counts) && i < 5; i++ {
		themes = append(themes, counts[i].word)
	}
	
	return themes
}

func calculateThemeSimilarity(themes1, themes2 []string) float64 {
	if len(themes1) == 0 || len(themes2) == 0 {
		return 0
	}
	
	set1 := make(map[string]bool)
	for _, t := range themes1 {
		set1[t] = true
	}
	
	overlap := 0
	for _, t := range themes2 {
		if set1[t] {
			overlap++
		}
	}
	
	return float64(overlap) / float64(max(len(themes1), len(themes2)))
}

func detectCoherenceIssues(text string, sentences []string) []CoherenceIssue {
	var issues []CoherenceIssue
	
	// Check for unclear pronoun references
	for i, sentence := range sentences {
		if hasUnclearPronouns(sentence, i, sentences) {
			issues = append(issues, CoherenceIssue{
				Type:        "unclear_reference",
				Location:    Position{Start: i, End: i},
				Description: "Pronoun without clear antecedent",
				Severity:    "moderate",
			})
		}
	}
	
	// Check for abrupt topic shifts
	for i := 1; i < len(sentences); i++ {
		if hasAbruptTopicShift(sentences[i-1], sentences[i]) {
			issues = append(issues, CoherenceIssue{
				Type:        "topic_shift",
				Location:    Position{Start: i-1, End: i},
				Description: "Abrupt topic change between sentences",
				Severity:    "minor",
			})
		}
	}
	
	return issues
}

func hasUnclearPronouns(sentence string, index int, allSentences []string) bool {
	pronouns := []string{"it", "this", "that", "they"}
	lower := strings.ToLower(sentence)
	
	// Check if sentence starts with a pronoun
	for _, p := range pronouns {
		if strings.HasPrefix(lower, p+" ") {
			// Check if previous sentence has a clear referent
			if index == 0 || !hasEntities(allSentences[index-1]) {
				return true
			}
		}
	}
	
	return false
}

func hasAbruptTopicShift(sent1, sent2 string) bool {
	topics1 := extractTopics(sent1)
	topics2 := extractTopics(sent2)
	
	// No overlap and no transition words
	if len(topics1) > 0 && len(topics2) > 0 {
		overlap := false
		for _, t1 := range topics1 {
			for _, t2 := range topics2 {
				if t1 == t2 {
					overlap = true
					break
				}
			}
		}
		
		if !overlap && !hasTransitionWords(sent2) {
			return true
		}
	}
	
	return false
}

func hasTransitionWords(sentence string) bool {
	transitions := []string{
		"however", "therefore", "moreover", "furthermore",
		"additionally", "meanwhile", "consequently",
	}
	
	lower := strings.ToLower(sentence)
	for _, trans := range transitions {
		if strings.Contains(lower, trans) {
			return true
		}
	}
	
	return false
}

func detectTransitionProblems(sentences []string) []TransitionProblem {
	var problems []TransitionProblem
	
	for i := 0; i < len(sentences)-1; i++ {
		quality := evaluateTransition(sentences[i], sentences[i+1])
		
		if quality < 0.3 {
			problem := TransitionProblem{
				FromSentence: i,
				ToSentence:   i + 1,
				Issue:        "Weak transition between sentences",
			}
			
			// Suggest improvement
			if !hasTransitionWords(sentences[i+1]) {
				problem.Suggestion = "Consider adding a transition word or phrase"
			} else {
				problem.Suggestion = "Consider rephrasing to improve flow"
			}
			
			problems = append(problems, problem)
		}
	}
	
	return problems
}

func rateCoherence(score float64) string {
	switch {
	case score >= 0.8:
		return "excellent"
	case score >= 0.6:
		return "good"
	case score >= 0.4:
		return "fair"
	default:
		return "poor"
	}
}

// Utility functions
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

func minFloat(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}

func maxFloat(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}