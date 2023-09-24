/* global GPUBufferUsage,GPUShaderStage */

const PADDING = 20
const WORKGROUP_SIZE = 64

export function gradient3 (device, format, a, width, height, vscale, midpoint, colour1, colour2, colour3) {
  const xscale = (width - 2 * PADDING) / width
  const yscale = (height - 2 * PADDING) / height
  const samples = a.audio.subarray(a.start, a.end)
  const pixels = Math.min(width - 2 * PADDING, samples.length)
  const stride = samples.length / pixels

  const vertices = new Float32Array([
    0.0, +1.0,
    0.0, +1.0 * midpoint,
    0.0, 0.0,
    0.0, -1.0 * midpoint,
    0.0, -1.0
  ])

  const vertexBuffer = device.createBuffer({
    label: 'waveform vertices',
    size: vertices.byteLength,
    usage: GPUBufferUsage.VERTEX | GPUBufferUsage.COPY_DST
  })

  device.queue.writeBuffer(vertexBuffer, 0, vertices)

  const vertexBufferLayout = {
    arrayStride: 8,
    attributes: [
      {
        format: 'float32x2',
        offset: 0,
        shaderLocation: 0
      }
    ]
  }

  const bindGroupLayout = device.createBindGroupLayout({
    label: 'waveform bind group Layout',
    entries: [
      {
        binding: 0,
        visibility: GPUShaderStage.VERTEX | GPUShaderStage.COMPUTE,
        buffer: { type: 'uniform' }
      },
      {
        binding: 1,
        visibility: GPUShaderStage.VERTEX | GPUShaderStage.COMPUTE,
        buffer: { type: 'read-only-storage' }
      },
      {
        binding: 2,
        visibility: GPUShaderStage.COMPUTE,
        buffer: { type: 'storage' }
      }
    ]
  })

  const pipelineLayout = device.createPipelineLayout({
    label: 'waveform pipeline layout',
    bindGroupLayouts: [bindGroupLayout]
  })

  const shader = device.createShaderModule({
    label: 'waveform shader',
    code: SHADER
  })

  const computer = device.createShaderModule({
    label: 'compute shader to accumulate and average audio samples',
    code: COMPUTE
  })

  const renderPipeline = device.createRenderPipeline({
    label: 'render pipeline',
    layout: pipelineLayout,
    vertex: {
      module: shader,
      entryPoint: 'vertexMain',
      buffers: [vertexBufferLayout]
    },
    fragment: {
      module: shader,
      entryPoint: 'fragmentMain',
      targets: [
        {
          format,
          blend: {
            operation: 'add',
            alpha: {
              srcFactor: 'one',
              dstFactor: 'one-minus-src-alpha'
            },
            color: {
              srcFactor: 'src-alpha',
              dstFactor: 'one-minus-src-alpha'
            }
          }
        }
      ]
    },
    primitive: {
      topology: 'line-strip'
    }
  })

  const computePipeline = device.createComputePipeline({
    label: 'compute pipeline',
    layout: pipelineLayout,
    compute: {
      module: computer,
      entryPoint: 'computeMain'
    }
  })

  const constants = pack({
    pixels,
    stride,
    samples: samples.length,
    xscale,
    yscale,
    vscale,
    colours: [colour3, colour2, colour1, colour2, colour3]
  })

  const waveform = new Float32Array(pixels * 2)
  const audio = new Float32Array(samples)

  const storage = {
    uniforms: device.createBuffer({
      label: 'waveform:constants',
      size: 16 * Math.ceil(constants.byteLength / 16),
      usage: GPUBufferUsage.UNIFORM | GPUBufferUsage.COPY_DST
    }),

    waveform: device.createBuffer({
      label: 'waveform:waveform',
      size: 16 * Math.ceil(waveform.byteLength / 16),
      usage: GPUBufferUsage.STORAGE | GPUBufferUsage.COPY_DST
    }),

    audio: device.createBuffer({
      label: 'waveform:audio',
      size: 16 * Math.ceil(audio.byteLength / 16),
      usage: GPUBufferUsage.STORAGE | GPUBufferUsage.COPY_DST
    })
  }

  device.queue.writeBuffer(storage.uniforms, 0, constants)
  device.queue.writeBuffer(storage.waveform, 0, waveform)
  device.queue.writeBuffer(storage.audio, 0, audio)

  const bindGroups = {
    render: device.createBindGroup({
      label: 'waveform renderer bind group',
      layout: bindGroupLayout,
      entries: [
        { binding: 0, resource: { buffer: storage.uniforms } },
        { binding: 1, resource: { buffer: storage.waveform } },
        { binding: 2, resource: { buffer: storage.audio } }
      ]
    }),

    compute: device.createBindGroup({
      label: 'waveform compute bind group',
      layout: bindGroupLayout,
      entries: [
        { binding: 0, resource: { buffer: storage.uniforms } },
        { binding: 1, resource: { buffer: storage.audio } },
        { binding: 2, resource: { buffer: storage.waveform } }
      ]
    })
  }

  const compute = function (pass) {
    const workgroups = Math.ceil(pixels / WORKGROUP_SIZE)
    const bindGroup = bindGroups.compute

    pass.setPipeline(computePipeline)
    pass.setBindGroup(0, bindGroup)
    pass.dispatchWorkgroups(workgroups)
  }

  const render = function (pass) {
    const bindGroup = bindGroups.render

    pass.setPipeline(renderPipeline)
    pass.setVertexBuffer(0, vertexBuffer)
    pass.setBindGroup(0, bindGroup)
    pass.draw(vertices.length / 2, pixels)
  }

  return { compute, render }
}

