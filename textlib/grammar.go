package textlib

import (
	"regexp"
	"strings"
	"unicode"
)

// GrammarAnalysis contains comprehensive grammar analysis results
type GrammarAnalysis struct {
	// Core issues
	Issues            []GrammarIssue
	IssueCount        int
	IssuesByType      map[string]int
	
	// Structural analysis
	SentenceStructure []SentenceStructure
	ClauseAnalysis    []ClauseAnalysis
	
	// Agreement checks
	SubjectVerbAgreement []AgreementIssue
	PronounAntecedent   []AgreementIssue
	NumberAgreement     []AgreementIssue
	
	// Tense and consistency
	TenseConsistency  []TenseIssue
	VoiceConsistency  []VoiceIssue
	
	// Punctuation analysis
	PunctuationIssues []PunctuationIssue
	QuotationBalance  BalanceCheck
	ParenthesesBalance BalanceCheck
	BracketBalance    BalanceCheck
	
	// Style and clarity
	PassiveVoiceUse   []PassiveVoiceInstance
	WordinessIssues   []WordinessIssue
	ClarityIssues     []ClarityIssue
	
	// Sentence quality
	RunOnSentences    []RunOnSentence
	SentenceFragments []SentenceFragment
	FusedSentences    []FusedSentence
	
	// Advanced features
	ParallelStructure []ParallelismIssue
	DanglingModifiers []ModifierIssue
	SplitInfinitives  []SplitInfinitive
	
	// Overall metrics
	GrammarScore      float64
	ComplexityLevel   string
	ReadabilityImpact string
}

type GrammarIssue struct {
	Type        string   // "agreement", "tense", "punctuation", "structure", etc.
	Severity    string   // "error", "warning", "suggestion"
	Message     string
	Position    Position
	Context     string   // Surrounding text
	Suggestions []string
}

type SentenceStructure struct {
	Text            string
	Type            string // "simple", "compound", "complex", "compound-complex"
	Position        Position
	ClauseCount     int
	MainClauses     []string
	SubordinateClauses []string
	CoordinatingConjunctions []string
	SubordinatingConjunctions []string
	Complexity      float64
}

type ClauseAnalysis struct {
	Text        string
	Type        string // "main", "subordinate", "relative", "noun", "adverbial"
	Position    Position
	Subject     string
	Predicate   string
	IsComplete  bool
}

type AgreementIssue struct {
	Type        string // "subject-verb", "pronoun-antecedent", "number"
	Subject     string
	Verb        string
	Position    Position
	Expected    string
	Actual      string
	Explanation string
}

type TenseIssue struct {
	Type        string // "inconsistent", "incorrect", "shift"
	Position    Position
	Context     string
	CurrentTense string
	ExpectedTense string
	Explanation string
}

type VoiceIssue struct {
	Type        string // "inconsistent", "unnecessary_passive"
	Position    Position
	Context     string
	CurrentVoice string
	SuggestedVoice string
}

type PunctuationIssue struct {
	Type        string // "missing", "incorrect", "extra", "misplaced"
	Character   string
	Position    Position
	Context     string
	Suggestion  string
}

type BalanceCheck struct {
	IsBalanced    bool
	OpenCount     int
	CloseCount    int
	UnmatchedOpen []Position
	UnmatchedClose []Position
}

type PassiveVoiceInstance struct {
	Position    Position
	Text        string
	AuxiliaryVerb string
	PastParticiple string
	ActiveSuggestion string
	IsAppropriate bool
}

type WordinessIssue struct {
	Type        string // "redundant", "verbose", "filler"
	Position    Position
	Text        string
	Suggestion  string
	WordCount   int
	SuggestedWordCount int
}

type ClarityIssue struct {
	Type        string // "ambiguous", "vague", "unclear_reference"
	Position    Position
	Text        string
	Problem     string
	Suggestions []string
}

type RunOnSentence struct {
	Position    Position
	Text        string
	Length      int
	ClauseCount int
	SuggestedBreaks []int
}

type SentenceFragment struct {
	Position    Position
	Text        string
	MissingElement string // "subject", "verb", "complete_thought"
	Suggestion  string
}

type FusedSentence struct {
	Position    Position
	Text        string
	JoinLocation int
	SuggestedFix string
}

type ParallelismIssue struct {
	Position    Position
	Text        string
	Elements    []string
	Problem     string
	Suggestion  string
}

type ModifierIssue struct {
	Type        string // "dangling", "misplaced", "squinting"
	Position    Position
	Modifier    string
	Modified    string
	Problem     string
	Suggestion  string
}

type SplitInfinitive struct {
	Position    Position
	Infinitive  string
	Adverb      string
	Suggestion  string
	Severity    string // Modern grammar is more lenient
}

