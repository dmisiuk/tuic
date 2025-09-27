package calculator

import (
	"math"
	"testing"
)

func TestNewParser(t *testing.T) {
	parser := NewParser()
	if parser == nil {
		t.Error("NewParser() should not return nil")
	}
}

func TestParseBasicExpressions(t *testing.T) {
	parser := NewParser()

	tests := []struct {
		expression string
		want       float64
		wantErr    bool
	}{
		// Basic arithmetic
		{"2+2", 4, false},
		{"5-3", 2, false},
		{"3*4", 12, false},
		{"10/2", 5, false},
		{"100/10", 10, false},

		// Decimal numbers
		{"2.5+1.5", 4, false},
		{"3.14*2", 6.28, false},
		{"10.5/2", 5.25, false},
		{"0.1+0.2", 0.3, false},

		// Operator precedence
		{"2+3*4", 14, false},     // 2 + (3*4) = 14
		{"2*3+4", 10, false},     // (2*3) + 4 = 10
		{"10-5*2", 0, false},     // 10 - (5*2) = 0
		{"20/5*2", 8, false},     // (20/5)*2 = 8

		// Multiple operations
		{"1+2+3+4", 10, false},
		{"10-5-2", 3, false},
		{"2*3*4", 24, false},
		{"100/10/2", 5, false},
		{"2+3*4-1", 13, false},
		{"10-2*3+4", 8, false},

		// Parentheses
		{"(2+3)*4", 20, false},
		{"2*(3+4)", 14, false},
		{"(10-5)*2", 10, false},
		{"(20/5)*2", 8, false},
		{"((2+3)*4)", 20, false},
		{"2*(3+(4*2))", 22, false},

		// Unary operators
		{"-5", -5, false},
		{"+5", 5, false},
		{"-2.5", -2.5, false},
		{"-(2+3)", -5, false},
		{"+(2*3)", 6, false},
		{"-3*4", -12, false},
		{"3*(-4)", -12, false},

		// Complex expressions
		{"2.5*(3+4.5)/2", 9.375, false},
		{"(10-3.5)*2+1", 14, false},
		{"5+3*2-4/2", 9, false},
		{"(5+3)*(2-4)/2", -8, false},

		// Edge cases
		{"0", 0, false},
		{"0.0", 0, false},
		{"123.456", 123.456, false},
		{"0.000001", 0.000001, false},
		{"999999.999", 999999.999, false},
	}

	for _, tt := range tests {
		result, err := parser.Parse(tt.expression)
		if tt.wantErr {
			if err == nil {
				t.Errorf("Parse(%q) expected error, got nil", tt.expression)
			}
			continue
		}
		if err != nil {
			t.Errorf("Parse(%q) returned error: %v", tt.expression, err)
			continue
		}
		// Allow for floating point imprecision
		if math.Abs(result-tt.want) > 1e-10 {
			t.Errorf("Parse(%q) = %f, want %f", tt.expression, result, tt.want)
		}
	}
}

func TestParseErrorCases(t *testing.T) {
	parser := NewParser()

	tests := []struct {
		expression string
		wantErr    bool
		errType    error
	}{
		{"", true, ErrEmptyExpression},
		{"abc", true, ErrInvalidExpression},
		{"2+", true, ErrInvalidExpression},
		{"2++3", true, ErrInvalidExpression},
		{"2*/3", true, ErrInvalidExpression},
		{"2..3", true, ErrInvalidNumber},
		{"2.3.4", true, ErrInvalidNumber},
		{"(2+3", true, ErrMismatchedParentheses},
		{"2+3)", true, ErrInvalidExpression},
		{"()", true, ErrInvalidExpression},
		{"(+)", true, ErrInvalidExpression},
		{"(-)", true, ErrInvalidExpression},
		{"2+*3", true, ErrInvalidExpression},
		{"*2", true, ErrInvalidExpression},
		{"2+", true, ErrInvalidExpression},
		{".5", true, ErrInvalidNumber},
		{"2.", true, ErrInvalidNumber},
		{"+", true, ErrInvalidExpression},
		{"-", true, ErrInvalidExpression},
		{"*", true, ErrInvalidExpression},
		{"/", true, ErrInvalidExpression},
		{"(2+3)*4)", true, ErrMismatchedParentheses},
		{"2+(3*4", true, ErrMismatchedParentheses},
	}

	for _, tt := range tests {
		_, err := parser.Parse(tt.expression)
		if !tt.wantErr {
			if err != nil {
				t.Errorf("Parse(%q) expected no error, got %v", tt.expression, err)
			}
			continue
		}
		if err == nil {
			t.Errorf("Parse(%q) expected error, got nil", tt.expression)
			continue
		}
		if tt.errType != nil && err != tt.errType {
			t.Errorf("Parse(%q) expected error type %v, got %v", tt.expression, tt.errType, err)
		}
	}
}

func TestParseDivisionByZero(t *testing.T) {
	parser := NewParser()

	_, err := parser.Parse("5/0")
	if err != ErrDivisionByZero {
		t.Errorf("Parse('5/0') expected ErrDivisionByZero, got %v", err)
	}

	// Test division by zero in complex expressions
	_, err = parser.Parse("2+3/0")
	if err != ErrDivisionByZero {
		t.Errorf("Parse('2+3/0') expected ErrDivisionByZero, got %v", err)
	}

	_, err = parser.Parse("(2+3)/(4-4)")
	if err != ErrDivisionByZero {
		t.Errorf("Parse('(2+3)/(4-4)') expected ErrDivisionByZero, got %v", err)
	}
}

