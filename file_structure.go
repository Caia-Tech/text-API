package textlib

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// File structure analysis data structures

type FileMetadata struct {
	Path         string
	Name         string
	Size         int64
	ModTime      time.Time
	IsDirectory  bool
	Permissions  os.FileMode
	Extension    string
	MimeType     string
	Owner        string
	Group        string
}

type StructureReport struct {
	TotalFiles      int
	TotalDirectories int
	TotalSize       int64
	FileTypes       map[string]int
	LargestFiles    []FileMetadata
	DeepestPath     string
	MaxDepth        int
	AverageFileSize float64
}

type FileType struct {
	Extension string
	MimeType  string
	Category  string
	IsBinary  bool
	IsText    bool
	IsImage   bool
	IsVideo   bool
	IsAudio   bool
	IsArchive bool
}

type IntegrityReport struct {
	IsValid        bool
	Errors         []string
	Warnings       []string
	Checksum       string
	Algorithm      string
	FileSize       int64
	LastModified   time.Time
	CorruptedBytes []int64
}

type ArchiveReport struct {
	Type            string
	TotalFiles      int
	CompressedSize  int64
	UncompressedSize int64
	CompressionRatio float64
	FileList        []FileEntry
	HasPassword     bool
	IsCorrupted     bool
}

type FileEntry struct {
	Path         string
	Size         int64
	CompressedSize int64
	ModTime      time.Time
	IsDirectory  bool
	CRC32        uint32
}

type DuplicateSet struct {
	Files    []string
	Size     int64
	Checksum string
	Count    int
}

// Basic file operations

func ExtractMetadata(filePath string) (FileMetadata, error) {
	metadata := FileMetadata{Path: filePath}
	
	info, err := os.Stat(filePath)
	if err != nil {
		return metadata, err
	}
	
	metadata.Name = info.Name()
	metadata.Size = info.Size()
	metadata.ModTime = info.ModTime()
	metadata.IsDirectory = info.IsDir()
	metadata.Permissions = info.Mode()
	metadata.Extension = strings.ToLower(filepath.Ext(filePath))
	
	if !metadata.IsDirectory {
		metadata.MimeType = detectMimeType(filePath)
	}
	
	return metadata, nil
}

func detectMimeType(filePath string) string {
	ext := strings.ToLower(filepath.Ext(filePath))
	
	mimeTypes := map[string]string{
		".txt":  "text/plain",
		".html": "text/html",
		".css":  "text/css",
		".js":   "application/javascript",
		".json": "application/json",
		".xml":  "application/xml",
		".pdf":  "application/pdf",
		".jpg":  "image/jpeg",
		".jpeg": "image/jpeg",
		".png":  "image/png",
		".gif":  "image/gif",
		".svg":  "image/svg+xml",
		".mp3":  "audio/mpeg",
		".wav":  "audio/wav",
		".mp4":  "video/mp4",
		".avi":  "video/x-msvideo",
		".zip":  "application/zip",
		".tar":  "application/x-tar",
		".gz":   "application/gzip",
		".7z":   "application/x-7z-compressed",
		".exe":  "application/x-executable",
		".dll":  "application/x-msdownload",
		".so":   "application/x-sharedlib",
		".go":   "text/x-go",
		".py":   "text/x-python",
		".java": "text/x-java",
		".c":    "text/x-c",
		".cpp":  "text/x-c++",
		".h":    "text/x-c-header",
	}
	
	if mimeType, exists := mimeTypes[ext]; exists {
		return mimeType
	}
	
	return "application/octet-stream"
}

