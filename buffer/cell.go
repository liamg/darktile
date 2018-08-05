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

func (cell *Cell) setRune(r rune) {
	cell.r = r
	cell.hasContent = true
}
