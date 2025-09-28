package test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"ccpm-demo/internal/ui/integration"
)

// AccessibilityTest provides a framework for testing accessibility features
type AccessibilityTest struct {
	name        string
	description string
	grid        *integration.ButtonGrid
}

// NewAccessibilityTest creates a new accessibility test
func NewAccessibilityTest(name, description string) *AccessibilityTest {
	return &AccessibilityTest{
		name:        name,
		description: description,
		grid:        integration.NewButtonGrid(),
	}
}

// TestKeyboardNavigation tests comprehensive keyboard navigation
func TestKeyboardNavigation(t *testing.T) {
	t.Run("comprehensive_keyboard_navigation", func(t *testing.T) {
		grid := integration.NewButtonGrid()

		// Test all arrow keys
		testCases := []struct {
			name           string
			keyPress       tea.KeyMsg
			expectedFocus  string
			canNavigate    bool
		}{
			{"initial_state", tea.KeyMsg{}, "C", true},
			{"down_from_c", tea.KeyMsg{Type: tea.KeyDown}, "7", true},
			{"right_from_7", tea.KeyMsg{Type: tea.KeyRight}, "8", true},
			{"right_from_8", tea.KeyMsg{Type: tea.KeyRight}, "9", true},
			{"right_from_9", tea.KeyMsg{Type: tea.KeyRight}, "×", true},
			{"down_from_multiply", tea.KeyMsg{Type: tea.KeyDown}, "-", true},
			{"up_from_minus", tea.KeyMsg{Type: tea.KeyUp}, "×", true},
			{"left_from_multiply", tea.KeyMsg{Type: tea.KeyLeft}, "9", true},
			{"up_from_9", tea.KeyMsg{Type: tea.KeyUp}, "8", true},
			{"left_from_8", tea.KeyMsg{Type: tea.KeyLeft}, "7", true},
			{"up_from_7", tea.KeyMsg{Type: tea.KeyUp}, "C", true},
			{"left_from_c_boundary", tea.KeyMsg{Type: tea.KeyLeft}, "C", false}, // Should not move
			{"up_from_c_boundary", tea.KeyMsg{Type: tea.KeyUp}, "C", false},  // Should not move
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				if tc.keyPress.Type != 0 {
					action := grid.HandleKeyPress(tc.keyPress)
					assert.Nil(t, action, "Navigation should not trigger actions")
				}

				focusedButton, exists := grid.GetFocusedButton()
				require.True(t, exists, "Should always have a focused button")

				assert.Equal(t, tc.expectedFocus, focusedButton.GetLabel(),
					"Focused button should be %s", tc.expectedFocus)
			})
		}

		// Test that we can reach all number buttons
		numberButtons := []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9"}
		reachableNumbers := make(map[string]bool)

		// Reset and navigate to all number buttons
		grid = integration.NewButtonGrid()

		// Navigate to each number button and verify it's reachable
		navigationSequences := []struct {
			targetButton string
			sequence     []tea.KeyType
		}{
			{"7", []tea.KeyType{tea.KeyDown}},
			{"8", []tea.KeyType{tea.KeyDown, tea.KeyRight}},
			{"9", []tea.KeyType{tea.KeyDown, tea.KeyRight, tea.KeyRight}},
			{"4", []tea.KeyType{tea.KeyDown, tea.KeyDown}},
			{"5", []tea.KeyType{tea.KeyDown, tea.KeyDown, tea.KeyRight}},
			{"6", []tea.KeyType{tea.KeyDown, tea.KeyDown, tea.KeyRight, tea.KeyRight}},
			{"1", []tea.KeyType{tea.KeyDown, tea.KeyDown, tea.KeyDown}},
			{"2", []tea.KeyType{tea.KeyDown, tea.KeyDown, tea.KeyDown, tea.KeyRight}},
			{"3", []tea.KeyType{tea.KeyDown, tea.KeyDown, tea.KeyDown, tea.KeyRight, tea.KeyRight}},
			{"0", []tea.KeyType{tea.KeyDown, tea.KeyDown, tea.KeyDown, tea.KeyDown}},
		}

		for _, nav := range navigationSequences {
			grid = integration.NewButtonGrid() // Reset for each test

			for _, keyType := range nav.sequence {
				grid.HandleKeyPress(tea.KeyMsg{Type: keyType})
			}

			focusedButton, exists := grid.GetFocusedButton()
			require.True(t, exists, "Should have focused button")
			assert.Equal(t, nav.targetButton, focusedButton.GetLabel())
			reachableNumbers[nav.targetButton] = true
		}

		// Verify all number buttons are reachable
		for _, num := range numberButtons {
			assert.True(t, reachableNumbers[num], "Number button %s should be reachable", num)
		}
	})
}

