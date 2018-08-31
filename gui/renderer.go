package gui

import (
	"fmt"
	"math"

	"github.com/go-gl/gl/all-core/gl"
	"gitlab.com/liamg/raft/buffer"
	"gitlab.com/liamg/raft/config"
	"gitlab.com/liamg/raft/glfont"
)

type Renderer interface {
	SetArea(areaX int, areaY int, areaWidth int, areaHeight int)
	DrawCellBg(cell buffer.Cell, col uint, row uint)
	DrawCellText(cell buffer.Cell, col uint, row uint)
	DrawCursor(col uint, row uint, colour config.Colour)
	GetTermSize() (uint, uint)
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
	termCols            uint
	termRows            uint
	cellPositions       map[[2]uint][2]float32
	rectangles          map[[2]uint]*rectangle
	config              *config.Config
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
	prog       uint32
}

func (r *OpenGLRenderer) newRectangle(x float32, y float32, colourAttr uint32) *rectangle {

	x = (x - float32(r.areaWidth/2)) / float32(r.areaWidth/2)
	y = -(y - float32(r.areaHeight/2)) / float32(r.areaHeight/2)
	w := r.cellWidth / float32(r.areaWidth/2)
	h := r.cellHeight / float32(r.areaHeight/2)

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

func NewOpenGLRenderer(config *config.Config, font *glfont.Font, fontScale int32, areaX int, areaY int, areaWidth int, areaHeight int, colourAttr uint32, program uint32) *OpenGLRenderer {
	r := &OpenGLRenderer{
		areaWidth:     areaWidth,
		areaHeight:    areaHeight,
		areaX:         areaX,
		areaY:         areaY,
		fontScale:     fontScale,
		cellPositions: map[[2]uint][2]float32{},
		rectangles:    map[[2]uint]*rectangle{},
		config:        config,
		colourAttr:    colourAttr,
		program:       program,
	}
	r.SetFont(font)
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
	r.SetFont(r.font)
}

func (r *OpenGLRenderer) SetFontScale(fontScale int32) {
	r.fontScale = fontScale
	r.SetFont(r.font)
}

func (r *OpenGLRenderer) SetFont(font *glfont.Font) { // @todo check for monospace and return error if not?
	r.font = font
	r.verticalCellPadding = (0.25 * float32(r.fontScale))
	r.cellWidth = font.Width(1, "X")
	r.cellHeight = font.Height(1, "X") + (r.verticalCellPadding * 2) // vertical padding
	r.termCols = uint(math.Floor(float64(float32(r.areaWidth) / r.cellWidth)))
	r.termRows = uint(math.Floor(float64(float32(r.areaHeight) / r.cellHeight)))
	r.calculatePositions()
	r.rectangles = map[[2]uint]*rectangle{}
}

func (r *OpenGLRenderer) calculatePositions() {
	for line := uint(0); line < r.termRows; line++ {
		for col := uint(0); col < r.termCols; col++ {
			// rounding to whole pixels makes everything nice
			x := float32(math.Round(float64((float32(col) * r.cellWidth))))
			y := float32(math.Round(float64(
				(float32(line) * r.cellHeight) + (r.cellHeight / 2) + r.verticalCellPadding,
			)))
			r.cellPositions[[2]uint{col, line}] = [2]float32{x, y}
		}
	}
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
	y := float32((float32(line) * r.cellHeight) + (r.cellHeight))
	r.rectangles[[2]uint{col, line}] = r.newRectangle(x, y, r.colourAttr)
	return r.rectangles[[2]uint{col, line}]
}

func (r *OpenGLRenderer) DrawCursor(col uint, row uint, colour config.Colour) {
	rect := r.getRectangle(col, row)
	rect.setColour(colour)
	rect.Draw()
}

func (r *OpenGLRenderer) DrawCellBg(cell buffer.Cell, col uint, row uint) {

	var bg [3]float32

	if cell.Attr().Reverse {
		bg = cell.Fg()
	} else {
		bg = cell.Bg()
	}

	if bg != r.config.ColourScheme.Background {
		rect := r.getRectangle(col, row)
		rect.setColour(bg)
		rect.Draw()
	}

}

func (r *OpenGLRenderer) DrawCellText(cell buffer.Cell, col uint, row uint) {

	var fg [3]float32

	if cell.Attr().Reverse {
		fg = cell.Bg()
	} else {
		fg = cell.Fg()
	}

	pos, ok := r.cellPositions[[2]uint{col, row}]
	if !ok {
		panic(fmt.Sprintf("Missing position data for cell at %d,%d", col, row))
	}

	var alpha float32 = 1
	if cell.Attr().Dim {
		alpha = 0.5
	}
	r.font.SetColor(fg[0], fg[1], fg[2], alpha)

	if cell.Attr().Bold { // bold means draw text again one pixel to right, so it's fatter
		r.font.Print(pos[0]+1, pos[1], 1, string(cell.Rune()))
	}
	r.font.Print(pos[0], pos[1], 1, string(cell.Rune()))

}
