package calculator

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

// Parser handles expression parsing and evaluation
type Parser struct {
	expression string
	position   int
}

// NewParser creates a new parser instance
func NewParser() *Parser {
	return &Parser{}
}

// Parse parses and evaluates a mathematical expression
func (p *Parser) Parse(expression string) (float64, error) {
	p.expression = strings.ReplaceAll(expression, " ", "")
	p.position = 0

	if len(p.expression) == 0 {
		return 0, ErrEmptyExpression
	}

	return p.parseExpression()
}

// parseExpression handles addition and subtraction (lowest precedence)
func (p *Parser) parseExpression() (float64, error) {
	left, err := p.parseTerm()
	if err != nil {
		return 0, err
	}

	for {
		op := p.peek()
		if op != '+' && op != '-' {
			break
		}

		p.consume() // consume the operator

		right, err := p.parseTerm()
		if err != nil {
			return 0, err
		}

		switch op {
		case '+':
			left += right
		case '-':
			left -= right
		}

		// Check for overflow/underflow
		if err := ValidateNumber(left); err != nil {
			return 0, err
		}
	}

	return left, nil
}

// parseTerm handles multiplication and division (higher precedence)
func (p *Parser) parseTerm() (float64, error) {
	left, err := p.parseFactor()
	if err != nil {
		return 0, err
	}

	for {
		op := p.peek()
		if op != '*' && op != '/' {
			break
		}

		p.consume() // consume the operator

		right, err := p.parseFactor()
		if err != nil {
			return 0, err
		}

		switch op {
		case '*':
			left *= right
		case '/':
			if right == 0 {
				return 0, ErrDivisionByZero
			}
			left /= right
		}

		// Check for overflow/underflow
		if err := ValidateNumber(left); err != nil {
			return 0, err
		}
	}

	return left, nil
}

// parseFactor handles numbers and parentheses
func (p *Parser) parseFactor() (float64, error) {
	// Handle unary plus and minus
	if p.peek() == '+' || p.peek() == '-' {
		op := p.peek()
		p.consume()

		value, err := p.parseFactor()
		if err != nil {
			return 0, err
		}

		if op == '-' {
			value = -value
		}

		return value, nil
	}

	// Handle parentheses
	if p.peek() == '(' {
		p.consume() // consume '('
		value, err := p.parseExpression()
		if err != nil {
			return 0, err
		}

		if p.peek() != ')' {
			return 0, ErrMismatchedParentheses
		}

		p.consume() // consume ')'
		return value, nil
	}

	// Handle numbers
	return p.parseNumber()
}

// parseNumber parses a numeric literal
func (p *Parser) parseNumber() (float64, error) {
	start := p.position

	// Parse integer part
	for p.position < len(p.expression) && unicode.IsDigit(rune(p.expression[p.position])) {
		p.position++
	}

	// Parse decimal part
	if p.position < len(p.expression) && p.expression[p.position] == '.' {
		p.position++

		// Must have at least one digit after decimal
		if p.position >= len(p.expression) || !unicode.IsDigit(rune(p.expression[p.position])) {
			return 0, ErrInvalidNumber
		}

		for p.position < len(p.expression) && unicode.IsDigit(rune(p.expression[p.position])) {
			p.position++
		}
	}

	if p.position == start {
		return 0, fmt.Errorf("%w: expected number at position %d", ErrInvalidExpression, p.position)
	}

	numberStr := p.expression[start:p.position]
	value, err := strconv.ParseFloat(numberStr, 64)
	if err != nil {
		return 0, fmt.Errorf("%w: %s", ErrInvalidNumber, err)
	}

	return value, nil
}

// peek returns the current character without consuming it
func (p *Parser) peek() byte {
	if p.position >= len(p.expression) {
		return 0
	}
	return p.expression[p.position]
}

// consume consumes the current character
func (p *Parser) consume() {
	if p.position < len(p.expression) {
		p.position++
	}
}

// EvaluateSimple evaluates a simple arithmetic expression
func EvaluateSimple(a, b float64, operator string) (float64, error) {
	switch operator {
	case "+":
		result := a + b
		if err := ValidateNumber(result); err != nil {
			return 0, err
		}
		return result, nil
	case "-":
		result := a - b
		if err := ValidateNumber(result); err != nil {
			return 0, err
		}
		return result, nil
	case "*":
		result := a * b
		if err := ValidateNumber(result); err != nil {
			return 0, err
		}
		return result, nil
	case "/":
		if b == 0 {
			return 0, ErrDivisionByZero
		}
		result := a / b
		if err := ValidateNumber(result); err != nil {
			return 0, err
		}
		return result, nil
	default:
		return 0, ErrInvalidOperator
	}
}
