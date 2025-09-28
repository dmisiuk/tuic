package components

import (
	"github.com/charmbracelet/lipgloss"
)

// LayoutType defines different layout strategies
type LayoutType int

const (
	LayoutFixed LayoutType = iota
	LayoutResponsive
	LayoutCompact
	LayoutWide
)

// LayoutConfig defines configuration for layout behavior
type LayoutConfig struct {
	Type           LayoutType
	MinWidth       int
	MaxWidth       int
	PreferredWidth int
	CellPadding    int
	CellSpacing    int
	GridPadding    int
	CenterHorizontally bool
	CenterVertically   bool
	AllowOverflow     bool
}

// ResponsiveLayout handles dynamic layout adjustments based on terminal size
type ResponsiveLayout struct {
	config LayoutConfig
	currentTermWidth int
	currentTermHeight int
	cachedLayouts     map[string]*GridLayout
	optimizationLevel int
}

// NewResponsiveLayout creates a new responsive layout manager
func NewResponsiveLayout() *ResponsiveLayout {
	return &ResponsiveLayout{
		config: LayoutConfig{
			Type:            LayoutResponsive,
			MinWidth:        60,
			MaxWidth:        80,
			PreferredWidth:  70,
			CellPadding:     1,
			CellSpacing:     1,
			GridPadding:     1,
			CenterHorizontally: true,
			CenterVertically:   false,
			AllowOverflow:      true,
		},
		cachedLayouts: make(map[string]*GridLayout),
		optimizationLevel: 1, // 0 = minimal, 1 = balanced, 2 = aggressive
	}
}

// WithConfig sets the layout configuration
func (rl *ResponsiveLayout) WithConfig(config LayoutConfig) *ResponsiveLayout {
	rl.config = config
	return rl
}

// WithLayoutType sets the layout type
func (rl *ResponsiveLayout) WithLayoutType(layoutType LayoutType) *ResponsiveLayout {
	rl.config.Type = layoutType
	return rl
}

// WithMinMaxWidth sets the minimum and maximum width constraints
func (rl *ResponsiveLayout) WithMinMaxWidth(min, max int) *ResponsiveLayout {
	rl.config.MinWidth = min
	rl.config.MaxWidth = max
	return rl
}

// WithPadding sets padding values
func (rl *ResponsiveLayout) WithPadding(cell, spacing, grid int) *ResponsiveLayout {
	rl.config.CellPadding = cell
	rl.config.CellSpacing = spacing
	rl.config.GridPadding = grid
	return rl
}

// WithCentering sets centering options
func (rl *ResponsiveLayout) WithCentering(horizontal, vertical bool) *ResponsiveLayout {
	rl.config.CenterHorizontally = horizontal
	rl.config.CenterVertically = vertical
	return rl
}

// WithOptimization sets the optimization level
func (rl *ResponsiveLayout) WithOptimization(level int) *ResponsiveLayout {
	rl.optimizationLevel = level
	return rl
}

// UpdateTermSize updates the terminal size information
func (rl *ResponsiveLayout) UpdateTermSize(width, height int) {
	if rl.currentTermWidth != width || rl.currentTermHeight != height {
		rl.currentTermWidth = width
		rl.currentTermHeight = height
		// Clear cache when terminal size changes
		if rl.optimizationLevel < 2 {
			rl.cachedLayouts = make(map[string]*GridLayout)
		}
	}
}

// CreateGridLayout creates a grid layout optimized for the current terminal size
func (rl *ResponsiveLayout) CreateGridLayout(cols, rows int) *GridLayout {
	cacheKey := rl.generateCacheKey(cols, rows)

	// Check cache first
	if cached, exists := rl.cachedLayouts[cacheKey]; exists && rl.optimizationLevel > 0 {
		return cached
	}

	grid := NewGridLayout()

	// Apply configuration based on layout type
	switch rl.config.Type {
	case LayoutFixed:
		rl.applyFixedLayout(grid, cols, rows)
	case LayoutCompact:
		rl.applyCompactLayout(grid, cols, rows)
	case LayoutWide:
		rl.applyWideLayout(grid, cols, rows)
	case LayoutResponsive:
		rl.applyResponsiveLayout(grid, cols, rows)
	}

	// Apply common settings
	grid.
		WithPadding(rl.config.CellPadding).
		WithSpacing(rl.config.CellSpacing).
		WithCentered(rl.config.CenterHorizontally).
		WithResponsive(true).
		WithMinMaxWidth(rl.config.MinWidth, rl.config.MaxWidth)

	// Cache the result
	if rl.optimizationLevel > 0 {
		rl.cachedLayouts[cacheKey] = grid
	}

	return grid
}

