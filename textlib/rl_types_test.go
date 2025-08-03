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
	"encoding/json"
	"testing"
	"time"
)

func TestProcessingMetrics(t *testing.T) {
	metrics := ProcessingMetrics{
		TimeElapsed:    100 * time.Millisecond,
		MemoryPeak:     1024 * 1024, // 1MB
		AlgorithmSteps: 42,
		CacheHits:      5,
	}
	
	// Test JSON serialization
	data, err := json.Marshal(metrics)
	if err != nil {
		t.Fatalf("Failed to marshal ProcessingMetrics: %v", err)
	}
	
	var decoded ProcessingMetrics
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Failed to unmarshal ProcessingMetrics: %v", err)
	}
	
	if decoded.AlgorithmSteps != metrics.AlgorithmSteps {
		t.Errorf("AlgorithmSteps mismatch: got %d, want %d", decoded.AlgorithmSteps, metrics.AlgorithmSteps)
	}
	
	if decoded.CacheHits != metrics.CacheHits {
		t.Errorf("CacheHits mismatch: got %d, want %d", decoded.CacheHits, metrics.CacheHits)
	}
}

func TestQualityMetrics(t *testing.T) {
	quality := QualityMetrics{
		Accuracy:   0.95,
		Confidence: 0.87,
		Coverage:   0.92,
	}
	
	// Validate ranges
	if quality.Accuracy < 0 || quality.Accuracy > 1 {
		t.Errorf("Accuracy out of range: %f", quality.Accuracy)
	}
	
	if quality.Confidence < 0 || quality.Confidence > 1 {
		t.Errorf("Confidence out of range: %f", quality.Confidence)
	}
	
	if quality.Coverage < 0 || quality.Coverage > 1 {
		t.Errorf("Coverage out of range: %f", quality.Coverage)
	}
}

func TestComplexityReport(t *testing.T) {
	report := ComplexityReport{
		LexicalComplexity:   0.65,
		SyntacticComplexity: 0.72,
		SemanticComplexity:  0.58,
		ReadabilityScores: map[string]float64{
			"flesch-kincaid": 8.2,
			"gunning-fog":    10.5,
		},
		ProcessingTime: 150 * time.Millisecond,
		MemoryUsed:     2 * 1024 * 1024, // 2MB
		AlgorithmUsed:  "complexity-deep",
		QualityMetrics: QualityMetrics{
			Accuracy:   0.90,
			Confidence: 0.85,
			Coverage:   0.88,
		},
	}
	
	// Test readability scores access
	if score, exists := report.ReadabilityScores["flesch-kincaid"]; !exists || score != 8.2 {
		t.Errorf("flesch-kincaid score incorrect: %f", score)
	}
	
	// Test algorithm name
	if report.AlgorithmUsed != "complexity-deep" {
		t.Errorf("Algorithm name incorrect: %s", report.AlgorithmUsed)
	}
}

func TestSentimentAnalysis(t *testing.T) {
	sentiment := SentimentAnalysis{
		OverallSentiment: Sentiment{
			Polarity:   0.75,
			Magnitude:  0.82,
			Label:      "positive",
			Confidence: 0.91,
		},
		SentenceLevel: []SentenceSentiment{
			{
				Text: "This is great!",
				Sentiment: Sentiment{
					Polarity:   0.9,
					Magnitude:  0.95,
					Label:      "positive",
					Confidence: 0.98,
				},
				Position: RLPosition{Start: 0, End: 14, Line: 1},
			},
		},
		AspectBased: map[string]Sentiment{
			"quality": {
				Polarity:   0.8,
				Magnitude:  0.85,
				Label:      "positive",
				Confidence: 0.88,
			},
		},
		EmotionProfile: EmotionProfile{
			Joy:      0.75,
			Anger:    0.05,
			Fear:     0.02,
			Sadness:  0.03,
			Surprise: 0.10,
			Trust:    0.70,
		},
		Confidence:         0.89,
		ProcessingApproach: "ml-advanced",
	}
	
	// Validate polarity range
	if sentiment.OverallSentiment.Polarity < -1 || sentiment.OverallSentiment.Polarity > 1 {
		t.Errorf("Polarity out of range: %f", sentiment.OverallSentiment.Polarity)
	}
	
	// Validate magnitude range
	if sentiment.OverallSentiment.Magnitude < 0 || sentiment.OverallSentiment.Magnitude > 1 {
		t.Errorf("Magnitude out of range: %f", sentiment.OverallSentiment.Magnitude)
	}
	
	// Check emotion profile sum (should be close to 1 for normalized emotions)
	emotionSum := sentiment.EmotionProfile.Joy + sentiment.EmotionProfile.Anger +
		sentiment.EmotionProfile.Fear + sentiment.EmotionProfile.Sadness +
		sentiment.EmotionProfile.Surprise + sentiment.EmotionProfile.Trust
	
	// Allow some tolerance for rounding
	if emotionSum < 0.9 || emotionSum > 2.0 {
		t.Errorf("Emotion profile sum unusual: %f", emotionSum)
	}
}

