/* global GPUBufferUsage,GPUShaderStage */

const PADDING = 20

export function grid (context, device, format, { colour, gridx, gridy }) {
  const width = context.canvas.width
  const height = context.canvas.height
  const xscale = (width - 2 * PADDING) / width
  const yscale = (height - 2 * PADDING) / height

  const vertices = new Float32Array([
    -1.0, -1.0,
    +1.0, -1.0,
    +1.0, +1.0,
    -1.0, +1.0,
    -1.0, -1.0
  ])

  const vertexBuffer = device.createBuffer({
    label: 'grid vertices',
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
    label: 'grid bind group layout',
    entries: [
      {
        binding: 0,
        visibility: GPUShaderStage.VERTEX,
        buffer: { type: 'uniform' }
      }
    ]
  })

  const pipelineLayout = device.createPipelineLayout({
    label: 'grid pipeline Layout',
    bindGroupLayouts: [bindGroupLayout]
  })

  const shader = device.createShaderModule({
    label: 'grid shader',
    code: SHADER
  })

  const pipeline = device.createRenderPipeline({
    label: 'grid pipeline',
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

  const constants = new Float32Array([gridx, gridy, xscale, yscale, ...colour])
  const uniformBuffer = device.createBuffer({
    label: 'grid constants',
    size: 16 * Math.ceil(constants.byteLength / 16),
    usage: GPUBufferUsage.UNIFORM | GPUBufferUsage.COPY_DST
  })

  device.queue.writeBuffer(uniformBuffer, 0, constants)

  const bindGroup = device.createBindGroup({
    label: 'grid renderer bind group',
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
    pass.draw(vertices.length / 2, gridx * gridy)
  }

  return { compute, render }
}

const SHADER = `
    struct constants {
      grid: vec2f,
      scale: vec2f,
      colour: vec4f
    };

    struct VertexInput {
        @location(0) pos: vec2f,
        @builtin(instance_index) instance: u32,
    };

    struct VertexOutput {
        @builtin(position) pos: vec4f,
        @location(0) colour: vec4f, 
    };

    @group(0) @binding(0) var<uniform> uconstants: constants;

    @vertex
    fn vertexMain(input: VertexInput) -> VertexOutput {
       var output: VertexOutput;

       let i = f32(input.instance);
       let grid = uconstants.grid;
       let scale = uconstants.scale;

       let cell = vec2f(i % grid.x, floor(i / grid.x));
       let offset = 2.0*cell / grid; 
       let pos = (input.pos + 1.0) / grid - 1.0 + offset;

       output.pos = vec4f(pos, 0.0, 1.0) * vec4f(scale,1.0,1.0);
       output.colour = uconstants.colour;

       return output;
    }

    @fragment
    fn fragmentMain(input: VertexOutput) -> @location(0) vec4f {
       return input.colour;
    }
`
