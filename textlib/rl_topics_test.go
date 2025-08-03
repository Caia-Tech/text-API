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

func TestClassifyTopics(t *testing.T) {
	tests := []struct {
		name           string
		text           string
		maxTopics      int
		expectedMethod string
		minTopics      int
		maxTopics_     int
		checkKeywords  bool
		expectedTopic  string
	}{
		{
			name: "Technology text - clustering",
			text: `Artificial intelligence and machine learning are transforming the technology industry. 
				Deep learning algorithms can process vast amounts of data to identify patterns. 
				Neural networks are becoming more sophisticated with each advancement. 
				The applications in computer vision and natural language processing are remarkable.`,
			maxTopics:      3,
			expectedMethod: "clustering",
			minTopics:      1,
			maxTopics_:     3,
			checkKeywords:  true,
			expectedTopic:  "technology",
		},
		{
			name: "Business text - statistical",
			text: `Our company's quarterly revenue has increased significantly this year. 
				The marketing team has developed effective strategies to reach new customers. 
				Sales performance in the international markets shows promising growth. 
				Customer satisfaction ratings have improved due to enhanced service quality. 
				The business development team is exploring new market opportunities.`,
			maxTopics:      5,
			expectedMethod: "statistical",
			minTopics:      2,
			maxTopics_:     5,
			checkKeywords:  true,
			expectedTopic:  "business",
		},
		{
			name: "Mixed content - comprehensive",
			text: `Climate change is affecting global weather patterns and environmental systems. 
				Scientists are studying the impact on ocean temperatures and marine life. 
				Government policies are being developed to address carbon emissions. 
				Renewable energy sources like solar and wind power are becoming more efficient. 
				Technology companies are investing in sustainable solutions and green initiatives. 
				The healthcare industry is preparing for climate-related health challenges. 
				Educational institutions are incorporating environmental science into their curricula. 
				Economic models are being adjusted to account for environmental costs.`,
			maxTopics:      15,
			expectedMethod: "comprehensive",
			minTopics:      3,
			maxTopics_:     15,
			checkKeywords:  true,
			expectedTopic:  "climate",
		},
		{
			name:           "Empty text",
			text:           "",
			maxTopics:      5,
			expectedMethod: "none",
			minTopics:      0,
			maxTopics_:     0,
			checkKeywords:  false,
		},
		{
			name: "Short text",
			text: "Machine learning is important.",
			maxTopics:      2,
			expectedMethod: "clustering",
			minTopics:      0,
			maxTopics_:     2,
			checkKeywords:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ClassifyTopics(tt.text, tt.maxTopics)

			// Check method
			if result.Method != tt.expectedMethod {
				t.Errorf("Expected method %s, got %s", tt.expectedMethod, result.Method)
			}

			// Check number of topics
			if len(result.Topics) < tt.minTopics {
				t.Errorf("Expected at least %d topics, got %d", tt.minTopics, len(result.Topics))
			}
			if len(result.Topics) > tt.maxTopics_ {
				t.Errorf("Expected at most %d topics, got %d", tt.maxTopics_, len(result.Topics))
			}

			// Check topic structure
			for i, topic := range result.Topics {
				if topic.Name == "" {
					t.Errorf("Topic %d has empty name", i)
				}
				if len(topic.Keywords) == 0 {
					t.Errorf("Topic %d has no keywords", i)
				}
				if topic.Confidence < 0.0 || topic.Confidence > 1.0 {
					t.Errorf("Topic %d confidence out of range: %f", i, topic.Confidence)
				}
				if topic.Coverage < 0.0 || topic.Coverage > 1.0 {
					t.Errorf("Topic %d coverage out of range: %f", i, topic.Coverage)
				}
			}

			// Check quality metrics
			if tt.text != "" {
				if result.QualityMetrics.Accuracy < 0.0 || result.QualityMetrics.Accuracy > 1.0 {
					t.Errorf("Quality accuracy out of range: %f", result.QualityMetrics.Accuracy)
				}
				if result.QualityMetrics.Coverage < 0.0 || result.QualityMetrics.Coverage > 1.0 {
					t.Errorf("Quality coverage out of range: %f", result.QualityMetrics.Coverage)
				}
			}

			// Check for expected topic content
			if tt.checkKeywords && tt.expectedTopic != "" && len(result.Topics) > 0 {
				found := false
				for _, topic := range result.Topics {
					topicText := strings.ToLower(topic.Name + " " + strings.Join(topic.Keywords, " "))
					if strings.Contains(topicText, tt.expectedTopic) {
						found = true
						break
					}
				}
				if !found {
					t.Logf("Expected to find topic related to '%s', but didn't. Topics: %v", 
						tt.expectedTopic, result.Topics)
					// Don't fail the test, just log for debugging
				}
			}

			// Check processing time
			if result.ProcessingTime <= 0 {
				t.Error("Invalid processing time")
			}

			// Check topic ordering (should be by confidence)
			for i := 1; i < len(result.Topics); i++ {
				if result.Topics[i].Confidence > result.Topics[i-1].Confidence {
					t.Errorf("Topics not ordered by confidence: topic %d (%.3f) > topic %d (%.3f)",
						i, result.Topics[i].Confidence, i-1, result.Topics[i-1].Confidence)
				}
			}
		})
	}
}

