import * as styles from '../styles.js'

export class Waveform extends HTMLElement {
  static get observedAttributes () {
    return ['vscale', 'colour', 'color']
  }

  constructor () {
    super()

    this.internal = {
      style: 'line',
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
    const styles = Array.from(shadow.querySelectorAll('#waveforms input[type="radio"]'))
    const swatches = Array.from(shadow.querySelectorAll('#settings wav2png-rgba'))
    const gradients = Array.from(shadow.querySelectorAll('#settings wav2png-gradient'))
    const vscale = shadow.querySelector('#settings input#vscale')

    styles.forEach((e) => {
      e.oninput = (event) => onStyle(this, e)
    })

    swatches.forEach((e) => {
      e.onchange = (event) => onChanged(this, e)
      e.onchanged = (event) => onChanged(this, e)
    })

    gradients.forEach((e) => {
      e.oninput = (event) => onChange(this, e)
      e.onchange = (event) => onChanged(this, e)
    })

    vscale.oninput = (event) => onChange(this, vscale)
    vscale.onchange = (event) => onChanged(this, vscale)
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

  /* eslint-disable-next-line accessor-pairs */
  set colour (v) {
    const shadow = this.shadowRoot

    // ... line
    { const settings = shadow.querySelector('div[for="line"]')
      const colour = settings.querySelector('wav2png-rgba')

      colour.value = v
    }

    // ... gradient
    { const settings = shadow.querySelector('div[for="gradient"]')
      const colour1 = settings.querySelector('#g2rgba1')
      const colour2 = settings.querySelector('#g2rgba2')

      colour1.value = rgba(rgb(v), 255)
      colour2.value = rgba(rgb(v), 128)
    }

    // ... gradient3
    { const settings = shadow.querySelector('div[for="gradient3"]')
      const colour1 = settings.querySelector('#g3rgba1')
      const colour2 = settings.querySelector('#g3rgba2')
      const colour3 = settings.querySelector('#g3rgba3')

      colour1.value = v
      colour2.value = rgba(rgb(v), 128)
      colour3.value = rgba(rgb(v), 192)
    }
  }

  get style () {
    const shadow = this.shadowRoot

    switch (this.internal.style) {
      case 'line': {
        const settings = shadow.querySelector('div[for="line"]')
        const rgba = settings.querySelector('wav2png-rgba').value
        const stop = settings.querySelector('#g1gradient1').value

        return styles.lineStyle(this.vscale, rgba, stop)
      }

      case 'gradient': {
        const settings = shadow.querySelector('div[for="gradient"]')
        const rgba1 = settings.querySelector('#g2rgba1').value
        const rgba2 = settings.querySelector('#g2rgba2').value
        const stop1 = settings.querySelector('#g2gradient1').value
        const stop2 = settings.querySelector('#g2gradient2').value

        return styles.gradientStyle(this.vscale, rgba1, rgba2, stop1, stop2)
      }

      case 'gradient3': {
        const settings = shadow.querySelector('div[for="gradient3"]')
        const rgba1 = settings.querySelector('#g3rgba1').value
        const rgba2 = settings.querySelector('#g3rgba2').value
        const rgba3 = settings.querySelector('#g3rgba3').value

        const stop1 = settings.querySelector('#g3gradient1').value
        const stop2 = settings.querySelector('#g3gradient2').value
        const stop3 = settings.querySelector('#g3gradient3').value

        return styles.gradient3Style(this.vscale, rgba1, rgba2, rgba3, stop1, stop2, stop3)
      }

      default: {
        const settings = shadow.querySelector('div[for="line"]')
        const rgba = settings.querySelector('wav2png-rgba').value
        const stop = settings.querySelector('#g1gradient1').value

        return styles.lineStyle(this.vscale, rgba, stop)
      }
    }
  }

  set style (v) {
    if (v === 'line') {
      this.internal.style = 'line'
    } else if (v === 'gradient') {
      this.internal.style = 'gradient'
    } else if (v === 'gradient3') {
      this.internal.style = 'gradient3'
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

function onStyle (component, e) {
  const shadow = component.shadowRoot
  const styles = Array.from(shadow.querySelectorAll('#waveforms input[type="radio"]'))
  const selected = styles.find((s) => s.checked)

  styles.forEach((s) => {
    const label = shadow.querySelector(`label[for="${s.id}"]`)
    const card = shadow.querySelector(`div[for="${s.value}"]`)

    label.classList.remove('selected')
    card.classList.add('hidden')
  })

  if (selected != null) {
    const label = shadow.querySelector(`label[for="${selected.id}"]`)
    const card = shadow.querySelector(`div[for="${selected.value}"]`)

    label.classList.add('selected')
    card.classList.remove('hidden')

    component.style = selected.value
  }

  if (component.onchanged) {
    component.onchanged(component.style)
  }
}

function recolour (component) {
  const shadow = component.shadowRoot

  // ... line
  {
    const settings = shadow.querySelector('div[for="line"]')
    const colour = settings.querySelector('wav2png-rgba').value

    const colours = [
      colour,
      colour,
      rgba(rgb(colour), 0)
    ]

    const midpoints = [
      Math.round(settings.querySelector('wav2png-gradient').value * 100)
    ]

    const stops = [
      `${colours[0]} 0%`,
      `${colours[1]} ${midpoints[0]}%`,
      `${colours[2]} 100%`
    ]

    const gradient = `linear-gradient(90deg, ${stops[0]}, ${stops[1]}, ${stops[2]})`

    settings.style.setProperty('--gradient', gradient)
  }

  // ... gradient
  {
    const settings = shadow.querySelector('div[for="gradient"]')

    const colours = [
      settings.querySelector('#g2rgba1').value,
      settings.querySelector('#g2rgba1').value,
      settings.querySelector('#g2rgba2').value,
      '#00000000'
    ]

    const midpoints = [
      Math.round(settings.querySelector('#g2gradient1').value * 100),
      Math.round(settings.querySelector('#g2gradient2').value * 100)
    ]

    const stops = [
      `${colours[0]} 0%`,
      `${colours[1]} ${midpoints[0]}%`,
      `${colours[2]} ${midpoints[1]}%`,
      `${colours[3]} 100%`
    ]

    const gradient = `linear-gradient(90deg, ${stops[0]}, ${stops[1]}, ${stops[2]}, ${stops[3]}`

    settings.style.setProperty('--gradient', gradient)
  }

  // ... gradient3
  {
    const settings = shadow.querySelector('div[for="gradient3"]')

    const colours = [
      settings.querySelector('#g3rgba1').value,
      settings.querySelector('#g3rgba1').value,
      settings.querySelector('#g3rgba2').value,
      settings.querySelector('#g3rgba3').value,
      '#00000000'
    ]

    const midpoints = [
      Math.round(settings.querySelector('#g3gradient1').value * 100),
      Math.round(settings.querySelector('#g3gradient2').value * 100),
      Math.round(settings.querySelector('#g3gradient3').value * 100)
    ]

    const stops = [
      `${colours[0]} 0%`,
      `${colours[1]} ${midpoints[0]}%`,
      `${colours[2]} ${midpoints[1]}%`,
      `${colours[3]} ${midpoints[2]}%`,
      `${colours[4]} 100%`
    ]

    const gradient = `linear-gradient(90deg, ${stops[0]}, ${stops[1]}, ${stops[2]}, ${stops[3]}, ${stops[4]})`

    settings.style.setProperty('--gradient', gradient)
  }
}

function rgb (rgba) {
  const match = `${rgba}`.match(/^#([a-fA-F0-9]{6})([a-fA-F0-9]{2})?$/)
  const u = Number.parseInt(match[1], 16).toString(16).padStart(6, '0')

  return `#${u}`
}

function rgba (rgb, alpha) {
  const match = `${rgb}`.match(/^#([a-fA-F0-9]{6})([a-fA-F0-9]{2})?$/)
  const u = Number.parseInt(match[1], 16).toString(16).padStart(6, '0')
  const v = Number.parseInt(`${alpha}`).toString(16).padStart(2, '0')

  return `#${u}${v}`
}

customElements.define('wav2png-waveform', Waveform)
