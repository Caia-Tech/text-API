package textlib

import (
	"math"
	"regexp"
	"sort"
	"strings"
	"unicode"
)

// Segment represents a text segment with metadata
type Segment struct {
	Text         string
	Start        int
	End          int
	Type         string // "paragraph", "section", "topic", "semantic", etc.
	TokenCount   int
	CharCount    int
	Metadata     map[string]interface{}
}

// ChunkingStrategy defines how text should be chunked
type ChunkingStrategy struct {
	Method          string  // "token", "sentence", "paragraph", "semantic", "sliding"
	MaxTokens       int
	MaxCharacters   int
	Overlap         int     // For sliding window
	PreserveWords   bool    // Don't split mid-word
	PreserveSentences bool  // Don't split mid-sentence
	SemanticThreshold float64 // For semantic chunking
}

// SegmentationResult contains segmented text
type SegmentationResult struct {
	Segments         []Segment
	SegmentationMethod string
	TotalSegments    int
	AverageSegmentSize int
	Boundaries       []int // Character positions of segment boundaries
}

// TextChunk represents a chunk of text for processing
type TextChunk struct {
	Content      string
	Index        int
	StartOffset  int
	EndOffset    int
	TokenCount   int
	Metadata     map[string]interface{}
	
	// For overlapping chunks
	OverlapStart int
	OverlapEnd   int
}

// SegmentText performs intelligent text segmentation
func SegmentText(text string, method string) *SegmentationResult {
	result := &SegmentationResult{
		Segments:           []Segment{},
		SegmentationMethod: method,
		Boundaries:         []int{},
	}
	
	switch method {
	case "paragraph":
		result.Segments = segmentByParagraphs(text)
	case "section":
		result.Segments = segmentBySections(text)
	case "topic":
		result.Segments = segmentByTopics(text)
	case "semantic":
		result.Segments = segmentBySemantic(text)
	case "sentence":
		result.Segments = segmentBySentences(text)
	case "fixed":
		result.Segments = segmentByFixedSize(text, 1000) // 1000 chars default
	default:
		result.Segments = segmentByParagraphs(text)
	}
	
	// Calculate statistics
	result.TotalSegments = len(result.Segments)
	if result.TotalSegments > 0 {
		totalSize := 0
		for _, seg := range result.Segments {
			totalSize += seg.CharCount
			result.Boundaries = append(result.Boundaries, seg.Start)
		}
		result.AverageSegmentSize = totalSize / result.TotalSegments
	}
	
	return result
}

// ChunkText divides text into chunks based on strategy
func ChunkText(text string, strategy ChunkingStrategy) []TextChunk {
	switch strategy.Method {
	case "token":
		return chunkByTokens(text, strategy)
	case "sentence":
		return chunkBySentences(text, strategy)
	case "paragraph":
		return chunkByParagraphs(text, strategy)
	case "semantic":
		return chunkBySemantic(text, strategy)
	case "sliding":
		return chunkBySliding(text, strategy)
	default:
		return chunkByTokens(text, strategy)
	}
}

// Segmentation methods

func segmentByParagraphs(text string) []Segment {
	segments := []Segment{}
	paragraphs := SplitIntoParagraphs(text)
	
	offset := 0
	for i, para := range paragraphs {
		// Find actual position in original text
		start := strings.Index(text[offset:], para)
		if start == -1 {
			continue
		}
		start += offset
		
		segment := Segment{
			Text:       para,
			Start:      start,
			End:        start + len(para),
			Type:       "paragraph",
			TokenCount: len(strings.Fields(para)),
			CharCount:  len(para),
			Metadata: map[string]interface{}{
				"index": i,
			},
		}
		
		segments = append(segments, segment)
		offset = start + len(para)
	}
	
	return segments
}

