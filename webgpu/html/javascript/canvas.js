import { background } from './background.js'
import { grid } from './grid.js'
import { waveform } from './waveform2.js'
import { black, green, transparent, rgba } from './colours.js'

class Canvas {
  constructor () {
    this.internal = {
      device: null,
      canvas: document.querySelector('#canvas canvas'),

      fill: black,

      grid: {
        colour: green,
        gridx: 8,
        gridy: 4
      },

      waveform: {
        vscale: 1.0,
        colours: [
          rgba('#80ccffff'),
          rgba('#80ccff40')
        ]
      },

      audio: new Float32Array(),
      fs: 44100,
      duration: 0,
      start: 0,
      end: 0
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
    this.internal.fs = fs
    this.internal.audio = audio
    this.internal.duration = audio.length / fs
    this.internal.start = 0
    this.internal.end = this.internal.duration

    this.redraw()
  }

  get fill () {
    return this.internal.fill
  }

  set fill (v) {
    this.internal.fill = v
  }

  get grid () {
    return this.internal.grid
  }

  set grid ({ colour }) {
    this.internal.grid.colour = colour
  }

  get waveform () {
    return this.internal.waveform
  }

  set waveform ({ vscale, colours }) {
    this.internal.waveform.vscale = vscale
    this.internal.waveform.colours = colours
  }

  /* eslint-disable-next-line accessor-pairs */
  set start (v) {
    const start = Number.parseFloat(`${v}`)

    if (!Number.isNaN(start)) {
      this.internal.start = constrain(start, 0, this.internal.duration)
    }
  }

  /* eslint-disable-next-line accessor-pairs */
  set end (v) {
    const end = Number.parseFloat(`${v}`)

    if (!Number.isNaN(end)) {
      this.internal.end = constrain(end, 0, this.internal.duration)
    }
  }

  get canvas () {
    return this.internal.canvas
  }

  redraw () {
    const duration = this.internal.duration
    const start = duration === 0 ? 0 : Math.floor(this.audio.length * this.internal.start / duration)
    const end = duration === 0 ? 0 : Math.floor(this.audio.length * this.internal.end / duration)

    const ctx = this.canvas.getContext('webgpu')
    const device = this.device
    const audio = this.audio.subarray(start, end)
    const format = navigator.gpu.getPreferredCanvasFormat()
    const layers = []

    ctx.configure({ device: this.device, format, alphaMode: 'premultiplied' })

    layers.push(background(ctx, device, format, this.fill))
    layers.push(grid(ctx, device, format, this.grid))

    if (audio.length > 0) {
      layers.push(waveform(ctx, device, format, audio, this.waveform))
    }

    draw(ctx, this.device, layers)
  }
}

export const canvas = new Canvas()

function draw (context, device, layers) {
  const encoder = device.createCommandEncoder()

  {
    const pass = encoder.beginComputePass()
    for (const layer of layers) {
      layer.compute(pass)
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
      layer.render(pass)
    }

    pass.end()
  }

  device.queue.submit([encoder.finish()])
}

function constrain (v, min, max) {
  return Math.max(Math.min(v, max), min)
}