// AnalyzeGrammar performs comprehensive grammar analysis
func AnalyzeGrammar(text string) *GrammarAnalysis {
	analysis := &GrammarAnalysis{
		Issues:               []GrammarIssue{},
		IssuesByType:        make(map[string]int),
		SentenceStructure:   []SentenceStructure{},
		ClauseAnalysis:      []ClauseAnalysis{},
		SubjectVerbAgreement: []AgreementIssue{},
		PronounAntecedent:   []AgreementIssue{},
		NumberAgreement:     []AgreementIssue{},
		TenseConsistency:    []TenseIssue{},
		VoiceConsistency:    []VoiceIssue{},
		PunctuationIssues:   []PunctuationIssue{},
		PassiveVoiceUse:     []PassiveVoiceInstance{},
		WordinessIssues:     []WordinessIssue{},
		ClarityIssues:       []ClarityIssue{},
		RunOnSentences:      []RunOnSentence{},
		SentenceFragments:   []SentenceFragment{},
		FusedSentences:      []FusedSentence{},
		ParallelStructure:   []ParallelismIssue{},
		DanglingModifiers:   []ModifierIssue{},
		SplitInfinitives:    []SplitInfinitive{},
	}
	
	// Analyze sentence structure
	analysis.SentenceStructure = analyzeSentenceStructure(text)
	
	// Check punctuation balance
	analysis.QuotationBalance = checkQuotationBalance(text)
	analysis.ParenthesesBalance = checkParenthesesBalance(text)
	analysis.BracketBalance = checkBracketBalance(text)
	
	// Detect grammar issues
	analysis.RunOnSentences = detectRunOnSentences(text)
	analysis.SentenceFragments = detectSentenceFragments(text)
	analysis.FusedSentences = detectFusedSentences(text)
	
	// Agreement checks
	analysis.SubjectVerbAgreement = checkSubjectVerbAgreement(text)
	
	// Voice and style analysis
	analysis.PassiveVoiceUse = detectPassiveVoice(text)
	analysis.WordinessIssues = detectWordiness(text)
	
	// Advanced grammar checks
	analysis.ParallelStructure = checkParallelStructure(text)
	analysis.DanglingModifiers = detectDanglingModifiers(text)
	analysis.SplitInfinitives = detectSplitInfinitives(text)
	
	// Punctuation issues
	analysis.PunctuationIssues = checkPunctuationIssues(text)
	
	// Tense consistency
	analysis.TenseConsistency = checkTenseConsistency(text)
	
	// Compile all issues
	compileAllIssues(analysis)
	
	// Calculate metrics
	analysis.GrammarScore = calculateGrammarScore(analysis)
	analysis.ComplexityLevel = assessComplexityLevel(analysis)
	analysis.ReadabilityImpact = assessReadabilityImpact(analysis)
	
	return analysis
}

func analyzeSentenceStructure(text string) []SentenceStructure {
	structures := []SentenceStructure{}
	sentences := SplitIntoSentences(text)
	
	charOffset := 0
	for _, sent := range sentences {
		// Find position in original text
		start := strings.Index(text[charOffset:], sent)
		if start == -1 {
			continue
		}
		start += charOffset
		
		structure := SentenceStructure{
			Text:     sent,
			Position: Position{Start: start, End: start + len(sent)},
		}
		
		// Analyze clause structure
		structure.MainClauses, structure.SubordinateClauses = identifyClauses(sent)
		structure.ClauseCount = len(structure.MainClauses) + len(structure.SubordinateClauses)
		
		// Identify conjunctions
		structure.CoordinatingConjunctions = findCoordinatingConjunctions(sent)
		structure.SubordinatingConjunctions = findSubordinatingConjunctions(sent)
		
		// Classify sentence type
		structure.Type = classifySentenceType(structure)
		structure.Complexity = calculateSentenceComplexity(structure)
		
		structures = append(structures, structure)
		charOffset = start + len(sent)
	}
	
	return structures
}

func identifyClauses(sentence string) ([]string, []string) {
	mainClauses := []string{}
	subordinateClauses := []string{}
	
	// Simple heuristic: split by coordinating conjunctions for main clauses
	coordConjPattern := `\b(and|but|or|nor|for|yet|so)\b`
	coordRe := regexp.MustCompile(coordConjPattern)
	
	parts := coordRe.Split(sentence, -1)
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part != "" && hasSubjectAndPredicate(part) {
			mainClauses = append(mainClauses, part)
		}
	}
	
	// Find subordinate clauses (starting with subordinating conjunctions)
	subordConjPattern := `\b(because|since|although|though|while|when|where|if|unless|until|after|before|as|that|which|who|whom|whose)\s+[^,.;!?]*`
	subordRe := regexp.MustCompile(subordConjPattern)
	matches := subordRe.FindAllString(sentence, -1)
	subordinateClauses = append(subordinateClauses, matches...)
	
	// If no main clauses found, treat whole sentence as one
	if len(mainClauses) == 0 {
		mainClauses = append(mainClauses, sentence)
	}
	
	return mainClauses, subordinateClauses
}

func hasSubjectAndPredicate(clause string) bool {
	words := strings.Fields(clause)
	if len(words) < 2 {
		return false
	}
	
	// Simple check: look for a verb
	hasVerb := false
	for _, word := range words {
		if isVerb(word) {
			hasVerb = true
			break
		}
	}
	
	return hasVerb
}

