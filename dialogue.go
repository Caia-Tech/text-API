package textlib

import (
	"regexp"
	"strings"
)

// DialogueAnalysis contains analysis of conversational text
type DialogueAnalysis struct {
	// Speaker identification
	Speakers         []Speaker
	Utterances       []Utterance
	
	// Turn-taking analysis
	TurnCount        int
	AverageTurnLength float64
	TurnDistribution map[string]int // Speaker -> turn count
	
	// Interaction patterns
	InteractionFlow  []Interaction
	ResponseLatency  []int // Approximated by text distance
	
	// Dialogue acts
	DialogueActs     []DialogueAct
	
	// Conversation dynamics
	TopicShifts      []TopicShift
	Interruptions    []Interruption
	Overlaps         []Overlap
	
	// Question-Answer pairs
	QAPairs          []QAPair
	
	// Sentiment flow
	SentimentFlow    []SentimentPoint
	
	// Conversation type
	ConversationType string // "interview", "debate", "casual", "formal", etc.
	Formality        string // "formal", "informal", "mixed"
}

type Speaker struct {
	ID          string
	Name        string
	FirstAppearance int
	Utterances  int
	WordCount   int
	Characteristics SpeakerCharacteristics
}

type SpeakerCharacteristics struct {
	AverageWordLength    float64
	VocabularyDiversity  float64
	QuestionFrequency    float64
	FormalityLevel       string
	DominantSentiment    string
}

type Utterance struct {
	Speaker      string
	Text         string
	Position     Position
	WordCount    int
	Type         string // "statement", "question", "exclamation", "command"
	DialogueAct  string
}

type Interaction struct {
	FromSpeaker  string
	ToSpeaker    string
	Type         string // "question-answer", "statement-response", "agreement", "disagreement"
	Position     Position
}

type DialogueAct struct {
	Speaker      string
	Act          string // "greeting", "question", "answer", "agreement", "disagreement", "clarification", etc.
	Text         string
	Position     Position
}

type TopicShift struct {
	Position     Position
	FromTopic    string
	ToTopic      string
	Initiator    string
}

type Interruption struct {
	Interrupter  string
	Interrupted  string
	Position     Position
	Type         string // "supportive", "competitive", "clarification"
}

type Overlap struct {
	Speakers     []string
	Position     Position
	Duration     int // Approximated by character count
}

type QAPair struct {
	Question     Utterance
	Answer       Utterance
	ResponseTime int // Distance in characters
	Completeness string // "complete", "partial", "evasive"
}

type SentimentPoint struct {
	Speaker      string
	Position     Position
	Sentiment    string // "positive", "negative", "neutral"
	Intensity    float64
}

// AnalyzeDialogue performs comprehensive dialogue analysis
func AnalyzeDialogue(text string) *DialogueAnalysis {
	analysis := &DialogueAnalysis{
		Speakers:         []Speaker{},
		Utterances:       []Utterance{},
		TurnDistribution: make(map[string]int),
		InteractionFlow:  []Interaction{},
		DialogueActs:     []DialogueAct{},
		TopicShifts:      []TopicShift{},
		Interruptions:    []Interruption{},
		Overlaps:         []Overlap{},
		QAPairs:          []QAPair{},
		SentimentFlow:    []SentimentPoint{},
	}
	
	// Extract utterances and speakers
	utterances := extractUtterances(text)
	analysis.Utterances = utterances
	
	// Identify speakers
	speakers := identifySpeakers(utterances)
	analysis.Speakers = speakers
	
	// Analyze turn-taking
	analyzeTurnTaking(analysis, utterances)
	
	// Identify dialogue acts
	analysis.DialogueActs = classifyDialogueActs(utterances)
	
	// Analyze interactions
	analysis.InteractionFlow = analyzeInteractions(utterances)
	
	// Find Q&A pairs
	analysis.QAPairs = findQAPairs(utterances)
	
	// Detect interruptions and overlaps
	analysis.Interruptions = detectInterruptions(text, utterances)
	analysis.Overlaps = detectOverlaps(text, utterances)
	
	// Analyze topic shifts
	analysis.TopicShifts = detectTopicShifts(utterances)
	
	// Analyze sentiment flow
	analysis.SentimentFlow = analyzeSentimentFlow(utterances)
	
	// Classify conversation type
	analysis.ConversationType = classifyConversationType(analysis)
	analysis.Formality = assessFormality(utterances)
	
	return analysis
}

