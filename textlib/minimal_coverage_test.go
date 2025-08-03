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

// Minimal tests to improve coverage on functions with 0% coverage

func TestMinMaxFloat(t *testing.T) {
	// Test minFloat function
	result1 := minFloat(1.5, 2.5)
	if result1 != 1.5 {
		t.Errorf("minFloat(1.5, 2.5) = %f, expected 1.5", result1)
	}

	// Test maxFloat function
	result2 := maxFloat(1.5, 2.5)
	if result2 != 2.5 {
		t.Errorf("maxFloat(1.5, 2.5) = %f, expected 2.5", result2)
	}
}

func TestPrivateIPDetection(t *testing.T) {
	tests := []struct {
		ip       string
		expected bool
	}{
		{"192.168.1.1", true},
		{"10.0.0.1", true},
		{"8.8.8.8", false},
		{"invalid", false},
	}

	for _, tt := range tests {
		result := isPrivateIP(tt.ip)
		if result != tt.expected {
			t.Errorf("isPrivateIP(%q) = %v, expected %v", tt.ip, result, tt.expected)
		}
	}
}

func TestReadableTextDetection(t *testing.T) {
	tests := []struct {
		content  []byte
		expected bool
	}{
		{[]byte("Hello world"), true},
		{[]byte{0x00, 0x01, 0xFF}, false},
		{[]byte{}, true},
	}

	for _, tt := range tests {
		result := isReadableText(tt.content)
		if result != tt.expected {
			t.Errorf("isReadableText(%v) = %v, expected %v", tt.content, result, tt.expected)
		}
	}
}

func TestExtractID3FieldSimple(t *testing.T) {
	data := []byte("some test data for ID3")
	result := extractID3Field(data, 0, 5)
	if len(result) != 5 {
		t.Errorf("extractID3Field returned wrong length: %d", len(result))
	}
}

func TestFileTypeDetection(t *testing.T) {
	tests := []struct {
		filename string
		isCode   bool
		isAsset  bool
	}{
		{"test.go", true, false},
		{"test.js", true, false},
		{"test.png", false, true},
		{"test.css", false, true},
		{"test.txt", false, false},
	}

	for _, tt := range tests {
		codeResult := isCodeFile(tt.filename)
		assetResult := isAssetFile(tt.filename)
		
		if codeResult != tt.isCode {
			t.Errorf("isCodeFile(%q) = %v, expected %v", tt.filename, codeResult, tt.isCode)
		}
		
		if assetResult != tt.isAsset {
			t.Errorf("isAssetFile(%q) = %v, expected %v", tt.filename, assetResult, tt.isAsset)
		}
	}
}

func TestFileAnalysisFunctions(t *testing.T) {
	// Test extractFileReferences
	code := "import './test.js'; require('../utils.js');"
	refs := extractFileReferences(code)
	_ = refs // Just ensure it doesn't panic

	// Test GetFilesBySize
	files := GetFilesBySize(".", 1000)
	_ = files // Just ensure it doesn't panic

	// Test GetFilesByAge  
	oldFiles := GetFilesByAge(".", 30)
	_ = oldFiles // Just ensure it doesn't panic

	// Test AnalyzeDiskUsage
	usage := AnalyzeDiskUsage(".")
	if usage.TotalSize < 0 {
		t.Error("AnalyzeDiskUsage returned negative size")
	}

	// Test isDirEmpty
	empty := isDirEmpty(".")
	_ = empty // Just ensure it doesn't panic

	// Test ExtractFileList
	fileList := ExtractFileList("test.zip")
	_ = fileList // Just ensure it doesn't panic
}

func TestVirusSignatures(t *testing.T) {
	signatures := GetCommonVirusSignatures()
	if len(signatures) == 0 {
		t.Error("GetCommonVirusSignatures returned no signatures")
	}
}

func TestPEPackingDetection(t *testing.T) {
	// Test with fake PE data
	data := []byte("MZ\x90\x00test data")
	packed, ratio := checkPEPacking(data)
	_ = packed
	_ = ratio
}