func segmentBySections(text string) []Segment {
	segments := []Segment{}
	
	// Look for section markers (headers, chapter markers, etc.)
	lines := strings.Split(text, "\n")
	currentSection := []string{}
	sectionStart := 0
	charOffset := 0
	
	for i, line := range lines {
		if isSectionHeader(line) && len(currentSection) > 0 {
			// End current section
			sectionText := strings.Join(currentSection, "\n")
			segment := Segment{
				Text:       sectionText,
				Start:      sectionStart,
				End:        charOffset,
				Type:       "section",
				TokenCount: len(strings.Fields(sectionText)),
				CharCount:  len(sectionText),
				Metadata: map[string]interface{}{
					"header": strings.TrimSpace(lines[i-len(currentSection)]),
				},
			}
			segments = append(segments, segment)
			
			// Start new section
			currentSection = []string{line}
			sectionStart = charOffset
		} else {
			currentSection = append(currentSection, line)
		}
		
		charOffset += len(line) + 1 // +1 for newline
	}
	
	// Add final section
	if len(currentSection) > 0 {
		sectionText := strings.Join(currentSection, "\n")
		segment := Segment{
			Text:       sectionText,
			Start:      sectionStart,
			End:        charOffset,
			Type:       "section",
			TokenCount: len(strings.Fields(sectionText)),
			CharCount:  len(sectionText),
			Metadata:   map[string]interface{}{},
		}
		segments = append(segments, segment)
	}
	
	// If no sections found, treat whole text as one section
	if len(segments) == 0 {
		segments = append(segments, Segment{
			Text:       text,
			Start:      0,
			End:        len(text),
			Type:       "section",
			TokenCount: len(strings.Fields(text)),
			CharCount:  len(text),
			Metadata:   map[string]interface{}{},
		})
	}
	
	return segments
}

func isSectionHeader(line string) bool {
	line = strings.TrimSpace(line)
	
	// Empty line cannot be header
	if line == "" {
		return false
	}
	
	// Check for markdown headers
	if strings.HasPrefix(line, "#") {
		return true
	}
	
	// Check for numbered sections
	if matched, _ := regexp.MatchString(`^\d+\.?\s+[A-Z]`, line); matched {
		return true
	}
	
	// Check for all caps (common section marker)
	if len(line) > 3 && strings.ToUpper(line) == line && !strings.ContainsAny(line, ".!?") {
		return true
	}
	
	// Check for common section keywords
	sectionKeywords := []string{"Chapter", "Section", "Part", "Introduction", "Conclusion"}
	for _, keyword := range sectionKeywords {
		if strings.HasPrefix(line, keyword) {
			return true
		}
	}
	
	return false
}

func segmentByTopics(text string) []Segment {
	segments := []Segment{}
	sentences := SplitIntoSentences(text)
	
	if len(sentences) == 0 {
		return segments
	}
	
	// Group sentences by topic similarity
	currentTopic := []string{sentences[0]}
	topicStart := 0
	charOffset := 0
	
	for i := 1; i < len(sentences); i++ {
		// Check topic similarity
		similarity := calculateTopicSimilarity(
			strings.Join(currentTopic, " "),
			sentences[i],
		)
		
		if similarity < 0.3 { // Topic shift threshold
			// Create segment for current topic
			topicText := strings.Join(currentTopic, " ")
			segment := Segment{
				Text:       topicText,
				Start:      topicStart,
				End:        charOffset,
				Type:       "topic",
				TokenCount: len(strings.Fields(topicText)),
				CharCount:  len(topicText),
				Metadata: map[string]interface{}{
					"sentenceCount": len(currentTopic),
				},
			}
			segments = append(segments, segment)
			
			// Start new topic
			currentTopic = []string{sentences[i]}
			topicStart = charOffset
		} else {
			currentTopic = append(currentTopic, sentences[i])
		}
		
		charOffset += len(sentences[i]) + 1
	}
	
	// Add final topic
	if len(currentTopic) > 0 {
		topicText := strings.Join(currentTopic, " ")
		segment := Segment{
			Text:       topicText,
			Start:      topicStart,
			End:        charOffset,
			Type:       "topic",
			TokenCount: len(strings.Fields(topicText)),
			CharCount:  len(topicText),
			Metadata: map[string]interface{}{
				"sentenceCount": len(currentTopic),
			},
		}
		segments = append(segments, segment)
	}
	
	return segments
}

