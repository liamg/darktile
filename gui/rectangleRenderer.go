package gui

import (
	"github.com/go-gl/gl/all-core/gl"
	"github.com/liamg/aminal/config"
)

const (
	rectangleRendererVertexShaderSource = `
		#version 330 core
		layout (location = 0) in vec2 position;
		uniform vec2 resolution;

		void main() {
			// convert from window coordinates to GL coordinates
			vec2 glCoordinates = ((position / resolution) * 2.0 - 1.0) * vec2(1, -1);

			gl_Position = vec4(glCoordinates, 0.0, 1.0);
		}` + "\x00"

	rectangleRendererFragmentShaderSource = `
		#version 330 core
		uniform vec4 inColor;
		out vec4 outColor;
		void main() {
			outColor = inColor;
		}` + "\x00"
)

type rectangleRenderer struct {
	program                   uint32
	vbo                       uint32
	vao                       uint32
	ibo                       uint32
	uniformLocationResolution int32
	uniformLocationInColor    int32
}

func createRectangleRendererProgram() (uint32, error) {
	vertexShader, err := compileShader(rectangleRendererVertexShaderSource, gl.VERTEX_SHADER)
	if err != nil {
		return 0, err
	}
	defer gl.DeleteShader(vertexShader)

	fragmentShader, err := compileShader(rectangleRendererFragmentShaderSource, gl.FRAGMENT_SHADER)
	if err != nil {
		return 0, err
	}
	defer gl.DeleteShader(fragmentShader)

	prog := gl.CreateProgram()
	gl.AttachShader(prog, vertexShader)
	gl.AttachShader(prog, fragmentShader)
	gl.LinkProgram(prog)

	return prog, nil
}

func newRectangleRenderer() (*rectangleRenderer, error) {
	prog, err := createRectangleRendererProgram()
	if err != nil {
		return nil, err
	}

	var vbo uint32
	var vao uint32
	var ibo uint32

	gl.GenBuffers(1, &vbo)
	gl.GenVertexArrays(1, &vao)
	gl.GenBuffers(1, &ibo)

	vertices := [12]float32{}

	indices := [...]uint32{
		0, 1, 2,
		2, 3, 0,
	}

	gl.BindVertexArray(vao)

	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4, gl.Ptr(&vertices[0]), gl.DYNAMIC_DRAW)

	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, ibo)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(indices)*4, gl.Ptr(&indices[0]), gl.DYNAMIC_DRAW)

	gl.VertexAttribPointer(0, 2, gl.FLOAT, false, 2*4, nil)
	gl.EnableVertexAttribArray(0)

	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	gl.BindVertexArray(0)

	return &rectangleRenderer{
		program:                   prog,
		vbo:                       vbo,
		vao:                       vao,
		ibo:                       ibo,
		uniformLocationResolution: gl.GetUniformLocation(prog, gl.Str("resolution\x00")),
		uniformLocationInColor:    gl.GetUniformLocation(prog, gl.Str("inColor\x00")),
	}, nil
}

func (rr *rectangleRenderer) Free() {
	if rr.program != 0 {
		gl.DeleteProgram(rr.program)
		rr.program = 0
	}

	if rr.vbo != 0 {
		gl.DeleteBuffers(1, &rr.vbo)
		rr.vbo = 0
	}

	if rr.vao != 0 {
		gl.DeleteBuffers(1, &rr.vao)
		rr.vao = 0
	}

	if rr.ibo != 0 {
		gl.DeleteBuffers(1, &rr.ibo)
		rr.ibo = 0
	}
}

func (rr *rectangleRenderer) render(left float32, top float32, width float32, height float32, colour config.Colour) {
	var savedProgram int32
	gl.GetIntegerv(gl.CURRENT_PROGRAM, &savedProgram)
	defer gl.UseProgram(uint32(savedProgram))

	currentViewport := [4]int32{}
	gl.GetIntegerv(gl.VIEWPORT, &currentViewport[0])

	gl.UseProgram(rr.program)
	gl.Uniform2f(rr.uniformLocationResolution, float32(currentViewport[2]), float32(currentViewport[3]))

	gl.Uniform4f(rr.uniformLocationInColor, colour[0], colour[1], colour[2], 1.0)

	vertices := [...]float32{
		left, top,
		left + width, top,
		left + width, top + height,
		left, top + height,
	}

	gl.NamedBufferSubData(rr.vbo, 0, len(vertices)*4, gl.Ptr(&vertices[0]))
	gl.BindVertexArray(rr.vao)

	gl.DrawElements(gl.TRIANGLES, 6, gl.UNSIGNED_INT, gl.PtrOffset(0))

	gl.BindVertexArray(0)
}
