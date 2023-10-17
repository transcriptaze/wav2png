export class RGBA extends HTMLElement {
  static get observedAttributes () {
    return ['colour', 'color']
  }

  constructor () {
    super()

    this.internal = {
      onChange: null,
      onChanged: null
    }

    const template = document.querySelector('#template-rgba')
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
    const rgb = shadow.querySelector('#rgb')
    const alpha = shadow.querySelector('#alpha')

    rgb.oninput = (event) => onChange(this, rgb)
    rgb.onchange = (event) => onChanged(this, rgb)

    alpha.oninput = (event) => onChange(this, alpha)
    alpha.onchange = (event) => onChanged(this, alpha)
  }

  disconnectedCallback () {
  }

  adoptedCallback () {
  }

  attributeChangedCallback (name, from, to) {
    if (name === 'colour' || name === 'color') {
      this.value = to
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

  get value () {
    const shadow = this.shadowRoot
    const rgb = shadow.querySelector('#rgb').value
    const alpha = shadow.querySelector('#alpha').value

    const match = `${rgb}`.match(/^#([a-fA-F0-9]{6})$/)
    const u = Number.parseInt(match[1], 16).toString(16).padStart(6, '0')
    const v = Number.parseInt(`${alpha}`).toString(16).padStart(2, '0')

    return `#${u}${v}`
  }

  set value (v) {
    const shadow = this.shadowRoot
    const rgb = shadow.querySelector('#rgb')
    const alpha = shadow.querySelector('#alpha')

    const match = `${v}`.match(/^#([a-fA-F0-9]{6})([a-fA-F0-9]{2})?$/)

    if (match != null && match.length > 1) {
      const u = Number.parseInt(match[1], 16).toString(16).padStart(6, '0')

      rgb.value = `#${u}`
    }

    if (match != null && match.length > 2) {
      const u = Number.parseInt(match[2], 16)

      alpha.value = u
    }
  }
}

function onChange (component, e) {
  if (component.onchange) {
    component.onchange(component.style)
  }
}

function onChanged (component, e) {
  if (component.onchanged) {
    component.onchanged(component.style)
  }
}

customElements.define('wav2png-rgba', RGBA)
