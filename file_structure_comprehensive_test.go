package textlib

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// Comprehensive tests for file structure analysis functions

func TestExtractMetadataComprehensive(t *testing.T) {
	// Create test files with different properties
	tmpDir, err := os.MkdirTemp("", "metadata_test_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	tests := []struct {
		name        string
		filename    string
		content     string
		expectedExt string
		isDir       bool
	}{
		{
			name:        "Text file",
			filename:    "test.txt",
			content:     "Hello, world!",
			expectedExt: ".txt",
			isDir:       false,
		},
		{
			name:        "JSON file",
			filename:    "data.json",
			content:     `{"key": "value"}`,
			expectedExt: ".json",
			isDir:       false,
		},
		{
			name:        "No extension",
			filename:    "README",
			content:     "This is a README file",
			expectedExt: "",
			isDir:       false,
		},
		{
			name:        "Directory",
			filename:    "testdir",
			content:     "",
			expectedExt: "",
			isDir:       true,
		},
		{
			name:        "Go source file",
			filename:    "main.go",
			content:     "package main\n\nfunc main() {}\n",
			expectedExt: ".go",
			isDir:       false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var filePath string
			
			if test.isDir {
				filePath = filepath.Join(tmpDir, test.filename)
				err := os.Mkdir(filePath, 0755)
				if err != nil {
					t.Fatalf("Failed to create directory: %v", err)
				}
			} else {
				filePath = filepath.Join(tmpDir, test.filename)
				err := os.WriteFile(filePath, []byte(test.content), 0644)
				if err != nil {
					t.Fatalf("Failed to create file: %v", err)
				}
			}

			metadata, err := ExtractMetadata(filePath)
			if err != nil {
				t.Errorf("ExtractMetadata failed: %v", err)
				return
			}

			if metadata.Name != test.filename {
				t.Errorf("Expected name %s, got %s", test.filename, metadata.Name)
			}

			if metadata.Extension != test.expectedExt {
				t.Errorf("Expected extension %s, got %s", test.expectedExt, metadata.Extension)
			}

			if metadata.IsDirectory != test.isDir {
				t.Errorf("Expected isDirectory %v, got %v", test.isDir, metadata.IsDirectory)
			}

			if !test.isDir && metadata.Size != int64(len(test.content)) {
				t.Errorf("Expected size %d, got %d", len(test.content), metadata.Size)
			}

			if metadata.Path != filePath {
				t.Errorf("Expected path %s, got %s", filePath, metadata.Path)
			}
		})
	}
}

