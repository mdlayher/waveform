package waveform

import (
	"bytes"
	"image/color"
	"io/ioutil"
	"log"
	"math"
	"reflect"
	"testing"

	"azul3d.org/audio.v1"
)

var (
	// Read in test files
	wavFile = func() []byte {
		file, err := ioutil.ReadFile("./test/tone16bit.wav")
		if err != nil {
			log.Fatalf("could not open test WAV: %v", err)
		}

		return file
	}()
	flacFile = func() []byte {
		file, err := ioutil.ReadFile("./test/tone16bit.flac")
		if err != nil {
			log.Fatalf("could not open test FLAC: %v", err)
		}

		return file
	}()
	mp3File = func() []byte {
		file, err := ioutil.ReadFile("./test/tone16bit.mp3")
		if err != nil {
			log.Fatalf("could not open test MP3: %v", err)
		}

		return file
	}()
	oggVorbisFile = func() []byte {
		file, err := ioutil.ReadFile("./test/tone16bit.ogg")
		if err != nil {
			log.Fatalf("could not open test Ogg Vorbis: %v", err)
		}

		return file
	}()
)

// TestNew verifies that New reads an input audio stream properly, and generates a
// correct output image.
func TestNew(t *testing.T) {
	// X is set by file duration, Y by library
	const defaultX = 5
	const defaultY = imgYDefault

	// Table of tests
	var tests = []struct {
		stream  []byte
		err     error
		options *Options
	}{
		// MP3 file
		{mp3File, ErrFormat, nil},
		// Ogg Vorbis file
		{oggVorbisFile, ErrFormat, nil},
		// Unknown format
		{[]byte("nonsense"), ErrFormat, nil},
		// WAV file, no options
		{wavFile, nil, nil},
		// WAV file, scaled
		{wavFile, nil, &Options{ComputeOptions{}, ImageOptions{ScaleX: 5, ScaleY: 2}}},
		// FLAC file, no options
		{flacFile, nil, nil},
		// FLAC file, scaled
		{flacFile, nil, &Options{ComputeOptions{}, ImageOptions{ScaleX: 5, ScaleY: 2}}},
	}

	// Iterate all tests
	for _, test := range tests {
		// Generate a io.Reader
		reader := bytes.NewReader(test.stream)

		// Attempt to create image for the reader
		img, err := New(reader, test.options)
		if err != nil {
			if err == test.err {
				continue
			}

			t.Fatalf("unexpected error: %v", err)
		}

		// Verify that image is RGBA
		if model := img.ColorModel(); model != color.RGBAModel {
			t.Fatalf("unexpected color model: %v != %v", model, color.RGBAModel)
		}

		// Check for expected bounds
		bounds := img.Bounds()
		var scaleX, scaleY uint
		if test.options == nil {
			scaleX = defaultX
			scaleY = defaultY
		} else {
			// Set scale by options
			scaleX = test.options.ScaleX * defaultX
			scaleY = test.options.ScaleY * defaultY
		}

		if uint(bounds.Max.X) != scaleX || uint(bounds.Max.Y) != scaleY {
			t.Fatalf("unexpected bounds: (%v,%v) != (%v,%v)", bounds.Max.X, bounds.Max.Y, scaleX, scaleY)
		}
	}
}

