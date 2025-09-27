package styles

import (
	"github.com/charmbracelet/lipgloss"
)

// RetroStyler provides retro Casio-specific styling patterns
type RetroStyler struct {
	palette        *ColorPalette
	baseStyle      lipgloss.Style
	borderStyle    lipgloss.Border
	cornerStyle    lipgloss.Border
	bevelStyle     lipgloss.Border
	shadowStyle    lipgloss.Style
	highlightStyle lipgloss.Style
}

// NewRetroStyler creates a new retro styler with default configurations
func NewRetroStyler() *RetroStyler {
	palette := NewColorPalette()
	return &RetroStyler{
		palette:     palette,
		baseStyle:   lipgloss.NewStyle(),
		borderStyle: lipgloss.NormalBorder(),
		cornerStyle: lipgloss.RoundedBorder(),
		bevelStyle:  lipgloss.DoubleBorder(),
		shadowStyle: lipgloss.NewStyle().
			Background(palette.GetShadow()).
			Foreground(palette.GetForeground()),
		highlightStyle: lipgloss.NewStyle().
			Background(palette.GetHighlight()).
			Foreground(palette.GetForeground()),
	}
}

// WithPalette sets the color palette
func (rs *RetroStyler) WithPalette(palette *ColorPalette) *RetroStyler {
	rs.palette = palette
	return rs
}

// GetPalette returns the current color palette
func (rs *RetroStyler) GetPalette() *ColorPalette {
	return rs.palette
}

// RetroButtonStyle creates a retro-styled button for the specified type and state
func (rs *RetroStyler) RetroButtonStyle(buttonType, state string) lipgloss.Style {
	colors := rs.palette.GetStateColors(buttonType, state)

	style := lipgloss.NewStyle().
		Foreground(colors.Foreground).
		Background(colors.Background).
		Border(rs.borderStyle, false).
		BorderForeground(colors.Border).
		Align(lipgloss.Center, lipgloss.Center).
		Padding(0, 1)

	// Apply retro-specific effects based on state
	switch state {
	case "focused", "Focused":
		style = rs.applyFocusEffects(style)
	case "pressed", "Pressed":
		style = rs.applyPressedEffects(style)
	case "disabled", "Disabled":
		style = rs.applyDisabledEffects(style)
	default:
		style = rs.applyNormalEffects(style)
	}

	return style
}

// applyNormalEffects applies normal state retro effects
func (rs *RetroStyler) applyNormalEffects(style lipgloss.Style) lipgloss.Style {
	// Add subtle bevel effect for normal state
	return style.BorderForeground(rs.palette.GetBorder())
}

// applyFocusEffects applies focus state retro effects
func (rs *RetroStyler) applyFocusEffects(style lipgloss.Style) lipgloss.Style {
	// Add bright border for focus state
	return style.BorderForeground(rs.palette.GetHighlight())
}

// applyPressedEffects applies pressed state retro effects
func (rs *RetroStyler) applyPressedEffects(style lipgloss.Style) lipgloss.Style {
	// Invert colors slightly and add emphasis
	return style.BorderForeground(lipgloss.Color("94")) // amber highlight
}

// applyDisabledEffects applies disabled state retro effects
func (rs *RetroStyler) applyDisabledEffects(style lipgloss.Style) lipgloss.Style {
	// Dim the colors
	return style.
		Foreground(lipgloss.Color("8")).
		BorderForeground(lipgloss.Color("244"))
}

// RetroGridStyle creates a retro-styled grid container
func (rs *RetroStyler) RetroGridStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Background(rs.palette.GetBackground()).
		Border(rs.cornerStyle).
		BorderForeground(rs.palette.GetBorder()).
		Padding(1).
		Margin(0, 1)
}

// RetroDisplayStyle creates a retro-styled display panel
func (rs *RetroStyler) RetroDisplayStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Background(lipgloss.Color("15")). // white background
		Foreground(lipgloss.Color("0")).  // black text
		Border(rs.bevelStyle).
		BorderForeground(lipgloss.Color("240")).
		Padding(1, 2).
		Align(lipgloss.Right, lipgloss.Center)
}

// RetroTitleStyle creates a retro-styled title text
func (rs *RetroStyler) RetroTitleStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(rs.palette.GetForeground()).
		Background(rs.palette.GetBackground()).
		Bold(true).
		Align(lipgloss.Center)
}

// RetroSubtitleStyle creates a retro-styled subtitle text
func (rs *RetroStyler) RetroSubtitleStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color("244")).
		Background(rs.palette.GetBackground()).
		Align(lipgloss.Center)
}

// RetroBorderStyle creates different retro border styles
func (rs *RetroStyler) RetroBorderStyle(styleType string) lipgloss.Border {
	switch styleType {
	case "normal", "Normal":
		return rs.borderStyle
	case "rounded", "Rounded":
		return rs.cornerStyle
	case "double", "Double":
		return rs.bevelStyle
	case "hidden", "Hidden":
		return lipgloss.HiddenBorder()
	default:
		return rs.borderStyle
	}
}

