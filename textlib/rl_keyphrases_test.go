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
	"strings"
	"testing"
)

func TestExtractKeyPhrases(t *testing.T) {
	tests := []struct {
		name           string
		text           string
		maxPhrases     int
		expectedAlgo   string
		minPhrases     int
		checkContent   bool
		expectedTerms  []string
	}{
		{
			name:         "Simple text with TF-IDF (max 5)",
			text:         "Machine learning is transforming the technology industry. Machine learning applications are everywhere.",
			maxPhrases:   5,
			expectedAlgo: "tf-idf",
			minPhrases:   2,
			checkContent: true,
			expectedTerms: []string{"machine", "learning"},
		},
		{
			name: "Medium text with statistical (max 20)",
			text: `Natural language processing enables computers to understand human language. 
			       Natural language processing uses various algorithms. These algorithms process text data. 
			       Text data analysis is crucial for natural language processing applications.`,
			maxPhrases:   20,
			expectedAlgo: "statistical",
			minPhrases:   5,
			checkContent: true,
			expectedTerms: []string{"natural language", "processing", "algorithms"},
		},
		{
			name: "Complex text with deep NLP (max 100)",
			text: `Artificial intelligence and machine learning are revolutionizing healthcare. 
			       Deep learning models can diagnose diseases with high accuracy. 
			       Neural networks process medical images to detect anomalies. 
			       Healthcare providers use AI systems for patient care optimization. 
			       Machine learning algorithms analyze patient data for personalized treatment.`,
			maxPhrases:   100,
			expectedAlgo: "deep-nlp",
			minPhrases:   10,
			checkContent: true,
			expectedTerms: []string{"machine learning", "artificial intelligence", "healthcare"},
		},
		{
			name:         "Empty text",
			text:         "",
			maxPhrases:   10,
			expectedAlgo: "",
			minPhrases:   0,
		},
		{
			name:         "Zero max phrases",
			text:         "Some text here",
			maxPhrases:   0,
			expectedAlgo: "",
			minPhrases:   0,
		},
		{
			name:         "Single word repeated",
			text:         "test test test test test",
			maxPhrases:   3,
			expectedAlgo: "tf-idf",
			minPhrases:   1,
			checkContent: true,
			expectedTerms: []string{"test"},
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			phrases := ExtractKeyPhrases(tt.text, tt.maxPhrases)
			
			// Check phrase count
			if len(phrases) < tt.minPhrases {
				t.Errorf("Got %d phrases, expected at least %d", len(phrases), tt.minPhrases)
			}
			
			if len(phrases) > tt.maxPhrases && tt.maxPhrases > 0 {
				t.Errorf("Got %d phrases, expected at most %d", len(phrases), tt.maxPhrases)
			}
			
			// Check content if required
			if tt.checkContent && len(phrases) > 0 {
				// Check if expected terms appear
				foundTerms := make(map[string]bool)
				for _, phrase := range phrases {
					for _, expectedTerm := range tt.expectedTerms {
						if strings.Contains(strings.ToLower(phrase.Text), expectedTerm) {
							foundTerms[expectedTerm] = true
						}
					}
				}
				
				for _, term := range tt.expectedTerms {
					if !foundTerms[term] {
						t.Errorf("Expected term '%s' not found in extracted phrases", term)
					}
				}
				
				// Validate phrase properties
				for i, phrase := range phrases {
					if phrase.Text == "" {
						t.Errorf("Phrase %d has empty text", i)
					}
					
					if phrase.Score < 0 {
						t.Errorf("Phrase %d has negative score: %f", i, phrase.Score)
					}
					
					if phrase.Confidence < 0 || phrase.Confidence > 1 {
						t.Errorf("Phrase %d confidence out of range: %f", i, phrase.Confidence)
					}
					
					if phrase.Category == "" {
						t.Errorf("Phrase %d has empty category", i)
					}
				}
			}
		})
	}
}

func TestKeyPhraseRanking(t *testing.T) {
	// Test that more frequent terms get higher scores
	text := `Data science is important. Data science uses statistics. 
	         Data science requires programming. Statistics is less common. 
	         Programming appears once.`
	
	phrases := ExtractKeyPhrases(text, 5)
	
	// Find "data science" - should be highest ranked
	var dataScienceScore float64
	var statisticsScore float64
	
	for _, phrase := range phrases {
		if strings.Contains(strings.ToLower(phrase.Text), "data") {
			dataScienceScore = phrase.Score
		}
		if strings.Contains(strings.ToLower(phrase.Text), "statistics") {
			statisticsScore = phrase.Score
		}
	}
	
	if dataScienceScore <= statisticsScore {
		t.Error("'data science' should have higher score than 'statistics'")
	}
}

