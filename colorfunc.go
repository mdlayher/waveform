package waveform

import (
	"image/color"
)

// ColorFunc is a function which accepts an input X and Y coordinate, and a
// default color value.  It will return the appropriate color which should be
// drawn at the specified X and Y coordinate.
type ColorFunc func(x int, y int, color color.Color) color.Color

// SolidColor is a ColorFunc which simply returns the input, default color
// as the color which should be drawn at all coordinates.  This is the default
// behavior of the waveform package.
func SolidColor(x int, y int, color color.Color) color.Color {
	return color
}
