import { rgba, lightblue } from './colours.js'
import { line } from './shaders/line.js'
import { gradient } from './shaders/gradient.js'
import { gradient3 } from './shaders/gradient3.js'

export function waveform (context, device, format, audio, style) {
  const width = context.canvas.width
  const height = context.canvas.height

  const duration = audio.duration
  const start = duration === 0 ? 0 : Math.floor(audio.audio.length * audio.start / duration)
  const end = duration === 0 ? 0 : Math.floor(audio.audio.length * audio.end / duration)
  const samples = audio.audio.subarray(start, end)

  // ... line?
  if (style.type === 'line') {
    const {
      vscale = '1.0',
      colour = '#80ccffff'
    } = style.line

    return line(device, format, samples, width, height, vscale, rgba(colour))
  }

  // ... gradient?
  if (style.type === 'gradient') {
    const {
      vscale = '1.0',
      colours = ['#80ccffff', '#80ccff80']
    } = style.gradient

    return gradient(device, format, samples, width, height, vscale, rgba(colours[0]), rgba(colours[1]))
  }

  // ... gradient3?
  if (style.type === 'gradient3') {
    const {
      vscale = '1.0',
      midpoint = 0.5,
      colours = ['#80ccffff', '#80ccff40', '#80ccff80']
    } = style.gradient3

    return gradient3(device, format, samples, width, height, vscale, midpoint, rgba(colours[0]), rgba(colours[1]), rgba(colours[2]))
  }

  // ... default
  return line(device, format, samples, width, height, 1.0, lightblue)
}