func isVerb(word string) bool {
	// Simple verb detection (can be enhanced)
	word = strings.ToLower(strings.Trim(word, ".,!?;:"))
	
	// Common verb patterns
	verbPatterns := []string{
		`.*ed$`,   // past tense
		`.*ing$`,  // present participle
		`.*s$`,    // third person singular
	}
	
	for _, pattern := range verbPatterns {
		matched, _ := regexp.MatchString(pattern, word)
		if matched {
			return true
		}
	}
	
	// Common irregular verbs and simple verbs
	irregularVerbs := []string{
		"is", "are", "was", "were", "be", "been", "being",
		"have", "has", "had", "do", "does", "did", "will",
		"would", "could", "should", "can", "may", "might",
		"go", "went", "come", "came", "see", "saw", "get", "got",
		"make", "made", "take", "took", "give", "gave", "know", "knew",
		"sit", "sat", "run", "ran", "walk", "walked", "talk", "talked",
		"eat", "ate", "drink", "drank", "sleep", "slept", "work", "worked",
		"play", "played", "read", "write", "wrote", "think", "thought",
		"find", "found", "tell", "told", "say", "said", "feel", "felt",
		"look", "looked", "listen", "heard", "move", "moved", "stop", "stopped",
	}
	
	for _, verb := range irregularVerbs {
		if word == verb {
			return true
		}
	}
	
	return false
}

func findCoordinatingConjunctions(sentence string) []string {
	conjunctions := []string{}
	coordConj := []string{"and", "but", "or", "nor", "for", "yet", "so"}
	
	words := strings.Fields(strings.ToLower(sentence))
	for _, word := range words {
		word = strings.Trim(word, ".,!?;:")
		for _, conj := range coordConj {
			if word == conj {
				conjunctions = append(conjunctions, conj)
			}
		}
	}
	
	return conjunctions
}

func findSubordinatingConjunctions(sentence string) []string {
	conjunctions := []string{}
	subordConj := []string{
		"because", "since", "although", "though", "while", "when",
		"where", "if", "unless", "until", "after", "before", "as",
	}
	
	words := strings.Fields(strings.ToLower(sentence))
	for _, word := range words {
		word = strings.Trim(word, ".,!?;:")
		for _, conj := range subordConj {
			if word == conj {
				conjunctions = append(conjunctions, conj)
			}
		}
	}
	
	return conjunctions
}

func classifySentenceType(structure SentenceStructure) string {
	mainCount := len(structure.MainClauses)
	subCount := len(structure.SubordinateClauses)
	
	if mainCount == 1 && subCount == 0 {
		return "simple"
	} else if mainCount > 1 && subCount == 0 {
		return "compound"
	} else if mainCount == 1 && subCount > 0 {
		return "complex"
	} else if mainCount > 1 && subCount > 0 {
		return "compound-complex"
	}
	
	return "simple"
}

func calculateSentenceComplexity(structure SentenceStructure) float64 {
	baseComplexity := float64(structure.ClauseCount) * 0.3
	
	// Add complexity for subordinate clauses
	baseComplexity += float64(len(structure.SubordinateClauses)) * 0.4
	
	// Add complexity for multiple conjunctions
	baseComplexity += float64(len(structure.CoordinatingConjunctions)+len(structure.SubordinatingConjunctions)) * 0.2
	
	// Normalize to 0-1 scale
	if baseComplexity > 3.0 {
		baseComplexity = 3.0
	}
	
	return baseComplexity / 3.0
}

func checkQuotationBalance(text string) BalanceCheck {
	check := BalanceCheck{
		UnmatchedOpen:  []Position{},
		UnmatchedClose: []Position{},
	}
	
	doubleQuotes := 0
	singleQuotes := 0
	
	runes := []rune(text)
	for _, r := range runes {
		switch r {
		case '"':
			doubleQuotes++
		case '\'', '\u2018', '\u2019':
			singleQuotes++
		}
	}
	
	check.IsBalanced = doubleQuotes%2 == 0 && singleQuotes%2 == 0
	check.OpenCount = (doubleQuotes + singleQuotes) / 2
	check.CloseCount = check.OpenCount
	
	return check
}

func checkParenthesesBalance(text string) BalanceCheck {
	return checkBalanceGeneric(text, '(', ')')
}

func checkBracketBalance(text string) BalanceCheck {
	return checkBalanceGeneric(text, '[', ']')
}

func checkBalanceGeneric(text string, open, close rune) BalanceCheck {
	check := BalanceCheck{
		UnmatchedOpen:  []Position{},
		UnmatchedClose: []Position{},
	}
	
	openStack := []int{}
	runes := []rune(text)
	
	for i, r := range runes {
		if r == open {
			openStack = append(openStack, i)
			check.OpenCount++
		} else if r == close {
			check.CloseCount++
			if len(openStack) > 0 {
				openStack = openStack[:len(openStack)-1]
			} else {
				check.UnmatchedClose = append(check.UnmatchedClose, Position{Start: i, End: i + 1})
			}
		}
	}
	
	// Remaining open positions are unmatched
	for _, pos := range openStack {
		check.UnmatchedOpen = append(check.UnmatchedOpen, Position{Start: pos, End: pos + 1})
	}
	
	check.IsBalanced = len(check.UnmatchedOpen) == 0 && len(check.UnmatchedClose) == 0
	
	return check
}