func extractUtterances(text string) []Utterance {
	utterances := []Utterance{}
	
	// Pattern for dialogue with speaker labels
	// Matches: "Speaker: dialogue" or "Speaker - dialogue" or "[Speaker] dialogue"
	patterns := []string{
		`(?m)^([A-Za-z][A-Za-z0-9\s]+):\s*(.+)$`,
		`(?m)^([A-Za-z][A-Za-z0-9\s]+)\s*-\s*(.+)$`,
		`(?m)^\[([A-Za-z][A-Za-z0-9\s]+)\]\s*(.+)$`,
	}
	
	combinedPattern := strings.Join(patterns, "|")
	re := regexp.MustCompile(combinedPattern)
	
	matches := re.FindAllStringSubmatchIndex(text, -1)
	
	if len(matches) == 0 {
		// Try to extract quoted dialogue
		utterances = extractQuotedDialogue(text)
	} else {
		// Extract labeled dialogue
		for _, match := range matches {
			var speaker, dialogue string
			
			// Determine which pattern matched
			if match[2] != -1 { // First pattern
				speaker = strings.TrimSpace(text[match[2]:match[3]])
				dialogue = strings.TrimSpace(text[match[4]:match[5]])
			} else if match[6] != -1 { // Second pattern
				speaker = strings.TrimSpace(text[match[6]:match[7]])
				dialogue = strings.TrimSpace(text[match[8]:match[9]])
			} else if match[10] != -1 { // Third pattern
				speaker = strings.TrimSpace(text[match[10]:match[11]])
				dialogue = strings.TrimSpace(text[match[12]:match[13]])
			}
			
			utterance := Utterance{
				Speaker:   speaker,
				Text:      dialogue,
				Position:  Position{Start: match[0], End: match[1]},
				WordCount: len(strings.Fields(dialogue)),
				Type:      classifyUtteranceType(dialogue),
			}
			
			utterances = append(utterances, utterance)
		}
	}
	
	return utterances
}

func extractQuotedDialogue(text string) []Utterance {
	utterances := []Utterance{}
	
	// Pattern for quoted dialogue with attribution
	// Matches: "dialogue," said Speaker or Speaker said, "dialogue"
	patterns := []string{
		`"([^"]+)"\s*,?\s*said\s+([A-Za-z][A-Za-z\s]+)`,
		`([A-Za-z][A-Za-z\s]+)\s+said\s*,?\s*"([^"]+)"`,
		`"([^"]+)"\s+([A-Za-z][A-Za-z\s]+)\s+(said|asked|replied|answered)`,
	}
	
	for _, pattern := range patterns {
		re := regexp.MustCompile(pattern)
		matches := re.FindAllStringSubmatchIndex(text, -1)
		
		for _, match := range matches {
			var speaker, dialogue string
			
			if strings.Contains(pattern, `said\s+([A-Za-z]`) {
				dialogue = text[match[2]:match[3]]
				speaker = text[match[4]:match[5]]
			} else {
				speaker = text[match[2]:match[3]]
				dialogue = text[match[4]:match[5]]
			}
			
			utterance := Utterance{
				Speaker:   strings.TrimSpace(speaker),
				Text:      strings.TrimSpace(dialogue),
				Position:  Position{Start: match[0], End: match[1]},
				WordCount: len(strings.Fields(dialogue)),
				Type:      classifyUtteranceType(dialogue),
			}
			
			utterances = append(utterances, utterance)
		}
	}
	
	// If no attributed quotes, extract all quotes as unknown speaker
	if len(utterances) == 0 {
		quoteRe := regexp.MustCompile(`"([^"]+)"`)
		matches := quoteRe.FindAllStringSubmatchIndex(text, -1)
		
		for i, match := range matches {
			dialogue := text[match[2]:match[3]]
			
			utterance := Utterance{
				Speaker:   "Speaker" + string(rune('A'+i%26)),
				Text:      dialogue,
				Position:  Position{Start: match[0], End: match[1]},
				WordCount: len(strings.Fields(dialogue)),
				Type:      classifyUtteranceType(dialogue),
			}
			
			utterances = append(utterances, utterance)
		}
	}
	
	return utterances
}

func classifyUtteranceType(text string) string {
	text = strings.TrimSpace(text)
	
	if strings.HasSuffix(text, "?") {
		return "question"
	} else if strings.HasSuffix(text, "!") {
		return "exclamation"
	} else if isCommand(text) {
		return "command"
	}
	
	return "statement"
}

