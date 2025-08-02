package textlib

import (
	"strings"
	"testing"
)

// Comprehensive tests for core textlib functions

func TestSplitIntoSentencesComprehensive(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:     "Empty string",
			input:    "",
			expected: []string{},
		},
		{
			name:     "Single sentence",
			input:    "Hello world.",
			expected: []string{"Hello world."},
		},
		{
			name:     "Multiple sentences",
			input:    "First sentence. Second sentence! Third question?",
			expected: []string{"First sentence.", "Second sentence!", "Third question?"},
		},
		{
			name:     "With abbreviations",
			input:    "Dr. Smith visited the U.S. yesterday. He met Mr. Johnson.",
			expected: []string{"Dr. Smith visited the U.S. yesterday.", "He met Mr. Johnson."},
		},
		{
			name:     "With Inc. and Corp.",
			input:    "Apple Inc. is a company. Microsoft Corp. is another company.",
			expected: []string{"Apple Inc. is a company.", "Microsoft Corp. is another company."},
		},
		{
			name:     "No punctuation",
			input:    "This has no punctuation",
			expected: []string{"This has no punctuation"},
		},
		{
			name:     "Only punctuation",
			input:    "...",
			expected: []string{"..."},
		},
		{
			name:     "Mixed punctuation",
			input:    "Question? Answer! Statement.",
			expected: []string{"Question?", "Answer!", "Statement."},
		},
		{
			name:     "Whitespace handling",
			input:    "First.   Second.    Third.",
			expected: []string{"First.", "Second.", "Third."},
		},
		{
			name:     "Newlines in text",
			input:    "First sentence.\nSecond sentence.",
			expected: []string{"First sentence.", "Second sentence."},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := SplitIntoSentences(test.input)
			
			if len(result) != len(test.expected) {
				t.Errorf("Expected %d sentences, got %d", len(test.expected), len(result))
				t.Errorf("Expected: %v", test.expected)
				t.Errorf("Got: %v", result)
				return
			}
			
			for i, expected := range test.expected {
				if result[i] != expected {
					t.Errorf("Sentence %d: expected %q, got %q", i, expected, result[i])
				}
			}
		})
	}
}

func TestSplitIntoParagraphsComprehensive(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:     "Empty string",
			input:    "",
			expected: []string{},
		},
		{
			name:     "Single paragraph",
			input:    "This is a single paragraph with multiple sentences. It continues here.",
			expected: []string{"This is a single paragraph with multiple sentences. It continues here."},
		},
		{
			name:     "Multiple paragraphs",
			input:    "First paragraph.\n\nSecond paragraph.\n\nThird paragraph.",
			expected: []string{"First paragraph.", "Second paragraph.", "Third paragraph."},
		},
		{
			name:     "Extra whitespace",
			input:    "First paragraph.\n\n\n  \nSecond paragraph.",
			expected: []string{"First paragraph.", "Second paragraph."},
		},
		{
			name:     "Tabs and spaces",
			input:    "First paragraph.\n\t\n  \n\nSecond paragraph.",
			expected: []string{"First paragraph.", "Second paragraph."},
		},
		{
			name:     "Single newlines should not split",
			input:    "First line\nSecond line\n\nNew paragraph",
			expected: []string{"First line\nSecond line", "New paragraph"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := SplitIntoParagraphs(test.input)
			
			if len(result) != len(test.expected) {
				t.Errorf("Expected %d paragraphs, got %d", len(test.expected), len(result))
				t.Errorf("Expected: %v", test.expected)
				t.Errorf("Got: %v", result)
				return
			}
			
			for i, expected := range test.expected {
				if result[i] != expected {
					t.Errorf("Paragraph %d: expected %q, got %q", i, expected, result[i])
				}
			}
		})
	}
}

