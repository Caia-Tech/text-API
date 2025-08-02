package textlib

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
)

// Security analysis data structures

type SecurityThreat struct {
	Type        string
	Severity    string
	Description string
	Position    int64
	Pattern     string
	Confidence  float64
	Mitigation  string
}

type ExecutableInfo struct {
	Format       string
	Architecture string
	Entrypoint   uint64
	Sections     []Section
	Imports      []ExecutableImport
	Exports      []Export
	HasDebugInfo bool
	IsPacked     bool
	Compiler     string
	Timestamp    uint32
}

type Section struct {
	Name         string
	VirtualAddr  uint64
	RawSize      uint32
	Permissions  string
	IsExecutable bool
	IsWritable   bool
	IsReadable   bool
}

type ExecutableImport struct {
	Library  string
	Function string
	Address  uint64
}

type Export struct {
	Name    string
	Address uint64
	Ordinal uint16
}

type EmbeddedFile struct {
	Type     string
	Offset   int64
	Size     int64
	Filename string
	Hash     string
}

type PermissionReport struct {
	Owner       PermissionSet
	Group       PermissionSet
	Others      PermissionSet
	Octal       string
	IsSecure    bool
	Warnings    []string
	Suggestions []string
}

type PermissionSet struct {
	Read    bool
	Write   bool
	Execute bool
}

type VirusSignature struct {
	Name        string
	Pattern     []byte
	Description string
	Severity    string
}

type Threat struct {
	Name        string
	Type        string
	Severity    string
	Description string
	Position    int64
	Confidence  float64
}

// Security scanning functions

func DetectMaliciousPatterns(filePath string) ([]SecurityThreat, error) {
	var threats []SecurityThreat
	
	file, err := os.Open(filePath)
	if err != nil {
		return threats, err
	}
	defer file.Close()
	
	content, err := io.ReadAll(file)
	if err != nil {
		return threats, err
	}
	
	// Check for various malicious patterns
	threats = append(threats, detectShellcodePatterns(content)...)
	threats = append(threats, detectSuspiciousStrings(content)...)
	threats = append(threats, detectObfuscationPatterns(content)...)
	threats = append(threats, detectNetworkPatterns(content)...)
	threats = append(threats, detectCryptoPatterns(content)...)
	threats = append(threats, detectAntiAnalysisPatterns(content)...)
	
	return threats, nil
}

func detectShellcodePatterns(content []byte) []SecurityThreat {
	var threats []SecurityThreat
	
	// Common shellcode patterns
	patterns := []struct {
		pattern     []byte
		description string
		severity    string
	}{
		{[]byte{0x90, 0x90, 0x90, 0x90}, "NOP sled detected", "high"},
		{[]byte{0xEB, 0xFE}, "Infinite loop (anti-debugging)", "medium"},
		{[]byte{0xCC}, "Debug breakpoint", "low"},
		{[]byte{0x31, 0xC0}, "XOR EAX, EAX (common in shellcode)", "medium"},
		{[]byte{0x64, 0x8B, 0x15}, "PEB access pattern", "high"},
	}
	
	for _, p := range patterns {
		offset := 0
		for {
			pos := bytes.Index(content[offset:], p.pattern)
			if pos == -1 {
				break
			}
			
			threat := SecurityThreat{
				Type:        "shellcode",
				Severity:    p.severity,
				Description: p.description,
				Position:    int64(offset + pos),
				Pattern:     fmt.Sprintf("%x", p.pattern),
				Confidence:  0.7,
				Mitigation:  "Review code for potential shellcode injection",
			}
			threats = append(threats, threat)
			offset += pos + len(p.pattern)
		}
	}
	
	return threats
}

