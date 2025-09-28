package components

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewResponsiveLayout(t *testing.T) {
	layout := NewResponsiveLayout()

	assert.NotNil(t, layout)
	assert.Equal(t, LayoutResponsive, layout.config.Type)
	assert.Equal(t, 60, layout.config.MinWidth)
	assert.Equal(t, 80, layout.config.MaxWidth)
	assert.Equal(t, 70, layout.config.PreferredWidth)
	assert.Equal(t, 1, layout.config.CellPadding)
	assert.Equal(t, 1, layout.config.CellSpacing)
	assert.Equal(t, 1, layout.config.GridPadding)
	assert.True(t, layout.config.CenterHorizontally)
	assert.False(t, layout.config.CenterVertically)
	assert.True(t, layout.config.AllowOverflow)
	assert.NotNil(t, layout.cachedLayouts)
	assert.Equal(t, 1, layout.optimizationLevel)
}

func TestResponsiveLayout_WithConfig(t *testing.T) {
	layout := NewResponsiveLayout()
	config := LayoutConfig{
		Type:            LayoutCompact,
		MinWidth:        40,
		MaxWidth:        100,
		PreferredWidth:  80,
		CellPadding:     0,
		CellSpacing:     0,
		GridPadding:     0,
		CenterHorizontally: false,
		CenterVertically:   true,
		AllowOverflow:      false,
	}

	modified := layout.WithConfig(config)

	assert.Same(t, layout, modified)
	assert.Equal(t, config.Type, layout.config.Type)
	assert.Equal(t, config.MinWidth, layout.config.MinWidth)
	assert.Equal(t, config.MaxWidth, layout.config.MaxWidth)
	assert.Equal(t, config.PreferredWidth, layout.config.PreferredWidth)
	assert.Equal(t, config.CellPadding, layout.config.CellPadding)
	assert.Equal(t, config.CellSpacing, layout.config.CellSpacing)
	assert.Equal(t, config.GridPadding, layout.config.GridPadding)
	assert.Equal(t, config.CenterHorizontally, layout.config.CenterHorizontally)
	assert.Equal(t, config.CenterVertically, layout.config.CenterVertically)
	assert.Equal(t, config.AllowOverflow, layout.config.AllowOverflow)
}

func TestResponsiveLayout_WithLayoutType(t *testing.T) {
	layout := NewResponsiveLayout()
	modified := layout.WithLayoutType(LayoutFixed)

	assert.Same(t, layout, modified)
	assert.Equal(t, LayoutFixed, layout.config.Type)
}

func TestResponsiveLayout_WithMinMaxWidth(t *testing.T) {
	layout := NewResponsiveLayout()
	modified := layout.WithMinMaxWidth(50, 90)

	assert.Same(t, layout, modified)
	assert.Equal(t, 50, layout.config.MinWidth)
	assert.Equal(t, 90, layout.config.MaxWidth)
}

func TestResponsiveLayout_WithPadding(t *testing.T) {
	layout := NewResponsiveLayout()
	modified := layout.WithPadding(2, 0, 0)

	assert.Same(t, layout, modified)
	assert.Equal(t, 2, layout.config.CellPadding)
	assert.Equal(t, 0, layout.config.CellSpacing)
	assert.Equal(t, 0, layout.config.GridPadding)
}

func TestResponsiveLayout_WithCentering(t *testing.T) {
	layout := NewResponsiveLayout()
	modified := layout.WithCentering(false, true)

	assert.Same(t, layout, modified)
	assert.False(t, layout.config.CenterHorizontally)
	assert.True(t, layout.config.CenterVertically)
}

func TestResponsiveLayout_WithOptimization(t *testing.T) {
	layout := NewResponsiveLayout()
	modified := layout.WithOptimization(2)

	assert.Same(t, layout, modified)
	assert.Equal(t, 2, layout.optimizationLevel)
}

