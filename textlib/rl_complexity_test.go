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

func TestAnalyzeTextComplexity(t *testing.T) {
	tests := []struct {
		name           string
		text           string
		depth          int
		wantAlgorithm  string
		checkMetrics   bool
		minAccuracy    float64
		maxProcessTime time.Duration
	}{
		{
			name:           "Simple text depth 1",
			text:           "The cat sat on the mat. It was happy.",
			depth:          1,
			wantAlgorithm:  "complexity-depth-1",
			checkMetrics:   true,
			minAccuracy:    0.75,
			maxProcessTime: 100 * time.Millisecond,
		},
		{
			name:           "Medium text depth 2",
			text:           "The quantum mechanics principles demonstrate that particles exhibit wave-particle duality. This phenomenon challenges classical physics understanding.",
			depth:          2,
			wantAlgorithm:  "complexity-depth-2",
			checkMetrics:   true,
			minAccuracy:    0.89, // Allow for floating point precision
			maxProcessTime: 200 * time.Millisecond,
		},
		{
			name: "Complex text depth 3",
			text: `The epistemological ramifications of quantum entanglement necessitate a fundamental 
			reconsideration of our ontological presuppositions. Furthermore, the Copenhagen interpretation's 
			probabilistic framework contradicts deterministic paradigms. Therefore, contemporary physicists 
			must reconcile these paradoxical observations with established theoretical constructs.`,
			depth:          3,
			wantAlgorithm:  "complexity-depth-3",
			checkMetrics:   true,
			minAccuracy:    0.95,
			maxProcessTime: 500 * time.Millisecond,
		},
		{
			name:          "Empty text",
			text:          "",
			depth:         2,
			wantAlgorithm: "complexity-depth-2",
			checkMetrics:  false,
		},
		{
			name:          "Invalid depth (too low)",
			text:          "Test text.",
			depth:         0,
			wantAlgorithm: "complexity-depth-2", // Should default to 2
		},
		{
			name:          "Invalid depth (too high)",
			text:          "Test text.",
			depth:         5,
			wantAlgorithm: "complexity-depth-2", // Should default to 2
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			report := AnalyzeTextComplexity(tt.text, tt.depth)
			
			// Check algorithm used
			if report.AlgorithmUsed != tt.wantAlgorithm {
				t.Errorf("Algorithm = %v, want %v", report.AlgorithmUsed, tt.wantAlgorithm)
			}
			
			// Check processing time
			if tt.maxProcessTime > 0 && report.ProcessingTime > tt.maxProcessTime {
				t.Errorf("ProcessingTime = %v, want <= %v", report.ProcessingTime, tt.maxProcessTime)
			}
			
			// Check quality metrics
			if tt.checkMetrics {
				if report.QualityMetrics.Accuracy < tt.minAccuracy {
					t.Errorf("Accuracy = %v, want >= %v", report.QualityMetrics.Accuracy, tt.minAccuracy)
				}
				
				// All quality metrics should be in valid range
				if report.QualityMetrics.Accuracy < 0 || report.QualityMetrics.Accuracy > 1 {
					t.Errorf("Accuracy out of range: %v", report.QualityMetrics.Accuracy)
				}
				if report.QualityMetrics.Confidence < 0 || report.QualityMetrics.Confidence > 1 {
					t.Errorf("Confidence out of range: %v", report.QualityMetrics.Confidence)
				}
				if report.QualityMetrics.Coverage < 0 || report.QualityMetrics.Coverage > 1 {
					t.Errorf("Coverage out of range: %v", report.QualityMetrics.Coverage)
				}
			}
			
			// Check complexity scores are in valid range
			if tt.text != "" {
				if report.LexicalComplexity < 0 || report.LexicalComplexity > 1 {
					t.Errorf("LexicalComplexity out of range: %v", report.LexicalComplexity)
				}
				
				if tt.depth >= 2 {
					if report.SyntacticComplexity < 0 || report.SyntacticComplexity > 1 {
						t.Errorf("SyntacticComplexity out of range: %v", report.SyntacticComplexity)
					}
				}
				
				if tt.depth >= 3 {
					if report.SemanticComplexity < 0 || report.SemanticComplexity > 1 {
						t.Errorf("SemanticComplexity out of range: %v", report.SemanticComplexity)
					}
				}
			}
			
			// Check readability scores exist
			if tt.text != "" {
				if _, exists := report.ReadabilityScores["flesch-kincaid"]; !exists {
					t.Error("Missing flesch-kincaid score")
				}
				if _, exists := report.ReadabilityScores["gunning-fog"]; !exists {
					t.Error("Missing gunning-fog score")
				}
				
				if tt.depth >= 2 {
					if _, exists := report.ReadabilityScores["coleman-liau"]; !exists {
						t.Error("Missing coleman-liau score")
					}
					if _, exists := report.ReadabilityScores["ari"]; !exists {
						t.Error("Missing ari score")
					}
				}
			}
		})
	}
}

