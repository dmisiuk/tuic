package input

import (
	"ccpm-demo/internal/ui"
)

// Focusable represents an element that can receive focus
type Focusable interface {
	// GetID returns a unique identifier for this focusable element
	GetID() string

	// GetPosition returns the position of the element in a grid layout
	GetPosition() (row, col int)

	// GetLabel returns the display label for the element
	GetLabel() string

	// IsEnabled returns true if the element can receive focus
	IsEnabled() bool

	// OnFocus is called when the element receives focus
	OnFocus()

	// OnBlur is called when the element loses focus
	OnBlur()

	// Activate is called when the focused element is activated (space/enter)
	Activate() (ui.Model, error)
}

// FocusManager manages focus state and navigation
type FocusManager struct {
	focusables []Focusable
	currentIndex int
	wrapNavigation bool

	// Focus history for restoration
	focusHistory []string
}

// NewFocusManager creates a new focus manager
func NewFocusManager() *FocusManager {
	return &FocusManager{
		focusables:      []Focusable{},
		currentIndex:    -1,
		wrapNavigation:  true,
		focusHistory:    []string{},
	}
}

// AddFocusable adds a focusable element to the manager
func (fm *FocusManager) AddFocusable(element Focusable) {
	fm.focusables = append(fm.focusables, element)
}

// RemoveFocusable removes a focusable element from the manager
func (fm *FocusManager) RemoveFocusable(id string) {
	for i, element := range fm.focusables {
		if element.GetID() == id {
			fm.focusables = append(fm.focusables[:i], fm.focusables[i+1:]...)

			// Adjust current index if needed
			if i < fm.currentIndex {
				fm.currentIndex--
			} else if i == fm.currentIndex {
				fm.currentIndex = -1
			}
			break
		}
	}
}

// SetFocusables sets the complete list of focusable elements
func (fm *FocusManager) SetFocusables(elements []Focusable) {
	fm.focusables = elements
	if fm.currentIndex >= len(fm.focusables) {
		fm.currentIndex = len(fm.focusables) - 1
	}
}

// SetFocus sets focus to a specific element by ID
func (fm *FocusManager) SetFocus(id string) bool {
	for i, element := range fm.focusables {
		if element.GetID() == id && element.IsEnabled() {
			fm.setFocusByIndex(i)
			return true
		}
	}
	return false
}

// SetFocusByPosition sets focus to the element at the specified grid position
func (fm *FocusManager) SetFocusByPosition(row, col int) bool {
	for i, element := range fm.focusables {
		if element.IsEnabled() {
			elemRow, elemCol := element.GetPosition()
			if elemRow == row && elemCol == col {
				fm.setFocusByIndex(i)
				return true
			}
		}
	}
	return false
}

// setFocusByIndex sets focus by index and handles focus events
func (fm *FocusManager) setFocusByIndex(index int) {
	// Blur current focused element
	if fm.currentIndex >= 0 && fm.currentIndex < len(fm.focusables) {
		fm.focusables[fm.currentIndex].OnBlur()
	}

	// Save to history
	if fm.currentIndex >= 0 && fm.currentIndex < len(fm.focusables) {
		currentID := fm.focusables[fm.currentIndex].GetID()
		fm.focusHistory = append(fm.focusHistory, currentID)
		// Keep history manageable
		if len(fm.focusHistory) > 10 {
			fm.focusHistory = fm.focusHistory[1:]
		}
	}

	// Set new focus
	fm.currentIndex = index

	// Focus new element
	if fm.currentIndex >= 0 && fm.currentIndex < len(fm.focusables) {
		fm.focusables[fm.currentIndex].OnFocus()
	}
}

// GetFocusedElement returns the currently focused element
func (fm *FocusManager) GetFocusedElement() Focusable {
	if fm.currentIndex >= 0 && fm.currentIndex < len(fm.focusables) {
		return fm.focusables[fm.currentIndex]
	}
	return nil
}

// GetFocusedID returns the ID of the currently focused element
func (fm *FocusManager) GetFocusedID() string {
	if element := fm.GetFocusedElement(); element != nil {
		return element.GetID()
	}
	return ""
}

