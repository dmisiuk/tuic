package visual

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"strings"

	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
)

// TerminalConfig represents terminal configuration for screenshot capture
type TerminalConfig struct {
	Width      int
	Height     int
	FontFace   font.Face
	FontSize   int
	Foreground color.Color
	Background color.Color
	CellWidth  int
	CellHeight int
}

// Screenshot represents a captured terminal screenshot
type Screenshot struct {
	Image    *image.RGBA
	Config   TerminalConfig
	Metadata ScreenshotMetadata
}

// ScreenshotMetadata contains metadata about the screenshot
type ScreenshotMetadata struct {
	Timestamp    string
	Application  string
	Version      string
	Theme        string
	TerminalType string
	Width        int
	Height       int
}

// NewDefaultConfig creates a default terminal configuration
func NewDefaultConfig() TerminalConfig {
	return TerminalConfig{
		Width:      80,
		Height:     24,
		FontFace:   basicfont.Face7x13,
		FontSize:   13,
		Foreground: color.White,
		Background: color.Black,
		CellWidth:  7,
		CellHeight: 13,
	}
}

// CaptureTerminal captures a terminal screenshot from text content
func CaptureTerminal(content string, config TerminalConfig) (*Screenshot, error) {
	// Create image with appropriate dimensions
	imgWidth := config.Width * config.CellWidth
	imgHeight := config.Height * config.CellHeight

	img := image.NewRGBA(image.Rect(0, 0, imgWidth, imgHeight))

	// Fill background
	for y := 0; y < imgHeight; y++ {
		for x := 0; x < imgWidth; x++ {
			img.Set(x, y, config.Background)
		}
	}

	// Parse and render content
	lines := strings.Split(content, "\n")
	for y, line := range lines {
		if y >= config.Height {
			break
		}
		renderLine(img, line, y, config)
	}

	metadata := ScreenshotMetadata{
		Timestamp:    fmt.Sprintf("%d", 0), // Placeholder
		Application:  "CCPM Calculator",
		Version:      "1.0.0",
		Theme:        "retro-casio",
		TerminalType: "xterm-256color",
		Width:        config.Width,
		Height:       config.Height,
	}

	return &Screenshot{
		Image:    img,
		Config:   config,
		Metadata: metadata,
	}, nil
}

// CaptureWithStyling captures terminal content with LipGloss styling
func CaptureWithStyling(styledContent string, config TerminalConfig) (*Screenshot, error) {
	// Strip ANSI codes and capture plain text for now
	// TODO: Implement proper ANSI code parsing and rendering
	plainText := stripANSICodes(styledContent)
	return CaptureTerminal(plainText, config)
}

// stripANSICodes removes ANSI escape codes from text
func stripANSICodes(text string) string {
	// Simple ANSI code stripper
	result := ""
	inEscape := false

	for _, char := range text {
		if char == '\x1b' {
			inEscape = true
			continue
		}

		if inEscape {
			if char >= 'A' && char <= 'Z' || char >= 'a' && char <= 'z' {
				inEscape = false
			}
			continue
		}

		result += string(char)
	}

	return result
}

// renderLine renders a single line of text to the image
func renderLine(img *image.RGBA, line string, lineNum int, config TerminalConfig) {
	drawer := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(config.Foreground),
		Face: config.FontFace,
		Dot:  fixed.P(0, (lineNum+1)*config.CellHeight-2),
	}

	// Handle text overflow
	if len(line) > config.Width {
		line = line[:config.Width]
	}

	drawer.DrawString(line)
}

// Save saves the screenshot to a file
func (s *Screenshot) Save(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	return png.Encode(file, s.Image)
}

// ToBytes returns the screenshot as PNG bytes
func (s *Screenshot) ToBytes() ([]byte, error) {
	var buf bytes.Buffer
	err := png.Encode(&buf, s.Image)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// Compare compares two screenshots and returns a diff image
func (s *Screenshot) Compare(other *Screenshot) (*image.RGBA, float64, error) {
	if s.Image.Bounds() != other.Image.Bounds() {
		return nil, 0, fmt.Errorf("screenshot dimensions don't match")
	}

	bounds := s.Image.Bounds()
	diff := image.NewRGBA(bounds)

	diffPixels := 0
	totalPixels := bounds.Dx() * bounds.Dy()

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			c1 := s.Image.RGBAAt(x, y)
			c2 := other.Image.RGBAAt(x, y)

			if c1 != c2 {
				// Mark differences in red
				diff.Set(x, y, color.RGBA{255, 0, 0, 255})
				diffPixels++
			} else {
				// Keep original pixel
				diff.Set(x, y, c1)
			}
		}
	}

	diffRatio := float64(diffPixels) / float64(totalPixels)
	return diff, diffRatio, nil
}

// NewScreenshotFromModel captures a screenshot from a Bubble Tea model
func NewScreenshotFromModel(model interface{}, config TerminalConfig) (*Screenshot, error) {
	// This would need to be implemented based on the actual model type
	// For now, we'll create a placeholder implementation

	// Get the view from the model
	var view string
	if m, ok := model.(interface{ View() string }); ok {
		view = m.View()
	} else {
		return nil, fmt.Errorf("model does not implement View() method")
	}

	return CaptureWithStyling(view, config)
}

// DecodePNG decodes a PNG image from a reader
func DecodePNG(r interface{}) (*image.RGBA, error) {
	// Placeholder implementation - would need actual PNG decoding
	return image.NewRGBA(image.Rect(0, 0, 80, 24)), nil
}

// SavePNG saves an image as PNG
func SavePNG(filename string, img *image.RGBA) error {
	// Placeholder implementation - would need actual PNG encoding
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// Would use png.Encode here
	return nil
}