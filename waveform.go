// Package waveform is capable of generating waveform images from audio streams.  MIT Licensed.
package waveform

import (
	"image"
	"image/color"
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

// Waveform is a struct which can be manipulated and used to generate
// audio waveform images from an input audio stream.
type Waveform struct {
	r io.Reader

	resolution uint
	sampleFn   SampleReduceFunc

	bgColorFn ColorFunc
	fgColorFn ColorFunc

	scaleX uint
	scaleY uint

	sharpness uint

	scaleClipping bool
}

// Generate immediately opens and reads an input audio stream, computes
// the values required for waveform generation, and returns a waveform image
// which is customized by zero or more, variadic, OptionsFunc parameters.
//
// Generate is equivalent to calling New, followed by the Compute and Draw
// methods of a Waveform struct.  In general, Generate should only be used
// for one-time waveform image generation.
func Generate(r io.Reader, options ...OptionsFunc) (image.Image, error) {
	w, err := New(r, options...)
	if err != nil {
		return nil, err
	}

	values, err := w.Compute()
	return w.Draw(values), err
}

// New generates a new Waveform struct, applying any input OptionsFunc
// on return.
func New(r io.Reader, options ...OptionsFunc) (*Waveform, error) {
	// Generate Waveform struct with sane defaults
	w := &Waveform{
		// Read from input stream
		r: r,

		// Read audio and compute values once per second of audio
		resolution: 1,

		// Use RMSF64Samples as a SampleReduceFunc
		sampleFn: RMSF64Samples,

		// Generate solid, black background color with solid, white
		// foreground color waveform using ColorFunc
		bgColorFn: SolidColor(color.White),
		fgColorFn: SolidColor(color.Black),

		// No scaling
		scaleX: 1,
		scaleY: 1,

		// Normal sharpness
		sharpness: 1,

		// Do not scale clipping values
		scaleClipping: false,
	}

	// Apply any input OptionsFunc on return
	return w, w.SetOptions(options...)
}

// Compute creates a slice of float64 values, computed using an input function.
//
// Compute is typically used once on an audio stream, to read and calculate the values
// used for subsequent waveform generations.  Its return value can be used with Draw to
// generate and customize multiple waveform images from a single stream.
func (w *Waveform) Compute() ([]float64, error) {
	return w.readAndComputeSamples()
}

// Draw creates a new image.Image from a slice of float64 values.
//
// Draw is typically used after a waveform has been computed one time, and a slice
// of computed values was returned from the first computation.  Subsequent calls to
// Draw may be used to customize a waveform using the same input values.
func (w *Waveform) Draw(values []float64) image.Image {
	return w.generateImage(values)
}

// readAndComputeSamples opens the input audio stream, computes samples according
// to an input function, and returns a slice of computed values and any errors
// which occurred during the computation.
func (w *Waveform) readAndComputeSamples() ([]float64, error) {
	// Validate struct members
	// These checks are also done when applying options, but verifying them here
	// will prevent a runtime panic if called on an empty Waveform instance.
	if w.sampleFn == nil {
		return nil, errSampleFunctionNil
	}
	if w.resolution == 0 {
		return nil, errResolutionZero
	}

	// Open audio decoder on input stream
	decoder, _, err := audio.NewDecoder(w.r)
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
	samples := make(audio.F64Samples, uint(config.SampleRate*config.Channels)/w.resolution)
	for {
		// Decode at specified resolution from options
		// On any error other than end-of-stream, return
		_, err := decoder.Read(samples)
		if err != nil && err != audio.EOS {
			return nil, err
		}

		// Apply SampleReduceFunc over float64 audio samples
		value = w.sampleFn(samples)

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

// generateImage takes a slice of computed values and generates
// a waveform image from the input.
func (w *Waveform) generateImage(computed []float64) image.Image {
	// Store integer scale values
	intScaleX := int(w.scaleX)
	intScaleY := int(w.scaleY)

	// Set image resolution
	imgX := len(computed) * intScaleX
	imgY := imgYDefault * intScaleY

	// Create output, rectangular image
	img := image.NewRGBA(image.Rect(0, 0, imgX, imgY))
	bounds := img.Bounds()

	// Calculate halfway point of Y-axis for image
	imgHalfY := bounds.Max.Y / 2

	// Calculate a peak value used for smoothing scaled X-axis images
	peak := int(math.Ceil(float64(w.scaleX)) / 2)

	// Calculate scaling factor, based upon maximum value computed by a SampleReduceFunc.
	// If option ScaleClipping is true, when maximum value is above certain thresholds
	// the scaling factor is reduced to show an accurate waveform with less clipping.
	imgScale := scaleDefault
	if w.scaleClipping {
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
	intBoundY := int(bounds.Max.Y)
	f64BoundY := float64(bounds.Max.Y)
	intSharpness := int(w.sharpness)

	// Begin iterating all computed values
	x := 0
	for count, c := range computed {
		// Scale computed value to an integer, using the height of the image and a constant
		// scaling factor
		scaleComputed = int(math.Floor(c * f64BoundY * imgScale))

		// Calculate the halfway point for the scaled computed value
		halfScaleComputed = scaleComputed / 2

		// Draw background color down the entire Y-axis
		for y := 0; y < intBoundY; y++ {
			// If X-axis is being scaled, draw background over several X coordinates
			for i := 0; i < intScaleX; i++ {
				img.Set(x+i, y, w.bgColorFn(count, x+i, y))
			}
		}

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

				// Retrieve and apply color function at specified computed value
				// count, and X and Y coordinates.
				// The output color is selected using the function, and is applied to
				// the resulting image.
				img.Set(x+i, y+adjust, w.fgColorFn(count, x+i, y+adjust))
			}
		}

		// Increase X by scaling factor, to continue drawing at next loop
		x += intScaleX
	}

	// Return generated image
	return img
}
