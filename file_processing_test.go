package textlib

import (
	"os"
	"testing"
	"time"
)

func TestExtractMetadata(t *testing.T) {
	// Create a temporary file for testing
	tmpFile, err := os.CreateTemp("", "test_*.txt")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()
	
	// Write some content
	content := "Hello, world!"
	tmpFile.WriteString(content)
	tmpFile.Close()
	
	metadata, err := ExtractMetadata(tmpFile.Name())
	if err != nil {
		t.Errorf("ExtractMetadata failed: %v", err)
	}
	
	if metadata.Size != int64(len(content)) {
		t.Errorf("Expected size %d, got %d", len(content), metadata.Size)
	}
	
	if metadata.Extension != ".txt" {
		t.Errorf("Expected extension .txt, got %s", metadata.Extension)
	}
	
	if metadata.IsDirectory {
		t.Errorf("Expected file, not directory")
	}
}

func TestDetectFileType(t *testing.T) {
	tests := []struct {
		content  []byte
		expected string
	}{
		{[]byte("%PDF-1.4"), "application/pdf"},
		{[]byte("\x89PNG\r\n\x1a\n"), "image/png"},
		{[]byte("\xFF\xD8\xFF"), "image/jpeg"},
		{[]byte("PK\x03\x04"), "application/zip"},
		{[]byte("Hello world"), "text/plain"},
		{[]byte("{\"key\": \"value\"}"), "application/json"},
	}
	
	for _, test := range tests {
		fileType := DetectFileType(test.content)
		if fileType.MimeType != test.expected {
			contentDesc := string(test.content)
			if len(contentDesc) > 10 {
				contentDesc = contentDesc[:10] + "..."
			}
			t.Errorf("DetectFileType(%s): expected %s, got %s", 
				contentDesc, test.expected, fileType.MimeType)
		}
	}
}

func TestCalculateChecksum(t *testing.T) {
	// Create a temporary file
	tmpFile, err := os.CreateTemp("", "checksum_test_*.txt")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()
	
	content := "Hello, world!"
	tmpFile.WriteString(content)
	tmpFile.Close()
	
	// Test different algorithms
	algorithms := []string{"md5", "sha1", "sha256"}
	
	for _, algo := range algorithms {
		checksum, err := CalculateChecksum(tmpFile.Name(), algo)
		if err != nil {
			t.Errorf("CalculateChecksum(%s) failed: %v", algo, err)
		}
		
		if checksum == "" {
			t.Errorf("CalculateChecksum(%s) returned empty checksum", algo)
		}
		
		// Verify checksum length
		expectedLengths := map[string]int{
			"md5":    32,
			"sha1":   40,
			"sha256": 64,
		}
		
		if len(checksum) != expectedLengths[algo] {
			t.Errorf("CalculateChecksum(%s): expected length %d, got %d", 
				algo, expectedLengths[algo], len(checksum))
		}
	}
}

func TestValidateFileIntegrity(t *testing.T) {
	// Test with existing file
	tmpFile, err := os.CreateTemp("", "integrity_test_*.txt")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()
	
	tmpFile.WriteString("Test content")
	tmpFile.Close()
	
	report := ValidateFileIntegrity(tmpFile.Name())
	if !report.IsValid {
		t.Errorf("Expected valid file, got invalid")
	}
	
	if report.FileSize == 0 {
		t.Errorf("Expected non-zero file size")
	}
	
	// Test with non-existing file
	report = ValidateFileIntegrity("/nonexistent/file.txt")
	if report.IsValid {
		t.Errorf("Expected invalid for non-existent file")
	}
	
	if len(report.Errors) == 0 {
		t.Errorf("Expected error for non-existent file")
	}
}

