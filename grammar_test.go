package textlib

import (
	"testing"
)

func TestGrammarAnalysis(t *testing.T) {
	tests := []struct {
		name        string
		text        string
		expectIssues bool
		issueTypes  []string
	}{
		{
			name: "Perfect grammar",
			text: "The cat sat on the mat. It was comfortable.",
			expectIssues: false,
			issueTypes:   []string{},
		},
		{
			name: "Run-on sentence",
			text: "This is a very long sentence that goes on and on and continues to provide information without proper breaks and becomes difficult to read because it has too many clauses and ideas packed into one sentence without appropriate punctuation to separate the different thoughts and concepts.",
			expectIssues: true,
			issueTypes:   []string{"structure"},
		},
		{
			name: "Sentence fragment",
			text: "Walking to the store. Because it was raining.",
			expectIssues: true,
			issueTypes:   []string{"structure"},
		},
		{
			name: "Subject-verb disagreement",
			text: "The cats is sleeping. They was tired.",
			expectIssues: true,
			issueTypes:   []string{"agreement"},
		},
		{
			name: "Punctuation issues",
			text: "Hello , world !What's up?",
			expectIssues: true,
			issueTypes:   []string{"punctuation"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			analysis := AnalyzeGrammar(tt.text)

			if tt.expectIssues {
				if analysis.IssueCount == 0 {
					t.Errorf("Expected issues but found none")
				}

				for _, expectedType := range tt.issueTypes {
					if count, exists := analysis.IssuesByType[expectedType]; !exists || count == 0 {
						t.Errorf("Expected issues of type %s but found none", expectedType)
					}
				}
			} else {
				if analysis.IssueCount > 0 {
					t.Errorf("Expected no issues but found %d", analysis.IssueCount)
				}
			}

			// Grammar score should be between 0 and 1
			if analysis.GrammarScore < 0 || analysis.GrammarScore > 1 {
				t.Errorf("Grammar score %f is outside valid range [0,1]", analysis.GrammarScore)
			}
		})
	}
}

func TestSentenceStructureAnalysis(t *testing.T) {
	tests := []struct {
		name         string
		text         string
		expectedType string
	}{
		{
			name:         "Simple sentence",
			text:         "The dog barks.",
			expectedType: "simple",
		},
		{
			name:         "Compound sentence",
			text:         "The dog barks, and the cat meows.",
			expectedType: "compound",
		},
		{
			name:         "Complex sentence",
			text:         "When the dog barks, the cat runs away.",
			expectedType: "complex",
		},
		{
			name:         "Compound-complex sentence",
			text:         "When the dog barks, the cat runs away, but the bird stays calm.",
			expectedType: "compound-complex",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			analysis := AnalyzeGrammar(tt.text)

			if len(analysis.SentenceStructure) == 0 {
				t.Fatal("No sentence structure analysis found")
			}

			structure := analysis.SentenceStructure[0]
			if structure.Type != tt.expectedType {
				t.Errorf("Expected sentence type %s, got %s", tt.expectedType, structure.Type)
			}
		})
	}
}

func TestPassiveVoiceDetection(t *testing.T) {
	tests := []struct {
		name           string
		text           string
		expectPassive  bool
	}{
		{
			name:          "Active voice",
			text:          "The team completed the project.",
			expectPassive: false,
		},
		{
			name:          "Passive voice",
			text:          "The project was completed by the team.",
			expectPassive: true,
		},
		{
			name:          "Multiple passive constructions",
			text:          "The report was written by John and was reviewed by the manager.",
			expectPassive: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			analysis := AnalyzeGrammar(tt.text)

			hasPassive := len(analysis.PassiveVoiceUse) > 0
			if hasPassive != tt.expectPassive {
				t.Errorf("Expected passive voice: %v, found: %v", tt.expectPassive, hasPassive)
			}
		})
	}
}

func TestPunctuationBalance(t *testing.T) {
	tests := []struct {
		name       string
		text       string
		expectBalance bool
	}{
		{
			name:          "Balanced quotes",
			text:          `He said "Hello world" to everyone.`,
			expectBalance: true,
		},
		{
			name:          "Unbalanced quotes",
			text:          `He said "Hello world to everyone.`,
			expectBalance: false,
		},
		{
			name:          "Balanced parentheses",
			text:          "This is a test (with parentheses) for checking.",
			expectBalance: true,
		},
		{
			name:          "Unbalanced parentheses",
			text:          "This is a test (with parentheses for checking.",
			expectBalance: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			analysis := AnalyzeGrammar(tt.text)

			if tt.name == "Balanced quotes" || tt.name == "Unbalanced quotes" {
				if analysis.QuotationBalance.IsBalanced != tt.expectBalance {
					t.Errorf("Expected quotation balance: %v, got: %v", 
						tt.expectBalance, analysis.QuotationBalance.IsBalanced)
				}
			}

			if tt.name == "Balanced parentheses" || tt.name == "Unbalanced parentheses" {
				if analysis.ParenthesesBalance.IsBalanced != tt.expectBalance {
					t.Errorf("Expected parentheses balance: %v, got: %v", 
						tt.expectBalance, analysis.ParenthesesBalance.IsBalanced)
				}
			}
		})
	}
}

func TestWordinessDetection(t *testing.T) {
	tests := []struct {
		name           string
		text           string
		expectWordiness bool
	}{
		{
			name:            "Concise text",
			text:            "To improve performance, we optimized the code.",
			expectWordiness: false,
		},
		{
			name:            "Wordy text",
			text:            "In order to improve performance, due to the fact that speed matters, we optimized the code.",
			expectWordiness: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			analysis := AnalyzeGrammar(tt.text)

			hasWordiness := len(analysis.WordinessIssues) > 0
			if hasWordiness != tt.expectWordiness {
				t.Errorf("Expected wordiness: %v, found: %v", tt.expectWordiness, hasWordiness)
			}
		})
	}
}

func TestComplexityAssessment(t *testing.T) {
	simpleText := "The cat sits. The dog runs."
	complexText := "Although the cat, which was sleeping peacefully on the windowsill that overlooked the garden, suddenly awoke when the dog, who had been chasing butterflies all morning, came running through the yard."

	simpleAnalysis := AnalyzeGrammar(simpleText)
	complexAnalysis := AnalyzeGrammar(complexText)

	if simpleAnalysis.ComplexityLevel == complexAnalysis.ComplexityLevel {
		t.Error("Simple and complex texts should have different complexity levels")
	}

	// Simple text should have higher grammar score (fewer issues)
	if simpleAnalysis.GrammarScore < complexAnalysis.GrammarScore {
		t.Error("Simple text should generally have a higher grammar score")
	}
}