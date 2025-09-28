package components

import (
	"testing"
	"github.com/charmbracelet/lipgloss"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewGridLayout(t *testing.T) {
	grid := NewGridLayout()

	assert.NotNil(t, grid)
	assert.Equal(t, 4, grid.dimensions.Columns)
	assert.Equal(t, 5, grid.dimensions.Rows)
	assert.Equal(t, 6, grid.cellWidth)
	assert.Equal(t, 3, grid.cellHeight)
	assert.Equal(t, 1, grid.padding)
	assert.Equal(t, 1, grid.spacing)
	assert.True(t, grid.responsive)
	assert.True(t, grid.centered)
	assert.Equal(t, 60, grid.minWidth)
	assert.Equal(t, 80, grid.maxWidth)
	assert.NotNil(t, grid.cells)
	assert.Empty(t, grid.cells)
}

func TestGridLayout_WithDimensions(t *testing.T) {
	grid := NewGridLayout()
	modified := grid.WithDimensions(3, 4)

	assert.Same(t, grid, modified)
	assert.Equal(t, 3, grid.dimensions.Columns)
	assert.Equal(t, 4, grid.dimensions.Rows)
}

func TestGridLayout_WithCellSize(t *testing.T) {
	grid := NewGridLayout()
	modified := grid.WithCellSize(8, 4)

	assert.Same(t, grid, modified)
	assert.Equal(t, 8, grid.cellWidth)
	assert.Equal(t, 4, grid.cellHeight)
}

func TestGridLayout_WithPadding(t *testing.T) {
	grid := NewGridLayout()
	modified := grid.WithPadding(2)

	assert.Same(t, grid, modified)
	assert.Equal(t, 2, grid.padding)
}

func TestGridLayout_WithSpacing(t *testing.T) {
	grid := NewGridLayout()
	modified := grid.WithSpacing(2)

	assert.Same(t, grid, modified)
	assert.Equal(t, 2, grid.spacing)
}

func TestGridLayout_WithResponsive(t *testing.T) {
	grid := NewGridLayout()
	modified := grid.WithResponsive(false)

	assert.Same(t, grid, modified)
	assert.False(t, grid.responsive)
}

func TestGridLayout_WithMinMaxWidth(t *testing.T) {
	grid := NewGridLayout()
	modified := grid.WithMinMaxWidth(50, 90)

	assert.Same(t, grid, modified)
	assert.Equal(t, 50, grid.minWidth)
	assert.Equal(t, 90, grid.maxWidth)
}

func TestGridLayout_WithCentered(t *testing.T) {
	grid := NewGridLayout()
	modified := grid.WithCentered(false)

	assert.Same(t, grid, modified)
	assert.False(t, grid.centered)
}

func TestGridLayout_AddCell(t *testing.T) {
	tests := []struct {
		name      string
		col, row  int
		content   string
		wantError error
	}{
		{
			name:     "valid position",
			col:      1,
			row:      1,
			content:  "test",
			wantError: nil,
		},
		{
			name:      "invalid column - negative",
			col:       -1,
			row:       1,
			content:   "test",
			wantError: ErrInvalidColumn,
		},
		{
			name:      "invalid column - too large",
			col:       4,
			row:       1,
			content:   "test",
			wantError: ErrInvalidColumn,
		},
		{
			name:      "invalid row - negative",
			col:       1,
			row:       -1,
			content:   "test",
			wantError: ErrInvalidRow,
		},
		{
			name:      "invalid row - too large",
			col:       1,
			row:       5,
			content:   "test",
			wantError: ErrInvalidRow,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			grid := NewGridLayout()
			style := lipgloss.NewStyle()

			err := grid.AddCell(tt.col, tt.row, tt.content, style)

			if tt.wantError != nil {
				assert.ErrorIs(t, err, tt.wantError)
			} else {
				assert.NoError(t, err)
				cell, err := grid.GetCell(tt.col, tt.row)
				require.NoError(t, err)
				assert.Equal(t, tt.content, cell.Content)
			}
		})
	}
}

