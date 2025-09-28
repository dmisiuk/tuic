package styles

import (
	"github.com/charmbracelet/lipgloss"
)

// ThemeManager manages different UI themes
type ThemeManager struct {
	retroStyler  *RetroStyler
	currentTheme string
	themes       map[string]*UITheme
}

// UITheme represents a complete UI theme
type UITheme struct {
	Name        string
	Description string
	Colors      *ColorPalette
	Styles      *ThemeStyles
	IsRetro     bool
}

// ThemeStyles contains all styled components for a theme
type ThemeStyles struct {
	Button       ButtonTheme
	Grid         GridTheme
	Display      DisplayTheme
	Text         TextTheme
	Border       BorderTheme
	Animation    AnimationTheme
}

// ButtonTheme defines styling for all button types and states
type ButtonTheme struct {
	Number   ButtonTypeTheme
	Operator ButtonTypeTheme
	Special  ButtonTypeTheme
}

// ButtonTypeTheme defines styling for a specific button type across all states
type ButtonTypeTheme struct {
	Normal   lipgloss.Style
	Focused  lipgloss.Style
	Pressed  lipgloss.Style
	Disabled lipgloss.Style
}

// GridTheme defines styling for grid components
type GridTheme struct {
	Container    lipgloss.Style
	Cell         lipgloss.Style
	CellFocused  lipgloss.Style
	CellPressed  lipgloss.Style
	CellDisabled lipgloss.Style
	Spacing      int
	Padding      int
}

// DisplayTheme defines styling for display components
type DisplayTheme struct {
	Main      lipgloss.Style
	Secondary lipgloss.Style
	Error     lipgloss.Style
	Info      lipgloss.Style
}

// TextTheme defines styling for text elements
type TextTheme struct {
	Title     lipgloss.Style
	Subtitle  lipgloss.Style
	Body      lipgloss.Style
	Caption   lipgloss.Style
	Error     lipgloss.Style
	Success   lipgloss.Style
	Warning   lipgloss.Style
}

// BorderTheme defines styling for borders
type BorderTheme struct {
	Normal   lipgloss.Border
	Focused  lipgloss.Border
	Pressed  lipgloss.Border
	Disabled lipgloss.Border
	Colors   BorderColors
}

// BorderColors defines colors for different border states
type BorderColors struct {
	Normal   lipgloss.Color
	Focused  lipgloss.Color
	Pressed  lipgloss.Color
	Disabled lipgloss.Color
}

// AnimationTheme defines styling for animations
type AnimationTheme struct {
	ButtonPress   []lipgloss.Style
	DisplayBlink  []lipgloss.Style
	Loader        []lipgloss.Style
	Highlight     lipgloss.Style
}

// NewThemeManager creates a new theme manager
func NewThemeManager() *ThemeManager {
	palette := NewColorPalette()
	retroStyler := NewRetroStyler().WithPalette(palette)

	tm := &ThemeManager{
		retroStyler:  retroStyler,
		currentTheme: "retro-casio",
		themes:       make(map[string]*UITheme),
	}

	// Initialize with default themes
	tm.initializeDefaultThemes()

	return tm
}

// initializeDefaultThemes initializes the built-in themes
func (tm *ThemeManager) initializeDefaultThemes() {
	// Retro Casio theme
	tm.themes["retro-casio"] = tm.createRetroCasioTheme()

	// Modern theme
	tm.themes["modern"] = tm.createModernTheme()

	// Minimal theme
	tm.themes["minimal"] = tm.createMinimalTheme()

	// Classic theme
	tm.themes["classic"] = tm.createClassicTheme()
}

// createRetroCasioTheme creates the retro Casio calculator theme
func (tm *ThemeManager) createRetroCasioTheme() *UITheme {
	palette := NewColorPalette()

	return &UITheme{
		Name:        "retro-casio",
		Description: "Classic retro Casio calculator styling",
		Colors:      palette,
		IsRetro:     true,
		Styles: &ThemeStyles{
			Button: tm.createRetroButtonTheme(palette),
			Grid:   tm.createRetroGridTheme(palette),
			Display: tm.createRetroDisplayTheme(palette),
			Text:    tm.createRetroTextTheme(palette),
			Border:  tm.createRetroBorderTheme(palette),
			Animation: tm.createRetroAnimationTheme(palette),
		},
	}
}