// generateCacheKey generates a unique key for caching layouts
func (rl *ResponsiveLayout) generateCacheKey(cols, rows int) string {
	return string(rune(cols)) + "_" + string(rune(rows)) + "_" +
		   string(rune(rl.currentTermWidth)) + "_" +
		   string(rune(rl.config.Type))
}

// applyFixedLayout applies fixed layout settings
func (rl *ResponsiveLayout) applyFixedLayout(grid *GridLayout, cols, rows int) {
	grid.
		WithDimensions(cols, rows).
		WithCellSize(6, 3) // Fixed cell size
}

// applyCompactLayout applies compact layout settings
func (rl *ResponsiveLayout) applyCompactLayout(grid *GridLayout, cols, rows int) {
	cellWidth := rl.calculateOptimalCellWidth(cols, 0.8) // 20% smaller
	cellHeight := 2 // Smaller height

	grid.
		WithDimensions(cols, rows).
		WithCellSize(cellWidth, cellHeight).
		WithSpacing(0) // No spacing between cells
}

// applyWideLayout applies wide layout settings
func (rl *ResponsiveLayout) applyWideLayout(grid *GridLayout, cols, rows int) {
	cellWidth := rl.calculateOptimalCellWidth(cols, 1.2) // 20% larger
	cellHeight := 4 // Taller cells

	grid.
		WithDimensions(cols, rows).
		WithCellSize(cellWidth, cellHeight).
		WithSpacing(2) // More spacing
}

// applyResponsiveLayout applies responsive layout settings
func (rl *ResponsiveLayout) applyResponsiveLayout(grid *GridLayout, cols, rows int) {
	cellWidth := rl.calculateOptimalCellWidth(cols, 1.0)
	cellHeight := 3 // Standard height

	// Adjust cell height based on terminal height
	if rl.currentTermHeight < 25 {
		cellHeight = 2 // Compact height for small terminals
	} else if rl.currentTermHeight > 35 {
		cellHeight = 4 // Taller height for large terminals
	}

	grid.
		WithDimensions(cols, rows).
		WithCellSize(cellWidth, cellHeight)
}

// calculateOptimalCellWidth calculates the optimal cell width
func (rl *ResponsiveLayout) calculateOptimalCellWidth(cols int, scale float64) int {
	if rl.currentTermWidth == 0 {
		return 6 // Default fallback
	}

	availableWidth := rl.currentTermWidth - (2 * rl.config.GridPadding)
	totalSpacing := (cols - 1) * rl.config.CellSpacing
	availableForCells := availableWidth - totalSpacing

	if availableForCells <= 0 {
		return 4 // Minimum viable width
	}

	baseWidth := availableForCells / cols
	scaledWidth := int(float64(baseWidth) * scale)

	// Apply constraints
	if scaledWidth < 4 {
		return 4
	}
	if scaledWidth > 12 {
		return 12
	}

	return scaledWidth
}

// CalculateSpacing calculates optimal spacing based on terminal size
func (rl *ResponsiveLayout) CalculateSpacing() int {
	if rl.currentTermWidth < 50 {
		return 0 // No spacing for very small terminals
	} else if rl.currentTermWidth < 70 {
		return 1 // Minimal spacing
	} else {
		return 2 // Normal spacing for larger terminals
	}
}

// CalculatePadding calculates optimal padding based on terminal size
func (rl *ResponsiveLayout) CalculatePadding() int {
	if rl.currentTermWidth < 50 {
		return 0 // No padding for very small terminals
	} else if rl.currentTermWidth < 60 {
		return 1 // Minimal padding
	} else {
		return 2 // Normal padding
	}
}

// GetLayoutMetrics returns metrics about the current layout
func (rl *ResponsiveLayout) GetLayoutMetrics() LayoutMetrics {
	return LayoutMetrics{
		TerminalWidth:  rl.currentTermWidth,
		TerminalHeight: rl.currentTermHeight,
		LayoutType:     rl.config.Type,
		CellWidth:      rl.calculateOptimalCellWidth(4, 1.0),
		CellHeight:     3,
		Spacing:        rl.CalculateSpacing(),
		Padding:        rl.CalculatePadding(),
		TotalWidth:     rl.calculateTotalWidth(4),
		TotalHeight:    rl.calculateTotalHeight(5),
	}
}

// calculateTotalWidth calculates the total width of the grid
func (rl *ResponsiveLayout) calculateTotalWidth(cols int) int {
	cellWidth := rl.calculateOptimalCellWidth(cols, 1.0)
	return (cellWidth * cols) +
		   (rl.config.CellSpacing * (cols - 1)) +
		   (2 * rl.config.GridPadding)
}

