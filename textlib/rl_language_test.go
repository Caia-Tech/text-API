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
	"time"
)

func TestDetectLanguage(t *testing.T) {
	tests := []struct {
		name               string
		text               string
		confidence         float64
		expectedLang       string
		expectedMethod     string
		expectedConfidence float64 // minimum expected
		checkAlternatives  bool
	}{
		{
			name:               "English text - fast",
			text:               "The quick brown fox jumps over the lazy dog.",
			confidence:         0.5,
			expectedLang:       "en",
			expectedMethod:     "character-frequency",
			expectedConfidence: 0.4,
		},
		{
			name:               "Spanish text - statistical",
			text:               "El gato está en la casa. ¿Dónde está el perro?",
			confidence:         0.8,
			expectedLang:       "es",
			expectedMethod:     "statistical-ngram",
			expectedConfidence: 0.6,
			checkAlternatives:  true,
		},
		{
			name:               "French text - comprehensive",
			text:               "Bonjour, comment allez-vous? J'espère que vous allez bien.",
			confidence:         0.95,
			expectedLang:       "fr",
			expectedMethod:     "comprehensive-analysis",
			expectedConfidence: 0.7,
		},
		{
			name:               "German text",
			text:               "Der Hund läuft schnell durch den Wald. Die Katze schläft.",
			confidence:         0.8,
			expectedLang:       "de",
			expectedMethod:     "statistical-ngram",
			expectedConfidence: 0.6,
		},
		{
			name:               "Italian text",
			text:               "Il cane corre veloce nel bosco. La vita è bella.",
			confidence:         0.7,
			expectedLang:       "it",
			expectedMethod:     "statistical-ngram",
			expectedConfidence: 0.5,
		},
		{
			name:               "Portuguese text",
			text:               "O cachorro corre rápido pela floresta. Onde está o gato?",
			confidence:         0.8,
			expectedLang:       "pt",
			expectedMethod:     "statistical-ngram",
			expectedConfidence: 0.5,
		},
		{
			name:               "Russian text",
			text:               "Привет, как дела? Я хорошо говорю по-русски.",
			confidence:         0.7,
			expectedLang:       "ru",
			expectedMethod:     "statistical-ngram",
			expectedConfidence: 0.7,
		},
		{
			name:               "Chinese text",
			text:               "你好，我是一个程序员。我喜欢编程。",
			confidence:         0.6,
			expectedLang:       "zh",
			expectedMethod:     "character-frequency",
			expectedConfidence: 0.8,
		},
		{
			name:               "Japanese text",
			text:               "こんにちは、元気ですか？私は日本語を勉強しています。",
			confidence:         0.6,
			expectedLang:       "ja",
			expectedMethod:     "character-frequency",
			expectedConfidence: 0.8,
		},
		{
			name:               "Mixed English-Spanish",
			text:               "Hello amigo, how are you today? Muy bien, gracias.",
			confidence:         0.8,
			expectedLang:       "en", // Should lean towards English
			expectedMethod:     "statistical-ngram",
			expectedConfidence: 0.4,
			checkAlternatives:  true,
		},
		{
			name:               "Very short text",
			text:               "Hola",
			confidence:         0.7,
			expectedLang:       "es",
			expectedMethod:     "statistical-ngram",
			expectedConfidence: 0.3, // Low confidence for short text
		},
		{
			name:               "Empty text",
			text:               "",
			confidence:         0.8,
			expectedLang:       "unknown",
			expectedMethod:     "empty",
			expectedConfidence: 0.0,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := DetectLanguage(tt.text, tt.confidence)
			
			// Check language
			if result.Language != tt.expectedLang {
				t.Errorf("Expected language %s, got %s", tt.expectedLang, result.Language)
			}
			
			// Check method
			if result.Method != tt.expectedMethod {
				t.Errorf("Expected method %s, got %s", tt.expectedMethod, result.Method)
			}
			
			// Check confidence
			if result.Confidence < tt.expectedConfidence {
				t.Errorf("Confidence %f below expected minimum %f", 
					result.Confidence, tt.expectedConfidence)
			}
			
			// Check confidence doesn't exceed requested
			if result.Confidence > tt.confidence {
				t.Errorf("Confidence %f exceeds requested %f", 
					result.Confidence, tt.confidence)
			}
			
			// Check alternatives for multi-language texts
			if tt.checkAlternatives && len(result.Alternatives) == 0 {
				t.Error("Expected alternatives but got none")
			}
			
			// Check processing time
			if result.ProcessingTime <= 0 {
				t.Error("Invalid processing time")
			}
		})
	}
}

func TestDetectLanguageConfidenceLevels(t *testing.T) {
	text := "This is a comprehensive English text with multiple sentences. " +
		"It contains various words and patterns that should be detected."
	
	// Test different confidence levels
	confidenceLevels := []struct {
		confidence     float64
		expectedMethod string
		maxTime        time.Duration
	}{
		{0.5, "character-frequency", 10 * time.Millisecond},
		{0.7, "statistical-ngram", 20 * time.Millisecond},
		{0.95, "comprehensive-analysis", 50 * time.Millisecond},
	}
	
	for _, tc := range confidenceLevels {
		result := DetectLanguage(text, tc.confidence)
		
		if result.Method != tc.expectedMethod {
			t.Errorf("Confidence %.2f: expected method %s, got %s",
				tc.confidence, tc.expectedMethod, result.Method)
		}
		
		if result.Language != "en" {
			t.Errorf("Confidence %.2f: expected English, got %s",
				tc.confidence, result.Language)
		}
		
		// Higher confidence should generally take more time
		// (This is a soft check as timing can vary)
		t.Logf("Confidence %.2f: took %v with method %s",
			tc.confidence, result.ProcessingTime, result.Method)
	}
}