func detectSuspiciousStrings(content []byte) []SecurityThreat {
	var threats []SecurityThreat
	
	suspiciousStrings := []struct {
		pattern     string
		description string
		severity    string
		confidence  float64
	}{
		{"cmd.exe", "Command shell execution", "high", 0.8},
		{"powershell", "PowerShell execution", "high", 0.8},
		{"CreateRemoteThread", "Process injection API", "high", 0.9},
		{"VirtualAlloc", "Memory allocation API", "medium", 0.7},
		{"WriteProcessMemory", "Memory writing API", "high", 0.9},
		{"GetProcAddress", "Dynamic API resolution", "medium", 0.6},
		{"LoadLibrary", "Dynamic library loading", "medium", 0.6},
		{"keylogger", "Keylogger functionality", "high", 0.9},
		{"password", "Password stealing", "medium", 0.5},
		{"bitcoin", "Cryptocurrency mining", "medium", 0.6},
		{"ransomware", "Ransomware functionality", "critical", 1.0},
		{"rootkit", "Rootkit functionality", "critical", 1.0},
		{"backdoor", "Backdoor functionality", "critical", 1.0},
	}
	
	contentStr := strings.ToLower(string(content))
	
	for _, s := range suspiciousStrings {
		pattern := strings.ToLower(s.pattern)
		offset := 0
		for {
			pos := strings.Index(contentStr[offset:], pattern)
			if pos == -1 {
				break
			}
			
			threat := SecurityThreat{
				Type:        "suspicious_string",
				Severity:    s.severity,
				Description: s.description,
				Position:    int64(offset + pos),
				Pattern:     s.pattern,
				Confidence:  s.confidence,
				Mitigation:  "Investigate context and legitimacy of this functionality",
			}
			threats = append(threats, threat)
			offset += pos + len(pattern)
		}
	}
	
	return threats
}

func detectObfuscationPatterns(content []byte) []SecurityThreat {
	var threats []SecurityThreat
	
	// Check for Base64 encoding
	base64Pattern := regexp.MustCompile(`[A-Za-z0-9+/]{20,}={0,2}`)
	matches := base64Pattern.FindAllIndex(content, -1)
	
	for _, match := range matches {
		threat := SecurityThreat{
			Type:        "obfuscation",
			Severity:    "medium",
			Description: "Base64 encoded data detected",
			Position:    int64(match[0]),
			Pattern:     "base64",
			Confidence:  0.6,
			Mitigation:  "Decode and analyze the base64 content",
		}
		threats = append(threats, threat)
	}
	
	// Check for hex encoding
	hexPattern := regexp.MustCompile(`(?i)[0-9a-f]{40,}`)
	hexMatches := hexPattern.FindAllIndex(content, -1)
	
	for _, match := range hexMatches {
		threat := SecurityThreat{
			Type:        "obfuscation",
			Severity:    "medium",
			Description: "Long hex string detected (possible encoded data)",
			Position:    int64(match[0]),
			Pattern:     "hex_encoding",
			Confidence:  0.5,
			Mitigation:  "Decode and analyze the hex content",
		}
		threats = append(threats, threat)
	}
	
	// Check for excessive string concatenation (JavaScript obfuscation)
	jsObfuscPattern := regexp.MustCompile(`["\'][^"\']{1,3}["\'][\s]*\+[\s]*["\'][^"\']{1,3}["\']`)
	jsMatches := jsObfuscPattern.FindAllIndex(content, -1)
	
	if len(jsMatches) > 10 {
		threat := SecurityThreat{
			Type:        "obfuscation",
			Severity:    "medium",
			Description: "JavaScript string obfuscation detected",
			Position:    int64(jsMatches[0][0]),
			Pattern:     "js_string_concat",
			Confidence:  0.7,
			Mitigation:  "Deobfuscate and analyze the JavaScript code",
		}
		threats = append(threats, threat)
	}
	
	return threats
}

func detectNetworkPatterns(content []byte) []SecurityThreat {
	var threats []SecurityThreat
	
	// IP address patterns
	ipPattern := regexp.MustCompile(`\b(?:\d{1,3}\.){3}\d{1,3}\b`)
	ipMatches := ipPattern.FindAllIndex(content, -1)
	
	for _, match := range ipMatches {
		ip := string(content[match[0]:match[1]])
		if isPrivateIP(ip) {
			continue // Skip private IPs
		}
		
		threat := SecurityThreat{
			Type:        "network",
			Severity:    "medium",
			Description: fmt.Sprintf("External IP address: %s", ip),
			Position:    int64(match[0]),
			Pattern:     ip,
			Confidence:  0.6,
			Mitigation:  "Verify if this IP connection is legitimate",
		}
		threats = append(threats, threat)
	}
	
	// URL patterns
	urlPattern := regexp.MustCompile(`https?://[^\s<>"{}|\\^` + "`" + `\[\]]+`)
	urlMatches := urlPattern.FindAllIndex(content, -1)
	
	for _, match := range urlMatches {
		url := string(content[match[0]:match[1]])
		
		threat := SecurityThreat{
			Type:        "network",
			Severity:    "medium",
			Description: fmt.Sprintf("External URL: %s", url),
			Position:    int64(match[0]),
			Pattern:     url,
			Confidence:  0.5,
			Mitigation:  "Verify the legitimacy of this URL",
		}
		threats = append(threats, threat)
	}
	
	return threats
}

