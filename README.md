# TextLib - Advanced Text Analysis Library for Go

[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![Go Reference](https://pkg.go.dev/badge/github.com/caiatech/textlib.svg)](https://pkg.go.dev/github.com/caiatech/textlib)
[![Go Report Card](https://goreportcard.com/badge/github.com/caiatech/textlib)](https://goreportcard.com/report/github.com/caiatech/textlib)

TextLib is a comprehensive Go library for advanced text analysis, processing, and machine learning applications. Built by [Caia Tech](https://caiatech.com), it provides deterministic text analysis functions essential for reinforcement learning and agentic AI systems.

## Features

### 🔤 Core Text Processing
- **Sentence Segmentation**: Smart sentence splitting with abbreviation handling
- **Entity Recognition**: Extract persons, organizations, locations, dates, money, percentages
- **Grammar Analysis**: Detect fragments, run-on sentences, subject-verb agreement
- **Readability Metrics**: Flesch scores, syllable counting, complexity analysis

### 📊 Advanced Text Analytics
- **Statistical Analysis**: Vocabulary diversity, type-token ratio, hapax legomena
- **Pattern Detection**: Repetitions, syntactic patterns, rhetorical devices
- **Text Similarity**: Levenshtein distance, Jaccard index, cosine similarity, TF-IDF
- **Dialogue Analysis**: Speaker identification, turn-taking, conversation flows
- **Text Segmentation**: Token-based, sentence-based, semantic, sliding window chunking

### 💻 Code Analysis
- **Function Extraction**: Multi-language function signature detection
- **Complexity Metrics**: Cyclomatic complexity, nesting depth analysis
- **Naming Conventions**: camelCase, PascalCase, snake_case, kebab-case validation
- **Security Scanning**: Hardcoded secrets, SQL injection patterns
- **Code Quality**: Duplicate detection, formatting analysis

### 🧮 Mathematical Expression Analysis
- **Expression Parsing**: Algebraic equations and mathematical formulas
- **Pattern Recognition**: Arithmetic sequences, geometric progressions
- **Validation**: Parentheses matching, operator precedence
- **Simplification**: Basic algebraic manipulation

### 📁 File Processing
- **Structure Analysis**: File metadata, type detection, integrity validation
- **Content Processing**: PDF text extraction, CSV parsing, image/audio metadata
- **Security Analysis**: Malicious pattern detection, executable analysis
- **Organization**: File categorization, duplicate detection, cleanup utilities

## Installation

```bash
go get github.com/caiatech/textlib
```

## Quick Start

```go
package main

import (
    "fmt"
    "github.com/caiatech/textlib"
)

func main() {
    text := "Dr. Smith went to the U.S. yesterday. He met Mr. Johnson."
    
    // Split into sentences
    sentences := textlib.SplitIntoSentences(text)
    fmt.Printf("Sentences: %v\n", sentences)
    
    // Extract named entities
    entities := textlib.ExtractNamedEntities(text)
    for _, entity := range entities {
        fmt.Printf("Entity: %s (%s)\n", entity.Text, entity.Type)
    }
    
    // Calculate readability
    score := textlib.CalculateFleschScore(text)
    fmt.Printf("Flesch Score: %.2f\n", score)
    
    // Analyze text statistics
    stats := textlib.CalculateTextStatistics(text)
    fmt.Printf("Words: %d, Unique: %d, TTR: %.2f\n", 
        stats.WordCount, stats.UniqueWords, stats.TypeTokenRatio)
}
```

## Advanced Usage

### Entity Recognition with Advanced Types

```go
// Extract comprehensive entities
entities := textlib.ExtractAdvancedEntities(text)
for _, entity := range entities {
    switch entity.Type {
    case "MONEY":
        fmt.Printf("Found money: %s\n", entity.Text)
    case "PHONE":
        fmt.Printf("Found phone: %s\n", entity.Text)
    case "EMAIL":
        fmt.Printf("Found email: %s\n", entity.Text)
    }
}
```

### Code Analysis

```go
code := `
func calculateSum(a, b int) int {
    return a + b
}
`

// Extract functions
functions := textlib.ExtractFunctionSignatures(code)
for _, fn := range functions {
    fmt.Printf("Function: %s (%s)\n", fn.Name, fn.Language)
}

// Calculate complexity
complexity := textlib.CalculateCyclomaticComplexity(code)
fmt.Printf("Cyclomatic Complexity: %d\n", complexity)

// Validate naming conventions
rules := textlib.NamingRules{RequireCamelCase: true}
issues := textlib.ValidateFunctionNames(code, rules)
```

### Mathematical Expression Analysis

```go
expression := "2x + 3 = 7"
result := textlib.ValidateMathExpression(expression)
if result.IsValid {
    fmt.Printf("Valid math expression: %s\n", expression)
}

// Extract mathematical patterns
patterns := textlib.DetectMathPatterns("The sequence 2, 4, 6, 8 follows a pattern")
for _, pattern := range patterns {
    fmt.Printf("Pattern: %s (%s)\n", pattern.Type, pattern.Description)
}
```

### File Processing

```go
// Analyze file structure
report, err := textlib.AnalyzeFileStructure("/path/to/directory")
if err == nil {
    fmt.Printf("Total files: %d, Total size: %d bytes\n", 
        report.TotalFiles, report.TotalSize)
}

// Extract PDF text
text, err := textlib.ExtractTextFromPDF("/path/to/document.pdf")
if err == nil {
    fmt.Printf("Extracted text: %s\n", text)
}

// Security analysis
threats, err := textlib.DetectMaliciousPatterns("/path/to/suspicious/file")
for _, threat := range threats {
    fmt.Printf("Threat: %s (severity: %s)\n", threat.Description, threat.Severity)
}
```

## API Reference

### Core Functions

| Function | Description |
|----------|-------------|
| `SplitIntoSentences(text string) []string` | Split text into sentences with abbreviation handling |
| `SplitIntoParagraphs(text string) []string` | Split text into paragraphs |
| `ExtractNamedEntities(text string) []Entity` | Extract basic named entities |
| `ExtractAdvancedEntities(text string) []Entity` | Extract comprehensive entity types |
| `CalculateFleschScore(text string) float64` | Calculate Flesch readability score |
| `CalculateTextStatistics(text string) TextStatistics` | Comprehensive text statistics |

### Analysis Functions

| Function | Description |
|----------|-------------|
| `DetectPatterns(text string) []Pattern` | Find repetitive patterns and structures |
| `CalculateCosineSimilarity(text1, text2 string) float64` | Text similarity using TF-IDF |
| `AnalyzeDialogue(text string) DialogueAnalysis` | Speaker and conversation analysis |
| `ChunkText(text string, strategy ChunkStrategy) []TextChunk` | Intelligent text segmentation |

### Code Analysis

| Function | Description |
|----------|-------------|
| `ExtractFunctionSignatures(code string) []FunctionSig` | Multi-language function detection |
| `CalculateCyclomaticComplexity(code string) int` | Code complexity metrics |
| `DetectDuplicateCode(code string, threshold int) []CodeBlock` | Find code duplication |
| `FindHardcodedSecrets(code string) []SecurityIssue` | Security vulnerability detection |

## Performance

TextLib is designed for high performance with:
- **Zero-copy operations** where possible
- **Parallel processing** for large text datasets
- **Memory-efficient algorithms** for statistical analysis
- **42.9% test coverage** with comprehensive test suites

## Use Cases

- **AI/ML Pipelines**: Text preprocessing for language models
- **Content Analysis**: Document classification and sentiment analysis
- **Code Quality Tools**: Static analysis and refactoring assistance
- **Security Scanning**: Malware detection and vulnerability assessment
- **Educational Software**: Grammar checking and readability analysis
- **Research Applications**: Linguistic analysis and corpus processing

## Contributing

We welcome contributions! Please see our [Contributing Guidelines](CONTRIBUTING.md) for details.

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## Testing

Run the full test suite:

```bash
go test -v ./...
```

Run with coverage:

```bash
go test -cover ./...
```

## License

This project is licensed under the Apache License 2.0 - see the [LICENSE](LICENSE) file for details.

## About Caia Tech

TextLib is developed and maintained by [Caia Tech](https://caiatech.com), a technology company focused on advanced AI and machine learning solutions.

**CEO**: Marvin Tutt  
**Contact**: [hello@caiatech.com](mailto:hello@caiatech.com)

## Acknowledgments

- Built with Go's standard library for maximum compatibility
- Inspired by leading NLP libraries and best practices
- Designed for the AI/ML community and enterprise applications

---

**⭐ Star this repository if you find it useful!**