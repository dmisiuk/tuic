package visual

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

// DemoAction represents a single action in a demo sequence
type DemoAction struct {
	Type        string        `json:"type"`
	Description string        `json:"description"`
	Delay       time.Duration `json:"delay"`
	KeyPress    tea.KeyMsg   `json:"keyPress,omitempty"`
	MouseClick  MouseClick   `json:"mouseClick,omitempty"`
	Screenshot  *Screenshot  `json:"screenshot,omitempty"`
}

// MouseClick represents a mouse click action
type MouseClick struct {
	X int `json:"x"`
	Y int `json:"y"`
}

// DemoSequence represents a sequence of demo actions
type DemoSequence struct {
	Name        string        `json:"name"`
	Description string        `json:"description"`
	Actions     []DemoAction  `json:"actions"`
	Duration    time.Duration `json:"duration"`
	CreatedAt   time.Time     `json:"createdAt"`
}

// DemoGenerator generates automated demos for calculator operations
type DemoGenerator struct {
	Model         interface{}
	Config        TerminalConfig
	OutputDir     string
	CurrentFrame  int
	Recording     bool
	Sequence      *DemoSequence
}

// NewDemoGenerator creates a new demo generator
func NewDemoGenerator(model interface{}, config TerminalConfig, outputDir string) *DemoGenerator {
	return &DemoGenerator{
		Model:     model,
		Config:    config,
		OutputDir: outputDir,
		Recording: false,
	}
}

// StartRecording starts recording a new demo sequence
func (dg *DemoGenerator) StartRecording(name, description string) error {
	if dg.Recording {
		return fmt.Errorf("already recording")
	}

	dg.Sequence = &DemoSequence{
		Name:        name,
		Description: description,
		Actions:     []DemoAction{},
		CreatedAt:   time.Now(),
	}
	dg.Recording = true
	dg.CurrentFrame = 0

	// Create output directory
	if err := os.MkdirAll(dg.OutputDir, 0755); err != nil {
		return err
	}

	return nil
}

// StopRecording stops the current recording
func (dg *DemoGenerator) StopRecording() error {
	if !dg.Recording {
		return fmt.Errorf("not recording")
	}

	dg.Recording = false
	dg.Sequence.Duration = time.Since(dg.Sequence.CreatedAt)

	// Save demo metadata
	return dg.saveDemoMetadata()
}

// CaptureFrame captures a frame with optional description
func (dg *DemoGenerator) CaptureFrame(description string) error {
	if !dg.Recording {
		return fmt.Errorf("not recording")
	}

	screenshot, err := NewScreenshotFromModel(dg.Model, dg.Config)
	if err != nil {
		return err
	}

	// Save screenshot
	filename := fmt.Sprintf("frame_%04d.png", dg.CurrentFrame)
	screenshotPath := filepath.Join(dg.OutputDir, filename)
	if err := screenshot.Save(screenshotPath); err != nil {
		return err
	}

	// Add action to sequence
	action := DemoAction{
		Type:        "screenshot",
		Description: description,
		Delay:       0, // Will be calculated later
		Screenshot:  screenshot,
	}

	dg.Sequence.Actions = append(dg.Sequence.Actions, action)
	dg.CurrentFrame++

	return nil
}

// AddKeyPress adds a key press action to the demo
func (dg *DemoGenerator) AddKeyPress(key tea.KeyMsg, description string) error {
	if !dg.Recording {
		return fmt.Errorf("not recording")
	}

	action := DemoAction{
		Type:        "keypress",
		Description: description,
		Delay:       100 * time.Millisecond, // Default delay
		KeyPress:    key,
	}

	dg.Sequence.Actions = append(dg.Sequence.Actions, action)
	return nil
}

// AddMouseClick adds a mouse click action to the demo
func (dg *DemoGenerator) AddMouseClick(x, y int, description string) error {
	if !dg.Recording {
		return fmt.Errorf("not recording")
	}

	action := DemoAction{
		Type:        "mouseclick",
		Description: description,
		Delay:       100 * time.Millisecond, // Default delay
		MouseClick:  MouseClick{X: x, Y: y},
	}

	dg.Sequence.Actions = append(dg.Sequence.Actions, action)
	return nil
}

