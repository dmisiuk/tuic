package visual

import (
	"fmt"
	"image"
	"image/color"
	"math"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// ComparisonResult represents the result of a visual comparison
type ComparisonResult struct {
	Identical      bool
	DiffRatio      float64
	DiffImage      *image.RGBA
	PixelDiffs     int
	TotalPixels    int
	DiffRegions    []DiffRegion
	Metrics        ComparisonMetrics
}

// DiffRegion represents a region where differences were found
type DiffRegion struct {
	X      int
	Y      int
	Width  int
	Height int
	Reason string
}

// ComparisonMetrics contains detailed comparison metrics
type ComparisonMetrics struct {
	ColorDistance    float64
	BrightnessDiff   float64
	ContrastDiff     float64
	StructuralDiff   float64
	PerceptualDiff   float64
}

// CompareConfig contains configuration for comparison
type CompareConfig struct {
	ColorTolerance    float64
	IgnoreAntiAliasing bool
	IgnoreMotionBlur  bool
	PerceptualMode    bool
}

// NewDefaultCompareConfig creates a default comparison configuration
func NewDefaultCompareConfig() CompareConfig {
	return CompareConfig{
		ColorTolerance:    0.1,
		IgnoreAntiAliasing: true,
		IgnoreMotionBlur:  true,
		PerceptualMode:    true,
	}
}

// CompareScreenshots compares two screenshots with detailed analysis
func CompareScreenshots(screenshot1, screenshot2 *Screenshot, config CompareConfig) (*ComparisonResult, error) {
	if screenshot1.Image.Bounds() != screenshot2.Image.Bounds() {
		return nil, fmt.Errorf("screenshot dimensions don't match")
	}

	bounds := screenshot1.Image.Bounds()
	result := &ComparisonResult{
		TotalPixels: bounds.Dx() * bounds.Dy(),
		DiffImage:   image.NewRGBA(bounds),
	}

	diffMap := make([]bool, result.TotalPixels)

	// Compare pixel by pixel
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			idx := (y-bounds.Min.Y)*bounds.Dx() + (x-bounds.Min.X)

			c1 := screenshot1.Image.RGBAAt(x, y)
			c2 := screenshot2.Image.RGBAAt(x, y)

			if pixelsEqual(c1, c2, config) {
				result.DiffImage.Set(x, y, c1)
			} else {
				result.DiffImage.Set(x, y, color.RGBA{255, 0, 0, 255})
				diffMap[idx] = true
				result.PixelDiffs++
			}
		}
	}

	result.DiffRatio = float64(result.PixelDiffs) / float64(result.TotalPixels)
	result.Identical = result.PixelDiffs == 0

	// Calculate detailed metrics
	result.Metrics = calculateMetrics(screenshot1.Image, screenshot2.Image, diffMap, bounds)

	// Find diff regions
	result.DiffRegions = findDiffRegions(diffMap, bounds.Dx(), bounds.Dy())

	return result, nil
}

// pixelsEqual checks if two pixels are considered equal based on tolerance
func pixelsEqual(c1, c2 color.RGBA, config CompareConfig) bool {
	if config.IgnoreAntiAliasing && isAntiAliased(c1, c2) {
		return true
	}

	distance := colorDistance(c1, c2)
	return distance <= config.ColorTolerance
}

// isAntiAliased checks if pixels differ only due to anti-aliasing
func isAntiAliased(c1, c2 color.RGBA) bool {
	// Simple anti-aliasing detection: check if colors are similar
	// and alpha channel suggests blending
	if c1.A != 255 || c2.A != 255 {
		return true
	}

	distance := colorDistance(c1, c2)
	return distance < 0.05
}

// colorDistance calculates the Euclidean distance between two colors
func colorDistance(c1, c2 color.RGBA) float64 {
	rDiff := float64(c1.R) - float64(c2.R)
	gDiff := float64(c1.G) - float64(c2.G)
	bDiff := float64(c1.B) - float64(c2.B)
	aDiff := float64(c1.A) - float64(c2.A)

	return math.Sqrt(rDiff*rDiff + gDiff*gDiff + bDiff*bDiff + aDiff*aDiff) / (255.0 * math.Sqrt(4))
}

// calculateMetrics calculates detailed comparison metrics
func calculateMetrics(img1, img2 *image.RGBA, diffMap []bool, bounds image.Rectangle) ComparisonMetrics {
	var colorSum, brightnessSum, contrastSum float64
	var pixelCount int

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			idx := (y-bounds.Min.Y)*bounds.Dx() + (x-bounds.Min.X)

			if !diffMap[idx] {
				c1 := img1.RGBAAt(x, y)
				c2 := img2.RGBAAt(x, y)

				colorSum += colorDistance(c1, c2)
				brightnessSum += math.Abs(brightness(c1) - brightness(c2))
				contrastSum += math.Abs(contrast(c1) - contrast(c2))
				pixelCount++
			}
		}
	}

	if pixelCount == 0 {
		return ComparisonMetrics{}
	}

	return ComparisonMetrics{
		ColorDistance:    colorSum / float64(pixelCount),
		BrightnessDiff:   brightnessSum / float64(pixelCount),
		ContrastDiff:     contrastSum / float64(pixelCount),
		StructuralDiff:   calculateStructuralDiff(img1, img2),
		PerceptualDiff:   calculatePerceptualDiff(img1, img2),
	}
}

