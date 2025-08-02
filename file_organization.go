package textlib

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"
)

// File organization data structures

type Category string

const (
	CategoryDocuments    Category = "documents"
	CategoryImages       Category = "images"
	CategoryVideos       Category = "videos"
	CategoryAudio        Category = "audio"
	CategoryArchives     Category = "archives"
	CategoryExecutables  Category = "executables"
	CategoryCode         Category = "code"
	CategoryData         Category = "data"
	CategoryTemp         Category = "temp"
	CategoryConfig       Category = "config"
	CategoryOther        Category = "other"
)

type File struct {
	Path     string
	Name     string
	Size     int64
	ModTime  time.Time
	Category Category
	Type     string
}

type VersionSet struct {
	BaseName string
	Files    []File
	Pattern  string
	Latest   File
}

type DateStrategy string

const (
	DateByYear       DateStrategy = "year"
	DateByMonth      DateStrategy = "month"
	DateByDay        DateStrategy = "day"
	DateByYearMonth  DateStrategy = "year_month"
)

type DirectoryStructure struct {
	Root        string
	Directories map[string][]File
	TotalFiles  int
	TotalSize   int64
}

type CleanupReport struct {
	DeletedFiles   []string
	DeletedDirs    []string
	FreedSpace     int64
	Errors         []string
	SkippedFiles   []string
	ProcessedCount int
}

// File management functions

func CategorizeFiles(directory string) (map[Category][]File, error) {
	categories := make(map[Category][]File)
	
	err := filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil // Skip errors and continue
		}
		
		if !info.IsDir() {
			file := File{
				Path:    path,
				Name:    info.Name(),
				Size:    info.Size(),
				ModTime: info.ModTime(),
			}
			
			file.Category = categorizeByExtension(filepath.Ext(file.Name))
			file.Type = filepath.Ext(file.Name)
			
			categories[file.Category] = append(categories[file.Category], file)
		}
		
		return nil
	})
	
	return categories, err
}

func categorizeByExtension(ext string) Category {
	ext = strings.ToLower(ext)
	
	documentExts := map[string]bool{
		".pdf": true, ".doc": true, ".docx": true, ".txt": true,
		".rtf": true, ".odt": true, ".pages": true, ".md": true,
		".tex": true, ".epub": true, ".mobi": true,
	}
	
	imageExts := map[string]bool{
		".jpg": true, ".jpeg": true, ".png": true, ".gif": true,
		".bmp": true, ".tiff": true, ".svg": true, ".webp": true,
		".ico": true, ".psd": true, ".ai": true, ".eps": true,
	}
	
	videoExts := map[string]bool{
		".mp4": true, ".avi": true, ".mkv": true, ".mov": true,
		".wmv": true, ".flv": true, ".webm": true, ".m4v": true,
		".3gp": true, ".ogv": true, ".ts": true,
	}
	
	audioExts := map[string]bool{
		".mp3": true, ".wav": true, ".flac": true, ".aac": true,
		".ogg": true, ".wma": true, ".m4a": true, ".opus": true,
		".aiff": true, ".ape": true,
	}
	
	archiveExts := map[string]bool{
		".zip": true, ".rar": true, ".7z": true, ".tar": true,
		".gz": true, ".bz2": true, ".xz": true, ".iso": true,
		".dmg": true, ".deb": true, ".rpm": true,
	}
	
	executableExts := map[string]bool{
		".exe": true, ".msi": true, ".app": true, ".deb": true,
		".rpm": true, ".dmg": true, ".pkg": true, ".run": true,
		".bin": true, ".com": true, ".scr": true,
	}
	
	codeExts := map[string]bool{
		".go": true, ".py": true, ".js": true, ".ts": true,
		".java": true, ".c": true, ".cpp": true, ".h": true,
		".cs": true, ".php": true, ".rb": true, ".rs": true,
		".swift": true, ".kt": true, ".scala": true, ".clj": true,
		".hs": true, ".ml": true, ".r": true, ".m": true,
		".pl": true, ".sh": true, ".bat": true, ".ps1": true,
	}
	
	dataExts := map[string]bool{
		".json": true, ".xml": true, ".csv": true, ".yaml": true,
		".yml": true, ".toml": true, ".sql": true, ".db": true,
		".sqlite": true, ".log": true, ".dat": true,
	}
	
	configExts := map[string]bool{
		".conf": true, ".config": true, ".cfg": true, ".ini": true,
		".properties": true, ".env": true, ".rc": true,
	}
	
	tempExts := map[string]bool{
		".tmp": true, ".temp": true, ".cache": true, ".bak": true,
		".backup": true, ".old": true, ".~": true,
	}
	
	switch {
	case documentExts[ext]:
		return CategoryDocuments
	case imageExts[ext]:
		return CategoryImages
	case videoExts[ext]:
		return CategoryVideos
	case audioExts[ext]:
		return CategoryAudio
	case archiveExts[ext]:
		return CategoryArchives
	case executableExts[ext]:
		return CategoryExecutables
	case codeExts[ext]:
		return CategoryCode
	case dataExts[ext]:
		return CategoryData
	case configExts[ext]:
		return CategoryConfig
	case tempExts[ext]:
		return CategoryTemp
	default:
		return CategoryOther
	}
}

