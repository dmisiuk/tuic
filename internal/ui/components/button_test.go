package components

import (
	"fmt"
	"strings"
	"testing"

	"github.com/charmbracelet/lipgloss"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestButtonState_String(t *testing.T) {
	tests := []struct {
		name     string
		state    ButtonState
		expected string
	}{
		{"Normal state", StateNormal, "normal"},
		{"Focused state", StateFocused, "focused"},
		{"Pressed state", StatePressed, "pressed"},
		{"Disabled state", StateDisabled, "disabled"},
		{"Unknown state", ButtonState(99), "unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.state.String()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestButtonType_String(t *testing.T) {
	tests := []struct {
		name     string
		btnType  ButtonType
		expected string
	}{
		{"Number type", TypeNumber, "number"},
		{"Operator type", TypeOperator, "operator"},
		{"Special type", TypeSpecial, "special"},
		{"Unknown type", ButtonType(99), "unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.btnType.String()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestNewButtonStateManager(t *testing.T) {
	config := ButtonConfig{
		Label:    "Test",
		Type:     TypeNumber,
		Value:    "1",
		Width:    10,
		Height:   3,
		Position: Position{Row: 0, Column: 0},
	}

	sm := NewButtonStateManager(TypeNumber, config)

	assert.NotNil(t, sm)
	assert.Equal(t, StateNormal, sm.State())
	assert.Equal(t, TypeNumber, sm.GetType())
	assert.Equal(t, config, sm.GetConfig())
}

func TestButtonStateManager_StateTransitions(t *testing.T) {
	config := ButtonConfig{
		Label: "Test",
		Type:  TypeNumber,
		Value: "1",
	}
	sm := NewButtonStateManager(TypeNumber, config)

	// Test initial state
	assert.Equal(t, StateNormal, sm.State())

	// Test Focus transition
	err := sm.Focus()
	assert.NoError(t, err)
	assert.Equal(t, StateFocused, sm.State())

	// Test Press transition from focused
	err = sm.Press()
	assert.NoError(t, err)
	assert.Equal(t, StatePressed, sm.State())

	// Test Release transition
	err = sm.Release()
	assert.NoError(t, err)
	assert.Equal(t, StateNormal, sm.State())

	// Test Blur from focused
	err = sm.Focus()
	assert.NoError(t, err)
	err = sm.Blur()
	assert.NoError(t, err)
	assert.Equal(t, StateNormal, sm.State())

	// Test Disable transition
	err = sm.Disable()
	assert.NoError(t, err)
	assert.Equal(t, StateDisabled, sm.State())

	// Test Enable transition
	err = sm.Enable()
	assert.NoError(t, err)
	assert.Equal(t, StateNormal, sm.State())
}

func TestButtonStateManager_InvalidTransitions(t *testing.T) {
	config := ButtonConfig{
		Label: "Test",
		Type:  TypeNumber,
		Value: "1",
	}
	sm := NewButtonStateManager(TypeNumber, config)

	// Test invalid transition from disabled to pressed
	err := sm.Disable()
	assert.NoError(t, err)
	err = sm.Press()
	assert.Error(t, err)
	assert.IsType(t, &InvalidStateTransitionError{}, err)

	// Test the error message
	invalidErr := err.(*InvalidStateTransitionError)
	assert.Equal(t, StateDisabled, invalidErr.From)
	assert.Equal(t, StatePressed, invalidErr.To)
	assert.Contains(t, err.Error(), "invalid state transition from disabled to pressed")
}

func TestButtonStateManager_SetState(t *testing.T) {
	config := ButtonConfig{
		Label: "Test",
		Type:  TypeNumber,
		Value: "1",
	}
	sm := NewButtonStateManager(TypeNumber, config)

	// Test valid state change
	err := sm.SetState(StateFocused)
	assert.NoError(t, err)
	assert.Equal(t, StateFocused, sm.State())

	// Test valid state change from focused to pressed
	err = sm.SetState(StatePressed)
	assert.NoError(t, err)
	assert.Equal(t, StatePressed, sm.State())

	// Test invalid state change from disabled to pressed
	err = sm.SetState(StateDisabled)
	assert.NoError(t, err)
	err = sm.SetState(StatePressed)
	assert.Error(t, err)
	assert.Equal(t, StateDisabled, sm.State()) // State should not change on error
}

func TestButtonStateManager_HelperMethods(t *testing.T) {
	config := ButtonConfig{
		Label: "Test",
		Type:  TypeNumber,
		Value: "1",
	}
	sm := NewButtonStateManager(TypeNumber, config)

	// Test initial state
	assert.True(t, sm.IsInteractive())
	assert.False(t, sm.IsFocused())
	assert.False(t, sm.IsPressed())

	// Test focused state
	err := sm.Focus()
	assert.NoError(t, err)
	assert.True(t, sm.IsInteractive())
	assert.True(t, sm.IsFocused())
	assert.False(t, sm.IsPressed())

	// Test pressed state
	err = sm.Press()
	assert.NoError(t, err)
	assert.True(t, sm.IsInteractive())
	assert.False(t, sm.IsFocused())
	assert.True(t, sm.IsPressed())

	// Test disabled state
	err = sm.SetState(StateDisabled)
	assert.NoError(t, err)
	assert.False(t, sm.IsInteractive())
	assert.False(t, sm.IsFocused())
	assert.False(t, sm.IsPressed())
}

func TestNewButton(t *testing.T) {
	config := ButtonConfig{
		Label:    "7",
		Type:     TypeNumber,
		Value:    "7",
		Width:    5,
		Height:   3,
		Position: Position{Row: 2, Column: 1},
	}

	button := NewButton(config)

	require.NotNil(t, button)
	assert.Equal(t, config.Label, button.GetLabel())
	assert.Equal(t, config.Type, button.GetType())
	assert.Equal(t, config.Value, button.GetValue())
	assert.Equal(t, config.Position, button.GetPosition())
	assert.Equal(t, StateNormal, button.GetState())
	assert.True(t, button.IsInteractive())
}

func TestNewButtonWithTheme(t *testing.T) {
	config := ButtonConfig{
		Label:    "+",
		Type:     TypeOperator,
		Value:    "+",
		Width:    5,
		Height:   3,
		Position: Position{Row: 0, Column: 3},
	}

	customTheme := ButtonTheme{
		Operator: ButtonTypeStyle{
			Normal: lipgloss.NewStyle().
				Foreground(lipgloss.Color("255")).
				Background(lipgloss.Color("100")),
		},
	}

	button := NewButtonWithTheme(config, customTheme)

	require.NotNil(t, button)
	assert.Equal(t, config.Label, button.GetLabel())
	assert.Equal(t, config.Type, button.GetType())
	assert.Equal(t, config.Value, button.GetValue())
}

func TestButton_StateManagement(t *testing.T) {
	config := ButtonConfig{
		Label: "C",
		Type:  TypeSpecial,
		Value: "C",
	}

	button := NewButton(config)

	// Test initial state
	assert.Equal(t, StateNormal, button.GetState())
	assert.True(t, button.IsInteractive())

	// Test focus
	err := button.Focus()
	assert.NoError(t, err)
	assert.Equal(t, StateFocused, button.GetState())
	assert.True(t, button.IsFocused())

	// Test press
	err = button.Press()
	assert.NoError(t, err)
	assert.Equal(t, StatePressed, button.GetState())
	assert.True(t, button.IsPressed())

	// Test release
	err = button.Release()
	assert.NoError(t, err)
	assert.Equal(t, StateNormal, button.GetState())

	// Test disable
	err = button.Disable()
	assert.NoError(t, err)
	assert.Equal(t, StateDisabled, button.GetState())
	assert.False(t, button.IsInteractive())

	// Test enable
	err = button.Enable()
	assert.NoError(t, err)
	assert.Equal(t, StateNormal, button.GetState())
	assert.True(t, button.IsInteractive())
}

func TestButton_Render(t *testing.T) {
	config := ButtonConfig{
		Label:  "9",
		Type:   TypeNumber,
		Value:  "9",
		Width:  10,
		Height: 3,
	}

	button := NewButton(config)

	// Test normal state rendering
	rendered := button.Render()
	assert.Contains(t, rendered, "9")
	assert.True(t, strings.HasPrefix(rendered, " "))
	assert.True(t, strings.HasSuffix(rendered, " "))

	// Test focused state rendering
	err := button.Focus()
	assert.NoError(t, err)
	rendered = button.Render()
	assert.Contains(t, rendered, "9")

	// Test pressed state rendering
	err = button.Press()
	assert.NoError(t, err)
	rendered = button.Render()
	assert.Contains(t, rendered, "9")

	// Test disabled state rendering
	err = button.Disable()
	assert.NoError(t, err)
	rendered = button.Render()
	assert.Contains(t, rendered, "9")
}

func TestButton_String(t *testing.T) {
	config := ButtonConfig{
		Label: "=",
		Type:  TypeSpecial,
		Value: "=",
	}

	button := NewButton(config)
	buttonStr := button.String()

	assert.Contains(t, buttonStr, "Button{")
	assert.Contains(t, buttonStr, "Label: =")
	assert.Contains(t, buttonStr, "Type: special")
	assert.Contains(t, buttonStr, "State: normal")
	assert.Contains(t, buttonStr, "Value: =")
}

func TestButtonAction(t *testing.T) {
	config := ButtonConfig{
		Label: "×",
		Type:  TypeOperator,
		Value: "*",
	}

	button := NewButton(config)

	// Test creating action
	action := button.Trigger("multiply")
	require.NotNil(t, action)
	assert.Equal(t, button, action.Button)
	assert.Equal(t, "multiply", action.Type)
	assert.Equal(t, "*", action.Value)
	assert.Nil(t, action.Context)

	// Test with context
	actionWithContext := action.WithContext("calculator")
	assert.Equal(t, "calculator", actionWithContext.Context)
}

func TestNewButtonAction(t *testing.T) {
	config := ButtonConfig{
		Label: "CE",
		Type:  TypeSpecial,
		Value: "CE",
	}

	button := NewButton(config)
	action := NewButtonAction(button, "clear_entry")

	assert.Equal(t, button, action.Button)
	assert.Equal(t, "clear_entry", action.Type)
	assert.Equal(t, "CE", action.Value)
}

func TestButtonRenderer_Render(t *testing.T) {
	config := ButtonConfig{
		Label:  "5",
		Type:   TypeNumber,
		Value:  "5",
		Width:  8,
		Height: 2,
	}

	button := NewButton(config)
	renderer := NewButtonRenderer(DefaultButtonTheme())

	// Test rendering different states
	for _, state := range testStates {
		t.Run(state.String(), func(t *testing.T) {
			// Set button state
			button.stateManager.currentState = state

			// Render with renderer
			rendered := renderer.Render(button)
			assert.Contains(t, rendered, "5")
			assert.True(t, len(rendered) > 0)
		})
	}
}

func TestButtonTypeStyle_getStyleForState(t *testing.T) {
	style := ButtonTypeStyle{
		Normal:   lipgloss.NewStyle().Background(lipgloss.Color("240")),
		Focused:  lipgloss.NewStyle().Background(lipgloss.Color("62")),
		Pressed:  lipgloss.NewStyle().Background(lipgloss.Color("94")),
		Disabled: lipgloss.NewStyle().Background(lipgloss.Color("240")),
	}

	tests := []struct {
		name     string
		state    ButtonState
		expected lipgloss.Color
	}{
		{"Normal state", StateNormal, lipgloss.Color("240")},
		{"Focused state", StateFocused, lipgloss.Color("62")},
		{"Pressed state", StatePressed, lipgloss.Color("94")},
		{"Disabled state", StateDisabled, lipgloss.Color("240")},
		{"Unknown state", ButtonState(99), lipgloss.Color("240")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := style.getStyleForState(tt.state)
			assert.Equal(t, tt.expected, result.GetBackground())
		})
	}
}

func TestDefaultButtonTheme(t *testing.T) {
	theme := DefaultButtonTheme()

	// Test that all button types have styles defined
	assert.NotNil(t, theme.Number)
	assert.NotNil(t, theme.Operator)
	assert.NotNil(t, theme.Special)

	// Test that all states have styles for each type
	states := []ButtonState{StateNormal, StateFocused, StatePressed, StateDisabled}

	for _, state := range states {
		t.Run(fmt.Sprintf("Number_%s", state.String()), func(t *testing.T) {
			style := theme.Number.getStyleForState(state)
			assert.NotNil(t, style)
		})

		t.Run(fmt.Sprintf("Operator_%s", state.String()), func(t *testing.T) {
			style := theme.Operator.getStyleForState(state)
			assert.NotNil(t, style)
		})

		t.Run(fmt.Sprintf("Special_%s", state.String()), func(t *testing.T) {
			style := theme.Special.getStyleForState(state)
			assert.NotNil(t, style)
		})
	}
}

func TestButton_RenderWithCustomDimensions(t *testing.T) {
	config := ButtonConfig{
		Label:  "0",
		Type:   TypeNumber,
		Value:  "0",
		Width:  15,
		Height: 4,
	}

	button := NewButton(config)
	rendered := button.Render()

	// Check that the rendered button respects dimensions
	// The exact length depends on styling and borders, so we just check it's reasonable
	assert.True(t, len(rendered) >= config.Width)
	assert.True(t, strings.Contains(rendered, "0"))
}

func TestInvalidStateTransitionError(t *testing.T) {
	err := &InvalidStateTransitionError{
		From: StateNormal,
		To:   StateDisabled,
	}

	assert.Equal(t, "invalid state transition from normal to disabled", err.Error())
	assert.True(t, err.Is(&InvalidStateTransitionError{}))

	// Test with a different error type
	differentErr := &InvalidStateTransitionError{From: StateFocused, To: StateNormal}
	assert.True(t, err.Is(differentErr)) // Both are InvalidStateTransitionError types

	// Test with non-matching error type
	var generalErr error = &InvalidStateTransitionError{}
	assert.True(t, err.Is(generalErr))
}

// testStates is a helper variable for testing all button states
var testStates = []ButtonState{StateNormal, StateFocused, StatePressed, StateDisabled}

// Benchmark tests for performance
func BenchmarkButton_Render(b *testing.B) {
	config := ButtonConfig{
		Label: "1",
		Type:  TypeNumber,
		Value: "1",
	}
	button := NewButton(config)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = button.Render()
	}
}

func BenchmarkButton_StateTransition(b *testing.B) {
	config := ButtonConfig{
		Label: "1",
		Type:  TypeNumber,
		Value: "1",
	}
	button := NewButton(config)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = button.Focus()
		_ = button.Press()
		_ = button.Release()
	}
}