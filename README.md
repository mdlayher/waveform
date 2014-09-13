waveform
========

Work-in-progress Go package to generate waveform images from audio files.  MIT Licensed.

The package itself is not yet implemented, but it will be derived from the work done on the
bundled `waveform` command.

Usage
=====

To install and use `waveform`, simply run:

```
$ go install github.com/mdlayher/waveform/...
```

The `waveform` binary is now installed in your `$GOPATH`.  It has several options available
for generating waveform images:

```
$ waveform -h
Usage of waveform:
  -color="#000000": hex color of output waveform image
  -in="": input audio file
  -out="": output PNG waveform image file
  -sharpness=1: sharpening factor used to add curvature to a scaled image
  -x=1: scaling factor for image X-axis
  -y=1: scaling factor for image Y-axis
```

The most basic usage requires the `-in` and `-out` parameters.  `waveform` currently supports
both WAV and FLAC audio files.

Example
=======

Use `waveform` to generate a waveform image from a FLAC audio file, and scale it both vertically
and horizontally.

```
$ waveform -in ~/Music/FLAC/Boston/1976\ -\ Boston/02\ -\ Peace\ Of\ Mind.flac -out ~/waveform.png -x 5 -y 2
waveform: 2014/09/13 15:02:23  audio: /home/matt/Music/FLAC/Boston/1976 - Boston/02 - Peace Of Mind.flac [flac, 44100Hz, 2ch]
waveform: 2014/09/13 15:02:27  scale: [1510x256]: x * 5, y * 2
waveform: 2014/09/13 15:02:27 maxRMS: 0.261 [scale: 3.000]
waveform: 2014/09/13 15:02:27  image: /home/matt/waveform.png
```

The result is a waveform image, located at `~/waveform.png`:

![waveform](https://cloud.githubusercontent.com/assets/1926905/4261650/b020c3c2-3b78-11e4-933c-c0b81e282973.png)
