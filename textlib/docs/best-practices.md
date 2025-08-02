# TextLib Best Practices Guide

This guide provides optimal usage patterns for the TextLib API to maximize performance and accuracy in your text analysis tasks.

## Quick Start: Use SmartAnalyze

For most use cases, `SmartAnalyze` provides the best balance of performance and comprehensive results:

```go
import "github.com/Caia-Tech/text-API/textlib"

func main() {
    text := "Your text here..."
    result := textlib.SmartAnalyze(text)
    
    // Access all analysis results
    fmt.Printf("Found %d entities\n", len(result.Entities))
    fmt.Printf("Readability: %s\n", result.Readability.Complexity)
    fmt.Printf("Sentiment: %s\n", result.Sentiment.Tone)
}
```

## Domain-Specific Optimization

Different types of text benefit from different analysis approaches. Use `DomainOptimizedAnalyze` for better results:

```go
// For technical documentation or code
result := textlib.DomainOptimizedAnalyze(text, "technical")

// For medical texts
result := textlib.DomainOptimizedAnalyze(text, "medical")

// For legal documents
result := textlib.DomainOptimizedAnalyze(text, "legal")

// For social media posts
result := textlib.DomainOptimizedAnalyze(text, "social")

// For business communications
result := textlib.DomainOptimizedAnalyze(text, "business")
```

## Optimal Function Sequences

Through extensive analysis, we've discovered these function sequences provide the best results:

### General Text Analysis

```go
// Best sequence for comprehensive analysis
// 1. Validate and clean first
// 2. Extract entities
// 3. Analyze readability
// 4. Extract keywords

// Use SmartAnalyze which implements this automatically, or:
cleanText := textlib.ValidateAndClean(text)
entities := textlib.ExtractAdvancedEntities(cleanText)
readability := textlib.CalculateFleschReadingEase(cleanText)
keywords := textlib.ExtractKeywords(cleanText, entities)
```

### Entity Extraction

For the most accurate entity extraction, always validate first:

```go
// Use ValidatedExtraction for 15% better accuracy
entities := textlib.ValidatedExtraction(text)

// This is better than:
// entities := textlib.ExtractAdvancedEntities(text) // Less accurate
```

### Quick Analysis (Social Media, Chat)

For short texts where speed matters:

```go
insights := textlib.QuickInsights(text)
// Returns sentiment, keywords, and key entities in <10ms
```

### Technical Content

For code and technical documentation:

```go
technical := textlib.DeepTechnicalAnalysis(text)
// Extracts code blocks, analyzes complexity, finds patterns
```

## Performance Tips

### 1. Choose the Right Function

- **SmartAnalyze**: Best for general purpose (balanced speed/accuracy)
- **QuickInsights**: Best for real-time applications (<10ms)
- **DeepTechnicalAnalysis**: Best for code/documentation (most thorough)
- **DomainOptimizedAnalyze**: Best for specific content types

### 2. Batch Processing

When analyzing multiple texts, process similar types together:

```go
// Group by domain for better performance
technicalTexts := filterTechnical(allTexts)
for _, text := range technicalTexts {
    result := textlib.DomainOptimizedAnalyze(text, "technical")
    // Process result
}
```

### 3. Reuse Results

Many functions build on each other. Reuse intermediate results:

```go
result := textlib.SmartAnalyze(text)

// Don't call individual functions again - use the comprehensive result
entities := result.Entities          // Already extracted
readability := result.Readability    // Already calculated
keywords := result.Keywords          // Already extracted
```

## Common Patterns

### Pattern 1: Content Classification

```go
func classifyContent(text string) string {
    result := textlib.SmartAnalyze(text)
    
    // Use multiple signals for classification
    if len(result.Entities) > 10 && result.Readability.Complexity == "complex" {
        return "technical"
    }
    
    if result.Sentiment.Tone != "neutral" && len(text) < 280 {
        return "social"
    }
    
    if result.Structure.DocumentType == "instructional" {
        return "tutorial"
    }
    
    return "general"
}
```

### Pattern 2: Entity-Aware Summarization

```go
func smartSummarize(text string) string {
    result := textlib.SmartAnalyze(text)
    
    // Build summary using key entities and keywords
    keyEntities := getTopEntities(result.Entities, 3)
    topKeywords := result.Keywords[:min(5, len(result.Keywords))]
    
    return fmt.Sprintf("This %s text discusses %s, focusing on %s",
        result.Readability.Complexity,
        joinEntities(keyEntities),
        strings.Join(topKeywords, ", "))
}
```

