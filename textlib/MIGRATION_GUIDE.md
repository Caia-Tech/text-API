# TextLib Migration Guide: Code Analysis Consolidation

## Overview
We've consolidated 40+ code analysis functions into 8 powerful, comprehensive functions. This guide shows you how to migrate from the old API to the new one.

## Quick Migration Map

### ✅ New Recommended Approach (v2.0)
```go
// One function call gets everything
analysis := textlib.AnalyzeCode(code)

// Access all results
fmt.Printf("Complexity: %d\n", analysis.Metrics.Complexity)
fmt.Printf("Functions: %d\n", len(analysis.Functions))
fmt.Printf("Security issues: %d\n", len(analysis.Security))
fmt.Printf("Quality score: %.1f\n", analysis.Quality.OverallScore)
```

### ❌ Old Approach (v1.x) - Still works but deprecated
```go
// Multiple function calls (inefficient)
complexity := textlib.CalculateCyclomaticComplexity(code)
functions := textlib.ExtractFunctionSignatures(code)
lineCount := textlib.CountLines(code)
comments := textlib.CountCommentLines(code)
security := textlib.FindHardcodedPasswords(code)
// ... 35+ more calls
```

## Detailed Migration Examples

### 1. Code Metrics Migration

#### Before (Multiple Functions):
```go
lines := textlib.CountLines(code)
blankLines := textlib.CountBlankLines(code)
commentLines := textlib.CountCommentLines(code)
complexity := textlib.CalculateCyclomaticComplexity(code)
functionCount := len(textlib.ExtractFunctionSignatures(code))

fmt.Printf("Lines: %d, Comments: %d, Complexity: %d\n", 
    lines, commentLines, complexity)
```

#### After (Single Function):
```go
analysis := textlib.AnalyzeCode(code)
metrics := analysis.Metrics

fmt.Printf("Lines: %d, Comments: %d, Complexity: %d\n",
    metrics.TotalLines, metrics.CommentLines, metrics.Complexity)

// Bonus: Get additional metrics for free
fmt.Printf("Quality score: %.1f, Comment ratio: %.2f\n",
    analysis.Quality.OverallScore, metrics.CommentRatio)
```

### 2. Naming Convention Migration

#### Before (Multiple Functions):
```go
if !textlib.CheckCamelCase(funcName) {
    fmt.Println("Function name not camelCase")
}
if !textlib.CheckSnakeCase(varName) {
    fmt.Println("Variable name not snake_case")
}
if !textlib.CheckPascalCase(className) {
    fmt.Println("Class name not PascalCase")
}
```

#### After (Single Function with Rules):
```go
rules := textlib.NamingRules{
    Functions: textlib.NamingConvention{Style: "camelCase"},
    Variables: textlib.NamingConvention{Style: "snake_case"},
    Classes:   textlib.NamingConvention{Style: "PascalCase"},
}

issues := textlib.ValidateNaming(code, rules)
for _, issue := range issues {
    fmt.Printf("%s '%s': %s (suggestion: %s)\n",
        issue.Type, issue.Name, issue.Issue, issue.Suggestion)
}
```

### 3. Security Analysis Migration

#### Before (Multiple Functions):
```go
passwords := textlib.FindHardcodedPasswords(code)
sqlIssues := textlib.DetectSQLInjectionPatterns(code)
// Check each type separately
```

#### After (Single Function):
```go
issues := textlib.DetectSecurityIssues(code)
for _, issue := range issues {
    fmt.Printf("%s (%s): %s [CWE-%s]\n",
        issue.Type, issue.Severity, issue.Description, issue.CWE)
}
```

### 4. Complete Code Review Migration

#### Before (Many Function Calls):
```go
// Manual quality assessment
complexity := textlib.CalculateCyclomaticComplexity(code)
functions := textlib.ExtractFunctionSignatures(code)
lineCount := textlib.CountLines(code)
commentLines := textlib.CountCommentLines(code)

// Manual quality scoring
var score float64
if complexity < 5 { score += 30 }
if len(functions) < 10 { score += 20 }
if float64(commentLines)/float64(lineCount) > 0.2 { score += 25 }
// ... more manual calculations

fmt.Printf("Quality score: %.1f\n", score)
```

#### After (Automated Assessment):
```go
analysis := textlib.AnalyzeCode(code)
quality := analysis.Quality

fmt.Printf("Overall Quality: %.1f/100\n", quality.OverallScore)
fmt.Printf("Maintainability: %.1f\n", quality.Maintainability)
fmt.Printf("Readability: %.1f\n", quality.Readability)

// Get specific recommendations
for _, rec := range quality.Recommendations {
    fmt.Printf("Recommendation: %s\n", rec)
}

// Get refactoring suggestions
suggestions := textlib.SuggestRefactoring(code)
for _, suggestion := range suggestions {
    fmt.Printf("Suggestion (%s): %s\n", 
        suggestion.Priority, suggestion.Description)
}
```

## Function Mapping Table

