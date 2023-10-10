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

describe.skip('audio pixel bucket logic', function () {
  it('100ms slice of 2s audio from 1.000s to 1.100ms', function () {
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

  it('100ms slice of 2s audio from 1.000.1s to 1.100.1ms', function () {
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

describe.only('audio pixel bucket logic', function () {
  it('100ms slice of 2s audio from 1.000s to 1.100ms, indexed from start of audio', function () {
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
    const expected = new Map([
      [0, { t: 0.000000, start: 0, end: 2, bucket: [0.0, 0.1] }],
      [1, { t: 0.000045, start: 2, end: 5, bucket: [0.2, 0.3, 0.4] }],
      [2, { t: 0.000113, start: 5, end: 7, bucket: [0.5, 0.6] }],
      [3, { t: 0.000159, start: 7, end: 9, bucket: [0.7, 0.8] }],
      [4, { t: 0.000204, start: 9, end: 12, bucket: [0.9, 1.0, 1.1] }]
    ])

    const compute = {
      stride: Math.fround(line.stride) // f32
    }

    const roundµs = function (t) {
      return Math.round(1000_000 * duration * t / N) / 1000_000
    }

    {
      const buckets = []
      let index = 0
      let start = Math.round(index * compute.stride)
      let end = Math.round((index + 1) * compute.stride)
      let t = { start: roundµs(start), end: roundµs(end) }
      let bucket = audio.audio.subarray(start, end)

      buckets.push({ index, t, start, end, bucket })

      while (t.start < from) {
        buckets.push({ index, t, start, end, bucket })

        if (expected.has(index)) {
          const e = expected.get(index)

          expect(t.start).to.be.closeTo(e.t, 0.000001)
          expect(start).to.equal(e.start)
          expect(end).to.equal(e.end)
          expect(bucket).to.eql(new Float64Array(e.bucket))
        }

        index++

        start = Math.round(index * compute.stride)
        end = Math.round((index + 1) * compute.stride)
        t = { start: roundµs(start), end: roundµs(end) }
        bucket = audio.audio.subarray(start, end)
      }

      const startʼ = buckets.at(-1)

      while (t.end <= to) {
        buckets.push({ index, t, start, end, bucket })

        if (expected.has(index)) {
          const e = expected.get(index)

          expect(t.start).to.be.closeTo(e.t, 0.000001)
          expect(start).to.equal(e.start)
          expect(end).to.equal(e.end)
          expect(bucket).to.eql(new Float64Array(e.bucket))
        }

        index++

        start = Math.round(index * compute.stride)
        end = Math.round((index + 1) * compute.stride)
        t = { start: roundµs(start), end: roundµs(end) }
        bucket = audio.audio.subarray(start, end)
      }

      const endʼ = buckets.at(-1)

      expect(startʼ.index).to.equal(18799)
      expect(startʼ.t.start).to.be.closeTo(from, 0.00005)
      expect(startʼ.start).to.equal(44098)
      expect(startʼ.end).to.equal(44100)
      expect(startʼ.bucket).to.eql(new Float64Array([4409.8, 4409.9]))

      expect(endʼ.index).to.equal(20679)
      expect(endʼ.t.end).to.be.closeTo(to, 0.00005)
      expect(endʼ.start).to.equal(48508)
      expect(endʼ.end).to.equal(48510)
      expect(endʼ.bucket).to.eql(new Float64Array([4850.8, 4850.9]))

      // for (let i = 1; i <= 10; i++) {
      //   const bucket = buckets[startʼ.index + i]
      //
      //   console.log(i, bucket.index, bucket.start, bucket.end)
      // }
    }
  })

  it('100ms slice of 2s audio from 1.000.1s to 1.100.1ms, indexed from start of audio', function () {
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
    const expected = new Map([
      [0, { t: 0.000000, start: 0, end: 2, bucket: [0.0, 0.1] }],
      [1, { t: 0.000045, start: 2, end: 5, bucket: [0.2, 0.3, 0.4] }],
      [2, { t: 0.000113, start: 5, end: 7, bucket: [0.5, 0.6] }],
      [3, { t: 0.000159, start: 7, end: 9, bucket: [0.7, 0.8] }],
      [4, { t: 0.000204, start: 9, end: 12, bucket: [0.9, 1.0, 1.1] }]
    ])

    const compute = {
      stride: Math.fround(line.stride) // f32
    }

    const roundµs = function (t) {
      return Math.round(1000_000 * duration * t / N) / 1000_000
    }

    {
      const buckets = []
      let index = 0
      let start = Math.round(index * compute.stride)
      let end = Math.round((index + 1) * compute.stride)
      let t = { start: roundµs(start), end: roundµs(end) }
      let bucket = audio.audio.subarray(start, end)

      buckets.push({ index, t, start, end, bucket })

      while (t.start < from) {
        buckets.push({ index, t, start, end, bucket })

        if (expected.has(index)) {
          const e = expected.get(index)

          expect(t.start).to.be.closeTo(e.t, 0.000001)
          expect(start).to.equal(e.start)
          expect(end).to.equal(e.end)
          expect(bucket).to.eql(new Float64Array(e.bucket))
        }

        index++

        start = Math.round(index * compute.stride)
        end = Math.round((index + 1) * compute.stride)
        t = { start: roundµs(start), end: roundµs(end) }
        bucket = audio.audio.subarray(start, end)
      }

      const startʼ = buckets.at(-1)

      while (t.end <= to) {
        buckets.push({ index, t, start, end, bucket })

        if (expected.has(index)) {
          const e = expected.get(index)

          expect(t.start).to.be.closeTo(e.t, 0.000001)
          expect(start).to.equal(e.start)
          expect(end).to.equal(e.end)
          expect(bucket).to.eql(new Float64Array(e.bucket))
        }

        index++

        start = Math.round(index * compute.stride)
        end = Math.round((index + 1) * compute.stride)
        t = { start: roundµs(start), end: roundµs(end) }
        bucket = audio.audio.subarray(start, end)
      }

      const endʼ = buckets.at(-1)

      expect(startʼ.index).to.equal(18801)
      expect(startʼ.t.start).to.be.closeTo(from, 0.00006)
      expect(startʼ.start).to.equal(44102)
      expect(startʼ.end).to.equal(44105)
      expect(startʼ.bucket).to.eql(new Float64Array([4410.2, 4410.3, 4410.4]))

      expect(endʼ.index).to.equal(20680)
      expect(endʼ.t.end).to.be.closeTo(to, 0.00006)
      expect(endʼ.start).to.equal(48510)
      expect(endʼ.end).to.equal(48512)
      expect(endʼ.bucket).to.eql(new Float64Array([4851.0, 4851.1]))

      expect(buckets.at(-2).index).to.equal(20679)
      expect(buckets.at(-2).t.end).to.be.closeTo(to - 0.0001, 0.00005)
      expect(buckets.at(-2).start).to.equal(48508)
      expect(buckets.at(-2).end).to.equal(48510)
      expect(buckets.at(-2).bucket).to.eql(new Float64Array([4850.8, 4850.9]))

      // for (let i = 1; i <= 10; i++) {
      //   const bucket = buckets[startʼ.index + i]
      //
      //   console.log(i, bucket.index, bucket.start, bucket.end)
      // }
    }
  })

  // 1 18799 44098 44100
  // 2 18800 44100 44102
  // 3 18801 44102 44105
  // 4 18802 44105 44107
  // 5 18803 44107 44109
  it('offsets: 100ms slice of 2s audio from 1.000s to 1.100ms, indexed from start of audio', function () {
    const start = 44098
    // const end = 48508
    const stride = Math.fround(2.345744680851064)
    const offset = 18799

    const expected = new Map([
      [0, { start: 0, end: 2 }],
      [1, { start: 2, end: 5 }],
      [2, { start: 5, end: 7 }],
      [3, { start: 7, end: 9 }],
      [4, { start: 9, end: 12 }]
    ])

    for (let i = 0; i < 1; i++) {
      const indices = {
        start: Math.round((offset + i) * stride) - start,
        end: Math.round((offset + i + 1) * stride) - start
      }

      expect(indices.start).to.equal(expected.get(i).start)
      expect(indices.end).to.equal(expected.get(i).end)
    }
  })

  // 1 18801 44102 44105
  // 2 18802 44105 44107
  // 3 18803 44107 44109
  // 4 18804 44109 44112
  // 5 18805 44112 44114
  it('offsets: 100ms slice of 2s audio from 1.000.1s to 1.100.1ms, indexed from start of audio', function () {
    const start = 44102
    const stride = Math.fround(2.345744680851064)
    const offset = 18801

    const expected = new Map([
      [0, { start: 0, end: 3 }],
      [1, { start: 3, end: 5 }],
      [2, { start: 5, end: 7 }],
      [3, { start: 7, end: 10 }],
      [4, { start: 10, end: 12 }]
    ])

    for (let i = 0; i < 5; i++) {
      const indices = {
        start: Math.round((offset + i) * stride) - start,
        end: Math.round((offset + i + 1) * stride) - start
      }

      // console.log(i, indices)

      expect(indices.start).to.equal(expected.get(i).start)
      expect(indices.end).to.equal(expected.get(i).end)
    }
  })
})

function clamp (v, min, max) {
  return (v < min) ? min : ((v > max) ? max : v)
}
