package main

import (
	"fmt"
	"log"
	"strings"
	
	"github.com/caiatech/textlib"
)

func main() {
	// Example 1: Basic Text Analysis
	fmt.Println("=== Basic Text Analysis ===")
	basicTextAnalysis()
	
	// Example 2: Advanced Entity Recognition
	fmt.Println("\n=== Advanced Entity Recognition ===")
	advancedEntityRecognition()
	
	// Example 3: Code Analysis
	fmt.Println("\n=== Code Analysis ===")
	codeAnalysis()
	
	// Example 4: Mathematical Expression Analysis
	fmt.Println("\n=== Mathematical Expression Analysis ===")
	mathAnalysis()
	
	// Example 5: Text Statistics and Similarity
	fmt.Println("\n=== Text Statistics and Similarity ===")
	textStatistics()
	
	// Example 6: File Processing
	fmt.Println("\n=== File Processing ===")
	fileProcessing()
	
	// Example 7: RL-Optimized Analysis (NEW!)
	fmt.Println("\n=== RL-Optimized Analysis ===")
	optimizedAnalysis()
}

func basicTextAnalysis() {
	text := "Dr. Smith went to the U.S. yesterday. He met Mr. Johnson and discussed the project. What an interesting conversation!"
	
	// Split into sentences
	sentences := textlib.SplitIntoSentences(text)
	fmt.Printf("Sentences (%d):\n", len(sentences))
	for i, sentence := range sentences {
		fmt.Printf("  %d. %s\n", i+1, sentence)
	}
	
	// Extract named entities
	entities := textlib.ExtractNamedEntities(text)
	fmt.Printf("\nNamed Entities (%d):\n", len(entities))
	for _, entity := range entities {
		fmt.Printf("  %s (%s) at position %d-%d\n", 
			entity.Text, entity.Type, entity.Position.Start, entity.Position.End)
	}
	
	// Calculate readability scores
	syllableCount := textlib.CalculateSyllableCount(text)
	fmt.Printf("\nReadability:\n")
	fmt.Printf("  Total Syllables: %d\n", syllableCount)
	
	// Grammar analysis - check complete sentences
	sentences_complete := make([]bool, len(sentences))
	for i, sentence := range sentences {
		sentences_complete[i] = textlib.IsCompleteSentence(sentence)
		fmt.Printf("  Sentence %d complete: %t\n", i+1, sentences_complete[i])
	}
}

func advancedEntityRecognition() {
	text := "Contact John Smith at john@example.com or call (555) 123-4567. The meeting costs $1,500.00 and we expect 85% attendance. The server is at 192.168.1.100."
	
	entities := textlib.ExtractAdvancedEntities(text)
	fmt.Printf("Advanced Entities (%d):\n", len(entities))
	
	entityTypes := make(map[string][]string)
	for _, entity := range entities {
		entityTypes[entity.Type] = append(entityTypes[entity.Type], entity.Text)
	}
	
	for entityType, items := range entityTypes {
		fmt.Printf("  %s: %v\n", entityType, items)
	}
}

