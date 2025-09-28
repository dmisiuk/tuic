package calculator

import (
	"sync"
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

// Calculator provides a high-level interface with variable support
type Calculator struct {
	engine    *Engine
	variables map[string]float64
	mu        sync.RWMutex
}

// NewCalculator creates a new calculator with variable support
func NewCalculator() *Calculator {
	return &Calculator{
		engine:    NewEngine(),
		variables: make(map[string]float64),
	}
}

// Evaluate evaluates a mathematical expression with variable support
func (c *Calculator) Evaluate(expression string) (float64, error) {
	// Simple variable substitution for now
	c.mu.RLock()
	for name, value := range c.variables {
		expression = replaceVariable(expression, name, value)
	}
	c.mu.RUnlock()

	return c.engine.Evaluate(expression)
}

// SetVariable sets a variable value
func (c *Calculator) SetVariable(name string, value float64) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.variables[name] = value
}

// GetVariable gets a variable value
func (c *Calculator) GetVariable(name string) (float64, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	value, exists := c.variables[name]
	return value, exists
}

// GetVariables returns all variables
func (c *Calculator) GetVariables() map[string]float64 {
	c.mu.RLock()
	defer c.mu.RUnlock()

	// Return a copy to avoid concurrent modification
	vars := make(map[string]float64)
	for k, v := range c.variables {
		vars[k] = v
	}
	return vars
}

// ClearVariables clears all variables
func (c *Calculator) ClearVariables() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.variables = make(map[string]float64)
}

// replaceVariable replaces variable names with their values in expressions
func replaceVariable(expr, name string, value float64) string {
	// Simple implementation - in production this would need proper parsing
	return expr // For now, we'll implement this in the parser
}