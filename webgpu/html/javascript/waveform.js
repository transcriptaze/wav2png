import { rgba, lightblue } from './colours.js'
import { line } from './shaders/line.js'
import { gradient } from './shaders/gradient.js'
import { gradient3 } from './shaders/gradient3.js'

export function waveform (context, device, format, audio, style) {
  const fs = Number.isNaN(audio.fs) ? 44100 : audio.fs
  const N = audio.audio.length
  const duration = clamp(audio.duration, 0, N / fs)
  const start = duration === 0 ? 0 : clamp(Math.floor(N * audio.start / duration), 0, N)
  const end = duration === 0 ? 0 : clamp(Math.floor(N * audio.end / duration), 0, N)

  if (duration <= 0 || end <= start) {
    return {}
  }

  const width = context.canvas.width
  const height = context.canvas.height
  const samples = audio.audio

  // ... line?
  if (style.type === 'line') {
    const {
      vscale = '1.0',
      colour = '#80ccffff'
    } = style.line

    return line(device, format, { start, end, audio: samples }, width, height, vscale, rgba(colour))
  }

  // ... gradient?
  if (style.type === 'gradient') {
    const {
      vscale = '1.0',
      colours = ['#80ccffff', '#80ccff80']
    } = style.gradient

    return gradient(device, format, { start, end, audio: samples }, width, height, vscale, rgba(colours[0]), rgba(colours[1]))
  }

  // ... gradient3?
  if (style.type === 'gradient3') {
    const {
      vscale = '1.0',
      colours = ['#80ccffff', '#80ccff40', '#80ccff80'],
      stops = [0.0, 0.5, 1.0]
    } = style.gradient3

    return gradient3(device, format, { start, end, audio: samples }, width, height, vscale, rgba(colours[0]), rgba(colours[1]), rgba(colours[2]), stops[1])
  }

  // ... default
  return line(device, format, { start, end, audio: samples }, width, height, 1.0, lightblue)
}

function clamp (v, min, max) {
  return (v < min) ? min : ((v > max) ? max : v)
}
