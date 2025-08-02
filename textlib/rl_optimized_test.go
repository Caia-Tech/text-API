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
	"time"
)

func TestSmartAnalyze(t *testing.T) {
	text := "Dr. Smith from OpenAI announced a $10M funding round. The new AI model shows 95% accuracy."
	
	result := SmartAnalyze(text)
	
	// Check entities were extracted
	if len(result.Entities) == 0 {
		t.Error("Expected entities to be extracted")
	}
	
	// Check readability was calculated
	if result.Readability.FleschScore == 0 {
		t.Error("Expected readability score to be calculated")
	}
	
	// Check keywords were extracted
	if len(result.Keywords) == 0 {
		t.Error("Expected keywords to be extracted")
	}
	
	// Check processing info
	if len(result.ProcessingInfo.FunctionsUsed) == 0 {
		t.Error("Expected functions used to be tracked")
	}
	
	// Verify optimization flag
	if result.ProcessingInfo.OptimizationUsed != "rl-optimized-general" {
		t.Errorf("Expected optimization 'rl-optimized-general', got '%s'", 
			result.ProcessingInfo.OptimizationUsed)
	}
}

func TestValidatedExtraction(t *testing.T) {
	// Text with potential entity overlap
	text := "OpenAI and OpenAI GPT-4 are different things. Dr. John Smith and Smith work here."
	
	entities := ValidatedExtraction(text)
	
	// Check that entities were extracted
	if len(entities) == 0 {
		t.Error("Expected entities to be extracted")
	}
	
	// Check for duplicates
	seen := make(map[string]bool)
	for _, entity := range entities {
		if seen[entity.Text] {
			t.Errorf("Duplicate entity found: %s", entity.Text)
		}
		seen[entity.Text] = true
	}
	
	// Should merge overlapping entities
	hasOpenAI := false
	for _, entity := range entities {
		if entity.Text == "OpenAI" || entity.Text == "OpenAI GPT-4" {
			hasOpenAI = true
		}
	}
	
	if !hasOpenAI {
		t.Error("Expected to find OpenAI entities")
	}
}

func TestDomainOptimizedAnalyze(t *testing.T) {
	tests := []struct {
		name     string
		text     string
		domain   string
		expected string
	}{
		{
			name:     "Technical domain",
			text:     "The function calculateSum() has O(n) complexity.",
			domain:   "technical",
			expected: "rl-optimized-technical",
		},
		{
			name:     "Medical domain",
			text:     "Patient shows symptoms of hypertension with BP 140/90.",
			domain:   "medical",
			expected: "rl-optimized-medical",
		},
		{
			name:     "Legal domain",
			text:     "Pursuant to Section 5.2, the party shall indemnify.",
			domain:   "legal",
			expected: "rl-optimized-legal",
		},
		{
			name:     "Social domain",
			text:     "Just launched our app! #startup @techcrunch",
			domain:   "social",
			expected: "rl-optimized-social",
		},
		{
			name:     "Business domain",
			text:     "Q3 revenue increased 23% YoY. Action: expand to EU.",
			domain:   "business",
			expected: "rl-optimized-business",
		},
		{
			name:     "Unknown domain",
			text:     "This is some text.",
			domain:   "unknown",
			expected: "rl-optimized-general",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := DomainOptimizedAnalyze(tt.text, tt.domain)
			
			if result.ProcessingInfo.OptimizationUsed != tt.expected {
				t.Errorf("Expected optimization '%s', got '%s'",
					tt.expected, result.ProcessingInfo.OptimizationUsed)
			}
			
			// Domain-specific checks
			switch tt.domain {
			case "social":
				// Should have quick processing
				if len(result.ProcessingInfo.FunctionsUsed) > 5 {
					t.Error("Social domain should use minimal functions for speed")
				}
			case "technical":
				// Should detect code-related entities
				hasCodeEntity := false
				for _, entity := range result.Entities {
					if entity.Type == "FUNCTION" {
						hasCodeEntity = true
						break
					}
				}
				if strings.Contains(tt.text, "function") && !hasCodeEntity {
					t.Error("Technical domain should detect function entities")
				}
			}
		})
	}
}

