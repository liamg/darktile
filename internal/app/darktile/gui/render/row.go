package render

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text"
	imagefont "golang.org/x/image/font"
)

func (r *Render) drawRow(viewY int, defaultBackgroundColour color.Color, defaultForegroundColour color.Color) {

	pixelY := r.font.CellSize.Y * viewY

	// draw a default colour background image across the entire row background
	ebitenutil.DrawRect(r.frame, 0, float64(pixelY), float64(r.pixelWidth), float64(r.font.CellSize.Y), defaultBackgroundColour)

	var colour color.Color

	// draw background for each cell in row
	for viewX := uint16(0); viewX < r.buffer.ViewWidth(); viewX++ {
		cell := r.buffer.GetCell(viewX, uint16(viewY))
		pixelX := r.font.CellSize.X * int(viewX)
		if cell != nil {
			colour = cell.Bg()
		}
		if colour == nil {
			colour = defaultBackgroundColour
		}

		ebitenutil.DrawRect(r.frame, float64(pixelX), float64(pixelY), float64(r.font.CellSize.X), float64(r.font.CellSize.Y), colour)
	}

	var useFace imagefont.Face
	var skipRunes int

	// draw text content of each cell in row
	for viewX := uint16(0); viewX < r.buffer.ViewWidth(); viewX++ {

		cell := r.buffer.GetCell(viewX, uint16(viewY))

		// we don't need to draw empty cells
		if cell == nil || cell.Rune().Rune == 0 {
			continue
		}
		colour = cell.Fg()
		if colour == nil {
			colour = defaultForegroundColour
		}

		// pick a font face for the cell
		if !cell.Bold() && !cell.Italic() {
			useFace = r.font.Regular
		} else if cell.Bold() && cell.Italic() {
			useFace = r.font.Italic
		} else if cell.Bold() {
			useFace = r.font.Bold
		} else if cell.Italic() {
			useFace = r.font.Italic
		}

		pixelX := r.font.CellSize.X * int(viewX)

		// underline the cell content if required
		if cell.Underline() {
			underlinePixelY := float64(pixelY + (r.font.DotDepth+r.font.CellSize.Y)/2)
			ebitenutil.DrawLine(r.frame, float64(pixelX), underlinePixelY, float64(pixelX+r.font.CellSize.X), underlinePixelY, colour)
		}

		// strikethrough the cell if required
		if cell.Strikethrough() {
			ebitenutil.DrawLine(
				r.frame,
				float64(pixelX),
				float64(pixelY+(r.font.CellSize.Y/2)),
				float64(pixelX+r.font.CellSize.X),
				float64(pixelY+(r.font.CellSize.Y/2)),
				colour,
			)
		}

		if r.enableLigatures && skipRunes == 0 {
			skipRunes = r.handleLigatures(viewX, uint16(viewY), useFace, colour)
		}

		if skipRunes > 0 {
			skipRunes--
			continue
		}

		// draw the text for the cell
		text.Draw(r.frame, string(cell.Rune().Rune), useFace, pixelX, pixelY+r.font.DotDepth, colour)
	}
}
