// Package waveform is capable of generating waveform images from audio streams.  MIT Licensed.
package waveform

import (
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
	// yDefault is the default height of the generated waveform image
	yDefault = 128

	// rmsScaleDefault is the default scaling factor used when scaling RMS value and waveform height
	// by the output image's height
	rmsScaleDefault = 3.00
)

// Options are used to customize properties about a waveform image.
type Options struct {
	// BackgroundColor and ForegroundColor specify the background and foreground
	// color of a waveform image, respectively.
	BackgroundColor color.Color
	ForegroundColor color.Color

	// ScaleX and ScaleY are scaling factors used to scale a waveform image on its
	// X or Y axis, respectively.
	ScaleX int
	ScaleY int

	// Sharpness is used to apply a curve to a waveform image, scaled on its X-axis.
	// A higher value results in steeper curves, and a lower value results in more
	// "blocky" curves.
	Sharpness int

	// ScaleRMS specifies if the waveform image should be scaled down on its Y-axis
	// when certain RMS thresholds are reached.  This can be used to show a more accurate
	// waveform when a song reaches very high RMS thresholds.
	ScaleRMS bool
}

// DefaultOptions is a set of sane defaults, which are applied when no options are
// passed to New.
var DefaultOptions = &Options{
	// Black waveform on white background
	BackgroundColor: color.White,
	ForegroundColor: color.Black,

	// No scaling
	ScaleX: 1,
	ScaleY: 1,

	// Normal sharpness
	Sharpness: 1,

	// Do not scale high RMS values
	ScaleRMS: false,
}

// New creates a new image.Image from a io.Reader, with optional an optional Options struct
// which may be set.  New reads the input io.Reader, processes its input into a waveform,
// and returns the resulting image.  On failure, New will return any errors which occur.
func New(r io.Reader, options *Options) (image.Image, error) {
	// If options are nil, set sane defaults
	if options == nil {
		options = DefaultOptions
	}

	// Open audio decoder on input stream
	decoder, _, err := audio.NewDecoder(r)
	if err != nil {
		return nil, err
	}

	// rms is a slice of computed RMS values from each second of audio samples
	rms := make([]float64, 0)

	// Track the maximum RMS value computed, optionally used for scaling later
	var maxRMS float64

	// samples is a slice of float64 audio samples, used to store decoded values
	config := decoder.Config()
	samples := make(audio.F64Samples, config.SampleRate*config.Channels)
	for {
		// Decode one second of audio
		if _, err := decoder.Read(samples); err != nil {
			// On end of stream, stop reading values
			if err == audio.EOS {
				break
			}

			// On all other errors, return
			return nil, err
		}

		// Calculate RMS from float64 audio samples
		rmsSample := rmsF64Samples(samples)

		// Track the highest RMS value
		if rmsSample > maxRMS {
			maxRMS = rmsSample
		}

		// Store computed value
		rms = append(rms, rmsSample)
	}

	// Set image resolution
	imgX := len(rms) * options.ScaleX
	imgY := yDefault * options.ScaleY

	// Create output image, fill image with white background
	img := image.NewRGBA(image.Rect(0, 0, imgX, imgY))
	draw.Draw(img, img.Bounds(), image.NewUniform(options.BackgroundColor), image.ZP, draw.Src)

	// Calculate halfway point of Y-axis for image
	imgHalfY := img.Bounds().Max.Y / 2

	// Calculate a peak value used for smoothing scaled X-axis images
	peak := int(math.Ceil(float64(options.ScaleX)) / 2)

	// Calculate RMS scaling factor, based upon maximum RMS value found
	// If option ScaleRMS is true, when maximum value is above certain thresholds
	// the scaling factor is reduced to show an accurate waveform with less clipping
	rmsScale := rmsScaleDefault
	if options.ScaleRMS {
		if maxRMS > 0.35 {
			rmsScale -= 0.5
		}
		if maxRMS > 0.40 {
			rmsScale -= 0.25
		}
	}

	// Begin iterating all gathered RMS values
	x := 0
	for _, r := range rms {
		// Scale RMS value to an integer, using the height of the image and a constant
		// scaling factor
		scaleRMS := int(math.Floor(r * float64(img.Bounds().Max.Y) * rmsScale))

		// Calculate the halfway point for the scaled RMS value
		halfScaleRMS := scaleRMS / 2

		// Iterate image coordinates on the Y-axis, generating a symmetrical waveform
		// image above and below the center of the image
		for y := imgHalfY - halfScaleRMS; y < scaleRMS+(imgHalfY-halfScaleRMS); y++ {
			// If X-axis is being scaled, draw RMS value over several X coordinates
			for i := 0; i < options.ScaleX; i++ {
				// When scaled, adjust RMS value to be lower on either side of the peak,
				// so that the image appears more smooth and less "blocky"
				var adjust int
				if i < peak {
					// Adjust downward
					adjust = (i - peak) * options.Sharpness
				} else if i == peak {
					// No adjustment at peak
					adjust = 0
				} else {
					// Adjust downward
					adjust = (peak - i) * options.Sharpness
				}

				// On top half of the image, invert adjustment to create symmetry between
				// top and bottom halves
				if y < imgHalfY {
					adjust = -1 * adjust
				}

				// Draw using specified color at specified X and Y coordinate
				img.Set(x+i, y+adjust, options.ForegroundColor)
			}
		}

		// Increase X by scaling factor, to continue drawing at next loop
		x += options.ScaleX
	}

	// Return generated image
	return img, nil
}

// rmsF64Samples calculates the root mean square of a slice of float64 audio samples,
// enabling the measurement of magnitude over the entire set of samples.
// Derived from: http://en.wikipedia.org/wiki/Root_mean_square
func rmsF64Samples(samples audio.F64Samples) float64 {
	// Square and sum all input samples
	var sumSquare float64
	for i := range samples {
		sumSquare += math.Pow(float64(samples.At(i)), 2)
	}

	// Multiply squared sum by (1/n) coefficient, return square root
	return math.Sqrt(float64((float64(1) / float64(samples.Len()))) * sumSquare)
}
