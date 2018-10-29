package buffer

import (
	"image"
)

type Cell struct {
	r     rune
	attr  CellAttributes
	image *image.RGBA
}

type CellAttributes struct {
	FgColour  [3]float32
	BgColour  [3]float32
	Bold      bool
	Dim       bool
	Underline bool
	Blink     bool
	Reverse   bool
	Hidden    bool
}

func (cell *Cell) Image() *image.RGBA {
	return cell.image
}

func (cell *Cell) SetImage(img *image.RGBA) {

	cell.image = img

}

func (cell *Cell) Attr() CellAttributes {
	return cell.attr
}

func (cell *Cell) Rune() rune {
	return cell.r
}

func (cell *Cell) Fg() [3]float32 {
	return cell.attr.FgColour
}

func (cell *Cell) Bg() [3]float32 {
	return cell.attr.BgColour
}

func (cell *Cell) erase() {
	cell.setRune(0)
}

func (cell *Cell) setRune(r rune) {
	cell.r = r
}

func NewBackgroundCell(colour [3]float32) Cell {
	return Cell{
		attr: CellAttributes{
			BgColour: colour,
		},
	}
}
