package gui

import (
	"github.com/go-gl/gl/all-core/gl"
	"github.com/liamg/glfont"
)

type Cell struct {
	vao        uint32
	vbo        uint32
	cv         uint32
	colourAttr uint32
	points     []float32
	bgColour   [3]float32
	fgColour   [3]float32
	hidden     bool
	r          rune
	font       *glfont.Font
	x          float32
	y          float32
}

func (gui *GUI) NewCell(font *glfont.Font, x float32, y float32, w float32, h float32, colourAttr uint32, bgColour [3]float32) Cell {
	cell := Cell{
		colourAttr: colourAttr,
		font:       font,
		x:          x + (float32(gui.width) / 2) - float32(gui.charWidth/2),
		y:          float32(gui.height) - (y + (float32(gui.height) / 2)) + (gui.charHeight / 2) - (gui.verticalPadding / 2),
	}

	cell.bgColour = bgColour

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

	cell.makeVao()

	return cell

}

func (cell *Cell) SetRune(r rune) {
	if cell.r == r {
		return
	}
	cell.r = r
}

func (cell *Cell) SetFgColour(r, g, b float32) {
	if cell.fgColour[0] != r || cell.fgColour[1] != g || cell.fgColour[2] != b {
		cell.fgColour = [3]float32{r, g, b}
	}
}

func (cell *Cell) SetBgColour(r float32, g float32, b float32) {

	if cell.bgColour[0] == r && cell.bgColour[1] == g && cell.bgColour[2] == b {
		return
	}

	cell.bgColour = [3]float32{r, g, b}
	//cell.Clean()
	cell.makeVao()
}

func (cell *Cell) Clean() {
	gl.DeleteVertexArrays(1, &cell.vao)
	gl.DeleteBuffers(1, &cell.vbo)
}

func (cell *Cell) makeVao() {

	gl.GenBuffers(1, &cell.vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, cell.vbo)
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(cell.points), gl.Ptr(cell.points), gl.STATIC_DRAW)

	gl.GenVertexArrays(1, &cell.vao)
	gl.BindVertexArray(cell.vao)
	gl.EnableVertexAttribArray(0)
	gl.BindBuffer(gl.ARRAY_BUFFER, cell.vbo)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 0, nil)

	// COLOUR
	gl.GenBuffers(1, &cell.cv)
	gl.BindBuffer(gl.ARRAY_BUFFER, cell.cv)
	triColor := []float32{
		cell.bgColour[0], cell.bgColour[1], cell.bgColour[2],
		cell.bgColour[0], cell.bgColour[1], cell.bgColour[2],
		cell.bgColour[0], cell.bgColour[1], cell.bgColour[2],
		cell.bgColour[0], cell.bgColour[1], cell.bgColour[2],
		cell.bgColour[0], cell.bgColour[1], cell.bgColour[2],
		cell.bgColour[0], cell.bgColour[1], cell.bgColour[2],
	}
	gl.BufferData(gl.ARRAY_BUFFER, len(triColor)*4, gl.Ptr(triColor), gl.STATIC_DRAW)
	gl.EnableVertexAttribArray(cell.colourAttr)
	gl.VertexAttribPointer(cell.colourAttr, 3, gl.FLOAT, false, 0, gl.PtrOffset(0))

	// END COLOUR

}

func (cell *Cell) DrawBg() {
	if cell.hidden {
		return
	}
	gl.BindVertexArray(cell.vao)
	gl.DrawArrays(gl.TRIANGLES, 0, 6)
}

func (cell *Cell) DrawText() {
	if cell.hidden || cell.r == ' ' {
		return
	}
	if cell.font != nil {
		cell.font.SetColor(cell.fgColour[0], cell.fgColour[1], cell.fgColour[2], 1.0)
		cell.font.Printf(cell.x, cell.y, 1, "%s", string(cell.r))
	}

}

func (cell *Cell) Show() {
	cell.hidden = false
}

func (cell *Cell) Hide() {
	cell.hidden = true
}
