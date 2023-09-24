# TODO

- [ ] Quantize audio on pixel boundaries (so that it doesn't change when shifting left and right)
- [ ] Distributable Go app

- [x] Github build workflow
- [x] Clip to grid (e.g. noise)
- [x] Two point gradient
      - (?) Draggable endpoints
- [ ] Three point gradient
      - [ ] Draggable midpoint
- [ ] Anti-aliasing
      - https://webkit.org/demos/webgpu/compute-blur.html

- [ ] Scrubbable left/right/plus/minus
- [ ] Render in Worker
- [ ] Restyle range controls       
- [ ] Fix Firefox flickering
- [ ] Futhark
      - https://futhark-lang.org/index.html
      
- [ ]  Better colour picker
       - https://rgbacolorpicker.com/


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

- https://web.dev/streams/
- https://rome.tools/

## NOTES

1. https://thebookofshaders.com/07/
2. https://github.com/AmesingFlank/taichi.js
3. https://eliemichel.github.io/LearnWebGPU/
4. https://en.wikipedia.org/wiki/Trochoidal_wave
5. https://gist.github.com/munrocket/30e645d584b5300ee69295e54674b3e4
6. https://wrfranklin.org/Research/Short_Notes/pnpoly.html