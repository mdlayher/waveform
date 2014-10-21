// Package waveform is capable of generating waveform images from audio streams.  MIT Licensed.
package waveform

import (
	"errors"
	"image"
	"image/color"
	"image/draw"
	"io"
	"math"

	"azul3d.org/audio.v1"

	// Import WAV and FLAC decoders
	_ "azul3d.org/audio/wav.v1"
	_ "github.com/azul3d/audio-flac"
)

const (
	// imgYDefault is the default height of the generated waveform image
	imgYDefault = 128

	// scaleDefault is the default scaling factor used when scaling computed
	// value and waveform height by the output image's height
	scaleDefault = 3.00
)

// Error values from azul3d/audio.v1 are wrapped, so that callers do not have to
// import an additional package to check for common errors.
var (
	// ErrFormat is returned when the input audio format is not a registered format
	// with the audio package.
	ErrFormat = struct{ error }{audio.ErrFormat}

	// ErrInvalidData is returned when the input audio format is recognized, but
	// the stream is invalid or corrupt in some way.
	ErrInvalidData = struct{ error }{audio.ErrInvalidData}

	// ErrUnexpectedEOS is returned when end-of-stream is encountered in the middle
	// of a fixed-size block or data structure.
	ErrUnexpectedEOS = struct{ error }{audio.ErrUnexpectedEOS}
)

var (
	// errNilFunction is returned when an input SampleReduceFunc is nil.
	errNilFunction = errors.New("waveform: nil SampleReduceFunc")

	// errZeroResolution is returned when input Resolution is 0.
	errZeroResolution = errors.New("waveform: zero Resolution")
)

// Options are used to customize properties about a waveform image.  It embeds both
// the ComputeOptions and ImageOptions types.  It is used to customize the New operation.
type Options struct {
	ComputeOptions
	ImageOptions
}

// ComputeOptions are used to customize the ComputeValues operation.
type ComputeOptions struct {
	// Resolution sets the number of times audio is read and drawn
	// as a waveform, per second of audio.
	Resolution uint

	// Function is used to specify an alternate SampleReduceFunc for use in waveform
	// generation.  The function is applied over a slice of float64 audio samples,
	// reducing them to a single value.
	Function SampleReduceFunc
}

// ImageOptions are used to customize the DrawImage operation.
type ImageOptions struct {
	// BackgroundColor and ForegroundColor specify the background and foreground
	// color of a waveform image, respectively.
	// AlternateColor specifies an optional secondary color which is alternated with
	// the foreground color to create a stripe effect in the image.  If not specified,
	// no alternate color will be used.
	BackgroundColor color.Color
	ForegroundColor color.Color
	AlternateColor  color.Color

	// ScaleX and ScaleY are scaling factors used to scale a waveform image on its
	// X or Y axis, respectively.
	ScaleX uint
	ScaleY uint

	// Sharpness is used to apply a curve to a waveform image, scaled on its X-axis.
	// A higher value results in steeper curves, and a lower value results in more
	// "blocky" curves.
	Sharpness uint

	// ScaleClipping specifies if the waveform image should be scaled down on its
	// Y-axis when clipping thresholds are reached.  This can be used to show a
	// more accurate waveform, when a waveform exhibits signs of audio clipping.
	ScaleClipping bool
}

// DefaultOptions is a set of sane defaults, which are applied when nil Options are
// passed to any function which accepts Options, ComputeOptions, or ImageOptions.
var DefaultOptions = Options{
	// Options specific to ComputeValues
	ComputeOptions{
		// Read audio and compute values once per second of audio
		Resolution: 1,

		// Use RMSF64Samples as a SampleReduceFunc
		Function: RMSF64Samples,
	},

	// Options specific to DrawImage
	ImageOptions{
		// Black waveform on white background
		// No alternate color
		BackgroundColor: color.White,
		ForegroundColor: color.Black,
		AlternateColor:  nil,

		// No scaling
		ScaleX: 1,
		ScaleY: 1,

		// Normal sharpness
		Sharpness: 1,

		// Do not scale clipping values
		ScaleClipping: false,
	},
}

