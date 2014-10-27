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
  -resolution=1: number of times audio is read and drawn per second of audio
  -sharpness=1: sharpening factor used to add curvature to a scaled image
  -x=1: scaling factor for image X-axis
  -y=1: scaling factor for image Y-axis
```

`waveform` currently supports both WAV and FLAC audio files.  An audio stream must
be passed on `stdin`, and the resulting, PNG-encoded image will be written to `stdout`.
Any errors which occur will be written to `stderr`.

Example
=======

Use `waveform` to generate a waveform image from a FLAC audio file, and scale it both
vertically and horizontally.

```
$ cat ~/Music/FLAC/Boston/1976\ -\ Boston/02\ -\ Peace\ Of\ Mind.flac | waveform -x 5 -y 2 > ~/waveform.png
```

The result is a waveform image, located at `~/waveform.png`:

![waveform](https://cloud.githubusercontent.com/assets/1926905/4261650/b020c3c2-3b78-11e4-933c-c0b81e282973.png)
