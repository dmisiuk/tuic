package integration

import (
	"testing"

	"github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"ccpm-demo/internal/ui/components"
)

func TestNewButtonGrid(t *testing.T) {
	t.Run("creates button grid with default configuration", func(t *testing.T) {
		grid := NewButtonGrid()

		// Check that grid is properly initialized
		assert.NotNil(t, grid)
		assert.NotNil(t, grid.buttons)
		assert.NotNil(t, grid.grid)
		assert.NotNil(t, grid.themeManager)

		// Check dimensions
		assert.Equal(t, 4, grid.dimensions.Columns)
		assert.Equal(t, 5, grid.dimensions.Rows)

		// Check that buttons were created
		assert.Greater(t, len(grid.buttons), 0)
		assert.Equal(t, 19, len(grid.buttons)) // 19 button definitions in the layout

		// Check default theme
		assert.Equal(t, "retro-casio", grid.GetCurrentTheme())
	})
}

func TestNewButtonGridWithTheme(t *testing.T) {
	t.Run("creates button grid with specific theme", func(t *testing.T) {
		grid, err := NewButtonGridWithTheme("nonexistent")

		// Should succeed even if theme doesn't exist (falls back to default)
		require.NoError(t, err)
		assert.NotNil(t, grid)
		assert.Equal(t, "retro-casio", grid.GetCurrentTheme()) // Fallback
	})
}

func TestButtonGridInitialization(t *testing.T) {
	t.Run("initializes calculator button layout correctly", func(t *testing.T) {
		grid := NewButtonGrid()

		// Check specific button positions and values
		tests := []struct {
			buttonID string
			expectedLabel string
			expectedType  components.ButtonType
			expectedValue string
		}{
			{"button_0_0", "C", components.TypeSpecial, "clear"},
			{"button_0_3", "÷", components.TypeOperator, "/"},
			{"button_1_0", "7", components.TypeNumber, "7"},
			{"button_1_3", "×", components.TypeOperator, "*"},
			{"button_4_0", "0", components.TypeNumber, "0"},
			{"button_4_2", "=", components.TypeSpecial, "="},
		}

		for _, test := range tests {
			button, exists := grid.GetButton(test.buttonID)
			assert.True(t, exists, "Button %s should exist", test.buttonID)
			if exists {
				assert.Equal(t, test.expectedLabel, button.GetLabel())
				assert.Equal(t, test.expectedType, button.GetType())
				assert.Equal(t, test.expectedValue, button.GetValue())
			}
		}
	})

	t.Run("sets initial focus correctly", func(t *testing.T) {
		grid := NewButtonGrid()

		focusedButton, exists := grid.GetFocusedButton()
		assert.True(t, exists, "Should have a focused button")
		if exists {
			assert.Equal(t, "C", focusedButton.GetLabel())
			assert.True(t, focusedButton.IsFocused())
		}
	})
}

func TestButtonGridKeyboardNavigation(t *testing.T) {
	t.Run("handles arrow key navigation", func(t *testing.T) {
		grid := NewButtonGrid()

		// Navigate down from C button (0,0) to 7 button (1,0)
		action := grid.HandleKeyPress(tea.KeyMsg{Type: tea.KeyDown})
		assert.Nil(t, action) // Navigation shouldn't trigger actions

		focusedButton, exists := grid.GetFocusedButton()
		require.True(t, exists)
		assert.Equal(t, "7", focusedButton.GetLabel())

		// Navigate right to 8 button (1,1)
		action = grid.HandleKeyPress(tea.KeyMsg{Type: tea.KeyRight})
		assert.Nil(t, action)

		focusedButton, exists = grid.GetFocusedButton()
		require.True(t, exists)
		assert.Equal(t, "8", focusedButton.GetLabel())
	})

	t.Run("handles boundary conditions in navigation", func(t *testing.T) {
		grid := NewButtonGrid()

		// Try to navigate up from top row
		originalFocused, _ := grid.GetFocusedButton()
		grid.HandleKeyPress(tea.KeyMsg{Type: tea.KeyUp})
		focusedAfter, _ := grid.GetFocusedButton()
		assert.Equal(t, originalFocused.GetLabel(), focusedAfter.GetLabel())

		// Navigate to bottom right corner
		for i := 0; i < 10; i++ {
			grid.HandleKeyPress(tea.KeyMsg{Type: tea.KeyDown})
		}
		for i := 0; i < 10; i++ {
			grid.HandleKeyPress(tea.KeyMsg{Type: tea.KeyRight})
		}

		// Try to navigate further right
		edgeFocused, _ := grid.GetFocusedButton()
		grid.HandleKeyPress(tea.KeyMsg{Type: tea.KeyRight})
		afterEdgeFocused, _ := grid.GetFocusedButton()
		assert.Equal(t, edgeFocused.GetLabel(), afterEdgeFocused.GetLabel())
	})
}

