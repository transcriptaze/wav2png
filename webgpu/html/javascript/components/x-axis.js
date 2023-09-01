const TIMESCALE = 10000
const DT = 10

export class XAxis extends HTMLElement {
  static get observedAttributes () {
    return []
  }

  constructor () {
    super()

    this.internal = {
      duration: 0,
      start: 0,
      end: 0,

      onChanged: null
    }

    const template = document.querySelector('#template-x-axis')
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
    const left = shadow.querySelector('#left')
    const right = shadow.querySelector('#right')
    const plus = shadow.querySelector('#plus')
    const minus = shadow.querySelector('#minus')

    left.onclick = (event) => onLeft(this, event)
    right.onclick = (event) => onRight(this, event)
    plus.onclick = (event) => onPlus(this, event)
    minus.onclick = (event) => onMinus(this, event)

    left.oncontextmenu = (event) => onLeft(this, event)
    right.oncontextmenu = (event) => onRight(this, event)
    plus.oncontextmenu = (event) => onPlus(this, event)
    minus.oncontextmenu = (event) => onMinus(this, event)
  }

  disconnectedCallback () {
  }

  adoptedCallback () {
  }

  attributeChangedCallback (name, from, to) {
  }

  get onchanged () {
    return this.internal.onChanged
  }

  set onchanged (v) {
    this.internal.onChanged = v
  }

  /* eslint-disable-next-line accessor-pairs */
  set audio ({ start, end, duration }) {
    const shadow = this.shadowRoot
    const left = shadow.querySelector('#left')
    const right = shadow.querySelector('#right')
    const plus = shadow.querySelector('#plus')
    const minus = shadow.querySelector('#minus')
    const t = Number.parseFloat(`${duration}`)

    if (!Number.isNaN(t)) {
      this.internal.duration = Math.round(t * TIMESCALE)
      this.start = start
      this.end = end

      left.disabled = this.internal.duration === 0
      right.disabled = this.internal.duration === 0
      plus.disabled = this.internal.duration === 0
      minus.disabled = this.internal.duration === 0
    }
  }

  get duration () {
    return this.internal.duration / TIMESCALE
  }

  get start () {
    return this.internal.start / TIMESCALE
  }

  set start (v) {
    const t = Number.parseFloat(`${v}`)

    if (Number.isNaN(t) || this.duration === 0) {
      this.internal.start = 0
    } else {
      this.internal.start = constrain(Math.round(t * TIMESCALE), 0, this.internal.end)
    }

    this.reselect()
  }

  get end () {
    return this.internal.end / TIMESCALE
  }

  set end (v) {
    const t = Number.parseFloat(`${v}`)

    if (Number.isNaN(t) || this.internal.duration === 0) {
      this.internal.end = 0
    } else {
      this.internal.end = constrain(Math.round(t * TIMESCALE), 0, this.internal.duration)
    }

    this.reselect()
  }

  /* eslint-disable-next-line accessor-pairs */
  reselect () {
    const shadow = this.shadowRoot
    const start = shadow.querySelector('#start')
    const end = shadow.querySelector('#end')
    const duration = shadow.querySelector('#duration')
    const from = this.internal.start
    const to = this.internal.end
    const delta = to - from

    start.querySelector('label').innerHTML = format(this.start)
    end.querySelector('label').innerHTML = format(this.end)
    duration.querySelector('label').innerHTML = format(delta / TIMESCALE)

    if (this.duration === 0) {
      start.classList.remove('visible')
      end.classList.remove('visible')
      duration.classList.remove('visible')
    } else {
      start.classList.add('visible')
      end.classList.add('visible')
      duration.classList.add('visible')
    }
  }
}

function onLeft (xaxis, event) {
  event.preventDefault()

  const dt = event.altKey ? DT / 10 : (event.ctrlKey ? DT * 10 : DT)
  const p = xaxis.internal.start
  const q = constrain(p + dt, 0, xaxis.internal.end)
  const delta = q - p

  xaxis.internal.start += delta
  xaxis.internal.end += delta

  xaxis.reselect()

  if (xaxis.onchanged != null) {
    xaxis.onchanged(xaxis.start, xaxis.end)
  }
}

function onRight (xaxis, event) {
  event.preventDefault()

  const dt = event.altKey ? DT / 10 : (event.ctrlKey ? DT * 10 : DT)
  const p = xaxis.internal.end
  const q = constrain(p - dt, xaxis.internal.start, xaxis.internal.duration)
  const delta = q - p

  xaxis.internal.start += delta
  xaxis.internal.end += delta

  xaxis.reselect()

  if (xaxis.onchanged != null) {
    xaxis.onchanged(xaxis.start, xaxis.end)
  }
}

function onPlus (xaxis, event) {
  event.preventDefault()

  const dt = event.altKey ? DT / 10 : (event.ctrlKey ? DT * 10 : DT)
  const delta = constrain(xaxis.internal.end - xaxis.internal.start + dt, 0, xaxis.internal.duration)

  const q = constrain(xaxis.internal.start + delta, xaxis.internal.start, xaxis.internal.duration)
  const p = constrain(q - delta, 0, xaxis.internal.end)

  xaxis.internal.start = p
  xaxis.internal.end = q

  xaxis.reselect()

  if (xaxis.onchanged != null) {
    xaxis.onchanged(xaxis.start, xaxis.end)
  }
}

function onMinus (xaxis, event) {
  event.preventDefault()

  const dt = event.altKey ? DT / 10 : (event.ctrlKey ? DT * 10 : DT)
  const delta = constrain(xaxis.internal.end - xaxis.internal.start - dt, 0, xaxis.internal.duration)

  const q = constrain(xaxis.internal.start + delta, xaxis.internal.start, xaxis.internal.duration)
  const p = constrain(q - delta, 0, xaxis.internal.end)

  xaxis.internal.start = p
  xaxis.internal.end = q

  xaxis.reselect()

  if (xaxis.onchanged != null) {
    xaxis.onchanged(xaxis.start, xaxis.end)
  }
}

function format (v) {
  const mm = Math.trunc(v / 60)
  const ss = Math.trunc(v) % 60
  const ms = Math.trunc(v * 1000) % 1000
  const µs = Math.trunc(v * 10000) % 10000 % 10

  const mmssms = mm > 0 ? `${mm}:${ss}.${ms.toString().padStart(3, '0')}` : `${ss}.${ms.toString().padStart(3, '0')}`

  return µs > 0 ? `${mmssms}.${µs}` : `${mmssms}`
}

function constrain (v, min, max) {
  return Math.min(Math.max(v, min), max)
}

customElements.define('wav2png-x-axis', XAxis)
