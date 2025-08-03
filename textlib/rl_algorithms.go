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
	"sync"
)

// Algorithm represents a text processing algorithm
type Algorithm interface {
	// Name returns the algorithm name
	Name() string
	
	// Description returns a human-readable description
	Description() string
	
	// Process executes the algorithm on the input
	Process(input interface{}, params map[string]interface{}) (interface{}, error)
	
	// EstimateComplexity returns computational complexity estimate
	EstimateComplexity(inputSize int) ComplexityEstimate
	
	// ValidateParams validates algorithm parameters
	ValidateParams(params map[string]interface{}) error
}

// ComplexityEstimate represents computational complexity
type ComplexityEstimate struct {
	TimeComplexity   string  // e.g., "O(n)", "O(n log n)", "O(n²)"
	SpaceComplexity  string  // e.g., "O(1)", "O(n)"
	EstimatedTimeMs  int64   // Estimated time in milliseconds
	EstimatedMemoryMB int64  // Estimated memory in megabytes
}

// AlgorithmRegistry manages available algorithms
type AlgorithmRegistry struct {
	mu         sync.RWMutex
	algorithms map[string]Algorithm
	metadata   map[string]AlgorithmMetadata
}

// AlgorithmMetadata stores additional algorithm information
type AlgorithmMetadata struct {
	Category        string   // e.g., "sentiment", "complexity", "entity"
	Version         string   // e.g., "1.0.0"
	Author          string   // e.g., "Caia Tech"
	Tags            []string // e.g., ["fast", "accurate", "memory-efficient"]
	QualityScore    float64  // 0-1, based on benchmarks
	SpeedScore      float64  // 0-1, based on benchmarks
	MemoryScore     float64  // 0-1, based on benchmarks
	SupportedInputs []string // e.g., ["text", "html", "markdown"]
}

// Global algorithm registry
var (
	algorithmRegistry = &AlgorithmRegistry{
		algorithms: make(map[string]Algorithm),
		metadata:   make(map[string]AlgorithmMetadata),
	}
)

// RegisterAlgorithm registers a new algorithm
func RegisterAlgorithm(name string, algo Algorithm, metadata AlgorithmMetadata) error {
	algorithmRegistry.mu.Lock()
	defer algorithmRegistry.mu.Unlock()
	
	if _, exists := algorithmRegistry.algorithms[name]; exists {
		return fmt.Errorf("algorithm %s already registered", name)
	}
	
	algorithmRegistry.algorithms[name] = algo
	algorithmRegistry.metadata[name] = metadata
	
	return nil
}

// GetAlgorithm retrieves an algorithm by name
func GetAlgorithm(name string) (Algorithm, error) {
	algorithmRegistry.mu.RLock()
	defer algorithmRegistry.mu.RUnlock()
	
	algo, exists := algorithmRegistry.algorithms[name]
	if !exists {
		return nil, fmt.Errorf("algorithm %s not found", name)
	}
	
	return algo, nil
}

// GetAlgorithmMetadata retrieves algorithm metadata
func GetAlgorithmMetadata(name string) (AlgorithmMetadata, error) {
	algorithmRegistry.mu.RLock()
	defer algorithmRegistry.mu.RUnlock()
	
	metadata, exists := algorithmRegistry.metadata[name]
	if !exists {
		return AlgorithmMetadata{}, fmt.Errorf("algorithm %s not found", name)
	}
	
	return metadata, nil
}

// ListAlgorithms returns all registered algorithm names
func ListAlgorithms() []string {
	algorithmRegistry.mu.RLock()
	defer algorithmRegistry.mu.RUnlock()
	
	names := make([]string, 0, len(algorithmRegistry.algorithms))
	for name := range algorithmRegistry.algorithms {
		names = append(names, name)
	}
	
	return names
}

// ListAlgorithmsByCategory returns algorithms in a specific category
func ListAlgorithmsByCategory(category string) []string {
	algorithmRegistry.mu.RLock()
	defer algorithmRegistry.mu.RUnlock()
	
	var names []string
	for name, metadata := range algorithmRegistry.metadata {
		if metadata.Category == category {
			names = append(names, name)
		}
	}
	
	return names
}

