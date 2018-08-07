package gui

import (
	"fmt"
	"math"

	"github.com/liamg/glfont"
	"gitlab.com/liamg/raft/buffer"
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
}

func NewOpenGLRenderer(font *glfont.Font, fontScale int32, areaX int, areaY int, areaWidth int, areaHeight int) *OpenGLRenderer {
	r := &OpenGLRenderer{
		areaWidth:     areaWidth,
		areaHeight:    areaHeight,
		areaX:         areaX,
		areaY:         areaY,
		fontScale:     fontScale,
		cellPositions: map[[2]int][2]float32{},
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

func (r *OpenGLRenderer) SetFont(font *glfont.Font) {
	r.font = font
	r.verticalCellPadding = (0.3 * float32(r.fontScale))
	r.cellWidth = font.Width(1, "X")
	r.cellHeight = font.Height(1, "X") + (r.verticalCellPadding * 2) // vertical padding
	r.termCols = int(math.Floor(float64(float32(r.areaWidth) / r.cellWidth)))
	r.termRows = int(math.Floor(float64(float32(r.areaHeight) / r.cellHeight)))
	r.calculatePositions()
}

func (r *OpenGLRenderer) calculatePositions() {
	for line := 0; line < r.termRows; line++ {
		for col := 0; col < r.termCols; col++ {
			// rounding to whole pixels makes everything nice
			x := float32(math.Floor(float64((float32(col) * r.cellWidth) + (r.cellWidth / 2))))
			y := float32(math.Floor(float64(
				(float32(line) * r.cellHeight) + (r.cellHeight / 2) + r.verticalCellPadding,
			)))
			r.cellPositions[[2]int{col, line}] = [2]float32{x, y}
		}
	}
}

func (r *OpenGLRenderer) DrawCell(cell *buffer.Cell, col int, row int) {

	if cell == nil {
		return
	}

	fg := cell.Fg()
	r.font.SetColor(fg[0], fg[1], fg[2], 1)

	pos, ok := r.cellPositions[[2]int{col, row}]
	if !ok {
		panic(fmt.Sprintf("Missing position data for cell at %d,%d", col, row))
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
