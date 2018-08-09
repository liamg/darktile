package gui

import (
	"fmt"
	"strings"

	"github.com/go-gl/gl/all-core/gl" // OR: github.com/go-gl/gl/v2.1/gl
)

const (
	vertexShaderSource = `
		#version 410
		in vec3 vp;
		in vec3 inColour;
		smooth out vec3 theColour;
		void main() {
			gl_Position = vec4(vp, 1.0);
			theColour = inColour;
		}
	` + "\x00"

	fragmentShaderSource = `
		#version 410
		smooth in vec3 theColour;
		out vec4 outColour;
		void main() {
			outColour = vec4(theColour, 1.0);
		}
	` + "\x00"
)

func compileShader(source string, shaderType uint32) (uint32, error) {
	shader := gl.CreateShader(shaderType)

	csources, free := gl.Strs(source)
	gl.ShaderSource(shader, 1, csources, nil)
	free()
	gl.CompileShader(shader)

	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log))

		return 0, fmt.Errorf("failed to compile %v: %v", source, log)
	}

	return shader, nil
}
