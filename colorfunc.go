package waveform

import (
	"image/color"
	"math"
	"math/rand"
	"time"
)

// ColorFunc is a function which accepts a variety of values which can be used
// to customize an output image.  These values include the current computed sample
// count (n), the current X coordinate (x), the current Y coordinate (y), and the
// maximum computed values for each of these. (maxN, maxX, maxY)
//
// A ColorFunc is applied during each image drawing iteration, and will
// return the appropriate color which should be drawn at the specified value
// for n, x, and y; possibly taking into account their maximum values.
type ColorFunc func(n int, x int, y int, maxN int, maxX int, maxY int) color.Color

// FuzzColor generates a ColorFunc which applies a random color on each call,
// selected from an input, variadic slice of colors.  This can be used to create
// a random fuzz or "static" effect in the resulting waveform image.
func FuzzColor(colors ...color.Color) ColorFunc {
	// Filter any nil values
	colors = filterNilColors(colors)

	// Seed RNG
	rand.Seed(time.Now().UnixNano())

	// Select a color at random on each call
	return func(n int, x int, y int, maxN int, maxX int, maxY int) color.Color {
		return colors[rand.Intn(len(colors))]
	}
}

// GradientColor generates a ColorFunc which produces a color gradient between two
// RGBA input colors.  The gradient attempts to gradually reduce the distance between
// two colors, creating a sweeping color change effect in the resulting waveform
// image.
func GradientColor(start color.RGBA, end color.RGBA) ColorFunc {
	// Current value starts at the first value, and gradually ascends or descends
	// to the second.
	current := start

	// Calculate absolute distances between RGB values
	distR := distanceUint8(start.R, end.R)
	distG := distanceUint8(start.G, end.G)
	distB := distanceUint8(start.B, end.B)

	// Store last calculated values from each call
	var lastR, lastG, lastB uint8

	// Store last seen N value to calculate gradient
	var lastN int
	return func(n int, x int, y int, maxN int, maxX int, maxY int) color.Color {
		// On first iteration, always return start value
		if n == 0 {
			return start
		}

		// On final iteration, always return end value
		if n == maxN {
			return end
		}

		// Unless a new n value is entered, use the same color
		if n <= lastN {
			return current
		}

		// New n value, store it now
		lastN = n

		// Return the current color the vast majority of the time, so that
		// the gradient is very gradual.
		//
		// This math is completely made up, but appears to work well enough
		// for our purposes.
		if n%(int(math.Pow(float64(maxN/255), 2))) != 0 {
			return current
		}

		// Recalculate current values, taking into account known starting points,
		// ending points, etc.
		current.R = stepValue(start.R > end.R, maxN, distR, current.R, lastR)
		lastR = current.R

		current.G = stepValue(start.R > end.R, maxN, distG, current.G, lastG)
		lastG = current.G

		current.B = stepValue(start.R > end.R, maxN, distB, current.B, lastB)
		lastB = current.B

		return current
	}
}

// SolidColor generates a ColorFunc which simply returns the input color
// as the color which should be drawn at all coordinates.
//
// This is the default behavior of the waveform package.
func SolidColor(inColor color.Color) ColorFunc {
	return func(n int, x int, y int, maxN int, maxX int, maxY int) color.Color {
		return inColor
	}
}

// StripeColor generates a ColorFunc which applies one color from the input,
// variadic slice at each computed value.  Each color is used in order, and
// the rotation will repeat until the image is complete. This creates a stripe
// effect in the resulting waveform image.
func StripeColor(colors ...color.Color) ColorFunc {
	// Filter any nil values
	colors = filterNilColors(colors)

	var lastN int
	return func(n int, x int, y int, maxN int, maxX int, maxY int) color.Color {
		// For each new n value, use the next color in the slice
		if n > lastN {
			lastN = n
		}

		return colors[lastN%len(colors)]
	}
}

// distanceUint8 calculates the absolute value of the distance between
// two uint8 values.
func distanceUint8(start uint8, end uint8) uint8 {
	if start > end {
		return start - end
	}

	return end - start
}

// stepValue returns a new uint8 color value, based upon a variety of input values.
// TODO(mdlayher): clean this up or break it into smaller pieces
func stepValue(greater bool, maxN int, dist uint8, current uint8, last uint8) uint8 {
	// If no distance, always return current value
	if dist == 0 {
		return current
	}

	// Calculate step based upon max N and distance
	step := uint8(math.Floor(float64(maxN) / float64(dist)))

	// If start is greater, prevent underflow
	if greater {
		if current > 0 && (current-step) < last {
			current -= step
		}
	} else {
		// If end is greater, prevent overflow
		if current < 255 && (current+step) > last {
			current += step
		}
	}

	return current
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
