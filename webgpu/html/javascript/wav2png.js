import { overview } from './overview.js'
import { canvas } from './canvas.js'
import { offscreen } from './offscreen.js'
import { rgba } from './colours.js'

const context = {
  loading: false,
  loaded: false
}

let stale = true
let duration = 0

export async function initialise () {
  if (!navigator.gpu) {
    console.error(new Error('WebGPU not supported on this browser'))
    window.location = '/unsupported.html'
    return
  }

  const adapter = await navigator.gpu.requestAdapter()
  if (!adapter) {
    console.error(new Error('No appropriate GPUAdapter found'))
    window.location = '/unsupported.html'
    return
  }

  const device = await adapter.requestDevice()

  overview.device = device
  canvas.device = device
  offscreen.device = device

  overview.redraw()
  canvas.redraw()

  const fill = document.getElementById('fill')
  const grid = document.getElementById('grid')
  const waveform = document.getElementById('waveform')
  const xaxis = document.getElementById('x-axis')

  fill.onchange = (c) => {
    stale = true
  }

  fill.onchanged = (c) => {
    canvas.fill = rgba(c)
    offscreen.fill = rgba(c)
    canvas.redraw()
  }

  grid.onchange = (c) => {
    stale = true
  }

  grid.onchanged = (c) => {
    canvas.grid = { colour: rgba(c) }
    offscreen.grid = { colour: rgba(c) }
    canvas.redraw()
  }

  waveform.onchange = (c) => {
    stale = true
  }

  waveform.onchanged = (w) => {
    canvas.waveform = w
    offscreen.waveform = w
    canvas.redraw()
  }

  overview.onchanged = (start, end) => {
    xaxis.start = start
    xaxis.end = end

    canvas.start = start
    canvas.end = end
    canvas.redraw()

    offscreen.start = start
    offscreen.end = end
  }

  xaxis.onchanged = (start, end) => {
    overview.start = start
    overview.end = end

    canvas.start = start
    canvas.end = end
    canvas.redraw()

    offscreen.start = start
    offscreen.end = end
  }

  const refresh = () => {
    if (stale) {
      redraw()
    }

    window.requestAnimationFrame(refresh)
  }

  window.requestAnimationFrame(refresh)
}

export function load (filename, blob) {
  const save = document.getElementById('save')
  const fill = document.getElementById('fill')
  const clear = document.getElementById('clear')
  const xaxis = document.getElementById('x-axis')

  save.disabled = true
  fill.disabled = true
  clear.disabled = true

  context.loading = true
  context.loaded = false

  busy()
    .then(b => blob.arrayBuffer())
    .then(b => transcode(b))
    .then(b => {
      return { fs: b.sampleRate, audio: b.getChannelData(0) }
    })
    .then(({ fs, audio }) => {
      context.loading = false
      context.loaded = true

      duration = audio.length / fs

      overview.audio = { fs, audio }
      canvas.audio = { fs, audio }
      offscreen.audio = { fs, audio }

      xaxis.audio = { start: 0, end: duration, duration }

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
  offscreen.render().then((blob) => {
    save(blob, false)
  })
}

export function trash () {
  const save = document.getElementById('save')
  const fill = document.getElementById('fill')
  const clear = document.getElementById('clear')
  const xaxis = document.getElementById('x-axis')

  context.loading = false
  context.loaded = false

  duration = 0

  overview.audio = { fs: 44100, audio: new Float32Array() }
  canvas.audio = { fs: 44100, audio: new Float32Array() }
  offscreen.audio = { fs: 44100, audio: new Float32Array() }
  xaxis.audio = { start: 0, end: 0, duration: 0 }

  save.disabled = true
  fill.disabled = true
  clear.disabled = true
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

function save (blob, timestamp) {
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

function redraw () {
  stale = false

  return new Promise(() => {
    const fill = document.getElementById('fill').colour
    const grid = document.getElementById('grid').colour
    const waveform = document.getElementById('waveform').style

    canvas.fill = rgba(fill)
    canvas.grid = { colour: rgba(grid) }
    canvas.waveform = waveform

    offscreen.fill = rgba(fill)
    offscreen.grid = { colour: rgba(grid) }
    offscreen.waveform = waveform

    canvas.redraw()
  })
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
