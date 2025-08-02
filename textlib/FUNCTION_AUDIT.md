# TextLib Function Audit & Consolidation Plan

## Current State: 114 Functions
**Problem**: We have too many functions, creating API complexity and user confusion.

## Function Breakdown by Category

### 🔥 Major Issues

#### 1. Code Analysis (40 functions) - **CRITICAL BLOAT**
```bash
$ grep "^func [A-Z]" code_analysis.go | wc -l
40
```

**Examples of redundancy:**
- `CheckCamelCase`, `CheckPascalCase`, `CheckSnakeCase`, `CheckKebabCase` → Should be 1 function
- `CountLines`, `CountBlankLines`, `CountCommentLines` → Should be 1 function returning struct
- Multiple security detection functions → Should be consolidated

#### 2. Math Analysis (13 functions) - **MODERATE BLOAT**
- Many specialized math functions that could be consolidated
- Complex pattern detection that overlaps

#### 3. File Processing (17 functions across 2 files) - **MODERATE BLOAT**
- Overlapping file type detection
- Redundant metadata extraction

### ✅ Well-Designed Modules

#### 1. RL-Optimized (5 functions) - **PERFECT**
- Each function serves distinct use case
- Clear value proposition
- Good parameter design

#### 2. Core TextLib (9 functions) - **GOOD**
- Essential text processing functions
- Well-scoped responsibilities

## Consolidation Strategy

### Phase 1: Code Analysis Consolidation (40 → 8 functions)

#### Before (40 functions):
```go
func CheckCamelCase(name string) bool
func CheckPascalCase(name string) bool  
func CheckSnakeCase(name string) bool
func CheckKebabCase(name string) bool
func CountLines(code string) int
func CountBlankLines(code string) int
func CountCommentLines(code string) int
func FindHardcodedPasswords(code string) []SecurityIssue
func DetectSQLInjectionPatterns(code string) []SecurityIssue
// ... 31 more functions
```

#### After (8 functions):
```go
// Core analysis
func AnalyzeCode(code string) CodeAnalysis
func ExtractFunctions(code string) []FunctionSignature
func CalculateComplexity(code string) ComplexityMetrics
func ValidateNaming(code string, rules NamingRules) []NamingIssue

// Security analysis  
func DetectSecurityIssues(code string) []SecurityIssue

// Code quality
func AnalyzeQuality(code string) QualityReport
func DetectDuplicates(code string, threshold int) []DuplicateBlock
func SuggestRefactoring(code string) []Suggestion
```

### Phase 2: Math Analysis Consolidation (13 → 4 functions)

#### After:
```go
func ValidateMathExpression(expr string) ValidationResult
func DetectMathPatterns(text string) []Pattern
func ParseEquations(text string) []Equation  
func AnalyzeMathComplexity(expr string) ComplexityScore
```

### Phase 3: File Processing Consolidation (17 → 6 functions)

#### After:
```go
func AnalyzeFile(filePath string) FileAnalysis
func ExtractMetadata(filePath string) Metadata
func ValidateFileStructure(filePath string) ValidationResult
func DetectFileType(content []byte) FileType
func ProcessArchive(archivePath string) ArchiveAnalysis
func OrganizeFiles(directory string, strategy OrgStrategy) OrganizationResult
```

## New API Structure: 3 Tiers

### Tier 1: Essential Functions (15 functions)
**Target users**: 90% of developers
```go
// RL-Optimized (recommended)
textlib.SmartAnalyze(text string) ComprehensiveResult
textlib.QuickInsights(text string) InsightSummary
textlib.ValidatedExtraction(text string) []Entity
textlib.DomainOptimizedAnalyze(text, domain string) ComprehensiveResult

// Core text processing
textlib.ExtractEntities(text string) []Entity
textlib.CalculateReadability(text string) ReadabilityScore
textlib.AnalyzeGrammar(text string) GrammarAnalysis
textlib.DetectPatterns(text string) PatternAnalysis
textlib.CalculateStatistics(text string) TextStatistics

// Code analysis
textlib.AnalyzeCode(code string) CodeAnalysis

// File processing
textlib.AnalyzeFile(filePath string) FileAnalysis

// Text comparison
textlib.CalculateSimilarity(text1, text2 string) SimilarityResult

// Segmentation
textlib.SplitIntoSentences(text string) []string
textlib.SplitIntoParagraphs(text string) []string
textlib.ChunkText(text string, strategy ChunkStrategy) []TextChunk
```

### Tier 2: Specialized Functions (20 functions)
**Target users**: Domain experts, specific use cases
```go
// Advanced code analysis
textlib.ExtractFunctions(code string) []FunctionSignature
textlib.CalculateComplexity(code string) ComplexityMetrics
textlib.DetectSecurityIssues(code string) []SecurityIssue
textlib.AnalyzeQuality(code string) QualityReport

// Advanced file processing
textlib.ExtractMetadata(filePath string) Metadata
textlib.ProcessArchive(archivePath string) ArchiveAnalysis
textlib.DetectFileType(content []byte) FileType

// Mathematical analysis
textlib.ValidateMathExpression(expr string) ValidationResult
textlib.DetectMathPatterns(text string) []Pattern

// Advanced text analysis
textlib.AnalyzeDialogue(text string) DialogueAnalysis
textlib.AnalyzeCoherence(text string) CoherenceAnalysis
textlib.DeepTechnicalAnalysis(text string) TechnicalResult

// Text transformation
textlib.NormalizeText(text string) string
textlib.ExtractKeywords(text string) []string

// Advanced entity recognition
textlib.ExtractAdvancedEntities(text string) []Entity
textlib.ExtractEntityRelationships(text string) []Relationship

// Content processing
textlib.ExtractTextFromPDF(filePath string) (string, error)
textlib.AnalyzeImageProperties(imagePath string) (ImageMetadata, error)
textlib.ParseCSVStructure(filePath string) (CSVSchema, error)
textlib.ValidateJSONStructure(jsonPath string) ([]ValidationError, error)
```

### Tier 3: Advanced/Experimental (10-15 functions)
**Target users**: Researchers, edge cases
- Experimental features
- Highly specialized functions
- Legacy functions (deprecated)

## Implementation Plan

### Step 1: Create New Consolidated Functions
- Don't break existing API yet
- Implement new consolidated functions alongside old ones
- Comprehensive testing

### Step 2: Migration Guide
- Document migration paths from old → new functions
- Show performance/usability benefits
- Provide automated migration tools where possible

### Step 3: Deprecation
- Mark old functions as deprecated
- Add warnings about upcoming removal
- Provide clear timelines

### Step 4: Removal (v2.0)
- Remove deprecated functions
- Clean API surface
- Update all documentation

## Success Metrics

### Before Consolidation:
- **114 functions** (overwhelming)
- **40 code analysis functions** (confusing)
- **Documentation**: 100+ pages
- **Learning curve**: High

### After Consolidation:
- **45-50 functions total** (manageable)
- **8 code analysis functions** (clear choices)
- **Documentation**: 30-40 pages
- **Learning curve**: Low

### Target: 80/20 Rule
- **15 essential functions** handle **80% of use cases**
- **20 specialized functions** handle **15% of use cases**  
- **10-15 advanced functions** handle **5% of use cases**

## Next Steps

1. Implement new consolidated code analysis functions
2. Create migration mappings
3. Add deprecation warnings
4. Update documentation to promote Tier 1 functions
5. Collect usage analytics to validate consolidation choices

The goal: Make TextLib easier to use while maintaining all capabilities through better-designed, more powerful functions.