package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"ccpm-demo/internal/calculator"
	"ccpm-demo/internal/testing/visual"
	"ccpm-demo/internal/ui"
)

func main() {
	var (
		updateMode   = flag.Bool("update", false, "Update baseline screenshots")
		outputDir    = flag.String("output", "test-results", "Output directory for test results")
		verbose      = flag.Bool("verbose", false, "Verbose output")
		tolerance    = flag.Float64("tolerance", 0.01, "Tolerance for visual differences (0.0-1.0)")
		demoMode     = flag.Bool("demo", false, "Generate demo screenshots instead of running tests")
		benchmark    = flag.Bool("benchmark", false, "Run benchmark tests")
		parallel     = flag.Int("parallel", 1, "Number of parallel test runs")
		theme        = flag.String("theme", "retro-casio", "Theme to test")
	)
	flag.Parse()

	if *verbose {
		log.SetFlags(log.LstdFlags | log.Lshortfile)
	} else {
		log.SetOutput(os.Stderr)
	}

	// Create calculator engine
	engine := calculator.NewEngine()

	// Create UI model
	model := ui.NewModel(engine)

	// Set theme if specified
	if *theme != "retro-casio" {
		if err := model.SetButtonGridTheme(*theme); err != nil {
			log.Printf("Warning: Failed to set theme '%s': %v", *theme, err)
		}
	}

	startTime := time.Now()

	if *demoMode {
		if err := runDemoMode(model, *outputDir, *verbose); err != nil {
			log.Fatalf("Demo mode failed: %v", err)
		}
	} else if *benchmark {
		if err := runBenchmarkMode(model, *outputDir, *verbose); err != nil {
			log.Fatalf("Benchmark mode failed: %v", err)
		}
	} else {
		if err := runTestMode(model, *outputDir, *updateMode, *tolerance, *verbose, *parallel); err != nil {
			log.Fatalf("Test mode failed: %v", err)
		}
	}

	duration := time.Since(startTime)
	fmt.Printf("\nTotal execution time: %s\n", duration)
}

func runTestMode(model ui.Model, outputDir string, updateMode bool, tolerance float64, verbose bool, parallel int) error {
	fmt.Printf("Running visual regression tests...\n")
	fmt.Printf("Output directory: %s\n", outputDir)
	fmt.Printf("Update mode: %v\n", updateMode)
	fmt.Printf("Tolerance: %.2f%%\n", tolerance*100)
	fmt.Printf("Parallel runs: %d\n", parallel)

	// Create test configuration
	config := visual.TestConfig{
		BaselineDir:   filepath.Join(outputDir, "baseline"),
		CurrentDir:    filepath.Join(outputDir, "current"),
		DiffDir:       filepath.Join(outputDir, "diff"),
		Tolerance:     tolerance,
		UpdateMode:    updateMode,
		ParallelRuns:  parallel,
		MaxDiffRatio:  0.1,
		MaxTestTime:   30 * time.Second,
		SaveScreenshots: true,
	}

	// Create visual regression test
	test := visual.NewVisualRegressionTest(
		"Calculator Visual Regression",
		"Comprehensive visual regression test for CCPM Calculator",
		model,
		config,
	)

	// Run tests
	if err := test.Run(); err != nil {
		return fmt.Errorf("test execution failed: %w", err)
	}

	// Generate and display report
	report := test.GenerateReport()
	fmt.Println(report)

	// Save results
	resultsFile := filepath.Join(outputDir, "results.json")
	if err := test.SaveResults(resultsFile); err != nil {
		return fmt.Errorf("failed to save results: %w", err)
	}

	// Save report
	reportFile := filepath.Join(outputDir, "report.txt")
	if err := os.WriteFile(reportFile, []byte(report), 0644); err != nil {
		return fmt.Errorf("failed to save report: %w", err)
	}

	fmt.Printf("Results saved to: %s\n", resultsFile)
	fmt.Printf("Report saved to: %s\n", reportFile)

	// Return error if tests failed
	if !test.Results.Passed {
		return fmt.Errorf("visual regression tests failed")
	}

	return nil
}

