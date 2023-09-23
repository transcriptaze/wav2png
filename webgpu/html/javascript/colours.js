export const transparent = [0, 0, 0, 0]
export const black = [0, 0, 0, 1]
export const green = [0, 1, 0, 1]
export const darkgreen = [0, 0.5, 0, 1]
export const lightblue = [0.5, 0.8, 1, 1]

export function rgba (colour) {
  const match = `${colour}`.match(/^#([a-fA-F0-9]{8})$/)

  if (match && match.length > 1) {
    const hex = match[1]
    const v = Number.parseInt(hex, 16)
    const r = ((v >>> 24) & 0x00ff) / 255
    const g = ((v >>> 16) & 0x00ff) / 255
    const b = ((v >>> 8) & 0x00ff) / 255
    const a = ((v >>> 0) & 0x00ff) / 255

    return [r * a, g * a, b * a, a]
  }

  return [0, 0, 0, 0]
}
