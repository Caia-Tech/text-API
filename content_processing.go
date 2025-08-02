package textlib

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// Content processing data structures

type CSVSchema struct {
	Headers      []string
	ColumnTypes  map[string]string
	RowCount     int
	HasHeaders   bool
	Delimiter    string
	QuoteChar    string
	Encoding     string
	DataSample   [][]string
}

type ImageMetadata struct {
	Width       int
	Height      int
	BitDepth    int
	ColorType   string
	Format      string
	FileSize    int64
	Compression string
	HasAlpha    bool
	DPI         int
	ColorSpace  string
	Created     time.Time
	Camera      CameraInfo
}

type CameraInfo struct {
	Make     string
	Model    string
	ISO      int
	Aperture string
	Shutter  string
	Flash    bool
	GPS      GPSInfo
}

type GPSInfo struct {
	Latitude  float64
	Longitude float64
	Altitude  float64
	HasGPS    bool
}

type AudioMetadata struct {
	Duration   time.Duration
	Bitrate    int
	SampleRate int
	Channels   int
	Format     string
	Codec      string
	FileSize   int64
	Title      string
	Artist     string
	Album      string
	Year       int
	Genre      string
	TrackNum   int
}

type LogSchema struct {
	Format      string
	TimeFormat  string
	Fields      []LogField
	SampleLines []string
	LineCount   int
	DateRange   DateRange
	LogLevel    []string
	IsStructured bool
}

type LogField struct {
	Name     string
	Type     string
	Pattern  string
	Required bool
	Position int
}

type DateRange struct {
	Start time.Time
	End   time.Time
}

type ValidationError struct {
	Line    int
	Column  int
	Message string
	Code    string
	Value   string
}

type XMLError struct {
	Line    int
	Column  int
	Message string
	Type    string
}

type CodecInfo struct {
	VideoCodec   string
	AudioCodec   string
	Container    string
	Resolution   string
	Framerate    float64
	Bitrate      int
	Duration     time.Duration
	HasSubtitles bool
	IsValid      bool
}

// Data extraction functions

func ExtractTextFromPDF(filePath string) (string, error) {
	// This would typically require a PDF library like github.com/ledongthuc/pdf
	// For now, provide a basic implementation that checks for PDF structure
	
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()
	
	// Check if it's a PDF file
	header := make([]byte, 4)
	_, err = file.Read(header)
	if err != nil {
		return "", err
	}
	
	if string(header) != "%PDF" {
		return "", fmt.Errorf("not a valid PDF file")
	}
	
	// Read the entire file for basic text extraction
	file.Seek(0, 0)
	content, err := io.ReadAll(file)
	if err != nil {
		return "", err
	}
	
	// Very basic text extraction - look for text between parentheses or brackets
	text := extractBasicPDFText(string(content))
	
	return text, nil
}

func extractBasicPDFText(content string) string {
	var extractedText strings.Builder
	
	// Look for text in PDF streams (very basic)
	textPattern := regexp.MustCompile(`\(([^)]+)\)`)
	matches := textPattern.FindAllStringSubmatch(content, -1)
	
	for _, match := range matches {
		if len(match) > 1 {
			// Clean up the text
			text := strings.ReplaceAll(match[1], "\\n", "\n")
			text = strings.ReplaceAll(text, "\\r", "\r")
			text = strings.ReplaceAll(text, "\\t", "\t")
			extractedText.WriteString(text + " ")
		}
	}
	
	// Also look for text in alternative formats
	altTextPattern := regexp.MustCompile(`\[([^\]]+)\]`)
	altMatches := altTextPattern.FindAllStringSubmatch(content, -1)
	
	for _, match := range altMatches {
		if len(match) > 1 && isReadableText(match[1]) {
			extractedText.WriteString(match[1] + " ")
		}
	}
	
	return strings.TrimSpace(extractedText.String())
}

func isReadableText(text string) bool {
	// Check if text contains mostly printable characters
	printableCount := 0
	for _, r := range text {
		if r >= 32 && r <= 126 {
			printableCount++
		}
	}
	return float64(printableCount)/float64(len(text)) > 0.7
}