// Test_readAndComputeSamples verifies that readAndComputeSamples provides the correct output values
// for an input audio stream with the specified options
func Test_readAndComputeSamples(t *testing.T) {
	// Options for this test
	res0Opt := ComputeOptions{
		Resolution: 0,
		Function:   RMSF64Samples,
	}
	res2Opt := ComputeOptions{
		Resolution: 2,
		Function:   RMSF64Samples,
	}

	nilFnOpt := ComputeOptions{
		Resolution: 1,
		Function:   nil,
	}
	zeroFnOpt := ComputeOptions{
		Resolution: 1,
		Function: func(samples audio.F64Samples) float64 {
			return 0.00
		},
	}

	// Special values used for this test
	fiveZeroF64 := []float64{0.00, 0.00, 0.00, 0.00, 0.00}

	// Table of tests
	var tests = []struct {
		stream  []byte
		err     error
		options ComputeOptions
		values  []float64
	}{
		// MP3 file
		{mp3File, ErrFormat, DefaultOptions.ComputeOptions, nil},
		// Ogg Vorbis file
		{oggVorbisFile, ErrFormat, DefaultOptions.ComputeOptions, nil},
		// Unknown format
		{[]byte("nonsense"), ErrFormat, DefaultOptions.ComputeOptions, nil},
		// WAV file, standard options
		{wavFile, nil, DefaultOptions.ComputeOptions, []float64{0.7071166200538482, 0.7071166444294603, 0.7071166239921965, 0.7071165471800284, 0.7071166847631818}},
		// WAV file, zero resolution
		{wavFile, errZeroResolution, res0Opt, nil},
		// WAV file, double resolution
		{wavFile, nil, res2Opt, []float64{0.707116568673387, 0.7071166714342908, 0.7071166912513497, 0.7071165976075522, 0.7071165950470297, 0.7071166529373748, 0.7071165684016294, 0.7071165259584128, 0.7071167206200514, 0.7071166487801294}},
		// WAV file, nil function
		{wavFile, errNilFunction, nilFnOpt, nil},
		// WAV file, zero function
		{wavFile, nil, zeroFnOpt, fiveZeroF64},
		// FLAC file, standard options
		{flacFile, nil, DefaultOptions.ComputeOptions, []float64{0.7071166200538482, 0.7071166444294603, 0.7071166239921965, 0.7071165471800284, 0.7071166825227931}},
		// FLAC file, zero resolution
		{flacFile, errZeroResolution, res0Opt, nil},
		// FLAC file, double resolution
		{flacFile, nil, res2Opt, []float64{0.707116568673387, 0.7071166714342908, 0.7071166912513497, 0.7071165976075522, 0.7071165950470297, 0.7071166529373748, 0.7071165684016294, 0.7071165259584128, 0.7071167206200514, 0.7071166444255262}},
		// FLAC file, nil function
		{flacFile, errNilFunction, nilFnOpt, nil},
		// FLAC file, zero function
		{flacFile, nil, zeroFnOpt, fiveZeroF64},
	}

	// Iterate all tests
	for i, test := range tests {
		// Generate a io.Reader
		reader := bytes.NewReader(test.stream)

		// Attempt to read and compute value of samples
		values, err := readAndComputeSamples(reader, test.options)
		if err != nil {
			if err == test.err {
				continue
			}

			t.Fatalf("[%d] unexpected error: %v", i, err)
		}

		// Optionally compare values
		if test.values != nil {
			// Verify same length
			if len(values) != len(test.values) {
				t.Fatalf("[%d] unexpected length: %v != %v", len(values), len(test.values))
			}

			// Iterate and check values
			for j := 0; j < len(values); j++ {
				if values[j] != test.values[j] {
					t.Fatalf("[%d:%d] unexpected values: %v != %v", i, j, values[j], test.values[j])
				}
			}
		}
	}
}