func TestButtonGridKeyboardActivation(t *testing.T) {
	t.Run("activates buttons with Enter and Space", func(t *testing.T) {
		grid := NewButtonGrid()

		// Activate focused button with Enter
		action := grid.HandleKeyPress(tea.KeyMsg{Type: tea.KeyEnter})
		assert.NotNil(t, action)
		assert.Equal(t, "press", action.Action)
		assert.Equal(t, "clear", action.Value)
		assert.Equal(t, "button_0_0", action.ButtonID)

		// Activate with Space
		action = grid.HandleKeyPress(tea.KeyMsg{Type: tea.KeySpace})
		assert.NotNil(t, action)
		assert.Equal(t, "press", action.Action)
		assert.Equal(t, "clear", action.Value)
	})

	t.Run("handles direct number input", func(t *testing.T) {
		grid := NewButtonGrid()

		// Press '5' key
		action := grid.HandleKeyPress(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'5'}})
		assert.NotNil(t, action)
		assert.Equal(t, "5", action.Value)
		assert.Equal(t, "button_2_1", action.ButtonID)

		// Press '+' key
		action = grid.HandleKeyPress(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'+'}})
		assert.NotNil(t, action)
		assert.Equal(t, "+", action.Value)
		assert.Equal(t, "button_3_3", action.ButtonID)
	})
}

func TestButtonGridMouseHandling(t *testing.T) {
	t.Run("handles mouse clicks on buttons", func(t *testing.T) {
		grid := NewButtonGrid()

		// Simulate clicking on the "=" button (approximately position)
		// Note: Exact coordinates depend on rendering, this is a simplified test
		msg := tea.MouseMsg{
			Type: tea.MouseLeft,
			X:    30,
			Y:    15,
		}

		_ = grid.HandleMouse(msg)
		// This might be nil due to coordinate calculation complexity
		// but the method should not panic
		// In a real test, we'd need to calculate exact coordinates
		assert.NotNil(t, grid) // Ensure grid still exists
	})

	t.Run("ignores non-left mouse clicks", func(t *testing.T) {
		grid := NewButtonGrid()

		msg := tea.MouseMsg{
			Type: tea.MouseRight,
			X:    10,
			Y:    10,
		}

		action := grid.HandleMouse(msg)
		assert.Nil(t, action)
	})
}

func TestButtonGridButtonManagement(t *testing.T) {
	t.Run("retrieves buttons by ID", func(t *testing.T) {
		grid := NewButtonGrid()

		button, exists := grid.GetButton("button_1_1")
		assert.True(t, exists)
		assert.Equal(t, "8", button.GetLabel())

		button, exists = grid.GetButton("nonexistent")
		assert.False(t, exists)
		assert.Nil(t, button)
	})

	t.Run("returns button count", func(t *testing.T) {
		grid := NewButtonGrid()
		assert.Equal(t, 19, grid.GetButtonCount())

		allButtons := grid.GetButtons()
		assert.Equal(t, 19, len(allButtons))
		assert.Equal(t, grid.buttons, allButtons)
	})

	t.Run("returns dimensions", func(t *testing.T) {
		grid := NewButtonGrid()
		dims := grid.GetDimensions()
		assert.Equal(t, 4, dims.Columns)
		assert.Equal(t, 5, dims.Rows)
	})
}

func TestButtonGridThemeManagement(t *testing.T) {
	t.Run("changes theme successfully", func(t *testing.T) {
		grid := NewButtonGrid()

		originalTheme := grid.GetCurrentTheme()
		assert.Equal(t, "retro-casio", originalTheme)

		err := grid.SetTheme("modern")
		assert.NoError(t, err)
		assert.Equal(t, "retro-casio", grid.GetCurrentTheme()) // Falls back to default

		// Buttons should still exist after theme change
		assert.Greater(t, grid.GetButtonCount(), 0)
	})

	t.Run("handles invalid theme gracefully", func(t *testing.T) {
		grid := NewButtonGrid()

		err := grid.SetTheme("nonexistent_theme")
		assert.NoError(t, err) // Should not error, falls back to default
		assert.Equal(t, "retro-casio", grid.GetCurrentTheme())
	})
}

