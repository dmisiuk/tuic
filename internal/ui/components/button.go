package components

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

// Button represents an interactive button component with state management
type Button struct {
	stateManager *ButtonStateManager
	styles       *lipgloss.Style
	theme        ButtonTheme
}

// ButtonTheme defines the styling theme for different button types
type ButtonTheme struct {
	Number   ButtonTypeStyle
	Operator ButtonTypeStyle
	Special  ButtonTypeStyle
}

// ButtonTypeStyle defines styling for a specific button type across all states
type ButtonTypeStyle struct {
	Normal   lipgloss.Style
	Focused  lipgloss.Style
	Pressed  lipgloss.Style
	Disabled lipgloss.Style
}

// ButtonRenderer handles the visual rendering of buttons
type ButtonRenderer struct {
	theme ButtonTheme
}

// NewButton creates a new button component with the specified configuration
func NewButton(config ButtonConfig) *Button {
	stateManager := NewButtonStateManager(config.Type, config)

	return &Button{
		stateManager: stateManager,
		styles:       &lipgloss.Style{},
		theme:        DefaultButtonTheme(),
	}
}

// NewButtonWithTheme creates a new button with a custom theme
func NewButtonWithTheme(config ButtonConfig, theme ButtonTheme) *Button {
	stateManager := NewButtonStateManager(config.Type, config)

	return &Button{
		stateManager: stateManager,
		styles:       &lipgloss.Style{},
		theme:        theme,
	}
}

// NewButtonRenderer creates a new button renderer with the specified theme
func NewButtonRenderer(theme ButtonTheme) *ButtonRenderer {
	return &ButtonRenderer{
		theme: theme,
	}
}

// GetConfig returns the button configuration
func (b *Button) GetConfig() ButtonConfig {
	return b.stateManager.GetConfig()
}

// GetState returns the current button state
func (b *Button) GetState() ButtonState {
	return b.stateManager.State()
}

// GetType returns the button type
func (b *Button) GetType() ButtonType {
	return b.stateManager.GetType()
}

// GetValue returns the button's value
func (b *Button) GetValue() string {
	return b.GetConfig().Value
}

// GetLabel returns the button's label
func (b *Button) GetLabel() string {
	return b.GetConfig().Label
}

// GetPosition returns the button's grid position
func (b *Button) GetPosition() Position {
	return b.GetConfig().Position
}

// Focus sets the button to focused state
func (b *Button) Focus() error {
	return b.stateManager.Focus()
}

// Press sets the button to pressed state
func (b *Button) Press() error {
	return b.stateManager.Press()
}

// Release returns the button to normal state from pressed
func (b *Button) Release() error {
	return b.stateManager.Release()
}

// Blur returns the button to normal state from focused
func (b *Button) Blur() error {
	return b.stateManager.Blur()
}

// Disable sets the button to disabled state
func (b *Button) Disable() error {
	return b.stateManager.Disable()
}

// Enable returns the button to normal state from disabled
func (b *Button) Enable() error {
	return b.stateManager.Enable()
}

// IsInteractive returns true if the button is in an interactive state
func (b *Button) IsInteractive() bool {
	return b.stateManager.IsInteractive()
}

// IsFocused returns true if the button has focus
func (b *Button) IsFocused() bool {
	return b.stateManager.IsFocused()
}

// IsPressed returns true if the button is being pressed
func (b *Button) IsPressed() bool {
	return b.stateManager.IsPressed()
}

// Render returns the styled string representation of the button
func (b *Button) Render() string {
	config := b.GetConfig()
	state := b.GetState()
	buttonType := b.GetType()

	var style lipgloss.Style

	// Get the appropriate style based on button type and state
	switch buttonType {
	case TypeNumber:
		style = b.getNumberStyle(state)
	case TypeOperator:
		style = b.getOperatorStyle(state)
	case TypeSpecial:
		style = b.getSpecialStyle(state)
	default:
		style = b.getDefaultStyle(state)
	}

	// Apply width and height from config
	if config.Width > 0 {
		style = style.Width(config.Width)
	}
	if config.Height > 0 {
		style = style.Height(config.Height)
	}

	return style.Render(config.Label)
}