// calculateTotalHeight calculates the total height of the grid
func (rl *ResponsiveLayout) calculateTotalHeight(rows int) int {
	cellHeight := 3 // Standard height
	if rl.currentTermHeight < 25 {
		cellHeight = 2
	} else if rl.currentTermHeight > 35 {
		cellHeight = 4
	}

	return (cellHeight * rows) +
		   (rl.config.CellSpacing * (rows - 1)) +
		   (2 * rl.config.GridPadding)
}

// OptimizeForTerminal applies optimizations based on terminal capabilities
func (rl *ResponsiveLayout) OptimizeForTerminal(termWidth, termHeight int) {
	rl.UpdateTermSize(termWidth, termHeight)

	// Adjust optimization level based on terminal size
	if termWidth < 50 || termHeight < 20 {
		rl.optimizationLevel = 0 // Minimal optimization for small terminals
	} else if termWidth > 120 && termHeight > 40 {
		rl.optimizationLevel = 2 // Aggressive optimization for large terminals
	} else {
		rl.optimizationLevel = 1 // Balanced optimization for medium terminals
	}

	// Adjust layout type based on terminal size
	if termWidth < 50 {
		rl.config.Type = LayoutCompact
	} else if termWidth > 100 {
		rl.config.Type = LayoutWide
	} else {
		rl.config.Type = LayoutResponsive
	}
}

// LayoutMetrics provides information about the current layout
type LayoutMetrics struct {
	TerminalWidth  int
	TerminalHeight int
	LayoutType     LayoutType
	CellWidth      int
	CellHeight     int
	Spacing        int
	Padding        int
	TotalWidth     int
	TotalHeight    int
}

// LayoutManager coordinates multiple layout components
type LayoutManager struct {
	responsive *ResponsiveLayout
	mainGrid   *GridLayout
	statusBar  *StatusBarLayout
	titleBar   *TitleBarLayout
}

// StatusBarLayout manages the status bar layout
type StatusBarLayout struct {
	width    int
	height   int
	style    lipgloss.Style
	content  string
	visible  bool
}

// TitleBarLayout manages the title bar layout
type TitleBarLayout struct {
	width    int
	height   int
	style    lipgloss.Style
	title    string
	subtitle string
	visible  bool
}

// NewLayoutManager creates a new layout manager
func NewLayoutManager() *LayoutManager {
	return &LayoutManager{
		responsive: NewResponsiveLayout(),
		statusBar: &StatusBarLayout{
			height:  1,
			visible: true,
		},
		titleBar: &TitleBarLayout{
			height:  1,
			visible: true,
		},
	}
}

// Initialize initializes the layout manager with terminal dimensions
func (lm *LayoutManager) Initialize(termWidth, termHeight int) {
	lm.responsive.OptimizeForTerminal(termWidth, termHeight)

	// Create main grid (4x5 for calculator)
	lm.mainGrid = lm.responsive.CreateGridLayout(4, 5)

	// Initialize status bar
	lm.statusBar.width = termWidth
	lm.statusBar.style = lipgloss.NewStyle().
		Background(lipgloss.Color("240")).
		Foreground(lipgloss.Color("15")).
		Width(termWidth)

	// Initialize title bar
	lm.titleBar.width = termWidth
	lm.titleBar.style = lipgloss.NewStyle().
		Background(lipgloss.Color("24")).
		Foreground(lipgloss.Color("15")).
		Bold(true).
		Width(termWidth)
}

// GetMainGrid returns the main grid layout
func (lm *LayoutManager) GetMainGrid() *GridLayout {
	return lm.mainGrid
}

// GetResponsiveLayout returns the responsive layout manager
func (lm *LayoutManager) GetResponsiveLayout() *ResponsiveLayout {
	return lm.responsive
}

// GetStatusBar returns the status bar layout
func (lm *LayoutManager) GetStatusBar() *StatusBarLayout {
	return lm.statusBar
}

// GetTitleBar returns the title bar layout
func (lm *LayoutManager) GetTitleBar() *TitleBarLayout {
	return lm.titleBar
}

// UpdateTermSize updates all layouts with new terminal dimensions
func (lm *LayoutManager) UpdateTermSize(width, height int) {
	lm.Initialize(width, height)
}

// GetMetrics returns layout metrics
func (lm *LayoutManager) GetMetrics() LayoutMetrics {
	return lm.responsive.GetLayoutMetrics()
}