// Navigate moves focus in the specified direction
func (fm *FocusManager) Navigate(direction string) bool {
	if len(fm.focusables) == 0 {
		return false
	}

	switch direction {
	case "up":
		return fm.navigateUp()
	case "down":
		return fm.navigateDown()
	case "left":
		return fm.navigateLeft()
	case "right":
		return fm.navigateRight()
	case "next":
		return fm.navigateNext()
	case "previous":
		return fm.navigatePrevious()
	default:
		return false
	}
}

// navigateUp moves focus to the element above the current one
func (fm *FocusManager) navigateUp() bool {
	if fm.currentIndex < 0 {
		return fm.navigateNext()
	}

	current := fm.focusables[fm.currentIndex]
	currentRow, currentCol := current.GetPosition()

	// Find the closest element above
	bestIndex := -1
	smallestDistance := -1

	for i, element := range fm.focusables {
		if i != fm.currentIndex && element.IsEnabled() {
			elemRow, elemCol := element.GetPosition()
			if elemRow < currentRow {
				// Calculate distance (Manhattan distance)
				distance := (currentRow - elemRow) + abs(currentCol - elemCol)
				if bestIndex == -1 || distance < smallestDistance {
					bestIndex = i
					smallestDistance = distance
				}
			}
		}
	}

	if bestIndex != -1 {
		fm.setFocusByIndex(bestIndex)
		return true
	}

	// Wrap around if enabled
	if fm.wrapNavigation {
		return fm.navigateLast()
	}

	return false
}

// navigateDown moves focus to the element below the current one
func (fm *FocusManager) navigateDown() bool {
	if fm.currentIndex < 0 {
		return fm.navigateNext()
	}

	current := fm.focusables[fm.currentIndex]
	currentRow, currentCol := current.GetPosition()

	// Find the closest element below
	bestIndex := -1
	smallestDistance := -1

	for i, element := range fm.focusables {
		if i != fm.currentIndex && element.IsEnabled() {
			elemRow, elemCol := element.GetPosition()
			if elemRow > currentRow {
				// Calculate distance (Manhattan distance)
				distance := (elemRow - currentRow) + abs(currentCol - elemCol)
				if bestIndex == -1 || distance < smallestDistance {
					bestIndex = i
					smallestDistance = distance
				}
			}
		}
	}

	if bestIndex != -1 {
		fm.setFocusByIndex(bestIndex)
		return true
	}

	// Wrap around if enabled
	if fm.wrapNavigation {
		return fm.navigateFirst()
	}

	return false
}

// navigateLeft moves focus to the element to the left of the current one
func (fm *FocusManager) navigateLeft() bool {
	if fm.currentIndex < 0 {
		return fm.navigateNext()
	}

	current := fm.focusables[fm.currentIndex]
	currentRow, currentCol := current.GetPosition()

	// Find the closest element to the left in the same row
	bestIndex := -1
	closestCol := -1

	for i, element := range fm.focusables {
		if i != fm.currentIndex && element.IsEnabled() {
			elemRow, elemCol := element.GetPosition()
			if elemRow == currentRow && elemCol < currentCol {
				if bestIndex == -1 || elemCol > closestCol {
					bestIndex = i
					closestCol = elemCol
				}
			}
		}
	}

	if bestIndex != -1 {
		fm.setFocusByIndex(bestIndex)
		return true
	}

	// Try to find the rightmost element in the row above
	if fm.wrapNavigation {
		aboveRow := currentRow - 1
		rightmostIndex := -1
		rightmostCol := -1

		for i, element := range fm.focusables {
			if i != fm.currentIndex && element.IsEnabled() {
				elemRow, elemCol := element.GetPosition()
				if elemRow == aboveRow {
					if rightmostIndex == -1 || elemCol > rightmostCol {
						rightmostIndex = i
						rightmostCol = elemCol
					}
				}
			}
		}

		if rightmostIndex != -1 {
			fm.setFocusByIndex(rightmostIndex)
			return true
		}
	}

	return false
}

