import { describe, it } from 'mocha'
import { expect } from 'chai'

describe('audio pixel bucket logic', function () {
  it.only('reference: complete buffer', function () {
    const duration = 5
    const fs = 44100
    const samples = new Float64Array(duration * fs)

    for (let i = 0; i < samples.length; i++) {
      samples[i] = i * 0.1
    }

    const width = 1920
    const padding = 20
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

    expect(pixels).to.equal(1880)
    expect(buckets.length).to.equal(pixels)

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
    for (let i = 0; i < buckets.length; i++) {
      const bucket = buckets[i]
      const N = bucket.end - bucket.start
      const v = N * (samples[bucket.start] + samples[bucket.end - 1]) / 2

      expect(bucket.value, `${i}`).to.be.closeTo(v, 0.00001)
    }
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