func TestResponsiveLayout_UpdateTermSize(t *testing.T) {
	layout := NewResponsiveLayout()
	assert.Equal(t, 0, layout.currentTermWidth)
	assert.Equal(t, 0, layout.currentTermHeight)

	layout.UpdateTermSize(100, 40)

	assert.Equal(t, 100, layout.currentTermWidth)
	assert.Equal(t, 40, layout.currentTermHeight)
}

func TestResponsiveLayout_UpdateTermSize_ClearsCache(t *testing.T) {
	layout := NewResponsiveLayout().WithOptimization(0) // Disable optimization to test cache clearing
	layout.cachedLayouts["test"] = &GridLayout{}

	layout.UpdateTermSize(100, 40)

	assert.Empty(t, layout.cachedLayouts)
}

func TestResponsiveLayout_CreateGridLayout(t *testing.T) {
	layout := NewResponsiveLayout()
	layout.UpdateTermSize(80, 30)

	grid := layout.CreateGridLayout(4, 5)

	require.NotNil(t, grid)
	assert.Equal(t, 4, grid.GetDimensions().Columns)
	assert.Equal(t, 5, grid.GetDimensions().Rows)
}

func TestResponsiveLayout_CreateGridLayout_Caching(t *testing.T) {
	layout := NewResponsiveLayout().WithOptimization(1)
	layout.UpdateTermSize(80, 30)

	// Create first grid
	grid1 := layout.CreateGridLayout(4, 5)

	// Create second grid with same parameters
	grid2 := layout.CreateGridLayout(4, 5)

	// Should be the same instance (cached)
	assert.Same(t, grid1, grid2)

	// Create grid with different parameters
	grid3 := layout.CreateGridLayout(3, 4)

	// Should be different instance
	assert.NotSame(t, grid1, grid3)
}

func TestResponsiveLayout_generateCacheKey(t *testing.T) {
	layout := NewResponsiveLayout()
	layout.UpdateTermSize(80, 30)

	key1 := layout.generateCacheKey(4, 5)
	key2 := layout.generateCacheKey(4, 5)
	key3 := layout.generateCacheKey(3, 4)

	assert.Equal(t, key1, key2)
	assert.NotEqual(t, key1, key3)
	// The cache key uses string(rune()) conversion, so we just verify it generates consistent keys
	assert.NotEmpty(t, key1)
	assert.NotEmpty(t, key2)
	assert.NotEmpty(t, key3)
}

func TestResponsiveLayout_applyFixedLayout(t *testing.T) {
	layout := NewResponsiveLayout()
	grid := NewGridLayout()

	layout.applyFixedLayout(grid, 4, 5)

	dimensions := grid.GetDimensions()
	assert.Equal(t, 4, dimensions.Columns)
	assert.Equal(t, 5, dimensions.Rows)
}

func TestResponsiveLayout_applyCompactLayout(t *testing.T) {
	layout := NewResponsiveLayout()
	layout.UpdateTermSize(80, 30)
	grid := NewGridLayout()

	layout.applyCompactLayout(grid, 4, 5)

	dimensions := grid.GetDimensions()
	assert.Equal(t, 4, dimensions.Columns)
	assert.Equal(t, 5, dimensions.Rows)
}

func TestResponsiveLayout_applyWideLayout(t *testing.T) {
	layout := NewResponsiveLayout()
	layout.UpdateTermSize(80, 30)
	grid := NewGridLayout()

	layout.applyWideLayout(grid, 4, 5)

	dimensions := grid.GetDimensions()
	assert.Equal(t, 4, dimensions.Columns)
	assert.Equal(t, 5, dimensions.Rows)
}

func TestResponsiveLayout_applyResponsiveLayout(t *testing.T) {
	layout := NewResponsiveLayout()
	layout.UpdateTermSize(80, 30)
	grid := NewGridLayout()

	layout.applyResponsiveLayout(grid, 4, 5)

	dimensions := grid.GetDimensions()
	assert.Equal(t, 4, dimensions.Columns)
	assert.Equal(t, 5, dimensions.Rows)
}

