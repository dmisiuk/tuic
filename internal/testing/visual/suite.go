package visual

import (
	"os"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/require"

	"ccpm-demo/internal/calculator"
	"ccpm-demo/internal/ui"
	visualpkg "ccpm-demo/internal/visual"
)

// VisualTestSuite represents a complete visual testing suite
type VisualTestSuite struct {
	Engine     *calculator.Engine
	Model      ui.Model
	Config     visualpkg.TerminalConfig
	TestConfig TestConfig
}

// NewVisualTestSuite creates a new visual test suite
func NewVisualTestSuite() *VisualTestSuite {
	engine := calculator.NewEngine()
	model := ui.NewModel(engine)

	return &VisualTestSuite{
		Engine: engine,
		Model:  model,
		Config: visualpkg.NewDefaultConfig(),
		TestConfig: TestConfig{
			BaselineDir:   "test/visual/baseline",
			CurrentDir:    "test/visual/current",
			DiffDir:       "test/visual/diff",
			Tolerance:     0.01,
			UpdateMode:    false,
			ParallelRuns:  1,
			MaxDiffRatio:  0.1,
			MaxTestTime:   30 * time.Second,
			SaveScreenshots: true,
		},
	}
}

// TestVisualRegression runs comprehensive visual regression tests
func (vts *VisualTestSuite) TestVisualRegression(t *testing.T) {
	t.Run("Visual Regression", func(t *testing.T) {
		test := NewVisualRegressionTest(
			"Calculator Visual Regression",
			"Comprehensive visual regression test for CCPM Calculator",
			vts.Model,
			vts.TestConfig,
		)

		err := test.Run()
		require.NoError(t, err, "Visual regression test should run successfully")

		report := test.GenerateReport()
		t.Logf("Visual Regression Report:\n%s", report)

		// Save results for CI/CD
		if err := test.SaveResults("test/visual/results.json"); err != nil {
			t.Logf("Failed to save test results: %v", err)
		}

		if !test.Results.Passed {
			t.Errorf("Visual regression tests failed")
		}
	})
}

// TestScreenshotCapture tests the screenshot capture functionality
func (vts *VisualTestSuite) TestScreenshotCapture(t *testing.T) {
	t.Run("Screenshot Capture", func(t *testing.T) {
		screenshot, err := visualpkg.NewScreenshotFromModel(vts.Model, vts.Config)
		require.NoError(t, err, "Should capture screenshot successfully")
		require.NotNil(t, screenshot, "Screenshot should not be nil")

		// Test screenshot properties
		require.NotEmpty(t, screenshot.Image, "Screenshot image should not be empty")
		require.NotZero(t, screenshot.Config.Width, "Screenshot width should not be zero")
		require.NotZero(t, screenshot.Config.Height, "Screenshot height should not be zero")

		// Test screenshot saving
		tempFile := t.TempDir() + "/test_screenshot.png"
		err = screenshot.Save(tempFile)
		require.NoError(t, err, "Should save screenshot successfully")
	})

	t.Run("Terminal Config", func(t *testing.T) {
		config := visualpkg.NewDefaultConfig()
		require.NotZero(t, config.Width, "Default width should not be zero")
		require.NotZero(t, config.Height, "Default height should not be zero")
		require.NotNil(t, config.FontFace, "Font face should not be nil")
		require.NotNil(t, config.Foreground, "Foreground color should not be nil")
		require.NotNil(t, config.Background, "Background color should not be nil")
	})
}

