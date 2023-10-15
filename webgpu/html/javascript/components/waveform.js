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
    const swatches = Array.from(shadow.querySelectorAll('#settings input.swatch'))
    const alphas = Array.from(shadow.querySelectorAll('#settings input.alpha'))
    const gradients = Array.from(shadow.querySelectorAll('#settings wav2png-gradient'))
    const vscale = shadow.querySelector('#settings input#vscale')

    styles.forEach((e) => {
      e.oninput = (event) => onStyle(this, e)
    })

    swatches.forEach((e) => {
      e.oninput = (event) => onChange(this, e)
      e.onchange = (event) => onChanged(this, e)
    })

    alphas.forEach((e) => {
      e.oninput = (event) => onChange(this, e)
      e.onchange = (event) => onChanged(this, e)
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
      const rgb = settings.querySelector('input#rgb')
      const alpha = settings.querySelector('input#alpha')

      rgb.value = v
      alpha.value = 255
    }

    // ... gradient
    { const settings = shadow.querySelector('div[for="gradient"]')

      { const rgb = settings.querySelector('input#rgb1')
        const alpha = settings.querySelector('input#alpha1')

        rgb.value = v
        alpha.value = 255
      }

      { const rgb = settings.querySelector('input#rgb2')
        const alpha = settings.querySelector('input#alpha2')

        rgb.value = v
        alpha.value = 128
      }
    }

    // ... gradient3
    { const settings = shadow.querySelector('div[for="gradient3"]')

      { const rgb = settings.querySelector('input#rgb1')
        const alpha = settings.querySelector('input#alpha1')

        rgb.value = v
        alpha.value = 255
      }

      { const rgb = settings.querySelector('input#rgb2')
        const alpha = settings.querySelector('input#alpha2')

        rgb.value = v
        alpha.value = 64
      }

      { const rgb = settings.querySelector('input#rgb3')
        const alpha = settings.querySelector('input#alpha3')

        rgb.value = v
        alpha.value = 128
      }
    }
  }

  get style () {
    const shadow = this.shadowRoot

    switch (this.internal.style) {
      case 'line': {
        const settings = shadow.querySelector('div[for="line"]')
        const rgb = settings.querySelector('input#rgb').value
        const alpha = settings.querySelector('input#alpha').value
        const midpoint = settings.querySelector('wav2png-gradient').value

        return styles.lineStyle(this.vscale, rgba(rgb, alpha), midpoint)
      }

      case 'gradient': {
        const settings = shadow.querySelector('div[for="gradient"]')
        const rgb1 = settings.querySelector('input#rgb1').value
        const alpha1 = settings.querySelector('input#alpha1').value
        const rgb2 = settings.querySelector('input#rgb2').value
        const alpha2 = settings.querySelector('input#alpha2').value
        const midpoint = settings.querySelector('wav2png-gradient').value

        return styles.gradientStyle(this.vscale, rgba(rgb1, alpha1), rgba(rgb2, alpha2), midpoint)
      }

      case 'gradient3': {
        const settings = shadow.querySelector('div[for="gradient3"]')
        const rgb1 = settings.querySelector('input#rgb1').value
        const alpha1 = settings.querySelector('input#alpha1').value
        const rgb2 = settings.querySelector('input#rgb2').value
        const alpha2 = settings.querySelector('input#alpha2').value
        const rgb3 = settings.querySelector('input#rgb3').value
        const alpha3 = settings.querySelector('input#alpha3').value
        const midpoint = settings.querySelector('wav2png-gradient').value

        return styles.gradient3Style(this.vscale, rgba(rgb1, alpha1), rgba(rgb2, alpha2), rgba(rgb3, alpha3), midpoint)
      }

      default: {
        const settings = shadow.querySelector('div[for="line"]')
        const rgb = settings.querySelector('input#rgb').value
        const alpha = settings.querySelector('input#alpha').value
        const midpoint = settings.querySelector('wav2png-gradient').value

        return styles.lineStyle(this.vscale, rgba(rgb, alpha), midpoint)
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
    const stops = ['#00000000', '#80ccffff', '#80ccffff']

    const rgb = settings.querySelector('input#rgb').value
    const alpha = settings.querySelector('input#alpha').value
    const midpoint = Math.round(settings.querySelector('wav2png-gradient').value * 100)

    stops[0] = rgba(rgb, 0)
    stops[1] = rgba(rgb, alpha)
    stops[2] = rgba(rgb, alpha)

    const gradient = `linear-gradient(90deg, ${stops[2]} 0%, ${stops[1]} ${midpoint}%, ${stops[0]} 100%)`

    settings.style.setProperty('--gradient', gradient)
  }

  // ... gradient
  {
    const settings = shadow.querySelector('div[for="gradient"]')
    const midpoint = Math.round(settings.querySelector('wav2png-gradient').value * 100)
    const stops = ['#00000000', '#80ccffff', '#80ccffff']

    {
      const rgb = settings.querySelector('input#rgb2').value
      const alpha = settings.querySelector('input#alpha2').value
      stops[1] = rgba(rgb, alpha)
    }

    {
      const rgb = settings.querySelector('input#rgb1').value
      const alpha = settings.querySelector('input#alpha1').value
      stops[2] = rgba(rgb, alpha)
    }

    const gradient = `linear-gradient(90deg, ${stops[2]} 0%, ${stops[1]} ${midpoint}%, ${stops[0]} 100%)`

    console.log('>>', midpoint)
    console.log(gradient)

    settings.style.setProperty('--gradient', gradient)
  }

  // ... gradient3
  {
    const settings = shadow.querySelector('div[for="gradient3"]')
    const stops = ['#80ccff80', '#80ccff40', '#80ccffff']
    const midpoint = Math.round(settings.querySelector('wav2png-gradient').value * 100)

    {
      const rgb = settings.querySelector('input#rgb3').value
      const alpha = settings.querySelector('input#alpha3').value
      stops[0] = rgba(rgb, alpha)
    }

    {
      const rgb = settings.querySelector('input#rgb2').value
      const alpha = settings.querySelector('input#alpha2').value
      stops[1] = rgba(rgb, alpha)
    }

    {
      const rgb = settings.querySelector('input#rgb1').value
      const alpha = settings.querySelector('input#alpha1').value
      stops[2] = rgba(rgb, alpha)
    }

    const gradient = `linear-gradient(90deg, ${stops[2]} 0%, ${stops[1]} ${midpoint}%, ${stops[0]} 100%)`

    settings.style.setProperty('--gradient', gradient)
  }
}

function rgba (rgb, alpha) {
  const match = `${rgb}`.match(/^#([a-fA-F0-9]{6})$/)

  const u = Number.parseInt(match[1], 16).toString(16).padStart(6, '0')
  const v = Number.parseInt(`${alpha}`).toString(16).padStart(2, '0')

  return `#${u}${v}`
}

customElements.define('wav2png-waveform', Waveform)
