# TODO

- [ ] Quantize audio on pixel boundaries (so that it doesn't change when shifting left and right)
      - [x] Calculate start/end by indexing from 0
      - [x] Set start offset in compute shader 
      - [ ] Fix jump at 1.002.2

- [ ] Draggable endpoints
      - [ ] Draggables
      - [x] Bottom half is not gradiented ... because compute shader is only generating extremums


- [ ] Distributable Go app

- [ ] Anti-aliasing
      - https://webkit.org/demos/webgpu/compute-blur.html

- [ ] Scrubbable left/right/plus/minus
- [ ] Render in Worker
- [ ] Restyle range controls
      - https://codepen.io/collection/DgYaMj/
      - https://www.smashingmagazine.com/2021/12/create-custom-range-input-consistent-browsers/
      - https://css-tricks.com/styling-cross-browser-compatible-range-inputs-css/
      - https://www.w3.org/WAI/ARIA/apg/patterns/slider/examples/slider-temperature/
      - https://codepen.io/ATC-test/pen/myPNqW

- [ ] Better colour picker
      - https://rgbacolorpicker.com/
- [ ] Fix Firefox flickering
- (?) Futhark
      - https://futhark-lang.org/index.html
      
- [ ] :hover for drag (with audio file only)
- [ ] lines:minmax
- [ ] lines:median
- [ ] lines:multi-point gradient
- [ ] lines:pixelated
- [ ] audio mixing

- (?) Create different bind group layouts for compute/render shaders to avoid confusion
- (?) Experiment with border image for windmill
      - https://developer.mozilla.org/en-US/docs/Web/CSS/border-image
      - https://web.dev/css-trig-functions/

- wav2mp4
  - https://developer.chrome.com/articles/webcodecs/

- https://web.dev/streams/
- https://rome.tools/

## NOTES

1. https://thebookofshaders.com/07/
2. https://github.com/AmesingFlank/taichi.js
3. https://eliemichel.github.io/LearnWebGPU/
4. https://en.wikipedia.org/wiki/Trochoidal_wave
5. https://gist.github.com/munrocket/30e645d584b5300ee69295e54674b3e4
6. https://wrfranklin.org/Research/Short_Notes/pnpoly.html