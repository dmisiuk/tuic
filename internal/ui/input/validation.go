package input

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
)

// InputValidator implements input validation for calculator operations
type InputValidator struct {
	maxInputLength   int
	maxDecimalPlaces int
	allowNegative    bool
	allowOperators   bool
	iv.lastValidationError string
}

// ValidationResult represents the result of input validation
type ValidationResult struct {
	IsValid   bool
	Value     string
	ErrorMsg  string
	Sanitized string
}

// NewInputValidator creates a new input validator with default settings
func NewInputValidator() *InputValidator {
	return &InputValidator{
		maxInputLength:   20,   // Maximum characters in input
		maxDecimalPlaces: 6,    // Maximum decimal places
		allowNegative:    true, // Allow negative numbers
		allowOperators:   true, // Allow operators
	}
}

// Validate implements the EventValidator interface
func (iv *InputValidator) Validate(event Event) bool {
	switch event.Type {
	case EventTypeKey:
		return iv.validateKeyEvent(event)
	case EventTypeMouse:
		return iv.validateMouseEvent(event)
	default:
		return true // Allow system and focus events by default
	}
}

// GetValidationError implements the EventValidator interface
func (iv *InputValidator) GetValidationError() string {
	// This will be set during the last validation operation
	return iv.iv.lastValidationError
}


// validateKeyEvent validates keyboard input events
func (iv *InputValidator) validateKeyEvent(event Event) bool {
	keyEvent, ok := event.Data.(KeyEvent)
	if !ok {
		return false
	}

	switch keyEvent.Action {
	case KeyActionNumber:
		return iv.validateNumberInput(keyEvent.Value)
	case KeyActionOperator:
		return iv.validateOperatorInput(keyEvent.Value)
	case KeyActionEquals:
		return iv.validateEqualsInput()
	case KeyActionBackspace:
		return iv.validateBackspaceInput()
	case KeyActionClear:
		return true // Clear is always valid
	default:
		return true // Allow other actions by default
	}
}

// validateMouseEvent validates mouse input events
func (iv *InputValidator) validateMouseEvent(event Event) bool {
	mouseEvent, ok := event.Data.(MouseEvent)
	if !ok {
		return false
	}

	switch mouseEvent.Action.Type {
	case "number":
		return iv.validateNumberInput(mouseEvent.Action.Value)
	case "operator":
		return iv.validateOperatorInput(mouseEvent.Action.Value)
	case "equals":
		return iv.validateEqualsInput()
	case "clear":
		return true // Clear is always valid
	case "backspace":
		return iv.validateBackspaceInput()
	default:
		return true // Allow other actions by default
	}
}

// validateNumberInput validates number input (digits and decimal point)
func (iv *InputValidator) validateNumberInput(value string) bool {
	if value == "" {
		iv.iv.lastValidationError = "Empty number input"
		return false
	}

	// Check if it's a single character (digit or decimal point)
	if len(value) > 1 {
		iv.iv.lastValidationError = "Number input must be a single character"
		return false
	}

	char := value[0]

	// Allow digits 0-9
	if char >= '0' && char <= '9' {
		iv.lastValidationError = ""
		return true
	}

	// Allow decimal point
	if char == '.' {
		iv.lastValidationError = ""
		return true
	}

	iv.lastValidationError = fmt.Sprintf("Invalid number input: %s", value)
	return false
}

// validateOperatorInput validates operator input
func (iv *InputValidator) validateOperatorInput(operator string) bool {
	if !iv.allowOperators {
		iv.lastValidationError = "Operators not allowed"
		return false
	}

	validOperators := []string{"+", "-", "*", "/"}
	for _, op := range validOperators {
		if operator == op {
			iv.lastValidationError = ""
			return true
		}
	}

	iv.lastValidationError = fmt.Sprintf("Invalid operator: %s", operator)
	return false
}

// validateEqualsInput validates equals operation
func (iv *InputValidator) validateEqualsInput() bool {
	// Equals is always valid, but we could add validation here
	// for example, checking if there's a valid expression to evaluate
	iv.lastValidationError = ""
	return true
}

// validateBackspaceInput validates backspace operation
func (iv *InputValidator) validateBackspaceInput() bool {
	// Backspace is always valid
	iv.lastValidationError = ""
	return true
}