func calculateTopicSimilarity(text1, text2 string) float64 {
	// Extract key terms
	terms1 := extractKeyTerms(text1)
	terms2 := extractKeyTerms(text2)
	
	if len(terms1) == 0 || len(terms2) == 0 {
		return 0
	}
	
	// Calculate Jaccard similarity
	intersection := 0
	set1 := make(map[string]bool)
	
	for _, term := range terms1 {
		set1[term] = true
	}
	
	for _, term := range terms2 {
		if set1[term] {
			intersection++
		}
	}
	
	union := len(terms1) + len(terms2) - intersection
	
	return float64(intersection) / float64(union)
}

func extractKeyTerms(text string) []string {
	// Simple key term extraction
	words := strings.Fields(strings.ToLower(text))
	terms := []string{}
	
	for _, word := range words {
		cleaned := strings.Trim(word, ".,!?;:")
		// Keep nouns and longer words
		if len(cleaned) > 4 && !isStopWord(cleaned) {
			terms = append(terms, cleaned)
		}
	}
	
	return terms
}

func segmentBySemantic(text string) []Segment {
	segments := []Segment{}
	sentences := SplitIntoSentences(text)
	
	if len(sentences) == 0 {
		return segments
	}
	
	// Use semantic coherence to group sentences
	currentGroup := []string{sentences[0]}
	groupStart := 0
	charOffset := len(sentences[0])
	
	for i := 1; i < len(sentences); i++ {
		// Calculate semantic coherence
		coherence := calculateSemanticCoherence(currentGroup, sentences[i])
		
		if coherence < 0.5 { // Semantic break threshold
			// Create segment
			groupText := strings.Join(currentGroup, " ")
			segment := Segment{
				Text:       groupText,
				Start:      groupStart,
				End:        charOffset,
				Type:       "semantic",
				TokenCount: len(strings.Fields(groupText)),
				CharCount:  len(groupText),
				Metadata: map[string]interface{}{
					"coherence": coherence,
				},
			}
			segments = append(segments, segment)
			
			// Start new group
			currentGroup = []string{sentences[i]}
			groupStart = charOffset + 1
		} else {
			currentGroup = append(currentGroup, sentences[i])
		}
		
		charOffset += len(sentences[i]) + 1
	}
	
	// Add final group
	if len(currentGroup) > 0 {
		groupText := strings.Join(currentGroup, " ")
		segment := Segment{
			Text:       groupText,
			Start:      groupStart,
			End:        charOffset,
			Type:       "semantic",
			TokenCount: len(strings.Fields(groupText)),
			CharCount:  len(groupText),
			Metadata:   map[string]interface{}{},
		}
		segments = append(segments, segment)
	}
	
	return segments
}

func calculateSemanticCoherence(group []string, newSentence string) float64 {
	// Simple coherence based on word overlap and entity continuity
	groupText := strings.Join(group, " ")
	
	// Word overlap
	groupWords := extractContentWords(groupText)
	sentWords := extractContentWords(newSentence)
	
	overlap := calculateWordOverlap(groupWords, sentWords)
	
	// Entity continuity
	groupEntities := ExtractNamedEntities(groupText)
	sentEntities := ExtractNamedEntities(newSentence)
	
	entityOverlap := 0.0
	if len(groupEntities) > 0 && len(sentEntities) > 0 {
		for _, e1 := range groupEntities {
			for _, e2 := range sentEntities {
				if e1.Text == e2.Text || e1.Type == e2.Type {
					entityOverlap = 1.0
					break
				}
			}
		}
	}
	
	// Combine metrics
	return (overlap + entityOverlap) / 2
}

func segmentBySentences(text string) []Segment {
	segments := []Segment{}
	sentences := SplitIntoSentences(text)
	
	charOffset := 0
	for i, sent := range sentences {
		start := strings.Index(text[charOffset:], sent)
		if start == -1 {
			continue
		}
		start += charOffset
		
		segment := Segment{
			Text:       sent,
			Start:      start,
			End:        start + len(sent),
			Type:       "sentence",
			TokenCount: len(strings.Fields(sent)),
			CharCount:  len(sent),
			Metadata: map[string]interface{}{
				"index": i,
			},
		}
		
		segments = append(segments, segment)
		charOffset = start + len(sent)
	}
	
	return segments
}

