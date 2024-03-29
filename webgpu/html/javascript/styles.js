export const LINE = {
  type: 'line',
  line: {
    vscale: 1.0,
    colour: '#80ccffff'
  }
}

export function lineStyle (vscale, colour, midpoint) {
  return {
    type: 'line',
    line: {
      vscale,
      colour,
      midpoint
    }
  }
}

export function gradientStyle (vscale, colour1, colour2, midpoint) {
  return {
    type: 'gradient',
    gradient: {
      vscale,
      colours: [colour1, colour2],
      midpoint
    }
  }
}

export function gradient3Style (vscale, colour1, colour2, colour3, midpoint) {
  return {
    type: 'gradient3',
    gradient3: {
      vscale,
      colours: [colour1, colour2, colour3],
      stops: [0.0, midpoint, 1.0]
    }
  }
}
