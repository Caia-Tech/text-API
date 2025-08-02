package textlib

import (
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
	"unicode"
)

// Extended entity types for RL and agentic systems
const (
	EntityPerson       = "PERSON"
	EntityOrganization = "ORGANIZATION"
	EntityLocation     = "LOCATION"
	EntityDate         = "DATE"
	EntityTime         = "TIME"
	EntityMoney        = "MONEY"
	EntityPercent      = "PERCENT"
	EntityProduct      = "PRODUCT"
	EntityEvent        = "EVENT"
	EntityEmail        = "EMAIL"
	EntityURL          = "URL"
	EntityPhone        = "PHONE"
	EntityNumber       = "NUMBER"
	EntityOrdinal      = "ORDINAL"
	EntityCardinal     = "CARDINAL"
	EntityCode         = "CODE"
	EntityAction       = "ACTION"
	EntityConcept      = "CONCEPT"
	EntityMetric       = "METRIC"
)

// AdvancedEntity represents an entity with confidence and context
type AdvancedEntity struct {
	Entity
	Confidence   float64
	Context      string   // Surrounding context
	Attributes   map[string]string
	Relationships []EntityRelation
}

// EntityRelation represents relationships between entities
type EntityRelation struct {
	FromEntity   string
	ToEntity     string
	RelationType string // "works_at", "located_in", "owns", etc.
	Confidence   float64
}

// Patterns for advanced entity detection
var (
	moneyPattern = regexp.MustCompile(`\$[\d,]+\.?\d*|\d+\s*(dollars?|cents?|euros?|pounds?|yen|yuan)`)
	percentPattern = regexp.MustCompile(`\d+\.?\d*\s*%|percent|percentage`)
	phonePattern = regexp.MustCompile(`(\+?\d{1,3}[-.\s]?)?\(?\d{3}\)?[-.\s]?\d{3}[-.\s]?\d{4}`)
	timePattern = regexp.MustCompile(`\d{1,2}:\d{2}(\s*(AM|PM|am|pm))?|\d{1,2}\s*(AM|PM|am|pm)`)
	codePattern = regexp.MustCompile("`[^`]+`|```[^`]+```")
	numberPattern = regexp.MustCompile(`\b\d+\.?\d*\b`)
	ordinalPattern = regexp.MustCompile(`\b\d+(st|nd|rd|th)\b|\b(first|second|third|fourth|fifth|sixth|seventh|eighth|ninth|tenth)\b`)
	metricPattern = regexp.MustCompile(`\b\d+\.?\d*\s*(ms|seconds?|minutes?|hours?|days?|weeks?|months?|years?|bytes?|KB|MB|GB|TB|km|m|cm|mm|kg|g|mg|°C|°F)\b`)
)

// ExtractAdvancedEntities performs comprehensive entity extraction
func ExtractAdvancedEntities(text string) []AdvancedEntity {
	var entities []AdvancedEntity
	
	// Extract basic entities first
	basicEntities := ExtractNamedEntities(text)
	for _, e := range basicEntities {
		entities = append(entities, AdvancedEntity{
			Entity:     e,
			Confidence: 0.8,
			Context:    extractContext(text, e.Position.Start, e.Position.End),
			Attributes: make(map[string]string),
		})
	}
	
	// Extract money entities
	moneyEntities := extractMoneyEntities(text)
	entities = append(entities, moneyEntities...)
	
	// Extract percentage entities
	percentEntities := extractPercentEntities(text)
	entities = append(entities, percentEntities...)
	
	// Extract phone numbers
	phoneEntities := extractPhoneEntities(text)
	entities = append(entities, phoneEntities...)
	
	// Extract time entities
	timeEntities := extractTimeEntities(text)
	entities = append(entities, timeEntities...)
	
	// Extract code snippets
	codeEntities := extractCodeEntities(text)
	entities = append(entities, codeEntities...)
	
	// Extract numbers and ordinals
	numberEntities := extractNumberEntities(text)
	entities = append(entities, numberEntities...)
	
	// Extract metrics
	metricEntities := extractMetricEntities(text)
	entities = append(entities, metricEntities...)
	
	// Extract actions and concepts
	actionEntities := extractActionEntities(text)
	entities = append(entities, actionEntities...)
	
	// Deduplicate and resolve conflicts
	entities = deduplicateAdvancedEntities(entities)
	
	// Extract relationships
	extractEntityRelationships(entities, text)
	
	return entities
}

