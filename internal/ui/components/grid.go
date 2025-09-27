package components

import (
	"errors"

	"github.com/charmbracelet/lipgloss"
)

var (
	ErrInvalidColumn = errors.New("invalid column index")
	ErrInvalidRow    = errors.New("invalid row index")
	ErrCellNotFound  = errors.New("cell not found")
)

// GridDimensions defines the size of the grid
type GridDimensions struct {
	Columns int
	Rows    int
}

// GridPosition represents a position in the grid
type GridPosition struct {
	Column int
	Row    int
}

// GridCell represents a single cell in the grid
type GridCell struct {
	Position GridPosition
	Content  string
	Style    lipgloss.Style
}

// GridLayout manages the layout and positioning of grid elements
type GridLayout struct {
	dimensions     GridDimensions
	cellWidth      int
	cellHeight     int
	padding        int
	spacing        int
	responsive     bool
	minWidth       int
	maxWidth       int
	centered       bool
	cells          map[GridPosition]*GridCell
	renderer       lipgloss.Style
	borderStyle    lipgloss.Style
	focusedStyle   lipgloss.Style
	pressedStyle   lipgloss.Style
}

// NewGridLayout creates a new grid layout manager
func NewGridLayout() *GridLayout {
	return &GridLayout{
		dimensions: GridDimensions{
			Columns: 4,
			Rows:    5,
		},
		cellWidth:    6,
		cellHeight:   3,
		padding:      1,
		spacing:      1,
		responsive:   true,
		minWidth:     60,
		maxWidth:     80,
		centered:     true,
		cells:        make(map[GridPosition]*GridCell),
		renderer:     lipgloss.NewStyle(),
		borderStyle:  lipgloss.NewStyle().Border(lipgloss.RoundedBorder()),
		focusedStyle: lipgloss.NewStyle().Background(lipgloss.Color("62")).Foreground(lipgloss.Color("15")),
		pressedStyle: lipgloss.NewStyle().Background(lipgloss.Color("94")).Foreground(lipgloss.Color("15")),
	}
}

// WithDimensions sets the grid dimensions
func (g *GridLayout) WithDimensions(cols, rows int) *GridLayout {
	g.dimensions.Columns = cols
	g.dimensions.Rows = rows
	return g
}

// WithCellSize sets the cell dimensions
func (g *GridLayout) WithCellSize(width, height int) *GridLayout {
	g.cellWidth = width
	g.cellHeight = height
	return g
}

// WithPadding sets the padding around cells
func (g *GridLayout) WithPadding(padding int) *GridLayout {
	g.padding = padding
	return g
}

// WithSpacing sets the spacing between cells
func (g *GridLayout) WithSpacing(spacing int) *GridLayout {
	g.spacing = spacing
	return g
}

// WithResponsive enables/disables responsive sizing
func (g *GridLayout) WithResponsive(responsive bool) *GridLayout {
	g.responsive = responsive
	return g
}

// WithMinMaxWidth sets the minimum and maximum width for responsive sizing
func (g *GridLayout) WithMinMaxWidth(min, max int) *GridLayout {
	g.minWidth = min
	g.maxWidth = max
	return g
}

// WithCentered sets whether the grid should be centered
func (g *GridLayout) WithCentered(centered bool) *GridLayout {
	g.centered = centered
	return g
}

// WithBorderStyle sets the border style for cells
func (g *GridLayout) WithBorderStyle(style lipgloss.Style) *GridLayout {
	g.borderStyle = style
	return g
}

// WithFocusedStyle sets the focused cell style
func (g *GridLayout) WithFocusedStyle(style lipgloss.Style) *GridLayout {
	g.focusedStyle = style
	return g
}

// WithPressedStyle sets the pressed cell style
func (g *GridLayout) WithPressedStyle(style lipgloss.Style) *GridLayout {
	g.pressedStyle = style
	return g
}

// AddCell adds a cell to the grid
func (g *GridLayout) AddCell(col, row int, content string, style lipgloss.Style) error {
	if col < 0 || col >= g.dimensions.Columns {
		return ErrInvalidColumn
	}
	if row < 0 || row >= g.dimensions.Rows {
		return ErrInvalidRow
	}

	pos := GridPosition{Column: col, Row: row}
	g.cells[pos] = &GridCell{
		Position: pos,
		Content:  content,
		Style:    style,
	}

	return nil
}

// GetCell retrieves a cell from the grid
func (g *GridLayout) GetCell(col, row int) (*GridCell, error) {
	pos := GridPosition{Column: col, Row: row}
	cell, exists := g.cells[pos]
	if !exists {
		return nil, ErrCellNotFound
	}
	return cell, nil
}

// RemoveCell removes a cell from the grid
func (g *GridLayout) RemoveCell(col, row int) error {
	pos := GridPosition{Column: col, Row: row}
	if _, exists := g.cells[pos]; !exists {
		return ErrCellNotFound
	}
	delete(g.cells, pos)
	return nil
}

// CalculateDimensions calculates the responsive dimensions based on terminal width
func (g *GridLayout) CalculateDimensions(termWidth int) (cellWidth, totalWidth int) {
	if !g.responsive {
		return g.cellWidth, g.calculateTotalWidth(g.cellWidth)
	}

	// Calculate optimal cell width based on terminal width
	optimalWidth := g.calculateOptimalCellWidth(termWidth)
	cellWidth = g.constrainWidth(optimalWidth)
	totalWidth = g.calculateTotalWidth(cellWidth)

	return cellWidth, totalWidth
}

