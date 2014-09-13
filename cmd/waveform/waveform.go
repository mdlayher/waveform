package main

import (
	"errors"
	"flag"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"log"
	"math"
	"os"
	"strconv"

	"azul3d.org/audio.v1"

	// Import WAV and FLAC decoders
	_ "azul3d.org/audio/wav.v1"
	_ "github.com/azul3d/audio-flac"
)

const (
	// app is the name of this application
	app = "waveform"

	// yDefault is the default height of the generated waveform image
	yDefault = 128

	// rmsScaleDefault is the default scaling factor used when scaling RMS value and waveform height
	// by the output image's height
	rmsScaleDefault = 3.00
)

var (
	// ErrMissingParameters is returned when required input and output filenames are not
	// passed via command-line flags.
	ErrMissingParameters = errors.New(app + ": missing required parameters: -in, -out")
)

var (
	// inFilename is the file name of the input audio file
	inFilename = flag.String("in", "", "input audio file")

	// inFilename is the file name of the output waveform PNG image file
	outFilename = flag.String("out", "", "output PNG waveform image file")

	// strColor is the hex color value used to color the waveform image
	strColor = flag.String("color", "#000000", "hex color of output waveform image")

	// scaleX is the scaling factor for the output waveform file's X-axis
	scaleX = flag.Int("x", 1, "scaling factor for image X-axis")

	// scaleY is the scaling factor for the output waveform file's Y-axis
	scaleY = flag.Int("y", 1, "scaling factor for image Y-axis")

	// sharpness is the factor used to add curvature to a scaled image, preventing
	// "blocky" images at higher scaling
	sharpness = flag.Int("sharpness", 1, "sharpening factor used to add curvature to a scaled image")
)

func main() {
	// Parse flags and check for required parameters
	flag.Parse()
	if *inFilename == "" || *outFilename == "" {
		log.Fatal(ErrMissingParameters)
	}

	// Open input audio file, exit if it is not valid
	audioFile, err := os.Open(*inFilename)
	if err != nil {
		log.Fatal(err)
	}

	// Open audio decoder, exit if decoder does not recognize input
	decoder, format, err := audio.NewDecoder(audioFile)
	if err != nil {
		log.Fatal(err)
	}

	// Log information regarding the input audio file
	log.SetPrefix(app + ": ")
	config := decoder.Config()
	log.Printf(" audio: %s [%s, %dHz, %dch]", audioFile.Name(), format, config.SampleRate, config.Channels)

	// rms is a slice of computed RMS values from each second of audio samples
	rms := make([]float64, 0)

	// Track the maximum RMS value computed, used for scaling later
	var maxRMS float64

	// samples is a slice of float64 audio samples, used to store decoded values
	samples := make(audio.F64Samples, config.SampleRate*config.Channels)
	for {
		// Decode one second of audio
		if _, err := decoder.Read(samples); err != nil {
			// On end of stream, stop reading values
			if err == audio.EOS {
				break
			}

			// On all other errors, panic
			panic(err)
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

	// Close audio file
	if err := audioFile.Close(); err != nil {
		panic(err)
	}

	// Set image resolution
	imgX := len(rms) * (*scaleX)
	imgY := yDefault * (*scaleY)
	log.Printf(" scale: [%dx%d]: x * %d, y * %d", imgX, imgY, *scaleX, *scaleY)

	// Create output image, fill image with white background
	img := image.NewRGBA(image.Rect(0, 0, imgX, imgY))
	draw.Draw(img, img.Bounds(), image.White, image.ZP, draw.Src)

	// Create image color from input hex color string, or default
	// to black if invalid
	colorR, colorG, colorB := hexToRGB(*strColor)
	waveformColor := color.RGBA{colorR, colorG, colorB, 255}

	// Calculate halfway point of Y-axis for image
	imgHalfY := img.Bounds().Max.Y / 2

	// Calculate a peak value used for smoothing scaled X-axis images
	peak := int(math.Ceil(float64(*scaleX)) / 2)

	// Calculate RMS scaling factor, based upon maximum RMS value found
	// If maximum value is above certain thresholds, the scaling factor is reduced
	// to show an accurate waveform with less clipping
	rmsScale := rmsScaleDefault
	if maxRMS > 0.35 {
		rmsScale -= 0.5
	}
	if maxRMS > 0.40 {
		rmsScale -= 0.25
	}
	log.Printf("maxRMS: %0.03f [scale: %0.03f]", maxRMS, rmsScale)

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
			for i := 0; i < *scaleX; i++ {
				// When scaled, adjust RMS value to be lower on either side of the peak,
				// so that the image appears more smooth and less "blocky"
				var adjust int
				if i < peak {
					// Adjust downward
					adjust = (i - peak) * (*sharpness)
				} else if i == peak {
					// No adjustment at peak
					adjust = 0
				} else {
					// Adjust downward
					adjust = (peak - i) * (*sharpness)
				}

				// On top half of the image, invert adjustment to create symmetry between
				// top and bottom halves
				if y < imgHalfY {
					adjust = -1 * adjust
				}

				// Draw using specified color at specified X and Y coordinate
				img.Set(x+i, y+adjust, waveformColor)
			}
		}

		// Increase X by scaling factor, to continue drawing at next loop
		x += *scaleX
	}

	// Attempt to create output image file
	imageFile, err := os.Create(*outFilename)
	if err != nil {
		panic(err)
	}
	defer imageFile.Close()

	// Encode results into output file
	log.Printf(" image: %s", imageFile.Name())
	if err := png.Encode(imageFile, img); err != nil {
		panic(err)
	}
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

// hexToRGB converts a hex string to a RGB triple.
// Credit: https://code.google.com/p/gorilla/source/browse/color/hex.go?r=ef489f63418265a7249b1d53bdc358b09a4a2ea0
func hexToRGB(h string) (uint8, uint8, uint8) {
	if len(h) > 0 && h[0] == '#' {
		h = h[1:]
	}
	if len(h) == 3 {
		h = h[:1] + h[:1] + h[1:2] + h[1:2] + h[2:] + h[2:]
	}
	if len(h) == 6 {
		if rgb, err := strconv.ParseUint(string(h), 16, 32); err == nil {
			return uint8(rgb >> 16), uint8((rgb >> 8) & 0xFF), uint8(rgb & 0xFF)
		}
	}
	return 0, 0, 0
}
