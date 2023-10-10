export function slice (audio, width, padding) {
  const N = audio.audio.length
  const duration = clamp(audio.duration, 0, N / audio.fs)
  const start = duration === 0 ? 0 : clamp(Math.floor(N * audio.start / duration), 0, N)
  const end = duration === 0 ? 0 : clamp(Math.floor(N * audio.end / duration), 0, N)

  const roundµs = function (t) {
    return Math.round(1000_000 * duration * t / N) / 1000_000
  }

  const pixels = Math.min(width - 2 * padding, end - start)
  const stride = Math.fround((end - start) / pixels)

  const slice = {
    start: 0,
    end,
    audio: audio.audio.subarray(start, end),
    offset: 0,
    pixels,
    stride
  }

  // ... calculate start/end indices and offset
  {
    let index = 0
    let start = Math.round(index * stride)
    let end = Math.round((index + 1) * stride)

    while (roundµs(start) < audio.start) {
      slice.start = start
      slice.offset = index

      index += 1
      start = Math.round(index * stride)
    }

    end = Math.round((index + 1) * stride)
    slice.end = end

    while (roundµs(end) <= audio.end) {
      slice.end = end

      index += 1
      end = Math.round((index + 1) * stride)
    }

    slice.audio = audio.audio.subarray(slice.start, slice.end)
  }

  return slice
}

function clamp (v, min, max) {
  return (v < min) ? min : ((v > max) ? max : v)
}
