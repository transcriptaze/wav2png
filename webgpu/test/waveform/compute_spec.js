import { describe, it } from 'mocha'
import { expect } from 'chai'

const width = 1920
const padding = 20

const audio = {
  fs: 44100,
  duration: 2,
  audio: new Float64Array(2 * 44100)
}

for (let i = 0; i < audio.audio.length; i++) {
  audio.audio[i] = Math.round(i * 0.1 * 10) / 10
}

describe('audio pixel bucket logic', function () {
  it.only('100ms slice of 2s audio from 1.000s to 1.100ms', function () {
    // ... setup
    const from = 1.000
    const to = 1.100
    const N = audio.audio.length
    const duration = clamp(audio.duration, 0, N / audio.fs)
    const start = duration === 0 ? 0 : clamp(Math.floor(N * from / duration), 0, N)
    const end = duration === 0 ? 0 : clamp(Math.floor(N * to / duration), 0, N)

    expect(N).to.equal(88200)
    expect(duration).to.equal(2)
    expect(start).to.equals(44100)
    expect(end).to.equal(48510)

    // ... line
    const line = {}

    line.start = start
    line.end = end
    line.N = line.end - line.start
    line.pixels = Math.min(width - 2 * padding, line.N)
    line.stride = line.N / line.pixels
    line.samples = audio.audio.subarray(line.start, line.end)

    expect(line.start).to.equal(44100)
    expect(line.end).to.equal(48510)
    expect(line.N).to.equal(4410)
    expect(line.pixels).to.equal(1880)
    expect(line.stride).to.be.closeTo(2.345744680851064, 0.0001)
    expect(line.samples.length).to.equal(4410)

    // ... compute shader
    const compute = {
      pixels: line.pixels,
      samples: line.samples.length,
      stride: Math.fround(line.stride) // f32
    }

    expect(compute.pixels).to.equal(1880)
    expect(compute.samples).to.equal(4410)
    expect(compute.stride).to.be.closeTo(2.3457446098327637, 0.0001)

    // ... pixel buckets
    const indices = new Map([
      [0, 0],
      [1, 2],
      [2, 5],
      [3, 7],
      [4, 9],
      [5, 12]
    ])

    const buckets = new Map([
      [0, [4410.0, 4410.1]],
      [1, [4410.2, 4410.3, 4410.4]],
      [2, [4410.5, 4410.6]],
      [3, [4410.7, 4410.8]],
      [4, [4410.9, 4411.0, 4411.1]]
    ])

    for (let x = 0; x < 5; x++) {
      const start = Math.round(x * compute.stride)
      const end = Math.round((x + 1) * compute.stride)
      const bucket = []

      for (let i = start; i < end; i++) {
        bucket.push(line.samples[i])
      }

      expect(start).to.equal(indices.get(x))
      expect(end).to.equal(indices.get(x + 1))
      expect(bucket).to.eql(buckets.get(x))
    }
  })

  it.only('100ms slice of 2s audio from 1.000.1s to 1.100.1ms', function () {
    // ... setup
    const from = 1.0001
    const to = 1.1001
    const N = audio.audio.length
    const duration = clamp(audio.duration, 0, N / audio.fs)
    const start = duration === 0 ? 0 : clamp(Math.floor(N * from / duration), 0, N)
    const end = duration === 0 ? 0 : clamp(Math.floor(N * to / duration), 0, N)

    expect(N).to.equal(88200)
    expect(duration).to.equal(2)
    expect(start).to.equals(44104)
    expect(end).to.equal(48514)

    // ... line
    const line = {}

    line.start = start
    line.end = end
    line.N = line.end - line.start
    line.pixels = Math.min(width - 2 * padding, line.N)
    line.stride = line.N / line.pixels
    line.samples = audio.audio.subarray(line.start, line.end)

    expect(line.start).to.equal(44104)
    expect(line.end).to.equal(48514)
    expect(line.N).to.equal(4410)
    expect(line.pixels).to.equal(1880)
    expect(line.stride).to.be.closeTo(2.345744680851064, 0.0001)
    expect(line.samples.length).to.equal(4410)

    // ... compute shader
    const compute = {
      pixels: line.pixels,
      samples: line.samples.length,
      stride: Math.fround(line.stride) // f32
    }

    expect(compute.pixels).to.equal(1880)
    expect(compute.samples).to.equal(4410)
    expect(compute.stride).to.be.closeTo(2.3457446098327637, 0.0001)

    const indices = new Map([
      [0, 0],
      [1, 2],
      [2, 5],
      [3, 7],
      [4, 9],
      [5, 12]
    ])

    const buckets = new Map([
      [0, [4410.4, 4410.5]],
      [1, [4410.6, 4410.7, 4410.8]],
      [2, [4410.9, 4411.0]],
      [3, [4411.1, 4411.2]],
      [4, [4411.3, 4411.4, 4411.5]]
    ])

    for (let x = 0; x < 5; x++) {
      const start = Math.round(x * compute.stride)
      const end = Math.round((x + 1) * compute.stride)
      const bucket = []

      for (let i = start; i < end; i++) {
        bucket.push(line.samples[i])
      }

      expect(start).to.equal(indices.get(x))
      expect(end).to.equal(indices.get(x + 1))
      expect(bucket).to.eql(buckets.get(x))
    }
  })
})

function clamp (v, min, max) {
  return (v < min) ? min : ((v > max) ? max : v)
}

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
