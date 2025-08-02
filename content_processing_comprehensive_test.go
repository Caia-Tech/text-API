package textlib

import (
	"os"
	"path/filepath"
	"testing"
)

// Comprehensive tests for content processing functions

func TestExtractTextFromPDFComprehensive(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "pdf_test_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	tests := []struct {
		name        string
		content     []byte
		expectError bool
		expectText  bool
	}{
		{
			name:        "Valid PDF header",
			content:     []byte("%PDF-1.4\n%âãÏÓ\nSome PDF content here"),
			expectError: false,
			expectText:  true,
		},
		{
			name:        "Invalid PDF",
			content:     []byte("Not a PDF file"),
			expectError: true,
			expectText:  false,
		},
		{
			name:        "Empty file",
			content:     []byte{},
			expectError: true,
			expectText:  false,
		},
		{
			name:        "PDF with binary content",
			content:     []byte("%PDF-1.4\n\x00\x01\x02\x03"),
			expectError: false,
			expectText:  true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			filePath := filepath.Join(tmpDir, test.name+".pdf")
			err := os.WriteFile(filePath, test.content, 0644)
			if err != nil {
				t.Fatalf("Failed to create test file: %v", err)
			}

			text, err := ExtractTextFromPDF(filePath)

			if test.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if test.expectText && text == "" {
				t.Errorf("Expected some text but got empty string")
			}
		})
	}

	// Test with non-existent file
	t.Run("Non-existent file", func(t *testing.T) {
		_, err := ExtractTextFromPDF("/nonexistent/file.pdf")
		if err == nil {
			t.Errorf("Expected error for non-existent file")
		}
	})
}

func TestParseCSVStructureComprehensive(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "csv_test_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	tests := []struct {
		name        string
		content     string
		expectError bool
		expectRows  int
		expectCols  int
	}{
		{
			name: "Simple CSV",
			content: `name,age,city
John,25,New York
Jane,30,Los Angeles
Bob,35,Chicago`,
			expectError: false,
			expectRows:  4, // Including header
			expectCols:  3,
		},
		{
			name: "CSV with quoted fields",
			content: `"Name","Age","Description"
"John Doe",25,"Software Engineer, Tech Company"
"Jane Smith",30,"Data Scientist, ""AI"" Company"`,
			expectError: false,
			expectRows:  3,
			expectCols:  3,
		},
		{
			name: "CSV with empty fields",
			content: `name,age,city
John,,New York
,30,
Bob,35,Chicago`,
			expectError: false,
			expectRows:  4,
			expectCols:  3,
		},
		{
			name: "Single column CSV",
			content: `items
apple
banana
orange`,
			expectError: false,
			expectRows:  4,
			expectCols:  1,
		},
		{
			name: "Empty CSV",
			content: ``,
			expectError: true,
			expectRows:  0,
			expectCols:  0,
		},
		{
			name: "CSV with different delimiters",
			content: `name;age;city
John;25;New York
Jane;30;Los Angeles`,
			expectError: false,
			expectRows:  3,
			expectCols:  3, // Should detect semicolon delimiter
		},
		{
			name: "Malformed CSV",
			content: `name,age,city
John,25,"Unclosed quote
Jane,30,Los Angeles`,
			expectError: true,
			expectRows:  0,
			expectCols:  0,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			filePath := filepath.Join(tmpDir, test.name+".csv")
			err := os.WriteFile(filePath, []byte(test.content), 0644)
			if err != nil {
				t.Fatalf("Failed to create test file: %v", err)
			}

			schema, err := ParseCSVStructure(filePath)

			if test.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if schema.RowCount != test.expectRows {
				t.Errorf("Expected %d rows, got %d", test.expectRows, schema.RowCount)
			}

			if test.expectRows > 0 && len(schema.Headers) != test.expectCols {
				t.Errorf("Expected %d columns, got %d", test.expectCols, len(schema.Headers))
			}

			// Check that column types are detected
			if len(schema.Headers) > 0 && len(schema.ColumnTypes) == 0 {
				t.Errorf("Expected column types to be detected")
			}
		})
	}
}

