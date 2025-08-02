package textlib

import (
	"strings"
	"testing"
)

// Comprehensive tests for code analysis functions

func TestCountLinesComprehensive(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected int
	}{
		{
			name:     "Empty string",
			input:    "",
			expected: 0,
		},
		{
			name:     "Single line",
			input:    "package main",
			expected: 1,
		},
		{
			name:     "Multiple lines",
			input:    "package main\n\nimport \"fmt\"\n\nfunc main() {\n    fmt.Println(\"Hello\")\n}",
			expected: 7,
		},
		{
			name:     "Only newlines",
			input:    "\n\n\n",
			expected: 4, // Split creates 4 empty strings
		},
		{
			name:     "Trailing newline",
			input:    "line1\nline2\n",
			expected: 3, // Split creates ["line1", "line2", ""]
		},
		{
			name:     "Windows line endings",
			input:    "line1\r\nline2\r\n",
			expected: 1, // Only counts \n
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := CountLines(test.input)
			if result != test.expected {
				t.Errorf("CountLines(%q): expected %d, got %d", test.input, test.expected, result)
			}
		})
	}
}

func TestCountBlankLinesComprehensive(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected int
	}{
		{
			name:     "No blank lines",
			input:    "line1\nline2\nline3",
			expected: 0,
		},
		{
			name:     "All blank lines",
			input:    "\n\n\n",
			expected: 3,
		},
		{
			name:     "Mixed content",
			input:    "line1\n\nline2\n\n\nline3",
			expected: 3,
		},
		{
			name:     "Whitespace only lines",
			input:    "line1\n   \n\t\nline2",
			expected: 2,
		},
		{
			name:     "Empty string",
			input:    "",
			expected: 0,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := CountBlankLines(test.input)
			if result != test.expected {
				t.Errorf("CountBlankLines(%q): expected %d, got %d", test.input, test.expected, result)
			}
		})
	}
}

func TestCountCommentLinesComprehensive(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected int
	}{
		{
			name:     "No comments",
			input:    "package main\nfunc test() {}",
			expected: 0,
		},
		{
			name:     "C++ style comments",
			input:    "// This is a comment\ncode here\n// Another comment",
			expected: 2,
		},
		{
			name:     "Python style comments",
			input:    "# Python comment\nprint('hello')\n# Another comment",
			expected: 2,
		},
		{
			name:     "C style single line",
			input:    "/* C comment */\ncode here",
			expected: 1,
		},
		{
			name:     "Mixed comment styles",
			input:    "// C++ comment\n# Python comment\n/* C comment */\ncode",
			expected: 3,
		},
		{
			name:     "SQL style comments",
			input:    "-- SQL comment\nSELECT * FROM table\n-- Another comment",
			expected: 2,
		},
		{
			name:     "MATLAB style comments",
			input:    "% MATLAB comment\nplot(x, y)\n% Another comment",
			expected: 2,
		},
		{
			name:     "Comments with indentation",
			input:    "    // Indented comment\n\t# Tab indented comment",
			expected: 2,
		},
		{
			name:     "Multi-line C style continuation",
			input:    "/*\n * Multi-line comment\n * continuation\n */",
			expected: 2, // Only lines starting with *
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := CountCommentLines(test.input)
			if result != test.expected {
				t.Errorf("CountCommentLines(%q): expected %d, got %d", test.input, test.expected, result)
			}
		})
	}
}

func TestExtractFunctionSignaturesComprehensive(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected int // number of functions expected
		language string
	}{
		{
			name:     "No functions",
			input:    "var x = 5;\nconsole.log(x);",
			expected: 0,
		},
		{
			name: "JavaScript functions",
			input: `function test() { return 1; }
function add(a, b) { return a + b; }`,
			expected: 2,
		},
		{
			name: "Python functions",
			input: `def hello():
    print("hello")

def add(x, y):
    return x + y`,
			expected: 2,
		},
		{
			name: "Go functions",
			input: `func main() {
    fmt.Println("hello")
}

func Add(a, b int) int {
    return a + b
}`,
			expected: 2,
		},
		{
			name: "Java methods",
			input: `public static void main(String[] args) {
    System.out.println("hello");
}

private int calculate(int x, int y) {
    return x + y;
}`,
			expected: 2,
		},
		{
			name: "C functions",
			input: `int main() {
    return 0;
}

void helper(int param) {
    // do something
}`,
			expected: 2,
		},
		{
			name: "Mixed languages",
			input: `function jsFunc() {}
def pyFunc():
    pass
func goFunc() {}`,
			expected: 3,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			functions := ExtractFunctionSignatures(test.input)
			if len(functions) != test.expected {
				t.Errorf("ExtractFunctionSignatures(%s): expected %d functions, got %d", 
					test.name, test.expected, len(functions))
				for i, fn := range functions {
					t.Logf("Function %d: %s (lang: %s)", i, fn.Name, fn.Language)
				}
			}
		})
	}
}

