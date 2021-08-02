package render

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/liamg/darktile/internal/app/darktile/termutil"
)

func (r *Render) drawCursor() {
	//draw cursor
	if !r.buffer.IsCursorVisible() {
		return
	}

	pixelX := float64(int(r.buffer.CursorColumn()) * r.font.CellSize.X)
	pixelY := float64(int(r.buffer.CursorLine()) * r.font.CellSize.Y)
	cell := r.buffer.GetCell(r.buffer.CursorColumn(), r.buffer.CursorLine())

	useFace := r.font.Regular
	if cell != nil {
		if cell.Bold() && cell.Italic() {
			useFace = r.font.BoldItalic
		} else if cell.Bold() {
			useFace = r.font.Bold
		} else if cell.Italic() {
			useFace = r.font.Italic
		}
	}

	pixelW, pixelH := float64(r.font.CellSize.X), float64(r.font.CellSize.Y)

	// empty rect without focus
	if !ebiten.IsFocused() {
		ebitenutil.DrawRect(r.frame, pixelX, pixelY, pixelW, pixelH, r.theme.CursorBackground())
		ebitenutil.DrawRect(r.frame, pixelX+1, pixelY+1, pixelW-2, pixelH-2, r.theme.CursorForeground())
		return
	}

	// draw the cursor shape
	switch r.buffer.GetCursorShape() {
	case termutil.CursorShapeBlinkingBar, termutil.CursorShapeSteadyBar:
		ebitenutil.DrawRect(r.frame, pixelX, pixelY, 2, pixelH, r.theme.CursorBackground())
	case termutil.CursorShapeBlinkingUnderline, termutil.CursorShapeSteadyUnderline:
		ebitenutil.DrawRect(r.frame, pixelX, pixelY+pixelH-2, pixelW, 2, r.theme.CursorBackground())
	default:
		// draw a custom cursor if we have one and there are no characters in the way
		if r.cursorImage != nil && (cell == nil || cell.Rune().Rune == 0) {
			opt := &ebiten.DrawImageOptions{}
			_, h := r.cursorImage.Size()
			ratio := 1 / (float64(h) / float64(r.font.CellSize.Y))
			actualHeight := float64(h) * ratio
			offsetY := (float64(r.font.CellSize.Y) - actualHeight) / 2
			opt.GeoM.Scale(ratio, ratio)
			opt.GeoM.Translate(pixelX, pixelY+offsetY)
			r.frame.DrawImage(r.cursorImage, opt)
			return
		}

		ebitenutil.DrawRect(r.frame, pixelX, pixelY, pixelW, pixelH, r.theme.CursorBackground())

		// we've drawn over the cell contents, so we need to draw it again in the cursor colours
		if cell != nil && cell.Rune().Rune > 0 {
			text.Draw(r.frame, string(cell.Rune().Rune), useFace, int(pixelX), int(pixelY)+r.font.DotDepth, r.theme.CursorForeground())
		}
	}
}
