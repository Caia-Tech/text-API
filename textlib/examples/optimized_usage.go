package main

import (
	"fmt"
	"log"
	"time"
	
	"github.com/Caia-Tech/text-API/textlib"
)

func main() {
	// Example texts
	technicalDoc := `The Redis persistence mechanism uses RDB snapshots and AOF logs. 
	RDB provides point-in-time snapshots of your dataset at specified intervals. 
	AOF logs every write operation received by the server. These operations can 
	be replayed to reconstruct the dataset.`
	
	socialPost := "Just launched our new AI product! 🚀 #MachineLearning #StartupLife 
	@TechCrunch check this out! Feeling excited about the future."
	
	businessEmail := `Dear Team,
	
	Following our Q3 review, I wanted to highlight our key achievements:
	- Revenue increased by 23% YoY
	- Customer satisfaction reached 94%
	- Launched 3 new product features
	
	Action items for next quarter:
	- Expand into European markets
	- Hire 5 additional engineers
	- Implement new CRM system
	
	Best regards,
	Sarah Johnson
	CEO, Tech Corp`
	
	fmt.Println("=== TextLib Optimized Usage Examples ===\n")
	
	// Example 1: Smart Analyze for general purpose
	example1SmartAnalyze(technicalDoc)
	
	// Example 2: Domain-optimized analysis
	example2DomainOptimized(businessEmail)
	
	// Example 3: Quick insights for social media
	example3QuickInsights(socialPost)
	
	// Example 4: Validated extraction
	example4ValidatedExtraction(technicalDoc)
	
	// Example 5: Performance comparison
	example5PerformanceComparison(businessEmail)
	
	// Example 6: Batch processing
	example6BatchProcessing()
}

func example1SmartAnalyze(text string) {
	fmt.Println("Example 1: Smart Analyze (General Purpose)")
	fmt.Println("-----------------------------------------")
	
	start := time.Now()
	result := textlib.SmartAnalyze(text)
	elapsed := time.Since(start)
	
	fmt.Printf("Analysis completed in %v\n", elapsed)
	fmt.Printf("Entities found: %d\n", len(result.Entities))
	for _, entity := range result.Entities {
		fmt.Printf("  - %s (%s)\n", entity.Text, entity.Type)
	}
	
	fmt.Printf("Readability: %s (Grade level: %.1f)\n", 
		result.Readability.Complexity, result.Readability.GradeLevel)
	
	fmt.Printf("Top keywords: %v\n", result.Keywords[:min(5, len(result.Keywords))])
	fmt.Printf("Document type: %s\n", result.Structure.DocumentType)
	fmt.Printf("Functions used: %v\n", result.ProcessingInfo.FunctionsUsed)
	fmt.Println()
}

func example2DomainOptimized(text string) {
	fmt.Println("Example 2: Domain-Optimized Analysis (Business)")
	fmt.Println("-----------------------------------------------")
	
	result := textlib.DomainOptimizedAnalyze(text, "business")
	
	fmt.Printf("Optimization used: %s\n", result.ProcessingInfo.OptimizationUsed)
	fmt.Printf("Sentiment: %s (confidence: %.2f)\n", 
		result.Sentiment.Tone, result.Sentiment.Confidence)
	
	// Business-specific insights
	fmt.Println("Key metrics found:")
	for _, entity := range result.Entities {
		if entity.Type == "MONEY" || entity.Type == "PERCENT" {
			fmt.Printf("  - %s\n", entity.Text)
		}
	}
	
	fmt.Printf("Action items detected: %d\n", countActionItems(result.Keywords))
	fmt.Println()
}

func example3QuickInsights(text string) {
	fmt.Println("Example 3: Quick Insights (Social Media)")
	fmt.Println("----------------------------------------")
	
	start := time.Now()
	insights := textlib.QuickInsights(text)
	elapsed := time.Since(start)
	
	fmt.Printf("Analysis completed in %v (optimized for speed)\n", elapsed)
	fmt.Printf("Sentiment: %s\n", insights.Sentiment.Tone)
	fmt.Printf("Top keywords: %v\n", insights.TopKeywords)
	fmt.Printf("Key entities: ")
	for _, entity := range insights.KeyEntities {
		fmt.Printf("%s ", entity.Text)
	}
	fmt.Println()
	fmt.Printf("Summary: %s\n", insights.Summary)
	fmt.Println()
}

