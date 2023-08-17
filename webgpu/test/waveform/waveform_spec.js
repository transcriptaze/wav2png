import { describe, it } from 'mocha'
import { expect } from 'chai'

describe('audio stride', function () {
  it('stride duration:5s, fs:44100, width:1920, padding:20', function () {
    const duration = 5
    const fs = 44100
    const width = 1920
    const padding = 20

    const pixels = width - 2*padding
    const samples = duration * fs
    const stride = Math.floor(samples/pixels)

    let last = 0
    for (let x=0; x<pixels; x++) {
        const start = x * stride;
        const end = start + stride

        last = end+1
    }

    const delta = samples - last
    const percentage = 100.0*delta/samples

    expect(Math.abs(percentage)).to.be.lt(0.1)
  })

  it('stride duration:0.1s, fs:44100, width:1920, padding:20', function () {
    const duration = 0.1
    const fs = 44100
    const width = 1920
    const padding = 20

    const pixels = width - 2*padding
    const samples = duration * fs
    const stride = Math.floor(samples/pixels)

    let last = 0
    for (let x=0; x<pixels; x++) {
        const start = x * stride;
        const end = start + stride

        last = end+1
    }

    const delta = samples - last
    const percentage = 100.0*delta/samples

    expect(Math.abs(percentage)).to.be.lt(0.1)
  })
})
