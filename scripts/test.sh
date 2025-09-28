#!/bin/bash

set -e

echo "Running CCPM tests..."
echo ""

# Run unit tests with coverage
echo "Running unit tests with coverage..."
go test -v -race -coverprofile=coverage.out ./...

# Generate coverage report
echo "Generating coverage report..."
go tool cover -html=coverage.out -o coverage.html
go tool cover -func=coverage.out > coverage.txt

echo "Coverage report generated:"
echo "  HTML: coverage.html"
echo "  Text: coverage.txt"
echo ""

# Run security scan
echo "Running security scan..."
if ! command -v govulncheck &> /dev/null; then
    echo "Installing govulncheck..."
    go install golang.org/x/vuln/cmd/govulncheck@latest
fi

govulncheck ./...

echo ""
# Run linting
echo "Running linting..."
if ! command -v golangci-lint &> /dev/null; then
    echo "Installing golangci-lint..."
    go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
fi

golangci-lint run

echo ""
echo "All tests completed successfully!"
echo ""

# Show coverage summary
echo "Coverage summary:"
grep -A 5 "total:" coverage.txt || echo "Coverage summary not available"