// Copyright 2025 Caia Tech
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package textlib

import (
	"time"
)

// Core types for RL-optimized API extensions

// ProcessingMetrics tracks performance characteristics of API calls
type ProcessingMetrics struct {
	TimeElapsed    time.Duration `json:"time_elapsed"`
	MemoryPeak     int64         `json:"memory_peak"`
	AlgorithmSteps int           `json:"algorithm_steps"`
	CacheHits      int           `json:"cache_hits"`
}

// QualityMetrics represents the quality assessment of results
type QualityMetrics struct {
	Accuracy   float64 `json:"accuracy"`
	Confidence float64 `json:"confidence"`
	Coverage   float64 `json:"coverage"`
}

// ResourceUsage tracks computational resources used
type ResourceUsage struct {
	MemoryUsedMB     int   `json:"memory_used_mb"`
	CPUTimeMs        int64 `json:"cpu_time_ms"`
	NetworkCallsMade int   `json:"network_calls_made"`
	CacheHits        int   `json:"cache_hits"`
}

// ProcessingCost represents the computational cost of an operation
type ProcessingCost struct {
	TimeMs    int64 `json:"time_ms"`
	MemoryKB  int64 `json:"memory_kb"`
	CPUCycles int64 `json:"cpu_cycles"`
}

// ComplexityReport represents detailed text complexity analysis
type ComplexityReport struct {
	LexicalComplexity   float64            `json:"lexical_complexity"`
	SyntacticComplexity float64            `json:"syntactic_complexity"`
	SemanticComplexity  float64            `json:"semantic_complexity"`
	ReadabilityScores   map[string]float64 `json:"readability_scores"`
	ProcessingTime      time.Duration      `json:"processing_time"`
	MemoryUsed          int64              `json:"memory_used"`
	AlgorithmUsed       string             `json:"algorithm_used"`
	QualityMetrics      QualityMetrics     `json:"quality_metrics"`
}

// KeyPhrase represents an extracted key phrase with metadata
type KeyPhrase struct {
	Text       string   `json:"text"`
	Score      float64  `json:"score"`
	Position   RLPosition `json:"position"`
	Category   string   `json:"category"`
	Context    string   `json:"context"`
	Confidence float64  `json:"confidence"`
}

// Position represents the location of text within a document
type RLPosition struct {
	Start int `json:"start"`
	End   int `json:"end"`
	Line  int `json:"line"`
}

// ReadabilityReport provides comprehensive readability analysis
type ReadabilityReport struct {
	Scores                 map[string]float64 `json:"scores"`
	Recommendation         string             `json:"recommendation"`
	TargetAudience         []string           `json:"target_audience"`
	ImprovementSuggestions []string           `json:"improvement_suggestions"`
	ProcessingCost         ProcessingCost     `json:"processing_cost"`
}

// LanguageResult represents language detection results
type LanguageResult struct {
	Language         string              `json:"language"`
	Confidence       float64             `json:"confidence"`
	Alternatives     []LanguageCandidate `json:"alternatives"`
	Method           string              `json:"method"`
	ProcessingTime   time.Duration       `json:"processing_time"`
}

// LanguageCandidate represents an alternative language detection
type LanguageCandidate struct {
	Language   string  `json:"language"`
	Confidence float64 `json:"confidence"`
	Reason     string  `json:"reason"`
}

// Summary represents text summarization results
type Summary struct {
	Text              string            `json:"text"`
	KeySentences      []string          `json:"key_sentences"`
	CompressionRatio  float64           `json:"compression_ratio"`
	QualityScore      float64           `json:"quality_score"`
	Method            string            `json:"method"`
	ProcessingMetrics ProcessingMetrics `json:"processing_metrics"`
}

// SentimentAnalysis represents comprehensive sentiment analysis
type SentimentAnalysis struct {
	OverallSentiment   Sentiment              `json:"overall_sentiment"`
	SentenceLevel      []SentenceSentiment    `json:"sentence_level"`
	AspectBased        map[string]Sentiment   `json:"aspect_based"`
	EmotionProfile     EmotionProfile         `json:"emotion_profile"`
	Confidence         float64                `json:"confidence"`
	ProcessingApproach string                 `json:"processing_approach"`
}