// SelectBestAlgorithm selects the best algorithm based on requirements
func SelectBestAlgorithm(category string, requirements AlgorithmRequirements) (string, error) {
	algorithmRegistry.mu.RLock()
	defer algorithmRegistry.mu.RUnlock()
	
	var bestAlgo string
	var bestScore float64
	
	for name, metadata := range algorithmRegistry.metadata {
		if metadata.Category != category {
			continue
		}
		
		// Calculate weighted score based on requirements
		score := calculateAlgorithmScore(metadata, requirements)
		
		if score > bestScore {
			bestScore = score
			bestAlgo = name
		}
	}
	
	if bestAlgo == "" {
		return "", fmt.Errorf("no suitable algorithm found for category %s", category)
	}
	
	return bestAlgo, nil
}

// AlgorithmRequirements specifies requirements for algorithm selection
type AlgorithmRequirements struct {
	MinQuality      float64  // Minimum quality score (0-1)
	MaxTimeMs       int64    // Maximum processing time in milliseconds
	MaxMemoryMB     int64    // Maximum memory usage in megabytes
	PreferredTags   []string // Preferred algorithm tags
	QualityWeight   float64  // Weight for quality in selection (0-1)
	SpeedWeight     float64  // Weight for speed in selection (0-1)
	MemoryWeight    float64  // Weight for memory efficiency (0-1)
}

// calculateAlgorithmScore calculates a weighted score for algorithm selection
func calculateAlgorithmScore(metadata AlgorithmMetadata, requirements AlgorithmRequirements) float64 {
	// Check minimum requirements
	if metadata.QualityScore < requirements.MinQuality {
		return 0
	}
	
	// Normalize weights
	totalWeight := requirements.QualityWeight + requirements.SpeedWeight + requirements.MemoryWeight
	if totalWeight == 0 {
		totalWeight = 3.0 // Equal weights if none specified
		requirements.QualityWeight = 1.0
		requirements.SpeedWeight = 1.0
		requirements.MemoryWeight = 1.0
	}
	
	// Calculate weighted score
	score := (metadata.QualityScore * requirements.QualityWeight +
		metadata.SpeedScore * requirements.SpeedWeight +
		metadata.MemoryScore * requirements.MemoryWeight) / totalWeight
	
	// Bonus for preferred tags
	tagBonus := 0.0
	for _, prefTag := range requirements.PreferredTags {
		for _, algoTag := range metadata.Tags {
			if prefTag == algoTag {
				tagBonus += 0.05
				break
			}
		}
	}
	
	// Cap tag bonus at 0.2
	if tagBonus > 0.2 {
		tagBonus = 0.2
	}
	
	return score + tagBonus
}

// BaseAlgorithm provides a base implementation for common algorithm functionality
type BaseAlgorithm struct {
	name        string
	description string
	processor   func(interface{}, map[string]interface{}) (interface{}, error)
	complexity  ComplexityEstimate
}

// Name returns the algorithm name
func (ba *BaseAlgorithm) Name() string {
	return ba.name
}

// Description returns the algorithm description
func (ba *BaseAlgorithm) Description() string {
	return ba.description
}

// Process executes the algorithm
func (ba *BaseAlgorithm) Process(input interface{}, params map[string]interface{}) (interface{}, error) {
	return ba.processor(input, params)
}

// EstimateComplexity returns the complexity estimate
func (ba *BaseAlgorithm) EstimateComplexity(inputSize int) ComplexityEstimate {
	// Adjust estimates based on input size
	estimate := ba.complexity
	
	// Simple linear scaling for now
	if ba.complexity.TimeComplexity == "O(n)" {
		estimate.EstimatedTimeMs = int64(float64(inputSize) / 1000.0 * float64(ba.complexity.EstimatedTimeMs))
	} else if ba.complexity.TimeComplexity == "O(n log n)" {
		factor := float64(inputSize) / 1000.0 * (1 + 0.3) // Approximate log factor
		estimate.EstimatedTimeMs = int64(factor * float64(ba.complexity.EstimatedTimeMs))
	} else if ba.complexity.TimeComplexity == "O(n²)" {
		factor := float64(inputSize*inputSize) / 1000000.0
		estimate.EstimatedTimeMs = int64(factor * float64(ba.complexity.EstimatedTimeMs))
	}
	
	return estimate
}