func ParseCSVStructure(filePath string) (CSVSchema, error) {
	schema := CSVSchema{
		ColumnTypes: make(map[string]string),
		DataSample:  make([][]string, 0),
	}
	
	file, err := os.Open(filePath)
	if err != nil {
		return schema, err
	}
	defer file.Close()
	
	// Detect delimiter by reading first few lines
	scanner := bufio.NewScanner(file)
	var firstLines []string
	for i := 0; i < 5 && scanner.Scan(); i++ {
		firstLines = append(firstLines, scanner.Text())
	}
	
	if len(firstLines) == 0 {
		return schema, fmt.Errorf("empty file")
	}
	
	// Detect delimiter
	schema.Delimiter = detectCSVDelimiter(firstLines)
	schema.QuoteChar = `"`
	
	// Reset file pointer
	file.Seek(0, 0)
	reader := csv.NewReader(file)
	reader.Comma = rune(schema.Delimiter[0])
	
	// Read all records
	records, err := reader.ReadAll()
	if err != nil {
		return schema, err
	}
	
	schema.RowCount = len(records)
	
	if len(records) == 0 {
		return schema, nil
	}
	
	// Determine if first row contains headers
	schema.HasHeaders = detectCSVHeaders(records)
	
	if schema.HasHeaders && len(records) > 0 {
		schema.Headers = records[0]
		records = records[1:]
	} else {
		// Generate column names
		for i := 0; i < len(records[0]); i++ {
			schema.Headers = append(schema.Headers, fmt.Sprintf("Column%d", i+1))
		}
	}
	
	// Analyze column types
	for i, header := range schema.Headers {
		colType := detectColumnType(records, i)
		schema.ColumnTypes[header] = colType
	}
	
	// Store sample data (first 5 rows)
	for i := 0; i < len(records) && i < 5; i++ {
		schema.DataSample = append(schema.DataSample, records[i])
	}
	
	schema.Encoding = "UTF-8" // Default assumption
	
	return schema, nil
}

func detectCSVDelimiter(lines []string) string {
	delimiters := []string{",", ";", "\t", "|"}
	scores := make(map[string]int)
	
	for _, line := range lines {
		for _, delim := range delimiters {
			count := strings.Count(line, delim)
			scores[delim] += count
		}
	}
	
	// Find delimiter with highest score
	maxScore := 0
	bestDelim := ","
	for delim, score := range scores {
		if score > maxScore {
			maxScore = score
			bestDelim = delim
		}
	}
	
	return bestDelim
}

func detectCSVHeaders(records [][]string) bool {
	if len(records) < 2 {
		return false
	}
	
	firstRow := records[0]
	secondRow := records[1]
	
	// Check if first row contains non-numeric values while second row contains numbers
	headerScore := 0
	for i := 0; i < len(firstRow) && i < len(secondRow); i++ {
		if _, err := strconv.ParseFloat(firstRow[i], 64); err != nil {
			headerScore++
		}
		if _, err := strconv.ParseFloat(secondRow[i], 64); err == nil {
			headerScore++
		}
	}
	
	return headerScore > len(firstRow)/2
}

func detectColumnType(records [][]string, colIndex int) string {
	if len(records) == 0 {
		return "unknown"
	}
	
	intCount := 0
	floatCount := 0
	dateCount := 0
	boolCount := 0
	totalCount := 0
	
	for _, record := range records {
		if colIndex >= len(record) {
			continue
		}
		
		value := strings.TrimSpace(record[colIndex])
		if value == "" {
			continue
		}
		
		totalCount++
		
		// Check if integer
		if _, err := strconv.Atoi(value); err == nil {
			intCount++
			continue
		}
		
		// Check if float
		if _, err := strconv.ParseFloat(value, 64); err == nil {
			floatCount++
			continue
		}
		
		// Check if boolean
		lower := strings.ToLower(value)
		if lower == "true" || lower == "false" || lower == "yes" || lower == "no" || lower == "1" || lower == "0" {
			boolCount++
			continue
		}
		
		// Check if date
		if isDateString(value) {
			dateCount++
			continue
		}
	}
	
	if totalCount == 0 {
		return "empty"
	}
	
	// Determine type based on majority
	threshold := totalCount / 2
	
	if intCount > threshold {
		return "integer"
	}
	if floatCount > threshold {
		return "float"
	}
	if dateCount > threshold {
		return "date"
	}
	if boolCount > threshold {
		return "boolean"
	}
	
	return "string"
}

