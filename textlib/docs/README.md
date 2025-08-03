# TextLib RL-Optimized API Documentation

## Overview

This directory contains comprehensive documentation for the TextLib RL-Optimized API, which provides advanced text processing capabilities with adaptive algorithm selection and performance optimization through reinforcement learning.

## Documentation Contents

### 1. [OpenAPI Specification](./openapi.yaml)
Complete OpenAPI 3.0 specification defining all endpoints, request/response schemas, and error codes. This specification can be used to:
- Generate client SDKs in multiple languages
- Create interactive API documentation
- Validate API requests and responses
- Set up API mocking for testing

### 2. [Interactive API Documentation](./api-docs.html)
HTML page that renders the OpenAPI specification using Swagger UI. Open this file in a web browser to:
- Explore all available endpoints
- View detailed request/response schemas
- Test API endpoints directly from the browser
- Download the OpenAPI specification

### 3. [API Implementation Guide](./API_IMPLEMENTATION_GUIDE.md)
Detailed guide for implementing the REST API server, including:
- Complete HTTP server implementation in Go
- Request handlers for all endpoints
- Error handling patterns
- Performance optimization strategies
- Docker and Kubernetes deployment configurations
- Monitoring and logging setup

### 4. [Client SDK Examples](./CLIENT_SDK_EXAMPLES.md)
Code examples for using the API in multiple programming languages:
- **Go**: Native client with retry and circuit breaker
- **Python**: Sync and async clients with type hints
- **JavaScript/TypeScript**: Promise-based client with TypeScript support
- **Java**: Client with Spring Boot integration
- **cURL**: Command-line examples for all endpoints

### 5. [RL-Optimized API Extensions](./RL_OPTIMIZED_API_EXTENSIONS.md)
Technical specification for the RL-optimized functions, including:
- Detailed type definitions
- Algorithm complexity analysis
- Performance benchmarks
- RL integration strategies
- Migration path from legacy APIs

### 6. [Implementation Roadmap](./IMPLEMENTATION_ROADMAP.md)
12-week phased implementation plan covering:
- Core infrastructure setup
- Priority function implementation
- Testing and benchmarking
- Production deployment
- RL system integration

### 7. [API Examples](./API_EXAMPLES.md)
Practical code examples demonstrating:
- Real-world use cases
- Integration patterns
- Performance optimization
- Error handling
- Batch processing

## Quick Start

### Running the API Documentation

```bash
# View interactive documentation
open api-docs.html

# Or serve it with a local server
python -m http.server 8000
# Then navigate to http://localhost:8000/api-docs.html
```

### Example API Usage

```bash
# Analyze text complexity
curl -X POST http://localhost:8080/api/v1/analyze/complexity \
  -H "Content-Type: application/json" \
  -d '{
    "text": "Your text here",
    "depth": 2
  }'

# Extract key phrases
curl -X POST http://localhost:8080/api/v1/extract/keyphrases \
  -H "Content-Type: application/json" \
  -d '{
    "text": "Your text here",
    "maxPhrases": 10
  }'
```

## API Features

### Adaptive Algorithm Selection
- **AnalyzeTextComplexity**: O(n), O(n log n), or O(n²) algorithms based on depth parameter
- **ExtractKeyPhrases**: TF-IDF, statistical, or deep NLP based on maxPhrases parameter
- **CalculateReadabilityMetrics**: Multiple readability formulas with configurable selection
- **DetectLanguage**: Fast heuristics or comprehensive analysis based on confidence requirement

### Performance Optimization
- Real-time metrics collection
- Intelligent caching strategies
- Batch processing support
- Configurable resource limits

### Quality Assurance
- Every result includes quality metrics (accuracy, confidence, coverage)
- Processing time and memory usage tracking
- Algorithm transparency (which algorithm was used)

## Architecture

```
┌─────────────────┐     ┌─────────────────┐     ┌─────────────────┐
│   Client Apps   │────▶│   REST API      │────▶│  TextLib Core   │
└─────────────────┘     └─────────────────┘     └─────────────────┘
                               │                          │
                               ▼                          ▼
                        ┌─────────────────┐     ┌─────────────────┐
                        │    Metrics      │     │   Algorithm     │
                        │   Collection    │     │    Registry     │
                        └─────────────────┘     └─────────────────┘
```

## Development

### Building the API Server

```bash
# Install dependencies
go mod download

# Build the server
go build -o textlib-api cmd/server/main.go

# Run tests
go test ./...

# Run with Docker
docker build -t textlib-api .
docker run -p 8080:8080 textlib-api
```

### Generating Client SDKs

```bash
# Generate Python client
openapi-generator generate \
  -i openapi.yaml \
  -g python \
  -o ../clients/python

# Generate TypeScript client
openapi-generator generate \
  -i openapi.yaml \
  -g typescript-axios \
  -o ../clients/typescript
```

## Contributing

When contributing to the API documentation:

1. Update the OpenAPI specification for any API changes
2. Regenerate client examples if endpoints change
3. Update implementation guides with new patterns
4. Add examples for new use cases
5. Keep performance benchmarks up to date

## License

Copyright 2025 Caia Tech. Licensed under the Apache License, Version 2.0.

## Support

For questions or issues:
- GitHub Issues: [https://github.com/Caia-Tech/text-API/issues](https://github.com/Caia-Tech/text-API/issues)
- Email: support@caiatech.com