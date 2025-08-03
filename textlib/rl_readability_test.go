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
	"math"
	"strings"
	"testing"
)

func TestCalculateReadabilityMetrics(t *testing.T) {
	tests := []struct {
		name              string
		text              string
		algorithms        []string
		expectedScores    map[string]float64 // Expected ranges
		expectedAudiences []string
		checkSuggestions  bool
	}{
		{
			name: "Simple text",
			text: "The cat sat on the mat. The dog ran fast.",
			algorithms: []string{"flesch", "gunning-fog"},
			expectedScores: map[string]float64{
				"flesch":      90.0, // Very easy
				"gunning-fog": 4.0,  // Elementary
			},
			expectedAudiences: []string{"elementary-school"},
			checkSuggestions:  true,
		},
		{
			name: "Complex academic text",
			text: `The epistemological ramifications of quantum entanglement necessitate 
			       a fundamental reconsideration of our ontological presuppositions regarding 
			       the nature of reality and causality.`,
			algorithms: []string{"flesch-kincaid", "gunning-fog", "coleman-liau"},
			expectedScores: map[string]float64{
				"flesch-kincaid": 25.0, // Graduate level
				"gunning-fog":    25.0, // Very difficult
				"coleman-liau":   25.0, // College+
			},
			expectedAudiences: []string{"graduate", "professional", "college"},
			checkSuggestions:  true,
		},
		{
			name: "Standard business text",
			text: `Our company provides innovative solutions for business challenges. 
			       We help organizations improve their operational efficiency through 
			       strategic consulting and technology implementation.`,
			algorithms: []string{"all"},
			expectedScores: map[string]float64{
				"flesch":         10.0, // Very difficult (could be 0)
				"flesch-kincaid": 18.0, // College
				"gunning-fog":    22.0, // Graduate
			},
			expectedAudiences: []string{"college", "graduate"},
			checkSuggestions:  true,
		},
		{
			name:       "Empty text",
			text:       "",
			algorithms: []string{"flesch"},
			expectedScores: map[string]float64{},
			expectedAudiences: []string{},
		},
		{
			name:       "Default algorithms",
			text:       "Simple test text.",
			algorithms: []string{}, // Should default to flesch and gunning-fog
			expectedScores: map[string]float64{
				"flesch":      80.0,
				"gunning-fog": 3.0,
			},
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			report := CalculateReadabilityMetrics(tt.text, tt.algorithms)
			
			// Check that we got scores
			if tt.text != "" && len(report.Scores) == 0 {
				t.Error("No scores calculated for non-empty text")
			}
			
			// Check score ranges (approximate)
			for algo, expectedScore := range tt.expectedScores {
				if score, exists := report.Scores[algo]; exists {
					// Allow reasonable variance
					tolerance := 0.4 * math.Abs(expectedScore)
					if tolerance < 8.0 {
						tolerance = 8.0
					}
					
					// Special case for Flesch score which can be 0 or negative
					if algo == "flesch" && score <= 0 && expectedScore <= 30 {
						// This is acceptable for very difficult text
					} else if math.Abs(score-expectedScore) > tolerance {
						t.Errorf("%s score %.1f not in expected range (expected ~%.1f)", 
							algo, score, expectedScore)
					}
				} else if algo != "smog" { // SMOG might not be calculated for short texts
					t.Errorf("Expected %s score not found", algo)
				}
			}
			
			// Check audiences
			if len(tt.expectedAudiences) > 0 {
				foundAudience := false
				for _, expected := range tt.expectedAudiences {
					for _, actual := range report.TargetAudience {
						if actual == expected {
							foundAudience = true
							break
						}
					}
				}
				if !foundAudience {
					t.Errorf("Expected audiences %v, got %v", 
						tt.expectedAudiences, report.TargetAudience)
				}
			}
			
			// Check suggestions
			if tt.checkSuggestions && tt.text != "" {
				if len(report.ImprovementSuggestions) == 0 {
					t.Error("No improvement suggestions provided")
				}
			}
			
			// Check recommendation
			if tt.text != "" && report.Recommendation == "" {
				t.Error("No recommendation provided")
			}
			
			// Check processing cost
			if tt.text != "" {
				if report.ProcessingCost.TimeMs <= 0 {
					t.Error("Invalid processing time")
				}
				if report.ProcessingCost.CPUCycles <= 0 {
					t.Error("Invalid CPU cycles estimate")
				}
			}
		})
	}
}

func TestFleschReadingEase(t *testing.T) {
	// Test with known examples
	tests := []struct {
		text          string
		expectedRange [2]float64 // min, max
		description   string
	}{
		{
			text:          "See Spot run. Run, Spot, run. Jane sees Spot.",
			expectedRange: [2]float64{90, 100},
			description:   "Very easy text",
		},
		{
			text: "The implementation of sophisticated algorithms requires " +
				"comprehensive understanding of computational complexity theory.",
			expectedRange: [2]float64{0, 30},
			description:   "Difficult text",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			report := CalculateReadabilityMetrics(tt.text, []string{"flesch"})
			
			if score, exists := report.Scores["flesch"]; exists {
				if score < tt.expectedRange[0] || score > tt.expectedRange[1] {
					t.Errorf("Flesch score %.1f not in expected range [%.1f, %.1f]",
						score, tt.expectedRange[0], tt.expectedRange[1])
				}
			} else {
				t.Error("Flesch score not calculated")
			}
		})
	}
}