func segmentByFixedSize(text string, size int) []Segment {
	segments := []Segment{}
	
	for i := 0; i < len(text); i += size {
		end := min(i+size, len(text))
		
		// Try to break at word boundary
		if end < len(text) && !unicode.IsSpace(rune(text[end])) {
			// Look for nearest space
			for j := end; j > i && j > end-50; j-- {
				if unicode.IsSpace(rune(text[j])) {
					end = j
					break
				}
			}
		}
		
		segmentText := text[i:end]
		segment := Segment{
			Text:       segmentText,
			Start:      i,
			End:        end,
			Type:       "fixed",
			TokenCount: len(strings.Fields(segmentText)),
			CharCount:  len(segmentText),
			Metadata: map[string]interface{}{
				"size": size,
			},
		}
		
		segments = append(segments, segment)
	}
	
	return segments
}

// Chunking methods

func chunkByTokens(text string, strategy ChunkingStrategy) []TextChunk {
	chunks := []TextChunk{}
	words := strings.Fields(text)
	
	if len(words) == 0 {
		return chunks
	}
	
	chunkSize := strategy.MaxTokens
	if chunkSize <= 0 {
		chunkSize = 100 // Default
	}
	
	for i := 0; i < len(words); i += chunkSize - strategy.Overlap {
		end := min(i+chunkSize, len(words))
		
		chunkWords := words[i:end]
		chunkText := strings.Join(chunkWords, " ")
		
		chunk := TextChunk{
			Content:     chunkText,
			Index:       len(chunks),
			StartOffset: i,
			EndOffset:   end,
			TokenCount:  len(chunkWords),
			Metadata:    map[string]interface{}{},
		}
		
		// Mark overlap region
		if i > 0 && strategy.Overlap > 0 {
			chunk.OverlapStart = 0
			chunk.OverlapEnd = min(strategy.Overlap, len(chunkWords))
		}
		
		chunks = append(chunks, chunk)
		
		// Stop if we've processed all words
		if end >= len(words) {
			break
		}
	}
	
	return chunks
}

func chunkBySentences(text string, strategy ChunkingStrategy) []TextChunk {
	chunks := []TextChunk{}
	sentences := SplitIntoSentences(text)
	
	if len(sentences) == 0 {
		return chunks
	}
	
	currentChunk := []string{}
	currentTokens := 0
	chunkStart := 0
	
	for i, sent := range sentences {
		sentTokens := len(strings.Fields(sent))
		
		// Check if adding this sentence would exceed limit
		if currentTokens+sentTokens > strategy.MaxTokens && len(currentChunk) > 0 {
			// Create chunk
			chunkText := strings.Join(currentChunk, " ")
			chunk := TextChunk{
				Content:     chunkText,
				Index:       len(chunks),
				StartOffset: chunkStart,
				EndOffset:   i,
				TokenCount:  currentTokens,
				Metadata: map[string]interface{}{
					"sentences": len(currentChunk),
				},
			}
			chunks = append(chunks, chunk)
			
			// Start new chunk
			currentChunk = []string{sent}
			currentTokens = sentTokens
			chunkStart = i
		} else {
			currentChunk = append(currentChunk, sent)
			currentTokens += sentTokens
		}
	}
	
	// Add final chunk
	if len(currentChunk) > 0 {
		chunkText := strings.Join(currentChunk, " ")
		chunk := TextChunk{
			Content:     chunkText,
			Index:       len(chunks),
			StartOffset: chunkStart,
			EndOffset:   len(sentences),
			TokenCount:  currentTokens,
			Metadata: map[string]interface{}{
				"sentences": len(currentChunk),
			},
		}
		chunks = append(chunks, chunk)
	}
	
	return chunks
}

