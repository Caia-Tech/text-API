package textlib

import (
	"regexp"
	"strconv"
	"strings"
)

// Mathematical expression analysis data structures

type MathExpression struct {
	Expression string
	Type       string // equation, inequality, function, polynomial, etc.
	Variables  []string
	Constants  []float64
	Operators  []string
	Position   Position
}

type AlgebraicTerm struct {
	Coefficient float64
	Variable    string
	Exponent    int
}

type Equation struct {
	LeftSide  []AlgebraicTerm
	RightSide []AlgebraicTerm
	Type      string // linear, quadratic, cubic, etc.
}

type MathPattern struct {
	Type        string
	Pattern     string
	Description string
	Examples    []string
}

type ValidationResult struct {
	IsValid     bool
	Errors      []string
	Warnings    []string
	Suggestions []string
}

// Mathematical Expression Parsing functions

func ExtractMathExpressions(text string) []MathExpression {
	expressions := []MathExpression{}
	
	// Patterns for mathematical expressions
	patterns := []*regexp.Regexp{
		regexp.MustCompile(`[a-zA-Z]\s*[=<>≤≥≠]\s*[^.!?]*`),  // Equations and inequalities
		regexp.MustCompile(`\b\d*[a-zA-Z]\^?\d*\s*[+\-*/]\s*\d*[a-zA-Z]\^?\d*`), // Algebraic expressions
		regexp.MustCompile(`\b[a-zA-Z]\([^)]*\)\s*[=<>]\s*[^.!?]*`), // Functions
		regexp.MustCompile(`∫|∑|∏|lim|sin|cos|tan|log|ln`), // Mathematical functions
		regexp.MustCompile(`\b\d+\s*[+\-*/]\s*\d+\s*[=]\s*\d+\b`), // Arithmetic
	}
	
	patternTypes := []string{"equation", "algebraic", "function", "calculus", "arithmetic"}
	
	for i, pattern := range patterns {
		matches := pattern.FindAllStringIndex(text, -1)
		for _, match := range matches {
			expr := MathExpression{
				Expression: text[match[0]:match[1]],
				Type:       patternTypes[i],
				Position:   Position{Start: match[0], End: match[1]},
				Variables:  ExtractVariables(text[match[0]:match[1]]),
				Constants:  ExtractConstants(text[match[0]:match[1]]),
				Operators:  ExtractOperators(text[match[0]:match[1]]),
			}
			expressions = append(expressions, expr)
		}
	}
	
	return expressions
}

func ExtractVariables(expression string) []string {
	variables := []string{}
	// Pattern to catch single letters that represent variables (including after digits)
	varPattern := regexp.MustCompile(`[a-zA-Z]`)
	matches := varPattern.FindAllString(expression, -1)
	
	// Remove duplicates
	seen := make(map[string]bool)
	for _, match := range matches {
		if !seen[match] && !isMathFunction(match) {
			variables = append(variables, match)
			seen[match] = true
		}
	}
	
	return variables
}

func ExtractConstants(expression string) []float64 {
	constants := []float64{}
	constPattern := regexp.MustCompile(`\b\d+\.?\d*\b`)
	matches := constPattern.FindAllString(expression, -1)
	
	for _, match := range matches {
		if val, err := strconv.ParseFloat(match, 64); err == nil {
			constants = append(constants, val)
		}
	}
	
	return constants
}

func ExtractOperators(expression string) []string {
	operators := []string{}
	opPattern := regexp.MustCompile(`[+\-*/^=<>≤≥≠]|sin|cos|tan|log|ln|sqrt`)
	matches := opPattern.FindAllString(expression, -1)
	
	// Remove duplicates while preserving order
	seen := make(map[string]bool)
	for _, match := range matches {
		if !seen[match] {
			operators = append(operators, match)
			seen[match] = true
		}
	}
	
	return operators
}

func isMathFunction(s string) bool {
	functions := []string{"sin", "cos", "tan", "log", "ln", "exp", "abs", "sqrt"}
	for _, fn := range functions {
		if s == fn {
			return true
		}
	}
	return false
}

func ParseEquation(equation string) Equation {
	eq := Equation{}
	
	// Split on equals sign
	parts := strings.Split(equation, "=")
	if len(parts) != 2 {
		return eq
	}
	
	eq.LeftSide = parseAlgebraicSide(strings.TrimSpace(parts[0]))
	eq.RightSide = parseAlgebraicSide(strings.TrimSpace(parts[1]))
	eq.Type = determineEquationType(eq)
	
	return eq
}

