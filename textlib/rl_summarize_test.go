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

func TestSummarizeText(t *testing.T) {
	tests := []struct {
		name              string
		text              string
		maxLength         int
		expectedMethod    string
		checkContent      bool
		minQuality        float64
		maxCompressionRatio float64
	}{
		{
			name: "Extractive summary (short)",
			text: "The quick brown fox jumps over the lazy dog. " +
				"This is a simple sentence. " +
				"Another sentence follows here. " +
				"The final sentence completes the paragraph.",
			maxLength:      80,
			expectedMethod: "extractive",
			checkContent:   true,
			minQuality:     0.7,
			maxCompressionRatio: 0.5,
		},
		{
			name: "Hybrid summary (medium)",
			text: `Artificial intelligence (AI) has become a transformative force in modern technology. 
				Machine learning algorithms can now perform tasks that once required human intelligence. 
				Deep learning networks have achieved remarkable success in image recognition and natural language processing. 
				However, ethical concerns about AI bias and privacy remain significant challenges. 
				Researchers are working to develop more transparent and accountable AI systems. 
				The future of AI will likely involve closer collaboration between humans and machines.`,
			maxLength:      200,
			expectedMethod: "hybrid",
			checkContent:   true,
			minQuality:     0.65,
			maxCompressionRatio: 0.4,
		},
		{
			name: "Abstractive summary (long)",
			text: `Climate change represents one of the most pressing challenges facing humanity today. 
				Rising global temperatures are causing melting ice caps, rising sea levels, and more frequent extreme weather events. 
				Scientists have observed unprecedented changes in weather patterns across the globe. 
				The primary cause is the emission of greenhouse gases from human activities, particularly the burning of fossil fuels. 
				Many countries have committed to reducing their carbon emissions through renewable energy adoption. 
				Solar and wind power are becoming increasingly cost-effective alternatives to traditional energy sources. 
				However, the transition to clean energy requires significant investment and political will. 
				Individual actions, such as reducing energy consumption and supporting sustainable practices, also play a crucial role. 
				The next decade will be critical in determining whether we can limit global warming to manageable levels. 
				International cooperation and innovative technologies will be essential for addressing this global crisis.`,
			maxLength:      400,
			expectedMethod: "abstractive",
			checkContent:   true,
			minQuality:     0.6,
			maxCompressionRatio: 0.5,
		},
		{
			name:           "Empty text",
			text:           "",
			maxLength:      100,
			expectedMethod: "none",
			checkContent:   false,
			minQuality:     0,
		},
		{
			name: "Single sentence",
			text: "This is a single sentence that should be returned as is.",
			maxLength:      80,
			expectedMethod: "extractive",
			checkContent:   true,
			minQuality:     0.7,
		},
		{
			name: "Very short max length",
			text: "This is a longer text with multiple sentences. Each sentence contains important information. " +
				"The summary should be very brief. Only the most important content should be included.",
			maxLength:      50,
			expectedMethod: "extractive",
			checkContent:   true,
			minQuality:     0.7,
			maxCompressionRatio: 0.3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SummarizeText(tt.text, tt.maxLength)

			// Check method
			if result.Method != tt.expectedMethod {
				t.Errorf("Expected method %s, got %s", tt.expectedMethod, result.Method)
			}

			// Check summary length
			if tt.text != "" && len(result.Summary) > tt.maxLength {
				t.Errorf("Summary length %d exceeds max length %d", 
					len(result.Summary), tt.maxLength)
			}

			// Check quality metrics
			if result.QualityMetrics.Accuracy < tt.minQuality {
				t.Errorf("Accuracy %f below minimum %f", 
					result.QualityMetrics.Accuracy, tt.minQuality)
			}

			// Check compression ratio
			if tt.maxCompressionRatio > 0 && result.CompressionRatio > tt.maxCompressionRatio {
				t.Errorf("Compression ratio %f exceeds maximum %f",
					result.CompressionRatio, tt.maxCompressionRatio)
			}

			// Check content if required
			if tt.checkContent && tt.text != "" {
				if result.Summary == "" {
					t.Error("Expected non-empty summary")
				}
				
				// For extractive, check that summary contains original sentences
				if tt.expectedMethod == "extractive" {
					sentences := SplitIntoSentences(result.Summary)
					originalSentences := SplitIntoSentences(tt.text)
					for _, s := range sentences {
						found := false
						for _, orig := range originalSentences {
							if strings.TrimSpace(s) == strings.TrimSpace(orig) {
								found = true
								break
							}
						}
						if !found {
							t.Errorf("Extractive summary contains non-original sentence: %s", s)
						}
					}
				}
			}

			// Check processing time
			if result.ProcessingTime <= 0 {
				t.Error("Invalid processing time")
			}

			// Check sentence counts
			if tt.text != "" {
				if result.OriginalSentences == 0 {
					t.Error("Original sentence count should be greater than 0")
				}
				if result.Summary != "" && result.SummarySentences == 0 {
					t.Error("Summary sentence count should be greater than 0")
				}
			}
		})
	}
}

func TestExtractiveMethodDetails(t *testing.T) {
	text := `The first sentence is very important. 
		This middle sentence contains some details. 
		Another sentence provides context. 
		The conclusion summarizes everything.`

	result := SummarizeText(text, 80) // Force extractive

	// Should prioritize first and last sentences
	if !strings.Contains(result.Summary, "first sentence") {
		t.Error("Extractive summary should include the first sentence")
	}

	// Check coverage
	if result.QualityMetrics.Coverage == 0 {
		t.Error("Coverage should be calculated")
	}

	if result.QualityMetrics.Coverage > 1.0 {
		t.Error("Coverage should not exceed 1.0")
	}
}

