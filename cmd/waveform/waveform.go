package main

import (
	"errors"
	"flag"
	"image/color"
	"image/png"
	"log"
	"os"
	"strconv"

	"github.com/mdlayher/waveform"
)

const (
	// app is the name of this application
	app = "waveform"
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

	// strBGColor is the hex color value used to color the background of the waveform image
	strBGColor = flag.String("bg-color", "#FFFFFF", "hex background color of output waveform image")

	// strFGColor is the hex color value used to color the foreground of the waveform image
	strFGColor = flag.String("fg-color", "#000000", "hex foreground color of output waveform image")

	// strAltColor is the hex color value used to set the alternate color of the waveform image
	strAltColor = flag.String("alt-color", "", "hex alternate color of output waveform image")

	// resolution is the number of times audio is read and the waveform is drawn,
	// per second of audio
	resolution = flag.Int("resolution", 1, "number of times audio is read and drawn per second of audio")

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

	log.SetPrefix(app + ": ")

	// Open input audio file, exit if it is not valid
	audioFile, err := os.Open(*inFilename)
	if err != nil {
		log.Fatal(err)
	}
	defer audioFile.Close()

	log.Printf("audio: %s", audioFile.Name())

	// Create image background color from input hex color string, or default
	// to black if invalid
	colorR, colorG, colorB := hexToRGB(*strBGColor)
	bgColor := color.RGBA{colorR, colorG, colorB, 255}

	// Create image foreground color from input hex color string, or default
	// to black if invalid
	colorR, colorG, colorB = hexToRGB(*strFGColor)
	fgColor := color.RGBA{colorR, colorG, colorB, 255}

	// Create image alternate color from input hex color string, or default
	// to foreground color if empty
	altColor := fgColor
	if *strAltColor != "" {
		colorR, colorG, colorB = hexToRGB(*strAltColor)
		altColor = color.RGBA{colorR, colorG, colorB, 255}
	}

	// Generate a waveform image from the input file, using values passed from
	// flags as options
	img, err := waveform.New(audioFile, &waveform.Options{
		BackgroundColor: bgColor,
		ForegroundColor: fgColor,
		AlternateColor:  altColor,

		Resolution: *resolution,

		ScaleX: *scaleX,
		ScaleY: *scaleY,

		Sharpness: *sharpness,

		ScaleRMS: true,
	})
	if err != nil {
		// Set of known errors
		knownErr := map[error]struct{}{
			waveform.ErrFormat:        struct{}{},
			waveform.ErrInvalidData:   struct{}{},
			waveform.ErrUnexpectedEOS: struct{}{},
		}

		// On known error, fatal log
		if _, ok := knownErr[err]; ok {
			log.Fatal(err)
		}

		// Unknown errors, panic
		panic(err)
	}

	// Attempt to create output image file
	imageFile, err := os.Create(*outFilename)
	if err != nil {
		panic(err)
	}
	defer imageFile.Close()

	// Encode results into output file
	log.Printf("image: %s", imageFile.Name())
	if err := png.Encode(imageFile, img); err != nil {
		panic(err)
	}
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