func TestGridLayout_GetCell(t *testing.T) {
	grid := NewGridLayout()
	style := lipgloss.NewStyle()

	// Test getting non-existent cell
	_, err := grid.GetCell(0, 0)
	assert.ErrorIs(t, err, ErrCellNotFound)

	// Add a cell and retrieve it
	err = grid.AddCell(0, 0, "test", style)
	require.NoError(t, err)

	cell, err := grid.GetCell(0, 0)
	require.NoError(t, err)
	assert.Equal(t, "test", cell.Content)
	assert.Equal(t, 0, cell.Position.Column)
	assert.Equal(t, 0, cell.Position.Row)
}

func TestGridLayout_RemoveCell(t *testing.T) {
	grid := NewGridLayout()
	style := lipgloss.NewStyle()

	// Test removing non-existent cell
	err := grid.RemoveCell(0, 0)
	assert.ErrorIs(t, err, ErrCellNotFound)

	// Add a cell and remove it
	err = grid.AddCell(0, 0, "test", style)
	require.NoError(t, err)

	err = grid.RemoveCell(0, 0)
	assert.NoError(t, err)

	// Verify cell is removed
	_, err = grid.GetCell(0, 0)
	assert.ErrorIs(t, err, ErrCellNotFound)
}

func TestGridLayout_CalculateDimensions(t *testing.T) {
	tests := []struct {
		name       string
		responsive bool
		termWidth  int
		wantWidth  int
		wantTotal  int
	}{
		{
			name:       "non-responsive",
			responsive: false,
			termWidth:  100,
			wantWidth:  6,
			wantTotal:  29, // (6*4) + (1*3) + (2*1) = 24 + 3 + 2 = 29
		},
		{
			name:       "responsive - small terminal",
			responsive: true,
			termWidth:  50,
			wantWidth:  11, // (50-2-3)/4 = 45/4 = 11.25 -> 11
			wantTotal:  49, // (11*4) + (1*3) + (2*1) = 44 + 3 + 2 = 49
		},
		{
			name:       "responsive - medium terminal",
			responsive: true,
			termWidth:  80,
			wantWidth:  12, // (80-2-3)/4 = 75/4 = 18.75 -> constrained to 12
			wantTotal:  53, // (12*4) + (1*3) + (2*1) = 48 + 3 + 2 = 53
		},
		{
			name:       "responsive - large terminal",
			responsive: true,
			termWidth:  120,
			wantWidth:  12, // max constraint applied
			wantTotal:  53, // (12*4) + (1*3) + (2*1) = 48 + 3 + 2 = 53
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			grid := NewGridLayout().WithResponsive(tt.responsive)
			cellWidth, totalWidth := grid.CalculateDimensions(tt.termWidth)

			assert.Equal(t, tt.wantWidth, cellWidth)
			assert.Equal(t, tt.wantTotal, totalWidth)
		})
	}
}

func TestGridLayout_calculateOptimalCellWidth(t *testing.T) {
	grid := NewGridLayout()

	tests := []struct {
		name      string
		termWidth int
		wantWidth int
	}{
		{
			name:      "very small terminal",
			termWidth: 30,
			wantWidth: 6, // fallback to default
		},
		{
			name:      "small terminal",
			termWidth: 50,
			wantWidth: 11, // (50-2-3)/4 = 45/4 = 11.25 -> 11
		},
		{
			name:      "medium terminal",
			termWidth: 80,
			wantWidth: 18, // (80-2-3)/4 = 75/4 = 18.75 -> 18
		},
		{
			name:      "large terminal",
			termWidth: 120,
			wantWidth: 28, // (120-2-3)/4 = 115/4 = 28.75 -> 28 (constrainWidth not applied in this test)
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			width := grid.calculateOptimalCellWidth(tt.termWidth)
			assert.Equal(t, tt.wantWidth, width)
		})
	}
}

func TestGridLayout_calculateTotalWidth(t *testing.T) {
	grid := NewGridLayout()

	tests := []struct {
		name      string
		cellWidth int
		wantTotal int
	}{
		{
			name:      "default cell width",
			cellWidth: 6,
			wantTotal: 29, // (6*4) + (1*3) + (2*1) = 24 + 3 + 2 = 29
		},
		{
			name:      "larger cell width",
			cellWidth: 8,
			wantTotal: 37, // (8*4) + (1*3) + (2*1) = 32 + 3 + 2 = 37
		},
		{
			name:      "smaller cell width",
			cellWidth: 4,
			wantTotal: 21, // (4*4) + (1*3) + (2*1) = 16 + 3 + 2 = 21
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			total := grid.calculateTotalWidth(tt.cellWidth)
			assert.Equal(t, tt.wantTotal, total)
		})
	}
}