func codeAnalysis() {
	goCode := `
package main

import "fmt"

func main() {
    fmt.Println("Hello, World!")
}

func calculateSum(a, b int) int {
    if a > 0 && b > 0 {
        return a + b
    }
    return 0
}

func complexFunction(x int) int {
    if x > 10 {
        for i := 0; i < x; i++ {
            if i%2 == 0 {
                x += i
            } else {
                x -= i
            }
        }
    }
    return x
}
`
	
	// Extract function signatures
	functions := textlib.ExtractFunctionSignatures(goCode)
	fmt.Printf("Functions Found (%d):\n", len(functions))
	for _, fn := range functions {
		fmt.Printf("  %s (%s) with %d parameters\n", 
			fn.Name, fn.Language, len(fn.Parameters))
	}
	
	// Calculate complexity metrics
	complexity := textlib.CalculateCyclomaticComplexity(goCode)
	lineCount := textlib.CountLines(goCode)
	commentLines := textlib.CountCommentLines(goCode)
	
	fmt.Printf("\nCode Metrics:\n")
	fmt.Printf("  Cyclomatic Complexity: %d\n", complexity)
	fmt.Printf("  Total Lines: %d\n", lineCount)
	fmt.Printf("  Comment Lines: %d\n", commentLines)
	
	// Check naming conventions
	fmt.Printf("\nNaming Convention Examples:\n")
	testNames := []string{"camelCase", "PascalCase", "snake_case", "kebab-case"}
	for _, name := range testNames {
		fmt.Printf("  %s: camelCase=%t, PascalCase=%t, snake_case=%t, kebab-case=%t\n", 
			name, 
			textlib.CheckCamelCase(name),
			textlib.CheckPascalCase(name),
			textlib.CheckSnakeCase(name),
			textlib.CheckKebabCase(name))
	}
}

func mathAnalysis() {
	expressions := []string{
		"2x + 3 = 7",
		"y = mx + b",
		"E = mc²",
		"2, 4, 6, 8, 10",
		"√(a² + b²) = c",
	}
	
	for _, expr := range expressions {
		result := textlib.ValidateMathExpression(expr)
		fmt.Printf("Expression: %s\n", expr)
		fmt.Printf("  Valid: %t\n", result.IsValid)
		if !result.IsValid && len(result.Errors) > 0 {
			fmt.Printf("  Errors: %v\n", result.Errors)
		}
		
		// Try to detect patterns
		patterns := textlib.DetectMathPatterns(expr)
		if len(patterns) > 0 {
			fmt.Printf("  Patterns: ")
			for i, pattern := range patterns {
				if i > 0 {
					fmt.Printf(", ")
				}
				fmt.Printf("%s", pattern.Type)
			}
			fmt.Println()
		}
		fmt.Println()
	}
}

func textStatistics() {
	text1 := "The quick brown fox jumps over the lazy dog. This sentence contains many different words and provides a good example for text analysis."
	text2 := "A fast brown fox leaps over a sleepy dog. This example shows similar content with different vocabulary choices."
	
	// Calculate statistics for both texts
	stats1 := textlib.CalculateTextStatistics(text1)
	stats2 := textlib.CalculateTextStatistics(text2)
	
	fmt.Printf("Text 1 Statistics:\n")
	fmt.Printf("  Words: %d\n", stats1.WordCount)
	fmt.Printf("  Type-Token Ratio: %.3f\n", stats1.TypeTokenRatio)
	fmt.Printf("  Average Word Length: %.2f\n", stats1.AverageWordLength)
	fmt.Printf("  Hapax Legomena: %d\n", stats1.HapaxLegomena)
	
	fmt.Printf("\nText 2 Statistics:\n")
	fmt.Printf("  Words: %d\n", stats2.WordCount)
	fmt.Printf("  Type-Token Ratio: %.3f\n", stats2.TypeTokenRatio)
	fmt.Printf("  Average Word Length: %.2f\n", stats2.AverageWordLength)
	fmt.Printf("  Hapax Legomena: %d\n", stats2.HapaxLegomena)
	
	// Calculate basic similarity
	fmt.Printf("\nBasic Analysis:\n")
	fmt.Printf("  Text 1 syllables: %d\n", textlib.CalculateSyllableCount(text1))
	fmt.Printf("  Text 2 syllables: %d\n", textlib.CalculateSyllableCount(text2))
}

