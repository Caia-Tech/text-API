# TextLib RL-Optimized API Implementation Guide

## Overview

This guide provides detailed instructions for implementing REST API endpoints for the TextLib RL-optimized functions. The API is designed to be deployed as a microservice with built-in performance tracking and adaptive algorithm selection.

## Table of Contents

1. [Architecture](#architecture)
2. [API Endpoints](#api-endpoints)
3. [Request/Response Examples](#requestresponse-examples)
4. [Error Handling](#error-handling)
5. [Performance Optimization](#performance-optimization)
6. [Deployment](#deployment)
7. [Monitoring](#monitoring)

## Architecture

### Component Structure

```
textlib-api/
├── cmd/
│   └── server/
│       └── main.go          # HTTP server entry point
├── internal/
│   ├── handlers/            # HTTP request handlers
│   │   ├── complexity.go
│   │   ├── keyphrases.go
│   │   ├── readability.go
│   │   └── language.go
│   ├── middleware/          # HTTP middleware
│   │   ├── auth.go
│   │   ├── logging.go
│   │   └── metrics.go
│   └── service/             # Business logic layer
│       └── textprocessor.go
├── pkg/
│   └── textlib/             # Core library (this package)
└── api/
    └── openapi.yaml         # OpenAPI specification
```

### HTTP Server Implementation

```go
// cmd/server/main.go
package main

import (
    "context"
    "log"
    "net/http"
    "os"
    "os/signal"
    "time"
    
    "github.com/gorilla/mux"
    "github.com/caiatech/textlib-api/internal/handlers"
    "github.com/caiatech/textlib-api/internal/middleware"
)

func main() {
    // Initialize router
    r := mux.NewRouter()
    
    // Apply middleware
    r.Use(middleware.Logging)
    r.Use(middleware.Metrics)
    r.Use(middleware.CORS)
    
    // API routes
    api := r.PathPrefix("/api/v1").Subrouter()
    
    // Complexity analysis
    api.HandleFunc("/analyze/complexity", handlers.AnalyzeComplexity).Methods("POST")
    
    // Key phrase extraction
    api.HandleFunc("/extract/keyphrases", handlers.ExtractKeyPhrases).Methods("POST")
    
    // Readability metrics
    api.HandleFunc("/calculate/readability", handlers.CalculateReadability).Methods("POST")
    
    // Language detection
    api.HandleFunc("/detect/language", handlers.DetectLanguage).Methods("POST")
    
    // Metrics endpoints
    api.HandleFunc("/metrics", handlers.GetMetrics).Methods("GET")
    api.HandleFunc("/metrics/reset", handlers.ResetMetrics).Methods("POST")
    
    // Health check
    api.HandleFunc("/health", handlers.HealthCheck).Methods("GET")
    
    // Server configuration
    srv := &http.Server{
        Addr:         ":8080",
        Handler:      r,
        ReadTimeout:  15 * time.Second,
        WriteTimeout: 15 * time.Second,
        IdleTimeout:  60 * time.Second,
    }
    
    // Start server
    go func() {
        log.Printf("Starting server on %s", srv.Addr)
        if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            log.Fatalf("Server failed to start: %v", err)
        }
    }()
    
    // Graceful shutdown
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, os.Interrupt)
    <-quit
    
    log.Println("Shutting down server...")
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()
    
    if err := srv.Shutdown(ctx); err != nil {
        log.Fatalf("Server forced to shutdown: %v", err)
    }
    
    log.Println("Server exited")
}
```

## API Endpoints

### 1. Analyze Text Complexity

```go
// internal/handlers/complexity.go
package handlers

import (
    "encoding/json"
    "net/http"
    "time"
    
    "github.com/caiatech/textlib"
)

type ComplexityRequest struct {
    Text  string `json:"text" validate:"required"`
    Depth int    `json:"depth,omitempty"`
}

type ComplexityResponse struct {
    LexicalComplexity   float64            `json:"lexicalComplexity"`
    SyntacticComplexity float64            `json:"syntacticComplexity"`
    SemanticComplexity  float64            `json:"semanticComplexity,omitempty"`
    ReadabilityScores   map[string]float64 `json:"readabilityScores"`
    ProcessingTime      string             `json:"processingTime"`
    MemoryUsed          int64              `json:"memoryUsed"`
    AlgorithmUsed       string             `json:"algorithmUsed"`
    QualityMetrics      QualityMetrics     `json:"qualityMetrics"`
}

func AnalyzeComplexity(w http.ResponseWriter, r *http.Request) {
    var req ComplexityRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        respondError(w, http.StatusBadRequest, "INVALID_JSON", err.Error())
        return
    }
    
    // Validate request
    if req.Text == "" {
        respondError(w, http.StatusBadRequest, "INVALID_INPUT", "Text parameter is required")
        return
    }
    
    // Default depth
    if req.Depth == 0 {
        req.Depth = 2
    }
    
    // Process text
    start := time.Now()
    report := textlib.AnalyzeTextComplexity(req.Text, req.Depth)
    
    // Build response
    resp := ComplexityResponse{
        LexicalComplexity:   report.LexicalComplexity,
        SyntacticComplexity: report.SyntacticComplexity,
        SemanticComplexity:  report.SemanticComplexity,
        ReadabilityScores:   report.ReadabilityScores,
        ProcessingTime:      report.ProcessingTime.String(),
        MemoryUsed:          report.MemoryUsed,
        AlgorithmUsed:       report.AlgorithmUsed,
        QualityMetrics: QualityMetrics{
            Accuracy:   report.QualityMetrics.Accuracy,
            Confidence: report.QualityMetrics.Confidence,
            Coverage:   report.QualityMetrics.Coverage,
        },
    }
    
    respondJSON(w, http.StatusOK, resp)
}
```

### 2. Extract Key Phrases

```go
// internal/handlers/keyphrases.go
package handlers

type KeyPhrasesRequest struct {
    Text       string `json:"text" validate:"required"`
    MaxPhrases int    `json:"maxPhrases" validate:"required,min=1,max=200"`
}

type KeyPhrasesResponse struct {
    Phrases  []KeyPhraseDTO `json:"phrases"`
    Metadata struct {
        Algorithm      string `json:"algorithm"`
        ProcessingTime string `json:"processingTime"`
    } `json:"metadata"`
}

type KeyPhraseDTO struct {
    Text       string     `json:"text"`
    Score      float64    `json:"score"`
    Position   PositionDTO `json:"position"`
    Category   string     `json:"category"`
    Context    string     `json:"context,omitempty"`
    Confidence float64    `json:"confidence"`
}

func ExtractKeyPhrases(w http.ResponseWriter, r *http.Request) {
    var req KeyPhrasesRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        respondError(w, http.StatusBadRequest, "INVALID_JSON", err.Error())
        return
    }
    
    // Validate
    if err := validate.Struct(req); err != nil {
        respondError(w, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
        return
    }
    
    // Process
    start := time.Now()
    phrases := textlib.ExtractKeyPhrases(req.Text, req.MaxPhrases)
    elapsed := time.Since(start)
    
    // Determine algorithm used
    algorithm := "tf-idf"
    if req.MaxPhrases > 10 && req.MaxPhrases <= 50 {
        algorithm = "statistical"
    } else if req.MaxPhrases > 50 {
        algorithm = "deep-nlp"
    }
    
    // Convert to DTOs
    phraseDTOs := make([]KeyPhraseDTO, len(phrases))
    for i, p := range phrases {
        phraseDTOs[i] = KeyPhraseDTO{
            Text:       p.Text,
            Score:      p.Score,
            Position:   PositionDTO{Start: p.Position.Start, End: p.Position.End},
            Category:   p.Category,
            Context:    p.Context,
            Confidence: p.Confidence,
        }
    }
    
    resp := KeyPhrasesResponse{
        Phrases: phraseDTOs,
    }
    resp.Metadata.Algorithm = algorithm
    resp.Metadata.ProcessingTime = elapsed.String()
    
    respondJSON(w, http.StatusOK, resp)
}
```

## Request/Response Examples

### Complexity Analysis

**Request:**
```bash
curl -X POST http://localhost:8080/api/v1/analyze/complexity \
  -H "Content-Type: application/json" \
  -d '{
    "text": "The epistemological ramifications of quantum entanglement necessitate a fundamental reconsideration of our ontological presuppositions.",
    "depth": 3
  }'
```

**Response:**
```json
{
  "lexicalComplexity": 0.89,
  "syntacticComplexity": 0.76,
  "semanticComplexity": 0.82,
  "readabilityScores": {
    "flesch-kincaid": 19.2,
    "gunning-fog": 22.8,
    "coleman-liau": 18.5,
    "ari": 20.1,
    "smog": 17.3
  },
  "processingTime": "287ms",
  "memoryUsed": 3145728,
  "algorithmUsed": "complexity-depth-3",
  "qualityMetrics": {
    "accuracy": 0.95,
    "confidence": 0.95,
    "coverage": 1.0
  }
}
```

### Key Phrase Extraction

**Request:**
```bash
curl -X POST http://localhost:8080/api/v1/extract/keyphrases \
  -H "Content-Type: application/json" \
  -d '{
    "text": "Machine learning algorithms are transforming data science. Deep learning neural networks process complex patterns in big data.",
    "maxPhrases": 5
  }'
```

**Response:**
```json
{
  "phrases": [
    {
      "text": "machine learning",
      "score": 0.92,
      "position": {"start": 0, "end": 16},
      "category": "phrase",
      "context": "Machine learning algorithms are...",
      "confidence": 0.88
    },
    {
      "text": "deep learning",
      "score": 0.87,
      "position": {"start": 58, "end": 71},
      "category": "phrase",
      "context": "...science. Deep learning neural...",
      "confidence": 0.85
    },
    {
      "text": "data science",
      "score": 0.79,
      "position": {"start": 45, "end": 57},
      "category": "phrase",
      "confidence": 0.82
    }
  ],
  "metadata": {
    "algorithm": "tf-idf",
    "processingTime": "52ms"
  }
}
```

## Error Handling

### Standard Error Response

```go
type ErrorResponse struct {
    Error   string                 `json:"error"`
    Message string                 `json:"message"`
    Details map[string]interface{} `json:"details,omitempty"`
}

func respondError(w http.ResponseWriter, code int, error, message string) {
    resp := ErrorResponse{
        Error:   error,
        Message: message,
    }
    respondJSON(w, code, resp)
}
```

### Common Error Codes

| HTTP Status | Error Code | Description |
|------------|------------|-------------|
| 400 | INVALID_INPUT | Missing or invalid parameters |
| 400 | INVALID_JSON | Malformed JSON in request body |
| 400 | VALIDATION_ERROR | Request validation failed |
| 401 | UNAUTHORIZED | Missing or invalid API key |
| 429 | RATE_LIMIT_EXCEEDED | Too many requests |
| 500 | INTERNAL_ERROR | Server error |
| 503 | SERVICE_UNAVAILABLE | Service temporarily unavailable |

## Performance Optimization

### 1. Request Batching

```go
type BatchRequest struct {
    Requests []json.RawMessage `json:"requests"`
    Parallel bool              `json:"parallel"`
}

func ProcessBatch(requests []json.RawMessage, parallel bool) []BatchResult {
    if parallel {
        return processParallel(requests)
    }
    return processSequential(requests)
}
```

### 2. Caching Strategy

```go
type CacheKey struct {
    Function   string
    Text       string
    Parameters string
}

var cache = NewLRUCache(1000) // 1000 entries

func getCached(key CacheKey) (interface{}, bool) {
    return cache.Get(key)
}

func setCache(key CacheKey, value interface{}, ttl time.Duration) {
    cache.SetWithTTL(key, value, ttl)
}
```

### 3. Rate Limiting

```go
type RateLimiter struct {
    limiter *rate.Limiter
    limits  map[string]*rate.Limiter // Per-API-key limits
}

func (rl *RateLimiter) Allow(apiKey string) bool {
    if limiter, exists := rl.limits[apiKey]; exists {
        return limiter.Allow()
    }
    return rl.limiter.Allow() // Default limiter
}
```

## Deployment

### Docker Configuration

```dockerfile
# Dockerfile
FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY . .

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o textlib-api cmd/server/main.go

FROM alpine:latest

RUN apk --no-cache add ca-certificates
WORKDIR /root/

COPY --from=builder /app/textlib-api .
COPY --from=builder /app/api ./api

EXPOSE 8080

CMD ["./textlib-api"]
```

### Kubernetes Deployment

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: textlib-api
  labels:
    app: textlib-api
spec:
  replicas: 3
  selector:
    matchLabels:
      app: textlib-api
  template:
    metadata:
      labels:
        app: textlib-api
    spec:
      containers:
      - name: textlib-api
        image: caiatech/textlib-api:1.0.0
        ports:
        - containerPort: 8080
        resources:
          requests:
            memory: "256Mi"
            cpu: "250m"
          limits:
            memory: "512Mi"
            cpu: "1000m"
        env:
        - name: LOG_LEVEL
          value: "info"
        - name: MAX_REQUEST_SIZE
          value: "10485760" # 10MB
        livenessProbe:
          httpGet:
            path: /api/v1/health
            port: 8080
          initialDelaySeconds: 30
          periodSeconds: 10
        readinessProbe:
          httpGet:
            path: /api/v1/health
            port: 8080
          initialDelaySeconds: 5
          periodSeconds: 5
```

## Monitoring

### Prometheus Metrics

```go
var (
    requestDuration = prometheus.NewHistogramVec(
        prometheus.HistogramOpts{
            Name: "textlib_api_request_duration_seconds",
            Help: "Request duration in seconds",
        },
        []string{"method", "endpoint", "status"},
    )
    
    requestCount = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "textlib_api_requests_total",
            Help: "Total number of requests",
        },
        []string{"method", "endpoint", "status"},
    )
    
    activeRequests = prometheus.NewGauge(
        prometheus.GaugeOpts{
            Name: "textlib_api_active_requests",
            Help: "Number of active requests",
        },
    )
)
```

### Grafana Dashboard

```json
{
  "dashboard": {
    "title": "TextLib API Monitoring",
    "panels": [
      {
        "title": "Request Rate",
        "targets": [{
          "expr": "rate(textlib_api_requests_total[5m])"
        }]
      },
      {
        "title": "Response Time P95",
        "targets": [{
          "expr": "histogram_quantile(0.95, rate(textlib_api_request_duration_seconds_bucket[5m]))"
        }]
      },
      {
        "title": "Algorithm Usage",
        "targets": [{
          "expr": "sum by (algorithm) (rate(textlib_algorithm_usage_total[5m]))"
        }]
      }
    ]
  }
}
```

### Logging

```go
type LogEntry struct {
    Timestamp   time.Time              `json:"timestamp"`
    Level       string                 `json:"level"`
    Method      string                 `json:"method"`
    Path        string                 `json:"path"`
    StatusCode  int                    `json:"status_code"`
    Duration    time.Duration          `json:"duration"`
    ClientIP    string                 `json:"client_ip"`
    UserAgent   string                 `json:"user_agent"`
    RequestID   string                 `json:"request_id"`
    Error       string                 `json:"error,omitempty"`
    Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

func LogRequest(entry LogEntry) {
    log.WithFields(log.Fields{
        "timestamp":   entry.Timestamp,
        "method":      entry.Method,
        "path":        entry.Path,
        "status_code": entry.StatusCode,
        "duration_ms": entry.Duration.Milliseconds(),
        "client_ip":   entry.ClientIP,
        "request_id":  entry.RequestID,
    }).Info("API request")
}
```

## Security Considerations

1. **Input Validation**: All inputs are validated for size and content
2. **Rate Limiting**: Prevent abuse through configurable rate limits
3. **API Authentication**: Optional API key authentication for production
4. **CORS Configuration**: Configurable CORS headers
5. **Request Size Limits**: Maximum request body size enforcement
6. **Timeout Protection**: Request timeout configuration
7. **SQL Injection**: Not applicable (no database queries)
8. **XSS Protection**: JSON responses with proper content-type headers

## Testing

### Integration Tests

```go
func TestComplexityEndpoint(t *testing.T) {
    // Setup
    router := setupRouter()
    
    // Test request
    body := `{"text": "Test text", "depth": 2}`
    req, _ := http.NewRequest("POST", "/api/v1/analyze/complexity", strings.NewReader(body))
    req.Header.Set("Content-Type", "application/json")
    
    // Execute
    rr := httptest.NewRecorder()
    router.ServeHTTP(rr, req)
    
    // Assert
    assert.Equal(t, http.StatusOK, rr.Code)
    
    var resp ComplexityResponse
    err := json.Unmarshal(rr.Body.Bytes(), &resp)
    assert.NoError(t, err)
    assert.NotZero(t, resp.LexicalComplexity)
}
```

### Load Testing

```bash
# Using Apache Bench
ab -n 10000 -c 100 -p request.json -T application/json \
   http://localhost:8080/api/v1/analyze/complexity

# Using k6
k6 run load-test.js
```

## Conclusion

This implementation guide provides a production-ready REST API for the TextLib RL-optimized functions. The API is designed for scalability, monitoring, and easy deployment in cloud environments.