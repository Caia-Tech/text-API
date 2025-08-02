# TextLib Usage Guide

## Installation

To use TextLib in your Go project, install it using `go get`:

```bash
go get github.com/Caia-Tech/text-API/textlib
```

## Basic Usage

### 1. Import the Library

```go
import "github.com/Caia-Tech/text-API/textlib"
```

### 2. Simple Example

Create a new Go file (e.g., `main.go`):

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/Caia-Tech/text-API/textlib"
)

func main() {
    // Example text
    text := "Dr. Smith went to New York on January 15, 2025. He paid $1,500 for the conference."
    
    // Split into sentences
    sentences := textlib.SplitIntoSentences(text)
    fmt.Printf("Found %d sentences:\n", len(sentences))
    for i, s := range sentences {
        fmt.Printf("%d: %s\n", i+1, s)
    }
    
    // Extract entities
    entities := textlib.ExtractNamedEntities(text)
    fmt.Printf("\nFound %d entities:\n", len(entities))
    for _, e := range entities {
        fmt.Printf("- %s (%s)\n", e.Text, e.Type)
    }
    
    // Calculate readability
    score := textlib.CalculateFleschReadingEase(text)
    fmt.Printf("\nReadability score: %.2f\n", score)
}
```

## Common Use Cases

### Text Analysis

```go
// Get text statistics
stats := textlib.CalculateTextStatistics(text)
fmt.Printf("Words: %d, Unique: %d, Avg length: %.2f\n", 
    stats.WordCount, stats.UniqueWords, stats.AverageWordLength)

// Check sentence completeness
if textlib.IsCompleteSentence(sentence) {
    fmt.Println("This is a complete sentence")
}

// Count syllables
syllables := textlib.CountSyllables(text)
fmt.Printf("Total syllables: %d\n", syllables)
```

### Advanced Entity Recognition

```go
// Extract emails, phones, URLs, dates, money, etc.
entities := textlib.ExtractAdvancedEntities(text)

// Group by type
byType := make(map[string][]string)
for _, e := range entities {
    byType[e.Type] = append(byType[e.Type], e.Text)
}

for entityType, items := range byType {
    fmt.Printf("%s: %v\n", entityType, items)
}
```

### Code Analysis

```go
sourceCode := `
func calculateSum(a, b int) int {
    if a > 0 && b > 0 {
        return a + b
    }
    return 0
}
`

// Extract function signatures
functions := textlib.ExtractFunctionSignatures(sourceCode)
for _, fn := range functions {
    fmt.Printf("Function: %s (language: %s)\n", fn.Name, fn.Language)
}

// Calculate complexity
complexity := textlib.CalculateCyclomaticComplexity(sourceCode)
fmt.Printf("Cyclomatic complexity: %d\n", complexity)

// Check naming conventions
if textlib.CheckCamelCase("myVariable") {
    fmt.Println("Variable follows camelCase convention")
}
```

### File Processing

```go
// Get file metadata
metadata, err := textlib.ExtractMetadata("/path/to/file.txt")
if err != nil {
    log.Fatal(err)
}
fmt.Printf("File size: %d bytes\n", metadata.Size)

// Detect file type
content, _ := os.ReadFile("document.pdf")
fileType := textlib.DetectFileType(content)
fmt.Printf("File type: %s\n", fileType.MimeType)

// Calculate checksum
checksum, err := textlib.CalculateChecksum("/path/to/file", "sha256")
if err != nil {
    log.Fatal(err)
}
fmt.Printf("SHA256: %s\n", checksum)
```

### Text Similarity

```go
text1 := "The quick brown fox jumps over the lazy dog"
text2 := "A fast brown fox leaps over a sleepy dog"

// Calculate Jaccard similarity
jaccard := textlib.CalculateJaccardSimilarity(text1, text2)
fmt.Printf("Jaccard similarity: %.2f\n", jaccard)

// Calculate Levenshtein distance
distance := textlib.CalculateLevenshteinDistance(text1, text2)
fmt.Printf("Levenshtein distance: %d\n", distance)