func TestCategorizeFiles(t *testing.T) {
	// Create a temporary directory
	tmpDir, err := os.MkdirTemp("", "categorize_test_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)
	
	// Create test files
	testFiles := map[string]Category{
		"document.pdf":   CategoryDocuments,
		"image.png":      CategoryImages,
		"video.mp4":      CategoryVideos,
		"audio.mp3":      CategoryAudio,
		"archive.zip":    CategoryArchives,
		"program.exe":    CategoryExecutables,
		"source.go":      CategoryCode,
		"data.json":      CategoryData,
		"config.conf":    CategoryConfig,
		"temp.tmp":       CategoryTemp,
		"unknown.xyz":    CategoryOther,
	}
	
	for filename := range testFiles {
		file, err := os.Create(tmpDir + "/" + filename)
		if err != nil {
			t.Fatalf("Failed to create test file %s: %v", filename, err)
		}
		file.Close()
	}
	
	categories, err := CategorizeFiles(tmpDir)
	if err != nil {
		t.Errorf("CategorizeFiles failed: %v", err)
	}
	
	// Check that files were categorized correctly
	for filename, expectedCategory := range testFiles {
		found := false
		for category, files := range categories {
			for _, file := range files {
				if file.Name == filename {
					if category != expectedCategory {
						t.Errorf("File %s: expected category %s, got %s", 
							filename, expectedCategory, category)
					}
					found = true
					break
				}
			}
		}
		
		if !found {
			t.Errorf("File %s not found in any category", filename)
		}
	}
}

func TestDetectMaliciousPatterns(t *testing.T) {
	// Create a test file with suspicious content
	tmpFile, err := os.CreateTemp("", "malicious_test_*.txt")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()
	
	// Write content with some suspicious patterns
	content := "cmd.exe /c echo hello"
	tmpFile.WriteString(content)
	tmpFile.Close()
	
	threats, err := DetectMaliciousPatterns(tmpFile.Name())
	if err != nil {
		t.Errorf("DetectMaliciousPatterns failed: %v", err)
	}
	
	// Should detect cmd.exe pattern
	found := false
	for _, threat := range threats {
		if threat.Pattern == "cmd.exe" {
			found = true
			break
		}
	}
	
	if !found {
		t.Errorf("Expected to detect cmd.exe pattern")
	}
}

func TestCleanupTempFiles(t *testing.T) {
	// Create a temporary directory with temp files
	tmpDir, err := os.MkdirTemp("", "cleanup_test_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)
	
	// Create various temp files
	tempFiles := []string{
		"file.tmp",
		"file.temp",
		"file.bak",
		"file.old",
		"file~",
		".DS_Store",
	}
	
	for _, filename := range tempFiles {
		file, err := os.Create(tmpDir + "/" + filename)
		if err != nil {
			t.Fatalf("Failed to create temp file %s: %v", filename, err)
		}
		file.WriteString("temp content")
		file.Close()
	}
	
	// Also create a normal file that shouldn't be deleted
	normalFile, err := os.Create(tmpDir + "/important.txt")
	if err != nil {
		t.Fatalf("Failed to create normal file: %v", err)
	}
	normalFile.WriteString("important content")
	normalFile.Close()
	
	report, err := CleanupTempFiles(tmpDir)
	if err != nil {
		t.Errorf("CleanupTempFiles failed: %v", err)
	}
	
	// Should have deleted temp files
	if len(report.DeletedFiles) == 0 {
		t.Errorf("Expected some files to be deleted")
	}
	
	// Important file should still exist
	if _, err := os.Stat(tmpDir + "/important.txt"); os.IsNotExist(err) {
		t.Errorf("Important file was incorrectly deleted")
	}
}

func TestOrganizeByDate(t *testing.T) {
	// Create test files with different dates
	files := []File{
		{
			Name:    "file1.txt",
			ModTime: time.Date(2023, 1, 15, 10, 0, 0, 0, time.UTC),
			Size:    100,
		},
		{
			Name:    "file2.txt",
			ModTime: time.Date(2023, 2, 20, 10, 0, 0, 0, time.UTC),
			Size:    200,
		},
		{
			Name:    "file3.txt",
			ModTime: time.Date(2024, 1, 10, 10, 0, 0, 0, time.UTC),
			Size:    300,
		},
	}
	
	// Test organization by year
	structure := OrganizeByDate(files, DateByYear)
	
	if len(structure.Directories) != 2 {
		t.Errorf("Expected 2 year directories, got %d", len(structure.Directories))
	}
	
	if structure.TotalFiles != 3 {
		t.Errorf("Expected 3 total files, got %d", structure.TotalFiles)
	}
	
	if structure.TotalSize != 600 {
		t.Errorf("Expected total size 600, got %d", structure.TotalSize)
	}
	
	// Check specific years
	if len(structure.Directories["2023"]) != 2 {
		t.Errorf("Expected 2 files in 2023, got %d", len(structure.Directories["2023"]))
	}
	
	if len(structure.Directories["2024"]) != 1 {
		t.Errorf("Expected 1 file in 2024, got %d", len(structure.Directories["2024"]))
	}
}

func TestDetectVersionedFiles(t *testing.T) {
	// Create a temporary directory
	tmpDir, err := os.MkdirTemp("", "version_test_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)
	
	// Create versioned files
	versionedFiles := []string{
		"document.txt",
		"document_v1.txt",
		"document_v2.txt",
		"photo.jpg",
		"photo(1).jpg",
		"photo(2).jpg",
	}
	
	for _, filename := range versionedFiles {
		file, err := os.Create(tmpDir + "/" + filename)
		if err != nil {
			t.Fatalf("Failed to create versioned file %s: %v", filename, err)
		}
		file.WriteString("content")
		file.Close()
	}
	
	versionSets, err := DetectVersionedFiles(tmpDir)
	if err != nil {
		t.Errorf("DetectVersionedFiles failed: %v", err)
	}
	
	// Should detect 2 version sets (document and photo)
	if len(versionSets) != 2 {
		t.Errorf("Expected 2 version sets, got %d", len(versionSets))
	}
	
	// Check that each set has multiple files
	for _, versionSet := range versionSets {
		if len(versionSet.Files) < 2 {
			t.Errorf("Version set %s should have at least 2 files, got %d", 
				versionSet.BaseName, len(versionSet.Files))
		}
	}
}