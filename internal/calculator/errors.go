package calculator

import (
	"fmt"
	"math"
)

// Error types for calculator operations
type CalculatorError string

func (e CalculatorError) Error() string {
	return string(e)
}

const (
	ErrDivisionByZero      CalculatorError = "division by zero"
	ErrInvalidExpression   CalculatorError = "invalid expression"
	ErrOverflow            CalculatorError = "arithmetic overflow"
	ErrUnderflow           CalculatorError = "arithmetic underflow"
	ErrEmptyExpression     CalculatorError = "empty expression"
	ErrInvalidNumber       CalculatorError = "invalid number format"
	ErrInvalidOperator     CalculatorError = "invalid operator"
	ErrMismatchedParentheses CalculatorError = "mismatched parentheses"
)

// IsOverflow checks if a calculation would result in overflow
func IsOverflow(value float64) bool {
	return math.IsInf(value, 1) || math.IsInf(value, -1)
}

// IsUnderflow checks if a calculation would result in underflow
func IsUnderflow(value float64) bool {
	return math.Abs(value) < math.SmallestNonzeroFloat64
}

// ValidateNumber checks if a number is valid for calculations
func ValidateNumber(value float64) error {
	if math.IsNaN(value) {
		return ErrInvalidNumber
	}
	if IsOverflow(value) {
		return ErrOverflow
	}
	return nil
}