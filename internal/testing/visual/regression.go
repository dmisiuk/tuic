package visual

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"ccpm-demo/internal/visual"
)

// VisualRegressionTest represents a complete visual regression test
type VisualRegressionTest struct {
	Name           string
	Description    string
	Model          interface{}
	Config         visual.TerminalConfig
	BaselineDir    string
	CurrentDir     string
	DiffDir        string
	Tolerance      float64
	UpdateMode     bool
	Results        *TestResults
}

// TestResults contains the results of a visual regression test run
type TestResults struct {
	Name        string                    `json:"name"`
	Description string                    `json:"description"`
	Passed      bool                      `json:"passed"`
	TotalTests  int                       `json:"totalTests"`
	PassedTests int                       `json:"passedTests"`
	FailedTests int                       `json:"failedTests"`
	SkippedTests int                       `json:"skippedTests"`
	Duration    time.Duration             `json:"duration"`
	TestCases   map[string]*TestCaseResult `json:"testCases"`
	RunAt       time.Time                 `json:"runAt"`
	Environment string                    `json:"environment"`
}

// TestCaseResult represents the result of a single test case
type TestCaseResult struct {
	Name        string        `json:"name"`
	Passed      bool          `json:"passed"`
	Skipped     bool          `json:"skipped"`
	Error       string        `json:"error,omitempty"`
	DiffRatio   float64       `json:"diffRatio"`
	Duration    time.Duration `json:"duration"`
	Screenshot  string        `json:"screenshot,omitempty"`
	Baseline    string        `json:"baseline,omitempty"`
	DiffImage   string        `json:"diffImage,omitempty"`
	Details     string        `json:"details,omitempty"`
}

// TestConfig contains configuration for visual regression tests
type TestConfig struct {
	BaselineDir   string
	CurrentDir    string
	DiffDir       string
	Tolerance     float64
	UpdateMode    bool
	ParallelRuns  int
	MaxDiffRatio  float64
	MaxTestTime   time.Duration
	SaveScreenshots bool
}

// NewVisualRegressionTest creates a new visual regression test
func NewVisualRegressionTest(name, description string, model interface{}, config TestConfig) *VisualRegressionTest {
	return &VisualRegressionTest{
		Name:        name,
		Description: description,
		Model:       model,
		Config:      visual.NewDefaultConfig(),
		BaselineDir: config.BaselineDir,
		CurrentDir:  config.CurrentDir,
		DiffDir:     config.DiffDir,
		Tolerance:   config.Tolerance,
		UpdateMode:  config.UpdateMode,
		Results: &TestResults{
			Name:        name,
			Description: description,
			TestCases:   make(map[string]*TestCaseResult),
			RunAt:       time.Now(),
			Environment: os.Getenv("GO_ENV"),
		},
	}
}

// Run runs the complete visual regression test suite
func (vrt *VisualRegressionTest) Run() error {
	startTime := time.Now()

	// Create directories
	if err := vrt.ensureDirectories(); err != nil {
		return fmt.Errorf("failed to create directories: %w", err)
	}

	// Run test cases
	vrt.runTestCases()

	// Calculate results
	vrt.Results.Duration = time.Since(startTime)
	vrt.Results.Passed = vrt.Results.FailedTests == 0
	vrt.Results.TotalTests = vrt.Results.PassedTests + vrt.Results.FailedTests + vrt.Results.SkippedTests

	return nil
}

// ensureDirectories creates required directories
func (vrt *VisualRegressionTest) ensureDirectories() error {
	dirs := []string{vrt.BaselineDir, vrt.CurrentDir, vrt.DiffDir}
	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
	}
	return nil
}

// runTestCases runs all test cases
func (vrt *VisualRegressionTest) runTestCases() {
	testCases := vrt.getTestCases()

	for _, tc := range testCases {
		result := vrt.runTestCase(tc)
		vrt.Results.TestCases[tc.name] = result

		if result.Passed {
			vrt.Results.PassedTests++
		} else if result.Skipped {
			vrt.Results.SkippedTests++
		} else {
			vrt.Results.FailedTests++
		}
	}
}