func FindUnusedFiles(projectPath string) ([]File, error) {
	var unusedFiles []File
	var allFiles []File
	
	// First, collect all files
	err := filepath.Walk(projectPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		
		if !info.IsDir() {
			file := File{
				Path:    path,
				Name:    info.Name(),
				Size:    info.Size(),
				ModTime: info.ModTime(),
			}
			allFiles = append(allFiles, file)
		}
		
		return nil
	})
	
	if err != nil {
		return unusedFiles, err
	}
	
	// Build reference map
	references := make(map[string]bool)
	
	for _, file := range allFiles {
		if isCodeFile(file.Path) {
			refs := extractFileReferences(file.Path)
			for _, ref := range refs {
				references[ref] = true
			}
		}
	}
	
	// Find unused files
	for _, file := range allFiles {
		if isAssetFile(file.Path) {
			relPath, _ := filepath.Rel(projectPath, file.Path)
			fileName := filepath.Base(file.Path)
			fileNameNoExt := strings.TrimSuffix(fileName, filepath.Ext(fileName))
			
			// Check if file is referenced
			if !references[relPath] && !references[fileName] && 
			   !references[fileNameNoExt] && !references[file.Path] {
				unusedFiles = append(unusedFiles, file)
			}
		}
	}
	
	return unusedFiles, nil
}

func isCodeFile(path string) bool {
	codeExts := []string{".go", ".py", ".js", ".ts", ".java", ".c", ".cpp", ".cs", ".php", ".rb", ".rs"}
	ext := strings.ToLower(filepath.Ext(path))
	
	for _, codeExt := range codeExts {
		if ext == codeExt {
			return true
		}
	}
	
	return false
}

func isAssetFile(path string) bool {
	assetExts := []string{".png", ".jpg", ".jpeg", ".gif", ".svg", ".css", ".js", ".json", ".xml", ".html"}
	ext := strings.ToLower(filepath.Ext(path))
	
	for _, assetExt := range assetExts {
		if ext == assetExt {
			return true
		}
	}
	
	return false
}

