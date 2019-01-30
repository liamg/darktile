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
	Inverse   bool
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
	if cell.Attr().Inverse {
		return cell.attr.BgColour
	}
	return cell.attr.FgColour
}

func (cell *Cell) Bg() [3]float32 {
	if cell.Attr().Inverse {
		return cell.attr.FgColour
	}
	return cell.attr.BgColour
}

func (cell *Cell) erase(bgColour [3]float32) {
	cell.setRune(0)
	cell.attr.BgColour = bgColour
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

func (cellAttr *CellAttributes) ReverseVideo() {
	oldFgColour := cellAttr.FgColour
	cellAttr.FgColour = cellAttr.BgColour
	cellAttr.BgColour = oldFgColour
}
