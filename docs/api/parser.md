# Parser API Reference

The `Parser` struct provides mathematical expression parsing and evaluation with support for operator precedence and parentheses.

## Overview

The parser handles expression parsing, operator precedence, and nested parentheses. It converts mathematical expressions into computed results using a recursive descent parsing approach.

## Table of Contents

- [Parser Struct](#parser-struct)
- [Constructor](#constructor)
- [Core Methods](#core-methods)
- [Parsing Methods](#parsing-methods)
- [Helper Methods](#helper-methods)
- [Error Handling](#error-handling)
- [Usage Examples](#usage-examples)

## Parser Struct

```go
type Parser struct {
    expression string  // The expression being parsed
    position   int     // Current position in the expression
}
```

### Fields

- `expression`: The mathematical expression string being parsed
- `position`: Current parsing position within the expression

## Constructor

### NewParser

```go
func NewParser() *Parser
```

Creates a new parser instance.

**Returns:**
- `*Parser`: Pointer to newly created parser instance

**Example:**
```go
parser := calculator.NewParser()
```

## Core Methods

### Parse

```go
func (p *Parser) Parse(expression string) (float64, error)
```

Parses and evaluates a mathematical expression.

**Parameters:**
- `expression`: Mathematical expression string (e.g., "2 + 3 * 4")

**Returns:**
- `float64`: Result of evaluation
- `error`: Error if parsing fails

**Supported Operations:**
- Addition (+), Subtraction (-)
- Multiplication (*), Division (/)
- Parentheses ( )
- Decimal numbers
- Operator precedence (multiplication/division before addition/subtraction)

**Example:**
```go
parser := calculator.NewParser()
result, err := parser.Parse("2 + 3 * 4")
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Result: %f\n", result) // Output: 14.000000
```

**Error Cases:**
- Empty expression
- Invalid expression syntax
- Division by zero
- Mismatched parentheses
- Invalid number format

## Parsing Methods

### parseExpression

```go
func (p *Parser) parseExpression() (float64, error)
```

Handles addition and subtraction operations (lowest precedence).

**Returns:**
- `float64`: Result of expression evaluation
- `error`: Error if parsing fails

**Behavior:**
- Parses terms separated by + or - operators
- Handles left-associative operations
- Validates results for overflow/underflow

**Example:**
```go
// Internal method used by Parse()
// Handles: 2 + 3 - 1 + 4
```

### parseTerm

```go
func (p *Parser) parseTerm() (float64, error)
```

Handles multiplication and division operations (higher precedence).

**Returns:**
- `float64`: Result of term evaluation
- `error`: Error if parsing fails

**Behavior:**
- Parses factors separated by * or / operators
- Handles left-associative operations
- Validates division by zero
- Validates results for overflow/underflow

**Example:**
```go
// Internal method used by Parse()
// Handles: 2 * 3 / 4 * 5
```

### parseFactor

```go
func (p *Parser) parseFactor() (float64, error)
```

Handles numbers, parentheses, and unary operators.

**Returns:**
- `float64`: Result of factor evaluation
- `error`: Error if parsing fails

**Behavior:**
- Handles unary plus and minus operators
- Parses parenthesized expressions
- Delegates number parsing to parseNumber()

**Example:**
```go
// Internal method used by Parse()
// Handles: -5, +3, (2 + 3), 42
```

### parseNumber

```go
func (p *Parser) parseNumber() (float64, error)
```

Parses numeric literals from the expression.

**Returns:**
- `float64`: Parsed number value
- `error`: Error if parsing fails

**Behavior:**
- Parses integer part
- Optionally parses decimal part
- Validates number format
- Converts string to float64

**Example:**
```go
// Parses: 123, 3.14, 0.5, 42
```

## Helper Methods

### peek

```go
func (p *Parser) peek() byte
```

Returns the current character without consuming it.

**Returns:**
- `byte`: Current character or 0 if at end of expression

**Example:**
```go
// If expression is "123+456" and position is 3
// peek() returns '+'
```

### consume

```go
func (p *Parser) consume()
```

Consumes the current character and advances position.

**Behavior:**
- Increments position if not at end
- No effect if already at end

**Example:**
```go
// If expression is "123+456" and position is 3
// consume() advances position to 4
```

## Utility Functions

### EvaluateSimple

```go
func EvaluateSimple(a, b float64, operator string) (float64, error)
```

Evaluates a simple arithmetic operation between two numbers.

**Parameters:**
- `a`: First operand
- `b`: Second operand
- `operator`: Operation to perform (+, -, *, /)

**Returns:**
- `float64`: Result of operation
- `error`: Error if operation fails

**Example:**
```go
result, err := calculator.EvaluateSimple(10, 5, "+")
// result = 15

result, err = calculator.EvaluateSimple(10, 2, "*")
// result = 20
```

## Error Handling

The parser returns specific error types:

- `ErrEmptyExpression`: Empty expression provided
- `ErrInvalidExpression`: Invalid expression syntax
- `ErrInvalidNumber`: Invalid number format
- `ErrDivisionByZero`: Division by zero attempted
- `ErrMismatchedParentheses`: Mismatched parentheses
- `ErrNumberOutOfRange`: Number out of valid range

### Error Handling Example

```go
parser := calculator.NewParser()

result, err := parser.Parse("(2 + 3")
if err != nil {
    switch {
    case errors.Is(err, ErrMismatchedParentheses):
        fmt.Println("Mismatched parentheses")
    case errors.Is(err, ErrInvalidExpression):
        fmt.Println("Invalid expression syntax")
    case errors.Is(err, ErrDivisionByZero):
        fmt.Println("Division by zero")
    default:
        fmt.Printf("Error: %v\n", err)
    }
    return
}
fmt.Printf("Result: %v\n", result)
```

## Usage Examples

### Basic Arithmetic

```go
parser := calculator.NewParser()

// Simple addition
result, _ := parser.Parse("2 + 3")           // 5

// Subtraction
result, _ = parser.Parse("10 - 4")            // 6

// Multiplication
result, _ = parser.Parse("7 * 8")             // 56

// Division
result, _ = parser.Parse("100 / 4")           // 25
```

### Operator Precedence

```go
parser := calculator.NewParser()

// Multiplication before addition
result, _ := parser.Parse("2 + 3 * 4")        // 14 (not 20)

// Same precedence, left to right
result, _ = parser.Parse("10 - 3 + 2")        // 9 (not 5)

// Division before addition
result, _ = parser.Parse("6 + 12 / 3")        // 10 (not 6)
```

### Parentheses

```go
parser := calculator.NewParser()

// Parentheses override precedence
result, _ := parser.Parse("(2 + 3) * 4")      // 20
result, _ = parser.Parse("2 + (3 * 4)")        // 14

// Nested parentheses
result, _ = parser.Parse("((1 + 2) * 3) + 4")  // 13
result, _ = parser.Parse("(10 + 5) * (8 - 3)") // 75
```

### Decimal Numbers

```go
parser := calculator.NewParser()

// Decimal calculations
result, _ := parser.Parse("3.14 * 2")          // 6.28
result, _ = parser.Parse("10.5 / 2.5")        // 4.2
result, _ = parser.Parse("0.5 + 0.25")         // 0.75
```

### Complex Expressions

```go
parser := calculator.NewParser()

// Complex expressions
result, _ := parser.Parse("2 * (3 + 4) / 2 + 1")  // 8
result, _ = parser.Parse("(10 + 5) * (8 - 3) / (2 + 1)")  // 25
result, _ = parser.Parse("2.5 * (4 + 6) - 1.5 / 3")  // 24.5
```

### Error Cases

```go
parser := calculator.NewParser()

// Division by zero
result, err := parser.Parse("5 / 0")
// err = ErrDivisionByZero

// Mismatched parentheses
result, err = parser.Parse("(2 + 3")
// err = ErrMismatchedParentheses

// Invalid expression
result, err = parser.Parse("2 + * 3")
// err = ErrInvalidExpression

// Empty expression
result, err = parser.Parse("")
// err = ErrEmptyExpression
```

## Performance Considerations

- The parser uses recursive descent for efficient parsing
- Position tracking ensures linear time complexity
- Minimal memory allocation during parsing
- Error detection happens early in the process

## Best Practices

### Input Validation

```go
func isValidExpression(expr string) bool {
    if len(expr) == 0 {
        return false
    }

    // Remove spaces for validation
    expr = strings.ReplaceAll(expr, " ", "")

    // Basic character validation
    for _, char := range expr {
        if !strings.Contains("0123456789+-*/.()", string(char)) {
            return false
        }
    }

    return true
}
```

### Safe Parsing

```go
func safeEvaluate(expr string) (float64, error) {
    parser := calculator.NewParser()

    result, err := parser.Parse(expr)
    if err != nil {
        return 0, fmt.Errorf("failed to evaluate expression '%s': %w", expr, err)
    }

    // Validate result
    if err := calculator.ValidateNumber(result); err != nil {
        return 0, fmt.Errorf("result validation failed: %w", err)
    }

    return result, nil
}
```

### Integration with Engine

```go
func calculateExpression(engine *calculator.Engine, expr string) (float64, error) {
    // Use engine's evaluate method which uses parser internally
    result, err := engine.Evaluate(expr)
    if err != nil {
        return 0, err
    }

    return result, nil
}
```

## Testing the Parser

```go
func TestParser(t *testing.T) {
    tests := []struct {
        name      string
        expression string
        want      float64
        wantErr   bool
    }{
        {
            name:       "simple addition",
            expression: "2 + 3",
            want:       5,
            wantErr:    false,
        },
        {
            name:       "operator precedence",
            expression: "2 + 3 * 4",
            want:       14,
            wantErr:    false,
        },
        {
            name:       "parentheses",
            expression: "(2 + 3) * 4",
            want:       20,
            wantErr:    false,
        },
        {
            name:       "division by zero",
            expression: "5 / 0",
            wantErr:    true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            parser := NewParser()
            got, err := parser.Parse(tt.expression)

            if tt.wantErr {
                if err == nil {
                    t.Errorf("Parser.Parse() expected error, got nil")
                    return
                }
                return
            }

            if err != nil {
                t.Errorf("Parser.Parse() unexpected error = %v", err)
                return
            }

            if got != tt.want {
                t.Errorf("Parser.Parse() = %v, want %v", got, tt.want)
            }
        })
    }
}
```

---

The parser provides a robust foundation for mathematical expression evaluation with proper error handling and operator precedence support.