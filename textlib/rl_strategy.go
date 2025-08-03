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
	"fmt"
	"math"
	"strings"
)

// StrategySelector selects optimal processing strategies
type StrategySelector struct {
	// Strategy history for learning
	history []StrategyOutcome
	
	// Learned preferences
	preferences map[string]float64
}

// StrategyOutcome represents the outcome of a strategy selection
type StrategyOutcome struct {
	Context    TextCharacteristics
	Strategy   ProcessingStrategy
	Metrics    OptimizationMetrics
	Successful bool
}

// TextCharacteristics describes text properties for strategy selection
type TextCharacteristics struct {
	Length      int     `json:"length"`
	Language    string  `json:"language"`
	Domain      string  `json:"domain"`
	Complexity  float64 `json:"complexity"`
	Structure   string  `json:"structure"`
}

// NewStrategySelector creates a new strategy selector
func NewStrategySelector() *StrategySelector {
	return &StrategySelector{
		history:     make([]StrategyOutcome, 0),
		preferences: make(map[string]float64),
	}
}

// SelectStrategy selects the optimal strategy for given text characteristics
func (ss *StrategySelector) SelectStrategy(characteristics TextCharacteristics, requirements AlgorithmRequirements) (ProcessingStrategy, error) {
	// Analyze text characteristics
	strategyType := ss.determineStrategyType(characteristics)
	
	// Build strategy based on requirements
	strategy := ProcessingStrategy{
		Name:        strategyType,
		Description: fmt.Sprintf("Optimized strategy for %s text", characteristics.Domain),
		Parameters:  make(map[string]interface{}),
	}
	
	// Set parameters based on text characteristics
	switch strategyType {
	case "fast":
		strategy.Parameters["depth"] = 1
		strategy.Parameters["algorithms"] = []string{"flesch", "gunning-fog"}
		strategy.Parameters["quality"] = 0.7
		strategy.ExpectedQuality = 0.7
		strategy.ExpectedSpeed = 0.95
		strategy.ResourceRequirements = ResourceRequirements{
			MinMemoryMB:      10,
			MaxMemoryMB:      50,
			EstimatedCPUTime: 50,
			NetworkRequired:  false,
			CacheRecommended: true,
		}
		
	case "balanced":
		strategy.Parameters["depth"] = 2
		strategy.Parameters["algorithms"] = []string{"flesch", "gunning-fog", "coleman-liau", "ari"}
		strategy.Parameters["quality"] = 0.85
		strategy.ExpectedQuality = 0.85
		strategy.ExpectedSpeed = 0.7
		strategy.ResourceRequirements = ResourceRequirements{
			MinMemoryMB:      50,
			MaxMemoryMB:      200,
			EstimatedCPUTime: 200,
			NetworkRequired:  false,
			CacheRecommended: true,
		}
		
	case "comprehensive":
		strategy.Parameters["depth"] = 3
		strategy.Parameters["algorithms"] = []string{"all"}
		strategy.Parameters["quality"] = 0.95
		strategy.ExpectedQuality = 0.95
		strategy.ExpectedSpeed = 0.4
		strategy.ResourceRequirements = ResourceRequirements{
			MinMemoryMB:      100,
			MaxMemoryMB:      500,
			EstimatedCPUTime: 500,
			NetworkRequired:  false,
			CacheRecommended: false, // Too many variations for effective caching
		}
		
	default:
		return strategy, fmt.Errorf("unknown strategy type: %s", strategyType)
	}
	
	// Adjust based on specific requirements
	ss.adjustStrategyForRequirements(&strategy, requirements)
	
	// Apply learned preferences
	ss.applyLearnedPreferences(&strategy)
	
	return strategy, nil
}

// determineStrategyType determines the base strategy type
func (ss *StrategySelector) determineStrategyType(characteristics TextCharacteristics) string {
	// Simple heuristics for strategy selection
	if characteristics.Length < 100 {
		return "fast"
	}
	
	if characteristics.Length > 10000 || characteristics.Complexity > 0.8 {
		return "comprehensive"
	}
	
	if characteristics.Domain == "technical" || characteristics.Domain == "academic" {
		return "comprehensive"
	}
	
	if characteristics.Domain == "social-media" || characteristics.Domain == "chat" {
		return "fast"
	}
	
	return "balanced"
}

