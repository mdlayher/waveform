waveform [![Build Status](https://travis-ci.org/mdlayher/waveform.svg?branch=master)](https://travis-ci.org/mdlayher/waveform) [![GoDoc](http://godoc.org/github.com/mdlayher/waveform?status.svg)](http://godoc.org/github.com/mdlayher/waveform)
========

Go package capable of generating waveform images from audio streams.  MIT Licensed.

This library supports any audio streams which the [azul3d/audio.v1](http://azul3d.org/audio.v1)
package is able to decode.  At the time of writing, this includes:
  - WAV
  - FLAC

An example binary called `waveform` is provided which show's the library's usage.
Please see [cmd/waveform/README.md](https://github.com/mdlayher/waveform/blob/master/cmd/waveform/README.md)
for details.

Here is an example, generated using `waveform` from Boston's "Peace of Mind".

```
$ cat ~/Music/FLAC/Boston/1976\ -\ Boston/02\ -\ Peace\ Of\ Mind.flac | waveform -x 5 -y 2 > ~/waveform.png
```

![waveform](https://cloud.githubusercontent.com/assets/1926905/4261650/b020c3c2-3b78-11e4-933c-c0b81e282973.png)
