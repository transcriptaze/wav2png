import { describe, it } from 'mocha'
import { expect } from 'chai'

const duration = 5
const fs = 44100
const samples = new Float64Array(duration * fs)
const width = 1920
const padding = 20

for (let i = 0; i < samples.length; i++) {
  samples[i] = i * 0.1
}

describe('audio pixel bucket logic', function () {
  it.only('reference: complete buffer', function () {
    const start = 0
    const end = 220500
    const N = end - start

    const pixels = width - 2 * padding
    const stride = N / pixels

    const buckets = []
    for (let i = 0; i < pixels; i++) {
      const from = Math.round(i * stride)
      const to = Math.round((i + 1) * stride)
      const value = samples.slice(from, to).reduce((v, a) => a + v)

      buckets.push({
        start: from,
        end: to,
        value
      })
    }

    // ... check setup
    expect(pixels).to.equal(1880)
    expect(buckets.length).to.equal(pixels)
    expect(stride).to.be.closeTo(117.2, 0.1)

    // ... check first/last indices
    expect(buckets.at(0).start).to.equal(start)
    expect(buckets.at(-1).end).to.equal(end)

    // ... check index continuity
    for (let i = 1; i < buckets.length; i++) {
      const p = buckets[i - 1]
      const q = buckets[i]

      expect(q[0]).to.equal(p[1])
    }

    // ... check index stride
    const floor = Math.floor(stride)
    const ceil = Math.ceil(stride)

    for (const v of buckets) {
      expect(v.end - v.start).to.be.oneOf([floor, ceil])
    }

    // ... check pixel values
    for (let i = 0; i < pixels; i++) {
      const pixel = buckets[i]
      const N = pixel.end - pixel.start
      const v = N * (samples[pixel.start] + samples[pixel.end - 1]) / 2

      expect(pixel.value, `${i}`).to.be.closeTo(v, 0.00001)
    }

    // ... sanity check
    expect(buckets[0].value).to.be.closeTo(678.6, 0.00001)
    expect(buckets[1879].value).to.be.closeTo(2579159.7, 0.00001)
  })

  it.only('reference: 1s slice', function () {
    const start = 0
    const end = 44100
    const N = end - start

    const pixels = width - 2 * padding
    const stride = N / pixels

    const buckets = []
    for (let i = 0; i < pixels; i++) {
      const from = Math.round(i * stride)
      const to = Math.round((i + 1) * stride)
      const value = samples.slice(from, to).reduce((v, a) => a + v)

      buckets.push({
        start: from,
        end: to,
        value
      })
    }

    // ... check setup
    expect(pixels).to.equal(1880)
    expect(buckets.length).to.equal(pixels)
    expect(stride).to.be.closeTo(23.4, 0.1)

    // ... check index continuity
    for (let i = 1; i < buckets.length; i++) {
      const p = buckets[i - 1]
      const q = buckets[i]

      expect(q[0]).to.equal(p[1])
    }

    // ... check index stride
    const floor = Math.floor(stride)
    const ceil = Math.ceil(stride)

    for (const v of buckets) {
      expect(v.end - v.start).to.be.oneOf([floor, ceil])
    }

    // ... check pixel values
    for (let i = 0; i < pixels; i++) {
      const pixel = buckets[i]
      const N = pixel.end - pixel.start
      const v = N * (samples[pixel.start] + samples[pixel.end - 1]) / 2

      expect(pixel.value, `${i}`).to.be.closeTo(v, 0.00001)
    }

    // ... sanity check
    expect(buckets.at(0).start).to.equal(start)
    expect(buckets.at(-1).end).to.equal(end)

    expect(buckets[0].value).to.be.closeTo(25.3, 0.00001)
    expect(buckets[1879].value).to.be.closeTo(101402.4, 0.00001)
  })

  it.only('slice @ 1ms', function () {
    const N = 44100
    const pixels = width - 2 * padding
    const stride = N / pixels

    // ... reference slice @ 0ms
    const slice0ms = []

    {
      const start = 0
      for (let i = 0; i < pixels; i++) {
        const from = Math.round(start + i * stride)
        const to = Math.round((i + 1) * stride)
        const value = samples.slice(from, to).reduce((v, a) => a + v)

        slice0ms.push({
          start: from,
          end: to,
          value
        })
      }
    }

    // ... slice @ 1ms
    const slice1ms = []

    {
      const start = 0
      for (let i = 0; i < pixels; i++) {
        const from = Math.round(start + i * stride)
        const to = Math.round((i + 1) * stride)
        const value = samples.slice(from, to).reduce((v, a) => a + v)

        slice1ms.push({
          start: from,
          end: to,
          value
        })
      }
    }

    // ... check setup
    expect(pixels).to.equal(1880)
    expect(slice0ms.length).to.equal(pixels)
    expect(slice1ms.length).to.equal(pixels)
    expect(stride).to.be.closeTo(23.4, 0.1)

    // ... check index continuity
    for (let i = 1; i < slice0ms.length; i++) {
      const p = slice0ms[i - 1]
      const q = slice0ms[i]

      expect(q[0]).to.equal(p[1])
    }

    for (let i = 1; i < slice1ms.length; i++) {
      const p = slice1ms[i - 1]
      const q = slice1ms[i]

      expect(q[0]).to.equal(p[1])
    }

    // ... check index stride
    const floor = Math.floor(stride)
    const ceil = Math.ceil(stride)

    for (const v of slice0ms) {
      expect(v.end - v.start).to.be.oneOf([floor, ceil])
    }

    for (const v of slice1ms) {
      expect(v.end - v.start).to.be.oneOf([floor, ceil])
    }

    // ... check pixel values
    for (let i = 0; i < pixels; i++) {
      const pixel = slice0ms[i]
      const N = pixel.end - pixel.start
      const v = N * (samples[pixel.start] + samples[pixel.end - 1]) / 2

      expect(pixel.value, `${i}`).to.be.closeTo(v, 0.00001)
    }

    for (let i = 0; i < pixels; i++) {
      const pixel = slice1ms[i]
      const N = pixel.end - pixel.start
      const v = N * (samples[pixel.start] + samples[pixel.end - 1]) / 2

      expect(pixel.value, `${i}`).to.be.closeTo(v, 0.00001)
    }

    for (let i = 0; i < pixels; i++) {
      const pixel0 = slice0ms[i]
      const pixel1 = slice1ms[i]

      expect(pixel1).to.deep.equal(pixel0)
    }

    // ... sanity check
    expect(slice0ms.at(0).start).to.equal(0)
    expect(slice0ms.at(-1).end).to.equal(44100)

    expect(slice0ms[0].value).to.be.closeTo(25.3, 0.00001)
    expect(slice0ms[1879].value).to.be.closeTo(101402.4, 0.00001)

    expect(slice1ms.at(0).start).to.equal(0)
    expect(slice1ms.at(-1).end).to.equal(44100)

    expect(slice1ms[0].value).to.be.closeTo(25.3, 0.00001)
    expect(slice1ms[1879].value).to.be.closeTo(101402.4, 0.00001)
  })
})

// COMPUTE SHADER
// let samples = u32(uconstants.samples);
// let pixels = u32(uconstants.pixels);
// let stride = f32(uconstants.stride);
// let start = u32(round(f32(pixel.x) * stride));
// let end = u32(round(f32(pixel.x + 1) * stride));
//
// var p = f32(0);
// var q = f32(0);
//
// var m = u32(0);
// var n = u32(0);
//
// for (var i: u32 = start; i < end; i++) {
//    let v = audio[i];
//
//    if (v > 0.0) {
//       p += v; m++;
//    } else if (v < 0.0) {
//       q += v; n++;
//    } else {
//       p += v; m++;
//       q += v; n++;
//    }
// }
//
// if (m > u32(0)) {
//    p = p/f32(m);
// }
//
// if (n > u32(0)) {
//    q = q/f32(n);
// }
//
// waveform[pixel.x] = vec2f(p,q);