func TestAnalyzeImagePropertiesComprehensive(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "image_test_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	tests := []struct {
		name        string
		filename    string
		content     []byte
		expectError bool
		expectType  string
	}{
		{
			name:        "JPEG image",
			filename:    "test.jpg",
			content:     []byte("\xFF\xD8\xFF\xE0\x00\x10JFIF\x00\x01\x01\x01\x00H\x00H"),
			expectError: false,
			expectType:  "JPEG",
		},
		{
			name:        "PNG image",
			filename:    "test.png",
			content:     []byte("\x89PNG\r\n\x1a\n\x00\x00\x00\rIHDR\x00\x00\x00\x10\x00\x00\x00\x10"),
			expectError: false,
			expectType:  "PNG",
		},
		{
			name:        "GIF image",
			filename:    "test.gif",
			content:     []byte("GIF89a\x10\x00\x10\x00"),
			expectError: false,
			expectType:  "GIF",
		},
		{
			name:        "BMP image",
			filename:    "test.bmp",
			content:     []byte("BM\x36\x00\x00\x00\x00\x00\x00\x00\x36\x00\x00\x00"),
			expectError: false,
			expectType:  "BMP",
		},
		{
			name:        "Not an image",
			filename:    "test.txt",
			content:     []byte("This is not an image file"),
			expectError: true,
			expectType:  "",
		},
		{
			name:        "Empty file",
			filename:    "empty.jpg",
			content:     []byte{},
			expectError: true,
			expectType:  "",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			filePath := filepath.Join(tmpDir, test.filename)
			err := os.WriteFile(filePath, test.content, 0644)
			if err != nil {
				t.Fatalf("Failed to create test file: %v", err)
			}

			metadata, err := AnalyzeImageProperties(filePath)

			if test.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if metadata.Format != test.expectType {
				t.Errorf("Expected format %s, got %s", test.expectType, metadata.Format)
			}

			// Check that basic metadata fields are populated
			if metadata.Width == 0 && test.expectType != "" {
				t.Errorf("Expected non-zero width for valid image")
			}

			if metadata.Height == 0 && test.expectType != "" {
				t.Errorf("Expected non-zero height for valid image")
			}
		})
	}
}

func TestExtractAudioMetadataComprehensive(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "audio_test_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	tests := []struct {
		name        string
		filename    string
		content     []byte
		expectError bool
		expectType  string
	}{
		{
			name:        "MP3 file",
			filename:    "test.mp3",
			content:     []byte("ID3\x03\x00\x00\x00\x00\x00\x00\xFF\xFB\x90\x00"),
			expectError: false,
			expectType:  "MP3",
		},
		{
			name:        "WAV file",
			filename:    "test.wav",
			content:     []byte("RIFF\x24\x00\x00\x00WAVEfmt \x10\x00\x00\x00"),
			expectError: false,
			expectType:  "WAV",
		},
		{
			name:        "FLAC file",
			filename:    "test.flac",
			content:     []byte("fLaC\x00\x00\x00\x22"),
			expectError: false,
			expectType:  "FLAC",
		},
		{
			name:        "OGG file",
			filename:    "test.ogg",
			content:     []byte("OggS\x00\x02\x00\x00\x00\x00\x00\x00\x00\x00"),
			expectError: false,
			expectType:  "OGG",
		},
		{
			name:        "Not an audio file",
			filename:    "test.txt",
			content:     []byte("This is not an audio file"),
			expectError: true,
			expectType:  "",
		},
		{
			name:        "Empty file",
			filename:    "empty.mp3",
			content:     []byte{},
			expectError: true,
			expectType:  "",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			filePath := filepath.Join(tmpDir, test.filename)
			err := os.WriteFile(filePath, test.content, 0644)
			if err != nil {
				t.Fatalf("Failed to create test file: %v", err)
			}

			metadata, err := ExtractAudioMetadata(filePath)

			if test.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if metadata.Format != test.expectType {
				t.Errorf("Expected format %s, got %s", test.expectType, metadata.Format)
			}

			// Check that basic metadata fields are present
			if metadata.SampleRate == 0 && test.expectType != "" {
				// Some formats might not have sample rate in minimal headers
				t.Logf("Sample rate is 0 for %s", test.expectType)
			}
		})
	}
}

