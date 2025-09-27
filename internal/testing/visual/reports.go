package visual

import (
	"encoding/json"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// ReportGenerator generates various types of visual test reports
type ReportGenerator struct {
	TestResults    *TestResults
	OutputDir      string
	TemplateDir    string
	IncludeScreenshots bool
	GenerateHTML    bool
	GenerateJSON    bool
	GenerateText    bool
}

// ReportConfig contains configuration for report generation
type ReportConfig struct {
	OutputDir      string
	TemplateDir    string
	IncludeScreenshots bool
	GenerateHTML    bool
	GenerateJSON    bool
	GenerateText    bool
	Theme          string
	ShowDiffImages  bool
	ShowThumbnails  bool
}

// NewReportGenerator creates a new report generator
func NewReportGenerator(results *TestResults, config ReportConfig) *ReportGenerator {
	return &ReportGenerator{
		TestResults:    results,
		OutputDir:      config.OutputDir,
		TemplateDir:    config.TemplateDir,
		IncludeScreenshots: config.IncludeScreenshots,
		GenerateHTML:    config.GenerateHTML,
		GenerateJSON:    config.GenerateJSON,
		GenerateText:    config.GenerateText,
	}
}

// GenerateReports generates all configured reports
func (rg *ReportGenerator) GenerateReports() error {
	if err := os.MkdirAll(rg.OutputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	var errors []string

	if rg.GenerateJSON {
		if err := rg.generateJSONReport(); err != nil {
			errors = append(errors, fmt.Sprintf("JSON report: %v", err))
		}
	}

	if rg.GenerateText {
		if err := rg.generateTextReport(); err != nil {
			errors = append(errors, fmt.Sprintf("Text report: %v", err))
		}
	}

	if rg.GenerateHTML {
		if err := rg.generateHTMLReport(); err != nil {
			errors = append(errors, fmt.Sprintf("HTML report: %v", err))
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("report generation errors: %s", strings.Join(errors, "; "))
	}

	return nil
}

// generateJSONReport generates a JSON report
func (rg *ReportGenerator) generateJSONReport() error {
	report := rg.enrichTestResults()

	data, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	filename := filepath.Join(rg.OutputDir, "visual-test-report.json")
	if err := os.WriteFile(filename, data, 0644); err != nil {
		return fmt.Errorf("failed to write JSON report: %w", err)
	}

	return nil
}

// generateTextReport generates a text report
func (rg *ReportGenerator) generateTextReport() error {
	report := rg.generateTextReportContent()

	filename := filepath.Join(rg.OutputDir, "visual-test-report.txt")
	if err := os.WriteFile(filename, []byte(report), 0644); err != nil {
		return fmt.Errorf("failed to write text report: %w", err)
	}

	return nil
}

// generateHTMLReport generates an HTML report
func (rg *ReportGenerator) generateHTMLReport() error {
	htmlTemplate := rg.getHTMLTemplate()

	// Create template with helper functions
	funcMap := template.FuncMap{
		"base": func(path string) string {
			return filepath.Base(path)
		},
		"gt": func(a, b interface{}) bool {
			switch av := a.(type) {
			case float64:
				switch bv := b.(type) {
				case float64:
					return av > bv
				case int:
					return av > float64(bv)
				}
			case int:
				switch bv := b.(type) {
				case float64:
					return float64(av) > bv
				case int:
					return av > bv
				}
			}
			return false
		},
		"formatFloat": func(f float64) string {
			return fmt.Sprintf("%.4f", f)
		},
		"mul": func(a, b int) int {
			return a * b
		},
	}

	tmpl, err := template.New("visual-report").Funcs(funcMap).Parse(htmlTemplate)
	if err != nil {
		return fmt.Errorf("failed to parse HTML template: %w", err)
	}

	report := rg.enrichTestResults()

	var buf strings.Builder
	if err := tmpl.Execute(&buf, report); err != nil {
		return fmt.Errorf("failed to execute HTML template: %w", err)
	}

	filename := filepath.Join(rg.OutputDir, "visual-test-report.html")
	if err := os.WriteFile(filename, []byte(buf.String()), 0644); err != nil {
		return fmt.Errorf("failed to write HTML report: %w", err)
	}

	// Copy CSS file
	cssContent := rg.getCSSContent()
	cssFilename := filepath.Join(rg.OutputDir, "visual-test-report.css")
	if err := os.WriteFile(cssFilename, []byte(cssContent), 0644); err != nil {
		return fmt.Errorf("failed to write CSS file: %w", err)
	}

	return nil
}

// HTMLReportData contains data for HTML template rendering
type HTMLReportData struct {
	*TestResults
	PassRate         float64
	AverageDuration  time.Duration
	MaxDuration      time.Duration
	Recommendations  []string
}

// enrichTestResults adds additional information to test results
func (rg *ReportGenerator) enrichTestResults() *HTMLReportData {
	baseReport := *rg.TestResults

	// Add summary statistics
	baseReport.TestCases = rg.enrichTestCases()

	return &HTMLReportData{
		TestResults:     &baseReport,
		PassRate:        rg.getPassRate(),
		AverageDuration:  rg.getAverageDuration(),
		MaxDuration:     rg.getMaxDuration(),
		Recommendations: rg.getRecommendations(),
	}
}

// enrichTestCases adds additional information to test cases
func (rg *ReportGenerator) enrichTestCases() map[string]*TestCaseResult {
	enriched := make(map[string]*TestCaseResult)

	for name, result := range rg.TestResults.TestCases {
		enriched[name] = rg.enrichTestCase(result)
	}

	return enriched
}

// enrichTestCase adds additional information to a test case
func (rg *ReportGenerator) enrichTestCase(result *TestCaseResult) *TestCaseResult {
	enriched := *result

	// Add status text
	if enriched.Passed {
		enriched.Details = "✅ PASSED"
	} else if enriched.Skipped {
		enriched.Details = "⏭️ SKIPPED"
	} else {
		enriched.Details = "❌ FAILED"
	}

	// Add performance rating
	if enriched.Duration > 0 {
		enriched.Details += fmt.Sprintf(" (%v)", enriched.Duration.Round(time.Millisecond))
	}

	return &enriched
}

// EnhancedTestCaseResult contains additional fields for HTML reporting
type EnhancedTestCaseResult struct {
	*TestCaseResult
	DiffRatioPercent  float64
	ScreenshotBase    string
	DiffImageBase     string
}

// generateTextReportContent generates the content for a text report
func (rg *ReportGenerator) generateTextReportContent() string {
	var report strings.Builder

	report.WriteString("=== Visual Test Report ===\n\n")
	report.WriteString(fmt.Sprintf("Test: %s\n", rg.TestResults.Name))
	report.WriteString(fmt.Sprintf("Description: %s\n", rg.TestResults.Description))
	report.WriteString(fmt.Sprintf("Run At: %s\n", rg.TestResults.RunAt.Format(time.RFC3339)))
	report.WriteString(fmt.Sprintf("Duration: %s\n", rg.TestResults.Duration))
	report.WriteString(fmt.Sprintf("Environment: %s\n\n", rg.TestResults.Environment))

	// Summary
	report.WriteString("--- Summary ---\n")
	report.WriteString(fmt.Sprintf("Status: %s\n", rg.getStatusString()))
	report.WriteString(fmt.Sprintf("Total Tests: %d\n", rg.TestResults.TotalTests))
	report.WriteString(fmt.Sprintf("Passed: %d\n", rg.TestResults.PassedTests))
	report.WriteString(fmt.Sprintf("Failed: %d\n", rg.TestResults.FailedTests))
	report.WriteString(fmt.Sprintf("Skipped: %d\n", rg.TestResults.SkippedTests))
	report.WriteString(fmt.Sprintf("Pass Rate: %.1f%%\n\n", rg.getPassRate()))

	// Failed Tests
	if rg.TestResults.FailedTests > 0 {
		report.WriteString("--- Failed Tests ---\n")
		for name, result := range rg.TestResults.TestCases {
			if !result.Passed && !result.Skipped {
				report.WriteString(fmt.Sprintf("\n❌ %s\n", name))
				report.WriteString(fmt.Sprintf("   Error: %s\n", result.Error))
				if result.DiffRatio > 0 {
					report.WriteString(fmt.Sprintf("   Diff Ratio: %.2f%%\n", result.DiffRatio*100))
				}
				report.WriteString(fmt.Sprintf("   Duration: %v\n", result.Duration))
				if result.Screenshot != "" {
					report.WriteString(fmt.Sprintf("   Screenshot: %s\n", filepath.Base(result.Screenshot)))
				}
				if result.DiffImage != "" {
					report.WriteString(fmt.Sprintf("   Diff Image: %s\n", filepath.Base(result.DiffImage)))
				}
			}
		}
		report.WriteString("\n")
	}

	// Passed Tests
	if rg.TestResults.PassedTests > 0 {
		report.WriteString("--- Passed Tests ---\n")
		for name, result := range rg.TestResults.TestCases {
			if result.Passed {
				report.WriteString(fmt.Sprintf("✅ %s (%v)\n", name, result.Duration))
			}
		}
		report.WriteString("\n")
	}

	// Skipped Tests
	if rg.TestResults.SkippedTests > 0 {
		report.WriteString("--- Skipped Tests ---\n")
		for name, result := range rg.TestResults.TestCases {
			if result.Skipped {
				report.WriteString(fmt.Sprintf("⏭️ %s (%s)\n", name, result.Error))
			}
		}
		report.WriteString("\n")
	}

	// Performance Summary
	report.WriteString("--- Performance Summary ---\n")
	avgDuration := rg.getAverageDuration()
	report.WriteString(fmt.Sprintf("Average Test Duration: %v\n", avgDuration))
	if maxDuration := rg.getMaxDuration(); maxDuration > 0 {
		report.WriteString(fmt.Sprintf("Max Test Duration: %v\n", maxDuration))
	}
	report.WriteString("\n")

	// Recommendations
	report.WriteString("--- Recommendations ---\n")
	recommendations := rg.getRecommendations()
	for _, rec := range recommendations {
		report.WriteString(fmt.Sprintf("• %s\n", rec))
	}

	return report.String()
}

// getStatusString returns a status string
func (rg *ReportGenerator) getStatusString() string {
	if rg.TestResults.Passed {
		return "✅ PASSED"
	}
	return "❌ FAILED"
}

// getPassRate calculates the pass rate
func (rg *ReportGenerator) getPassRate() float64 {
	if rg.TestResults.TotalTests == 0 {
		return 0
	}
	return float64(rg.TestResults.PassedTests) / float64(rg.TestResults.TotalTests) * 100
}

// getAverageDuration calculates average test duration
func (rg *ReportGenerator) getAverageDuration() time.Duration {
	if rg.TestResults.TotalTests == 0 {
		return 0
	}

	var total time.Duration
	for _, result := range rg.TestResults.TestCases {
		total += result.Duration
	}

	return total / time.Duration(rg.TestResults.TotalTests)
}

// getMaxDuration finds the maximum test duration
func (rg *ReportGenerator) getMaxDuration() time.Duration {
	var max time.Duration
	for _, result := range rg.TestResults.TestCases {
		if result.Duration > max {
			max = result.Duration
		}
	}
	return max
}

// getRecommendations generates recommendations based on test results
func (rg *ReportGenerator) getRecommendations() []string {
	var recommendations []string

	// Performance recommendations
	if avgDuration := rg.getAverageDuration(); avgDuration > 5*time.Second {
		recommendations = append(recommendations, "Some tests are taking longer than expected. Consider optimizing test performance.")
	}

	// Failure rate recommendations
	passRate := rg.getPassRate()
	if passRate < 90 {
		recommendations = append(recommendations, "Test failure rate is high. Consider investigating failing tests.")
	}

	// Visual diff recommendations
	hasHighDiffs := false
	for _, result := range rg.TestResults.TestCases {
		if !result.Passed && result.DiffRatio > 0.05 {
			hasHighDiffs = true
			break
		}
	}
	if hasHighDiffs {
		recommendations = append(recommendations, "Some tests have high visual diff ratios. Check if changes are intentional or if baselines need updating.")
	}

	// Environment recommendations
	if rg.TestResults.Environment == "" {
		recommendations = append(recommendations, "Test environment information is missing. Consider adding environment details.")
	}

	// Default recommendations
	if len(recommendations) == 0 {
		recommendations = append(recommendations, "All tests passed! Consider adding more comprehensive visual tests.")
	}

	return recommendations
}

// getHTMLTemplate returns the HTML template for the report
func (rg *ReportGenerator) getHTMLTemplate() string {
	return `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{{.Name}} - Visual Test Report</title>
    <link rel="stylesheet" href="visual-test-report.css">
</head>
<body>
    <div class="container">
        <header>
            <h1>{{.Name}}</h1>
            <p class="description">{{.Description}}</p>
            <div class="meta">
                <span class="run-at">Run At: {{.RunAt}}</span>
                <span class="duration">Duration: {{.Duration}}</span>
                <span class="environment">Environment: {{.Environment}}</span>
            </div>
        </header>

        <section class="summary">
            <h2>Summary</h2>
            <div class="summary-grid">
                <div class="summary-item {{if .Passed}}passed{{else}}failed{{end}}">
                    <div class="summary-value">{{if .Passed}}✅ PASSED{{else}}❌ FAILED{{end}}</div>
                    <div class="summary-label">Status</div>
                </div>
                <div class="summary-item">
                    <div class="summary-value">{{.TotalTests}}</div>
                    <div class="summary-label">Total Tests</div>
                </div>
                <div class="summary-item">
                    <div class="summary-value">{{.PassedTests}}</div>
                    <div class="summary-label">Passed</div>
                </div>
                <div class="summary-item">
                    <div class="summary-value">{{.FailedTests}}</div>
                    <div class="summary-label">Failed</div>
                </div>
                <div class="summary-item">
                    <div class="summary-value">{{.SkippedTests}}</div>
                    <div class="summary-label">Skipped</div>
                </div>
                <div class="summary-item">
                    <div class="summary-value">{{printf "%.1f" .PassRate}}%</div>
                    <div class="summary-label">Pass Rate</div>
                </div>
            </div>
        </section>

        {{if gt .FailedTests 0}}
        <section class="failed-tests">
            <h2>Failed Tests</h2>
            {{range $name, $result := .TestCases}}
                {{if and (not $result.Passed) (not $result.Skipped)}}
                <div class="test-case failed">
                    <h3>{{$name}}</h3>
                    <div class="test-details">
                        <div class="test-error">
                            <strong>Error:</strong> {{$result.Error}}
                        </div>
                        {{if gt $result.DiffRatio 0}}
                        <div class="test-diff">
                            <strong>Diff Ratio:</strong> {{printf "%.2f" .DiffRatioPercent}}%
                        </div>
                        {{end}}
                        <div class="test-duration">
                            <strong>Duration:</strong> {{$result.Duration}}
                        </div>
                        {{if $result.Screenshot}}
                        <div class="test-screenshot">
                            <strong>Screenshot:</strong> <a href="{{$result.Screenshot}}">{{$result.ScreenshotBase}}</a>
                        </div>
                        {{end}}
                        {{if $result.DiffImage}}
                        <div class="test-diff-image">
                            <strong>Diff Image:</strong> <a href="{{$result.DiffImage}}">{{$result.DiffImageBase}}</a>
                        </div>
                        {{end}}
                    </div>
                </div>
                {{end}}
            {{end}}
        </section>
        {{end}}

        <section class="all-tests">
            <h2>All Test Cases</h2>
            <table class="test-table">
                <thead>
                    <tr>
                        <th>Test Name</th>
                        <th>Status</th>
                        <th>Duration</th>
                        <th>Diff Ratio</th>
                        <th>Details</th>
                    </tr>
                </thead>
                <tbody>
                    {{range $name, $result := .TestCases}}
                    <tr class="{{if $result.Passed}}passed{{else if $result.Skipped}}skipped{{else}}failed{{end}}">
                        <td>{{$name}}</td>
                        <td>
                            {{if $result.Passed}}✅ Passed
                            {{else if $result.Skipped}}⏭️ Skipped
                            {{else}}❌ Failed
                            {{end}}
                        </td>
                        <td>{{$result.Duration}}</td>
                        <td>{{if gt $result.DiffRatio 0}}{{printf "%.2f" .DiffRatioPercent}}%{{else}}-{{end}}</td>
                        <td>{{$result.Details}}</td>
                    </tr>
                    {{end}}
                </tbody>
            </table>
        </section>

        <section class="performance">
            <h2>Performance</h2>
            <div class="performance-stats">
                <div class="stat-item">
                    <div class="stat-value">{{.AverageDuration}}</div>
                    <div class="stat-label">Average Duration</div>
                </div>
                <div class="stat-item">
                    <div class="stat-value">{{.MaxDuration}}</div>
                    <div class="stat-label">Max Duration</div>
                </div>
            </div>
        </section>

        <section class="recommendations">
            <h2>Recommendations</h2>
            <ul class="recommendations-list">
                {{range .Recommendations}}
                <li>{{.}}</li>
                {{end}}
            </ul>
        </section>
    </div>
</body>
</html>`
}

// getCSSContent returns the CSS content for the HTML report
func (rg *ReportGenerator) getCSSContent() string {
	return `/* Visual Test Report CSS */
* {
    margin: 0;
    padding: 0;
    box-sizing: border-box;
}

body {
    font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
    line-height: 1.6;
    color: #333;
    background-color: #f5f5f5;
}

.container {
    max-width: 1200px;
    margin: 0 auto;
    padding: 20px;
}

header {
    background: white;
    padding: 2rem;
    border-radius: 8px;
    box-shadow: 0 2px 4px rgba(0,0,0,0.1);
    margin-bottom: 2rem;
}

header h1 {
    font-size: 2rem;
    margin-bottom: 0.5rem;
    color: #2c3e50;
}

.description {
    color: #666;
    margin-bottom: 1rem;
}

.meta {
    display: flex;
    gap: 1rem;
    font-size: 0.9rem;
    color: #666;
}

.summary {
    background: white;
    padding: 2rem;
    border-radius: 8px;
    box-shadow: 0 2px 4px rgba(0,0,0,0.1);
    margin-bottom: 2rem;
}

.summary-grid {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(150px, 1fr));
    gap: 1rem;
    margin-top: 1rem;
}

.summary-item {
    text-align: center;
    padding: 1rem;
    border-radius: 4px;
    background: #f8f9fa;
}

.summary-item.passed {
    background: #d4edda;
    color: #155724;
}

.summary-item.failed {
    background: #f8d7da;
    color: #721c24;
}

.summary-value {
    font-size: 1.5rem;
    font-weight: bold;
    margin-bottom: 0.25rem;
}

.summary-label {
    font-size: 0.9rem;
    color: #666;
}

.test-case {
    background: white;
    margin-bottom: 1rem;
    border-radius: 8px;
    overflow: hidden;
    box-shadow: 0 2px 4px rgba(0,0,0,0.1);
}

.test-case.failed {
    border-left: 4px solid #dc3545;
}

.test-case h3 {
    padding: 1rem;
    background: #f8f9fa;
    border-bottom: 1px solid #dee2e6;
}

.test-details {
    padding: 1rem;
}

.test-details > div {
    margin-bottom: 0.5rem;
}

.test-table {
    width: 100%;
    border-collapse: collapse;
    background: white;
    border-radius: 8px;
    overflow: hidden;
    box-shadow: 0 2px 4px rgba(0,0,0,0.1);
}

.test-table th,
.test-table td {
    padding: 1rem;
    text-align: left;
    border-bottom: 1px solid #dee2e6;
}

.test-table th {
    background: #f8f9fa;
    font-weight: 600;
}

.test-table tr.passed {
    background: #d4edda;
}

.test-table tr.failed {
    background: #f8d7da;
}

.test-table tr.skipped {
    background: #fff3cd;
}

.performance-stats {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
    gap: 1rem;
}

.stat-item {
    text-align: center;
    padding: 1rem;
    background: white;
    border-radius: 8px;
    box-shadow: 0 2px 4px rgba(0,0,0,0.1);
}

.stat-value {
    font-size: 1.2rem;
    font-weight: bold;
    color: #2c3e50;
}

.stat-label {
    font-size: 0.9rem;
    color: #666;
}

.recommendations-list {
    list-style: none;
    padding: 0;
}

.recommendations-list li {
    background: white;
    padding: 1rem;
    margin-bottom: 0.5rem;
    border-radius: 4px;
    border-left: 4px solid #007bff;
    box-shadow: 0 2px 4px rgba(0,0,0,0.1);
}

h2 {
    font-size: 1.5rem;
    margin-bottom: 1rem;
    color: #2c3e50;
}

a {
    color: #007bff;
    text-decoration: none;
}

a:hover {
    text-decoration: underline;
}

@media (max-width: 768px) {
    .container {
        padding: 10px;
    }

    .meta {
        flex-direction: column;
        gap: 0.5rem;
    }

    .summary-grid {
        grid-template-columns: repeat(2, 1fr);
    }

    .test-table {
        font-size: 0.9rem;
    }

    .test-table th,
    .test-table td {
        padding: 0.5rem;
    }
}`
}