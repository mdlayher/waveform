package waveform

import (
	"math"
	"testing"

	"azul3d.org/audio.v1"
)

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
