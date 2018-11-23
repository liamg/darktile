package glfont

import (
	"fmt"
	"image"
	"image/draw"
	"io"
	"math"

	"github.com/go-gl/gl/all-core/gl"
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

const DPI = 72

// A Font allows rendering of text to an OpenGL context.
type Font struct {
	characters  map[rune]*character
	vao         uint32
	vbo         uint32
	program     uint32
	texture     uint32 // Holds the glyph texture id.
	color       color
	ttf         *truetype.Font
	scale       float32
	linePadding float32
	lineHeight  float32
}

type color struct {
	r float32
	g float32
	b float32
	a float32
}

//LoadFont loads the specified font at the given scale.
func LoadFont(reader io.Reader, scale float32, windowWidth int, windowHeight int) (*Font, error) {

	// Configure the default font vertex and fragment shaders
	program, err := newProgram(vertexFontShader, fragmentFontShader)
	if err != nil {
		panic(err)
	}

	// Activate corresponding render state
	gl.UseProgram(program)

	//set screen resolution
	resUniform := gl.GetUniformLocation(program, gl.Str("resolution\x00"))
	gl.Uniform2f(resUniform, float32(windowWidth), float32(windowHeight))

	return LoadTrueTypeFont(program, reader, scale)
}

//SetColor allows you to set the text color to be used when you draw the text
func (f *Font) SetColor(red float32, green float32, blue float32, alpha float32) {
	f.color.r = red
	f.color.g = green
	f.color.b = blue
	f.color.a = alpha
}

func (f *Font) UpdateResolution(windowWidth int, windowHeight int) {
	gl.UseProgram(f.program)
	resUniform := gl.GetUniformLocation(f.program, gl.Str("resolution\x00"))
	gl.Uniform2f(resUniform, float32(windowWidth), float32(windowHeight))
	gl.UseProgram(0)
	//f.characters = map[rune]*character{}
}

func (f *Font) LineHeight() float32 {
	return f.lineHeight
}

func (f *Font) LinePadding() float32 {
	return f.linePadding
}

//Printf draws a string to the screen, takes a list of arguments like printf
func (f *Font) Print(x, y float32, text string) error {

	x = float32(math.Round(float64(x)))
	y = float32(math.Round(float64(y)))

	indices := []rune(text)

	if len(indices) == 0 {
		return nil
	}

	//setup blending mode
	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)

	// Activate corresponding render state
	gl.UseProgram(f.program)
	//set text color
	gl.Uniform4f(gl.GetUniformLocation(f.program, gl.Str("textColor\x00")), f.color.r, f.color.g, f.color.b, f.color.a)
	//set screen resolution
	//resUniform := gl.GetUniformLocation(f.program, gl.Str("resolution\x00"))
	//gl.Uniform2f(resUniform, float32(2560), float32(1440))

	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindVertexArray(f.vao)

	// Iterate through all characters in string
	for i := range indices {

		//get rune
		runeIndex := indices[i]

		//find rune in fontChar list
		ch, err := f.GetRune(runeIndex)
		if err != nil {
			return err // @todo ignore errors?
		}

		//calculate position and size for current rune
		xpos := x + float32(ch.bearingH)
		ypos := y - float32(+ch.height-ch.bearingV)
		w := float32(ch.width)
		h := float32(ch.height)

		//set quad positions
		var x1 = xpos
		var x2 = xpos + w
		var y1 = ypos
		var y2 = ypos + h

		//setup quad array
		var vertices = []float32{
			//  X, Y, Z, U, V
			// Front
			x1, y1, 0.0, 0.0,
			x2, y1, 1.0, 0.0,
			x1, y2, 0.0, 1.0,
			x1, y2, 0.0, 1.0,
			x2, y1, 1.0, 0.0,
			x2, y2, 1.0, 1.0}

		// Render glyph texture over quad
		gl.BindTexture(gl.TEXTURE_2D, ch.textureID)
		// Update content of VBO memory
		gl.BindBuffer(gl.ARRAY_BUFFER, f.vbo)

		//BufferSubData(target Enum, offset int, data []byte)
		gl.BufferSubData(gl.ARRAY_BUFFER, 0, len(vertices)*4, gl.Ptr(vertices)) // Be sure to use glBufferSubData and not glBufferData
		// Render quad
		gl.DrawArrays(gl.TRIANGLES, 0, 24)

		gl.BindBuffer(gl.ARRAY_BUFFER, 0)
		// Now advance cursors for next glyph (note that advance is number of 1/64 pixels)
		x += float32((ch.advance >> 6)) // Bitshift by 6 to get value in pixels (2^6 = 64 (divide amount of 1/64th pixels by 64 to get amount of pixels))

	}

	//clear opengl textures and programs
	gl.BindVertexArray(0)
	gl.BindTexture(gl.TEXTURE_2D, 0)
	gl.UseProgram(0)
	gl.Disable(gl.BLEND)

	return nil
}