func detectCryptoPatterns(content []byte) []SecurityThreat {
	var threats []SecurityThreat
	
	cryptoPatterns := []struct {
		pattern     string
		description string
		severity    string
	}{
		{"CryptAcquireContext", "Cryptographic context creation", "medium"},
		{"CryptGenKey", "Cryptographic key generation", "medium"},
		{"CryptEncrypt", "Data encryption", "medium"},
		{"CryptDecrypt", "Data decryption", "medium"},
		{"BCryptGenRandom", "Random number generation", "low"},
		{"Crypt32.dll", "Windows crypto library", "medium"},
	}
	
	contentStr := strings.ToLower(string(content))
	
	for _, p := range cryptoPatterns {
		pattern := strings.ToLower(p.pattern)
		if strings.Contains(contentStr, pattern) {
			pos := strings.Index(contentStr, pattern)
			
			threat := SecurityThreat{
				Type:        "crypto",
				Severity:    p.severity,
				Description: p.description,
				Position:    int64(pos),
				Pattern:     p.pattern,
				Confidence:  0.6,
				Mitigation:  "Review cryptographic usage for security",
			}
			threats = append(threats, threat)
		}
	}
	
	return threats
}

func detectAntiAnalysisPatterns(content []byte) []SecurityThreat {
	var threats []SecurityThreat
	
	antiAnalysisPatterns := []struct {
		pattern     string
		description string
		severity    string
	}{
		{"IsDebuggerPresent", "Debugger detection", "high"},
		{"CheckRemoteDebuggerPresent", "Remote debugger detection", "high"},
		{"NtQueryInformationProcess", "Process information query", "medium"},
		{"ZwSetInformationThread", "Thread information manipulation", "high"},
		{"GetTickCount", "Timing analysis (sandbox evasion)", "medium"},
		{"Sleep", "Delay execution (sandbox evasion)", "medium"},
		{"VirtualBox", "Virtual machine detection", "high"},
		{"VMware", "Virtual machine detection", "high"},
		{"QEMU", "Virtual machine detection", "high"},
		{"SandBox", "Sandbox detection", "high"},
	}
	
	contentStr := strings.ToLower(string(content))
	
	for _, p := range antiAnalysisPatterns {
		pattern := strings.ToLower(p.pattern)
		if strings.Contains(contentStr, pattern) {
			pos := strings.Index(contentStr, pattern)
			
			threat := SecurityThreat{
				Type:        "anti_analysis",
				Severity:    p.severity,
				Description: p.description,
				Position:    int64(pos),
				Pattern:     p.pattern,
				Confidence:  0.8,
				Mitigation:  "This may indicate evasion techniques",
			}
			threats = append(threats, threat)
		}
	}
	
	return threats
}

func isPrivateIP(ip string) bool {
	privateRanges := []string{
		"10.", "172.16.", "172.17.", "172.18.", "172.19.",
		"172.20.", "172.21.", "172.22.", "172.23.", "172.24.",
		"172.25.", "172.26.", "172.27.", "172.28.", "172.29.",
		"172.30.", "172.31.", "192.168.", "127.", "169.254.",
	}
	
	for _, prefix := range privateRanges {
		if strings.HasPrefix(ip, prefix) {
			return true
		}
	}
	
	return false
}

func AnalyzeExecutableHeaders(exePath string) (ExecutableInfo, error) {
	info := ExecutableInfo{
		Sections: make([]Section, 0),
		Imports:  make([]ExecutableImport, 0),
		Exports:  make([]Export, 0),
	}
	
	file, err := os.Open(exePath)
	if err != nil {
		return info, err
	}
	defer file.Close()
	
	// Read file header
	header := make([]byte, 64)
	_, err = file.Read(header)
	if err != nil {
		return info, err
	}
	
	// Detect executable format
	if string(header[:4]) == "\x7fELF" {
		info.Format = "ELF"
		return analyzeELFHeaders(file, info)
	} else if string(header[:2]) == "MZ" {
		info.Format = "PE"
		return analyzePEHeaders(file, info)
	} else if string(header[:4]) == "\xFE\xED\xFA\xCE" || string(header[:4]) == "\xCE\xFA\xED\xFE" {
		info.Format = "Mach-O"
		return analyzeMachOHeaders(file, info)
	}
	
	return info, fmt.Errorf("unknown executable format")
}

