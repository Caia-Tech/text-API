# TextLib API Client SDK Examples

This document provides client SDK examples in multiple programming languages for the TextLib RL-Optimized API.

## Table of Contents

- [Go Client](#go-client)
- [Python Client](#python-client)
- [JavaScript/TypeScript Client](#javascripttypescript-client)
- [Java Client](#java-client)
- [cURL Examples](#curl-examples)

## Go Client

### Installation

```bash
go get github.com/caiatech/textlib-client-go
```

### Usage Example

```go
package main

import (
    "context"
    "fmt"
    "log"
    
    "github.com/caiatech/textlib-client-go"
)

func main() {
    // Create client
    client := textlib.NewClient(
        textlib.WithBaseURL("https://api.textlib.caiatech.com/v1"),
        textlib.WithAPIKey("your-api-key"),
        textlib.WithTimeout(30 * time.Second),
    )
    
    ctx := context.Background()
    
    // Analyze text complexity
    complexityReq := &textlib.ComplexityRequest{
        Text:  "The quantum mechanics principles demonstrate wave-particle duality.",
        Depth: 2,
    }
    
    complexity, err := client.AnalyzeComplexity(ctx, complexityReq)
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Lexical Complexity: %.2f\n", complexity.LexicalComplexity)
    fmt.Printf("Readability (Flesch-Kincaid): %.1f\n", 
        complexity.ReadabilityScores["flesch-kincaid"])
    
    // Extract key phrases
    phrasesReq := &textlib.KeyPhrasesRequest{
        Text:       "Machine learning transforms business operations through automation.",
        MaxPhrases: 5,
    }
    
    phrases, err := client.ExtractKeyPhrases(ctx, phrasesReq)
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Println("\nKey Phrases:")
    for _, phrase := range phrases.Phrases {
        fmt.Printf("- %s (score: %.2f, confidence: %.2f)\n", 
            phrase.Text, phrase.Score, phrase.Confidence)
    }
}
```

### Advanced Usage

```go
// Batch processing
requests := []textlib.BatchItem{
    {
        Type: "complexity",
        Data: textlib.ComplexityRequest{Text: "First text", Depth: 1},
    },
    {
        Type: "keyphrases",
        Data: textlib.KeyPhrasesRequest{Text: "Second text", MaxPhrases: 10},
    },
}

results, err := client.ProcessBatch(ctx, requests, true) // parallel=true

// With retry and circuit breaker
client := textlib.NewClient(
    textlib.WithRetry(3, 1*time.Second),
    textlib.WithCircuitBreaker(10, 5*time.Minute),
)
```

## Python Client

### Installation

```bash
pip install textlib-client
```

### Usage Example

```python
from textlib import TextLibClient
from textlib.models import ComplexityRequest, KeyPhrasesRequest

# Initialize client
client = TextLibClient(
    base_url="https://api.textlib.caiatech.com/v1",
    api_key="your-api-key"
)

# Analyze complexity
complexity = client.analyze_complexity(
    text="The epistemological implications require careful consideration.",
    depth=3
)

print(f"Lexical Complexity: {complexity.lexical_complexity:.2f}")
print(f"Semantic Complexity: {complexity.semantic_complexity:.2f}")
print(f"Algorithm Used: {complexity.algorithm_used}")

# Extract key phrases with different algorithms
for max_phrases in [5, 25, 60]:
    result = client.extract_key_phrases(
        text="Artificial intelligence and machine learning revolutionize healthcare.",
        max_phrases=max_phrases
    )
    
    print(f"\nAlgorithm: {result.metadata.algorithm} ({max_phrases} phrases)")
    for phrase in result.phrases[:3]:
        print(f"  - {phrase.text}: {phrase.score:.2f}")

# Detect language
language = client.detect_language(
    text="Bonjour, comment allez-vous?",
    confidence=0.9
)

print(f"\nDetected Language: {language.language} ({language.confidence:.2f})")
```

### Async Support

```python
import asyncio
from textlib import AsyncTextLibClient

async def process_texts(texts):
    async with AsyncTextLibClient(api_key="your-api-key") as client:
        tasks = [
            client.analyze_complexity(text, depth=2)
            for text in texts
        ]
        
        results = await asyncio.gather(*tasks)
        return results

# Run async
texts = ["Text 1", "Text 2", "Text 3"]
results = asyncio.run(process_texts(texts))
```

## JavaScript/TypeScript Client

### Installation

```bash
npm install @caiatech/textlib-client
```

### Usage Example

```javascript
import { TextLibClient } from '@caiatech/textlib-client';

// Initialize client
const client = new TextLibClient({
    baseURL: 'https://api.textlib.caiatech.com/v1',
    apiKey: 'your-api-key',
});

// Analyze complexity
async function analyzeText() {
    try {
        const complexity = await client.analyzeComplexity({
            text: 'Complex scientific theories require deep understanding.',
            depth: 2,
        });
        
        console.log('Complexity Analysis:');
        console.log(`- Lexical: ${complexity.lexicalComplexity.toFixed(2)}`);
        console.log(`- Syntactic: ${complexity.syntacticComplexity.toFixed(2)}`);
        console.log(`- Processing Time: ${complexity.processingTime}`);
        
    } catch (error) {
        console.error('Error:', error.message);
    }
}

// Extract key phrases with progress callback
async function extractPhrases() {
    const result = await client.extractKeyPhrases({
        text: 'Natural language processing enables computers to understand text.',
        maxPhrases: 20,
    }, {
        onProgress: (progress) => {
            console.log(`Processing: ${progress}%`);
        }
    });
    
    console.log('\nTop Key Phrases:');
    result.phrases
        .slice(0, 5)
        .forEach(phrase => {
            console.log(`- "${phrase.text}" (${phrase.category})`);
        });
}
```

### TypeScript Example

```typescript
import { 
    TextLibClient, 
    ComplexityRequest, 
    ComplexityResponse,
    KeyPhrase 
} from '@caiatech/textlib-client';

class TextAnalyzer {
    private client: TextLibClient;
    
    constructor(apiKey: string) {
        this.client = new TextLibClient({ apiKey });
    }
    
    async analyzeDocument(text: string): Promise<DocumentAnalysis> {
        // Parallel requests
        const [complexity, keyPhrases, language] = await Promise.all([
            this.client.analyzeComplexity({ text, depth: 2 }),
            this.client.extractKeyPhrases({ text, maxPhrases: 10 }),
            this.client.detectLanguage({ text, confidence: 0.8 }),
        ]);
        
        return {
            complexity,
            keyPhrases: keyPhrases.phrases,
            language: language.language,
            processingTime: this.sumProcessingTimes([
                complexity.processingTime,
                keyPhrases.metadata.processingTime,
                language.processingTime,
            ]),
        };
    }
    
    private sumProcessingTimes(times: string[]): number {
        return times.reduce((sum, time) => {
            const ms = parseInt(time.replace('ms', ''));
            return sum + ms;
        }, 0);
    }
}

interface DocumentAnalysis {
    complexity: ComplexityResponse;
    keyPhrases: KeyPhrase[];
    language: string;
    processingTime: number;
}
```

## Java Client

### Maven Dependency

```xml
<dependency>
    <groupId>com.caiatech</groupId>
    <artifactId>textlib-client</artifactId>
    <version>1.0.0</version>
</dependency>
```

### Usage Example

```java
import com.caiatech.textlib.TextLibClient;
import com.caiatech.textlib.models.*;

public class TextAnalysisExample {
    public static void main(String[] args) {
        // Create client
        TextLibClient client = TextLibClient.builder()
            .baseUrl("https://api.textlib.caiatech.com/v1")
            .apiKey("your-api-key")
            .connectionTimeout(10000)
            .readTimeout(30000)
            .build();
        
        try {
            // Analyze complexity
            ComplexityRequest request = ComplexityRequest.builder()
                .text("Advanced algorithms require careful implementation.")
                .depth(2)
                .build();
            
            ComplexityResponse response = client.analyzeComplexity(request);
            
            System.out.printf("Lexical Complexity: %.2f%n", 
                response.getLexicalComplexity());
            System.out.printf("Flesch-Kincaid Score: %.1f%n",
                response.getReadabilityScores().get("flesch-kincaid"));
            
            // Extract key phrases with streaming
            KeyPhrasesRequest phrasesReq = KeyPhrasesRequest.builder()
                .text("Machine learning applications in healthcare.")
                .maxPhrases(15)
                .build();
            
            client.extractKeyPhrasesStream(phrasesReq)
                .forEach(phrase -> {
                    System.out.printf("Found: %s (%.2f)%n", 
                        phrase.getText(), phrase.getScore());
                });
                
        } catch (TextLibException e) {
            System.err.println("API Error: " + e.getMessage());
        }
    }
}
```

### Spring Boot Integration

```java
@RestController
@RequestMapping("/api/text-analysis")
public class TextAnalysisController {
    
    private final TextLibClient textLibClient;
    
    @Autowired
    public TextAnalysisController(TextLibClient textLibClient) {
        this.textLibClient = textLibClient;
    }
    
    @PostMapping("/analyze")
    public AnalysisResult analyzeText(@RequestBody AnalysisRequest request) {
        // Analyze with caching
        String cacheKey = DigestUtils.md5Hex(request.getText());
        
        return cacheManager.get(cacheKey, () -> {
            ComplexityResponse complexity = textLibClient.analyzeComplexity(
                ComplexityRequest.builder()
                    .text(request.getText())
                    .depth(request.getDepth())
                    .build()
            );
            
            KeyPhrasesResponse phrases = textLibClient.extractKeyPhrases(
                KeyPhrasesRequest.builder()
                    .text(request.getText())
                    .maxPhrases(10)
                    .build()
            );
            
            return new AnalysisResult(complexity, phrases);
        });
    }
}
```

## cURL Examples

### Basic Examples

```bash
# Analyze text complexity (depth 1 - fast)
curl -X POST https://api.textlib.caiatech.com/v1/analyze/complexity \
  -H "Content-Type: application/json" \
  -H "X-API-Key: your-api-key" \
  -d '{
    "text": "Simple text for analysis.",
    "depth": 1
  }'

# Extract key phrases (TF-IDF algorithm)
curl -X POST https://api.textlib.caiatech.com/v1/extract/keyphrases \
  -H "Content-Type: application/json" \
  -H "X-API-Key: your-api-key" \
  -d '{
    "text": "Artificial intelligence and machine learning are transforming industries.",
    "maxPhrases": 5
  }'

# Calculate readability metrics
curl -X POST https://api.textlib.caiatech.com/v1/calculate/readability \
  -H "Content-Type: application/json" \
  -H "X-API-Key: your-api-key" \
  -d '{
    "text": "The implementation requires careful consideration of multiple factors.",
    "algorithms": ["flesch", "gunning-fog", "coleman-liau"]
  }'

# Detect language with high confidence
curl -X POST https://api.textlib.caiatech.com/v1/detect/language \
  -H "Content-Type: application/json" \
  -H "X-API-Key: your-api-key" \
  -d '{
    "text": "Hola, ¿cómo estás?",
    "confidence": 0.95
  }'
```

### Advanced Examples

```bash
# Analyze long document with deep analysis
curl -X POST https://api.textlib.caiatech.com/v1/analyze/complexity \
  -H "Content-Type: application/json" \
  -H "X-API-Key: your-api-key" \
  -d @- << EOF
{
  "text": "$(cat document.txt)",
  "depth": 3
}
EOF

# Extract many phrases (uses deep NLP)
curl -X POST https://api.textlib.caiatech.com/v1/extract/keyphrases \
  -H "Content-Type: application/json" \
  -H "X-API-Key: your-api-key" \
  -d '{
    "text": "'"$(cat article.txt | jq -Rs .)"'",
    "maxPhrases": 100
  }'

# Get performance metrics
curl -X GET https://api.textlib.caiatech.com/v1/metrics \
  -H "X-API-Key: your-api-key" | jq '.'

# Batch processing with jq
cat texts.json | jq '.texts[] | {text: ., depth: 2}' | \
while read -r request; do
  curl -X POST https://api.textlib.caiatech.com/v1/analyze/complexity \
    -H "Content-Type: application/json" \
    -H "X-API-Key: your-api-key" \
    -d "$request"
  sleep 0.1 # Rate limiting
done
```

### Error Handling

```bash
# Handle errors with proper status codes
response=$(curl -s -w "\n%{http_code}" -X POST \
  https://api.textlib.caiatech.com/v1/analyze/complexity \
  -H "Content-Type: application/json" \
  -H "X-API-Key: your-api-key" \
  -d '{"text": "", "depth": 2}')

body=$(echo "$response" | head -n -1)
status=$(echo "$response" | tail -n 1)

if [ "$status" -ne 200 ]; then
  echo "Error $status: $(echo "$body" | jq -r '.message')"
else
  echo "$body" | jq '.'
fi
```

## Rate Limiting and Best Practices

### Rate Limiting

All clients should implement exponential backoff:

```python
import time
from typing import Callable, Any

def with_retry(func: Callable, max_retries: int = 3) -> Any:
    for attempt in range(max_retries):
        try:
            return func()
        except RateLimitError as e:
            if attempt == max_retries - 1:
                raise
            wait_time = 2 ** attempt + 0.1 * random.random()
            time.sleep(wait_time)
```

### Best Practices

1. **Batch Requests**: Combine multiple analyses when possible
2. **Cache Results**: Implement client-side caching for repeated texts
3. **Async Processing**: Use async/await for better performance
4. **Error Handling**: Always handle API errors gracefully
5. **Monitoring**: Track API usage and performance metrics
6. **Text Size**: Keep individual texts under 50KB for optimal performance
7. **Connection Pooling**: Reuse HTTP connections

## SDK Development

To create your own SDK:

1. Use the OpenAPI specification at `/api/openapi.yaml`
2. Generate client code using OpenAPI Generator
3. Add retry logic and circuit breakers
4. Implement response caching
5. Add comprehensive error handling
6. Include usage examples and documentation

Example generation:

```bash
# Generate Python client
openapi-generator generate \
  -i https://api.textlib.caiatech.com/v1/openapi.yaml \
  -g python \
  -o ./python-client \
  --additional-properties packageName=textlib_client

# Generate TypeScript client
openapi-generator generate \
  -i https://api.textlib.caiatech.com/v1/openapi.yaml \
  -g typescript-axios \
  -o ./typescript-client \
  --additional-properties npmName=@caiatech/textlib-client
```