//Width returns the width of a piece of text in pixels
func (f *Font) Size(text string) (float32, float32) {

	var width float32
	var height float32

	indices := []rune(text)

	if len(indices) == 0 {
		return 0, 0
	}

	// Iterate through all characters in string
	for i := range indices {

		//get rune
		runeIndex := indices[i]

		//find rune in fontChar list
		ch, err := f.GetRune(runeIndex)
		if err != nil {
			continue
		}

		// Now advance cursors for next glyph (note that advance is number of 1/64 pixels)
		width += float32((ch.advance >> 6)) // Bitshift by 6 to get value in pixels (2^6 = 64 (divide amount of 1/64th pixels by 64 to get amount of pixels))

		// Now advance cursors for next glyph (note that advance is number of 1/64 pixels)
		if float32(ch.height)*f.scale > height {
			height = float32(ch.height)
		}
	}

	return width, height
}

func(f *Font) MaxSize() (float32, float32){
	b:= f.ttf.Bounds(fixed.Int26_6(f.scale))
	return float32(b.Max.X - b.Min.X),float32(b.Max.Y - b.Min.Y)
}

func(f *Font) MinY() float32 {
	b:= f.ttf.Bounds(fixed.Int26_6(f.scale))
	return float32(b.Min.Y)
}

func(f *Font) MaxY() float32 {
	b:= f.ttf.Bounds(fixed.Int26_6(f.scale))
	return float32(b.Max.Y)
}

func (f *Font) GetRune(r rune) (*character, error) {

	cc, ok := f.characters[r]
	if ok {
		return cc, nil
	}

	char := new(character)

	//create new face to measure glyph diamensions
	ttfFace := truetype.NewFace(f.ttf, &truetype.Options{
		Size:    float64(f.scale),
		DPI:     DPI,
		Hinting: font.HintingFull,
	})

	gBnd, gAdv, ok := ttfFace.GlyphBounds(r)
	if ok != true {
		return nil, fmt.Errorf("ttf face glyphBounds error")
	}

	gh := int32((gBnd.Max.Y - gBnd.Min.Y) >> 6)
	gw := int32((gBnd.Max.X - gBnd.Min.X) >> 6)

	//if gylph has no diamensions set to a max value
	if gw == 0 || gh == 0 {
		gBnd = f.ttf.Bounds(fixed.Int26_6(f.scale))
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
	c.SetDPI(DPI)
	c.SetFont(f.ttf)
	c.SetFontSize(float64(f.scale))
	c.SetClip(rgba.Bounds())
	c.SetDst(rgba)
	c.SetSrc(fg)
	c.SetHinting(font.HintingFull)

	//set the glyph dot
	px := 0 - (int(gBnd.Min.X) >> 6)
	py := (gAscent)
	pt := freetype.Pt(px, py)

	// Draw the text from mask to image
	if _, err := c.DrawString(string(r), pt); err != nil {
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

	f.characters[r] = char

	return char, nil
}
