![build](https://github.com/transcriptaze/wav2png/workflows/go/badge.svg)

# wav2png

Command line version of wav2png to renders a WAV file as a PNG image, with options to draw a grid, custom colouring and anti-aliasing. The command line version currently supports the following WAV encodings:

- 16-bit signed PCM 
- 32-bit floating point PCM
- EXTENSIBLE 32-bit floating point PCM 

An online version implemented by compiling this library to WASM can be found [here](https://transcriptaze.github.io/W2P.html) (the online verson supports any audio format supported by the browser).

The command line version includes two utilities:

- `wav2png`, which renders a WAV file to a single PNG image
- `wav2mp4`, which renders a WAV file to a set of PNG frames that can be used with _ffmpeg_ to create an MP4
   rendering of the WAV file

## Raison d'être

wav2png was initially created as a Go utility library to render an audio file as an anti-aliased waveform for a
WASM project - it just seemed like a good idea to add a standalone command line version.

## Installation

Platform specific executables can be downloaded from the [releases](https://github.com/transcriptaze/wav2png/releases) 
page. Installation is straightforward - download and extract the archive for your platform and place the executables in 
a directory in your PATH. 

### Building from source

Assuming you have `Go` and `make` installed:

```
git clone https://github.com/transcriptaze/wav2png.git
cd wav2png/go
make build
```

If you prefer not to use `make`:
```
git clone https://github.com/transcriptaze/wav2png.git
cd wav2png/go
mkdir bin
go build -o bin ./...
```

#### Dependencies

_NONE_

## wav2png

Command line:
```
wav2png [--debug] [options] [--out <path>] <wav>

  <wav>         WAV file to render.

  --out <path>  File path for PNG file - if <path> is a directory, the WAV file name is
                used. Defaults to the WAV file base path.

  --debug       Displays occasionally useful diagnostic information.

Options:
  --settings <file>      JSON file with the default settings for the height, width, etc. Defaults to 
                         .settings.json if not specified, falling back to internal default values if 
                         .settings.json does not exist.

  --width <pixels>       Width (in pixels) of the PNG image. Valid values are in the range 32 to 8192, 
                         defaults to 645px.

  --height <pixels>      Height (in pixels) of the PNG image. Valid values are in the range 32 to 8192, 
                         defaults to 395px.
  
  --padding <pixels>     Padding (in pixels) between the border of the PNG and the extent of the rendered
                         waveform. Valid values are -16 to +32, defaults to 2px.

  --palette <palette>    Palette used to colour the waveform. May be the name of one of the internal colour
                         palettes or a user provided PNG file. Defaults to 'ice'
  
  --fill <fillspec>      Fill specification for the background colour, in the form type:colour 
                         e.g. solid:#0000ffff. Currently the only fill type supported is 'solid', defaults
                         to solid:#000000ff.

  --grid <gridspec>      Grid specification for an optional rectilinear grid, in the form 
                         type:colour:size:overlay, e.g.
                         - none
                         - square:#008000ff:~64
                         - rectangle:#008000ff:~64x48:overlay
                         
                         The size may preceded by a 'fit':
                         - ~  approximate
                         - =  exact
                         - ≥  at least
                         - ≤  at most
                         - >  greater than
                         - <  less than

                         If gridspec includes :overlay, the grid is rendered 'in front' of the waveform.

                         The default gridspec is 'square:#008000ff:~64'

  --antialias <kernel>   The antialising kernel applied to soften the rendered PNG. Valid values are:
                         - none        no antialiasing
                         - horizontal  blurs horizontal edges
                         - vertical    blurs vertical edges
                         - soft        blurs both horizontal and vertical edges

                         The default kernel is 'vertical'.

  --scale <scale>        A vertical scaling factor to size the height of the rendered waveform. The valid 
                         range is 0.2 to 5.0, defaults to 1.0.

  --mix  <mixspec>       Specifies how to combine channels from a stereo WAV file. Valid values are:
                         - 'L'    Renders the left channel only
                         - 'R'    Renders the right channel only
                         - 'L+R'  Combines the left and right channels
                         
                         Defaults to 'L+R'.

  --start <time>         The start time of the segment of audio to render, in Go time format (e.g. 10s or
                         1m5s). Defaults to 0s.

  --end <time>           The end time of the segment of audio to render, in Go time format (e.g. 10s or
                         1m5s). Defaults to the end of the audio.


Example:

wav2png --debug                                        \
          --settings 'settings.json'                     \
          --height 390                                   \
          --width 641                                    \
          --padding 0                                    \
          --palette 'amber.png'                          \
          --fill 'solid:#0000ffff'                       \
          --grid 'rectangular:#800000ff:~32x128:overlay' \
          --antialias 'soft'                             \
          --scale 0.5                                    \
          --start 0.5s                                   \
          --end 1.5s                                     \
          --mix 'L+R'                                    \
          --out example.png                              \
          example.wav
```

## wav2mp4

Command line:
```
wav2mp4 [--debug] [options] [--out <path>] --window <duration> --fps <frame rate> --cursor <cursor spec> <wav>

  <wav>                  WAV file to render.

  --out <path>           File path for MP4 file - if <path> is a directory, the WAV file name is
                         used and defaults to the WAV file base path. wav2mp4 generates a set of ffmpeg
                         frames files in the 'frames' subdirectory of the out file directory. 

  --window <duration>    The interval allotted to a single frame, in Go time format e.g. --window 1s.
                         The window interval must be less than the duration of the MP4.

  --fps <frame rate>     Frame rate for the MP4 in frames per second e.g. --fps 30

  --cursor <cursorspec>  Cursor to indicate for the current play position. A cursor is specified by the 
                         image source and dynamic:

                         --cursor <image>:<dynamic>

                         where image may be:
                         - none
                         - red  (internal 'red' cursor)
                         - blue (internal 'blue' cursor)
                         - a PNG file with a custom cursor image

                         The cursor 'dynamic' defaults to 'sweep' if not specified, but may be one of
                         the following:
                         - sweep  Moves linearly from left to right over the duration of the MP4
                         - left   Fixed on left side
                         - right  Fixed on right side
                         - center Fixed in center of frame
                         - ease   Migrates from the left to center of the frame, before moving to the
                                  right side to finish
                         - erf    Moves 'sigmoidally' from left to right over the duration of the MP4, 
                                  with the sigmoid defined by the inverse error function

  --debug                Displays occasionally useful diagnostic information.

Options:
  --settings <file>      JSON file with the default settings for the height, width, etc. Defaults to 
                         .settings.json if not specified, falling back to internal default values if 
                         .settings.json does not exist.

  --width <pixels>       Width (in pixels) of the PNG image. Valid values are in the range 32 to 8192, 
                         defaults to 645px.

  --height <pixels>      Height (in pixels) of the PNG image. Valid values are in the range 32 to 8192, 
                         defaults to 395px.
  
  --padding <pixels>     Padding (in pixels) between the border of the PNG and the extent of the 
                         rendered waveform. Valid values are -16 to +32, defaults to 2px.

  --palette <palette>    Palette used to colour the waveform. May be the name of one of the internal 
                         colour palettes or a user provided PNG file. Defaults to 'ice'
  
  --fill <fillspec>      Fill specification for the background colour, in the form type:colour 
                         e.g. solid:#0000ffff. Currently the only fill type supported is 'solid', 
                         defaults to solid:#000000ff.

  --grid <gridspec>      Grid specification for an optional rectilinear grid, in the form 
                         type:colour:size:overlay, e.g.
                         - none
                         - square:#008000ff:~64
                         - rectangle:#008000ff:~64x48:overlay
                         
                         The size may preceded by a 'fit':
                         - ~  approximate
                         - =  exact
                         - ≥  at least
                         - ≤  at most
                         - >  greater than
                         - <  less than

                         If gridspec includes :overlay, the grid is rendered 'in front' of the 
                         waveform.

                         The default gridspec is 'square:#008000ff:~64'

  --antialias <kernel>   The antialising kernel applied to soften the rendered PNG. Valid values are:
                         - none        no antialiasing
                         - horizontal  blurs horizontal edges
                         - vertical    blurs vertical edges
                         - soft        blurs both horizontal and vertical edges

                         The default kernel is 'vertical'.

  --scale <scale>        A vertical scaling factor to size the height of the rendered waveform. The 
                         valid range is 0.2 to 5.0, defaults to 1.0.

  --mix  <mixspec>       Specifies how to combine channels from a stereo WAV file. Valid values are:
                         - 'L'    Renders the left channel only
                         - 'R'    Renders the right channel only
                         - 'L+R'  Combines the left and right channels
                         
                         Defaults to 'L+R'.

  --start <time>         The start time of the segment of audio to render, in Go time format (e.g. 10s
                         or 1m5s). Defaults to 0s.

  --end <time>           The end time of the segment of audio to render, in Go time format (e.g. 10s
                         or 1m5s). Defaults to the end of the audio.


Example:

  wav2mp4 --debug --window 1s --fps 30 --cursor red:sweep --out ./example/chirp.mp4 ./samples/chirp.wav
  cd ./example/frames
  ffmpeg -framerate 30 -i frame-%05d.png -c:v libx264 -pix_fmt yuv420p -crf 23 -y out.mp4
  ffmpeg -i out.mp4    -i ./samples/chirp.wav -c:v copy -c:a aac -y ../chirp.mp4
```
