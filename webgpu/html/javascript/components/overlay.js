export class Overlay extends HTMLElement {
  static get observedAttributes () {
    return ['width', 'height', 'padding']
  }

  constructor () {
    super()

    this.internal = {
      duration: 60,
      padding: 20,
      start: 0,
      end: 1920 - 2 * 20,

      handles: {
        left: new Drag(this, () => this.start, (v) => { this.start = v }),
        right: new Drag(this, () => this.end, (v) => { this.end = v })
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
    const start = Number.parseInt(`${v}`)

    if (!Number.isNaN(start)) {
      this.internal.start = Math.min(Math.max(start, 0), this.end)
    }
  }

  get end () {
    return this.internal.end
  }

  set end (v) {
    const end = Number.parseInt(`${v}`)

    if (!Number.isNaN(end)) {
      this.internal.end = Math.max(Math.min(end, 1920 - 2 * 20), this.start)
    }
  }

  reset () {
    const shadow = this.shadowRoot
    const canvas = shadow.querySelector('canvas')

    this.duration = 0
    this.start = 0
    this.end = this.width - 2 * this.padding

    redraw(this, canvas)
  }
}

class Drag {
  constructor (overlay, getX, setX) {
    this.overlay = overlay
    this.getX = getX
    this.setX = setX
    this.dragging = false

    this.internal = {
      hscale: 2, // FIXME calculate from client width
      vscale: 2, // FIXME calculate from client height
      origin: { x: 0 },
      start: { x: 0 }
    }
  }

  get origin () {
    return this.internal.origin
  }

  get startXY () {
    return this.internal.start
  }

  start (event, canvas) {
    this.dragging = true
    this.internal.origin = { x: this.getX() }
    this.internal.start = { x: this.internal.hscale * event.offsetX, y: this.internal.vscale * event.offsetY }

    canvas.onpointermove = (event) => this.onPointerMove(event, canvas)
    canvas.onpointerup = (event) => this.onPointerUp(event, canvas)

    canvas.setPointerCapture(event.pointerId)
  }

  onPointerMove (event, canvas) {
    if (this.dragging) {
      const hscale = 2 // FIXME calculate from client width
      const vscale = 2 // FIXME calculate from client height
      const xy = { x: hscale * event.offsetX, y: vscale * event.offsetY }
      const dx = xy.x - this.startXY.x

      this.setX(this.origin.x + dx)

      redraw(this.overlay, canvas)
    }
  }

  onPointerUp (event, canvas, drag) {
    canvas.onpointermove = null
    canvas.onpointerup = null
    canvas.releasePointerCapture(event.pointerId)

    if (this.dragging) {
      const hscale = 2 // FIXME calculate from client width
      const vscale = 2 // FIXME calculate from client height
      const xy = { x: hscale * event.offsetX, y: vscale * event.offsetY }
      const dx = xy.x - this.startXY.x

      this.dragging = false
      this.setX(this.origin.x + dx)

      redraw(this.overlay, canvas)
    }
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
  ctx.lineWidth = 0

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

// function onChange (component) {
// }

// function onChanged (component) {
// }

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
