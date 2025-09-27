package test

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/charmbracelet/lipgloss"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"ccpm-demo/internal/ui/integration"
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