// adjustStrategyForRequirements adjusts strategy based on requirements
func (ss *StrategySelector) adjustStrategyForRequirements(strategy *ProcessingStrategy, requirements AlgorithmRequirements) {
	// Adjust quality parameter
	if requirements.MinQuality > strategy.ExpectedQuality {
		// Need to increase quality
		if depth, ok := strategy.Parameters["depth"].(int); ok && depth < 3 {
			strategy.Parameters["depth"] = depth + 1
			strategy.ExpectedQuality = math.Min(strategy.ExpectedQuality+0.1, 0.95)
		}
		
		if quality, ok := strategy.Parameters["quality"].(float64); ok {
			strategy.Parameters["quality"] = math.Max(quality, requirements.MinQuality)
		}
	}
	
	// Adjust for time constraints
	if requirements.MaxTimeMs > 0 && strategy.ResourceRequirements.EstimatedCPUTime > requirements.MaxTimeMs {
		// Need to reduce processing time
		if depth, ok := strategy.Parameters["depth"].(int); ok && depth > 1 {
			strategy.Parameters["depth"] = depth - 1
			strategy.ExpectedSpeed = math.Min(strategy.ExpectedSpeed+0.2, 0.95)
			strategy.ResourceRequirements.EstimatedCPUTime = 
				int64(float64(strategy.ResourceRequirements.EstimatedCPUTime) * 0.5)
		}
		
		// Reduce algorithm count
		if algorithms, ok := strategy.Parameters["algorithms"].([]string); ok && len(algorithms) > 2 {
			strategy.Parameters["algorithms"] = algorithms[:2]
		}
	}
	
	// Adjust for memory constraints
	if requirements.MaxMemoryMB > 0 && int64(strategy.ResourceRequirements.MaxMemoryMB) > requirements.MaxMemoryMB {
		strategy.ResourceRequirements.MaxMemoryMB = int(requirements.MaxMemoryMB)
		strategy.Parameters["streaming"] = true
		strategy.Parameters["batch_size"] = 100
	}
}

// applyLearnedPreferences applies learned preferences from history
func (ss *StrategySelector) applyLearnedPreferences(strategy *ProcessingStrategy) {
	// Apply preferences based on strategy name
	if preference, exists := ss.preferences[strategy.Name]; exists {
		// Adjust expected metrics based on historical performance
		adjustment := preference / 100.0 // Convert to adjustment factor
		strategy.ExpectedQuality *= (1 + adjustment)
		strategy.ExpectedSpeed *= (1 - adjustment*0.5) // Quality improvements often reduce speed
	}
	
	// Apply domain-specific preferences
	if _, ok := strategy.Parameters["algorithms"].([]string); ok {
		domainKey := fmt.Sprintf("algorithms_%s", strategy.Name)
		if _, exists := ss.preferences[domainKey]; exists {
			// Could adjust algorithm selection based on domain preferences
		}
	}
}

// RecordOutcome records the outcome of a strategy selection
func (ss *StrategySelector) RecordOutcome(characteristics TextCharacteristics, strategy ProcessingStrategy, metrics OptimizationMetrics, successful bool) {
	outcome := StrategyOutcome{
		Context:    characteristics,
		Strategy:   strategy,
		Metrics:    metrics,
		Successful: successful,
	}
	
	ss.history = append(ss.history, outcome)
	
	// Update preferences based on outcome
	ss.updatePreferences(outcome)
	
	// Limit history size
	if len(ss.history) > 1000 {
		ss.history = ss.history[len(ss.history)-1000:]
	}
}

// updatePreferences updates learned preferences based on outcomes
func (ss *StrategySelector) updatePreferences(outcome StrategyOutcome) {
	// Simple preference update based on weighted score
	score := outcome.Metrics.WeightedTotal
	
	// Update strategy preference
	currentPref := ss.preferences[outcome.Strategy.Name]
	if outcome.Successful {
		// Positive reinforcement
		ss.preferences[outcome.Strategy.Name] = currentPref + (score-0.7)*10
	} else {
		// Negative reinforcement
		ss.preferences[outcome.Strategy.Name] = currentPref - 5
	}
	
	// Update domain-specific preferences
	domainKey := fmt.Sprintf("domain_%s_%s", outcome.Context.Domain, outcome.Strategy.Name)
	domainPref := ss.preferences[domainKey]
	if outcome.Successful {
		ss.preferences[domainKey] = domainPref + (score-0.7)*5
	} else {
		ss.preferences[domainKey] = domainPref - 2
	}
}

// GetRecommendations provides strategy recommendations
func (ss *StrategySelector) GetRecommendations(characteristics TextCharacteristics) []StrategyRecommendation {
	recommendations := []StrategyRecommendation{}
	
	// Analyze history for similar contexts
	similarOutcomes := ss.findSimilarOutcomes(characteristics)
	
	// Group by strategy and calculate average performance
	strategyPerformance := make(map[string][]float64)
	for _, outcome := range similarOutcomes {
		if outcome.Successful {
			strategyPerformance[outcome.Strategy.Name] = append(
				strategyPerformance[outcome.Strategy.Name],
				outcome.Metrics.WeightedTotal,
			)
		}
	}
	
	// Create recommendations
	for strategyName, scores := range strategyPerformance {
		if len(scores) > 0 {
			avgScore := 0.0
			for _, score := range scores {
				avgScore += score
			}
			avgScore /= float64(len(scores))
			
			recommendation := StrategyRecommendation{
				StrategyName: strategyName,
				Confidence:   math.Min(float64(len(scores))/10.0, 1.0), // More data = higher confidence
				ExpectedScore: avgScore,
				Rationale:    fmt.Sprintf("Based on %d similar successful outcomes", len(scores)),
			}
			
			recommendations = append(recommendations, recommendation)
		}
	}
	
	return recommendations
}

