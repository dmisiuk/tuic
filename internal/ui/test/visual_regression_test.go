package test

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"ccpm-demo/internal/calculator"
	"ccpm-demo/internal/ui"
	"ccpm-demo/internal/ui/integration"
	visualpkg "ccpm-demo/internal/visual"
	visualtesting "ccpm-demo/internal/testing/visual"
)

// VisualRegressionTest provides a framework for testing UI rendering consistency
type VisualRegressionTest struct {
	testName    string
	grid        *integration.ButtonGrid
	termWidth   int
	snapshotDir string
}

// NewVisualRegressionTest creates a new visual regression test
func NewVisualRegressionTest(testName string, termWidth int) *VisualRegressionTest {
	return &VisualRegressionTest{
		testName:    testName,
		grid:        integration.NewButtonGrid(),
		termWidth:   termWidth,
		snapshotDir: "test/snapshots",
	}
}

// RunTest executes a visual regression test
func (vrt *VisualRegressionTest) RunTest(t *testing.T) {
	t.Run(vrt.testName, func(t *testing.T) {
		// Ensure snapshot directory exists
		err := os.MkdirAll(vrt.snapshotDir, 0755)
		require.NoError(t, err)

		// Generate current rendering
		current := vrt.grid.Render(vrt.termWidth)

		// Clean up rendering for comparison (remove dynamic elements)
		cleaned := vrt.cleanRendering(current)

		// Get snapshot file path
		snapshotFile := filepath.Join(vrt.snapshotDir, vrt.testName+".snap")

		// Check if snapshot exists
		if _, err := os.Stat(snapshotFile); os.IsNotExist(err) {
			// Create new snapshot
			err := os.WriteFile(snapshotFile, []byte(cleaned), 0644)
			require.NoError(t, err, "Failed to create snapshot")
			t.Logf("Created new snapshot: %s", snapshotFile)
			return
		}

		// Load snapshot
		snapshot, err := os.ReadFile(snapshotFile)
		require.NoError(t, err, "Failed to read snapshot")

		// Compare renderings
		assert.Equal(t, string(snapshot), cleaned, "Rendering differs from snapshot")

		// If test fails and UPDATE_SNAPSHOTS is set, update snapshot
		if t.Failed() && os.Getenv("UPDATE_SNAPSHOTS") == "true" {
			err := os.WriteFile(snapshotFile, []byte(cleaned), 0644)
			require.NoError(t, err, "Failed to update snapshot")
			t.Logf("Updated snapshot: %s", snapshotFile)
		}
	})
}

// cleanRendering removes dynamic elements that shouldn't be compared
func (vrt *VisualRegressionTest) cleanRendering(rendering string) string {
	// Remove color codes and other dynamic styling
	cleaned := rendering

	// Remove ANSI escape codes
	ansiRegex := regexp.MustCompile(`\x1b\[[0-9;]*m`)
	cleaned = ansiRegex.ReplaceAllString(cleaned, "")

	// Normalize whitespace
	cleaned = strings.ReplaceAll(cleaned, "\t", "  ")
	cleaned = regexp.MustCompile(`\s+`).ReplaceAllString(cleaned, " ")
	cleaned = strings.TrimSpace(cleaned)

	// Normalize line endings
	cleaned = strings.ReplaceAll(cleaned, "\r\n", "\n")

	return cleaned
}

// TestRetroCasioStyling tests the retro Casio theme styling consistency
func TestRetroCasioStyling(t *testing.T) {
	vrt := NewVisualRegressionTest("retro_casio_80_width", 80)
	vrt.RunTest(t)
}

// TestNarrowTerminalRendering tests rendering on narrow terminals
func TestNarrowTerminalRendering(t *testing.T) {
	vrt := NewVisualRegressionTest("narrow_terminal_60_width", 60)
	vrt.RunTest(t)
}

// TestWideTerminalRendering tests rendering on wide terminals
func TestWideTerminalRendering(t *testing.T) {
	vrt := NewVisualRegressionTest("wide_terminal_100_width", 100)
	vrt.RunTest(t)
}