func parseAlgebraicSide(side string) []AlgebraicTerm {
	terms := []AlgebraicTerm{}
	
	// Simple term parsing - handle basic polynomial terms
	termPattern := regexp.MustCompile(`([+\-]?\s*\d*)\s*([a-zA-Z]?)(\^?\d*)?`)
	matches := termPattern.FindAllStringSubmatch(side, -1)
	
	for _, match := range matches {
		if len(match) >= 4 {
			term := AlgebraicTerm{}
			
			// Parse coefficient
			coeffStr := strings.ReplaceAll(match[1], " ", "")
			if coeffStr == "" || coeffStr == "+" {
				term.Coefficient = 1
			} else if coeffStr == "-" {
				term.Coefficient = -1
			} else {
				if coeff, err := strconv.ParseFloat(coeffStr, 64); err == nil {
					term.Coefficient = coeff
				}
			}
			
			// Parse variable
			term.Variable = match[2]
			
			// Parse exponent
			expStr := strings.TrimPrefix(match[3], "^")
			if expStr == "" && term.Variable != "" {
				term.Exponent = 1
			} else if exp, err := strconv.Atoi(expStr); err == nil {
				term.Exponent = exp
			}
			
			if term.Coefficient != 0 || term.Variable != "" {
				terms = append(terms, term)
			}
		}
	}
	
	return terms
}

func determineEquationType(eq Equation) string {
	maxDegree := 0
	
	allTerms := append(eq.LeftSide, eq.RightSide...)
	for _, term := range allTerms {
		if term.Exponent > maxDegree {
			maxDegree = term.Exponent
		}
	}
	
	switch maxDegree {
	case 0:
		return "constant"
	case 1:
		return "linear"
	case 2:
		return "quadratic"
	case 3:
		return "cubic"
	default:
		return "polynomial"
	}
}

// Algebraic Operations functions

func SimplifyExpression(expression string) string {
	// Basic expression simplification
	expr := strings.ReplaceAll(expression, " ", "")
	
	// Combine like terms (basic implementation)
	expr = combineLikeTerms(expr)
	
	// Remove redundant operations
	expr = strings.ReplaceAll(expr, "+-", "-")
	expr = strings.ReplaceAll(expr, "--", "+")
	expr = strings.ReplaceAll(expr, "*1", "")
	expr = strings.ReplaceAll(expr, "1*", "")
	
	return expr
}

func combineLikeTerms(expression string) string {
	// Simple like term combination - could be expanded
	termMap := make(map[string]float64)
	
	// Extract terms with their coefficients
	termPattern := regexp.MustCompile(`([+\-]?\d*\.?\d*)([a-zA-Z]+\^?\d*)`)
	matches := termPattern.FindAllStringSubmatch(expression, -1)
	
	for _, match := range matches {
		if len(match) >= 3 {
			coeffStr := match[1]
			variable := match[2]
			
			if coeffStr == "" || coeffStr == "+" {
				termMap[variable] += 1
			} else if coeffStr == "-" {
				termMap[variable] -= 1
			} else {
				if coeff, err := strconv.ParseFloat(coeffStr, 64); err == nil {
					termMap[variable] += coeff
				}
			}
		}
	}
	
	// Reconstruct expression
	result := ""
	first := true
	for variable, coeff := range termMap {
		if coeff == 0 {
			continue
		}
		
		if !first && coeff > 0 {
			result += "+"
		}
		
		if coeff == 1 && variable != "" {
			result += variable
		} else if coeff == -1 && variable != "" {
			result += "-" + variable
		} else {
			result += strconv.FormatFloat(coeff, 'g', -1, 64)
			if variable != "" {
				result += variable
			}
		}
		
		first = false
	}
	
	if result == "" {
		return "0"
	}
	
	return result
}

func ExpandExpression(expression string) string {
	// Basic expansion for simple cases like (a+b)(c+d)
	expr := strings.ReplaceAll(expression, " ", "")
	
	// Pattern for (term1+term2)(term3+term4)
	expansionPattern := regexp.MustCompile(`\(([^)]+)\)\(([^)]+)\)`)
	match := expansionPattern.FindStringSubmatch(expr)
	
	if len(match) >= 3 {
		first := strings.Split(match[1], "+")
		second := strings.Split(match[2], "+")
		
		expanded := ""
		for i, term1 := range first {
			for j, term2 := range second {
				if i > 0 || j > 0 {
					expanded += "+"
				}
				expanded += term1 + "*" + term2
			}
		}
		
		return expanded
	}
	
	return expression
}