func isCommand(text string) bool {
	// Simple heuristic for commands
	words := strings.Fields(text)
	if len(words) == 0 {
		return false
	}
	
	firstWord := strings.ToLower(words[0])
	commandVerbs := []string{
		"go", "come", "stop", "wait", "look", "listen",
		"tell", "give", "take", "bring", "leave", "stay",
		"please", "let", "help", "show", "make", "do",
	}
	
	for _, verb := range commandVerbs {
		if firstWord == verb {
			return true
		}
	}
	
	return false
}

func identifySpeakers(utterances []Utterance) []Speaker {
	speakerMap := make(map[string]*Speaker)
	
	for i, utterance := range utterances {
		if _, exists := speakerMap[utterance.Speaker]; !exists {
			speakerMap[utterance.Speaker] = &Speaker{
				ID:              utterance.Speaker,
				Name:            utterance.Speaker,
				FirstAppearance: i,
				Utterances:      0,
				WordCount:       0,
			}
		}
		
		speakerMap[utterance.Speaker].Utterances++
		speakerMap[utterance.Speaker].WordCount += utterance.WordCount
	}
	
	// Calculate speaker characteristics
	for _, speaker := range speakerMap {
		speaker.Characteristics = calculateSpeakerCharacteristics(speaker.Name, utterances)
	}
	
	// Convert to slice
	speakers := []Speaker{}
	for _, speaker := range speakerMap {
		speakers = append(speakers, *speaker)
	}
	
	return speakers
}

func calculateSpeakerCharacteristics(speakerName string, utterances []Utterance) SpeakerCharacteristics {
	chars := SpeakerCharacteristics{}
	
	totalWords := 0
	totalChars := 0
	questions := 0
	words := []string{}
	
	for _, u := range utterances {
		if u.Speaker == speakerName {
			uWords := strings.Fields(u.Text)
			words = append(words, uWords...)
			
			for _, word := range uWords {
				totalChars += len(word)
			}
			
			if u.Type == "question" {
				questions++
			}
			
			totalWords += len(uWords)
		}
	}
	
	if totalWords > 0 {
		chars.AverageWordLength = float64(totalChars) / float64(totalWords)
		
		// Vocabulary diversity
		uniqueWords := make(map[string]bool)
		for _, w := range words {
			uniqueWords[strings.ToLower(w)] = true
		}
		chars.VocabularyDiversity = float64(len(uniqueWords)) / float64(totalWords)
	}
	
	speakerUtterances := 0
	for _, u := range utterances {
		if u.Speaker == speakerName {
			speakerUtterances++
		}
	}
	
	if speakerUtterances > 0 {
		chars.QuestionFrequency = float64(questions) / float64(speakerUtterances)
	}
	
	// Assess formality (simplified)
	informalWords := 0
	for _, w := range words {
		if isInformalWord(strings.ToLower(w)) {
			informalWords++
		}
	}
	
	informalRatio := float64(informalWords) / float64(max(totalWords, 1))
	if informalRatio > 0.1 {
		chars.FormalityLevel = "informal"
	} else if informalRatio < 0.02 {
		chars.FormalityLevel = "formal"
	} else {
		chars.FormalityLevel = "neutral"
	}
	
	return chars
}

func isInformalWord(word string) bool {
	informal := []string{
		"yeah", "yep", "nah", "nope", "gonna", "wanna",
		"gotta", "kinda", "sorta", "ain't", "y'all",
		"ok", "okay", "um", "uh", "like", "totally",
	}
	
	for _, inf := range informal {
		if word == inf {
			return true
		}
	}
	
	return false
}

func analyzeTurnTaking(analysis *DialogueAnalysis, utterances []Utterance) {
	if len(utterances) == 0 {
		return
	}
	
	analysis.TurnCount = len(utterances)
	
	totalLength := 0
	for _, u := range utterances {
		totalLength += u.WordCount
		analysis.TurnDistribution[u.Speaker]++
	}
	
	analysis.AverageTurnLength = float64(totalLength) / float64(len(utterances))
	
	// Calculate response latency (approximated)
	for i := 1; i < len(utterances); i++ {
		if utterances[i].Speaker != utterances[i-1].Speaker {
			latency := utterances[i].Position.Start - utterances[i-1].Position.End
			analysis.ResponseLatency = append(analysis.ResponseLatency, latency)
		}
	}
}