func TestDocumentAnalysis(t *testing.T) {
	doc := DocumentAnalysis{
		TextAnalysis: ComplexityReport{
			LexicalComplexity: 0.7,
			AlgorithmUsed:     "complexity-basic",
		},
		StructureAnalysis: StructureAnalysis{
			DocumentType: "article",
			Sections: []RLSection{
				{
					Title:     "Introduction",
					Level:     1,
					Content:   "This is the introduction.",
					Position:  RLPosition{Start: 0, End: 25},
					WordCount: 4,
				},
			},
			Tables: []Table{
				{
					Caption: "Results",
					Headers: []string{"Name", "Value"},
					Rows: [][]string{
						{"Test1", "100"},
						{"Test2", "200"},
					},
				},
			},
		},
		MetadataExtraction: Metadata{
			Title:    "Test Document",
			Author:   "Test Author",
			Language: "en",
			Keywords: []string{"test", "document", "analysis"},
		},
		QualityAssessment: QualityAssessment{
			OverallScore:      0.85,
			ReadabilityScore:  0.82,
			CompletenessScore: 0.88,
			ConsistencyScore:  0.90,
			Issues: []QualityIssue{
				{
					Type:        "spelling",
					Severity:    "low",
					Description: "Minor spelling issue",
					Location:    RLPosition{Start: 50, End: 55},
					Suggestion:  "Check spelling",
				},
			},
			Recommendations: []string{"Add more examples", "Clarify terminology"},
		},
		ProcessingStrategy: "balanced",
		Performance: PerformanceMetrics{
			TotalTime: 250 * time.Millisecond,
			StepTimings: map[string]time.Duration{
				"parsing":    50 * time.Millisecond,
				"analysis":   150 * time.Millisecond,
				"formatting": 50 * time.Millisecond,
			},
			MemoryUsage:      5 * 1024 * 1024, // 5MB
			CacheUtilization: 0.65,
		},
	}
	
	// Test document type
	if doc.StructureAnalysis.DocumentType != "article" {
		t.Errorf("Document type incorrect: %s", doc.StructureAnalysis.DocumentType)
	}
	
	// Test table structure
	if len(doc.StructureAnalysis.Tables) != 1 {
		t.Errorf("Expected 1 table, got %d", len(doc.StructureAnalysis.Tables))
	} else if len(doc.StructureAnalysis.Tables[0].Rows) != 2 {
		t.Errorf("Expected 2 rows, got %d", len(doc.StructureAnalysis.Tables[0].Rows))
	}
	
	// Test quality scores
	if doc.QualityAssessment.OverallScore < 0 || doc.QualityAssessment.OverallScore > 1 {
		t.Errorf("Overall score out of range: %f", doc.QualityAssessment.OverallScore)
	}
}