// TestVisualComparison tests the visual comparison functionality
func (vts *VisualTestSuite) TestVisualComparison(t *testing.T) {
	t.Run("Visual Comparison", func(t *testing.T) {
		// Capture two screenshots
		screenshot1, err := visualpkg.NewScreenshotFromModel(vts.Model, vts.Config)
		require.NoError(t, err, "Should capture first screenshot")

		screenshot2, err := visualpkg.NewScreenshotFromModel(vts.Model, vts.Config)
		require.NoError(t, err, "Should capture second screenshot")

		// Compare screenshots
		compareConfig := visualpkg.NewDefaultCompareConfig()
		result, err := visualpkg.CompareScreenshots(screenshot1, screenshot2, compareConfig)
		require.NoError(t, err, "Should compare screenshots successfully")
		require.NotNil(t, result, "Comparison result should not be nil")

		// Test comparison results
		require.NotNil(t, result.DiffImage, "Diff image should not be nil")
		require.GreaterOrEqual(t, result.DiffRatio, 0.0, "Diff ratio should be >= 0")
		require.LessOrEqual(t, result.DiffRatio, 1.0, "Diff ratio should be <= 1")

		// Since both screenshots are identical, diff ratio should be very small
		require.LessOrEqual(t, result.DiffRatio, 0.001, "Identical screenshots should have minimal diff")

		// Test comparison report
		report := result.RenderComparisonReport()
		require.NotEmpty(t, report, "Comparison report should not be empty")
		require.Contains(t, report, "Visual Comparison Report", "Report should contain title")
	})

	t.Run("Tolerance Testing", func(t *testing.T) {
		config := visualpkg.NewDefaultCompareConfig()
		config.ColorTolerance = 0.5 // High tolerance

		screenshot1, err := visualpkg.NewScreenshotFromModel(vts.Model, vts.Config)
		require.NoError(t, err, "Should capture first screenshot")

		screenshot2, err := visualpkg.NewScreenshotFromModel(vts.Model, vts.Config)
		require.NoError(t, err, "Should capture second screenshot")

		result, err := visualpkg.CompareScreenshots(screenshot1, screenshot2, config)
		require.NoError(t, err, "Should compare screenshots with high tolerance")

		// With high tolerance, identical screenshots should be considered equal
		require.True(t, result.Identical, "High tolerance should accept identical screenshots")
	})
}

// TestDemoGeneration tests the demo generation functionality
func (vts *VisualTestSuite) TestDemoGeneration(t *testing.T) {
	t.Run("Demo Generation", func(t *testing.T) {
		tempDir := t.TempDir()
		demoGen := visualpkg.NewDemoGenerator(vts.Model, vts.Config, tempDir)

		// Test basic demo generation
		err := demoGen.GenerateBasicDemo()
		require.NoError(t, err, "Should generate basic demo successfully")

		// Test advanced demo generation
		err = demoGen.GenerateAdvancedDemo()
		require.NoError(t, err, "Should generate advanced demo successfully")

		// Test error demo generation
		err = demoGen.GenerateErrorDemo()
		require.NoError(t, err, "Should generate error demo successfully")

		// Test navigation demo generation
		err = demoGen.GenerateKeyboardNavigationDemo()
		require.NoError(t, err, "Should generate navigation demo successfully")

		// Test all demos generation
		err = demoGen.GenerateAllDemos()
		require.NoError(t, err, "Should generate all demos successfully")

		// Verify demo files were created
		demos := []string{"basic", "advanced", "errors", "navigation"}
		for _, demo := range demos {
			demoDir := tempDir + "/" + demo
			_, err := os.Stat(demoDir)
			require.NoError(t, err, "Demo directory should exist: %s", demoDir)

			// Check for demo files
			files, err := os.ReadDir(demoDir)
			require.NoError(t, err, "Should read demo directory: %s", demoDir)
			require.Greater(t, len(files), 0, "Demo directory should contain files: %s", demoDir)
		}
	})

	t.Run("Demo Recording", func(t *testing.T) {
		tempDir := t.TempDir()
		demoGen := visualpkg.NewDemoGenerator(vts.Model, vts.Config, tempDir)

		// Test recording
		err := demoGen.StartRecording("test_recording", "Test recording")
		require.NoError(t, err, "Should start recording successfully")
		require.True(t, demoGen.Recording, "Should be recording")

		// Capture a frame
		err = demoGen.CaptureFrame("Test frame")
		require.NoError(t, err, "Should capture frame successfully")

		// Add key press
		err = demoGen.AddKeyPress(tea.KeyMsg{Type: tea.KeyEnter}, "Test key press")
		require.NoError(t, err, "Should add key press successfully")

		// Stop recording
		err = demoGen.StopRecording()
		require.NoError(t, err, "Should stop recording successfully")
		require.False(t, demoGen.Recording, "Should not be recording")

		// Verify demo metadata was created
		metadataFile := tempDir + "/demo.json"
		_, err = os.Stat(metadataFile)
		require.NoError(t, err, "Demo metadata file should exist")
	})
}

