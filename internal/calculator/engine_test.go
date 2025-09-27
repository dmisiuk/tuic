package calculator

import (
	"math"
	"testing"
)

func TestNewEngine(t *testing.T) {
	engine := NewEngine()
	if engine == nil {
		t.Error("NewEngine() should not return nil")
	}
	if engine.GetValue() != 0 {
		t.Errorf("Expected initial value 0, got %f", engine.GetValue())
	}
}

func TestBasicArithmetic(t *testing.T) {
	engine := NewEngine()

	// Test addition
	result, err := engine.Add(5)
	if err != nil {
		t.Errorf("Add(5) returned error: %v", err)
	}
	if result != 5 {
		t.Errorf("Add(5) = %f, want 5", result)
	}

	// Test subtraction
	result, err = engine.Subtract(2)
	if err != nil {
		t.Errorf("Subtract(2) returned error: %v", err)
	}
	if result != 3 {
		t.Errorf("Subtract(2) = %f, want 3", result)
	}

	// Test multiplication
	result, err = engine.Multiply(4)
	if err != nil {
		t.Errorf("Multiply(4) returned error: %v", err)
	}
	if result != 12 {
		t.Errorf("Multiply(4) = %f, want 12", result)
	}

	// Test division
	result, err = engine.Divide(3)
	if err != nil {
		t.Errorf("Divide(3) returned error: %v", err)
	}
	if result != 4 {
		t.Errorf("Divide(3) = %f, want 4", result)
	}
}

func TestDivisionByZero(t *testing.T) {
	engine := NewEngine()
	engine.SetValue(10)

	_, err := engine.Divide(0)
	if err != ErrDivisionByZero {
		t.Errorf("Expected ErrDivisionByZero, got %v", err)
	}
}

func TestClearOperations(t *testing.T) {
	engine := NewEngine()
	engine.SetValue(100)

	// Test ClearEntry
	engine.InputNumber(5)
	engine.ClearEntry()
	if engine.GetEntryValue() != 0 {
		t.Errorf("ClearEntry() should reset entry value to 0")
	}

	// Test Clear
	engine.Clear()
	if engine.GetValue() != 0 {
		t.Errorf("Clear() should reset current value to 0")
	}
	if engine.GetEntryValue() != 0 {
		t.Errorf("Clear() should reset entry value to 0")
	}
	if engine.ShouldClear() {
		t.Errorf("Clear() should reset shouldClear to false")
	}
}

func TestInputNumber(t *testing.T) {
	engine := NewEngine()

	// Test digit input
	err := engine.InputNumber(1)
	if err != nil {
		t.Errorf("InputNumber(1) returned error: %v", err)
	}
	if engine.GetEntryValue() != 1 {
		t.Errorf("InputNumber(1) = %f, want 1", engine.GetEntryValue())
	}

	err = engine.InputNumber(2)
	if err != nil {
		t.Errorf("InputNumber(2) returned error: %v", err)
	}
	if engine.GetEntryValue() != 12 {
		t.Errorf("InputNumber(2) = %f, want 12", engine.GetEntryValue())
	}

	// Test invalid digit
	err = engine.InputNumber(-1)
	if err != ErrInvalidNumber {
		t.Errorf("InputNumber(-1) should return ErrInvalidNumber, got %v", err)
	}

	err = engine.InputNumber(10)
	if err != ErrInvalidNumber {
		t.Errorf("InputNumber(10) should return ErrInvalidNumber, got %v", err)
	}
}

func TestSetValue(t *testing.T) {
	engine := NewEngine()

	err := engine.SetValue(42.5)
	if err != nil {
		t.Errorf("SetValue(42.5) returned error: %v", err)
	}
	if engine.GetValue() != 42.5 {
		t.Errorf("SetValue(42.5) = %f, want 42.5", engine.GetValue())
	}
	if !engine.ShouldClear() {
		t.Errorf("SetValue() should set shouldClear to true")
	}
}

func TestSetValueWithInvalidValue(t *testing.T) {
	engine := NewEngine()

	// Test NaN
	err := engine.SetValue(math.NaN())
	if err != ErrInvalidNumber {
		t.Errorf("SetValue(NaN) should return ErrInvalidNumber, got %v", err)
	}

	// Test Infinity
	err = engine.SetValue(math.Inf(1))
	if err != ErrOverflow {
		t.Errorf("SetValue(Inf) should return ErrOverflow, got %v", err)
	}
}