func example4ValidatedExtraction(text string) {
	fmt.Println("Example 4: Validated Extraction (Higher Accuracy)")
	fmt.Println("-------------------------------------------------")
	
	// Compare validated vs non-validated
	regularEntities := textlib.ExtractAdvancedEntities(text)
	validatedEntities := textlib.ValidatedExtraction(text)
	
	fmt.Printf("Regular extraction: %d entities\n", len(regularEntities))
	fmt.Printf("Validated extraction: %d entities (cleaned and merged)\n", len(validatedEntities))
	
	fmt.Println("Validated entities:")
	for _, entity := range validatedEntities {
		fmt.Printf("  - %s (%s)\n", entity.Text, entity.Type)
	}
	fmt.Println()
}

func example5PerformanceComparison(text string) {
	fmt.Println("Example 5: Performance Comparison")
	fmt.Println("--------------------------------")
	
	// Method 1: Individual function calls (inefficient)
	start1 := time.Now()
	entities := textlib.ExtractAdvancedEntities(text)
	readability := textlib.CalculateFleschReadingEase(text)
	stats := textlib.CalculateTextStatistics(text)
	keywords := extractSimpleKeywords(text)
	method1Time := time.Since(start1)
	
	// Method 2: SmartAnalyze (optimized)
	start2 := time.Now()
	result := textlib.SmartAnalyze(text)
	method2Time := time.Since(start2)
	
	fmt.Printf("Method 1 (Individual calls): %v\n", method1Time)
	fmt.Printf("Method 2 (SmartAnalyze): %v\n", method2Time)
	fmt.Printf("Speed improvement: %.1fx faster\n", 
		float64(method1Time.Nanoseconds())/float64(method2Time.Nanoseconds()))
	
	// Results are equivalent
	fmt.Printf("\nResults comparison:\n")
	fmt.Printf("Entities - Method 1: %d, Method 2: %d\n", 
		len(entities), len(result.Entities))
	fmt.Printf("Readability - Method 1: %.1f, Method 2: %.1f\n", 
		readability, result.Readability.FleschScore)
	fmt.Println()
}

func example6BatchProcessing() {
	fmt.Println("Example 6: Batch Processing")
	fmt.Println("--------------------------")
	
	// Sample documents of different types
	documents := []struct {
		text   string
		domain string
	}{
		{"Technical documentation about APIs...", "technical"},
		{"Patient shows symptoms of...", "medical"},
		{"Pursuant to section 5.2 of the agreement...", "legal"},
		{"Check out our summer sale! #shopping", "social"},
		{"Q4 revenue projections indicate...", "business"},
	}
	
	results := make([]textlib.ComprehensiveResult, len(documents))
	
	start := time.Now()
	for i, doc := range documents {
		results[i] = textlib.DomainOptimizedAnalyze(doc.text, doc.domain)
	}
	elapsed := time.Since(start)
	
	fmt.Printf("Processed %d documents in %v\n", len(documents), elapsed)
	fmt.Printf("Average time per document: %v\n", elapsed/time.Duration(len(documents)))
	
	// Summary of results
	for i, result := range results {
		fmt.Printf("Document %d (%s): %d entities, %s complexity\n",
			i+1, documents[i].domain, len(result.Entities), result.Readability.Complexity)
	}
}

// Helper functions

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func countActionItems(keywords []string) int {
	count := 0
	actionWords := []string{"implement", "expand", "hire", "launch", "review", "analyze"}
	
	for _, keyword := range keywords {
		for _, action := range actionWords {
			if strings.Contains(strings.ToLower(keyword), action) {
				count++
				break
			}
		}
	}
	
	return count
}

func extractSimpleKeywords(text string) []string {
	// Simple keyword extraction for comparison
	words := strings.Fields(strings.ToLower(text))
	freq := make(map[string]int)
	
	for _, word := range words {
		cleaned := strings.Trim(word, ".,!?;:")
		if len(cleaned) > 4 {
			freq[cleaned]++
		}
	}
	
	keywords := make([]string, 0, len(freq))
	for word := range freq {
		keywords = append(keywords, word)
	}
	
	return keywords
}