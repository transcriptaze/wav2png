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
