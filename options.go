package waveform

import "image/color"

// OptionsError is an error which is returned when invalid input
// options are set on a Waveform struct.
type OptionsError struct {
	Option string
	Reason string
}

// Error returns the string representation of an OptionsError.
func (e *OptionsError) Error() string {
	return e.Option + ": " + e.Reason
}

// OptionsFunc is a function which is applied to an input Waveform
// struct, and can manipulate its properties.
type OptionsFunc func(*Waveform) error

// SetOptions applies zero or more OptionsFunc to the receiving Waveform
// struct, manipulating its properties.
func (w *Waveform) SetOptions(options ...OptionsFunc) error {
	for _, o := range options {
		if err := o(w); err != nil {
			return err
		}
	}

	return nil
}

// Colors generates an OptionsFunc which applies the input foreground,
// background, and alternate colors to an input Waveform struct.
//
// If an alternate color is specified, it will be alternated with the
// foreground color to create a stripe effect in the image.
// If alternate color is set to nil, no alternate color will be used.
func Colors(fg color.Color, bg color.Color, alt color.Color) func(*Waveform) error {
	return func(w *Waveform) error {
		return w.setColors(fg, bg, alt)
	}
}

// SetColors applies the input foreground, background, and alternate colors
// to the receiving Waveform struct.
func (w *Waveform) SetColors(fg color.Color, bg color.Color, alt color.Color) error {
	return w.SetOptions(Colors(fg, bg, alt))
}

// setColors directly modifies the foreground, background, and alternate color
// members of the receiving struct.
func (w *Waveform) setColors(fg color.Color, bg color.Color, alt color.Color) error {
	// Foreground color cannot be nil
	if fg == nil {
		return &OptionsError{
			Option: "colors",
			Reason: "foreground color cannot be nil",
		}
	}

	// Background color cannot be nil
	if bg == nil {
		return &OptionsError{
			Option: "colors",
			Reason: "background color cannot be nil",
		}
	}

	w.fg = fg
	w.bg = bg
	w.alt = alt

	return nil
}

// Function generates an OptionsFunc which applies the input function
// value to an input Waveform struct.
//
// This function is used to compute values from audio samples, for use in
// waveform generation.  The function is applied over a slice of float64
// audio samples, reducing them to a single value.
func Function(function SampleReduceFunc) func(*Waveform) error {
	return func(w *Waveform) error {
		return w.setFunction(function)
	}
}

// SetFunction applies the input function to the receiving Waveform struct.
func (w *Waveform) SetFunction(function SampleReduceFunc) error {
	return w.SetOptions(Function(function))
}

// setFunction directly sets the function member of the receiving Waveform
// struct.
func (w *Waveform) setFunction(function SampleReduceFunc) error {
	// Function cannot be nil
	if function == nil {
		return &OptionsError{
			Option: "function",
			Reason: "function cannot be nil",
		}
	}

	w.function = function

	return nil
}

// Resolution generates an OptionsFunc which applies the input resolution
// value to an input Waveform struct.
//
// This value indicates the number of times audio is read and drawn
// as a waveform, per second of audio.
func Resolution(resolution uint) func(*Waveform) error {
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
		return &OptionsError{
			Option: "resolution",
			Reason: "resolution cannot be 0",
		}
	}

	w.resolution = resolution

	return nil
}

// Scale generates an OptionsFunc which applies the input X and Y axis scaling
// factors to an input Waveform struct.
//
// This value indicates how a generated waveform image will be scaled, for both
// its X and Y axes.
func Scale(x uint, y uint) func(*Waveform) error {
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
		return &OptionsError{
			Option: "scale",
			Reason: "X scale cannot be 0",
		}
	}

	// Y scale cannot be zero
	if y == 0 {
		return &OptionsError{
			Option: "scale",
			Reason: "Y scale cannot be 0",
		}
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
func ScaleClipping() func(*Waveform) error {
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
func Sharpness(sharpness uint) func(*Waveform) error {
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