// New creates a new image.Image from a io.Reader.  An Options struct may be passed
// to enable further customization; else, DefaultOptions is used.
//
// New is typically used for one-time-only waveform generation.  If the
// same audio stream will be processed multiple times during a single program run, it
// is better to use a combination of ComputeValues and DrawImage, instead of New.
//
// New is equivalent to calling ComputeValues, and passing the return value to DrawImage.
func New(r io.Reader, options *Options) (image.Image, error) {
	// Perform validation and corrections on options
	var opt Options
	if options == nil {
		opt = DefaultOptions
	} else {
		opt = validateOptions(*options)
	}

	// Compute values from input stream
	computed, err := readAndComputeSamples(r, opt.ComputeOptions)
	if err != nil {
		return nil, err
	}

	// Generate and return output image
	return generateImage(computed, opt.ImageOptions), nil
}

// ComputeValues creates a slice of float64 values, computed using an input function.  A ComputeOptions
// struct may be passed to enable further customization; else, DefaultOptions.ComputeOptions is used.
//
// ComputeValues is typically used once on an audio stream, to read and calculate the values
// used for subsequent waveform generations.  Its return value can be used with DrawImage to
// generate and customize multiple waveform images from a single stream.
func ComputeValues(r io.Reader, options *ComputeOptions) ([]float64, error) {
	// Perform validation and corrections on options
	var opt Options
	if options == nil {
		opt = DefaultOptions
	} else {
		opt = validateOptions(Options{
			ComputeOptions: *options,
		})
	}

	return readAndComputeSamples(r, opt.ComputeOptions)
}

// DrawImage creates a new image.Image from a slice of float64 values.  An ImageOptions struct
// may be passed to enable further customization; else, DefaultOptions.ImageOptions is used.
//
// DrawImage is typically used after a waveform has been computed one time, and a slice of
// computed values was returned from the first computation.  Subsequent waveform generations from
// the same audio stream can be generated and customized much more quickly using this function,
// rather than calling New again.
func DrawImage(values []float64, options *ImageOptions) image.Image {
	// Perform validation and corrections on options
	var opt Options
	if options == nil {
		opt = DefaultOptions
	} else {
		opt = validateOptions(Options{
			ImageOptions: *options,
		})
	}

	// Generate and return output image
	return generateImage(values, opt.ImageOptions)
}

// SampleReduceFunc is a function which reduces a set of float64 audio samples
// into a single float64 value.
type SampleReduceFunc func(samples audio.F64Samples) float64

// RMSF64Samples is a SampleReduceFunc which calculates the root mean square
// of a slice of float64 audio samples, enabling the measurement of magnitude
// over the entire set of samples.
//
// Derived from: http://en.wikipedia.org/wiki/Root_mean_square.
func RMSF64Samples(samples audio.F64Samples) float64 {
	// Square and sum all input samples
	var sumSquare float64
	for i := range samples {
		sumSquare += math.Pow(float64(samples.At(i)), 2)
	}

	// Multiply squared sum by (1/n) coefficient, return square root
	return math.Sqrt(sumSquare / float64(samples.Len()))
}

// readAndComputeSamples opens an input audio stream, computes samples according
// to an input function from options, and returns a slice of computed values and
// any errors which occurred.
func readAndComputeSamples(r io.Reader, options ComputeOptions) ([]float64, error) {
	// Validate input
	if options.Function == nil {
		return nil, errNilFunction
	}
	if options.Resolution == 0 {
		return nil, errZeroResolution
	}

	// Open audio decoder on input stream
	decoder, _, err := audio.NewDecoder(r)
	if err != nil {
		// Unknown format
		if err == audio.ErrFormat {
			return nil, ErrFormat
		}

		// Invalid data
		if err == audio.ErrInvalidData {
			return nil, ErrInvalidData
		}

		// Unexpected end-of-stream
		if err == audio.ErrUnexpectedEOS {
			return nil, ErrUnexpectedEOS
		}

		// All other errors
		return nil, err
	}

	// computed is a slice of computed values by a SampleReduceFunc, from each
	// slice of audio samples
	var computed []float64

	// Track the current computed value
	var value float64

	// samples is a slice of float64 audio samples, used to store decoded values
	config := decoder.Config()
	samples := make(audio.F64Samples, uint(config.SampleRate*config.Channels)/options.Resolution)
	for {
		// Decode at specified resolution from options
		// On any error other than end-of-stream, return
		_, err := decoder.Read(samples)
		if err != nil && err != audio.EOS {
			return nil, err
		}

		// Apply SampleReduceFunc over float64 audio samples
		value = options.Function(samples)

		// Store computed value
		computed = append(computed, value)

		// On end of stream, stop reading values
		if err == audio.EOS {
			break
		}
	}

	// Return slice of computed values
	return computed, nil
}