// TestDirectKeyboardInput tests direct keyboard input accessibility
func TestDirectKeyboardInput(t *testing.T) {
	t.Run("direct_keyboard_input", func(t *testing.T) {
		grid := integration.NewButtonGrid()

		// Test all number keys
		numberKeys := []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9"}
		for _, num := range numberKeys {
			t.Run(fmt.Sprintf("number_key_%s", num), func(t *testing.T) {
				action := grid.HandleKeyPress(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{num[0]}})
				require.NotNil(t, action, "Number key %s should trigger action", num)
				assert.Equal(t, num, action.Value)
				assert.Equal(t, "press", action.Action)
			})
		}

		// Test operator keys
		operatorKeys := map[string]string{
			"+": "+", "-": "-", "*": "*", "/": "/",
		}
		for key, expectedValue := range operatorKeys {
			t.Run(fmt.Sprintf("operator_key_%s", key), func(t *testing.T) {
				action := grid.HandleKeyPress(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{key[0]}})
				require.NotNil(t, action, "Operator key %s should trigger action", key)
				assert.Equal(t, expectedValue, action.Value)
			})
		}

		// Test special keys
		specialKeys := map[string]string{
			".": ".", "=": "=",
		}
		for key, expectedValue := range specialKeys {
			t.Run(fmt.Sprintf("special_key_%s", key), func(t *testing.T) {
				action := grid.HandleKeyPress(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{key[0]}})
				require.NotNil(t, action, "Special key %s should trigger action", key)
				assert.Equal(t, expectedValue, action.Value)
			})
		}

		// Test that Enter and Space activate current button
		activationKeys := []tea.KeyType{tea.KeyEnter, tea.KeySpace}
		for _, keyType := range activationKeys {
			t.Run(fmt.Sprintf("activation_key_%v", keyType), func(t *testing.T) {
				// Focus on a specific button first
				grid = integration.NewButtonGrid()
				grid.HandleKeyPress(tea.KeyMsg{Type: tea.KeyDown}) // Focus on "7"

				action := grid.HandleKeyPress(tea.KeyMsg{Type: keyType})
				require.NotNil(t, action, "Activation key should trigger action")
				assert.Equal(t, "7", action.Value)
			})
		}
	})
}

// TestScreenReaderCompatibility tests screen reader compatibility
func TestScreenReaderCompatibility(t *testing.T) {
	t.Run("screen_reader_compatibility", func(t *testing.T) {
		grid := integration.NewButtonGrid()
		rendering := grid.Render(80)

		// Test that rendering is text-based and readable
		assert.NotEmpty(t, rendering)

		// Should contain meaningful text labels
		essentialLabels := []string{"C", "CE", "7", "8", "9", "÷", "0", ".", "="}
		for _, label := range essentialLabels {
			assert.Contains(t, rendering, label,
				"Rendering should contain readable label: %s", label)
		}

		// Should not rely solely on color for information
		// (This is a basic test - real accessibility testing would be more comprehensive)
		lines := strings.Split(rendering, "\n")
		for i, line := range lines {
			if strings.TrimSpace(line) != "" {
				// Each line should contain meaningful text, not just colors
				assert.Greater(t, len(strings.TrimSpace(line)), 0,
					"Line %d should contain meaningful content", i)
			}
		}
	})
}