func TestGridLayout_GetCellPosition(t *testing.T) {
	// Test without centering first
	grid := NewGridLayout().WithCentered(false)

	tests := []struct {
		name        string
		col, row    int
		cellWidth   int
		wantX, wantY int
	}{
		{
			name:     "top-left cell",
			col:      0,
			row:      0,
			cellWidth: 6,
			wantX:    1, // padding
			wantY:    1, // padding
		},
		{
			name:     "top-right cell",
			col:      3,
			row:      0,
			cellWidth: 6,
			wantX:    22, // 1 + (3 * (6 + 1)) = 1 + 21 = 22
			wantY:    1,  // padding
		},
		{
			name:     "bottom-left cell",
			col:      0,
			row:      4,
			cellWidth: 6,
			wantX:    1, // padding
			wantY:    17, // 1 + (4 * (3 + 1)) = 1 + 16 = 17
		},
		{
			name:     "bottom-right cell",
			col:      3,
			row:      4,
			cellWidth: 6,
			wantX:    22, // 1 + (3 * (6 + 1)) = 1 + 21 = 22
			wantY:    17, // 1 + (4 * (3 + 1)) = 1 + 16 = 17
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			x, y := grid.GetCellPosition(tt.col, tt.row, tt.cellWidth)
			assert.Equal(t, tt.wantX, x)
			assert.Equal(t, tt.wantY, y)
		})
	}
}

func TestGridLayout_GetCellAtPosition(t *testing.T) {
	// Test without centering to avoid complex calculations
	grid := NewGridLayout().WithCentered(false)
	cellWidth := 6

	tests := []struct {
		name           string
		x, y           int
		wantCol, wantRow int
		wantFound      bool
	}{
		{
			name:      "within first cell",
			x:         1,
			y:         1,
			wantCol:   0,
			wantRow:   0,
			wantFound: true,
		},
		{
			name:      "within last cell",
			x:         22,
			y:         17,
			wantCol:   3,
			wantRow:   4,
			wantFound: true,
		},
		{
			name:      "between cells horizontally",
			x:         7,
			y:         1,
			wantCol:   -1,
			wantRow:   -1,
			wantFound: false,
		},
		{
			name:      "between cells vertically",
			x:         1,
			y:         4,
			wantCol:   -1,
			wantRow:   -1,
			wantFound: false,
		},
		{
			name:      "outside grid bounds",
			x:         100,
			y:         100,
			wantCol:   -1,
			wantRow:   -1,
			wantFound: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			col, row, found := grid.GetCellAtPosition(tt.x, tt.y, cellWidth)
			assert.Equal(t, tt.wantCol, col)
			assert.Equal(t, tt.wantRow, row)
			assert.Equal(t, tt.wantFound, found)
		})
	}
}

func TestGridLayout_GetDimensions(t *testing.T) {
	grid := NewGridLayout()
	dimensions := grid.GetDimensions()

	assert.Equal(t, 4, dimensions.Columns)
	assert.Equal(t, 5, dimensions.Rows)
}

func TestGridLayout_GetCellCount(t *testing.T) {
	grid := NewGridLayout()
	style := lipgloss.NewStyle()

	// Initially empty
	assert.Equal(t, 0, grid.GetCellCount())

	// Add some cells
	grid.AddCell(0, 0, "test1", style)
	assert.Equal(t, 1, grid.GetCellCount())

	grid.AddCell(1, 1, "test2", style)
	assert.Equal(t, 2, grid.GetCellCount())

	// Remove a cell
	grid.RemoveCell(0, 0)
	assert.Equal(t, 1, grid.GetCellCount())
}

