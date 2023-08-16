import { Drag } from './drag.js'

export class Overlay extends HTMLElement {
  static get observedAttributes () {
    return ['width', 'height', 'padding']
  }

  constructor () {
    super()

    this.internal = {
      duration: 2,
      padding: 20,
      start: 0,
      end: 1920 - 2 * 20,

      handles: {
        left: new Drag(this, () => getStartXY(this), (x, y, dragging) => setStartXY(this, x, y, dragging)),
        right: new Drag(this, () => getEndXY(this), (x, y, dragging) => setEndXY(this, x, y, dragging))
      },

      onChange: null,
      onChanged: null
    }

    const template = document.querySelector('#template-overlay')
    const stylesheet = document.createElement('link')
    const content = template.content
    const shadow = this.attachShadow({ mode: 'open' })
    const clone = content.cloneNode(true)

    stylesheet.setAttribute('rel', 'stylesheet')
    stylesheet.setAttribute('href', '/css/components.css')

    shadow.appendChild(stylesheet)
    shadow.appendChild(clone)
  }

  connectedCallback () {
    const shadow = this.shadowRoot
    const canvas = shadow.querySelector('canvas')
    const handles = this.internal.handles

    canvas.onpointerdown = (event) => onPointerDown(this, canvas, handles, event)

    redraw(this, canvas)
  }

  disconnectedCallback () {
  }

  adoptedCallback () {
  }

  attributeChangedCallback (name, from, to) {
    if (name === 'width') {
      this.width = to
    }

    if (name === 'height') {
      this.height = to
    }

    if (name === 'padding') {
      this.padding = to
    }
  }

  get onchange () {
    return this.internal.onChange
  }

  set onchange (v) {
    this.internal.onChange = v
  }

  get onchanged () {
    return this.internal.onChanged
  }

  set onchanged (v) {
    this.internal.onChanged = v
  }

  get width () {
    const shadow = this.shadowRoot
    const canvas = shadow.querySelector('canvas')

    return canvas.width
  }

  set width (v) {
    const shadow = this.shadowRoot
    const canvas = shadow.querySelector('canvas')
    const w = Number.parseInt(`${v}`)

    if (!Number.isNaN(w)) {
      canvas.width = w
    }
  }

  get height () {
    const shadow = this.shadowRoot
    const canvas = shadow.querySelector('canvas')

    return canvas.height
  }

  set height (v) {
    const shadow = this.shadowRoot
    const canvas = shadow.querySelector('canvas')
    const w = Number.parseInt(`${v}`)

    if (!Number.isNaN(w)) {
      canvas.height = w
    }
  }

  get padding () {
    return this.internal.padding
  }

  set padding (v) {
    const p = Number.parseInt(`${v}`)

    if (!Number.isNaN(p)) {
      this.internal.padding = p
    }
  }

  get duration () {
    return this.internal.duration
  }

  set duration (v) {
    const duration = Number.parseFloat(`${v}`)

    if (!Number.isNaN(duration)) {
      this.internal.duration = duration
    }
  }

  get start () {
    return this.internal.start
  }

  set start (v) {
    const shadow = this.shadowRoot
    const canvas = shadow.querySelector('canvas')
    const start = Number.parseInt(`${v}`)

    if (!Number.isNaN(start)) {
      this.internal.start = Math.min(Math.max(start, 0), this.end)
      redraw(this, canvas)
    }
  }

  get end () {
    return this.internal.end
  }

  set end (v) {
    const shadow = this.shadowRoot
    const canvas = shadow.querySelector('canvas')
    const end = Number.parseInt(`${v}`)

    if (!Number.isNaN(end)) {
      this.internal.end = Math.max(Math.min(end, 1920 - 2 * 20), this.start)
      redraw(this, canvas)
    }
  }

  reset () {
    const shadow = this.shadowRoot
    const canvas = shadow.querySelector('canvas')

    this.start = 0
    this.end = this.width - 2 * this.padding

    redraw(this, canvas)
  }

  format (v) {
    const ms = Math.trunc(v * 1000) % 1000
    const ss = Math.trunc(v) % 60
    const mm = Math.trunc(v / 60)

    if (mm > 0) {
      return `${mm}:${ss}.${ms.toString().padStart(3, '0')}`
    } else {
      return `${ss}.${ms.toString().padStart(3, '0')}`
    }
  }
}

function getStartXY (overlay) {
  return { x: overlay.start, y: 0 }
}

function setStartXY (overlay, x, y, dragging) {
  overlay.start = x

  if (dragging) {
    onChange(overlay)
  } else {
    onChanged(overlay)
  }
}

function getEndXY (overlay) {
  return { x: overlay.end, y: 0 }
}

function setEndXY (overlay, x, y, dragging) {
  overlay.end = x

  if (dragging) {
    onChange(overlay)
  } else {
    onChanged(overlay)
  }
}

function onChange (overlay) {
  if (overlay.onchange != null) {
    const s = overlay.duration * overlay.start / (1920 - 40)
    const t = overlay.duration * overlay.end / (1920 - 40)

    overlay.onchange(s, t)
  }
}

function onChanged (overlay) {
  if (overlay.onchanged != null) {
    const s = overlay.duration * overlay.start / (1920 - 40)
    const t = overlay.duration * overlay.end / (1920 - 40)

    overlay.onchanged(s, t)
  }
}