func detectRunOnSentences(text string) []RunOnSentence {
	runOns := []RunOnSentence{}
	sentences := SplitIntoSentences(text)
	
	charOffset := 0
	for _, sent := range sentences {
		start := strings.Index(text[charOffset:], sent)
		if start == -1 {
			continue
		}
		start += charOffset
		
		wordCount := len(strings.Fields(sent))
		commaCount := strings.Count(sent, ",")
		clauseCount := len(findCoordinatingConjunctions(sent)) + len(findSubordinatingConjunctions(sent)) + 1
		
		// Heuristic: long sentence with many clauses or commas
		if wordCount > 25 && (clauseCount > 3 || commaCount > 5) {
			runOn := RunOnSentence{
				Position:    Position{Start: start, End: start + len(sent)},
				Text:        sent,
				Length:      wordCount,
				ClauseCount: clauseCount,
			}
			
			// Suggest break points at conjunctions and strong punctuation
			runOn.SuggestedBreaks = findSentenceBreakPoints(sent)
			runOns = append(runOns, runOn)
		}
		
		charOffset = start + len(sent)
	}
	
	return runOns
}

func findSentenceBreakPoints(sentence string) []int {
	breakPoints := []int{}
	
	// Look for coordinating conjunctions that could start new sentences
	words := strings.Fields(sentence)
	position := 0
	
	for i, word := range words {
		cleanWord := strings.ToLower(strings.Trim(word, ".,!?;:"))
		
		// Coordinating conjunctions that often start new sentences
		if cleanWord == "and" || cleanWord == "but" || cleanWord == "so" {
			if i > 3 { // Don't break very short initial clauses
				breakPoints = append(breakPoints, position)
			}
		}
		
		position += len(word) + 1
	}
	
	return breakPoints
}

func detectSentenceFragments(text string) []SentenceFragment {
	fragments := []SentenceFragment{}
	sentences := SplitIntoSentences(text)
	
	charOffset := 0
	for _, sent := range sentences {
		start := strings.Index(text[charOffset:], sent)
		if start == -1 {
			continue
		}
		start += charOffset
		
		if isSentenceFragment(sent) {
			fragment := SentenceFragment{
				Position: Position{Start: start, End: start + len(sent)},
				Text:     sent,
			}
			
			fragment.MissingElement = identifyMissingElement(sent)
			fragment.Suggestion = suggestFragmentFix(sent, fragment.MissingElement)
			
			fragments = append(fragments, fragment)
		}
		
		charOffset = start + len(sent)
	}
	
	return fragments
}

func isSentenceFragment(sentence string) bool {
	sentence = strings.TrimSpace(sentence)
	if len(sentence) < 3 {
		return true
	}
	
	words := strings.Fields(sentence)
	if len(words) < 2 {
		return true
	}
	
	// Check for subject and predicate
	hasSubject := false
	hasPredicate := false
	hasDeterminer := false
	
	for _, word := range words {
		if isPronoun(word) || isNoun(word) {
			hasSubject = true
		}
		if isDeterminer(word) {
			hasDeterminer = true
		}
		if isVerb(word) {
			hasPredicate = true
		}
	}
	
	// A determiner followed by a noun can form a subject
	if hasDeterminer && !hasSubject {
		// Look for determiner + noun pattern
		for i := 0; i < len(words)-1; i++ {
			if isDeterminer(words[i]) && isNoun(words[i+1]) {
				hasSubject = true
				break
			}
		}
	}
	
	// Check if it starts with subordinating conjunction without main clause
	firstWord := strings.ToLower(strings.Trim(words[0], ".,!?;:"))
	subordConj := []string{"because", "since", "although", "when", "while", "if", "unless"}
	
	startsWithSubord := false
	for _, conj := range subordConj {
		if firstWord == conj {
			startsWithSubord = true
			break
		}
	}
	
	// Fragment if: starts with subordinating conjunction but no main clause follows
	if startsWithSubord && !strings.Contains(sentence, ",") {
		return true
	}
	
	// Fragment if missing essential elements
	return !hasSubject || !hasPredicate
}

func identifyMissingElement(sentence string) string {
	words := strings.Fields(sentence)
	hasSubject := false
	hasPredicate := false
	
	for _, word := range words {
		if isPronoun(word) || isNoun(word) {
			hasSubject = true
		}
		if isVerb(word) {
			hasPredicate = true
		}
	}
	
	if !hasSubject && !hasPredicate {
		return "complete_thought"
	} else if !hasSubject {
		return "subject"
	} else if !hasPredicate {
		return "verb"
	}
	
	return "complete_thought"
}

