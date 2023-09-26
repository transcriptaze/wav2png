import { describe, it } from 'mocha'
import { expect } from 'chai'

describe('audio start/end', function () {
  it.only('debugging audio start/end logic', function () {
    // const duration = 5
    // const fs = 44100
    const width = 1920
    const padding = 20
    const start = 0
    const end = 220500
    const N = end - start

    // const samples = duration * fs
    const pixels = width - 2 * padding
    const stride = N / pixels

    const list = []
    for (let i = 0; i < pixels; i++) {
      const _start = Math.round(i * stride)
      const _end = Math.round((i + 1) * stride)

      list.push([_start, _end, _end - _start])
    }

    expect(pixels).to.equal(1880)
    expect(list.length).to.equal(pixels)

    // ... check first/last indices
    expect(list.at(0)[0]).to.equal(start)
    expect(list.at(-1)[1]).to.equal(end)

    // ... check index continuity
    for (let i = 1; i < list.length; i++) {
      const p = list[i - 1]
      const q = list[i]

      expect(q[0]).to.equal(p[1])
    }

    // ... check index stride
    const floor = Math.floor(stride)
    const ceil = Math.ceil(stride)

    for (const v of list) {
      expect(v[2]).to.be.oneOf([floor, ceil])
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
