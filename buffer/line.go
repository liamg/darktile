package buffer

type Line struct {
	wrapped bool // whether line was wrapped onto from the previous one
	cells   []Cell
}

func newLine() Line {
	return Line{
		wrapped: false,
		cells:   []Cell{},
	}
}

func (line *Line) setWrapped(wrapped bool) {
	line.wrapped = wrapped
}

func (line *Line) String() string {
	runes := []rune{}
	for _, cell := range line.cells {
		runes = append(runes, cell.r)
	}
	return string(runes)
}

// @todo test these (ported from legacy) ------------------
func (line *Line) CutCellsAfter(n int) []Cell {
	cut := line.cells[n:]
	line.cells = line.cells[:n]
	return cut
}

func (line *Line) CutCellsFromBeginning(n int) []Cell {
	if n > len(line.cells) {
		n = len(line.cells)
	}
	cut := line.cells[:n]
	line.cells = line.cells[n:]
	return cut
}

func (line *Line) CutCellsFromEnd(n int) []Cell {
	cut := line.cells[len(line.cells)-n:]
	line.cells = line.cells[:len(line.cells)-n]
	return cut
}

// -------------------------------------------------------