// ValidateExpression validates a complete mathematical expression
func (iv *InputValidator) ValidateExpression(expression string) ValidationResult {
	result := ValidationResult{
		IsValid:   false,
		Value:     expression,
		ErrorMsg:  "",
		Sanitized: "",
	}

	// Remove extra whitespace
	sanitized := strings.TrimSpace(expression)
	sanitized = strings.ReplaceAll(sanitized, "  ", " ")

	// Check maximum length
	if utf8.RuneCountInString(sanitized) > iv.maxInputLength {
		result.ErrorMsg = fmt.Sprintf("Input too long (max %d characters)", iv.maxInputLength)
		return result
	}

	// Validate expression structure
	if !iv.isValidExpressionStructure(sanitized) {
		result.ErrorMsg = iv.lastValidationError
		return result
	}

	// Validate each component
	if !iv.validateExpressionComponents(sanitized) {
		result.ErrorMsg = iv.lastValidationError
		return result
	}

	result.IsValid = true
	result.Sanitized = sanitized
	return result
}

// isValidExpressionStructure validates the overall structure of an expression
func (iv *InputValidator) isValidExpressionStructure(expression string) bool {
	if expression == "" {
		return true // Empty expression is valid (clear state)
	}

	// Check for valid start and end
	if !iv.isValidExpressionStart(expression) {
		iv.lastValidationError = "Invalid expression start"
		return false
	}

	if !iv.isValidExpressionEnd(expression) {
		iv.lastValidationError = "Invalid expression end"
		return false
	}

	// Check for balanced parentheses (if supported)
	if !iv.hasBalancedParentheses(expression) {
		iv.lastValidationError = "Unbalanced parentheses"
		return false
	}

	return true
}

// isValidExpressionStart checks if the expression starts with a valid token
func (iv *InputValidator) isValidExpressionStart(expression string) bool {
	if expression == "" {
		return true
	}

	firstChar := expression[0]

	// Can start with: digit, decimal point, minus sign (for negative numbers)
	return (firstChar >= '0' && firstChar <= '9') ||
		firstChar == '.' ||
		(firstChar == '-' && iv.allowNegative)
}

// isValidExpressionEnd checks if the expression ends with a valid token
func (iv *InputValidator) isValidExpressionEnd(expression string) bool {
	if expression == "" {
		return true
	}

	lastChar := expression[len(expression)-1]

	// Can end with: digit, decimal point
	return (lastChar >= '0' && lastChar <= '9') || lastChar == '.'
}

// hasBalancedParentheses checks if parentheses are balanced
func (iv *InputValidator) hasBalancedParentheses(expression string) bool {
	balance := 0
	for _, char := range expression {
		switch char {
		case '(':
			balance++
		case ')':
			balance--
			if balance < 0 {
				return false
			}
		}
	}
	return balance == 0
}

// validateExpressionComponents validates each component of the expression
func (iv *InputValidator) validateExpressionComponents(expression string) bool {
	tokens := iv.tokenizeExpression(expression)

	for i, token := range tokens {
		if !iv.isValidToken(token) {
			return false
		}

		// Check operator placement (not at start/end unless it's a negative sign)
		if iv.isOperator(token) {
			if i == 0 && token != "-" {
				iv.lastValidationError = "Operator cannot be at start"
				return false
			}
			if i == len(tokens)-1 {
				iv.lastValidationError = "Operator cannot be at end"
				return false
			}
		}
	}

	return true
}

// tokenizeExpression splits an expression into tokens
func (iv *InputValidator) tokenizeExpression(expression string) []string {
	var tokens []string
	var currentToken strings.Builder

	for _, char := range expression {
		if unicode.IsSpace(char) {
			if currentToken.Len() > 0 {
				tokens = append(tokens, currentToken.String())
				currentToken.Reset()
			}
		} else if iv.isOperatorToken(char) {
			if currentToken.Len() > 0 {
				tokens = append(tokens, currentToken.String())
				currentToken.Reset()
			}
			tokens = append(tokens, string(char))
		} else {
			currentToken.WriteRune(char)
		}
	}

	if currentToken.Len() > 0 {
		tokens = append(tokens, currentToken.String())
	}

	return tokens
}

// isOperatorToken checks if a character is an operator
func (iv *InputValidator) isOperatorToken(char rune) bool {
	return char == '+' || char == '-' || char == '*' || char == '/'
}

// isOperator checks if a token is an operator
func (iv *InputValidator) isOperator(token string) bool {
	return token == "+" || token == "-" || token == "*" || token == "/"
}