func TestDetectFileTypeComprehensive(t *testing.T) {
	tests := []struct {
		name         string
		content      []byte
		expectedMime string
		expectedCat  string
		expectedBin  bool
		expectedText bool
	}{
		{
			name:         "PDF file",
			content:      []byte("%PDF-1.4\n%âãÏÓ"),
			expectedMime: "application/pdf",
			expectedCat:  "document",
			expectedBin:  true,
			expectedText: false,
		},
		{
			name:         "PNG image",
			content:      []byte("\x89PNG\r\n\x1a\n\x00\x00\x00\rIHDR"),
			expectedMime: "image/png",
			expectedCat:  "image",
			expectedBin:  true,
			expectedText: false,
		},
		{
			name:         "JPEG image",
			content:      []byte("\xFF\xD8\xFF\xE0\x00\x10JFIF"),
			expectedMime: "image/jpeg",
			expectedCat:  "image",
			expectedBin:  true,
			expectedText: false,
		},
		{
			name:         "ZIP archive",
			content:      []byte("PK\x03\x04\x14\x00\x00\x00"),
			expectedMime: "application/zip",
			expectedCat:  "archive",
			expectedBin:  true,
			expectedText: false,
		},
		{
			name:         "ELF executable",
			content:      []byte("\x7fELF\x02\x01\x01\x00"),
			expectedMime: "application/x-executable",
			expectedCat:  "executable",
			expectedBin:  true,
			expectedText: false,
		},
		{
			name:         "Windows PE executable",
			content:      []byte("MZ\x90\x00\x03\x00\x00\x00"),
			expectedMime: "application/x-msdownload",
			expectedCat:  "executable",
			expectedBin:  true,
			expectedText: false,
		},
		{
			name:         "Plain text",
			content:      []byte("Hello, world! This is plain text."),
			expectedMime: "text/plain",
			expectedCat:  "text",
			expectedBin:  false,
			expectedText: true,
		},
		{
			name:         "HTML document",
			content:      []byte("<!DOCTYPE html>\n<html><head><title>Test</title></head></html>"),
			expectedMime: "text/html",
			expectedCat:  "text",
			expectedBin:  false,
			expectedText: true,
		},
		{
			name:         "JSON data",
			content:      []byte(`{"name": "test", "value": 123, "array": [1, 2, 3]}`),
			expectedMime: "application/json",
			expectedCat:  "text",
			expectedBin:  false,
			expectedText: true,
		},
		{
			name:         "XML document",
			content:      []byte(`<?xml version="1.0" encoding="UTF-8"?>\n<root><item>test</item></root>`),
			expectedMime: "application/xml",
			expectedCat:  "text",
			expectedBin:  false,
			expectedText: true,
		},
		{
			name:         "Binary data with null bytes",
			content:      []byte("\x00\x01\x02\x03\x04\x05\x06\x07"),
			expectedMime: "application/octet-stream",
			expectedCat:  "binary",
			expectedBin:  true,
			expectedText: false,
		},
		{
			name:         "Empty file",
			content:      []byte{},
			expectedMime: "",
			expectedCat:  "",
			expectedBin:  false,
			expectedText: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			fileType := DetectFileType(test.content)

			if fileType.MimeType != test.expectedMime {
				t.Errorf("Expected mime type %s, got %s", test.expectedMime, fileType.MimeType)
			}

			if fileType.Category != test.expectedCat {
				t.Errorf("Expected category %s, got %s", test.expectedCat, fileType.Category)
			}

			if fileType.IsBinary != test.expectedBin {
				t.Errorf("Expected binary %v, got %v", test.expectedBin, fileType.IsBinary)
			}

			if fileType.IsText != test.expectedText {
				t.Errorf("Expected text %v, got %v", test.expectedText, fileType.IsText)
			}
		})
	}
}

