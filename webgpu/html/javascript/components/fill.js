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

    const rgba = shadow.querySelector('wav2png-rgba')
    if (rgba) {
      rgba.onchange = (event) => onChange(this)
      rgba.onchanged = (event) => onChanged(this)
    }
  }

  disconnectedCallback () {
  }

  adoptedCallback () {
  }

  attributeChangedCallback (name, from, to) {
    const shadow = this.shadowRoot

    const rgba = shadow.querySelector('wav2png-rgba')
    if (rgba) {
      if (name === 'colour' || name === 'color') {
        this.colour = to
      }
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
    const rgba = shadow.querySelector('wav2png-rgba')

    return rgba.value
  }

  set colour (v) {
    const shadow = this.shadowRoot
    const rgba = shadow.querySelector('wav2png-rgba')

    rgba.value = v
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