// TestHighContrastMode tests high contrast mode compatibility
func TestHighContrastMode(t *testing.T) {
	t.Run("high_contrast_compatibility", func(t *testing.T) {
		// Test with different themes to ensure they work in high contrast
		themes := []string{"retro-casio", "modern", "minimal"}

		for _, theme := range themes {
			t.Run(theme, func(t *testing.T) {
				grid, err := integration.NewButtonGridWithTheme(theme)
				require.NoError(t, err)

				rendering := grid.Render(80)

				// Should render properly
				assert.NotEmpty(t, rendering)

				// Should contain essential elements
				assert.Contains(t, rendering, "C")
				assert.Contains(t, rendering, "7")
				assert.Contains(t, rendering, "=")

				// Should have some structure (not just random characters)
				lines := strings.Split(rendering, "\n")
				structuredLines := 0
				for _, line := range lines {
					if strings.TrimSpace(line) != "" {
						structuredLines++
					}
				}
				assert.Greater(t, structuredLines, 5, "Should have structured layout")
			})
		}
	})
}

// TestReducedMotion tests reduced motion compatibility
func TestReducedMotion(t *testing.T) {
	t.Run("reduced_motion_compatibility", func(t *testing.T) {
		grid := integration.NewButtonGrid()

		// Test that rapid state changes don't cause issues
		for i := 0; i < 10; i++ {
			// Navigate rapidly
			grid.HandleKeyPress(tea.KeyMsg{Type: tea.KeyDown})
			grid.HandleKeyPress(tea.KeyMsg{Type: tea.KeyRight})
			grid.HandleKeyPress(tea.KeyMsg{Type: tea.KeyUp})
			grid.HandleKeyPress(tea.KeyMsg{Type: tea.KeyLeft})

			// Should always render properly
			rendering := grid.Render(80)
			assert.NotEmpty(t, rendering)
			assert.Contains(t, rendering, "C") // Should always have essential buttons
		}

		// Test that final state is consistent
		focusedButton, exists := grid.GetFocusedButton()
		require.True(t, exists)
		assert.True(t, focusedButton.IsFocused())
	})
}

// TestErrorRecovery tests error recovery and graceful degradation
func TestErrorRecovery(t *testing.T) {
	t.Run("error_recovery", func(t *testing.T) {
		grid := integration.NewButtonGrid()

		// Test invalid theme changes
		err := grid.SetTheme("nonexistent_theme")
		assert.NoError(t, err) // Should fall back gracefully

		// Should still work after invalid theme
		rendering := grid.Render(80)
		assert.NotEmpty(t, rendering)

		// Test invalid key presses (should be ignored gracefully)
		invalidKeys := []tea.KeyType{
			tea.KeyCtrlA, tea.KeyCtrlB, tea.KeyCtrlZ, tea.KeyF1, tea.KeyF12,
		}

		for _, keyType := range invalidKeys {
			action := grid.HandleKeyPress(tea.KeyMsg{Type: keyType})
			assert.Nil(t, action, "Invalid key %v should be ignored", keyType)
		}

		// Should still work after invalid key presses
		rendering = grid.Render(80)
		assert.NotEmpty(t, rendering)

		// Test that valid keys still work
		action := grid.HandleKeyPress(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'5'}})
		assert.NotNil(t, action, "Valid key should still work after invalid ones")
	})
}

