---
name: tuic
description: Terminal-based calculator with retro Casio styling, cross-platform Go implementation
status: backlog
created: 2025-09-27T19:24:49Z
---

# PRD: TUIC (Terminal UI Calculator)

## Executive Summary

TUIC is a terminal-based calculator application that recreates the classic retro Casio calculator experience in a modern cross-platform implementation. Built with Go for rapid development and easy deployment, TUIC provides basic mathematical operations through an intuitive grid-based interface that supports both mouse and keyboard navigation with visual and audio feedback.

**Value Proposition**: Developers and terminal users get a familiar, efficient calculator without leaving their command-line environment, combining nostalgic design with modern terminal capabilities.

## Problem Statement

### What problem are we solving?
- Terminal users currently lack a visually appealing, interactive calculator that works seamlessly in their command-line environment
- Existing terminal calculators are either too complex (scientific) or too simplistic (command-line only)
- No current solution provides the familiar tactile experience of physical calculators in terminal environments

### Why is this important now?
- Growing terminal/CLI-first developer community needs better productivity tools
- Cross-platform compatibility is essential for modern development workflows
- Single-binary deployment aligns with modern DevOps practices

## User Stories

### Primary User Personas

**Terminal Power User (Alex)**
- Works primarily in terminal environments
- Values efficiency and keyboard shortcuts
- Wants calculator access without GUI context switching
- Appreciates retro/nostalgic interfaces

**Developer (Sam)**
- Frequently needs quick calculations during coding
- Uses multiple operating systems (macOS, Linux, Windows)
- Prefers lightweight, fast-loading tools
- Values applications that integrate well with development workflow

### Detailed User Journeys

**Story 1: Quick Calculation**
```
As Alex, I want to quickly calculate 25 * 4 + 100
So that I can verify a calculation without leaving my terminal

Acceptance Criteria:
- Can launch calculator with `tuic` command
- Can input "25 * 4 + 100 =" using keyboard
- See both the expression being built and final result
- Calculator responds within 100ms to each input
```

**Story 2: Mouse Interaction**
```
As Sam, I want to use my mouse to click calculator buttons
So that I can interact naturally when my hands are already on the mouse

Acceptance Criteria:
- Mouse clicks register on calculator buttons
- Buttons show visual feedback when clicked
- Mouse hover shows button focus state
- Can complete full calculations using only mouse
```

**Story 3: Keyboard Navigation**
```
As Alex, I want to navigate the calculator using Tab and arrow keys
So that I can operate it efficiently without a mouse

Acceptance Criteria:
- Tab/Shift+Tab cycles through all buttons
- Arrow keys navigate the button grid
- Enter/Space activates the focused button
- Number keys directly input numbers regardless of focus
- Current focus is clearly visible
```

### Pain Points Being Addressed
- Context switching between terminal and GUI calculators
- Lack of tactile feedback in terminal calculations
- Inconsistent calculator availability across platforms
- Complex installation processes for simple tools

## Requirements

### Functional Requirements

**Core Mathematical Operations**
- Addition (+), Subtraction (-), Multiplication (*), Division (/)
- Decimal number support with proper floating-point handling
- Clear (C) and Clear Entry (CE) functions
- Equals (=) operation for result calculation

**User Interface**
- Traditional calculator grid layout (4x5 buttons recommended)
- Display area showing:
  - Current expression being built
  - Final calculation result
  - Error messages in human-readable format
- Retro Casio-style visual design with terminal colors

**Input Methods**
- **Keyboard Input**:
  - Number keys (0-9) for direct number entry
  - Operator keys (+, -, *, /) for operations
  - Enter key for equals operation
  - Tab/Shift+Tab for button navigation
  - Arrow keys for grid navigation
  - Space for button activation
- **Mouse Input**:
  - Click buttons to activate
  - Hover for focus indication

**Visual Feedback**
- Button color changes when pressed/clicked
- Different visual styles for:
  - Number buttons
  - Operator buttons
  - Special buttons (equals, clear)
- Focus indicator for keyboard navigation

**Audio Feedback**
- Button press sounds
- Error notification sounds
- Distinct sounds for different button types

### Non-Functional Requirements

**Performance**
- Application startup time: < 500ms
- Input response time: < 100ms
- Memory usage: < 10MB at runtime
- Single binary deployment

**Compatibility**
- **Operating Systems**: macOS, Linux, Windows
- **Terminal Types**: Support all major terminals (Terminal.app, iTerm2, GNOME Terminal, Windows Terminal, etc.)
- **Terminal Features**: Colors, cursor positioning, mouse support

**Reliability**
- Graceful error handling for all edge cases
- No crashes on invalid input
- Proper cleanup on exit

**Usability**
- Intuitive operation matching physical calculator behavior
- Clear visual hierarchy and button grouping
- Consistent behavior across all platforms

**Testing Requirements**
- **End-to-End Testing**: Complete user workflow testing from application launch to calculation completion
- **Visual Testing**: Automated visual regression testing to ensure UI consistency across releases
- **Visual Demo Creation**: Each PR must include visual demonstration (screenshots/recordings) showing:
  - New features being added
  - Before/after comparisons for UI changes
  - Cross-platform behavior demonstration
  - Interactive feature showcases (mouse/keyboard navigation)
- **Cross-Platform Testing**: Automated testing on all target platforms (macOS, Linux, Windows)
- **Terminal Compatibility Testing**: Testing across multiple terminal emulators

## Success Criteria