func TestExtractNamedEntitiesComprehensive(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected map[string]int // entity type -> count
	}{
		{
			name:     "Empty string",
			input:    "",
			expected: map[string]int{},
		},
		{
			name:     "No entities",
			input:    "This text has no named entities.",
			expected: map[string]int{},
		},
		{
			name:     "Person names",
			input:    "John Smith and Mary Johnson met yesterday.",
			expected: map[string]int{"PERSON": 2},
		},
		{
			name:     "Organizations",
			input:    "Apple Inc and Microsoft Corporation are tech companies.",
			expected: map[string]int{"ORGANIZATION": 2},
		},
		{
			name:     "Locations",
			input:    "I visited New York, NY and Los Angeles, CA last week.",
			expected: map[string]int{"LOCATION": 2},
		},
		{
			name:     "Dates",
			input:    "The meeting is on January 15, 2024 and February 20, 2024.",
			expected: map[string]int{"DATE": 2},
		},
		{
			name:     "Money amounts",
			input:    "The price is $100.50 and the tax is $10.25.",
			expected: map[string]int{"MONEY": 2},
		},
		{
			name:     "Percentages",
			input:    "The growth was 15% last year and 20% this year.",
			expected: map[string]int{"PERCENT": 2},
		},
		{
			name:     "Time expressions",
			input:    "The meeting is at 10:30 AM and ends at 2:45 PM.",
			expected: map[string]int{"TIME": 2},
		},
		{
			name:     "Mixed entities",
			input:    "John Smith from Apple Inc visited New York, NY on January 15, 2024 at 10:30 AM for a $1,000,000 deal with 15% commission.",
			expected: map[string]int{
				"PERSON":       1,
				"ORGANIZATION": 1,
				"LOCATION":     1,
				"DATE":         1,
				"TIME":         1,
				"MONEY":        1,
				"PERCENT":      1,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			entities := ExtractNamedEntities(test.input)
			
			// Count entities by type
			counts := make(map[string]int)
			for _, entity := range entities {
				counts[entity.Type]++
			}
			
			// Check counts match expected
			for expectedType, expectedCount := range test.expected {
				if counts[expectedType] != expectedCount {
					t.Errorf("Entity type %s: expected %d, got %d", 
						expectedType, expectedCount, counts[expectedType])
				}
			}
			
			// Check no unexpected entities
			for actualType, actualCount := range counts {
				if test.expected[actualType] != actualCount {
					t.Errorf("Unexpected entity type %s with count %d", actualType, actualCount)
				}
			}
		})
	}
}

func TestCalculateSyllableCountComprehensive(t *testing.T) {
	tests := []struct {
		word     string
		expected int
	}{
		// Basic cases
		{"", 0},
		{"a", 1},
		{"I", 1},
		{"the", 1},
		{"cat", 1},
		{"dog", 1},
		
		// Two syllables
		{"hello", 2},
		{"water", 2},
		{"table", 2},
		{"simple", 2},
		{"apple", 2},
		
		// Three syllables
		{"beautiful", 3},
		{"computer", 3},
		{"telephone", 3},
		{"elephant", 3},
		
		// Four syllables
		{"education", 4},
		{"information", 4},
		{"television", 4},
		
		// Silent e cases
		{"make", 1},
		{"take", 1},
		{"like", 1},
		{"bike", 1},
		{"home", 1},
		{"time", 1},
		
		// -le endings
		{"castle", 2},
		{"little", 2},
		{"middle", 2},
		{"circle", 2},
		{"people", 2},
		
		// Edge cases
		{"rhythm", 1}, // No vowels except y
		{"strengths", 1}, // Consonant cluster
		{"eye", 1}, // Diphthong
		
		// Numbers and punctuation should return 0
		{"123", 0},
		{"!!!", 0},
		{"@#$", 0},
		
		// Mixed alphanumeric
		{"test123", 1}, // Should extract letters only
	}

	for _, test := range tests {
		t.Run(test.word, func(t *testing.T) {
			result := CalculateSyllableCount(test.word)
			if result != test.expected {
				t.Errorf("CalculateSyllableCount(%q): expected %d, got %d", 
					test.word, test.expected, result)
			}
		})
	}
}