func TestCalculateCyclomaticComplexityComprehensive(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected int
	}{
		{
			name:     "Simple function",
			input:    "func simple() { return 1 }",
			expected: 1, // Base complexity
		},
		{
			name: "Function with if",
			input: `func withIf(x int) int {
    if x > 0 {
        return x
    }
    return 0
}`,
			expected: 2, // Base + if
		},
		{
			name: "Function with if-else",
			input: `func withIfElse(x int) int {
    if x > 0 {
        return x
    } else if x < 0 {
        return -x
    }
    return 0
}`,
			expected: 3, // Base + if + elif
		},
		{
			name: "Function with loop",
			input: `func withLoop() {
    for i := 0; i < 10; i++ {
        fmt.Println(i)
    }
}`,
			expected: 2, // Base + for
		},
		{
			name: "Complex function",
			input: `func complex(x int) int {
    if x > 10 {
        for i := 0; i < x; i++ {
            if i%2 == 0 {
                return i
            }
        }
    } else if x < 0 {
        while x < 0 {
            x++
        }
    }
    return x
}`,
			expected: 6, // Base + if + for + if + elif + while
		},
		{
			name: "Function with logical operators",
			input: `func withLogical(a, b bool) bool {
    return a && b || c && d
}`,
			expected: 3, // Base + 2 logical operators
		},
		{
			name: "Function with switch",
			input: `func withSwitch(x int) string {
    switch x {
    case 1:
        return "one"
    case 2:
        return "two"
    default:
        return "other"
    }
}`,
			expected: 3, // Base + 2 cases (default doesn't count typically)
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := CalculateCyclomaticComplexity(test.input)
			if result != test.expected {
				t.Errorf("CalculateCyclomaticComplexity(%s): expected %d, got %d", 
					test.name, test.expected, result)
			}
		})
	}
}

func TestCheckCamelCaseComprehensive(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		// Valid camelCase
		{"Simple camelCase", "camelCase", true},
		{"Single letter", "a", true},
		{"Multiple words", "thisIsCamelCase", true},
		{"With numbers", "camelCase123", true},
		{"CamelCase with acronym", "parseXMLData", true},
		
		// Invalid camelCase
		{"Empty string", "", false},
		{"PascalCase", "PascalCase", false},
		{"snake_case", "snake_case", false},
		{"kebab-case", "kebab-case", false},
		{"UPPER_CASE", "UPPER_CASE", false},
		{"With spaces", "camel Case", false},
		{"With dots", "camel.Case", false},
		{"Starting with number", "2camelCase", false},
		{"All lowercase", "alllowercase", false}, // No uppercase letters
		{"Mixed with underscore", "camel_Case", false},
		{"Mixed with hyphen", "camel-Case", false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := CheckCamelCase(test.input)
			if result != test.expected {
				t.Errorf("CheckCamelCase(%q): expected %v, got %v", 
					test.input, test.expected, result)
			}
		})
	}
}

func TestCheckSnakeCaseComprehensive(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		// Valid snake_case
		{"Simple snake_case", "snake_case", true},
		{"Single letter", "a", true},
		{"Multiple words", "this_is_snake_case", true},
		{"With numbers", "snake_case_123", true},
		{"Numbers in middle", "test_123_case", true},
		
		// Invalid snake_case
		{"Empty string", "", false},
		{"camelCase", "camelCase", false},
		{"PascalCase", "PascalCase", false},
		{"kebab-case", "kebab-case", false},
		{"UPPER_CASE", "UPPER_CASE", false},
		{"With spaces", "snake case", false},
		{"With dots", "snake.case", false},
		{"Starting with underscore", "_snake_case", false},
		{"Ending with underscore", "snake_case_", false},
		{"Double underscore", "snake__case", false},
		{"Mixed case", "Snake_Case", false},
		{"Starting with number", "2_snake_case", false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := CheckSnakeCase(test.input)
			if result != test.expected {
				t.Errorf("CheckSnakeCase(%q): expected %v, got %v", 
					test.input, test.expected, result)
			}
		})
	}
}

func TestCheckPascalCaseComprehensive(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		// Valid PascalCase
		{"Simple PascalCase", "PascalCase", true},
		{"Single letter", "A", true},
		{"Multiple words", "ThisIsPascalCase", true},
		{"With numbers", "PascalCase123", true},
		{"Acronym", "XMLHttpRequest", true},
		
		// Invalid PascalCase
		{"Empty string", "", false},
		{"camelCase", "camelCase", false},
		{"snake_case", "snake_case", false},
		{"kebab-case", "kebab-case", false},
		{"UPPER_CASE", "UPPER_CASE", false},
		{"With spaces", "Pascal Case", false},
		{"With underscore", "Pascal_Case", false},
		{"With hyphen", "Pascal-Case", false},
		{"Starting with lowercase", "pascalCase", false},
		{"Starting with number", "2PascalCase", false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := CheckPascalCase(test.input)
			if result != test.expected {
				t.Errorf("CheckPascalCase(%q): expected %v, got %v", 
					test.input, test.expected, result)
			}
		})
	}
}

