import { background } from './background.js'
import { grid } from './grid.js'
import { waveform } from './waveform.js'
import { black, green, transparent } from './colours.js'

class Offscreen {
  constructor () {
    const canvas = document.querySelector('#canvas canvas')
    const width = canvas.width
    const height = canvas.height

    this.internal = {
      device: null,
      audio: new Float32Array(),
      canvas: new OffscreenCanvas(width, height),
      fill: black,
      grid: {
        colour: green
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

  set audio (v) {
    this.internal.audio = v
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

  get canvas () {
    return this.internal.canvas
  }

  render () {
    const ctx = this.canvas.getContext('webgpu')
    const device = this.device
    const audio = this.audio
    const format = navigator.gpu.getPreferredCanvasFormat()
    const layers = []

    ctx.configure({ device: this.device, format, alphaMode: 'premultiplied' })

    layers.push(background(ctx, device, format, this.fill))
    layers.push(grid(ctx, device, format, this.grid))

    if (audio.length > 0) {
      layers.push(waveform(ctx, device, format, audio))
    }

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