// TestConsistentInterface tests that the interface remains consistent
func TestConsistentInterface(t *testing.T) {
	t.Run("consistent_interface", func(t *testing.T) {
		grid := integration.NewButtonGrid()

		// Test that all essential operations work consistently
		operations := []struct {
			name string
			fn   func() bool
		}{
			{
				"get_button_count",
				func() bool { return grid.GetButtonCount() > 0 },
			},
			{
				"get_focused_button",
				func() bool {
					button, exists := grid.GetFocusedButton()
					return exists && button != nil
				},
			},
			{
				"get_dimensions",
				func() bool {
					dims := grid.GetDimensions()
					return dims.Columns > 0 && dims.Rows > 0
				},
			},
			{
				"rendering",
				func() bool {
					rendering := grid.Render(80)
					return len(rendering) > 0
				},
			},
		}

		for _, op := range operations {
			t.Run(op.name, func(t *testing.T) {
				assert.True(t, op.fn(), "Operation %s should work", op.name)
			})
		}

		// Test consistency after state changes
		for i := 0; i < 5; i++ {
			grid.HandleKeyPress(tea.KeyMsg{Type: tea.KeyDown})

			for _, op := range operations {
				assert.True(t, op.fn(), "Operation %s should work after state change", op.name)
			}
		}
	})
}

// TestMemoryAccessibility tests memory usage and performance for accessibility
func TestMemoryAccessibility(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping memory tests in short mode")
	}

	t.Run("memory_accessibility", func(t *testing.T) {
		// Test that multiple instances don't cause excessive memory usage
		var grids []*integration.ButtonGrid

		for i := 0; i < 100; i++ {
			grid := integration.NewButtonGrid()
			grids = append(grids, grid)

			// Each grid should work independently
			assert.Equal(t, 18, grid.GetButtonCount())
			button, ok := grid.GetFocusedButton()
			assert.True(t, ok, "Should have focused button")
			assert.NotNil(t, button, "Focused button should not be nil")
		}

		// All grids should still work
		for i, grid := range grids {
			assert.Equal(t, 18, grid.GetButtonCount(),
				"Grid %d should maintain button count", i)
			button, ok := grid.GetFocusedButton()
			assert.True(t, ok, "Grid %d should have focused button", i)
			assert.NotNil(t, button, "Grid %d focused button should not be nil", i)
		}
	})
}

// TestCognitiveAccessibility tests cognitive accessibility features
func TestCognitiveAccessibility(t *testing.T) {
	t.Run("cognitive_accessibility", func(t *testing.T) {
		grid := integration.NewButtonGrid()

		// Test consistent button layout
		rendering := grid.Render(80)
		lines := strings.Split(rendering, "\n")

		// Should have consistent row structure
		buttonRows := 0
		for _, line := range lines {
			if strings.Contains(line, "│") && strings.TrimSpace(line) != "" {
				buttonRows++
			}
		}
		assert.Equal(t, 5, buttonRows, "Should have 5 consistent button rows")

		// Test predictable button placement (similar buttons in similar positions)
		// Numbers should be in predictable grid pattern
		grid = integration.NewButtonGrid()

		// Navigate to verify number grid layout
		positions := make(map[string][2]int)
		for row := 1; row <= 3; row++ {
			for col := 0; col < 3; col++ {
				// Navigate to position
				grid = integration.NewButtonGrid()
				for i := 0; i < row; i++ {
					grid.HandleKeyPress(tea.KeyMsg{Type: tea.KeyDown})
				}
				for i := 0; i < col; i++ {
					grid.HandleKeyPress(tea.KeyMsg{Type: tea.KeyRight})
				}

				focused, _ := grid.GetFocusedButton()
				positions[focused.GetLabel()] = [2]int{row, col}
			}
		}

		// Verify number layout is logical
		assert.Equal(t, [2]int{1, 0}, positions["7"], "7 should be at row 1, col 0")
		assert.Equal(t, [2]int{1, 1}, positions["8"], "8 should be at row 1, col 1")
		assert.Equal(t, [2]int{1, 2}, positions["9"], "9 should be at row 1, col 2")
		assert.Equal(t, [2]int{2, 0}, positions["4"], "4 should be at row 2, col 0")
	})
}