func analyzeELFHeaders(file *os.File, info ExecutableInfo) (ExecutableInfo, error) {
	file.Seek(0, 0)
	
	// Read ELF header
	elfHeader := make([]byte, 64)
	_, err := file.Read(elfHeader)
	if err != nil {
		return info, err
	}
	
	// Parse ELF header
	class := elfHeader[4] // 32-bit or 64-bit
	if class == 1 {
		info.Architecture = "32-bit"
	} else if class == 2 {
		info.Architecture = "64-bit"
	}
	
	// Read entry point
	if class == 2 { // 64-bit
		info.Entrypoint = binary.LittleEndian.Uint64(elfHeader[24:32])
	} else { // 32-bit
		info.Entrypoint = uint64(binary.LittleEndian.Uint32(elfHeader[24:28]))
	}
	
	// Read section headers (simplified)
	shoff := binary.LittleEndian.Uint64(elfHeader[40:48])
	shentsize := binary.LittleEndian.Uint16(elfHeader[58:60])
	shnum := binary.LittleEndian.Uint16(elfHeader[60:62])
	
	// Read sections
	file.Seek(int64(shoff), 0)
	for i := 0; i < int(shnum); i++ {
		sectionHeader := make([]byte, shentsize)
		_, err := file.Read(sectionHeader)
		if err != nil {
			break
		}
		
		section := Section{
			Name:        fmt.Sprintf("section_%d", i),
			VirtualAddr: binary.LittleEndian.Uint64(sectionHeader[16:24]),
			RawSize:     binary.LittleEndian.Uint32(sectionHeader[32:36]),
		}
		
		// Parse section flags
		flags := binary.LittleEndian.Uint64(sectionHeader[8:16])
		section.IsExecutable = (flags & 0x4) != 0
		section.IsWritable = (flags & 0x1) != 0
		section.IsReadable = true // ELF sections are generally readable
		
		info.Sections = append(info.Sections, section)
	}
	
	return info, nil
}

func analyzePEHeaders(file *os.File, info ExecutableInfo) (ExecutableInfo, error) {
	file.Seek(0, 0)
	
	// Read DOS header
	dosHeader := make([]byte, 64)
	_, err := file.Read(dosHeader)
	if err != nil {
		return info, err
	}
	
	// Get PE header offset
	peOffset := binary.LittleEndian.Uint32(dosHeader[60:64])
	
	// Read PE header
	file.Seek(int64(peOffset), 0)
	peSignature := make([]byte, 4)
	_, err = file.Read(peSignature)
	if err != nil {
		return info, err
	}
	
	if string(peSignature) != "PE\x00\x00" {
		return info, fmt.Errorf("invalid PE signature")
	}
	
	// Read COFF header
	coffHeader := make([]byte, 20)
	_, err = file.Read(coffHeader)
	if err != nil {
		return info, err
	}
	
	machine := binary.LittleEndian.Uint16(coffHeader[0:2])
	switch machine {
	case 0x014c:
		info.Architecture = "i386"
	case 0x8664:
		info.Architecture = "x86_64"
	case 0x01c0:
		info.Architecture = "ARM"
	default:
		info.Architecture = fmt.Sprintf("unknown(0x%x)", machine)
	}
	
	numberOfSections := binary.LittleEndian.Uint16(coffHeader[2:4])
	info.Timestamp = binary.LittleEndian.Uint32(coffHeader[4:8])
	
	// Read Optional header
	optionalHeaderSize := binary.LittleEndian.Uint16(coffHeader[16:18])
	optionalHeader := make([]byte, optionalHeaderSize)
	_, err = file.Read(optionalHeader)
	if err != nil {
		return info, err
	}
	
	// Get entry point
	info.Entrypoint = uint64(binary.LittleEndian.Uint32(optionalHeader[16:20]))
	
	// Read sections
	for i := 0; i < int(numberOfSections); i++ {
		sectionHeader := make([]byte, 40)
		_, err := file.Read(sectionHeader)
		if err != nil {
			break
		}
		
		// Extract section name
		nameBytes := sectionHeader[0:8]
		name := string(bytes.TrimRight(nameBytes, "\x00"))
		
		section := Section{
			Name:        name,
			VirtualAddr: uint64(binary.LittleEndian.Uint32(sectionHeader[12:16])),
			RawSize:     binary.LittleEndian.Uint32(sectionHeader[16:20]),
		}
		
		// Parse section characteristics
		characteristics := binary.LittleEndian.Uint32(sectionHeader[36:40])
		section.IsExecutable = (characteristics & 0x20000000) != 0
		section.IsWritable = (characteristics & 0x80000000) != 0
		section.IsReadable = (characteristics & 0x40000000) != 0
		
		info.Sections = append(info.Sections, section)
	}
	
	// Check for packing indicators
	info.IsPacked = checkPEPacking(info.Sections)
	
	return info, nil
}

