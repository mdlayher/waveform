package waveform

import (
	"fmt"
	"image/color"
	"testing"
)

// TestOptionsError verifies that the format of OptionsError.Error does
// not change.
func TestOptionsError(t *testing.T) {
	var tests = []struct {
		option string
		reason string
	}{
		{"foo", "bar"},
		{"baz", "qux"},
		{"one", "two"},
	}

	for _, test := range tests {
		// Generate options error
		opErr := &OptionsError{
			Option: test.option,
			Reason: test.reason,
		}

		// Verify correct format
		if opErr.Error() != fmt.Sprintf("%s: %s", test.option, test.reason) {
			t.Fatalf("unexpected Error string: %v", opErr.Error())
		}
	}
}

// TestOptionColorsOK verifies that Colors returns no error with acceptable input.
func TestOptionColorsOK(t *testing.T) {
	testWaveformOptionFunc(t, Colors(color.Black, color.Black, color.Black), nil)
}

// TestOptionColorsNilForeground verifies that Colors does not accept a nil
// foreground color.
func TestOptionColorsNilForeground(t *testing.T) {
	testWaveformOptionFunc(t, Colors(nil, color.Black, color.Black), errColorsNilForeground)
}

// TestOptionColorsNilBackground verifies that Colors does not accept a nil
// backround color.
func TestOptionColorsNilBackground(t *testing.T) {
	testWaveformOptionFunc(t, Colors(color.Black, nil, color.Black), errColorsNilBackground)
}

// TestOptionFunctionOK verifies that Function returns no error with acceptable input.
func TestOptionFunctionOK(t *testing.T) {
	testWaveformOptionFunc(t, Function(RMSF64Samples), nil)
}

// TestOptionFunctionNil verifies that Function does not accept a nil SampleReduceFunc.
func TestOptionFunctionNil(t *testing.T) {
	testWaveformOptionFunc(t, Function(nil), errFunctionNil)
}

// TestOptionResolutionOK verifies that Resolution returns no error with acceptable input.
func TestOptionResolutionOK(t *testing.T) {
	testWaveformOptionFunc(t, Resolution(1), nil)
}

// TestOptionResolutionZero verifies that Resolution does not accept integer 0.
func TestOptionResolutionZero(t *testing.T) {
	testWaveformOptionFunc(t, Resolution(0), errResolutionZero)
}

// TestOptionScaleOK verifies that Scale returns no error with acceptable input.
func TestOptionScaleOK(t *testing.T) {
	testWaveformOptionFunc(t, Scale(1, 1), nil)
}

// TestOptionScaleXZero verifies that Scale does not accept an X value integer 0.
func TestOptionScaleXZero(t *testing.T) {
	testWaveformOptionFunc(t, Scale(0, 1), errScaleXZero)
}

// TestOptionScaleYZero verifies that Scale does not accept an Y value integer 0.
func TestOptionScaleYZero(t *testing.T) {
	testWaveformOptionFunc(t, Scale(1, 0), errScaleYZero)
}

// TestOptionScaleClippingOK verifies that ScaleClipping returns no error.
func TestOptionScaleClippingOK(t *testing.T) {
	testWaveformOptionFunc(t, ScaleClipping(), nil)
}

// TestOptionSharpnessOK verifies that Sharpness returns no error.
func TestOptionSharpnessOK(t *testing.T) {
	testWaveformOptionFunc(t, Sharpness(0), nil)
}

// TestWaveformSetColors verifies that the Waveform.SetColors method properly
// modifies struct members.
func TestWaveformSetColors(t *testing.T) {
	// Predefined test values
	fg := color.Black
	bg := color.White
	alt := color.White

	// Generate empty Waveform, apply parameters
	w := &Waveform{}
	if err := w.SetColors(fg, bg, alt); err != nil {
		t.Fatal(err)
	}

	// Validate that struct members are set properly
	if w.fg != fg {
		t.Fatalf("unexpected foreground color: %v != %v", w.fg, fg)
	}
	if w.bg != bg {
		t.Fatalf("unexpected background color: %v != %v", w.bg, bg)
	}
	if w.alt != alt {
		t.Fatalf("unexpected alternate color: %v != %v", w.alt, alt)
	}
}

// TestWaveformSetFunction verifies that the Waveform.SetFunction method properly
// modifies struct members.
func TestWaveformSetFunction(t *testing.T) {
	// Generate empty Waveform, apply parameters
	w := &Waveform{}
	if err := w.SetFunction(RMSF64Samples); err != nil {
		t.Fatal(err)
	}

	// Validate that struct members are set properly
	if w.function == nil {
		t.Fatalf("SetFunction failed, nil function member")
	}
}

// TestWaveformSetResolution verifies that the Waveform.SetResolution method properly
// modifies struct members.
func TestWaveformSetResolution(t *testing.T) {
	// Predefined test values
	res := uint(1)

	// Generate empty Waveform, apply parameters
	w := &Waveform{}
	if err := w.SetResolution(res); err != nil {
		t.Fatal(err)
	}

	// Validate that struct members are set properly
	if w.resolution != res {
		t.Fatalf("unexpected resolution: %v != %v", w.resolution, res)
	}
}

// TestWaveformSetScale verifies that the Waveform.SetScale method properly
// modifies struct members.
func TestWaveformSetScale(t *testing.T) {
	// Predefined test values
	x := uint(1)
	y := uint(1)

	// Generate empty Waveform, apply parameters
	w := &Waveform{}
	if err := w.SetScale(x, y); err != nil {
		t.Fatal(err)
	}

	// Validate that struct members are set properly
	if w.scaleX != x {
		t.Fatalf("unexpected scale X: %v != %v", w.scaleX, x)
	}
	if w.scaleY != y {
		t.Fatalf("unexpected scale Y: %v != %v", w.scaleY, y)
	}
}

// TestWaveformSetScaleClipping verifies that the Waveform.SetScaleClipping method properly
// modifies struct members.
func TestWaveformSetScaleClipping(t *testing.T) {
	// Generate empty Waveform, apply function
	w := &Waveform{}
	if err := w.SetScaleClipping(); err != nil {
		t.Fatal(err)
	}

	// Validate that struct members are set properly
	if !w.scaleClipping {
		t.Fatalf("SetScaleClipping failed, false scaleClipping member")
	}
}

// TestWaveformSetSharpness verifies that the Waveform.SetSharpness method properly
// modifies struct members.
func TestWaveformSetSharpness(t *testing.T) {
	// Predefined test values
	sharpness := uint(1)

	// Generate empty Waveform, apply parameters
	w := &Waveform{}
	if err := w.SetSharpness(sharpness); err != nil {
		t.Fatal(err)
	}

	// Validate that struct members are set properly
	if w.sharpness != sharpness {
		t.Fatalf("unexpected sharpness: %v != %v", w.sharpness, sharpness)
	}
}

// testWaveformOptionFunc is a test helper which verifies that applying the
// input OptionsFunc to a new Waveform struct generates the appropriate
// error output.
func testWaveformOptionFunc(t *testing.T, fn OptionsFunc, err error) {
	if _, wErr := New(nil, fn); wErr != err {
		t.Fatalf("unexpected error: %v != %v", wErr, err)
	}
}
