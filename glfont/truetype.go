package glfont

import (
	"fmt"
	"github.com/go-gl/gl/all-core/gl"
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
	"image"
	"image/draw"
	"io"
	"io/ioutil"
)

type character struct {
	textureID uint32 // ID handle of the glyph texture
	width     int    //glyph width
	height    int    //glyph height
	advance   int    //glyph advance
	bearingH  int    //glyph bearing horizontal
	bearingV  int    //glyph bearing vertical
}

//LoadTrueTypeFont builds a set of textures based on a ttf files gylphs
func LoadTrueTypeFont(program uint32, r io.Reader, scale int32, low, high rune, dir Direction) (*Font, error) {
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	// Read the truetype font.
	ttf, err := truetype.Parse(data)
	if err != nil {
		return nil, err
	}

	//make Font stuct type
	f := new(Font)
	f.fontChar = make([]*character, 0, high-low+1)
	f.program = program            //set shader program
	f.SetColor(1.0, 1.0, 1.0, 1.0) //set default white

	//make each gylph
	for ch := low; ch <= high; ch++ {

		char := new(character)

		//create new face to measure glyph diamensions
		ttfFace := truetype.NewFace(ttf, &truetype.Options{
			Size:    float64(scale),
			DPI:     72,
			Hinting: font.HintingFull,
		})

		gBnd, gAdv, ok := ttfFace.GlyphBounds(ch)
		if ok != true {
			return nil, fmt.Errorf("ttf face glyphBounds error")
		}

		gh := int32((gBnd.Max.Y - gBnd.Min.Y) >> 6)
		gw := int32((gBnd.Max.X - gBnd.Min.X) >> 6)

		//if gylph has no diamensions set to a max value
		if gw == 0 || gh == 0 {
			gBnd = ttf.Bounds(fixed.Int26_6(scale))
			gw = int32((gBnd.Max.X - gBnd.Min.X) >> 6)
			gh = int32((gBnd.Max.Y - gBnd.Min.Y) >> 6)

			//above can sometimes yield 0 for font smaller than 48pt, 1 is minimum
			if gw == 0 || gh == 0 {
				gw = 1
				gh = 1
			}
		}

		//The glyph's ascent and descent equal -bounds.Min.Y and +bounds.Max.Y.
		gAscent := int(-gBnd.Min.Y) >> 6
		gdescent := int(gBnd.Max.Y) >> 6

		//set w,h and adv, bearing V and bearing H in char
		char.width = int(gw)
		char.height = int(gh)
		char.advance = int(gAdv)
		char.bearingV = gdescent
		char.bearingH = (int(gBnd.Min.X) >> 6)

		//create image to draw glyph
		fg, bg := image.White, image.Black
		rect := image.Rect(0, 0, int(gw), int(gh))
		rgba := image.NewRGBA(rect)
		draw.Draw(rgba, rgba.Bounds(), bg, image.ZP, draw.Src)

		//create a freetype context for drawing
		c := freetype.NewContext()
		c.SetDPI(72)
		c.SetFont(ttf)
		c.SetFontSize(float64(scale))
		c.SetClip(rgba.Bounds())
		c.SetDst(rgba)
		c.SetSrc(fg)
		c.SetHinting(font.HintingFull)

		//set the glyph dot
		px := 0 - (int(gBnd.Min.X) >> 6)
		py := (gAscent)
		pt := freetype.Pt(px, py)

		// Draw the text from mask to image
		_, err = c.DrawString(string(ch), pt)
		if err != nil {
			return nil, err
		}

		// Generate texture
		var texture uint32
		gl.GenTextures(1, &texture)
		gl.BindTexture(gl.TEXTURE_2D, texture)
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
		gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, int32(rgba.Rect.Dx()), int32(rgba.Rect.Dy()), 0,
			gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(rgba.Pix))

		char.textureID = texture

		//add char to fontChar list
		f.fontChar = append(f.fontChar, char)

	}

	gl.BindTexture(gl.TEXTURE_2D, 0)

	// Configure VAO/VBO for texture quads
	gl.GenVertexArrays(1, &f.vao)
	gl.GenBuffers(1, &f.vbo)
	gl.BindVertexArray(f.vao)
	gl.BindBuffer(gl.ARRAY_BUFFER, f.vbo)

	gl.BufferData(gl.ARRAY_BUFFER, 6*4*4, nil, gl.STATIC_DRAW)

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

	return f, nil
}