func isDateString(value string) bool {
	dateFormats := []string{
		"2006-01-02",
		"01/02/2006",
		"2006-01-02 15:04:05",
		"01/02/2006 15:04:05",
		"2006-01-02T15:04:05Z",
		"Jan 2, 2006",
		"January 2, 2006",
	}
	
	for _, format := range dateFormats {
		if _, err := time.Parse(format, value); err == nil {
			return true
		}
	}
	
	return false
}

func AnalyzeImageProperties(imagePath string) (ImageMetadata, error) {
	metadata := ImageMetadata{}
	
	file, err := os.Open(imagePath)
	if err != nil {
		return metadata, err
	}
	defer file.Close()
	
	info, err := file.Stat()
	if err != nil {
		return metadata, err
	}
	metadata.FileSize = info.Size()
	metadata.Created = info.ModTime()
	
	// Read file header to determine format
	header := make([]byte, 20)
	_, err = file.Read(header)
	if err != nil {
		return metadata, err
	}
	
	// Detect image format and extract basic properties
	if len(header) >= 8 && string(header[:8]) == "\x89PNG\r\n\x1a\n" {
		metadata.Format = "PNG"
		return analyzePNGProperties(file, metadata)
	} else if len(header) >= 2 && header[0] == 0xFF && header[1] == 0xD8 {
		metadata.Format = "JPEG"
		return analyzeJPEGProperties(file, metadata)
	} else if len(header) >= 6 && string(header[:6]) == "GIF87a" || string(header[:6]) == "GIF89a" {
		metadata.Format = "GIF"
		return analyzeGIFProperties(file, metadata)
	}
	
	return metadata, fmt.Errorf("unsupported image format")
}

func analyzePNGProperties(file *os.File, metadata ImageMetadata) (ImageMetadata, error) {
	// PNG analysis would require detailed chunk parsing
	// This is a simplified version
	file.Seek(16, 0) // Skip PNG signature and IHDR chunk header
	
	// Read IHDR data
	ihdr := make([]byte, 13)
	_, err := file.Read(ihdr)
	if err != nil {
		return metadata, err
	}
	
	// Extract width and height (big-endian)
	metadata.Width = int(ihdr[0])<<24 | int(ihdr[1])<<16 | int(ihdr[2])<<8 | int(ihdr[3])
	metadata.Height = int(ihdr[4])<<24 | int(ihdr[5])<<16 | int(ihdr[6])<<8 | int(ihdr[7])
	metadata.BitDepth = int(ihdr[8])
	
	colorType := ihdr[9]
	switch colorType {
	case 0:
		metadata.ColorType = "Grayscale"
	case 2:
		metadata.ColorType = "RGB"
	case 3:
		metadata.ColorType = "Palette"
	case 4:
		metadata.ColorType = "Grayscale+Alpha"
		metadata.HasAlpha = true
	case 6:
		metadata.ColorType = "RGBA"
		metadata.HasAlpha = true
	}
	
	metadata.Compression = "Deflate"
	metadata.ColorSpace = "sRGB"
	
	return metadata, nil
}

func analyzeJPEGProperties(file *os.File, metadata ImageMetadata) (ImageMetadata, error) {
	// JPEG analysis requires parsing markers
	// This is a simplified version
	metadata.ColorType = "RGB"
	metadata.Compression = "JPEG"
	metadata.ColorSpace = "YCbCr"
	
	// For full JPEG analysis, would need to parse EXIF data
	// This is a placeholder
	metadata.Width = 0  // Would be extracted from SOF marker
	metadata.Height = 0 // Would be extracted from SOF marker
	
	return metadata, nil
}