func classifyDialogueActs(utterances []Utterance) []DialogueAct {
	acts := []DialogueAct{}
	
	for _, u := range utterances {
		act := DialogueAct{
			Speaker:  u.Speaker,
			Text:     u.Text,
			Position: u.Position,
		}
		
		// Classify the act
		lower := strings.ToLower(u.Text)
		
		switch {
		case isGreeting(lower):
			act.Act = "greeting"
		case isFarewell(lower):
			act.Act = "farewell"
		case u.Type == "question":
			act.Act = "question"
		case isAnswer(lower):
			act.Act = "answer"
		case isAgreement(lower):
			act.Act = "agreement"
		case isDisagreement(lower):
			act.Act = "disagreement"
		case isClarification(lower):
			act.Act = "clarification"
		case isAcknowledgment(lower):
			act.Act = "acknowledgment"
		default:
			act.Act = "statement"
		}
		
		acts = append(acts, act)
	}
	
	return acts
}

func isGreeting(text string) bool {
	greetings := []string{
		"hello", "hi", "hey", "good morning", "good afternoon",
		"good evening", "howdy", "greetings",
	}
	
	for _, g := range greetings {
		if strings.Contains(text, g) {
			return true
		}
	}
	
	return false
}

func isFarewell(text string) bool {
	farewells := []string{
		"goodbye", "bye", "see you", "farewell", "take care",
		"good night", "later", "catch you later",
	}
	
	for _, f := range farewells {
		if strings.Contains(text, f) {
			return true
		}
	}
	
	return false
}

func isAnswer(text string) bool {
	// Simple heuristics
	answerStarts := []string{
		"yes", "no", "maybe", "i think", "i believe",
		"in my opinion", "well", "actually", "it's",
	}
	
	for _, start := range answerStarts {
		if strings.HasPrefix(text, start) {
			return true
		}
	}
	
	return false
}

func isAgreement(text string) bool {
	agreements := []string{
		"agree", "right", "exactly", "correct", "yes",
		"absolutely", "definitely", "of course", "sure",
		"i agree", "you're right", "that's right",
	}
	
	for _, a := range agreements {
		if strings.Contains(text, a) {
			return true
		}
	}
	
	return false
}

func isDisagreement(text string) bool {
	disagreements := []string{
		"disagree", "wrong", "no", "not really", "actually",
		"but", "however", "on the contrary", "i don't think",
		"that's not", "i don't agree",
	}
	
	for _, d := range disagreements {
		if strings.Contains(text, d) {
			return true
		}
	}
	
	return false
}

func isClarification(text string) bool {
	clarifications := []string{
		"what do you mean", "could you explain", "i don't understand",
		"what", "sorry", "pardon", "can you clarify",
		"in other words", "you mean",
	}
	
	for _, c := range clarifications {
		if strings.Contains(text, c) {
			return true
		}
	}
	
	return false
}

func isAcknowledgment(text string) bool {
	acknowledgments := []string{
		"i see", "uh huh", "okay", "got it", "understood",
		"mm hmm", "right", "i understand",
	}
	
	for _, a := range acknowledgments {
		if strings.Contains(text, a) {
			return true
		}
	}
	
	return false
}

func analyzeInteractions(utterances []Utterance) []Interaction {
	interactions := []Interaction{}
	
	for i := 1; i < len(utterances); i++ {
		if utterances[i].Speaker != utterances[i-1].Speaker {
			interaction := Interaction{
				FromSpeaker: utterances[i-1].Speaker,
				ToSpeaker:   utterances[i].Speaker,
				Position: Position{
					Start: utterances[i-1].Position.Start,
					End:   utterances[i].Position.End,
				},
			}
			
			// Classify interaction type
			if utterances[i-1].Type == "question" && utterances[i].Type != "question" {
				interaction.Type = "question-answer"
			} else if isAgreement(strings.ToLower(utterances[i].Text)) {
				interaction.Type = "agreement"
			} else if isDisagreement(strings.ToLower(utterances[i].Text)) {
				interaction.Type = "disagreement"
			} else {
				interaction.Type = "statement-response"
			}
			
			interactions = append(interactions, interaction)
		}
	}
	
	return interactions
}