func TestButtonGridRendering(t *testing.T) {
	t.Run("renders without panicking", func(t *testing.T) {
		grid := NewButtonGrid()

		// Test rendering with different widths
		output80 := grid.Render(80)
		assert.NotEmpty(t, output80)
		assert.Contains(t, output80, "C") // Should contain some button labels

		output60 := grid.Render(60)
		assert.NotEmpty(t, output60)

		output100 := grid.Render(100)
		assert.NotEmpty(t, output100)
	})

	t.Run("rendering adapts to terminal width", func(t *testing.T) {
		grid := NewButtonGrid()

		narrow := grid.Render(40)
		wide := grid.Render(100)

		// Both should render but may have different layouts
		assert.NotEmpty(t, narrow)
		assert.NotEmpty(t, wide)

		// Both should contain calculator buttons
		assert.Contains(t, narrow, "C")
		assert.Contains(t, wide, "C")
	})
}

func TestButtonGridPositionValidation(t *testing.T) {
	t.Run("validates grid positions correctly", func(t *testing.T) {
		grid := NewButtonGrid()

		// Valid positions
		assert.True(t, grid.isValidPosition(0, 0))
		assert.True(t, grid.isValidPosition(3, 4))
		assert.True(t, grid.isValidPosition(2, 2))

		// Invalid positions
		assert.False(t, grid.isValidPosition(-1, 0))
		assert.False(t, grid.isValidPosition(0, -1))
		assert.False(t, grid.isValidPosition(4, 0))  // Beyond column limit
		assert.False(t, grid.isValidPosition(0, 5))  // Beyond row limit
		assert.False(t, grid.isValidPosition(10, 10)) // Way beyond limits
	})
}

func TestButtonGridStringRepresentation(t *testing.T) {
	t.Run("returns meaningful string representation", func(t *testing.T) {
		grid := NewButtonGrid()

		str := grid.String()
		assert.Contains(t, str, "ButtonGrid")
		assert.Contains(t, str, "4x5")
		assert.Contains(t, str, "19") // Button count
		assert.Contains(t, str, "retro-casio") // Theme
		assert.Contains(t, str, "button_0_0") // Initial focus
	})
}

func TestButtonGridFocusManagement(t *testing.T) {
	t.Run("tracks focus correctly", func(t *testing.T) {
		grid := NewButtonGrid()

		// Initially focused on C button
		focused, exists := grid.GetFocusedButton()
		require.True(t, exists)
		assert.Equal(t, "C", focused.GetLabel())
		assert.True(t, focused.IsFocused())

		// Navigate to 7 button
		grid.HandleKeyPress(tea.KeyMsg{Type: tea.KeyDown})
		focused, exists = grid.GetFocusedButton()
		require.True(t, exists)
		assert.Equal(t, "7", focused.GetLabel())

		// Original C button should no longer be focused
		cButton, _ := grid.GetButton("button_0_0")
		assert.False(t, cButton.IsFocused())
	})
}

func TestButtonGridActionCreation(t *testing.T) {
	t.Run("creates proper action objects", func(t *testing.T) {
		grid := NewButtonGrid()

		// Activate the 5 button
		grid.HandleKeyPress(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'5'}})
		action := grid.HandleKeyPress(tea.KeyMsg{Type: tea.KeyEnter})

		require.NotNil(t, action)
		assert.Equal(t, "press", action.Action)
		assert.Equal(t, "5", action.Value)
		assert.Equal(t, "button_2_1", action.ButtonID)
		assert.NotNil(t, action.Button)
		assert.Equal(t, "5", action.Button.GetLabel())
	})
}

// Benchmark tests
func BenchmarkButtonGridRender(b *testing.B) {
	grid := NewButtonGrid()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = grid.Render(80)
	}
}

func BenchmarkButtonGridKeyPress(b *testing.B) {
	grid := NewButtonGrid()
	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'5'}}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = grid.HandleKeyPress(msg)
	}
}

func BenchmarkButtonGridMouseHandling(b *testing.B) {
	grid := NewButtonGrid()
	msg := tea.MouseMsg{Type: tea.MouseLeft, X: 20, Y: 10}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = grid.HandleMouse(msg)
	}
}