func TestCheckKebabCaseComprehensive(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		// Valid kebab-case
		{"Simple kebab-case", "kebab-case", true},
		{"Single letter", "a", true},
		{"Multiple words", "this-is-kebab-case", true},
		{"With numbers", "kebab-case-123", true},
		{"Numbers in middle", "test-123-case", true},
		
		// Invalid kebab-case
		{"Empty string", "", false},
		{"camelCase", "camelCase", false},
		{"PascalCase", "PascalCase", false},
		{"snake_case", "snake_case", false},
		{"UPPER-CASE", "UPPER-CASE", false},
		{"With spaces", "kebab case", false},
		{"With underscore", "kebab_case", false},
		{"Starting with hyphen", "-kebab-case", false},
		{"Ending with hyphen", "kebab-case-", false},
		{"Double hyphen", "kebab--case", false},
		{"Mixed case", "Kebab-Case", false},
		{"Starting with number", "2-kebab-case", false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := CheckKebabCase(test.input)
			if result != test.expected {
				t.Errorf("CheckKebabCase(%q): expected %v, got %v", 
					test.input, test.expected, result)
			}
		})
	}
}

func TestValidateFunctionNamesComprehensive(t *testing.T) {
	code := `func calculateTotal() int { return 0 }
func a() {} // Too short
func thisIsAVeryVeryVeryLongFunctionNameThatExceedsTheLimit() {} // Too long
func snake_case_function() {}
func PascalCaseFunction() {}
func invalidFunction() {} // Missing verb
func getName() string { return "" } // Good with verb`

	tests := []struct {
		name  string
		rules NamingRules
		expectedIssueTypes []string
	}{
		{
			name: "Length constraints",
			rules: NamingRules{
				MinLength: 3,
				MaxLength: 20,
			},
			expectedIssueTypes: []string{"short", "long"},
		},
		{
			name: "camelCase style",
			rules: NamingRules{
				FunctionStyle: "camelCase",
			},
			expectedIssueTypes: []string{"convention"},
		},
		{
			name: "snake_case style",
			rules: NamingRules{
				FunctionStyle: "snake_case",
			},
			expectedIssueTypes: []string{"convention"},
		},
		{
			name: "Require verb",
			rules: NamingRules{
				RequireVerb: true,
			},
			expectedIssueTypes: []string{"verb"},
		},
		{
			name: "All rules",
			rules: NamingRules{
				FunctionStyle: "camelCase",
				RequireVerb:   true,
				MinLength:     3,
				MaxLength:     20,
			},
			expectedIssueTypes: []string{"short", "long", "convention", "verb"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			issues := ValidateFunctionNames(code, test.rules)
			
			// Check that we get some issues for this problematic code
			if len(issues) == 0 {
				t.Errorf("Expected some validation issues, got none")
				return
			}
			
			// Check that issues contain expected types
			issueText := strings.ToLower(strings.Join(func() []string {
				var descriptions []string
				for _, issue := range issues {
					descriptions = append(descriptions, issue.Description)
				}
				return descriptions
			}(), " "))
			
			for _, expectedType := range test.expectedIssueTypes {
				found := false
				switch expectedType {
				case "short":
					found = strings.Contains(issueText, "short")
				case "long":
					found = strings.Contains(issueText, "long")
				case "convention":
					found = strings.Contains(issueText, "convention")
				case "verb":
					found = strings.Contains(issueText, "verb")
				}
				
				if !found {
					t.Errorf("Expected issue type %s not found in: %s", expectedType, issueText)
				}
			}
		})
	}
}

func TestDetectTyposComprehensive(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected int // number of typos expected
	}{
		{
			name:     "No typos",
			input:    "This code has no spelling mistakes.",
			expected: 0,
		},
		{
			name:     "Single typo",
			input:    "The lenght of the array is important.",
			expected: 1,
		},
		{
			name:     "Multiple typos",
			input:    "I will recieve the arguement and seperate it.",
			expected: 3,
		},
		{
			name:     "Typo in variable name",
			input:    "var beginingIndex = 0;",
			expected: 1,
		},
		{
			name:     "Typo in comment",
			input:    "// This function will accomodate the request",
			expected: 1,
		},
		{
			name:     "Case insensitive",
			input:    "The LENGHT and Recieve functions.",
			expected: 2,
		},
		{
			name:     "Typos in strings",
			input:    `printf("Definately occured in maintainence");`,
			expected: 3,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			typos := DetectTypos(test.input)
			if len(typos) != test.expected {
				t.Errorf("DetectTypos(%s): expected %d typos, got %d", 
					test.name, test.expected, len(typos))
				for i, typo := range typos {
					t.Logf("Typo %d: %s at position %d", i, typo.Description, typo.Position.Start)
				}
			}
		})
	}
}

