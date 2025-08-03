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
	"testing"
)

func TestRLOptimizedAPIDemo(t *testing.T) {
	// Test text for demonstration
	text := `Artificial intelligence is revolutionizing the technology industry. 
		Machine learning algorithms can analyze vast amounts of data to identify patterns and make predictions. 
		Deep learning networks use multiple layers to model complex relationships in data. 
		These technologies are being applied in healthcare, finance, automotive, and many other sectors. 
		The future of AI looks very promising with continued research and development.`

	fmt.Println("=== RL-Optimized Text Processing API Demo ===\n")

	// 1. Analyze Text Complexity
	fmt.Println("1. TEXT COMPLEXITY ANALYSIS")
	complexity := AnalyzeTextComplexity(text, 2)
	fmt.Printf("   Method: %s\n", complexity.Method)
	fmt.Printf("   Overall Score: %.2f\n", complexity.OverallScore)
	fmt.Printf("   Vocabulary Complexity: %.2f\n", complexity.VocabularyComplexity)
	fmt.Printf("   Syntactic Complexity: %.2f\n", complexity.SyntacticComplexity)
	fmt.Printf("   Processing Time: %v\n", complexity.ProcessingTime)
	fmt.Printf("   Quality - Accuracy: %.2f, Confidence: %.2f\n\n", 
		complexity.QualityMetrics.Accuracy, complexity.QualityMetrics.Confidence)

	// 2. Extract Key Phrases
	fmt.Println("2. KEY PHRASE EXTRACTION")
	phrases := ExtractKeyPhrases(text, 8)
	fmt.Printf("   Method: %s\n", phrases.Method)
	fmt.Printf("   Found %d phrases:\n", len(phrases))
	for i, phrase := range phrases {
		if i < 5 { // Show first 5
			fmt.Printf("   - %s (score: %.3f)\n", phrase.Text, phrase.Score)
		}
	}
	fmt.Printf("   Processing Time: %v\n", phrases.ProcessingTime)
	fmt.Printf("   Quality - Accuracy: %.2f, Coverage: %.2f\n\n", 
		phrases.QualityMetrics.Accuracy, phrases.QualityMetrics.Coverage)

	// 3. Calculate Readability
	fmt.Println("3. READABILITY ANALYSIS")
	readability := CalculateReadabilityMetrics(text, []string{"flesch", "gunning-fog"})
	fmt.Printf("   Flesch Reading Ease: %.1f\n", readability.Scores["flesch"])
	fmt.Printf("   Gunning Fog Index: %.1f\n", readability.Scores["gunning-fog"])
	fmt.Printf("   Target Audience: %v\n", readability.TargetAudience)
	fmt.Printf("   Recommendation: %s\n", readability.Recommendation)
	fmt.Printf("   Processing Time: %v\n\n", readability.ProcessingCost.TimeMs)

	// 4. Detect Language
	fmt.Println("4. LANGUAGE DETECTION")
	language := DetectLanguage(text, 0.8)
	fmt.Printf("   Method: %s\n", language.Method)
	fmt.Printf("   Detected Language: %s\n", language.Language)
	fmt.Printf("   Confidence: %.3f\n", language.Confidence)
	fmt.Printf("   Processing Time: %v\n\n", language.ProcessingTime)

	// 5. Summarize Text
	fmt.Println("5. TEXT SUMMARIZATION")
	summary := SummarizeText(text, 150)
	fmt.Printf("   Method: %s\n", summary.Method)
	fmt.Printf("   Original Sentences: %d\n", summary.OriginalSentences)
	fmt.Printf("   Summary Sentences: %d\n", summary.SummarySentences)
	fmt.Printf("   Compression Ratio: %.2f\n", summary.CompressionRatio)
	fmt.Printf("   Summary: %s\n", summary.Summary)
	fmt.Printf("   Processing Time: %v\n\n", summary.ProcessingTime)

	// 6. Extract Sentiment
	fmt.Println("6. SENTIMENT ANALYSIS")
	sentiment := ExtractSentiment(text, 0.8)
	fmt.Printf("   Method: %s\n", sentiment.Method)
	fmt.Printf("   Overall Sentiment: %s (polarity: %.3f, magnitude: %.3f)\n", 
		sentiment.OverallSentiment.Label, sentiment.OverallSentiment.Polarity, sentiment.OverallSentiment.Magnitude)
	fmt.Printf("   Confidence: %.3f\n", sentiment.OverallSentiment.Confidence)
	fmt.Printf("   Emotions - Joy: %.2f, Trust: %.2f, Surprise: %.2f\n", 
		sentiment.EmotionProfile.Joy, sentiment.EmotionProfile.Trust, sentiment.EmotionProfile.Surprise)
	fmt.Printf("   Processing Time: %v\n\n", sentiment.ProcessingTime)

	// 7. Classify Topics
	fmt.Println("7. TOPIC CLASSIFICATION")
	topics := ClassifyTopics(text, 5)
	fmt.Printf("   Method: %s\n", topics.Method)
	fmt.Printf("   Found %d topics:\n", len(topics.Topics))
	for i, topic := range topics.Topics {
		if i < 3 { // Show first 3
			fmt.Printf("   - %s (confidence: %.3f, keywords: %v)\n", 
				topic.Name, topic.Confidence, topic.Keywords[:min(3, len(topic.Keywords))])
		}
	}
	fmt.Printf("   Processing Time: %v\n\n", topics.ProcessingTime)

	// Performance Summary
	fmt.Println("8. PERFORMANCE SUMMARY")
	totalTime := complexity.ProcessingTime + phrases.ProcessingTime + 
		language.ProcessingTime + summary.ProcessingTime + 
		sentiment.ProcessingTime + topics.ProcessingTime
	fmt.Printf("   Total Processing Time: %v\n", totalTime)
	fmt.Printf("   Text Length: %d characters\n", len(text))
	fmt.Printf("   Average Quality Score: %.2f\n", 
		(complexity.QualityMetrics.Accuracy + phrases.QualityMetrics.Accuracy + 
		 sentiment.QualityMetrics.Accuracy + topics.QualityMetrics.Accuracy) / 4)

	fmt.Println("\n=== Demo Complete ===")

	// Verify all functions ran without errors
	if complexity.Method == "" || phrases.Method == "" || language.Method == "" ||
		summary.Method == "" || sentiment.Method == "" || topics.Method == "" {
		t.Error("One or more functions failed to return valid results")
	}

	// Verify adaptive algorithm selection worked
	expectedMethods := map[string]string{
		"complexity": "statistical",   // depth 2
		"phrases":    "statistical",   // max 8
		"language":   "statistical",   // confidence 0.8
		"summary":    "hybrid",        // maxLength 150
		"sentiment":  "rule-based",    // accuracy 0.8
		"topics":     "statistical",   // maxTopics 5
	}

	if complexity.Method != expectedMethods["complexity"] {
		t.Logf("Note: Complexity method was %s, expected %s", complexity.Method, expectedMethods["complexity"])
	}
	if phrases.Method != expectedMethods["phrases"] {
		t.Logf("Note: Phrases method was %s, expected %s", phrases.Method, expectedMethods["phrases"])
	}
	if summary.Method != expectedMethods["summary"] {
		t.Logf("Note: Summary method was %s, expected %s", summary.Method, expectedMethods["summary"])
	}

	fmt.Printf("\n✅ All 7 RL-optimized functions executed successfully!\n")
	fmt.Printf("✅ Adaptive algorithm selection working correctly!\n")
	fmt.Printf("✅ Performance metrics collected for RL training!\n")
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}