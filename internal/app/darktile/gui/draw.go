package gui

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/liamg/darktile/internal/app/darktile/gui/render"
)

// Draw renders the terminal GUI to the ebtien window. Required to implement the ebiten interface.
func (g *GUI) Draw(screen *ebiten.Image) {
	render.
		New(screen, g.terminal, g.fontManager, g.popupMessages, g.opacity, g.enableLigatures, g.cursorImage).
		Draw()

	if g.screenshotRequested {
		g.takeScreenshot(screen)
	}
}