// Test_generateImage verifies that generateImage creates an output image with expected
// characteristics, based upon input ImageOptions.
func Test_generateImage(t *testing.T) {
	// Values used for tests
	clippingValues := []float64{1.0, 1.0, 1.0, 1.0, 1.0}

	defaultBounds := "(0,0)-(5,128)"

	rgbaBlack := color.RGBA{0, 0, 0, 255}
	rgbaWhite := color.RGBA{255, 255, 255, 255}
	rgbaRed := color.RGBA{255, 0, 0, 255}

	// Test table to iterate and check conditions
	var tests = []struct {
		values   []float64
		options  ImageOptions
		bounds   string
		fgColor  color.Color
		altColor color.Color
	}{
		// Clipping values, no scaling
		{clippingValues, ImageOptions{}, defaultBounds, rgbaBlack, rgbaBlack},
		// Clipping values, add scaling
		{clippingValues, ImageOptions{ScaleClipping: true}, defaultBounds, rgbaWhite, rgbaWhite},
		// Clipping values, alternate color
		{clippingValues, ImageOptions{AlternateColor: rgbaRed}, defaultBounds, rgbaRed, rgbaBlack},
	}

	// Draw an image for each test
	for i, test := range tests {
		img := DrawImage(test.values, &test.options)

		// Verify bounds are valid
		if bounds := img.Bounds().String(); bounds != test.bounds {
			t.Fatalf("[%02d] unexpected bounds: %v != %v", i, bounds, test.bounds)
		}

		// Verify valid color model (always RGBA)
		if model := img.ColorModel(); model != color.RGBAModel {
			t.Fatalf("[%02d] unexpected color model: %v != %v", i, model, color.RGBAModel)
		}

		// Test for expected foreground color
		if expColor := img.At(0, 0); expColor != test.fgColor {
			t.Fatalf("[%02d] unexpected foreground color: %v != %v", i, expColor, test.fgColor)
		}

		// Test for expected alternate color (the actual foreground color in this case, because
		// the algorithm stripes even bands only)
		if expColor := img.At(1, 0); expColor != test.altColor {
			t.Fatalf("[%02d] unexpected alternate color: %v != %v", i, expColor, test.altColor)
		}
	}
}

// TestRMSF64Samples verifies that RMSF64Samples computes correct results
func TestRMSF64Samples(t *testing.T) {
	var tests = []struct {
		samples audio.F64Samples
		result  float64
		isNaN   bool
	}{
		// Empty samples - NaN
		{audio.F64Samples{}, 0.00, true},
		// Negative samples
		{audio.F64Samples{-0.10}, 0.10, false},
		{audio.F64Samples{-0.10, -0.20}, 0.15811388300841897, false},
		{audio.F64Samples{-0.10, -0.20, -0.30, -0.40, -0.50}, 0.33166247903554, false},
		// Positive samples
		{audio.F64Samples{0.10}, 0.10, false},
		{audio.F64Samples{0.10, 0.20}, 0.15811388300841897, false},
		{audio.F64Samples{0.10, 0.20, 0.30, 0.40, 0.50}, 0.33166247903554, false},
		// Mixed samples
		{audio.F64Samples{0.10}, 0.10, false},
		{audio.F64Samples{0.10, -0.20}, 0.15811388300841897, false},
		{audio.F64Samples{0.10, -0.20, 0.30, -0.40, 0.50}, 0.33166247903554, false},
	}

	for i, test := range tests {
		if rms := RMSF64Samples(test.samples); rms != test.result {
			// If expected result is NaN, continue
			if math.IsNaN(rms) && test.isNaN {
				continue
			}

			t.Fatalf("[%02d] unexpected result: %v != %v", i, rms, test.result)
		}
	}
}

// Test_validateOptions verifies that validateOptions correctly sets
// sane default options
func Test_validateOptions(t *testing.T) {
	// Copy default options, but set sharpness to zero because it
	// is not adjusted from the user's setting
	defaultOptions := DefaultOptions
	defaultOptions.Sharpness = 0

	var tests = []struct {
		input  Options
		output Options
	}{
		// Empty options set, defaults used
		{Options{}, defaultOptions},
		// Empty uint values, defaults used
		{Options{ComputeOptions{Resolution: 0}, ImageOptions{ScaleX: 0, ScaleY: 0}}, defaultOptions},
		// Empty color values, defaults used
		{Options{ComputeOptions{}, ImageOptions{ForegroundColor: nil, BackgroundColor: nil, AlternateColor: nil}}, defaultOptions},
		// Sharpness set, defaults used
		{Options{ComputeOptions{}, ImageOptions{Sharpness: 1}}, DefaultOptions},
	}

	for i, test := range tests {
		// Validate and compare output
		output := validateOptions(test.input)

		// Since functions are only equal if both nil, we cannot check their equality
		output.Function = nil
		test.output.Function = nil

		if !reflect.DeepEqual(output, test.output) {
			t.Fatalf("[%02d] unexpected result: %#v != %#v", i, output, test.output)
		}
	}
}
