package waveform

import (
	"fmt"
	"image/color"
)

var (
	// errColorsNilForeground is returned when a nil foreground color is
	// used in a call to Colors.
	errColorsNilForeground = &OptionsError{
		Option: "colors",
		Reason: "foreground color cannot be nil",
	}

	// errColorsNilBackground is returned when a nil background color is
	// used in a call to Colors.
	errColorsNilBackground = &OptionsError{
		Option: "colors",
		Reason: "background color cannot be nil",
	}

	// errFGColorFunctionNil is returned when a nil ColorFunc is used in
	// a call to FGColorFunction.
	errFGColorFunctionNil = &OptionsError{
		Option: "fgColorFunction",
		Reason: "function cannot be nil",
	}

	// errSampleFunctionNil is returned when a nil SampleReduceFunc is used in
	// a call to SampleFunc.
	errSampleFunctionNil = &OptionsError{
		Option: "sampleFunction",
		Reason: "function cannot be nil",
	}

	// errResolutionZero is returned when integer 0 is used in a call
	// to Resolution.
	errResolutionZero = &OptionsError{
		Option: "resolution",
		Reason: "resolution cannot be 0",
	}

	// errScaleXZero is returned when integer 0 is used as the X value
	// in a call to Scale.
	errScaleXZero = &OptionsError{
		Option: "scale",
		Reason: "X scale cannot be 0",
	}

	// errScaleYZero is returned when integer 0 is used as the Y value
	// in a call to Scale.
	errScaleYZero = &OptionsError{
		Option: "scale",
		Reason: "Y scale cannot be 0",
	}
)

// OptionsError is an error which is returned when invalid input
// options are set on a Waveform struct.
type OptionsError struct {
	Option string
	Reason string
}

// Error returns the string representation of an OptionsError.
func (e *OptionsError) Error() string {
	return fmt.Sprintf("%s: %s", e.Option, e.Reason)
}

// OptionsFunc is a function which is applied to an input Waveform
// struct, and can manipulate its properties.
type OptionsFunc func(*Waveform) error

// SetOptions applies zero or more OptionsFunc to the receiving Waveform
// struct, manipulating its properties.
func (w *Waveform) SetOptions(options ...OptionsFunc) error {
	for _, o := range options {
		// Do not apply nil function arguments
		if o == nil {
			continue
		}

		if err := o(w); err != nil {
			return err
		}
	}

	return nil
}

// Colors generates an OptionsFunc which applies the input foreground
// and background colors to an input Waveform struct.
func Colors(fg color.Color, bg color.Color) OptionsFunc {
	return func(w *Waveform) error {
		return w.setColors(fg, bg)
	}
}

// SetColors applies the input foreground and background color
// to the receiving Waveform struct.
func (w *Waveform) SetColors(fg color.Color, bg color.Color) error {
	return w.SetOptions(Colors(fg, bg))
}

// setColors directly modifies the foreground and background color
// members of the receiving struct.
func (w *Waveform) setColors(fg color.Color, bg color.Color) error {
	// Foreground color cannot be nil
	if fg == nil {
		return errColorsNilForeground
	}

	// Background color cannot be nil
	if bg == nil {
		return errColorsNilBackground
	}

	w.fg = fg
	w.bg = bg

	return nil
}

// FGColorFunction generates an OptionsFunc which applies the input foreground
// ColorFunc to an input Waveform struct.
//
// This function is used to apply a variety of color schemes to the foreground
// of a waveform image, and is called during each drawing loop of the foreground
// image.
func FGColorFunction(function ColorFunc) OptionsFunc {
	return func(w *Waveform) error {
		return w.setFGColorFunction(function)
	}
}

// SetFGColorFunction applies the input ColorFunc to the receiving Waveform
// struct for foreground use.
func (w *Waveform) SetFGColorFunction(function ColorFunc) error {
	return w.SetOptions(FGColorFunction(function))
}

// setFGColorFunction directly sets the foreground ColorFunc member of the
// receiving Waveform struct.
func (w *Waveform) setFGColorFunction(function ColorFunc) error {
	// Function cannot be nil
	if function == nil {
		return errFGColorFunctionNil
	}

	w.fgColorFn = function

	return nil
}

