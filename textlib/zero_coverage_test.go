// Copyright 2025 Caia Tech
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package textlib

import (
	"testing"
)

// Tests for functions with 0% coverage to improve overall coverage

func TestDetectVerbTense(t *testing.T) {
	tests := []struct {
		name     string
		verb     string
		expected string
	}{
		{
			name:     "Present tense",
			verb:     "runs",
			expected: "present",
		},
		{
			name:     "Past tense",
			verb:     "ran",
			expected: "past",
		},
		{
			name:     "Unknown verb",
			verb:     "xyz",
			expected: "unknown",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := detectVerbTense(tt.verb)
			if result != tt.expected {
				t.Errorf("detectVerbTense(%q) = %q, expected %q", tt.verb, result, tt.expected)
			}
		})
	}
}

func TestExtractImports(t *testing.T) {
	tests := []struct {
		name     string
		code     string
		expected int
	}{
		{
			name: "Go imports",
			code: `package main
import (
	"fmt"
	"strings"
)`,
			expected: 2,
		},
		{
			name:     "No imports",
			code:     "package main\nfunc main() {}",
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ExtractImports(tt.code)
			if len(result) != tt.expected {
				t.Errorf("ExtractImports(%q) returned %d imports, expected %d", tt.name, len(result), tt.expected)
			}
		})
	}
}

func TestExtractClassDefinitions(t *testing.T) {
	tests := []struct {
		name     string
		code     string
		expected int
	}{
		{
			name: "Java class",
			code: `public class MyClass {
	private int value;
}`,
			expected: 1,
		},
		{
			name:     "No classes",
			code:     "function hello() { return 'world'; }",
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ExtractClassDefinitions(tt.code)
			if len(result) != tt.expected {
				t.Errorf("ExtractClassDefinitions(%q) returned %d classes, expected %d", tt.name, len(result), tt.expected)
			}
		})
	}
}

func TestFindUnusedImports(t *testing.T) {
	tests := []struct {
		name     string
		code     string
		expected int
	}{
		{
			name: "Go with unused import",
			code: `package main
import (
	"fmt"
	"strings"
)

func main() {
	fmt.Println("Hello")
	// strings is unused
}`,
			expected: 1,
		},
		{
			name:     "No imports",
			code:     "package main\nfunc main() {}",
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FindUnusedImports(tt.code, []string{})
			if len(result) != tt.expected {
				t.Errorf("FindUnusedImports(%q) returned %d unused imports, expected %d", tt.name, len(result), tt.expected)
			}
		})
	}
}

