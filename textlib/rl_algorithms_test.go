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

func TestAlgorithmRegistry(t *testing.T) {
	// Clear registry for test
	algorithmRegistry = &AlgorithmRegistry{
		algorithms: make(map[string]Algorithm),
		metadata:   make(map[string]AlgorithmMetadata),
	}
	
	// Create test algorithm
	testAlgo := NewBaseAlgorithm(
		"test-algo",
		"Test algorithm for unit tests",
		func(input interface{}, params map[string]interface{}) (interface{}, error) {
			text := input.(string)
			return strings.ToUpper(text), nil
		},
		ComplexityEstimate{
			TimeComplexity:    "O(n)",
			SpaceComplexity:   "O(1)",
			EstimatedTimeMs:   10,
			EstimatedMemoryMB: 1,
		},
	)
	
	// Register algorithm
	metadata := AlgorithmMetadata{
		Category:     "text-transform",
		Version:      "1.0.0",
		Author:       "Test",
		Tags:         []string{"fast", "simple"},
		QualityScore: 1.0,
		SpeedScore:   0.95,
		MemoryScore:  0.98,
		SupportedInputs: []string{"text"},
	}
	
	err := RegisterAlgorithm("test-algo", testAlgo, metadata)
	if err != nil {
		t.Fatalf("Failed to register algorithm: %v", err)
	}
	
	// Try to register again (should fail)
	err = RegisterAlgorithm("test-algo", testAlgo, metadata)
	if err == nil {
		t.Error("Expected error when registering duplicate algorithm")
	}
	
	// Retrieve algorithm
	algo, err := GetAlgorithm("test-algo")
	if err != nil {
		t.Fatalf("Failed to get algorithm: %v", err)
	}
	
	if algo.Name() != "test-algo" {
		t.Errorf("Algorithm name mismatch: got %s, want test-algo", algo.Name())
	}
	
	// Test processing
	result, err := algo.Process("hello", nil)
	if err != nil {
		t.Fatalf("Algorithm processing failed: %v", err)
	}
	
	if result != "HELLO" {
		t.Errorf("Unexpected result: got %v, want HELLO", result)
	}
	
	// Get metadata
	meta, err := GetAlgorithmMetadata("test-algo")
	if err != nil {
		t.Fatalf("Failed to get metadata: %v", err)
	}
	
	if meta.Category != "text-transform" {
		t.Errorf("Category mismatch: got %s, want text-transform", meta.Category)
	}
	
	// List algorithms
	algos := ListAlgorithms()
	if len(algos) != 1 || algos[0] != "test-algo" {
		t.Errorf("ListAlgorithms incorrect: %v", algos)
	}
	
	// List by category
	categoryAlgos := ListAlgorithmsByCategory("text-transform")
	if len(categoryAlgos) != 1 || categoryAlgos[0] != "test-algo" {
		t.Errorf("ListAlgorithmsByCategory incorrect: %v", categoryAlgos)
	}
}

func TestComplexityEstimation(t *testing.T) {
	// Test linear complexity
	linearAlgo := NewBaseAlgorithm(
		"linear",
		"Linear algorithm",
		nil,
		ComplexityEstimate{
			TimeComplexity:    "O(n)",
			SpaceComplexity:   "O(1)",
			EstimatedTimeMs:   10, // for 1000 units
			EstimatedMemoryMB: 1,
		},
	)
	
	estimate := linearAlgo.EstimateComplexity(2000)
	if estimate.EstimatedTimeMs != 20 {
		t.Errorf("Linear complexity estimation wrong: got %d, want 20", estimate.EstimatedTimeMs)
	}
	
	// Test n log n complexity
	nlogAlgo := NewBaseAlgorithm(
		"nlogn",
		"N log N algorithm",
		nil,
		ComplexityEstimate{
			TimeComplexity:    "O(n log n)",
			SpaceComplexity:   "O(n)",
			EstimatedTimeMs:   10, // for 1000 units
			EstimatedMemoryMB: 5,
		},
	)
	
	estimate = nlogAlgo.EstimateComplexity(2000)
	// Should be roughly 2 * 1.3 * 10 = 26
	if estimate.EstimatedTimeMs < 20 || estimate.EstimatedTimeMs > 30 {
		t.Errorf("N log N complexity estimation out of range: %d", estimate.EstimatedTimeMs)
	}
	
	// Test quadratic complexity
	quadAlgo := NewBaseAlgorithm(
		"quadratic",
		"Quadratic algorithm",
		nil,
		ComplexityEstimate{
			TimeComplexity:    "O(n²)",
			SpaceComplexity:   "O(n²)",
			EstimatedTimeMs:   10, // for 1000 units
			EstimatedMemoryMB: 10,
		},
	)
	
	estimate = quadAlgo.EstimateComplexity(2000)
	// Should be (2000*2000)/(1000*1000) * 10 = 40
	if estimate.EstimatedTimeMs != 40 {
		t.Errorf("Quadratic complexity estimation wrong: got %d, want 40", estimate.EstimatedTimeMs)
	}
}