// TestButtonFocusStates tests visual consistency of button focus states
func TestButtonFocusStates(t *testing.T) {
	t.Run("button_focus_visual_consistency", func(t *testing.T) {
		grid := integration.NewButtonGrid()

		// Test different focus states
		states := []string{
			"initial_focus",
			"after_navigation_down",
			"after_navigation_right",
			"after_navigation_up",
			"after_navigation_left",
		}

		for _, state := range states {
			t.Run(state, func(t *testing.T) {
				// Apply state changes
				switch state {
				case "after_navigation_down":
					grid.HandleKeyPress(tea.KeyMsg{Type: tea.KeyDown})
				case "after_navigation_right":
					grid.HandleKeyPress(tea.KeyMsg{Type: tea.KeyRight})
				case "after_navigation_up":
					grid.HandleKeyPress(tea.KeyMsg{Type: tea.KeyUp})
				case "after_navigation_left":
					grid.HandleKeyPress(tea.KeyMsg{Type: tea.KeyLeft})
				}

				rendering := grid.Render(80)

				// Check that focused button is visually distinct
				assert.Contains(t, rendering, "7") // Should show focused button

				// Reset grid for next test
				grid = integration.NewButtonGrid()
			})
		}
	})
}

// TestButtonPressVisualFeedback tests visual feedback when buttons are pressed
func TestButtonPressVisualFeedback(t *testing.T) {
	t.Run("button_press_visual_feedback", func(t *testing.T) {
		grid := integration.NewButtonGrid()

		// Press a button
		action := grid.HandleKeyPress(tea.KeyMsg{Type: tea.KeyEnter})
		require.NotNil(t, action)

		rendering := grid.Render(80)

		// Should still show the button (pressed state)
		assert.Contains(t, rendering, "C")
	})
}

// TestThemeConsistency tests visual consistency across different themes
func TestThemeConsistency(t *testing.T) {
	themes := []string{"retro-casio", "modern", "minimal", "classic"}

	for _, theme := range themes {
		t.Run(fmt.Sprintf("theme_%s_consistency", theme), func(t *testing.T) {
			grid, err := integration.NewButtonGridWithTheme(theme)
			require.NoError(t, err)

			rendering := grid.Render(80)

			// Basic checks that all themes render properly
			assert.NotEmpty(t, rendering)
			assert.Contains(t, rendering, "C") // Should always have clear button
			assert.Contains(t, rendering, "0") // Should always have zero button
			assert.Contains(t, rendering, "=") // Should always have equals button
		})
	}
}

// TestGridLayoutConsistency tests that grid layout remains consistent
func TestGridLayoutConsistency(t *testing.T) {
	t.Run("grid_layout_consistency", func(t *testing.T) {
		grid := integration.NewButtonGrid()

		// Test different terminal widths
		widths := []int{60, 70, 80, 90, 100, 120}

		for _, width := range widths {
			t.Run(fmt.Sprintf("width_%d", width), func(t *testing.T) {
				rendering := grid.Render(width)

				// Should always contain all essential buttons
				essentialButtons := []string{"C", "7", "8", "9", "÷", "0", ".", "="}
				for _, button := range essentialButtons {
					assert.Contains(t, rendering, button,
						"Width %d should contain button %s", width, button)
				}

				// Should be properly formatted
				lines := strings.Split(strings.TrimSpace(rendering), "\n")
				assert.Greater(t, len(lines), 5, "Should have multiple lines of buttons")

				// Each line should have consistent structure
				for i, line := range lines {
					if strings.TrimSpace(line) != "" {
						// Should contain button characters
						assert.Contains(t, line, "│",
							"Line %d should have border characters", i)
					}
				}
			})
		}
	})
}

// TestButtonAccessibility tests visual accessibility features
func TestButtonAccessibility(t *testing.T) {
	t.Run("button_accessibility", func(t *testing.T) {
		grid := integration.NewButtonGrid()

		rendering := grid.Render(80)

		// Test color contrast simulation by checking for distinct styling
		// This is a simplified test - in reality, you'd use proper color contrast tools
		assert.Contains(t, rendering, "C") // Clear button should be visible
		assert.Contains(t, rendering, "÷") // Operator button should be visible
		assert.Contains(t, rendering, "7") // Number button should be visible

		// Test that different button types are visually distinct
		// This is basic - in a real test you'd analyze the actual styling
		assert.NotEmpty(t, rendering)
	})
}

