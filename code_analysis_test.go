package textlib

import (
	"testing"
)

func TestCountLines(t *testing.T) {
	code := `package main

import "fmt"

func main() {
    fmt.Println("Hello, world!")
}`
	
	expected := 7
	result := CountLines(code)
	
	if result != expected {
		t.Errorf("Expected %d lines, got %d", expected, result)
	}
}

func TestCountBlankLines(t *testing.T) {
	code := `package main

import "fmt"


func main() {
    fmt.Println("Hello, world!")
}`
	
	expected := 3
	result := CountBlankLines(code)
	
	if result != expected {
		t.Errorf("Expected %d blank lines, got %d", expected, result)
	}
}

func TestCountCommentLines(t *testing.T) {
	code := `package main

// This is a comment
import "fmt"

/* This is also a comment */
func main() {
    // Another comment
    fmt.Println("Hello, world!")
}`
	
	expected := 3
	result := CountCommentLines(code)
	
	if result != expected {
		t.Errorf("Expected %d comment lines, got %d", expected, result)
	}
}

func TestExtractFunctionSignatures(t *testing.T) {
	code := `package main

func main() {
    fmt.Println("Hello")
}

func add(a int, b int) int {
    return a + b
}`
	
	functions := ExtractFunctionSignatures(code)
	
	if len(functions) != 2 {
		t.Errorf("Expected 2 functions, got %d", len(functions))
	}
	
	if functions[0].Name != "main" {
		t.Errorf("Expected first function to be 'main', got '%s'", functions[0].Name)
	}
	
	if functions[1].Name != "add" {
		t.Errorf("Expected second function to be 'add', got '%s'", functions[1].Name)
	}
}

func TestCalculateCyclomaticComplexity(t *testing.T) {
	code := `func complex(x int) int {
    if x > 10 {
        for i := 0; i < x; i++ {
            if i%2 == 0 {
                return i
            }
        }
    }
    return x
}`
	
	complexity := CalculateCyclomaticComplexity(code)
	
	// Base complexity (1) + if (1) + for (1) + if (1) = 4
	expected := 4
	if complexity != expected {
		t.Errorf("Expected complexity %d, got %d", expected, complexity)
	}
}

func TestCheckCamelCase(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"camelCase", true},
		{"CamelCase", false}, // Should start with lowercase
		{"snake_case", false},
		{"kebab-case", false},
		{"", false},
		{"a", true},
	}
	
	for _, test := range tests {
		result := CheckCamelCase(test.input)
		if result != test.expected {
			t.Errorf("CheckCamelCase(%s): expected %v, got %v", test.input, test.expected, result)
		}
	}
}

func TestCheckSnakeCase(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"snake_case", true},
		{"camelCase", false},
		{"PascalCase", false},
		{"kebab-case", false},
		{"", false},
		{"a", true},
	}
	
	for _, test := range tests {
		result := CheckSnakeCase(test.input)
		if result != test.expected {
			t.Errorf("CheckSnakeCase(%s): expected %v, got %v", test.input, test.expected, result)
		}
	}
}

func TestFindHardcodedPasswords(t *testing.T) {
	code := `const password = "secret123"
const apiKey = "abc123xyz"
const config = { host: "localhost" }`
	
	issues := FindHardcodedPasswords(code)
	
	if len(issues) < 2 {
		t.Errorf("Expected at least 2 security issues, got %d", len(issues))
	}
}

func TestDetectSQLInjectionPatterns(t *testing.T) {
	code := `query := "SELECT * FROM users WHERE id = " + userID + " AND active = 1"`
	
	issues := DetectSQLInjectionPatterns(code)
	
	if len(issues) == 0 {
		t.Errorf("Expected SQL injection pattern to be detected")
	}
}

func TestNormalizeWhitespace(t *testing.T) {
	code := "func test() {\t\n\t\treturn true   \n\n\n\n}"
	
	normalized := NormalizeWhitespace(code)
	
	// Should convert tabs to spaces and limit blank lines
	if !contains(normalized, "    ") { // Should contain spaces instead of tabs
		t.Errorf("Expected tabs to be converted to spaces")
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && s[:len(substr)] == substr
}

func TestDetectIndentationStyle(t *testing.T) {
	tabCode := "func test() {\n\tif true {\n\t\treturn\n\t}\n}"
	spaceCode := "func test() {\n    if true {\n        return\n    }\n}"
	
	tabReport := DetectIndentationStyle(tabCode)
	spaceReport := DetectIndentationStyle(spaceCode)
	
	if tabReport.Style != "tabs" {
		t.Errorf("Expected 'tabs', got '%s'", tabReport.Style)
	}
	
	if spaceReport.Style != "spaces" {
		t.Errorf("Expected 'spaces', got '%s'", spaceReport.Style)
	}
}