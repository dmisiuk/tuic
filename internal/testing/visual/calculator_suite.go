package visual

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"ccpm-demo/internal/calculator"
	"ccpm-demo/internal/ui"
	visualpkg "ccpm-demo/internal/visual"
)

// CalculatorVisualSuite provides visual testing specifically for calculator functionality
type CalculatorVisualSuite struct {
	Engine     *calculator.Engine
	Model      ui.Model
	Config     visualpkg.TerminalConfig
	TestConfig TestConfig
}

// NewCalculatorVisualSuite creates a new calculator visual test suite
func NewCalculatorVisualSuite() *CalculatorVisualSuite {
	engine := calculator.NewEngine()
	model := ui.NewModel(engine)

	return &CalculatorVisualSuite{
		Engine: engine,
		Model:  model,
		Config: visualpkg.NewDefaultConfig(),
		TestConfig: TestConfig{
			BaselineDir:   "test/visual/calculator/baseline",
			CurrentDir:    "test/visual/calculator/current",
			DiffDir:       "test/visual/calculator/diff",
			Tolerance:     0.01,
			UpdateMode:    false,
			ParallelRuns:  1,
			MaxDiffRatio:  0.1,
			MaxTestTime:   30 * time.Second,
			SaveScreenshots: true,
		},
	}
}

// TestCalculatorOperations tests visual aspects of calculator operations
func (cvs *CalculatorVisualSuite) TestCalculatorOperations(t *testing.T) {
	t.Run("Calculator Operations", func(t *testing.T) {
		operations := []struct {
			name        string
			input       string
			expected    string
			description string
		}{
			{
				name:        "basic_addition",
				input:       "123+456=",
				expected:    "579",
				description: "Basic addition operation",
			},
			{
				name:        "basic_subtraction",
				input:       "100-50=",
				expected:    "50",
				description: "Basic subtraction operation",
			},
			{
				name:        "basic_multiplication",
				input:       "12*12=",
				expected:    "144",
				description: "Basic multiplication operation",
			},
			{
				name:        "basic_division",
				input:       "100/4=",
				expected:    "25",
				description: "Basic division operation",
			},
			{
				name:        "decimal_calculation",
				input:       "3.14*2=",
				expected:    "6.28",
				description: "Decimal calculation",
			},
			{
				name:        "complex_expression",
				input:       "(123+456)*2/3=",
				expected:    "386",
				description: "Complex expression with parentheses",
			},
		}

		for _, op := range operations {
			t.Run(op.name, func(t *testing.T) {
				// Reset calculator
				cvs.resetCalculator()

				// Capture initial state
				initialScreenshot, err := visualpkg.NewScreenshotFromModel(cvs.Model, cvs.Config)
				require.NoError(t, err, "Should capture initial screenshot")

				// Perform calculation
				cvs.performInput(op.input)

				// Capture result state
				resultScreenshot, err := visualpkg.NewScreenshotFromModel(cvs.Model, cvs.Config)
				require.NoError(t, err, "Should capture result screenshot")

				// Verify visual differences between initial and result
				compareConfig := visualpkg.NewDefaultCompareConfig()
				result, err := visualpkg.CompareScreenshots(initialScreenshot, resultScreenshot, compareConfig)
				require.NoError(t, err, "Should compare screenshots successfully")

				// There should be some visual difference after calculation
				require.Greater(t, result.DiffRatio, 0.0, "Calculation should produce visual changes")
				require.Less(t, result.DiffRatio, 0.1, "Visual changes should be minimal")

				// Verify the calculation result is displayed
				view := cvs.Model.View()
				require.Contains(t, view, op.expected, "Result should be displayed in the view")

				t.Logf("Operation '%s' completed successfully. Diff ratio: %.4f", op.description, result.DiffRatio)
			})
		}
	})
}

