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

// Tests for basic functions to improve coverage safely

func TestMinFloat(t *testing.T) {
	tests := []struct {
		name     string
		a, b     float64
		expected float64
	}{
		{
			name:     "First smaller",
			a:        1.5,
			b:        2.5,
			expected: 1.5,
		},
		{
			name:     "Second smaller",
			a:        3.5,
			b:        2.5,
			expected: 2.5,
		},
		{
			name:     "Equal values",
			a:        2.5,
			b:        2.5,
			expected: 2.5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := minFloat(tt.a, tt.b)
			if result != tt.expected {
				t.Errorf("minFloat(%f, %f) = %f, expected %f", tt.a, tt.b, result, tt.expected)
			}
		})
	}
}

func TestMaxFloat(t *testing.T) {
	tests := []struct {
		name     string
		a, b     float64
		expected float64
	}{
		{
			name:     "First larger",
			a:        2.5,
			b:        1.5,
			expected: 2.5,
		},
		{
			name:     "Second larger",
			a:        1.5,
			b:        3.5,
			expected: 3.5,
		},
		{
			name:     "Equal values",
			a:        2.5,
			b:        2.5,
			expected: 2.5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := maxFloat(tt.a, tt.b)
			if result != tt.expected {
				t.Errorf("maxFloat(%f, %f) = %f, expected %f", tt.a, tt.b, result, tt.expected)
			}
		})
	}
}

func TestIsReadableText(t *testing.T) {
	tests := []struct {
		name     string
		content  []byte
		expected bool
	}{
		{
			name:     "Plain text",
			content:  []byte("Hello, this is readable text."),
			expected: true,
		},
		{
			name:     "Binary data",
			content:  []byte{0x00, 0x01, 0x02, 0xFF, 0xFE},
			expected: false,
		},
		{
			name:     "Empty content",
			content:  []byte{},
			expected: true, // Empty is considered readable
		},
		{
			name:     "Text with some control chars",
			content:  []byte("Hello\nWorld\t!"),
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isReadableText(tt.content)
			if result != tt.expected {
				t.Errorf("isReadableText(%q) = %v, expected %v", string(tt.content), result, tt.expected)
			}
		})
	}
}

func TestExtractID3Field(t *testing.T) {
	tests := []struct {
		name   string
		data   []byte
		offset int
		length int
	}{
		{
			name:   "Valid ID3 data",
			data:   []byte("ID3\x03\x00\x00\x00\x00\x00\x00TITLE"),
			offset: 10,
			length: 5,
		},
		{
			name:   "Empty data",
			data:   []byte{},
			offset: 0,
			length: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extractID3Field(tt.data, tt.offset, tt.length)
			// Just test that it doesn't panic
			_ = result
		})
	}
}

func TestIsPrivateIP(t *testing.T) {
	tests := []struct {
		name     string
		ip       string
		expected bool
	}{
		{
			name:     "Private IP 192.168.x.x",
			ip:       "192.168.1.1",
			expected: true,
		},
		{
			name:     "Private IP 10.x.x.x",
			ip:       "10.0.0.1",
			expected: true,
		},
		{
			name:     "Private IP 172.16-31.x.x",
			ip:       "172.16.0.1",
			expected: true,
		},
		{
			name:     "Public IP",
			ip:       "8.8.8.8",
			expected: false,
		},
		{
			name:     "Localhost",
			ip:       "127.0.0.1",
			expected: true,
		},
		{
			name:     "Invalid IP",
			ip:       "not.an.ip",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isPrivateIP(tt.ip)
			if result != tt.expected {
				t.Errorf("isPrivateIP(%q) = %v, expected %v", tt.ip, result, tt.expected)
			}
		})
	}
}

