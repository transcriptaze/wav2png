export class Gradient extends HTMLElement {
  static get observedAttributes () {
    return []
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
    const slider = shadow.querySelector('#gradient-slider')

    slider.oninput = (event) => onChange(this, slider)
    slider.onchange = (event) => onChanged(this, slider)
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

  // /* eslint-disable-next-line accessor-pairs */
  // set colour (v) {
  // }

  get value () {
    return 1.0
  }

  set value (v) {
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

// function rgba (rgb, alpha) {
//   const match = `${rgb}`.match(/^#([a-fA-F0-9]{6})$/)
//
//   const u = Number.parseInt(match[1], 16).toString(16).padStart(6, '0')
//   const v = Number.parseInt(`${alpha}`).toString(16).padStart(2, '0')
//
//   return `#${u}${v}`
// }

customElements.define('wav2png-gradient', Gradient)
