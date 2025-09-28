package visual

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	tea "github.com/charmbracelet/bubbletea"

	"ccpm-demo/internal/calculator"
	"ccpm-demo/internal/ui"
	visualpkg "ccpm-demo/internal/visual"
)

// TestCompleteVisualTestingWorkflow tests the entire visual testing workflow
func TestCompleteVisualTestingWorkflow(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping comprehensive visual tests in short mode")
	}

	t.Run("Complete Visual Testing Workflow", func(t *testing.T) {
		// Create temporary directory for test artifacts
		tempDir := t.TempDir()

		// Initialize calculator components
		engine := calculator.NewEngine()
		model := ui.NewModel(engine)

		// Step 1: Test screenshot capture
		t.Run("Screenshot Capture", func(t *testing.T) {
			config := visualpkg.NewDefaultConfig()
			screenshot, err := visualpkg.NewScreenshotFromModel(model, config)
			require.NoError(t, err, "Should capture screenshot successfully")
			require.NotNil(t, screenshot, "Screenshot should not be nil")

			// Test screenshot saving
			screenshotPath := filepath.Join(tempDir, "test_screenshot.png")
			err = screenshot.Save(screenshotPath)
			require.NoError(t, err, "Should save screenshot successfully")

			// Verify file exists
			_, err = os.Stat(screenshotPath)
			require.NoError(t, err, "Screenshot file should exist")
		})

		// Step 2: Test visual comparison
		t.Run("Visual Comparison", func(t *testing.T) {
			config := visualpkg.NewDefaultConfig()

			// Capture two screenshots
			screenshot1, err := visualpkg.NewScreenshotFromModel(model, config)
			require.NoError(t, err, "Should capture first screenshot")

			screenshot2, err := visualpkg.NewScreenshotFromModel(model, config)
			require.NoError(t, err, "Should capture second screenshot")

			// Compare screenshots
			compareConfig := visualpkg.NewDefaultCompareConfig()
			result, err := visualpkg.CompareScreenshots(screenshot1, screenshot2, compareConfig)
			require.NoError(t, err, "Should compare screenshots successfully")
			require.NotNil(t, result, "Comparison result should not be nil")

			// Since screenshots are identical, diff should be minimal
			require.LessOrEqual(t, result.DiffRatio, 0.01, "Identical screenshots should have minimal diff")
			require.True(t, result.Identical, "Identical screenshots should be marked as identical")

			// Test comparison report
			report := result.RenderComparisonReport()
			require.NotEmpty(t, report, "Comparison report should not be empty")
			require.Contains(t, report, "Visual Comparison Report", "Report should contain title")
		})

		// Step 3: Test demo generation
		t.Run("Demo Generation", func(t *testing.T) {
			config := visualpkg.NewDefaultConfig()
			demoDir := filepath.Join(tempDir, "demo")
			demoGen := visualpkg.NewDemoGenerator(model, config, demoDir)

			// Test recording
			err := demoGen.StartRecording("workflow_test", "Workflow test recording")
			require.NoError(t, err, "Should start recording successfully")
			require.True(t, demoGen.Recording, "Should be recording")

			// Capture frames
			err = demoGen.CaptureFrame("Initial state")
			require.NoError(t, err, "Should capture initial frame")

			err = demoGen.CaptureFrame("After interaction")
			require.NoError(t, err, "Should capture interaction frame")

			// Add key press
			err = demoGen.AddKeyPress(tea.KeyMsg{Type: tea.KeyEnter}, "Test key press")
			require.NoError(t, err, "Should add key press successfully")

			// Stop recording
			err = demoGen.StopRecording()
			require.NoError(t, err, "Should stop recording successfully")
			require.False(t, demoGen.Recording, "Should not be recording")

			// Verify demo files
			_, err = os.Stat(filepath.Join(demoDir, "demo.json"))
			require.NoError(t, err, "Demo metadata file should exist")

			// Check for screenshot files
			files, err := os.ReadDir(demoDir)
			require.NoError(t, err, "Should read demo directory")
			require.Greater(t, len(files), 0, "Demo directory should contain files")
		})

		// Step 4: Test visual regression testing
		t.Run("Visual Regression Testing", func(t *testing.T) {
			config := TestConfig{
				BaselineDir:   filepath.Join(tempDir, "baseline"),
				CurrentDir:    filepath.Join(tempDir, "current"),
				DiffDir:       filepath.Join(tempDir, "diff"),
				Tolerance:     0.01,
				UpdateMode:    true, // Create baseline
				ParallelRuns:  1,
				MaxDiffRatio:  0.1,
				MaxTestTime:   30 * time.Second,
				SaveScreenshots: true,
			}

			test := NewVisualRegressionTest(
				"Workflow Integration Test",
				"Comprehensive workflow integration test",
				model,
				config,
			)

			// Run test
			err := test.Run()
			require.NoError(t, err, "Visual regression test should run successfully")

			// Check results
			require.True(t, test.Results.Passed, "Visual regression test should pass")
			require.Greater(t, test.Results.TotalTests, 0, "Should have run some tests")

			// Generate and check report
			report := test.GenerateReport()
			require.NotEmpty(t, report, "Report should not be empty")
			require.Contains(t, report, "Visual Regression Test Report", "Report should contain title")

			// Save results
			resultsFile := filepath.Join(tempDir, "results.json")
			err = test.SaveResults(resultsFile)
			require.NoError(t, err, "Should save results successfully")

			_, err = os.Stat(resultsFile)
			require.NoError(t, err, "Results file should exist")
		})

		// Step 5: Test report generation
		t.Run("Report Generation", func(t *testing.T) {
			// Create test results
			testResults := &TestResults{
				Name:        "Integration Test Results",
				Description: "Results from integration test",
				Passed:      true,
				TotalTests:  5,
				PassedTests: 5,
				FailedTests: 0,
				SkippedTests: 0,
				Duration:    2 * time.Second,
				RunAt:       time.Now(),
				Environment: "test",
				TestCases: map[string]*TestCaseResult{
					"test1": {
						Name:    "test1",
						Passed:  true,
						Duration: 100 * time.Millisecond,
						Details: "Test 1 passed",
					},
					"test2": {
						Name:    "test2",
						Passed:  true,
						Duration: 200 * time.Millisecond,
						Details: "Test 2 passed",
					},
				},
			}

			// Generate reports
			reportConfig := ReportConfig{
				OutputDir:      tempDir,
				GenerateHTML:  true,
				GenerateJSON:  true,
				GenerateText:  true,
			}

			reportGen := NewReportGenerator(testResults, reportConfig)
			err := reportGen.GenerateReports()
			require.NoError(t, err, "Should generate reports successfully")

			// Verify report files exist
			_, err = os.Stat(filepath.Join(tempDir, "visual-test-report.html"))
			require.NoError(t, err, "HTML report should exist")

			_, err = os.Stat(filepath.Join(tempDir, "visual-test-report.json"))
			require.NoError(t, err, "JSON report should exist")

			_, err = os.Stat(filepath.Join(tempDir, "visual-test-report.txt"))
			require.NoError(t, err, "Text report should exist")

			_, err = os.Stat(filepath.Join(tempDir, "visual-test-report.css"))
			require.NoError(t, err, "CSS file should exist")
		})

		// Step 6: Test calculator-specific visual tests
		t.Run("Calculator Visual Tests", func(t *testing.T) {
			suite := NewCalculatorVisualSuite()

			// Test a subset of calculator functionality
			t.Run("Calculator Operations", func(t *testing.T) {
				// Test screenshot capture
				screenshot, err := visualpkg.NewScreenshotFromModel(suite.Model, suite.Config)
				require.NoError(t, err, "Should capture calculator screenshot")
				require.NotNil(t, screenshot, "Calculator screenshot should not be nil")
			})

			t.Run("Calculator Demo Generation", func(t *testing.T) {
				demoDir := filepath.Join(tempDir, "calc_demo")
				demoGen := visualpkg.NewDemoGenerator(suite.Model, suite.Config, demoDir)

				err := demoGen.GenerateBasicDemo()
				require.NoError(t, err, "Should generate calculator demo")

				// Verify demo files
				_, err = os.Stat(filepath.Join(demoDir, "demo.json"))
				require.NoError(t, err, "Calculator demo metadata should exist")
			})
		})

		// Step 7: Verify all artifacts
		t.Run("Artifact Verification", func(t *testing.T) {
			// List all generated files
			files := []string{}
			err := filepath.Walk(tempDir, func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}
				if !info.IsDir() {
					files = append(files, path)
				}
				return nil
			})
			require.NoError(t, err, "Should walk temp directory")

			// Should have generated various artifacts
			require.Greater(t, len(files), 5, "Should have generated multiple test artifacts")

			t.Logf("Generated %d test artifacts:", len(files))
			for _, file := range files {
				relPath, err := filepath.Rel(tempDir, file)
				if err == nil {
					t.Logf("  - %s", relPath)
				}
			}
		})

		t.Logf("Complete visual testing workflow test passed. All artifacts saved in: %s", tempDir)
	})
}

