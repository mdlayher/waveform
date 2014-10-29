package waveform

import (
	"bytes"
	"fmt"
	"image/png"
	"os"
)

// ExampleGenerate provides example usage of Generate, using a media file from the filesystem.
// Generate is typically used for one-time, direct creation of an image.Image from
// an input audio stream.
func ExampleGenerate() {
	// Generate accepts io.Reader, so we will use a media file in the filesystem
	file, err := os.Open("./test/tone16bit.flac")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("open:", file.Name())
	defer file.Close()

	// Directly generate waveform image from audio file
	img, err := Generate(file)
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