func TestParseLogFileStructureComprehensive(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "log_test_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	tests := []struct {
		name           string
		content        string
		expectError    bool
		expectEntries  int
		expectErrorCnt int
		expectWarnCnt  int
		expectInfoCnt  int
	}{
		{
			name: "Standard log format",
			content: `2024-01-01 10:00:00 INFO Application started
2024-01-01 10:01:00 WARN Configuration file not found, using defaults
2024-01-01 10:02:00 ERROR Failed to connect to database
2024-01-01 10:03:00 INFO Connection restored
2024-01-01 10:04:00 DEBUG Processing request #123`,
			expectError:    false,
			expectEntries:  5,
			expectErrorCnt: 1,
			expectWarnCnt:  1,
			expectInfoCnt:  2,
		},
		{
			name: "Apache access log format",
			content: `127.0.0.1 - - [01/Jan/2024:10:00:00 +0000] "GET /index.html HTTP/1.1" 200 1234
127.0.0.1 - - [01/Jan/2024:10:01:00 +0000] "POST /api/data HTTP/1.1" 500 567
127.0.0.1 - - [01/Jan/2024:10:02:00 +0000] "GET /favicon.ico HTTP/1.1" 404 0`,
			expectError:   false,
			expectEntries: 3,
		},
		{
			name: "JSON log format",
			content: `{"timestamp":"2024-01-01T10:00:00Z","level":"info","message":"App started"}
{"timestamp":"2024-01-01T10:01:00Z","level":"error","message":"DB error"}
{"timestamp":"2024-01-01T10:02:00Z","level":"warn","message":"High memory usage"}`,
			expectError:    false,
			expectEntries:  3,
			expectErrorCnt: 1,
			expectWarnCnt:  1,
			expectInfoCnt:  1,
		},
		{
			name: "Mixed format log",
			content: `INFO: Application starting up
[2024-01-01 10:00:00] ERROR: Database connection failed
WARNING - Memory usage high
DEBUG: Processing completed
FATAL System crash detected`,
			expectError:    false,
			expectEntries:  5,
			expectErrorCnt: 1,
			expectWarnCnt:  1,
			expectInfoCnt:  1,
		},
		{
			name: "Empty log file",
			content: ``,
			expectError:   false,
			expectEntries: 0,
		},
		{
			name: "Single line log",
			content: `2024-01-01 10:00:00 INFO Single log entry`,
			expectError:   false,
			expectEntries: 1,
			expectInfoCnt: 1,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			filePath := filepath.Join(tmpDir, test.name+".log")
			err := os.WriteFile(filePath, []byte(test.content), 0644)
			if err != nil {
				t.Fatalf("Failed to create test file: %v", err)
			}

			schema, err := ParseLogFileStructure(filePath)

			if test.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if schema.LineCount != test.expectEntries {
				t.Errorf("Expected %d entries, got %d", test.expectEntries, schema.LineCount)
			}

			// Check that log levels are detected if expected
			if test.expectErrorCnt > 0 || test.expectWarnCnt > 0 || test.expectInfoCnt > 0 {
				if len(schema.LogLevel) == 0 {
					t.Errorf("Expected log levels to be detected")
				}
			}

			// Check time range for non-empty logs
			if test.expectEntries > 0 {
				if schema.DateRange.Start.IsZero() && schema.DateRange.End.IsZero() {
					t.Logf("No date range detected for %s format", schema.Format)
				}

				if !schema.DateRange.Start.IsZero() && !schema.DateRange.End.IsZero() && 
				   schema.DateRange.Start.After(schema.DateRange.End) {
					t.Errorf("Start time should not be after end time")
				}
			}
		})
	}
}

