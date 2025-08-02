package textlib

import (
	"math"
	"regexp"
	"sort"
	"strings"
	"unicode"
)

// TextStatistics holds comprehensive statistics about a text
type TextStatistics struct {
	// Basic counts
	CharacterCount       int
	CharacterCountNoSpace int
	WordCount           int
	UniqueWordCount     int
	SentenceCount       int
	ParagraphCount      int
	
	// Vocabulary metrics
	TypeTokenRatio      float64  // Unique words / total words
	HapaxLegomena       int      // Words appearing only once
	DisLegomena         int      // Words appearing exactly twice
	VocabularyRichness  float64  // Yule's K measure
	LexicalDiversity    float64  // MTLD score
	
	// Word length distribution
	AverageWordLength   float64
	WordLengthStdDev    float64
	ShortWords          int  // <= 3 chars
	MediumWords         int  // 4-7 chars
	LongWords           int  // 8-12 chars
	VeryLongWords       int  // > 12 chars
	
	// Sentence metrics
	AverageSentenceLength     float64
	SentenceLengthStdDev      float64
	ShortSentences           int  // < 10 words
	MediumSentences          int  // 10-25 words
	LongSentences            int  // 26-40 words
	VeryLongSentences        int  // > 40 words
	
	// Complexity indicators
	AverageSyllablesPerWord  float64
	ComplexWordCount         int      // Words with 3+ syllables
	ComplexWordRatio         float64
	
	// Punctuation analysis
	PunctuationCount         int
	PunctuationDensity       float64  // Punctuation per 100 words
	ExclamationCount         int
	QuestionCount            int
	CommaCount              int
	SemicolonCount          int
	
	// Capitalization patterns
	CapitalizedWords         int
	AllCapsWords            int
	TitleCaseWords          int
	
	// Most frequent items
	MostFrequentWords       []WordFrequency
	MostFrequentBigrams     []BigramFrequency
	MostFrequentTrigrams    []TrigramFrequency
}

type WordFrequency struct {
	Word      string
	Count     int
	Frequency float64
}

type BigramFrequency struct {
	Bigram    string
	Count     int
	Frequency float64
}

type TrigramFrequency struct {
	Trigram   string
	Count     int
	Frequency float64
}

// CalculateTextStatistics computes comprehensive statistics for a text
func CalculateTextStatistics(text string) *TextStatistics {
	stats := &TextStatistics{}
	
	// Basic counts
	stats.CharacterCount = len([]rune(text))
	stats.CharacterCountNoSpace = countNonSpaceChars(text)
	
	// Process sentences and words
	sentences := SplitIntoSentences(text)
	stats.SentenceCount = len(sentences)
	
	paragraphs := SplitIntoParagraphs(text)
	stats.ParagraphCount = len(paragraphs)
	
	// Word analysis
	words := extractWords(text)
	stats.WordCount = len(words)
	
	// Word frequency analysis
	wordFreq := make(map[string]int)
	for _, word := range words {
		lower := strings.ToLower(word)
		wordFreq[lower]++
	}
	
	stats.UniqueWordCount = len(wordFreq)
	
	// Calculate vocabulary metrics
	stats.TypeTokenRatio = float64(stats.UniqueWordCount) / float64(stats.WordCount)
	stats.HapaxLegomena = countWordsWithFrequency(wordFreq, 1)
	stats.DisLegomena = countWordsWithFrequency(wordFreq, 2)
	stats.VocabularyRichness = calculateYulesK(wordFreq, stats.WordCount)
	stats.LexicalDiversity = calculateMTLD(words)
	
	// Word length analysis
	stats.AverageWordLength, stats.WordLengthStdDev = calculateWordLengthStats(words)
	stats.ShortWords, stats.MediumWords, stats.LongWords, stats.VeryLongWords = classifyWordLengths(words)
	
	// Sentence length analysis
	sentenceLengths := make([]int, len(sentences))
	for i, sent := range sentences {
		sentenceLengths[i] = len(strings.Fields(sent))
	}
	stats.AverageSentenceLength, stats.SentenceLengthStdDev = calculateMeanStdDev(sentenceLengths)
	stats.ShortSentences, stats.MediumSentences, stats.LongSentences, stats.VeryLongSentences = classifySentenceLengths(sentenceLengths)
	
	// Complexity analysis
	totalSyllables := 0
	stats.ComplexWordCount = 0
	for _, word := range words {
		syllables := CountSyllables(word)
		totalSyllables += syllables
		if syllables >= 3 {
			stats.ComplexWordCount++
		}
	}
	stats.AverageSyllablesPerWord = float64(totalSyllables) / float64(stats.WordCount)
	stats.ComplexWordRatio = float64(stats.ComplexWordCount) / float64(stats.WordCount)
	
	// Punctuation analysis
	stats.PunctuationCount = countPunctuation(text)
	stats.PunctuationDensity = (float64(stats.PunctuationCount) / float64(stats.WordCount)) * 100
	stats.ExclamationCount = strings.Count(text, "!")
	stats.QuestionCount = strings.Count(text, "?")
	stats.CommaCount = strings.Count(text, ",")
	stats.SemicolonCount = strings.Count(text, ";")
	
	// Capitalization analysis
	stats.CapitalizedWords, stats.AllCapsWords, stats.TitleCaseWords = analyzeCapitalization(words)
	
	// Most frequent items
	stats.MostFrequentWords = findMostFrequentWords(wordFreq, 10)
	stats.MostFrequentBigrams = findMostFrequentBigrams(words, 10)
	stats.MostFrequentTrigrams = findMostFrequentTrigrams(words, 10)
	
	return stats
}