// StrategyRecommendation represents a strategy recommendation
type StrategyRecommendation struct {
	StrategyName  string  `json:"strategy_name"`
	Confidence    float64 `json:"confidence"`
	ExpectedScore float64 `json:"expected_score"`
	Rationale     string  `json:"rationale"`
}

// findSimilarOutcomes finds historically similar text processing outcomes
func (ss *StrategySelector) findSimilarOutcomes(characteristics TextCharacteristics) []StrategyOutcome {
	similar := []StrategyOutcome{}
	
	for _, outcome := range ss.history {
		similarity := ss.calculateSimilarity(characteristics, outcome.Context)
		if similarity > 0.7 { // 70% similarity threshold
			similar = append(similar, outcome)
		}
	}
	
	return similar
}

// calculateSimilarity calculates similarity between text characteristics
func (ss *StrategySelector) calculateSimilarity(a, b TextCharacteristics) float64 {
	similarity := 0.0
	weights := 0.0
	
	// Length similarity (logarithmic scale)
	if a.Length > 0 && b.Length > 0 {
		lengthRatio := math.Min(float64(a.Length), float64(b.Length)) / 
			math.Max(float64(a.Length), float64(b.Length))
		similarity += lengthRatio * 0.3
		weights += 0.3
	}
	
	// Language match
	if a.Language == b.Language {
		similarity += 0.2
	}
	weights += 0.2
	
	// Domain match
	if a.Domain == b.Domain {
		similarity += 0.3
	} else if strings.Contains(a.Domain, b.Domain) || strings.Contains(b.Domain, a.Domain) {
		similarity += 0.15
	}
	weights += 0.3
	
	// Complexity similarity
	complexityDiff := math.Abs(a.Complexity - b.Complexity)
	similarity += (1 - complexityDiff) * 0.2
	weights += 0.2
	
	if weights > 0 {
		return similarity / weights
	}
	
	return 0
}

// GetInsights provides insights from strategy selection history
func (ss *StrategySelector) GetInsights() StrategyInsights {
	insights := StrategyInsights{
		TotalSelections: len(ss.history),
		StrategyUsage:   make(map[string]int),
		DomainPatterns:  make(map[string]string),
		PerformanceTrends: make(map[string][]float64),
	}
	
	// Analyze strategy usage
	for _, outcome := range ss.history {
		insights.StrategyUsage[outcome.Strategy.Name]++
		
		// Track performance trends
		if outcome.Successful {
			insights.PerformanceTrends[outcome.Strategy.Name] = append(
				insights.PerformanceTrends[outcome.Strategy.Name],
				outcome.Metrics.WeightedTotal,
			)
		}
	}
	
	// Identify domain patterns
	domainStrategy := make(map[string]map[string]int)
	for _, outcome := range ss.history {
		if outcome.Successful {
			if _, exists := domainStrategy[outcome.Context.Domain]; !exists {
				domainStrategy[outcome.Context.Domain] = make(map[string]int)
			}
			domainStrategy[outcome.Context.Domain][outcome.Strategy.Name]++
		}
	}
	
	// Find most common strategy per domain
	for domain, strategies := range domainStrategy {
		maxCount := 0
		bestStrategy := ""
		for strategy, count := range strategies {
			if count > maxCount {
				maxCount = count
				bestStrategy = strategy
			}
		}
		insights.DomainPatterns[domain] = bestStrategy
	}
	
	// Calculate success rate
	successCount := 0
	for _, outcome := range ss.history {
		if outcome.Successful {
			successCount++
		}
	}
	if insights.TotalSelections > 0 {
		insights.SuccessRate = float64(successCount) / float64(insights.TotalSelections)
	}
	
	return insights
}

// StrategyInsights represents insights from strategy selection
type StrategyInsights struct {
	TotalSelections   int                      `json:"total_selections"`
	SuccessRate       float64                  `json:"success_rate"`
	StrategyUsage     map[string]int           `json:"strategy_usage"`
	DomainPatterns    map[string]string        `json:"domain_patterns"`
	PerformanceTrends map[string][]float64     `json:"performance_trends"`
}