func extractContext(text string, start, end int) string {
	contextWindow := 50
	contextStart := start - contextWindow
	contextEnd := end + contextWindow
	
	if contextStart < 0 {
		contextStart = 0
	}
	if contextEnd > len(text) {
		contextEnd = len(text)
	}
	
	return text[contextStart:contextEnd]
}

func extractMoneyEntities(text string) []AdvancedEntity {
	var entities []AdvancedEntity
	matches := moneyPattern.FindAllStringSubmatchIndex(text, -1)
	
	for _, match := range matches {
		value := text[match[0]:match[1]]
		amount := parseMoneyAmount(value)
		currency := parseCurrency(value)
		
		entity := AdvancedEntity{
			Entity: Entity{
				Text:     value,
				Type:     EntityMoney,
				Position: Position{Start: match[0], End: match[1]},
			},
			Confidence: 0.95,
			Context:    extractContext(text, match[0], match[1]),
			Attributes: map[string]string{
				"amount":   strconv.FormatFloat(amount, 'f', 2, 64),
				"currency": currency,
			},
		}
		entities = append(entities, entity)
	}
	
	return entities
}

func parseMoneyAmount(text string) float64 {
	// Remove currency symbols and words
	cleaned := strings.TrimSpace(text)
	cleaned = strings.ReplaceAll(cleaned, "$", "")
	cleaned = strings.ReplaceAll(cleaned, ",", "")
	cleaned = strings.ReplaceAll(cleaned, "dollars", "")
	cleaned = strings.ReplaceAll(cleaned, "dollar", "")
	cleaned = strings.ReplaceAll(cleaned, "cents", "")
	cleaned = strings.ReplaceAll(cleaned, "cent", "")
	cleaned = strings.TrimSpace(cleaned)
	
	amount, _ := strconv.ParseFloat(cleaned, 64)
	return amount
}

func parseCurrency(text string) string {
	lower := strings.ToLower(text)
	switch {
	case strings.Contains(text, "$") || strings.Contains(lower, "dollar"):
		return "USD"
	case strings.Contains(lower, "euro"):
		return "EUR"
	case strings.Contains(lower, "pound"):
		return "GBP"
	case strings.Contains(lower, "yen"):
		return "JPY"
	case strings.Contains(lower, "yuan"):
		return "CNY"
	default:
		return "USD"
	}
}

func extractPercentEntities(text string) []AdvancedEntity {
	var entities []AdvancedEntity
	matches := percentPattern.FindAllStringSubmatchIndex(text, -1)
	
	for _, match := range matches {
		value := text[match[0]:match[1]]
		percent := parsePercentage(value)
		
		entity := AdvancedEntity{
			Entity: Entity{
				Text:     value,
				Type:     EntityPercent,
				Position: Position{Start: match[0], End: match[1]},
			},
			Confidence: 0.9,
			Context:    extractContext(text, match[0], match[1]),
			Attributes: map[string]string{
				"value": strconv.FormatFloat(percent, 'f', 2, 64),
			},
		}
		entities = append(entities, entity)
	}
	
	return entities
}

func parsePercentage(text string) float64 {
	cleaned := strings.TrimSpace(text)
	cleaned = strings.ReplaceAll(cleaned, "%", "")
	cleaned = strings.ReplaceAll(cleaned, "percent", "")
	cleaned = strings.ReplaceAll(cleaned, "percentage", "")
	cleaned = strings.TrimSpace(cleaned)
	
	value, _ := strconv.ParseFloat(cleaned, 64)
	return value
}

func extractPhoneEntities(text string) []AdvancedEntity {
	var entities []AdvancedEntity
	matches := phonePattern.FindAllStringSubmatchIndex(text, -1)
	
	for _, match := range matches {
		value := text[match[0]:match[1]]
		
		entity := AdvancedEntity{
			Entity: Entity{
				Text:     value,
				Type:     EntityPhone,
				Position: Position{Start: match[0], End: match[1]},
			},
			Confidence: 0.85,
			Context:    extractContext(text, match[0], match[1]),
			Attributes: map[string]string{
				"normalized": normalizePhoneNumber(value),
			},
		}
		entities = append(entities, entity)
	}
	
	return entities
}

