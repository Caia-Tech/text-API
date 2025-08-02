# Contributing to TextLib

Thank you for your interest in contributing to TextLib! We welcome contributions from the community and appreciate your help in making this library better.

## Code of Conduct

By participating in this project, you agree to abide by our Code of Conduct. Please be respectful and professional in all interactions.

## Getting Started

### Prerequisites

- Go 1.19 or later
- Git
- Basic understanding of text processing and Go programming

### Setting Up Your Development Environment

1. Fork the repository on GitHub
2. Clone your fork locally:
   ```bash
   git clone https://github.com/yourusername/textlib.git
   cd textlib
   ```

3. Add the upstream repository:
   ```bash
   git remote add upstream https://github.com/caiatech/textlib.git
   ```

4. Create a new branch for your feature:
   ```bash
   git checkout -b feature/your-feature-name
   ```

## Development Guidelines

### Code Style

- Follow standard Go conventions and formatting
- Use `gofmt` to format your code
- Use `golint` and `go vet` to check for issues
- Write clear, self-documenting code with meaningful variable names
- Add comments for complex algorithms or business logic

### Testing

- Write comprehensive tests for all new functionality
- Maintain or improve test coverage (currently at 42.9%)
- Use table-driven tests where appropriate
- Include edge cases and error scenarios
- Test file names should end with `_test.go`

Example test structure:
```go
func TestYourFunction(t *testing.T) {
    tests := []struct {
        name     string
        input    string
        expected string
        wantErr  bool
    }{
        {
            name:     "valid input",
            input:    "test input",
            expected: "expected output",
            wantErr:  false,
        },
        // Add more test cases
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result, err := YourFunction(tt.input)
            if (err != nil) != tt.wantErr {
                t.Errorf("YourFunction() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            if result != tt.expected {
                t.Errorf("YourFunction() = %v, want %v", result, tt.expected)
            }
        })
    }
}
```

### Documentation

- Document all public functions, types, and methods
- Use Go's standard documentation format
- Include examples in documentation when helpful
- Update README.md if adding new major features

### Performance

- Consider performance implications of your changes
- Use benchmarks for performance-critical code
- Avoid unnecessary allocations in hot paths
- Use appropriate data structures for the use case

## Types of Contributions

### Bug Fixes

- Search existing issues before creating a new one
- Include a clear description of the bug and steps to reproduce
- Provide the expected vs. actual behavior
- Include relevant system information (Go version, OS, etc.)

### Feature Requests

- Discuss major changes in an issue before implementing
- Ensure the feature aligns with the library's goals
- Consider backward compatibility
- Update documentation and examples

### Documentation Improvements

- Fix typos and grammatical errors
- Improve clarity and completeness
- Add examples and use cases
- Update API documentation

## Submitting Changes

### Before Submitting

1. Ensure all tests pass:
   ```bash
   go test -v ./...
   ```

2. Check test coverage:
   ```bash
   go test -cover ./...
   ```

3. Run linting tools:
   ```bash
   go vet ./...
   golint ./...
   ```

4. Update documentation if necessary

### Pull Request Process

1. Update your branch with the latest upstream changes:
   ```bash
   git fetch upstream
   git rebase upstream/main
   ```

2. Push your changes to your fork:
   ```bash
   git push origin feature/your-feature-name
   ```

3. Create a pull request with:
   - Clear title and description
   - Reference any related issues
   - List of changes made
   - Testing information

4. Be responsive to code review feedback

### Pull Request Template

```markdown
## Description
Brief description of changes made.

## Type of Change
- [ ] Bug fix
- [ ] New feature
- [ ] Documentation update
- [ ] Performance improvement
- [ ] Code refactoring

## Testing
- [ ] Added new tests
- [ ] All existing tests pass
- [ ] Manual testing completed

## Checklist
- [ ] Code follows style guidelines
- [ ] Self-review completed
- [ ] Documentation updated
- [ ] No breaking changes (or marked as such)
```

## Areas for Contribution

### High Priority
- Improve test coverage (target: 60%+)
- Performance optimizations
- Memory usage improvements
- Better error handling and messages

### Medium Priority
- Additional language support for code analysis
- More sophisticated entity recognition
- Advanced mathematical expression parsing
- Enhanced security pattern detection

### Low Priority
- Additional file format support
- More text similarity algorithms
- Extended readability metrics
- Performance benchmarking suite

## Development Setup

### Running Tests

```bash
# Run all tests
go test ./...

# Run with verbose output
go test -v ./...

# Run with coverage
go test -cover ./...

# Run specific test
go test -run TestFunctionName
```

### Building

```bash
# Build the library
go build ./...

# Build with race detection
go build -race ./...
```

### Benchmarking

```bash
# Run benchmarks
go test -bench=. ./...

# Run specific benchmark
go test -bench=BenchmarkFunctionName
```

## Release Process

TextLib follows semantic versioning (SemVer):
- MAJOR: Incompatible API changes
- MINOR: New functionality (backward compatible)
- PATCH: Bug fixes (backward compatible)

## Getting Help

- Open an issue for bugs and feature requests
- Join discussions in existing issues
- Reach out to maintainers for guidance
- Check the documentation and examples first

## Recognition

Contributors are recognized in our:
- CONTRIBUTORS.md file
- Release notes for significant contributions
- GitHub contributors page

## License

By contributing to TextLib, you agree that your contributions will be licensed under the Apache License 2.0.

Thank you for contributing to TextLib! 🚀