func analyzeGIFProperties(file *os.File, metadata ImageMetadata) (ImageMetadata, error) {
	// GIF analysis
	file.Seek(6, 0) // Skip GIF header
	
	// Read logical screen descriptor
	lsd := make([]byte, 7)
	_, err := file.Read(lsd)
	if err != nil {
		return metadata, err
	}
	
	metadata.Width = int(lsd[0]) | int(lsd[1])<<8
	metadata.Height = int(lsd[2]) | int(lsd[3])<<8
	metadata.ColorType = "Palette"
	metadata.Compression = "LZW"
	
	return metadata, nil
}

func ExtractAudioMetadata(audioPath string) (AudioMetadata, error) {
	metadata := AudioMetadata{}
	
	info, err := os.Stat(audioPath)
	if err != nil {
		return metadata, err
	}
	metadata.FileSize = info.Size()
	
	// Detect format from extension
	ext := strings.ToLower(filepath.Ext(audioPath))
	switch ext {
	case ".mp3":
		metadata.Format = "MP3"
		metadata.Codec = "MPEG-1 Layer 3"
		return extractMP3Metadata(audioPath, metadata)
	case ".wav":
		metadata.Format = "WAV"
		metadata.Codec = "PCM"
		return extractWAVMetadata(audioPath, metadata)
	case ".flac":
		metadata.Format = "FLAC"
		metadata.Codec = "FLAC"
		return extractFLACMetadata(audioPath, metadata)
	case ".ogg":
		metadata.Format = "OGG"
		metadata.Codec = "Vorbis"
		return extractOGGMetadata(audioPath, metadata)
	default:
		return metadata, fmt.Errorf("unsupported audio format: %s", ext)
	}
}

func extractMP3Metadata(audioPath string, metadata AudioMetadata) (AudioMetadata, error) {
	file, err := os.Open(audioPath)
	if err != nil {
		return metadata, err
	}
	defer file.Close()
	
	// Look for ID3v2 tag at beginning
	header := make([]byte, 10)
	_, err = file.Read(header)
	if err != nil {
		return metadata, err
	}
	
	if string(header[:3]) == "ID3" {
		// Parse ID3v2 tag
		tagSize := int(header[6])<<21 | int(header[7])<<14 | int(header[8])<<7 | int(header[9])
		
		// Read ID3 tag data
		tagData := make([]byte, tagSize)
		_, err = file.Read(tagData)
		if err == nil {
			metadata = parseID3Tag(tagData, metadata)
		}
	}
	
	// Basic MP3 properties (would need frame parsing for accuracy)
	metadata.Bitrate = 128   // Default assumption
	metadata.SampleRate = 44100
	metadata.Channels = 2
	
	return metadata, nil
}

func parseID3Tag(tagData []byte, metadata AudioMetadata) AudioMetadata {
	// Simplified ID3 parsing - would need full implementation for production
	content := string(tagData)
	
	// Look for common fields (this is very basic)
	if titlePos := strings.Index(content, "TIT2"); titlePos != -1 {
		// Extract title (simplified)
		metadata.Title = extractID3Field(content, titlePos)
	}
	
	if artistPos := strings.Index(content, "TPE1"); artistPos != -1 {
		metadata.Artist = extractID3Field(content, artistPos)
	}
	
	if albumPos := strings.Index(content, "TALB"); albumPos != -1 {
		metadata.Album = extractID3Field(content, albumPos)
	}
	
	return metadata
}

func extractID3Field(content string, pos int) string {
	// Very basic field extraction - production would need proper frame parsing
	start := pos + 10 // Skip frame header
	if start >= len(content) {
		return ""
	}
	
	end := start
	for end < len(content) && content[end] != 0 && end < start+100 {
		end++
	}
	
	if end > start {
		field := content[start:end]
		// Clean up the field
		if len(field) > 0 && field[0] == 1 {
			field = field[1:] // Skip encoding byte
		}
		return strings.TrimSpace(field)
	}
	
	return ""
}