func normalizePhoneNumber(phone string) string {
	// Remove all non-digit characters
	digits := ""
	for _, r := range phone {
		if unicode.IsDigit(r) {
			digits += string(r)
		}
	}
	return digits
}

func extractTimeEntities(text string) []AdvancedEntity {
	var entities []AdvancedEntity
	matches := timePattern.FindAllStringSubmatchIndex(text, -1)
	
	for _, match := range matches {
		value := text[match[0]:match[1]]
		
		entity := AdvancedEntity{
			Entity: Entity{
				Text:     value,
				Type:     EntityTime,
				Position: Position{Start: match[0], End: match[1]},
			},
			Confidence: 0.9,
			Context:    extractContext(text, match[0], match[1]),
			Attributes: map[string]string{
				"normalized": normalizeTime(value),
			},
		}
		entities = append(entities, entity)
	}
	
	return entities
}

func normalizeTime(timeStr string) string {
	// Simple normalization to 24-hour format
	timeStr = strings.ToLower(strings.TrimSpace(timeStr))
	
	// Try to parse and normalize
	layouts := []string{
		"3:04pm", "3:04 pm", "3pm", "3 pm",
		"15:04", "15",
	}
	
	for _, layout := range layouts {
		if t, err := time.Parse(layout, timeStr); err == nil {
			return t.Format("15:04")
		}
	}
	
	return timeStr
}

func extractCodeEntities(text string) []AdvancedEntity {
	var entities []AdvancedEntity
	matches := codePattern.FindAllStringSubmatchIndex(text, -1)
	
	for _, match := range matches {
		value := text[match[0]:match[1]]
		
		entity := AdvancedEntity{
			Entity: Entity{
				Text:     value,
				Type:     EntityCode,
				Position: Position{Start: match[0], End: match[1]},
			},
			Confidence: 1.0,
			Context:    extractContext(text, match[0], match[1]),
			Attributes: map[string]string{
				"language": detectCodeLanguage(value),
			},
		}
		entities = append(entities, entity)
	}
	
	return entities
}

func detectCodeLanguage(code string) string {
	// Simple heuristic-based detection
	lower := strings.ToLower(code)
	
	switch {
	case strings.Contains(lower, "func ") || strings.Contains(lower, "package "):
		return "go"
	case strings.Contains(lower, "def ") || strings.Contains(lower, "import "):
		return "python"
	case strings.Contains(lower, "function") || strings.Contains(lower, "const "):
		return "javascript"
	case strings.Contains(lower, "public class") || strings.Contains(lower, "private "):
		return "java"
	default:
		return "unknown"
	}
}

func extractNumberEntities(text string) []AdvancedEntity {
	var entities []AdvancedEntity
	
	// Extract ordinals
	ordinalMatches := ordinalPattern.FindAllStringSubmatchIndex(text, -1)
	for _, match := range ordinalMatches {
		value := text[match[0]:match[1]]
		
		entity := AdvancedEntity{
			Entity: Entity{
				Text:     value,
				Type:     EntityOrdinal,
				Position: Position{Start: match[0], End: match[1]},
			},
			Confidence: 0.95,
			Context:    extractContext(text, match[0], match[1]),
			Attributes: map[string]string{
				"numeric": convertOrdinalToNumber(value),
			},
		}
		entities = append(entities, entity)
	}
	
	// Extract plain numbers (avoiding overlap with other entities)
	numberMatches := numberPattern.FindAllStringSubmatchIndex(text, -1)
	for _, match := range numberMatches {
		// Check if this number is part of another entity
		if !isPartOfOtherEntity(match[0], match[1], entities) {
			value := text[match[0]:match[1]]
			
			entity := AdvancedEntity{
				Entity: Entity{
					Text:     value,
					Type:     EntityCardinal,
					Position: Position{Start: match[0], End: match[1]},
				},
				Confidence: 0.9,
				Context:    extractContext(text, match[0], match[1]),
				Attributes: map[string]string{
					"value": value,
				},
			}
			entities = append(entities, entity)
		}
	}
	
	return entities
}