func TestQuickInsights(t *testing.T) {
	text := "Amazing product! 🚀 #innovation @company"
	
	start := time.Now()
	insights := QuickInsights(text)
	elapsed := time.Since(start)
	
	// Should be fast
	if elapsed > 50*time.Millisecond {
		t.Errorf("QuickInsights took too long: %v", elapsed)
	}
	
	// Should have sentiment
	if insights.Sentiment.Tone == "" {
		t.Error("Expected sentiment to be analyzed")
	}
	
	// Should have keywords
	if len(insights.TopKeywords) == 0 {
		t.Error("Expected keywords to be extracted")
	}
	
	// Should have summary
	if insights.Summary == "" {
		t.Error("Expected summary to be generated")
	}
}

func TestDeepTechnicalAnalysis(t *testing.T) {
	text := `
Here's a function to calculate factorial:

` + "```go" + `
func factorial(n int) int {
    if n <= 1 {
        return 1
    }
    return n * factorial(n-1)
}
` + "```" + `

This function has O(n) time complexity.
`
	
	result := DeepTechnicalAnalysis(text)
	
	// Should detect code blocks
	if len(result.CodeBlocks) == 0 {
		t.Error("Expected code blocks to be detected")
	}
	
	// Should analyze code
	if len(result.CodeAnalysis) == 0 {
		t.Error("Expected code analysis")
	}
	
	// Should extract technical terms
	if len(result.TechnicalTerms) == 0 {
		t.Error("Expected technical terms to be extracted")
	}
	
	// Should have complexity score
	if result.Complexity.Score == 0 {
		t.Error("Expected complexity to be calculated")
	}
}

func TestPerformanceComparison(t *testing.T) {
	text := strings.Repeat("This is a test sentence. ", 100) // ~2500 chars
	
	// Method 1: Individual calls
	start1 := time.Now()
	_ = ExtractAdvancedEntities(text)
	_ = CalculateFleschReadingEase(text)
	_ = CalculateTextStatistics(text)
	method1Time := time.Since(start1)
	
	// Method 2: SmartAnalyze
	start2 := time.Now()
	result := SmartAnalyze(text)
	method2Time := time.Since(start2)
	
	// SmartAnalyze should be faster due to optimized processing
	if method2Time > method1Time {
		t.Logf("Warning: SmartAnalyze (%v) was slower than individual calls (%v)",
			method2Time, method1Time)
	}
	
	// But should provide comprehensive results
	if len(result.Entities) == 0 || result.Readability.FleschScore == 0 {
		t.Error("SmartAnalyze should provide comprehensive results")
	}
}

func TestValidateAndClean(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Smart quotes",
			input:    ""Hello" and 'world'",
			expected: "\"Hello\" and 'world'",
		},
		{
			name:     "Multiple spaces",
			input:    "Too    many     spaces",
			expected: "Too many spaces",
		},
		{
			name:     "Trim whitespace",
			input:    "  text with spaces  ",
			expected: "text with spaces",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := validateAndClean(tt.input)
			if result != tt.expected {
				t.Errorf("Expected '%s', got '%s'", tt.expected, result)
			}
		})
	}
}

func TestAnalyzeReadabilityEnhanced(t *testing.T) {
	text := "The quick brown fox jumps over the lazy dog. This is a simple sentence."
	
	score := analyzeReadabilityEnhanced(text)
	
	// Check all fields are populated
	if score.FleschScore == 0 {
		t.Error("Expected Flesch score to be calculated")
	}
	
	if score.GradeLevel == 0 {
		t.Error("Expected grade level to be calculated")
	}
	
	if score.Complexity == "" {
		t.Error("Expected complexity to be determined")
	}
	
	if score.ReadingTime == 0 {
		t.Error("Expected reading time to be calculated")
	}
}

