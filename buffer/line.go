package buffer

import (
	"strings"
)

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

func (line *Line) Cells() []Cell {
	return line.cells
}

func (line *Line) ReverseVideo() {
	for i, _ := range line.cells {
		line.cells[i].attr.ReverseVideo()
	}
}

// Cleanse removes null bytes from the end of the row
func (line *Line) Cleanse() {
	cut := 0
	for i := len(line.cells) - 1; i >= 0; i-- {
		if line.cells[i].r != 0 {
			break
		}
		cut++
	}
	if cut == 0 {
		return
	}
	line.cells = line.cells[:len(line.cells)-cut]
}

func (line *Line) setWrapped(wrapped bool) {
	line.wrapped = wrapped
}

func (line *Line) String() string {
	runes := []rune{}
	for _, cell := range line.cells {
		runes = append(runes, cell.r)
	}
	return strings.TrimRight(string(runes), "\x00 ")
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

func (line *Line) Append(cells ...Cell) {
	line.cells = append(line.cells, cells...)
}

// -------------------------------------------------------