function pack ({ pixels, stride, samples, xscale, yscale, vscale, colours }) {
  const pad = 0
  const buffer = new ArrayBuffer(32 + colours.length * 16)
  const view = new DataView(buffer)

  view.setUint32(0, pixels, true)
  view.setFloat32(4, stride, true)
  view.setUint32(8, samples, true)
  view.setUint32(12, pad, true)
  view.setFloat32(16, xscale, true) // vec2f: must be on a 16-byte boundary
  view.setFloat32(20, yscale, true) //
  view.setFloat32(24, vscale, true)

  for (let i = 0; i < colours.length; i++) {
    const colour = colours[i]
    const ix = 32 + i * 16

    view.setFloat32(ix + 0, colour[0], true) // vec4f: must be on a 16-byte boundary
    view.setFloat32(ix + 4, colour[1], true) //
    view.setFloat32(ix + 8, colour[2], true) //
    view.setFloat32(ix + 12, colour[3], true) //
  }

  return new Uint8Array(buffer)
}

const SHADER = `
    struct constants {
      pixels: u32,
      stride: f32,
      samples: u32,
      pad: f32,
      scale: vec2f,
      vscale: f32,
      colours: array<vec4f,5>,
    };

    struct VertexInput {
        @location(0)             pos: vec2f,
        @builtin(instance_index) instance: u32,
        @builtin(vertex_index)   vertex: u32,
    };

    struct VertexOutput {
        @builtin(position) pos: vec4f,
        @location(0) colour: vec4f,
    };

    @group(0) @binding(0) var<uniform> uconstants: constants;
    @group(0) @binding(1) var<storage> waveform: array<vec2<f32>>;

    @vertex
    fn vertexMain(input: VertexInput) -> VertexOutput {
       var output: VertexOutput;

       let i = f32(input.instance);
       let scale = uconstants.scale;
       let vscale = uconstants.vscale;
       let w = f32(uconstants.pixels - u32(1));

       let height = vscale * abs(waveform[input.instance]);
       let origin = vec2f(-1.0, 0.0);
       let offset = origin + 2.0*i/w;
       let x = input.pos.x + offset.x;
       let y = clamp(input.pos.y*height[input.vertex/4],-1.0,1.0);

       output.pos = vec4f(scale.x*x, scale.y*y, 0.0, 1.0);
       output.colour = uconstants.colours[input.vertex];

       return output;
    }

    @fragment
    fn fragmentMain(input: VertexOutput) -> @location(0) vec4f {
       return input.colour;
    }
`

const COMPUTE = `
    struct constants {
      pixels: u32,
      stride: f32,
      samples: u32,
      pad: f32,
      scale: vec2f,
      vscale: f32,
      colours: array<vec4f,5>
    };

    @group(0) @binding(0) var<uniform> uconstants: constants;
    @group(0) @binding(1) var<storage> audio: array<f32>;
    @group(0) @binding(2) var<storage, read_write> waveform: array<vec2<f32>>;

    @compute  @workgroup_size(${WORKGROUP_SIZE})
    fn computeMain(@builtin(global_invocation_id) pixel: vec3u) {
       let samples = u32(uconstants.samples);
       let pixels = u32(uconstants.pixels);
       let stride = f32(uconstants.stride);
       let start = u32(round(f32(pixel.x) * stride));
       let end = u32(round(f32(pixel.x + 1) * stride));

       var p = f32(0);
       var q = f32(0);

       var m = u32(0);
       var n = u32(0);

       for (var i: u32 = start; i < end; i++) {
          let v = audio[i];

          if (v > 0.0) {
             p += v; m++;
          } else if (v < 0.0) {
             q += v; n++;
          } else {
             p += v; m++;
             q += v; n++;
          }
       }

       if (m > u32(0)) {
          p = p/f32(m);
       }

       if (n > u32(0)) {
          q = q/f32(n);
       }

       waveform[pixel.x] = vec2f(p,q);
    }
`
