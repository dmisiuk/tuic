# Engine API Reference

The `Engine` struct provides the core calculator functionality with state management and arithmetic operations.

## Overview

The engine maintains calculator state including current value, entry value, and clear status. It handles basic arithmetic operations, expression evaluation, and number input.

## Table of Contents

- [Engine Struct](#engine-struct)
- [Constructor](#constructor)
- [Core Methods](#core-methods)
- [Arithmetic Operations](#arithmetic-operations)
- [State Management](#state-management)
- [Number Input](#number-input)
- [Helper Methods](#helper-methods)

## Engine Struct

```go
type Engine struct {
    currentValue float64  // Current display value
    entryValue   float64  // Current entry being typed
    shouldClear  bool     // Whether to clear display before next input
}
```

### Fields

- `currentValue`: The result of the last operation, displayed on screen
- `entryValue`: The number being currently input by the user
- `shouldClear`: Flag indicating if display should clear before next input

## Constructor

### NewEngine

```go
func NewEngine() *Engine
```

Creates a new calculator engine with initial state set to zero.

**Returns:**
- `*Engine`: Pointer to newly created engine instance

**Example:**
```go
engine := calculator.NewEngine()
```

## Core Methods

### Evaluate

```go
func (e *Engine) Evaluate(expression string) (float64, error)
```

Evaluates a mathematical expression and returns the result.

**Parameters:**
- `expression`: Mathematical expression string (e.g., "2 + 3 * 4")

**Returns:**
- `float64`: Result of evaluation
- `error`: Error if evaluation fails

**Supported Operations:**
- Addition (+), Subtraction (-)
- Multiplication (*), Division (/)
- Parentheses ( )
- Decimal numbers
- Operator precedence

**Example:**
```go
result, err := engine.Evaluate("(10 + 5) * (8 - 3) / 2")
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Result: %.2f\n", result) // Output: 37.50
```

**Error Cases:**
- Empty expression
- Invalid expression syntax
- Division by zero
- Mismatched parentheses
- Number overflow/underflow

### SetValue

```go
func (e *Engine) SetValue(value float64) error
```

Sets the current value directly.

**Parameters:**
- `value`: Value to set

**Returns:**
- `error`: Error if value is invalid

**Example:**
```go
err := engine.SetValue(3.14159)
if err != nil {
    log.Fatal(err)
}
```

## Arithmetic Operations

### Add

```go
func (e *Engine) Add(value float64) (float64, error)
```

Adds a number to the current value.

**Parameters:**
- `value`: Number to add

**Returns:**
- `float64`: Result of addition
- `error`: Error if operation fails

**Example:**
```go
engine.SetValue(10)
result, err := engine.Add(5)
// result = 15
```

### Subtract

```go
func (e *Engine) Subtract(value float64) (float64, error)
```

Subtracts a number from the current value.

**Parameters:**
- `value`: Number to subtract

**Returns:**
- `float64`: Result of subtraction
- `error`: Error if operation fails

**Example:**
```go
engine.SetValue(10)
result, err := engine.Subtract(3)
// result = 7
```

### Multiply

```go
func (e *Engine) Multiply(value float64) (float64, error)
```

Multiplies the current value by a number.

**Parameters:**
- `value`: Number to multiply by

**Returns:**
- `float64`: Result of multiplication
- `error`: Error if operation fails

**Example:**
```go
engine.SetValue(6)
result, err := engine.Multiply(7)
// result = 42
```

### Divide

```go
func (e *Engine) Divide(value float64) (float64, error)
```

Divides the current value by a number.

**Parameters:**
- `value`: Number to divide by

**Returns:**
- `float64`: Result of division
- `error`: Error if operation fails

**Example:**
```go
engine.SetValue(20)
result, err := engine.Divide(4)
// result = 5
```

**Special Cases:**
- Returns `ErrDivisionByZero` if `value` is 0
- Validates result for overflow/underflow

## State Management

### Clear

```go
func (e *Engine) Clear()
```

Clears all calculator values (C functionality).

**Example:**
```go
engine.SetValue(123.45)
engine.Clear()
// currentValue = 0, entryValue = 0, shouldClear = false
```

### ClearEntry

```go
func (e *Engine) ClearEntry()
```

Clears the current entry value (CE functionality).

**Example:**
```go
engine.InputNumber(1)
engine.InputNumber(2)
engine.InputNumber(3)
// entryValue = 123

engine.ClearEntry()
// entryValue = 0
```

## Number Input

### InputNumber

```go
func (e *Engine) InputNumber(digit int) error
```

Inputs a digit into the calculator with proper state management.

**Parameters:**
- `digit`: Digit to input (0-9)

**Returns:**
- `error`: Error if digit is invalid

**Behavior:**
- Clears entry if `shouldClear` is true
- Appends digit to current entry value
- Handles multi-digit input correctly

**Example:**
```go
// Input 123
engine.InputNumber(1)
engine.InputNumber(2)
engine.InputNumber(3)

entryValue := engine.GetEntryValue()
// entryValue = 123
```

### PerformOperation

```go
func (e *Engine) PerformOperation(op string) (float64, error)
```

Performs an arithmetic operation between current and entry values.

**Parameters:**
- `op`: Operation to perform (+, -, *, /)

**Returns:**
- `float64`: Result of operation
- `error`: Error if operation fails

**Behavior:**
- If `shouldClear` is true, returns current value
- Otherwise performs operation between current and entry values
- Updates entry value with result

**Example:**
```go
engine.SetValue(10)
engine.InputNumber(1)
engine.InputNumber(5)
// entryValue = 15

result, err := engine.PerformOperation("+")
// result = 25 (10 + 15)
```

## Helper Methods

### GetValue

```go
func (e *Engine) GetValue() float64
```

Returns the current display value.

**Returns:**
- `float64`: Current value

**Example:**
```go
engine.SetValue(42)
value := engine.GetValue()
// value = 42
```

### GetEntryValue

```go
func (e *Engine) GetEntryValue() float64
```

Returns the current entry value.

**Returns:**
- `float64`: Entry value

**Example:**
```go
engine.InputNumber(1)
engine.InputNumber(2)
entry := engine.GetEntryValue()
// entry = 12
```

### ShouldClear

```go
func (e *Engine) ShouldClear() bool
```

Returns whether the display should be cleared before next input.

**Returns:**
- `bool`: True if display should clear

**Example:**
```go
if engine.ShouldClear() {
    engine.ClearEntry()
}
```

## Error Handling

The engine methods return specific error types:

- `ErrEmptyExpression`: Empty expression provided
- `ErrInvalidExpression`: Invalid expression syntax
- `ErrInvalidNumber`: Invalid number format
- `ErrDivisionByZero`: Division by zero attempted
- `ErrInvalidOperator`: Invalid operator
- `ErrNumberOutOfRange`: Number out of valid range

### Error Handling Example

```go
result, err := engine.Evaluate("5 / 0")
if err != nil {
    switch {
    case errors.Is(err, ErrDivisionByZero):
        fmt.Println("Cannot divide by zero")
    case errors.Is(err, ErrInvalidExpression):
        fmt.Println("Invalid expression syntax")
    default:
        fmt.Printf("Error: %v\n", err)
    }
    return
}
fmt.Printf("Result: %v\n", result)
```

## Usage Patterns

### Basic Calculator Flow

```go
engine := calculator.NewEngine()

// User enters: 1 2 + 3 4 =
engine.InputNumber(1)
engine.InputNumber(2)          // Entry: 12

result, err := engine.PerformOperation("+")
// result = 12, currentValue = 12

engine.InputNumber(3)
engine.InputNumber(4)          // Entry: 34

result, err = engine.PerformOperation("=")
// result = 46 (12 + 34)
```

### Expression Evaluation

```go
engine := calculator.NewEngine()

// Evaluate complex expression
result, err := engine.Evaluate("2 * (3 + 4) / 2")
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Result: %v\n", result) // Output: 7
```

### Error Recovery

```go
engine := calculator.NewEngine()

// Try invalid operation
result, err := engine.Evaluate("5 / 0")
if err != nil {
    fmt.Printf("Error: %v\n", err)
    engine.Clear()  // Reset calculator state
}

// Continue with valid operation
result, err = engine.Evaluate("10 + 5")
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Result: %v\n", result)
```

## Performance Considerations

- The engine uses float64 for calculations
- All operations include overflow/underflow validation
- State management is optimized for typical calculator usage patterns
- Expression parsing is handled by a separate parser component