// Helper functions

func countNonSpaceChars(text string) int {
	count := 0
	for _, r := range text {
		if !unicode.IsSpace(r) {
			count++
		}
	}
	return count
}

func extractWords(text string) []string {
	var words []string
	fields := strings.Fields(text)
	for _, field := range fields {
		word := cleanWord(field)
		if word != "" {
			words = append(words, word)
		}
	}
	return words
}

func cleanWord(word string) string {
	// Remove punctuation but keep letters and numbers
	return regexp.MustCompile(`[^a-zA-Z0-9]`).ReplaceAllString(word, "")
}

func countWordsWithFrequency(freq map[string]int, target int) int {
	count := 0
	for _, f := range freq {
		if f == target {
			count++
		}
	}
	return count
}

func calculateYulesK(wordFreq map[string]int, totalWords int) float64 {
	// Yule's K = 10^4 * (M1 * M1 - M2) / (M1 * M1)
	// where M1 = total words, M2 = sum of squares of frequencies
	M1 := float64(totalWords)
	M2 := 0.0
	
	for _, freq := range wordFreq {
		M2 += float64(freq * freq)
	}
	
	if M1 == 0 {
		return 0
	}
	
	return 10000 * (M1*M1 - M2) / (M1 * M1)
}

func calculateMTLD(words []string) float64 {
	// Measure of Textual Lexical Diversity
	if len(words) < 10 {
		return float64(len(words))
	}
	
	forward := calculateMTLDForward(words, 0.72)
	backward := calculateMTLDBackward(words, 0.72)
	
	return (forward + backward) / 2.0
}

func calculateMTLDForward(words []string, threshold float64) float64 {
	factors := 0.0
	factor := 0.0
	types := make(map[string]bool)
	
	for _, word := range words {
		types[strings.ToLower(word)] = true
		factor++
		
		ttr := float64(len(types)) / factor
		if ttr <= threshold {
			factors++
			factor = 0.0
			types = make(map[string]bool)
		}
	}
	
	if factor > 0 {
		factors += factor / float64(len(words))
	}
	
	if factors == 0 {
		return float64(len(words))
	}
	
	return float64(len(words)) / factors
}

func calculateMTLDBackward(words []string, threshold float64) float64 {
	// Reverse the words and calculate
	reversed := make([]string, len(words))
	for i, word := range words {
		reversed[len(words)-1-i] = word
	}
	return calculateMTLDForward(reversed, threshold)
}

func calculateWordLengthStats(words []string) (mean, stdDev float64) {
	lengths := make([]int, len(words))
	for i, word := range words {
		lengths[i] = len([]rune(word))
	}
	return calculateMeanStdDev(lengths)
}

func calculateMeanStdDev(values []int) (mean, stdDev float64) {
	if len(values) == 0 {
		return 0, 0
	}
	
	sum := 0
	for _, v := range values {
		sum += v
	}
	mean = float64(sum) / float64(len(values))
	
	variance := 0.0
	for _, v := range values {
		diff := float64(v) - mean
		variance += diff * diff
	}
	variance /= float64(len(values))
	stdDev = math.Sqrt(variance)
	
	return mean, stdDev
}

func classifyWordLengths(words []string) (short, medium, long, veryLong int) {
	for _, word := range words {
		length := len([]rune(word))
		switch {
		case length <= 3:
			short++
		case length <= 7:
			medium++
		case length <= 12:
			long++
		default:
			veryLong++
		}
	}
	return
}

func classifySentenceLengths(lengths []int) (short, medium, long, veryLong int) {
	for _, length := range lengths {
		switch {
		case length < 10:
			short++
		case length <= 25:
			medium++
		case length <= 40:
			long++
		default:
			veryLong++
		}
	}
	return
}