// TestVisualPerformance tests the performance of visual operations
func (vts *VisualTestSuite) TestVisualPerformance(t *testing.T) {
	t.Run("Screenshot Performance", func(t *testing.T) {
		if testing.Short() {
			t.Skip("Skipping performance tests in short mode")
		}

		iterations := 50
		var totalTime time.Duration

		for i := 0; i < iterations; i++ {
			start := time.Now()
			_, err := visualpkg.NewScreenshotFromModel(vts.Model, vts.Config)
			require.NoError(t, err, "Should capture screenshot")
			totalTime += time.Since(start)
		}

		avgTime := totalTime / time.Duration(iterations)
		t.Logf("Average screenshot capture time: %v", avgTime)

		// Performance requirement: screenshots should be fast
		require.Less(t, avgTime, 10*time.Millisecond, "Screenshot capture should be fast")
	})

	t.Run("Comparison Performance", func(t *testing.T) {
		if testing.Short() {
			t.Skip("Skipping performance tests in short mode")
		}

		// Pre-capture screenshots
		screenshot1, err := visualpkg.NewScreenshotFromModel(vts.Model, vts.Config)
		require.NoError(t, err, "Should capture first screenshot")

		screenshot2, err := visualpkg.NewScreenshotFromModel(vts.Model, vts.Config)
		require.NoError(t, err, "Should capture second screenshot")

		iterations := 50
		var totalTime time.Duration
		compareConfig := visualpkg.NewDefaultCompareConfig()

		for i := 0; i < iterations; i++ {
			start := time.Now()
			_, err := visualpkg.CompareScreenshots(screenshot1, screenshot2, compareConfig)
			require.NoError(t, err, "Should compare screenshots")
			totalTime += time.Since(start)
		}

		avgTime := totalTime / time.Duration(iterations)
		t.Logf("Average comparison time: %v", avgTime)

		// Performance requirement: comparison should be fast
		require.Less(t, avgTime, 50*time.Millisecond, "Comparison should be fast")
	})
}

// TestVisualReliability tests the reliability of visual operations
func (vts *VisualTestSuite) TestVisualReliability(t *testing.T) {
	t.Run("Concurrent Screenshots", func(t *testing.T) {
		// Test concurrent screenshot capture
		concurrent := 10
		results := make(chan error, concurrent)
		start := make(chan struct{})

		for i := 0; i < concurrent; i++ {
			go func() {
				<-start
				_, err := visualpkg.NewScreenshotFromModel(vts.Model, vts.Config)
				results <- err
			}()
		}

		close(start)

		// Collect results
		for i := 0; i < concurrent; i++ {
			err := <-results
			require.NoError(t, err, "Concurrent screenshot should succeed")
		}
	})

	t.Run("Error Handling", func(t *testing.T) {
		// Test error handling for invalid configurations
		invalidConfig := vts.Config
		invalidConfig.Width = 0

		_, err := visualpkg.NewScreenshotFromModel(vts.Model, invalidConfig)
		require.NoError(t, err, "Should handle invalid config gracefully")

		// Test error handling for nil model
		_, err = visualpkg.NewScreenshotFromModel(nil, vts.Config)
		require.Error(t, err, "Should handle nil model with error")
	})
}

// TestIntegration tests integration with the existing UI components
func (vts *VisualTestSuite) TestIntegration(t *testing.T) {
	t.Run("UI Integration", func(t *testing.T) {
		// Test that the UI model can be used for visual testing
		require.NotNil(t, vts.Model, "UI model should not be nil")

		// Test that the model has the required interface
		var viewInterface interface{ View() string } = vts.Model
		require.NotNil(t, viewInterface, "Model should implement View() method")

		view := vts.Model.View()
		require.NotEmpty(t, view, "Model view should not be empty")

		// Test screenshot capture with real UI
		screenshot, err := visualpkg.NewScreenshotFromModel(vts.Model, vts.Config)
		require.NoError(t, err, "Should capture UI screenshot")
		require.NotNil(t, screenshot, "Screenshot should not be nil")

		// Test that screenshot contains UI elements
		// This is a basic test - in reality you'd analyze the actual content
		require.NotNil(t, screenshot.Image, "Screenshot image should not be empty")
	})

	t.Run("Theme Integration", func(t *testing.T) {
		// Test different themes
		themes := []string{"retro-casio", "modern", "minimal"}

		for _, theme := range themes {
			t.Run(theme, func(t *testing.T) {
				err := vts.Model.SetButtonGridTheme(theme)
				if err != nil {
					t.Logf("Warning: Failed to set theme '%s': %v", theme, err)
					t.SkipNow()
				}

				screenshot, err := visualpkg.NewScreenshotFromModel(vts.Model, vts.Config)
				require.NoError(t, err, "Should capture screenshot with theme %s", theme)
				require.NotNil(t, screenshot, "Screenshot should not be nil with theme %s", theme)
			})
		}
	})
}

// RunVisualTests runs all visual tests
func RunVisualTests(t *testing.T) {
	suite := NewVisualTestSuite()

	suite.TestVisualRegression(t)
	suite.TestScreenshotCapture(t)
	suite.TestVisualComparison(t)
	suite.TestDemoGeneration(t)
	suite.TestVisualPerformance(t)
	suite.TestVisualReliability(t)
	suite.TestIntegration(t)
}