func findQAPairs(utterances []Utterance) []QAPair {
	pairs := []QAPair{}
	
	for i := 0; i < len(utterances)-1; i++ {
		if utterances[i].Type == "question" {
			// Look for answer in next few utterances
			for j := i + 1; j < len(utterances) && j < i+4; j++ {
				if utterances[j].Speaker != utterances[i].Speaker {
					// Potential answer
					pair := QAPair{
						Question:     utterances[i],
						Answer:       utterances[j],
						ResponseTime: utterances[j].Position.Start - utterances[i].Position.End,
					}
					
					// Assess completeness
					if utterances[j].Type == "question" {
						pair.Completeness = "evasive"
					} else if utterances[j].WordCount < 5 {
						pair.Completeness = "partial"
					} else {
						pair.Completeness = "complete"
					}
					
					pairs = append(pairs, pair)
					break
				}
			}
		}
	}
	
	return pairs
}

func detectInterruptions(text string, utterances []Utterance) []Interruption {
	interruptions := []Interruption{}
	
	// Look for patterns like "..." or "--" at end of utterances
	for i := 0; i < len(utterances)-1; i++ {
		u1 := utterances[i]
		u2 := utterances[i+1]
		
		if u1.Speaker != u2.Speaker {
			// Check if first utterance appears incomplete
			if strings.HasSuffix(u1.Text, "...") || strings.HasSuffix(u1.Text, "--") ||
			   strings.HasSuffix(u1.Text, "—") {
				interruption := Interruption{
					Interrupter: u2.Speaker,
					Interrupted: u1.Speaker,
					Position:    u2.Position,
				}
				
				// Classify interruption type
				if isAgreement(strings.ToLower(u2.Text)) {
					interruption.Type = "supportive"
				} else if isClarification(strings.ToLower(u2.Text)) {
					interruption.Type = "clarification"
				} else {
					interruption.Type = "competitive"
				}
				
				interruptions = append(interruptions, interruption)
			}
		}
	}
	
	return interruptions
}

func detectOverlaps(text string, utterances []Utterance) []Overlap {
	overlaps := []Overlap{}
	
	// Check for position overlaps
	for i := 0; i < len(utterances)-1; i++ {
		for j := i + 1; j < len(utterances); j++ {
			if utterances[i].Position.End > utterances[j].Position.Start &&
			   utterances[i].Position.Start < utterances[j].Position.End {
				// Overlap detected
				overlap := Overlap{
					Speakers: []string{utterances[i].Speaker, utterances[j].Speaker},
					Position: Position{
						Start: max(utterances[i].Position.Start, utterances[j].Position.Start),
						End:   min(utterances[i].Position.End, utterances[j].Position.End),
					},
				}
				
				overlap.Duration = overlap.Position.End - overlap.Position.Start
				overlaps = append(overlaps, overlap)
			}
		}
	}
	
	return overlaps
}

func detectTopicShifts(utterances []Utterance) []TopicShift {
	shifts := []TopicShift{}
	
	if len(utterances) < 2 {
		return shifts
	}
	
	// Extract topics from utterances
	prevTopics := extractTopics(utterances[0].Text)
	
	for i := 1; i < len(utterances); i++ {
		currTopics := extractTopics(utterances[i].Text)
		
		// Check for topic change
		overlap := calculateTopicOverlap(prevTopics, currTopics)
		
		if overlap < 0.3 && len(prevTopics) > 0 && len(currTopics) > 0 {
			shift := TopicShift{
				Position:  utterances[i].Position,
				FromTopic: strings.Join(prevTopics[:min(2, len(prevTopics))], ", "),
				ToTopic:   strings.Join(currTopics[:min(2, len(currTopics))], ", "),
				Initiator: utterances[i].Speaker,
			}
			
			shifts = append(shifts, shift)
		}
		
		if len(currTopics) > 0 {
			prevTopics = currTopics
		}
	}
	
	return shifts
}

func calculateTopicOverlap(topics1, topics2 []string) float64 {
	if len(topics1) == 0 || len(topics2) == 0 {
		return 0
	}
	
	set1 := make(map[string]bool)
	for _, t := range topics1 {
		set1[t] = true
	}
	
	overlap := 0
	for _, t := range topics2 {
		if set1[t] {
			overlap++
		}
	}
	
	return float64(overlap) / float64(max(len(topics1), len(topics2)))
}

func analyzeSentimentFlow(utterances []Utterance) []SentimentPoint {
	flow := []SentimentPoint{}
	
	for _, u := range utterances {
		sentiment := analyzeSentiment(u.Text)
		
		point := SentimentPoint{
			Speaker:   u.Speaker,
			Position:  u.Position,
			Sentiment: sentiment.sentiment,
			Intensity: sentiment.intensity,
		}
		
		flow = append(flow, point)
	}
	
	return flow
}

