package styles

import (
	"github.com/charmbracelet/lipgloss"
)

// ColorPalette defines the retro Casio-inspired color scheme
type ColorPalette struct {
	// Primary colors for different button types
	NumberColors    ButtonColorSet
	OperatorColors  ButtonColorSet
	SpecialColors   ButtonColorSet

	// General UI colors
	Background      lipgloss.Color
	Foreground      lipgloss.Color
	Border          lipgloss.Color
	Shadow          lipgloss.Color
	Highlight       lipgloss.Color

	// State-specific colors
	FocusColors     ButtonStateColors
	DisabledColors  ButtonStateColors
}

// ButtonColorSet defines colors for a button type across different states
type ButtonColorSet struct {
	Normal   ButtonStateColors
	Focused  ButtonStateColors
	Pressed  ButtonStateColors
	Disabled ButtonStateColors
}

// ButtonStateColors defines foreground and background colors for a button state
type ButtonStateColors struct {
	Foreground lipgloss.Color
	Background lipgloss.Color
	Border     lipgloss.Color
}

// NewColorPalette creates a new retro Casio-inspired color palette
func NewColorPalette() *ColorPalette {
	return &ColorPalette{
		// Number buttons - classic gray/white scheme
		NumberColors: ButtonColorSet{
			Normal: ButtonStateColors{
				Foreground: lipgloss.Color("15"),   // white
				Background: lipgloss.Color("240"),  // dark gray
				Border:     lipgloss.Color("244"),  // light gray
			},
			Focused: ButtonStateColors{
				Foreground: lipgloss.Color("15"),   // white
				Background: lipgloss.Color("246"),  // lighter gray
				Border:     lipgloss.Color("62"),   // blue highlight
			},
			Pressed: ButtonStateColors{
				Foreground: lipgloss.Color("15"),   // white
				Background: lipgloss.Color("248"),  // light gray
				Border:     lipgloss.Color("94"),   // amber highlight
			},
			Disabled: ButtonStateColors{
				Foreground: lipgloss.Color("8"),    // dark gray
				Background: lipgloss.Color("240"),  // dark gray
				Border:     lipgloss.Color("244"),  // light gray
			},
		},

		// Operator buttons - classic orange/amber scheme
		OperatorColors: ButtonColorSet{
			Normal: ButtonStateColors{
				Foreground: lipgloss.Color("15"),   // white
				Background: lipgloss.Color("208"),  // orange
				Border:     lipgloss.Color("202"),  // bright orange
			},
			Focused: ButtonStateColors{
				Foreground: lipgloss.Color("15"),   // white
				Background: lipgloss.Color("214"),  // light orange
				Border:     lipgloss.Color("62"),   // blue highlight
			},
			Pressed: ButtonStateColors{
				Foreground: lipgloss.Color("15"),   // white
				Background: lipgloss.Color("220"),  // light amber
				Border:     lipgloss.Color("94"),   // amber highlight
			},
			Disabled: ButtonStateColors{
				Foreground: lipgloss.Color("8"),    // dark gray
				Background: lipgloss.Color("208"),  // orange (dimmed)
				Border:     lipgloss.Color("202"),  // bright orange
			},
		},

		// Special buttons - classic red scheme
		SpecialColors: ButtonColorSet{
			Normal: ButtonStateColors{
				Foreground: lipgloss.Color("15"),   // white
				Background: lipgloss.Color("196"),  // red
				Border:     lipgloss.Color("160"),  // dark red
			},
			Focused: ButtonStateColors{
				Foreground: lipgloss.Color("15"),   // white
				Background: lipgloss.Color("203"),  // light red
				Border:     lipgloss.Color("62"),   // blue highlight
			},
			Pressed: ButtonStateColors{
				Foreground: lipgloss.Color("15"),   // white
				Background: lipgloss.Color("210"),  // pink
				Border:     lipgloss.Color("94"),   // amber highlight
			},
			Disabled: ButtonStateColors{
				Foreground: lipgloss.Color("8"),    // dark gray
				Background: lipgloss.Color("196"),  // red (dimmed)
				Border:     lipgloss.Color("160"),  // dark red
			},
		},

		// General UI colors
		Background: lipgloss.Color("235"),      // dark background
		Foreground: lipgloss.Color("15"),       // white text
		Border:     lipgloss.Color("244"),      // light gray borders
		Shadow:     lipgloss.Color("238"),      // shadow color
		Highlight:  lipgloss.Color("62"),       // blue highlight

		// State-specific fallback colors
		FocusColors: ButtonStateColors{
			Foreground: lipgloss.Color("15"),   // white
			Background: lipgloss.Color("62"),   // blue
			Border:     lipgloss.Color("94"),   // amber
		},

		DisabledColors: ButtonStateColors{
			Foreground: lipgloss.Color("8"),    // dark gray
			Background: lipgloss.Color("240"),  // dark gray
			Border:     lipgloss.Color("244"),  // light gray
		},
	}
}