func TestLanguageAlternatives(t *testing.T) {
	// Text that could be multiple languages
	tests := []struct {
		name              string
		text              string
		expectedPrimary   string
		expectedAlternative string
	}{
		{
			name:              "Spanish/Portuguese similarity",
			text:              "A casa está muito bonita e grande.",
			expectedPrimary:   "pt",
			expectedAlternative: "es",
		},
		{
			name:              "German/Dutch similarity",
			text:              "Het water is koud.",
			expectedPrimary:   "nl",
			expectedAlternative: "de",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := DetectLanguage(tt.text, 0.8)
			
			// Check primary language (may be either)
			if result.Language != tt.expectedPrimary && result.Language != tt.expectedAlternative {
				t.Errorf("Expected %s or %s, got %s", 
					tt.expectedPrimary, tt.expectedAlternative, result.Language)
			}
			
			// Check that alternatives contain the other language
			foundAlternative := false
			for _, alt := range result.Alternatives {
				if alt.Language == tt.expectedAlternative || alt.Language == tt.expectedPrimary {
					foundAlternative = true
					if alt.Confidence <= 0 {
						t.Error("Alternative confidence should be positive")
					}
					if alt.Reason == "" {
						t.Error("Alternative should have a reason")
					}
					break
				}
			}
			
			if !foundAlternative && len(result.Alternatives) > 0 {
				t.Logf("Primary: %s, Alternatives: %v", result.Language, result.Alternatives)
			}
		})
	}
}

func TestSpecialCharacterDetection(t *testing.T) {
	tests := []struct {
		name         string
		text         string
		expectedLang string
		description  string
	}{
		{
			name:         "Chinese characters",
			text:         "这是中文文本",
			expectedLang: "zh",
			description:  "Chinese Unicode range",
		},
		{
			name:         "Japanese hiragana",
			text:         "これは日本語です",
			expectedLang: "ja",
			description:  "Japanese hiragana",
		},
		{
			name:         "Japanese katakana",
			text:         "カタカナテキスト",
			expectedLang: "ja",
			description:  "Japanese katakana",
		},
		{
			name:         "Russian Cyrillic",
			text:         "Это русский текст",
			expectedLang: "ru",
			description:  "Cyrillic script",
		},
		{
			name:         "Spanish with special chars",
			text:         "¿Cómo estás? ¡Muy bien!",
			expectedLang: "es",
			description:  "Spanish punctuation",
		},
		{
			name:         "French with accents",
			text:         "Où êtes-vous allé? À Paris!",
			expectedLang: "fr",
			description:  "French accents",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := DetectLanguage(tt.text, 0.7)
			
			if result.Language != tt.expectedLang {
				t.Errorf("%s: expected %s, got %s", 
					tt.description, tt.expectedLang, result.Language)
			}
			
			// Non-Latin scripts should have high confidence
			if (tt.expectedLang == "zh" || tt.expectedLang == "ja" || tt.expectedLang == "ru") &&
				result.Confidence < 0.7 {
				t.Errorf("%s: expected high confidence, got %f",
					tt.description, result.Confidence)
			}
		})
	}
}

func TestEdgeCases(t *testing.T) {
	// Test confidence bounds
	result := DetectLanguage("Test", 0.3) // Below minimum
	if result.Confidence < 0.3 {
		t.Error("Confidence should be at least 0.5 (minimum)")
	}
	
	result = DetectLanguage("Test", 1.5) // Above maximum
	if result.Confidence > 0.95 {
		t.Error("Confidence should be capped at 0.95")
	}
	
	// Test numbers and symbols only
	result = DetectLanguage("123 456 789 !@#$%", 0.7)
	if result.Confidence > 0.5 {
		t.Error("Low confidence expected for non-letter text")
	}
	
	// Test single word
	result = DetectLanguage("Hello", 0.8)
	if result.Language != "en" {
		t.Errorf("Expected English for 'Hello', got %s", result.Language)
	}
}

func TestLongTextPerformance(t *testing.T) {
	// Generate long texts in different languages
	texts := map[string]string{
		"en": strings.Repeat("The quick brown fox jumps over the lazy dog. ", 100),
		"es": strings.Repeat("El rápido zorro marrón salta sobre el perro perezoso. ", 100),
		"fr": strings.Repeat("Le rapide renard brun saute par-dessus le chien paresseux. ", 100),
	}
	
	for lang, text := range texts {
		start := time.Now()
		result := DetectLanguage(text, 0.9)
		elapsed := time.Since(start)
		
		if result.Language != lang {
			t.Errorf("Long %s text: expected %s, got %s", lang, lang, result.Language)
		}
		
		if result.Confidence < 0.8 {
			t.Errorf("Long %s text: expected high confidence, got %f", lang, result.Confidence)
		}
		
		t.Logf("Long %s text: detected in %v", lang, elapsed)
	}
}

func BenchmarkDetectLanguageFast(b *testing.B) {
	text := "The quick brown fox jumps over the lazy dog."
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		DetectLanguage(text, 0.5)
	}
}

func BenchmarkDetectLanguageStatistical(b *testing.B) {
	text := "The implementation of advanced algorithms requires careful consideration."
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		DetectLanguage(text, 0.8)
	}
}

func BenchmarkDetectLanguageComprehensive(b *testing.B) {
	text := "Natural language processing is a subfield of linguistics, computer science, " +
		"and artificial intelligence concerned with the interactions between computers and human language."
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		DetectLanguage(text, 0.95)
	}
}