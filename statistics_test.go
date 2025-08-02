package textlib

import (
	"testing"
)

func TestCalculateTextStatistics(t *testing.T) {
	text := "The quick brown fox jumps over the lazy dog. This sentence contains every letter of the alphabet at least once. It's a pangram!"

	stats := CalculateTextStatistics(text)

	// Basic validations
	if stats.CharacterCount == 0 {
		t.Error("Character count should be greater than 0")
	}

	if stats.WordCount == 0 {
		t.Error("Word count should be greater than 0")
	}

	if stats.SentenceCount == 0 {
		t.Error("Sentence count should be greater than 0")
	}

	if stats.ParagraphCount == 0 {
		t.Error("Paragraph count should be greater than 0")
	}

	// Ratio validations
	if stats.TypeTokenRatio < 0 || stats.TypeTokenRatio > 1 {
		t.Errorf("Type-token ratio %f should be between 0 and 1", stats.TypeTokenRatio)
	}

	if stats.VocabularyRichness < 0 || stats.VocabularyRichness > 1 {
		t.Errorf("Vocabulary richness %f should be between 0 and 1", stats.VocabularyRichness)
	}

	if stats.LexicalDiversity < 0 {
		t.Errorf("Lexical diversity %f should be non-negative", stats.LexicalDiversity)
	}

	// Unique words should not exceed total words
	if stats.UniqueWordCount > stats.WordCount {
		t.Errorf("Unique words (%d) cannot exceed total words (%d)", 
			stats.UniqueWordCount, stats.WordCount)
	}

	// Most frequent words should be populated
	if len(stats.MostFrequentWords) == 0 {
		t.Error("Most frequent words should be populated")
	}

	// Test with empty text
	emptyStats := CalculateTextStatistics("")
	if emptyStats.CharacterCount != 0 || emptyStats.WordCount != 0 {
		t.Error("Empty text should have zero counts")
	}
}

func TestWordLengthDistribution(t *testing.T) {
	text := "I am happy. This is a longer sentence with various word lengths."
	stats := CalculateTextStatistics(text)

	// Check that word length categories sum up correctly
	total := stats.ShortWords + stats.MediumWords + stats.LongWords + stats.VeryLongWords
	if total != stats.WordCount {
		t.Errorf("Word length categories (%d) don't sum to total words (%d)", 
			total, stats.WordCount)
	}

	// Average word length should be reasonable
	if stats.AverageWordLength < 1 || stats.AverageWordLength > 20 {
		t.Errorf("Average word length %f seems unreasonable", stats.AverageWordLength)
	}
}

func TestSentenceLengthDistribution(t *testing.T) {
	text := "Short. This is a medium sentence. This is a much longer sentence with many words and clauses that extends beyond typical length."
	stats := CalculateTextStatistics(text)

	// Check that sentence length categories sum up correctly
	total := stats.ShortSentences + stats.MediumSentences + stats.LongSentences + stats.VeryLongSentences
	if total != stats.SentenceCount {
		t.Errorf("Sentence length categories (%d) don't sum to total sentences (%d)", 
			total, stats.SentenceCount)
	}

	// Average sentence length should be reasonable
	if stats.AverageSentenceLength < 1 {
		t.Errorf("Average sentence length %f should be at least 1", stats.AverageSentenceLength)
	}
}

func TestBigramAndTrigramAnalysis(t *testing.T) {
	text := "The cat sat on the mat. The dog ran in the park."
	stats := CalculateTextStatistics(text)

	// Should have some bigrams and trigrams
	if len(stats.MostFrequentBigrams) == 0 {
		t.Error("Should have found some bigrams")
	}

	if len(stats.MostFrequentTrigrams) == 0 {
		t.Error("Should have found some trigrams")
	}

	// Bigram frequencies should be positive
	for _, bigram := range stats.MostFrequentBigrams {
		if bigram.Frequency <= 0 {
			t.Errorf("Bigram frequency should be positive, got %f", bigram.Frequency)
		}
	}

	// Trigram frequencies should be positive
	for _, trigram := range stats.MostFrequentTrigrams {
		if trigram.Frequency <= 0 {
			t.Errorf("Trigram frequency should be positive, got %f", trigram.Frequency)
		}
	}
}

func TestExtractWords(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected int // expected word count
	}{
		{
			name:     "Simple sentence",
			input:    "Hello world test",
			expected: 3,
		},
		{
			name:     "With punctuation",
			input:    "Hello, world! How are you?",
			expected: 5,
		},
		{
			name:     "Empty string",
			input:    "",
			expected: 0,
		},
		{
			name:     "Only punctuation",
			input:    "!@#$%^&*()",
			expected: 0,
		},
		{
			name:     "Mixed alphanumeric",
			input:    "Test123 word2 hello",
			expected: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			words := extractWords(tt.input)
			if len(words) != tt.expected {
				t.Errorf("Expected %d words, got %d", tt.expected, len(words))
			}
		})
	}
}

func TestYulesK(t *testing.T) {
	// Create a text with known word frequency distribution
	wordFreq := map[string]int{
		"the":  10,
		"cat":  5,
		"dog":  3,
		"run":  2,
		"fast": 1,
	}
	totalWords := 21

	k := calculateYulesK(wordFreq, totalWords)

	// K should be positive
	if k <= 0 {
		t.Errorf("Yule's K should be positive, got %f", k)
	}

	// Test edge case with empty frequency
	emptyK := calculateYulesK(map[string]int{}, 0)
	if emptyK != 0 {
		t.Errorf("Yule's K for empty frequency should be 0, got %f", emptyK)
	}
}

func TestMTLD(t *testing.T) {
	// Test with a simple word list
	words := []string{"the", "cat", "sat", "on", "the", "mat", "the", "dog", "ran"}
	mtld := calculateMTLD(words)

	// MTLD should be positive
	if mtld <= 0 {
		t.Errorf("MTLD should be positive, got %f", mtld)
	}

	// Test with empty word list
	emptyMTLD := calculateMTLD([]string{})
	if emptyMTLD != 0 {
		t.Errorf("MTLD for empty list should be 0, got %f", emptyMTLD)
	}

	// Test with very short list
	shortMTLD := calculateMTLD([]string{"word"})
	if shortMTLD <= 0 {
		t.Errorf("MTLD for single word should be positive, got %f", shortMTLD)
	}
}