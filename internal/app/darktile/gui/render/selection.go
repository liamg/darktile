package render

import (
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text"
)

func (r *Render) drawSelection() {
	_, selection := r.buffer.GetSelection()
	if selection == nil {
		// nothing selected
		return
	}

	bg, fg := r.theme.SelectionBackground(), r.theme.SelectionForeground()

	for y := selection.Start.Line; y <= selection.End.Line; y++ {
		xStart, xEnd := 0, int(r.buffer.ViewWidth())
		if y == selection.Start.Line {
			xStart = int(selection.Start.Col)
		}
		if y == selection.End.Line {
			xEnd = int(selection.End.Col)
		}
		for x := xStart; x <= xEnd; x++ {
			pX, pY := float64(x*r.font.CellSize.X), float64(y*uint64(r.font.CellSize.Y))
			ebitenutil.DrawRect(r.frame, pX, pY, float64(r.font.CellSize.X), float64(r.font.CellSize.Y), bg)
			cell := r.buffer.GetCell(uint16(x), uint16(y))
			if cell == nil || cell.Rune().Rune == 0 {
				continue
			}
			text.Draw(r.frame, string(cell.Rune().Rune), r.font.Regular, int(pX), int(pY)+r.font.DotDepth, fg)
		}
	}
}
