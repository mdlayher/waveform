package waveform

import (
	"image/color"
	"testing"
)

// TestSolidColor verifies that SolidColor always returns the same input
// color, for all input values.
func TestSolidColor(t *testing.T) {
	colors := []color.Color{
		color.RGBA{0, 0, 0, 255},
		color.RGBA{255, 0, 0, 255},
		color.RGBA{0, 255, 0, 255},
		color.RGBA{0, 255, 255, 255},
	}

	for i, c := range colors {
		if out := SolidColor(i, i, i, c); out != c {
			t.Fatalf("unexpected SolidColor color: %v != %v", out, c)
		}
	}
}

// TestAlternateColor verifies that AlternateColor returns appropriate
// alternated color values for each iteration of a counter.
func TestAlternateColor(t *testing.T) {
	// Set up known values
	black := color.RGBA{0, 0, 0, 255}
	white := color.RGBA{255, 255, 255, 255}

	// Generate color function, start loop
	colorFn := AlternateColor(black)
	for i := 0; i < 10000; i++ {
		out := colorFn(i, i, i, white)

		// On even iterations, alternate color is used
		if i%2 == 0 {
			if out != black {
				t.Fatalf("[%02d] unexpected even AlternateColor color: %v != %v", i, out, black)
			}
		} else {
			// On odd iterations, primary color is used
			if out != white {
				t.Fatalf("[%02d] unexpected odd AlternateColor color: %v != %v", i, out, white)
			}
		}
	}
}