type sentimentResult struct {
	sentiment string
	intensity float64
}

func analyzeSentiment(text string) sentimentResult {
	// Simple sentiment analysis
	lower := strings.ToLower(text)
	
	positiveWords := []string{
		"good", "great", "excellent", "wonderful", "happy",
		"love", "like", "enjoy", "pleased", "glad",
		"fantastic", "amazing", "beautiful", "perfect",
	}
	
	negativeWords := []string{
		"bad", "terrible", "awful", "hate", "dislike",
		"angry", "sad", "disappointed", "frustrated",
		"horrible", "disgusting", "annoying", "worst",
	}
	
	positiveCount := 0
	negativeCount := 0
	
	words := strings.Fields(lower)
	for _, word := range words {
		word = strings.Trim(word, ".,!?;:")
		
		for _, pos := range positiveWords {
			if word == pos {
				positiveCount++
			}
		}
		
		for _, neg := range negativeWords {
			if word == neg {
				negativeCount++
			}
		}
	}
	
	totalSentiment := positiveCount - negativeCount
	
	result := sentimentResult{}
	
	if totalSentiment > 0 {
		result.sentiment = "positive"
		result.intensity = float64(totalSentiment) / float64(len(words))
	} else if totalSentiment < 0 {
		result.sentiment = "negative"
		result.intensity = float64(-totalSentiment) / float64(len(words))
	} else {
		result.sentiment = "neutral"
		result.intensity = 0
	}
	
	// Cap intensity at 1.0
	if result.intensity > 1.0 {
		result.intensity = 1.0
	}
	
	return result
}

func classifyConversationType(analysis *DialogueAnalysis) string {
	// Based on various factors
	
	questionRatio := 0.0
	totalUtterances := len(analysis.Utterances)
	
	if totalUtterances > 0 {
		questions := 0
		for _, u := range analysis.Utterances {
			if u.Type == "question" {
				questions++
			}
		}
		questionRatio = float64(questions) / float64(totalUtterances)
	}
	
	// Check speaker balance
	maxUtterances := 0
	for _, count := range analysis.TurnDistribution {
		if count > maxUtterances {
			maxUtterances = count
		}
	}
	
	speakerDominance := float64(maxUtterances) / float64(totalUtterances)
	
	// Classify based on patterns
	if questionRatio > 0.4 && speakerDominance > 0.6 {
		return "interview"
	} else if len(analysis.Speakers) == 2 && questionRatio > 0.3 {
		
		// Check for debate patterns
		disagreements := 0
		for _, act := range analysis.DialogueActs {
			if act.Act == "disagreement" {
				disagreements++
			}
		}
		
		if float64(disagreements)/float64(totalUtterances) > 0.1 {
			return "debate"
		}
		
		return "discussion"
	} else if analysis.AverageTurnLength < 10 {
		return "casual"
	} else if analysis.Formality == "formal" {
		return "formal"
	}
	
	return "conversation"
}

func assessFormality(utterances []Utterance) string {
	if len(utterances) == 0 {
		return "neutral"
	}
	
	formalCount := 0
	informalCount := 0
	
	for _, u := range utterances {
		words := strings.Fields(strings.ToLower(u.Text))
		
		for _, word := range words {
			if isInformalWord(word) {
				informalCount++
			} else if isFormalWord(word) {
				formalCount++
			}
		}
		
		// Check for contractions
		if strings.Contains(u.Text, "'") {
			informalCount++
		}
		
		// Check for formal address
		if strings.Contains(u.Text, "Mr.") || strings.Contains(u.Text, "Mrs.") ||
		   strings.Contains(u.Text, "Dr.") || strings.Contains(u.Text, "Professor") {
			formalCount++
		}
	}
	
	if formalCount > informalCount*2 {
		return "formal"
	} else if informalCount > formalCount*2 {
		return "informal"
	}
	
	return "mixed"
}

func isFormalWord(word string) bool {
	formal := []string{
		"therefore", "however", "furthermore", "nevertheless",
		"consequently", "accordingly", "whereas", "albeit",
		"pursuant", "regarding", "concerning",
	}
	
	for _, f := range formal {
		if word == f {
			return true
		}
	}
	
	return false
}