package gui

import (
	v41 "github.com/4ydx/gltext/v4.1"
	"github.com/go-gl/mathgl/mgl32"
)

type Cell struct {
	text *v41.Text
}

func NewCell(font *v41.Font, x float32, y float32, w float32, h float32) Cell {
	cell := Cell{
		text: v41.NewText(font, 1.0, 1.1),
	}

	cell.text.SetPosition(mgl32.Vec2{x, y})

	return cell

}

func (cell *Cell) Draw() {

	if cell.text != nil {
		cell.text.Draw()
	}
}

func (cell *Cell) Show() {
	if cell.text != nil {
		cell.text.Show()
	}
}

func (cell *Cell) Hide() {
	if cell.text != nil {
		cell.text.Hide()
	}
}

func (cell *Cell) SetRune(r rune) {
	if cell.text != nil {
		cell.text.SetString(string(r))
	}
}

func (cell *Cell) SetColour(r float32, g float32, b float32) {
	if cell.text != nil {
		cell.text.SetColor(mgl32.Vec3{r, g, b})
	}
}

func (cell *Cell) Release() {
	if cell.text != nil {
		cell.text.Release()
	}
}
