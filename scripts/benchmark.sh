#!/bin/bash

set -e

echo "Running CCPM benchmarks..."
echo ""

# Run benchmarks
echo "Running Go benchmarks..."
go test -bench=. -benchmem -benchtime=1s ./... | tee benchmark-results.txt

echo ""
echo "Benchmark results saved to benchmark-results.txt"
echo ""

# Run memory profiling
echo "Running memory profiling..."
mkdir -p profiles

go test -bench=. -memprofile=profiles/mem.prof ./...
go tool pprof -png profiles/mem.prof > profiles/memory-profile.png

echo "Memory profile saved to profiles/memory-profile.png"
echo ""

# Run CPU profiling
echo "Running CPU profiling..."
go test -bench=. -cpuprofile=profiles/cpu.prof ./...
go tool pprof -png profiles/cpu.prof > profiles/cpu-profile.png

echo "CPU profile saved to profiles/cpu-profile.png"
echo ""

# Generate benchmark report
echo "Generating benchmark report..."
cat > benchmark-report.md << EOF
# CCPM Benchmark Report

Generated: $(date)
Commit: $(git rev-parse --short HEAD 2>/dev/null || echo "unknown")

## Benchmark Results

\`\`\`
$(cat benchmark-results.txt)
\`\`\`

## Performance Analysis

This report contains performance benchmarks for the CCPM application.

### Key Metrics
- Binary size validation (< 5MB per platform)
- Memory usage profiling
- CPU performance testing
- Calculator engine performance

### Files Generated
- \`benchmark-results.txt\`: Raw benchmark output
- \`profiles/memory-profile.png\`: Memory usage visualization
- \`profiles/cpu-profile.png\`: CPU usage visualization

EOF

echo "Benchmark report generated: benchmark-report.md"