import { background } from './background.js'
import { grid } from './grid.js'
import { waveform } from './waveform.js'
import { black } from './colours.js'

class Overview {
  constructor () {
    this.internal = {
      device: null,
      audio: new Float32Array(),
      canvas: document.querySelector('#overview canvas')
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

  set audio (v) {
    this.internal.audio = v
    this.redraw()
  }

  get canvas () {
    return this.internal.canvas
  }

  redraw () {
    const ctx = this.canvas.getContext('webgpu')
    const device = this.device
    const audio = this.audio
    const format = navigator.gpu.getPreferredCanvasFormat()
    const layers = []

    ctx.configure({ device: this.device, format })

    layers.push(background(ctx, device, format, black))
    layers.push(grid(ctx, device, format))

    if (audio.length > 0) {
      layers.push(waveform(ctx, device, format, audio))
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