// generateImage takes a slice of computed values and options, and generates
// a waveform image from the input.  Options are applied as set by the caller.
func generateImage(computed []float64, options ImageOptions) image.Image {
	// Store integer scale values
	intScaleX := int(options.ScaleX)
	intScaleY := int(options.ScaleY)

	// Set image resolution
	imgX := len(computed) * intScaleX
	imgY := imgYDefault * intScaleY

	// Create output image, fill image with specified background color
	img := image.NewRGBA(image.Rect(0, 0, imgX, imgY))
	draw.Draw(img, img.Bounds(), image.NewUniform(options.BackgroundColor), image.ZP, draw.Src)

	// Calculate halfway point of Y-axis for image
	imgHalfY := img.Bounds().Max.Y / 2

	// Calculate a peak value used for smoothing scaled X-axis images
	peak := int(math.Ceil(float64(options.ScaleX)) / 2)

	// Calculate scaling factor, based upon maximum value computed by a SampleReduceFunc.
	// If option ScaleClipping is true, when maximum value is above certain thresholds
	// the scaling factor is reduced to show an accurate waveform with less clipping.
	imgScale := scaleDefault
	if options.ScaleClipping {
		// Find maximum value from input slice
		var maxValue float64
		for _, c := range computed {
			if c > maxValue {
				maxValue = c
			}
		}

		// For each 0.05 maximum increment at 0.30 and above, reduce the scaling
		// factor by 0.25.  This is a rough estimate and may be tweaked in the future.
		for i := 0.30; i < maxValue; i += 0.05 {
			imgScale -= 0.25
		}
	}

	// Values to be used for repeated computations
	var scaleComputed, halfScaleComputed, adjust int
	f64BoundY := float64(img.Bounds().Max.Y)
	intSharpness := int(options.Sharpness)

	// Begin iterating all computed values
	x := 0
	for count, c := range computed {
		// Scale computed value to an integer, using the height of the image and a constant
		// scaling factor
		scaleComputed = int(math.Floor(c * f64BoundY * imgScale))

		// Calculate the halfway point for the scaled computed value
		halfScaleComputed = scaleComputed / 2

		// Iterate image coordinates on the Y-axis, generating a symmetrical waveform
		// image above and below the center of the image
		for y := imgHalfY - halfScaleComputed; y < scaleComputed+(imgHalfY-halfScaleComputed); y++ {
			// If X-axis is being scaled, draw computed value over several X coordinates
			for i := 0; i < intScaleX; i++ {
				// When scaled, adjust computed value to be lower on either side of the peak,
				// so that the image appears more smooth and less "blocky"
				if i < peak {
					// Adjust downward
					adjust = (i - peak) * intSharpness
				} else if i == peak {
					// No adjustment at peak
					adjust = 0
				} else {
					// Adjust downward
					adjust = (peak - i) * intSharpness
				}

				// On top half of the image, invert adjustment to create symmetry between
				// top and bottom halves
				if y < imgHalfY {
					adjust = -1 * adjust
				}

				// On odd iterations (or if no alternate set), draw using specified
				// foreground color at specified X and Y coordinate
				if count%2 != 0 || options.AlternateColor == nil {
					img.Set(x+i, y+adjust, options.ForegroundColor)
				} else {
					// On even iterations, draw using specified alternate color at
					// specified X and Y coordinate
					img.Set(x+i, y+adjust, options.AlternateColor)
				}
			}
		}

		// Increase X by scaling factor, to continue drawing at next loop
		x += intScaleX
	}

	// Return generated image
	return img
}

// validateOptions verifies that an input Options struct is correct, and
// sets sane defaults for fields which are not specified
func validateOptions(options Options) Options {
	// If resolution is 0, set it to default to avoid divide-by-zero panic
	if options.Resolution == 0 {
		options.Resolution = DefaultOptions.Resolution
	}

	// If either scale is 0, set to default to avoid empty image
	if options.ScaleX == 0 {
		options.ScaleX = DefaultOptions.ScaleX
	}
	if options.ScaleY == 0 {
		options.ScaleY = DefaultOptions.ScaleY
	}

	// If color options are nil, set sane defaults to prevent panic
	if options.BackgroundColor == nil {
		options.BackgroundColor = DefaultOptions.BackgroundColor
	}
	if options.ForegroundColor == nil {
		options.ForegroundColor = DefaultOptions.ForegroundColor
	}
	if options.AlternateColor == nil {
		options.AlternateColor = DefaultOptions.AlternateColor
	}

	// If no SampleReduceFunc is specified, use default
	if options.Function == nil {
		options.Function = DefaultOptions.Function
	}

	return options
}
