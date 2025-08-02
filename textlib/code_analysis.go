package textlib

import (
	"regexp"
	"strings"
	"unicode"
)

// Code analysis data structures

type FunctionSig struct {
	Name       string
	Parameters []Parameter
	ReturnType string
	Visibility string // public, private, protected
	Position   Position
	Language   string
}

type Parameter struct {
	Name string
	Type string
}

type Import struct {
	Path     string
	Alias    string
	Language string
	Position Position
}

type Class struct {
	Name       string
	SuperClass string
	Methods    []FunctionSig
	Fields     []Field
	Position   Position
	Language   string
}

type Field struct {
	Name       string
	Type       string
	Visibility string
}

type DuplicateBlock struct {
	Code       string
	Locations  []Position
	Lines      int
	Similarity float64
}

type Function struct {
	Name     string
	Lines    int
	Position Position
}

type CodeBlock struct {
	Code        string
	NestingLevel int
	Position    Position
	BlockType   string // if, for, while, try, etc.
}

type IndentationReport struct {
	Style           string // tabs, spaces, mixed
	ConsistentLevel bool
	Issues          []IndentationIssue
	AverageIndent   float64
}

type IndentationIssue struct {
	Line        int
	Expected    int
	Actual      int
	Description string
}

type FileSize struct {
	Lines      int
	Characters int
	Words      int
	Functions  int
	Classes    int
}

type NamingRules struct {
	FunctionStyle string // camelCase, snake_case, PascalCase
	RequireVerb   bool
	MaxLength     int
	MinLength     int
}

type BraceStyle int

const (
	SameLine BraceStyle = iota
	NextLine
	AllmanStyle
)

type IndentStyle struct {
	Type  string // "tabs" or "spaces"
	Size  int    // number of spaces per indent level
}

type SecurityIssue struct {
	Type        string
	Description string
	Position    Position
	Severity    string
	Pattern     string
}

type Vulnerability struct {
	Type        string
	Description string
	Position    Position
	RiskLevel   string
	Mitigation  string
}

type Issue struct {
	Type        string
	Description string
	Position    Position
	Suggestion  string
}

type PermissionIssue struct {
	FileOperation string
	Permission    string
	Position      Position
	Risk          string
}

// Basic parsing functions

// Deprecated: Use AnalyzeCode().Metrics.TotalLines instead for comprehensive analysis
func CountLines(code string) int {
	if code == "" {
		return 0
	}
	// Count only \n characters, don't normalize line endings
	return len(strings.Split(code, "\n"))
}

// Deprecated: Use AnalyzeCode().Metrics.BlankLines instead for comprehensive analysis
func CountBlankLines(code string) int {
	if code == "" {
		return 0
	}
	// Normalize line endings
	code = strings.ReplaceAll(code, "\r\n", "\n")
	code = strings.ReplaceAll(code, "\r", "\n")
	
	lines := strings.Split(code, "\n")
	count := 0
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			count++
		}
	}
	// Adjust if the last element is empty due to trailing newline
	if len(lines) > 0 && lines[len(lines)-1] == "" {
		count--
	}
	return count
}

// Deprecated: Use AnalyzeCode().Metrics.CommentLines instead for comprehensive analysis
func CountCommentLines(code string) int {
	lines := strings.Split(code, "\n")
	count := 0
	
	// Support multiple comment styles
	commentPatterns := []*regexp.Regexp{
		regexp.MustCompile(`^\s*//`),        // C++/Java/JS style
		regexp.MustCompile(`^\s*#`),         // Python/Shell style
		regexp.MustCompile(`^\s*/\*.*\*/`),  // C style single line
		regexp.MustCompile(`^\s*\*`),        // C style multi-line continuation
		regexp.MustCompile(`^\s*--`),        // SQL style
		regexp.MustCompile(`^\s*%`),         // MATLAB style
	}
	
	for _, line := range lines {
		for _, pattern := range commentPatterns {
			if pattern.MatchString(line) {
				count++
				break
			}
		}
	}
	
	return count
}

func ExtractFunctionSignatures(code string) []FunctionSig {
	// Auto-detect language and extract functions for that language
	language := detectProgrammingLanguage(code)
	return ExtractFunctionSignaturesForLanguage(code, language)
}

func ExtractFunctionSignaturesForLanguage(code string, language string) []FunctionSig {
	functions := []FunctionSig{}
	
	// Function patterns for different languages
	patterns := map[string]*regexp.Regexp{
		"javascript": regexp.MustCompile(`function\s+(\w+)\s*\(([^)]*)\)`),
		"python":     regexp.MustCompile(`def\s+(\w+)\s*\(([^)]*)\):`),
		"java":       regexp.MustCompile(`(public|private|protected)?\s*(static)?\s*\w+\s+(\w+)\s*\(([^)]*)\)`),
		"go":         regexp.MustCompile(`func\s+(\w+)\s*\(([^)]*)\)(?:\s*\(([^)]*)\)|\s*(\w+))?`),
		"c":          regexp.MustCompile(`\w+\s+(\w+)\s*\(([^)]*)\)`),
	}
	
	pattern, exists := patterns[language]
	if !exists {
		return functions
	}
	
	lines := strings.Split(code, "\n")
	
	for lineNum, line := range lines {
		matches := pattern.FindStringSubmatch(line)
		if len(matches) > 0 {
			function := FunctionSig{
				Name:     extractFunctionName(matches, language),
				Position: Position{Start: lineNum, End: lineNum},
				Language: language,
			}
			
			// Extract parameters
			if len(matches) > 2 && matches[2] != "" {
				function.Parameters = parseParameters(matches[2], language)
			}
			
			functions = append(functions, function)
		}
	}
	
	return functions
}