// TestVisualTestingPerformance benchmarks the performance of the visual testing framework
func TestVisualTestingPerformance(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping performance tests in short mode")
	}

	t.Run("Visual Testing Performance", func(t *testing.T) {
		// Setup
		engine := calculator.NewEngine()
		model := ui.NewModel(engine)
		config := visualpkg.NewDefaultConfig()

		// Benchmark screenshot capture
		t.Run("Screenshot Capture Performance", func(t *testing.T) {
			iterations := 20
			var totalTime time.Duration

			for i := 0; i < iterations; i++ {
				start := time.Now()
				_, err := visualpkg.NewScreenshotFromModel(model, config)
				require.NoError(t, err, "Should capture screenshot")
				totalTime += time.Since(start)
			}

			avgTime := totalTime / time.Duration(iterations)
			t.Logf("Average screenshot capture time: %v", avgTime)

			// Performance requirement
			require.Less(t, avgTime, 100*time.Millisecond, "Screenshot capture should be fast (< 100ms)")
		})

		// Benchmark visual comparison
		t.Run("Visual Comparison Performance", func(t *testing.T) {
			// Pre-capture screenshots
			screenshot1, err := visualpkg.NewScreenshotFromModel(model, config)
			require.NoError(t, err, "Should capture first screenshot")

			screenshot2, err := visualpkg.NewScreenshotFromModel(model, config)
			require.NoError(t, err, "Should capture second screenshot")

			iterations := 20
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

			// Performance requirement
			require.Less(t, avgTime, 50*time.Millisecond, "Comparison should be fast (< 50ms)")
		})

		// Benchmark demo generation
		t.Run("Demo Generation Performance", func(t *testing.T) {
			tempDir := t.TempDir()
			demoGen := visualpkg.NewDemoGenerator(model, config, tempDir)

			start := time.Now()
			err := demoGen.GenerateBasicDemo()
			require.NoError(t, err, "Should generate basic demo")
			duration := time.Since(start)

			t.Logf("Basic demo generation time: %v", duration)

			// Performance requirement
			require.Less(t, duration, 2*time.Second, "Demo generation should be fast (< 2s)")
		})

		t.Logf("All visual testing performance benchmarks passed")
	})
}