func analyzeMachOHeaders(file *os.File, info ExecutableInfo) (ExecutableInfo, error) {
	file.Seek(0, 0)
	
	// Read Mach-O header
	header := make([]byte, 32)
	_, err := file.Read(header)
	if err != nil {
		return info, err
	}
	
	magic := binary.BigEndian.Uint32(header[0:4])
	switch magic {
	case 0xFEEDFACE:
		info.Architecture = "32-bit"
	case 0xFEEDFACF:
		info.Architecture = "64-bit"
	}
	
	// Basic Mach-O analysis (simplified)
	info.Entrypoint = 0 // Would need to parse load commands for actual entry point
	
	return info, nil
}

func checkPEPacking(sections []Section) bool {
	// Check for common packing indicators
	for _, section := range sections {
		// Very small number of sections
		if len(sections) < 3 {
			return true
		}
		
		// Sections with suspicious names
		suspiciousNames := []string{"UPX", "ASPack", "PECompact", "FSG", "Themida"}
		for _, name := range suspiciousNames {
			if strings.Contains(strings.ToUpper(section.Name), name) {
				return true
			}
		}
		
		// Very high virtual size vs raw size ratio
		if section.RawSize > 0 && float64(section.VirtualAddr)/float64(section.RawSize) > 10 {
			return true
		}
	}
	
	return false
}

func FindEmbeddedFiles(filePath string) ([]EmbeddedFile, error) {
	var embedded []EmbeddedFile
	
	file, err := os.Open(filePath)
	if err != nil {
		return embedded, err
	}
	defer file.Close()
	
	content, err := io.ReadAll(file)
	if err != nil {
		return embedded, err
	}
	
	// Look for common file signatures
	signatures := map[string][]byte{
		"ZIP":  {0x50, 0x4B, 0x03, 0x04},
		"PDF":  {0x25, 0x50, 0x44, 0x46},
		"PNG":  {0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A},
		"JPEG": {0xFF, 0xD8, 0xFF},
		"EXE":  {0x4D, 0x5A},
		"GIF":  {0x47, 0x49, 0x46, 0x38},
	}
	
	for fileType, signature := range signatures {
		offset := 0
		for {
			pos := bytes.Index(content[offset:], signature)
			if pos == -1 {
				break
			}
			
			actualPos := offset + pos
			
			// Try to determine file size (simplified)
			size := estimateEmbeddedFileSize(content, actualPos, fileType)
			
			embeddedFile := EmbeddedFile{
				Type:     fileType,
				Offset:   int64(actualPos),
				Size:     size,
				Filename: fmt.Sprintf("embedded_%s_%d", strings.ToLower(fileType), actualPos),
			}
			
			// Calculate hash of embedded file
			if size > 0 && actualPos+int(size) <= len(content) {
				// Would hash the embedded data in production
				embeddedFile.Hash = "placeholder_hash"
			}
			
			embedded = append(embedded, embeddedFile)
			offset = actualPos + len(signature)
		}
	}
	
	return embedded, nil
}

func estimateEmbeddedFileSize(content []byte, offset int, fileType string) int64 {
	// Simplified size estimation
	switch fileType {
	case "ZIP":
		// For ZIP, we'd need to parse the directory structure
		return 1024 // Placeholder
	case "PDF":
		// Look for %%EOF marker
		eofMarker := []byte("%%EOF")
		pos := bytes.Index(content[offset:], eofMarker)
		if pos != -1 {
			return int64(pos + len(eofMarker))
		}
		return 1024
	case "PNG":
		// Look for IEND chunk
		iendMarker := []byte("IEND")
		pos := bytes.Index(content[offset:], iendMarker)
		if pos != -1 {
			return int64(pos + len(iendMarker) + 4) // +4 for CRC
		}
		return 1024
	default:
		return 1024 // Default size
	}
}

