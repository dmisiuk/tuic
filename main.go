package main

import (
	"bufio"
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"

	"ccpm-demo/internal/calculator"
)

// Version information
var (
	Version   = "dev"
	BuildTime = "unknown"
	CommitHash = "unknown"
)

func main() {
	// Handle command line arguments
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "--version", "-v":
			printVersion()
			return
		case "--help", "-h":
			printHelp()
			return
		case "--eval":
			if len(os.Args) < 3 {
				fmt.Println("Error: --eval requires an expression")
				os.Exit(1)
			}
			evalExpression(strings.Join(os.Args[2:], " "))
			return
		}
	}

	// Interactive mode
	fmt.Printf("CCPM Calculator v%s\n", Version)
	fmt.Printf("Type 'help' for commands, 'quit' to exit\n\n")

	calc := calculator.NewCalculator()
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("> ")
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("Error reading input: %v\n", err)
			break
		}

		input = strings.TrimSpace(input)
		if input == "" {
			continue
		}

		switch input {
		case "quit", "exit", "q":
			fmt.Println("Goodbye!")
			return
		case "help", "h":
			printInteractiveHelp()
		case "version", "v":
			printVersion()
		case "vars":
			printVariables(calc)
		case "clear":
			calc.ClearVariables()
			fmt.Println("Variables cleared")
		default:
			if strings.HasPrefix(input, "set ") {
				handleVariableSet(calc, input[4:])
			} else {
				evalExpressionWithCalc(calc, input)
			}
		}
	}
}

func evalExpression(expr string) {
	calc := calculator.NewCalculator()
	evalExpressionWithCalc(calc, expr)
}

func evalExpressionWithCalc(calc *calculator.Calculator, expr string) {
	result, err := calc.Evaluate(expr)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("= %g\n", result)
}

func handleVariableSet(calc *calculator.Calculator, input string) {
	parts := strings.SplitN(input, "=", 2)
	if len(parts) != 2 {
		fmt.Println("Usage: set variable = value")
		return
	}

	varName := strings.TrimSpace(parts[0])
	valueStr := strings.TrimSpace(parts[1])

	value, err := strconv.ParseFloat(valueStr, 64)
	if err != nil {
		// Try to evaluate as expression
		result, evalErr := calc.Evaluate(valueStr)
		if evalErr != nil {
			fmt.Printf("Error parsing value: %v\n", err)
			return
		}
		value = result
	}

	calc.SetVariable(varName, value)
	fmt.Printf("Set %s = %g\n", varName, value)
}

func printVariables(calc *calculator.Calculator) {
	vars := calc.GetVariables()
	if len(vars) == 0 {
		fmt.Println("No variables defined")
		return
	}

	fmt.Println("Variables:")
	for name, value := range vars {
		fmt.Printf("  %s = %g\n", name, value)
	}
}

func printVersion() {
	fmt.Printf("CCPM Calculator v%s\n", Version)
	fmt.Printf("Build: %s\n", CommitHash)
	fmt.Printf("Time: %s\n", BuildTime)
	fmt.Printf("Go: %s\n", runtime.Version())
	fmt.Printf("OS/Arch: %s/%s\n", runtime.GOOS, runtime.GOARCH)
}

func printHelp() {
	fmt.Printf("CCPM Calculator v%s\n\n", Version)
	fmt.Printf("Usage: ccpm [options]\n\n")
	fmt.Printf("Options:\n")
	fmt.Printf("  -v, --version    Show version information\n")
	fmt.Printf("  -h, --help       Show this help message\n")
	fmt.Printf("  --eval EXPR      Evaluate expression and exit\n\n")
	fmt.Printf("Interactive Commands:\n")
	fmt.Printf("  help, h          Show interactive help\n")
	fmt.Printf("  version, v       Show version\n")
	fmt.Printf("  quit, exit, q    Exit calculator\n")
	fmt.Printf("  vars             Show all variables\n")
	fmt.Printf("  clear            Clear all variables\n")
	fmt.Printf("  set var = value  Set variable\n")
}

func printInteractiveHelp() {
	fmt.Println("Interactive Commands:")
	fmt.Println("  help, h          Show this help")
	fmt.Println("  version, v       Show version")
	fmt.Println("  quit, exit, q    Exit calculator")
	fmt.Println("  vars             Show all variables")
	fmt.Println("  clear            Clear all variables")
	fmt.Println("  set var = value  Set variable")
	fmt.Println("")
	fmt.Println("Mathematical Operations:")
	fmt.Println("  + - * /          Basic arithmetic")
	fmt.Println("  ^                Power")
	fmt.Println("  ( )              Grouping")
	fmt.Println("  sin, cos, tan    Trigonometric functions")
	fmt.Println("  sqrt             Square root")
	fmt.Println("  Variables can be used in expressions")
}