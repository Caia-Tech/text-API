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
	"time"
)

func TestMetricsCollector(t *testing.T) {
	// Start metrics collection
	mc := StartMetricsCollection()
	if mc == nil {
		t.Fatal("StartMetricsCollection returned nil")
	}
	
	// Simulate some processing
	time.Sleep(10 * time.Millisecond)
	
	// Record memory usage
	mc.RecordMemoryUsage()
	
	// Record processing steps
	mc.RecordProcessingTime("step1")
	time.Sleep(5 * time.Millisecond)
	mc.RecordProcessingTime("step2")
	
	// Increment algorithm steps
	for i := 0; i < 5; i++ {
		mc.IncrementAlgorithmSteps()
	}
	
	// Record cache hits
	mc.RecordCacheHit()
	mc.RecordCacheHit()
	
	// Get metrics
	metrics := mc.GetMetrics()
	
	// Validate metrics
	if metrics.TimeElapsed < 15*time.Millisecond {
		t.Errorf("TimeElapsed too short: %v", metrics.TimeElapsed)
	}
	
	if metrics.AlgorithmSteps != 5 {
		t.Errorf("Expected 5 algorithm steps, got %d", metrics.AlgorithmSteps)
	}
	
	if metrics.CacheHits != 2 {
		t.Errorf("Expected 2 cache hits, got %d", metrics.CacheHits)
	}
	
	// Memory peak should be non-negative
	if metrics.MemoryPeak < 0 {
		t.Errorf("MemoryPeak is negative: %d", metrics.MemoryPeak)
	}
}

func TestGlobalMetrics(t *testing.T) {
	// Reset metrics for clean test
	ResetGlobalMetrics()
	
	// Record some function calls
	params1 := map[string]interface{}{
		"depth": 2,
		"quality": 0.85,
	}
	
	metrics1 := ProcessingMetrics{
		TimeElapsed:    50 * time.Millisecond,
		MemoryPeak:     1024 * 1024, // 1MB
		AlgorithmSteps: 10,
		CacheHits:      1,
	}
	
	quality1 := QualityMetrics{
		Accuracy:   0.90,
		Confidence: 0.85,
		Coverage:   0.88,
	}
	
	// Record multiple calls
	RecordFunctionCall("AnalyzeText", params1, metrics1, &quality1)
	
	// Second call with different metrics
	metrics2 := ProcessingMetrics{
		TimeElapsed:    30 * time.Millisecond,
		MemoryPeak:     512 * 1024, // 0.5MB
		AlgorithmSteps: 5,
		CacheHits:      2,
	}
	
	RecordFunctionCall("AnalyzeText", params1, metrics2, &quality1)
	
	// Get global metrics
	global := GetGlobalMetrics()
	
	// Validate function calls
	if global.FunctionCalls["AnalyzeText"] != 2 {
		t.Errorf("Expected 2 calls to AnalyzeText, got %d", global.FunctionCalls["AnalyzeText"])
	}
	
	// Validate performance stats
	stats, exists := global.PerformanceStats["AnalyzeText"]
	if !exists {
		t.Fatal("No performance stats for AnalyzeText")
	}
	
	if stats.CallCount != 2 {
		t.Errorf("Expected 2 calls in stats, got %d", stats.CallCount)
	}
	
	if stats.MinTime != 30*time.Millisecond {
		t.Errorf("Expected min time 30ms, got %v", stats.MinTime)
	}
	
	if stats.MaxTime != 50*time.Millisecond {
		t.Errorf("Expected max time 50ms, got %v", stats.MaxTime)
	}
	
	// Check average time
	expectedAverage := 40 * time.Millisecond
	if stats.AverageTime != expectedAverage {
		t.Errorf("Expected average time %v, got %v", expectedAverage, stats.AverageTime)
	}
	
	// Validate quality distribution
	qualDist, exists := global.QualityScores["AnalyzeText"]
	if !exists {
		t.Fatal("No quality scores for AnalyzeText")
	}
	
	if qualDist.Min != 0.90 || qualDist.Max != 0.90 {
		t.Errorf("Quality score range incorrect: min=%f, max=%f", qualDist.Min, qualDist.Max)
	}
}