// createRetroButtonTheme creates retro button styling
func (tm *ThemeManager) createRetroButtonTheme(palette *ColorPalette) ButtonTheme {
	return ButtonTheme{
		Number: ButtonTypeTheme{
			Normal: lipgloss.NewStyle().
				Foreground(palette.GetNumberColors().Normal.Foreground).
				Background(palette.GetNumberColors().Normal.Background).
				Border(lipgloss.NormalBorder(), false).
				BorderForeground(palette.GetNumberColors().Normal.Border).
				Align(lipgloss.Center, lipgloss.Center).
				Padding(0, 1),
			Focused: lipgloss.NewStyle().
				Foreground(palette.GetNumberColors().Focused.Foreground).
				Background(palette.GetNumberColors().Focused.Background).
				Border(lipgloss.NormalBorder(), false).
				BorderForeground(palette.GetNumberColors().Focused.Border).
				Align(lipgloss.Center, lipgloss.Center).
				Padding(0, 1),
			Pressed: lipgloss.NewStyle().
				Foreground(palette.GetNumberColors().Pressed.Foreground).
				Background(palette.GetNumberColors().Pressed.Background).
				Border(lipgloss.NormalBorder(), false).
				BorderForeground(palette.GetNumberColors().Pressed.Border).
				Align(lipgloss.Center, lipgloss.Center).
				Padding(0, 1),
			Disabled: lipgloss.NewStyle().
				Foreground(palette.GetNumberColors().Disabled.Foreground).
				Background(palette.GetNumberColors().Disabled.Background).
				Border(lipgloss.NormalBorder(), false).
				BorderForeground(palette.GetNumberColors().Disabled.Border).
				Align(lipgloss.Center, lipgloss.Center).
				Padding(0, 1),
		},
		Operator: ButtonTypeTheme{
			Normal: lipgloss.NewStyle().
				Foreground(palette.GetOperatorColors().Normal.Foreground).
				Background(palette.GetOperatorColors().Normal.Background).
				Border(lipgloss.NormalBorder(), false).
				BorderForeground(palette.GetOperatorColors().Normal.Border).
				Align(lipgloss.Center, lipgloss.Center).
				Padding(0, 1),
			Focused: lipgloss.NewStyle().
				Foreground(palette.GetOperatorColors().Focused.Foreground).
				Background(palette.GetOperatorColors().Focused.Background).
				Border(lipgloss.NormalBorder(), false).
				BorderForeground(palette.GetOperatorColors().Focused.Border).
				Align(lipgloss.Center, lipgloss.Center).
				Padding(0, 1),
			Pressed: lipgloss.NewStyle().
				Foreground(palette.GetOperatorColors().Pressed.Foreground).
				Background(palette.GetOperatorColors().Pressed.Background).
				Border(lipgloss.NormalBorder(), false).
				BorderForeground(palette.GetOperatorColors().Pressed.Border).
				Align(lipgloss.Center, lipgloss.Center).
				Padding(0, 1),
			Disabled: lipgloss.NewStyle().
				Foreground(palette.GetOperatorColors().Disabled.Foreground).
				Background(palette.GetOperatorColors().Disabled.Background).
				Border(lipgloss.NormalBorder(), false).
				BorderForeground(palette.GetOperatorColors().Disabled.Border).
				Align(lipgloss.Center, lipgloss.Center).
				Padding(0, 1),
		},
		Special: ButtonTypeTheme{
			Normal: lipgloss.NewStyle().
				Foreground(palette.GetSpecialColors().Normal.Foreground).
				Background(palette.GetSpecialColors().Normal.Background).
				Border(lipgloss.NormalBorder(), false).
				BorderForeground(palette.GetSpecialColors().Normal.Border).
				Align(lipgloss.Center, lipgloss.Center).
				Padding(0, 1),
			Focused: lipgloss.NewStyle().
				Foreground(palette.GetSpecialColors().Focused.Foreground).
				Background(palette.GetSpecialColors().Focused.Background).
				Border(lipgloss.NormalBorder(), false).
				BorderForeground(palette.GetSpecialColors().Focused.Border).
				Align(lipgloss.Center, lipgloss.Center).
				Padding(0, 1),
			Pressed: lipgloss.NewStyle().
				Foreground(palette.GetSpecialColors().Pressed.Foreground).
				Background(palette.GetSpecialColors().Pressed.Background).
				Border(lipgloss.NormalBorder(), false).
				BorderForeground(palette.GetSpecialColors().Pressed.Border).
				Align(lipgloss.Center, lipgloss.Center).
				Padding(0, 1),
			Disabled: lipgloss.NewStyle().
				Foreground(palette.GetSpecialColors().Disabled.Foreground).
				Background(palette.GetSpecialColors().Disabled.Background).
				Border(lipgloss.NormalBorder(), false).
				BorderForeground(palette.GetSpecialColors().Disabled.Border).
				Align(lipgloss.Center, lipgloss.Center).
				Padding(0, 1),
		},
	}
}

