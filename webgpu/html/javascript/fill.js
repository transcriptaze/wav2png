/* global GPUBufferUsage,GPUShaderStage */

export function fill (context, device, format, colour) {
  const vertices = new Float32Array([
    -1.0, -1.0,
    +1.0, -1.0,
    +1.0, +1.0,

    -1.0, -1.0,
    +1.0, +1.0,
    -1.0, +1.0
  ])

  const vertexBuffer = device.createBuffer({
    label: 'background vertices',
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
    label: 'background bind group layout',
    entries: [
      {
        binding: 0,
        visibility: GPUShaderStage.VERTEX | GPUShaderStage.FRAGMENT,
        buffer: { type: 'uniform' }
      }
    ]
  })

  const pipelineLayout = device.createPipelineLayout({
    label: 'background pipeline Layout',
    bindGroupLayouts: [bindGroupLayout]
  })

  const shader = device.createShaderModule({
    label: 'background shader',
    code: SHADER
  })

  const pipeline = device.createRenderPipeline({
    label: 'background pipeline',
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
        { format }
      ]
    },
    primitive: {
      topology: 'triangle-list'
    }
  })

  const constants = new Float32Array(colour)
  const uniformBuffer = device.createBuffer({
    label: 'background constants',
    size: constants.byteLength,
    usage: GPUBufferUsage.UNIFORM | GPUBufferUsage.COPY_DST
  })

  device.queue.writeBuffer(uniformBuffer, 0, constants)

  const bindGroup = device.createBindGroup({
    label: 'background renderer bind group',
    layout: bindGroupLayout,
    entries: [
      { binding: 0, resource: { buffer: uniformBuffer } }
    ]
  })

  const compute = function (pass) {
  }

  const render = function (pass) {
    pass.setPipeline(pipeline)
    pass.setVertexBuffer(0, vertexBuffer)
    pass.setBindGroup(0, bindGroup)
    pass.draw(vertices.length / 2)
  }

  return { compute, render }
}

const SHADER = `
    struct constants {
      colour: vec4f
    };

    struct VertexInput {
        @location(0) xy: vec2f
    };

    struct VertexOutput {
        @builtin(position) pos: vec4f
    };

    @group(0) @binding(0) var<uniform> uconstants: constants;

    @vertex
    fn vertexMain(input: VertexInput) -> VertexOutput {
       var output: VertexOutput;

       output.pos = vec4f(input.xy, 0.0, 1.0);

       return output;
    }

    @fragment
    fn fragmentMain(input: VertexOutput) -> @location(0) vec4f {
       return uconstants.colour;
    }
`