// saveDemoMetadata saves the demo metadata to JSON
func (dg *DemoGenerator) saveDemoMetadata() error {
	metadataPath := filepath.Join(dg.OutputDir, "demo.json")
	data, err := json.MarshalIndent(dg.Sequence, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(metadataPath, data, 0644)
}

// GenerateBasicDemo generates a basic calculator demo
func (dg *DemoGenerator) GenerateBasicDemo() error {
	if err := dg.StartRecording("basic_calculator", "Basic calculator operations"); err != nil {
		return err
	}
	defer dg.StopRecording()

	// Initial state
	if err := dg.CaptureFrame("Initial calculator state"); err != nil {
		return err
	}

	// Simple calculation: 123 + 456 = 579
	calculations := []struct {
		keys       []rune
		description string
	}{
		{[]rune{'1'}, "Press '1'"},
		{[]rune{'2'}, "Press '2'"},
		{[]rune{'3'}, "Press '3'"},
		{[]rune{'+'}, "Press '+'"},
		{[]rune{'4'}, "Press '4'"},
		{[]rune{'5'}, "Press '5'"},
		{[]rune{'6'}, "Press '6'"},
		{[]rune{'='}, "Press '=' to calculate"},
	}

	for _, calc := range calculations {
		for _, key := range calc.keys {
			if err := dg.AddKeyPress(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{key}}, calc.description); err != nil {
				return err
			}
		}
		if err := dg.CaptureFrame(calc.description); err != nil {
			return err
		}
	}

	return nil
}

// GenerateAdvancedDemo generates an advanced calculator demo
func (dg *DemoGenerator) GenerateAdvancedDemo() error {
	if err := dg.StartRecording("advanced_calculator", "Advanced calculator operations"); err != nil {
		return err
	}
	defer dg.StopRecording()

	// Initial state
	if err := dg.CaptureFrame("Initial calculator state"); err != nil {
		return err
	}

	// Complex calculation: (123 + 456) * 2 / 3
	operations := []struct {
		keys       []rune
		description string
	}{
		{[]rune{'('}, "Start parentheses"},
		{[]rune{'1'}, "Press '1'"},
		{[]rune{'2'}, "Press '2'"},
		{[]rune{'3'}, "Press '3'"},
		{[]rune{')'}, "End parentheses"},
		{[]rune{'+'}, "Press '+'"},
		{[]rune{'4'}, "Press '4'"},
		{[]rune{'5'}, "Press '5'"},
		{[]rune{'6'}, "Press '6'"},
		{[]rune{'*'}, "Press '*'"},
		{[]rune{'2'}, "Press '2'"},
		{[]rune{'/'}, "Press '/'"},
		{[]rune{'3'}, "Press '3'"},
		{[]rune{'='}, "Calculate result"},
	}

	for _, op := range operations {
		for _, key := range op.keys {
			if err := dg.AddKeyPress(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{key}}, op.description); err != nil {
				return err
			}
		}
		if err := dg.CaptureFrame(op.description); err != nil {
			return err
		}
	}

	return nil
}

// GenerateErrorDemo generates error handling demo
func (dg *DemoGenerator) GenerateErrorDemo() error {
	if err := dg.StartRecording("error_handling", "Error handling demonstration"); err != nil {
		return err
	}
	defer dg.StopRecording()

	// Initial state
	if err := dg.CaptureFrame("Initial calculator state"); err != nil {
		return err
	}

	// Division by zero
	errorOps := []struct {
		keys       []rune
		description string
	}{
		{[]rune{'1'}, "Press '1'"},
		{[]rune{'0'}, "Press '0'"},
		{[]rune{'/'}, "Press '/'"},
		{[]rune{'0'}, "Press '0' (division by zero)"},
		{[]rune{'='}, "Attempt calculation"},
		{[]rune{'C'}, "Clear error"},
	}

	for _, op := range errorOps {
		for _, key := range op.keys {
			if err := dg.AddKeyPress(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{key}}, op.description); err != nil {
				return err
			}
		}
		if err := dg.CaptureFrame(op.description); err != nil {
			return err
		}
	}

	return nil
}

