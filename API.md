# TextLib API Documentation

This document provides detailed API reference for all functions and types in TextLib.

## Table of Contents

- [Core Types](#core-types)
- [Text Processing](#text-processing)
- [Entity Recognition](#entity-recognition)
- [Grammar Analysis](#grammar-analysis)
- [Text Statistics](#text-statistics)
- [Code Analysis](#code-analysis)
- [Mathematical Analysis](#mathematical-analysis)
- [File Processing](#file-processing)
- [Security Analysis](#security-analysis)

## Core Types

### Position
```go
type Position struct {
    Start int
    End   int
}
```
Represents a position range in text.

### Entity
```go
type Entity struct {
    Text     string
    Type     string
    Position Position
}
```
Represents a named entity found in text.

### Issue
```go
type Issue struct {
    Type        string
    Description string
    Position    Position
    Severity    string
}
```
Represents a grammar or code issue.

## Text Processing

### SplitIntoSentences
```go
func SplitIntoSentences(text string) []string
```
Splits text into sentences, handling common abbreviations like "Dr.", "U.S.", etc.

**Parameters:**
- `text`: Input text to split

**Returns:**
- `[]string`: Array of sentences

**Example:**
```go
sentences := textlib.SplitIntoSentences("Dr. Smith went to the U.S. He met Mr. Johnson.")
// Returns: ["Dr. Smith went to the U.S.", "He met Mr. Johnson."]
```

### SplitIntoParagraphs
```go
func SplitIntoParagraphs(text string) []string
```
Splits text into paragraphs based on double line breaks.

**Parameters:**
- `text`: Input text to split

**Returns:**
- `[]string`: Array of paragraphs

### ExtractNamedEntities
```go
func ExtractNamedEntities(text string) []Entity
```
Extracts basic named entities (PERSON, ORGANIZATION, LOCATION, etc.).

**Parameters:**
- `text`: Input text to analyze

**Returns:**
- `[]Entity`: Array of detected entities

**Entity Types:**
- `PERSON`: Names of people
- `ORGANIZATION`: Company and organization names
- `LOCATION`: Places and geographical locations
- `DATE`: Date expressions
- `TIME`: Time expressions
- `MONEY`: Monetary amounts
- `PERCENT`: Percentage values

## Entity Recognition

### ExtractAdvancedEntities
```go
func ExtractAdvancedEntities(text string) []Entity
```
Extracts comprehensive entity types including emails, phones, IPs, etc.

**Additional Entity Types:**
- `EMAIL`: Email addresses
- `PHONE`: Phone numbers
- `URL`: Web URLs
- `IP`: IP addresses
- `CREDIT_CARD`: Credit card numbers
- `SSN`: Social Security numbers
- `CODE`: Code snippets

## Grammar Analysis

### CheckGrammar
```go
func CheckGrammar(text string) []Issue
```
Performs comprehensive grammar checking.

**Checks Include:**
- Sentence fragments
- Run-on sentences
- Subject-verb agreement
- Capitalization errors
- Punctuation balance

### IsCompleteSentence
```go
func IsCompleteSentence(text string) bool
```
Determines if text forms a complete sentence.

### CheckParenthesesBalance
```go
func CheckParenthesesBalance(text string) bool
```
Checks if parentheses are properly balanced.

### CheckQuoteBalance
```go
func CheckQuoteBalance(text string) bool
```
Checks if quotes are properly balanced.

## Text Statistics

### TextStatistics
```go
type TextStatistics struct {
    WordCount         int
    UniqueWords       int
    CharCount         int
    SentenceCount     int
    ParagraphCount    int
    TypeTokenRatio    float64
    AverageWordLength float64
    HapaxLegomena     int
    Vocabulary        map[string]int
}
```

### CalculateTextStatistics
```go
func CalculateTextStatistics(text string) TextStatistics
```
Calculates comprehensive text statistics.

### CalculateFleschScore
```go
func CalculateFleschScore(text string) float64
```
Calculates Flesch readability score (0-100, higher = more readable).

### CalculateSyllableCount
```go
func CalculateSyllableCount(text string) int
```
Counts total syllables in text.

### CalculateTypeTokenRatio
```go
func CalculateTypeTokenRatio(text string) float64
```
Calculates vocabulary diversity (unique words / total words).

## Code Analysis

### FunctionSig
```go
type FunctionSig struct {
    Name       string
    Parameters []Parameter
    ReturnType string
    Visibility string
    Position   Position
    Language   string
}
```

### ExtractFunctionSignatures
```go
func ExtractFunctionSignatures(code string) []FunctionSig
```
Extracts function signatures from code (supports multiple languages).

**Supported Languages:**
- Go
- JavaScript
- Python
- Java
- C/C++

### CalculateCyclomaticComplexity
```go
func CalculateCyclomaticComplexity(code string) int
```
Calculates cyclomatic complexity of code.

### CountLines
```go
func CountLines(code string) int
```
Counts total lines in code.

### CountBlankLines
```go
func CountBlankLines(code string) int
```
Counts blank lines in code.

### CountCommentLines
```go
func CountCommentLines(code string) int
```
Counts comment lines in code.

### Naming Convention Validation

#### CheckCamelCase
```go
func CheckCamelCase(name string) bool
```
Validates camelCase naming convention.

#### CheckPascalCase
```go
func CheckPascalCase(name string) bool
```
Validates PascalCase naming convention.

#### CheckSnakeCase
```go
func CheckSnakeCase(name string) bool
```
Validates snake_case naming convention.

#### CheckKebabCase
```go
func CheckKebabCase(name string) bool
```
Validates kebab-case naming convention.

## Mathematical Analysis

### MathExpression
```go
type MathExpression struct {
    Expression string
    Type       string
    Variables  []string
    Position   Position
}
```

### ValidationResult
```go
type ValidationResult struct {
    IsValid bool
    Errors  []string
    Type    string
}
```

### ValidateMathExpression
```go
func ValidateMathExpression(expression string) ValidationResult
```
Validates mathematical expressions and equations.

### ExtractMathExpressions
```go
func ExtractMathExpressions(text string) []MathExpression
```
Extracts mathematical expressions from text.

### DetectMathPatterns
```go
func DetectMathPatterns(text string) []MathPattern
```
Detects mathematical patterns (sequences, progressions, etc.).

## File Processing

### FileMetadata
```go
type FileMetadata struct {
    Name        string
    Size        int64
    Extension   string
    MimeType    string
    Created     time.Time
    Modified    time.Time
    Permissions string
    IsDirectory bool
    Path        string
}
```

### ExtractMetadata
```go
func ExtractMetadata(filePath string) (FileMetadata, error)
```
Extracts comprehensive file metadata.

### DetectFileType
```go
func DetectFileType(content []byte) FileType
```
Detects file type from content (magic bytes).

### CalculateChecksum
```go
func CalculateChecksum(filePath string, algorithm string) (string, error)
```
Calculates file checksum using specified algorithm.

**Supported Algorithms:**
- `md5`
- `sha1`
- `sha256`
- `sha512`

### Content Processing

#### ExtractTextFromPDF
```go
func ExtractTextFromPDF(filePath string) (string, error)
```
Extracts text content from PDF files.

#### ParseCSVStructure
```go
func ParseCSVStructure(filePath string) (CSVSchema, error)
```
Analyzes CSV file structure and schema.

#### AnalyzeImageProperties
```go
func AnalyzeImageProperties(imagePath string) (ImageMetadata, error)
```
Extracts image metadata and properties.

#### ExtractAudioMetadata
```go
func ExtractAudioMetadata(audioPath string) (AudioMetadata, error)
```
Extracts audio file metadata.

## Security Analysis

### SecurityThreat
```go
type SecurityThreat struct {
    Type        string
    Severity    string
    Description string
    Position    int64
    Pattern     string
    Confidence  float64
    Mitigation  string
}
```

### DetectMaliciousPatterns
```go
func DetectMaliciousPatterns(filePath string) ([]SecurityThreat, error)
```
Scans files for malicious patterns and security threats.

### AnalyzeExecutableHeaders
```go
func AnalyzeExecutableHeaders(exePath string) (ExecutableInfo, error)
```
Analyzes executable file headers and structure.

### ScanForViruses
```go
func ScanForViruses(filePath string, signatures []VirusSignature) ([]Threat, error)
```
Scans files using provided virus signatures.

### CheckFilePermissions
```go
func CheckFilePermissions(filePath string) (PermissionReport, error)
```
Analyzes file permissions for security issues.

## Error Handling

All functions that can fail return an error as the last return value. Always check errors:

```go
result, err := textlib.SomeFunction(input)
if err != nil {
    log.Printf("Error: %v", err)
    return
}
// Use result safely
```

## Performance Notes

- Functions are optimized for performance with large texts
- Memory usage is minimized through efficient algorithms
- Some functions support parallel processing for large datasets
- Use appropriate data structures for your specific use case

## Version Compatibility

TextLib follows semantic versioning (SemVer):
- Patch versions (x.x.1): Bug fixes, backward compatible
- Minor versions (x.1.x): New features, backward compatible  
- Major versions (1.x.x): Breaking changes, may require code updates

Current version: 1.0.0