// Sentiment represents sentiment polarity and magnitude
type Sentiment struct {
	Polarity   float64 `json:"polarity"`   // -1 to 1
	Magnitude  float64 `json:"magnitude"`  // 0 to 1
	Label      string  `json:"label"`      // positive/negative/neutral
	Confidence float64 `json:"confidence"`
}

// SentenceSentiment represents sentiment for a single sentence
type SentenceSentiment struct {
	Text      string    `json:"text"`
	Sentiment Sentiment `json:"sentiment"`
	Position  RLPosition  `json:"position"`
}

// EmotionProfile represents emotional content analysis
type EmotionProfile struct {
	Joy      float64 `json:"joy"`
	Anger    float64 `json:"anger"`
	Fear     float64 `json:"fear"`
	Sadness  float64 `json:"sadness"`
	Surprise float64 `json:"surprise"`
	Trust    float64 `json:"trust"`
}

// Topic represents a discovered topic in text
type Topic struct {
	Name       string   `json:"name"`
	Keywords   []string `json:"keywords"`
	Confidence float64  `json:"confidence"`
	Coverage   float64  `json:"coverage"` // % of text related to this topic
	Examples   []string `json:"examples"`
}

// SummaryResult represents the output of text summarization
type SummaryResult struct {
	Summary             string          `json:"summary"`
	Method              string          `json:"method"`           // extractive/hybrid/abstractive
	OriginalSentences   int             `json:"original_sentences"`
	SummarySentences    int             `json:"summary_sentences"`
	CompressionRatio    float64         `json:"compression_ratio"`
	ProcessingTime      time.Duration   `json:"processing_time"`
	QualityMetrics      QualityMetrics  `json:"quality_metrics"`
}

// SentimentResult represents the output of sentiment analysis
type SentimentResult struct {
	OverallSentiment   Sentiment           `json:"overall_sentiment"`
	SentenceSentiments []SentenceSentiment `json:"sentence_sentiments"`
	EmotionProfile     EmotionProfile      `json:"emotion_profile"`
	Method             string              `json:"method"`       // lexicon-based/rule-based/contextual-analysis
	ProcessingTime     time.Duration       `json:"processing_time"`
	QualityMetrics     QualityMetrics      `json:"quality_metrics"`
}

// TopicResult represents the output of topic classification
type TopicResult struct {
	Topics         []Topic        `json:"topics"`
	Method         string         `json:"method"`         // clustering/statistical/comprehensive
	ProcessingTime time.Duration  `json:"processing_time"`
	QualityMetrics QualityMetrics `json:"quality_metrics"`
}

// DocumentAnalysis represents comprehensive document analysis
type DocumentAnalysis struct {
	TextAnalysis       ComplexityReport   `json:"text_analysis"`
	StructureAnalysis  StructureAnalysis  `json:"structure_analysis"`
	MetadataExtraction Metadata           `json:"metadata"`
	QualityAssessment  QualityAssessment  `json:"quality_assessment"`
	ProcessingStrategy string             `json:"processing_strategy"`
	Performance        PerformanceMetrics `json:"performance"`
}

// StructureAnalysis represents document structure analysis
type StructureAnalysis struct {
	DocumentType string            `json:"document_type"`
	Sections     []RLSection         `json:"sections"`
	Tables       []Table           `json:"tables"`
	Images       []ImageReference  `json:"images"`
	Links        []Link            `json:"links"`
}

// Section represents a document section
type RLSection struct {
	Title      string   `json:"title"`
	Level      int      `json:"level"`
	Content    string   `json:"content"`
	Position   RLPosition `json:"position"`
	WordCount  int      `json:"word_count"`
}

// Table represents a table in a document
type Table struct {
	Caption  string     `json:"caption"`
	Headers  []string   `json:"headers"`
	Rows     [][]string `json:"rows"`
	Position RLPosition   `json:"position"`
}