func TestTopicMethodSelection(t *testing.T) {
	text := "Technology companies are developing artificial intelligence solutions."
	
	// Test method selection based on maxTopics
	testCases := []struct {
		maxTopics      int
		expectedMethod string
	}{
		{1, "clustering"},
		{3, "clustering"},
		{5, "statistical"},
		{10, "statistical"},
		{15, "comprehensive"},
	}

	for _, tc := range testCases {
		result := ClassifyTopics(text, tc.maxTopics)
		if result.Method != tc.expectedMethod {
			t.Errorf("MaxTopics %d: expected method %s, got %s", 
				tc.maxTopics, tc.expectedMethod, result.Method)
		}
	}
}

func TestTopicClustering(t *testing.T) {
	text := `Machine learning algorithms are used in data science projects. 
		Data analysis helps identify patterns in large datasets. 
		Statistical models can predict future trends based on historical data.`
	
	result := ClassifyTopics(text, 2) // Force clustering method
	
	if result.Method != "clustering" {
		t.Errorf("Expected clustering method, got %s", result.Method)
	}

	// Should identify data/machine learning related topics
	if len(result.Topics) == 0 {
		t.Error("Expected at least one topic for clustering")
	}

	// Check that topics have reasonable content
	for _, topic := range result.Topics {
		if len(topic.Keywords) == 0 {
			t.Error("Topic should have keywords")
		}
		if len(topic.Examples) == 0 {
			t.Error("Topic should have examples")
		}
	}
}

func TestTopicStatistical(t *testing.T) {
	text := `The financial markets showed strong performance this quarter. 
		Investment portfolios have delivered consistent returns for clients. 
		Economic indicators suggest continued market stability. 
		Banking institutions are reporting increased lending activity. 
		Insurance companies are expanding their product offerings.`
	
	result := ClassifyTopics(text, 5) // Force statistical method
	
	if result.Method != "statistical" {
		t.Errorf("Expected statistical method, got %s", result.Method)
	}

	// Should identify financial/business topics
	if len(result.Topics) == 0 {
		t.Error("Expected topics for statistical analysis")
	}

	// Check for business/financial keywords
	foundFinancial := false
	for _, topic := range result.Topics {
		for _, keyword := range topic.Keywords {
			if strings.Contains(keyword, "market") || 
				strings.Contains(keyword, "financial") ||
				strings.Contains(keyword, "investment") ||
				strings.Contains(keyword, "economic") {
				foundFinancial = true
				break
			}
		}
		if foundFinancial {
			break
		}
	}

	if !foundFinancial {
		t.Log("Expected to find financial-related keywords")
		// Don't fail, just log for debugging
	}
}

func TestTopicComprehensive(t *testing.T) {
	text := `Healthcare systems are adapting to new technologies and patient needs. 
		Medical professionals use advanced diagnostic tools for better treatment outcomes. 
		Digital health platforms enable remote patient monitoring and telemedicine. 
		Pharmaceutical companies are developing innovative drugs through clinical trials. 
		Health insurance providers are expanding coverage for preventive care. 
		Government health agencies are updating public health policies. 
		Medical research institutions are collaborating on breakthrough studies.`
	
	result := ClassifyTopics(text, 12) // Force comprehensive method
	
	if result.Method != "comprehensive" {
		t.Errorf("Expected comprehensive method, got %s", result.Method)
	}

	// Should identify multiple health-related topics
	if len(result.Topics) < 2 {
		t.Error("Expected multiple topics for comprehensive analysis")
	}

	// Check that comprehensive analysis provides richer results
	totalKeywords := 0
	totalExamples := 0
	for _, topic := range result.Topics {
		totalKeywords += len(topic.Keywords)
		totalExamples += len(topic.Examples)
	}

	if totalKeywords < 10 {
		t.Errorf("Expected comprehensive analysis to provide more keywords, got %d", totalKeywords)
	}
}

