export class Gradient extends HTMLElement {
  static get observedAttributes () {
    return ['value']
  }

  constructor () {
    super()

    this.internal = {
      onChange: null,
      onChanged: null
    }

    const template = document.querySelector('#template-gradient')
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
    const slider = shadow.querySelector('#slider1')

    slider.oninput = (event) => onChange(this, slider)
    slider.onchange = (event) => onChanged(this, slider)
  }

  disconnectedCallback () {
  }

  adoptedCallback () {
  }

  attributeChangedCallback (name, from, to) {
    if (name === 'value') {
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
    const slider1 = shadow.querySelector('#slider1').value
    const slider2 = shadow.querySelector('#slider2').value
    const slider3 = shadow.querySelector('#slider3').value

    return [slider1, slider2, slider3]
  }

  set value (v) {
    const shadow = this.shadowRoot
    const slider = shadow.querySelector('#slider1')
    const f = Number.parseFloat(`${v}`)

    if (!Number.isNaN(f)) {
      slider.value = clamp(f, 0.0, 1.0)
    }
  }
}

function onChange (component, e) {
  recolour(component)

  if (component.onchange) {
    component.onchange(component.style)
  }
}

function onChanged (component, e) {
  recolour(component)

  if (component.onchanged) {
    component.onchanged(component.style)
  }
}

function recolour (component) {
}

function clamp (v, min, max) {
  return (v < min) ? min : ((v > max) ? max : v)
}

customElements.define('wav2png-gradient', Gradient)
