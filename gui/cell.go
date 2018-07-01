package gui

import (
	v41 "github.com/4ydx/gltext/v4.1"
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

type Cell struct {
	text       *v41.Text
	vao        uint32
	vbo        uint32
	cv         uint32
	colourAttr uint32
	points     []float32
	colour     [3]float32
	hidden     bool
}

func (gui *GUI) NewCell(font *v41.Font, x float32, y float32, w float32, h float32, colourAttr uint32, bgColour [3]float32) Cell {
	cell := Cell{
		text:       v41.NewText(font, 1.0, 1.1),
		colourAttr: colourAttr,
	}

	cell.colour = bgColour
	cell.text.SetPosition(mgl32.Vec2{x, y})

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

func (cell *Cell) SetFgColour(r, g, b float32) {
	if cell.text != nil {
		cell.text.SetColor(mgl32.Vec3{r, g, b})
	}
}

func (cell *Cell) SetBgColour(r float32, g float32, b float32) {

	if cell.colour[0] == r && cell.colour[1] == g && cell.colour[2] == b {
		return
	}

	cell.colour = [3]float32{r, g, b}
	cell.Clean()
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
		cell.colour[0], cell.colour[1], cell.colour[2],
		cell.colour[0], cell.colour[1], cell.colour[2],
		cell.colour[0], cell.colour[1], cell.colour[2],
		cell.colour[0], cell.colour[1], cell.colour[2],
		cell.colour[0], cell.colour[1], cell.colour[2],
		cell.colour[0], cell.colour[1], cell.colour[2],
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
	if cell.hidden {
		return
	}
	if cell.text != nil {
		cell.text.Draw()
	}

}

func (cell *Cell) Show() {
	cell.hidden = false
}

func (cell *Cell) Hide() {
	cell.hidden = true
}

func (cell *Cell) SetRune(r rune) {
	if cell.text != nil {
		if r == '%' {
			cell.text.SetString("%%")
		} else {
			cell.text.SetString(string(r))
		}

	}
}

func (cell *Cell) Release() {
	if cell.text != nil {
		cell.text.Release()
	}
}