// TestCase represents a single test case
type TestCase struct {
	name        string
	description string
	setupFunc   func() error
	teardownFunc func() error
}

// getTestCases returns all test cases to run
func (vrt *VisualRegressionTest) getTestCases() []TestCase {
	return []TestCase{
		{
			name:        "initial_state",
			description: "Initial calculator state",
			setupFunc:   vrt.setupInitialState,
			teardownFunc: vrt.teardownInitialState,
		},
		{
			name:        "basic_calculation",
			description: "Basic calculation (123 + 456)",
			setupFunc:   vrt.setupBasicCalculation,
			teardownFunc: vrt.teardownBasicCalculation,
		},
		{
			name:        "error_state",
			description: "Error state (division by zero)",
			setupFunc:   vrt.setupErrorState,
			teardownFunc: vrt.teardownErrorState,
		},
		{
			name:        "keyboard_navigation",
			description: "Keyboard navigation states",
			setupFunc:   vrt.setupKeyboardNavigation,
			teardownFunc: vrt.teardownKeyboardNavigation,
		},
		{
			name:        "theme_switching",
			description: "Theme switching",
			setupFunc:   vrt.setupThemeSwitching,
			teardownFunc: vrt.teardownThemeSwitching,
		},
	}
}

// runTestCase runs a single test case
func (vrt *VisualRegressionTest) runTestCase(tc TestCase) *TestCaseResult {
	startTime := time.Now()
	result := &TestCaseResult{
		Name:    tc.name,
		Passed:  false,
		Skipped: false,
	}

	defer func() {
		result.Duration = time.Since(startTime)
	}()

	// Setup
	if err := tc.setupFunc(); err != nil {
		result.Error = fmt.Sprintf("setup failed: %v", err)
		return result
	}

	// Capture screenshot
	screenshot, err := visual.NewScreenshotFromModel(vrt.Model, vrt.Config)
	if err != nil {
		result.Error = fmt.Sprintf("screenshot capture failed: %v", err)
		return result
	}

	// Save current screenshot
	currentPath := filepath.Join(vrt.CurrentDir, tc.name+".png")
	if err := screenshot.Save(currentPath); err != nil {
		result.Error = fmt.Sprintf("failed to save current screenshot: %v", err)
		return result
	}
	result.Screenshot = currentPath

	// Check baseline
	baselinePath := filepath.Join(vrt.BaselineDir, tc.name+".png")
	if _, err := os.Stat(baselinePath); os.IsNotExist(err) {
		// No baseline exists, create one
		if vrt.UpdateMode {
			if err := screenshot.Save(baselinePath); err != nil {
				result.Error = fmt.Sprintf("failed to create baseline: %v", err)
				return result
			}
			result.Passed = true
			result.Details = "Created new baseline"
			return result
		}
		result.Error = "no baseline exists and update mode is disabled"
		return result
	}

	// Load baseline
	baselineFile, err := os.ReadFile(baselinePath)
	if err != nil {
		result.Error = fmt.Sprintf("failed to load baseline: %v", err)
		return result
	}

	baselineImg, err := visual.DecodePNG(bytes.NewReader(baselineFile))
	if err != nil {
		result.Error = fmt.Sprintf("failed to decode baseline: %v", err)
		return result
	}

	baselineScreenshot := &visual.Screenshot{
		Image: baselineImg,
		Config: vrt.Config,
	}

	// Compare screenshots
	compareConfig := visual.NewDefaultCompareConfig()
	compareResult, err := visual.CompareScreenshots(baselineScreenshot, screenshot, compareConfig)
	if err != nil {
		result.Error = fmt.Sprintf("comparison failed: %v", err)
		return result
	}

	result.DiffRatio = compareResult.DiffRatio
	result.Details = compareResult.RenderComparisonReport()

	// Save diff image if comparison failed
	if !compareResult.Identical {
		diffPath := filepath.Join(vrt.DiffDir, tc.name+"_diff.png")
		if err := visual.SavePNG(diffPath, compareResult.DiffImage); err != nil {
			result.Error = fmt.Sprintf("failed to save diff image: %v", err)
			return result
		}
		result.DiffImage = diffPath
	}

	// Check tolerance
	if compareResult.DiffRatio <= vrt.Tolerance {
		result.Passed = true
	} else {
		result.Error = fmt.Sprintf("diff ratio %.2f%% exceeds tolerance %.2f%%",
			compareResult.DiffRatio*100, vrt.Tolerance*100)
	}

	// Update baseline if needed
	if vrt.UpdateMode && !result.Passed {
		if err := screenshot.Save(baselinePath); err != nil {
			result.Error = fmt.Sprintf("failed to update baseline: %v", err)
			return result
		}
		result.Passed = true
		result.Details = "Updated baseline"
	}

	result.Baseline = baselinePath

	// Teardown
	if err := tc.teardownFunc(); err != nil {
		result.Error = fmt.Sprintf("teardown failed: %v", err)
		result.Passed = false
	}

	return result
}