// ValidateParams validates algorithm parameters
func (ba *BaseAlgorithm) ValidateParams(params map[string]interface{}) error {
	// Base implementation - can be overridden
	return nil
}

// NewBaseAlgorithm creates a new base algorithm
func NewBaseAlgorithm(name, description string, processor func(interface{}, map[string]interface{}) (interface{}, error), complexity ComplexityEstimate) *BaseAlgorithm {
	return &BaseAlgorithm{
		name:        name,
		description: description,
		processor:   processor,
		complexity:  complexity,
	}
}

// InitializeDefaultAlgorithms registers default algorithms
func InitializeDefaultAlgorithms() {
	// Register sentiment analysis algorithms
	RegisterAlgorithm("sentiment-lexicon", 
		NewBaseAlgorithm(
			"sentiment-lexicon",
			"Fast lexicon-based sentiment analysis",
			nil, // Processor would be implemented
			ComplexityEstimate{
				TimeComplexity:    "O(n)",
				SpaceComplexity:   "O(1)",
				EstimatedTimeMs:   10,
				EstimatedMemoryMB: 1,
			},
		),
		AlgorithmMetadata{
			Category:     "sentiment",
			Version:      "1.0.0",
			Author:       "Caia Tech",
			Tags:         []string{"fast", "simple", "lexicon"},
			QualityScore: 0.7,
			SpeedScore:   0.95,
			MemoryScore:  0.98,
			SupportedInputs: []string{"text"},
		},
	)
	
	RegisterAlgorithm("sentiment-ml-basic",
		NewBaseAlgorithm(
			"sentiment-ml-basic",
			"Machine learning based sentiment analysis",
			nil, // Processor would be implemented
			ComplexityEstimate{
				TimeComplexity:    "O(n)",
				SpaceComplexity:   "O(n)",
				EstimatedTimeMs:   50,
				EstimatedMemoryMB: 10,
			},
		),
		AlgorithmMetadata{
			Category:     "sentiment",
			Version:      "1.0.0",
			Author:       "Caia Tech",
			Tags:         []string{"ml", "accurate"},
			QualityScore: 0.85,
			SpeedScore:   0.7,
			MemoryScore:  0.6,
			SupportedInputs: []string{"text"},
		},
	)
	
	// Register complexity analysis algorithms
	RegisterAlgorithm("complexity-basic",
		NewBaseAlgorithm(
			"complexity-basic",
			"Basic text complexity analysis",
			nil, // Processor would be implemented
			ComplexityEstimate{
				TimeComplexity:    "O(n)",
				SpaceComplexity:   "O(1)",
				EstimatedTimeMs:   20,
				EstimatedMemoryMB: 2,
			},
		),
		AlgorithmMetadata{
			Category:     "complexity",
			Version:      "1.0.0",
			Author:       "Caia Tech",
			Tags:         []string{"fast", "basic"},
			QualityScore: 0.75,
			SpeedScore:   0.9,
			MemoryScore:  0.95,
			SupportedInputs: []string{"text"},
		},
	)
	
	RegisterAlgorithm("complexity-deep",
		NewBaseAlgorithm(
			"complexity-deep",
			"Deep semantic complexity analysis",
			nil, // Processor would be implemented
			ComplexityEstimate{
				TimeComplexity:    "O(n²)",
				SpaceComplexity:   "O(n)",
				EstimatedTimeMs:   200,
				EstimatedMemoryMB: 50,
			},
		),
		AlgorithmMetadata{
			Category:     "complexity",
			Version:      "1.0.0",
			Author:       "Caia Tech",
			Tags:         []string{"deep", "semantic", "comprehensive"},
			QualityScore: 0.95,
			SpeedScore:   0.3,
			MemoryScore:  0.4,
			SupportedInputs: []string{"text", "html", "markdown"},
		},
	)
}