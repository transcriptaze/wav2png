# TODO

- [ ] Columns renderer
      - [x] Refactor default renderer out as 'lines'
      - [x] Delegate grid, background etc to compositor
      - [x] --style <JSON> argument
      - [x] columns.json
      - [x] wav2mp4
      - [x] Clean up wav2png.options
      - [x] Move cursors to own package
      - [x] Restructure palettes
      - [ ] Load palette from PNG file
      - [ ] Stretch palette horizontally

- [ ] Optimise antialiasing
- [ ] Gallery
      - [ ] black and white

- [ ] wav2mp4 does not seperate L/R audio with --mix option
- [ ] wav2mp4 does not trim audio to size either

- [ ] [smoothstep cursor](https://iquilezles.org/www/articles/smoothstepintegral/smoothstepintegral.htm)
- [ ] [exponential step cursor](https://iquilezles.org/www/articles/functions/functions.htm)
- [ ] https://sound.stackexchange.com/questions/52076/whats-the-terminology-of-these-two-waveforms

### Future

1. WAV decoding for large files (mmap ?)
2. Add support for other audio formats
3. GPU/WebGL implementation
4. SVG out
5. Spectrogram
6. Load anti-aliasing kernel from file
7. Use different 'flag' package with separate -- and - options
8. https://softwarerecs.stackexchange.com/questions/87339/audio-to-waveform-images
9. https://graphicriver.net/item/waveform-artist-mp3-to-waveform-poster/20644757


### NOTES

1. binpac for WAV parser ?
2. mmap
