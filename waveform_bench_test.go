package waveform

import (
	"bytes"
	"testing"
)

// BenchmarkNewWAV checks the performance of the New() function with a WAV file
func BenchmarkNewWAV(b *testing.B) {
	for i := 0; i < b.N; i++ {
		New(bytes.NewReader(wavFile), nil)
	}
}

// BenchmarkNewFLAC checks the performance of the New() function with a FLAC file
func BenchmarkNewFLAC(b *testing.B) {
	for i := 0; i < b.N; i++ {
		New(bytes.NewReader(flacFile), nil)
	}
}