// RetroShadowStyle creates a retro shadow effect
func (rs *RetroStyler) RetroShadowStyle() lipgloss.Style {
	return rs.shadowStyle
}

// RetroHighlightStyle creates a retro highlight effect
func (rs *RetroStyler) RetroHighlightStyle() lipgloss.Style {
	return rs.highlightStyle
}

// RetroBevelStyle creates a 3D bevel effect
func (rs *RetroStyler) RetroBevelStyle(inset bool) lipgloss.Style {
	if inset {
		return lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("238")) // darker shadow
	}
	return lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("248")) // lighter highlight
}

// RetroButtonFrame creates a retro button frame with 3D effects
func (rs *RetroStyler) RetroButtonFrame(buttonType, state string, width, height int) lipgloss.Style {
	baseStyle := rs.RetroButtonStyle(buttonType, state)

	// Apply dimensions
	if width > 0 {
		baseStyle = baseStyle.Width(width)
	}
	if height > 0 {
		baseStyle = baseStyle.Height(height)
	}

	// Add 3D bevel effect for normal state
	if state == "normal" || state == "Normal" {
		bevelStyle := rs.RetroBevelStyle(false)
		baseStyle = baseStyle.Inherit(bevelStyle)
	}

	return baseStyle
}

// RetroCalculatorTheme creates a complete retro calculator theme
func (rs *RetroStyler) RetroCalculatorTheme() *CalculatorTheme {
	return &CalculatorTheme{
		Background:    rs.RetroGridStyle(),
		Display:       rs.RetroDisplayStyle(),
		Title:         rs.RetroTitleStyle(),
		Subtitle:      rs.RetroSubtitleStyle(),
		NumberButton:  rs.RetroButtonStyle("number", "normal"),
		OperatorButton: rs.RetroButtonStyle("operator", "normal"),
		SpecialButton: rs.RetroButtonStyle("special", "normal"),
		FocusedButton: rs.RetroButtonStyle("number", "focused"),
		PressedButton: rs.RetroButtonStyle("number", "pressed"),
		DisabledButton: rs.RetroButtonStyle("number", "disabled"),
		GridBorder:    rs.RetroBorderStyle("rounded"),
		Shadow:        rs.RetroShadowStyle(),
		Highlight:      rs.RetroHighlightStyle(),
	}
}

// CalculatorTheme represents a complete retro calculator theme
type CalculatorTheme struct {
	Background    lipgloss.Style
	Display       lipgloss.Style
	Title         lipgloss.Style
	Subtitle      lipgloss.Style
	NumberButton  lipgloss.Style
	OperatorButton lipgloss.Style
	SpecialButton lipgloss.Style
	FocusedButton lipgloss.Style
	PressedButton lipgloss.Style
	DisabledButton lipgloss.Style
	GridBorder    lipgloss.Border
	Shadow        lipgloss.Style
	Highlight      lipgloss.Style
}

// GetButtonStyle returns the appropriate button style based on type and state
func (ct *CalculatorTheme) GetButtonStyle(buttonType, state string) lipgloss.Style {
	switch buttonType {
	case "number", "Number":
		return ct.NumberButton
	case "operator", "Operator":
		return ct.OperatorButton
	case "special", "Special":
		return ct.SpecialButton
	default:
		return ct.NumberButton
	}
}

// GetStateStyle returns the appropriate style based on state
func (ct *CalculatorTheme) GetStateStyle(state string) lipgloss.Style {
	switch state {
	case "focused", "Focused":
		return ct.FocusedButton
	case "pressed", "Pressed":
		return ct.PressedButton
	case "disabled", "Disabled":
		return ct.DisabledButton
	default:
		return ct.NumberButton
	}
}

// RetroAnimationFrames returns frames for simple retro animations
func (rs *RetroStyler) RetroAnimationFrames(animationType string) []lipgloss.Style {
	switch animationType {
	case "button_press":
		return []lipgloss.Style{
			rs.RetroButtonStyle("number", "normal"),
			rs.RetroButtonStyle("number", "focused"),
			rs.RetroButtonStyle("number", "pressed"),
		}
	case "display_blink":
		return []lipgloss.Style{
			rs.RetroDisplayStyle(),
			rs.RetroDisplayStyle().Foreground(lipgloss.Color("240")),
		}
	default:
		return []lipgloss.Style{rs.RetroButtonStyle("number", "normal")}
	}
}

// ValidateTheme validates that the retro theme is consistent
func (rs *RetroStyler) ValidateTheme() bool {
	// Check that palette is valid
	if !rs.palette.Validate() {
		return false
	}

	// Test button style creation
	testStyle := rs.RetroButtonStyle("number", "normal")
	if testStyle == (lipgloss.Style{}) {
		return false
	}

	// Test calculator theme
	theme := rs.RetroCalculatorTheme()
	if theme.Background == (lipgloss.Style{}) {
		return false
	}

	return true
}