/* global GPUBufferUsage,GPUShaderStage */

const GRIDX = 8
const GRIDY = 4
const PADDING = 20

export function grid (context, device, format) {
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
    code: `
    struct constants {
      grid: vec2f,
      scale: vec2f
    };

    struct VertexInput {
        @location(0) pos: vec2f,
        @builtin(instance_index) instance: u32,
    };

    struct VertexOutput {
        @builtin(position) pos: vec4f,
        @location(0) cell: vec2f, 
    };

    struct FragInput {
        @location(0) cell: vec2f,
    };

    @group(0) @binding(0) var<uniform> uconstants: constants;

    @vertex
    fn vertexMain(input: VertexInput) -> VertexOutput {
       var output: VertexOutput;

       let i = f32(input.instance);
       let grid = uconstants.grid;
       let scale = uconstants.scale;

       let cell = vec2f(i % grid.x, floor(i / grid.x));
       let cellOffset = 2*cell / grid; 
       let gridPos = (input.pos + 1) / grid - 1 + cellOffset;

       output.pos = vec4f(gridPos, 0, 1) * vec4f(scale,1,1);
       output.cell = cell/grid;

       return output;
    }

    @fragment
    fn fragmentMain(input: FragInput) -> @location(0) vec4f {
       return vec4f(0, 1, 0, 1);
    }
`
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
        { format }
      ]
    },
    primitive: {
      topology: 'line-strip'
    }
  })

  const constants = new Float32Array([GRIDX, GRIDY, xscale, yscale])
  const uniformBuffer = device.createBuffer({
    label: 'grid constants',
    size: constants.byteLength,
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
    pass.draw(vertices.length / 2, GRIDX * GRIDY)
  }

  return { compute, render }
}