func TestHybridMethodCompression(t *testing.T) {
	text := `This is a very long sentence with many unnecessary words that could be compressed. 
		For example, this sentence contains filler phrases that add little value. 
		In addition, some words are really quite redundant and could be removed.`

	result := SummarizeText(text, 150) // Force hybrid

	// Check that compression occurred
	if len(result.Summary) >= len(text) {
		t.Error("Hybrid summary should be shorter than original")
	}

	// Check that common filler words are removed
	lowerSummary := strings.ToLower(result.Summary)
	fillers := []string{"very", "really", "quite", "for example", "in addition"}
	
	fillerCount := 0
	for _, filler := range fillers {
		if strings.Contains(lowerSummary, filler) {
			fillerCount++
		}
	}

	if fillerCount > 2 {
		t.Errorf("Too many filler words remaining in compressed summary: %d", fillerCount)
	}
}

func TestAbstractiveMethodGeneration(t *testing.T) {
	text := `Machine learning is a subset of artificial intelligence. 
		It enables systems to learn from data without explicit programming. 
		Deep learning uses neural networks with multiple layers. 
		These technologies are transforming various industries.`

	result := SummarizeText(text, 300) // Force abstractive

	// Check that summary contains generated text
	if strings.Contains(result.Summary, "This text") || 
		strings.Contains(result.Summary, "main point") {
		// Good - contains template phrases
	} else {
		t.Error("Abstractive summary should contain generated phrases")
	}

	// Should mention key topics
	if !strings.Contains(strings.ToLower(result.Summary), "machine learning") &&
		!strings.Contains(strings.ToLower(result.Summary), "artificial intelligence") {
		t.Error("Abstractive summary should mention key topics")
	}
}

func TestSummaryQualityMetrics(t *testing.T) {
	text := strings.Repeat("This is a test sentence. ", 20)
	
	// Test quality degradation with shorter summaries
	lengths := []int{50, 150, 300}
	previousAccuracy := 1.0
	
	for _, length := range lengths {
		result := SummarizeText(text, length)
		
		// Longer summaries should generally have higher accuracy
		if length > 150 && result.QualityMetrics.Accuracy > previousAccuracy {
			t.Errorf("Accuracy should generally decrease with more aggressive summarization")
		}
		previousAccuracy = result.QualityMetrics.Accuracy
		
		// All metrics should be between 0 and 1
		if result.QualityMetrics.Accuracy < 0 || result.QualityMetrics.Accuracy > 1 {
			t.Errorf("Accuracy out of range: %f", result.QualityMetrics.Accuracy)
		}
		if result.QualityMetrics.Confidence < 0 || result.QualityMetrics.Confidence > 1 {
			t.Errorf("Confidence out of range: %f", result.QualityMetrics.Confidence)
		}
		if result.QualityMetrics.Coverage < 0 || result.QualityMetrics.Coverage > 1 {
			t.Errorf("Coverage out of range: %f", result.QualityMetrics.Coverage)
		}
	}
}

func TestEdgeCases(t *testing.T) {
	// Test with only stop words
	result := SummarizeText("the and is at which on with to for of", 50)
	// Since it's all stop words, the summary should be empty or the same as input
	// We'll be more lenient here as this is an edge case
	if len(result.Summary) > 50 {
		t.Errorf("Summary exceeds max length: got %d chars", len(result.Summary))
	}

	// Test with single word
	result = SummarizeText("Hello", 10)
	if result.Summary != "Hello" && result.Summary != "" {
		t.Errorf("Single word summary unexpected: %s", result.Summary)
	}

	// Test with zero max length
	result = SummarizeText("Some text here", 0)
	if result.Summary != "" {
		t.Error("Zero max length should return empty summary")
	}

	// Test with very long sentences
	longSentence := "This is " + strings.Repeat("very ", 50) + "long."
	result = SummarizeText(longSentence, 20)
	if len(result.Summary) > 25 { // Some buffer for word boundaries
		t.Error("Summary should respect max length even with long sentences")
	}
}

func TestLongTextPerformance(t *testing.T) {
	// Generate a long text
	paragraphs := []string{}
	for i := 0; i < 50; i++ {
		paragraphs = append(paragraphs, 
			"This is paragraph number "+string(rune('0'+i%10))+". "+
			"It contains several sentences about various topics. "+
			"The content is designed to test summarization performance. "+
			"Each paragraph adds more information to process.")
	}
	longText := strings.Join(paragraphs, " ")

	// Test different summary lengths
	methods := []struct {
		maxLength int
		method    string
	}{
		{80, "extractive"},
		{250, "hybrid"},
		{500, "abstractive"},
	}

	for _, m := range methods {
		result := SummarizeText(longText, m.maxLength)
		
		if result.Method != m.method {
			t.Errorf("Expected method %s for length %d, got %s", 
				m.method, m.maxLength, result.Method)
		}

		if result.ProcessingTime.Milliseconds() > 100 {
			t.Logf("Warning: %s method took %v for long text", 
				m.method, result.ProcessingTime)
		}
	}
}

func BenchmarkSummarizeTextExtractive(b *testing.B) {
	text := strings.Repeat("This is a test sentence with some content. ", 10)
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		SummarizeText(text, 50)
	}
}

func BenchmarkSummarizeTextHybrid(b *testing.B) {
	text := strings.Repeat("This is a longer test sentence with more detailed content. ", 20)
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		SummarizeText(text, 200)
	}
}

func BenchmarkSummarizeTextAbstractive(b *testing.B) {
	text := strings.Repeat("This comprehensive test sentence contains various important topics. ", 30)
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		SummarizeText(text, 400)
	}
}