func FactorExpression(expression string) string {
	// Basic factoring - extract common factors
	expr := strings.ReplaceAll(expression, " ", "")
	
	// Find common factors in polynomial terms
	terms := strings.FieldsFunc(expr, func(r rune) bool {
		return r == '+' || r == '-'
	})
	
	if len(terms) < 2 {
		return expression
	}
	
	// Find greatest common factor
	gcf := findGCF(terms)
	if gcf != "" && gcf != "1" {
		factored := gcf + "("
		first := true
		for _, term := range terms {
			if !first {
				factored += "+"
			}
			factored += strings.TrimPrefix(term, gcf)
			first = false
		}
		factored += ")"
		return factored
	}
	
	return expression
}

func findGCF(terms []string) string {
	if len(terms) == 0 {
		return ""
	}
	
	// Simple GCF finding for numeric coefficients
	coefficients := []int{}
	for _, term := range terms {
		coeffPattern := regexp.MustCompile(`^\d+`)
		match := coeffPattern.FindString(term)
		if match != "" {
			if coeff, err := strconv.Atoi(match); err == nil {
				coefficients = append(coefficients, coeff)
			}
		}
	}
	
	if len(coefficients) > 1 {
		gcf := gcd(coefficients[0], coefficients[1])
		for i := 2; i < len(coefficients); i++ {
			gcf = gcd(gcf, coefficients[i])
		}
		if gcf > 1 {
			return strconv.Itoa(gcf)
		}
	}
	
	return "1"
}

func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

// Validation functions

func ValidateMathExpression(expression string) ValidationResult {
	result := ValidationResult{IsValid: true}
	
	// Check for balanced parentheses
	if !hasBalancedParentheses(expression) {
		result.IsValid = false
		result.Errors = append(result.Errors, "Unbalanced parentheses")
	}
	
	// Check for valid operators
	invalidOps := findInvalidOperators(expression)
	if len(invalidOps) > 0 {
		result.Warnings = append(result.Warnings, "Potentially invalid operators found")
	}
	
	// Check for consecutive operators
	if hasConsecutiveOperators(expression) {
		result.IsValid = false
		result.Errors = append(result.Errors, "Consecutive operators found")
	}
	
	// Check for empty expressions
	if strings.TrimSpace(expression) == "" {
		result.IsValid = false
		result.Errors = append(result.Errors, "Empty expression")
	}
	
	return result
}

func hasBalancedParentheses(expression string) bool {
	count := 0
	for _, char := range expression {
		if char == '(' {
			count++
		} else if char == ')' {
			count--
			if count < 0 {
				return false
			}
		}
	}
	return count == 0
}

func findInvalidOperators(expression string) []string {
	validOps := map[string]bool{
		"+": true, "-": true, "*": true, "/": true, "^": true,
		"=": true, "<": true, ">": true, "≤": true, "≥": true, "≠": true,
		"sin": true, "cos": true, "tan": true, "log": true, "ln": true,
		"sqrt": true, "abs": true, "exp": true,
	}
	
	opPattern := regexp.MustCompile(`[+\-*/^=<>≤≥≠]|sin|cos|tan|log|ln|sqrt|abs|exp`)
	matches := opPattern.FindAllString(expression, -1)
	
	invalid := []string{}
	for _, match := range matches {
		if !validOps[match] {
			invalid = append(invalid, match)
		}
	}
	
	return invalid
}

func hasConsecutiveOperators(expression string) bool {
	opPattern := regexp.MustCompile(`[+\-*/^]{2,}`)
	return opPattern.MatchString(expression)
}

// Pattern Recognition functions