func CheckFilePermissions(filePath string) (PermissionReport, error) {
	report := PermissionReport{
		Warnings:    make([]string, 0),
		Suggestions: make([]string, 0),
	}
	
	info, err := os.Stat(filePath)
	if err != nil {
		return report, err
	}
	
	mode := info.Mode()
	perm := mode.Perm()
	
	// Parse permissions
	report.Owner.Read = (perm & 0400) != 0
	report.Owner.Write = (perm & 0200) != 0
	report.Owner.Execute = (perm & 0100) != 0
	
	report.Group.Read = (perm & 0040) != 0
	report.Group.Write = (perm & 0020) != 0
	report.Group.Execute = (perm & 0010) != 0
	
	report.Others.Read = (perm & 0004) != 0
	report.Others.Write = (perm & 0002) != 0
	report.Others.Execute = (perm & 0001) != 0
	
	// Convert to octal
	report.Octal = fmt.Sprintf("%o", perm)
	
	// Security analysis
	report.IsSecure = true
	
	// Check for overly permissive permissions
	if report.Others.Write {
		report.IsSecure = false
		report.Warnings = append(report.Warnings, "File is world-writable")
		report.Suggestions = append(report.Suggestions, "Remove write permission for others")
	}
	
	if report.Others.Execute && !info.IsDir() {
		report.Warnings = append(report.Warnings, "File is executable by others")
		report.Suggestions = append(report.Suggestions, "Consider removing execute permission for others")
	}
	
	// Check for group write on sensitive files
	if report.Group.Write && (strings.HasSuffix(filePath, ".conf") || 
		strings.HasSuffix(filePath, ".key") || 
		strings.HasSuffix(filePath, ".pem")) {
		report.Warnings = append(report.Warnings, "Sensitive file is group-writable")
		report.Suggestions = append(report.Suggestions, "Remove group write permission for sensitive files")
	}
	
	// Check for missing execute on directories
	if info.IsDir() && (!report.Owner.Execute || !report.Group.Execute) {
		report.Warnings = append(report.Warnings, "Directory lacks execute permission")
		report.Suggestions = append(report.Suggestions, "Add execute permission for directory access")
	}
	
	return report, nil
}

func ScanForViruses(filePath string, signatures []VirusSignature) ([]Threat, error) {
	var threats []Threat
	
	file, err := os.Open(filePath)
	if err != nil {
		return threats, err
	}
	defer file.Close()
	
	// Read file in chunks for memory efficiency
	const chunkSize = 64 * 1024
	buffer := make([]byte, chunkSize)
	offset := int64(0)
	
	for {
		n, err := file.Read(buffer)
		if err == io.EOF {
			break
		}
		if err != nil {
			return threats, err
		}
		
		// Scan current chunk against all signatures
		for _, signature := range signatures {
			pos := bytes.Index(buffer[:n], signature.Pattern)
			if pos != -1 {
				threat := Threat{
					Name:        signature.Name,
					Type:        "virus",
					Severity:    signature.Severity,
					Description: signature.Description,
					Position:    offset + int64(pos),
					Confidence:  1.0, // Exact signature match
				}
				threats = append(threats, threat)
			}
		}
		
		offset += int64(n)
		
		// If we read less than chunk size, we're at EOF
		if n < chunkSize {
			break
		}
	}
	
	return threats, nil
}

// Helper function to create common virus signatures
func GetCommonVirusSignatures() []VirusSignature {
	return []VirusSignature{
		{
			Name:        "EICAR-Test-File",
			Pattern:     []byte("X5O!P%@AP[4\\PZX54(P^)7CC)7}$EICAR-STANDARD-ANTIVIRUS-TEST-FILE!$H+H*"),
			Description: "EICAR antivirus test file",
			Severity:    "test",
		},
		{
			Name:        "Suspicious-Shellcode-Pattern",
			Pattern:     []byte{0x90, 0x90, 0x90, 0x90, 0x31, 0xC0},
			Description: "Common shellcode pattern (NOP sled + XOR)",
			Severity:    "medium",
		},
		{
			Name:        "PE-Injection-Pattern",
			Pattern:     []byte("CreateRemoteThread"),
			Description: "Process injection API usage",
			Severity:    "high",
		},
	}
}