func TestMeasureExecution(t *testing.T) {
	// Reset metrics
	ResetGlobalMetrics()
	
	// Define a test function
	testFunc := func() (interface{}, error) {
		time.Sleep(20 * time.Millisecond)
		return "test result", nil
	}
	
	// Measure execution
	params := map[string]interface{}{"test": true}
	result, metrics, err := MeasureExecution("TestFunction", params, testFunc)
	
	// Validate results
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	
	if result != "test result" {
		t.Errorf("Expected 'test result', got %v", result)
	}
	
	// Check metrics
	if metrics.TimeElapsed < 20*time.Millisecond {
		t.Errorf("TimeElapsed too short: %v", metrics.TimeElapsed)
	}
	
	// Check global metrics were updated
	global := GetGlobalMetrics()
	if global.FunctionCalls["TestFunction"] != 1 {
		t.Errorf("Expected 1 call to TestFunction, got %d", global.FunctionCalls["TestFunction"])
	}
}

func TestQualityDistribution(t *testing.T) {
	// Reset metrics
	ResetGlobalMetrics()
	
	// Record multiple quality scores
	qualities := []float64{0.1, 0.3, 0.5, 0.7, 0.9, 0.95}
	
	for _, q := range qualities {
		quality := QualityMetrics{
			Accuracy:   q,
			Confidence: q,
			Coverage:   q,
		}
		
		metrics := ProcessingMetrics{
			TimeElapsed: 10 * time.Millisecond,
		}
		
		RecordFunctionCall("TestQuality", nil, metrics, &quality)
	}
	
	// Get global metrics
	global := GetGlobalMetrics()
	
	// Check quality distribution
	dist, exists := global.QualityScores["TestQuality"]
	if !exists {
		t.Fatal("No quality distribution for TestQuality")
	}
	
	// Check buckets
	if len(dist.Buckets) != 6 {
		t.Errorf("Expected 6 buckets, got %d", len(dist.Buckets))
	}
	
	// Check counts (should have one in each bucket except last)
	expectedCounts := []int64{1, 1, 1, 1, 2, 0} // 0.9 and 0.95 both go in bucket 4
	for i, expected := range expectedCounts {
		if i < len(dist.Counts) && dist.Counts[i] != expected {
			t.Errorf("Bucket %d: expected %d, got %d", i, expected, dist.Counts[i])
		}
	}
	
	// Check min/max
	if dist.Min != 0.1 {
		t.Errorf("Expected min 0.1, got %f", dist.Min)
	}
	
	if dist.Max != 0.95 {
		t.Errorf("Expected max 0.95, got %f", dist.Max)
	}
}

func TestResourceUsageStats(t *testing.T) {
	// Reset metrics
	ResetGlobalMetrics()
	
	// Record multiple calls with resource usage
	for i := 0; i < 5; i++ {
		metrics := ProcessingMetrics{
			TimeElapsed: time.Duration(i+1) * 10 * time.Millisecond,
			MemoryPeak:  int64(i+1) * 1024 * 1024, // 1MB, 2MB, etc.
		}
		
		RecordFunctionCall("ResourceTest", nil, metrics, nil)
	}
	
	// Get global metrics
	global := GetGlobalMetrics()
	
	// Check resource usage
	// Total memory: 1+2+3+4+5 = 15MB
	expectedMemory := int64(15)
	if global.ResourceUsageStats.TotalMemoryMB != expectedMemory {
		t.Errorf("Expected total memory %dMB, got %dMB", 
			expectedMemory, global.ResourceUsageStats.TotalMemoryMB)
	}
	
	// Total CPU time: 10+20+30+40+50 = 150ms
	expectedCPU := int64(150)
	if global.ResourceUsageStats.TotalCPUTimeMs != expectedCPU {
		t.Errorf("Expected total CPU time %dms, got %dms",
			expectedCPU, global.ResourceUsageStats.TotalCPUTimeMs)
	}
}

func TestConcurrentMetricsRecording(t *testing.T) {
	// Reset metrics
	ResetGlobalMetrics()
	
	// Test concurrent access
	done := make(chan bool, 10)
	
	for i := 0; i < 10; i++ {
		go func(id int) {
			metrics := ProcessingMetrics{
				TimeElapsed:    time.Duration(id) * time.Millisecond,
				MemoryPeak:     int64(id) * 1024,
				AlgorithmSteps: id,
				CacheHits:      id % 3,
			}
			
			RecordFunctionCall("ConcurrentTest", nil, metrics, nil)
			done <- true
		}(i)
	}
	
	// Wait for all goroutines
	for i := 0; i < 10; i++ {
		<-done
	}
	
	// Check results
	global := GetGlobalMetrics()
	if global.FunctionCalls["ConcurrentTest"] != 10 {
		t.Errorf("Expected 10 concurrent calls, got %d", global.FunctionCalls["ConcurrentTest"])
	}
}