func extractFileReferences(filePath string) []string {
	var references []string
	
	file, err := os.Open(filePath)
	if err != nil {
		return references
	}
	defer file.Close()
	
	content, err := io.ReadAll(file)
	if err != nil {
		return references
	}
	
	contentStr := string(content)
	
	// Common import/include patterns
	patterns := []*regexp.Regexp{
		regexp.MustCompile(`import\s+["']([^"']+)["']`),           // ES6 imports
		regexp.MustCompile(`require\s*\(\s*["']([^"']+)["']\s*\)`), // CommonJS require
		regexp.MustCompile(`#include\s*[<"]([^>"]+)[>"]`),        // C/C++ includes
		regexp.MustCompile(`from\s+["']([^"']+)["']`),            // Python imports
		regexp.MustCompile(`import\s+["']([^"']+)["']`),          // Python imports
		regexp.MustCompile(`src\s*=\s*["']([^"']+)["']`),         // HTML src attributes
		regexp.MustCompile(`href\s*=\s*["']([^"']+)["']`),        // HTML href attributes
		regexp.MustCompile(`url\s*\(\s*["']?([^"')]+)["']?\s*\)`), // CSS url()
		regexp.MustCompile(`["']([^"']*\.(png|jpg|jpeg|gif|svg|css|js|json|xml|html))["']`), // File references
	}
	
	for _, pattern := range patterns {
		matches := pattern.FindAllStringSubmatch(contentStr, -1)
		for _, match := range matches {
			if len(match) > 1 {
				ref := strings.TrimSpace(match[1])
				if ref != "" && !strings.HasPrefix(ref, "http") && !strings.HasPrefix(ref, "//") {
					references = append(references, ref)
				}
			}
		}
	}
	
	return references
}

func DetectVersionedFiles(directory string) ([]VersionSet, error) {
	var versionSets []VersionSet
	fileMap := make(map[string][]File)
	
	err := filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		
		if !info.IsDir() {
			file := File{
				Path:    path,
				Name:    info.Name(),
				Size:    info.Size(),
				ModTime: info.ModTime(),
			}
			
			fileMap[extractBaseName(file.Name)] = append(fileMap[extractBaseName(file.Name)], file)
		}
		
		return nil
	})
	
	if err != nil {
		return versionSets, err
	}
	
	// Find files with multiple versions
	for baseName, files := range fileMap {
		if len(files) > 1 {
			// Sort by modification time
			sort.Slice(files, func(i, j int) bool {
				return files[i].ModTime.After(files[j].ModTime)
			})
			
			versionSet := VersionSet{
				BaseName: baseName,
				Files:    files,
				Pattern:  detectVersionPattern(files),
				Latest:   files[0], // Most recent
			}
			
			versionSets = append(versionSets, versionSet)
		}
	}
	
	return versionSets, nil
}

func extractBaseName(filename string) string {
	// Remove common version patterns
	patterns := []*regexp.Regexp{
		regexp.MustCompile(`_v\d+(\.\d+)*`),           // file_v1.2.3
		regexp.MustCompile(`_\d{4}-\d{2}-\d{2}`),      // file_2023-01-01
		regexp.MustCompile(`_\d+`),                    // file_1, file_2
		regexp.MustCompile(`\(\d+\)`),                 // file(1), file(2)
		regexp.MustCompile(`\s-\s\w+`),                // file - Copy
		regexp.MustCompile(`\s\(\d+\)`),               // file (1)
		regexp.MustCompile(`\.backup$`),               // file.backup
		regexp.MustCompile(`\.bak$`),                  // file.bak
		regexp.MustCompile(`\.old$`),                  // file.old
		regexp.MustCompile(`~$`),                      // file~
	}
	
	ext := filepath.Ext(filename)
	nameWithoutExt := strings.TrimSuffix(filename, ext)
	
	for _, pattern := range patterns {
		nameWithoutExt = pattern.ReplaceAllString(nameWithoutExt, "")
	}
	
	return nameWithoutExt + ext
}