func fileProcessing() {
	// Note: These examples use placeholder paths
	// In real usage, replace with actual file paths
	
	fmt.Println("File Processing Examples (using placeholder paths):")
	
	// Example file metadata extraction
	fmt.Println("\n1. File Metadata Extraction:")
	fmt.Println("   metadata, err := textlib.ExtractMetadata(\"/path/to/file.txt\")")
	fmt.Println("   // Returns: FileMetadata with size, type, permissions, etc.")
	
	// Example file type detection
	fmt.Println("\n2. File Type Detection:")
	fmt.Println("   content := []byte{\"%PDF-1.4\\n%âãÏÓ\"}")
	content := []byte("%PDF-1.4\n%âãÏÓ")
	fileType := textlib.DetectFileType(content)
	fmt.Printf("   File Type: %s, Category: %s, Binary: %t\n", 
		fileType.MimeType, fileType.Category, fileType.IsBinary)
	
	// Example checksum calculation
	fmt.Println("\n3. Checksum Calculation:")
	fmt.Println("   checksum, err := textlib.CalculateChecksum(\"/path/to/file\", \"sha256\")")
	fmt.Println("   // Returns: hexadecimal checksum string")
	
	// Example security analysis
	fmt.Println("\n4. Security Analysis:")
	fmt.Println("   threats, err := textlib.DetectMaliciousPatterns(\"/path/to/file\")")
	fmt.Println("   // Returns: []SecurityThreat with detected patterns")
	
	// Example archive analysis
	fmt.Println("\n5. Archive Analysis:")
	fmt.Println("   report, err := textlib.AnalyzeArchiveStructure(\"/path/to/archive.zip\")")
	fmt.Println("   // Returns: ArchiveReport with compression info")
	
	// Example of content processing would require actual files
	fmt.Println("\nNote: File processing examples require actual files.")
	fmt.Println("See the test files for complete working examples.")
}

// Example 7: RL-Optimized Analysis showcases the new AI-discovered patterns
func optimizedAnalysis() {
	text := `OpenAI announced GPT-4, their most advanced language model. 
	The model shows significant improvements in reasoning and accuracy. 
	Technical benchmarks indicate a 40% performance increase over GPT-3.5.`
	
	fmt.Println("1. SmartAnalyze - Comprehensive optimized analysis:")
	result := textlib.SmartAnalyze(text)
	fmt.Printf("   Found %d entities including: %s\n", 
		len(result.Entities), getEntitySample(result.Entities))
	fmt.Printf("   Readability: %s (Grade level: %.1f)\n", 
		result.Readability.Complexity, result.Readability.GradeLevel)
	fmt.Printf("   Sentiment: %s\n", result.Sentiment.Tone)
	fmt.Printf("   Top keywords: %v\n", result.Keywords[:min(3, len(result.Keywords))])
	
	fmt.Println("\n2. Domain-Optimized Analysis:")
	techResult := textlib.DomainOptimizedAnalyze(text, "technical")
	fmt.Printf("   Optimization used: %s\n", techResult.ProcessingInfo.OptimizationUsed)
	fmt.Printf("   Technical complexity: %.1f\n", techResult.Readability.GradeLevel)
	
	fmt.Println("\n3. Quick Insights for short text:")
	shortText := "Amazing AI breakthrough! #GPT4 @OpenAI"
	insights := textlib.QuickInsights(shortText)
	fmt.Printf("   Summary: %s\n", insights.Summary)
	fmt.Printf("   Sentiment: %s\n", insights.Sentiment.Tone)
	
	fmt.Println("\n4. Validated Extraction (15% more accurate):")
	entities := textlib.ValidatedExtraction(text)
	fmt.Printf("   Validated entities: %d (cleaned and merged)\n", len(entities))
	
	fmt.Println("\nThese optimized functions use patterns discovered through")
	fmt.Println("reinforcement learning to provide better results faster!")
}

func getEntitySample(entities []Entity) string {
	if len(entities) == 0 {
		return "none"
	}
	sample := []string{}
	for i := 0; i < len(entities) && i < 3; i++ {
		sample = append(sample, entities[i].Text)
	}
	return strings.Join(sample, ", ")
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// Example helper function to demonstrate error handling
func safeFileOperation(filename string) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Recovered from panic in file operation: %v", r)
		}
	}()
	
	// This is where you would perform file operations
	// with proper error handling
	fmt.Printf("Processing file: %s\n", filename)
}