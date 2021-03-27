## v0.0

### IN PROGRESS

- [x] Split fill/grid
- [ ] GridSpec
      - fit to border
      - spacing 
      - lines
      - colour
- [ ] Move WAV decoding to cmd


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

- [x] Fix grid size
- [x] Fix offsets for odd heights
- [x] Negative height to invert
- [x] Normalize to 16 bit PCM
- [x] Anti-alias full image
- [x] Check anti-aliasing

### TODO

1. Remove external decoder dependency
2. Use different 'flag' package with separate -- and - options
3. Online renderer
4. GPU/WebGL implementation
5. Add support for other audio formats
6. SVG out
7. Spectrogram
8. MP4

### NOTES

1. binpac for WAV parser ?
2. mmap
