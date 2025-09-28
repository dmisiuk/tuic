---
name: tuic
status: in-progress
created: 2025-09-27T19:40:11Z
progress: 77%
updated: 2025-09-28T09:32:10Z
prd: ccpm/prds/tuic.md
github: https://github.com/dmisiuk/tuic/issues/4
---

# Epic: TUIC (Terminal UI Calculator)

## Overview

TUIC is a Go-based terminal calculator that recreates the classic Casio calculator experience using modern TUI libraries. The implementation focuses on simplicity and performance, leveraging Bubble Tea for UI rendering and cross-platform compatibility through Go's single-binary deployment model.

**Technical Summary**: A self-contained terminal application with grid-based button layout, dual input modes (keyboard/mouse), visual/audio feedback, and comprehensive testing including visual regression testing.

## Architecture Decisions

### Core Technology Stack
- **Language**: Go 1.19+ for cross-platform single-binary compilation
- **TUI Framework**: Bubble Tea + Lip Gloss for modern terminal UI with styling
- **Audio**: Beep library for cross-platform sound feedback
- **Testing**: Go's built-in testing framework + custom visual testing utilities
- **Build**: GitHub Actions with cross-compilation for Windows, macOS, Linux

### Key Design Patterns
- **Model-View-Update (MVU)**: Bubble Tea's architecture for predictable state management
- **Component-Based UI**: Separate components for display, button grid, and input handling
- **Command Pattern**: Button actions as discrete commands for easy testing
- **Strategy Pattern**: Different input handlers for keyboard vs mouse interaction

### Performance Optimizations
- **Minimal Dependencies**: Only essential libraries to keep binary size < 5MB
- **Efficient Rendering**: Selective screen updates using Bubble Tea's optimization
- **Memory Management**: Stateless calculations to minimize memory footprint

## Technical Approach

### Frontend Components

**Display Component**
- Current expression state management
- Result formatting and display
- Error message handling with graceful degradation
- Retro styling using terminal colors and box drawing characters

**Button Grid Component**
- 4x5 grid layout (numbers, operators, special functions)
- Visual state management (normal, focused, pressed)
- Different button types: numbers (0-9), operators (+,-,*,/), special (=,C,CE)
- Keyboard navigation with focus tracking

**Input Handler Component**
- Dual input processing: direct key input + navigation
- Mouse event handling with hover states
- Tab/Shift+Tab cycling through button grid
- Arrow key navigation with wraparound

### Backend Services

**Calculator Engine**
- Basic arithmetic operations with proper operator precedence
- Floating-point precision handling
- Division by zero and overflow error handling
- Expression parsing and evaluation

**Audio Service**
- Cross-platform sound feedback using Beep library
- Different sound profiles for number, operator, and special buttons
- Graceful degradation when audio unavailable

**State Management**
- Centralized state using Bubble Tea model
- Immutable state updates for predictability
- Clear separation between UI state and calculation state

### Infrastructure

**Build System**
- Go modules for dependency management
- Cross-compilation for 3 platforms via GitHub Actions
- Single binary output with no external dependencies
- Automated release pipeline with artifact generation

**Testing Infrastructure**
- Unit tests for calculator engine
- Integration tests for full user workflows
- Visual regression testing using custom screenshot utilities
- Cross-platform testing matrix in CI/CD

## Implementation Strategy

### Development Phases

**Phase 1: Core Engine (Week 1)**
- Calculator logic implementation
- Basic Bubble Tea app structure
- Simple text-based interface
- Unit tests for calculation engine

**Phase 2: UI Implementation (Weeks 2-3)**
- Grid layout with Bubble Tea components
- Visual styling with Lip Gloss
- Keyboard and mouse input handling
- Focus management and navigation

**Phase 3: Polish & Testing (Weeks 4-5)**
- Audio feedback integration
- Visual regression testing setup
- Cross-platform testing
- Performance optimization

**Phase 4: CI/CD & Release (Week 6)**
- GitHub Actions pipeline
- Automated visual demo generation
- Release automation
- Documentation and demo creation

### Risk Mitigation

**Terminal Compatibility**: Early testing across major terminals, fallback rendering for unsupported features
**Audio Dependencies**: Optional audio with graceful degradation, no hard dependencies
**Cross-Platform Issues**: Automated testing matrix, early platform-specific testing
**Performance**: Regular benchmarking, profiling integration in CI/CD