// RenderWithRenderer renders a button using an external renderer
func (br *ButtonRenderer) Render(button *Button) string {
	config := button.GetConfig()
	state := button.GetState()
	buttonType := button.GetType()

	var style lipgloss.Style

	// Get the appropriate style based on button type and state
	switch buttonType {
	case TypeNumber:
		style = br.theme.Number.getStyleForState(state)
	case TypeOperator:
		style = br.theme.Operator.getStyleForState(state)
	case TypeSpecial:
		style = br.theme.Special.getStyleForState(state)
	default:
		style = br.getDefaultStyle(state)
	}

	// Apply width and height from config
	if config.Width > 0 {
		style = style.Width(config.Width)
	}
	if config.Height > 0 {
		style = style.Height(config.Height)
	}

	return style.Render(config.Label)
}

// getNumberStyle returns the style for number buttons based on state
func (b *Button) getNumberStyle(state ButtonState) lipgloss.Style {
	return b.theme.Number.getStyleForState(state)
}

// getOperatorStyle returns the style for operator buttons based on state
func (b *Button) getOperatorStyle(state ButtonState) lipgloss.Style {
	return b.theme.Operator.getStyleForState(state)
}

// getSpecialStyle returns the style for special buttons based on state
func (b *Button) getSpecialStyle(state ButtonState) lipgloss.Style {
	return b.theme.Special.getStyleForState(state)
}

// getDefaultStyle returns a default style for unknown button types
func (b *Button) getDefaultStyle(state ButtonState) lipgloss.Style {
	switch state {
	case StateNormal:
		return lipgloss.NewStyle().
			Foreground(lipgloss.Color("15")).
			Background(lipgloss.Color("240")).
			Align(lipgloss.Center)
	case StateFocused:
		return lipgloss.NewStyle().
			Foreground(lipgloss.Color("15")).
			Background(lipgloss.Color("62")).
			Align(lipgloss.Center)
	case StatePressed:
		return lipgloss.NewStyle().
			Foreground(lipgloss.Color("15")).
			Background(lipgloss.Color("94")).
			Align(lipgloss.Center)
	case StateDisabled:
		return lipgloss.NewStyle().
			Foreground(lipgloss.Color("8")).
			Background(lipgloss.Color("240")).
			Align(lipgloss.Center)
	default:
		return lipgloss.NewStyle().Align(lipgloss.Center)
	}
}

// getDefaultStyle returns a default style for the renderer
func (br *ButtonRenderer) getDefaultStyle(state ButtonState) lipgloss.Style {
	switch state {
	case StateNormal:
		return lipgloss.NewStyle().
			Foreground(lipgloss.Color("15")).
			Background(lipgloss.Color("240")).
			Align(lipgloss.Center)
	case StateFocused:
		return lipgloss.NewStyle().
			Foreground(lipgloss.Color("15")).
			Background(lipgloss.Color("62")).
			Align(lipgloss.Center)
	case StatePressed:
		return lipgloss.NewStyle().
			Foreground(lipgloss.Color("15")).
			Background(lipgloss.Color("94")).
			Align(lipgloss.Center)
	case StateDisabled:
		return lipgloss.NewStyle().
			Foreground(lipgloss.Color("8")).
			Background(lipgloss.Color("240")).
			Align(lipgloss.Center)
	default:
		return lipgloss.NewStyle().Align(lipgloss.Center)
	}
}

// getStyleForState returns the appropriate style for a given state
func (bts ButtonTypeStyle) getStyleForState(state ButtonState) lipgloss.Style {
	switch state {
	case StateNormal:
		return bts.Normal
	case StateFocused:
		return bts.Focused
	case StatePressed:
		return bts.Pressed
	case StateDisabled:
		return bts.Disabled
	default:
		return bts.Normal
	}
}