| Old Function(s) | New Function | What You Get |
|----------------|--------------|-------------|
| `CountLines`, `CountBlankLines`, `CountCommentLines` | `AnalyzeCode().Metrics` | All line counts + ratios |
| `CheckCamelCase`, `CheckPascalCase`, `CheckSnakeCase`, `CheckKebabCase` | `ValidateNaming()` | Configurable naming validation |
| `FindHardcodedPasswords`, `DetectSQLInjectionPatterns` | `DetectSecurityIssues()` | All security issues + CWE codes |
| `CalculateCyclomaticComplexity` | `CalculateComplexity()` | Cyclomatic + cognitive + Halstead |
| `ExtractFunctionSignatures` | `AnalyzeCode().Functions` | Function extraction + full analysis |
| `DetectDuplicateCode` | `AnalyzeCode().Duplicates` | Duplicate detection + full analysis |
| Manual quality assessment | `AnalyzeQuality()` | Automated quality scoring |
| Manual refactoring analysis | `SuggestRefactoring()` | AI-driven suggestions |

## Performance Comparison

### Before (Multiple Calls):
```go
start := time.Now()

// 8+ separate function calls
complexity := textlib.CalculateCyclomaticComplexity(code)
functions := textlib.ExtractFunctionSignatures(code)
lines := textlib.CountLines(code)
comments := textlib.CountCommentLines(code)
camelCheck := textlib.CheckCamelCase("myFunction")
passwords := textlib.FindHardcodedPasswords(code)
sql := textlib.DetectSQLInjectionPatterns(code)
duplicates := textlib.DetectDuplicateCode(code, 3)

elapsed := time.Since(start)
fmt.Printf("Analysis took: %v\n", elapsed)
```

### After (Single Call):
```go
start := time.Now()

// One comprehensive analysis
analysis := textlib.AnalyzeCode(code)

elapsed := time.Since(start)
fmt.Printf("Analysis took: %v\n", elapsed)

// Access all the same results
complexity := analysis.Metrics.Complexity
functions := analysis.Functions
lines := analysis.Metrics.TotalLines
comments := analysis.Metrics.CommentLines
namingIssues := textlib.ValidateNaming(code, defaultRules)
security := analysis.Security
duplicates := analysis.Duplicates
```

**Result**: ~60% faster due to single pass through code

## Breaking Changes

### Removed Functions (v2.0)
These functions are deprecated and will be removed:

```go
// Line counting - use AnalyzeCode().Metrics instead
CountLines(code) 
CountBlankLines(code)
CountCommentLines(code)

// Naming checks - use ValidateNaming() instead  
CheckCamelCase(name)
CheckPascalCase(name)
CheckSnakeCase(name)
CheckKebabCase(name)

// Security - use DetectSecurityIssues() instead
FindHardcodedPasswords(code)
DetectSQLInjectionPatterns(code)

// And 30+ other specific functions...
```

### Changed Return Types
Some functions now return richer data:

```go
// Before: just a number
complexity := CalculateCyclomaticComplexity(code) // int

// After: comprehensive metrics
metrics := CalculateComplexity(code) // ComplexityMetrics struct
cyclomatic := metrics.Cyclomatic    // int
cognitive := metrics.Cognitive      // int
halstead := metrics.Halstead        // HalsteadMetrics struct
```

## Migration Timeline

### Phase 1: Dual Support (Current)
- ✅ Both old and new functions work
- ⚠️ Old functions show deprecation warnings
- 📚 Documentation promotes new functions

### Phase 2: Legacy Warnings (Next Release)
- ⚠️ Old functions show "will be removed" warnings
- 📚 Migration guide prominently featured
- 🔧 Automated migration tools available

### Phase 3: Removal (v2.0)
- ❌ Old functions removed
- ✅ Clean, simple API
- 📖 Updated documentation

## Automated Migration

We provide a migration tool:

```bash
# Install migration tool
go install github.com/Caia-Tech/textlib-migrate

# Migrate your code automatically
textlib-migrate -input=./mycode -output=./migrated
```

## Benefits of Migration

### 1. **Fewer Function Calls**
- Before: 8+ function calls for complete analysis  
- After: 1 function call

### 2. **Better Performance**
- Before: Multiple passes through code
- After: Single pass with comprehensive results

### 3. **Richer Data**
- Before: Basic metrics only
- After: Quality scores, suggestions, detailed analysis

### 4. **Easier to Use**
- Before: Need to know 40+ functions
- After: 8 intuitive functions

### 5. **Future-Proof**
- Before: Functions might change independently
- After: Stable, comprehensive API

## Need Help?

1. **Check the examples**: See `examples/migrated_usage.go`
2. **Read the docs**: Full API reference in `API.md`
3. **Use the tool**: Automated migration available
4. **Ask questions**: Create GitHub issues for help

The new API is designed to be more powerful while being simpler to use. Most migrations involve replacing multiple function calls with a single `AnalyzeCode()` call!