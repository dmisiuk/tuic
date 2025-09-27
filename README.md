# CCPM Demo Project

A comprehensive calculator application built with Go, featuring advanced mathematical operations and a robust parsing engine.

## 📋 Overview

This project implements a powerful calculator with support for complex mathematical expressions, proper error handling, and a clean architecture. The calculator engine handles basic arithmetic operations, expression parsing, and provides a solid foundation for building calculator applications.

## ✨ Features

- **Basic Operations**: Addition, subtraction, multiplication, division
- **Expression Parsing**: Support for complex mathematical expressions
- **Error Handling**: Comprehensive error handling with specific error types
- **State Management**: Proper calculator state management (C, CE functionality)
- **Number Validation**: Overflow/underflow protection
- **Parentheses Support**: Full support for nested expressions
- **Decimal Numbers**: Accurate floating-point arithmetic

## 🚀 Installation

### Prerequisites

- Go 1.25.1 or higher
- Git

### Clone the Repository

```bash
git clone https://github.com/your-username/ccpm-demo.git
cd ccpm-demo
```

### Build the Project

```bash
go mod tidy
go build -o calculator
```

### Run Tests

```bash
go test ./...
```

## 📖 Usage

### Basic Usage

```go
package main

import (
    "fmt"
    "ccpm-demo/internal/calculator"
)

func main() {
    engine := calculator.NewEngine()

    // Simple operations
    result, err := engine.Evaluate("2 + 3 * 4")
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        return
    }
    fmt.Printf("Result: %f\n", result) // Output: 14.000000

    // Complex expressions
    result, err = engine.Evaluate("(10 + 5) * (8 - 3) / 2")
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        return
    }
    fmt.Printf("Result: %f\n", result) // Output: 37.500000
}
```

### Engine State Management

```go
engine := calculator.NewEngine()

// Clear all values (C)
engine.Clear()

// Clear current entry (CE)
engine.ClearEntry()

// Input numbers digit by digit
err := engine.InputNumber(1)
err = engine.InputNumber(2)
err = engine.InputNumber(3)

// Get current value
value := engine.GetValue()
```

## 🏗️ Architecture

### Core Components

1. **Engine** (`/internal/calculator/engine.go`)
   - Main calculator state management
   - Basic arithmetic operations
   - Expression evaluation

2. **Parser** (`/internal/calculator/parser.go`)
   - Expression parsing and evaluation
   - Operator precedence handling
   - Parentheses support

3. **Error Handling** (`/internal/calculator/errors.go`)
   - Custom error types
   - Number validation
   - Error recovery

### Project Structure

```
ccpm-demo/
├── internal/
│   └── calculator/
│       ├── engine.go      # Main calculator engine
│       ├── parser.go      # Expression parser
│       ├── errors.go      # Error definitions
│       ├── engine_test.go # Engine tests
│       └── parser_test.go # Parser tests
├── docs/                  # Documentation
├── install/               # Installation scripts
├── ccpm/                  # CCPM framework
├── .claude/              # Claude configuration
└── README.md            # This file
```

## 🧪 Testing

The project includes comprehensive unit tests for all components:

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run specific test file
go test ./internal/calculator/engine_test.go

# Run tests with verbose output
go test -v ./...
```

## 🔧 API Reference

### Engine Methods

#### `NewEngine() *Engine`
Creates a new calculator engine instance.

#### `Evaluate(expression string) (float64, error)`
Evaluates a mathematical expression and returns the result.

**Parameters:**
- `expression`: Mathematical expression string

**Returns:**
- `float64`: Result of evaluation
- `error`: Error if evaluation fails

#### `Clear()`
Clears all calculator values (C functionality).

#### `ClearEntry()`
Clears the current entry (CE functionality).

#### `InputNumber(digit int) error`
Inputs a digit into the calculator.

**Parameters:**
- `digit`: Digit to input (0-9)

#### `PerformOperation(op string) (float64, error)`
Performs an arithmetic operation.

**Parameters:**
- `op`: Operation (+, -, *, /)

### Error Types

- `ErrEmptyExpression`: Empty expression provided
- `ErrInvalidExpression`: Invalid expression syntax
- `ErrInvalidNumber`: Invalid number format
- `ErrDivisionByZero`: Division by zero attempted
- `ErrInvalidOperator`: Invalid operator
- `ErrMismatchedParentheses`: Mismatched parentheses
- `ErrNumberOutOfRange`: Number out of valid range

## 🛠️ Development

### Code Style

The project follows Go standard conventions:
- Use `gofmt` for formatting
- Follow effective Go practices
- Write comprehensive tests
- Document public APIs

### Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests for new functionality
5. Ensure all tests pass
6. Submit a pull request

## 🐛 Troubleshooting

### Common Issues

**Build Errors**
```bash
# Ensure dependencies are up to date
go mod tidy

# Clean and rebuild
go clean -cache
go build
```

**Test Failures**
```bash
# Run tests with verbose output
go test -v ./...

# Check test coverage
go test -cover ./...
```

**Import Errors**
```bash
# Verify module path
go mod verify

# Update dependencies
go get -u ./...
```

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## 📞 Support

For support, please open an issue in the GitHub repository or contact the development team.

---

**CCPM Demo Project** - Built with ❤️ using Go