// GetNumberColors returns the color set for number buttons
func (cp *ColorPalette) GetNumberColors() ButtonColorSet {
	return cp.NumberColors
}

// GetOperatorColors returns the color set for operator buttons
func (cp *ColorPalette) GetOperatorColors() ButtonColorSet {
	return cp.OperatorColors
}

// GetSpecialColors returns the color set for special buttons
func (cp *ColorPalette) GetSpecialColors() ButtonColorSet {
	return cp.SpecialColors
}

// GetBackground returns the background color
func (cp *ColorPalette) GetBackground() lipgloss.Color {
	return cp.Background
}

// GetForeground returns the foreground color
func (cp *ColorPalette) GetForeground() lipgloss.Color {
	return cp.Foreground
}

// GetBorder returns the border color
func (cp *ColorPalette) GetBorder() lipgloss.Color {
	return cp.Border
}

// GetShadow returns the shadow color
func (cp *ColorPalette) GetShadow() lipgloss.Color {
	return cp.Shadow
}

// GetHighlight returns the highlight color
func (cp *ColorPalette) GetHighlight() lipgloss.Color {
	return cp.Highlight
}

// GetFocusColors returns the focus state color set
func (cp *ColorPalette) GetFocusColors() ButtonStateColors {
	return cp.FocusColors
}

// GetDisabledColors returns the disabled state color set
func (cp *ColorPalette) GetDisabledColors() ButtonStateColors {
	return cp.DisabledColors
}

// GetButtonColors returns the color set for a specific button type
func (cp *ColorPalette) GetButtonColors(buttonType string) ButtonColorSet {
	switch buttonType {
	case "number", "Number":
		return cp.NumberColors
	case "operator", "Operator":
		return cp.OperatorColors
	case "special", "Special":
		return cp.SpecialColors
	default:
		return cp.NumberColors // default to number colors
	}
}

// GetStateColors returns the colors for a specific state within a button type
func (cp *ColorPalette) GetStateColors(buttonType, state string) ButtonStateColors {
	colorSet := cp.GetButtonColors(buttonType)

	switch state {
	case "normal", "Normal", "":
		return colorSet.Normal
	case "focused", "Focused":
		return colorSet.Focused
	case "pressed", "Pressed":
		return colorSet.Pressed
	case "disabled", "Disabled":
		return colorSet.Disabled
	default:
		return colorSet.Normal // default to normal state
	}
}

// ColorNames returns human-readable color names for debugging
func (cp *ColorPalette) ColorNames() map[string]string {
	return map[string]string{
		"background":    "dark gray",
		"foreground":    "white",
		"border":        "light gray",
		"shadow":        "shadow gray",
		"highlight":     "blue",
		"number_normal": "dark gray",
		"number_focus":  "light gray",
		"number_press":  "lightest gray",
		"operator_norm": "orange",
		"operator_focus": "light orange",
		"operator_press": "light amber",
		"special_normal": "red",
		"special_focus":  "light red",
		"special_press":  "pink",
	}
}

// Validate checks if all colors are valid ANSI colors
func (cp *ColorPalette) Validate() bool {
	// Check main colors
	mainColors := []lipgloss.Color{
		cp.Background, cp.Foreground, cp.Border, cp.Shadow, cp.Highlight,
	}

	for _, color := range mainColors {
		if !cp.isValidColor(color) {
			return false
		}
	}

	// Check button colors
	buttonSets := []ButtonColorSet{cp.NumberColors, cp.OperatorColors, cp.SpecialColors}
	for _, set := range buttonSets {
		states := []ButtonStateColors{set.Normal, set.Focused, set.Pressed, set.Disabled}
		for _, state := range states {
			if !cp.isValidColor(state.Foreground) || !cp.isValidColor(state.Background) || !cp.isValidColor(state.Border) {
				return false
			}
		}
	}

	return true
}

// isValidColor checks if a color is valid
func (cp *ColorPalette) isValidColor(color lipgloss.Color) bool {
	return color != "" && string(color) != "0"
}