// createRetroGridTheme creates retro grid styling
func (tm *ThemeManager) createRetroGridTheme(palette *ColorPalette) GridTheme {
	return GridTheme{
		Container: lipgloss.NewStyle().
			Background(palette.GetBackground()).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(palette.GetBorder()).
			Padding(1).
			Margin(0, 1),
		Cell: lipgloss.NewStyle().
			Background(palette.GetBackground()).
			Border(lipgloss.NormalBorder()).
			BorderForeground(palette.GetBorder()).
			Align(lipgloss.Center, lipgloss.Center),
		CellFocused: lipgloss.NewStyle().
			Background(palette.GetBackground()).
			Border(lipgloss.NormalBorder()).
			BorderForeground(palette.GetHighlight()).
			Align(lipgloss.Center, lipgloss.Center),
		CellPressed: lipgloss.NewStyle().
			Background(palette.GetBackground()).
			Border(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("94")).
			Align(lipgloss.Center, lipgloss.Center),
		CellDisabled: lipgloss.NewStyle().
			Background(palette.GetBackground()).
			Border(lipgloss.HiddenBorder()).
			Align(lipgloss.Center, lipgloss.Center),
		Spacing: 1,
		Padding: 1,
	}
}

// createRetroDisplayTheme creates retro display styling
func (tm *ThemeManager) createRetroDisplayTheme(palette *ColorPalette) DisplayTheme {
	return DisplayTheme{
		Main: lipgloss.NewStyle().
			Background(lipgloss.Color("15")). // white background
			Foreground(lipgloss.Color("0")).  // black text
			Border(lipgloss.DoubleBorder()).
			BorderForeground(lipgloss.Color("240")).
			Padding(1, 2).
			Align(lipgloss.Right, lipgloss.Center),
		Secondary: lipgloss.NewStyle().
			Background(palette.GetBackground()).
			Foreground(palette.GetForeground()).
			Align(lipgloss.Left, lipgloss.Center),
		Error: lipgloss.NewStyle().
			Background(lipgloss.Color("196")). // red background
			Foreground(lipgloss.Color("15")).  // white text
			Align(lipgloss.Center, lipgloss.Center),
		Info: lipgloss.NewStyle().
			Background(palette.GetBackground()).
			Foreground(lipgloss.Color("14")). // cyan text
			Align(lipgloss.Center, lipgloss.Center),
	}
}

// createRetroTextTheme creates retro text styling
func (tm *ThemeManager) createRetroTextTheme(palette *ColorPalette) TextTheme {
	return TextTheme{
		Title: lipgloss.NewStyle().
			Foreground(palette.GetForeground()).
			Background(palette.GetBackground()).
			Bold(true).
			Align(lipgloss.Center),
		Subtitle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("244")).
			Background(palette.GetBackground()).
			Align(lipgloss.Center),
		Body: lipgloss.NewStyle().
			Foreground(palette.GetForeground()).
			Background(palette.GetBackground()),
		Caption: lipgloss.NewStyle().
			Foreground(lipgloss.Color("8")).
			Background(palette.GetBackground()).
			Align(lipgloss.Center),
		Error: lipgloss.NewStyle().
			Foreground(lipgloss.Color("196")).
			Background(palette.GetBackground()).
			Bold(true),
		Success: lipgloss.NewStyle().
			Foreground(lipgloss.Color("46")).
			Background(palette.GetBackground()).
			Bold(true),
		Warning: lipgloss.NewStyle().
			Foreground(lipgloss.Color("226")).
			Background(palette.GetBackground()).
			Bold(true),
	}
}

// createRetroBorderTheme creates retro border styling
func (tm *ThemeManager) createRetroBorderTheme(palette *ColorPalette) BorderTheme {
	return BorderTheme{
		Normal:   lipgloss.NormalBorder(),
		Focused:  lipgloss.NormalBorder(),
		Pressed:  lipgloss.DoubleBorder(),
		Disabled: lipgloss.HiddenBorder(),
		Colors: BorderColors{
			Normal:   palette.GetBorder(),
			Focused:  palette.GetHighlight(),
			Pressed:  lipgloss.Color("94"),
			Disabled: lipgloss.Color("244"),
		},
	}
}

