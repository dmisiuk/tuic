package calculator

import (
	"fmt"
	"math/rand"
	"testing"
)

func BenchmarkCalculator_Evaluate_Simple(b *testing.B) {
	calc := NewCalculator()
	expr := "2 + 3 * 4"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := calc.Evaluate(expr)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkCalculator_Evaluate_Complex(b *testing.B) {
	calc := NewCalculator()
	expr := "((100 + 50) * 2 - 30) / 4 + 10 ^ 2"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := calc.Evaluate(expr)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkCalculator_Evaluate_Variables(b *testing.B) {
	calc := NewCalculator()
	calc.SetVariable("x", 10)
	calc.SetVariable("y", 20)
	calc.SetVariable("z", 5)

	expr := "x * y + z / 2"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := calc.Evaluate(expr)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkParser_Parse(b *testing.B) {
	parser := NewParser()
	expr := "((100 + 50) * 2 - 30) / 4 + 10 ^ 2"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := parser.Parse(expr)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkCalculator_ExpressionGenerator(b *testing.B) {
	calc := NewCalculator()

	// Generate random expressions
	ops := []string{"+", "-", "*", "/", "^"}
	numbers := []int{1, 2, 3, 4, 5, 10, 20, 50, 100}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Generate a random expression
		expr := fmt.Sprintf("%d %s %d %s %d",
			numbers[rand.Intn(len(numbers))],
			ops[rand.Intn(len(ops))],
			numbers[rand.Intn(len(numbers))],
			ops[rand.Intn(len(ops))],
			numbers[rand.Intn(len(numbers))],
		)

		_, err := calc.Evaluate(expr)
		if err != nil {
			// Some expressions might be invalid, that's okay for benchmarking
			continue
		}
	}
}

func BenchmarkCalculator_VariableOperations(b *testing.B) {
	calc := NewCalculator()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		varName := fmt.Sprintf("var%d", i%100)
		calc.SetVariable(varName, float64(i))

		if i%100 == 0 {
			calc.GetVariable(varName)
		}
	}
}

func BenchmarkCalculator_MemoryUsage(b *testing.B) {
	calc := NewCalculator()

	// Set many variables
	for i := 0; i < 1000; i++ {
		calc.SetVariable(fmt.Sprintf("x%d", i), float64(i))
	}

	expr := "x1 + x2 + x3 + x4 + x5"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := calc.Evaluate(expr)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkCalculator_LargeExpression(b *testing.B) {
	calc := NewCalculator()

	// Build a very large expression
	var expr string
	for i := 0; i < 100; i++ {
		if i > 0 {
			expr += " + "
		}
		expr += fmt.Sprintf("%d", i+1)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := calc.Evaluate(expr)
		if err != nil {
			b.Fatal(err)
		}
	}
}