func detectProgrammingLanguage(code string) string {
	// Simple language detection based on keywords and syntax
	if strings.Contains(code, "function ") && strings.Contains(code, "{") {
		return "javascript"
	}
	if strings.Contains(code, "def ") && strings.Contains(code, ":") {
		return "python"
	}
	if strings.Contains(code, "func ") && strings.Contains(code, "{") {
		return "go"
	}
	if strings.Contains(code, "public ") || strings.Contains(code, "private ") || strings.Contains(code, "class ") {
		return "java"
	}
	if strings.Contains(code, "#include") || strings.Contains(code, "int main") {
		return "c"
	}
	return "unknown"
}

func extractFunctionName(matches []string, lang string) string {
	switch lang {
	case "java":
		if len(matches) >= 4 {
			return matches[3]
		}
	case "javascript", "python", "go", "c":
		if len(matches) >= 2 {
			return matches[1]
		}
	}
	return ""
}

func parseParameters(paramStr string, lang string) []Parameter {
	parameters := []Parameter{}
	
	if strings.TrimSpace(paramStr) == "" {
		return parameters
	}
	
	parts := strings.Split(paramStr, ",")
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		
		param := Parameter{}
		
		switch lang {
		case "python":
			// Python: param_name or param_name: type
			if strings.Contains(part, ":") {
				typeParts := strings.Split(part, ":")
				param.Name = strings.TrimSpace(typeParts[0])
				if len(typeParts) > 1 {
					param.Type = strings.TrimSpace(typeParts[1])
				}
			} else {
				param.Name = part
			}
		case "java", "c":
			// Java/C: type param_name
			words := strings.Fields(part)
			if len(words) >= 2 {
				param.Type = words[0]
				param.Name = words[1]
			} else if len(words) == 1 {
				param.Name = words[0]
			}
		case "go":
			// Go: param_name type or just param_name
			words := strings.Fields(part)
			if len(words) >= 2 {
				param.Name = words[0]
				param.Type = words[1]
			} else {
				param.Name = words[0]
			}
		default:
			param.Name = part
		}
		
		parameters = append(parameters, param)
	}
	
	return parameters
}

func ExtractImports(code string) []Import {
	imports := []Import{}
	lines := strings.Split(code, "\n")
	
	// Import patterns for different languages
	patterns := map[string]*regexp.Regexp{
		"python":     regexp.MustCompile(`^\s*(?:from\s+(\S+)\s+)?import\s+([^#\n]+)`),
		"javascript": regexp.MustCompile(`^\s*import\s+(?:\{([^}]+)\}|\*\s+as\s+(\w+)|(\w+))\s+from\s+['"]([^'"]+)['"]`),
		"java":       regexp.MustCompile(`^\s*import\s+(?:static\s+)?([^;]+);`),
		"go":         regexp.MustCompile(`^\s*import\s+(?:"([^"]+)"|(\w+)\s+"([^"]+)")`),
		"c":          regexp.MustCompile(`^\s*#include\s+[<"]([^>"]+)[>"]`),
	}
	
	for lineNum, line := range lines {
		for lang, pattern := range patterns {
			matches := pattern.FindStringSubmatch(line)
			if len(matches) > 0 {
				imp := Import{
					Language: lang,
					Position: Position{Start: lineNum, End: lineNum},
				}
				
				switch lang {
				case "python":
					if matches[1] != "" {
						imp.Path = matches[1]
						imp.Alias = strings.TrimSpace(matches[2])
					} else {
						imp.Path = strings.TrimSpace(matches[2])
					}
				case "javascript":
					if len(matches) >= 5 {
						imp.Path = matches[4]
						if matches[2] != "" {
							imp.Alias = matches[2]
						} else if matches[3] != "" {
							imp.Alias = matches[3]
						}
					}
				case "java", "c":
					if len(matches) >= 2 {
						imp.Path = strings.TrimSpace(matches[1])
					}
				case "go":
					if matches[1] != "" {
						imp.Path = matches[1]
					} else if matches[3] != "" {
						imp.Path = matches[3]
						imp.Alias = matches[2]
					}
				}
				
				imports = append(imports, imp)
			}
		}
	}
	
	return imports
}

func ExtractClassDefinitions(code string) []Class {
	classes := []Class{}
	lines := strings.Split(code, "\n")
	
	// Class patterns for different languages
	patterns := map[string]*regexp.Regexp{
		"python": regexp.MustCompile(`^\s*class\s+(\w+)(?:\(([^)]+)\))?:`),
		"java":   regexp.MustCompile(`^\s*(?:public|private|protected)?\s*class\s+(\w+)(?:\s+extends\s+(\w+))?`),
		"javascript": regexp.MustCompile(`^\s*class\s+(\w+)(?:\s+extends\s+(\w+))?`),
		"c++":    regexp.MustCompile(`^\s*class\s+(\w+)(?:\s*:\s*(?:public|private|protected)\s+(\w+))?`),
	}
	
	for lineNum, line := range lines {
		for lang, pattern := range patterns {
			matches := pattern.FindStringSubmatch(line)
			if len(matches) > 0 {
				class := Class{
					Name:     matches[1],
					Language: lang,
					Position: Position{Start: lineNum, End: lineNum},
					Methods:  []FunctionSig{},
					Fields:   []Field{},
				}
				
				if len(matches) > 2 && matches[2] != "" {
					class.SuperClass = matches[2]
				}
				
				classes = append(classes, class)
			}
		}
	}
	
	return classes
}

