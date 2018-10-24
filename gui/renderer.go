package gui

import (
	"math"

	"github.com/go-gl/gl/all-core/gl"
	"github.com/liamg/aminal/buffer"
	"github.com/liamg/aminal/config"
	"github.com/liamg/aminal/glfont"
)

type OpenGLRenderer struct {
	font          *glfont.Font
	boldFont      *glfont.Font
	areaWidth     int
	areaHeight    int
	areaX         int
	areaY         int
	cellWidth     float32
	cellHeight    float32
	termCols      uint
	termRows      uint
	cellPositions map[[2]uint][2]float32
	rectangles    map[[2]uint]*rectangle
	config        *config.Config
	colourAttr    uint32
	program       uint32
}

type rectangle struct {
	vao        uint32
	vbo        uint32
	cv         uint32
	colourAttr uint32
	colour     [3]float32
	points     []float32
	prog       uint32
}

func (r *OpenGLRenderer) CellWidth() float32 {
	return r.cellWidth
}

func (r *OpenGLRenderer) CellHeight() float32 {
	return r.cellHeight
}

func (r *OpenGLRenderer) newRectangle(x float32, y float32, colourAttr uint32) *rectangle {

	halfAreaWidth := float32(r.areaWidth / 2)
	halfAreaHeight := float32(r.areaHeight / 2)

	x = (x - halfAreaWidth) / halfAreaWidth
	y = -(y - halfAreaHeight) / halfAreaHeight
	w := r.cellWidth / halfAreaWidth
	h := r.cellHeight / halfAreaHeight

	rect := &rectangle{
		points: []float32{
			x, y, 0,
			x, y + h, 0,
			x + w, y + h, 0,

			x + w, y, 0,
			x, y, 0,
			x + w, y + h, 0,
		},
		colourAttr: colourAttr,
		prog:       r.program,
	}

	gl.UseProgram(rect.prog)

	// SHAPE
	gl.GenBuffers(1, &rect.vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, rect.vbo)
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(rect.points), gl.Ptr(rect.points), gl.STATIC_DRAW)

	gl.GenVertexArrays(1, &rect.vao)
	gl.BindVertexArray(rect.vao)
	gl.EnableVertexAttribArray(0)

	gl.BindBuffer(gl.ARRAY_BUFFER, rect.vbo)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 0, nil)

	// colour
	gl.GenBuffers(1, &rect.cv)

	rect.setColour([3]float32{0, 1, 0})

	return rect
}

func (rect *rectangle) Draw() {
	gl.UseProgram(rect.prog)
	gl.BindVertexArray(rect.vao)
	gl.DrawArrays(gl.TRIANGLES, 0, 6)
}

func (rect *rectangle) setColour(colour [3]float32) {
	if rect.colour == colour {
		return
	}

	c := []float32{
		colour[0], colour[1], colour[2],
		colour[0], colour[1], colour[2],
		colour[0], colour[1], colour[2],
		colour[0], colour[1], colour[2],
		colour[0], colour[1], colour[2],
		colour[0], colour[1], colour[2],
	}

	gl.UseProgram(rect.prog)
	gl.BindBuffer(gl.ARRAY_BUFFER, rect.cv)
	gl.BufferData(gl.ARRAY_BUFFER, len(c)*4, gl.Ptr(c), gl.STATIC_DRAW)
	gl.EnableVertexAttribArray(rect.colourAttr)
	gl.VertexAttribPointer(rect.colourAttr, 3, gl.FLOAT, false, 0, gl.PtrOffset(0))

	rect.colour = colour
}

func (rect *rectangle) Free() {
	gl.UseProgram(rect.prog)
	gl.DeleteVertexArrays(1, &rect.vao)
	gl.DeleteBuffers(1, &rect.vbo)
	gl.DeleteBuffers(1, &rect.cv)
}