func convertOrdinalToNumber(ordinal string) string {
	ordinalMap := map[string]string{
		"first": "1", "second": "2", "third": "3", "fourth": "4", "fifth": "5",
		"sixth": "6", "seventh": "7", "eighth": "8", "ninth": "9", "tenth": "10",
	}
	
	lower := strings.ToLower(ordinal)
	if num, ok := ordinalMap[lower]; ok {
		return num
	}
	
	// Extract number from ordinals like "1st", "2nd", etc.
	for _, r := range ordinal {
		if unicode.IsDigit(r) {
			continue
		} else {
			break
		}
	}
	
	return ordinal
}

func extractMetricEntities(text string) []AdvancedEntity {
	var entities []AdvancedEntity
	matches := metricPattern.FindAllStringSubmatchIndex(text, -1)
	
	for _, match := range matches {
		value := text[match[0]:match[1]]
		number, unit := parseMetric(value)
		
		entity := AdvancedEntity{
			Entity: Entity{
				Text:     value,
				Type:     EntityMetric,
				Position: Position{Start: match[0], End: match[1]},
			},
			Confidence: 0.95,
			Context:    extractContext(text, match[0], match[1]),
			Attributes: map[string]string{
				"value": strconv.FormatFloat(number, 'f', -1, 64),
				"unit":  unit,
				"type":  classifyMetricType(unit),
			},
		}
		entities = append(entities, entity)
	}
	
	return entities
}

func parseMetric(text string) (float64, string) {
	// Split number and unit
	parts := strings.Fields(text)
	if len(parts) >= 2 {
		number, _ := strconv.ParseFloat(parts[0], 64)
		return number, parts[1]
	}
	
	// Try to extract inline (e.g., "10ms")
	for i, r := range text {
		if !unicode.IsDigit(r) && r != '.' {
			if i > 0 {
				number, _ := strconv.ParseFloat(text[:i], 64)
				return number, text[i:]
			}
			break
		}
	}
	
	return 0, text
}

func classifyMetricType(unit string) string {
	unit = strings.ToLower(unit)
	
	switch {
	case strings.Contains(unit, "second") || unit == "ms" || unit == "s":
		return "time"
	case strings.Contains(unit, "minute") || strings.Contains(unit, "hour") || 
		 strings.Contains(unit, "day") || strings.Contains(unit, "week") ||
		 strings.Contains(unit, "month") || strings.Contains(unit, "year"):
		return "duration"
	case strings.Contains(unit, "byte") || unit == "kb" || unit == "mb" || 
		 unit == "gb" || unit == "tb":
		return "data"
	case unit == "km" || unit == "m" || unit == "cm" || unit == "mm":
		return "distance"
	case unit == "kg" || unit == "g" || unit == "mg":
		return "weight"
	case strings.Contains(unit, "°"):
		return "temperature"
	default:
		return "other"
	}
}

func extractActionEntities(text string) []AdvancedEntity {
	var entities []AdvancedEntity
	
	// Action verbs that often indicate important actions in text
	actionVerbs := []string{
		"create", "delete", "update", "modify", "send", "receive", "process",
		"analyze", "generate", "transform", "execute", "run", "start", "stop",
		"enable", "disable", "configure", "install", "deploy", "build", "test",
		"validate", "verify", "authenticate", "authorize", "connect", "disconnect",
	}
	
	sentences := SplitIntoSentences(text)
	position := 0
	
	for _, sentence := range sentences {
		sentenceStart := strings.Index(text[position:], sentence)
		if sentenceStart == -1 {
			continue
		}
		sentenceStart += position
		
		words := strings.Fields(sentence)
		wordPos := sentenceStart
		
		for _, word := range words {
			wordStart := strings.Index(text[wordPos:], word)
			if wordStart == -1 {
				continue
			}
			wordStart += wordPos
			
			cleanWord := strings.ToLower(strings.Trim(word, ".,!?;:"))
			for _, action := range actionVerbs {
				if cleanWord == action || strings.HasPrefix(cleanWord, action) {
					entity := AdvancedEntity{
						Entity: Entity{
							Text:     word,
							Type:     EntityAction,
							Position: Position{Start: wordStart, End: wordStart + len(word)},
						},
						Confidence: 0.7,
						Context:    extractContext(text, wordStart, wordStart+len(word)),
						Attributes: map[string]string{
							"base_form": action,
							"tense":     detectVerbTense(word),
						},
					}
					entities = append(entities, entity)
					break
				}
			}
			
			wordPos = wordStart + len(word)
		}
		
		position = sentenceStart + len(sentence)
	}
	
	return entities
}