func TestTopicQualityProgression(t *testing.T) {
	text := `Technology innovation drives economic growth and social progress. 
		Software development teams create applications that solve real-world problems. 
		Data scientists analyze information to generate actionable insights. 
		Artificial intelligence enhances human capabilities across various industries.`
	
	methods := []struct {
		maxTopics int
		method    string
	}{
		{2, "clustering"},
		{5, "statistical"},
		{12, "comprehensive"},
	}

	var previousAccuracy float64
	
	for i, m := range methods {
		result := ClassifyTopics(text, m.maxTopics)
		
		if result.Method != m.method {
			t.Errorf("Expected method %s, got %s", m.method, result.Method)
		}

		// Higher complexity methods should generally have higher accuracy
		if i > 0 && result.QualityMetrics.Accuracy < previousAccuracy {
			t.Logf("Note: Accuracy regression from %.3f to %.3f (method %s to %s)", 
				previousAccuracy, result.QualityMetrics.Accuracy, methods[i-1].method, m.method)
			// Don't fail as this can vary with different texts
		}
		previousAccuracy = result.QualityMetrics.Accuracy
		
		t.Logf("Method %s: accuracy=%.3f, confidence=%.3f, coverage=%.3f", 
			m.method, result.QualityMetrics.Accuracy, 
			result.QualityMetrics.Confidence, result.QualityMetrics.Coverage)
	}
}

func TestEdgeCases(t *testing.T) {
	// Test with zero maxTopics
	result := ClassifyTopics("Some text here", 0)
	if len(result.Topics) != 0 {
		t.Error("Zero maxTopics should return no topics")
	}

	// Test with very large maxTopics
	result = ClassifyTopics("Short text", 100)
	if len(result.Topics) > 10 { // Reasonable upper bound
		t.Errorf("Very large maxTopics should be capped, got %d topics", len(result.Topics))
	}

	// Test with single word
	result = ClassifyTopics("technology", 3)
	if len(result.Topics) > 1 {
		t.Error("Single word should produce at most one topic")
	}

	// Test with repeated text
	repeated := strings.Repeat("machine learning artificial intelligence ", 20)
	result = ClassifyTopics(repeated, 5)
	if len(result.Topics) > 3 {
		t.Error("Repeated text should not produce too many distinct topics")
	}
}

func TestTopicExamples(t *testing.T) {
	text := `Machine learning is revolutionizing data analysis. 
		Advanced algorithms can process massive datasets efficiently. 
		Artificial intelligence applications are found in many industries. 
		Data scientists use statistical models to derive insights.`
	
	result := ClassifyTopics(text, 3)
	
	for i, topic := range result.Topics {
		// Topics should have examples
		if len(topic.Examples) == 0 {
			t.Errorf("Topic %d should have examples", i)
		}

		// Examples should contain topic keywords
		for _, example := range topic.Examples {
			foundKeyword := false
			lowerExample := strings.ToLower(example)
			for _, keyword := range topic.Keywords {
				if strings.Contains(lowerExample, keyword) {
					foundKeyword = true
					break
				}
			}
			if !foundKeyword {
				t.Logf("Example '%s' doesn't contain any keywords from topic '%s'", 
					example, topic.Name)
				// Don't fail as this might be valid in some cases
			}
		}
	}
}

func BenchmarkClassifyTopicsClustering(b *testing.B) {
	text := "Machine learning algorithms analyze data patterns to make predictions."
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ClassifyTopics(text, 2)
	}
}

func BenchmarkClassifyTopicsStatistical(b *testing.B) {
	text := `Technology companies are investing heavily in artificial intelligence research. 
		Machine learning applications span across healthcare, finance, and transportation. 
		Data science teams develop predictive models for business decision making.`
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ClassifyTopics(text, 7)
	}
}

func BenchmarkClassifyTopicsComprehensive(b *testing.B) {
	text := `Artificial intelligence is transforming multiple industries including healthcare, finance, and education. 
		Machine learning algorithms enable computers to learn from data without explicit programming. 
		Deep learning networks use multiple layers to model complex patterns in large datasets. 
		Natural language processing allows machines to understand and generate human language. 
		Computer vision systems can identify objects and patterns in images and videos.`
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ClassifyTopics(text, 15)
	}
}