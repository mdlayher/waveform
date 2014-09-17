package waveform

import (
	"bytes"
	"fmt"
	"image/png"
	"math/rand"
	"os"
	"time"
)

// ExampleNew provides example usage of New, using a media file from the filesystem.
func ExampleNew() {
	// New accepts io.Reader, so we will use a media file in the filesystem
	file, err := os.Open("./test/tone16bit.flac")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("open:", file.Name())
	defer file.Close()

	// Directly generate waveform image from audio file
	img, err := New(file, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Encode image as PNG into buffer
	buf := bytes.NewBuffer(nil)
	if err := png.Encode(buf, img); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("encoded: %d bytes", buf.Len())

	// Output:
	// open: ./test/tone16bit.flac
	// encoded: 88 bytes
}

// ExampleComputeValues provides example usage of ComputeValues, using a media
// file from the filesystem.
func ExampleComputeValues() {
	// ComputeValues accepts io.Reader, so we will use a media file in the filesystem
	file, err := os.Open("./test/tone16bit.flac")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("open:", file.Name())
	defer file.Close()

	// Compute values from samples in audio file
	values, err := ComputeValues(file, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Show the number of values computed
	fmt.Println("values:", len(values))

	// Output:
	// open: ./test/tone16bit.flac
	// values: 5
}

// ExampleDrawImage provides example usage of DrawImage, using randomly
// generated values.
func ExampleDrawImage() {
	rand.Seed(time.Now().UnixNano())

	// Generate random float64 values
	values := make([]float64, 5)
	for i := 0; i < len(values); i++ {
		values[i] = rand.Float64()
	}
	fmt.Println("values:", len(values))

	// Directly generate waveform image from random float64 values
	img := DrawImage(values, nil)
	fmt.Println("bounds:", img.Bounds())

	// Output:
	// values: 5
	// bounds: (0,0)-(5,128)
}