// DefaultButtonTheme returns the default retro Casio-inspired button theme
func DefaultButtonTheme() ButtonTheme {
	return ButtonTheme{
		Number: ButtonTypeStyle{
			Normal: lipgloss.NewStyle().
				Foreground(lipgloss.Color("15")).
				Background(lipgloss.Color("240")).
				Align(lipgloss.Center).
				Border(lipgloss.NormalBorder(), false).
				BorderForeground(lipgloss.Color("244")),
			Focused: lipgloss.NewStyle().
				Foreground(lipgloss.Color("15")).
				Background(lipgloss.Color("246")).
				Align(lipgloss.Center).
				Border(lipgloss.NormalBorder(), false).
				BorderForeground(lipgloss.Color("62")),
			Pressed: lipgloss.NewStyle().
				Foreground(lipgloss.Color("15")).
				Background(lipgloss.Color("248")).
				Align(lipgloss.Center).
				Border(lipgloss.NormalBorder(), false).
				BorderForeground(lipgloss.Color("94")),
			Disabled: lipgloss.NewStyle().
				Foreground(lipgloss.Color("8")).
				Background(lipgloss.Color("240")).
				Align(lipgloss.Center).
				Border(lipgloss.NormalBorder(), false).
				BorderForeground(lipgloss.Color("244")),
		},
		Operator: ButtonTypeStyle{
			Normal: lipgloss.NewStyle().
				Foreground(lipgloss.Color("15")).
				Background(lipgloss.Color("208")).
				Align(lipgloss.Center).
				Border(lipgloss.NormalBorder(), false).
				BorderForeground(lipgloss.Color("202")),
			Focused: lipgloss.NewStyle().
				Foreground(lipgloss.Color("15")).
				Background(lipgloss.Color("214")).
				Align(lipgloss.Center).
				Border(lipgloss.NormalBorder(), false).
				BorderForeground(lipgloss.Color("62")),
			Pressed: lipgloss.NewStyle().
				Foreground(lipgloss.Color("15")).
				Background(lipgloss.Color("220")).
				Align(lipgloss.Center).
				Border(lipgloss.NormalBorder(), false).
				BorderForeground(lipgloss.Color("94")),
			Disabled: lipgloss.NewStyle().
				Foreground(lipgloss.Color("8")).
				Background(lipgloss.Color("208")).
				Align(lipgloss.Center).
				Border(lipgloss.NormalBorder(), false).
				BorderForeground(lipgloss.Color("202")),
		},
		Special: ButtonTypeStyle{
			Normal: lipgloss.NewStyle().
				Foreground(lipgloss.Color("15")).
				Background(lipgloss.Color("196")).
				Align(lipgloss.Center).
				Border(lipgloss.NormalBorder(), false).
				BorderForeground(lipgloss.Color("160")),
			Focused: lipgloss.NewStyle().
				Foreground(lipgloss.Color("15")).
				Background(lipgloss.Color("203")).
				Align(lipgloss.Center).
				Border(lipgloss.NormalBorder(), false).
				BorderForeground(lipgloss.Color("62")),
			Pressed: lipgloss.NewStyle().
				Foreground(lipgloss.Color("15")).
				Background(lipgloss.Color("210")).
				Align(lipgloss.Center).
				Border(lipgloss.NormalBorder(), false).
				BorderForeground(lipgloss.Color("94")),
			Disabled: lipgloss.NewStyle().
				Foreground(lipgloss.Color("8")).
				Background(lipgloss.Color("196")).
				Align(lipgloss.Center).
				Border(lipgloss.NormalBorder(), false).
				BorderForeground(lipgloss.Color("160")),
		},
	}
}

// String returns a string representation of the button
func (b *Button) String() string {
	return fmt.Sprintf("Button{Label: %s, Type: %s, State: %s, Value: %s}",
		b.GetLabel(), b.GetType(), b.GetState(), b.GetValue())
}

// ButtonAction represents an action that can be performed by a button
type ButtonAction struct {
	Button  *Button
	Type    string
	Value   string
	Context interface{}
}

// NewButtonAction creates a new button action
func NewButtonAction(button *Button, actionType string) *ButtonAction {
	return &ButtonAction{
		Button: button,
		Type:   actionType,
		Value:  button.GetValue(),
	}
}

// WithContext adds context to the button action
func (ba *ButtonAction) WithContext(context interface{}) *ButtonAction {
	ba.Context = context
	return ba
}

// Trigger creates an action when a button is activated
func (b *Button) Trigger(actionType string) *ButtonAction {
	return NewButtonAction(b, actionType)
}