func AnalyzeFileStructure(dirPath string) (StructureReport, error) {
	report := StructureReport{
		FileTypes:    make(map[string]int),
		LargestFiles: make([]FileMetadata, 0),
	}
	
	var totalSize int64
	var maxDepth int
	var deepestPath string
	var fileSizes []int64
	
	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil // Skip errors and continue
		}
		
		// Calculate depth
		relPath, _ := filepath.Rel(dirPath, path)
		depth := strings.Count(relPath, string(filepath.Separator))
		if depth > maxDepth {
			maxDepth = depth
			deepestPath = path
		}
		
		if info.IsDir() {
			report.TotalDirectories++
		} else {
			report.TotalFiles++
			size := info.Size()
			totalSize += size
			fileSizes = append(fileSizes, size)
			
			// Track file types
			ext := strings.ToLower(filepath.Ext(path))
			if ext == "" {
				ext = "no-extension"
			}
			report.FileTypes[ext]++
			
			// Track largest files (keep top 10)
			metadata, _ := ExtractMetadata(path)
			if len(report.LargestFiles) < 10 {
				report.LargestFiles = append(report.LargestFiles, metadata)
			} else {
				// Find smallest in top 10 and replace if current is larger
				smallestIdx := 0
				for i, file := range report.LargestFiles {
					if file.Size < report.LargestFiles[smallestIdx].Size {
						smallestIdx = i
					}
				}
				if size > report.LargestFiles[smallestIdx].Size {
					report.LargestFiles[smallestIdx] = metadata
				}
			}
		}
		
		return nil
	})
	
	if err != nil {
		return report, err
	}
	
	report.TotalSize = totalSize
	report.MaxDepth = maxDepth
	report.DeepestPath = deepestPath
	
	if len(fileSizes) > 0 {
		report.AverageFileSize = float64(totalSize) / float64(len(fileSizes))
	}
	
	return report, nil
}

func DetectFileType(content []byte) FileType {
	fileType := FileType{}
	
	if len(content) == 0 {
		return fileType
	}
	
	// Check magic bytes for common file types
	if len(content) >= 4 {
		header := content[:4]
		
		// PDF
		if string(content[:4]) == "%PDF" {
			fileType.MimeType = "application/pdf"
			fileType.Extension = ".pdf"
			fileType.Category = "document"
			fileType.IsBinary = true
		}
		
		// PNG
		if len(content) >= 8 && string(content[:8]) == "\x89PNG\r\n\x1a\n" {
			fileType.MimeType = "image/png"
			fileType.Extension = ".png"
			fileType.Category = "image"
			fileType.IsImage = true
			fileType.IsBinary = true
		}
		
		// JPEG
		if len(content) >= 2 && content[0] == 0xFF && content[1] == 0xD8 {
			fileType.MimeType = "image/jpeg"
			fileType.Extension = ".jpg"
			fileType.Category = "image"
			fileType.IsImage = true
			fileType.IsBinary = true
		}
		
		// ZIP
		if string(header[:2]) == "PK" {
			fileType.MimeType = "application/zip"
			fileType.Extension = ".zip"
			fileType.Category = "archive"
			fileType.IsArchive = true
			fileType.IsBinary = true
		}
		
		// ELF (Linux executable)
		if string(header) == "\x7fELF" {
			fileType.MimeType = "application/x-executable"
			fileType.Extension = ""
			fileType.Category = "executable"
			fileType.IsBinary = true
		}
		
		// Windows PE executable
		if string(content[:2]) == "MZ" {
			fileType.MimeType = "application/x-msdownload"
			fileType.Extension = ".exe"
			fileType.Category = "executable"
			fileType.IsBinary = true
		}
	}
	
	// Check if it's text by examining content
	if fileType.MimeType == "" {
		isText := true
		for i := 0; i < len(content) && i < 512; i++ {
			if content[i] == 0 {
				isText = false
				break
			}
		}
		
		if isText {
			fileType.IsText = true
			fileType.MimeType = "text/plain"
			fileType.Category = "text"
			
			// Try to detect specific text formats
			contentStr := string(content[:minFileSize(1024, len(content))])
			if strings.Contains(contentStr, "<html") || strings.Contains(contentStr, "<!DOCTYPE") {
				fileType.MimeType = "text/html"
				fileType.Extension = ".html"
			} else if strings.HasPrefix(contentStr, "{") || strings.HasPrefix(contentStr, "[") {
				fileType.MimeType = "application/json"
				fileType.Extension = ".json"
			} else if strings.Contains(contentStr, "<?xml") {
				fileType.MimeType = "application/xml"
				fileType.Extension = ".xml"
			}
		} else {
			fileType.IsBinary = true
			fileType.MimeType = "application/octet-stream"
			fileType.Category = "binary"
		}
	}
	
	return fileType
}