// brightness calculates the perceived brightness of a color
func brightness(c color.RGBA) float64 {
	return (0.299*float64(c.R) + 0.587*float64(c.G) + 0.114*float64(c.B)) / 255.0
}

// contrast calculates the contrast of a color
func contrast(c color.RGBA) float64 {
	lum := brightness(c)
	return (lum + 0.05) / (0.05)
}

// calculateStructuralDiff calculates structural similarity
func calculateStructuralDiff(img1, img2 *image.RGBA) float64 {
	// Simplified structural similarity calculation
	// In a real implementation, you'd use SSIM or similar algorithms
	return 0.0
}

// calculatePerceptualDiff calculates perceptual difference
func calculatePerceptualDiff(img1, img2 *image.RGBA) float64 {
	// Simplified perceptual difference calculation
	// In a real implementation, you'd use more sophisticated metrics
	return 0.0
}

// findDiffRegions identifies contiguous regions of differences
func findDiffRegions(diffMap []bool, width, height int) []DiffRegion {
	visited := make([]bool, len(diffMap))
	var regions []DiffRegion

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			idx := y*width + x
			if diffMap[idx] && !visited[idx] {
				region := exploreRegion(diffMap, visited, x, y, width, height)
				if region.Width > 0 && region.Height > 0 {
					regions = append(regions, region)
				}
			}
		}
	}

	return regions
}

// exploreRegion explores a contiguous diff region
func exploreRegion(diffMap, visited []bool, startX, startY, width, height int) DiffRegion {
	minX, maxX := startX, startX
	minY, maxY := startY, startY

	stack := []struct{ x, y int }{{startX, startY}}

	for len(stack) > 0 {
		current := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		if current.x < 0 || current.x >= width || current.y < 0 || current.y >= height {
			continue
		}

		idx := current.y*width + current.x
		if visited[idx] || !diffMap[idx] {
			continue
		}

		visited[idx] = true

		minX = min(minX, current.x)
		maxX = max(maxX, current.x)
		minY = min(minY, current.y)
		maxY = max(maxY, current.y)

		// Add neighbors
		stack = append(stack,
			struct{ x, y int }{current.x - 1, current.y},
			struct{ x, y int }{current.x + 1, current.y},
			struct{ x, y int }{current.x, current.y - 1},
			struct{ x, y int }{current.x, current.y + 1},
		)
	}

	return DiffRegion{
		X:      minX,
		Y:      minY,
		Width:  maxX - minX + 1,
		Height: maxY - minY + 1,
		Reason: "visual_difference",
	}
}

// min and max helpers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// RenderComparisonReport renders a human-readable comparison report
func (cr *ComparisonResult) RenderComparisonReport() string {
	var report strings.Builder

	report.WriteString("=== Visual Comparison Report ===\n\n")

	if cr.Identical {
		report.WriteString("✅ Screenshots are identical\n")
	} else {
		report.WriteString("❌ Screenshots differ\n")
		report.WriteString(fmt.Sprintf("   Diff Ratio: %.2f%% (%d/%d pixels)\n",
			cr.DiffRatio*100, cr.PixelDiffs, cr.TotalPixels))
		report.WriteString(fmt.Sprintf("   Diff Regions: %d\n", len(cr.DiffRegions)))

		// Show detailed metrics
		report.WriteString("\n--- Detailed Metrics ---\n")
		report.WriteString(fmt.Sprintf("Color Distance: %.4f\n", cr.Metrics.ColorDistance))
		report.WriteString(fmt.Sprintf("Brightness Diff: %.4f\n", cr.Metrics.BrightnessDiff))
		report.WriteString(fmt.Sprintf("Contrast Diff: %.4f\n", cr.Metrics.ContrastDiff))
		report.WriteString(fmt.Sprintf("Structural Diff: %.4f\n", cr.Metrics.StructuralDiff))
		report.WriteString(fmt.Sprintf("Perceptual Diff: %.4f\n", cr.Metrics.PerceptualDiff))

		// Show diff regions
		if len(cr.DiffRegions) > 0 {
			report.WriteString("\n--- Diff Regions ---\n")
			for i, region := range cr.DiffRegions {
				if i < 5 { // Limit to first 5 regions
					report.WriteString(fmt.Sprintf("Region %d: (%d,%d) %dx%d %s\n",
						i+1, region.X, region.Y, region.Width, region.Height, region.Reason))
				}
			}
			if len(cr.DiffRegions) > 5 {
				report.WriteString(fmt.Sprintf("... and %d more regions\n", len(cr.DiffRegions)-5))
			}
		}
	}

	return report.String()
}