// Code organization functions

func FindUnusedImports(code string, usage []string) []string {
	imports := ExtractImports(code)
	unused := []string{}
	
	// Create usage map for faster lookup
	usageMap := make(map[string]bool)
	for _, u := range usage {
		usageMap[u] = true
	}
	
	for _, imp := range imports {
		// Extract the module/package name from path
		moduleName := getModuleName(imp.Path)
		
		// Check if module or alias is used
		if imp.Alias != "" {
			if !usageMap[imp.Alias] {
				unused = append(unused, imp.Path)
			}
		} else if !usageMap[moduleName] {
			unused = append(unused, imp.Path)
		}
	}
	
	return unused
}

func getModuleName(path string) string {
	// Extract module name from import path
	parts := strings.Split(path, "/")
	if len(parts) > 0 {
		return parts[len(parts)-1]
	}
	
	parts = strings.Split(path, ".")
	if len(parts) > 0 {
		return parts[len(parts)-1]
	}
	
	return path
}

func DetectDuplicateCode(code string, threshold int) []DuplicateBlock {
	duplicates := []DuplicateBlock{}
	lines := strings.Split(code, "\n")
	
	// Create sliding window of lines to check for duplicates
	for i := 0; i < len(lines)-threshold; i++ {
		block := strings.Join(lines[i:i+threshold], "\n")
		normalizedBlock := normalizeCodeBlock(block)
		
		// Skip empty or trivial blocks
		if strings.TrimSpace(normalizedBlock) == "" || isTrivialBlock(normalizedBlock) {
			continue
		}
		
		locations := []Position{{Start: i, End: i + threshold}}
		
		// Look for duplicates of this block
		for j := i + threshold; j < len(lines)-threshold; j++ {
			compareBlock := strings.Join(lines[j:j+threshold], "\n")
			normalizedCompare := normalizeCodeBlock(compareBlock)
			
			similarity := calculateCodeSimilarity(normalizedBlock, normalizedCompare)
			if similarity > 0.8 { // 80% similarity threshold
				locations = append(locations, Position{Start: j, End: j + threshold})
			}
		}
		
		if len(locations) > 1 {
			duplicate := DuplicateBlock{
				Code:       block,
				Locations:  locations,
				Lines:      threshold,
				Similarity: calculateAverageSimilarity(locations, lines),
			}
			duplicates = append(duplicates, duplicate)
		}
	}
	
	return deduplicateDuplicateBlocks(duplicates)
}

func normalizeCodeBlock(code string) string {
	// Remove extra whitespace and normalize formatting
	lines := strings.Split(code, "\n")
	normalized := []string{}
	
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" {
			// Remove comments
			if idx := strings.Index(line, "//"); idx != -1 {
				line = line[:idx]
			}
			if idx := strings.Index(line, "#"); idx != -1 {
				line = line[:idx]
			}
			line = strings.TrimSpace(line)
			if line != "" {
				normalized = append(normalized, line)
			}
		}
	}
	
	return strings.Join(normalized, "\n")
}

func isTrivialBlock(code string) bool {
	lines := strings.Split(code, "\n")
	nonEmptyLines := 0
	
	for _, line := range lines {
		if strings.TrimSpace(line) != "" {
			nonEmptyLines++
		}
	}
	
	return nonEmptyLines < 3
}

func calculateCodeSimilarity(code1, code2 string) float64 {
	if code1 == code2 {
		return 1.0
	}
	
	// Use token-based similarity
	tokens1 := tokenizeCode(code1)
	tokens2 := tokenizeCode(code2)
	
	if len(tokens1) == 0 || len(tokens2) == 0 {
		return 0
	}
	
	// Calculate Jaccard similarity
	set1 := make(map[string]bool)
	set2 := make(map[string]bool)
	
	for _, token := range tokens1 {
		set1[token] = true
	}
	
	for _, token := range tokens2 {
		set2[token] = true
	}
	
	intersection := 0
	for token := range set1 {
		if set2[token] {
			intersection++
		}
	}
	
	union := len(set1) + len(set2) - intersection
	if union == 0 {
		return 0
	}
	
	return float64(intersection) / float64(union)
}

func tokenizeCode(code string) []string {
	// Simple tokenization - split on common separators
	tokenRegex := regexp.MustCompile(`[a-zA-Z_]\w*|[0-9]+\.?[0-9]*|[+\-*/=<>!&|]+|[{}()\[\];,.]`)
	return tokenRegex.FindAllString(code, -1)
}

func calculateAverageSimilarity(locations []Position, lines []string) float64 {
	if len(locations) < 2 {
		return 0
	}
	
	totalSimilarity := 0.0
	comparisons := 0
	
	for i := 0; i < len(locations); i++ {
		for j := i + 1; j < len(locations); j++ {
			block1 := strings.Join(lines[locations[i].Start:locations[i].End], "\n")
			block2 := strings.Join(lines[locations[j].Start:locations[j].End], "\n")
			
			similarity := calculateCodeSimilarity(normalizeCodeBlock(block1), normalizeCodeBlock(block2))
			totalSimilarity += similarity
			comparisons++
		}
	}
	
	if comparisons == 0 {
		return 0
	}
	
	return totalSimilarity / float64(comparisons)
}

