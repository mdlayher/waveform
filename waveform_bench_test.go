package waveform

import (
	"bytes"
	"math/rand"
	"testing"
	"time"

	"azul3d.org/engine/audio"
)

// BenchmarkGenerateWAV checks the performance of the Generate() function with a WAV file
func BenchmarkGenerateWAV(b *testing.B) {
	benchmarkGenerate(b, wavFile)
}

// BenchmarkGenerateFLAC checks the performance of the Generate() function with a FLAC file
func BenchmarkGenerateFLAC(b *testing.B) {
	benchmarkGenerate(b, flacFile)
}

// BenchmarkWaveformComputeWAV checks the performance of the WaveformCompute() function with a WAV file
func BenchmarkWaveformComputeWAV(b *testing.B) {
	benchmarkWaveformCompute(b, wavFile)
}

// BenchmarkWaveformComputeFLAC checks the performance of the WaveformCompute() function with a FLAC file
func BenchmarkWaveformComputeFLAC(b *testing.B) {
	benchmarkWaveformCompute(b, flacFile)
}

// BenchmarkWaveformDraw60 checks the performance of the WaveformDraw() function
// with approximately 60 seconds of computed values
func BenchmarkWaveformDraw60(b *testing.B) {
	benchmarkWaveformDraw(b, 60)
}

// BenchmarkWaveformDraw120 checks the performance of the WaveformDraw() function
// with approximately 120 seconds of computed values
func BenchmarkWaveformDraw120(b *testing.B) {
	benchmarkWaveformDraw(b, 120)
}

// BenchmarkWaveformDraw240 checks the performance of the WaveformDraw() function
// with approximately 240 seconds of computed values
func BenchmarkWaveformDraw240(b *testing.B) {
	benchmarkWaveformDraw(b, 240)
}

// BenchmarkWaveformDraw480 checks the performance of the WaveformDraw() function
// with approximately 480 seconds of computed values
func BenchmarkWaveformDraw480(b *testing.B) {
	benchmarkWaveformDraw(b, 480)
}

// BenchmarkWaveformDraw960 checks the performance of the WaveformDraw() function
// with approximately 960 seconds of computed values
func BenchmarkWaveformDraw960(b *testing.B) {
	benchmarkWaveformDraw(b, 960)
}

// BenchmarkRMSF64Samples22050 checks the performance of the RMSF64Samples() function
// with 22050 samples
func BenchmarkRMSF64Samples22050(b *testing.B) {
	benchmarkRMSF64Samples(b, 22050)
}

// BenchmarkRMSF64Samples44100 checks the performance of the RMSF64Samples() function
// with 44100 samples
func BenchmarkRMSF64Samples44100(b *testing.B) {
	benchmarkRMSF64Samples(b, 44100)
}

// BenchmarkRMSF64Samples88200 checks the performance of the RMSF64Samples() function
// with 88200 samples
func BenchmarkRMSF64Samples88200(b *testing.B) {
	benchmarkRMSF64Samples(b, 88200)
}

// BenchmarkRMSF64Samples176400 checks the performance of the RMSF64Samples() function
// with 176400 samples
func BenchmarkRMSF64Samples176400(b *testing.B) {
	benchmarkRMSF64Samples(b, 176400)
}

// benchmarkGenerate contains common logic for benchmarking Generate
func benchmarkGenerate(b *testing.B, data []byte) {
	for i := 0; i < b.N; i++ {
		Generate(bytes.NewReader(data))
	}
}

// benchmarkWaveformCompute contains common logic for benchmarking Waveform.Compute
func benchmarkWaveformCompute(b *testing.B, data []byte) {
	for i := 0; i < b.N; i++ {
		w, err := New(bytes.NewReader(data))
		if err != nil {
			panic(err)
		}

		w.Compute()
	}
}

// benchmarkWaveformDraw contains common logic for benchmarking Waveform.Draw
func benchmarkWaveformDraw(b *testing.B, count int) {
	values := make([]float64, count)
	for i := 0; i < b.N; i++ {
		w, err := New(nil)
		if err != nil {
			panic(err)
		}

		w.Draw(values)
	}
}

// benchmarkRMSF64Samples contains common logic for benchmarking RMSF64Samples
func benchmarkRMSF64Samples(b *testing.B, count int) {
	// Generate slice of samples
	rand.Seed(time.Now().UnixNano())
	var samples audio.Float64
	for i := 0; i < count; i++ {
		samples = append(samples, rand.Float64())
	}

	// Reset timer and start benchmark
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		RMSF64Samples(samples)
	}
}