func countPunctuation(text string) int {
	count := 0
	for _, r := range text {
		if unicode.IsPunct(r) {
			count++
		}
	}
	return count
}

func analyzeCapitalization(words []string) (capitalized, allCaps, titleCase int) {
	for _, word := range words {
		if len(word) == 0 {
			continue
		}
		
		runes := []rune(word)
		if unicode.IsUpper(runes[0]) {
			if isAllCaps(word) {
				allCaps++
			} else if isTitleCase(word) {
				titleCase++
			} else {
				capitalized++
			}
		}
	}
	return
}

func isAllCaps(word string) bool {
	for _, r := range word {
		if unicode.IsLetter(r) && !unicode.IsUpper(r) {
			return false
		}
	}
	return true
}

func isTitleCase(word string) bool {
	foundLower := false
	for i, r := range word {
		if i == 0 {
			if !unicode.IsUpper(r) {
				return false
			}
		} else if unicode.IsLetter(r) {
			if unicode.IsUpper(r) {
				return false
			}
			foundLower = true
		}
	}
	return foundLower
}

func findMostFrequentWords(wordFreq map[string]int, limit int) []WordFrequency {
	// Filter out common stop words
	stopWords := map[string]bool{
		"the": true, "a": true, "an": true, "and": true, "or": true,
		"but": true, "in": true, "on": true, "at": true, "to": true,
		"for": true, "of": true, "with": true, "by": true, "from": true,
		"is": true, "are": true, "was": true, "were": true, "be": true,
		"been": true, "being": true, "have": true, "has": true, "had": true,
		"do": true, "does": true, "did": true, "will": true, "would": true,
		"could": true, "should": true, "may": true, "might": true, "must": true,
		"shall": true, "can": true, "this": true, "that": true, "these": true,
		"those": true, "i": true, "you": true, "he": true, "she": true,
		"it": true, "we": true, "they": true, "them": true, "their": true,
	}
	
	var frequencies []WordFrequency
	totalWords := 0
	for _, count := range wordFreq {
		totalWords += count
	}
	
	for word, count := range wordFreq {
		if !stopWords[word] && len(word) > 2 {
			frequencies = append(frequencies, WordFrequency{
				Word:      word,
				Count:     count,
				Frequency: float64(count) / float64(totalWords),
			})
		}
	}
	
	sort.Slice(frequencies, func(i, j int) bool {
		return frequencies[i].Count > frequencies[j].Count
	})
	
	if len(frequencies) > limit {
		frequencies = frequencies[:limit]
	}
	
	return frequencies
}

func findMostFrequentBigrams(words []string, limit int) []BigramFrequency {
	bigramFreq := make(map[string]int)
	
	for i := 0; i < len(words)-1; i++ {
		bigram := strings.ToLower(words[i]) + " " + strings.ToLower(words[i+1])
		bigramFreq[bigram]++
	}
	
	var frequencies []BigramFrequency
	totalBigrams := len(words) - 1
	
	for bigram, count := range bigramFreq {
		if count > 1 { // Only include bigrams that appear more than once
			frequencies = append(frequencies, BigramFrequency{
				Bigram:    bigram,
				Count:     count,
				Frequency: float64(count) / float64(totalBigrams),
			})
		}
	}
	
	sort.Slice(frequencies, func(i, j int) bool {
		return frequencies[i].Count > frequencies[j].Count
	})
	
	if len(frequencies) > limit {
		frequencies = frequencies[:limit]
	}
	
	return frequencies
}

func findMostFrequentTrigrams(words []string, limit int) []TrigramFrequency {
	trigramFreq := make(map[string]int)
	
	for i := 0; i < len(words)-2; i++ {
		trigram := strings.ToLower(words[i]) + " " + 
				   strings.ToLower(words[i+1]) + " " + 
				   strings.ToLower(words[i+2])
		trigramFreq[trigram]++
	}
	
	var frequencies []TrigramFrequency
	totalTrigrams := len(words) - 2
	
	for trigram, count := range trigramFreq {
		if count > 1 { // Only include trigrams that appear more than once
			frequencies = append(frequencies, TrigramFrequency{
				Trigram:   trigram,
				Count:     count,
				Frequency: float64(count) / float64(totalTrigrams),
			})
		}
	}
	
	sort.Slice(frequencies, func(i, j int) bool {
		return frequencies[i].Count > frequencies[j].Count
	})
	
	if len(frequencies) > limit {
		frequencies = frequencies[:limit]
	}
	
	return frequencies
}