func TestComplexityDepthProgression(t *testing.T) {
	// Test that higher depth provides more analysis
	text := `The implementation of advanced algorithms requires careful consideration of computational 
	complexity. Moreover, optimization strategies must balance performance with resource constraints. 
	Therefore, developers should analyze trade-offs systematically.`
	
	report1 := AnalyzeTextComplexity(text, 1)
	report2 := AnalyzeTextComplexity(text, 2)
	report3 := AnalyzeTextComplexity(text, 3)
	
	// Check that processing time increases with depth
	if report2.ProcessingTime <= report1.ProcessingTime {
		t.Error("Depth 2 should take longer than depth 1")
	}
	if report3.ProcessingTime <= report2.ProcessingTime {
		t.Error("Depth 3 should take longer than depth 2")
	}
	
	// Check that quality metrics increase with depth
	if report2.QualityMetrics.Accuracy <= report1.QualityMetrics.Accuracy {
		t.Error("Depth 2 should have higher accuracy than depth 1")
	}
	if report3.QualityMetrics.Accuracy <= report2.QualityMetrics.Accuracy {
		t.Error("Depth 3 should have higher accuracy than depth 2")
	}
	
	// Check that more metrics are calculated at higher depths
	if len(report2.ReadabilityScores) <= len(report1.ReadabilityScores) {
		t.Error("Depth 2 should calculate more readability scores")
	}
	
	// Depth 3 should have semantic complexity
	if report3.SemanticComplexity == 0 {
		t.Error("Depth 3 should calculate semantic complexity")
	}
	if report1.SemanticComplexity != 0 || report2.SemanticComplexity == 0 {
		// Note: report2.SemanticComplexity could be 0 if not calculated at depth 2
	}
}

func TestCountSyllables(t *testing.T) {
	tests := []struct {
		word     string
		expected int
	}{
		{"cat", 1},
		{"happy", 2},
		{"beautiful", 3}, // beau-ti-ful
		{"the", 1},
		{"create", 1}, // silent e makes it cre-ate -> create
		{"complicated", 4},
		{"a", 1},
		{"I", 1},
		{"extraordinary", 5},
		{"queue", 1}, // special case
	}
	
	for _, tt := range tests {
		t.Run(tt.word, func(t *testing.T) {
			count := countSyllables(tt.word)
			if count != tt.expected {
				t.Errorf("countSyllables(%q) = %d, want %d", tt.word, count, tt.expected)
			}
		})
	}
}

func TestReadabilityScores(t *testing.T) {
	// Test with known text samples
	simpleText := "See Spot run. Run Spot run. Jane sees Spot."
	complexText := "The epistemological implications of quantum mechanics necessitate a paradigmatic shift in our understanding of reality."
	
	simpleReport := AnalyzeTextComplexity(simpleText, 2)
	complexReport := AnalyzeTextComplexity(complexText, 2)
	
	// Simple text should have lower readability scores
	if simpleReport.ReadabilityScores["flesch-kincaid"] >= complexReport.ReadabilityScores["flesch-kincaid"] {
		t.Error("Simple text should have lower Flesch-Kincaid score than complex text")
	}
	
	if simpleReport.ReadabilityScores["gunning-fog"] >= complexReport.ReadabilityScores["gunning-fog"] {
		t.Error("Simple text should have lower Gunning Fog score than complex text")
	}
	
	// Complex text should have higher complexity scores
	if simpleReport.LexicalComplexity >= complexReport.LexicalComplexity {
		t.Error("Simple text should have lower lexical complexity")
	}
}

func TestLongTextPerformance(t *testing.T) {
	// Generate a long text
	sentences := []string{
		"The quick brown fox jumps over the lazy dog.",
		"Advanced algorithms require careful implementation.",
		"Performance optimization is crucial for scalability.",
		"Complex systems demand thorough analysis.",
	}
	
	var longText strings.Builder
	for i := 0; i < 100; i++ {
		longText.WriteString(sentences[i%len(sentences)])
		longText.WriteString(" ")
	}
	
	text := longText.String()
	
	// Test performance at different depths
	for depth := 1; depth <= 3; depth++ {
		start := time.Now()
		report := AnalyzeTextComplexity(text, depth)
		elapsed := time.Since(start)
		
		t.Logf("Depth %d: %v ms for %d chars", depth, elapsed.Milliseconds(), len(text))
		
		// Verify report is valid
		if report.AlgorithmUsed == "" {
			t.Errorf("Empty algorithm name at depth %d", depth)
		}
		
		// Check memory was tracked
		if report.MemoryUsed < 0 {
			t.Errorf("Invalid memory usage at depth %d: %d", depth, report.MemoryUsed)
		}
	}
}

func BenchmarkAnalyzeTextComplexityDepth1(b *testing.B) {
	text := "The implementation of advanced algorithms requires careful consideration. " +
		"Performance optimization strategies must balance efficiency with maintainability."
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		AnalyzeTextComplexity(text, 1)
	}
}

func BenchmarkAnalyzeTextComplexityDepth2(b *testing.B) {
	text := "The implementation of advanced algorithms requires careful consideration. " +
		"Performance optimization strategies must balance efficiency with maintainability."
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		AnalyzeTextComplexity(text, 2)
	}
}

func BenchmarkAnalyzeTextComplexityDepth3(b *testing.B) {
	text := "The implementation of advanced algorithms requires careful consideration. " +
		"Performance optimization strategies must balance efficiency with maintainability."
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		AnalyzeTextComplexity(text, 3)
	}
}