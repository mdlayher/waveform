package waveform

import (
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"testing"
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

// TestWaveformComputeWAVOK verifies that the Waveform.Compute method produces
// appropriate computed samples and error for an input audio stream.
// The input stream is in WAV format, and no errors should occur.
func TestWaveformComputeWAVOK(t *testing.T) {
	testWaveformCompute(t, bytes.NewReader(wavFile), nil,
		[]float64{
			0.7071166200538482,
			0.7071166444294603,
			0.7071166239921965,
			0.7071165471800284,
			0.7071166825227931,
			0.7071166825227931,
		},
		nil,
	)
}

// TestWaveformComputeWAVErrInvalidData verifies that the Waveform.Compute method produces
// appropriate computed samples and error for an input audio stream.
// The input stream is in WAV format, but contains invalid data.
func TestWaveformComputeWAVErrInvalidData(t *testing.T) {
	testWaveformCompute(t, bytes.NewReader([]byte{'R', 'I', 'F', 'F', 'W', 'A', 'V', '3'}), ErrInvalidData, nil, nil)
}

// TestWaveformComputeWAVEOF verifies that the Waveform.Compute method produces
// appropriate computed samples and error for an input audio stream.
// The input stream is in WAV format, but reaches EOF during decoding.
func TestWaveformComputeWAVEOF(t *testing.T) {
	testWaveformCompute(t, bytes.NewReader([]byte{'R', 'I', 'F', 'F'}), io.EOF, nil, nil)
}

// TestWaveformComputeFLACOK verifies that the Waveform.Compute method produces
// appropriate computed samples and error for an input audio stream.
// The input stream is in FLAC format, and no errors should occur.
func TestWaveformComputeFLACOK(t *testing.T) {
	testWaveformCompute(t, bytes.NewReader(flacFile), nil,
		[]float64{
			0.7071166200538482,
			0.7071166444294603,
			0.7071166239921965,
			0.7071165471800284,
			0.7071166825227931,
		},
		nil,
	)
}

// TestWaveformComputeFLACErrInvalidData verifies that the Waveform.Compute method produces
// appropriate computed samples and error for an input audio stream.
// The input stream is in FLAC format, but contains invalid data.
func TestWaveformComputeFLACErrInvalidData(t *testing.T) {
	testWaveformCompute(t, bytes.NewReader([]byte{'f', 'L', 'a', 'C'}), ErrInvalidData, nil, nil)
}

// TestWaveformComputeMP3ErrFormat verifies that the Waveform.Compute method produces
// appropriate computed samples and error for an input audio stream.
// The input stream is in MP3 format, and should produce an unsupported format error.
func TestWaveformComputeMP3ErrFormat(t *testing.T) {
	testWaveformCompute(t, bytes.NewReader(mp3File), ErrFormat, nil, nil)
}

// TestWaveformComputeOggVorbisErrFormat verifies that the Waveform.Compute method produces
// appropriate computed samples and error for an input audio stream.
// The input stream is in Ogg Vorbis format, and should produce an unsupported format error.
func TestWaveformComputeOggVorbisErrFormat(t *testing.T) {
	testWaveformCompute(t, bytes.NewReader(oggVorbisFile), ErrFormat, nil, nil)
}

// TestWaveformComputeSampleFuncFunctionNil verifies that the Waveform.Compute method returns an error
// if a nil SampleReduceFunc member is set.
func TestWaveformComputeSampleFuncFunctionNil(t *testing.T) {
	if _, err := new(Waveform).Compute(); err != errSampleFunctionNil {
		t.Fatalf("unexpected Compute error: %v != %v", err, errSampleFunctionNil)
	}
}

// TestWaveformComputeResolutionZero verifies that the Waveform.Compute method returns an error
// if the resolution member is 0.
func TestWaveformComputeResolutionZero(t *testing.T) {
	w := &Waveform{
		sampleFn: RMSF64Samples,
	}
	if _, err := w.Compute(); err != errResolutionZero {
		t.Fatalf("unexpected Compute error: %v != %v", err, errResolutionZero)
	}
}

// testWaveformCompute is a test helper which verifies that generating a Waveform
// from an input io.Reader, applying the appropriate OptionsFunc, and calling its
// Compute method, will produce the appropriate computed values and error.
func testWaveformCompute(t *testing.T, r io.Reader, err error, values []float64, fn []OptionsFunc) {
	// Generate new Waveform, apply any functions
	w, wErr := New(r, fn...)
	if wErr != nil {
		t.Fatal(err)
	}

	// Compute values from waveform
	computed, cErr := w.Compute()
	if cErr != err {
		t.Fatalf("unexpected Compute error: %v != %v", cErr, err)
	}

	// Ensure values slices match in length
	if len(values) != len(computed) {
		t.Fatalf("unexpected Compute values length: %v != %v [%v != %v]", len(values), len(computed), values, computed)
	}

	// Iterate all values and check for equality
	for i := range values {
		if values[i] != computed[i] {
			t.Fatalf("unexpected Compute value at index %d: %v != %v", i, values[i], computed[i])
		}
	}
}