func TestCheckPEPacking(t *testing.T) {
	tests := []struct {
		name string
		data []byte
	}{
		{
			name: "PE header",
			data: []byte("MZ\x90\x00\x03\x00\x00\x00\x04\x00\x00\x00\xFF\xFF\x00\x00"),
		},
		{
			name: "Non-PE data",
			data: []byte("Not a PE file"),
		},
		{
			name: "Empty data",
			data: []byte{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			packed, ratio := checkPEPacking(tt.data)
			// Just test that it doesn't panic
			_ = packed
			_ = ratio
		})
	}
}

func TestGetCommonVirusSignatures(t *testing.T) {
	t.Run("Get virus signatures", func(t *testing.T) {
		signatures := GetCommonVirusSignatures()
		// Should return some signatures
		if len(signatures) == 0 {
			t.Error("GetCommonVirusSignatures() returned empty slice")
		}
	})
}

func TestIsDirEmpty(t *testing.T) {
	tests := []struct {
		name string
		path string
	}{
		{
			name: "Current directory",
			path: ".",
		},
		{
			name: "Non-existent directory",
			path: "/path/that/does/not/exist",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isDirEmpty(tt.path)
			// Just test that it doesn't panic
			_ = result
		})
	}
}

func TestIsCodeFile(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		expected bool
	}{
		{
			name:     "Go file",
			filename: "main.go",
			expected: true,
		},
		{
			name:     "JavaScript file",
			filename: "script.js",
			expected: true,
		},
		{
			name:     "Python file",
			filename: "script.py",
			expected: true,
		},
		{
			name:     "Text file",
			filename: "readme.txt",
			expected: false,
		},
		{
			name:     "Image file",
			filename: "image.png",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isCodeFile(tt.filename)
			if result != tt.expected {
				t.Errorf("isCodeFile(%q) = %v, expected %v", tt.filename, result, tt.expected)
			}
		})
	}
}

func TestIsAssetFile(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		expected bool
	}{
		{
			name:     "Image file",
			filename: "image.png",
			expected: true,
		},
		{
			name:     "CSS file",
			filename: "style.css",
			expected: true,
		},
		{
			name:     "JavaScript file",
			filename: "script.js",
			expected: false, // JS is code, not asset
		},
		{
			name:     "Font file",
			filename: "font.ttf",
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isAssetFile(tt.filename)
			if result != tt.expected {
				t.Errorf("isAssetFile(%q) = %v, expected %v", tt.filename, result, tt.expected)
			}
		})
	}
}

func TestExtractFileReferences(t *testing.T) {
	tests := []struct {
		name string
		code string
	}{
		{
			name: "Code with file references",
			code: `import "./module.js";
require("../utils.js");
#include "header.h"`,
		},
		{
			name: "Code without file references",
			code: "var x = 5;\nfunction test() { return x; }",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extractFileReferences(tt.code)
			// Just test that it doesn't panic
			_ = result
		})
	}
}

func TestGetFilesBySize(t *testing.T) {
	t.Run("Get files by size", func(t *testing.T) {
		result := GetFilesBySize(".", 1000) // Files larger than 1KB
		// Just test that it doesn't panic
		_ = result
	})
}

func TestGetFilesByAge(t *testing.T) {
	t.Run("Get files by age", func(t *testing.T) {
		result := GetFilesByAge(".", 30) // Files older than 30 days
		// Just test that it doesn't panic
		_ = result
	})
}

func TestAnalyzeDiskUsage(t *testing.T) {
	t.Run("Analyze disk usage", func(t *testing.T) {
		result := AnalyzeDiskUsage(".")
		// Just test that it doesn't panic and returns something reasonable
		if result.TotalSize < 0 {
			t.Error("AnalyzeDiskUsage returned negative total size")
		}
	})
}

func TestExtractFileList(t *testing.T) {
	tests := []struct {
		name     string
		filename string
	}{
		{
			name:     "ZIP file extension",
			filename: "archive.zip",
		},
		{
			name:     "TAR file extension",
			filename: "archive.tar",
		},
		{
			name:     "Regular file",
			filename: "document.txt",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ExtractFileList(tt.filename)
			// Just test that it doesn't panic
			_ = result
		})
	}
}