func suggestFragmentFix(sentence, missingElement string) string {
	switch missingElement {
	case "subject":
		return "Add a subject (who or what performs the action)"
	case "verb":
		return "Add a main verb or predicate"
	case "complete_thought":
		return "Connect to a complete sentence or add missing elements"
	default:
		return "Revise to form a complete sentence"
	}
}

func detectFusedSentences(text string) []FusedSentence {
	fused := []FusedSentence{}
	sentences := SplitIntoSentences(text)
	
	charOffset := 0
	for _, sent := range sentences {
		start := strings.Index(text[charOffset:], sent)
		if start == -1 {
			continue
		}
		start += charOffset
		
		// Look for comma splices and fused sentences
		if strings.Contains(sent, ",") {
			parts := strings.Split(sent, ",")
			if len(parts) == 2 {
				part1 := strings.TrimSpace(parts[0])
				part2 := strings.TrimSpace(parts[1])
				
				// Check if both parts could be independent clauses
				if hasSubjectAndPredicate(part1) && hasSubjectAndPredicate(part2) &&
				   !startsWithConjunction(part2) {
					fusedSent := FusedSentence{
						Position:     Position{Start: start, End: start + len(sent)},
						Text:         sent,
						JoinLocation: len(part1),
						SuggestedFix: "Use a semicolon, add a conjunction, or split into two sentences",
					}
					fused = append(fused, fusedSent)
				}
			}
		}
		
		charOffset = start + len(sent)
	}
	
	return fused
}

func startsWithConjunction(text string) bool {
	words := strings.Fields(text)
	if len(words) == 0 {
		return false
	}
	
	firstWord := strings.ToLower(strings.Trim(words[0], ".,!?;:"))
	conjunctions := []string{"and", "but", "or", "so", "yet", "because", "although", "when", "while"}
	
	for _, conj := range conjunctions {
		if firstWord == conj {
			return true
		}
	}
	
	return false
}

func checkSubjectVerbAgreement(text string) []AgreementIssue {
	issues := []AgreementIssue{}
	sentences := SplitIntoSentences(text)
	
	// Simple pattern matching for common agreement errors
	patterns := []struct {
		pattern     string
		issue       string
		suggestion  string
	}{
		{`\b(he|she|it)\s+(are|were)\b`, "subject-verb", "Use 'is' or 'was' with singular subjects"},
		{`\b(they|we|you)\s+(is|was)\b`, "subject-verb", "Use 'are' or 'were' with plural subjects"},
		{`\b(cats|dogs|birds|people|children)\s+is\b`, "subject-verb", "Use 'are' with plural nouns"},
		{`\b(cat|dog|bird|person|child)\s+are\b`, "subject-verb", "Use 'is' with singular nouns"},
		{`\bthe\s+\w+s\s+is\b`, "subject-verb", "Plural nouns usually take 'are'"},
		{`\bthe\s+\w+\s+are\b`, "number", "Check if noun is singular or plural"},
	}
	
	charOffset := 0
	for _, sent := range sentences {
		start := strings.Index(text[charOffset:], sent)
		if start == -1 {
			continue
		}
		start += charOffset
		
		for _, p := range patterns {
			re := regexp.MustCompile(p.pattern)
			matches := re.FindAllStringIndex(sent, -1)
			
			for _, match := range matches {
				issue := AgreementIssue{
					Type:        p.issue,
					Position:    Position{Start: start + match[0], End: start + match[1]},
					Explanation: p.suggestion,
				}
				issues = append(issues, issue)
			}
		}
		
		charOffset = start + len(sent)
	}
	
	return issues
}

func detectPassiveVoice(text string) []PassiveVoiceInstance {
	instances := []PassiveVoiceInstance{}
	
	// Pattern for passive voice: be verb + past participle
	passivePattern := `\b(is|are|was|were|be|been|being)\s+\w*ed\b`
	re := regexp.MustCompile(passivePattern)
	
	matches := re.FindAllStringSubmatchIndex(text, -1)
	for _, match := range matches {
		instance := PassiveVoiceInstance{
			Position: Position{Start: match[0], End: match[1]},
			Text:     text[match[0]:match[1]],
		}
		
		// Extract components
		parts := strings.Fields(instance.Text)
		if len(parts) >= 2 {
			instance.AuxiliaryVerb = parts[0]
			instance.PastParticiple = parts[1]
		}
		
		instance.ActiveSuggestion = "Consider using active voice for clarity"
		instance.IsAppropriate = assessPassiveAppropriate(text, match[0], match[1])
		
		instances = append(instances, instance)
	}
	
	return instances
}

func assessPassiveAppropriate(text string, start, end int) bool {
	// Simple heuristic: passive voice might be appropriate in formal/scientific writing
	context := ""
	contextStart := max(0, start-50)
	contextEnd := min(len(text), end+50)
	context = text[contextStart:contextEnd]
	
	// Look for scientific/formal indicators
	formalIndicators := []string{"research", "study", "analysis", "method", "result", "data"}
	for _, indicator := range formalIndicators {
		if strings.Contains(strings.ToLower(context), indicator) {
			return true
		}
	}
	
	return false
}