// calculateOptimalCellWidth calculates the optimal cell width for the given terminal width
func (g *GridLayout) calculateOptimalCellWidth(termWidth int) int {
	availableWidth := termWidth - 2*g.padding // Account for padding
	totalSpacing := (g.dimensions.Columns - 1) * g.spacing
	availableForCells := availableWidth - totalSpacing

	if availableForCells <= 0 {
		return g.cellWidth // Fallback to default
	}

	return availableForCells / g.dimensions.Columns
}

// constrainWidth constrains the cell width within min/max bounds
func (g *GridLayout) constrainWidth(width int) int {
	if width < g.cellWidth {
		return g.cellWidth
	}
	if width > g.cellWidth*2 {
		return g.cellWidth * 2 // Don't make cells too large
	}
	return width
}

// calculateTotalWidth calculates the total width of the grid
func (g *GridLayout) calculateTotalWidth(cellWidth int) int {
	return (cellWidth * g.dimensions.Columns) +
		   (g.spacing * (g.dimensions.Columns - 1)) +
		   (2 * g.padding)
}

// GetCellPosition returns the screen position for a grid cell
func (g *GridLayout) GetCellPosition(col, row int, cellWidth int) (x, y int) {
	// Calculate horizontal position
	x = g.padding + (col * (cellWidth + g.spacing))

	// Calculate vertical position
	y = g.padding + (row * (g.cellHeight + g.spacing))

	// Center horizontally if requested
	if g.centered {
		totalWidth := g.calculateTotalWidth(cellWidth)
		centerX := (g.maxWidth - totalWidth) / 2
		if centerX > 0 {
			x += centerX
		}
	}

	return x, y
}

// Render renders the grid to a string
func (g *GridLayout) Render(termWidth int) string {
	cellWidth, totalWidth := g.CalculateDimensions(termWidth)

	// Build grid rows
	var rows []string
	for row := 0; row < g.dimensions.Rows; row++ {
		var rowCells []string
		for col := 0; col < g.dimensions.Columns; col++ {
			pos := GridPosition{Column: col, Row: row}
			cell, exists := g.cells[pos]

			var cellContent string
			var cellStyle lipgloss.Style

			if exists {
				cellContent = cell.Content
				cellStyle = cellStyle.
					Width(cellWidth).
					Height(g.cellHeight).
					Align(lipgloss.Center, lipgloss.Center).
					Border(lipgloss.RoundedBorder()).
					BorderForeground(lipgloss.Color("8"))

				// Merge with cell's style
				cellStyle = cellStyle.Inherit(cell.Style)
			} else {
				// Empty cell
				cellContent = ""
				cellStyle = lipgloss.NewStyle().
					Width(cellWidth).
					Height(g.cellHeight)
			}

			rowCells = append(rowCells, cellStyle.Render(cellContent))
		}

		// Join cells in row with spacing
		rowString := lipgloss.JoinHorizontal(lipgloss.Top, rowCells...)
		rows = append(rows, rowString)
	}

	// Join rows with vertical spacing
	gridContent := lipgloss.JoinVertical(lipgloss.Left, rows...)

	// Apply container styling
	containerStyle := lipgloss.NewStyle().
		Width(totalWidth).
		Padding(g.padding)

	return containerStyle.Render(gridContent)
}

// GetCellAtPosition returns the grid cell at the given screen position
func (g *GridLayout) GetCellAtPosition(x, y int, cellWidth int) (col, row int, found bool) {
	for rowIdx := 0; rowIdx < g.dimensions.Rows; rowIdx++ {
		for colIdx := 0; colIdx < g.dimensions.Columns; colIdx++ {
			cellX, cellY := g.GetCellPosition(colIdx, rowIdx, cellWidth)

			// Check if the click/touch is within this cell's bounds
			if x >= cellX && x < cellX+cellWidth &&
			   y >= cellY && y < cellY+g.cellHeight {
				return colIdx, rowIdx, true
			}
		}
	}

	return -1, -1, false
}

// GetDimensions returns the grid dimensions
func (g *GridLayout) GetDimensions() GridDimensions {
	return g.dimensions
}

// GetCellCount returns the number of cells in the grid
func (g *GridLayout) GetCellCount() int {
	return len(g.cells)
}

// Clear removes all cells from the grid
func (g *GridLayout) Clear() {
	g.cells = make(map[GridPosition]*GridCell)
}

// IsValidPosition checks if a position is valid within the grid
func (g *GridLayout) IsValidPosition(col, row int) bool {
	return col >= 0 && col < g.dimensions.Columns &&
	       row >= 0 && row < g.dimensions.Rows
}

// GetAdjacentPosition returns the position adjacent to the given position in the specified direction
func (g *GridLayout) GetAdjacentPosition(col, row int, direction Direction) (newCol, newRow int, valid bool) {
	switch direction {
	case DirectionUp:
		newCol, newRow = col, row-1
	case DirectionDown:
		newCol, newRow = col, row+1
	case DirectionLeft:
		newCol, newRow = col-1, row
	case DirectionRight:
		newCol, newRow = col+1, row
	default:
		return col, row, false
	}

	valid = g.IsValidPosition(newCol, newRow)
	return newCol, newRow, valid
}

// Direction represents navigation direction
type Direction int

const (
	DirectionUp Direction = iota
	DirectionDown
	DirectionLeft
	DirectionRight
)