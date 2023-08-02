export class Fill extends HTMLElement {
  static get observedAttributes () {
    return ['colour', 'color']
  }

  constructor () {
    super()

    this.internal = {
      onChange: null,
      onChanged: null
    }

    const template = document.querySelector('#template-fill')
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

    rgb.oninput = (event) => onChange(this)
    alpha.oninput = (event) => onChange(this)

    rgb.onchange = (event) => onChanged(this)
    alpha.onchange = (event) => onChanged(this)
  }

  disconnectedCallback () {
  }

  adoptedCallback () {
  }

  attributeChangedCallback (name, from, to) {
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
}

function onChange (component) {
  const colour = component.colour
  if (component.onchange) {
    component.onchange(colour)
  }
}

function onChanged (component) {
  const colour = component.colour
  if (component.onchanged) {
    component.onchanged(colour)
  }
}

customElements.define('wav2png-fill', Fill)
