package waveform

import (
	"bytes"
	"testing"
)

// BenchmarkNewWAV checks the performance of the New() function with a WAV file
func BenchmarkNewWAV(b *testing.B) {
	benchmarkNew(b, wavFile, nil)
}

// BenchmarkNewWAVResolution2 checks the performance of the New() function with a WAV file
// at 2x resolution
func BenchmarkNewWAVResolution2(b *testing.B) {
	benchmarkNew(b, wavFile, &Options{
		Resolution: 2,
	})
}

// BenchmarkNewWAVResolution4 checks the performance of the New() function with a WAV file
// at 4x resolution
func BenchmarkNewWAVResolution4(b *testing.B) {
	benchmarkNew(b, wavFile, &Options{
		Resolution: 4,
	})
}

// BenchmarkNewWAVResolution8 checks the performance of the New() function with a WAV file
// at 8x resolution
func BenchmarkNewWAVResolution8(b *testing.B) {
	benchmarkNew(b, wavFile, &Options{
		Resolution: 8,
	})
}

// BenchmarkNewWAVResolution16 checks the performance of the New() function with a WAV file
// at 16x resolution
func BenchmarkNewWAVResolution16(b *testing.B) {
	benchmarkNew(b, wavFile, &Options{
		Resolution: 16,
	})
}

// BenchmarkNewWAVScaleX2 checks the performance of the New() function with a WAV file
// at 2x X scale
func BenchmarkNewWAVScaleX2(b *testing.B) {
	benchmarkNew(b, wavFile, &Options{
		ScaleX: 2,
	})
}

// BenchmarkNewWAVScaleX4 checks the performance of the New() function with a WAV file
// at 4x X scale
func BenchmarkNewWAVScaleX4(b *testing.B) {
	benchmarkNew(b, wavFile, &Options{
		ScaleX: 4,
	})
}

// BenchmarkNewWAVScaleX8 checks the performance of the New() function with a WAV file
// at 8x X scale
func BenchmarkNewWAVScaleX8(b *testing.B) {
	benchmarkNew(b, wavFile, &Options{
		ScaleX: 8,
	})
}

// BenchmarkNewWAVScaleX16 checks the performance of the New() function with a WAV file
// at 16x X scale
func BenchmarkNewWAVScaleX16(b *testing.B) {
	benchmarkNew(b, wavFile, &Options{
		ScaleX: 16,
	})
}

// BenchmarkNewWAVScaleY2 checks the performance of the New() function with a WAV file
// at 2x Y scale
func BenchmarkNewWAVScaleY2(b *testing.B) {
	benchmarkNew(b, wavFile, &Options{
		ScaleY: 2,
	})
}

// BenchmarkNewWAVScaleY4 checks the performance of the New() function with a WAV file
// at 4x Y scale
func BenchmarkNewWAVScaleY4(b *testing.B) {
	benchmarkNew(b, wavFile, &Options{
		ScaleY: 4,
	})
}

// BenchmarkNewWAVScaleY8 checks the performance of the New() function with a WAV file
// at 8x Y scale
func BenchmarkNewWAVScaleY8(b *testing.B) {
	benchmarkNew(b, wavFile, &Options{
		ScaleY: 8,
	})
}

// BenchmarkNewWAVScaleY16 checks the performance of the New() function with a WAV file
// at 16x Y scale
func BenchmarkNewWAVScaleY16(b *testing.B) {
	benchmarkNew(b, wavFile, &Options{
		ScaleY: 16,
	})
}

// BenchmarkNewFLAC checks the performance of the New() function with a FLAC file
func BenchmarkNewFLAC(b *testing.B) {
	benchmarkNew(b, flacFile, nil)
}

// BenchmarkNewFLACResolution2 checks the performance of the New() function with a FLAC file
// at 2x resolution
func BenchmarkNewFLACResolution2(b *testing.B) {
	benchmarkNew(b, flacFile, &Options{
		Resolution: 2,
	})
}

// BenchmarkNewFLACResolution4 checks the performance of the New() function with a FLAC file
// at 4x resolution
func BenchmarkNewFLACResolution4(b *testing.B) {
	benchmarkNew(b, flacFile, &Options{
		Resolution: 4,
	})
}

// BenchmarkNewFLACResolution8 checks the performance of the New() function with a FLAC file
// at 8x resolution
func BenchmarkNewFLACResolution8(b *testing.B) {
	benchmarkNew(b, flacFile, &Options{
		Resolution: 8,
	})
}

// BenchmarkNewFLACResolution16 checks the performance of the New() function with a FLAC file
// at 16x resolution
func BenchmarkNewFLACResolution16(b *testing.B) {
	benchmarkNew(b, flacFile, &Options{
		Resolution: 16,
	})
}

// BenchmarkNewFLACScaleX2 checks the performance of the New() function with a FLAC file
// at 2x X scale
func BenchmarkNewFLACScaleX2(b *testing.B) {
	benchmarkNew(b, flacFile, &Options{
		ScaleX: 2,
	})
}

// BenchmarkNewFLACScaleX4 checks the performance of the New() function with a FLAC file
// at 4x X scale
func BenchmarkNewFLACScaleX4(b *testing.B) {
	benchmarkNew(b, flacFile, &Options{
		ScaleX: 4,
	})
}

// BenchmarkNewFLACScaleX8 checks the performance of the New() function with a FLAC file
// at 8x X scale
func BenchmarkNewFLACScaleX8(b *testing.B) {
	benchmarkNew(b, flacFile, &Options{
		ScaleX: 8,
	})
}

// BenchmarkNewFLACScaleX16 checks the performance of the New() function with a FLAC file
// at 16x X scale
func BenchmarkNewFLACScaleX16(b *testing.B) {
	benchmarkNew(b, flacFile, &Options{
		ScaleX: 16,
	})
}

// BenchmarkNewFLACScaleY2 checks the performance of the New() function with a FLAC file
// at 2x Y scale
func BenchmarkNewFLACScaleY2(b *testing.B) {
	benchmarkNew(b, flacFile, &Options{
		ScaleY: 2,
	})
}

// BenchmarkNewFLACScaleY4 checks the performance of the New() function with a FLAC file
// at 4x Y scale
func BenchmarkNewFLACScaleY4(b *testing.B) {
	benchmarkNew(b, flacFile, &Options{
		ScaleY: 4,
	})
}

// BenchmarkNewFLACScaleY8 checks the performance of the New() function with a FLAC file
// at 8x Y scale
func BenchmarkNewFLACScaleY8(b *testing.B) {
	benchmarkNew(b, flacFile, &Options{
		ScaleY: 8,
	})
}

// BenchmarkNewFLACScaleY16 checks the performance of the New() function with a FLAC file
// at 16x Y scale
func BenchmarkNewFLACScaleY16(b *testing.B) {
	benchmarkNew(b, flacFile, &Options{
		ScaleY: 16,
	})
}

// benchmarkNew contains common logic for benchmarking New
func benchmarkNew(b *testing.B, data []byte, options *Options) {
	for i := 0; i < b.N; i++ {
		New(bytes.NewReader(data), options)
	}
}
