// Copyright 2012 The go-gl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package v41

import (
	"github.com/4ydx/gltext"
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"image"
)

var fontVertexShaderSource string = `
#version 330

uniform mat4 scale_matrix;
uniform mat4 orthographic_matrix;
uniform vec2 final_position;

in vec4 centered_position;
in vec2 uv;

out vec2 fragment_uv;

// The orthographic projection uses a lower left-hand point of (0,0)
// 1) We center the text on screen.
// 2) We perform othographic transformation and then scaling.
// 3) We move the text to its final resting place.
// This is all pretty standard I would imagine, but it took me a bit to sort out what has to happen :P

void main() {
  fragment_uv = uv;
  vec4 scaled = scale_matrix * orthographic_matrix * centered_position;
  gl_Position = vec4(scaled.x + final_position.x, scaled.y + final_position.y, scaled.z, scaled.w);
}
` + "\x00"

var fontFragmentShaderSource string = `
#version 330

uniform sampler2D fragment_texture;
uniform float fadeout;
uniform vec4 fragment_color_adjustment;

in vec2 fragment_uv;
out vec4 fragment_color;

void main() {
  vec4 color     = texture(fragment_texture, fragment_uv);
  color.xyz      = fragment_color_adjustment.xyz;
	color.w        = color.w - fadeout;
  fragment_color = color;
}
` + "\x00"

type Font struct {
	Config         *gltext.FontConfig // Character set for this font.
	textureID      uint32             // Holds the glyph texture id.
	maxGlyphWidth  int                // Largest glyph width.
	maxGlyphHeight int                // Largest glyph height.
	program        uint32             // program compiled from shaders

	// attributes
	centeredPositionAttribute uint32 // vertex centered_position required for scaling around the orthographic projections center
	uvAttribute               uint32 // texture position

	// The final screen position post-scaling
	finalPositionUniform int32

	// Position of the shaders fragment texture variable
	fragmentTextureUniform int32

	// The desired color of the text
	colorUniform   int32
	fadeoutUniform int32

	// View matrix
	orthographicMatrixUniform int32
	OrthographicMatrix        mgl32.Mat4

	// Scale the resulting text
	scaleMatrixUniform int32

	textureWidth  float32
	textureHeight float32
	WindowWidth   float32
	WindowHeight  float32
}

func (f *Font) GetTextureWidth() float32 {
	return f.textureWidth
}

func (f *Font) GetTextureHeight() float32 {
	return f.textureHeight
}

func NewFont(config *gltext.FontConfig) (f *Font, err error) {
	if config == nil {
		panic("Nil config")
	}
	f = &Font{}
	f.Config = config

	// Resize image to next power-of-two.
	config.Image = gltext.Pow2Image(config.Image).(*image.NRGBA)
	ib := config.Image.Bounds()

	f.textureWidth = float32(ib.Dx())
	f.textureHeight = float32(ib.Dy())

	for _, glyph := range config.Glyphs {
		if glyph.Width > f.maxGlyphWidth {
			f.maxGlyphWidth = glyph.Width
		}
		if glyph.Height > f.maxGlyphHeight {
			f.maxGlyphHeight = glyph.Height
		}
	}

	// save to disk for testing
	if gltext.IsDebug {
		err = gltext.SaveImage(".", "Debug", config.Image)
		if err != nil {
			return f, err
		}
	}

	// generate texture
	gl.GenTextures(1, &f.textureID)
	gl.BindTexture(gl.TEXTURE_2D, f.textureID)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexImage2D(
		gl.TEXTURE_2D,
		0,
		gl.RGBA,
		int32(ib.Dx()),
		int32(ib.Dy()),
		0,
		gl.RGBA,
		gl.UNSIGNED_BYTE,
		gl.Ptr(config.Image.Pix),
	)
	gl.BindTexture(gl.TEXTURE_2D, 0)

	// create shader program and define attributes and uniforms
	f.program, err = NewProgram(fontVertexShaderSource, fontFragmentShaderSource)
	if err != nil {
		return f, err
	}

	// attributes
	f.centeredPositionAttribute = uint32(gl.GetAttribLocation(f.program, gl.Str("centered_position\x00")))
	f.uvAttribute = uint32(gl.GetAttribLocation(f.program, gl.Str("uv\x00")))

	// uniforms
	f.finalPositionUniform = gl.GetUniformLocation(f.program, gl.Str("final_position\x00"))
	f.orthographicMatrixUniform = gl.GetUniformLocation(f.program, gl.Str("orthographic_matrix\x00"))
	f.scaleMatrixUniform = gl.GetUniformLocation(f.program, gl.Str("scale_matrix\x00"))
	f.fragmentTextureUniform = gl.GetUniformLocation(f.program, gl.Str("fragment_texture\x00"))
	f.colorUniform = gl.GetUniformLocation(f.program, gl.Str("fragment_color_adjustment\x00"))
	f.fadeoutUniform = gl.GetUniformLocation(f.program, gl.Str("fadeout\x00"))

	return f, nil
}

func (f *Font) ResizeWindow(width float32, height float32) {
	f.WindowWidth = width
	f.WindowHeight = height
	f.OrthographicMatrix = mgl32.Ortho2D(-f.WindowWidth/2, f.WindowWidth/2, -f.WindowHeight/2, f.WindowHeight/2)
}

func (f *Font) Release() {
	gl.DeleteTextures(1, &f.textureID)
}
