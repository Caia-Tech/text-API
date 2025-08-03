# TextLib RL-Optimized API Examples

## Table of Contents
1. [Basic Usage Examples](#basic-usage-examples)
2. [Performance Optimization Scenarios](#performance-optimization-scenarios)
3. [Quality vs Speed Trade-offs](#quality-vs-speed-trade-offs)
4. [Batch Processing Examples](#batch-processing-examples)
5. [Real-World Use Cases](#real-world-use-cases)
6. [RL Integration Examples](#rl-integration-examples)
7. [Migration Examples](#migration-examples)

---

## Basic Usage Examples

### 1. Multi-Algorithm Entity Extraction

```go
package main

import (
    "fmt"
    "github.com/caiatech/textlib"
)

func main() {
    text := "Apple Inc. was founded by Steve Jobs in Cupertino, California in 1976."
    
    // Fast extraction for real-time applications
    fastStrategy := textlib.EntityExtractionStrategy{
        Algorithm:           "fast",
        ConfidenceThreshold: 0.6,
        MaxEntities:        20,
        EnableCaching:      true,
        EnableParallel:     false,
    }
    
    fastResult := textlib.ExtractEntitiesOptimized(text, fastStrategy)
    fmt.Printf("Fast extraction: %d entities in %dms\n", 
        len(fastResult.Entities), fastResult.ProcessingMetrics.ProcessingTimeMs)
    
    // Accurate extraction for high-quality analysis
    accurateStrategy := textlib.EntityExtractionStrategy{
        Algorithm:           "accurate",
        ConfidenceThreshold: 0.9,
        MaxEntities:        50,
        EnableCaching:      true,
        EnableParallel:     true,
    }
    
    accurateResult := textlib.ExtractEntitiesOptimized(text, accurateStrategy)
    fmt.Printf("Accurate extraction: %d entities in %dms (quality: %.2f)\n", 
        len(accurateResult.Entities), 
        accurateResult.ProcessingMetrics.ProcessingTimeMs,
        accurateResult.QualityAssessment.OverallScore)
    
    // Compare alternatives
    for _, alt := range accurateResult.Alternatives {
        fmt.Printf("Alternative %s: estimated %dms, quality %.2f\n",
            alt.Strategy.Algorithm, 
            alt.PredictedMetrics.ProcessingTimeMs,
            alt.PredictedMetrics.QualityScore)
    }
}
```

### 2. Tiered Text Analysis

```go
func demonstrateTieredAnalysis() {
    text := `Artificial Intelligence (AI) has revolutionized numerous industries, 
    from healthcare and finance to transportation and entertainment. Machine learning 
    algorithms enable computers to learn from data without explicit programming, 
    while deep learning networks mimic the human brain's neural structure.`
    
    strategy := textlib.ProcessingStrategy{
        Priority:           "balanced",
        MaxMemoryMB:        100,
        MaxTimeoutMs:       5000,
        QualityThreshold:   0.8,
        CachingEnabled:     true,
    }
    
    // Level 1: Basic analysis (fastest)
    basicResult := textlib.AnalyzeTextTiered(text, textlib.LevelBasic, strategy)
    fmt.Printf("Basic Analysis (%.0fms):\n", float64(basicResult.ProcessingMetrics.ProcessingTimeMs))
    fmt.Printf("  Words: %d, Sentences: %d\n", 
        basicResult.BasicStats.WordCount, 
        basicResult.BasicStats.SentenceCount)
    fmt.Printf("  Reading Level: %.1f\n", basicResult.BasicStats.FleschReadingEase)
    
    // Show preview of next level
    if basicResult.NextLevelPreview != nil {
        fmt.Printf("  Next level would take ~%dms and provide: %s\n",
            basicResult.NextLevelPreview.EstimatedTime.Milliseconds(),
            strings.Join(basicResult.NextLevelPreview.AdditionalInsights, ", "))
    }
    
    // Level 2: Standard analysis (includes entities)
    standardResult := textlib.AnalyzeTextTiered(text, textlib.LevelStandard, strategy)
    fmt.Printf("\nStandard Analysis (%.0fms):\n", 
        float64(standardResult.ProcessingMetrics.ProcessingTimeMs))
    fmt.Printf("  Entities found: %d\n", len(standardResult.Entities))
    for _, entity := range standardResult.Entities[:3] { // Show first 3
        fmt.Printf("    %s (%s): %.2f confidence\n", 
            entity.Text, entity.Type, entity.Confidence)
    }
    
    // Level 3: Advanced analysis (includes sentiment)
    advancedResult := textlib.AnalyzeTextTiered(text, textlib.LevelAdvanced, strategy)
    fmt.Printf("\nAdvanced Analysis (%.0fms):\n", 
        float64(advancedResult.ProcessingMetrics.ProcessingTimeMs))
    if advancedResult.Sentiment != nil {
        fmt.Printf("  Overall sentiment: %s (%.2f polarity)\n",
            advancedResult.Sentiment.OverallSentiment.Label,
            advancedResult.Sentiment.OverallSentiment.Polarity)
    }
}
```

### 3. Quality-Configurable Processing

```go
func demonstrateQualityLevels() {
    text := "The quick brown fox jumps over the lazy dog. This sentence contains most letters."
    
    qualityLevels := []textlib.QualityLevel{
        textlib.QualityFast,        // 0.1 - 10x speed, 70% accuracy
        textlib.QualityBalanced,    // 0.5 - 3x speed, 85% accuracy  
        textlib.QualityAccurate,    // 0.8 - 1.5x speed, 95% accuracy
        textlib.QualityPrecision,   // 1.0 - 1x speed, 99% accuracy
    }
    
    qualityNames := []string{"Fast", "Balanced", "Accurate", "Precision"}
    
    for i, quality := range qualityLevels {
        result := textlib.ProcessWithQuality(text, quality, "entity_extraction")
        
        fmt.Printf("%s Quality (%.1f):\n", qualityNames[i], float64(quality))
        fmt.Printf("  Processing time: %dms\n", result.ProcessingMetrics.ProcessingTimeMs)
        fmt.Printf("  Memory used: %dMB\n", result.ProcessingMetrics.MemoryUsageMB)
        fmt.Printf("  Quality achieved: %.2f\n", result.ActualQuality)
        fmt.Printf("  Trade-offs:\n")
        fmt.Printf("    Speed gain: %.1fx\n", result.QualityTradeoffs.SpeedGain)
        fmt.Printf("    Memory saving: %.1f%%\n", result.QualityTradeoffs.MemorySaving)
        fmt.Printf("    Accuracy loss: %.1f%%\n", result.QualityTradeoffs.AccuracyLoss)
        fmt.Println()
    }
}
```

---

## Performance Optimization Scenarios

### 1. Real-Time Chat Analysis

```go
// Scenario: Real-time chat moderation system
func chatModerationExample() {
    chatMessage := "Hey everyone! Check out this amazing new product: bit.ly/suspicious-link"
    
    // Ultra-fast analysis for real-time systems
    strategy := textlib.EntityExtractionStrategy{
        Algorithm:           "fast",
        ConfidenceThreshold: 0.5,  // Lower confidence for speed
        MaxEntities:        10,    // Limit entities for speed
        EnableCaching:      true,  // Cache common patterns
        EnableParallel:     false, // Avoid overhead for small text
    }
    
    start := time.Now()
    result := textlib.ExtractEntitiesOptimized(chatMessage, strategy)
    elapsed := time.Since(start)
    
    // Check for suspicious content
    hasLinks := false
    for _, entity := range result.Entities {
        if entity.Type == "URL" || entity.Type == "EMAIL" {
            hasLinks = true
            break
        }
    }
    
    fmt.Printf("Chat analysis completed in %dms\n", elapsed.Milliseconds())
    fmt.Printf("Suspicious links detected: %v\n", hasLinks)
    fmt.Printf("Cache hit rate: %.2f\n", result.ProcessingMetrics.CacheHitRate)
    
    // For production, aim for <50ms total processing time
    if elapsed.Milliseconds() > 50 {
        fmt.Printf("WARNING: Processing too slow for real-time chat (%dms)\n", elapsed.Milliseconds())
    }
}
```

### 2. Document Analysis Pipeline

```go
// Scenario: Large document processing with quality requirements
func documentAnalysisPipeline() {
    document := loadLargeDocument() // 50,000+ words
    
    // Multi-stage processing for optimal efficiency
    
    // Stage 1: Quick assessment
    quickStrategy := textlib.ProcessingStrategy{
        Priority:         "speed",
        MaxMemoryMB:      50,
        MaxTimeoutMs:     1000,
        QualityThreshold: 0.6,
        CachingEnabled:   true,
    }
    
    quickResult := textlib.AnalyzeTextTiered(document, textlib.LevelBasic, quickStrategy)
    
    // Decision point: determine if full analysis is needed
    complexity := calculateComplexity(quickResult.BasicStats)
    
    var finalResult textlib.TieredAnalysisResult
    
    if complexity > 0.7 {
        // High complexity document - use comprehensive analysis
        comprehensiveStrategy := textlib.ProcessingStrategy{
            Priority:         "accuracy",
            MaxMemoryMB:      500,
            MaxTimeoutMs:     30000,
            QualityThreshold: 0.9,
            CachingEnabled:   true,
        }
        
        finalResult = textlib.AnalyzeTextTiered(document, textlib.LevelComprehensive, comprehensiveStrategy)
        fmt.Printf("Comprehensive analysis: %dms, quality %.2f\n",
            finalResult.ProcessingMetrics.ProcessingTimeMs,
            finalResult.QualityScore)
    } else {
        // Simple document - standard analysis sufficient
        standardStrategy := textlib.ProcessingStrategy{
            Priority:         "balanced",
            MaxMemoryMB:      200,
            MaxTimeoutMs:     10000,
            QualityThreshold: 0.8,
            CachingEnabled:   true,
        }
        
        finalResult = textlib.AnalyzeTextTiered(document, textlib.LevelStandard, standardStrategy)
        fmt.Printf("Standard analysis: %dms, quality %.2f\n",
            finalResult.ProcessingMetrics.ProcessingTimeMs,
            finalResult.QualityScore)
    }
    
    // Generate insights based on analysis level
    generateDocumentInsights(finalResult)
}
```

### 3. Batch Processing Optimization

```go
// Scenario: Processing thousands of customer reviews
func batchReviewProcessing() {
    reviews := loadCustomerReviews() // 10,000 reviews
    
    // Optimize batch processing based on system resources
    systemMemoryMB := getAvailableMemoryMB()
    cpuCores := runtime.NumCPU()
    
    var strategy textlib.BatchStrategy
    
    if systemMemoryMB > 8000 && cpuCores >= 8 {
        // High-end system: aggressive parallel processing
        strategy = textlib.BatchStrategy{
            Parallel:        true,
            BatchSize:       100,
            WorkerCount:     cpuCores - 1,
            LoadBalancing:   "adaptive",
            MemorySharing:   true,
            CacheStrategy:   "aggressive",
            ErrorHandling:   "skip_errors",
        }
    } else if systemMemoryMB > 4000 {
        // Medium system: conservative parallel processing
        strategy = textlib.BatchStrategy{
            Parallel:        true,
            BatchSize:       50,
            WorkerCount:     cpuCores / 2,
            LoadBalancing:   "round_robin",
            MemorySharing:   false,
            CacheStrategy:   "selective",
            ErrorHandling:   "skip_errors",
        }
    } else {
        // Low-end system: sequential processing
        strategy = textlib.BatchStrategy{
            Parallel:        false,
            BatchSize:       20,
            WorkerCount:     1,
            LoadBalancing:   "round_robin",
            MemorySharing:   false,
            CacheStrategy:   "none",
            ErrorHandling:   "fail_fast",
        }
    }
    
    fmt.Printf("Processing %d reviews with strategy: %+v\n", len(reviews), strategy)
    
    start := time.Now()
    result := textlib.AnalyzeBatchOptimized(reviews, strategy, "sentiment_analysis")
    elapsed := time.Since(start)
    
    fmt.Printf("Batch processing completed in %v\n", elapsed)
    fmt.Printf("Success rate: %.1f%% (%d/%d)\n", 
        float64(result.OverallMetrics.SuccessfulItems)/float64(result.OverallMetrics.TotalItems)*100,
        result.OverallMetrics.SuccessfulItems,
        result.OverallMetrics.TotalItems)
    fmt.Printf("Average time per item: %.1fms\n", result.OverallMetrics.AverageTimePerItem)
    fmt.Printf("Peak memory usage: %dMB\n", result.OverallMetrics.PeakMemoryMB)
    fmt.Printf("CPU efficiency: %.1f%%\n", result.OverallMetrics.CPUEfficiency*100)
    
    // Apply optimization recommendations
    for _, tip := range result.OptimizationTips {
        fmt.Printf("💡 %s\n", tip)
    }
}
```

---

## Quality vs Speed Trade-offs

### 1. Email Validation Service

```go
// Scenario: Email validation with different accuracy requirements
func emailValidationService() {
    emails := []string{
        "user@example.com",
        "invalid.email",
        "test@nonexistent-domain.com",
        "user+tag@gmail.com",
    }
    
    // Fast validation for signup forms (real-time)
    fmt.Println("=== FAST VALIDATION (Real-time signup) ===")
    for _, email := range emails {
        result := textlib.ValidateEmail(email, false) // No MX check
        fmt.Printf("%s: %v (%dms)\n", 
            email, result.IsValid, result.ProcessingCost)
    }
    
    fmt.Println("\n=== THOROUGH VALIDATION (Batch processing) ===")
    for _, email := range emails {
        result := textlib.ValidateEmail(email, true) // With MX check
        fmt.Printf("%s: %v (%dms)", email, result.IsValid, result.ProcessingCost)
        if result.MXRecordValid {
            fmt.Printf(" ✓ MX valid")
        }
        if len(result.Issues) > 0 {
            fmt.Printf(" Issues: %v", result.Issues)
        }
        fmt.Println()
    }
}
```

### 2. Content Summarization

```go
// Scenario: News article summarization with different compression ratios
func contentSummarizationExample() {
    article := loadNewsArticle() // 2000 words
    
    scenarios := []struct {
        name       string
        ratio      float64
        algorithm  string
        useCase    string
    }{
        {"Tweet Summary", 0.05, "extractive", "Social media sharing"},
        {"Email Summary", 0.15, "abstractive", "Email newsletters"},
        {"Executive Summary", 0.25, "hybrid", "Business reports"},
        {"Detailed Summary", 0.40, "statistical", "Research purposes"},
    }
    
    for _, scenario := range scenarios {
        fmt.Printf("=== %s (%.0f%% compression) ===\n", scenario.name, scenario.ratio*100)
        
        start := time.Now()
        summary := textlib.SummarizeText(article, scenario.ratio, scenario.algorithm)
        elapsed := time.Since(start)
        
        fmt.Printf("Use case: %s\n", scenario.useCase)
        fmt.Printf("Original: %d words → Summary: %d words\n", 
            len(strings.Fields(article)), len(strings.Fields(summary.Text)))
        fmt.Printf("Processing time: %v\n", elapsed)
        fmt.Printf("Quality score: %.2f\n", summary.QualityScore)
        fmt.Printf("Compression achieved: %.1f%%\n", summary.CompressionRatio*100)
        fmt.Printf("Key sentences: %d\n", len(summary.KeySentences))
        
        // Show cache metrics
        fmt.Printf("Cache hits: %d\n", summary.ProcessingMetrics.CacheHits)
        
        fmt.Printf("Summary: %s...\n", summary.Text[:min(200, len(summary.Text))])
        fmt.Println()
    }
}
```

### 3. Language Detection Accuracy

```go
// Scenario: Multi-language support with accuracy requirements
func languageDetectionExample() {
    texts := map[string]string{
        "English":    "The quick brown fox jumps over the lazy dog.",
        "Spanish":    "El rápido zorro marrón salta sobre el perro perezoso.",
        "French":     "Le renard brun rapide saute par-dessus le chien paresseux.",
        "German":     "Der schnelle braune Fuchs springt über den faulen Hund.",
        "Mixed":      "Hello world! Como estas? Bonjour le monde!",
        "Short":      "OK",
    }
    
    confidenceLevels := []float64{0.5, 0.7, 0.9, 0.95}
    
    for language, text := range texts {
        fmt.Printf("=== %s Text ===\n", language)
        fmt.Printf("Text: %s\n", text)
        
        for _, confidence := range confidenceLevels {
            result := textlib.DetectLanguage(text, confidence)
            
            fmt.Printf("Confidence %.2f: %s (%.2f confidence, %dms)\n",
                confidence,
                result.Language,
                result.Confidence,
                result.ProcessingTime.Milliseconds())
                
            // Show alternatives for uncertain cases
            if result.Confidence < 0.9 && len(result.Alternatives) > 0 {
                fmt.Printf("  Alternatives: ")
                for _, alt := range result.Alternatives[:min(2, len(result.Alternatives))] {
                    fmt.Printf("%s(%.2f) ", alt.Language, alt.Confidence)
                }
                fmt.Println()
            }
        }
        fmt.Println()
    }
}
```

---

## Real-World Use Cases

### 1. Social Media Monitoring

```go
// Scenario: Real-time social media sentiment monitoring
func socialMediaMonitoring() {
    // Simulate social media posts
    posts := []string{
        "Love the new iPhone! Best camera ever! 📸",
        "Terrible customer service experience. Very disappointed 😞",
        "Just bought @Apple stock. Great quarterly results!",
        "Anyone else having issues with iOS update?",
    }
    
    // Real-time processing requirements: <100ms per post
    strategy := textlib.SentimentStrategy{
        Granularity: "document",  // Document-level for speed
        Model:       "lexicon",   // Fastest model
        EnableCache: true,
        MaxTimeMs:   100,
    }
    
    fmt.Println("=== Real-time Social Media Sentiment ===")
    
    totalPosts := 0
    totalTime := int64(0)
    positiveCount := 0
    negativeCount := 0
    
    for _, post := range posts {
        start := time.Now()
        sentiment := textlib.ExtractSentimentOptimized(post, strategy)
        elapsed := time.Since(start)
        
        totalPosts++
        totalTime += elapsed.Milliseconds()
        
        if sentiment.OverallSentiment.Polarity > 0.1 {
            positiveCount++
        } else if sentiment.OverallSentiment.Polarity < -0.1 {
            negativeCount++
        }
        
        fmt.Printf("Post: %s\n", post[:min(50, len(post))])
        fmt.Printf("Sentiment: %s (%.2f polarity, %dms)\n",
            sentiment.OverallSentiment.Label,
            sentiment.OverallSentiment.Polarity,
            elapsed.Milliseconds())
        
        // Alert on negative sentiment
        if sentiment.OverallSentiment.Polarity < -0.5 {
            fmt.Printf("🚨 ALERT: Strong negative sentiment detected!\n")
        }
        fmt.Println()
    }
    
    fmt.Printf("Summary: %d posts processed in %dms average\n",
        totalPosts, totalTime/int64(totalPosts))
    fmt.Printf("Sentiment distribution: %d positive, %d negative, %d neutral\n",
        positiveCount, negativeCount, totalPosts-positiveCount-negativeCount)
}
```

### 2. Customer Support Automation

```go
// Scenario: Intelligent customer support ticket routing
func customerSupportRouting() {
    tickets := []struct {
        id      string
        subject string
        content string
        urgency string
    }{
        {"T001", "Login Issue", "Cannot access my account after password reset", "high"},
        {"T002", "Billing Question", "Why was I charged twice this month?", "medium"},
        {"T003", "Feature Request", "Would love to see dark mode support", "low"},
        {"T004", "System Down", "Payment processing is completely broken!", "critical"},
    }
    
    for _, ticket := range tickets {
        fmt.Printf("=== Ticket %s ===\n", ticket.id)
        fmt.Printf("Subject: %s\n", ticket.subject)
        fmt.Printf("Urgency: %s\n", ticket.urgency)
        
        // Analyze ticket content
        fullText := ticket.subject + " " + ticket.content
        
        // Multi-tiered analysis based on urgency
        var analysisLevel textlib.AnalysisLevel
        switch ticket.urgency {
        case "critical":
            analysisLevel = textlib.LevelComprehensive // Full analysis
        case "high":
            analysisLevel = textlib.LevelAdvanced     // Include sentiment
        default:
            analysisLevel = textlib.LevelStandard     // Basic + entities
        }
        
        strategy := textlib.ProcessingStrategy{
            Priority:         "balanced",
            MaxMemoryMB:      100,
            MaxTimeoutMs:     5000,
            QualityThreshold: 0.8,
            CachingEnabled:   true,
        }
        
        result := textlib.AnalyzeTextTiered(fullText, analysisLevel, strategy)
        
        // Extract key information
        var department string
        var sentiment string
        var confidence float64
        
        // Determine department from entities
        for _, entity := range result.Entities {
            switch entity.Type {
            case "PAYMENT", "BILLING":
                department = "Finance"
            case "TECHNICAL", "ERROR":
                department = "Engineering"
            case "FEATURE", "REQUEST":
                department = "Product"
            default:
                department = "General Support"
            }
        }
        
        // Get sentiment if available
        if result.Sentiment != nil {
            sentiment = result.Sentiment.OverallSentiment.Label
            confidence = result.Sentiment.Confidence
        }
        
        fmt.Printf("Routing recommendation: %s\n", department)
        if sentiment != "" {
            fmt.Printf("Customer sentiment: %s (confidence: %.2f)\n", sentiment, confidence)
        }
        fmt.Printf("Processing time: %dms\n", result.ProcessingMetrics.ProcessingTimeMs)
        fmt.Printf("Analysis quality: %.2f\n", result.QualityScore)
        
        // Generate action recommendations
        if sentiment == "negative" && confidence > 0.8 {
            fmt.Printf("🔔 Recommend priority escalation due to negative sentiment\n")
        }
        
        fmt.Println()
    }
}
```

### 3. Content Quality Assessment

```go
// Scenario: Blog post quality assessment before publishing
func contentQualityAssessment() {
    blogPost := `
    Artificial Intelligence: The Future is Now
    
    AI is changing everything. Machine learning algorithms are everywhere. 
    They help us with recommendations on Netflix and Spotify. They power 
    search engines. They enable self-driving cars.
    
    Deep learning is a subset of machine learning. It uses neural networks 
    with many layers. These networks can learn complex patterns in data.
    
    The future of AI looks bright. We'll see more automation. More 
    personalization. More intelligent systems.
    `
    
    // Comprehensive analysis for content quality
    strategy := textlib.ProcessingStrategy{
        Priority:         "accuracy",
        MaxMemoryMB:      200,
        MaxTimeoutMs:     10000,
        QualityThreshold: 0.9,
        CachingEnabled:   true,
    }
    
    result := textlib.AnalyzeTextTiered(blogPost, textlib.LevelComprehensive, strategy)
    
    fmt.Println("=== Content Quality Assessment ===")
    
    // Basic readability
    fmt.Printf("Readability Score: %.1f\n", result.BasicStats.FleschReadingEase)
    readabilityLevel := categorizeReadability(result.BasicStats.FleschReadingEase)
    fmt.Printf("Reading Level: %s\n", readabilityLevel)
    
    // Content structure
    fmt.Printf("Word Count: %d\n", result.BasicStats.WordCount)
    fmt.Printf("Sentence Count: %d\n", result.BasicStats.SentenceCount)
    fmt.Printf("Average Sentence Length: %.1f words\n", 
        float64(result.BasicStats.WordCount)/float64(result.BasicStats.SentenceCount))
    
    // Entity analysis
    fmt.Printf("Key Topics/Entities: %d\n", len(result.Entities))
    topicMap := make(map[string]int)
    for _, entity := range result.Entities {
        topicMap[entity.Type]++
    }
    for topic, count := range topicMap {
        fmt.Printf("  %s: %d mentions\n", topic, count)
    }
    
    // Sentiment analysis
    if result.Sentiment != nil {
        fmt.Printf("Tone: %s (%.2f polarity)\n", 
            result.Sentiment.OverallSentiment.Label,
            result.Sentiment.OverallSentiment.Polarity)
    }
    
    // Quality recommendations
    fmt.Println("\n=== Improvement Recommendations ===")
    
    if result.BasicStats.FleschReadingEase < 60 {
        fmt.Println("• Consider simplifying language for better readability")
    }
    
    if result.BasicStats.WordCount < 300 {
        fmt.Println("• Content may be too short for good SEO performance")
    }
    
    avgSentenceLength := float64(result.BasicStats.WordCount) / float64(result.BasicStats.SentenceCount)
    if avgSentenceLength > 20 {
        fmt.Println("• Consider breaking up long sentences")
    }
    
    if len(result.Entities) < 3 {
        fmt.Println("• Add more specific topics and entities for better topical authority")
    }
    
    fmt.Printf("\nOverall Quality Score: %.1f/10\n", result.QualityScore*10)
}
```

---

## RL Integration Examples

### 1. Adaptive Performance Learning

```go
// Scenario: System learns optimal parameters based on usage patterns
func adaptivePerformanceLearning() {
    // Simulate RL training data collection
    trainingCollector := textlib.NewTrainingDataCollector()
    
    // Different text types for learning
    textTypes := map[string]string{
        "technical":    loadTechnicalDocument(),
        "creative":     loadCreativeWriting(),
        "news":         loadNewsArticle(),
        "social":       loadSocialMediaPosts(),
        "legal":        loadLegalDocument(),
    }
    
    fmt.Println("=== RL-Adaptive Performance Learning ===")
    
    for textType, text := range textTypes {
        fmt.Printf("\nLearning optimal parameters for %s text...\n", textType)
        
        // Try different strategies and collect performance data
        strategies := []textlib.EntityExtractionStrategy{
            {Algorithm: "fast", ConfidenceThreshold: 0.5, MaxEntities: 20},
            {Algorithm: "balanced", ConfidenceThreshold: 0.7, MaxEntities: 50},
            {Algorithm: "accurate", ConfidenceThreshold: 0.9, MaxEntities: 100},
        }
        
        bestStrategy := textlib.EntityExtractionStrategy{}
        bestReward := 0.0
        
        for _, strategy := range strategies {
            result := textlib.ExtractEntitiesOptimized(text, strategy)
            
            // Calculate multi-objective reward
            reward := calculateReward(result.ProcessingMetrics, result.QualityAssessment)
            
            fmt.Printf("  %s algorithm: reward %.3f (%.0fms, quality %.2f)\n",
                strategy.Algorithm, reward, 
                float64(result.ProcessingMetrics.ProcessingTimeMs),
                result.QualityAssessment.OverallScore)
            
            // Record training data
            trainingCollector.RecordFunctionCall("ExtractEntitiesOptimized", 
                map[string]interface{}{
                    "algorithm": strategy.Algorithm,
                    "confidence": strategy.ConfidenceThreshold,
                    "max_entities": strategy.MaxEntities,
                    "text_type": textType,
                }, result)
            
            if reward > bestReward {
                bestReward = reward
                bestStrategy = strategy
            }
        }
        
        fmt.Printf("  → Best strategy for %s: %s (reward: %.3f)\n",
            textType, bestStrategy.Algorithm, bestReward)
    }
    
    // Generate training dataset
    dataset := trainingCollector.GenerateTrainingDataset()
    fmt.Printf("\nTraining dataset generated: %d samples\n", len(dataset.Features))
    
    // Simulate learned optimal strategies
    learnedStrategies := map[string]textlib.EntityExtractionStrategy{
        "technical": {Algorithm: "accurate", ConfidenceThreshold: 0.9, MaxEntities: 50},
        "creative":  {Algorithm: "balanced", ConfidenceThreshold: 0.7, MaxEntities: 30},
        "news":      {Algorithm: "fast", ConfidenceThreshold: 0.6, MaxEntities: 25},
        "social":    {Algorithm: "fast", ConfidenceThreshold: 0.5, MaxEntities: 15},
        "legal":     {Algorithm: "accurate", ConfidenceThreshold: 0.95, MaxEntities: 75},
    }
    
    fmt.Println("\n=== Learned Optimal Strategies ===")
    for textType, strategy := range learnedStrategies {
        fmt.Printf("%s: %s algorithm, %.1f confidence, %d max entities\n",
            textType, strategy.Algorithm, strategy.ConfidenceThreshold, strategy.MaxEntities)
    }
}
```

### 2. Dynamic Resource Allocation

```go
// Scenario: RL optimizes resource allocation based on system state
func dynamicResourceAllocation() {
    // Simulate different system states
    systemStates := []struct {
        name           string
        cpuLoad        float64
        memoryUsage    float64
        networkLatency int64
        userPriority   string
    }{
        {"peak_hours", 0.85, 0.75, 200, "high"},
        {"normal_hours", 0.45, 0.40, 50, "medium"},
        {"off_hours", 0.20, 0.25, 25, "low"},
        {"maintenance", 0.95, 0.90, 500, "low"},
    }
    
    text := "Sample text for processing analysis and optimization testing."
    
    fmt.Println("=== Dynamic RL Resource Allocation ===")
    
    for _, state := range systemStates {
        fmt.Printf("\nSystem State: %s\n", state.name)
        fmt.Printf("CPU Load: %.0f%%, Memory: %.0f%%, Latency: %dms\n",
            state.cpuLoad*100, state.memoryUsage*100, state.networkLatency)
        
        // RL-learned strategy selection based on system state
        var strategy textlib.ProcessingStrategy
        
        if state.cpuLoad > 0.8 || state.memoryUsage > 0.8 {
            // High resource usage: prioritize efficiency
            strategy = textlib.ProcessingStrategy{
                Priority:         "speed",
                MaxMemoryMB:      30,
                MaxTimeoutMs:     1000,
                QualityThreshold: 0.6,
                CachingEnabled:   true,
            }
        } else if state.userPriority == "high" {
            // High priority user: prioritize quality
            strategy = textlib.ProcessingStrategy{
                Priority:         "accuracy",
                MaxMemoryMB:      200,
                MaxTimeoutMs:     10000,
                QualityThreshold: 0.9,
                CachingEnabled:   true,
            }
        } else {
            // Balanced approach
            strategy = textlib.ProcessingStrategy{
                Priority:         "balanced",
                MaxMemoryMB:      100,
                MaxTimeoutMs:     5000,
                QualityThreshold: 0.8,
                CachingEnabled:   true,
            }
        }
        
        // Apply RL-optimized strategy
        result := textlib.AnalyzeTextTiered(text, textlib.LevelStandard, strategy)
        
        fmt.Printf("Strategy: %s priority\n", strategy.Priority)
        fmt.Printf("Resource allocation: %dMB max memory, %dms timeout\n",
            strategy.MaxMemoryMB, strategy.MaxTimeoutMs)
        fmt.Printf("Results: %dms processing, %dMB used, quality %.2f\n",
            result.ProcessingMetrics.ProcessingTimeMs,
            result.ProcessingMetrics.MemoryUsageMB,
            result.QualityScore)
        
        // Calculate efficiency score
        resourceEfficiency := calculateResourceEfficiency(result.ProcessingMetrics, strategy)
        fmt.Printf("Resource efficiency: %.2f\n", resourceEfficiency)
        
        if resourceEfficiency < 0.7 {
            fmt.Printf("⚠️  Low efficiency detected - RL should adjust strategy\n")
        }
    }
}
```

---

## Migration Examples

### 1. Gradual Migration from Legacy API

```go
// Before: Legacy API usage
func legacyAPIUsage() {
    text := "Sample text for analysis"
    
    // Old way - no optimization control
    entities := textlib.ExtractAdvancedEntities(text)
    complexity := textlib.AnalyzeComplexity(text)
    
    fmt.Printf("Legacy API: Found %d entities\n", len(entities))
    fmt.Printf("Legacy complexity score: %.2f\n", complexity.FleschScore)
}

// After: Gradual migration to optimized API
func gradualMigration() {
    text := "Sample text for analysis"
    
    // Step 1: Use new API with default strategies (drop-in replacement)
    defaultStrategy := textlib.EntityExtractionStrategy{
        Algorithm:           "balanced",
        ConfidenceThreshold: 0.7,
        MaxEntities:        50,
        EnableCaching:      true,
    }
    
    entityResult := textlib.ExtractEntitiesOptimized(text, defaultStrategy)
    entities := entityResult.Entities
    
    // Get same result as legacy API but with metrics
    fmt.Printf("New API (default): Found %d entities in %dms\n", 
        len(entities), entityResult.ProcessingMetrics.ProcessingTimeMs)
    
    // Step 2: Gradually adopt optimization features
    if entityResult.ProcessingMetrics.ProcessingTimeMs > 500 {
        // Use faster algorithm if processing is slow
        fastStrategy := textlib.EntityExtractionStrategy{
            Algorithm:           "fast",
            ConfidenceThreshold: 0.6,
            MaxEntities:        30,
            EnableCaching:      true,
        }
        
        fastResult := textlib.ExtractEntitiesOptimized(text, fastStrategy)
        fmt.Printf("Optimized: Found %d entities in %dms (%.1fx faster)\n",
            len(fastResult.Entities),
            fastResult.ProcessingMetrics.ProcessingTimeMs,
            float64(entityResult.ProcessingMetrics.ProcessingTimeMs)/float64(fastResult.ProcessingMetrics.ProcessingTimeMs))
    }
    
    // Step 3: Use compatibility functions during transition
    legacyComplexity := textlib.AnalyzeComplexity(text) // Still works
    
    // But can also use new tiered analysis
    tieredResult := textlib.AnalyzeTextTiered(text, textlib.LevelBasic, textlib.ProcessingStrategy{
        Priority: "balanced",
    })
    
    fmt.Printf("Legacy complexity: %.2f\n", legacyComplexity.FleschScore)
    fmt.Printf("New tiered analysis: %.2f (quality: %.2f)\n", 
        tieredResult.BasicStats.FleschReadingEase, tieredResult.QualityScore)
}
```

### 2. Performance Optimization Migration

```go
// Migration example: Optimizing a high-traffic application
func performanceOptimizationMigration() {
    // Simulate high-traffic scenario
    texts := generateTestTexts(1000) // 1000 texts to process
    
    fmt.Println("=== Performance Optimization Migration ===")
    
    // Phase 1: Baseline with legacy API
    fmt.Println("\nPhase 1: Legacy API Baseline")
    start := time.Now()
    legacyResults := make([]interface{}, len(texts))
    
    for i, text := range texts {
        legacyResults[i] = textlib.ExtractAdvancedEntities(text)
    }
    
    legacyTime := time.Since(start)
    fmt.Printf("Legacy processing: %v total, %.1fms average per text\n",
        legacyTime, float64(legacyTime.Milliseconds())/float64(len(texts)))
    
    // Phase 2: Drop-in replacement with metrics
    fmt.Println("\nPhase 2: Drop-in Replacement (with metrics)")
    start = time.Now()
    
    defaultStrategy := textlib.EntityExtractionStrategy{
        Algorithm:           "balanced",
        ConfidenceThreshold: 0.7,
        MaxEntities:        50,
        EnableCaching:      true,
    }
    
    newResults := make([]textlib.EntityResult, len(texts))
    for i, text := range texts {
        newResults[i] = textlib.ExtractEntitiesOptimized(text, defaultStrategy)
    }
    
    newTime := time.Since(start)
    fmt.Printf("New API (default): %v total, %.1fms average per text\n",
        newTime, float64(newTime.Milliseconds())/float64(len(texts)))
    
    // Analyze performance characteristics
    avgQuality := 0.0
    for _, result := range newResults {
        avgQuality += result.QualityAssessment.OverallScore
    }
    avgQuality /= float64(len(newResults))
    
    fmt.Printf("Average quality: %.2f\n", avgQuality)
    
    // Phase 3: Optimize based on collected metrics
    fmt.Println("\nPhase 3: RL-Optimized Processing")
    
    // Analyze collected metrics to optimize strategy
    if avgQuality > 0.85 && newTime > legacyTime*110/100 {
        // Quality is high but speed decreased - try faster algorithm
        fmt.Println("Optimization: Quality sufficient, optimizing for speed")
        
        optimizedStrategy := textlib.EntityExtractionStrategy{
            Algorithm:           "fast",
            ConfidenceThreshold: 0.6,
            MaxEntities:        30,
            EnableCaching:      true,
            EnableParallel:     true,
        }
        
        start = time.Now()
        for i, text := range texts[:100] { // Test on subset first
            textlib.ExtractEntitiesOptimized(text, optimizedStrategy)
        }
        optimizedTime := time.Since(start)
        
        projectedFullTime := time.Duration(int64(optimizedTime) * int64(len(texts)) / 100)
        speedup := float64(newTime) / float64(projectedFullTime)
        
        fmt.Printf("Optimized strategy: %.1fx speedup projected\n", speedup)
        
        if speedup > 1.5 {
            fmt.Println("✅ Applying optimized strategy to production")
        }
    }
    
    // Phase 4: Batch processing optimization
    fmt.Println("\nPhase 4: Batch Processing Optimization")
    
    batchStrategy := textlib.BatchStrategy{
        Parallel:        true,
        BatchSize:       50,
        WorkerCount:     runtime.NumCPU() - 1,
        LoadBalancing:   "adaptive",
        CacheStrategy:   "aggressive",
        ErrorHandling:   "skip_errors",
    }
    
    start = time.Now()
    batchResult := textlib.AnalyzeBatchOptimized(texts[:200], batchStrategy, "entity_extraction")
    batchTime := time.Since(start)
    
    projectedBatchTime := time.Duration(int64(batchTime) * int64(len(texts)) / 200)
    batchSpeedup := float64(newTime) / float64(projectedBatchTime)
    
    fmt.Printf("Batch processing: %.1fx speedup, %.1f%% CPU efficiency\n",
        batchSpeedup, batchResult.OverallMetrics.CPUEfficiency*100)
    
    fmt.Println("\n=== Migration Summary ===")
    fmt.Printf("Legacy → New API: %.1fx speed change\n", 
        float64(legacyTime)/float64(newTime))
    fmt.Printf("New API → Optimized: %.1fx potential speedup\n", speedup)
    fmt.Printf("New API → Batch: %.1fx potential speedup\n", batchSpeedup)
    fmt.Printf("Total optimization potential: %.1fx\n", speedup*batchSpeedup)
}
```

---

## Utility Functions

```go
// Helper functions used in examples

func calculateReward(metrics textlib.OptimizationMetrics, quality textlib.QualityAssessment) float64 {
    // Multi-objective reward: balance speed, memory, and quality
    speedScore := 1.0 / (1.0 + float64(metrics.ProcessingTimeMs)/1000.0)
    memoryScore := 1.0 / (1.0 + float64(metrics.MemoryUsageMB)/100.0)
    qualityScore := quality.OverallScore
    
    return (speedScore*0.3 + memoryScore*0.2 + qualityScore*0.5)
}

func calculateResourceEfficiency(metrics textlib.OptimizationMetrics, strategy textlib.ProcessingStrategy) float64 {
    timeEfficiency := 1.0 - float64(metrics.ProcessingTimeMs)/float64(strategy.MaxTimeoutMs)
    memoryEfficiency := 1.0 - float64(metrics.MemoryUsageMB)/float64(strategy.MaxMemoryMB)
    return (timeEfficiency + memoryEfficiency) / 2.0
}

func categorizeReadability(score float64) string {
    switch {
    case score >= 90:
        return "Very Easy"
    case score >= 80:
        return "Easy"
    case score >= 70:
        return "Fairly Easy"
    case score >= 60:
        return "Standard"
    case score >= 50:
        return "Fairly Difficult"
    case score >= 30:
        return "Difficult"
    default:
        return "Very Difficult"
    }
}

func generateTestTexts(count int) []string {
    // Generate variety of test texts for benchmarking
    texts := make([]string, count)
    templates := []string{
        "Short text for testing.",
        "This is a medium-length text that contains several sentences and should provide adequate complexity for testing purposes.",
        "This is a much longer text that contains multiple paragraphs, complex sentence structures, and various types of entities including organizations, people, locations, and technical terms that should challenge the processing algorithms.",
    }
    
    for i := 0; i < count; i++ {
        texts[i] = templates[i%len(templates)]
    }
    return texts
}

func min(a, b int) int {
    if a < b {
        return a
    }
    return b
}
```

These comprehensive examples demonstrate the practical application of the RL-optimized TextLib API across various real-world scenarios, showcasing the flexibility and power of the optimization framework while maintaining ease of use and backward compatibility.