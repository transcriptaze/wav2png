import { rgba } from '../colours.js'

export class Waveform extends HTMLElement {
  static get observedAttributes () {
    return ['vscale', 'colour', 'color']
  }

  constructor () {
    super()

    this.internal = {
      onChange: null,
      onChanged: null
    }

    const template = document.querySelector('#template-waveform')
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
    const vscale = shadow.querySelector('input#vscale')
    const rgb = shadow.querySelector('input#rgb')
    const alpha = shadow.querySelector('input#alpha')

    vscale.oninput = (event) => onChange(this)
    rgb.oninput = (event) => onChange(this)
    alpha.oninput = (event) => onChange(this)

    vscale.onchange = (event) => onChanged(this)
    rgb.onchange = (event) => onChanged(this)
    alpha.onchange = (event) => onChanged(this)
  }

  disconnectedCallback () {
  }

  adoptedCallback () {
  }

  attributeChangedCallback (name, from, to) {
    if (name === 'vscale') {
      this.vscale = to
    }

    if (name === 'colour' || name === 'color') {
      this.colour = to
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

  get vscale () {
    const shadow = this.shadowRoot
    const input = shadow.querySelector('input#vscale')
    const vscale = Number.parseFloat(input.value)

    return !Number.isNaN(vscale) ? vscale : 1.0
  }

  set vscale (v) {
    const shadow = this.shadowRoot
    const input = shadow.querySelector('input#vscale')
    const min = Number.parseFloat(input.min)
    const max = Number.parseFloat(input.max)
    const vscale = Number.parseFloat(`${v}`)

    input.value = Math.min(Math.max(vscale, min), max)
  }

  get colour () {
    const shadow = this.shadowRoot
    const input = shadow.querySelector('input#rgb')
    const range = shadow.querySelector('input#alpha')
    const rgb = input.value
    const alpha = Math.trunc(range.value * 255).toString(16).padStart(2, '0')

    return `${rgb}${alpha}`
  }

  set colour (v) {
    const shadow = this.shadowRoot
    const input = shadow.querySelector('input#rgb')
    const range = shadow.querySelector('input#alpha')
    const match = `${v}`.match(/^#([a-fA-F0-9]{6})([a-fA-F0-9]{2})?$/)

    if (match != null) {
      const rgb = Number.parseInt(match[1], 16).toString(16).padStart(6, '0')
      const alpha = Number.parseInt(match[2], 16)

      input.value = `#${rgb}`

      if (!Number.isNaN(alpha)) {
        range.value = Math.min(Math.max(alpha, 0), 255) / 255
      } else {
        range.value = 1.0
      }
    }
  }

  get waveform () {
    return {
      vscale: this.vscale,
      colour: rgba(this.colour)
    }
  }
}

function onChange (component) {
  const waveform = component.waveform
  if (component.onchange) {
    component.onchange(waveform)
  }
}

function onChanged (component) {
  const waveform = component.waveform
  if (component.onchanged) {
    component.onchanged(waveform)
  }
}

customElements.define('wav2png-waveform', Waveform)