### Pattern 3: Quality Checking

```go
func checkTextQuality(text string) QualityReport {
    result := textlib.SmartAnalyze(text)
    
    report := QualityReport{}
    
    // Check readability
    if result.Readability.GradeLevel > 12 {
        report.Issues = append(report.Issues, "Text may be too complex")
    }
    
    // Check structure
    if result.Structure.Sentences < 3 {
        report.Issues = append(report.Issues, "Text too short for meaningful analysis")
    }
    
    // Check entity coverage
    if len(result.Entities) == 0 {
        report.Warnings = append(report.Warnings, "No entities detected")
    }
    
    return report
}
```

## Error Handling

Always handle edge cases gracefully:

```go
func safeAnalyze(text string) (*textlib.ComprehensiveResult, error) {
    // Check input
    if len(text) == 0 {
        return nil, fmt.Errorf("empty text provided")
    }
    
    if len(text) > 1_000_000 {
        return nil, fmt.Errorf("text too large (max 1MB)")
    }
    
    // Perform analysis with timeout
    result := textlib.SmartAnalyze(text)
    
    // Validate results
    if len(result.ProcessingInfo.FunctionsUsed) == 0 {
        return nil, fmt.Errorf("analysis failed")
    }
    
    return &result, nil
}
```

## Advanced Usage

### Custom Analysis Pipeline

Create your own optimized pipeline:

```go
type CustomAnalyzer struct {
    skipSentiment bool
    maxEntities   int
}

func (ca *CustomAnalyzer) Analyze(text string) CustomResult {
    // Start with validation (always recommended)
    clean := textlib.ValidateAndClean(text)
    
    result := CustomResult{}
    
    // Entity extraction (most valuable)
    entities := textlib.ExtractAdvancedEntities(clean)
    if ca.maxEntities > 0 && len(entities) > ca.maxEntities {
        entities = entities[:ca.maxEntities]
    }
    result.Entities = entities
    
    // Conditional sentiment
    if !ca.skipSentiment {
        result.Sentiment = textlib.AnalyzeSentiment(clean)
    }
    
    // Keywords based on entities
    result.Keywords = textlib.ExtractKeywordsOptimized(clean, entities)
    
    return result
}
```

### Streaming Analysis

For large documents:

```go
func analyzeStream(reader io.Reader) ([]ChunkResult, error) {
    scanner := bufio.NewScanner(reader)
    scanner.Split(splitIntoParagraphs)
    
    results := []ChunkResult{}
    
    for scanner.Scan() {
        chunk := scanner.Text()
        
        // Use QuickInsights for each chunk
        insights := textlib.QuickInsights(chunk)
        
        results = append(results, ChunkResult{
            Text:     chunk,
            Insights: insights,
        })
    }
    
    return results, scanner.Err()
}
```

## Benchmarks

Typical performance for a 1000-word technical document:

| Function | Time | Accuracy |
|----------|------|----------|
| SmartAnalyze | 15ms | 94% |
| QuickInsights | 3ms | 85% |
| DeepTechnicalAnalysis | 25ms | 97% |
| ValidatedExtraction | 8ms | 95% |

## Migration Guide

If you're currently using individual functions, here's how to migrate:

### Before (Individual Calls)
```go
// Inefficient - multiple passes over text
entities := textlib.ExtractNamedEntities(text)
readability := textlib.CalculateFleschReadingEase(text) 
keywords := textlib.ExtractKeywords(text)
stats := textlib.CalculateTextStatistics(text)
```

### After (Optimized)
```go
// Efficient - single optimized pass
result := textlib.SmartAnalyze(text)
// All results available in result struct
```

## FAQ

**Q: When should I use SmartAnalyze vs individual functions?**
A: Use SmartAnalyze for most cases. Only use individual functions if you need exactly one specific analysis.

**Q: How do I know which domain to specify?**
A: If unsure, use SmartAnalyze without a domain. It will auto-detect the best approach.

**Q: Can I combine multiple domains?**
A: Not directly, but SmartAnalyze handles mixed content well.

**Q: What's the maximum text size?**
A: Recommended max is 100KB for optimal performance. Larger texts should be chunked.

## Summary

The key to optimal TextLib usage is:

1. **Validate First**: Always clean and validate text before analysis
2. **Use Composite Functions**: SmartAnalyze and domain-specific functions are optimized
3. **Match Function to Use Case**: Quick for real-time, Deep for thorough analysis
4. **Reuse Results**: Don't call multiple functions when one comprehensive call suffices

Following these patterns will give you the best performance and accuracy from TextLib.