func TestIsCompleteSentenceComprehensive(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		// Valid complete sentences
		{"Simple declarative", "The cat sits.", true},
		{"With exclamation", "What a day!", true},
		{"Question", "How are you?", true},
		{"Longer sentence", "The quick brown fox jumps over the lazy dog.", true},
		{"With common verb", "She is happy.", true},
		{"Past tense verb", "He walked home.", true},
		{"Multiple words with verb", "They are running fast.", true},
		
		// Invalid/incomplete sentences
		{"Empty string", "", false},
		{"Too short", "Hi", false},
		{"No punctuation", "This has no end punctuation", false},
		{"No verb", "The big red car.", false},
		{"Fragment", "Running fast.", false},
		{"Just punctuation", "...", false},
		{"Single word", "Hello.", false},
		
		// Edge cases
		{"Imperative", "Go!", true}, // "Go" can be considered a verb
		{"With auxiliary", "She has been working.", true},
		{"Contraction", "It's raining.", true},
		{"Question word", "What happened?", true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := IsCompleteSentence(test.input)
			if result != test.expected {
				t.Errorf("IsCompleteSentence(%q): expected %v, got %v", 
					test.input, test.expected, result)
			}
		})
	}
}

func TestCountWordsComprehensive(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected int
	}{
		{"Empty string", "", 0},
		{"Single word", "hello", 1},
		{"Two words", "hello world", 2},
		{"Multiple spaces", "hello    world", 2},
		{"Leading/trailing spaces", "  hello world  ", 2},
		{"Tabs and newlines", "hello\tworld\ntest", 3},
		{"Punctuation", "hello, world!", 2},
		{"Numbers", "I have 5 cats", 3},
		{"Mixed content", "The 2023 year was great!", 5},
		{"Only spaces", "   ", 0},
		{"Only punctuation", "...", 1}, // Fields counts punctuation as words
		{"Contractions", "don't can't won't", 3},
		{"Hyphenated", "twenty-five years old", 3},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := CountWords(test.input)
			if result != test.expected {
				t.Errorf("CountWords(%q): expected %d, got %d", 
					test.input, test.expected, result)
			}
		})
	}
}

func TestCountSentencesComprehensive(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected int
	}{
		{"Empty string", "", 0},
		{"Single sentence", "Hello world.", 1},
		{"Multiple sentences", "First. Second! Third?", 3},
		{"With abbreviations", "Dr. Smith visited the U.S. yesterday.", 1},
		{"No punctuation", "This has no punctuation", 1},
		{"Only punctuation", "...", 1},
		{"Mixed punctuation", "Really? Yes! Absolutely.", 3},
		{"Newlines", "First sentence.\nSecond sentence.", 2},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := CountSentences(test.input)
			if result != test.expected {
				t.Errorf("CountSentences(%q): expected %d, got %d", 
					test.input, test.expected, result)
			}
		})
	}
}

func TestCountSyllablesComprehensive(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected int
	}{
		{"Empty string", "", 0},
		{"Single syllable word", "cat", 1},
		{"Multiple single syllable", "the cat sat", 3},
		{"Mixed syllables", "beautiful computer", 6}, // 3 + 3
		{"With punctuation", "Hello, world!", 3}, // 2 + 1
		{"Numbers ignored", "I have 5 cats", 3}, // "I"(1) + "have"(1) + "cats"(1)
		{"Complex sentence", "The beautiful butterfly flew gracefully.", 11},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := CountSyllables(test.input)
			if result != test.expected {
				t.Errorf("CountSyllables(%q): expected %d, got %d", 
					test.input, test.expected, result)
			}
		})
	}
}

func TestCalculateFleschReadingEaseComprehensive(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		minScore float64
		maxScore float64
	}{
		{
			name:     "Empty string",
			input:    "",
			minScore: 0,
			maxScore: 0,
		},
		{
			name:     "Simple text",
			input:    "The cat sat on the mat.",
			minScore: 80, // Should be fairly easy to read
			maxScore: 120,
		},
		{
			name:     "Complex text",
			input:    "The implementation of sophisticated algorithms requires comprehensive understanding of computational complexity theory.",
			minScore: 0,  // Should be harder to read
			maxScore: 50,
		},
		{
			name:     "Very simple",
			input:    "I am. You are. We go.",
			minScore: 90,
			maxScore: 150,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := CalculateFleschReadingEase(test.input)
			
			if result < test.minScore || result > test.maxScore {
				t.Errorf("CalculateFleschReadingEase(%q): expected score between %.1f and %.1f, got %.1f", 
					test.input, test.minScore, test.maxScore, result)
			}
		})
	}
}

