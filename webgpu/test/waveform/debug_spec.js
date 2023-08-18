import { describe, it } from 'mocha'
import { expect } from 'chai'

describe('audio stride', function () {
  it.skip('stride (old) duration:5s, fs:44100, width:1920, padding:20', function () {
    const duration = 5
    const fs = 44100
    const N = duration * fs
    const width = 1920
    const padding = 20

    const samples = new Uint32Array(N)
    const pixels = width - 2 * padding
    const stride = samples.length / pixels

    let start = 0
    let end = 0
    for (let x = 0; x < pixels; x++) {
      start = x * stride
      end = start + stride

      for (let i = start; i < end; i++) {
        samples[i] += 1
      }
    }

    const delta = samples.length - end
    const percentage = 100.0 * delta / samples.length

    console.log(`delta: ${delta}  ${percentage.toFixed(2)}%`)

    expect(samples.slice(0, end).every((e) => e === 1)).to.equal(true)
    expect(end).to.equal(N)
    expect(Math.abs(percentage)).to.be.lt(0.1)
  })

  it('stride duration:5s, fs:44100, width:1920, padding:20', function () {
    const duration = 5
    const fs = 44100
    const N = duration * fs
    const width = 1920
    const padding = 20

    const samples = new Uint32Array(N)
    const pixels = width - 2 * padding
    const stride = samples.length / pixels

    let start = 0
    let end = 0
    for (let x = 0; x < pixels; x++) {
      start = Math.round(x * stride)
      end = Math.round((x + 1) * stride)

      for (let i = start; i < end; i++) {
        samples[i] += 1
      }
    }

    const delta = samples.length - end
    const percentage = 100.0 * delta / samples.length

    console.log(`    - delta: ${delta}  ${percentage.toFixed(2)}%`)

    expect(samples.slice(0, end).every((e) => e === 1)).to.equal(true)
    expect(end).to.equal(N)
    expect(Math.abs(percentage)).to.be.lt(0.1)
  })

  it.skip('stride (old) duration:0.1s, fs:44100, width:1920, padding:20', function () {
    const duration = 0.1
    const fs = 44100
    const N = duration * fs
    const width = 1920
    const padding = 20

    const samples = new Uint32Array(N)
    const pixels = width - 2 * padding
    const stride = Math.floor(samples.length / pixels)

    let start = 0
    let end = 0
    for (let x = 0; x < pixels; x++) {
      start = x * stride
      end = start + stride

      for (let i = start; i < end; i++) {
        samples[i] += 1
      }
    }

    const delta = samples.length - end
    const percentage = 100.0 * delta / samples.length

    console.log(`    - delta: ${delta}  ${percentage.toFixed(2)}%`)

    expect(samples.slice(0, end).every((e) => e === 1)).to.equal(true)
    expect(end).to.equal(N)
    expect(Math.abs(percentage)).to.be.lt(0.1)
  })

  it('stride duration:0.1s, fs:44100, width:1920, padding:20', function () {
    const duration = 0.1
    const fs = 44100
    const N = duration * fs
    const width = 1920
    const padding = 20

    const samples = new Uint32Array(N)
    const pixels = width - 2 * padding
    const stride = samples.length / pixels

    let start = 0
    let end = 0
    for (let x = 0; x < pixels; x++) {
      start = Math.round(x * stride)
      end = Math.round((x + 1) * stride)

      for (let i = start; i < end; i++) {
        samples[i] += 1
      }
    }

    const delta = samples.length - end
    const percentage = 100.0 * delta / samples.length

    console.log(`    - delta: ${delta}  ${percentage.toFixed(2)}%`)

    expect(samples.slice(0, end).every((e) => e === 1)).to.equal(true)
    expect(end).to.equal(N)
    expect(Math.abs(percentage)).to.be.lt(0.1)
  })
})
