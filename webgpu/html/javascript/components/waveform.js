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
    const rgbas = Array.from(shadow.querySelectorAll('#settings wav2png-rgba'))
    const swatches = Array.from(shadow.querySelectorAll('#settings input.rgb'))
    const alphas = Array.from(shadow.querySelectorAll('#settings input.alpha'))
    const gradients = Array.from(shadow.querySelectorAll('#settings wav2png-gradient'))
    const vscale = shadow.querySelector('#settings input#vscale')

    styles.forEach((e) => {
      e.oninput = (event) => onStyle(this, e)
    })

    rgbas.forEach((e) => {
      e.onchange = (event) => onChanged(this, e)
      e.onchanged = (event) => onChanged(this, e)
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
      const rgba = settings.querySelector('wav2png-rgba')

      rgba.colour = v
    }

    // ... gradient
    { const settings = shadow.querySelector('div[for="gradient"]')
      const rgba1 = settings.querySelector('#g2rgba1')
      const rgba2 = settings.querySelector('#g2rgba2')

      rgba1.colour = rgba(rgb(v), 255)
      rgba2.colour = rgba(rgb(v), 128)
    }

    // ... gradient3
    { const settings = shadow.querySelector('div[for="gradient3"]')

      { const rgb = settings.querySelector('#g3rgb1')
        const alpha = settings.querySelector('#g3alpha1')

        rgb.value = v
        alpha.value = 255
      }

      { const rgb = settings.querySelector('#g3rgb2')
        const alpha = settings.querySelector('#g3alpha2')

        rgb.value = v
        alpha.value = 64
      }

      { const rgb = settings.querySelector('#g3rgb3')
        const alpha = settings.querySelector('#g3alpha3')

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
        const rgb1 = settings.querySelector('#g3rgb1').value
        const alpha1 = settings.querySelector('#g3alpha1').value
        const rgb2 = settings.querySelector('#g3rgb2').value
        const alpha2 = settings.querySelector('#g3alpha2').value
        const rgb3 = settings.querySelector('#g3rgb3').value
        const alpha3 = settings.querySelector('#g3alpha3').value
        const stop1 = settings.querySelector('#g3gradient1').value
        const stop2 = settings.querySelector('#g3gradient2').value
        const stop3 = settings.querySelector('#g3gradient3').value

        return styles.gradient3Style(this.vscale, rgba(rgb1, alpha1), rgba(rgb2, alpha2), rgba(rgb3, alpha3), stop1, stop2, stop3)
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
    const midpoint = settings.querySelector('wav2png-gradient').value

    const colours = [colour, colour, rgba(rgb(colour), 0)]
    const midpoints = [Math.round(midpoint * 100)]
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
    const colour1 = settings.querySelector('#g2rgba1').value
    const colour2 = settings.querySelector('#g2rgba2').value
    const midpoint1 = settings.querySelector('#g2gradient1').value
    const midpoint2 = settings.querySelector('#g2gradient2').value

    const colours = [colour1, colour1, colour2, '#00000000']
    const midpoints = [
      Math.round(midpoint1 * 100),
      Math.round(midpoint2 * 100)
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
    const midpoint1 = Math.round(settings.querySelector('#g3gradient1').value * 100)
    const midpoint2 = Math.round(settings.querySelector('#g3gradient2').value * 100)
    const midpoint3 = Math.round(settings.querySelector('#g3gradient3').value * 100)
    const stops = ['#80ccffff', '#80ccffff', '#80ccff80', '#80ccff40', '#00000000']

    {
      const rgb = settings.querySelector('#g3rgb1').value
      const alpha = settings.querySelector('#g3alpha1').value
      stops[0] = rgba(rgb, alpha)
      stops[1] = rgba(rgb, alpha)
    }

    {
      const rgb = settings.querySelector('#g3rgb2').value
      const alpha = settings.querySelector('#g3alpha2').value
      stops[2] = rgba(rgb, alpha)
    }

    {
      const rgb = settings.querySelector('#g3rgb3').value
      const alpha = settings.querySelector('#g3alpha3').value
      stops[3] = rgba(rgb, alpha)
      stops[4] = '#00000000'
    }

    const gradient = `linear-gradient(90deg, ${stops[0]} 0%, ${stops[1]} ${midpoint1}%, ${stops[2]} ${midpoint2}%, ${stops[3]} ${midpoint3}%, ${stops[4]} 100%)`

    settings.style.setProperty('--gradient', gradient)
  }
}

function rgb (rgba) {
  const match = `${rgba}`.match(/^#([a-fA-F0-9]{6})([a-fA-F0-9]{2})?$/)
  const u = Number.parseInt(match[1], 16).toString(16).padStart(6, '0')

  return `#${u}`
}

// function alpha (rgba) {
//   const match = `${rgba}`.match(/^#([a-fA-F0-9]{6})([a-fA-F0-9]{2})$/)
//   const v = Number.parseInt(match[2], 16).toString(16).padStart(2, '0')
//
//   return `${v}`
// }

function rgba (rgb, alpha) {
  const match = `${rgb}`.match(/^#([a-fA-F0-9]{6})([a-fA-F0-9]{2})?$/)
  const u = Number.parseInt(match[1], 16).toString(16).padStart(6, '0')
  const v = Number.parseInt(`${alpha}`).toString(16).padStart(2, '0')

  return `#${u}${v}`
}

customElements.define('wav2png-waveform', Waveform)