func NewOpenGLRenderer(config *config.Config, font *glfont.Font, boldFont *glfont.Font, areaX int, areaY int, areaWidth int, areaHeight int, colourAttr uint32, program uint32) *OpenGLRenderer {
	r := &OpenGLRenderer{
		areaWidth:     areaWidth,
		areaHeight:    areaHeight,
		areaX:         areaX,
		areaY:         areaY,
		cellPositions: map[[2]uint][2]float32{},
		rectangles:    map[[2]uint]*rectangle{},
		config:        config,
		colourAttr:    colourAttr,
		program:       program,
	}
	r.SetFont(font, boldFont)
	return r
}

func (r *OpenGLRenderer) GetTermSize() (uint, uint) {
	return r.termCols, r.termRows
}

func (r *OpenGLRenderer) SetArea(areaX int, areaY int, areaWidth int, areaHeight int) {
	r.areaWidth = areaWidth
	r.areaHeight = areaHeight
	r.areaX = areaX
	r.areaY = areaY
	r.SetFont(r.font, r.boldFont)
}

func (r *OpenGLRenderer) SetFont(font *glfont.Font, bold *glfont.Font) { // @todo check for monospace and return error if not?
	r.font = font
	r.boldFont = bold
	r.cellWidth, _ = font.Size("X")
	r.cellHeight = font.LineHeight() // vertical padding
	r.termCols = uint(math.Floor(float64(float32(r.areaWidth) / r.cellWidth)))
	r.termRows = uint(math.Floor(float64(float32(r.areaHeight) / r.cellHeight)))
	r.rectangles = map[[2]uint]*rectangle{}
}

func (r *OpenGLRenderer) getRectangle(col uint, row uint) *rectangle {
	if rect, ok := r.rectangles[[2]uint{col, row}]; ok {
		return rect
	}
	return r.generateRectangle(col, row)
}

func (r *OpenGLRenderer) generateRectangle(col uint, line uint) *rectangle {

	rect, ok := r.rectangles[[2]uint{col, line}]
	if ok {
		rect.Free()
	}

	// rounding to whole pixels makes everything nice
	x := float32(float32(col) * r.cellWidth)
	y := float32((float32(line) * r.cellHeight)) + r.cellHeight
	r.rectangles[[2]uint{col, line}] = r.newRectangle(x, y, r.colourAttr)
	return r.rectangles[[2]uint{col, line}]
}

func (r *OpenGLRenderer) DrawCursor(col uint, row uint, colour config.Colour) {
	rect := r.getRectangle(col, row)
	rect.setColour(colour)
	rect.Draw()
}

func (r *OpenGLRenderer) DrawCellBg(cell buffer.Cell, col uint, row uint, cursor bool, colour *config.Colour) {

	var bg [3]float32

	if colour != nil {
		bg = *colour
	} else {

		if cursor {
			bg = r.config.ColourScheme.Cursor
		} else if cell.Attr().Reverse {
			bg = cell.Fg()
		} else {
			bg = cell.Bg()
		}
	}

	if bg != r.config.ColourScheme.Background {
		rect := r.getRectangle(col, row)
		rect.setColour(bg)
		rect.Draw()
	}
}

func (r *OpenGLRenderer) DrawCellText(cell buffer.Cell, col uint, row uint, colour *config.Colour) {

	var fg [3]float32

	if colour != nil {
		fg = *colour
	} else {
		if cell.Attr().Reverse {
			fg = cell.Bg()
		} else {
			fg = cell.Fg()
		}
	}

	var alpha float32 = 1
	if cell.Attr().Dim {
		alpha = 0.5
	}
	r.font.SetColor(fg[0], fg[1], fg[2], alpha)

	x := float32(r.areaX) + float32(col)*r.cellWidth
	y := float32(r.areaY) + (float32(row+1) * r.cellHeight) - (r.font.LinePadding())

	if cell.Attr().Bold { // bold means draw text again one pixel to right, so it's fatter
		if r.boldFont != nil {
			y := float32(r.areaY) + (float32(row+1) * r.cellHeight) - (r.boldFont.LinePadding())
			r.boldFont.SetColor(fg[0], fg[1], fg[2], alpha)
			r.boldFont.Print(x, y, string(cell.Rune()))
			return
		}
		r.font.Print(x+1, y, string(cell.Rune()))
	}

	r.font.Print(x, y, string(cell.Rune()))

}
