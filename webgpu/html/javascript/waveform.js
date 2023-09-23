import { rgba, lightblue } from './colours.js'
import { line } from './shaders/line.js'

export function waveform (context, device, format, samples, style) {
  const { type = 'line' } = style

  // ... line?
  if (type === 'line') {
    const { vscale = '1.0', colour = '#80ccffff' } = style.line

    return line(device, format, samples, context.canvas.width, context.canvas.height, vscale, rgba(colour))
  }

  // ... default
  return line(device, format, samples, context.canvas.width, context.canvas.height, 1.0, lightblue)
}
