import { describe, it } from 'mocha'
import { expect } from 'chai'

const width = 1920
const padding = 20

const audio = {
  fs: 44100,
  duration: 2,
  audio: new Float64Array(2 * 44100),
  start: 1.0,
  end: 1.1
}

for (let i = 0; i < audio.audio.length; i++) {
  audio.audio[i] = Math.round(i * 0.1 * 10) / 10
}

describe('audio pixel bucket logic', function () {
  it.only('100ms slice of 2s audio from 1.000s to 1.100ms', function () {
    // ... setup
    const N = audio.audio.length
    const duration = clamp(audio.duration, 0, N / audio.fs)
    const start = duration === 0 ? 0 : clamp(Math.floor(N * audio.start / duration), 0, N)
    const end = duration === 0 ? 0 : clamp(Math.floor(N * audio.end / duration), 0, N)

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

    // ... pixel 1
    {
      const x = 0
      const start = Math.round(x * compute.stride)
      const end = Math.round((x + 1) * compute.stride)
      const bucket = []

      for (let i = start; i < end; i++) {
        bucket.push(line.samples[i])
      }

      expect(start).to.equal(0)
      expect(end).to.equal(2)
      expect(bucket).to.eql([4410.0, 4410.1])
    }

    // ... pixel 2
    {
      const x = 1
      const start = Math.round(x * compute.stride)
      const end = Math.round((x + 1) * compute.stride)
      const bucket = []

      for (let i = start; i < end; i++) {
        bucket.push(line.samples[i])
      }

      expect(start).to.equal(2)
      expect(end).to.equal(5)
      expect(bucket).to.eql([4410.2, 4410.3, 4410.4])
    }

    // ... pixel 3
    {
      const x = 2
      const start = Math.round(x * compute.stride)
      const end = Math.round((x + 1) * compute.stride)
      const bucket = []

      for (let i = start; i < end; i++) {
        bucket.push(line.samples[i])
      }

      expect(start).to.equal(5)
      expect(end).to.equal(7)
      expect(bucket).to.eql([4410.5, 4410.6])
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
