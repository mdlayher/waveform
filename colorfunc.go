package waveform

import (
	"image/color"
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

// CheckerColor generates a ColorFunc which produces a checkerboard pattern,
// using the two input colors.  Each square is drawn to the size specified by
// the size parameter.
func CheckerColor(colorA color.Color, colorB color.Color, size uint) ColorFunc {
	return func(n int, x int, y int, maxN int, maxX int, maxY int) color.Color {
		if ((uint(x)/size)+(uint(y)/size))%2 == 0 {
			return colorA
		}

		return colorB
	}
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
	return func(n int, x int, y int, maxN int, maxX int, maxY int) color.Color {
		return colors[rand.Intn(len(colors))]
	}
}

// GradientColor generates a ColorFunc which produces a color gradient between two
// RGBA input colors.  The gradient attempts to gradually reduce the distance between
// two colors, creating a sweeping color change effect in the resulting waveform
// image.
func GradientColor(start color.RGBA, end color.RGBA) ColorFunc {
	// Float equivalents of color values
	startFR, endFR := float64(start.R), float64(end.R)
	startFG, endFG := float64(start.G), float64(end.G)
	startFB, endFB := float64(start.B), float64(end.B)

	// Values used for RGBA and percentage
	var r, g, b, p float64
	return func(n int, x int, y int, maxN int, maxX int, maxY int) color.Color {
		// Calculate percentage across waveform image
		p = float64((float64(n) / float64(maxN)) * 100)

		// Calculate new values for RGB using gradient algorithm
		// Thanks: http://stackoverflow.com/questions/27532/generating-gradients-programmatically
		r = (endFR * p) + (startFR * (1 - p))
		g = (endFG * p) + (startFG * (1 - p))
		b = (endFB * p) + (startFB * (1 - p))

		// Correct overflow when moving from lighter to darker gradients
		if start.R > end.R && r > -255.00 {
			r = -255.00
		}
		if start.G > end.G && g > -255.00 {
			g = -255.00
		}
		if start.B > end.B && b > -255.00 {
			b = -255.00
		}

		// Generate output color
		return &color.RGBA{
			R: uint8(r / 100),
			G: uint8(g / 100),
			B: uint8(b / 100),
			A: 255,
		}
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