func TestAnalyzeFileStructureComprehensive(t *testing.T) {
	// Create a complex directory structure
	tmpDir, err := os.MkdirTemp("", "structure_test_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create directory structure
	dirs := []string{
		"level1",
		"level1/level2",
		"level1/level2/level3",
		"images",
		"documents",
	}

	for _, dir := range dirs {
		err := os.MkdirAll(filepath.Join(tmpDir, dir), 0755)
		if err != nil {
			t.Fatalf("Failed to create directory %s: %v", dir, err)
		}
	}

	// Create files with different sizes and types
	files := []struct {
		path    string
		content string
	}{
		{"file1.txt", "Small file"},
		{"file2.log", strings.Repeat("a", 1000)}, // 1KB
		{"file3.json", `{"large": "` + strings.Repeat("x", 5000) + `"}`}, // ~5KB
		{"level1/nested.go", "package main\n\nfunc main() {}"},
		{"level1/level2/deep.py", "print('hello from deep level')"},
		{"level1/level2/level3/deepest.txt", "This is the deepest file"},
		{"images/photo.jpg", "fake-jpeg-content"},
		{"documents/report.pdf", "fake-pdf-content"},
	}

	for _, file := range files {
		filePath := filepath.Join(tmpDir, file.path)
		err := os.WriteFile(filePath, []byte(file.content), 0644)
		if err != nil {
			t.Fatalf("Failed to create file %s: %v", file.path, err)
		}
	}

	report, err := AnalyzeFileStructure(tmpDir)
	if err != nil {
		t.Errorf("AnalyzeFileStructure failed: %v", err)
		return
	}

	// Check basic counts
	if report.TotalFiles != len(files) {
		t.Errorf("Expected %d files, got %d", len(files), report.TotalFiles)
	}

	if report.TotalDirectories != len(dirs)+1 { // +1 for tmpDir itself
		t.Errorf("Expected %d directories, got %d", len(dirs)+1, report.TotalDirectories)
	}

	// Check depth calculation
	if report.MaxDepth < 3 {
		t.Errorf("Expected max depth at least 3, got %d", report.MaxDepth)
	}

	// Check file types
	expectedExtensions := []string{".txt", ".log", ".json", ".go", ".py", ".jpg", ".pdf"}
	for _, ext := range expectedExtensions {
		if report.FileTypes[ext] == 0 {
			t.Errorf("Expected to find files with extension %s", ext)
		}
	}

	// Check total size is reasonable
	if report.TotalSize == 0 {
		t.Errorf("Expected non-zero total size")
	}

	// Check average file size
	if report.AverageFileSize == 0 {
		t.Errorf("Expected non-zero average file size")
	}

	// Check largest files tracking
	if len(report.LargestFiles) == 0 {
		t.Errorf("Expected some largest files to be tracked")
	}
}

func TestCalculateChecksumComprehensive(t *testing.T) {
	// Create test file
	tmpFile, err := os.CreateTemp("", "checksum_test_*.txt")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	content := "The quick brown fox jumps over the lazy dog"
	tmpFile.WriteString(content)
	tmpFile.Close()

	tests := []struct {
		algorithm    string
		expectedLen  int
		shouldError  bool
	}{
		{"md5", 32, false},
		{"sha1", 40, false},
		{"sha256", 64, false},
		{"MD5", 32, false},    // Test case insensitive
		{"SHA256", 64, false}, // Test case insensitive
		{"invalid", 0, true},  // Unsupported algorithm
	}

	for _, test := range tests {
		t.Run(test.algorithm, func(t *testing.T) {
			checksum, err := CalculateChecksum(tmpFile.Name(), test.algorithm)

			if test.shouldError {
				if err == nil {
					t.Errorf("Expected error for algorithm %s", test.algorithm)
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error for algorithm %s: %v", test.algorithm, err)
				return
			}

			if len(checksum) != test.expectedLen {
				t.Errorf("Expected checksum length %d for %s, got %d", 
					test.expectedLen, test.algorithm, len(checksum))
			}

			// Verify checksum is hexadecimal
			for _, char := range checksum {
				if !((char >= '0' && char <= '9') || (char >= 'a' && char <= 'f')) {
					t.Errorf("Checksum contains non-hex character: %c", char)
					break
				}
			}
		})
	}

	// Test with non-existent file
	_, err = CalculateChecksum("/nonexistent/file.txt", "md5")
	if err == nil {
		t.Errorf("Expected error for non-existent file")
	}
}

func TestValidateFileIntegrityComprehensive(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "integrity_test_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	tests := []struct {
		name        string
		filename    string
		content     []byte
		expectValid bool
		expectError bool
		expectWarn  bool
	}{
		{
			name:        "Valid text file",
			filename:    "valid.txt",
			content:     []byte("This is a valid text file"),
			expectValid: true,
			expectError: false,
			expectWarn:  false,
		},
		{
			name:        "Empty file",
			filename:    "empty.txt",
			content:     []byte{},
			expectValid: true,
			expectError: false,
			expectWarn:  true, // Empty file warning
		},
		{
			name:        "Text file with null bytes",
			filename:    "corrupted.txt",
			content:     []byte("Hello\x00world\x00test"),
			expectValid: true,
			expectError: false,
			expectWarn:  true, // Null byte warning
		},
		{
			name:        "Binary file",
			filename:    "binary.bin",
			content:     []byte{0x00, 0x01, 0x02, 0x03, 0xFF, 0xFE},
			expectValid: true,
			expectError: false,
			expectWarn:  false,
		},
		{
			name:        "JSON file",
			filename:    "data.json",
			content:     []byte(`{"valid": "json", "number": 123}`),
			expectValid: true,
			expectError: false,
			expectWarn:  false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			filePath := filepath.Join(tmpDir, test.filename)
			err := os.WriteFile(filePath, test.content, 0644)
			if err != nil {
				t.Fatalf("Failed to create test file: %v", err)
			}

			report := ValidateFileIntegrity(filePath)

			if report.IsValid != test.expectValid {
				t.Errorf("Expected valid %v, got %v", test.expectValid, report.IsValid)
			}

			if test.expectError && len(report.Errors) == 0 {
				t.Errorf("Expected errors but got none")
			}

			if !test.expectError && len(report.Errors) > 0 {
				t.Errorf("Expected no errors but got: %v", report.Errors)
			}

			if test.expectWarn && len(report.Warnings) == 0 {
				t.Errorf("Expected warnings but got none")
			}

			if report.FileSize != int64(len(test.content)) {
				t.Errorf("Expected file size %d, got %d", len(test.content), report.FileSize)
			}

			if !report.LastModified.IsZero() {
				// Check that modification time is recent (within last minute)
				if time.Since(report.LastModified) > time.Minute {
					t.Errorf("File modification time seems too old: %v", report.LastModified)
				}
			}
		})
	}

	// Test with non-existent file
	t.Run("Non-existent file", func(t *testing.T) {
		report := ValidateFileIntegrity("/nonexistent/file.txt")
		if report.IsValid {
			t.Errorf("Expected invalid for non-existent file")
		}
		if len(report.Errors) == 0 {
			t.Errorf("Expected error for non-existent file")
		}
	})
}

