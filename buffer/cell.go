package buffer

type Cell struct {
	r          rune
	attr       CellAttributes
	hasContent bool
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

func newCell() Cell {
	return Cell{}
}

func (cell *Cell) Rune() rune {
	if cell.r == 0 {
		return '~'
	}
	return cell.r
}

func (cell *Cell) Fg() [3]float32 {
	if cell.r == 0 {
		return [3]float32{0.5, 0.5, 0.5}
	}
	return [3]float32{0.9, 0.9, 0.9} // @todo fix this
	//return cell.attr.FgColour
}

func (cell *Cell) Bg() [3]float32 {
	return cell.attr.BgColour
}

func (cell *Cell) erase() {
	cell.setRune(0)
}

func (cell *Cell) setRune(r rune) {
	cell.r = r
	cell.hasContent = r > 0
}
