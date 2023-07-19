# TODO

- [x] Download image
- [x] Textured page background
- [x] Logo
- [x] Render audio
      - [x] 'copy' audio to waveform
      - [x] Accumulate and average audio
      - [x] Edge case: stride x pixels > samples.length
      - [x] Pack constants properly i.e. not using all float32.
      - [x] Render real audio
      - [x] Accept any transcodable audio format

- [ ] Set background colour
- [ ] Distributable Go app
- [ ] :hover for drag (with audio file only)
- [ ] lines:minmax
- [ ] lines:median
- [ ] lines:multi-point gradient
- [ ] lines:pixelated
- (?) Create different bind group layouts for compute/render shaders to avoid confusion
- (?) Experiment with border image for windmill
      - https://developer.mozilla.org/en-US/docs/Web/CSS/border-image

## NOTES

1. https://thebookofshaders.com/07/
2. https://github.com/AmesingFlank/taichi.js
3. https://eliemichel.github.io/LearnWebGPU/
4. https://en.wikipedia.org/wiki/Trochoidal_wave
5. https://gist.github.com/munrocket/30e645d584b5300ee69295e54674b3e4