.PHONY: build test benchmark clean release security-check lint help

# Version information
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_TIME := $(shell date -u '+%Y-%m-%d %H:%M:%S UTC')
COMMIT_HASH := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")

# Build flags
LDFLAGS := -s -w -X 'main.Version=$(VERSION)' -X 'main.BuildTime=$(BUILD_TIME)' -X 'main.CommitHash=$(COMMIT_HASH)'

# Default target
all: build

# Build the application for current platform
build:
	go build -ldflags="$(LDFLAGS)" -o ccpm .

# Build for all platforms
build-all:
	./scripts/build.sh

# Run tests
test:
	./scripts/test.sh

# Run benchmarks
benchmark:
	./scripts/benchmark.sh

# Run security checks
security-check:
	govulncheck ./...

# Run linting
lint:
	golangci-lint run

# Clean build artifacts
clean:
	rm -rf build/ coverage.html coverage.out coverage.txt benchmark-results.txt profiles/

# Create release
release: build-all
	@echo "Release created in build/ directory"

# Run the application
run: build
	./ccpm

# Install dependencies
deps:
	go mod download
	go mod tidy

# Generate documentation
docs:
	@echo "Documentation files:"
	@find . -name "*.md" -not -path "./.git/*" | sort

# Show help
help:
	@echo "CCPM Makefile"
	@echo ""
	@echo "Targets:"
	@echo "  build        - Build for current platform"
	@echo "  build-all    - Build for all platforms"
	@echo "  test         - Run tests"
	@echo "  benchmark    - Run benchmarks"
	@echo "  security-check - Run security checks"
	@echo "  lint         - Run linting"
	@echo "  clean        - Clean build artifacts"
	@echo "  release      - Create release"
	@echo "  run          - Run the application"
	@echo "  deps         - Install dependencies"
	@echo "  docs         - Show documentation files"
	@echo "  help         - Show this help"