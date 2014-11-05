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
  -alt="": hex alternate color of output waveform image
  -bg="#FFFFFF": hex background color of output waveform image
  -fg="#000000": hex foreground color of output waveform image
  -fn="solid": function used to color output waveform image [options: fuzz, gradient, solid, stripe]
  -resolution=1: number of times audio is read and drawn per second of audio
  -sharpness=1: sharpening factor used to add curvature to a scaled image
  -x=1: scaling factor for image X-axis
  -y=1: scaling factor for image Y-axis
```

`waveform` currently supports both WAV and FLAC audio files.  An audio stream must
be passed on `stdin`, and the resulting, PNG-encoded image will be written to `stdout`.
Any errors which occur will be written to `stderr`.