func TestResponsiveLayout_calculateOptimalCellWidth(t *testing.T) {
	layout := NewResponsiveLayout()
	layout.UpdateTermSize(80, 30)

	width := layout.calculateOptimalCellWidth(4, 1.0)
	assert.True(t, width > 0)
	assert.True(t, width <= 20) // Reasonable upper bound
}

func TestResponsiveLayout_calculateOptimalCellWidth_ZeroTermWidth(t *testing.T) {
	layout := NewResponsiveLayout()
	layout.UpdateTermSize(0, 30)

	width := layout.calculateOptimalCellWidth(4, 1.0)
	assert.Equal(t, 6, width) // Default fallback
}

func TestResponsiveLayout_calculateOptimalCellWidth_SmallTermWidth(t *testing.T) {
	layout := NewResponsiveLayout()
	layout.UpdateTermSize(30, 30)

	width := layout.calculateOptimalCellWidth(4, 1.0)
	assert.Equal(t, 6, width) // Default fallback (30-2-3-3)/4 = 22/4 = 5.5 -> 6
}

func TestResponsiveLayout_CalculateSpacing(t *testing.T) {
	layout := NewResponsiveLayout()

	// Test small terminal
	layout.UpdateTermSize(40, 30)
	spacing := layout.CalculateSpacing()
	assert.Equal(t, 0, spacing)

	// Test medium terminal
	layout.UpdateTermSize(60, 30)
	spacing = layout.CalculateSpacing()
	assert.Equal(t, 1, spacing)

	// Test large terminal
	layout.UpdateTermSize(100, 30)
	spacing = layout.CalculateSpacing()
	assert.Equal(t, 2, spacing)
}

func TestResponsiveLayout_CalculatePadding(t *testing.T) {
	layout := NewResponsiveLayout()

	// Test small terminal
	layout.UpdateTermSize(40, 30)
	padding := layout.CalculatePadding()
	assert.Equal(t, 0, padding)

	// Test medium terminal
	layout.UpdateTermSize(60, 30)
	padding = layout.CalculatePadding()
	assert.Equal(t, 2, padding) // Medium terminals get normal padding

	// Test large terminal
	layout.UpdateTermSize(100, 30)
	padding = layout.CalculatePadding()
	assert.Equal(t, 2, padding)
}

func TestResponsiveLayout_GetLayoutMetrics(t *testing.T) {
	layout := NewResponsiveLayout()
	layout.UpdateTermSize(80, 30)

	metrics := layout.GetLayoutMetrics()

	assert.Equal(t, 80, metrics.TerminalWidth)
	assert.Equal(t, 30, metrics.TerminalHeight)
	assert.Equal(t, LayoutResponsive, metrics.LayoutType)
	assert.True(t, metrics.CellWidth > 0)
	assert.True(t, metrics.CellHeight > 0)
	assert.True(t, metrics.TotalWidth > 0)
	assert.True(t, metrics.TotalHeight > 0)
}

func TestResponsiveLayout_calculateTotalWidth(t *testing.T) {
	layout := NewResponsiveLayout()

	width := layout.calculateTotalWidth(4)
	assert.True(t, width > 0)
}

func TestResponsiveLayout_calculateTotalHeight(t *testing.T) {
	layout := NewResponsiveLayout()

	height := layout.calculateTotalHeight(5)
	assert.True(t, height > 0)
}

func TestResponsiveLayout_OptimizeForTerminal_Small(t *testing.T) {
	layout := NewResponsiveLayout()

	layout.OptimizeForTerminal(40, 15)

	assert.Equal(t, 40, layout.currentTermWidth)
	assert.Equal(t, 15, layout.currentTermHeight)
	assert.Equal(t, 0, layout.optimizationLevel)
	assert.Equal(t, LayoutCompact, layout.config.Type)
}

func TestResponsiveLayout_OptimizeForTerminal_Medium(t *testing.T) {
	layout := NewResponsiveLayout()

	layout.OptimizeForTerminal(80, 30)

	assert.Equal(t, 80, layout.currentTermWidth)
	assert.Equal(t, 30, layout.currentTermHeight)
	assert.Equal(t, 1, layout.optimizationLevel)
	assert.Equal(t, LayoutResponsive, layout.config.Type)
}

