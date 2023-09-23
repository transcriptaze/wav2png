import { background } from './background.js'
import { grid } from './grid.js'
import { waveform } from './waveform.js'
import { black, rgba } from './colours.js'

class Overview {
  constructor () {
    this.internal = {
      device: null,
      audio: new Float32Array(),
      canvas: document.querySelector('#overview canvas'),
      overlay: document.querySelector('#overview wav2png-overlay'),

      fill: black,

      grid: {
        colour: rgba('#ffff00a0'),
        gridx: 64,
        gridy: 4
      },

      waveform: {
        type: 'line',
        line: {
          vscale: 1.0,
          colour: '#80ccffff'
        }
      }
    }

    this.internal.overlay.onchanged = (start, end) => {
      if (this.internal.onChange != null) {
        this.internal.onChange(start, end)
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
    this.internal.audio = audio
    this.internal.overlay.audio = { start: 0, end: audio.length / fs, duration: audio.length / fs }

    this.redraw()
  }

  get fill () {
    return this.internal.fill
  }

  set fill (v) {
    this.internal.fill = v
    this.redraw()
  }

  get grid () {
    return this.internal.grid
  }

  set grid ({ colour }) {
    this.internal.grid.colour = colour
    this.redraw()
  }

  get waveform () {
    return this.internal.waveform
  }

  set waveform ({ vscale }) {
    this.internal.waveform.line.vscale = vscale
    this.redraw()
  }

  get canvas () {
    return this.internal.canvas
  }

  /* eslint-disable-next-line accessor-pairs */
  set start (v) {
    this.internal.overlay.start = v
  }

  /* eslint-disable-next-line accessor-pairs */
  set end (v) {
    this.internal.overlay.end = v
  }

  /* eslint-disable-next-line accessor-pairs */
  set onchange (v) {
    this.internal.overlay.onchange = v
  }

  /* eslint-disable-next-line accessor-pairs */
  set onchanged (v) {
    this.internal.overlay.onchanged = v
  }

  redraw () {
    const ctx = this.canvas.getContext('webgpu')
    const device = this.device
    const audio = this.audio
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

export const overview = new Overview()

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
        clearValue: { r: 0, g: 0, b: 0, a: 1 }
      }]
    })

    for (const layer of layers) {
      layer.render(pass)
    }

    pass.end()
  }

  device.queue.submit([encoder.finish()])
}