func detectWordiness(text string) []WordinessIssue {
	issues := []WordinessIssue{}
	
	// Common wordy phrases and their concise alternatives
	wordyPhrases := map[string]string{
		"in order to":           "to",
		"due to the fact that":  "because",
		"at this point in time": "now",
		"for the reason that":   "because",
		"in the event that":     "if",
		"it is important to note that": "",
		"it should be noted that": "",
		"the fact that":         "that",
		"in spite of the fact that": "although",
		"in view of the fact that": "because",
	}
	
	for wordy, concise := range wordyPhrases {
		pattern := regexp.MustCompile(`(?i)\b` + regexp.QuoteMeta(wordy) + `\b`)
		matches := pattern.FindAllStringIndex(text, -1)
		
		for _, match := range matches {
			issue := WordinessIssue{
				Type:       "verbose",
				Position:   Position{Start: match[0], End: match[1]},
				Text:       text[match[0]:match[1]],
				Suggestion: concise,
				WordCount:  len(strings.Fields(text[match[0]:match[1]])),
				SuggestedWordCount: len(strings.Fields(concise)),
			}
			
			if concise == "" {
				issue.Suggestion = "Consider removing this phrase"
				issue.SuggestedWordCount = 0
			}
			
			issues = append(issues, issue)
		}
	}
	
	return issues
}

func checkParallelStructure(text string) []ParallelismIssue {
	issues := []ParallelismIssue{}
	
	// Look for series with "and" or "or"
	seriesPattern := `\w+,\s*\w+,?\s*(and|or)\s*\w+`
	re := regexp.MustCompile(seriesPattern)
	
	matches := re.FindAllStringSubmatchIndex(text, -1)
	for _, match := range matches {
		seriesText := text[match[0]:match[1]]
		
		// Simple check: look for inconsistent verb forms
		if hasParallelismIssue(seriesText) {
			issue := ParallelismIssue{
				Position:   Position{Start: match[0], End: match[1]},
				Text:       seriesText,
				Problem:    "Inconsistent parallel structure",
				Suggestion: "Use consistent grammatical forms",
			}
			issues = append(issues, issue)
		}
	}
	
	return issues
}

func hasParallelismIssue(series string) bool {
	// Simple heuristic: check for mixed verb forms (very basic)
	parts := regexp.MustCompile(`[,\s]+(and|or)\s+`).Split(series, -1)
	
	if len(parts) < 2 {
		return false
	}
	
	forms := make(map[string]int)
	for _, part := range parts {
		words := strings.Fields(strings.TrimSpace(part))
		if len(words) > 0 {
			lastWord := strings.ToLower(strings.Trim(words[len(words)-1], ".,!?;:"))
			if strings.HasSuffix(lastWord, "ing") {
				forms["gerund"]++
			} else if strings.HasSuffix(lastWord, "ed") {
				forms["past"]++
			} else if strings.HasPrefix(lastWord, "to ") {
				forms["infinitive"]++
			} else {
				forms["other"]++
			}
		}
	}
	
	// Issue if more than one form type is present
	return len(forms) > 1
}

func detectDanglingModifiers(text string) []ModifierIssue {
	issues := []ModifierIssue{}
	sentences := SplitIntoSentences(text)
	
	// Pattern for potential dangling modifiers (participial phrases at start)
	danglingPattern := `^(.*ing[^,]*),\s*(.+)$`
	re := regexp.MustCompile(danglingPattern)
	
	charOffset := 0
	for _, sent := range sentences {
		start := strings.Index(text[charOffset:], sent)
		if start == -1 {
			continue
		}
		start += charOffset
		
		matches := re.FindStringSubmatch(sent)
		if len(matches) == 3 {
			modifier := matches[1]
			mainClause := matches[2]
			
			// Check if the subject of main clause logically performs the action in modifier
			if hasDanglingModifier(modifier, mainClause) {
				issue := ModifierIssue{
					Type:       "dangling",
					Position:   Position{Start: start, End: start + len(modifier)},
					Modifier:   modifier,
					Modified:   mainClause,
					Problem:    "Modifier may not clearly relate to the intended subject",
					Suggestion: "Ensure the modifier clearly relates to the subject",
				}
				issues = append(issues, issue)
			}
		}
		
		charOffset = start + len(sent)
	}
	
	return issues
}

func hasDanglingModifier(modifier, mainClause string) bool {
	// Simple heuristic: if main clause starts with "the", "it", or "there", 
	// it might not be performing the action in the modifier
	mainWords := strings.Fields(mainClause)
	if len(mainWords) > 0 {
		firstWord := strings.ToLower(strings.Trim(mainWords[0], ".,!?;:"))
		problematicStarts := []string{"the", "it", "there", "this", "that"}
		
		for _, start := range problematicStarts {
			if firstWord == start {
				return true
			}
		}
	}
	
	return false
}