func deduplicateDuplicateBlocks(duplicates []DuplicateBlock) []DuplicateBlock {
	// Remove overlapping duplicate blocks, keeping the ones with more locations
	unique := []DuplicateBlock{}
	
	for _, dup := range duplicates {
		isOverlapping := false
		for j, existing := range unique {
			if blocksOverlap(dup.Locations, existing.Locations) {
				// Keep the one with more occurrences
				if len(dup.Locations) > len(existing.Locations) {
					unique[j] = dup
				}
				isOverlapping = true
				break
			}
		}
		
		if !isOverlapping {
			unique = append(unique, dup)
		}
	}
	
	return unique
}

func blocksOverlap(locations1, locations2 []Position) bool {
	for _, loc1 := range locations1 {
		for _, loc2 := range locations2 {
			if loc1.Start < loc2.End && loc2.Start < loc1.End {
				return true
			}
		}
	}
	return false
}

func FindLongFunctions(code string, maxLines int) []Function {
	longFunctions := []Function{}
	functions := ExtractFunctionSignatures(code)
	lines := strings.Split(code, "\n")
	
	for _, function := range functions {
		// Find the end of the function
		startLine := function.Position.Start
		endLine := findFunctionEnd(lines, startLine)
		
		functionLines := endLine - startLine + 1
		if functionLines > maxLines {
			longFunc := Function{
				Name:     function.Name,
				Lines:    functionLines,
				Position: Position{Start: startLine, End: endLine},
			}
			longFunctions = append(longFunctions, longFunc)
		}
	}
	
	return longFunctions
}

func findFunctionEnd(lines []string, startLine int) int {
	// Simple heuristic: find matching braces or indentation
	braceCount := 0
	indentLevel := -1
	
	for i := startLine; i < len(lines); i++ {
		line := lines[i]
		
		// Count braces
		for _, char := range line {
			if char == '{' {
				braceCount++
			} else if char == '}' {
				braceCount--
				if braceCount == 0 && i > startLine {
					return i
				}
			}
		}
		
		// For Python-style indentation
		if braceCount == 0 {
			currentIndent := getIndentationLevel(line)
			if indentLevel == -1 && strings.TrimSpace(line) != "" {
				indentLevel = currentIndent
			} else if currentIndent <= indentLevel && strings.TrimSpace(line) != "" && i > startLine {
				return i - 1
			}
		}
	}
	
	return len(lines) - 1
}

func getIndentationLevel(line string) int {
	count := 0
	for _, char := range line {
		if char == ' ' {
			count++
		} else if char == '\t' {
			count += 4 // Treat tab as 4 spaces
		} else {
			break
		}
	}
	return count
}

func FindDeepNesting(code string, maxDepth int) []CodeBlock {
	deepBlocks := []CodeBlock{}
	lines := strings.Split(code, "\n")
	
	currentDepth := 0
	braceStack := []int{}
	
	for lineNum, line := range lines {
		// Count nesting depth
		for _, char := range line {
			if char == '{' || char == '(' || char == '[' {
				currentDepth++
				braceStack = append(braceStack, lineNum)
			} else if char == '}' || char == ')' || char == ']' {
				if currentDepth > 0 {
					currentDepth--
					if len(braceStack) > 0 {
						braceStack = braceStack[:len(braceStack)-1]
					}
				}
			}
		}
		
		// Check for control structures that add nesting
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "if ") || 
		   strings.HasPrefix(trimmed, "for ") ||
		   strings.HasPrefix(trimmed, "while ") ||
		   strings.HasPrefix(trimmed, "try ") ||
		   strings.HasPrefix(trimmed, "switch ") {
			currentDepth++
		}
		
		if currentDepth > maxDepth {
			block := CodeBlock{
				Code:         line,
				NestingLevel: currentDepth,
				Position:     Position{Start: lineNum, End: lineNum},
				BlockType:    detectBlockType(line),
			}
			deepBlocks = append(deepBlocks, block)
		}
	}
	
	return deepBlocks
}

func detectBlockType(line string) string {
	trimmed := strings.TrimSpace(strings.ToLower(line))
	
	if strings.HasPrefix(trimmed, "if") {
		return "if"
	} else if strings.HasPrefix(trimmed, "for") {
		return "for"
	} else if strings.HasPrefix(trimmed, "while") {
		return "while"
	} else if strings.HasPrefix(trimmed, "try") {
		return "try"
	} else if strings.HasPrefix(trimmed, "switch") {
		return "switch"
	} else if strings.Contains(trimmed, "{") {
		return "block"
	}
	
	return "unknown"
}

// Metrics/Measurements functions

func CalculateCyclomaticComplexity(code string) int {
	complexity := 1 // Base complexity
	lines := strings.Split(code, "\n")
	
	// Patterns that increase complexity
	patterns := []*regexp.Regexp{
		regexp.MustCompile(`\b(if|elif|else if)\b`),
		regexp.MustCompile(`\b(for|while|do)\b`),
		regexp.MustCompile(`\bcase\b`),
		regexp.MustCompile(`\b(catch|except)\b`),
		regexp.MustCompile(`\b(&&|\|\|)\b`),
		regexp.MustCompile(`\?.*:`), // Ternary operator
	}
	
	for _, line := range lines {
		for _, pattern := range patterns {
			matches := pattern.FindAllString(line, -1)
			complexity += len(matches)
		}
	}
	
	return complexity
}

func CountParameters(functionSig string) int {
	// Extract parameter list from function signature
	parenStart := strings.Index(functionSig, "(")
	parenEnd := strings.LastIndex(functionSig, ")")
	
	if parenStart == -1 || parenEnd == -1 || parenEnd <= parenStart {
		return 0
	}
	
	paramStr := strings.TrimSpace(functionSig[parenStart+1 : parenEnd])
	if paramStr == "" {
		return 0
	}
	
	// Count parameters by commas, handling nested parentheses
	paramCount := 1
	parenDepth := 0
	
	for _, char := range paramStr {
		switch char {
		case '(':
			parenDepth++
		case ')':
			parenDepth--
		case ',':
			if parenDepth == 0 {
				paramCount++
			}
		}
	}
	
	return paramCount
}