func TestParseWithSpaces(t *testing.T) {
	parser := NewParser()

	// Test that spaces are properly removed
	result, err := parser.Parse("2 + 3 * 4")
	if err != nil {
		t.Errorf("Parse('2 + 3 * 4') returned error: %v", err)
	}
	if result != 14 {
		t.Errorf("Parse('2 + 3 * 4') = %f, want 14", result)
	}

	result, err = parser.Parse(" ( 2 + 3 ) * 4 ")
	if err != nil {
		t.Errorf("Parse(' ( 2 + 3 ) * 4 ') returned error: %v", err)
	}
	if result != 20 {
		t.Errorf("Parse(' ( 2 + 3 ) * 4 ') = %f, want 20", result)
	}
}

func TestParseLargeNumbers(t *testing.T) {
	parser := NewParser()

	// Test large numbers that might cause overflow
	largeNum := math.MaxFloat64 / 10
	result, err := parser.Parse("1e300")
	if err != nil {
		t.Errorf("Parse('1e300') returned error: %v", err)
	}
	if !math.IsInf(result, 1) {
		t.Errorf("Parse('1e300') should return infinity for overflow")
	}

	// Test underflow
	result, err = parser.Parse("1e-350")
	if err != nil {
		t.Errorf("Parse('1e-350') returned error: %v", err)
	}
	if result != 0 {
		t.Errorf("Parse('1e-350') should return 0 for underflow")
	}
}

func TestEvaluateSimple(t *testing.T) {
	tests := []struct {
		a         float64
		b         float64
		operator  string
		want      float64
		wantErr   bool
		errType   error
	}{
		{2, 3, "+", 5, false, nil},
		{5, 3, "-", 2, false, nil},
		{3, 4, "*", 12, false, nil},
		{10, 2, "/", 5, false, nil},
		{2.5, 2, "*", 5, false, nil},
		{10, 4, "/", 2.5, false, nil},
		{5, 0, "/", 0, true, ErrDivisionByZero},
		{5, 3, "%", 0, true, ErrInvalidOperator},
		{math.Inf(1), 1, "+", 0, true, ErrOverflow},
		{math.NaN(), 1, "+", 0, true, ErrInvalidNumber},
	}

	for _, tt := range tests {
		result, err := EvaluateSimple(tt.a, tt.b, tt.operator)
		if tt.wantErr {
			if err == nil {
				t.Errorf("EvaluateSimple(%f, %f, %q) expected error, got nil", tt.a, tt.b, tt.operator)
			}
			if tt.errType != nil && err != tt.errType {
				t.Errorf("EvaluateSimple(%f, %f, %q) expected error type %v, got %v", tt.a, tt.b, tt.operator, tt.errType, err)
			}
			continue
		}
		if err != nil {
			t.Errorf("EvaluateSimple(%f, %f, %q) returned error: %v", tt.a, tt.b, tt.operator, err)
			continue
		}
		if math.Abs(result-tt.want) > 1e-10 {
			t.Errorf("EvaluateSimple(%f, %f, %q) = %f, want %f", tt.a, tt.b, tt.operator, result, tt.want)
		}
	}
}

func TestParserState(t *testing.T) {
	parser := NewParser()

	// Test that parser resets properly for each parse
	result, err := parser.Parse("2+3")
	if err != nil {
		t.Errorf("First parse returned error: %v", err)
	}
	if result != 5 {
		t.Errorf("First parse result = %f, want 5", result)
	}

	// Second parse should work normally
	result, err = parser.Parse("5*2")
	if err != nil {
		t.Errorf("Second parse returned error: %v", err)
	}
	if result != 10 {
		t.Errorf("Second parse result = %f, want 10", result)
	}
}

func TestParserEdgeCases(t *testing.T) {
	parser := NewParser()

	// Test very long expressions
	longExpr := "1+2+3+4+5+6+7+8+9+10"
	result, err := parser.Parse(longExpr)
	if err != nil {
		t.Errorf("Parse(%q) returned error: %v", longExpr, err)
	}
	if result != 55 {
		t.Errorf("Parse(%q) = %f, want 55", longExpr, result)
	}

	// Test nested parentheses
	nestedExpr := "((1+2)*((3+4)*(5+6)))"
	result, err = parser.Parse(nestedExpr)
	if err != nil {
		t.Errorf("Parse(%q) returned error: %v", nestedExpr, err)
	}
	if result != 231 {
		t.Errorf("Parse(%q) = %f, want 231", nestedExpr, result)
	}
}

func BenchmarkParser(b *testing.B) {
	parser := NewParser()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		parser.Parse("2+3*4-1/2")
	}
}

func BenchmarkComplexParser(b *testing.B) {
	parser := NewParser()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		parser.Parse("(2.5*(3+4.5))/2+1.5*3")
	}
}

func BenchmarkParserWithParentheses(b *testing.B) {
	parser := NewParser()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		parser.Parse("((1+2)*(3+4))/(5+6)")
	}
}