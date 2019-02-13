package glfont

import (
	"io"
	"io/ioutil"

	"github.com/go-gl/gl/all-core/gl"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
)

type character struct {
	textureID uint32 // ID handle of the glyph texture
	width     int    //glyph width
	height    int    //glyph height
	advance   int    //glyph advance
	bearingH  int    //glyph bearing horizontal
	bearingV  int    //glyph bearing vertical
}

//LoadTrueTypeFont builds a set of textures based on a ttf files glyphs
func LoadTrueTypeFont(program uint32, r io.Reader, scale float32) (*Font, error) {
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	//make Font struct type
	f := new(Font)
	f.scale = scale
	f.characters = map[rune]*character{}
	f.program = program //set shader program
	// Read the truetype font.
	f.ttf, err = truetype.Parse(data)
	if err != nil {
		return nil, err
	}
	f.SetColor(1.0, 1.0, 1.0, 1.0) //set default white

	_, h := f.MaxSize()
	f.lineHeight = h

	gl.BindTexture(gl.TEXTURE_2D, 0)

	// Configure VAO/VBO for texture quads
	gl.GenVertexArrays(1, &f.vao)
	gl.GenBuffers(1, &f.vbo)
	gl.BindVertexArray(f.vao)
	gl.BindBuffer(gl.ARRAY_BUFFER, f.vbo)

	gl.BufferData(gl.ARRAY_BUFFER, 6*4*4, nil, gl.DYNAMIC_DRAW)

	vertAttrib := uint32(gl.GetAttribLocation(f.program, gl.Str("vert\x00")))
	gl.EnableVertexAttribArray(vertAttrib)
	gl.VertexAttribPointer(vertAttrib, 2, gl.FLOAT, false, 4*4, gl.PtrOffset(0))
	defer gl.DisableVertexAttribArray(vertAttrib)

	texCoordAttrib := uint32(gl.GetAttribLocation(f.program, gl.Str("vertTexCoord\x00")))
	gl.EnableVertexAttribArray(texCoordAttrib)
	gl.VertexAttribPointer(texCoordAttrib, 2, gl.FLOAT, false, 4*4, gl.PtrOffset(2*4))
	defer gl.DisableVertexAttribArray(texCoordAttrib)

	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	gl.BindVertexArray(0)

	//create new face to measure glyph dimensions
	f.ttfFace = truetype.NewFace(f.ttf, &truetype.Options{
		Size:    float64(f.scale),
		DPI:     DPI,
		Hinting: font.HintingFull,
	})

	return f, nil
}