func TestExtractKeywordsOptimized(t *testing.T) {
	text := "OpenAI released GPT-4. OpenAI is an AI research company. GPT-4 is impressive."
	entities := []Entity{
		{Type: "ORGANIZATION", Text: "OpenAI"},
		{Type: "PRODUCT", Text: "GPT-4"},
	}
	
	keywords := extractKeywordsOptimized(text, entities)
	
	// Entity words should be boosted
	openAIRank := -1
	gpt4Rank := -1
	
	for i, keyword := range keywords {
		if keyword == "openai" {
			openAIRank = i
		}
		if keyword == "gpt-4" {
			gpt4Rank = i
		}
	}
	
	// Entity keywords should rank high (in top 5)
	if openAIRank > 4 || openAIRank == -1 {
		t.Error("Entity 'OpenAI' should rank high in keywords")
	}
}

func TestAnalyzeSentiment(t *testing.T) {
	tests := []struct {
		name         string
		text         string
		expectedTone string
	}{
		{
			name:         "Positive",
			text:         "This is a great product! I love it.",
			expectedTone: "positive",
		},
		{
			name:         "Negative",
			text:         "Terrible experience. Very disappointed.",
			expectedTone: "negative",
		},
		{
			name:         "Neutral",
			text:         "The meeting is at 3pm.",
			expectedTone: "neutral",
		},
		{
			name:         "Mixed",
			text:         "Good features but terrible support.",
			expectedTone: "mixed",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := analyzeSentiment(tt.text)
			
			if result.Tone != tt.expectedTone {
				t.Errorf("Expected tone '%s', got '%s'", 
					tt.expectedTone, result.Tone)
			}
			
			// Check confidence is reasonable
			if result.Confidence < 0 || result.Confidence > 1 {
				t.Errorf("Confidence should be between 0 and 1, got %f",
					result.Confidence)
			}
		})
	}
}

func TestDomainSpecificFeatures(t *testing.T) {
	t.Run("Technical code detection", func(t *testing.T) {
		text := "```go\nfunc main() {}\n```"
		blocks := detectCodeBlocks(text)
		
		if len(blocks) != 1 {
			t.Errorf("Expected 1 code block, got %d", len(blocks))
		}
		
		if blocks[0].Language != "go" {
			t.Errorf("Expected language 'go', got '%s'", blocks[0].Language)
		}
	})
	
	t.Run("Medical term extraction", func(t *testing.T) {
		text := "Patient takes 50mg aspirin. BP is stable."
		entities := extractMedicalTerms(text)
		
		hasDosage := false
		hasAbbreviation := false
		
		for _, e := range entities {
			if e.Type == "DOSAGE" {
				hasDosage = true
			}
			if e.Type == "ABBREVIATION" && e.Text == "BP" {
				hasAbbreviation = true
			}
		}
		
		if !hasDosage {
			t.Error("Expected to find dosage entity")
		}
		if !hasAbbreviation {
			t.Error("Expected to find medical abbreviation")
		}
	})
	
	t.Run("Social media entities", func(t *testing.T) {
		text := "#AI is trending! Follow @techguru"
		entities := extractSocialEntities(text)
		
		hasHashtag := false
		hasMention := false
		
		for _, e := range entities {
			if e.Type == "HASHTAG" {
				hasHashtag = true
			}
			if e.Type == "MENTION" {
				hasMention = true
			}
		}
		
		if !hasHashtag {
			t.Error("Expected to find hashtag")
		}
		if !hasMention {
			t.Error("Expected to find mention")
		}
	})
	
	t.Run("Business action items", func(t *testing.T) {
		text := "Action: Implement new CRM. TODO: Review Q3 reports."
		actions := extractActionItems(text)
		
		if len(actions) < 2 {
			t.Errorf("Expected at least 2 action items, got %d", len(actions))
		}
	})
}

// Benchmark tests

func BenchmarkSmartAnalyze(b *testing.B) {
	text := "Dr. Smith from OpenAI announced a breakthrough in artificial intelligence. " +
		"The new model achieves 95% accuracy on complex reasoning tasks."
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = SmartAnalyze(text)
	}
}

func BenchmarkQuickInsights(b *testing.B) {
	text := "Amazing product launch! #innovation @techcrunch"
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = QuickInsights(text)
	}
}

func BenchmarkValidatedExtraction(b *testing.B) {
	text := "OpenAI, Google, and Microsoft are leading AI research."
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = ValidatedExtraction(text)
	}
}