func TestAlgorithmSelection(t *testing.T) {
	// Clear and initialize registry
	algorithmRegistry = &AlgorithmRegistry{
		algorithms: make(map[string]Algorithm),
		metadata:   make(map[string]AlgorithmMetadata),
	}
	
	// Register multiple algorithms
	algorithms := []struct {
		name     string
		category string
		quality  float64
		speed    float64
		memory   float64
		tags     []string
	}{
		{"fast-low-quality", "analysis", 0.6, 0.95, 0.98, []string{"fast", "basic"}},
		{"balanced", "analysis", 0.8, 0.7, 0.8, []string{"balanced"}},
		{"high-quality", "analysis", 0.95, 0.4, 0.5, []string{"accurate", "comprehensive"}},
		{"memory-efficient", "analysis", 0.75, 0.6, 0.95, []string{"memory-efficient"}},
	}
	
	for _, alg := range algorithms {
		algo := NewBaseAlgorithm(alg.name, "Test algorithm", nil, ComplexityEstimate{})
		metadata := AlgorithmMetadata{
			Category:     alg.category,
			QualityScore: alg.quality,
			SpeedScore:   alg.speed,
			MemoryScore:  alg.memory,
			Tags:         alg.tags,
		}
		RegisterAlgorithm(alg.name, algo, metadata)
	}
	
	// Test 1: Prioritize quality
	requirements := AlgorithmRequirements{
		MinQuality:    0.7,
		QualityWeight: 0.8,
		SpeedWeight:   0.1,
		MemoryWeight:  0.1,
	}
	
	best, err := SelectBestAlgorithm("analysis", requirements)
	if err != nil {
		t.Fatalf("Failed to select algorithm: %v", err)
	}
	
	if best != "high-quality" {
		t.Errorf("Expected high-quality algorithm, got %s", best)
	}
	
	// Test 2: Prioritize speed
	requirements = AlgorithmRequirements{
		MinQuality:    0.5,
		QualityWeight: 0.1,
		SpeedWeight:   0.8,
		MemoryWeight:  0.1,
	}
	
	best, err = SelectBestAlgorithm("analysis", requirements)
	if err != nil {
		t.Fatalf("Failed to select algorithm: %v", err)
	}
	
	if best != "fast-low-quality" {
		t.Errorf("Expected fast-low-quality algorithm, got %s", best)
	}
	
	// Test 3: Prioritize memory efficiency
	requirements = AlgorithmRequirements{
		MinQuality:    0.7,
		QualityWeight: 0.2,
		SpeedWeight:   0.2,
		MemoryWeight:  0.6,
	}
	
	best, err = SelectBestAlgorithm("analysis", requirements)
	if err != nil {
		t.Fatalf("Failed to select algorithm: %v", err)
	}
	
	if best != "memory-efficient" {
		t.Errorf("Expected memory-efficient algorithm, got %s", best)
	}
	
	// Test 4: With preferred tags
	requirements = AlgorithmRequirements{
		MinQuality:     0.7,
		QualityWeight:  0.3,
		SpeedWeight:    0.3,
		MemoryWeight:   0.4,
		PreferredTags:  []string{"balanced"},
	}
	
	best, err = SelectBestAlgorithm("analysis", requirements)
	if err != nil {
		t.Fatalf("Failed to select algorithm: %v", err)
	}
	
	if best != "balanced" {
		t.Errorf("Expected balanced algorithm (with tag preference), got %s", best)
	}
}

func TestDefaultAlgorithms(t *testing.T) {
	// Clear registry
	algorithmRegistry = &AlgorithmRegistry{
		algorithms: make(map[string]Algorithm),
		metadata:   make(map[string]AlgorithmMetadata),
	}
	
	// Initialize default algorithms
	InitializeDefaultAlgorithms()
	
	// Check sentiment algorithms
	sentimentAlgos := ListAlgorithmsByCategory("sentiment")
	if len(sentimentAlgos) < 2 {
		t.Errorf("Expected at least 2 sentiment algorithms, got %d", len(sentimentAlgos))
	}
	
	// Check complexity algorithms
	complexityAlgos := ListAlgorithmsByCategory("complexity")
	if len(complexityAlgos) < 2 {
		t.Errorf("Expected at least 2 complexity algorithms, got %d", len(complexityAlgos))
	}
	
	// Verify sentiment-lexicon
	lexicon, err := GetAlgorithm("sentiment-lexicon")
	if err != nil {
		t.Fatalf("Failed to get sentiment-lexicon: %v", err)
	}
	
	if lexicon.Name() != "sentiment-lexicon" {
		t.Errorf("Algorithm name mismatch: %s", lexicon.Name())
	}
	
	// Verify metadata quality scores
	meta, err := GetAlgorithmMetadata("sentiment-lexicon")
	if err != nil {
		t.Fatalf("Failed to get metadata: %v", err)
	}
	
	if meta.SpeedScore <= meta.QualityScore {
		t.Error("Lexicon algorithm should be faster than accurate")
	}
	
	// Verify ML algorithm is more accurate
	mlMeta, err := GetAlgorithmMetadata("sentiment-ml-basic")
	if err != nil {
		t.Fatalf("Failed to get ML metadata: %v", err)
	}
	
	if mlMeta.QualityScore <= meta.QualityScore {
		t.Error("ML algorithm should be more accurate than lexicon")
	}
}