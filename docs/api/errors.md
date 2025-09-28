# Errors API Reference

The calculator package provides comprehensive error handling with specific error types for different failure scenarios.

## Overview

The error handling system includes specific error types for various calculator operations, making it easy to identify and handle different types of errors appropriately.

## Table of Contents

- [Error Types](#error-types)
- [Validation Functions](#validation-functions)
- [Error Handling Patterns](#error-handling-patterns)
- [Best Practices](#best-practices)
- [Examples](#examples)

## Error Types

### ErrEmptyExpression

```go
var ErrEmptyExpression = errors.New("empty expression")
```

**Description**: Returned when an empty expression is provided for evaluation.

**When it occurs:**
- Calling `Evaluate("")`
- Calling `Parse("")`

**Example:**
```go
result, err := engine.Evaluate("")
// err = ErrEmptyExpression
```

### ErrInvalidExpression

```go
var ErrInvalidExpression = errors.New("invalid expression")
```

**Description**: Returned when the expression syntax is invalid.

**When it occurs:**
- Invalid operator sequences (e.g., "2 + * 3")
- Invalid character combinations
- Malformed expressions

**Example:**
```go
result, err := engine.Evaluate("2 + * 3")
// err = ErrInvalidExpression
```

### ErrInvalidNumber

```go
var ErrInvalidNumber = errors.New("invalid number")
```

**Description**: Returned when a number format is invalid.

**When it occurs:**
- Invalid decimal format (e.g., "3.14.15")
- Numbers with invalid characters
- Malformed numeric literals

**Example:**
```go
result, err := engine.Evaluate("3.14.15 + 2")
// err = ErrInvalidNumber
```

### ErrDivisionByZero

```go
var ErrDivisionByZero = errors.New("division by zero")
```

**Description**: Returned when division by zero is attempted.

**When it occurs:**
- Division by zero (e.g., "5 / 0")
- Division by a number that evaluates to zero

**Example:**
```go
result, err := engine.Evaluate("5 / 0")
// err = ErrDivisionByZero
```

### ErrInvalidOperator

```go
var ErrInvalidOperator = errors.New("invalid operator")
```

**Description**: Returned when an invalid operator is used.

**When it occurs:**
- Invalid arithmetic operator (e.g., "5 x 3")
- Unsupported operations

**Example:**
```go
result, err := engine.PerformOperation("x")
// err = ErrInvalidOperator
```

### ErrMismatchedParentheses

```go
var ErrMismatchedParentheses = errors.New("mismatched parentheses")
```

**Description**: Returned when parentheses are mismatched in expressions.

**When it occurs:**
- Unclosed parentheses (e.g., "(2 + 3")
- Extra closing parentheses (e.g., "2 + 3)")
- Unbalanced parentheses

**Example:**
```go
result, err := engine.Evaluate("(2 + 3")
// err = ErrMismatchedParentheses
```

### ErrNumberOutOfRange

```go
var ErrNumberOutOfRange = errors.New("number out of range")
```

**Description**: Returned when a number is outside the valid range.

**When it occurs:**
- Numbers that would cause overflow
- Numbers that would cause underflow
- Extremely large or small values

**Example:**
```go
result, err := engine.Evaluate("1e300 * 1e300")
// err = ErrNumberOutOfRange
```

### ErrInvalidDigit

```go
var ErrInvalidDigit = errors.New("invalid digit")
```

**Description**: Returned when an invalid digit is provided to InputNumber.

**When it occurs:**
- Digits outside 0-9 range
- Negative numbers
- Non-digit values

**Example:**
```go
err := engine.InputNumber(12)
// err = ErrInvalidDigit
```

## Validation Functions

### ValidateNumber

```go
func ValidateNumber(value float64) error
```

Validates that a number is within the acceptable range.

**Parameters:**
- `value`: Number to validate

**Returns:**
- `error`: ErrNumberOutOfRange if validation fails, nil otherwise

**Validation Criteria:**
- Not infinity or NaN
- Within reasonable floating-point range
- No overflow/underflow conditions

**Example:**
```go
err := calculator.ValidateNumber(123.45)
// err = nil (valid)

err = calculator.ValidateNumber(math.Inf(1))
// err = ErrNumberOutOfRange
```

## Error Handling Patterns

### Basic Error Handling

```go
func calculate(expression string) (float64, error) {
    engine := calculator.NewEngine()
    result, err := engine.Evaluate(expression)
    if err != nil {
        return 0, fmt.Errorf("calculation failed: %w", err)
    }
    return result, nil
}
```

### Specific Error Handling

```go
func safeCalculate(expression string) error {
    engine := calculator.NewEngine()
    result, err := engine.Evaluate(expression)

    if err != nil {
        switch {
        case errors.Is(err, ErrEmptyExpression):
            return fmt.Errorf("please enter an expression")
        case errors.Is(err, ErrDivisionByZero):
            return fmt.Errorf("cannot divide by zero")
        case errors.Is(err, ErrMismatchedParentheses):
            return fmt.Errorf("check your parentheses")
        case errors.Is(err, ErrInvalidExpression):
            return fmt.Errorf("invalid expression format")
        case errors.Is(err, ErrInvalidNumber):
            return fmt.Errorf("invalid number format")
        case errors.Is(err, ErrNumberOutOfRange):
            return fmt.Errorf("number too large or small")
        default:
            return fmt.Errorf("calculation error: %w", err)
        }
    }

    fmt.Printf("Result: %f\n", result)
    return nil
}
```

### Error Recovery

```go
func calculateWithRecovery(expression string) {
    engine := calculator.NewEngine()

    result, err := engine.Evaluate(expression)
    if err != nil {
        fmt.Printf("Error: %v\n", err)

        // Reset calculator state on error
        engine.Clear()

        // Provide user guidance
        switch {
        case errors.Is(err, ErrDivisionByZero):
            fmt.Println("Hint: Cannot divide by zero")
        case errors.Is(err, ErrMismatchedParentheses):
            fmt.Println("Hint: Check that all parentheses are properly closed")
        case errors.Is(err, ErrInvalidExpression):
            fmt.Println("Hint: Ensure expression format is correct")
        }
        return
    }

    fmt.Printf("Result: %f\n", result)
}
```

### Custom Error Messages

```go
func getErrorMessage(err error) string {
    switch {
    case errors.Is(err, ErrEmptyExpression):
        return "Please enter a mathematical expression"
    case errors.Is(err, ErrDivisionByZero):
        return "Division by zero is not allowed"
    case errors.Is(err, ErrMismatchedParentheses):
        return "Parentheses are mismatched. Please check your expression"
    case errors.Is(err, ErrInvalidExpression):
        return "Invalid expression syntax. Please check your input"
    case errors.Is(err, ErrInvalidNumber):
        return "Invalid number format. Please use valid numbers"
    case errors.Is(err, ErrNumberOutOfRange):
        return "Number is too large or too small for calculation"
    case errors.Is(err, ErrInvalidOperator):
        return "Invalid operator. Use +, -, *, or /"
    case errors.Is(err, ErrInvalidDigit):
        return "Invalid digit. Please use numbers 0-9"
    default:
        return fmt.Sprintf("Calculation error: %v", err)
    }
}
```

## Best Practices

### Always Check Errors

```go
// Good - Always check errors
result, err := engine.Evaluate("2 + 3")
if err != nil {
    // Handle error
}
fmt.Printf("Result: %f\n", result)

// Bad - Ignoring errors
result, _ := engine.Evaluate("2 + 3") // Potential panic
fmt.Printf("Result: %f\n", result)
```

### Use Specific Error Handling

```go
// Good - Specific error handling
result, err := engine.Evaluate(expression)
if errors.Is(err, ErrDivisionByZero) {
    fmt.Println("Cannot divide by zero")
    return
}

// Bad - Generic error handling
result, err := engine.Evaluate(expression)
if err != nil {
    fmt.Println("Some error occurred")
    return
}
```

### Wrap Errors with Context

```go
// Good - Error wrapping
result, err := engine.Evaluate(userInput)
if err != nil {
    return fmt.Errorf("failed to calculate '%s': %w", userInput, err)
}

// Acceptable - Direct error
result, err := engine.Evaluate(userInput)
if err != nil {
    return err
}
```

### Validate Input Before Processing

```go
func calculateSafely(expression string) (float64, error) {
    // Validate input first
    if len(strings.TrimSpace(expression)) == 0 {
        return 0, ErrEmptyExpression
    }

    // Basic character validation
    for _, char := range expression {
        if !strings.Contains("0123456789+-*/.() ", string(char)) {
            return 0, ErrInvalidExpression
        }
    }

    // Proceed with calculation
    engine := calculator.NewEngine()
    return engine.Evaluate(expression)
}
```

### Reset State on Error

```go
func calculateWithReset(expression string) error {
    engine := calculator.NewEngine()

    result, err := engine.Evaluate(expression)
    if err != nil {
        // Reset calculator state on error
        engine.Clear()
        return fmt.Errorf("calculation failed: %w", err)
    }

    fmt.Printf("Result: %f\n", result)
    return nil
}
```

## Examples

### User Input Handling

```go
func handleUserInput(input string) {
    engine := calculator.NewEngine()

    // Trim whitespace
    input = strings.TrimSpace(input)

    if input == "" {
        fmt.Println("Please enter an expression")
        return
    }

    result, err := engine.Evaluate(input)
    if err != nil {
        fmt.Printf("Error: %s\n", getErrorMessage(err))

        // Show help based on error type
        if errors.Is(err, ErrMismatchedParentheses) {
            fmt.Println("Example: (2 + 3) * 4")
        } else if errors.Is(err, ErrDivisionByZero) {
            fmt.Println("Example: 10 / 2")
        }
        return
    }

    fmt.Printf("= %f\n", result)
}
```

### Batch Processing

```go
func processExpressions(expressions []string) map[string]error {
    results := make(map[string]error)
    engine := calculator.NewEngine()

    for _, expr := range expressions {
        _, err := engine.Evaluate(expr)
        if err != nil {
            results[expr] = err
            // Reset engine for next calculation
            engine.Clear()
        } else {
            results[expr] = nil
        }
    }

    return results
}
```

### Interactive Calculator

```go
func interactiveCalculator() {
    engine := calculator.NewEngine()
    scanner := bufio.NewScanner(os.Stdin)

    for {
        fmt.Print("> ")
        if !scanner.Scan() {
            break
        }

        input := strings.TrimSpace(scanner.Text())
        if input == "exit" || input == "quit" {
            break
        }

        if input == "clear" || input == "c" {
            engine.Clear()
            fmt.Println("Calculator cleared")
            continue
        }

        result, err := engine.Evaluate(input)
        if err != nil {
            fmt.Printf("Error: %s\n", getErrorMessage(err))
            continue
        }

        fmt.Printf("= %f\n", result)
    }
}
```

### Error Logging

```go
func logCalculationError(expr string, err error) {
    logEntry := map[string]interface{}{
        "expression": expr,
        "error":      err.Error(),
        "error_type": getErrorType(err),
        "timestamp":  time.Now(),
    }

    // Log to file or monitoring system
    jsonData, _ := json.Marshal(logEntry)
    log.Printf("Calculation error: %s", string(jsonData))
}

func getErrorType(err error) string {
    switch {
    case errors.Is(err, ErrEmptyExpression):
        return "empty_expression"
    case errors.Is(err, ErrDivisionByZero):
        return "division_by_zero"
    case errors.Is(err, ErrMismatchedParentheses):
        return "mismatched_parentheses"
    // ... other error types
    default:
        return "unknown_error"
    }
}
```

## Testing Error Conditions

```go
func TestErrorHandling(t *testing.T) {
    engine := calculator.NewEngine()

    tests := []struct {
        name        string
        expression  string
        expectedErr error
    }{
        {
            name:        "empty expression",
            expression:  "",
            expectedErr: ErrEmptyExpression,
        },
        {
            name:        "division by zero",
            expression:  "5 / 0",
            expectedErr: ErrDivisionByZero,
        },
        {
            name:        "mismatched parentheses",
            expression:  "(2 + 3",
            expectedErr: ErrMismatchedParentheses,
        },
        {
            name:        "invalid expression",
            expression:  "2 + * 3",
            expectedErr: ErrInvalidExpression,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            _, err := engine.Evaluate(tt.expression)

            if err == nil {
                t.Errorf("Expected error %v, got nil", tt.expectedErr)
                return
            }

            if !errors.Is(err, tt.expectedErr) {
                t.Errorf("Expected error %v, got %v", tt.expectedErr, err)
            }
        })
    }
}
```

---

The error handling system provides comprehensive coverage of all possible error scenarios in the calculator, making it easy to write robust applications that handle errors gracefully.