import { fill } from './fill.js'
import { grid } from './grid.js'
import { waveform } from './waveform.js'
import { black, rgba } from './colours.js'

class Overview {
  constructor () {
    this.internal = {
      device: null,
      canvas: document.querySelector('#overview canvas'),
      overlay: document.querySelector('#overview wav2png-overlay'),

      audio: {
        start: 0,
        end: 0,
        duration: 0,
        fs: 44100,
        audio: new Float32Array()
      },

      styles: {
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
    this.internal.audio = {
      start: 0,
      end: audio.length / fs,
      duration: audio.length / fs,
      fs,
      audio
    }

    this.internal.overlay.audio = {
      start: 0,
      end: audio.length / fs,
      duration: audio.length / fs
    }

    this.redraw()
  }

  get styles () {
    return this.internal.styles
  }

  /* eslint-disable-next-line accessor-pairs */
  set fill (v) {
    this.internal.styles.fill = v
    this.redraw()
  }

  /* eslint-disable-next-line accessor-pairs */
  set grid ({ colour }) {
    this.internal.styles.grid.colour = colour
    this.redraw()
  }

  /* eslint-disable-next-line accessor-pairs */
  set waveform ({ vscale }) {
    this.internal.styles.waveform.line.vscale = vscale
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
    const format = navigator.gpu.getPreferredCanvasFormat()
    const audio = this.audio
    const styles = this.styles
    const layers = []

    ctx.configure({ device: this.device, format, alphaMode: 'premultiplied' })

    layers.push(fill(ctx, device, format, styles.fill))
    layers.push(waveform(ctx, device, format, audio, styles.waveform))
    layers.push(grid(ctx, device, format, styles.grid))

    draw(ctx, this.device, layers)
  }
}

export const overview = new Overview()

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
        clearValue: { r: 0, g: 0, b: 0, a: 1 }
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