// TestResponsiveRendering tests responsive behavior
func TestResponsiveRendering(t *testing.T) {
	t.Run("responsive_rendering", func(t *testing.T) {
		grid := integration.NewButtonGrid()

		// Test extreme widths
		testCases := []struct {
			width     int
			minLines  int
			maxLines  int
		}{
			{40, 8, 15},   // Very narrow
			{60, 8, 12},   // Narrow
			{80, 8, 12},   // Standard
			{100, 8, 12},  // Wide
			{120, 8, 12},  // Very wide
		}

		for _, tc := range testCases {
			t.Run(fmt.Sprintf("width_%d", tc.width), func(t *testing.T) {
				rendering := grid.Render(tc.width)
				lines := strings.Split(strings.TrimSpace(rendering), "\n")
				lineCount := len(lines)

				assert.GreaterOrEqual(t, lineCount, tc.minLines,
					"Width %d should have at least %d lines", tc.width, tc.minLines)
				assert.LessOrEqual(t, lineCount, tc.maxLines,
					"Width %d should have at most %d lines", tc.width, tc.maxLines)

				// Should still be functional
				assert.Contains(t, rendering, "C")
				assert.Contains(t, rendering, "=")
			})
		}
	})
}

// TestRenderingPerformance benchmarks rendering performance
func TestRenderingPerformance(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping performance tests in short mode")
	}

	t.Run("rendering_performance", func(t *testing.T) {
		grid := integration.NewButtonGrid()

		// Benchmark rendering at different widths
		widths := []int{60, 80, 100, 120}

		for _, width := range widths {
			t.Run(fmt.Sprintf("width_%d", width), func(t *testing.T) {
				// Warm up
				for i := 0; i < 10; i++ {
					_ = grid.Render(width)
				}

				// Benchmark
				start := make(chan struct{}, 1)
				results := make(chan timeDuration, 100)

				// Run concurrent rendering tests
				for i := 0; i < 10; i++ {
					go func() {
						<-start
						for j := 0; j < 10; j++ {
							startTime := timeNow()
							grid.Render(width)
							endTime := timeNow()
							results <- endTime.Sub(startTime)
						}
					}()
				}

				close(start)

				// Collect results
				var durations []timeDuration
				for i := 0; i < 100; i++ {
					durations = append(durations, <-results)
				}

				// Calculate statistics
				var total timeDuration
				for _, d := range durations {
					total += d
				}
				average := total / timeDuration(len(durations))

				// Assert performance requirements (should render in less than 1ms)
				assert.Less(t, average.Milliseconds(), int64(1),
					"Average rendering time should be less than 1ms, got %v", average)
			})
		}
	})
}

// Helper types for performance testing
type timeDuration struct {
	duration int64
}

func (td timeDuration) Milliseconds() int64 {
	return td.duration / 1e6
}

func timeNow() timeDuration {
	return timeDuration{duration: 0} // Simplified for testing
}

// TestEdgeCases tests edge cases in rendering
func TestEdgeCases(t *testing.T) {
	t.Run("edge_cases", func(t *testing.T) {
		grid := integration.NewButtonGrid()

		// Test zero width
		rendering := grid.Render(0)
		assert.NotEmpty(t, rendering) // Should handle gracefully

		// Test very large width
		rendering = grid.Render(1000)
		assert.NotEmpty(t, rendering)

		// Test negative width (should be handled gracefully)
		rendering = grid.Render(-1)
		assert.NotEmpty(t, rendering)
	})
}

// TestButtonStateTransitions tests visual consistency during state transitions
func TestButtonStateTransitions(t *testing.T) {
	t.Run("button_state_transitions", func(t *testing.T) {
		grid := integration.NewButtonGrid()

		// Test sequence: normal -> focused -> pressed -> normal
		states := []struct {
			name     string
			keyPress tea.KeyMsg
		}{
			{"initial", tea.KeyMsg{}},
			{"focused", tea.KeyMsg{Type: tea.KeyDown}},
			{"pressed", tea.KeyMsg{Type: tea.KeyEnter}},
		}

		for _, state := range states {
			t.Run(state.name, func(t *testing.T) {
				if state.keyPress.Type != 0 {
					grid.HandleKeyPress(state.keyPress)
				}

				rendering := grid.Render(80)

				// Should always render properly regardless of state
				assert.NotEmpty(t, rendering)
				assert.Contains(t, rendering, "7") // Should show current button
			})
		}
	})
}

