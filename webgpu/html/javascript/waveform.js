import { rgba, lightblue } from './colours.js'
import { line } from './shaders/line.js'
import { gradient } from './shaders/gradient.js'

export function waveform (context, device, format, samples, style) {
  const width = context.canvas.width
  const height = context.canvas.height
  const { type = 'line' } = style

  // ... line?
  if (type === 'line') {
    const { vscale = '1.0', colour = '#80ccffff' } = style.line

    return line(device, format, samples, width, height, vscale, rgba(colour))
  }

  // ... gradient?
  if (type === 'gradient') {
    const { vscale = '1.0', colours = ['#80ccffff', '#80ccff40'] } = style.gradient

    return gradient(device, format, samples, width, height, vscale, rgba(colours[0]), rgba(colours[1]))
  }

  // ... default
  return line(device, format, samples, width, height, 1.0, lightblue)
}