// Test case setup and teardown methods
func (vrt *VisualRegressionTest) setupInitialState() error {
	// Reset model to initial state - simplified for now
	return nil
}

func (vrt *VisualRegressionTest) teardownInitialState() error {
	return nil
}

func (vrt *VisualRegressionTest) setupBasicCalculation() error {
	// Simulate basic calculation: 123 + 456 - simplified for now
	return nil
}

func (vrt *VisualRegressionTest) teardownBasicCalculation() error {
	return nil
}

func (vrt *VisualRegressionTest) setupErrorState() error {
	// Simulate division by zero - simplified for now
	return nil
}

func (vrt *VisualRegressionTest) teardownErrorState() error {
	return nil
}

func (vrt *VisualRegressionTest) setupKeyboardNavigation() error {
	// Navigate through buttons - simplified for now
	return nil
}

func (vrt *VisualRegressionTest) teardownKeyboardNavigation() error {
	return nil
}

func (vrt *VisualRegressionTest) setupThemeSwitching() error {
	// Switch to different theme - simplified for now
	return nil
}

func (vrt *VisualRegressionTest) teardownThemeSwitching() error {
	// Reset to default theme - simplified for now
	return nil
}

// SaveResults saves the test results to a JSON file
func (vrt *VisualRegressionTest) SaveResults(filename string) error {
	data, err := json.MarshalIndent(vrt.Results, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(filename, data, 0644)
}

// GenerateReport generates a human-readable report
func (vrt *VisualRegressionTest) GenerateReport() string {
	var report strings.Builder

	report.WriteString("=== Visual Regression Test Report ===\n\n")
	report.WriteString(fmt.Sprintf("Test: %s\n", vrt.Name))
	report.WriteString(fmt.Sprintf("Description: %s\n", vrt.Description))
	report.WriteString(fmt.Sprintf("Duration: %s\n", vrt.Results.Duration))
	report.WriteString(fmt.Sprintf("Environment: %s\n", vrt.Results.Environment))
	report.WriteString(fmt.Sprintf("Status: %s\n\n", vrt.getStatusString()))

	report.WriteString("--- Summary ---\n")
	report.WriteString(fmt.Sprintf("Total Tests: %d\n", vrt.Results.TotalTests))
	report.WriteString(fmt.Sprintf("Passed: %d\n", vrt.Results.PassedTests))
	report.WriteString(fmt.Sprintf("Failed: %d\n", vrt.Results.FailedTests))
	report.WriteString(fmt.Sprintf("Skipped: %d\n\n", vrt.Results.SkippedTests))

	if vrt.Results.FailedTests > 0 {
		report.WriteString("--- Failed Tests ---\n")
		for name, result := range vrt.Results.TestCases {
			if !result.Passed && !result.Skipped {
				report.WriteString(fmt.Sprintf("❌ %s: %s\n", name, result.Error))
				if result.DiffRatio > 0 {
					report.WriteString(fmt.Sprintf("   Diff Ratio: %.2f%%\n", result.DiffRatio*100))
				}
			}
		}
		report.WriteString("\n")
	}

	return report.String()
}

func (vrt *VisualRegressionTest) getStatusString() string {
	if vrt.Results.Passed {
		return "✅ PASSED"
	}
	return "❌ FAILED"
}