func TestAlgorithmSelection(t *testing.T) {
	text := "Test text for algorithm selection based on maxPhrases parameter"
	
	// Test TF-IDF selection
	phrases1 := ExtractKeyPhrases(text, 5)
	if len(phrases1) == 0 || phrases1[0].Confidence < 0.7 {
		t.Error("TF-IDF algorithm issue")
	}
	
	// Test statistical selection
	phrases2 := ExtractKeyPhrases(text, 30)
	if len(phrases2) == 0 {
		t.Error("Statistical algorithm issue")
	}
	
	// Test deep NLP selection
	phrases3 := ExtractKeyPhrases(text, 60)
	if len(phrases3) == 0 {
		t.Error("Deep NLP algorithm issue")
	}
}

func TestStopWordFiltering(t *testing.T) {
	text := "The machine learning and the artificial intelligence are the future"
	
	phrases := ExtractKeyPhrases(text, 10)
	
	// Check that stop words are not included as standalone phrases
	for _, phrase := range phrases {
		singleWord := !strings.Contains(phrase.Text, " ")
		if singleWord && isStopWord(phrase.Text) {
			t.Errorf("Stop word '%s' included as key phrase", phrase.Text)
		}
	}
}

func TestPhraseCategories(t *testing.T) {
	text := `Machine learning algorithms are complex. Natural language processing is fascinating. 
	         Deep neural network architectures evolve rapidly.`
	
	phrases := ExtractKeyPhrases(text, 20)
	
	categories := make(map[string]int)
	for _, phrase := range phrases {
		categories[phrase.Category]++
	}
	
	// Should have multiple categories
	if len(categories) < 2 {
		t.Errorf("Expected multiple categories, got %d", len(categories))
	}
	
	// Check valid categories
	for cat := range categories {
		validCategories := map[string]bool{
			"term": true, "phrase": true, "multi-phrase": true, "noun-phrase": true,
		}
		if !validCategories[cat] {
			t.Errorf("Invalid category: %s", cat)
		}
	}
}

func TestPhraseDiversity(t *testing.T) {
	// Test that deep NLP produces diverse phrases
	text := strings.Repeat("Advanced machine learning techniques for data analysis. ", 10)
	
	phrases := ExtractKeyPhrases(text, 60) // Use deep NLP
	
	// Count unique phrases
	uniquePhrases := make(map[string]bool)
	for _, phrase := range phrases {
		uniquePhrases[strings.ToLower(phrase.Text)] = true
	}
	
	// Should have reasonable diversity despite repetitive text
	if len(uniquePhrases) < 5 {
		t.Errorf("Expected more phrase diversity, got %d unique phrases", len(uniquePhrases))
	}
}

func TestLongTextPerformance(t *testing.T) {
	// Generate a long document
	sentences := []string{
		"Artificial intelligence transforms business operations.",
		"Machine learning models predict customer behavior.",
		"Deep learning revolutionizes image recognition.",
		"Natural language processing enables chatbots.",
		"Computer vision applications detect objects.",
	}
	
	var longText strings.Builder
	for i := 0; i < 50; i++ {
		longText.WriteString(sentences[i%len(sentences)])
		longText.WriteString(" ")
	}
	
	text := longText.String()
	
	// Test different algorithms
	testCases := []struct {
		maxPhrases int
		name       string
	}{
		{5, "TF-IDF"},
		{25, "Statistical"},
		{75, "Deep NLP"},
	}
	
	for _, tc := range testCases {
		phrases := ExtractKeyPhrases(text, tc.maxPhrases)
		
		if len(phrases) == 0 {
			t.Errorf("%s produced no phrases", tc.name)
		}
		
		// Check phrases are sorted by score
		for i := 1; i < len(phrases); i++ {
			if phrases[i].Score > phrases[i-1].Score {
				t.Errorf("%s: phrases not sorted by score", tc.name)
			}
		}
	}
}

func TestContextExtraction(t *testing.T) {
	text := "The quick brown fox jumps over the lazy dog"
	
	testCases := []struct {
		position    int
		contextSize int
		expected    string
	}{
		{10, 5, "brown"},      // Position 10 is in "brown", context of 5 chars each side
		{0, 10, "The quick"}, // Start position, no leading ellipsis
		{35, 10, "lazy dog"},  // Near end position
	}
	for _, tc := range testCases {
		context := extractContext(text, tc.position, tc.contextSize)
		if !strings.Contains(context, strings.TrimPrefix(strings.TrimSuffix(tc.expected, "..."), "...")) {
			t.Errorf("Context extraction failed: got %q, expected to contain %q", context, tc.expected)
		}
	}
}

func BenchmarkExtractKeyPhrasesTFIDF(b *testing.B) {
	text := "Machine learning algorithms process data efficiently. Data processing requires optimization."
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ExtractKeyPhrases(text, 5)
	}
}

func BenchmarkExtractKeyPhrasesStatistical(b *testing.B) {
	text := strings.Repeat("Natural language processing enables advanced text analysis. ", 5)
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ExtractKeyPhrases(text, 25)
	}
}

func BenchmarkExtractKeyPhrasesDeep(b *testing.B) {
	text := strings.Repeat("Deep learning neural networks revolutionize artificial intelligence. ", 10)
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ExtractKeyPhrases(text, 60)
	}
}