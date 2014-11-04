package waveform

import (
	"image/color"
	"math/rand"
	"time"
)

// ColorFunc is a function which accepts a computed value count, X and Y
// coordinates, and a default color value.
//
// A ColorFunc is applied during each image drawing iteration, and will
// return the appropriate color which should be drawn at the specified X and Y
// coordinate, based upon the return of the function.
type ColorFunc func(count int, x int, y int, color color.Color) color.Color

// SolidColor is a ColorFunc which simply returns the input, default color
// as the color which should be drawn at all coordinates.  This is the default
// behavior of the waveform package.
func SolidColor(count int, x int, y int, color color.Color) color.Color {
	return color
}

// FuzzColor generates a ColorFunc which applies a random color on each call,
// selected from an input, variadic slice of colors.  This can be used to create
// a random fuzz or "static" effect in the resulting waveform image.
func FuzzColor(colors ...color.Color) ColorFunc {
	// Filter any nil values
	colors = filterNilColors(colors)

	// Seed RNG
	rand.Seed(time.Now().UnixNano())

	// Select a color at random on each call
	return func(count int, x int, y int, inColor color.Color) color.Color {
		return colors[rand.Intn(len(colors))]
	}
}

// StripeColor generates a ColorFunc which applies one color from the input,
// variadic slice at each computed value.  Each color is used in order, and
// the rotation will repeat until the image is complete. This creates a stripe
// effect in the resulting waveform image.
func StripeColor(colors ...color.Color) ColorFunc {
	// Filter any nil values
	colors = filterNilColors(colors)

	var lastCount int
	return func(count int, x int, y int, inColor color.Color) color.Color {
		// For each new count value, use the next color in the slice
		if count > lastCount {
			lastCount = count
		}

		return colors[lastCount%len(colors)]
	}
}

// filterNilColors strips any nil color.Color values from the input slice.
func filterNilColors(colors []color.Color) []color.Color {
	var cleanColors []color.Color
	for _, c := range colors {
		if c != nil {
			cleanColors = append(cleanColors, c)
		}
	}

	return cleanColors
}