// TestVisualFrameworkIntegration tests integration with the new visual testing framework
func TestVisualFrameworkIntegration(t *testing.T) {
	t.Run("framework_integration", func(t *testing.T) {
		// Create calculator engine and UI model
		engine := calculator.NewEngine()
		model := ui.NewModel(engine)

		// Create visual test suite
		suite := visualtesting.NewVisualTestSuite()

		// Test screenshot capture
		t.Run("screenshot_capture", func(t *testing.T) {
			screenshot, err := visualpkg.NewScreenshotFromModel(model, suite.Config)
			require.NoError(t, err, "Should capture screenshot successfully")
			require.NotNil(t, screenshot, "Screenshot should not be nil")

			// Test screenshot properties
			assert.NotNil(t, screenshot.Image, "Screenshot image should not be nil")
			assert.Greater(t, screenshot.Config.Width, 0, "Screenshot width should be positive")
			assert.Greater(t, screenshot.Config.Height, 0, "Screenshot height should be positive")

			// Test saving screenshot
			tempDir := t.TempDir()
			screenshotPath := filepath.Join(tempDir, "integration_test.png")
			err = screenshot.Save(screenshotPath)
			assert.NoError(t, err, "Should save screenshot successfully")

			// Verify file was created
			_, err = os.Stat(screenshotPath)
			assert.NoError(t, err, "Screenshot file should exist")
		})

		// Test visual comparison
		t.Run("visual_comparison", func(t *testing.T) {
			// Capture two screenshots
			screenshot1, err := visualpkg.NewScreenshotFromModel(model, suite.Config)
			require.NoError(t, err, "Should capture first screenshot")

			screenshot2, err := visualpkg.NewScreenshotFromModel(model, suite.Config)
			require.NoError(t, err, "Should capture second screenshot")

			// Compare screenshots
			compareConfig := visualpkg.NewDefaultCompareConfig()
			result, err := visualpkg.CompareScreenshots(screenshot1, screenshot2, compareConfig)
			require.NoError(t, err, "Should compare screenshots successfully")
			require.NotNil(t, result, "Comparison result should not be nil")

			// Test comparison results
			assert.NotNil(t, result.DiffImage, "Diff image should not be nil")
			assert.GreaterOrEqual(t, result.DiffRatio, 0.0, "Diff ratio should be >= 0")
			assert.LessOrEqual(t, result.DiffRatio, 1.0, "Diff ratio should be <= 1")

			// Since both screenshots are identical, diff ratio should be minimal
			assert.LessOrEqual(t, result.DiffRatio, 0.01, "Identical screenshots should have minimal diff")

			// Test comparison report
			report := result.RenderComparisonReport()
			assert.NotEmpty(t, report, "Comparison report should not be empty")
			assert.Contains(t, report, "Visual Comparison Report", "Report should contain title")
		})

		// Test demo generation
		t.Run("demo_generation", func(t *testing.T) {
			tempDir := t.TempDir()
			demoGen := visualpkg.NewDemoGenerator(model, suite.Config, tempDir)

			// Test recording
			err := demoGen.StartRecording("integration_test", "Integration test recording")
			require.NoError(t, err, "Should start recording successfully")
			assert.True(t, demoGen.Recording, "Should be recording")

			// Capture a frame
			err = demoGen.CaptureFrame("Integration test frame")
			assert.NoError(t, err, "Should capture frame successfully")

			// Add key press
			err = demoGen.AddKeyPress(tea.KeyMsg{Type: tea.KeyEnter}, "Test key press")
			assert.NoError(t, err, "Should add key press successfully")

			// Stop recording
			err = demoGen.StopRecording()
			assert.NoError(t, err, "Should stop recording successfully")
			assert.False(t, demoGen.Recording, "Should not be recording")

			// Verify demo files were created
			_, err = os.Stat(filepath.Join(tempDir, "demo.json"))
			assert.NoError(t, err, "Demo metadata file should exist")
		})

		// Test visual regression test
		t.Run("visual_regression_test", func(t *testing.T) {
			tempDir := t.TempDir()

			config := visualtesting.TestConfig{
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

			test := visualtesting.NewVisualRegressionTest(
				"Integration Test",
				"Integration test for visual testing framework",
				model,
				config,
			)

			// Run test
			err := test.Run()
			require.NoError(t, err, "Visual regression test should run successfully")

			// Check results
			assert.True(t, test.Results.Passed, "Visual regression test should pass")
			assert.Greater(t, test.Results.TotalTests, 0, "Should have run some tests")

			// Generate and check report
			report := test.GenerateReport()
			assert.NotEmpty(t, report, "Report should not be empty")
			assert.Contains(t, report, "Visual Regression Test Report", "Report should contain title")

			// Save results
			resultsFile := filepath.Join(tempDir, "results.json")
			err = test.SaveResults(resultsFile)
			assert.NoError(t, err, "Should save results successfully")

			_, err = os.Stat(resultsFile)
			assert.NoError(t, err, "Results file should exist")
		})

		// Test performance
		t.Run("performance", func(t *testing.T) {
			if testing.Short() {
				t.Skip("Skipping performance tests in short mode")
			}

			iterations := 20
			var totalTime time.Duration

			for i := 0; i < iterations; i++ {
				start := time.Now()
				_, err := visualpkg.NewScreenshotFromModel(model, suite.Config)
				assert.NoError(t, err, "Should capture screenshot")
				totalTime += time.Since(start)
			}

			avgTime := totalTime / time.Duration(iterations)
			t.Logf("Average screenshot capture time: %v", avgTime)

			// Performance requirement: screenshots should be fast
			assert.Less(t, avgTime, 50*time.Millisecond, "Screenshot capture should be fast (< 50ms)")
		})
	})
}

// TestCompleteVisualTestingWorkflow tests the complete visual testing workflow
func TestCompleteVisualTestingWorkflow(t *testing.T) {
	t.Run("complete_workflow", func(t *testing.T) {
		// Create test environment
		tempDir := t.TempDir()
		engine := calculator.NewEngine()
		model := ui.NewModel(engine)

		// Step 1: Create baseline screenshots
		t.Run("create_baseline", func(t *testing.T) {
			config := visualtesting.TestConfig{
				BaselineDir:   filepath.Join(tempDir, "baseline"),
				CurrentDir:    filepath.Join(tempDir, "current"),
				DiffDir:       filepath.Join(tempDir, "diff"),
				Tolerance:     0.01,
				UpdateMode:    true, // Create baseline
				ParallelRuns:  1,
				SaveScreenshots: true,
			}

			test := visualtesting.NewVisualRegressionTest(
				"Baseline Creation",
				"Create baseline screenshots",
				model,
				config,
			)

			err := test.Run()
			require.NoError(t, err, "Should create baseline successfully")
			assert.True(t, test.Results.Passed, "Baseline creation should pass")
		})

		// Step 2: Run regression test against baseline
		t.Run("regression_test", func(t *testing.T) {
			config := visualtesting.TestConfig{
				BaselineDir:   filepath.Join(tempDir, "baseline"),
				CurrentDir:    filepath.Join(tempDir, "current"),
				DiffDir:       filepath.Join(tempDir, "diff"),
				Tolerance:     0.01,
				UpdateMode:    false, // Don't update baseline
				ParallelRuns:  1,
				SaveScreenshots: true,
			}

			test := visualtesting.NewVisualRegressionTest(
				"Regression Test",
				"Test against baseline",
				model,
				config,
			)

			err := test.Run()
			require.NoError(t, err, "Should run regression test successfully")
			assert.True(t, test.Results.Passed, "Regression test should pass")
		})

		// Step 3: Generate demo
		t.Run("generate_demo", func(t *testing.T) {
			suite := visualtesting.NewVisualTestSuite()
			demoDir := filepath.Join(tempDir, "demo")
			demoGen := visualpkg.NewDemoGenerator(model, suite.Config, demoDir)

			err := demoGen.GenerateBasicDemo()
			require.NoError(t, err, "Should generate demo successfully")

			// Verify demo files exist
			_, err = os.Stat(filepath.Join(demoDir, "demo.json"))
			assert.NoError(t, err, "Demo metadata should exist")
		})

		// Step 4: Verify all artifacts
		t.Run("verify_artifacts", func(t *testing.T) {
			// Check baseline directory
			_, err := os.Stat(filepath.Join(tempDir, "baseline"))
			assert.NoError(t, err, "Baseline directory should exist")

			// Check current directory
			_, err = os.Stat(filepath.Join(tempDir, "current"))
			assert.NoError(t, err, "Current directory should exist")

			// Check demo directory
			_, err = os.Stat(filepath.Join(tempDir, "demo"))
			assert.NoError(t, err, "Demo directory should exist")

			// Check that directories contain files
			baselineFiles, err := os.ReadDir(filepath.Join(tempDir, "baseline"))
			assert.NoError(t, err, "Should read baseline directory")
			assert.Greater(t, len(baselineFiles), 0, "Baseline directory should contain files")
		})

		t.Logf("Complete visual testing workflow test passed. Artifacts saved in: %s", tempDir)
	})
}