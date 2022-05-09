package gui

import (
	"image"
)

// Layout provides the terminal gui size in pixels. Required to implement the ebiten interface.
func (g *GUI) Layout(outsideWidth, outsideHeight int) (int, int) {

	w, h := outsideWidth, outsideHeight

	if g.size.X != w || g.size.Y != h {
		g.size = image.Point{
			X: w,
			Y: h,
		}
		g.resize(w, h)
	}

	return w, h
}

func (g *GUI) resize(w, h int) {

	if g.fontManager.CharSize().X == 0 || g.fontManager.CharSize().Y == 0 || g.terminal == nil {
		return
	}

	cols := uint16(w / g.fontManager.CharSize().X)
	rows := uint16(h / g.fontManager.CharSize().Y)

	g.terminal.Lock()
	defer g.terminal.Unlock()

	if g.terminal.IsRunning() {
		_ = g.terminal.SetSize(rows, cols)
	}
}
