// Command waveform is a simple utility which reads an audio file from stdin,
// processes it into a waveform image using input flags, and writes a PNG image
// of the generated waveform to stdout.
package main

import (
	"flag"
	"fmt"
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

	// Names of available color functions
	fnFuzz     = "fuzz"
	fnGradient = "gradient"
	fnSolid    = "solid"
	fnStripe   = "stripe"
)

var (
	// strBGColor is the hex color value used to color the background of the waveform image
	strBGColor = flag.String("bg", "#FFFFFF", "hex background color of output waveform image")

	// strFGColor is the hex color value used to color the foreground of the waveform image
	strFGColor = flag.String("fg", "#000000", "hex foreground color of output waveform image")

	// strAltColor is the hex color value used to set the alternate color of the waveform image
	strAltColor = flag.String("alt", "", "hex alternate color of output waveform image")

	// resolution is the number of times audio is read and the waveform is drawn,
	// per second of audio
	resolution = flag.Uint("resolution", 1, "number of times audio is read and drawn per second of audio")

	// scaleX is the scaling factor for the output waveform file's X-axis
	scaleX = flag.Uint("x", 1, "scaling factor for image X-axis")

	// scaleY is the scaling factor for the output waveform file's Y-axis
	scaleY = flag.Uint("y", 1, "scaling factor for image Y-axis")

	// sharpness is the factor used to add curvature to a scaled image, preventing
	// "blocky" images at higher scaling
	sharpness = flag.Uint("sharpness", 1, "sharpening factor used to add curvature to a scaled image")

	// strFn is an identifier which selects the ColorFunc used to color the waveform image
	strFn = flag.String("fn", fnSolid, "function used to color output waveform image "+fnOptions)
)

// fnOptions is the help string which lists available options
var fnOptions = fmt.Sprintf("[options: %s, %s, %s, %s]", fnFuzz, fnGradient, fnSolid, fnStripe)

func main() {
	// Parse flags
	flag.Parse()

	// Move all logging output to stderr, as output image will occupy
	// the stdout stream
	log.SetOutput(os.Stderr)
	log.SetPrefix(app + ": ")

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

	// Set of available functions
	fnSet := map[string]waveform.ColorFunc{
		fnFuzz:     waveform.FuzzColor(fgColor, altColor),
		fnGradient: waveform.GradientColor(fgColor, altColor),
		fnSolid:    waveform.SolidColor(fgColor),
		fnStripe:   waveform.StripeColor(fgColor, altColor),
	}

	// Validate user-selected function
	colorFn, ok := fnSet[*strFn]
	if !ok {
		log.Fatalf("unknown function: %q %s", *strFn, fnOptions)
	}

	// Generate a waveform image from stdin, using values passed from
	// flags as options
	img, err := waveform.Generate(os.Stdin,
		waveform.BGColorFunction(waveform.SolidColor(bgColor)),
		waveform.FGColorFunction(colorFn),
		waveform.Resolution(*resolution),
		waveform.Scale(*scaleX, *scaleY),
		waveform.ScaleClipping(),
		waveform.Sharpness(*sharpness),
	)
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

	// Encode results as PNG to stdout
	if err := png.Encode(os.Stdout, img); err != nil {
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