func TestIsVowelHelper(t *testing.T) {
	vowels := []rune{'a', 'e', 'i', 'o', 'u', 'y', 'A', 'E', 'I', 'O', 'U', 'Y'}
	consonants := []rune{'b', 'c', 'd', 'f', 'g', 'h', 'j', 'k', 'l', 'm', 'n', 'p', 'q', 'r', 's', 't', 'v', 'w', 'x', 'z'}

	for _, vowel := range vowels {
		if !isVowel(vowel) {
			t.Errorf("Expected %c to be a vowel", vowel)
		}
	}

	for _, consonant := range consonants {
		if isVowel(consonant) {
			t.Errorf("Expected %c to be a consonant", consonant)
		}
	}
}

func TestIsVerbSimpleHelper(t *testing.T) {
	tests := []struct {
		word     string
		expected bool
	}{
		// Common verbs
		{"is", true},
		{"are", true},
		{"was", true},
		{"were", true},
		{"have", true},
		{"do", true},
		{"will", true},
		{"can", true},
		
		// Verb patterns
		{"walked", true},  // -ed ending
		{"running", true}, // -ing ending
		{"goes", true},    // -s ending
		{"talks", true},   // -s ending
		
		// Non-verbs
		{"cat", false},
		{"book", false},
		{"red", false},
		{"quickly", false},
		
		// Edge cases
		{"", false},
		{"a", false},
		{"the", false},
		
		// With punctuation
		{"walked.", true},
		{"running!", true},
		{"is?", true},
	}

	for _, test := range tests {
		t.Run(test.word, func(t *testing.T) {
			result := isVerbSimple(test.word)
			if result != test.expected {
				t.Errorf("isVerbSimple(%q): expected %v, got %v", 
					test.word, test.expected, result)
			}
		})
	}
}

// Test error conditions and edge cases
func TestErrorConditions(t *testing.T) {
	t.Run("Nil input handling", func(t *testing.T) {
		// Test functions with empty/nil inputs
		result := SplitIntoSentences("")
		if len(result) != 0 {
			t.Errorf("Expected empty slice for empty input")
		}
		
		entities := ExtractNamedEntities("")
		if len(entities) != 0 {
			t.Errorf("Expected no entities for empty input")
		}
		
		if CountWords("") != 0 {
			t.Errorf("Expected 0 words for empty input")
		}
		
		if CountSentences("") != 0 {
			t.Errorf("Expected 0 sentences for empty input")
		}
	})

	t.Run("Very long input", func(t *testing.T) {
		// Test with very long input
		longText := strings.Repeat("This is a test sentence. ", 1000)
		
		sentences := SplitIntoSentences(longText)
		if len(sentences) != 1000 {
			t.Errorf("Expected 1000 sentences, got %d", len(sentences))
		}
		
		wordCount := CountWords(longText)
		expectedWords := 6 * 1000 // 6 words per sentence * 1000 sentences
		if wordCount != expectedWords {
			t.Errorf("Expected %d words, got %d", expectedWords, wordCount)
		}
	})

	t.Run("Unicode and special characters", func(t *testing.T) {
		unicodeText := "Héllo wörld! 你好世界。Привет мир!"
		
		// Should handle unicode without crashing
		sentences := SplitIntoSentences(unicodeText)
		if len(sentences) == 0 {
			t.Errorf("Expected to handle unicode text")
		}
		
		words := CountWords(unicodeText)
		if words == 0 {
			t.Errorf("Expected to count unicode words")
		}
	})
}

// Benchmark tests
func BenchmarkSplitIntoSentences(b *testing.B) {
	text := "This is a test sentence. This is another test sentence! Is this a question? Yes, it is."
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		SplitIntoSentences(text)
	}
}

func BenchmarkExtractNamedEntities(b *testing.B) {
	text := "John Smith from Apple Inc visited New York, NY on January 15, 2024 at 10:30 AM for a $1,000,000 deal."
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ExtractNamedEntities(text)
	}
}

func BenchmarkCalculateFleschReadingEase(b *testing.B) {
	text := "The quick brown fox jumps over the lazy dog. This is a test sentence for benchmarking purposes. It contains multiple sentences with varying complexity levels."
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		CalculateFleschReadingEase(text)
	}
}