func chunkByParagraphs(text string, strategy ChunkingStrategy) []TextChunk {
	chunks := []TextChunk{}
	paragraphs := SplitIntoParagraphs(text)
	
	if len(paragraphs) == 0 {
		return chunks
	}
	
	currentChunk := []string{}
	currentTokens := 0
	chunkStart := 0
	
	for i, para := range paragraphs {
		paraTokens := len(strings.Fields(para))
		
		// Check if adding this paragraph would exceed limit
		if currentTokens+paraTokens > strategy.MaxTokens && len(currentChunk) > 0 {
			// Create chunk
			chunkText := strings.Join(currentChunk, "\n\n")
			chunk := TextChunk{
				Content:     chunkText,
				Index:       len(chunks),
				StartOffset: chunkStart,
				EndOffset:   i,
				TokenCount:  currentTokens,
				Metadata: map[string]interface{}{
					"paragraphs": len(currentChunk),
				},
			}
			chunks = append(chunks, chunk)
			
			// Start new chunk
			currentChunk = []string{para}
			currentTokens = paraTokens
			chunkStart = i
		} else {
			currentChunk = append(currentChunk, para)
			currentTokens += paraTokens
		}
	}
	
	// Add final chunk
	if len(currentChunk) > 0 {
		chunkText := strings.Join(currentChunk, "\n\n")
		chunk := TextChunk{
			Content:     chunkText,
			Index:       len(chunks),
			StartOffset: chunkStart,
			EndOffset:   len(paragraphs),
			TokenCount:  currentTokens,
			Metadata: map[string]interface{}{
				"paragraphs": len(currentChunk),
			},
		}
		chunks = append(chunks, chunk)
	}
	
	return chunks
}

func chunkBySemantic(text string, strategy ChunkingStrategy) []TextChunk {
	chunks := []TextChunk{}
	sentences := SplitIntoSentences(text)
	
	if len(sentences) == 0 {
		return chunks
	}
	
	// Build semantic groups
	semanticGroups := [][]string{}
	currentGroup := []string{sentences[0]}
	
	for i := 1; i < len(sentences); i++ {
		coherence := calculateSemanticCoherence(currentGroup, sentences[i])
		
		if coherence < strategy.SemanticThreshold {
			// Start new group
			semanticGroups = append(semanticGroups, currentGroup)
			currentGroup = []string{sentences[i]}
		} else {
			currentGroup = append(currentGroup, sentences[i])
		}
	}
	
	// Add final group
	if len(currentGroup) > 0 {
		semanticGroups = append(semanticGroups, currentGroup)
	}
	
	// Convert groups to chunks respecting token limit
	currentChunk := []string{}
	currentTokens := 0
	
	for _, group := range semanticGroups {
		groupText := strings.Join(group, " ")
		groupTokens := len(strings.Fields(groupText))
		
		if currentTokens+groupTokens > strategy.MaxTokens && len(currentChunk) > 0 {
			// Create chunk
			chunkText := strings.Join(currentChunk, " ")
			chunk := TextChunk{
				Content:    chunkText,
				Index:      len(chunks),
				TokenCount: currentTokens,
				Metadata: map[string]interface{}{
					"semanticGroups": len(currentChunk),
				},
			}
			chunks = append(chunks, chunk)
			
			currentChunk = []string{groupText}
			currentTokens = groupTokens
		} else {
			currentChunk = append(currentChunk, groupText)
			currentTokens += groupTokens
		}
	}
	
	// Add final chunk
	if len(currentChunk) > 0 {
		chunkText := strings.Join(currentChunk, " ")
		chunk := TextChunk{
			Content:    chunkText,
			Index:      len(chunks),
			TokenCount: currentTokens,
			Metadata: map[string]interface{}{
				"semanticGroups": len(currentChunk),
			},
		}
		chunks = append(chunks, chunk)
	}
	
	return chunks
}

func chunkBySliding(text string, strategy ChunkingStrategy) []TextChunk {
	chunks := []TextChunk{}
	words := strings.Fields(text)
	
	if len(words) == 0 {
		return chunks
	}
	
	windowSize := strategy.MaxTokens
	if windowSize <= 0 {
		windowSize = 100
	}
	
	stride := windowSize - strategy.Overlap
	if stride <= 0 {
		stride = 1
	}
	
	for i := 0; i < len(words); i += stride {
		end := min(i+windowSize, len(words))
		
		chunkWords := words[i:end]
		chunkText := strings.Join(chunkWords, " ")
		
		chunk := TextChunk{
			Content:     chunkText,
			Index:       len(chunks),
			StartOffset: i,
			EndOffset:   end,
			TokenCount:  len(chunkWords),
			Metadata: map[string]interface{}{
				"window": windowSize,
				"stride": stride,
			},
		}
		
		// Mark overlap regions
		if i > 0 {
			chunk.OverlapStart = 0
			chunk.OverlapEnd = min(strategy.Overlap, len(chunkWords))
		}
		
		chunks = append(chunks, chunk)
		
		// Stop if we've reached the end
		if end >= len(words) {
			break
		}
	}
	
	return chunks
}

