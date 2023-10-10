import * as quantize from '../../html/javascript/shaders/quantize.js'
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

describe.only('shader audio slice logic', function () {
  it('validate slice offsets for 100ms slice of 2s audio from 1.000.1s to 1.100.1ms, indexed from start of audio', function () {
    const slice = quantize.slice({
      fs: audio.fs,
      start: 1.0001,
      end: 1.1001,
      duration: 2,
      audio: audio.audio
    }, width, padding)

    // console.log(s)
    expect(slice.start).to.equal(44102)
    expect(slice.end).to.equal(48512)
    expect(slice.audio.length).to.equal(4410)
    expect(slice.offset).to.equal(18801)
    expect(slice.pixels).to.equal(1880)
    expect(slice.stride).to.be.closeTo(Math.fround(2.345744680851064), 0.00001)
  })
})