func DetectMathPatterns(text string) []MathPattern {
	patterns := []MathPattern{}
	
	// Arithmetic sequences
	if detectArithmeticSequence(text) {
		patterns = append(patterns, MathPattern{
			Type:        "arithmetic_sequence",
			Pattern:     "a, a+d, a+2d, ...",
			Description: "Arithmetic sequence detected",
			Examples:    []string{"2, 4, 6, 8, 10"},
		})
	}
	
	// Geometric sequences
	if detectGeometricSequence(text) {
		patterns = append(patterns, MathPattern{
			Type:        "geometric_sequence",
			Pattern:     "a, ar, ar², ...",
			Description: "Geometric sequence detected",
			Examples:    []string{"2, 4, 8, 16, 32"},
		})
	}
	
	// Quadratic patterns
	if detectQuadraticPattern(text) {
		patterns = append(patterns, MathPattern{
			Type:        "quadratic",
			Pattern:     "ax² + bx + c",
			Description: "Quadratic expression pattern",
			Examples:    []string{"x² + 2x + 1", "2x² - 3x + 1"},
		})
	}
	
	return patterns
}

func detectArithmeticSequence(text string) bool {
	// Look for sequences like "1, 3, 5, 7" or "2, 4, 6, 8"
	seqPattern := regexp.MustCompile(`\b\d+,\s*\d+,\s*\d+`)
	matches := seqPattern.FindAllString(text, -1)
	
	for _, match := range matches {
		nums := regexp.MustCompile(`\d+`).FindAllString(match, -1)
		if len(nums) >= 3 {
			// Check if differences are constant
			diff1, _ := strconv.Atoi(nums[1])
			diff0, _ := strconv.Atoi(nums[0])
			diff2, _ := strconv.Atoi(nums[2])
			
			if (diff1 - diff0) == (diff2 - diff1) {
				return true
			}
		}
	}
	
	return false
}

func detectGeometricSequence(text string) bool {
	// Look for sequences where each term is multiplied by a constant
	seqPattern := regexp.MustCompile(`\b\d+,\s*\d+,\s*\d+`)
	matches := seqPattern.FindAllString(text, -1)
	
	for _, match := range matches {
		nums := regexp.MustCompile(`\d+`).FindAllString(match, -1)
		if len(nums) >= 3 {
			n0, _ := strconv.ParseFloat(nums[0], 64)
			n1, _ := strconv.ParseFloat(nums[1], 64)
			n2, _ := strconv.ParseFloat(nums[2], 64)
			
			if n0 != 0 && n1 != 0 && (n1/n0) == (n2/n1) {
				return true
			}
		}
	}
	
	return false
}

func detectQuadraticPattern(text string) bool {
	// Look for patterns like ax² + bx + c
	quadPattern := regexp.MustCompile(`\b\d*[a-zA-Z]\^?2\s*[+\-]\s*\d*[a-zA-Z]\s*[+\-]\s*\d+`)
	return quadPattern.MatchString(text)
}

func ClassifyEquation(equation string) string {
	eq := ParseEquation(equation)
	return eq.Type
}

func FindMathConstants(text string) []string {
	constants := []string{}
	
	// Common mathematical constants
	constPatterns := map[string]*regexp.Regexp{
		"π":     regexp.MustCompile(`π|pi|3\.14159`),
		"e":     regexp.MustCompile(`\be\b|2\.71828`),
		"φ":     regexp.MustCompile(`φ|phi|1\.618`),
		"√2":    regexp.MustCompile(`√2|sqrt\(2\)|1\.414`),
		"ln(2)": regexp.MustCompile(`ln\(2\)|0\.693`),
	}
	
	for constant, pattern := range constPatterns {
		if pattern.MatchString(text) {
			constants = append(constants, constant)
		}
	}
	
	return constants
}

func DetectMathNotation(text string) []string {
	notations := []string{}
	
	// Mathematical notation patterns
	notationPatterns := map[string]*regexp.Regexp{
		"integral":    regexp.MustCompile(`∫|\\int`),
		"summation":   regexp.MustCompile(`∑|\\sum`),
		"product":     regexp.MustCompile(`∏|\\prod`),
		"limit":       regexp.MustCompile(`lim|\\lim`),
		"derivative":  regexp.MustCompile(`d/dx|\\frac{d}{dx}|'`),
		"partial":     regexp.MustCompile(`∂|\\partial`),
		"infinity":    regexp.MustCompile(`∞|\\infty`),
		"set":         regexp.MustCompile(`∈|∉|⊂|⊃|∪|∩`),
		"logic":       regexp.MustCompile(`∀|∃|∧|∨|¬|⟹|⟺`),
		"inequality":  regexp.MustCompile(`≤|≥|≠|≈|≡`),
	}
	
	for notation, pattern := range notationPatterns {
		if pattern.MatchString(text) {
			notations = append(notations, notation)
		}
	}
	
	return notations
}