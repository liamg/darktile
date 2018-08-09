package gui

import (
	"fmt"
	"math"

	"github.com/go-gl/gl/all-core/gl"
	"github.com/liamg/glfont"
	"gitlab.com/liamg/raft/buffer"
	"gitlab.com/liamg/raft/config"
)

type Renderer interface {
	SetArea(areaX int, areaY int, areaWidth int, areaHeight int)
	DrawCell(cell *buffer.Cell, col int, row int)
	GetTermSize() (int, int)
}

type OpenGLRenderer struct {
	font                *glfont.Font
	areaWidth           int
	areaHeight          int
	areaX               int
	areaY               int
	fontScale           int32
	cellWidth           float32
	cellHeight          float32
	verticalCellPadding float32
	termCols            int
	termRows            int
	cellPositions       map[[2]int][2]float32
	rectangles          map[[2]int]*rectangle
	config              config.Config
	colourAttr          uint32
	program             uint32
}

type rectangle struct {
	vao        uint32
	vbo        uint32
	cv         uint32
	colourAttr uint32
	colour     [3]float32
	points     []float32
}

func (r *OpenGLRenderer) newRectangle(x float32, y float32, colourAttr uint32) *rectangle {

	x = (x - float32(r.areaWidth/2)) / float32(r.areaWidth/2)
	y = -(y - float32(r.areaHeight/2)) / float32(r.areaHeight/2)
	w := r.cellWidth / float32(r.areaWidth/2)
	h := r.cellHeight / float32(r.areaHeight/2)

	rect := &rectangle{
		colour: [3]float32{0, 0, 1},
		points: []float32{
			x, y + h, 0,
			x, y, 0,
			x + w, y, 0,
			x, y + h, 0,
			x + w, y + h, 0,
			x + w, y, 0,
		},
		colourAttr: colourAttr,
	}

	rect.gen()

	return rect
}

func (rect *rectangle) gen() {

	colour := []float32{
		rect.colour[0], rect.colour[1], rect.colour[2],
		rect.colour[0], rect.colour[1], rect.colour[2],
		rect.colour[0], rect.colour[1], rect.colour[2],
		rect.colour[0], rect.colour[1], rect.colour[2],
		rect.colour[0], rect.colour[1], rect.colour[2],
		rect.colour[0], rect.colour[1], rect.colour[2],
	}

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
	gl.BindBuffer(gl.ARRAY_BUFFER, rect.cv)
	gl.BufferData(gl.ARRAY_BUFFER, len(colour)*4, gl.Ptr(colour), gl.STATIC_DRAW)
	gl.EnableVertexAttribArray(rect.colourAttr)
	gl.VertexAttribPointer(rect.colourAttr, 3, gl.FLOAT, false, 0, gl.PtrOffset(0))
}

func (rect *rectangle) setColour(colour [3]float32) {
	if rect.colour == colour {
		return
	}
	rect.Free()
	rect.colour = colour
	rect.gen()
}

func (rect *rectangle) Free() {
	gl.DeleteVertexArrays(1, &rect.vao)
	gl.DeleteBuffers(1, &rect.vbo)
	gl.DeleteBuffers(1, &rect.cv)
}

func NewOpenGLRenderer(config config.Config, font *glfont.Font, fontScale int32, areaX int, areaY int, areaWidth int, areaHeight int, colourAttr uint32, program uint32) *OpenGLRenderer {
	r := &OpenGLRenderer{
		areaWidth:     areaWidth,
		areaHeight:    areaHeight,
		areaX:         areaX,
		areaY:         areaY,
		fontScale:     fontScale,
		cellPositions: map[[2]int][2]float32{},
		rectangles:    map[[2]int]*rectangle{},
		config:        config,
		colourAttr:    colourAttr,
		program:       program,
	}
	r.SetFont(font)
	return r
}

func (r *OpenGLRenderer) GetTermSize() (int, int) {
	return r.termCols, r.termRows
}

func (r *OpenGLRenderer) SetArea(areaX int, areaY int, areaWidth int, areaHeight int) {
	r.areaWidth = areaWidth
	r.areaHeight = areaHeight
	r.areaX = areaX
	r.areaY = areaY
	r.SetFont(r.font)
}

func (r *OpenGLRenderer) SetFontScale(fontScale int32) {
	r.fontScale = fontScale
	r.SetFont(r.font)
}