// Resolution generates an OptionsFunc which applies the input resolution
// value to an input Waveform struct.
//
// This value indicates the number of times audio is read and drawn
// as a waveform, per second of audio.
func Resolution(resolution uint) OptionsFunc {
	return func(w *Waveform) error {
		return w.setResolution(resolution)
	}
}

// SetResolution applies the input resolution to the receiving Waveform struct.
func (w *Waveform) SetResolution(resolution uint) error {
	return w.SetOptions(Resolution(resolution))
}

// setResolution directly sets the resolution member of the receiving Waveform
// struct.
func (w *Waveform) setResolution(resolution uint) error {
	// Resolution cannot be zero
	if resolution == 0 {
		return errResolutionZero
	}

	w.resolution = resolution

	return nil
}

// SampleFunc generates an OptionsFunc which applies the input SampleReduceFunc
// to an input Waveform struct.
//
// This function is used to compute values from audio samples, for use in
// waveform generation.  The function is applied over a slice of float64
// audio samples, reducing them to a single value.
func SampleFunction(function SampleReduceFunc) OptionsFunc {
	return func(w *Waveform) error {
		return w.setSampleFunction(function)
	}
}

// SetSampleFunction applies the input SampleReduceFunc to the receiving Waveform
// struct.
func (w *Waveform) SetSampleFunction(function SampleReduceFunc) error {
	return w.SetOptions(SampleFunction(function))
}

// setSampleFunction directly sets the SampleReduceFunc member of the receiving
// Waveform struct.
func (w *Waveform) setSampleFunction(function SampleReduceFunc) error {
	// Function cannot be nil
	if function == nil {
		return errSampleFunctionNil
	}

	w.sampleFn = function

	return nil
}

// Scale generates an OptionsFunc which applies the input X and Y axis scaling
// factors to an input Waveform struct.
//
// This value indicates how a generated waveform image will be scaled, for both
// its X and Y axes.
func Scale(x uint, y uint) OptionsFunc {
	return func(w *Waveform) error {
		return w.setScale(x, y)
	}
}

// SetScale applies the input X and Y axis scaling to the receiving Waveform
// struct.
func (w *Waveform) SetScale(x uint, y uint) error {
	return w.SetOptions(Scale(x, y))
}

// setScale directly sets the scaleX and scaleY members of the receiving Waveform
// struct.
func (w *Waveform) setScale(x uint, y uint) error {
	// X scale cannot be zero
	if x == 0 {
		return errScaleXZero
	}

	// Y scale cannot be zero
	if y == 0 {
		return errScaleYZero

	}

	w.scaleX = x
	w.scaleY = y

	return nil
}

// ScaleClipping generates an OptionsFunc which sets the scaleClipping member
// to true on an input Waveform struct.
//
// This value indicates if the waveform image should be scaled down on its Y-axis
// when clipping thresholds are reached.  This can be used to show a more accurate
// waveform when the input audio stream exhibits signs of clipping.
func ScaleClipping() OptionsFunc {
	return func(w *Waveform) error {
		return w.setScaleClipping(true)
	}
}

// SetScaleClipping applies sets the scaleClipping member true for the receiving
// Waveform struct.
func (w *Waveform) SetScaleClipping() error {
	return w.SetOptions(ScaleClipping())
}

// setScaleClipping directly sets the scaleClipping member of the receiving Waveform
// struct.
func (w *Waveform) setScaleClipping(scaleClipping bool) error {
	w.scaleClipping = scaleClipping

	return nil
}

// Sharpness generates an OptionsFunc which applies the input sharpness
// value to an input Waveform struct.
//
// This value indicates the amount of curvature which is applied to a
// waveform image, scaled on its X-axis.  A higher value results in steeper
// curves, and a lower value results in more "blocky" curves.
func Sharpness(sharpness uint) OptionsFunc {
	return func(w *Waveform) error {
		return w.setSharpness(sharpness)
	}
}

// SetSharpness applies the input sharpness to the receiving Waveform struct.
func (w *Waveform) SetSharpness(sharpness uint) error {
	return w.SetOptions(Sharpness(sharpness))
}

// setSharpness directly sets the sharpness member of the receiving Waveform
// struct.
func (w *Waveform) setSharpness(sharpness uint) error {
	w.sharpness = sharpness

	return nil
}