func detectVersionPattern(files []File) string {
	if len(files) < 2 {
		return "unknown"
	}
	
	// Analyze filenames to detect pattern
	first := files[0].Name
	second := files[1].Name
	
	patterns := []struct {
		regex   *regexp.Regexp
		pattern string
	}{
		{regexp.MustCompile(`_v\d+(\.\d+)*`), "semantic_version"},
		{regexp.MustCompile(`_\d{4}-\d{2}-\d{2}`), "date_version"},
		{regexp.MustCompile(`_\d+`), "numeric_suffix"},
		{regexp.MustCompile(`\(\d+\)`), "parentheses_number"},
		{regexp.MustCompile(`\s-\s\w+`), "copy_suffix"},
	}
	
	for _, p := range patterns {
		if p.regex.MatchString(first) && p.regex.MatchString(second) {
			return p.pattern
		}
	}
	
	return "custom"
}

func OrganizeByDate(files []File, strategy DateStrategy) DirectoryStructure {
	structure := DirectoryStructure{
		Directories: make(map[string][]File),
	}
	
	for _, file := range files {
		var dirName string
		
		switch strategy {
		case DateByYear:
			dirName = file.ModTime.Format("2006")
		case DateByMonth:
			dirName = file.ModTime.Format("2006-01")
		case DateByDay:
			dirName = file.ModTime.Format("2006-01-02")
		case DateByYearMonth:
			dirName = file.ModTime.Format("2006/01")
		default:
			dirName = "unknown"
		}
		
		structure.Directories[dirName] = append(structure.Directories[dirName], file)
		structure.TotalFiles++
		structure.TotalSize += file.Size
	}
	
	return structure
}

func CleanupTempFiles(directory string) (CleanupReport, error) {
	report := CleanupReport{
		DeletedFiles: make([]string, 0),
		DeletedDirs:  make([]string, 0),
		Errors:       make([]string, 0),
		SkippedFiles: make([]string, 0),
	}
	
	// Patterns for temporary files
	tempPatterns := []*regexp.Regexp{
		regexp.MustCompile(`\.tmp$`),
		regexp.MustCompile(`\.temp$`),
		regexp.MustCompile(`\.cache$`),
		regexp.MustCompile(`\.bak$`),
		regexp.MustCompile(`\.backup$`),
		regexp.MustCompile(`\.old$`),
		regexp.MustCompile(`~$`),
		regexp.MustCompile(`^\.(DS_Store|Thumbs\.db)$`),
		regexp.MustCompile(`^\.#`), // Emacs lock files
		regexp.MustCompile(`#.*#$`), // Emacs backup files
		regexp.MustCompile(`\.swp$`), // Vim swap files
		regexp.MustCompile(`\.swo$`), // Vim swap files
	}
	
	// Temporary directory patterns
	tempDirPatterns := []*regexp.Regexp{
		regexp.MustCompile(`^tmp`),
		regexp.MustCompile(`^temp`),
		regexp.MustCompile(`cache$`),
		regexp.MustCompile(`\.cache$`),
	}
	
	var filesToDelete []string
	var dirsToDelete []string
	
	err := filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			report.Errors = append(report.Errors, err.Error())
			return nil
		}
		
		name := info.Name()
		report.ProcessedCount++
		
		if info.IsDir() {
			// Check if directory should be deleted
			for _, pattern := range tempDirPatterns {
				if pattern.MatchString(name) {
					// Check if directory is empty or contains only temp files
					if isEmpty, _ := isDirEmpty(path); isEmpty {
						dirsToDelete = append(dirsToDelete, path)
					}
					break
				}
			}
		} else {
			// Check if file should be deleted
			shouldDelete := false
			
			for _, pattern := range tempPatterns {
				if pattern.MatchString(name) {
					shouldDelete = true
					break
				}
			}
			
			// Check file age (delete files older than 7 days in temp directories)
			if strings.Contains(path, "tmp") || strings.Contains(path, "temp") {
				if time.Since(info.ModTime()) > 7*24*time.Hour {
					shouldDelete = true
				}
			}
			
			// Check for zero-byte files
			if info.Size() == 0 && shouldDelete {
				shouldDelete = true
			}
			
			if shouldDelete {
				// Safety check - don't delete important files
				if isImportantFile(path) {
					report.SkippedFiles = append(report.SkippedFiles, path)
				} else {
					filesToDelete = append(filesToDelete, path)
					report.FreedSpace += info.Size()
				}
			}
		}
		
		return nil
	})
	
	if err != nil {
		return report, err
	}
	
	// Delete files
	for _, file := range filesToDelete {
		err := os.Remove(file)
		if err != nil {
			report.Errors = append(report.Errors, fmt.Sprintf("Failed to delete %s: %v", file, err))
		} else {
			report.DeletedFiles = append(report.DeletedFiles, file)
		}
	}
	
	// Delete directories (in reverse order to handle nested directories)
	for i := len(dirsToDelete) - 1; i >= 0; i-- {
		dir := dirsToDelete[i]
		err := os.Remove(dir)
		if err != nil {
			report.Errors = append(report.Errors, fmt.Sprintf("Failed to delete directory %s: %v", dir, err))
		} else {
			report.DeletedDirs = append(report.DeletedDirs, dir)
		}
	}
	
	return report, nil
}

