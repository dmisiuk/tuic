# Common Issues and Solutions

This guide covers common issues you might encounter while using the CCPM Calculator and their solutions.

## 🚀 Installation Issues

### Go Version Compatibility

**Problem**: Build fails with Go version errors

```bash
# Error message
go version must be 1.25.1 or higher
```

**Solution**:
```bash
# Check Go version
go version

# If outdated, update Go
# macOS: brew upgrade go
# Ubuntu: sudo apt-get update && sudo apt-get install golang-go
# Windows: Download from https://golang.org/dl/
```

### Module Dependencies

**Problem**: `go mod tidy` fails or cannot find dependencies

```bash
# Error message
cannot find module "ccpm-demo/internal/calculator"
```

**Solution**:
```bash
# Clean module cache
go clean -modcache
go mod tidy

# Verify module path
go mod verify

# If issues persist, reinitialize
rm go.mod go.sum
go mod init ccpm-demo
go mod tidy
```

### Build Errors

**Problem**: Compilation fails with syntax errors

```bash
# Error message
syntax error: unexpected newline, expecting comma or }
```

**Solution**:
```bash
# Format code
go fmt ./...

# Check for syntax errors
go build -v ./...

# If still failing, check Go version compatibility
```

## 🧮 Calculation Issues

### Division by Zero

**Problem**: Returns error when dividing by zero

```go
// Code
result, err := engine.Evaluate("5 / 0")
// Error: division by zero
```

**Solution**: Handle the error gracefully

```go
result, err := engine.Evaluate("5 / 0")
if err != nil {
    if errors.Is(err, calculator.ErrDivisionByZero) {
        fmt.Println("Cannot divide by zero")
        engine.Clear()  // Reset calculator
        return
    }
    fmt.Printf("Error: %v\n", err)
    return
}
fmt.Printf("Result: %v\n", result)
```

### Invalid Expressions

**Problem**: Expression parsing fails with syntax errors

```go
// Problem cases
engine.Evaluate("2 + * 3")      // Invalid operator sequence
engine.Evaluate("(2 + 3")       // Mismatched parentheses
engine.Evaluate("2 ++ 3")       // Duplicate operators
```

**Solution**: Validate expressions before evaluation

```go
func isValidExpression(expr string) bool {
    // Remove spaces
    expr = strings.ReplaceAll(expr, " ", "")

    // Check for empty expression
    if len(expr) == 0 {
        return false
    }

    // Check for invalid characters
    for _, char := range expr {
        if !strings.Contains("0123456789+-*/.()", string(char)) {
            return false
        }
    }

    return true
}

// Usage
if !isValidExpression("2 + * 3") {
    fmt.Println("Invalid expression")
    return
}
```

### Number Overflow

**Problem**: Calculations result in very large or very small numbers

```go
// Example
result, _ := engine.Evaluate("1e300 * 1e300")  // Potential overflow
```

**Solution**: The engine automatically validates numbers

```go
result, err := engine.Evaluate("1e300 * 1e300")
if err != nil {
    if errors.Is(err, calculator.ErrNumberOutOfRange) {
        fmt.Println("Number out of range")
        return
    }
    fmt.Printf("Error: %v\n", err)
    return
}
```

## 🏗️ Development Issues

### Test Failures

**Problem**: Unit tests are failing

```bash
# Error message
--- FAIL: TestEngine_Evaluate (0.00s)
    engine_test.go:45: Engine.Evaluate() = 6, want 5
```

**Solution**:
```bash
# Run tests with verbose output
go test -v ./internal/calculator

# Run specific test
go test -run TestEngine_Evaluate ./internal/calculator

# Generate coverage report
go test -coverprofile=coverage.out ./internal/calculator
go tool cover -html=coverage.out
```

### Import Path Issues

**Problem**: Cannot import internal packages

```go
// Error
import "ccpm-demo/internal/calculator" // not found
```

**Solution**: Ensure you're in the correct directory and module is properly initialized

```bash
# Check current directory
pwd

# Initialize module if needed
go mod init ccpm-demo

# Check go.mod content
cat go.mod
```

### Race Conditions

**Problem**: Concurrent access to engine causes data races

```go
// Problem: Multiple goroutines accessing same engine
go func() {
    engine.SetValue(10)
}()

go func() {
    engine.Add(5)
}()
```

