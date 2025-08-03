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

// Simple tests for coverage improvement that will definitely compile

func TestDetectCodeLanguageSimple(t *testing.T) {
	tests := []struct {
		name     string
		code     string
		expected string
	}{
		{
			name:     "Go code",
			code:     "func main() { fmt.Println(\"Hello\") }",
			expected: "go",
		},
		{
			name:     "JavaScript code",
			code:     "function hello() { console.log(\"Hello\"); }",
			expected: "javascript",
		},
		{
			name:     "Empty code",
			code:     "",
			expected: "unknown",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := detectCodeLanguage(tt.code)
			if result != tt.expected {
				t.Errorf("detectCodeLanguage(%q) = %q, expected %q", tt.code, result, tt.expected)
			}
		})
	}
}

func TestParseCurrencySimple(t *testing.T) {
	tests := []struct {
		name     string
		text     string
		expected string
	}{
		{
			name:     "USD symbol",
			text:     "$100",
			expected: "USD",
		},
		{
			name:     "No currency",
			text:     "just text",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseCurrency(tt.text)
			if result != tt.expected {
				t.Errorf("parseCurrency(%q) = %q, expected %q", tt.text, result, tt.expected)
			}
		})
	}
}

func TestNormalizeTimeSimple(t *testing.T) {
	tests := []struct {
		name     string
		time     string
		expected string
	}{
		{
			name:     "24-hour format",
			time:     "14:45",
			expected: "14:45",
		},
		{
			name:     "Empty time",
			time:     "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := normalizeTime(tt.time)
			if result != tt.expected {
				t.Errorf("normalizeTime(%q) = %q, expected %q", tt.time, result, tt.expected)
			}
		})
	}
}

func TestConvertOrdinalToNumberSimple(t *testing.T) {
	tests := []struct {
		name     string
		ordinal  string
		expected string
	}{
		{
			name:     "First",
			ordinal:  "first",
			expected: "1",
		},
		{
			name:     "Invalid ordinal",
			ordinal:  "random",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := convertOrdinalToNumber(tt.ordinal)
			if result != tt.expected {
				t.Errorf("convertOrdinalToNumber(%q) = %q, expected %q", tt.ordinal, result, tt.expected)
			}
		})
	}
}

func TestExtractCodeEntitiesSimple(t *testing.T) {
	tests := []struct {
		name     string
		text     string
		expected int
	}{
		{
			name:     "Function calls",
			text:     "Call the function calculateSum() and then processResult().",
			expected: 2,
		},
		{
			name:     "No code entities",
			text:     "This is just regular text without any code references.",
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extractCodeEntities(tt.text)
			if len(result) != tt.expected {
				t.Errorf("extractCodeEntities(%q) returned %d entities, expected %d", tt.name, len(result), tt.expected)
			}
		})
	}
}

func TestExtractMetricEntitiesSimple(t *testing.T) {
	tests := []struct {
		name     string
		text     string
		expected int
	}{
		{
			name:     "Distance and weight",
			text:     "The package weighs 2.5 kg and traveled 150 km.",
			expected: 2,
		},
		{
			name:     "No metrics",
			text:     "This text has no measurements.",
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extractMetricEntities(tt.text)
			if len(result) != tt.expected {
				t.Errorf("extractMetricEntities(%q) returned %d entities, expected %d", tt.name, len(result), tt.expected)
			}
		})
	}
}

func TestNormalizeCodeBlockSimple(t *testing.T) {
	tests := []struct {
		name     string
		code     string
	}{
		{
			name: "Mixed indentation",
			code: "  function test() {\n\t\treturn true;\n  }",
		},
		{
			name: "Clean code",
			code: "clean code",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := normalizeCodeBlock(tt.code)
			// Just test that it doesn't panic and returns something
			if len(result) == 0 && len(tt.code) > 0 {
				t.Errorf("normalizeCodeBlock(%q) returned empty string for non-empty input", tt.code)
			}
		})
	}
}

func TestGetModuleNameSimple(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		expected string
	}{
		{
			name:     "Python module path",
			path:     "mypackage.submodule.module",
			expected: "module",
		},
		{
			name:     "Simple name",
			path:     "utils",
			expected: "utils",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getModuleName(tt.path)
			if result != tt.expected {
				t.Errorf("getModuleName(%q) = %q, expected %q", tt.path, result, tt.expected)
			}
		})
	}
}