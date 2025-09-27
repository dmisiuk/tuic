# Quick Start Guide

Welcome to the CCPM Calculator! This guide will help you get started quickly with the calculator application.

## 🚀 Installation

### Prerequisites

- Go 1.25.1 or higher
- Git

### Quick Setup

```bash
# Clone the repository
git clone https://github.com/your-username/ccpm-demo.git
cd ccpm-demo

# Install dependencies
go mod tidy

# Build the application
go build -o calculator

# Run tests (optional)
go test ./...
```

## 🎯 Your First Calculation

### Using the Engine Directly

```go
package main

import (
    "fmt"
    "ccpm-demo/internal/calculator"
)

func main() {
    // Create a new calculator engine
    engine := calculator.NewEngine()

    // Basic addition
    result, err := engine.Evaluate("2 + 3")
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        return
    }
    fmt.Printf("2 + 3 = %v\n", result) // Output: 5

    // More complex expression
    result, err = engine.Evaluate("10 * (5 + 3) / 2")
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        return
    }
    fmt.Printf("10 * (5 + 3) / 2 = %v\n", result) // Output: 40
}
```

### Step-by-Step Number Input

```go
engine := calculator.NewEngine()

// Enter numbers digit by digit
engine.InputNumber(1)  // Entry: 1
engine.InputNumber(2)  // Entry: 12
engine.InputNumber(3)  // Entry: 123

// Perform addition
result, err := engine.PerformOperation("+")
fmt.Printf("123 + ? = %v\n", result) // Output: 123 + ? = 123

// Enter second number
engine.InputNumber(4)
engine.InputNumber(5)  // Entry: 45

// Get final result
result, err = engine.PerformOperation("=")
fmt.Printf("123 + 45 = %v\n", result) // Output: 123 + 45 = 168
```

## 📝 Common Operations

### Basic Arithmetic

```go
engine := calculator.NewEngine()

// Addition
result, _ := engine.Evaluate("15 + 25")      // 40

// Subtraction
result, _ := engine.Evaluate("50 - 15")      // 35

// Multiplication
result, _ := engine.Evaluate("7 * 8")        // 56

// Division
result, _ := engine.Evaluate("100 / 4")      // 25
```

### Working with Decimals

```go
// Decimal numbers are automatically handled
result, _ := engine.Evaluate("3.14 * 2")     // 6.28
result, _ := engine.Evaluate("10.5 / 2.5")   // 4.2
```

### Using Parentheses

```go
// Parentheses control order of operations
result, _ := engine.Evaluate("(2 + 3) * 4")     // 20 (not 14)
result, _ := engine.Evaluate("2 + (3 * 4)")     // 14
result, _ := engine.Evaluate("((1 + 2) * 3) + 4") // 13
```

## 🛠️ Calculator Controls

### Clear Functions

```go
engine := calculator.NewEngine()

// Enter some numbers
engine.InputNumber(1)
engine.InputNumber(2)
engine.InputNumber(3)

// Clear current entry (CE)
engine.ClearEntry()  // Clears just the entry value

// Clear everything (C)
engine.Clear()       // Resets calculator to initial state
```

### Getting Current Values

```go
engine := calculator.NewEngine()
engine.SetValue(42)

// Get current display value
current := engine.GetValue()           // 42

// Get entry value
entry := engine.GetEntryValue()        // 0

// Check if display should clear
shouldClear := engine.ShouldClear()     // true/false
```

## 🚨 Error Handling

### Handling Common Errors

```go
engine := calculator.NewEngine()

result, err := engine.Evaluate("5 / 0")
if err != nil {
    fmt.Printf("Calculation failed: %v\n", err)
    // Handle the error appropriately
    engine.Clear()  // Reset calculator state
}

// Continue with next calculation
result, err = engine.Evaluate("10 + 5")
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Result: %v\n", result)
```

### Error Types

- **Division by zero**: "5 / 0"
- **Invalid expression**: "2 + * 3"
- **Mismatched parentheses**: "(2 + 3"
- **Empty expression**: ""

## 🎨 Advanced Features

### Complex Expressions

```go
// You can chain operations
result, _ := engine.Evaluate("2 * 3 + 4 * 5 - 6 / 2") // 24

// Nested parentheses
result, _ := engine.Evaluate("((10 + 5) * (8 - 3)) / (2 + 1)") // 25
```

### State Management

```go
engine := calculator.NewEngine()

// The engine maintains state between operations
engine.SetValue(10)
result, _ := engine.Add(5)      // 15
result, _ = engine.Multiply(2)  // 30
result, _ = engine.Subtract(8)  // 22

// The current value is preserved
current := engine.GetValue()    // 22
```

## 📚 Next Steps

- Read the [Basic Usage Guide](basic-usage.md) for detailed examples
- Check the [API Reference](../api/engine.md) for complete method documentation
- Explore the [Advanced Features Guide](advanced-features.md)
- Learn about [Keyboard Shortcuts](keyboard-shortcuts.md)

## 🔧 Troubleshooting

**Build Issues**
```bash
# Clean and rebuild
go clean -cache
go mod tidy
go build
```

**Import Errors**
```bash
# Verify module path
go mod verify

# Update dependencies
go get -u ./...
```

For more troubleshooting help, see the [Troubleshooting Guide](../troubleshooting/common-issues.md).

---

Happy calculating! 🧮