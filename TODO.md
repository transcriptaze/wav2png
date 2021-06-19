## v0.0

### IN PROGRESS

- [x] Add fill to command line options
- [ ] Load palette from file
      - 'fire' palette
      - 'green' palette
      - 'amber' palette
      - 'blues' palette
      - 'inverted' palette
      
- [ ] Get default settings from file
- [ ] Add references to README

- [ ] Render should take a reader parameter
- [ ] Grid
      -- horizontal lines start 1 pixel above/below 0
- [ ] Load anti-aliasing kernel from file
- [ ] Select channel to render
      - 1
      - 2
      - 1+2
- [ ] Add anti-aliasing to command line options
- [ ] Optimize antialias to use NGRBA.Pix and NRGBA.Stride values
- [ ] Platform executables
- [ ] README

- [x] Add grid to command line options
- [x] Convert PCM16 to float32 symmetrically
- [x] from/to
- [x] float32 WAV file
- [x] Grid spec
- [x] Set background colour
- [x] Move WAV decoding to cmd
- [x] Replace WAV decoder (cf. https://github.com/go-audio/wav/issues/8)
- [x] FillSpec
- [x] GridSpec
- [x] Split fill/grid
- [x] Fix offsets for odd heights
- [x] Negative height to invert
- [x] Normalize to 16 bit PCM
- [x] Anti-alias full image
- [x] Check anti-aliasing

### TODO

1. Use different 'flag' package with separate -- and - options
2. GPU/WebGL implementation
3. Add support for other audio formats
4. SVG out
5. Spectrogram
6. MP4

### NOTES

1. binpac for WAV parser ?
2. mmap