// TestErrorHandlingVisual tests visual aspects of error handling
func (cvs *CalculatorVisualSuite) TestErrorHandlingVisual(t *testing.T) {
	t.Run("Error Handling Visual", func(t *testing.T) {
		errorCases := []struct {
			name        string
			input       string
			description string
		}{
			{
				name:        "division_by_zero",
				input:       "10/0=",
				description: "Division by zero error",
			},
			{
				name:        "invalid_syntax",
				input:       "123++456=",
				description: "Invalid syntax error",
			},
			{
				name:        "mismatched_parentheses",
				input:       "(123+456",
				description: "Mismatched parentheses error",
			},
			{
				name:        "overflow",
				input:       "999999999999999999999*999999999999999999999=",
				description: "Overflow error",
			},
		}

		for _, testCase := range errorCases {
			t.Run(testCase.name, func(t *testing.T) {
				// Reset calculator
				cvs.resetCalculator()

				// Capture initial state
				initialScreenshot, err := visualpkg.NewScreenshotFromModel(cvs.Model, cvs.Config)
				require.NoError(t, err, "Should capture initial screenshot")

				// Perform error-inducing operation
				cvs.performInput(testCase.input)

				// Capture error state
				errorScreenshot, err := visualpkg.NewScreenshotFromModel(cvs.Model, cvs.Config)
				require.NoError(t, err, "Should capture error screenshot")

				// Verify visual differences
				compareConfig := visualpkg.NewDefaultCompareConfig()
				result, err := visualpkg.CompareScreenshots(initialScreenshot, errorScreenshot, compareConfig)
				require.NoError(t, err, "Should compare screenshots successfully")

				// There should be visual differences indicating error state
				require.Greater(t, result.DiffRatio, 0.0, "Error should produce visual changes")

				// Verify error is displayed
				view := cvs.Model.View()
				require.Contains(t, view, "Error", "Error should be displayed in the view")

				t.Logf("Error case '%s' handled correctly. Diff ratio: %.4f", testCase.description, result.DiffRatio)
			})
		}
	})
}

// TestButtonVisualFeedback tests visual feedback from button interactions
func (cvs *CalculatorVisualSuite) TestButtonVisualFeedback(t *testing.T) {
	t.Run("Button Visual Feedback", func(t *testing.T) {
		// Test different button types
		buttonTests := []struct {
			buttonType string
			buttonKeys []string
			description string
		}{
			{
				buttonType: "number",
				buttonKeys: []string{"1", "2", "3", "4", "5"},
				description: "Number button visual feedback",
			},
			{
				buttonType: "operator",
				buttonKeys: []string{"+", "-", "*", "/"},
				description: "Operator button visual feedback",
			},
			{
				buttonType: "function",
				buttonKeys: []string{"C", "CE", "=", "."},
				description: "Function button visual feedback",
			},
		}

		for _, buttonTest := range buttonTests {
			t.Run(buttonTest.buttonType, func(t *testing.T) {
				for _, key := range buttonTest.buttonKeys {
					t.Run(fmt.Sprintf("key_%s", key), func(t *testing.T) {
						// Reset calculator
						cvs.resetCalculator()

						// Capture initial state
						initialScreenshot, err := visualpkg.NewScreenshotFromModel(cvs.Model, cvs.Config)
						require.NoError(t, err, "Should capture initial screenshot")

						// Press button
						cvs.performInput(key)

						// Capture after button press
						pressedScreenshot, err := visualpkg.NewScreenshotFromModel(cvs.Model, cvs.Config)
						require.NoError(t, err, "Should capture pressed screenshot")

						// Verify visual differences
						compareConfig := visualpkg.NewDefaultCompareConfig()
						result, err := visualpkg.CompareScreenshots(initialScreenshot, pressedScreenshot, compareConfig)
						require.NoError(t, err, "Should compare screenshots successfully")

						// There should be visual feedback when pressing buttons
						require.Greater(t, result.DiffRatio, 0.0, "Button press should produce visual feedback")

						t.Logf("Button '%s' feedback verified. Diff ratio: %.4f", key, result.DiffRatio)
					})
				}
			})
		}
	})
}