function redraw (component, canvas) {
  const width = canvas.width
  const height = canvas.height
  const padding = component.padding
  const start = component.start
  const end = Math.min(width - 2 * padding, component.end)
  const ctx = canvas.getContext('2d')

  ctx.clearRect(0, 0, width, height)

  // ... draw sizing handles
  ctx.fillStyle = '#80ccff80'
  ctx.strokeStyle = '#80ccff80'
  ctx.lineWidth = 2

  ctx.beginPath()
  ctx.moveTo(0, 0)
  ctx.lineTo(padding + start, 0)
  ctx.lineTo(padding + start, height)
  ctx.lineTo(0, height)
  ctx.fill()

  ctx.beginPath()
  ctx.moveTo(padding + start - 32, height / 2)
  ctx.lineTo(padding + start, height / 2 - 32)
  ctx.lineTo(padding + start, height / 2 + 32)
  ctx.fill()

  ctx.beginPath()
  ctx.moveTo(width, 0)
  ctx.lineTo(padding + end, 0)
  ctx.lineTo(padding + end, height)
  ctx.lineTo(width, height)
  ctx.fill()

  ctx.beginPath()
  ctx.moveTo(padding + end + 32, height / 2)
  ctx.lineTo(padding + end, height / 2 - 32)
  ctx.lineTo(padding + end, height / 2 + 32)
  ctx.fill()

  ctx.beginPath()
  ctx.moveTo(padding + start, 0)
  ctx.lineTo(padding + start, height)
  ctx.stroke()

  ctx.beginPath()
  ctx.moveTo(padding + end, 0)
  ctx.lineTo(padding + end, height)
  ctx.stroke()

  // ... labels
  ctx.font = '18px sans-serif'
  ctx.fillStyle = '#80ccffc0'
  ctx.strokeStyle = '#80ccffc0'

  const labels = {
    start: component.format(component.duration * component.start / component.width),
    end: component.format(component.duration * component.end / component.width),
    duration: component.format(component.duration * (component.end - component.start) / component.width)
  }

  const fm = {
    start: ctx.measureText(labels.start),
    end: ctx.measureText(labels.end),
    duration: ctx.measureText(labels.duration)
  }

  const w = 8 + fm.start.width + 8 + fm.end.width + 8
  const middle = (start + end) / 2

  const x = {
    start: Math.min(padding + start + 8, padding + middle - w / 2),
    end: Math.max(padding + end - 8, padding + middle + w / 2)
  }

  ctx.textAlign = 'left'
  ctx.fillText(labels.start, x.start, 0 + fm.start.fontBoundingBoxAscent)

  ctx.textAlign = 'right'
  ctx.fillText(labels.end, x.end, 0 + fm.end.fontBoundingBoxAscent)

  ctx.textAlign = 'center'
  ctx.fillText(labels.duration, padding + middle, height - fm.duration.fontBoundingBoxDescent)

  const dl = middle - start - fm.duration.width / 2 - 16
  const dr = middle - end + fm.duration.width / 2 + 16

  if (dl > 0 && dr < 0) {
    ctx.beginPath()
    ctx.moveTo(padding + start + 8, height - fm.duration.fontBoundingBoxAscent / 2)
    ctx.lineTo(padding + start + 8 + dl, height - fm.duration.fontBoundingBoxAscent / 2)
    ctx.moveTo(padding + end - 8, height - fm.duration.fontBoundingBoxAscent / 2)
    ctx.lineTo(padding + end - 8 + dr, height - fm.duration.fontBoundingBoxAscent / 2)
    ctx.stroke()
  }
}

function onPointerDown (component, canvas, handles, event) {
  const width = canvas.width
  const height = canvas.height
  const padding = component.padding

  if (event.button === 0) {
    event.preventDefault()

    const hscale = 2 // FIXME calculate from client width
    const vscale = 2 // FIXME calculate from client height
    const xy = { x: hscale * event.offsetX, y: vscale * event.offsetY }
    const start = Math.max(component.start, 0)
    const end = Math.min(component.end, width - 2 * padding)

    const left = [
      { x: padding + start, y: 0 },
      { x: padding + start, y: height },
      { x: padding + start - 48, y: height },
      { x: padding + start - 48, y: 0 }
    ]

    const right = [
      { x: padding + end, y: 0 },
      { x: padding + end, y: height },
      { x: padding + end + 48, y: height },
      { x: padding + end + 48, y: 0 }
    ]

    if (hittest(xy, left)) {
      handles.left.start(event, canvas)
    } else if (hittest(xy, right)) {
      handles.right.start(event, canvas)
    }
  }
}

// Ref. https://wrfranklin.org/Research/Short_Notes/pnpoly.html
function hittest (xy, polygon) {
  const N = polygon.length
  let hit = false

  for (let i = 0, j = N - 1; i < N; j = i++) {
    if (((polygon[i].y > xy.y) !== (polygon[j].y > xy.y)) &&
       (xy.x < (polygon[j].x - polygon[i].x) * (xy.y - polygon[i].y) / (polygon[j].y - polygon[i].y) + polygon[i].x)) {
      hit = !hit
    }
  }

  return hit
}

customElements.define('wav2png-overlay', Overlay)