func TestFindDuplicateFilesComprehensive(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "duplicates_test_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create duplicate files
	testContent := []struct {
		name    string
		content string
	}{
		{"file1.txt", "identical content"},
		{"file2.txt", "identical content"}, // Duplicate of file1
		{"file3.txt", "identical content"}, // Another duplicate
		{"unique1.txt", "unique content 1"},
		{"unique2.txt", "unique content 2"},
		{"empty1.txt", ""},
		{"empty2.txt", ""}, // Duplicate empty file
	}

	for _, file := range testContent {
		filePath := filepath.Join(tmpDir, file.name)
		err := os.WriteFile(filePath, []byte(file.content), 0644)
		if err != nil {
			t.Fatalf("Failed to create test file %s: %v", file.name, err)
		}
	}

	duplicates, err := FindDuplicateFiles(tmpDir)
	if err != nil {
		t.Errorf("FindDuplicateFiles failed: %v", err)
		return
	}

	// Should find duplicate sets
	expectedSets := 1 // One set with "identical content" (3 files)
	// Note: empty files are ignored (size 0)

	if len(duplicates) != expectedSets {
		t.Errorf("Expected %d duplicate sets, got %d", expectedSets, len(duplicates))
		for i, dup := range duplicates {
			t.Logf("Duplicate set %d: %d files, size %d", i, dup.Count, dup.Size)
			for _, file := range dup.Files {
				t.Logf("  - %s", file)
			}
		}
		return
	}

	// Check the duplicate set
	if len(duplicates) > 0 {
		dupSet := duplicates[0]
		if dupSet.Count != 3 {
			t.Errorf("Expected 3 duplicate files, got %d", dupSet.Count)
		}

		if dupSet.Size != int64(len("identical content")) {
			t.Errorf("Expected size %d, got %d", len("identical content"), dupSet.Size)
		}

		if dupSet.Checksum == "" {
			t.Errorf("Expected non-empty checksum")
		}

		// Verify all files in the set exist
		for _, filePath := range dupSet.Files {
			if _, err := os.Stat(filePath); os.IsNotExist(err) {
				t.Errorf("Duplicate file does not exist: %s", filePath)
			}
		}
	}
}

