package gui

import (
	"github.com/go-gl/gl/all-core/gl"
	"github.com/liamg/aminal/hints"
)

type annotation struct {
	hint *hints.Hint
}

func newAnnotation(it *hints.Hint) *annotation {
	return &annotation{
		hint: it,
	}
}

func (a *annotation) render(gui *GUI) {

	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT | gl.STENCIL_BUFFER_BIT)

	lines := gui.terminal.GetVisibleLines()
	for y := 0; y < len(lines); y++ {
		cells := lines[y].Cells()
		for x := 0; x < len(cells); x++ {
			if int(x) >= len(cells) {
				break
			}
			cell := cells[x]

			var colour *[3]float32
			var alpha float32 = 0.6

			if y == int(a.hint.StartY) {
				if x >= int(a.hint.StartX) && x <= int(a.hint.StartX+uint16(len(a.hint.Word))) {
					colour = &[3]float32{0.2, 1.0, 0.2}
					alpha = 1.0
				}
			}
			gui.renderer.DrawCellText(cell, uint(x), uint(y), alpha, colour)
		}
	}

	gui.textbox(a.hint.StartX+1, a.hint.StartY+3, a.hint.Description)

}