func TestFindLongFunctionsComprehensive(t *testing.T) {
	code := `func shortFunction() {
    return 1
}

func mediumFunction() {
    line1
    line2
    line3
    line4
    line5
}

func longFunction() {
    line1
    line2
    line3
    line4
    line5
    line6
    line7
    line8
    line9
    line10
    line11
    line12
}`

	tests := []struct {
		name     string
		maxLines int
		expected int // number of long functions
	}{
		{
			name:     "Very strict limit",
			maxLines: 3,
			expected: 2, // medium and long functions
		},
		{
			name:     "Medium limit",
			maxLines: 7,
			expected: 1, // only long function
		},
		{
			name:     "Generous limit",
			maxLines: 20,
			expected: 0, // no functions exceed limit
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			longFunctions := FindLongFunctions(code, test.maxLines)
			if len(longFunctions) != test.expected {
				t.Errorf("FindLongFunctions(maxLines=%d): expected %d functions, got %d", 
					test.maxLines, test.expected, len(longFunctions))
				for i, fn := range longFunctions {
					t.Logf("Long function %d: %s (%d lines)", i, fn.Name, fn.Lines)
				}
			}
		})
	}
}

func TestDetectDuplicateCodeComprehensive(t *testing.T) {
	code := `line1
line2
line3
duplicate_line_a
duplicate_line_b
duplicate_line_c
line7
line8
duplicate_line_a
duplicate_line_b
duplicate_line_c
line12`

	tests := []struct {
		name      string
		threshold int
		expected  int // number of duplicate blocks
	}{
		{
			name:      "Small threshold",
			threshold: 2,
			expected:  1, // Should find the 3-line duplicate
		},
		{
			name:      "Exact threshold",
			threshold: 3,
			expected:  1, // Should find the 3-line duplicate
		},
		{
			name:      "Large threshold",
			threshold: 5,
			expected:  0, // No duplicates of 5+ lines
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			duplicates := DetectDuplicateCode(code, test.threshold)
			if len(duplicates) != test.expected {
				t.Errorf("DetectDuplicateCode(threshold=%d): expected %d duplicates, got %d", 
					test.threshold, test.expected, len(duplicates))
				for i, dup := range duplicates {
					t.Logf("Duplicate %d: %d locations, similarity %.2f", 
						i, len(dup.Locations), dup.Similarity)
				}
			}
		})
	}
}

// Test error conditions
func TestCodeAnalysisErrorConditions(t *testing.T) {
	t.Run("Empty input handling", func(t *testing.T) {
		if CountLines("") != 0 {
			t.Errorf("Expected 0 lines for empty input")
		}
		
		if CountBlankLines("") != 0 {
			t.Errorf("Expected 0 blank lines for empty input")
		}
		
		if CountCommentLines("") != 0 {
			t.Errorf("Expected 0 comment lines for empty input")
		}
		
		functions := ExtractFunctionSignatures("")
		if len(functions) != 0 {
			t.Errorf("Expected no functions for empty input")
		}
	})

	t.Run("Invalid input handling", func(t *testing.T) {
		// Test with non-code content
		nonCode := "This is not code at all, just regular text."
		
		complexity := CalculateCyclomaticComplexity(nonCode)
		if complexity != 1 {
			t.Errorf("Expected base complexity 1 for non-code, got %d", complexity)
		}
		
		functions := ExtractFunctionSignatures(nonCode)
		if len(functions) != 0 {
			t.Errorf("Expected no functions in non-code text")
		}
	})
}

// Benchmark tests
func BenchmarkCountLines(b *testing.B) {
	code := strings.Repeat("line of code\n", 1000)
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		CountLines(code)
	}
}

func BenchmarkExtractFunctionSignatures(b *testing.B) {
	code := `function test1() { return 1; }
function test2(a, b) { return a + b; }
def python_func():
    return "hello"
func goFunc() int {
    return 42
}`
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ExtractFunctionSignatures(code)
	}
}

func BenchmarkCalculateCyclomaticComplexity(b *testing.B) {
	code := `func complex(x int) int {
    if x > 10 {
        for i := 0; i < x; i++ {
            if i%2 == 0 {
                return i
            }
        }
    } else if x < 0 {
        while x < 0 {
            x++
        }
    }
    return x
}`
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		CalculateCyclomaticComplexity(code)
	}
}