func CalculateFileSize(code string) FileSize {
	lines := strings.Split(code, "\n")
	words := strings.Fields(code)
	functions := ExtractFunctionSignatures(code)
	classes := ExtractClassDefinitions(code)
	
	return FileSize{
		Lines:      len(lines),
		Characters: len(code),
		Words:      len(words),
		Functions:  len(functions),
		Classes:    len(classes),
	}
}

func CountLoops(code string) int {
	loopPatterns := []*regexp.Regexp{
		regexp.MustCompile(`\b(for|while|do)\b`),
		regexp.MustCompile(`\.forEach\b`),
		regexp.MustCompile(`\.map\b`),
		regexp.MustCompile(`\.filter\b`),
	}
	
	count := 0
	for _, pattern := range loopPatterns {
		matches := pattern.FindAllString(code, -1)
		count += len(matches)
	}
	
	return count
}

func CountConditionals(code string) int {
	conditionalPatterns := []*regexp.Regexp{
		regexp.MustCompile(`\bif\b`),
		regexp.MustCompile(`\belse\s+if\b`),
		regexp.MustCompile(`\belif\b`),
		regexp.MustCompile(`\bswitch\b`),
		regexp.MustCompile(`\bcase\b`),
		regexp.MustCompile(`\?.*:`), // Ternary
	}
	
	count := 0
	for _, pattern := range conditionalPatterns {
		matches := pattern.FindAllString(code, -1)
		count += len(matches)
	}
	
	return count
}

func CalculateAverageFunctionLength(code string) float64 {
	functions := ExtractFunctionSignatures(code)
	if len(functions) == 0 {
		return 0
	}
	
	lines := strings.Split(code, "\n")
	totalLines := 0
	
	for _, function := range functions {
		startLine := function.Position.Start
		endLine := findFunctionEnd(lines, startLine)
		functionLines := endLine - startLine + 1
		totalLines += functionLines
	}
	
	return float64(totalLines) / float64(len(functions))
}

func CountReturnStatements(code string) int {
	returnPattern := regexp.MustCompile(`\breturn\b`)
	return len(returnPattern.FindAllString(code, -1))
}

// Naming/Convention Checks functions

// Deprecated: Use ValidateNaming() with appropriate NamingRules instead
func CheckCamelCase(name string) bool {
	if name == "" {
		return false
	}
	
	// First character should be lowercase
	if !unicode.IsLower(rune(name[0])) {
		return false
	}
	
	// Should contain at least one uppercase letter (for compound words)
	hasUppercase := false
	for _, char := range name[1:] {
		if unicode.IsUpper(char) {
			hasUppercase = true
			break
		}
	}
	
	// No underscores or hyphens
	if strings.Contains(name, "_") || strings.Contains(name, "-") {
		return false
	}
	
	return hasUppercase || len(name) == 1
}

// Deprecated: Use ValidateNaming() with appropriate NamingRules instead
func CheckSnakeCase(name string) bool {
	if name == "" {
		return false
	}
	
	// Should be all lowercase with underscores
	snakePattern := regexp.MustCompile(`^[a-z][a-z0-9_]*[a-z0-9]$`)
	return snakePattern.MatchString(name) || regexp.MustCompile(`^[a-z]$`).MatchString(name)
}

func CheckPascalCase(name string) bool {
	if name == "" {
		return false
	}
	
	// First character should be uppercase
	if !unicode.IsUpper(rune(name[0])) {
		return false
	}
	
	// No underscores or hyphens
	if strings.Contains(name, "_") || strings.Contains(name, "-") {
		return false
	}
	
	return true
}

func CheckKebabCase(name string) bool {
	if name == "" {
		return false
	}
	
	// Check for double hyphens (not allowed)
	if strings.Contains(name, "--") {
		return false
	}
	
	// Should be all lowercase with single hyphens, not starting or ending with hyphen
	kebabPattern := regexp.MustCompile(`^[a-z][a-z0-9]*(-[a-z0-9]+)*$`)
	return kebabPattern.MatchString(name)
}

func ValidateFunctionNames(code string, rules NamingRules) []Issue {
	issues := []Issue{}
	functions := ExtractFunctionSignatures(code)
	
	for _, function := range functions {
		name := function.Name
		
		// Check length constraints
		if rules.MaxLength > 0 && len(name) > rules.MaxLength {
			issues = append(issues, Issue{
				Type:        "naming",
				Description: "Function name too long: " + name,
				Position:    function.Position,
				Suggestion:  "Consider shortening the function name",
			})
		}
		
		if rules.MinLength > 0 && len(name) < rules.MinLength {
			issues = append(issues, Issue{
				Type:        "naming",
				Description: "Function name too short: " + name,
				Position:    function.Position,
				Suggestion:  "Consider using a more descriptive name",
			})
		}
		
		// Check naming style
		var isValidStyle bool
		switch rules.FunctionStyle {
		case "camelCase":
			isValidStyle = CheckCamelCase(name)
		case "snake_case":
			isValidStyle = CheckSnakeCase(name)
		case "PascalCase":
			isValidStyle = CheckPascalCase(name)
		case "kebab-case":
			isValidStyle = CheckKebabCase(name)
		default:
			isValidStyle = true
		}
		
		if !isValidStyle {
			issues = append(issues, Issue{
				Type:        "naming",
				Description: "Function name doesn't follow " + rules.FunctionStyle + " convention: " + name,
				Position:    function.Position,
				Suggestion:  "Rename to follow " + rules.FunctionStyle + " convention",
			})
		}
		
		// Check if function name should contain a verb
		if rules.RequireVerb && !containsVerb(name) {
			issues = append(issues, Issue{
				Type:        "naming",
				Description: "Function name should contain a verb: " + name,
				Position:    function.Position,
				Suggestion:  "Add a verb to make the function's purpose clear",
			})
		}
	}
	
	return issues
}