// ImageReference represents an image in a document
type ImageReference struct {
	Path        string   `json:"path"`
	Caption     string   `json:"caption"`
	AltText     string   `json:"alt_text"`
	Position    RLPosition `json:"position"`
	Dimensions  Dimensions `json:"dimensions"`
}

// Dimensions represents image dimensions
type Dimensions struct {
	Width  int `json:"width"`
	Height int `json:"height"`
}

// Link represents a hyperlink in a document
type Link struct {
	Text     string   `json:"text"`
	URL      string   `json:"url"`
	Position RLPosition `json:"position"`
	Type     string   `json:"type"` // internal/external
}

// Metadata represents extracted document metadata
type Metadata struct {
	Title       string    `json:"title"`
	Author      string    `json:"author"`
	Date        time.Time `json:"date"`
	Keywords    []string  `json:"keywords"`
	Language    string    `json:"language"`
	Description string    `json:"description"`
}

// QualityAssessment represents document quality analysis
type QualityAssessment struct {
	OverallScore       float64          `json:"overall_score"`
	ReadabilityScore   float64          `json:"readability_score"`
	CompletenessScore  float64          `json:"completeness_score"`
	ConsistencyScore   float64          `json:"consistency_score"`
	Issues             []QualityIssue   `json:"issues"`
	Recommendations    []string         `json:"recommendations"`
}

// QualityIssue represents a specific quality problem
type QualityIssue struct {
	Type        string   `json:"type"`
	Severity    string   `json:"severity"` // low/medium/high
	Description string   `json:"description"`
	Location    RLPosition `json:"location"`
	Suggestion  string   `json:"suggestion"`
}

// PerformanceMetrics tracks API performance
type PerformanceMetrics struct {
	TotalTime      time.Duration          `json:"total_time"`
	StepTimings    map[string]time.Duration `json:"step_timings"`
	MemoryUsage    int64                  `json:"memory_usage"`
	CacheUtilization float64              `json:"cache_utilization"`
}

// BatchResult represents results from batch processing
type BatchResult struct {
	Results            []interface{}     `json:"results"`
	OverallMetrics     OverallMetrics    `json:"overall_metrics"`
	ProcessingStrategy BatchStrategy     `json:"processing_strategy"`
	Errors             []ProcessingError `json:"errors"`
}

// OverallMetrics represents aggregate metrics for batch processing
type OverallMetrics struct {
	TotalProcessed   int           `json:"total_processed"`
	SuccessCount     int           `json:"success_count"`
	ErrorCount       int           `json:"error_count"`
	AverageTime      time.Duration `json:"average_time"`
	TotalTime        time.Duration `json:"total_time"`
	ThroughputPerSec float64       `json:"throughput_per_sec"`
}

// BatchStrategy represents batch processing configuration
type BatchStrategy struct {
	Parallel   bool   `json:"parallel"`
	BatchSize  int    `json:"batch_size"`
	Workers    int    `json:"workers"`
	Priority   string `json:"priority"` // "speed", "memory", "accuracy"
}

// ProcessingError represents an error during processing
type ProcessingError struct {
	Index       int    `json:"index"`
	Error       string `json:"error"`
	InputSample string `json:"input_sample"`
	Timestamp   time.Time `json:"timestamp"`
}

// Translation represents translation results
type Translation struct {
	OriginalText   string        `json:"original_text"`
	TranslatedText string        `json:"translated_text"`
	SourceLanguage string        `json:"source_language"`
	TargetLanguage string        `json:"target_language"`
	Confidence     float64       `json:"confidence"`
	Method         string        `json:"method"`
	CacheUsed      bool          `json:"cache_used"`
	ProcessingTime time.Duration `json:"processing_time"`
}

// ValidationResult represents validation results
type RLValidationResult struct {
	IsValid         bool              `json:"is_valid"`
	ValidationLevel string            `json:"validation_level"`
	Issues          []ValidationIssue `json:"issues"`
	Suggestions     []string          `json:"suggestions"`
	MXRecordValid   bool              `json:"mx_record_valid"`
	ProcessingCost  int               `json:"processing_cost"`
}