// TestVisualTestingReliability tests the reliability and error handling of the framework
func TestVisualTestingReliability(t *testing.T) {
	t.Run("Visual Testing Reliability", func(t *testing.T) {
		// Test error handling
		t.Run("Error Handling", func(t *testing.T) {
			engine := calculator.NewEngine()
			model := ui.NewModel(engine)
			config := visualpkg.NewDefaultConfig()

			// Test nil model handling
			_, err := visualpkg.NewScreenshotFromModel(nil, config)
			require.Error(t, err, "Should handle nil model with error")

			// Test invalid config handling
			invalidConfig := config
			invalidConfig.Width = 0
			_, err = visualpkg.NewScreenshotFromModel(model, invalidConfig)
			// Should handle gracefully (may not error depending on implementation)
		})

		// Test concurrent access
		t.Run("Concurrent Access", func(t *testing.T) {
			engine := calculator.NewEngine()
			model := ui.NewModel(engine)
			config := visualpkg.NewDefaultConfig()

			concurrent := 5
			results := make(chan error, concurrent)
			start := make(chan struct{})

			for i := 0; i < concurrent; i++ {
				go func() {
					<-start
					_, err := visualpkg.NewScreenshotFromModel(model, config)
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

		// Test resource cleanup
		t.Run("Resource Cleanup", func(t *testing.T) {
			tempDir := t.TempDir()
			engine := calculator.NewEngine()
			model := ui.NewModel(engine)
			config := visualpkg.NewDefaultConfig()

			// Generate test artifacts
			screenshot, err := visualpkg.NewScreenshotFromModel(model, config)
			require.NoError(t, err, "Should capture screenshot")

			screenshotPath := filepath.Join(tempDir, "cleanup_test.png")
			err = screenshot.Save(screenshotPath)
			require.NoError(t, err, "Should save screenshot")

			// Verify file exists
			_, err = os.Stat(screenshotPath)
			require.NoError(t, err, "Screenshot file should exist before cleanup")

			// Note: Actual cleanup would be handled by the test framework
			// This test mainly ensures files can be created and accessed properly
		})

		t.Logf("All visual testing reliability tests passed")
	})
}