func containsVerb(name string) bool {
	commonVerbs := []string{
		"get", "set", "is", "has", "can", "should", "will",
		"create", "make", "build", "generate", "produce",
		"find", "search", "locate", "discover",
		"add", "insert", "append", "push",
		"remove", "delete", "pop", "clear",
		"update", "modify", "change", "edit",
		"validate", "check", "verify", "test",
		"load", "save", "read", "write",
		"send", "receive", "fetch", "post",
		"start", "stop", "run", "execute",
		"calculate", "compute", "process", "handle",
	}
	
	lowerName := strings.ToLower(name)
	for _, verb := range commonVerbs {
		if strings.HasPrefix(lowerName, verb) {
			return true
		}
	}
	
	return false
}

func CheckConstantNaming(code string) []Issue {
	issues := []Issue{}
	lines := strings.Split(code, "\n")
	
	// Patterns for constants in different languages
	constPatterns := []*regexp.Regexp{
		regexp.MustCompile(`\bconst\s+(\w+)\s*=`),        // JavaScript/Go
		regexp.MustCompile(`\bfinal\s+\w+\s+(\w+)\s*=`), // Java
		regexp.MustCompile(`^\s*(\w+)\s*=\s*["'0-9]`),   // Python (simple heuristic)
	}
	
	for lineNum, line := range lines {
		for _, pattern := range constPatterns {
			matches := pattern.FindStringSubmatch(line)
			if len(matches) > 1 {
				constName := matches[1]
				
				// Constants should be UPPER_CASE
				if !isUpperSnakeCase(constName) {
					issues = append(issues, Issue{
						Type:        "naming",
						Description: "Constant should be UPPER_CASE: " + constName,
						Position:    Position{Start: lineNum, End: lineNum},
						Suggestion:  "Use UPPER_CASE naming for constants",
					})
				}
			}
		}
	}
	
	return issues
}

func isUpperSnakeCase(name string) bool {
	upperSnakePattern := regexp.MustCompile(`^[A-Z][A-Z0-9_]*$`)
	return upperSnakePattern.MatchString(name)
}

func DetectTypos(code string) []Issue {
	issues := []Issue{}
	lines := strings.Split(code, "\n")
	
	// Common typos in code
	typos := map[string]string{
		"lenght":     "length",
		"recieve":    "receive",
		"seperate":   "separate",
		"definately": "definitely",
		"occured":    "occurred",
		"begining":   "beginning",
		"existance":  "existence",
		"maintainence": "maintenance",
		"accomodate": "accommodate",
		"arguement":  "argument",
	}
	
	for lineNum, line := range lines {
		for typo, correction := range typos {
			if strings.Contains(strings.ToLower(line), typo) {
				issues = append(issues, Issue{
					Type:        "typo",
					Description: "Possible typo: " + typo,
					Position:    Position{Start: lineNum, End: lineNum},
					Suggestion:  "Did you mean: " + correction + "?",
				})
			}
		}
	}
	
	return issues
}

// Formatting/Whitespace functions

func NormalizeWhitespace(code string) string {
	lines := strings.Split(code, "\n")
	normalized := []string{}
	
	for _, line := range lines {
		// Remove trailing whitespace
		line = strings.TrimRight(line, " \t")
		
		// Normalize tabs to spaces (4 spaces per tab)
		line = strings.ReplaceAll(line, "\t", "    ")
		
		normalized = append(normalized, line)
	}
	
	// Remove excessive blank lines (more than 2 consecutive)
	result := []string{}
	blankCount := 0
	
	for _, line := range normalized {
		if strings.TrimSpace(line) == "" {
			blankCount++
			if blankCount <= 2 {
				result = append(result, line)
			}
		} else {
			blankCount = 0
			result = append(result, line)
		}
	}
	
	return strings.Join(result, "\n")
}

func FixIndentation(code string, style IndentStyle) string {
	lines := strings.Split(code, "\n")
	fixed := []string{}
	
	indentLevel := 0
	indentString := getIndentString(style)
	
	for _, line := range lines {
		trimmed := strings.TrimLeft(line, " \t")
		
		// Skip empty lines
		if strings.TrimSpace(line) == "" {
			fixed = append(fixed, "")
			continue
		}
		
		// Adjust indent level based on closing braces
		if strings.HasPrefix(trimmed, "}") || strings.HasPrefix(trimmed, "]") || strings.HasPrefix(trimmed, ")") {
			if indentLevel > 0 {
				indentLevel--
			}
		}
		
		// Apply current indentation
		indented := strings.Repeat(indentString, indentLevel) + trimmed
		fixed = append(fixed, indented)
		
		// Adjust indent level based on opening braces
		if strings.HasSuffix(trimmed, "{") || strings.HasSuffix(trimmed, "[") || strings.HasSuffix(trimmed, "(") {
			indentLevel++
		}
		
		// Handle control structures
		if strings.HasPrefix(trimmed, "if ") || strings.HasPrefix(trimmed, "for ") || 
		   strings.HasPrefix(trimmed, "while ") || strings.HasPrefix(trimmed, "else") {
			if !strings.HasSuffix(trimmed, "{") {
				indentLevel++
			}
		}
	}
	
	return strings.Join(fixed, "\n")
}