func detectSplitInfinitives(text string) []SplitInfinitive {
	infinitives := []SplitInfinitive{}
	
	// Pattern: to + adverb + verb
	splitPattern := `\bto\s+(\w+ly)\s+(\w+)\b`
	re := regexp.MustCompile(splitPattern)
	
	matches := re.FindAllStringSubmatchIndex(text, -1)
	for _, match := range matches {
		infinitive := SplitInfinitive{
			Position:   Position{Start: match[0], End: match[1]},
			Infinitive: text[match[0]:match[1]],
		}
		
		if len(match) >= 6 {
			infinitive.Adverb = text[match[2]:match[3]]
			infinitive.Suggestion = "Consider: 'to " + text[match[4]:match[5]] + " " + infinitive.Adverb + "'"
		}
		
		infinitive.Severity = "suggestion" // Modern grammar is more lenient
		
		infinitives = append(infinitives, infinitive)
	}
	
	return infinitives
}

func checkPunctuationIssues(text string) []PunctuationIssue {
	issues := []PunctuationIssue{}
	
	// Common punctuation issues
	patterns := []struct {
		pattern     string
		issueType   string
		suggestion  string
	}{
		{`\s+,`, "spacing", "Remove space before comma"},
		{`,[^\s]`, "spacing", "Add space after comma"},
		{`\s+\.`, "spacing", "Remove space before period"},
		{`\.[A-Z]`, "spacing", "Add space after period"},
		{`\?\?+`, "extra", "Use single question mark"},
		{`!!+`, "extra", "Use single exclamation mark"},
		{`\.\.\.\.+`, "extra", "Use three dots for ellipsis"},
	}
	
	for _, p := range patterns {
		re := regexp.MustCompile(p.pattern)
		matches := re.FindAllStringIndex(text, -1)
		
		for _, match := range matches {
			issue := PunctuationIssue{
				Type:       p.issueType,
				Position:   Position{Start: match[0], End: match[1]},
				Context:    text[match[0]:match[1]],
				Suggestion: p.suggestion,
			}
			issues = append(issues, issue)
		}
	}
	
	return issues
}

func checkTenseConsistency(text string) []TenseIssue {
	issues := []TenseIssue{}
	sentences := SplitIntoSentences(text)
	
	if len(sentences) < 2 {
		return issues
	}
	
	prevTense := ""
	charOffset := 0
	
	for i, sent := range sentences {
		start := strings.Index(text[charOffset:], sent)
		if start == -1 {
			continue
		}
		start += charOffset
		
		currentTense := identifyTense(sent)
		
		if i > 0 && prevTense != "" && currentTense != "" && 
		   prevTense != currentTense && !isValidTenseShift(prevTense, currentTense) {
			issue := TenseIssue{
				Type:          "inconsistent",
				Position:      Position{Start: start, End: start + len(sent)},
				Context:       sent,
				CurrentTense:  currentTense,
				ExpectedTense: prevTense,
				Explanation:   "Inconsistent tense usage",
			}
			issues = append(issues, issue)
		}
		
		if currentTense != "" {
			prevTense = currentTense
		}
		
		charOffset = start + len(sent)
	}
	
	return issues
}

func identifyTense(sentence string) string {
	// Simple tense identification
	words := strings.Fields(strings.ToLower(sentence))
	
	for _, word := range words {
		word = strings.Trim(word, ".,!?;:")
		
		// Past tense indicators
		if strings.HasSuffix(word, "ed") || word == "was" || word == "were" {
			return "past"
		}
		
		// Present tense indicators
		if word == "is" || word == "are" || word == "am" {
			return "present"
		}
		
		// Future tense indicators
		if word == "will" || word == "shall" {
			return "future"
		}
	}
	
	return ""
}

func isValidTenseShift(from, to string) bool {
	// Some tense shifts are acceptable (e.g., for dialogue, reported speech)
	validShifts := map[string][]string{
		"past":    {"present"}, // For dialogue or general statements
		"present": {"past"},    // For examples or anecdotes
	}
	
	if validTenses, exists := validShifts[from]; exists {
		for _, valid := range validTenses {
			if to == valid {
				return true
			}
		}
	}
	
	return false
}

func compileAllIssues(analysis *GrammarAnalysis) {
	// Convert all specific issues to general issues list
	
	// Add run-on sentences
	for _, runOn := range analysis.RunOnSentences {
		issue := GrammarIssue{
			Type:     "structure",
			Severity: "warning",
			Message:  "Run-on sentence detected",
			Position: runOn.Position,
			Context:  truncateContext(runOn.Text, 100),
			Suggestions: []string{"Break into shorter sentences", "Use appropriate punctuation"},
		}
		analysis.Issues = append(analysis.Issues, issue)
		analysis.IssuesByType["structure"]++
	}
	
	// Add fragments
	for _, fragment := range analysis.SentenceFragments {
		issue := GrammarIssue{
			Type:     "structure",
			Severity: "error",
			Message:  "Sentence fragment detected",
			Position: fragment.Position,
			Context:  fragment.Text,
			Suggestions: []string{fragment.Suggestion},
		}
		analysis.Issues = append(analysis.Issues, issue)
		analysis.IssuesByType["structure"]++
	}
	
	// Add agreement issues
	for _, agreement := range analysis.SubjectVerbAgreement {
		issue := GrammarIssue{
			Type:     "agreement",
			Severity: "error",
			Message:  "Subject-verb agreement error",
			Position: agreement.Position,
			Context:  agreement.Explanation,
			Suggestions: []string{agreement.Explanation},
		}
		analysis.Issues = append(analysis.Issues, issue)
		analysis.IssuesByType["agreement"]++
	}
	
	// Add punctuation issues
	for _, punct := range analysis.PunctuationIssues {
		issue := GrammarIssue{
			Type:     "punctuation",
			Severity: "warning",
			Message:  "Punctuation issue",
			Position: punct.Position,
			Context:  punct.Context,
			Suggestions: []string{punct.Suggestion},
		}
		analysis.Issues = append(analysis.Issues, issue)
		analysis.IssuesByType["punctuation"]++
	}
	
	// Count total issues
	analysis.IssueCount = len(analysis.Issues)
}

