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
	"testing"
)

func TestStrategySelector(t *testing.T) {
	selector := NewStrategySelector()
	
	// Test short text strategy
	shortText := TextCharacteristics{
		Length:     50,
		Language:   "en",
		Domain:     "general",
		Complexity: 0.3,
		Structure:  "simple",
	}
	
	requirements := AlgorithmRequirements{
		MinQuality: 0.7,
	}
	
	strategy, err := selector.SelectStrategy(shortText, requirements)
	if err != nil {
		t.Fatalf("Failed to select strategy: %v", err)
	}
	
	if strategy.Name != "fast" {
		t.Errorf("Expected fast strategy for short text, got %s", strategy.Name)
	}
	
	// Test long complex text strategy
	longText := TextCharacteristics{
		Length:     15000,
		Language:   "en",
		Domain:     "technical",
		Complexity: 0.9,
		Structure:  "complex",
	}
	
	strategy, err = selector.SelectStrategy(longText, requirements)
	if err != nil {
		t.Fatalf("Failed to select strategy: %v", err)
	}
	
	if strategy.Name != "comprehensive" {
		t.Errorf("Expected comprehensive strategy for long text, got %s", strategy.Name)
	}
	
	// Test balanced strategy
	mediumText := TextCharacteristics{
		Length:     2000,
		Language:   "en",
		Domain:     "business",
		Complexity: 0.5,
		Structure:  "standard",
	}
	
	strategy, err = selector.SelectStrategy(mediumText, requirements)
	if err != nil {
		t.Fatalf("Failed to select strategy: %v", err)
	}
	
	if strategy.Name != "balanced" {
		t.Errorf("Expected balanced strategy for medium text, got %s", strategy.Name)
	}
}

func TestStrategyAdjustments(t *testing.T) {
	selector := NewStrategySelector()
	
	text := TextCharacteristics{
		Length:     5000,
		Language:   "en",
		Domain:     "general",
		Complexity: 0.6,
		Structure:  "standard",
	}
	
	// Test quality adjustment
	highQualityReq := AlgorithmRequirements{
		MinQuality:    0.9,
		QualityWeight: 0.8,
		SpeedWeight:   0.1,
		MemoryWeight:  0.1,
	}
	
	strategy, err := selector.SelectStrategy(text, highQualityReq)
	if err != nil {
		t.Fatalf("Failed to select strategy: %v", err)
	}
	
	if quality, ok := strategy.Parameters["quality"].(float64); !ok || quality < 0.9 {
		t.Errorf("Strategy quality not adjusted properly: %v", strategy.Parameters["quality"])
	}
	
	// Test time constraint adjustment
	fastReq := AlgorithmRequirements{
		MinQuality:    0.6,
		MaxTimeMs:     100,
		QualityWeight: 0.1,
		SpeedWeight:   0.8,
		MemoryWeight:  0.1,
	}
	
	strategy, err = selector.SelectStrategy(text, fastReq)
	if err != nil {
		t.Fatalf("Failed to select strategy: %v", err)
	}
	
	if strategy.ResourceRequirements.EstimatedCPUTime > 100 {
		t.Errorf("Strategy time not adjusted properly: %d ms", 
			strategy.ResourceRequirements.EstimatedCPUTime)
	}
	
	// Test memory constraint adjustment
	memoryReq := AlgorithmRequirements{
		MinQuality:    0.7,
		MaxMemoryMB:   100,
		QualityWeight: 0.3,
		SpeedWeight:   0.3,
		MemoryWeight:  0.4,
	}
	
	strategy, err = selector.SelectStrategy(text, memoryReq)
	if err != nil {
		t.Fatalf("Failed to select strategy: %v", err)
	}
	
	if strategy.ResourceRequirements.MaxMemoryMB > 100 {
		t.Errorf("Strategy memory not adjusted properly: %d MB",
			strategy.ResourceRequirements.MaxMemoryMB)
	}
	
	if _, ok := strategy.Parameters["streaming"]; !ok {
		t.Error("Expected streaming parameter for memory-constrained strategy")
	}
}