### Measurable Outcomes
- **Adoption**: 1000+ downloads within first 3 months
- **Performance**: 95% of operations complete within 100ms
- **Compatibility**: Successfully runs on 3 major OS platforms
- **User Satisfaction**: 4.5+ star rating on release platforms
- **Visual Consistency**: Zero visual regression failures in CI/CD pipeline

### Key Metrics and KPIs
- Binary size: < 5MB
- Cross-platform test coverage: 100%
- User retention: 70% weekly active users
- Issue resolution time: < 48 hours for critical bugs
- Visual test coverage: 100% of UI components

## Constraints & Assumptions

### Technical Limitations
- Terminal-only interface (no GUI dependencies)
- Limited to basic arithmetic operations (no scientific functions)
- Single-user application (no multi-user features)
- Text-based display (no graphics beyond terminal characters)

### Timeline Constraints
- First release target: 4-6 weeks from start
- MVP must include all core mathematical operations
- Cross-platform testing required before release

### Resource Limitations
- Single developer for initial implementation
- Must use free/open-source libraries only
- CI/CD through GitHub Actions (free tier limits)

### Assumptions
- Users have basic terminal knowledge
- Go toolchain available for development
- GitHub suitable for project hosting and CI/CD
- Terminal mouse support available on target platforms

## Out of Scope

### Explicitly NOT Building (First Release)
- Scientific calculator functions (sin, cos, log, etc.)
- Memory functions (M+, M-, MR, MC)
- Calculation history/previous results
- Settings/preferences persistence
- Multiple calculator modes
- Programmable functions
- Unit conversions
- Graph plotting capabilities
- Multi-precision arithmetic
- Plugin system

### Future Considerations (Not First Release)
- History functionality
- Memory operations
- Preferences saving
- Multiple themes
- Calculation export features

## Dependencies

### External Dependencies
- **Go Programming Language**: Version 1.19+ for development
- **Terminal Libraries**:
  - TUI library (tcell, bubbletea, or similar)
  - Mouse support library
  - Audio library for sound feedback
- **GitHub**: Repository hosting and CI/CD
- **GitHub Actions**: Automated testing and release builds
- **Visual Testing Tools**: Screenshot/recording tools for demo creation

### Internal Team Dependencies
- Single developer responsible for:
  - Go implementation
  - Cross-platform testing
  - Visual demo creation
  - Documentation
  - Release management

### Platform Dependencies
- Terminal emulator with:
  - Color support
  - Mouse input capability
  - Cursor positioning
  - Audio output capability

## Technical Architecture

### Language Choice: Go
**Rationale**:
- Single binary compilation
- Excellent cross-platform support
- Strong standard library
- Fast compilation for rapid development
- Growing ecosystem of terminal UI libraries
- Easy to learn and maintain

### Recommended Libraries
- **TUI Framework**: Bubble Tea (modern, well-maintained)
- **Terminal Rendering**: Lip Gloss (styling companion to Bubble Tea)
- **Audio**: Beep library for cross-platform sound
- **Testing**: Built-in Go testing framework
- **Visual Testing**: Custom screenshot/recording utilities

### Build & Distribution
- Single binary per platform (Windows .exe, macOS/Linux executables)
- GitHub Releases for distribution
- GitHub Actions for automated building
- Cross-compilation using Go's built-in capabilities
- Automated visual demo generation in CI/CD

## Implementation Phases

### Phase 1: Core Functionality (Weeks 1-2)
- Basic calculator engine
- Simple terminal UI
- Keyboard number input
- Basic operations (+, -, *, /, =)
- Initial end-to-end test framework

### Phase 2: Enhanced UI (Weeks 3-4)
- Grid-based button layout
- Visual feedback and styling
- Mouse support
- Keyboard navigation
- Visual testing infrastructure
- Visual demo creation tools

### Phase 3: Polish & Distribution (Weeks 5-6)
- Audio feedback
- Error handling
- Cross-platform testing
- CI/CD setup with visual testing
- Comprehensive visual demo library
- First release with complete visual documentation

## Testing Strategy

### End-to-End Testing Requirements
- **Full User Workflows**: Test complete calculation sequences from launch to result
- **Cross-Platform E2E**: Same test suite running on all target platforms
- **Input Method Testing**: E2E tests for both keyboard and mouse interaction paths
- **Error Scenario Testing**: E2E tests for all error conditions and recovery

### Visual Testing & Demo Requirements
- **Automated Visual Regression**: Screenshots compared across builds to detect UI changes
- **PR Visual Requirements**: Every PR must include:
  - Before/after screenshots for UI changes
  - Screen recordings demonstrating new features
  - Cross-platform visual comparison when applicable
  - Interactive demo showing mouse/keyboard navigation
- **Demo Library**: Build comprehensive visual library showing:
  - All calculator operations
  - Navigation methods
  - Error states
  - Platform-specific behaviors
- **CI/CD Integration**: Visual tests and demo generation automated in build pipeline

## Risk Mitigation

### Technical Risks
- **Terminal compatibility issues**: Early testing across multiple terminals
- **Audio support limitations**: Graceful degradation if audio unavailable
- **Mouse support inconsistency**: Ensure keyboard-only operation works perfectly
- **Visual testing complexity**: Start with simple screenshot comparison, evolve as needed

### Project Risks
- **Scope creep**: Strict adherence to "basic operations only" for first release
- **Cross-platform testing**: Automated CI/CD testing on multiple platforms
- **Performance issues**: Regular benchmarking during development
- **Visual demo maintenance**: Automated generation to reduce manual overhead