// Detect repeated patterns
patterns := textlib.DetectRepeatedPatterns(text, 3)
for _, p := range patterns {
    fmt.Printf("Pattern '%s' appears %d times\n", p.Pattern, p.Count)
}
```

### Mathematical Expression Analysis

```go
expression := "2x + 3 = 7"

// Validate expression
result := textlib.ValidateMathExpression(expression)
if result.IsValid {
    fmt.Println("Valid mathematical expression")
}

// Detect patterns
patterns := textlib.DetectMathPatterns("2, 4, 6, 8, 10")
for _, p := range patterns {
    fmt.Printf("Found pattern: %s\n", p.Type)
}
```

## Complete Program Example

Here's a complete program that demonstrates multiple features:

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/Caia-Tech/text-API/textlib"
)

func main() {
    // Sample text for analysis
    article := `
    Dr. Jane Smith, CEO of Tech Corp, announced a groundbreaking discovery 
    on March 15, 2025. The new algorithm improves processing speed by 85% 
    and reduces costs to $500 per unit. Contact her at jane@techcorp.com 
    or call (555) 123-4567 for more information.
    
    The implementation uses advanced machine learning techniques. Initial 
    tests show promising results with 99.5% accuracy. The product will 
    launch at 2:30 PM EST next Monday.
    `
    
    // 1. Basic text analysis
    fmt.Println("=== BASIC TEXT ANALYSIS ===")
    sentences := textlib.SplitIntoSentences(article)
    fmt.Printf("Sentences: %d\n", len(sentences))
    
    stats := textlib.CalculateTextStatistics(article)
    fmt.Printf("Words: %d\n", stats.WordCount)
    fmt.Printf("Unique words: %d\n", stats.UniqueWords)
    fmt.Printf("Average word length: %.2f\n", stats.AverageWordLength)
    
    // 2. Entity extraction
    fmt.Println("\n=== ENTITY EXTRACTION ===")
    entities := textlib.ExtractAdvancedEntities(article)
    
    entityMap := make(map[string][]string)
    for _, e := range entities {
        entityMap[e.Type] = append(entityMap[e.Type], e.Text)
    }
    
    for entityType, items := range entityMap {
        fmt.Printf("%s: %v\n", entityType, items)
    }
    
    // 3. Readability
    fmt.Println("\n=== READABILITY ===")
    readability := textlib.CalculateFleschReadingEase(article)
    fmt.Printf("Flesch Reading Ease: %.2f\n", readability)
    
    if readability >= 60 {
        fmt.Println("Easy to read")
    } else if readability >= 30 {
        fmt.Println("Moderately difficult")
    } else {
        fmt.Println("Difficult to read")
    }
    
    // 4. Grammar check
    fmt.Println("\n=== GRAMMAR CHECK ===")
    for i, sentence := range sentences {
        if textlib.IsCompleteSentence(sentence) {
            fmt.Printf("Sentence %d: ✓ Complete\n", i+1)
        } else {
            fmt.Printf("Sentence %d: ✗ Fragment\n", i+1)
        }
    }
}
```

## Error Handling

Always handle errors when using file operations:

```go
// File processing with error handling
metadata, err := textlib.ExtractMetadata("/path/to/file")
if err != nil {
    log.Printf("Error extracting metadata: %v", err)
    return
}

// PDF text extraction
text, err := textlib.ExtractTextFromPDF("/path/to/document.pdf")
if err != nil {
    log.Printf("Error extracting PDF text: %v", err)
    return
}
```

## Performance Considerations

1. **Large Files**: When processing large files, use streaming where possible
2. **Batch Processing**: Process multiple texts in parallel using goroutines
3. **Memory Usage**: The library is designed to be memory-efficient, but monitor usage for very large texts

## Build and Run

1. Create your Go module:
```bash
go mod init myproject
```

2. Get the library:
```bash
go get github.com/Caia-Tech/text-API/textlib
```

3. Build and run:
```bash
go build -o myapp
./myapp
```

## Support

For issues, feature requests, or contributions, visit:
https://github.com/Caia-Tech/text-API

## License

TextLib is licensed under the Apache License 2.0.