func getIndentString(style IndentStyle) string {
	if style.Type == "tabs" {
		return "\t"
	}
	return strings.Repeat(" ", style.Size)
}

func DetectIndentationStyle(code string) IndentationReport {
	lines := strings.Split(code, "\n")
	tabCount := 0
	spaceCount := 0
	mixedCount := 0
	indentLevels := []int{}
	issues := []IndentationIssue{}
	
	for lineNum, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		
		// Count leading whitespace
		leadingSpaces := 0
		leadingTabs := 0
		
		for _, char := range line {
			if char == ' ' {
				leadingSpaces++
			} else if char == '\t' {
				leadingTabs++
			} else {
				break
			}
		}
		
		// Determine style used on this line
		if leadingTabs > 0 && leadingSpaces > 0 {
			mixedCount++
			issues = append(issues, IndentationIssue{
				Line:        lineNum,
				Description: "Mixed tabs and spaces",
			})
		} else if leadingTabs > 0 {
			tabCount++
			indentLevels = append(indentLevels, leadingTabs)
		} else if leadingSpaces > 0 {
			spaceCount++
			indentLevels = append(indentLevels, leadingSpaces)
		}
	}
	
	// Determine predominant style
	var style string
	if tabCount > spaceCount {
		style = "tabs"
	} else if spaceCount > tabCount {
		style = "spaces"
	} else {
		style = "mixed"
	}
	
	// Calculate average indent
	averageIndent := 0.0
	if len(indentLevels) > 0 {
		total := 0
		for _, level := range indentLevels {
			total += level
		}
		averageIndent = float64(total) / float64(len(indentLevels))
	}
	
	return IndentationReport{
		Style:           style,
		ConsistentLevel: mixedCount == 0,
		Issues:          issues,
		AverageIndent:   averageIndent,
	}
}

func RemoveTrailingWhitespace(code string) string {
	lines := strings.Split(code, "\n")
	for i, line := range lines {
		lines[i] = strings.TrimRight(line, " \t")
	}
	return strings.Join(lines, "\n")
}

func ConvertTabsToSpaces(code string, spaceCount int) string {
	spaceString := strings.Repeat(" ", spaceCount)
	return strings.ReplaceAll(code, "\t", spaceString)
}

func ConvertSpacesToTabs(code string, spaceCount int) string {
	spaceString := strings.Repeat(" ", spaceCount)
	return strings.ReplaceAll(code, spaceString, "\t")
}

func EnsureNewlineAtEOF(code string) string {
	if !strings.HasSuffix(code, "\n") {
		return code + "\n"
	}
	return code
}

func DetectBraceStyle(code string) BraceStyle {
	lines := strings.Split(code, "\n")
	sameLineCount := 0
	nextLineCount := 0
	
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		
		// Look for opening braces
		if strings.Contains(trimmed, "{") {
			// Check if brace is at end of line (same line style)
			if strings.HasSuffix(trimmed, "{") {
				sameLineCount++
			}
		}
		
		// Check if this line is just an opening brace (next line style)
		if trimmed == "{" {
			nextLineCount++
		}
	}
	
	if nextLineCount > sameLineCount {
		return NextLine
	}
	return SameLine
}

// Security (Rule-based) functions

func FindHardcodedPasswords(code string) []SecurityIssue {
	issues := []SecurityIssue{}
	lines := strings.Split(code, "\n")
	
	// Patterns that might indicate hardcoded passwords
	patterns := []*regexp.Regexp{
		regexp.MustCompile(`(?i)(password|pwd|pass)\s*[=:]\s*["'][^"']+["']`),
		regexp.MustCompile(`(?i)(secret|key|token)\s*[=:]\s*["'][^"']+["']`),
		regexp.MustCompile(`(?i)api[_-]?key\s*[=:]\s*["'][^"']+["']`),
	}
	
	for lineNum, line := range lines {
		for _, pattern := range patterns {
			if pattern.MatchString(line) {
				issues = append(issues, SecurityIssue{
					Type:        "hardcoded-secret",
					Description: "Possible hardcoded password or secret",
					Position:    Position{Start: lineNum, End: lineNum},
					Severity:    "high",
					Pattern:     line,
				})
			}
		}
	}
	
	return issues
}

func DetectSQLInjectionPatterns(code string) []SecurityIssue {
	issues := []SecurityIssue{}
	lines := strings.Split(code, "\n")
	
	// Patterns that might indicate SQL injection vulnerabilities
	patterns := []*regexp.Regexp{
		regexp.MustCompile(`(?i)select\s+.*\+.*["']`),         // String concatenation in SELECT
		regexp.MustCompile(`(?i)insert\s+.*\+.*["']`),         // String concatenation in INSERT
		regexp.MustCompile(`(?i)update\s+.*\+.*["']`),         // String concatenation in UPDATE
		regexp.MustCompile(`(?i)delete\s+.*\+.*["']`),         // String concatenation in DELETE
		regexp.MustCompile(`(?i)where\s+.*\+.*["']`),          // String concatenation in WHERE
		regexp.MustCompile(`["']\s*\+\s*\w+\s*\+\s*["']`),    // General string concatenation
	}
	
	for lineNum, line := range lines {
		for _, pattern := range patterns {
			if pattern.MatchString(line) {
				issues = append(issues, SecurityIssue{
					Type:        "sql-injection",
					Description: "Possible SQL injection vulnerability",
					Position:    Position{Start: lineNum, End: lineNum},
					Severity:    "high",
					Pattern:     line,
				})
			}
		}
	}
	
	return issues
}