func isDirEmpty(path string) (bool, error) {
	f, err := os.Open(path)
	if err != nil {
		return false, err
	}
	defer f.Close()
	
	_, err = f.Readdirnames(1)
	if err == io.EOF {
		return true, nil
	}
	return false, err
}

func isImportantFile(path string) bool {
	importantPatterns := []*regexp.Regexp{
		regexp.MustCompile(`\.git`),
		regexp.MustCompile(`\.ssh`),
		regexp.MustCompile(`\.config`),
		regexp.MustCompile(`license`),
		regexp.MustCompile(`readme`),
		regexp.MustCompile(`makefile`),
		regexp.MustCompile(`dockerfile`),
		regexp.MustCompile(`\.env`),
	}
	
	lowerPath := strings.ToLower(path)
	
	for _, pattern := range importantPatterns {
		if pattern.MatchString(lowerPath) {
			return true
		}
	}
	
	// Check if file is in an important directory
	importantDirs := []string{"src", "lib", "bin", "etc", "usr", "opt"}
	for _, dir := range importantDirs {
		if strings.Contains(lowerPath, "/"+dir+"/") {
			return true
		}
	}
	
	return false
}

// Additional utility functions

func GetFilesBySize(directory string, minSize int64) ([]File, error) {
	var largeFiles []File
	
	err := filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		
		if !info.IsDir() && info.Size() >= minSize {
			file := File{
				Path:    path,
				Name:    info.Name(),
				Size:    info.Size(),
				ModTime: info.ModTime(),
			}
			largeFiles = append(largeFiles, file)
		}
		
		return nil
	})
	
	// Sort by size (largest first)
	sort.Slice(largeFiles, func(i, j int) bool {
		return largeFiles[i].Size > largeFiles[j].Size
	})
	
	return largeFiles, err
}

func GetFilesByAge(directory string, olderThan time.Duration) ([]File, error) {
	var oldFiles []File
	cutoff := time.Now().Add(-olderThan)
	
	err := filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		
		if !info.IsDir() && info.ModTime().Before(cutoff) {
			file := File{
				Path:    path,
				Name:    info.Name(),
				Size:    info.Size(),
				ModTime: info.ModTime(),
			}
			oldFiles = append(oldFiles, file)
		}
		
		return nil
	})
	
	// Sort by age (oldest first)
	sort.Slice(oldFiles, func(i, j int) bool {
		return oldFiles[i].ModTime.Before(oldFiles[j].ModTime)
	})
	
	return oldFiles, err
}

func AnalyzeDiskUsage(directory string) (map[string]int64, error) {
	usage := make(map[string]int64)
	
	err := filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		
		if !info.IsDir() {
			dir := filepath.Dir(path)
			usage[dir] += info.Size()
		}
		
		return nil
	})
	
	return usage, err
}