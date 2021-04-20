## v0.0

### IN PROGRESS

- [ ] Grid spec
      - none/square/rectangular
      - 'smart' e.g. flip between no. of lines/size
      - colour
      - overlay/underlay
      - return array of intercepts

- [ ] Render should take a reader parameter
- [ ] Move WAV decoding to cmd

- [ ] Grid size
      -- horizontal lines start 1 pixel above/below 0
      -- size to number of lines and/or grid size
- [ ] Decoder for extended format/float32 WAV
      -- https://github.com/go-audio/wav/issues/8

- [ ] Add references to README
- [ ] Define palette in terms of gradients
- [ ] Load palette from file
- [ ] Load anti-aliasing kernel from file
- [ ] Select channel to render
- [ ] Mix channels
- [ ] 'fire' palette
- [ ] 'green' palette
- [ ] 'amber' palette
- [ ] 'blues' palette
- [ ] 'inverted' palette
- [ ] Add anti-aliasing to command line options
- [ ] Add grid to command line options
- [ ] Optimize antialias to use NGRBA.Pix and NRGBA.Stride values
- [ ] start/end/duration command line options
- [ ] Platform executables
- [ ] Set background colour

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