**Solution**: Use mutex for concurrent access

```go
type SafeEngine struct {
    mu    sync.Mutex
    *calculator.Engine
}

func (se *SafeEngine) Evaluate(expr string) (float64, error) {
    se.mu.Lock()
    defer se.mu.Unlock()
    return se.Engine.Evaluate(expr)
}
```

## 🎯 Performance Issues

### Memory Usage

**Problem**: High memory usage with many calculations

```bash
# Check memory usage
go test -bench=. -benchmem ./internal/calculator
```

**Solution**: Reset engine state periodically

```go
// After complex calculations
engine.Clear()

// Or create new engine
engine = calculator.NewEngine()
```

### Slow Calculations

**Problem**: Complex expressions take too long

```bash
# Profile performance
go test -cpuprofile=cpu.out -bench=. ./internal/calculator
go tool pprof cpu.out
```

**Solution**: Break down complex expressions

```go
// Instead of: "((1+2)*(3+4)*(5+6))/(7+8)"
// Break into steps:
result, _ := engine.Evaluate("1+2")  // 3
engine.SetValue(result)
result, _ = engine.Evaluate("3+4")   // 7
result, _ = engine.Multiply(7)       // 21
```

## 📱 UI/UX Issues

### Input Handling

**Problem**: User input not properly validated

```go
// Invalid input
engine.InputNumber(12)  // Error: digit must be 0-9
```

**Solution**: Validate input before processing

```go
func isValidDigit(digit int) bool {
    return digit >= 0 && digit <= 9
}

// Usage
if !isValidDigit(digit) {
    fmt.Println("Invalid digit")
    return
}
err := engine.InputNumber(digit)
```

### State Management

**Problem**: Calculator state becomes inconsistent

```go
// Common issue
engine.SetValue(10)
engine.InputNumber(5)          // Should clear first
result, _ := engine.Add(5)     // Unexpected result
```

**Solution**: Always check shouldClear flag

```go
if engine.ShouldClear() {
    engine.ClearEntry()
}
err := engine.InputNumber(digit)
```

## 🛠️ Debugging Tips

### Enable Verbose Logging

```go
// Create a wrapper with logging
type LoggingEngine struct {
    *calculator.Engine
}

func (le *LoggingEngine) Evaluate(expr string) (float64, error) {
    fmt.Printf("Evaluating: %s\n", expr)
    result, err := le.Engine.Evaluate(expr)
    if err != nil {
        fmt.Printf("Error: %v\n", err)
    } else {
        fmt.Printf("Result: %f\n", result)
    }
    return result, err
}
```

### State Inspection

```go
// Helper function to inspect engine state
func inspectEngine(engine *calculator.Engine) {
    fmt.Printf("Current: %f\n", engine.GetValue())
    fmt.Printf("Entry: %f\n", engine.GetEntryValue())
    fmt.Printf("ShouldClear: %t\n", engine.ShouldClear())
}
```

### Test with Known Values

```go
// Test cases for debugging
testCases := []struct {
    expr string
    want float64
}{
    {"2 + 2", 4},
    {"10 / 2", 5},
    {"3 * 4", 12},
    {"15 - 7", 8},
}

for _, tc := range testCases {
    got, err := engine.Evaluate(tc.expr)
    if err != nil {
        fmt.Printf("Error evaluating %s: %v\n", tc.expr, err)
        continue
    }
    if got != tc.want {
        fmt.Printf("%s = %f, want %f\n", tc.expr, got, tc.want)
    }
}
```

## 🆘 Getting Help

### Check Resources

1. **API Documentation**: Check [API Reference](../api/engine.md)
2. **User Guide**: Read [User Guide](../user-guide/quickstart.md)
3. **Examples**: See [Examples](../examples/basic-examples.md)

### Search Issues

```bash
# Search existing issues
gh issue list --state open --search "calculation error"

# Create new issue
gh issue create --title "Problem with division" --body "Description of issue"
```

### Community Support

- **GitHub Discussions**: Ask questions
- **Stack Overflow**: Use tag `ccpm-calculator`
- **Email**: Contact development team

---

Remember: Most issues can be resolved by checking input validation, handling errors properly, and understanding the calculator's state management.