func (r *OpenGLRenderer) SetFont(font *glfont.Font) { // @todo check for monospace and return error if not?
	r.font = font
	r.verticalCellPadding = (0.3 * float32(r.fontScale))
	r.cellWidth = font.Width(1, "X")
	r.cellHeight = font.Height(1, "X") + (r.verticalCellPadding * 2) // vertical padding
	r.termCols = int(math.Floor(float64(float32(r.areaWidth) / r.cellWidth)))
	r.termRows = int(math.Floor(float64(float32(r.areaHeight) / r.cellHeight)))
	r.calculatePositions()
	r.generateRectangles()
}

func (r *OpenGLRenderer) calculatePositions() {
	for line := 0; line < r.termRows; line++ {
		for col := 0; col < r.termCols; col++ {
			// rounding to whole pixels makes everything nice
			x := float32(math.Round(float64((float32(col) * r.cellWidth))))
			y := float32(math.Round(float64(
				(float32(line) * r.cellHeight) + (r.cellHeight / 2) + r.verticalCellPadding,
			)))
			r.cellPositions[[2]int{col, line}] = [2]float32{x, y}
		}
	}
}

func (r *OpenGLRenderer) generateRectangles() {
	gl.UseProgram(r.program)
	for line := 0; line < r.termRows; line++ {
		for col := 0; col < r.termCols; col++ {

			rect, ok := r.rectangles[[2]int{col, line}]
			if ok {
				rect.Free()
			}

			// rounding to whole pixels makes everything nice
			x := float32(float64((float32(col) * r.cellWidth)))
			y := float32(float64((float32(line) * r.cellHeight) + (r.cellHeight)))
			r.rectangles[[2]int{col, line}] = r.newRectangle(x, y, r.colourAttr)
		}
	}
}

func (r *OpenGLRenderer) DrawCell(cell *buffer.Cell, col int, row int) {

	if cell == nil || cell.Attr().Hidden || cell.Rune() == 0x00 {
		return
	}

	var fg [3]float32
	var bg [3]float32

	if cell.Attr().Reverse {
		fg = cell.Bg()
		bg = cell.Fg()
	} else {
		fg = cell.Fg()
		bg = cell.Bg()
	}

	var alpha float32 = 1
	if cell.Attr().Dim {
		alpha = 0.5
	}
	r.font.SetColor(fg[0], fg[1], fg[2], alpha)

	pos, ok := r.cellPositions[[2]int{col, row}]
	if !ok {
		panic(fmt.Sprintf("Missing position data for cell at %d,%d", col, row))
	}

	rect, ok := r.rectangles[[2]int{col, row}]
	if !ok {
		panic(fmt.Sprintf("Missing rectangle data for cell at %d,%d", col, row))
	}

	gl.UseProgram(r.program)

	// don't bother rendering rectangles that are the same colour as the background
	if bg != r.config.ColourScheme.DefaultBg {
		rect.setColour(bg)
		gl.BindVertexArray(rect.vao)
		gl.DrawArrays(gl.TRIANGLES, 0, 6)
	}

	if cell.Attr().Bold { // bold means draw text again one pixel to right, so it's fatter
		r.font.Print(pos[0]+1, pos[1], 1, string(cell.Rune()))
	}
	r.font.Print(pos[0], pos[1], 1, string(cell.Rune()))

	/*
		this was passed into cell
		x := ((float32(col) * gui.charWidth) - (float32(gui.width) / 2)) + (gui.charWidth / 2)
		y := -(((float32(row) * gui.charHeight) - (float32(gui.height) / 2)) + (gui.charHeight / 2))

		this was in cell:
		x:          x + (float32(gui.width) / 2) - float32(gui.charWidth/2),
		y:          float32(gui.height) - (y + (float32(gui.height) / 2)) + (gui.charHeight / 2) - (gui.verticalPadding / 2),


		and then points:

		x = (x - (w / 2)) / (float32(gui.width) / 2)
		y = (y - (h / 2)) / (float32(gui.height) / 2)
		w = (w / float32(gui.width/2))
		h = (h / float32(gui.height/2))
		cell.points = []float32{
			x, y + h, 0,
			x, y, 0,
			x + w, y, 0,
			x, y + h, 0,
			x + w, y + h, 0,
			x + w, y, 0,
		}

	*/

}