func extractWAVMetadata(audioPath string, metadata AudioMetadata) (AudioMetadata, error) {
	file, err := os.Open(audioPath)
	if err != nil {
		return metadata, err
	}
	defer file.Close()
	
	// Read WAV header
	header := make([]byte, 44)
	_, err = file.Read(header)
	if err != nil {
		return metadata, err
	}
	
	// Check RIFF header
	if string(header[:4]) != "RIFF" || string(header[8:12]) != "WAVE" {
		return metadata, fmt.Errorf("invalid WAV file")
	}
	
	// Extract format information
	if string(header[12:16]) == "fmt " {
		metadata.SampleRate = int(header[24]) | int(header[25])<<8 | int(header[26])<<16 | int(header[27])<<24
		metadata.Channels = int(header[22]) | int(header[23])<<8
		metadata.Bitrate = int(header[28]) | int(header[29])<<8 | int(header[30])<<16 | int(header[31])<<24
	}
	
	return metadata, nil
}

func extractFLACMetadata(audioPath string, metadata AudioMetadata) (AudioMetadata, error) {
	// FLAC metadata extraction would require specialized parsing
	// This is a placeholder
	metadata.Bitrate = 0 // FLAC is lossless
	metadata.SampleRate = 44100
	metadata.Channels = 2
	return metadata, nil
}

func extractOGGMetadata(audioPath string, metadata AudioMetadata) (AudioMetadata, error) {
	// OGG metadata extraction would require specialized parsing
	// This is a placeholder
	metadata.Bitrate = 160
	metadata.SampleRate = 44100
	metadata.Channels = 2
	return metadata, nil
}

func ParseLogFileStructure(logPath string) (LogSchema, error) {
	schema := LogSchema{
		Fields:      make([]LogField, 0),
		SampleLines: make([]string, 0),
		LogLevel:    make([]string, 0),
	}
	
	file, err := os.Open(logPath)
	if err != nil {
		return schema, err
	}
	defer file.Close()
	
	scanner := bufio.NewScanner(file)
	lineCount := 0
	var lines []string
	
	// Read first 100 lines for analysis
	for scanner.Scan() && lineCount < 100 {
		line := scanner.Text()
		lines = append(lines, line)
		lineCount++
		
		if len(schema.SampleLines) < 5 {
			schema.SampleLines = append(schema.SampleLines, line)
		}
	}
	
	// Count total lines
	for scanner.Scan() {
		lineCount++
	}
	schema.LineCount = lineCount
	
	if len(lines) == 0 {
		return schema, nil
	}
	
	// Detect log format
	schema.Format = detectLogFormat(lines)
	schema.TimeFormat = detectTimeFormat(lines)
	schema.IsStructured = detectStructuredFormat(lines)
	
	// Extract log levels
	schema.LogLevel = extractLogLevels(lines)
	
	// Detect fields based on format
	schema.Fields = detectLogFields(lines, schema.Format)
	
	// Extract date range
	schema.DateRange = extractDateRange(lines, schema.TimeFormat)
	
	return schema, nil
}

func detectLogFormat(lines []string) string {
	if len(lines) == 0 {
		return "unknown"
	}
	
	// Check for common log formats
	firstLine := lines[0]
	
	// Apache Common Log Format
	if matched, _ := regexp.MatchString(`^\d+\.\d+\.\d+\.\d+ - - \[`, firstLine); matched {
		return "apache_common"
	}
	
	// Apache Combined Log Format
	if matched, _ := regexp.MatchString(`^\d+\.\d+\.\d+\.\d+ - - \[.*\] ".*" \d+ \d+ ".*" ".*"`, firstLine); matched {
		return "apache_combined"
	}
	
	// Syslog format
	if matched, _ := regexp.MatchString(`^[A-Za-z]{3}\s+\d{1,2}\s+\d{2}:\d{2}:\d{2}`, firstLine); matched {
		return "syslog"
	}
	
	// JSON format
	if strings.HasPrefix(strings.TrimSpace(firstLine), "{") {
		return "json"
	}
	
	// CSV format
	if strings.Contains(firstLine, ",") && strings.Count(firstLine, ",") > 2 {
		return "csv"
	}
	
	// Generic timestamp format
	if matched, _ := regexp.MatchString(`^\d{4}-\d{2}-\d{2}`, firstLine); matched {
		return "timestamp_generic"
	}
	
	return "custom"
}