func FindInsecureRandomUsage(code string) []SecurityIssue {
	issues := []SecurityIssue{}
	lines := strings.Split(code, "\n")
	
	// Patterns for insecure random number generation
	patterns := []*regexp.Regexp{
		regexp.MustCompile(`\bMath\.random\(\)`),           // JavaScript
		regexp.MustCompile(`\brandom\.random\(\)`),         // Python
		regexp.MustCompile(`\brand\(\)`),                   // C/C++
		regexp.MustCompile(`\bnew Random\(\)`),             // Java
	}
	
	for lineNum, line := range lines {
		for _, pattern := range patterns {
			if pattern.MatchString(line) {
				issues = append(issues, SecurityIssue{
					Type:        "weak-crypto",
					Description: "Insecure random number generation",
					Position:    Position{Start: lineNum, End: lineNum},
					Severity:    "medium",
					Pattern:     line,
				})
			}
		}
	}
	
	return issues
}

func DetectUnsafeDeserializationPatterns(code string) []Vulnerability {
	vulnerabilities := []Vulnerability{}
	lines := strings.Split(code, "\n")
	
	// Patterns for unsafe deserialization
	patterns := []*regexp.Regexp{
		regexp.MustCompile(`\bpickle\.loads?\(`),           // Python pickle
		regexp.MustCompile(`\beval\(`),                     // JavaScript/Python eval
		regexp.MustCompile(`\bJSON\.parse\(`),              // JavaScript JSON
		regexp.MustCompile(`\bunserialize\(`),              // PHP
		regexp.MustCompile(`\bObjectInputStream\(`),        // Java
	}
	
	for lineNum, line := range lines {
		for _, pattern := range patterns {
			if pattern.MatchString(line) {
				vulnerabilities = append(vulnerabilities, Vulnerability{
					Type:        "unsafe-deserialization",
					Description: "Potentially unsafe deserialization",
					Position:    Position{Start: lineNum, End: lineNum},
					RiskLevel:   "high",
					Mitigation:  "Validate input before deserialization or use safe alternatives",
				})
			}
		}
	}
	
	return vulnerabilities
}

func FindPathTraversalRisks(code string) []Vulnerability {
	vulnerabilities := []Vulnerability{}
	lines := strings.Split(code, "\n")
	
	// Patterns that might indicate path traversal vulnerabilities
	patterns := []*regexp.Regexp{
		regexp.MustCompile(`\.\./`),                       // Relative path traversal
		regexp.MustCompile(`\\\.\.\\\`),                     // Windows path traversal
		regexp.MustCompile(`(?i)file\s*\+\s*["']`),        // File concatenation
		regexp.MustCompile(`(?i)path\s*\+\s*["']`),        // Path concatenation
	}
	
	for lineNum, line := range lines {
		for _, pattern := range patterns {
			if pattern.MatchString(line) {
				vulnerabilities = append(vulnerabilities, Vulnerability{
					Type:        "path-traversal",
					Description: "Possible path traversal vulnerability",
					Position:    Position{Start: lineNum, End: lineNum},
					RiskLevel:   "medium",
					Mitigation:  "Sanitize and validate file paths, use path.join() or similar",
				})
			}
		}
	}
	
	return vulnerabilities
}

func DetectFilePermissionIssues(code string) []PermissionIssue {
	issues := []PermissionIssue{}
	lines := strings.Split(code, "\n")
	
	// Patterns for overly permissive file operations
	patterns := []*regexp.Regexp{
		regexp.MustCompile(`chmod\s+777`),                  // Unix overly permissive
		regexp.MustCompile(`chmod\s+666`),                  // Unix world writable
		regexp.MustCompile(`(?i)file.*create.*777`),        // Generic file creation
		regexp.MustCompile(`(?i)permission.*777`),          // Permission setting
	}
	
	for lineNum, line := range lines {
		for _, pattern := range patterns {
			if pattern.MatchString(line) {
				issues = append(issues, PermissionIssue{
					FileOperation: "chmod",
					Permission:    "777",
					Position:      Position{Start: lineNum, End: lineNum},
					Risk:          "Overly permissive file permissions",
				})
			}
		}
	}
	
	return issues
}

func FindXSSVulnerabilities(code string) []SecurityIssue {
	issues := []SecurityIssue{}
	lines := strings.Split(code, "\n")
	
	// Patterns that might indicate XSS vulnerabilities
	patterns := []*regexp.Regexp{
		regexp.MustCompile(`(?i)innerHTML\s*[=+]`),         // JavaScript innerHTML
		regexp.MustCompile(`(?i)outerHTML\s*[=+]`),         // JavaScript outerHTML
		regexp.MustCompile(`(?i)document\.write\(`),         // JavaScript document.write
		regexp.MustCompile(`(?i)eval\s*\(`),                // JavaScript eval
		regexp.MustCompile(`(?i)\$\{.*\}`),                 // Template literal injection
	}
	
	for lineNum, line := range lines {
		for _, pattern := range patterns {
			if pattern.MatchString(line) {
				issues = append(issues, SecurityIssue{
					Type:        "xss",
					Description: "Possible XSS vulnerability",
					Position:    Position{Start: lineNum, End: lineNum},
					Severity:    "high",
					Pattern:     line,
				})
			}
		}
	}
	
	return issues
}