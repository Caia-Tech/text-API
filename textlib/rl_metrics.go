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
	"runtime"
	"sync"
	"time"
)

// MetricsCollector collects performance metrics for API calls
type MetricsCollector struct {
	mu             sync.RWMutex
	startTime      time.Time
	startMemory    runtime.MemStats
	algorithmSteps int
	cacheHits      int
	peakMemory     int64
	recordings     []MetricRecording
}

// MetricRecording represents a single metric recording
type MetricRecording struct {
	Timestamp   time.Time
	MemoryUsed  int64
	CPUTime     time.Duration
	Description string
}

// APIMetrics represents aggregate API usage metrics
type APIMetrics struct {
	FunctionCalls         map[string]int64                `json:"function_calls"`
	ParameterDistribution map[string]interface{}          `json:"parameter_distribution"`
	PerformanceStats      map[string]PerformanceStats     `json:"performance_stats"`
	QualityScores         map[string]QualityDistribution  `json:"quality_scores"`
	ResourceUsageStats    ResourceUsageStats              `json:"resource_usage"`
}

// PerformanceStats represents performance statistics for a function
type PerformanceStats struct {
	AverageTime   time.Duration `json:"average_time"`
	MinTime       time.Duration `json:"min_time"`
	MaxTime       time.Duration `json:"max_time"`
	P50Time       time.Duration `json:"p50_time"`
	P95Time       time.Duration `json:"p95_time"`
	P99Time       time.Duration `json:"p99_time"`
	CallCount     int64         `json:"call_count"`
	ErrorCount    int64         `json:"error_count"`
	CacheHitRate  float64       `json:"cache_hit_rate"`
}

// QualityDistribution represents quality score distribution
type QualityDistribution struct {
	Mean     float64   `json:"mean"`
	StdDev   float64   `json:"std_dev"`
	Min      float64   `json:"min"`
	Max      float64   `json:"max"`
	Buckets  []float64 `json:"buckets"`
	Counts   []int64   `json:"counts"`
}

// ResourceUsageStats represents aggregate resource usage
type ResourceUsageStats struct {
	TotalMemoryMB    int64   `json:"total_memory_mb"`
	AverageMemoryMB  int64   `json:"average_memory_mb"`
	PeakMemoryMB     int64   `json:"peak_memory_mb"`
	TotalCPUTimeMs   int64   `json:"total_cpu_time_ms"`
	AverageCPUTimeMs int64   `json:"average_cpu_time_ms"`
	CacheHitRate     float64 `json:"cache_hit_rate"`
	NetworkCalls     int64   `json:"network_calls"`
}

// Global metrics storage
var (
	globalMetrics     = &APIMetrics{
		FunctionCalls:         make(map[string]int64),
		ParameterDistribution: make(map[string]interface{}),
		PerformanceStats:      make(map[string]PerformanceStats),
		QualityScores:         make(map[string]QualityDistribution),
		ResourceUsageStats:    ResourceUsageStats{},
	}
	globalMetricsMu sync.RWMutex
)

// StartMetricsCollection creates a new metrics collector
func StartMetricsCollection() *MetricsCollector {
	mc := &MetricsCollector{
		startTime:  time.Now(),
		recordings: make([]MetricRecording, 0),
	}
	
	// Record initial memory state
	runtime.ReadMemStats(&mc.startMemory)
	mc.peakMemory = int64(mc.startMemory.Alloc)
	
	return mc
}

// RecordMemoryUsage records current memory usage
func (mc *MetricsCollector) RecordMemoryUsage() {
	mc.mu.Lock()
	defer mc.mu.Unlock()
	
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	
	currentMemory := int64(memStats.Alloc)
	if currentMemory > mc.peakMemory {
		mc.peakMemory = currentMemory
	}
	
	mc.recordings = append(mc.recordings, MetricRecording{
		Timestamp:   time.Now(),
		MemoryUsed:  currentMemory,
		CPUTime:     time.Since(mc.startTime),
		Description: "memory_checkpoint",
	})
}

// RecordProcessingTime records the time taken for a processing step
func (mc *MetricsCollector) RecordProcessingTime(stepName string) {
	mc.mu.Lock()
	defer mc.mu.Unlock()
	
	mc.recordings = append(mc.recordings, MetricRecording{
		Timestamp:   time.Now(),
		CPUTime:     time.Since(mc.startTime),
		Description: stepName,
	})
}

// IncrementAlgorithmSteps increments the algorithm step counter
func (mc *MetricsCollector) IncrementAlgorithmSteps() {
	mc.mu.Lock()
	defer mc.mu.Unlock()
	mc.algorithmSteps++
}

// RecordCacheHit records a cache hit
func (mc *MetricsCollector) RecordCacheHit() {
	mc.mu.Lock()
	defer mc.mu.Unlock()
	mc.cacheHits++
}

// GetMetrics returns the collected metrics
func (mc *MetricsCollector) GetMetrics() ProcessingMetrics {
	mc.mu.RLock()
	defer mc.mu.RUnlock()
	
	var endMemStats runtime.MemStats
	runtime.ReadMemStats(&endMemStats)
	
	memoryUsed := mc.peakMemory - int64(mc.startMemory.Alloc)
	if memoryUsed < 0 {
		memoryUsed = 0
	}
	
	return ProcessingMetrics{
		TimeElapsed:    time.Since(mc.startTime),
		MemoryPeak:     memoryUsed,
		AlgorithmSteps: mc.algorithmSteps,
		CacheHits:      mc.cacheHits,
	}
}

