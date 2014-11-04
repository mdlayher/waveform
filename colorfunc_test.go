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

// TestSolidColor verifies that SolidColor always returns the same input
// color, for all input values.
func TestSolidColor(t *testing.T) {
	colors := []color.Color{
		black,
		white,
	}

	for i, c := range colors {
		if out := SolidColor(i, i, i, c); out != c {
			t.Fatalf("unexpected SolidColor color: %v != %v", out, c)
		}
	}
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
		if out, ok := set[fn(i, i, i, nil).(color.RGBA)]; !ok {
			t.Fatalf("color not in set: %v", out)
		}
	}
}

// testStripeColor is a test helper which aids in testing the StripeColor function.
func testStripeColor(t *testing.T, in []color.Color, out []color.Color) {
	// Validate that StripeColor produces expected output at each index
	fn := StripeColor(in...)
	for i := 0; i < len(out); i++ {
		if c := fn(i, 0, 0, nil); c != out[i] {
			t.Fatalf("[%02d] unexpected output color: %v != %v", i, c, out[i])
		}
	}
}