func detectTimeFormat(lines []string) string {
	timeFormats := []string{
		"2006-01-02 15:04:05",
		"2006/01/02 15:04:05",
		"Jan 02 15:04:05",
		"2006-01-02T15:04:05Z",
		"02/Jan/2006:15:04:05 -0700",
	}
	
	for _, line := range lines {
		for _, format := range timeFormats {
			// Try to find timestamp in line
			if matched, _ := regexp.MatchString(regexp.QuoteMeta(format[:10]), line); matched {
				return format
			}
		}
	}
	
	return "unknown"
}

func detectStructuredFormat(lines []string) bool {
	if len(lines) == 0 {
		return false
	}
	
	// Check if it's JSON
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "{") && strings.HasSuffix(trimmed, "}") {
			var js json.RawMessage
			if json.Unmarshal([]byte(trimmed), &js) == nil {
				return true
			}
		}
	}
	
	// Check if it's consistently delimited
	delimiters := []string{"\t", "|", ","}
	for _, delim := range delimiters {
		consistent := true
		firstCount := strings.Count(lines[0], delim)
		if firstCount == 0 {
			continue
		}
		
		for _, line := range lines[1:5] { // Check first few lines
			if strings.Count(line, delim) != firstCount {
				consistent = false
				break
			}
		}
		
		if consistent {
			return true
		}
	}
	
	return false
}

func extractLogLevels(lines []string) []string {
	levels := make(map[string]bool)
	levelPattern := regexp.MustCompile(`(?i)\b(DEBUG|INFO|WARN|WARNING|ERROR|FATAL|TRACE|CRITICAL)\b`)
	
	for _, line := range lines {
		matches := levelPattern.FindAllString(line, -1)
		for _, match := range matches {
			levels[strings.ToUpper(match)] = true
		}
	}
	
	var result []string
	for level := range levels {
		result = append(result, level)
	}
	
	return result
}

func detectLogFields(lines []string, format string) []LogField {
	fields := []LogField{}
	
	switch format {
	case "apache_common":
		fields = []LogField{
			{Name: "ip", Type: "string", Pattern: `\d+\.\d+\.\d+\.\d+`, Required: true, Position: 0},
			{Name: "timestamp", Type: "datetime", Pattern: `\[.*\]`, Required: true, Position: 1},
			{Name: "method", Type: "string", Pattern: `"[A-Z]+ `, Required: true, Position: 2},
			{Name: "status", Type: "integer", Pattern: `" \d{3} `, Required: true, Position: 3},
			{Name: "size", Type: "integer", Pattern: ` \d+$`, Required: false, Position: 4},
		}
	case "json":
		// For JSON, would need to parse actual structure
		fields = append(fields, LogField{Name: "json_object", Type: "object", Required: true, Position: 0})
	case "syslog":
		fields = []LogField{
			{Name: "timestamp", Type: "datetime", Pattern: `^[A-Za-z]{3}\s+\d{1,2}\s+\d{2}:\d{2}:\d{2}`, Required: true, Position: 0},
			{Name: "hostname", Type: "string", Required: true, Position: 1},
			{Name: "service", Type: "string", Required: true, Position: 2},
			{Name: "message", Type: "string", Required: true, Position: 3},
		}
	default:
		// Generic fields
		fields = append(fields, LogField{Name: "line", Type: "string", Required: true, Position: 0})
	}
	
	return fields
}

func extractDateRange(lines []string, timeFormat string) DateRange {
	dateRange := DateRange{}
	
	if timeFormat == "unknown" || len(lines) == 0 {
		return dateRange
	}
	
	// Try to parse first and last timestamps
	timePattern := regexp.MustCompile(`\d{4}-\d{2}-\d{2}[ T]\d{2}:\d{2}:\d{2}`)
	
	// Find first timestamp
	for _, line := range lines {
		if timestamp := timePattern.FindString(line); timestamp != "" {
			if t, err := time.Parse("2006-01-02 15:04:05", timestamp); err == nil {
				dateRange.Start = t
				break
			}
		}
	}
	
	// Find last timestamp (check from end)
	for i := len(lines) - 1; i >= 0; i-- {
		if timestamp := timePattern.FindString(lines[i]); timestamp != "" {
			if t, err := time.Parse("2006-01-02 15:04:05", timestamp); err == nil {
				dateRange.End = t
				break
			}
		}
	}
	
	return dateRange
}

