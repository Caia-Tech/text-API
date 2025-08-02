package textlib

import (
	"testing"
)

func TestSplitIntoSentences(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:     "Simple sentences",
			input:    "Hello world. How are you? I'm fine!",
			expected: []string{"Hello world.", "How are you?", "I'm fine!"},
		},
		{
			name:     "With abbreviations",
			input:    "Dr. Smith went to the U.S. yesterday. He met Mr. Johnson.",
			expected: []string{"Dr. Smith went to the U.S. yesterday.", "He met Mr. Johnson."},
		},
		{
			name:     "Empty string",
			input:    "",
			expected: []string{},
		},
		{
			name:     "Single sentence",
			input:    "This is one sentence.",
			expected: []string{"This is one sentence."},
		},
		{
			name:     "No punctuation",
			input:    "This has no ending punctuation",
			expected: []string{"This has no ending punctuation"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SplitIntoSentences(tt.input)
			if len(result) != len(tt.expected) {
				t.Errorf("Expected %d sentences, got %d", len(tt.expected), len(result))
				return
			}
			for i, expected := range tt.expected {
				if result[i] != expected {
					t.Errorf("Sentence %d: expected %q, got %q", i, expected, result[i])
				}
			}
		})
	}
}

func TestSplitIntoParagraphs(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:     "Multiple paragraphs",
			input:    "First paragraph.\n\nSecond paragraph.\n\nThird paragraph.",
			expected: []string{"First paragraph.", "Second paragraph.", "Third paragraph."},
		},
		{
			name:     "Single paragraph",
			input:    "Just one paragraph with multiple sentences. Here's another sentence.",
			expected: []string{"Just one paragraph with multiple sentences. Here's another sentence."},
		},
		{
			name:     "Empty string",
			input:    "",
			expected: []string{},
		},
		{
			name:     "Extra whitespace",
			input:    "Para one.\n\n\n\nPara two.\n\n   \n\nPara three.",
			expected: []string{"Para one.", "Para two.", "Para three."},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SplitIntoParagraphs(tt.input)
			if len(result) != len(tt.expected) {
				t.Errorf("Expected %d paragraphs, got %d", len(tt.expected), len(result))
				return
			}
			for i, expected := range tt.expected {
				if result[i] != expected {
					t.Errorf("Paragraph %d: expected %q, got %q", i, expected, result[i])
				}
			}
		})
	}
}

func TestExtractNamedEntities(t *testing.T) {
	text := "John Smith works at Apple Inc. He lives in New York, NY and was born on January 15, 1990. He earns $75,000 per year."
	entities := ExtractNamedEntities(text)

	// Should find at least some entities
	if len(entities) == 0 {
		t.Error("Expected to find named entities but found none")
	}

	// Check for specific entity types
	foundTypes := make(map[string]bool)
	for _, entity := range entities {
		foundTypes[entity.Type] = true
	}

	expectedTypes := []string{"PERSON", "ORGANIZATION", "LOCATION", "DATE", "MONEY"}
	for _, expectedType := range expectedTypes {
		if !foundTypes[expectedType] {
			t.Errorf("Expected to find entity type %s", expectedType)
		}
	}
}

func TestCalculateSyllableCount(t *testing.T) {
	tests := []struct {
		word     string
		expected int
	}{
		{"cat", 1},
		{"hello", 2},
		{"beautiful", 3},
		{"education", 4},
		{"a", 1},
		{"the", 1},
		{"", 0},
		{"simple", 2},
		{"castle", 2}, // silent e
	}

	for _, tt := range tests {
		t.Run(tt.word, func(t *testing.T) {
			result := CalculateSyllableCount(tt.word)
			if result != tt.expected {
				t.Errorf("Word %q: expected %d syllables, got %d", tt.word, tt.expected, result)
			}
		})
	}
}

func TestIsCompleteSentence(t *testing.T) {
	tests := []struct {
		text     string
		expected bool
	}{
		{"The cat sat on the mat.", true},
		{"Hello world!", true},
		{"What are you doing?", true},
		{"Running fast", false},
		{"Because it was raining", false},
		{"", false},
		{"Yes.", false}, // Too short, no clear subject-predicate
		{"The dog barks loudly.", true},
	}

	for _, tt := range tests {
		t.Run(tt.text, func(t *testing.T) {
			result := IsCompleteSentence(tt.text)
			if result != tt.expected {
				t.Errorf("Text %q: expected %v, got %v", tt.text, tt.expected, result)
			}
		})
	}
}

func TestCalculateFleschReadingEase(t *testing.T) {
	// Test with simple text
	simpleText := "The cat sat. The dog ran."
	complexText := "Notwithstanding the aforementioned considerations, the implementation of sophisticated algorithms necessitates comprehensive evaluation."

	simpleScore := CalculateFleschReadingEase(simpleText)
	complexScore := CalculateFleschReadingEase(complexText)

	// Simple text should have higher readability score
	if simpleScore <= complexScore {
		t.Errorf("Simple text should have higher Flesch score than complex text. Simple: %f, Complex: %f", 
			simpleScore, complexScore)
	}

	// Scores should be reasonable (though can go negative for very complex text)
	if simpleScore < 0 || simpleScore > 120 {
		t.Errorf("Simple text Flesch score %f is outside reasonable range", simpleScore)
	}
}

func TestCountFunctions(t *testing.T) {
	text := "Hello world. How are you today? I am fine, thank you!"

	wordCount := CountWords(text)
	sentenceCount := CountSentences(text)
	syllableCount := CountSyllables(text)

	if wordCount == 0 {
		t.Error("Word count should be greater than 0")
	}

	if sentenceCount != 3 {
		t.Errorf("Expected 3 sentences, got %d", sentenceCount)
	}

	if syllableCount == 0 {
		t.Error("Syllable count should be greater than 0")
	}

	// Test edge cases
	emptyWordCount := CountWords("")
	if emptyWordCount != 0 {
		t.Errorf("Empty string should have 0 words, got %d", emptyWordCount)
	}

	emptySentenceCount := CountSentences("")
	if emptySentenceCount != 0 {
		t.Errorf("Empty string should have 0 sentences, got %d", emptySentenceCount)
	}
}