// ValidationIssue represents a specific validation problem
type ValidationIssue struct {
	Field       string `json:"field"`
	Issue       string `json:"issue"`
	Severity    string `json:"severity"`
	Suggestion  string `json:"suggestion"`
}

// DeepAnalysis represents comprehensive deep analysis
type DeepAnalysis struct {
	BasicAnalysis   interface{}     `json:"basic_analysis"`
	MLInsights      MLInsights      `json:"ml_insights"`
	AdvancedMetrics AdvancedMetrics `json:"advanced_metrics"`
	ResourceUsage   ResourceUsage   `json:"resource_usage"`
	QualityScore    float64         `json:"quality_score"`
}

// MLInsights represents machine learning-based insights
type MLInsights struct {
	TopicModeling      []Topic             `json:"topic_modeling"`
	SemanticSimilarity []SimilarityPair    `json:"semantic_similarity"`
	EntityRelations    []RLEntityRelation    `json:"entity_relations"`
	WritingStyle       WritingStyleProfile `json:"writing_style"`
}

// SimilarityPair represents semantic similarity between text segments
type SimilarityPair struct {
	Text1      string  `json:"text1"`
	Text2      string  `json:"text2"`
	Similarity float64 `json:"similarity"`
	Method     string  `json:"method"`
}

// EntityRelation represents relationships between entities
type RLEntityRelation struct {
	Entity1    string `json:"entity1"`
	Entity2    string `json:"entity2"`
	Relation   string `json:"relation"`
	Confidence float64 `json:"confidence"`
	Context    string `json:"context"`
}

// WritingStyleProfile represents writing style characteristics
type WritingStyleProfile struct {
	Formality       float64          `json:"formality"`
	Tone            string           `json:"tone"`
	VocabularyLevel string           `json:"vocabulary_level"`
	SentenceVariety float64          `json:"sentence_variety"`
	Characteristics []string         `json:"characteristics"`
}

// AdvancedMetrics represents advanced text metrics
type AdvancedMetrics struct {
	SemanticDensity   float64            `json:"semantic_density"`
	InformationValue  float64            `json:"information_value"`
	Coherence         float64            `json:"coherence"`
	Consistency       float64            `json:"consistency"`
	CustomMetrics     map[string]float64 `json:"custom_metrics"`
}

// OptimizationMetrics for RL integration
type OptimizationMetrics struct {
	QualityScore     float64       `json:"quality_score"`     // 0-1, higher is better
	PerformanceScore float64       `json:"performance_score"` // 0-1, higher is better (inverse of time)
	ResourceScore    float64       `json:"resource_score"`    // 0-1, higher is better (inverse of memory)
	UserSatisfaction float64       `json:"user_satisfaction"` // 0-1, based on user feedback
	WeightedTotal    float64       `json:"weighted_total"`    // Combined score
}

// ProcessingStrategy represents a strategy for text processing
type ProcessingStrategy struct {
	Name             string                 `json:"name"`
	Description      string                 `json:"description"`
	Parameters       map[string]interface{} `json:"parameters"`
	ExpectedQuality  float64                `json:"expected_quality"`
	ExpectedSpeed    float64                `json:"expected_speed"`
	ResourceRequirements ResourceRequirements `json:"resource_requirements"`
}

// ResourceRequirements specifies resource needs
type ResourceRequirements struct {
	MinMemoryMB      int     `json:"min_memory_mb"`
	MaxMemoryMB      int     `json:"max_memory_mb"`
	EstimatedCPUTime int64   `json:"estimated_cpu_time"`
	NetworkRequired  bool    `json:"network_required"`
	CacheRecommended bool    `json:"cache_recommended"`
}

// OptimizationResult represents the result of an optimization
type OptimizationResult struct {
	SelectedStrategy ProcessingStrategy  `json:"selected_strategy"`
	PredictedMetrics OptimizationMetrics `json:"predicted_metrics"`
	ActualMetrics    OptimizationMetrics `json:"actual_metrics"`
	Recommendation   string              `json:"recommendation"`
}