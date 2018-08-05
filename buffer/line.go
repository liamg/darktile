package buffer

type line struct {
	wrapped bool // whether line was wrapped onto from the previous one
	cells   []Cell
}

func newLine() line {
	return line{
		wrapped: false,
		cells:   []Cell{},
	}
}