func TestStrategyOutcomeRecording(t *testing.T) {
	selector := NewStrategySelector()
	
	// Record some outcomes
	characteristics := TextCharacteristics{
		Length:     1000,
		Language:   "en",
		Domain:     "news",
		Complexity: 0.5,
		Structure:  "article",
	}
	
	strategy := ProcessingStrategy{
		Name:            "balanced",
		ExpectedQuality: 0.85,
		ExpectedSpeed:   0.70,
	}
	
	metrics := OptimizationMetrics{
		QualityScore:     0.88,
		PerformanceScore: 0.75,
		ResourceScore:    0.82,
		UserSatisfaction: 0.90,
		WeightedTotal:    0.84,
	}
	
	// Record successful outcome
	selector.RecordOutcome(characteristics, strategy, metrics, true)
	
	// Check history
	if len(selector.history) != 1 {
		t.Errorf("Expected 1 outcome in history, got %d", len(selector.history))
	}
	
	// Check preferences updated
	if pref, exists := selector.preferences["balanced"]; !exists || pref <= 0 {
		t.Error("Preferences not updated for successful outcome")
	}
	
	// Record failed outcome
	poorMetrics := OptimizationMetrics{
		QualityScore:     0.45,
		PerformanceScore: 0.30,
		ResourceScore:    0.40,
		UserSatisfaction: 0.35,
		WeightedTotal:    0.38,
	}
	
	selector.RecordOutcome(characteristics, strategy, poorMetrics, false)
	
	// Check preferences decreased
	if pref, exists := selector.preferences["balanced"]; !exists || pref >= 0 {
		t.Error("Preferences not decreased for failed outcome")
	}
}

func TestStrategyRecommendations(t *testing.T) {
	selector := NewStrategySelector()
	
	// Record multiple outcomes for different strategies
	characteristics := TextCharacteristics{
		Length:     2000,
		Language:   "en",
		Domain:     "technical",
		Complexity: 0.7,
		Structure:  "documentation",
	}
	
	// Record successful comprehensive strategy
	for i := 0; i < 5; i++ {
		strategy := ProcessingStrategy{Name: "comprehensive"}
		metrics := OptimizationMetrics{
			WeightedTotal: 0.85 + float64(i)*0.01,
		}
		selector.RecordOutcome(characteristics, strategy, metrics, true)
	}
	
	// Record less successful fast strategy
	for i := 0; i < 3; i++ {
		strategy := ProcessingStrategy{Name: "fast"}
		metrics := OptimizationMetrics{
			WeightedTotal: 0.65 + float64(i)*0.01,
		}
		selector.RecordOutcome(characteristics, strategy, metrics, true)
	}
	
	// Get recommendations
	recommendations := selector.GetRecommendations(characteristics)
	
	if len(recommendations) < 2 {
		t.Errorf("Expected at least 2 recommendations, got %d", len(recommendations))
	}
	
	// Check comprehensive strategy is recommended higher
	var comprehensiveRec, fastRec *StrategyRecommendation
	for i := range recommendations {
		if recommendations[i].StrategyName == "comprehensive" {
			comprehensiveRec = &recommendations[i]
		} else if recommendations[i].StrategyName == "fast" {
			fastRec = &recommendations[i]
		}
	}
	
	if comprehensiveRec == nil || fastRec == nil {
		t.Fatal("Expected both comprehensive and fast recommendations")
	}
	
	if comprehensiveRec.ExpectedScore <= fastRec.ExpectedScore {
		t.Error("Comprehensive strategy should have higher expected score")
	}
	
	if comprehensiveRec.Confidence <= fastRec.Confidence {
		t.Error("Comprehensive strategy should have higher confidence (more data)")
	}
}