// navigateRight moves focus to the element to the right of the current one
func (fm *FocusManager) navigateRight() bool {
	if fm.currentIndex < 0 {
		return fm.navigateNext()
	}

	current := fm.focusables[fm.currentIndex]
	currentRow, currentCol := current.GetPosition()

	// Find the closest element to the right in the same row
	bestIndex := -1
	closestCol := -1

	for i, element := range fm.focusables {
		if i != fm.currentIndex && element.IsEnabled() {
			elemRow, elemCol := element.GetPosition()
			if elemRow == currentRow && elemCol > currentCol {
				if bestIndex == -1 || elemCol < closestCol {
					bestIndex = i
					closestCol = elemCol
				}
			}
		}
	}

	if bestIndex != -1 {
		fm.setFocusByIndex(bestIndex)
		return true
	}

	// Try to find the leftmost element in the row below
	if fm.wrapNavigation {
		belowRow := currentRow + 1
		leftmostIndex := -1
		leftmostCol := -1

		for i, element := range fm.focusables {
			if i != fm.currentIndex && element.IsEnabled() {
				elemRow, elemCol := element.GetPosition()
				if elemRow == belowRow {
					if leftmostIndex == -1 || elemCol < leftmostCol {
						leftmostIndex = i
						leftmostCol = elemCol
					}
				}
			}
		}

		if leftmostIndex != -1 {
			fm.setFocusByIndex(leftmostIndex)
			return true
		}
	}

	return false
}

// navigateNext moves focus to the next element in the list (Tab navigation)
func (fm *FocusManager) navigateNext() bool {
	if len(fm.focusables) == 0 {
		return false
	}

	if fm.currentIndex < 0 {
		// Start with first enabled element
		for i, element := range fm.focusables {
			if element.IsEnabled() {
				fm.setFocusByIndex(i)
				return true
			}
		}
		return false
	}

	// Find next enabled element
	startIndex := fm.currentIndex
	for i := 1; i <= len(fm.focusables); i++ {
		nextIndex := (startIndex + i) % len(fm.focusables)
		if fm.focusables[nextIndex].IsEnabled() {
			fm.setFocusByIndex(nextIndex)
			return true
		}
	}

	return false
}

// navigatePrevious moves focus to the previous element in the list (Shift+Tab navigation)
func (fm *FocusManager) navigatePrevious() bool {
	if len(fm.focusables) == 0 {
		return false
	}

	if fm.currentIndex < 0 {
		// Start with last enabled element
		for i := len(fm.focusables) - 1; i >= 0; i-- {
			if fm.focusables[i].IsEnabled() {
				fm.setFocusByIndex(i)
				return true
			}
		}
		return false
	}

	// Find previous enabled element
	startIndex := fm.currentIndex
	for i := 1; i <= len(fm.focusables); i++ {
		prevIndex := (startIndex - i + len(fm.focusables)) % len(fm.focusables)
		if fm.focusables[prevIndex].IsEnabled() {
			fm.setFocusByIndex(prevIndex)
			return true
		}
	}

	return false
}

// navigateFirst moves focus to the first enabled element
func (fm *FocusManager) navigateFirst() bool {
	for i, element := range fm.focusables {
		if element.IsEnabled() {
			fm.setFocusByIndex(i)
			return true
		}
	}
	return false
}

// navigateLast moves focus to the last enabled element
func (fm *FocusManager) navigateLast() bool {
	for i := len(fm.focusables) - 1; i >= 0; i-- {
		if fm.focusables[i].IsEnabled() {
			fm.setFocusByIndex(i)
			return true
		}
	}
	return false
}

// Activate activates the currently focused element
func (fm *FocusManager) Activate(model ui.Model) (ui.Model, error) {
	if element := fm.GetFocusedElement(); element != nil {
		return element.Activate(model)
	}
	return model, nil
}

// RestoreFocus restores focus to the previously focused element
func (fm *FocusManager) RestoreFocus() bool {
	if len(fm.focusHistory) == 0 {
		return false
	}

	// Get the last focused ID from history
	lastID := fm.focusHistory[len(fm.focusHistory)-1]
	fm.focusHistory = fm.focusHistory[:len(fm.focusHistory)-1]

	return fm.SetFocus(lastID)
}

// ClearFocus removes focus from all elements
func (fm *FocusManager) ClearFocus() {
	if fm.currentIndex >= 0 && fm.currentIndex < len(fm.focusables) {
		fm.focusables[fm.currentIndex].OnBlur()
		fm.currentIndex = -1
	}
}

// SetWrapNavigation enables or disables wrap-around navigation
func (fm *FocusManager) SetWrapNavigation(wrap bool) {
	fm.wrapNavigation = wrap
}

// GetFocusables returns all focusable elements
func (fm *FocusManager) GetFocusables() []Focusable {
	return fm.focusables
}

// abs returns the absolute value of an integer
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}