// GenerateKeyboardNavigationDemo generates keyboard navigation demo
func (dg *DemoGenerator) GenerateKeyboardNavigationDemo() error {
	if err := dg.StartRecording("keyboard_navigation", "Keyboard navigation demonstration"); err != nil {
		return err
	}
	defer dg.StopRecording()

	// Initial state
	if err := dg.CaptureFrame("Initial calculator state"); err != nil {
		return err
	}

	// Navigation keys
	navOps := []struct {
		keyType    tea.KeyType
		description string
	}{
		{tea.KeyTab, "Tab to focus first button"},
		{tea.KeyDown, "Navigate down"},
		{tea.KeyRight, "Navigate right"},
		{tea.KeyUp, "Navigate up"},
		{tea.KeyLeft, "Navigate left"},
		{tea.KeyEnter, "Press focused button"},
	}

	for _, op := range navOps {
		if err := dg.AddKeyPress(tea.KeyMsg{Type: op.keyType}, op.description); err != nil {
			return err
		}
		if err := dg.CaptureFrame(op.description); err != nil {
			return err
		}
	}

	return nil
}

// GenerateAllDemos generates all demo types
func (dg *DemoGenerator) GenerateAllDemos() error {
	// Create base directory
	baseDir := dg.OutputDir
	demos := []struct {
		name    string
		genFunc func() error
	}{
		{"basic", dg.GenerateBasicDemo},
		{"advanced", dg.GenerateAdvancedDemo},
		{"errors", dg.GenerateErrorDemo},
		{"navigation", dg.GenerateKeyboardNavigationDemo},
	}

	for _, demo := range demos {
		dg.OutputDir = filepath.Join(baseDir, demo.name)
		if err := demo.genFunc(); err != nil {
			return fmt.Errorf("failed to generate %s demo: %w", demo.name, err)
		}
	}

	return nil
}

// RenderDemoScript renders a demo script for playback
func (dg *DemoGenerator) RenderDemoScript(sequence *DemoSequence) (string, error) {
	var script strings.Builder

	script.WriteString(fmt.Sprintf("# Demo Script: %s\n", sequence.Name))
	script.WriteString(fmt.Sprintf("# %s\n", sequence.Description))
	script.WriteString(fmt.Sprintf("# Duration: %s\n\n", sequence.Duration))

	for i, action := range sequence.Actions {
		script.WriteString(fmt.Sprintf("# Step %d: %s\n", i+1, action.Description))

		switch action.Type {
		case "keypress":
			script.WriteString(fmt.Sprintf("keypress %s\n", keyToString(action.KeyPress)))
		case "mouseclick":
			script.WriteString(fmt.Sprintf("mouseclick %d %d\n", action.MouseClick.X, action.MouseClick.Y))
		case "screenshot":
			script.WriteString("screenshot\n")
		}

		script.WriteString(fmt.Sprintf("delay %s\n\n", action.Delay))
	}

	return script.String(), nil
}

// keyToString converts a KeyMsg to string representation
func keyToString(key tea.KeyMsg) string {
	switch key.Type {
	case tea.KeyEnter:
		return "enter"
	case tea.KeyTab:
		return "tab"
	case tea.KeySpace:
		return "space"
	case tea.KeyBackspace:
		return "backspace"
	case tea.KeyUp:
		return "up"
	case tea.KeyDown:
		return "down"
	case tea.KeyLeft:
		return "left"
	case tea.KeyRight:
		return "right"
	case tea.KeyEsc:
		return "escape"
	case tea.KeyRunes:
		if len(key.Runes) > 0 {
			return string(key.Runes[0])
		}
		return "runes_empty"
	default:
		return fmt.Sprintf("unknown_%d", int(key.Type))
	}
}