func TestResponsiveLayout_OptimizeForTerminal_Large(t *testing.T) {
	layout := NewResponsiveLayout()

	layout.OptimizeForTerminal(140, 50)

	assert.Equal(t, 140, layout.currentTermWidth)
	assert.Equal(t, 50, layout.currentTermHeight)
	assert.Equal(t, 2, layout.optimizationLevel)
	assert.Equal(t, LayoutWide, layout.config.Type)
}

func TestNewLayoutManager(t *testing.T) {
	manager := NewLayoutManager()

	assert.NotNil(t, manager)
	assert.NotNil(t, manager.responsive)
	assert.NotNil(t, manager.statusBar)
	assert.NotNil(t, manager.titleBar)
	assert.Nil(t, manager.mainGrid)
}

func TestLayoutManager_Initialize(t *testing.T) {
	manager := NewLayoutManager()

	manager.Initialize(100, 40)

	assert.NotNil(t, manager.mainGrid)
	assert.Equal(t, 100, manager.statusBar.width)
	assert.Equal(t, 100, manager.titleBar.width)
	assert.Equal(t, 40, manager.responsive.currentTermHeight)
}

func TestLayoutManager_GetMainGrid(t *testing.T) {
	manager := NewLayoutManager()
	manager.Initialize(100, 40)

	grid := manager.GetMainGrid()

	assert.NotNil(t, grid)
	assert.Same(t, manager.mainGrid, grid)
}

func TestLayoutManager_GetResponsiveLayout(t *testing.T) {
	manager := NewLayoutManager()

	responsive := manager.GetResponsiveLayout()

	assert.NotNil(t, responsive)
	assert.Same(t, manager.responsive, responsive)
}

func TestLayoutManager_GetStatusBar(t *testing.T) {
	manager := NewLayoutManager()

	statusBar := manager.GetStatusBar()

	assert.NotNil(t, statusBar)
	assert.Same(t, manager.statusBar, statusBar)
}

func TestLayoutManager_GetTitleBar(t *testing.T) {
	manager := NewLayoutManager()

	titleBar := manager.GetTitleBar()

	assert.NotNil(t, titleBar)
	assert.Same(t, manager.titleBar, titleBar)
}

func TestLayoutManager_UpdateTermSize(t *testing.T) {
	manager := NewLayoutManager()
	manager.Initialize(100, 40)

	manager.UpdateTermSize(120, 50)

	assert.Equal(t, 120, manager.statusBar.width)
	assert.Equal(t, 120, manager.titleBar.width)
	assert.Equal(t, 50, manager.responsive.currentTermHeight)
}

func TestLayoutManager_GetMetrics(t *testing.T) {
	manager := NewLayoutManager()
	manager.Initialize(100, 40)

	metrics := manager.GetMetrics()

	assert.Equal(t, 100, metrics.TerminalWidth)
	assert.Equal(t, 40, metrics.TerminalHeight)
	assert.Equal(t, LayoutResponsive, metrics.LayoutType)
}

func TestStatusBarLayout(t *testing.T) {
	statusBar := &StatusBarLayout{
		width:   80,
		height:  1,
		visible: true,
	}

	assert.Equal(t, 80, statusBar.width)
	assert.Equal(t, 1, statusBar.height)
	assert.True(t, statusBar.visible)
}

func TestTitleBarLayout(t *testing.T) {
	titleBar := &TitleBarLayout{
		width:    80,
		height:   1,
		title:    "Calculator",
		subtitle: "Retro Edition",
		visible:  true,
	}

	assert.Equal(t, 80, titleBar.width)
	assert.Equal(t, 1, titleBar.height)
	assert.Equal(t, "Calculator", titleBar.title)
	assert.Equal(t, "Retro Edition", titleBar.subtitle)
	assert.True(t, titleBar.visible)
}