func TestAllAlgorithms(t *testing.T) {
	text := strings.Repeat("The quick brown fox jumps over the lazy dog. ", 10)
	
	report := CalculateReadabilityMetrics(text, []string{"all"})
	
	// Should have all algorithms
	expectedAlgos := []string{"flesch", "flesch-kincaid", "gunning-fog", 
		"coleman-liau", "ari", "smog"}
	
	for _, algo := range expectedAlgos {
		if _, exists := report.Scores[algo]; !exists {
			t.Errorf("Algorithm %s not calculated with 'all' option", algo)
		}
	}
	
	// All scores should be reasonable
	for algo, score := range report.Scores {
		if algo == "flesch" {
			if score < 0 || score > 100 {
				t.Errorf("%s score out of range: %.1f", algo, score)
			}
		} else if algo != "smog-estimated" {
			// Grade level scores
			if score < 0 || score > 30 {
				t.Errorf("%s score unreasonable: %.1f", algo, score)
			}
		}
	}
}

func TestSMOGCalculation(t *testing.T) {
	// SMOG requires 30+ sentences
	shortText := "Short text. Only two sentences."
	longText := strings.Repeat("This is a sentence with some polysyllabic words. ", 35)
	
	// Short text should use estimation
	shortReport := CalculateReadabilityMetrics(shortText, []string{"smog"})
	if _, exists := shortReport.Scores["smog-estimated"]; !exists {
		t.Error("SMOG should be estimated for short texts")
	}
	
	// Long text should calculate properly
	longReport := CalculateReadabilityMetrics(longText, []string{"smog"})
	if _, exists := longReport.Scores["smog-estimated"]; exists {
		t.Error("SMOG should not be estimated for long texts")
	}
	
	if score, exists := longReport.Scores["smog"]; exists {
		if score <= 0 {
			t.Errorf("Invalid SMOG score: %.1f", score)
		}
	} else {
		t.Error("SMOG score not calculated for long text")
	}
}

func TestImprovementSuggestions(t *testing.T) {
	tests := []struct {
		name               string
		text               string
		expectedSuggestion string
	}{
		{
			name: "Long sentences",
			text: "This is an extremely long sentence that contains many words and " +
				"continues on and on without any breaks which makes it very difficult " +
				"to read and understand for most readers.",
			expectedSuggestion: "sentence length",
		},
		{
			name: "Complex vocabulary",
			text: "The anthropomorphic characteristics of the protagonist exemplify " +
				"the dichotomous relationship between humanity and technology.",
			expectedSuggestion: "simpler vocabulary",
		},
		{
			name: "Good readability",
			text: "The cat sat on the mat. The dog played in the yard.",
			expectedSuggestion: "excellent readability",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			report := CalculateReadabilityMetrics(tt.text, []string{"flesch"})
			
			found := false
			for _, suggestion := range report.ImprovementSuggestions {
				if strings.Contains(strings.ToLower(suggestion), tt.expectedSuggestion) {
					found = true
					break
				}
			}
			
			if !found {
				t.Errorf("Expected suggestion about '%s' not found. Got: %v",
					tt.expectedSuggestion, report.ImprovementSuggestions)
			}
		})
	}
}

func TestTargetAudienceIdentification(t *testing.T) {
	tests := []struct {
		name             string
		text             string
		expectedAudience string
	}{
		{
			name:             "Elementary level",
			text:             "I like dogs. Dogs are fun. Cats are nice too.",
			expectedAudience: "elementary-school",
		},
		{
			name: "College level",
			text: "The theoretical framework necessitates a comprehensive analysis " +
				"of the underlying assumptions inherent in the methodological approach.",
			expectedAudience: "graduate",
		},
		{
			name: "General public",
			text: "Climate change affects our daily lives in many ways. " +
				"We can help by reducing energy use and recycling more.",
			expectedAudience: "general",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			report := CalculateReadabilityMetrics(tt.text, []string{"all"})
			
			found := false
			for _, audience := range report.TargetAudience {
				if strings.Contains(audience, tt.expectedAudience) {
					found = true
					break
				}
			}
			
			if !found {
				t.Errorf("Expected audience '%s' not found. Got: %v",
					tt.expectedAudience, report.TargetAudience)
			}
		})
	}
}

func BenchmarkCalculateReadabilityMetricsSimple(b *testing.B) {
	text := "The quick brown fox jumps over the lazy dog."
	algorithms := []string{"flesch", "gunning-fog"}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		CalculateReadabilityMetrics(text, algorithms)
	}
}

func BenchmarkCalculateReadabilityMetricsAll(b *testing.B) {
	text := strings.Repeat("This is a sample sentence for testing readability. ", 10)
	algorithms := []string{"all"}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		CalculateReadabilityMetrics(text, algorithms)
	}
}

func BenchmarkCalculateReadabilityMetricsLong(b *testing.B) {
	// Generate long text (40+ sentences for SMOG)
	text := strings.Repeat("The comprehensive analysis demonstrates significant findings. ", 50)
	algorithms := []string{"all"}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		CalculateReadabilityMetrics(text, algorithms)
	}
}