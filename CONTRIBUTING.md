# Contributing to CCPM Demo Project

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

### Documentation Standards
- Document all public functions
- Include parameter descriptions
- Provide usage examples
- Document error cases

## 🧪 Testing

### Test Requirements
- Write tests for all new functionality
- Maintain test coverage above 80%
- Test both success and error cases
- Use table-driven tests for multiple test cases

### Running Tests
```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run specific test file
go test ./internal/calculator/engine_test.go
```

## 🐛 Bug Reports

### Reporting Bugs
1. Use GitHub Issues
2. Include reproduction steps
3. Provide expected vs actual behavior
4. Include environment information

## ✨ Feature Requests

### Requesting Features
1. Check existing issues first
2. Use GitHub Issues with "Feature Request" label
3. Describe the use case
4. Suggest implementation approach if possible

## 📞 Getting Help

- **GitHub Discussions**: For general questions and discussions
- **GitHub Issues**: For bug reports and feature requests
- **Documentation**: See [docs/](docs/) for detailed guides

## 📄 Code of Conduct

Please note that this project adheres to a Code of Conduct. By participating in this project, you agree to abide by its terms.

---

Happy coding! 🚀