## Task Breakdown Preview

High-level task categories (≤10 tasks total):

- [ ] **Calculator Engine**: Core arithmetic logic, error handling, expression parsing
- [ ] **TUI Foundation**: Bubble Tea app structure, basic rendering, state management
- [ ] **Button Grid UI**: 4x5 layout, visual styling, button components with Lip Gloss
- [ ] **Input System**: Keyboard/mouse handlers, focus management, navigation logic
- [ ] **Audio Integration**: Beep library integration, sound profiles, graceful degradation
- [ ] **Visual Testing**: Screenshot utilities, regression testing, demo generation tools
- [ ] **Cross-Platform Build**: GitHub Actions, multi-platform compilation, release automation
- [ ] **End-to-End Testing**: Full workflow tests, error scenario coverage, performance validation
- [ ] **Documentation & Demos**: Visual demos, usage documentation, PR demo requirements

## Dependencies

### External Dependencies
- **Go 1.19+**: Required for modern module support and cross-compilation
- **Bubble Tea**: TUI framework for terminal applications
- **Lip Gloss**: Styling library companion to Bubble Tea
- **Beep**: Cross-platform audio library for button feedback
- **GitHub Actions**: CI/CD pipeline for automated building and testing

### Internal Dependencies
- **Visual Testing Tools**: Custom utilities for screenshot capture and comparison
- **Demo Generation**: Automated visual demonstration creation for PRs
- **Cross-Platform Testing**: Validation across Windows, macOS, Linux environments

### Platform Dependencies
- **Terminal Emulator**: Modern terminal with color, mouse, and cursor support
- **Audio Output**: System audio capability (optional, graceful degradation)
- **Go Toolchain**: Available for local development and testing

## Success Criteria (Technical)

### Performance Benchmarks
- **Startup Time**: < 500ms from command execution to UI ready
- **Input Response**: < 100ms from input to visual feedback
- **Memory Usage**: < 10MB runtime memory footprint
- **Binary Size**: < 5MB for each platform binary

### Quality Gates
- **Test Coverage**: 90%+ code coverage for calculator engine
- **Cross-Platform**: 100% feature parity across all 3 platforms
- **Visual Consistency**: Zero visual regressions in automated testing
- **Error Handling**: Graceful handling of all edge cases (division by zero, overflow, invalid input)

### Acceptance Criteria
- **All User Stories**: Complete implementation of all PRD user journeys
- **Performance Targets**: Meeting all specified performance benchmarks
- **Visual Requirements**: Comprehensive visual demo library for all features
- **Platform Support**: Verified functionality on all target platforms

## Estimated Effort

### Overall Timeline
- **Total Duration**: 6 weeks (4-6 week target from PRD)
- **Core Development**: 4 weeks
- **Testing & Polish**: 1.5 weeks
- **CI/CD & Release**: 0.5 weeks

### Resource Requirements
- **Single Developer**: Full-stack development, testing, and release management
- **Platform Access**: Windows, macOS, Linux environments for testing
- **GitHub Actions**: Free tier sufficient for build and release automation

### Critical Path Items
1. **Bubble Tea Integration**: Learning curve and architecture setup
2. **Cross-Platform Audio**: Beep library integration and fallback handling
3. **Visual Testing**: Custom screenshot and comparison utilities
4. **GitHub Actions Setup**: Multi-platform build and release pipeline

**Risk Buffer**: Built-in 1-week buffer for platform-specific issues and testing

## Tasks Created
- [x] #5 - Calculator Engine (parallel: true) ✅
- [x] #8 - TUI Foundation (parallel: true) ✅
- [x] #9 - Button Grid UI (parallel: false) ✅
- [x] #10 - Input System (parallel: false) ✅
- [x] #11 - Audio Integration (parallel: false) ✅
- [x] #12 - Visual Testing (parallel: true) ✅
- [x] #13 - Cross-Platform Build (parallel: true) ✅
- [ ] #6 - End-to-End Testing (parallel: false)
- [ ] #7 - Documentation & Demos (parallel: true)

Total tasks: 9
Parallel tasks: 5 (#5, #8, #12, #13, #7)
Sequential tasks: 4 (#9, #10, #11, #6)
Estimated total effort: 114-146 hours (approximately 6 weeks for single developer)