func TestValidateJSONStructureComprehensive(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "json_test_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	tests := []struct {
		name        string
		content     string
		expectValid bool
		expectError bool
	}{
		{
			name: "Valid JSON object",
			content: `{
				"name": "John",
				"age": 30,
				"active": true,
				"address": {
					"street": "123 Main St",
					"city": "New York"
				}
			}`,
			expectValid: true,
			expectError: false,
		},
		{
			name: "Valid JSON array",
			content: `[
				{"id": 1, "name": "Item 1"},
				{"id": 2, "name": "Item 2"},
				{"id": 3, "name": "Item 3"}
			]`,
			expectValid: true,
			expectError: false,
		},
		{
			name:        "Valid simple JSON",
			content:     `{"message": "Hello, World!"}`,
			expectValid: true,
			expectError: false,
		},
		{
			name:        "Invalid JSON - missing quote",
			content:     `{"name: "John", "age": 30}`,
			expectValid: false,
			expectError: false,
		},
		{
			name:        "Invalid JSON - trailing comma",
			content:     `{"name": "John", "age": 30,}`,
			expectValid: false,
			expectError: false,
		},
		{
			name:        "Invalid JSON - unescaped characters",
			content:     `{"message": "Hello\nWorld"}`,
			expectValid: false,
			expectError: false,
		},
		{
			name:        "Empty JSON object",
			content:     `{}`,
			expectValid: true,
			expectError: false,
		},
		{
			name:        "Empty JSON array",
			content:     `[]`,
			expectValid: true,
			expectError: false,
		},
		{
			name:        "Empty file",
			content:     ``,
			expectValid: false,
			expectError: false,
		},
		{
			name:        "Not JSON at all",
			content:     `This is not JSON content`,
			expectValid: false,
			expectError: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			filePath := filepath.Join(tmpDir, test.name+".json")
			err := os.WriteFile(filePath, []byte(test.content), 0644)
			if err != nil {
				t.Fatalf("Failed to create test file: %v", err)
			}

			errors, err := ValidateJSONStructure(filePath)

			if test.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			isValid := len(errors) == 0
			if isValid != test.expectValid {
				t.Errorf("Expected valid %v, got %v", test.expectValid, isValid)
				if !isValid && len(errors) > 0 {
					t.Logf("Validation errors: %v", errors)
				}
			}

			// For valid JSON, no errors should be present
			if test.expectValid && len(errors) > 0 {
				t.Errorf("Expected no errors for valid JSON, got: %v", errors)
			}

			// For invalid JSON, check that errors are reported
			if !test.expectValid && len(errors) == 0 {
				t.Errorf("Expected validation errors for invalid JSON")
			}
		})
	}
}

func TestCheckXMLWellFormednessComprehensive(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "xml_test_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	tests := []struct {
		name        string
		content     string
		expectValid bool
		expectError bool
	}{
		{
			name: "Valid XML document",
			content: `<?xml version="1.0" encoding="UTF-8"?>
<root>
	<person id="1">
		<name>John Doe</name>
		<age>30</age>
		<email>john@example.com</email>
	</person>
	<person id="2">
		<name>Jane Smith</name>
		<age>25</age>
		<email>jane@example.com</email>
	</person>
</root>`,
			expectValid: true,
			expectError: false,
		},
		{
			name: "Valid simple XML",
			content: `<?xml version="1.0"?>
<message>Hello, World!</message>`,
			expectValid: true,
			expectError: false,
		},
		{
			name: "XML with namespaces",
			content: `<?xml version="1.0"?>
<root xmlns:ns="http://example.com/namespace">
	<ns:element>Content</ns:element>
</root>`,
			expectValid: true,
			expectError: false,
		},
		{
			name: "XML with CDATA",
			content: `<?xml version="1.0"?>
<root>
	<script><![CDATA[
		function test() {
			return "Hello";
		}
	]]></script>
</root>`,
			expectValid: true,
			expectError: false,
		},
		{
			name:        "Invalid XML - unclosed tag",
			content:     `<root><item>Content</root>`,
			expectValid: false,
			expectError: false,
		},
		{
			name:        "Invalid XML - malformed tag",
			content:     `<root><item>Content</item<>/root>`,
			expectValid: false,
			expectError: false,
		},
		{
			name:        "Invalid XML - missing closing tag",
			content:     `<root><item>Content</item>`,
			expectValid: false,
			expectError: false,
		},
		{
			name:        "Empty XML",
			content:     ``,
			expectValid: false,
			expectError: false,
		},
		{
			name:        "Not XML",
			content:     `This is not XML content`,
			expectValid: false,
			expectError: false,
		},
		{
			name:        "Well-formed but no declaration",
			content:     `<root><item>Content</item></root>`,
			expectValid: true,
			expectError: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			filePath := filepath.Join(tmpDir, test.name+".xml")
			err := os.WriteFile(filePath, []byte(test.content), 0644)
			if err != nil {
				t.Fatalf("Failed to create test file: %v", err)
			}

			errors, err := CheckXMLWellFormedness(filePath)

			if test.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			isValid := len(errors) == 0
			if isValid != test.expectValid {
				t.Errorf("Expected valid %v, got %v", test.expectValid, isValid)
				if !isValid && len(errors) > 0 {
					t.Logf("Validation errors: %v", errors)
				}
			}

			// For invalid XML, check that errors are reported
			if !test.expectValid && len(errors) == 0 {
				t.Errorf("Expected validation errors for invalid XML")
			}
		})
	}
}

