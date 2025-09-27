# Keyboard Shortcuts Guide

This guide provides keyboard shortcuts and interaction patterns for using the CCPM Calculator effectively.

## 🎯 Overview

The calculator supports both direct number input and expression evaluation through keyboard interactions. This guide covers all available shortcuts and input methods.

## 📝 Basic Input

### Number Keys
```
0-9    - Input digits
.      - Decimal point
+      - Addition
-      - Subtraction
*      - Multiplication
/      - Division
( )    - Parentheses
```

### Expression Examples
```
2+3*4         - Simple arithmetic
(10+5)*2      - Parentheses for precedence
3.14*2        - Decimal numbers
1+2+3+4       - Chain operations
10/2.5        - Division with decimals
```

## ⌨️ Keyboard Shortcuts

### Basic Operations
| Key | Action | Example |
|-----|--------|---------|
| `0-9` | Input digit | `1` `2` `3` |
| `.` | Decimal point | `3.14` |
| `+` | Addition | `2+3` |
| `-` | Subtraction | `10-4` |
| `*` | Multiplication | `6*7` |
| `/` | Division | `20/4` |
| `(` `)` | Parentheses | `(2+3)*4` |

### Control Keys
| Key | Action | Description |
|-----|--------|-------------|
| `Enter` | Calculate/Execute | Evaluate the current expression |
| `Escape` | Clear (C) | Clear all values and reset calculator |
| `Backspace` | Delete last digit | Remove last entered digit |
| `Delete` | Clear Entry (CE) | Clear current entry only |
| `Space` | Optional separator | Can be used for readability |

### Advanced Operations
| Key | Action | Example |
|-----|--------|---------|
| `=` | Equals/Calculate | Same as Enter |
| `^` | Power (if supported) | `2^3` = 8 |
| `%` | Percentage | `50%` = 0.5 |

## 🖱️ Mouse Interactions

### If using a GUI interface:
- **Click buttons**: Same as keyboard equivalents
- **Clear button**: Same as Escape key
- **Equals button**: Same as Enter key
- **Backspace button**: Same as Backspace key

## 📖 Usage Examples

### Basic Calculations
```
Input: 123+456
Press: Enter
Result: 579

Input: (10+5)*3
Press: Enter
Result: 45

Input: 3.14*2
Press: Enter
Result: 6.28
```

### Chain Operations
```
Step 1: Input 10
Step 2: Press +
Step 3: Input 5
Step 4: Press Enter
Result: 15

Step 5: Press *
Step 6: Input 2
Step 7: Press Enter
Result: 30
```

### Complex Expressions
```
Input: 2*(3+4)-5/2
Press: Enter
Result: 11.5

Input: ((1+2)*3+4)*2
Press: Enter
Result: 28
```

## 🚨 Error Handling

### Common Input Errors
- **Division by zero**: `5/0` → Error message
- **Mismatched parentheses**: `(2+3` → Error message
- **Invalid expression**: `2++3` → Error message
- **Empty expression**: Press Enter with no input → Error message

### Error Recovery
1. **Escape**: Clear everything and start over
2. **Backspace**: Remove last character and continue
3. **Continue**: Fix the expression and press Enter again

## 💡 Tips and Tricks

### Expression Building
- Use parentheses to control order: `(2+3)*4` vs `2+3*4`
- Chain operations: `10+5+3+2` works as expected
- Mixed operations: `2*3+4/2` evaluates correctly with precedence

### Efficiency Tips
- Use keyboard shortcuts faster than mouse clicks
- Learn common expression patterns
- Use parentheses for complex calculations
- Press Enter to quickly get results

### Advanced Usage
- Nested parentheses: `((1+2)*(3+4))/(5+6)`
- Decimal precision: `3.14159*2.71828`
- Large numbers: `123456789*987654321`

## 🔧 Technical Details

### Input Processing
- All input is processed as mathematical expressions
- Spaces are automatically removed
- Expressions are parsed with proper operator precedence
- Results are calculated with floating-point precision

### Supported Characters
- Digits: `0-9`
- Operators: `+ - * /`
- Parentheses: `( )`
- Decimal point: `.`
- Spaces: ` ` (optional)

### Expression Rules
- Must start with a number, `+`, `-`, or `(`
- Must end with a number or `)`
- Operators must be between numbers or parentheses
- Parentheses must be balanced
- Decimal points must be followed by digits

## 🎮 Interactive Calculator Example

```go
package main

import (
    "bufio"
    "fmt"
    "os"
    "strings"
    "ccpm-demo/internal/calculator"
)

func main() {
    engine := calculator.NewEngine()
    scanner := bufio.NewScanner(os.Stdin)

    fmt.Println("CCPM Calculator - Type 'exit' to quit")
    fmt.Println("Example expressions: 2+3*4, (10+5)*2, 3.14*2")
    fmt.Println("")

    for {
        fmt.Print("> ")
        if !scanner.Scan() {
            break
        }

        input := strings.TrimSpace(scanner.Text())
        if input == "exit" || input == "quit" {
            break
        }

        if input == "clear" || input == "c" {
            engine.Clear()
            fmt.Println("Calculator cleared")
            continue
        }

        result, err := engine.Evaluate(input)
        if err != nil {
            fmt.Printf("Error: %v\n", err)
            continue
        }

        fmt.Printf("= %f\n", result)
    }
}
```

## 📚 Related Documentation

- [Quick Start Guide](quickstart.md) - Basic usage instructions
- [Basic Usage](basic-usage.md) - Detailed usage examples
- [API Reference](../api/engine.md) - Technical API documentation
- [Troubleshooting](../troubleshooting/common-issues.md) - Common issues and solutions

---

**Pro Tip**: Practice with simple expressions first, then gradually move to more complex ones. The calculator follows standard mathematical precedence rules!