// createRetroAnimationTheme creates retro animation styling
func (tm *ThemeManager) createRetroAnimationTheme(palette *ColorPalette) AnimationTheme {
	return AnimationTheme{
		ButtonPress: []lipgloss.Style{
			tm.retroStyler.RetroButtonStyle("number", "normal"),
			tm.retroStyler.RetroButtonStyle("number", "focused"),
			tm.retroStyler.RetroButtonStyle("number", "pressed"),
		},
		DisplayBlink: []lipgloss.Style{
			tm.retroStyler.RetroDisplayStyle(),
			tm.retroStyler.RetroDisplayStyle().Foreground(lipgloss.Color("240")),
		},
		Loader: []lipgloss.Style{
			lipgloss.NewStyle().Foreground(lipgloss.Color("240")),
			lipgloss.NewStyle().Foreground(lipgloss.Color("244")),
			lipgloss.NewStyle().Foreground(lipgloss.Color("248")),
			lipgloss.NewStyle().Foreground(lipgloss.Color("15")),
		},
		Highlight: lipgloss.NewStyle().
			Background(palette.GetHighlight()).
			Foreground(palette.GetForeground()),
	}
}

// createModernTheme creates a modern theme (fallback)
func (tm *ThemeManager) createModernTheme() *UITheme {
	palette := NewColorPalette()
	return &UITheme{
		Name:        "modern",
		Description: "Modern clean styling",
		Colors:      palette,
		IsRetro:     false,
		Styles:      &ThemeStyles{}, // Simplified for brevity
	}
}

// createMinimalTheme creates a minimal theme (fallback)
func (tm *ThemeManager) createMinimalTheme() *UITheme {
	palette := NewColorPalette()
	return &UITheme{
		Name:        "minimal",
		Description: "Minimal styling",
		Colors:      palette,
		IsRetro:     false,
		Styles:      &ThemeStyles{}, // Simplified for brevity
	}
}

// createClassicTheme creates a classic theme (fallback)
func (tm *ThemeManager) createClassicTheme() *UITheme {
	palette := NewColorPalette()
	return &UITheme{
		Name:        "classic",
		Description: "Classic styling",
		Colors:      palette,
		IsRetro:     false,
		Styles:      &ThemeStyles{}, // Simplified for brevity
	}
}

// GetTheme returns a theme by name
func (tm *ThemeManager) GetTheme(name string) (*UITheme, error) {
	theme, exists := tm.themes[name]
	if !exists {
		return tm.themes["retro-casio"], nil // Default to retro Casio
	}
	return theme, nil
}

// GetCurrentTheme returns the current active theme
func (tm *ThemeManager) GetCurrentTheme() *UITheme {
	return tm.themes[tm.currentTheme]
}

// SetTheme sets the current theme
func (tm *ThemeManager) SetTheme(name string) error {
	if _, exists := tm.themes[name]; !exists {
		return &ThemeNotFoundError{Name: name}
	}
	tm.currentTheme = name
	return nil
}

// ListThemes returns a list of available theme names
func (tm *ThemeManager) ListThemes() []string {
	var names []string
	for name := range tm.themes {
		names = append(names, name)
	}
	return names
}

// GetButtonTheme returns the button theme for the current theme
func (tm *ThemeManager) GetButtonTheme() ButtonTheme {
	return tm.GetCurrentTheme().Styles.Button
}

// GetButtonStyle returns a button style for the specified type and state
func (tm *ThemeManager) GetButtonStyle(buttonType, state string) lipgloss.Style {
	buttonTheme := tm.GetButtonTheme()

	switch buttonType {
	case "number", "Number":
		return tm.getButtonTypeStyle(buttonTheme.Number, state)
	case "operator", "Operator":
		return tm.getButtonTypeStyle(buttonTheme.Operator, state)
	case "special", "Special":
		return tm.getButtonTypeStyle(buttonTheme.Special, state)
	default:
		return tm.getButtonTypeStyle(buttonTheme.Number, state)
	}
}

// getButtonTypeStyle returns the style for a button type based on state
func (tm *ThemeManager) getButtonTypeStyle(theme ButtonTypeTheme, state string) lipgloss.Style {
	switch state {
	case "normal", "Normal", "":
		return theme.Normal
	case "focused", "Focused":
		return theme.Focused
	case "pressed", "Pressed":
		return theme.Pressed
	case "disabled", "Disabled":
		return theme.Disabled
	default:
		return theme.Normal
	}
}

// ThemeNotFoundError represents an error when a theme is not found
type ThemeNotFoundError struct {
	Name string
}

// Error implements the error interface
func (e *ThemeNotFoundError) Error() string {
	return "theme not found: " + e.Name
}

// Is compares error types for error handling
func (e *ThemeNotFoundError) Is(target error) bool {
	_, ok := target.(*ThemeNotFoundError)
	return ok
}