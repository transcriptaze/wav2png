/* global GPUBufferUsage,GPUShaderStage */

const GRID_SIZE = 32
const REDRAW = 2500

export function cells (context, device, format) {
  const background = { r: 0.2, g: 0.5, b: 0.4, a: 1 }

  const vertices = new Float32Array([
    -0.8, -0.8,
    0.8, -0.8,
    0.8, 0.8,

    -0.8, -0.8,
    0.8, 0.8,
    -0.8, 0.8
  ])

  const vertexBuffer = device.createBuffer({
    label: 'Cell vertices',
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
    label: 'Cell Bind Group Layout',
    entries: [
      {
        binding: 0,
        visibility: GPUShaderStage.VERTEX | GPUShaderStage.COMPUTE,
        buffer: { type: 'uniform' } // Grid uniform buffer
      },
      {
        binding: 1,
        visibility: GPUShaderStage.VERTEX | GPUShaderStage.COMPUTE,
        buffer: { type: 'read-only-storage' } // Cell state input buffer
      },
      {
        binding: 2,
        visibility: GPUShaderStage.COMPUTE,
        buffer: { type: 'storage' } // Cell state output buffer
      }
    ]
  })

  const pipelineLayout = device.createPipelineLayout({
    label: 'Cell Pipeline Layout',
    bindGroupLayouts: [bindGroupLayout]
  })

  const cellShaderModule = device.createShaderModule({
    label: 'Cell shader',
    code: `
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

    @group(0) @binding(0) var<uniform> grid: vec2f;
    @group(0) @binding(1) var<storage> cellState: array<u32>;

    @vertex
    fn vertexMain(input: VertexInput) -> VertexOutput {
       var output: VertexOutput;

       let i = f32(input.instance);
       let cell = vec2f(i % grid.x, floor(i / grid.x));
       let state = f32(cellState[input.instance]);
       let cellOffset = cell / grid * 2; 
       let gridPos = (input.pos*state + 1) / grid - 1 + cellOffset;

       output.pos = vec4f(gridPos, 0, 1);
       output.cell = cell/grid;

       return output;
    }

    @fragment
    fn fragmentMain(input: FragInput) -> @location(0) vec4f {
       return vec4f(input.cell, 1.0 - input.cell.x, 1);
    }
`
  })

  const cellPipeline = device.createRenderPipeline({
    label: 'Cell pipeline',
    layout: pipelineLayout,
    vertex: {
      module: cellShaderModule,
      entryPoint: 'vertexMain',
      buffers: [vertexBufferLayout]
    },
    fragment: {
      module: cellShaderModule,
      entryPoint: 'fragmentMain',
      targets: [
        { format }
      ]
    }
  })

  const uniformArray = new Float32Array([GRID_SIZE, GRID_SIZE])
  const uniformBuffer = device.createBuffer({
    label: 'Grid Uniforms',
    size: uniformArray.byteLength,
    usage: GPUBufferUsage.UNIFORM | GPUBufferUsage.COPY_DST
  })

  device.queue.writeBuffer(uniformBuffer, 0, uniformArray)

  const cellStateArray = new Uint32Array(GRID_SIZE * GRID_SIZE)
  const cellStateStorage = [
    device.createBuffer({
      label: 'Cell State A',
      size: cellStateArray.byteLength,
      usage: GPUBufferUsage.STORAGE | GPUBufferUsage.COPY_DST
    }),
    device.createBuffer({
      label: 'Cell State B',
      size: cellStateArray.byteLength,
      usage: GPUBufferUsage.STORAGE | GPUBufferUsage.COPY_DST
    })
  ]

  for (let i = 0; i < cellStateArray.length; i += 3) {
    cellStateArray[i] = 1
  }

  device.queue.writeBuffer(cellStateStorage[0], 0, cellStateArray)

  for (let i = 0; i < cellStateArray.length; i++) {
    cellStateArray[i] = i % 2
  }

  device.queue.writeBuffer(cellStateStorage[1], 0, cellStateArray)

  const bindGroups = [
    device.createBindGroup({
      label: 'Cell renderer bind group A',
      layout: bindGroupLayout, // cellPipeline.getBindGroupLayout(0),
      entries: [
        { binding: 0, resource: { buffer: uniformBuffer } },
        { binding: 1, resource: { buffer: cellStateStorage[0] } },
        { binding: 2, resource: { buffer: cellStateStorage[1] } }
      ]
    }),
    device.createBindGroup({
      label: 'Cell renderer bind group B',
      layout: bindGroupLayout, // cellPipeline.getBindGroupLayout(0),
      entries: [
        { binding: 0, resource: { buffer: uniformBuffer } },
        { binding: 1, resource: { buffer: cellStateStorage[1] } },
        { binding: 2, resource: { buffer: cellStateStorage[0] } }
      ]
    })
  ]

  const compute = function (pass, step) {
  }

  const render = function (pass, step) {
    const index = step % 2

    pass.setPipeline(cellPipeline)
    pass.setVertexBuffer(0, vertexBuffer)
    pass.setBindGroup(0, bindGroups[index])
    pass.draw(vertices.length / 2, GRID_SIZE * GRID_SIZE)
  }

  return { background, compute, render, interval: REDRAW }
}