// isValidToken validates a single token
func (iv *InputValidator) isValidToken(token string) bool {
	if iv.isOperator(token) {
		return true
	}

	// Check if it's a valid number
	if _, err := strconv.ParseFloat(token, 64); err == nil {
		return true
	}

	// Check if it has valid decimal places
	if strings.Contains(token, ".") {
		parts := strings.Split(token, ".")
		if len(parts) != 2 {
			iv.lastValidationError = fmt.Sprintf("Invalid number format: %s", token)
			return false
		}

		// Check decimal places
		if len(parts[1]) > iv.maxDecimalPlaces {
			iv.lastValidationError = fmt.Sprintf("Too many decimal places (max %d)", iv.maxDecimalPlaces)
			return false
		}
	}

	iv.lastValidationError = fmt.Sprintf("Invalid token: %s", token)
	return false
}

// SanitizeInput sanitizes input by removing invalid characters
func (iv *InputValidator) SanitizeInput(input string) string {
	var sanitized strings.Builder

	for _, char := range input {
		if iv.isValidChar(char) {
			sanitized.WriteRune(char)
		}
	}

	return sanitized.String()
}

// isValidChar checks if a character is valid for calculator input
func (iv *InputValidator) isValidChar(char rune) bool {
	// Allow digits
	if char >= '0' && char <= '9' {
		return true
	}

	// Allow decimal point
	if char == '.' {
		return true
	}

	// Allow operators
	if iv.isOperatorToken(char) {
		return true
	}

	// Allow negative sign at start
	if char == '-' && iv.allowNegative {
		return true
	}

	// Allow whitespace
	if unicode.IsSpace(char) {
		return true
	}

	return false
}

// ValidateNumberInput validates number input with current context
func (iv *InputValidator) ValidateNumberInput(currentInput, newChar string) ValidationResult {
	result := ValidationResult{
		IsValid:   false,
		Value:     currentInput + newChar,
		ErrorMsg:  "",
		Sanitized: "",
	}

	// Check maximum length
	if utf8.RuneCountInString(currentInput)+1 > iv.maxInputLength {
		result.ErrorMsg = fmt.Sprintf("Input too long (max %d characters)", iv.maxInputLength)
		return result
	}

	// Validate the new character
	if !iv.validateNumberInput(newChar) {
		result.ErrorMsg = iv.lastValidationError
		return result
	}

	// Check for multiple decimal points
	if newChar == "." && strings.Contains(currentInput, ".") {
		result.ErrorMsg = "Multiple decimal points not allowed"
		return result
	}

	// Check for leading zero issues
	if iv.hasLeadingZeroIssue(currentInput, newChar) {
		result.ErrorMsg = "Invalid number format"
		return result
	}

	result.IsValid = true
	result.Sanitized = currentInput + newChar
	return result
}

// hasLeadingZeroIssue checks for invalid leading zero patterns
func (iv *InputValidator) hasLeadingZeroIssue(currentInput, newChar string) bool {
	if newChar == "0" && currentInput == "" {
		return false // Single zero is fine
	}

	if newChar == "0" && currentInput == "0" {
		return true // Multiple leading zeros
	}

	if newChar >= "1" && newChar <= "9" && currentInput == "0" {
		return true // Leading zero before other digits
	}

	return false
}

// SetMaxInputLength sets the maximum input length
func (iv *InputValidator) SetMaxInputLength(length int) {
	iv.maxInputLength = length
}

// SetMaxDecimalPlaces sets the maximum decimal places
func (iv *InputValidator) SetMaxDecimalPlaces(places int) {
	iv.maxDecimalPlaces = places
}

// SetAllowNegative sets whether negative numbers are allowed
func (iv *InputValidator) SetAllowNegative(allow bool) {
	iv.allowNegative = allow
}

// SetAllowOperators sets whether operators are allowed
func (iv *InputValidator) SetAllowOperators(allow bool) {
	iv.allowOperators = allow
}

// GetMaxInputLength returns the maximum input length
func (iv *InputValidator) GetMaxInputLength() int {
	return iv.maxInputLength
}

// GetMaxDecimalPlaces returns the maximum decimal places
func (iv *InputValidator) GetMaxDecimalPlaces() int {
	return iv.maxDecimalPlaces
}

// GetAllowNegative returns whether negative numbers are allowed
func (iv *InputValidator) GetAllowNegative() bool {
	return iv.allowNegative
}

// GetAllowOperators returns whether operators are allowed
func (iv *InputValidator) GetAllowOperators() bool {
	return iv.allowOperators
}