func truncateContext(text string, maxLen int) string {
	if len(text) <= maxLen {
		return text
	}
	return text[:maxLen] + "..."
}

func calculateGrammarScore(analysis *GrammarAnalysis) float64 {
	if analysis.IssueCount == 0 {
		return 1.0
	}
	
	// Weight different issue types
	errorWeight := 3.0
	warningWeight := 1.0
	
	totalPenalty := 0.0
	for _, issue := range analysis.Issues {
		if issue.Severity == "error" {
			totalPenalty += errorWeight
		} else {
			totalPenalty += warningWeight
		}
	}
	
	// Normalize based on text length (approximate)
	textLength := 0
	for _, structure := range analysis.SentenceStructure {
		textLength += len(structure.Text)
	}
	
	if textLength == 0 {
		return 1.0
	}
	
	// Score based on errors per 1000 characters
	errorsPerK := (totalPenalty * 1000.0) / float64(textLength)
	score := 1.0 - (errorsPerK / 10.0) // 10 errors per 1000 chars = 0 score
	
	if score < 0 {
		score = 0
	}
	
	return score
}

func assessComplexityLevel(analysis *GrammarAnalysis) string {
	avgComplexity := 0.0
	if len(analysis.SentenceStructure) > 0 {
		for _, structure := range analysis.SentenceStructure {
			avgComplexity += structure.Complexity
		}
		avgComplexity /= float64(len(analysis.SentenceStructure))
	}
	
	if avgComplexity < 0.3 {
		return "simple"
	} else if avgComplexity < 0.6 {
		return "moderate"
	} else {
		return "complex"
	}
}

func assessReadabilityImpact(analysis *GrammarAnalysis) string {
	criticalIssues := 0
	for _, issue := range analysis.Issues {
		if issue.Severity == "error" {
			criticalIssues++
		}
	}
	
	if criticalIssues == 0 {
		return "minimal"
	} else if criticalIssues < 3 {
		return "moderate"
	} else {
		return "significant"
	}
}

// Utility functions

func isNoun(word string) bool {
	// Simple heuristic for noun detection
	word = strings.ToLower(strings.Trim(word, ".,!?;:"))
	
	// Common noun patterns
	nounSuffixes := []string{"tion", "sion", "ness", "ment", "ity", "er", "or", "ist"}
	for _, suffix := range nounSuffixes {
		if strings.HasSuffix(word, suffix) {
			return true
		}
	}
	
	// Common simple nouns (for better detection)
	commonNouns := []string{
		"cat", "dog", "house", "car", "book", "table", "chair", "door", "window",
		"man", "woman", "child", "boy", "girl", "person", "people",
		"day", "night", "time", "year", "month", "week",
		"water", "food", "money", "work", "school", "home",
		"mat", "tree", "bird", "fish", "hand", "head", "eye", "foot",
	}
	
	for _, noun := range commonNouns {
		if word == noun {
			return true
		}
	}
	
	// Capitalize words are often nouns (proper nouns)
	original := strings.Trim(word, ".,!?;:")
	if len(original) > 0 && unicode.IsUpper(rune(original[0])) {
		return true
	}
	
	return false
}

func isDeterminer(word string) bool {
	word = strings.ToLower(strings.Trim(word, ".,!?;:"))
	determiners := []string{"the", "a", "an", "this", "that", "these", "those", "my", "your", "his", "her", "its", "our", "their"}
	
	for _, det := range determiners {
		if word == det {
			return true
		}
	}
	
	return false
}

func isPronoun(word string) bool {
	word = strings.ToLower(strings.Trim(word, ".,!?;:"))
	pronouns := []string{
		"i", "you", "he", "she", "it", "we", "they",
		"me", "him", "her", "us", "them",
		"my", "your", "his", "her", "its", "our", "their",
		"mine", "yours", "hers", "ours", "theirs",
		"this", "that", "these", "those",
		"who", "whom", "whose", "which", "what",
	}
	
	for _, pronoun := range pronouns {
		if word == pronoun {
			return true
		}
	}
	
	return false
}