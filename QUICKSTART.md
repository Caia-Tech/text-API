# Quick Start Guide for TextLib

## Step 1: Create a New Go Project

```bash
# Create a new directory for your project
mkdir my-text-analyzer
cd my-text-analyzer

# Initialize a Go module
go mod init my-text-analyzer
```

## Step 2: Install TextLib

```bash
go get github.com/Caia-Tech/text-API/textlib
```

## Step 3: Create Your First Program

Create a file named `main.go`:

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/Caia-Tech/text-API/textlib"
)

func main() {
    // Your text to analyze
    text := "Hello! My name is John Smith. I work at Tech Corp in San Francisco. You can reach me at john@techcorp.com or (555) 123-4567."
    
    // Basic text analysis
    fmt.Println("=== TEXT ANALYSIS ===")
    
    // Count words
    wordCount := textlib.CountWords(text)
    fmt.Printf("Word count: %d\n", wordCount)
    
    // Split into sentences
    sentences := textlib.SplitIntoSentences(text)
    fmt.Printf("Number of sentences: %d\n", len(sentences))
    
    // Extract entities
    fmt.Println("\n=== ENTITIES FOUND ===")
    entities := textlib.ExtractAdvancedEntities(text)
    
    for _, entity := range entities {
        fmt.Printf("%s: %s\n", entity.Type, entity.Text)
    }
    
    // Calculate readability
    fmt.Println("\n=== READABILITY ===")
    score := textlib.CalculateFleschReadingEase(text)
    fmt.Printf("Flesch Reading Ease Score: %.2f\n", score)
    
    if score >= 60 {
        fmt.Println("Reading level: Easy")
    } else if score >= 30 {
        fmt.Println("Reading level: Moderate")
    } else {
        fmt.Println("Reading level: Difficult")
    }
}
```

## Step 4: Run Your Program

```bash
go run main.go
```

## Expected Output

```
=== TEXT ANALYSIS ===
Word count: 22
Number of sentences: 3

=== ENTITIES FOUND ===
PERSON: John Smith
ORGANIZATION: Tech Corp
LOCATION: San Francisco
EMAIL: john@techcorp.com
PHONE: (555) 123-4567

=== READABILITY ===
Flesch Reading Ease Score: 65.73
Reading level: Easy
```

## Next Steps

Now that you have TextLib working, try these features:

### 1. Analyze Code

```go
sourceCode := `
func add(a, b int) int {
    return a + b
}
`

functions := textlib.ExtractFunctionSignatures(sourceCode)
for _, fn := range functions {
    fmt.Printf("Found function: %s\n", fn.Name)
}
```

### 2. Check Text Similarity

```go
text1 := "The quick brown fox"
text2 := "The fast brown fox"

similarity := textlib.CalculateJaccardSimilarity(text1, text2)
fmt.Printf("Similarity: %.2f%%\n", similarity * 100)
```

### 3. Process Files

```go
// Get file metadata
metadata, err := textlib.ExtractMetadata("document.txt")
if err == nil {
    fmt.Printf("File size: %d bytes\n", metadata.Size)
    fmt.Printf("Created: %s\n", metadata.CreatedAt)
}
```

## Common Issues

### Import Error
If you get an import error, make sure you're using the correct import path:
```go
import "github.com/Caia-Tech/text-API/textlib"
```

### Module Issues
If you have module issues, try:
```bash
go mod tidy
go mod download
```

## Full Documentation

For complete documentation and advanced examples, see:
- [USAGE.md](textlib/USAGE.md) - Detailed usage guide
- [API.md](textlib/API.md) - Complete API reference
- [examples/](textlib/examples/) - Example programs

## Support

- GitHub Issues: https://github.com/Caia-Tech/text-API/issues
- Email: support@caiatech.com