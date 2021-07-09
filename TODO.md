## v1.x

### IN PROGRESS

- [ ] wav2mp4
      - [x] window
      - [x] FPS
      - [x] cursor
      - [ ] cursor - centre
            - fix offset so that it doesn't discard first and last left/right bits
      - [x] cursor - left
      - [ ] cursor - right
            - fix offset so that it doesn't discard first and last left/right bits
      - [ ] cursor - ease-in/out
      - [ ] cursor - PNG file
      - [ ] cursor - stretch middle
      - [ ] README
- [ ] Copy wav2mp4 render fixes across to wav2png
- [ ] Rebuild webapp with wav2png fixes

### TODO

1. WAV decoding for large files (mmap ?)
2. Add support for other audio formats
3. GPU/WebGL implementation
4. SVG out
5. Spectrogram
6. Load anti-aliasing kernel from file
7. Use different 'flag' package with separate -- and - options

### NOTES

1. binpac for WAV parser ?
2. mmap