func TestBatchProcessing(t *testing.T) {
	batch := BatchResult{
		Results: []interface{}{
			"Result 1",
			"Result 2",
			"Result 3",
		},
		OverallMetrics: OverallMetrics{
			TotalProcessed:   3,
			SuccessCount:     2,
			ErrorCount:       1,
			AverageTime:      100 * time.Millisecond,
			TotalTime:        300 * time.Millisecond,
			ThroughputPerSec: 10.0,
		},
		ProcessingStrategy: BatchStrategy{
			Parallel:  true,
			BatchSize: 10,
			Workers:   4,
			Priority:  "speed",
		},
		Errors: []ProcessingError{
			{
				Index:       2,
				Error:       "Processing failed",
				InputSample: "Sample input",
				Timestamp:   time.Now(),
			},
		},
	}
	
	// Test metrics consistency
	if batch.OverallMetrics.TotalProcessed != batch.OverallMetrics.SuccessCount+batch.OverallMetrics.ErrorCount {
		t.Errorf("Metrics inconsistency: total=%d, success=%d, error=%d",
			batch.OverallMetrics.TotalProcessed,
			batch.OverallMetrics.SuccessCount,
			batch.OverallMetrics.ErrorCount)
	}
	
	// Test strategy
	if !batch.ProcessingStrategy.Parallel {
		t.Error("Expected parallel processing")
	}
	
	if batch.ProcessingStrategy.Workers < 1 {
		t.Errorf("Invalid worker count: %d", batch.ProcessingStrategy.Workers)
	}
}

func TestOptimizationMetrics(t *testing.T) {
	metrics := OptimizationMetrics{
		QualityScore:     0.85,
		PerformanceScore: 0.75,
		ResourceScore:    0.90,
		UserSatisfaction: 0.88,
		WeightedTotal:    0.845,
	}
	
	// Calculate expected weighted total (assuming equal weights)
	expectedTotal := (metrics.QualityScore + metrics.PerformanceScore + 
		metrics.ResourceScore + metrics.UserSatisfaction) / 4.0
	
	// Allow small tolerance for floating point
	tolerance := 0.01
	if diff := expectedTotal - metrics.WeightedTotal; diff > tolerance {
		t.Errorf("Weighted total mismatch: expected ~%f, got %f", expectedTotal, metrics.WeightedTotal)
	}
	
	// Validate all scores are in range
	if metrics.QualityScore < 0 || metrics.QualityScore > 1 {
		t.Errorf("QualityScore out of range: %f", metrics.QualityScore)
	}
	
	if metrics.PerformanceScore < 0 || metrics.PerformanceScore > 1 {
		t.Errorf("PerformanceScore out of range: %f", metrics.PerformanceScore)
	}
	
	if metrics.ResourceScore < 0 || metrics.ResourceScore > 1 {
		t.Errorf("ResourceScore out of range: %f", metrics.ResourceScore)
	}
	
	if metrics.UserSatisfaction < 0 || metrics.UserSatisfaction > 1 {
		t.Errorf("UserSatisfaction out of range: %f", metrics.UserSatisfaction)
	}
}

func TestProcessingStrategy(t *testing.T) {
	strategy := ProcessingStrategy{
		Name:        "balanced",
		Description: "Balanced processing for general text",
		Parameters: map[string]interface{}{
			"depth":      2,
			"algorithms": []string{"flesch", "gunning-fog"},
			"quality":    0.85,
		},
		ExpectedQuality: 0.85,
		ExpectedSpeed:   0.70,
		ResourceRequirements: ResourceRequirements{
			MinMemoryMB:      50,
			MaxMemoryMB:      200,
			EstimatedCPUTime: 150,
			NetworkRequired:  false,
			CacheRecommended: true,
		},
	}
	
	// Test parameter access
	if depth, ok := strategy.Parameters["depth"].(int); !ok || depth != 2 {
		t.Errorf("Depth parameter incorrect: %v", strategy.Parameters["depth"])
	}
	
	// Test resource requirements
	if strategy.ResourceRequirements.MinMemoryMB > strategy.ResourceRequirements.MaxMemoryMB {
		t.Errorf("Invalid memory requirements: min=%d, max=%d",
			strategy.ResourceRequirements.MinMemoryMB,
			strategy.ResourceRequirements.MaxMemoryMB)
	}
	
	// Test expected metrics
	if strategy.ExpectedQuality+strategy.ExpectedSpeed > 2.0 {
		t.Error("Expected metrics out of range")
	}
}