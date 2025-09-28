# WARP.md

This file provides guidance to WARP (warp.dev) when working with code in this repository.

Project: CCPM Demo (Go 1.25.1)

- Primary language: Go
- Entrypoint: main.go (CLI calculator with interactive and one-shot eval modes)
- Core packages: internal/calculator (engine, parser, errors), internal/ui (Bubble Tea TUI), internal/audio (beep integration)
- Tooling: Makefile + scripts/, golangci-lint, govulncheck

Common commands

- Dependencies
  - make deps
  - go mod tidy

- Build
  - make build
  - Output binary: ./ccpm (ldflags inject Version, BuildTime, CommitHash)
  - Cross-platform: make build-all (scripts/build.sh writes to build/ and generates checksums)

- Run
  - make run
  - ./ccpm --help
  - ./ccpm --version
  - Evaluate once: ./ccpm --eval "(10 + 5) * 3"
  - Interactive REPL: ./ccpm, then type expressions or commands: help | vars | clear | set var = value | quit

- Test
  - All tests: go test ./...
  - Verbose + race + coverage (script): ./scripts/test.sh
  - Single package: go test ./internal/calculator -v
  - Single file: go test ./internal/calculator -run TestEngine
  - Regex match: go test ./internal/calculator -run 'TestParser.*'
  - Benchmarks: make benchmark (writes benchmark-results.txt and profiles/*)
  - Coverage artifacts (from script): coverage.out, coverage.html, coverage.txt

- Lint and security
  - Lint: make lint (uses golangci-lint)
  - Security scan: make security-check (uses govulncheck)
  - The test script will auto-install golangci-lint and govulncheck if missing

- Cleaning
  - make clean (removes build/ coverage.* benchmark-results.txt profiles/)

High-level architecture

- CLI application (main.go)
  - Modes: one-shot evaluation via --eval and interactive REPL
  - Commands: --help | --version, interactive: help | version | vars | clear | set var = value | quit
  - Versioning: Makefile passes ldflags to set Version, BuildTime, CommitHash

- Calculator core (internal/calculator)
  - Engine (engine.go)
    - Stateful arithmetic engine with Clear, ClearEntry, Add/Subtract/Multiply/Divide, InputNumber, PerformOperation
    - Evaluate delegates to Parser and validates numeric ranges
    - Maintains currentValue/entryValue/shouldClear with simple input model
  - Parser (parser.go)
    - Whitespace-stripping recursive-descent parser
    - Grammar: expression (+/-) → term (*//) → factor (number | (expr) | unary +/-)
    - Validates numbers, handles division by zero and mismatched parentheses
  - Errors and validation (errors.go)
    - Error taxonomy: ErrEmptyExpression, ErrInvalidExpression, ErrInvalidNumber, ErrDivisionByZero, ErrInvalidOperator, ErrMismatchedParentheses, ErrNumberOutOfRange

- TUI (internal/ui)
  - Bubble Tea-based MVU Model (model.go, terminal.go)
  - State: displayValue, operator, previousValue, isWaitingForOperand, input/output/history
  - Styling with lipgloss; optional ButtonGrid and audio hooks

- Audio integration (internal/audio)
  - Optional UX enhancement using github.com/faiface/beep and audio plumbing
  - Non-fatal initialization; emits UI events via audio.EventHandler

- Commands and utilities (cmd/)
  - cmd/tuic: TUI client entry; cmd/test-visual: visual test harness

- Tooling and automation
  - Makefile
    - build, build-all (scripts/build.sh), test (scripts/test.sh), benchmark (scripts/benchmark.sh), lint, security-check, run, deps, clean, release
  - scripts/
    - build.sh: multi-OS/arch builds into build/, checksum generation
    - test.sh: unit tests (+race +coverage), coverage outputs, govulncheck, golangci-lint
    - benchmark.sh: benchmarks + mem/cpu profiles and a markdown report

Important references from repository docs

- README.md
  - Confirms Go version, quick build/test instructions, engine usage examples, and directory overview for internal/calculator
  - Testing commands: go test ./..., -cover, -v, and per-file tests

Notes about rules and agents

- An AGENTS.md exists at the repo root that describes specialized agents (code-analyzer, file-analyzer, test-runner, parallel-worker) and how they should summarize work. When coordinating larger refactors or investigations, prefer using these agents’ patterns: perform heavy analysis in isolation and return concise summaries.

UI and debugging tips specific to this codebase

- REPL quick checks: use ./ccpm --eval "expression" to validate parser/engine behavior without entering the interactive loop
- LDFLAGS-driven versioning: verify with ./ccpm --version after make build to ensure Makefile ldflags are applied
- TUI concerns: Bubble Tea rendering and key handling live under internal/ui; when adding UI features, track state in Model and update view() accordingly
