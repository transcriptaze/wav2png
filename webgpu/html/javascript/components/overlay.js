export class Overlay extends HTMLElement {
  static get observedAttributes () {
    return ['width', 'height', 'padding']
  }

  constructor () {
    super()

    this.internal = {
      padding: 20,
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

    canvas.onpointerdown = (event) => onPointerDown(this, canvas, event)
    canvas.onpointerup = (event) => onPointerUp(this, canvas, event)

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
}

function redraw (component, canvas) {
  const width = canvas.width
  const height = canvas.height
  const padding = component.padding
  const ctx = canvas.getContext('2d')

  // ... draw sizing handles
  ctx.fillStyle = '#80ccffa0'
  ctx.strokeStyle = '#80ccffa0'
  ctx.lineWidth = 0

  ctx.beginPath()
  ctx.moveTo(0, 0)
  ctx.lineTo(padding, 0)
  ctx.lineTo(padding + 32, height / 2)
  ctx.lineTo(padding, height)
  ctx.lineTo(0, height)
  ctx.fill()

  ctx.beginPath()
  ctx.moveTo(width, 0)
  ctx.lineTo(width - padding, 0)
  ctx.lineTo(width - padding - 32, height / 2)
  ctx.lineTo(width - padding, height)
  ctx.lineTo(width, height)
  ctx.fill()
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

function onPointerDown (component, canvas, event) {
  const width = canvas.width
  const height = canvas.height
  const padding = component.padding

  if (event.button === 0) {
    event.preventDefault()

    const hscale = 2 // FIXME calculate from client width
    const vscale = 2 // FIXME calculate from client height
    const xy = { x: hscale * event.offsetX, y: vscale * event.offsetY }
    const start = 0
    const end = width

    const left = [
      { x: start, y: 0 },
      { x: start + padding, y: 0 },
      { x: start + padding + 32, y: height / 2 },
      { x: start + padding, y: height },
      { x: start, y: height }
    ]

    const right = [
      { x: end, y: 0 },
      { x: end - padding, y: 0 },
      { x: end - padding - 32, y: height / 2 },
      { x: end - padding, y: height },
      { x: end, y: height }
    ]

    if (hittest(xy, left) || hittest(xy, right)) {
      canvas.onpointermove = (event) => onPointerMove(component, canvas, event)
      canvas.setPointerCapture(event.pointerId)

      //   const gotcha = function (origin, p, icon) {
      //     drag.dragging = true
      //     drag.origin = origin
      //     drag.start = { x: event.offsetX, y: event.offsetY }
      //     drag.inflection = p
      //     drag.icon = icon
      //     drag.context.xscale = xscale
      //     drag.context.yscale = yscale
      //     drag.context.X = X
      //     drag.context.Xʼ = Xʼ
      //     drag.context.Y = Y
      //     drag.context.Yʼ = Yʼ
      //   }
    }
  }
}

function onPointerUp (component, canvas, event) {
  canvas.onpointermove = null
  canvas.releasePointerCapture(event.pointerId)

  // if (this.drag.dragging) {
  //   onMouseUp(this, event, this.drag)
  //   this.redraw()
  // }
}

function onPointerMove (component, canvas, event) {
}

// function onMouseUp (editor, event, drag) {
//   drag.dragging = false
//
//   const p = drag.inflection
//   const canvas = event.currentTarget
//
//   const xscale = drag.context.xscale
//   const yscale = drag.context.yscale
//   const X = drag.origin.x
//   const Y = drag.context.Y

//   const x = X + xscale * p.at
//   const dx = (canvas.width / 600) * (event.offsetX - drag.start.x)
//   const xʼ = x + dx
//   const at = (xʼ - X) / xscale
//
//   const y = Y + yscale * p.level
//   const dy = (canvas.height / 308) * (event.offsetY - drag.start.y)
//   const yʼ = y + dy
//   const level = (yʼ - Y) / yscale
//
//   const evt = new CustomEvent('changed', {
//     detail: {
//       tag: drag.inflection.tag,
//       at,
//       level
//     }
//   })
//
//   editor.dispatchEvent(evt)
// }

// function onMouseMove (editor, event, envelope, drag) {
//   const p = drag.inflection
//   const canvas = event.currentTarget
//
//   const xscale = drag.context.xscale
//   const yscale = drag.context.yscale
//   const X = drag.origin.x
//   const Y = drag.context.Y
//
//   const x = X + xscale * p.at
//   const dx = (canvas.width / 600) * (event.offsetX - drag.start.x)
//   const xʼ = x + dx
//   const at = (xʼ - X) / xscale
//
//   const y = Y + yscale * p.level
//   const dy = (canvas.height / 308) * (event.offsetY - drag.start.y)
//   const yʼ = y + dy
//   const level = (yʼ - Y) / yscale
//
//   const evt = new CustomEvent('change', {
//     detail: {
//       tag: drag.inflection.tag,
//       at,
//       level
//     }
//   })

//   editor.dispatchEvent(evt)
// }

// function onChange (component) {
// }

// function onChanged (component) {
// }

customElements.define('wav2png-overlay', Overlay)
