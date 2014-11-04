package waveform

import (
	"image/color"
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

// AlternateColor generates a ColorFunc which applies an alternate color
// on alternating X-axis values.  This can be used to create a stripe effect
// in the resulting waveform image.
func AlternateColor(alt color.Color) ColorFunc {
	return func(count int, x int, y int, color color.Color) color.Color {
		// On odd iterations (or if no alternate set), draw using specified
		// foreground color at specified X and Y coordinate
		if count%2 != 0 || alt == nil {
			return color
		} else {
			// On even iterations, draw using specified alternate color
			return alt
		}
	}
}