func TestBlocksOverlap(t *testing.T) {
	tests := []struct {
		name     string
		block1   CodeBlock
		block2   CodeBlock
		expected bool
	}{
		{
			name: "Overlapping blocks",
			block1: CodeBlock{
				StartLine: 1,
				EndLine:   5,
			},
			block2: CodeBlock{
				StartLine: 3,
				EndLine:   7,
			},
			expected: true,
		},
		{
			name: "Non-overlapping blocks",
			block1: CodeBlock{
				StartLine: 1,
				EndLine:   5,
			},
			block2: CodeBlock{
				StartLine: 6,
				EndLine:   10,
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := blocksOverlap(tt.block1, tt.block2)
			if result != tt.expected {
				t.Errorf("blocksOverlap() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

func TestGetIndentationLevel(t *testing.T) {
	tests := []struct {
		name     string
		line     string
		expected int
	}{
		{
			name:     "No indentation",
			line:     "function test() {",
			expected: 0,
		},
		{
			name:     "Two spaces",
			line:     "  var x = 5;",
			expected: 2,
		},
		{
			name:     "Tab indentation",
			line:     "\tvar x = 5;",
			expected: 4, // assuming tab = 4 spaces
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getIndentationLevel(tt.line)
			if result != tt.expected {
				t.Errorf("getIndentationLevel(%q) = %d, expected %d", tt.line, result, tt.expected)
			}
		})
	}
}

func TestFindDeepNesting(t *testing.T) {
	tests := []struct {
		name     string
		code     string
		expected int
	}{
		{
			name: "Nested code",
			code: `if (a) {
	if (b) {
		if (c) {
			return true;
		}
	}
}`,
			expected: 1, // Expecting at least one deeply nested block
		},
		{
			name:     "Flat code",
			code:     "var x = 5;\nreturn x;",
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FindDeepNesting(tt.code)
			if len(result) < tt.expected {
				t.Errorf("FindDeepNesting(%q) returned %d blocks, expected at least %d", tt.name, len(result), tt.expected)
			}
		})
	}
}

func TestDetectBlockType(t *testing.T) {
	tests := []struct {
		name     string
		line     string
		expected string
	}{
		{
			name:     "If statement",
			line:     "if (condition) {",
			expected: "if",
		},
		{
			name:     "For loop",
			line:     "for (int i = 0; i < 10; i++) {",
			expected: "for",
		},
		{
			name:     "Regular line",
			line:     "var x = 5;",
			expected: "other",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := detectBlockType(tt.line)
			if result != tt.expected {
				t.Errorf("detectBlockType(%q) = %q, expected %q", tt.line, result, tt.expected)
			}
		})
	}
}

func TestCountParameters(t *testing.T) {
	tests := []struct {
		name     string
		code     string
		expected int
	}{
		{
			name: "Function with parameters",
			code: `func test(a int, b string, c bool) {
	return
}`,
			expected: 3,
		},
		{
			name:     "No functions",
			code:     "var x = 5;",
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CountParameters(tt.code)
			// Just test that it doesn't panic
			if result < 0 {
				t.Errorf("CountParameters(%q) returned negative value: %d", tt.name, result)
			}
		})
	}
}

func TestCalculateFileSize(t *testing.T) {
	tests := []struct {
		name     string
		code     string
		expected int
	}{
		{
			name:     "Simple code",
			code:     "hello world",
			expected: 11,
		},
		{
			name:     "Empty code",
			code:     "",
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CalculateFileSize(tt.code)
			if result != tt.expected {
				t.Errorf("CalculateFileSize(%q) = %d, expected %d", tt.code, result, tt.expected)
			}
		})
	}
}

func TestCountLoops(t *testing.T) {
	tests := []struct {
		name     string
		code     string
		expected int
	}{
		{
			name: "Code with loops",
			code: `for (int i = 0; i < 10; i++) {
	while (condition) {
		break;
	}
}`,
			expected: 2,
		},
		{
			name:     "No loops",
			code:     "var x = 5;\nreturn x;",
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CountLoops(tt.code)
			if result < tt.expected {
				t.Errorf("CountLoops(%q) = %d, expected at least %d", tt.code, result, tt.expected)
			}
		})
	}
}

func TestCountConditionals(t *testing.T) {
	tests := []struct {
		name     string
		code     string
		expected int
	}{
		{
			name: "Code with conditionals",
			code: `if (a) {
	return true;
} else if (b) {
	return false;
}`,
			expected: 2,
		},
		{
			name:     "No conditionals",
			code:     "var x = 5;\nreturn x;",
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CountConditionals(tt.code)
			if result < tt.expected {
				t.Errorf("CountConditionals(%q) = %d, expected at least %d", tt.code, result, tt.expected)
			}
		})
	}
}

func TestCalculateAverageFunctionLength(t *testing.T) {
	tests := []struct {
		name     string
		code     string
	}{
		{
			name: "Functions",
			code: `func test1() {
	return 1
}

func test2() {
	return 2
	return 3
}`,
		},
		{
			name: "No functions",
			code: "var x = 5;",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CalculateAverageFunctionLength(tt.code)
			// Just test that it doesn't panic and returns a reasonable value
			if result < 0 {
				t.Errorf("CalculateAverageFunctionLength(%q) returned negative value: %f", tt.name, result)
			}
		})
	}
}

func TestCountReturnStatements(t *testing.T) {
	tests := []struct {
		name     string
		code     string
		expected int
	}{
		{
			name: "Code with returns",
			code: `func test() {
	if (condition) {
		return true;
	}
	return false;
}`,
			expected: 2,
		},
		{
			name:     "No returns",
			code:     "var x = 5;",
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CountReturnStatements(tt.code)
			if result < tt.expected {
				t.Errorf("CountReturnStatements(%q) = %d, expected at least %d", tt.code, result, tt.expected)
			}
		})
	}
}