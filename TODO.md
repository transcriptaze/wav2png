## v0.0

### IN PROGRESS

- [ ] Load palette from file
      - 'fire' palette
      - 'green' palette
      - 'amber' palette
      - 'blues' palette
      - 'inverted' palette
      
- [ ] Add grid to command line options
- [ ] Get settings from command line/file
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

1. Remove external decoder dependency
2. Use different 'flag' package with separate -- and - options
3. GPU/WebGL implementation
4. Add support for other audio formats
5. SVG out
6. Spectrogram
7. MP4

### NOTES

1. binpac for WAV parser ?
2. mmap
