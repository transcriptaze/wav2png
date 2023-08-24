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

      onChange: null,
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

    left.onclick = (event) => onLeft(this, event)
    right.onclick = (event) => onRight(this, event)
  }

  disconnectedCallback () {
  }

  adoptedCallback () {
  }

  attributeChangedCallback (name, from, to) {
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

  /* eslint-disable-next-line accessor-pairs */
  set audio ({ start, end, duration }) {
    const shadow = this.shadowRoot
    const left = shadow.querySelector('#left')
    const right = shadow.querySelector('#right')
    const t = Number.parseFloat(`${duration}`)

    if (!Number.isNaN(t)) {
      this.internal.duration = Math.round(t * 1000)
      this.start = start
      this.end = end

      left.disabled = this.internal.duration === 0
      right.disabled = this.internal.duration === 0
    }
  }

  get duration () {
    return this.internal.duration / 1000
  }

  get start () {
    return this.internal.start / 1000
  }

  set start (v) {
    const shadow = this.shadowRoot
    const start = shadow.querySelector('#start')
    const t = Number.parseFloat(`${v}`)

    if (Number.isNaN(t) || this.duration === 0) {
      this.internal.start = 0
    } else {
      this.internal.start = constrain(Math.round(t * 1000), 0, this.internal.end)
    }

    if (Number.isNaN(t) || this.duration === 0) {
      start.innerHTML = ''
    } else {
      start.innerHTML = format(this.start)
    }

    this.reselect()
  }

  get end () {
    return this.internal.end / 1000
  }

  set end (v) {
    const shadow = this.shadowRoot
    const end = shadow.querySelector('#end')
    const t = Number.parseFloat(`${v}`)

    if (Number.isNaN(t) || this.internal.duration === 0) {
      this.internal.end = 0
    } else {
      this.internal.end = constrain(Math.round(t * 1000), 0, this.internal.duration)
    }

    if (Number.isNaN(t) || this.internal.duration === 0) {
      end.innerHTML = ''
    } else {
      end.innerHTML = format(this.end)
    }

    this.reselect()
  }

  /* eslint-disable-next-line accessor-pairs */
  reselect () {
    const shadow = this.shadowRoot
    const duration = shadow.querySelector('#duration')
    const start = this.internal.start
    const end = this.internal.end
    const dt = end - start

    if (this.internal.duration === 0) {
      duration.innerHTML = ''
    } else {
      duration.innerHTML = format(dt / 1000)
    }
  }
}

function onLeft (xaxis, event) {
}

function onRight (xaxis, event) {
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

function constrain (v, min, max) {
  return Math.min(Math.max(v, min), max)
}

customElements.define('wav2png-x-axis', XAxis)
