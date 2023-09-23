export const LINE = {
  type: 'line',
  line: {
    vscale: 1.0,
    colour: '#80ccffff'
  }
}

export function lineStyle (vscale, colour) {
  return {
    type: 'line',
    line: {
      vscale,
      colour
    }
  }
}

export function gradientStyle (vscale, colour1, colour2) {
  return {
    type: 'gradient',
    gradient: {
      vscale,
      colours: [colour1, colour2]
    }
  }
}
