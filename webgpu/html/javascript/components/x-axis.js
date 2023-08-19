export class XAxis extends HTMLElement {
  static get observedAttributes () {
    return []
  }

  constructor () {
    super()

    this.internal = {
      duration: 0,
      start: 0,
      end: 1
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
  }

  disconnectedCallback () {
  }

  adoptedCallback () {
  }

  attributeChangedCallback (name, from, to) {
  }

  /* eslint-disable-next-line accessor-pairs */
  set audio ({ fs, audio }) {
    const duration = audio.length / fs

    this.internal.duration = Math.max(duration, 0)
    this.selected = { start: 0, end: 1 }
  }

  get duration () {
    return this.internal.duration
  }

  /* eslint-disable-next-line accessor-pairs */
  set selected ({ start, end }) {
    const shadow = this.shadowRoot
    const from = Number.parseFloat(`${start}`)
    const to = Number.parseFloat(`${end}`)

    if (Number.isNaN(from) || Number.isNaN(to) || this.internal.duration === 0) {
      shadow.querySelector('#start').innerHTML = '&nbsp;'
      shadow.querySelector('#end').innerHTML = '&nbsp;'
      shadow.querySelector('#duration').innerHTML = '&nbsp;'
    } else {
      const p = Math.min(Math.max(from, 0), Math.min(to, 1))
      const q = Math.max(Math.min(to, 1), Math.max(from, 0))

      const s = round2ms(p * this.duration)
      const t = round2ms(q * this.duration)
      const dt = t - s

      shadow.querySelector('#start').innerHTML = format(s / 1000)
      shadow.querySelector('#end').innerHTML = format(t / 1000)
      shadow.querySelector('#duration').innerHTML = format(dt / 1000)
    }
  }
}

function round2ms (v) {
  const ms = Math.trunc(v * 1000) % 1000
  const ss = Math.trunc(v) % 60
  const mm = Math.trunc(v / 60)

  return mm * 60 * 1000 + ss * 1000 + ms
}

function format (v) {
  const ms = Math.trunc(v * 1000) % 1000
  const ss = Math.trunc(v) % 60
  const mm = Math.trunc(v / 60)

  if (mm > 0) {
    return `${mm}:${ss}.${ms.toString().padStart(3, '0')}`
  } else {
    return `${ss}.${ms.toString().padStart(3, '0')}`
  }
}

customElements.define('wav2png-x-axis', XAxis)
