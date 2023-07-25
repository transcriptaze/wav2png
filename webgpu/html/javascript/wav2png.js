import { background } from './background.js'
import { grid } from './grid.js'
import { waveform } from './waveform.js'

const transparent = rgba('#00000000')

const context = {
  device: null,
  loading: false,
  loaded: false,
  audio: null,
  fill: '#000000ff'
}

export async function initialise () {
  const canvas = document.querySelector('#canvas canvas')

  if (!navigator.gpu) {
    throw new Error('WebGPU not supported on this browser.')
  }

  const adapter = await navigator.gpu.requestAdapter()
  if (!adapter) {
    throw new Error('No appropriate GPUAdapter found.')
  }

  context.device = await adapter.requestDevice()
  context.audio = new Float32Array()

  render(context.device, canvas, context.audio)
}

export function load (filename, blob) {
  const device = context.device
  const canvas = document.querySelector('#canvas canvas')
  const save = document.getElementById('save')
  const fill = document.getElementById('fill')
  const clear = document.getElementById('clear')

  save.disabled = true
  fill.disabled = true
  clear.disabled = true

  context.loading = true
  context.loaded = false
  context.audio = new Float32Array()

  busy()
    .then(b => blob.arrayBuffer())
    .then(b => transcode(b))
    .then(b => b.getChannelData(0))
    .then(audio => {
      render(device, canvas, audio)
      return audio
    })
    .then(audio => {
      context.audio = audio
      context.loading = false
      context.loaded = true

      save.disabled = false
      fill.disabled = false
      clear.disabled = false
    })
    .catch((err) => {
      console.error(err)
    })
    .finally(unbusy)
}

// NTS: canvas.toBlob only works if invoked immediately after redraw (apparently needs to happen
//      before unconfigure() or getCurrentTexture() are invoked ?). Workaround is to render it
//      to an offscreen buffer for download.
export function download () {
  const canvas = document.querySelector('#canvas canvas')
  const width = canvas.width
  const height = canvas.height
  const offscreen = new OffscreenCanvas(width, height)
  const device = context.device
  const audio = context.audio

  render(device, offscreen, audio)

  offscreen.convertToBlob().then((blob) => {
    save(blob, false)
  })
}

export function trash () {
  const device = context.device
  const canvas = document.querySelector('#canvas canvas')
  const save = document.getElementById('save')
  const fill = document.getElementById('fill')
  const clear = document.getElementById('clear')

  context.loading = false
  context.loaded = false
  context.audio = new Float32Array()

  save.disabled = true
  fill.disabled = true
  clear.disabled = true

  render(device, canvas, context.audio)
}

export function fill (colour) {
  context.fill = colour

  const canvas = document.querySelector('#canvas canvas')

  render(context.device, canvas, context.audio)
}

async function transcode (bytes) {
  const AudioContext = window.AudioContext || window.webkitAudioContext
  const ctx = new AudioContext()
  const buffer = await ctx.decodeAudioData(bytes)
  const offline = new OfflineAudioContext(1, 44100 * buffer.duration, 44100)
  const src = offline.createBufferSource()

  src.buffer = buffer
  src.connect(offline.destination)
  src.start()

  return offline.startRendering()
}

function render (device, canvas, audio) {
  const ctx = canvas.getContext('webgpu')
  const format = navigator.gpu.getPreferredCanvasFormat()

  ctx.configure({ device, format })

  const layers = []

  layers.push(background(ctx, device, format, rgba(context.fill)))
  layers.push(grid(ctx, device, format))

  if (audio.length > 0) {
    layers.push(waveform(ctx, device, format, audio))
  }

  draw(ctx, device, layers)
}

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

function save (blob, timestamp) {
  // if (context.loaded) {
  const now = new Date()
  const year = `${now.getFullYear()}`.padStart(4, '0')
  const month = `${now.getMonth() + 1}`.padStart(2, '0')
  const day = `${now.getDate()}`.padStart(2, '0')
  const hour = `${now.getHours()}`.padStart(2, '0')
  const minute = `${now.getMinutes()}`.padStart(2, '0')
  const second = `${now.getSeconds()}`.padStart(2, '0')
  const filename = timestamp ? `wav2png ${year}-${month}-${day} ${hour}.${minute}.${second}.png` : 'wav2png.png'

  if (window.showSaveFilePicker) {
    saveWithPicker(blob, filename)
  } else {
    const url = URL.createObjectURL(blob)
    const anchor = document.getElementById('download')

    anchor.href = url
    anchor.download = 'wav2png.png'
    anchor.click()

    URL.revokeObjectURL(url)
  }
  // }
}

async function saveWithPicker (blob, filename) {
  try {
    const options = {
      suggestedName: filename,
      types: [
        {
          description: 'wav2png PNG file',
          accept: { 'image/png': ['.png'] }
        }
      ]
    }

    const handle = await window.showSaveFilePicker(options)
    const stream = await handle.createWritable()

    await stream.write(blob)
    await stream.close()
  } catch (err) {
    if (err.name !== 'AbortError') {
      console.error(err)
    }
  }
}

function busy () {
  const windmill = document.getElementById('windmill')

  return new Promise((resolve) => {
    if (context.loading) {
      windmill.classList.add('visible')
    }

    // A 100ms delay let things like radio buttons get updated and the windmill displaying
    // before the redraw is complete
    setTimeout(resolve, 100)
  })
}

function unbusy () {
  const windmill = document.getElementById('windmill')

  windmill.classList.remove('visible')
}

function rgba (colour) {
  const match = `${colour}`.match(/^#([a-fA-F0-9]+)$/)

  if (match && match.length > 1) {
    const hex = match[1]
    const v = Number.parseInt(hex, 16)
    const r = (v >>> 24) & 0x00ff
    const g = (v >>> 16) & 0x00ff
    const b = (v >>> 8) & 0x00ff
    const a = (v >>> 0) & 0x00ff

    return [r / 255, g / 255, b / 255, a / 255]
  }

  return [0, 0, 0, 0]
}