func TestArchiveAnalysisComprehensive(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "archive_test_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	tests := []struct {
		name        string
		filename    string
		content     []byte
		expectType  string
		expectError bool
	}{
		{
			name:        "ZIP file",
			filename:    "test.zip",
			content:     []byte("PK\x03\x04\x14\x00\x00\x00"),
			expectType:  "ZIP",
			expectError: false,
		},
		{
			name:        "TAR file",
			filename:    "test.tar",
			content:     make([]byte, 1024), // TAR needs to be multiple of 512
			expectType:  "TAR",
			expectError: false,
		},
		{
			name:        "GZIP file",
			filename:    "test.gz",
			content:     []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00"),
			expectType:  "GZIP",
			expectError: false,
		},
		{
			name:        "TAR.GZ file",
			filename:    "test.tar.gz",
			content:     []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00"),
			expectType:  "TAR.GZ",
			expectError: false,
		},
		{
			name:        "7Z file",
			filename:    "test.7z",
			content:     []byte("7z\xBC\xAF\x27\x1C"),
			expectType:  "7Z",
			expectError: false,
		},
		{
			name:        "Unsupported format",
			filename:    "test.rar",
			content:     []byte("Rar!\x1a\x07\x00"),
			expectType:  "",
			expectError: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			filePath := filepath.Join(tmpDir, test.filename)
			err := os.WriteFile(filePath, test.content, 0644)
			if err != nil {
				t.Fatalf("Failed to create test file: %v", err)
			}

			report, err := AnalyzeArchiveStructure(filePath)

			if test.expectError {
				if err == nil {
					t.Errorf("Expected error for %s", test.filename)
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error for %s: %v", test.filename, err)
				return
			}

			if report.Type != test.expectType {
				t.Errorf("Expected type %s, got %s", test.expectType, report.Type)
			}

			if report.CompressedSize != int64(len(test.content)) {
				t.Errorf("Expected compressed size %d, got %d", 
					len(test.content), report.CompressedSize)
			}
		})
	}
}

func TestDetectCompressionRatioComprehensive(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "compression_test_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	tests := []struct {
		filename     string
		expectedRatio float64
	}{
		{"test.zip", 0.6},
		{"test.gz", 0.3},
		{"test.7z", 0.25},
		{"test.tar", 1.0},
		{"test.txt", 1.0}, // No compression
	}

	for _, test := range tests {
		t.Run(test.filename, func(t *testing.T) {
			filePath := filepath.Join(tmpDir, test.filename)
			err := os.WriteFile(filePath, []byte("test content"), 0644)
			if err != nil {
				t.Fatalf("Failed to create test file: %v", err)
			}

			ratio, err := DetectCompressionRatio(filePath)
			if err != nil {
				t.Errorf("DetectCompressionRatio failed: %v", err)
				return
			}

			if ratio != test.expectedRatio {
				t.Errorf("Expected ratio %.2f, got %.2f", test.expectedRatio, ratio)
			}
		})
	}

	// Test with non-existent file
	t.Run("Non-existent file", func(t *testing.T) {
		_, err := DetectCompressionRatio("/nonexistent/file.zip")
		if err == nil {
			t.Errorf("Expected error for non-existent file")
		}
	})
}

// Benchmark tests
func BenchmarkExtractMetadata(b *testing.B) {
	tmpFile, err := os.CreateTemp("", "benchmark_*.txt")
	if err != nil {
		b.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	tmpFile.WriteString("benchmark content")
	tmpFile.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ExtractMetadata(tmpFile.Name())
	}
}

func BenchmarkDetectFileType(b *testing.B) {
	content := []byte("This is sample text content for benchmarking file type detection")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		DetectFileType(content)
	}
}

func BenchmarkCalculateChecksum(b *testing.B) {
	tmpFile, err := os.CreateTemp("", "checksum_benchmark_*.txt")
	if err != nil {
		b.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	// Create a larger file for meaningful benchmark
	content := strings.Repeat("benchmark content ", 1000)
	tmpFile.WriteString(content)
	tmpFile.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		CalculateChecksum(tmpFile.Name(), "md5")
	}
}