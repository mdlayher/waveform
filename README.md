waveform [![Build Status](https://travis-ci.org/mdlayher/waveform.svg?branch=master)](https://travis-ci.org/mdlayher/waveform) [![GoDoc](http://godoc.org/github.com/mdlayher/waveform?status.svg)](http://godoc.org/github.com/mdlayher/waveform)
========

Go package capable of generating waveform images from audio streams.  MIT Licensed.

This library supports any audio streams which the [azul3d/engine/audio](http://azul3d.org/engine/audio)
package is able to decode.  At the time of writing, this includes:
  - WAV
  - FLAC

An example binary called `waveform` is provided which show's the library's usage.
Please see [cmd/waveform/README.md](https://github.com/mdlayher/waveform/blob/master/cmd/waveform/README.md)
for details.

Examples
========

Here are several example images generated using `waveform`.  Enjoy!

Generate a waveform image, and scale it both vertically and horizontally.

```
$ cat ~/Music/02\ -\ Peace\ Of\ Mind.flac | waveform -x 5 -y 2 > ~/waveform.png
```

![waveform](https://cloud.githubusercontent.com/assets/1926905/4910038/6ce9f5d0-647a-11e4-8a93-ed54812d114d.png)

Apply a foreground and background color, to make things more interesting.

```
cat ~/Music/02\ -\ Peace\ Of\ Mind.flac | waveform -fg=#FF3300 -bg=#0099CC -x 5 -y 2 > ~/waveform_color.png
```

![waveform_color](https://cloud.githubusercontent.com/assets/1926905/4910043/757b0edc-647a-11e4-8ebd-73175246421d.png)

Apply an alternate foreground color, draw using a stripe pattern.

```
cat ~/Music/02\ -\ Peace\ Of\ Mind.flac | waveform -fg=#FF3300 -bg=#0099CC -alt=#FF9933 -fn stripe -x 5 -y 2 > ~/waveform_stripe.png
```

![waveform_stripe](https://cloud.githubusercontent.com/assets/1926905/4910067/a560f76a-647a-11e4-8562-c430134c1187.png)

Apply an alternate foreground color, draw using a random fuzz pattern.

```
cat ~/Music/02\ -\ Peace\ Of\ Mind.flac | waveform -fg=#FF3300 -bg=#0099CC -alt=#FF9933 -fn fuzz -x 5 -y 2 > ~/waveform_fuzz.png
```

![waveform_fuzz](https://cloud.githubusercontent.com/assets/1926905/4910076/c6aa0e70-647a-11e4-8385-754960c9f074.png)

Apply a new set of colors, draw using a gradient pattern.

```
cat ~/Music/02\ -\ Peace\ Of\ Mind.flac | waveform -fg=#FF0000 -bg=#00FF00 -alt=#0000FF -fn gradient -x 5 -y 2 > ~/waveform_gradient.png
```

![waveform_gradient](https://cloud.githubusercontent.com/assets/1926905/5416955/c5592f10-8202-11e4-943d-d86214b26b18.png)

Apply a checkerboard color set, draw using a checkerboard pattern.

```
cat ~/Music/02\ -\ Peace\ Of\ Mind.flac | waveform -fg=#000000 -bg=#222222 -alt=#FFFFFF -fn checker -x 5 -y 2 > ~/waveform_checker.png
```

![waveform_checker](https://cloud.githubusercontent.com/assets/1926905/4961769/e3280c96-66d2-11e4-8e3c-d0b843230589.png)
