package styles

import (
	"github.com/charmbracelet/lipgloss"
)

// StyleSystem manages all styling for the UI components
type StyleSystem struct {
	colors   *ColorPalette
	themes   *ThemeManager
	retro    *RetroStyler
	renderer *StyleRenderer
}

// NewStyleSystem creates a new style system with default configurations
func NewStyleSystem() *StyleSystem {
	return &StyleSystem{
		colors:   NewColorPalette(),
		themes:   NewThemeManager(),
		retro:    NewRetroStyler(),
		renderer: NewStyleRenderer(),
	}
}

// GetColors returns the color palette
func (ss *StyleSystem) GetColors() *ColorPalette {
	return ss.colors
}

// GetThemes returns the theme manager
func (ss *StyleSystem) GetThemes() *ThemeManager {
	return ss.themes
}

// GetRetro returns the retro styler
func (ss *StyleSystem) GetRetro() *RetroStyler {
	return ss.retro
}

// GetRenderer returns the style renderer
func (ss *StyleSystem) GetRenderer() *StyleRenderer {
	return ss.renderer
}

// StyleRenderer handles the rendering of styled components
type StyleRenderer struct {
	baseStyle lipgloss.Style
}

// NewStyleRenderer creates a new style renderer
func NewStyleRenderer() *StyleRenderer {
	return &StyleRenderer{
		baseStyle: lipgloss.NewStyle(),
	}
}

// ApplyBaseStyle applies the base styling to a given style
func (sr *StyleRenderer) ApplyBaseStyle(style lipgloss.Style) lipgloss.Style {
	return style.Inherit(sr.baseStyle)
}

// WithPadding applies padding to a style
func (sr *StyleRenderer) WithPadding(style lipgloss.Style, top, right, bottom, left int) lipgloss.Style {
	return style.Padding(top, right, bottom, left)
}

// WithMargin applies margin to a style
func (sr *StyleRenderer) WithMargin(style lipgloss.Style, top, right, bottom, left int) lipgloss.Style {
	return style.Margin(top, right, bottom, left)
}

// WithBorder applies border styling to a style
func (sr *StyleRenderer) WithBorder(style lipgloss.Style, border lipgloss.Border) lipgloss.Style {
	return style.Border(border)
}

// WithBorderForeground sets the border color
func (sr *StyleRenderer) WithBorderForeground(style lipgloss.Style, color lipgloss.Color) lipgloss.Style {
	return style.BorderForeground(color)
}

// WithAlignment sets text alignment
func (sr *StyleRenderer) WithAlignment(style lipgloss.Style, horizontal, vertical lipgloss.Position) lipgloss.Style {
	return style.Align(horizontal, vertical)
}

// WithDimensions sets width and height
func (sr *StyleRenderer) WithDimensions(style lipgloss.Style, width, height int) lipgloss.Style {
	result := style
	if width > 0 {
		result = result.Width(width)
	}
	if height > 0 {
		result = result.Height(height)
	}
	return result
}

// WithColors sets foreground and background colors
func (sr *StyleRenderer) WithColors(style lipgloss.Style, fg, bg lipgloss.Color) lipgloss.Style {
	result := style
	if fg != "" {
		result = result.Foreground(fg)
	}
	if bg != "" {
		result = result.Background(bg)
	}
	return result
}

// StyleConfig represents configuration for creating styles
type StyleConfig struct {
	Foreground     lipgloss.Color
	Background     lipgloss.Color
	Border         lipgloss.Border
	BorderForegnd  lipgloss.Color
	Width          int
	Height         int
	AlignHorizontal lipgloss.Position
	AlignVertical   lipgloss.Position
	PaddingTop     int
	PaddingRight   int
	PaddingBottom  int
	PaddingLeft    int
	MarginTop      int
	MarginRight    int
	MarginBottom   int
	MarginLeft     int
}

// NewStyle creates a new style from configuration
func (sr *StyleRenderer) NewStyle(config StyleConfig) lipgloss.Style {
	style := lipgloss.NewStyle()

	// Apply colors
	if config.Foreground != "" {
		style = style.Foreground(config.Foreground)
	}
	if config.Background != "" {
		style = style.Background(config.Background)
	}

	// Apply border
	if config.Border != (lipgloss.Border{}) {
		style = style.Border(config.Border)
	}
	if config.BorderForegnd != "" {
		style = style.BorderForeground(config.BorderForegnd)
	}

	// Apply dimensions
	if config.Width > 0 {
		style = style.Width(config.Width)
	}
	if config.Height > 0 {
		style = style.Height(config.Height)
	}

	// Apply alignment
	if config.AlignHorizontal != lipgloss.Position(0) || config.AlignVertical != lipgloss.Position(0) {
		style = style.Align(config.AlignHorizontal, config.AlignVertical)
	}

	// Apply padding
	if config.PaddingTop > 0 || config.PaddingRight > 0 || config.PaddingBottom > 0 || config.PaddingLeft > 0 {
		style = style.Padding(config.PaddingTop, config.PaddingRight, config.PaddingBottom, config.PaddingLeft)
	}

	// Apply margin
	if config.MarginTop > 0 || config.MarginRight > 0 || config.MarginBottom > 0 || config.MarginLeft > 0 {
		style = style.Margin(config.MarginTop, config.MarginRight, config.MarginBottom, config.MarginLeft)
	}

	return style
}

// DefaultStyleConfig returns a default style configuration
func (sr *StyleRenderer) DefaultStyleConfig() StyleConfig {
	return StyleConfig{
		Foreground:     lipgloss.Color("15"), // white
		Background:     lipgloss.Color("240"), // dark gray
		Border:         lipgloss.NormalBorder(),
		BorderForegnd:  lipgloss.Color("244"), // light gray
		AlignHorizontal: lipgloss.Center,
		AlignVertical:   lipgloss.Center,
		PaddingTop:     0,
		PaddingRight:   1,
		PaddingBottom:  0,
		PaddingLeft:    1,
	}
}