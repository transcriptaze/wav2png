import { background } from './background.js'
import { grid } from './grid.js'
import { waveform } from './waveform.js'
import { LINE } from './styles.js'
import { black, green, transparent } from './colours.js'

class Offscreen {
  constructor () {
    const canvas = document.querySelector('#canvas canvas')
    const width = canvas.width
    const height = canvas.height

    this.internal = {
      device: null,
      canvas: new OffscreenCanvas(width, height),

      styles: {
        fill: black,
        grid: {
          colour: green,
          gridx: 8,
          gridy: 4
        },
        waveform: LINE
      },

      audio: {
        start: 0,
        end: 0,
        duration: 0,
        fs: 44100,
        audio: new Float32Array()
      }
    }
  }

  get device () {
    return this.internal.device
  }

  set device (v) {
    this.internal.device = v
  }

  get audio () {
    return this.internal.audio
  }

  set audio ({ fs, audio }) {
    this.internal.audio = {
      start: 0,
      end: audio.length / fs,
      duration: audio.length / fs,
      fs,
      audio
    }
  }

  get styles () {
    return this.internal.styles
  }

  /* eslint-disable-next-line accessor-pairs */
  set fill (v) {
    this.internal.styles.fill = v
  }

  /* eslint-disable-next-line accessor-pairs */
  set grid ({ colour }) {
    this.internal.styles.grid.colour = colour
  }

  /* eslint-disable-next-line accessor-pairs */
  set waveform (v) {
    this.internal.styles.waveform = v
  }

  /* eslint-disable-next-line accessor-pairs */
  set start (v) {
    const start = Number.parseFloat(`${v}`)

    if (!Number.isNaN(start)) {
      this.audio.start = constrain(start, 0, this.audio.duration)
    }
  }

  /* eslint-disable-next-line accessor-pairs */
  set end (v) {
    const end = Number.parseFloat(`${v}`)

    if (!Number.isNaN(end)) {
      this.audio.end = constrain(end, 0, this.audio.duration)
    }
  }

  get canvas () {
    return this.internal.canvas
  }

  render () {
    const ctx = this.canvas.getContext('webgpu')
    const device = this.device
    const format = navigator.gpu.getPreferredCanvasFormat()
    const audio = this.audio
    const styles = this.styles
    const layers = []

    ctx.configure({ device: this.device, format, alphaMode: 'premultiplied' })

    layers.push(background(ctx, device, format, this.styles.fill))
    layers.push(waveform(ctx, device, format, audio, styles.waveform))
    layers.push(grid(ctx, device, format, this.styles.grid))

    draw(ctx, this.device, layers)

    return this.canvas.convertToBlob()
  }
}

export const offscreen = new Offscreen()

function draw (context, device, layers) {
  const encoder = device.createCommandEncoder()

  {
    const pass = encoder.beginComputePass()
    for (const layer of layers) {
      if (layer.compute != null) {
        layer.compute(pass)
      }
    }
    pass.end()
  }

  {
    const pass = encoder.beginRenderPass({
      colorAttachments: [{
        view: context.getCurrentTexture().createView(),
        loadOp: 'clear',
        storeOp: 'store',
        clearValue: transparent
      }]
    })

    for (const layer of layers) {
      if (layer.render != null) {
        layer.render(pass)
      }
    }

    pass.end()
  }

  device.queue.submit([encoder.finish()])
}

function constrain (v, min, max) {
  return Math.max(Math.min(v, max), min)
}
