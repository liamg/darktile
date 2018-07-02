package terminal

type Cell struct {
	r    rune
	attr CellAttributes
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

func (terminal *Terminal) NewCell() Cell {
	return Cell{
		attr: CellAttributes{
			FgColour: terminal.colourScheme.DefaultFg,
			BgColour: terminal.colourScheme.DefaultBg,
		},
	}
}

func (cell *Cell) GetRune() rune {
	return cell.r
}

func (cell *Cell) IsHidden() bool {
	return cell.attr.Hidden
}

func (cell *Cell) GetFgColour() (r float32, g float32, b float32) {

	if cell.attr.Reverse {
		return cell.attr.BgColour[0], cell.attr.BgColour[1], cell.attr.BgColour[2]
	}
	return cell.attr.FgColour[0], cell.attr.FgColour[1], cell.attr.FgColour[2]
}

func (cell *Cell) GetBgColour() (r float32, g float32, b float32) {

	if cell.attr.Reverse {
		return cell.attr.FgColour[0], cell.attr.FgColour[1], cell.attr.FgColour[2]
	}
	return cell.attr.BgColour[0], cell.attr.BgColour[1], cell.attr.BgColour[2]
}
