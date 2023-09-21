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
    const rgb = shadow.querySelector('input#rgb')
    const alpha = shadow.querySelector('input#alpha')
    const vscale = shadow.querySelector('input#vscale')
    const pickers = Array.from(shadow.querySelectorAll('#waveforms input[type="radio"]'))

    rgb.oninput = (event) => onChange(this)
    alpha.oninput = (event) => onChange(this)
    vscale.oninput = (event) => onChange(this)

    rgb.onchange = (event) => onChanged(this)
    alpha.onchange = (event) => onChanged(this)
    vscale.onchange = (event) => onChanged(this)

    pickers.forEach((e) => {
      e.oninput = (event) => onWaveForm(this, event)
    })
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
    const exemplar = shadow.querySelector('svg.exemplar #bar')
    const match = `${v}`.match(/^#([a-fA-F0-9]{6})([a-fA-F0-9]{2})?$/)

    if (match != null) {
      const rgb = Number.parseInt(match[1], 16).toString(16).padStart(6, '0')
      const alpha = Number.parseInt(match[2], 16)

      input.value = `#${rgb}`
      exemplar.style.fill = `"#${rgb}"`

      if (!Number.isNaN(alpha)) {
        range.value = Math.min(Math.max(alpha, 0), 255) / 255
      } else {
        range.value = 1.0
      }
    }
  }

  get waveform () {
    const colour = rgba(this.colour)

    return {
      vscale: this.vscale,
      colours: [
        colour,
        [colour[0], colour[1], colour[2], 0.25]
      ]
    }
  }
}

function onChange (component) {
  const shadow = component.shadowRoot
  const waveform = component.waveform
  const exemplar = shadow.querySelector('svg.exemplar #bar')

  exemplar.style.fill = component.colour

  if (component.onchange) {
    component.onchange(waveform)
  }
}

function onChanged (component) {
  const shadow = component.shadowRoot
  const waveform = component.waveform
  const exemplar = shadow.querySelector('svg.exemplar #bar')

  exemplar.style.fill = component.colour

  if (component.onchanged) {
    component.onchanged(waveform)
  }
}

function onWaveForm (component, event) {
  const shadow = component.shadowRoot
  const pickers = Array.from(shadow.querySelectorAll('#waveforms input[type="radio"]'))

  pickers.forEach((e) => {
    const label = shadow.querySelector(`label[for="${e.id}"`)

    if (e.id === event.target.id) {
      label.classList.add('selected')
    } else {
      label.classList.remove('selected')
    }
  })
}

customElements.define('wav2png-waveform', Waveform)