func runDemoMode(model ui.Model, outputDir string, verbose bool) error {
	fmt.Printf("Generating demo screenshots...\n")
	fmt.Printf("Output directory: %s\n", outputDir)

	// Create demo generator
	config := visual.NewDefaultConfig()
	demoGen := visual.NewDemoGenerator(model, config, outputDir)

	// Generate all demos
	if err := demoGen.GenerateAllDemos(); err != nil {
		return fmt.Errorf("demo generation failed: %w", err)
	}

	fmt.Printf("Demo screenshots generated successfully!\n")
	fmt.Printf("Output directory: %s\n", outputDir)

	// List generated demos
	demos, err := os.ReadDir(outputDir)
	if err != nil {
		return fmt.Errorf("failed to list demos: %w", err)
	}

	fmt.Printf("\nGenerated demos:\n")
	for _, demo := range demos {
		if demo.IsDir() {
			fmt.Printf("  - %s\n", demo.Name())
		}
	}

	return nil
}

func runBenchmarkMode(model ui.Model, outputDir string, verbose bool) error {
	fmt.Printf("Running visual performance benchmarks...\n")
	fmt.Printf("Output directory: %s\n", outputDir)

	// Create benchmark configuration
	config := visual.TestConfig{
		BaselineDir:   filepath.Join(outputDir, "benchmark"),
		CurrentDir:    filepath.Join(outputDir, "benchmark"),
		DiffDir:       filepath.Join(outputDir, "benchmark"),
		Tolerance:     0.0,
		UpdateMode:    true,
		ParallelRuns:  1,
		MaxDiffRatio:  0.0,
		MaxTestTime:   60 * time.Second,
		SaveScreenshots: true,
	}

	// Create visual regression test
	test := visual.NewVisualRegressionTest(
		"Calculator Visual Benchmark",
		"Performance benchmark for visual operations",
		model,
		config,
	)

	// Run benchmark iterations
	iterations := 10
	var totalTime time.Duration
	var minTime, maxTime time.Duration

	for i := 0; i < iterations; i++ {
		fmt.Printf("Running iteration %d/%d...\n", i+1, iterations)

		startTime := time.Now()
		if err := test.Run(); err != nil {
			return fmt.Errorf("benchmark iteration %d failed: %w", i+1, err)
		}

		iterationTime := time.Since(startTime)
		totalTime += iterationTime

		if i == 0 {
			minTime = iterationTime
			maxTime = iterationTime
		} else {
			if iterationTime < minTime {
				minTime = iterationTime
			}
			if iterationTime > maxTime {
				maxTime = iterationTime
			}
		}
	}

	avgTime := totalTime / time.Duration(iterations)

	fmt.Printf("\nBenchmark Results:\n")
	fmt.Printf("Iterations: %d\n", iterations)
	fmt.Printf("Total time: %s\n", totalTime)
	fmt.Printf("Average time: %s\n", avgTime)
	fmt.Printf("Min time: %s\n", minTime)
	fmt.Printf("Max time: %s\n", maxTime)

	// Save benchmark results
	benchmarkResults := fmt.Sprintf(`{
  "name": "Calculator Visual Benchmark",
  "iterations": %d,
  "total_time_ms": %d,
  "average_time_ms": %d,
  "min_time_ms": %d,
  "max_time_ms": %d,
  "run_at": "%s"
}`, iterations, totalTime.Milliseconds(), avgTime.Milliseconds(),
   minTime.Milliseconds(), maxTime.Milliseconds(), time.Now().Format(time.RFC3339))

	benchmarkFile := filepath.Join(outputDir, "benchmark.json")
	if err := os.WriteFile(benchmarkFile, []byte(benchmarkResults), 0644); err != nil {
		return fmt.Errorf("failed to save benchmark results: %w", err)
	}

	fmt.Printf("Benchmark results saved to: %s\n", benchmarkFile)

	return nil
}

func init() {
	// Set up logging
	log.SetPrefix("[visual-test] ")
	log.SetFlags(log.LstdFlags)
}