func TestEvaluate(t *testing.T) {
	engine := NewEngine()

	tests := []struct {
		expression string
		want       float64
		wantErr    bool
	}{
		{"2+2", 4, false},
		{"5-3", 2, false},
		{"3*4", 12, false},
		{"10/2", 5, false},
		{"2.5*2", 5, false},
		{"1+2*3", 7, false}, // Operator precedence test
		{"(1+2)*3", 9, false}, // Parentheses test
		{"-5", -5, false}, // Unary minus test
		{"+5", 5, false}, // Unary plus test
		{"", 0, true}, // Empty expression
		{"abc", 0, true}, // Invalid expression
		{"2/0", 0, true}, // Division by zero
	}

	for _, tt := range tests {
		result, err := engine.Evaluate(tt.expression)
		if tt.wantErr {
			if err == nil {
				t.Errorf("Evaluate(%q) expected error, got nil", tt.expression)
			}
			continue
		}
		if err != nil {
			t.Errorf("Evaluate(%q) returned error: %v", tt.expression, err)
			continue
		}
		if math.Abs(result-tt.want) > 1e-10 {
			t.Errorf("Evaluate(%q) = %f, want %f", tt.expression, result, tt.want)
		}
	}
}

func TestOverflowDetection(t *testing.T) {
	engine := NewEngine()

	// Test with very large numbers that might cause overflow
	largeNum := math.MaxFloat64 / 2
	engine.SetValue(largeNum)

	_, err := engine.Multiply(3)
	if err != ErrOverflow {
		t.Errorf("Expected overflow error for large multiplication, got %v", err)
	}
}

func TestUnderflowDetection(t *testing.T) {
	engine := NewEngine()

	// Test with very small numbers that might cause underflow
	smallNum := math.SmallestNonzeroFloat64
	engine.SetValue(smallNum)

	result, err := engine.Divide(2)
	if err != nil {
		t.Errorf("Unexpected error for small division: %v", err)
	}
	if result == 0 {
		t.Errorf("Expected non-zero result for small division, got %f", result)
	}
}

func TestPerformOperation(t *testing.T) {
	engine := NewEngine()

	// Test operation with shouldClear = true
	engine.SetValue(10)
	engine.ShouldClear = true

	result, err := engine.PerformOperation("+")
	if err != nil {
		t.Errorf("PerformOperation(+) with shouldClear=true returned error: %v", err)
	}
	if result != 10 {
		t.Errorf("PerformOperation(+) with shouldClear=true = %f, want 10", result)
	}

	// Test normal operation
	engine.ShouldClear = false
	engine.InputNumber(5)

	result, err = engine.PerformOperation("+")
	if err != nil {
		t.Errorf("PerformOperation(+) returned error: %v", err)
	}
	if result != 15 {
		t.Errorf("PerformOperation(+) = %f, want 15", result)
	}
}

func TestFloatingPointPrecision(t *testing.T) {
	engine := NewEngine()

	// Test decimal arithmetic
	result, err := engine.Evaluate("0.1 + 0.2")
	if err != nil {
		t.Errorf("Evaluate('0.1 + 0.2') returned error: %v", err)
	}
	// Allow for floating point imprecision
	if math.Abs(result-0.3) > 1e-10 {
		t.Errorf("Evaluate('0.1 + 0.2') = %f, want ~0.3", result)
	}

	// Test more complex decimal operations
	result, err = engine.Evaluate("3.14159 * 2")
	if err != nil {
		t.Errorf("Evaluate('3.14159 * 2') returned error: %v", err)
	}
	if math.Abs(result-6.28318) > 1e-5 {
		t.Errorf("Evaluate('3.14159 * 2') = %f, want ~6.28318", result)
	}
}

func BenchmarkBasicOperations(b *testing.B) {
	engine := NewEngine()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		engine.Add(5)
		engine.Subtract(2)
		engine.Multiply(4)
		engine.Divide(2)
	}
}

func BenchmarkExpressionEvaluation(b *testing.B) {
	engine := NewEngine()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		engine.Evaluate("2+3*4-1/2")
	}
}