func detectVerbTense(verb string) string {
	lower := strings.ToLower(verb)
	
	switch {
	case strings.HasSuffix(lower, "ing"):
		return "progressive"
	case strings.HasSuffix(lower, "ed"):
		return "past"
	case strings.HasSuffix(lower, "s") && !strings.HasSuffix(lower, "ss"):
		return "present_third"
	default:
		return "base"
	}
}

func isPartOfOtherEntity(start, end int, entities []AdvancedEntity) bool {
	for _, entity := range entities {
		if start >= entity.Position.Start && end <= entity.Position.End {
			return true
		}
	}
	return false
}

func deduplicateAdvancedEntities(entities []AdvancedEntity) []AdvancedEntity {
	// Sort by position
	sort.Slice(entities, func(i, j int) bool {
		return entities[i].Position.Start < entities[j].Position.Start
	})
	
	// Remove overlapping entities, keeping the one with higher confidence
	var result []AdvancedEntity
	
	for i := 0; i < len(entities); i++ {
		keep := true
		
		for j := i + 1; j < len(entities); j++ {
			// Check for overlap
			if entities[i].Position.End > entities[j].Position.Start {
				// Keep the one with higher confidence
				if entities[j].Confidence > entities[i].Confidence {
					keep = false
					break
				}
			} else {
				// No more overlaps possible
				break
			}
		}
		
		if keep {
			result = append(result, entities[i])
		}
	}
	
	return result
}

func extractEntityRelationships(entities []AdvancedEntity, text string) {
	// Simple relationship extraction based on proximity and patterns
	for i := 0; i < len(entities); i++ {
		for j := i + 1; j < len(entities); j++ {
			// Check if entities are in the same sentence
			if areInSameSentence(entities[i], entities[j], text) {
				// Look for relationship patterns
				relation := detectRelationship(entities[i], entities[j], text)
				if relation.RelationType != "" {
					entities[i].Relationships = append(entities[i].Relationships, relation)
				}
			}
		}
	}
}

func areInSameSentence(e1, e2 AdvancedEntity, text string) bool {
	// Simple check: look for sentence ending between entities
	start := e1.Position.End
	end := e2.Position.Start
	
	if start > end {
		start, end = end, start
	}
	
	between := text[start:end]
	return !strings.ContainsAny(between, ".!?")
}

func detectRelationship(e1, e2 AdvancedEntity, text string) EntityRelation {
	// Extract text between entities
	start := e1.Position.End
	end := e2.Position.Start
	
	if start > end {
		start, end = end, start
		e1, e2 = e2, e1
	}
	
	between := strings.ToLower(text[start:end])
	
	// Pattern matching for relationships
	relation := EntityRelation{
		FromEntity: e1.Text,
		ToEntity:   e2.Text,
		Confidence: 0.5,
	}
	
	switch {
	case strings.Contains(between, " at ") && e1.Type == EntityPerson && e2.Type == EntityOrganization:
		relation.RelationType = "works_at"
		relation.Confidence = 0.8
	case strings.Contains(between, " in ") && e2.Type == EntityLocation:
		relation.RelationType = "located_in"
		relation.Confidence = 0.8
	case strings.Contains(between, " of ") && e1.Type == EntityPerson && e2.Type == EntityOrganization:
		relation.RelationType = "member_of"
		relation.Confidence = 0.7
	case strings.Contains(between, " from ") && e2.Type == EntityLocation:
		relation.RelationType = "from_location"
		relation.Confidence = 0.7
	case strings.Contains(between, " to ") && e2.Type == EntityLocation:
		relation.RelationType = "to_location"
		relation.Confidence = 0.7
	case strings.Contains(between, " owns ") || strings.Contains(between, " has "):
		relation.RelationType = "owns"
		relation.Confidence = 0.6
	}
	
	return relation
}