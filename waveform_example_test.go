package waveform

import (
	"bytes"
	"fmt"
	"image/color"
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

	// Directly generate waveform image from audio file, applying any number
	// of options functions along the way
	img, err := Generate(file,
		// Solid white background
		BGColorFunction(SolidColor(color.White)),
		// Striped red, green, and blue foreground
		FGColorFunction(StripeColor(
			color.RGBA{255, 0, 0, 255},
			color.RGBA{0, 255, 0, 255},
			color.RGBA{0, 0, 255, 255},
		)),
		// Scaled 10x horizontally, 2x vertically
		Scale(10, 2),
	)
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
	fmt.Printf("encoded: %d bytes\nresolution: %s", buf.Len(), img.Bounds().Max)

	// Output:
	// open: ./test/tone16bit.flac
	// encoded: 344 bytes
	// resolution: (50,256)
}
