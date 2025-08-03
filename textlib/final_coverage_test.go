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

// Final attempt at simple coverage tests

func TestAnalyzeComplexityBasic(t *testing.T) {
	result := AnalyzeComplexity("The cat sat on the mat.")
	if result == nil {
		t.Error("AnalyzeComplexity returned nil")
	}
}

func TestSegmentTextBasic(t *testing.T) {
	result := SegmentText("First paragraph.\n\nSecond paragraph.", "paragraph")
	if len(result.Segments) == 0 {
		t.Error("SegmentText returned no segments")
	}
}

func TestChunkTextBasic(t *testing.T) {
	result := ChunkText("The quick brown fox jumps over the lazy dog.", "tokens", 5)
	if len(result.Chunks) == 0 {
		t.Error("ChunkText returned no chunks")
	}
}

func TestCalculateSimilarityBasic(t *testing.T) {
	result := CalculateSimilarity("hello world", "hello there", "jaccard")
	if result < 0 || result > 1 {
		t.Errorf("CalculateSimilarity returned out of range value: %f", result)
	}
}

func TestCalculateDiffBasic(t *testing.T) {
	result := CalculateDiff("hello world", "hello there")
	if len(result.Operations) == 0 {
		t.Error("CalculateDiff returned no operations")
	}
}

func TestDetectPatternsBasic(t *testing.T) {
	result := DetectPatterns("The cat sat. The dog ran. The cat jumped.")
	if len(result.WordRepetitions) == 0 {
		t.Error("DetectPatterns found no word repetitions")
	}
}

func TestAnalyzeDialogueBasic(t *testing.T) {
	result := AnalyzeDialogue("'Hello there,' said John. 'How are you?' Mary replied.")
	if len(result.Utterances) == 0 {
		t.Error("AnalyzeDialogue found no utterances")
	}
}

func TestAnalyzeCoherenceBasic(t *testing.T) {
	result := AnalyzeCoherence("The cat sat on the mat. It was black and white.")
	if result.OverallScore < 0 || result.OverallScore > 1 {
		t.Errorf("AnalyzeCoherence returned out of range score: %f", result.OverallScore)
	}
}

func TestFindUnusedFilesBasic(t *testing.T) {
	result := FindUnusedFiles(".")
	// Just test that it doesn't panic
	_ = result
}

func TestExpandExpressionBasic(t *testing.T) {
	result := ExpandExpression("(x + 2)(x - 1)")
	if result == "" {
		t.Error("ExpandExpression returned empty string")
	}
}

func TestFactorExpressionBasic(t *testing.T) {
	result := FactorExpression("x^2 + 3x + 2")
	if result == "" {
		t.Error("FactorExpression returned empty string")
	}
}

func TestFindGCFBasic(t *testing.T) {
	result := findGCF([]int{12, 18, 24})
	if result <= 0 {
		t.Errorf("findGCF returned invalid result: %d", result)
	}
}

func TestDetectMathPatternsBasic(t *testing.T) {
	result := DetectMathPatterns("1, 3, 5, 7, 9")
	if len(result.ArithmeticSequences) == 0 {
		t.Error("DetectMathPatterns found no arithmetic sequences")
	}
}

func TestVerifyImageIntegrityBasic(t *testing.T) {
	// Test with a simple path
	result := VerifyImageIntegrity("nonexistent.jpg")
	// Should return an error for non-existent file
	_ = result
}

func TestValidateVideoCodecBasic(t *testing.T) {
	// Test with a simple path  
	result := ValidateVideoCodec("nonexistent.mp4")
	// Should return an error for non-existent file
	_ = result
}

func TestOptimalChunkSizeBasic(t *testing.T) {
	result := OptimalChunkSize("This is a test text for chunking analysis.")
	if result <= 0 {
		t.Errorf("OptimalChunkSize returned invalid size: %d", result)
	}
}

func TestFindNaturalBoundariesBasic(t *testing.T) {
	result := FindNaturalBoundaries("Sentence one. Sentence two. Sentence three.")
	if len(result) == 0 {
		t.Error("FindNaturalBoundaries found no boundaries")
	}
}

func TestMergeSmallSegmentsBasic(t *testing.T) {
	segments := []TextSegment{
		{Content: "Short", WordCount: 1},
		{Content: "Also short", WordCount: 2},
		{Content: "This is a longer segment with more words", WordCount: 8},
	}
	result := MergeSmallSegments(segments, 5)
	if len(result) == 0 {
		t.Error("MergeSmallSegments returned no segments")
	}
}

func TestBalanceChunksBasic(t *testing.T) {
	chunks := []TextChunk{
		{Content: "Short chunk", TokenCount: 10},
		{Content: "This is a much longer chunk with many more tokens and words", TokenCount: 50},
	}
	result := BalanceChunks(chunks, 30)
	if len(result) == 0 {
		t.Error("BalanceChunks returned no chunks")
	}
}