func TestSimilarityCalculation(t *testing.T) {
	selector := NewStrategySelector()
	
	// Test identical characteristics
	char1 := TextCharacteristics{
		Length:     1000,
		Language:   "en",
		Domain:     "technical",
		Complexity: 0.7,
		Structure:  "article",
	}
	
	char2 := TextCharacteristics{
		Length:     1000,
		Language:   "en",
		Domain:     "technical",
		Complexity: 0.7,
		Structure:  "article",
	}
	
	similarity := selector.calculateSimilarity(char1, char2)
	if similarity != 1.0 {
		t.Errorf("Identical characteristics should have similarity 1.0, got %f", similarity)
	}
	
	// Test different language
	char3 := TextCharacteristics{
		Length:     1000,
		Language:   "es",
		Domain:     "technical",
		Complexity: 0.7,
		Structure:  "article",
	}
	
	similarity = selector.calculateSimilarity(char1, char3)
	if similarity >= 1.0 {
		t.Errorf("Different language should reduce similarity, got %f", similarity)
	}
	
	// Test very different length
	char4 := TextCharacteristics{
		Length:     100,
		Language:   "en",
		Domain:     "technical",
		Complexity: 0.7,
		Structure:  "article",
	}
	
	similarity = selector.calculateSimilarity(char1, char4)
	if similarity >= 0.8 {
		t.Errorf("Very different length should reduce similarity significantly, got %f", similarity)
	}
	
	// Test subdomain matching
	char5 := TextCharacteristics{
		Length:     1000,
		Language:   "en",
		Domain:     "technical-documentation",
		Complexity: 0.7,
		Structure:  "article",
	}
	
	similarity = selector.calculateSimilarity(char1, char5)
	if similarity < 0.7 || similarity >= 0.95 {
		t.Errorf("Subdomain should have partial similarity, got %f", similarity)
	}
}

func TestStrategyInsights(t *testing.T) {
	selector := NewStrategySelector()
	
	// Record various outcomes
	domains := []string{"technical", "news", "social-media", "technical", "news"}
	strategies := []string{"comprehensive", "balanced", "fast", "comprehensive", "balanced"}
	
	for i, domain := range domains {
		char := TextCharacteristics{
			Length:     1000 * (i + 1),
			Language:   "en",
			Domain:     domain,
			Complexity: 0.5,
			Structure:  "standard",
		}
		
		strategy := ProcessingStrategy{
			Name: strategies[i],
		}
		
		metrics := OptimizationMetrics{
			WeightedTotal: 0.7 + float64(i)*0.05,
		}
		
		selector.RecordOutcome(char, strategy, metrics, true)
	}
	
	// Get insights
	insights := selector.GetInsights()
	
	if insights.TotalSelections != 5 {
		t.Errorf("Expected 5 total selections, got %d", insights.TotalSelections)
	}
	
	if insights.SuccessRate != 1.0 {
		t.Errorf("Expected 100%% success rate, got %f", insights.SuccessRate)
	}
	
	// Check strategy usage
	if insights.StrategyUsage["comprehensive"] != 2 {
		t.Errorf("Expected 2 uses of comprehensive strategy, got %d", 
			insights.StrategyUsage["comprehensive"])
	}
	
	// Check domain patterns
	if insights.DomainPatterns["technical"] != "comprehensive" {
		t.Errorf("Expected comprehensive strategy for technical domain, got %s",
			insights.DomainPatterns["technical"])
	}
	
	// Check performance trends
	if len(insights.PerformanceTrends["balanced"]) != 2 {
		t.Errorf("Expected 2 performance records for balanced strategy, got %d",
			len(insights.PerformanceTrends["balanced"]))
	}
}

func TestHistoryLimit(t *testing.T) {
	selector := NewStrategySelector()
	
	// Record more than limit (1000) outcomes
	for i := 0; i < 1100; i++ {
		char := TextCharacteristics{
			Length:     1000,
			Language:   "en",
			Domain:     "test",
			Complexity: 0.5,
			Structure:  "standard",
		}
		
		strategy := ProcessingStrategy{
			Name: "test",
		}
		
		metrics := OptimizationMetrics{
			WeightedTotal: 0.8,
		}
		
		selector.RecordOutcome(char, strategy, metrics, true)
	}
	
	// Check history is limited
	if len(selector.history) != 1000 {
		t.Errorf("Expected history limited to 1000, got %d", len(selector.history))
	}
}