// File validation functions

func ValidateJSONStructure(jsonPath string) ([]ValidationError, error) {
	var errors []ValidationError
	
	file, err := os.Open(jsonPath)
	if err != nil {
		return errors, err
	}
	defer file.Close()
	
	content, err := io.ReadAll(file)
	if err != nil {
		return errors, err
	}
	
	var js json.RawMessage
	err = json.Unmarshal(content, &js)
	if err != nil {
		// Parse JSON error
		if syntaxErr, ok := err.(*json.SyntaxError); ok {
			line, col := getLineColumn(content, syntaxErr.Offset)
			validationErr := ValidationError{
				Line:    line,
				Column:  col,
				Message: syntaxErr.Error(),
				Code:    "SYNTAX_ERROR",
			}
			errors = append(errors, validationErr)
		} else {
			validationErr := ValidationError{
				Line:    1,
				Column:  1,
				Message: err.Error(),
				Code:    "PARSE_ERROR",
			}
			errors = append(errors, validationErr)
		}
	}
	
	return errors, nil
}

func getLineColumn(content []byte, offset int64) (int, int) {
	line := 1
	col := 1
	
	for i := int64(0); i < offset && i < int64(len(content)); i++ {
		if content[i] == '\n' {
			line++
			col = 1
		} else {
			col++
		}
	}
	
	return line, col
}

func CheckXMLWellFormedness(xmlPath string) ([]XMLError, error) {
	var errors []XMLError
	
	file, err := os.Open(xmlPath)
	if err != nil {
		return errors, err
	}
	defer file.Close()
	
	decoder := xml.NewDecoder(file)
	
	for {
		_, err := decoder.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			if xmlErr, ok := err.(*xml.SyntaxError); ok {
				error := XMLError{
					Line:    int(xmlErr.Line),
					Message: xmlErr.Msg,
					Type:    "SYNTAX_ERROR",
				}
				errors = append(errors, error)
			}
			// Break on any error to prevent infinite loops
			break
		}
	}
	
	return errors, nil
}

func VerifyImageIntegrity(imagePath string) bool {
	metadata, err := AnalyzeImageProperties(imagePath)
	if err != nil {
		return false
	}
	
	// Basic integrity checks
	if metadata.Width <= 0 || metadata.Height <= 0 {
		return false
	}
	
	if metadata.FileSize <= 0 {
		return false
	}
	
	// Format-specific checks would go here
	return true
}

func ValidateVideoCodec(videoPath string) (CodecInfo, error) {
	codec := CodecInfo{}
	
	// This would typically require a media library like ffmpeg
	// For now, provide basic format detection
	
	ext := strings.ToLower(filepath.Ext(videoPath))
	switch ext {
	case ".mp4":
		codec.Container = "MP4"
		codec.VideoCodec = "H.264"
		codec.AudioCodec = "AAC"
		codec.IsValid = true
	case ".avi":
		codec.Container = "AVI"
		codec.VideoCodec = "DivX/XviD"
		codec.AudioCodec = "MP3"
		codec.IsValid = true
	case ".mkv":
		codec.Container = "Matroska"
		codec.VideoCodec = "H.264"
		codec.AudioCodec = "AAC"
		codec.IsValid = true
	case ".mov":
		codec.Container = "QuickTime"
		codec.VideoCodec = "H.264"
		codec.AudioCodec = "AAC"
		codec.IsValid = true
	default:
		codec.IsValid = false
		return codec, fmt.Errorf("unsupported video format: %s", ext)
	}
	
	// Get file info
	info, err := os.Stat(videoPath)
	if err != nil {
		codec.IsValid = false
		return codec, err
	}
	
	// Basic file validation
	if info.Size() == 0 {
		codec.IsValid = false
		return codec, fmt.Errorf("empty video file")
	}
	
	return codec, nil
}