func TestGridLayout_Clear(t *testing.T) {
	grid := NewGridLayout()
	style := lipgloss.NewStyle()

	// Add some cells
	grid.AddCell(0, 0, "test1", style)
	grid.AddCell(1, 1, "test2", style)
	assert.Equal(t, 2, grid.GetCellCount())

	// Clear all cells
	grid.Clear()
	assert.Equal(t, 0, grid.GetCellCount())
	assert.Empty(t, grid.cells)
}

func TestGridLayout_IsValidPosition(t *testing.T) {
	grid := NewGridLayout()

	tests := []struct {
		name      string
		col, row  int
		wantValid bool
	}{
		{
			name:      "valid position",
			col:       0,
			row:       0,
			wantValid: true,
		},
		{
			name:      "valid edge position",
			col:       3,
			row:       4,
			wantValid: true,
		},
		{
			name:      "invalid negative column",
			col:       -1,
			row:       0,
			wantValid: false,
		},
		{
			name:      "invalid negative row",
			col:       0,
			row:       -1,
			wantValid: false,
		},
		{
			name:      "invalid column too large",
			col:       4,
			row:       0,
			wantValid: false,
		},
		{
			name:      "invalid row too large",
			col:       0,
			row:       5,
			wantValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			valid := grid.IsValidPosition(tt.col, tt.row)
			assert.Equal(t, tt.wantValid, valid)
		})
	}
}

func TestGridLayout_GetAdjacentPosition(t *testing.T) {
	grid := NewGridLayout()

	tests := []struct {
		name           string
		startCol, startRow int
		direction      Direction
		wantCol, wantRow  int
		wantValid       bool
	}{
		{
			name:      "up from center",
			startCol:  2,
			startRow:  2,
			direction: DirectionUp,
			wantCol:   2,
			wantRow:   1,
			wantValid: true,
		},
		{
			name:      "down from center",
			startCol:  2,
			startRow:  2,
			direction: DirectionDown,
			wantCol:   2,
			wantRow:   3,
			wantValid: true,
		},
		{
			name:      "left from center",
			startCol:  2,
			startRow:  2,
			direction: DirectionLeft,
			wantCol:   1,
			wantRow:   2,
			wantValid: true,
		},
		{
			name:      "right from center",
			startCol:  2,
			startRow:  2,
			direction: DirectionRight,
			wantCol:   3,
			wantRow:   2,
			wantValid: true,
		},
		{
			name:      "up from top edge",
			startCol:  2,
			startRow:  0,
			direction: DirectionUp,
			wantCol:   2,
			wantRow:   -1,
			wantValid: false,
		},
		{
			name:      "down from bottom edge",
			startCol:  2,
			startRow:  4,
			direction: DirectionDown,
			wantCol:   2,
			wantRow:   5,
			wantValid: false,
		},
		{
			name:      "left from left edge",
			startCol:  0,
			startRow:  2,
			direction: DirectionLeft,
			wantCol:   -1,
			wantRow:   2,
			wantValid: false,
		},
		{
			name:      "right from right edge",
			startCol:  3,
			startRow:  2,
			direction: DirectionRight,
			wantCol:   4,
			wantRow:   2,
			wantValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			col, row, valid := grid.GetAdjacentPosition(tt.startCol, tt.startRow, tt.direction)
			assert.Equal(t, tt.wantCol, col)
			assert.Equal(t, tt.wantRow, row)
			assert.Equal(t, tt.wantValid, valid)
		})
	}
}

func TestGridLayout_Render(t *testing.T) {
	grid := NewGridLayout()
	style := lipgloss.NewStyle()

	// Add some test cells
	grid.AddCell(0, 0, "1", style)
	grid.AddCell(1, 0, "2", style)
	grid.AddCell(2, 0, "3", style)
	grid.AddCell(3, 0, "+", style)

	// Test rendering
	output := grid.Render(80)

	assert.NotEmpty(t, output)
	assert.Contains(t, output, "1")
	assert.Contains(t, output, "2")
	assert.Contains(t, output, "3")
	assert.Contains(t, output, "+")
}

func TestGridLayout_RenderEmpty(t *testing.T) {
	grid := NewGridLayout()

	// Test rendering empty grid
	output := grid.Render(80)

	assert.NotEmpty(t, output)
	// Should contain grid structure but no cell content
}