// Test error conditions
func TestContentProcessingErrorConditions(t *testing.T) {
	t.Run("Non-existent file handling", func(t *testing.T) {
		nonExistentFile := "/nonexistent/path/file.txt"

		// Test all functions with non-existent files
		_, err := ExtractTextFromPDF(nonExistentFile)
		if err == nil {
			t.Errorf("ExtractTextFromPDF should error on non-existent file")
		}

		_, err = ParseCSVStructure(nonExistentFile)
		if err == nil {
			t.Errorf("ParseCSVStructure should error on non-existent file")
		}

		_, err = AnalyzeImageProperties(nonExistentFile)
		if err == nil {
			t.Errorf("AnalyzeImageProperties should error on non-existent file")
		}

		_, err = ExtractAudioMetadata(nonExistentFile)
		if err == nil {
			t.Errorf("ExtractAudioMetadata should error on non-existent file")
		}

		_, err = ParseLogFileStructure(nonExistentFile)
		if err == nil {
			t.Errorf("ParseLogFileStructure should error on non-existent file")
		}

		_, err = ValidateJSONStructure(nonExistentFile)
		if err == nil {
			t.Errorf("ValidateJSONStructure should error on non-existent file")
		}

		_, err = CheckXMLWellFormedness(nonExistentFile)
		if err == nil {
			t.Errorf("CheckXMLWellFormedness should error on non-existent file")
		}
	})

	t.Run("Permission denied handling", func(t *testing.T) {
		// Create a file and remove read permissions
		tmpFile, err := os.CreateTemp("", "permission_test_*.txt")
		if err != nil {
			t.Fatalf("Failed to create temp file: %v", err)
		}
		defer os.Remove(tmpFile.Name())
		tmpFile.Close()

		// Remove read permissions
		err = os.Chmod(tmpFile.Name(), 0000)
		if err != nil {
			t.Fatalf("Failed to change permissions: %v", err)
		}
		defer os.Chmod(tmpFile.Name(), 0644) // Restore for cleanup

		// Test functions should handle permission errors gracefully
		_, err = ParseCSVStructure(tmpFile.Name())
		if err == nil {
			t.Errorf("ParseCSVStructure should error on permission denied")
		}
	})
}

// Benchmark tests
func BenchmarkParseCSVStructure(b *testing.B) {
	tmpFile, err := os.CreateTemp("", "csv_benchmark_*.csv")
	if err != nil {
		b.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	// Create a CSV with many rows
	content := "name,age,city,occupation\n"
	for i := 0; i < 1000; i++ {
		content += "Person" + string(rune(i)) + ",25,City" + string(rune(i)) + ",Job" + string(rune(i)) + "\n"
	}
	tmpFile.WriteString(content)
	tmpFile.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ParseCSVStructure(tmpFile.Name())
	}
}

func BenchmarkParseLogFileStructure(b *testing.B) {
	tmpFile, err := os.CreateTemp("", "log_benchmark_*.log")
	if err != nil {
		b.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	// Create a log file with many entries
	content := ""
	levels := []string{"INFO", "WARN", "ERROR", "DEBUG"}
	for i := 0; i < 1000; i++ {
		level := levels[i%len(levels)]
		content += "2024-01-01 10:00:00 " + level + " Log entry number " + string(rune(i)) + "\n"
	}
	tmpFile.WriteString(content)
	tmpFile.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ParseLogFileStructure(tmpFile.Name())
	}
}

func BenchmarkValidateJSONStructure(b *testing.B) {
	tmpFile, err := os.CreateTemp("", "json_benchmark_*.json")
	if err != nil {
		b.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	// Create a large JSON file
	content := `{"data": [`
	for i := 0; i < 1000; i++ {
		if i > 0 {
			content += ","
		}
		content += `{"id": ` + string(rune(i)) + `, "name": "Item` + string(rune(i)) + `"}`
	}
	content += `]}`
	tmpFile.WriteString(content)
	tmpFile.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ValidateJSONStructure(tmpFile.Name())
	}
}