# Contributing Guidelines

Thank you for your interest in contributing to the CCPM Demo Project! This document provides guidelines and instructions for contributors.

## 🤝 How to Contribute

### 1. Getting Started

1. **Fork the Repository**
   ```bash
   # Fork the repository on GitHub
   # Clone your fork locally
   git clone https://github.com/your-username/ccpm-demo.git
   cd ccpm-demo
   ```

2. **Set Up Development Environment**
   ```bash
   # Install dependencies
   go mod tidy

   # Run tests to verify setup
   go test ./...
   ```

### 2. Development Workflow

#### Branch Naming Convention
- Feature branches: `feature/your-feature-name`
- Bug fixes: `fix/issue-description`
- Documentation: `docs/documentation-update`
- Hot fixes: `hotfix/critical-bug-fix`

#### Example Branch Names
```
feature/add-scientific-notation
fix/division-precision-issue
docs/update-api-documentation
hotfix/memory-leak-fix
```

#### Making Changes
1. Create a new branch:
   ```bash
   git checkout -b feature/your-feature-name
   ```

2. Make your changes and write tests:
   ```bash
   # Edit files
   # Write tests for new functionality
   ```

3. Run tests and linting:
   ```bash
   go test ./...
   go fmt ./...
   ```

4. Commit your changes:
   ```bash
   git add .
   git commit -m "Issue #123: Add scientific notation support"
   ```

5. Push to your fork:
   ```bash
   git push origin feature/your-feature-name
   ```

6. Create a Pull Request

### 3. Pull Request Process

#### PR Requirements
- **Title**: Use format "Issue #123: Your change description"
- **Description**: Include what was changed and why
- **Tests**: All tests must pass
- **Documentation**: Update relevant documentation
- **Code Style**: Follow Go conventions

#### PR Template
```markdown
## Description
Brief description of what this PR does and why it's needed.

## Changes
- List of changes made
- Files modified

## Testing
- How was this tested?
- Test cases added

## Related Issues
- Fixes #123
- Related to #456

## Checklist
- [ ] Tests pass
- [ ] Code follows style guide
- [ ] Documentation updated
- [ ] No breaking changes
```

## 📝 Code Style

### Go Conventions
- Follow [Effective Go](https://golang.org/doc/effective_go.html)
- Use `gofmt` for formatting
- Keep lines under 120 characters
- Use meaningful variable names

### Example Code Style
```go
// Good
func calculateTotal(items []Item) (float64, error) {
    if len(items) == 0 {
        return 0, ErrEmptyItems
    }

    var total float64
    for _, item := range items {
        total += item.Price
    }

    return total, nil
}

// Bad
func calc(it []Item) (float64, error) {
    var t float64
    for i := range it {
        t += it[i].Price
    }
    return t, nil
}
```

### Documentation Standards
- Document all public functions
- Include parameter descriptions
- Provide usage examples
- Document error cases

```go
// Evaluate evaluates a mathematical expression and returns the result.
//
// Parameters:
//   - expression: Mathematical expression to evaluate (e.g., "2 + 3 * 4")
//
// Returns:
//   - float64: Result of the evaluation
//   - error: Error if evaluation fails (e.g., invalid expression, division by zero)
//
// Example:
//   result, err := engine.Evaluate("(10 + 5) * 3")
//   if err != nil {
//       log.Fatal(err)
//   }
//   fmt.Println(result) // Output: 45
func (e *Engine) Evaluate(expression string) (float64, error) {
    // Implementation
}
```

## 🧪 Testing

### Test Requirements
- Write tests for all new functionality
- Maintain test coverage above 80%
- Test both success and error cases
- Use table-driven tests for multiple test cases

### Test Structure
```go
func TestEngine_Evaluate(t *testing.T) {
    tests := []struct {
        name        string
        expression  string
        want        float64
        wantErr     bool
        errorType   error
    }{
        {
            name:       "simple addition",
            expression: "2 + 3",
            want:       5,
            wantErr:    false,
        },
        {
            name:       "division by zero",
            expression: "5 / 0",
            wantErr:    true,
            errorType:  ErrDivisionByZero,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            engine := NewEngine()
            got, err := engine.Evaluate(tt.expression)

            if tt.wantErr {
                if err == nil {
                    t.Errorf("Engine.Evaluate() expected error, got nil")
                    return
                }
                if tt.errorType != nil && !errors.Is(err, tt.errorType) {
                    t.Errorf("Engine.Evaluate() expected error type %v, got %v", tt.errorType, err)
                }
                return
            }

            if err != nil {
                t.Errorf("Engine.Evaluate() unexpected error = %v", err)
                return
            }

            if got != tt.want {
                t.Errorf("Engine.Evaluate() = %v, want %v", got, tt.want)
            }
        })
    }
}
```

### Running Tests
```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests with verbose output
go test -v ./...

# Run specific test
go test -run TestEngine_Evaluate ./internal/calculator

# Generate coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## 🐛 Bug Reports

### Reporting Bugs
1. Use GitHub Issues
2. Include reproduction steps
3. Provide expected vs actual behavior
4. Include environment information

### Bug Report Template
```markdown
## Bug Description
Brief description of the bug.

## Steps to Reproduce
1. Step one
2. Step two
3. Step three

## Expected Behavior
What should happen.

## Actual Behavior
What actually happens.

## Environment
- OS: [e.g., macOS 12.0]
- Go version: [e.g., 1.18.0]
- Project version: [e.g., v1.0.0]

## Additional Information
Any additional context, logs, or screenshots.
```

## ✨ Feature Requests

### Requesting Features
1. Check existing issues first
2. Use GitHub Issues with "Feature Request" label
3. Describe the use case
4. Suggest implementation approach if possible

### Feature Request Template
```markdown
## Feature Description
What feature would you like to see added?

## Use Case
Describe the problem this feature would solve.

## Proposed Solution
How do you envision this feature working?

## Alternatives Considered
What alternatives have you considered?

## Additional Context
Any other context, mockups, or examples.
```

## 📞 Getting Help

- **GitHub Discussions**: For general questions and discussions
- **GitHub Issues**: For bug reports and feature requests
- **Email**: For private questions and security issues

## 📄 Code of Conduct

Please note that this project adheres to a [Code of Conduct](CODE_OF_CONDUCT.md). By participating in this project, you agree to abide by its terms.

---

Happy coding! 🚀