// RecordFunctionCall records a function call for global metrics
func RecordFunctionCall(functionName string, params map[string]interface{}, metrics ProcessingMetrics, quality *QualityMetrics) {
	globalMetricsMu.Lock()
	defer globalMetricsMu.Unlock()
	
	// Increment call count
	globalMetrics.FunctionCalls[functionName]++
	
	// Update performance stats
	stats, exists := globalMetrics.PerformanceStats[functionName]
	if !exists {
		stats = PerformanceStats{
			MinTime: metrics.TimeElapsed,
			MaxTime: metrics.TimeElapsed,
		}
	}
	
	// Update timing stats
	stats.CallCount++
	if metrics.TimeElapsed < stats.MinTime {
		stats.MinTime = metrics.TimeElapsed
	}
	if metrics.TimeElapsed > stats.MaxTime {
		stats.MaxTime = metrics.TimeElapsed
	}
	
	// Simple average calculation (would use running average in production)
	if stats.AverageTime == 0 {
		stats.AverageTime = metrics.TimeElapsed
	} else {
		stats.AverageTime = (stats.AverageTime*time.Duration(stats.CallCount-1) + metrics.TimeElapsed) / time.Duration(stats.CallCount)
	}
	
	// Update cache hit rate
	if metrics.CacheHits > 0 {
		stats.CacheHitRate = float64(metrics.CacheHits) / float64(stats.CallCount)
	}
	
	globalMetrics.PerformanceStats[functionName] = stats
	
	// Update resource usage
	globalMetrics.ResourceUsageStats.TotalMemoryMB += metrics.MemoryPeak / (1024 * 1024)
	globalMetrics.ResourceUsageStats.TotalCPUTimeMs += int64(metrics.TimeElapsed.Milliseconds())
	
	// Update quality scores if provided
	if quality != nil {
		updateQualityDistribution(functionName, quality)
	}
}

// updateQualityDistribution updates quality score distribution
func updateQualityDistribution(functionName string, quality *QualityMetrics) {
	dist, exists := globalMetrics.QualityScores[functionName]
	if !exists {
		dist = QualityDistribution{
			Min:     quality.Accuracy,
			Max:     quality.Accuracy,
			Buckets: []float64{0.0, 0.2, 0.4, 0.6, 0.8, 1.0},
			Counts:  make([]int64, 6),
		}
	}
	
	// Update min/max
	if quality.Accuracy < dist.Min {
		dist.Min = quality.Accuracy
	}
	if quality.Accuracy > dist.Max {
		dist.Max = quality.Accuracy
	}
	
	// Update bucket counts
	bucketIndex := int(quality.Accuracy * 5)
	if bucketIndex >= len(dist.Counts) {
		bucketIndex = len(dist.Counts) - 1
	}
	dist.Counts[bucketIndex]++
	
	// Simple mean calculation
	totalCalls := int64(0)
	for _, count := range dist.Counts {
		totalCalls += count
	}
	
	if totalCalls > 0 {
		weightedSum := float64(0)
		for i, count := range dist.Counts {
			if i < len(dist.Buckets)-1 {
				bucketMid := (dist.Buckets[i] + dist.Buckets[i+1]) / 2
				weightedSum += bucketMid * float64(count)
			}
		}
		dist.Mean = weightedSum / float64(totalCalls)
	}
	
	globalMetrics.QualityScores[functionName] = dist
}

// GetGlobalMetrics returns the global API metrics
func GetGlobalMetrics() APIMetrics {
	globalMetricsMu.RLock()
	defer globalMetricsMu.RUnlock()
	
	// Create a copy to avoid race conditions
	metricsCopy := APIMetrics{
		FunctionCalls:         make(map[string]int64),
		ParameterDistribution: make(map[string]interface{}),
		PerformanceStats:      make(map[string]PerformanceStats),
		QualityScores:         make(map[string]QualityDistribution),
		ResourceUsageStats:    globalMetrics.ResourceUsageStats,
	}
	
	// Copy maps
	for k, v := range globalMetrics.FunctionCalls {
		metricsCopy.FunctionCalls[k] = v
	}
	for k, v := range globalMetrics.PerformanceStats {
		metricsCopy.PerformanceStats[k] = v
	}
	for k, v := range globalMetrics.QualityScores {
		metricsCopy.QualityScores[k] = v
	}
	
	return metricsCopy
}

// ResetGlobalMetrics resets the global metrics
func ResetGlobalMetrics() {
	globalMetricsMu.Lock()
	defer globalMetricsMu.Unlock()
	
	globalMetrics = &APIMetrics{
		FunctionCalls:         make(map[string]int64),
		ParameterDistribution: make(map[string]interface{}),
		PerformanceStats:      make(map[string]PerformanceStats),
		QualityScores:         make(map[string]QualityDistribution),
		ResourceUsageStats:    ResourceUsageStats{},
	}
}

// MeasureExecution is a helper function to measure function execution
func MeasureExecution(functionName string, params map[string]interface{}, fn func() (interface{}, error)) (interface{}, ProcessingMetrics, error) {
	collector := StartMetricsCollection()
	
	// Execute the function
	result, err := fn()
	
	// Get final metrics
	metrics := collector.GetMetrics()
	
	// Record the call
	RecordFunctionCall(functionName, params, metrics, nil)
	
	return result, metrics, err
}