func minFileSize(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func ValidateFileIntegrity(filePath string) IntegrityReport {
	report := IntegrityReport{IsValid: true}
	
	info, err := os.Stat(filePath)
	if err != nil {
		report.IsValid = false
		report.Errors = append(report.Errors, "File does not exist: "+err.Error())
		return report
	}
	
	report.FileSize = info.Size()
	report.LastModified = info.ModTime()
	
	// Check if file is readable
	file, err := os.Open(filePath)
	if err != nil {
		report.IsValid = false
		report.Errors = append(report.Errors, "Cannot read file: "+err.Error())
		return report
	}
	defer file.Close()
	
	// Calculate checksum
	checksum, err := CalculateChecksum(filePath, "sha256")
	if err != nil {
		report.Warnings = append(report.Warnings, "Could not calculate checksum: "+err.Error())
	} else {
		report.Checksum = checksum
		report.Algorithm = "sha256"
	}
	
	// Basic integrity checks
	if info.Size() == 0 {
		report.Warnings = append(report.Warnings, "File is empty")
	}
	
	// Check for typical corruption patterns (null bytes in text files)
	ext := strings.ToLower(filepath.Ext(filePath))
	textExtensions := map[string]bool{
		".txt": true, ".md": true, ".json": true, ".xml": true,
		".html": true, ".css": true, ".js": true, ".py": true,
		".go": true, ".java": true, ".c": true, ".cpp": true,
	}
	
	if textExtensions[ext] {
		buffer := make([]byte, 1024)
		n, err := file.Read(buffer)
		if err != nil && err != io.EOF {
			report.Warnings = append(report.Warnings, "Error reading file content")
		} else {
			for i := 0; i < n; i++ {
				if buffer[i] == 0 {
					report.CorruptedBytes = append(report.CorruptedBytes, int64(i))
					report.Warnings = append(report.Warnings, fmt.Sprintf("Null byte found at position %d", i))
				}
			}
		}
	}
	
	return report
}

func CalculateChecksum(filePath string, algorithm string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()
	
	switch strings.ToLower(algorithm) {
	case "md5":
		h := md5.New()
		if _, err := io.Copy(h, file); err != nil {
			return "", err
		}
		return hex.EncodeToString(h.Sum(nil)), nil
		
	case "sha1":
		h := sha1.New()
		if _, err := io.Copy(h, file); err != nil {
			return "", err
		}
		return hex.EncodeToString(h.Sum(nil)), nil
		
	case "sha256":
		h := sha256.New()
		if _, err := io.Copy(h, file); err != nil {
			return "", err
		}
		return hex.EncodeToString(h.Sum(nil)), nil
		
	default:
		return "", fmt.Errorf("unsupported algorithm: %s", algorithm)
	}
}

// Archive/compression functions

func AnalyzeArchiveStructure(archivePath string) (ArchiveReport, error) {
	report := ArchiveReport{}
	
	fileInfo, err := os.Stat(archivePath)
	if err != nil {
		return report, err
	}
	
	report.CompressedSize = fileInfo.Size()
	
	// Determine archive type from extension
	ext := strings.ToLower(filepath.Ext(archivePath))
	switch ext {
	case ".zip":
		report.Type = "ZIP"
		return analyzeZipStructure(archivePath, report)
	case ".tar":
		report.Type = "TAR"
		return analyzeTarStructure(archivePath, report)
	case ".gz":
		if strings.HasSuffix(strings.ToLower(archivePath), ".tar.gz") {
			report.Type = "TAR.GZ"
		} else {
			report.Type = "GZIP"
		}
		return analyzeGzipStructure(archivePath, report)
	case ".7z":
		report.Type = "7Z"
		return analyze7zStructure(archivePath, report)
	default:
		return report, fmt.Errorf("unsupported archive type: %s", ext)
	}
}

func analyzeZipStructure(archivePath string, report ArchiveReport) (ArchiveReport, error) {
	// Basic ZIP analysis without external dependencies
	file, err := os.Open(archivePath)
	if err != nil {
		return report, err
	}
	defer file.Close()
	
	// Read ZIP header to check basic structure
	header := make([]byte, 4)
	_, err = file.Read(header)
	if err != nil {
		return report, err
	}
	
	// Check for ZIP signature
	if string(header[:2]) == "PK" {
		report.IsCorrupted = false
		// For full ZIP analysis, would need zip package
		// This is a basic structure check
		report.TotalFiles = 1 // Placeholder
	} else {
		report.IsCorrupted = true
	}
	
	return report, nil
}

func analyzeTarStructure(archivePath string, report ArchiveReport) (ArchiveReport, error) {
	// Basic TAR analysis
	file, err := os.Open(archivePath)
	if err != nil {
		return report, err
	}
	defer file.Close()
	
	// TAR files don't have a magic signature, so we'll do basic validation
	stat, err := file.Stat()
	if err != nil {
		return report, err
	}
	
	// TAR blocks are 512 bytes
	if stat.Size()%512 == 0 {
		report.IsCorrupted = false
	} else {
		report.IsCorrupted = true
	}
	
	return report, nil
}

func analyzeGzipStructure(archivePath string, report ArchiveReport) (ArchiveReport, error) {
	file, err := os.Open(archivePath)
	if err != nil {
		return report, err
	}
	defer file.Close()
	
	// Check GZIP header
	header := make([]byte, 3)
	_, err = file.Read(header)
	if err != nil {
		return report, err
	}
	
	// GZIP magic number: 0x1f, 0x8b
	if header[0] == 0x1f && header[1] == 0x8b {
		report.IsCorrupted = false
	} else {
		report.IsCorrupted = true
	}
	
	return report, nil
}

func analyze7zStructure(archivePath string, report ArchiveReport) (ArchiveReport, error) {
	file, err := os.Open(archivePath)
	if err != nil {
		return report, err
	}
	defer file.Close()
	
	// Check 7z header
	header := make([]byte, 6)
	_, err = file.Read(header)
	if err != nil {
		return report, err
	}
	
	// 7z signature: "7z¼¯'"
	if string(header) == "7z\xBC\xAF\x27\x1C" {
		report.IsCorrupted = false
	} else {
		report.IsCorrupted = true
	}
	
	return report, nil
}

func ExtractFileList(archivePath string) ([]FileEntry, error) {
	// This would require specific archive handling libraries
	// For now, return a basic implementation
	entries := []FileEntry{}
	
	ext := strings.ToLower(filepath.Ext(archivePath))
	if ext == ".zip" || ext == ".tar" || ext == ".gz" || ext == ".7z" {
		// Placeholder - would need proper archive libraries
		entry := FileEntry{
			Path:        "example.txt",
			Size:        1024,
			ModTime:     time.Now(),
			IsDirectory: false,
		}
		entries = append(entries, entry)
	}
	
	return entries, nil
}

func DetectCompressionRatio(filePath string) (float64, error) {
	_, err := os.Stat(filePath)
	if err != nil {
		return 0, err
	}
	
	// For actual compression ratio, we'd need to decompress or analyze archive metadata
	// This is a placeholder implementation
	ext := strings.ToLower(filepath.Ext(filePath))
	
	// Typical compression ratios for different formats
	switch ext {
	case ".zip":
		return 0.6, nil // 60% of original size
	case ".gz":
		return 0.3, nil // 30% of original size
	case ".7z":
		return 0.25, nil // 25% of original size
	case ".tar":
		return 1.0, nil // No compression
	default:
		return 1.0, nil // No compression
	}
}

func FindDuplicateFiles(directory string) ([]DuplicateSet, error) {
	fileHashes := make(map[string][]string)
	duplicates := []DuplicateSet{}
	
	err := filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil // Skip errors
		}
		
		if !info.IsDir() && info.Size() > 0 {
			checksum, err := CalculateChecksum(path, "md5")
			if err != nil {
				return nil // Skip files we can't read
			}
			
			fileHashes[checksum] = append(fileHashes[checksum], path)
		}
		
		return nil
	})
	
	if err != nil {
		return duplicates, err
	}
	
	// Find duplicates
	for checksum, files := range fileHashes {
		if len(files) > 1 {
			// Get size from first file
			info, err := os.Stat(files[0])
			var size int64
			if err == nil {
				size = info.Size()
			}
			
			duplicate := DuplicateSet{
				Files:    files,
				Size:     size,
				Checksum: checksum,
				Count:    len(files),
			}
			duplicates = append(duplicates, duplicate)
		}
	}
	
	return duplicates, nil
}