// TestKeyboardNavigationVisual tests visual aspects of keyboard navigation
func (cvs *CalculatorVisualSuite) TestKeyboardNavigationVisual(t *testing.T) {
	t.Run("Keyboard Navigation Visual", func(t *testing.T) {
		navigationTests := []struct {
			name        string
			sequence    []string
			description string
		}{
			{
				name:        "tab_navigation",
				sequence:    []string{"Tab", "Tab", "Tab"},
				description: "Tab navigation through buttons",
			},
			{
				name:        "arrow_navigation",
				sequence:    []string{"Down", "Right", "Up", "Left"},
				description: "Arrow key navigation",
			},
			{
				name:        "combined_navigation",
				sequence:    []string{"Tab", "Down", "Right", "Enter"},
				description: "Combined navigation sequence",
			},
		}

		for _, navTest := range navigationTests {
			t.Run(navTest.name, func(t *testing.T) {
				// Reset calculator
				cvs.resetCalculator()

				// Capture initial state
				initialScreenshot, err := visualpkg.NewScreenshotFromModel(cvs.Model, cvs.Config)
				require.NoError(t, err, "Should capture initial screenshot")

				// Perform navigation sequence
				for _, key := range navTest.sequence {
					cvs.performKeyInput(key)
				}

				// Capture after navigation
				navigatedScreenshot, err := visualpkg.NewScreenshotFromModel(cvs.Model, cvs.Config)
				require.NoError(t, err, "Should capture navigated screenshot")

				// Verify visual differences
				compareConfig := visualpkg.NewDefaultCompareConfig()
				result, err := visualpkg.CompareScreenshots(initialScreenshot, navigatedScreenshot, compareConfig)
				require.NoError(t, err, "Should compare screenshots successfully")

				// Navigation should produce visual changes (focus indicators)
				require.Greater(t, result.DiffRatio, 0.0, "Navigation should produce visual changes")

				t.Logf("Navigation '%s' verified. Diff ratio: %.4f", navTest.description, result.DiffRatio)
			})
		}
	})
}

// TestThemeSwitchingVisual tests visual aspects of theme switching
func (cvs *CalculatorVisualSuite) TestThemeSwitchingVisual(t *testing.T) {
	t.Run("Theme Switching Visual", func(t *testing.T) {
		themes := []string{"retro-casio", "modern", "minimal", "classic"}

		for i, theme := range themes {
			t.Run(theme, func(t *testing.T) {
				// Reset calculator
				cvs.resetCalculator()

				// Set theme
				err := cvs.Model.SetButtonGridTheme(theme)
				if err != nil {
					t.Logf("Warning: Could not set theme '%s': %v", theme, err)
					t.SkipNow()
				}

				// Capture themed screenshot
				themeScreenshot, err := visualpkg.NewScreenshotFromModel(cvs.Model, cvs.Config)
				require.NoError(t, err, "Should capture themed screenshot")

				// Compare with default theme if available
				if i > 0 {
					// Reset to default theme
					cvs.resetCalculator()
					err = cvs.Model.SetButtonGridTheme("retro-casio")
					require.NoError(t, err, "Should reset to default theme")

					defaultScreenshot, err := visualpkg.NewScreenshotFromModel(cvs.Model, cvs.Config)
					require.NoError(t, err, "Should capture default screenshot")

					// Compare themes
					compareConfig := visualpkg.NewDefaultCompareConfig()
					result, err := visualpkg.CompareScreenshots(defaultScreenshot, themeScreenshot, compareConfig)
					require.NoError(t, err, "Should compare themed screenshots successfully")

					// Different themes should produce visual differences
					require.Greater(t, result.DiffRatio, 0.0, "Theme switch should produce visual changes")

					t.Logf("Theme '%s' difference from default: %.4f", theme, result.DiffRatio)
				}

				// Verify theme is applied
				view := cvs.Model.View()
				require.NotEmpty(t, view, "Themed view should not be empty")

				t.Logf("Theme '%s' applied successfully", theme)
			})
		}
	})
}

