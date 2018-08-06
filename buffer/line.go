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