// Utility functions for optimal chunking

// OptimalChunkSize calculates optimal chunk size based on text characteristics
func OptimalChunkSize(text string, targetChunks int) int {
	tokens := len(strings.Fields(text))
	
	if targetChunks <= 0 {
		targetChunks = 10
	}
	
	optimalSize := tokens / targetChunks
	
	// Round to reasonable boundaries
	if optimalSize < 50 {
		return 50
	} else if optimalSize > 1000 {
		return 1000
	}
	
	// Round to nearest 50
	return ((optimalSize + 25) / 50) * 50
}

// FindNaturalBoundaries identifies natural breaking points in text
func FindNaturalBoundaries(text string) []int {
	boundaries := []int{0}
	
	// Paragraph boundaries
	paragraphs := SplitIntoParagraphs(text)
	offset := 0
	
	for _, para := range paragraphs {
		pos := strings.Index(text[offset:], para)
		if pos != -1 {
			offset += pos + len(para)
			boundaries = append(boundaries, offset)
		}
	}
	
	// Section boundaries (if any)
	lines := strings.Split(text, "\n")
	charPos := 0
	
	for _, line := range lines {
		if isSectionHeader(line) {
			boundaries = append(boundaries, charPos)
		}
		charPos += len(line) + 1
	}
	
	// Sort and deduplicate
	uniqueBoundaries := make(map[int]bool)
	for _, b := range boundaries {
		uniqueBoundaries[b] = true
	}
	
	boundaries = []int{}
	for b := range uniqueBoundaries {
		boundaries = append(boundaries, b)
	}
	
	// Sort boundaries
	sort.Ints(boundaries)
	
	return boundaries
}

// MergeSmallSegments combines segments that are too small
func MergeSmallSegments(segments []Segment, minSize int) []Segment {
	if len(segments) <= 1 {
		return segments
	}
	
	merged := []Segment{}
	current := segments[0]
	
	for i := 1; i < len(segments); i++ {
		if current.CharCount < minSize && segments[i].Type == current.Type {
			// Merge with next segment
			current.Text += " " + segments[i].Text
			current.End = segments[i].End
			current.CharCount += segments[i].CharCount + 1
			current.TokenCount += segments[i].TokenCount
		} else {
			merged = append(merged, current)
			current = segments[i]
		}
	}
	
	// Add final segment
	merged = append(merged, current)
	
	return merged
}

// BalanceChunks redistributes content to create more evenly sized chunks
func BalanceChunks(chunks []TextChunk, targetVariance float64) []TextChunk {
	if len(chunks) <= 1 {
		return chunks
	}
	
	// Calculate current variance
	sizes := make([]float64, len(chunks))
	totalSize := 0.0
	
	for i, chunk := range chunks {
		sizes[i] = float64(chunk.TokenCount)
		totalSize += sizes[i]
	}
	
	mean := totalSize / float64(len(chunks))
	variance := 0.0
	
	for _, size := range sizes {
		diff := size - mean
		variance += diff * diff
	}
	variance /= float64(len(chunks))
	stdDev := math.Sqrt(variance)
	
	// If variance is acceptable, return as is
	if stdDev/mean < targetVariance {
		return chunks
	}
	
	// Rebalance chunks
	targetSize := int(mean)
	balanced := []TextChunk{}
	
	// Combine all content
	allWords := []string{}
	for _, chunk := range chunks {
		words := strings.Fields(chunk.Content)
		allWords = append(allWords, words...)
	}
	
	// Redistribute
	for i := 0; i < len(allWords); i += targetSize {
		end := min(i+targetSize, len(allWords))
		
		chunkWords := allWords[i:end]
		chunk := TextChunk{
			Content:     strings.Join(chunkWords, " "),
			Index:       len(balanced),
			StartOffset: i,
			EndOffset:   end,
			TokenCount:  len(chunkWords),
			Metadata: map[string]interface{}{
				"rebalanced": true,
			},
		}
		
		balanced = append(balanced, chunk)
	}
	
	return balanced
}