package waveform

import (
	"io/ioutil"
	"log"
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
