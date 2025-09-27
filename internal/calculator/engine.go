package calculator

import (
	"math"
)

// Engine represents the calculator engine state
type Engine struct {
	currentValue float64
	entryValue   float64
	shouldClear  bool
}

// NewEngine creates a new calculator engine
func NewEngine() *Engine {
	return &Engine{
		currentValue: 0,
		entryValue:   0,
		shouldClear:  false,
	}
}

// Evaluate evaluates a mathematical expression and returns the result
func (e *Engine) Evaluate(expression string) (float64, error) {
	if expression == "" {
		return 0, ErrEmptyExpression
	}

	parser := NewParser()
	result, err := parser.Parse(expression)
	if err != nil {
		return 0, err
	}

	// Validate the result
	if err := ValidateNumber(result); err != nil {
		return 0, err
	}

	e.currentValue = result
	e.shouldClear = true
	return result, nil
}

// Clear clears all values (C functionality)
func (e *Engine) Clear() {
	e.currentValue = 0
	e.entryValue = 0
	e.shouldClear = false
}

// ClearEntry clears the current entry (CE functionality)
func (e *Engine) ClearEntry() {
	e.entryValue = 0
	e.shouldClear = false
}

// Add adds a number to the current value
func (e *Engine) Add(value float64) (float64, error) {
	if err := ValidateNumber(value); err != nil {
		return 0, err
	}

	result := e.currentValue + value
	if err := ValidateNumber(result); err != nil {
		return 0, err
	}

	e.currentValue = result
	e.shouldClear = true
	return result, nil
}

// Subtract subtracts a number from the current value
func (e *Engine) Subtract(value float64) (float64, error) {
	if err := ValidateNumber(value); err != nil {
		return 0, err
	}

	result := e.currentValue - value
	if err := ValidateNumber(result); err != nil {
		return 0, err
	}

	e.currentValue = result
	e.shouldClear = true
	return result, nil
}

// Multiply multiplies the current value by a number
func (e *Engine) Multiply(value float64) (float64, error) {
	if err := ValidateNumber(value); err != nil {
		return 0, err
	}

	result := e.currentValue * value
	if err := ValidateNumber(result); err != nil {
		return 0, err
	}

	e.currentValue = result
	e.shouldClear = true
	return result, nil
}

// Divide divides the current value by a number
func (e *Engine) Divide(value float64) (float64, error) {
	if err := ValidateNumber(value); err != nil {
		return 0, err
	}

	if value == 0 {
		return 0, ErrDivisionByZero
	}

	result := e.currentValue / value
	if err := ValidateNumber(result); err != nil {
		return 0, err
	}

	e.currentValue = result
	e.shouldClear = true
	return result, nil
}

// SetValue sets the current value directly
func (e *Engine) SetValue(value float64) error {
	if err := ValidateNumber(value); err != nil {
		return err
	}

	e.currentValue = value
	e.shouldClear = true
	return nil
}

// GetValue returns the current value
func (e *Engine) GetValue() float64 {
	return e.currentValue
}

// GetEntryValue returns the entry value
func (e *Engine) GetEntryValue() float64 {
	return e.entryValue
}

// ShouldClear returns whether the display should be cleared before next input
func (e *Engine) ShouldClear() bool {
	return e.shouldClear
}

// InputNumber handles number input with proper state management
func (e *Engine) InputNumber(digit int) error {
	if digit < 0 || digit > 9 {
		return ErrInvalidNumber
	}

	if e.shouldClear {
		e.entryValue = 0
		e.shouldClear = false
	}

	e.entryValue = e.entryValue*10 + float64(digit)
	return nil
}

// PerformOperation performs an arithmetic operation between current and entry values
func (e *Engine) PerformOperation(op string) (float64, error) {
	if e.shouldClear {
		e.shouldClear = false
		return e.currentValue, nil
	}

	var result float64
	var err error

	switch op {
	case "+":
		result, err = e.Add(e.entryValue)
	case "-":
		result, err = e.Subtract(e.entryValue)
	case "*":
		result, err = e.Multiply(e.entryValue)
	case "/":
		result, err = e.Divide(e.entryValue)
	default:
		return 0, ErrInvalidOperator
	}

	if err != nil {
		return 0, err
	}

	e.entryValue = result
	return result, nil
}