package waveform

import (
	"image/color"
	"testing"
)

// Named colors for easy testing
var (
	black = color.RGBA{0, 0, 0, 255}
	white = color.RGBA{255, 255, 255, 255}
	red   = color.RGBA{255, 0, 0, 255}
	green = color.RGBA{0, 255, 0, 255}
	blue  = color.RGBA{0, 0, 255, 255}
)

// TestCheckerColorOneColor verifies that CheckerColor produces only the single
// color used in its input.
func TestCheckerColorOneColor(t *testing.T) {
	testCheckerColor(t, black, black)
}

// TestCheckerColorTwoColors verifies that CheckerColor produces only colors
// which are used in its input.
func TestCheckerColorTwoColors(t *testing.T) {
	testCheckerColor(t, black, white)
}

// TestFuzzColorOneColor verifies that FuzzColor produces only the single
// color used in its input.
func TestFuzzColorOneColor(t *testing.T) {
	testFuzzColor(t, []color.Color{black})
}

// TestFuzzColorMultipleColors verifies that FuzzColor produces only colors
// which are used in its input.
func TestFuzzColorMultipleColors(t *testing.T) {
	testFuzzColor(t, []color.Color{black, white, red, green, blue})
}

// TestGradientColorOneColor verifies that GradientColor produces only the single
// color used in its input.
func TestGradientColorOneColor(t *testing.T) {
	testGradientColor(t, black, black)
}

// TestGradientColorTwoColors verifies that GradientColor produces a correct
// gradient between two colors.
func TestGradientColorTwoColors(t *testing.T) {
	testGradientColor(t, black, white)
}

// TestSolidColor verifies that SolidColor always returns the same input
// color, for all input values.
func TestSolidColor(t *testing.T) {
	colors := []color.Color{
		black,
		white,
	}

	for i, c := range colors {
		fn := SolidColor(c)
		if out := fn(i, i, i, i, i, i); out != c {
			t.Fatalf("unexpected SolidColor color: %v != %v", out, c)
		}
	}
}

// TestStripeColorOneColor verifies that StripeColor produces a correct
// color sequence with a single input color.
func TestStripeColorOneColor(t *testing.T) {
	testStripeColor(t, []color.Color{black}, []color.Color{
		black, black, black, black,
	})
}

// TestStripeColorMultipleColors verifies that StripeColor produces a correct
// color sequence with multiple input colors.
func TestStripeColorMultipleColors(t *testing.T) {
	testStripeColor(t, []color.Color{
		black, white, white, red, green, green, green, blue,
	}, []color.Color{
		black, white, white, red, green, green, green, blue,
		black, white, white, red, green, green, green, blue,
	})
}

// testCheckerColor is a test helper which aids in testing the CheckerColor function.
func testCheckerColor(t *testing.T, colorA color.Color, colorB color.Color) {
	// Predefined values for test
	const maxX, maxY, size = 1000, 1000, 10

	// Generate checker function with input values
	fn := CheckerColor(colorA, colorB, size)

	// Iterate all coordinates and check color at each
	for x := 0; x < maxX; x++ {
		for y := 0; y < maxY; y++ {
			// Check color at specified coordinate
			c := fn(0, x, y, 0, maxX, maxY)

			// Apply checker algorithm to determine if color A or B should be used
			if ((uint(x)/size)+(uint(y)/size))%2 == 0 {
				if c != colorA {
					t.Fatalf("unexpected color: %v != %v", c, colorA)
				}
			} else {
				if c != colorB {
					t.Fatalf("unexpected color: %v != %v", c, colorB)
				}
			}
		}
	}
}

// testFuzzColor is a test helper which aids in testing the FuzzColor function.
func testFuzzColor(t *testing.T, in []color.Color) {
	// Make a set of colors from input slice
	set := make(map[color.RGBA]struct{})
	for _, c := range in {
		set[c.(color.RGBA)] = struct{}{}
	}

	// Validate that FuzzColor only produces colors which are present in
	// the input slice.
	fn := FuzzColor(in...)
	for i := 0; i < 10000; i++ {
		if out, ok := set[fn(i, i, i, i, i, i).(color.RGBA)]; !ok {
			t.Fatalf("color not in set: %v", out)
		}
	}
}

// testGradientColor is a test helper which aids in testing the GradientColor
// function.
func testGradientColor(t *testing.T, start color.RGBA, end color.RGBA) {
	const maxN = 100

	// Generate function with defined values
	fn := GradientColor(start, end)

	// Check edges
	for i, n := range []int{0, maxN} {
		// Get color at point, get RGBA equivalent
		c := fn(n, 0, 0, maxN, 0, 0)
		r, g, b, _ := c.RGBA()

		// First iteration, use start; second, use end
		var testColor color.RGBA
		if i == 0 {
			testColor = start
		} else {
			testColor = end
		}

		// Compare values to ensure correctness
		if testColor.R != uint8(r) || testColor.G != uint8(g) || testColor.B != uint8(b) {
			t.Fatalf("unexpected color at %d%%: %v != %v", n, c, testColor)
		}
	}
}

// testStripeColor is a test helper which aids in testing the StripeColor function.
func testStripeColor(t *testing.T, in []color.Color, out []color.Color) {
	// Validate that StripeColor produces expected output at each index
	fn := StripeColor(in...)
	for i := 0; i < len(out); i++ {
		if c := fn(i, 0, 0, 0, 0, 0); c != out[i] {
			t.Fatalf("[%02d] unexpected output color: %v != %v", i, c, out[i])
		}
	}
}