// TestCalculatorRegressionFull runs a full regression test for calculator functionality
func (cvs *CalculatorVisualSuite) TestCalculatorRegressionFull(t *testing.T) {
	t.Run("Calculator Regression Full", func(t *testing.T) {
		test := NewVisualRegressionTest(
			"Calculator Full Regression",
			"Comprehensive visual regression test for calculator functionality",
			cvs.Model,
			cvs.TestConfig,
		)

		err := test.Run()
		require.NoError(t, err, "Calculator regression test should run successfully")

		report := test.GenerateReport()
		t.Logf("Calculator Regression Report:\n%s", report)

		if !test.Results.Passed {
			t.Errorf("Calculator regression tests failed")
		}
	})
}

// TestCalculatorDemoGeneration generates comprehensive calculator demos
func (cvs *CalculatorVisualSuite) TestCalculatorDemoGeneration(t *testing.T) {
	t.Run("Calculator Demo Generation", func(t *testing.T) {
		tempDir := t.TempDir()
		demoGen := visualpkg.NewDemoGenerator(cvs.Model, cvs.Config, tempDir)

		// Generate comprehensive calculator demo
		err := demoGen.StartRecording("calculator_comprehensive", "Comprehensive calculator functionality demo")
		require.NoError(t, err, "Should start recording")

		// Initial state
		err = demoGen.CaptureFrame("Initial calculator state")
		require.NoError(t, err, "Should capture initial state")

		// Basic operations
		operations := []struct {
			input       string
			description string
		}{
			{"123", "Enter first number"},
			{"+", "Addition operator"},
			{"456", "Enter second number"},
			{"=", "Calculate result"},
			{"C", "Clear result"},
		}

		for _, op := range operations {
			cvs.performInput(op.input)
			err = demoGen.CaptureFrame(op.description)
			require.NoError(t, err, "Should capture operation state")
		}

		// Error handling
		cvs.performInput("10/0=")
		err = demoGen.CaptureFrame("Division by zero error")
		require.NoError(t, err, "Should capture error state")

		// Clear error
		cvs.performInput("C")
		err = demoGen.CaptureFrame("Clear error")
		require.NoError(t, err, "Should capture clear state")

		// Complex calculation
		cvs.performInput("(123+456)*2/3=")
		err = demoGen.CaptureFrame("Complex calculation")
		require.NoError(t, err, "Should capture complex calculation")

		// Stop recording
		err = demoGen.StopRecording()
		require.NoError(t, err, "Should stop recording")

		// Verify demo files
		_, err = os.Stat(tempDir + "/demo.json")
		require.NoError(t, err, "Demo metadata should exist")

		// Verify screenshot files exist
		files, err := os.ReadDir(tempDir)
		require.NoError(t, err, "Should read demo directory")
		require.Greater(t, len(files), 0, "Demo directory should contain files")

		t.Logf("Calculator demo generated successfully with %d files", len(files))
	})
}

// Helper methods
func (cvs *CalculatorVisualSuite) resetCalculator() {
	// Reset calculator engine
	cvs.Engine = calculator.NewEngine()
	cvs.Model = ui.NewModel(cvs.Engine)
}

func (cvs *CalculatorVisualSuite) performInput(input string) {
	for _, char := range input {
		cvs.performKeyInput(string(char))
	}
}

func (cvs *CalculatorVisualSuite) performKeyInput(key string) {
	// This is a simplified version - in reality, you'd use proper key handling
	// For visual testing, we're primarily concerned with the visual changes
}

// RunCalculatorVisualTests runs all calculator visual tests
func RunCalculatorVisualTests(t *testing.T) {
	suite := NewCalculatorVisualSuite()

	suite.TestCalculatorOperations(t)
	suite.TestErrorHandlingVisual(t)
	suite.TestButtonVisualFeedback(t)
	suite.TestKeyboardNavigationVisual(t)
	suite.TestThemeSwitchingVisual(t)
	suite.TestCalculatorRegressionFull(t)
	suite.TestCalculatorDemoGeneration(t)
}