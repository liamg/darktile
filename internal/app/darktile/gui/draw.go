package gui

import (
	"image/color"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/liamg/darktile/internal/app/darktile/termutil"
	imagefont "golang.org/x/image/font"
)

// Draw renders the terminal GUI to the ebtien window. Required to implement the ebiten interface.
func (g *GUI) Draw(screen *ebiten.Image) {

	tmp := ebiten.NewImage(g.size.X, g.size.Y)

	cellSize := g.fontManager.CharSize()
	dotDepth := g.fontManager.DotDepth()

	buffer := g.terminal.GetActiveBuffer()

	regularFace := g.fontManager.RegularFontFace()
	boldFace := g.fontManager.BoldFontFace()
	italicFace := g.fontManager.ItalicFontFace()
	boldItalicFace := g.fontManager.BoldItalicFontFace()

	var useFace imagefont.Face

	defBg := g.terminal.Theme().DefaultBackground()
	defFg := g.terminal.Theme().DefaultForeground()

	var colour color.Color

	endX := float64(cellSize.X * int(buffer.ViewWidth()))
	endY := float64(cellSize.Y * int(buffer.ViewHeight()))
	extraW := float64(g.size.X) - endX
	extraH := float64(g.size.Y) - endY
	if extraW > 0 {
		ebitenutil.DrawRect(tmp, endX, 0, extraW, endY, defBg)
	}
	if extraH > 0 {
		ebitenutil.DrawRect(tmp, 0, endY, float64(g.size.X), extraH, defBg)
	}

	var inHighlight bool
	var highlightRendered bool
	var highlightMin termutil.Position
	highlightMin.Col = uint16(g.size.X)
	highlightMin.Line = uint64(g.size.Y)
	var highlightMax termutil.Position

	for y := int(buffer.ViewHeight() - 1); y >= 0; y-- {
		py := cellSize.Y * y

		ebitenutil.DrawRect(tmp, 0, float64(py), float64(g.size.X), float64(cellSize.Y), defBg)
		inHighlight = false
		for x := uint16(0); x < buffer.ViewWidth(); x++ {
			cell := buffer.GetCell(x, uint16(y))
			px := cellSize.X * int(x)
			if cell != nil {
				colour = cell.Bg()
			} else {
				colour = defBg
			}
			isCursor := g.terminal.GetActiveBuffer().IsCursorVisible() && int(buffer.CursorLine()) == y && buffer.CursorColumn() == x
			if isCursor {
				colour = g.terminal.Theme().CursorBackground()
			} else if buffer.InSelection(termutil.Position{
				Line: uint64(y),
				Col:  x,
			}) {
				colour = g.terminal.Theme().SelectionBackground()
			} else if colour == nil {
				colour = defBg
			}

			ebitenutil.DrawRect(tmp, float64(px), float64(py), float64(cellSize.X), float64(cellSize.Y), colour)

			if buffer.IsHighlighted(termutil.Position{
				Line: uint64(y),
				Col:  x,
			}) {

				if !inHighlight {
					highlightRendered = true
				}

				if uint64(y) < highlightMin.Line {
					highlightMin.Col = uint16(g.size.X)
					highlightMin.Line = uint64(y)
				}
				if uint64(y) > highlightMax.Line {
					highlightMax.Line = uint64(y)
				}
				if uint64(y) == highlightMax.Line && x > highlightMax.Col {
					highlightMax.Col = x
				}
				if uint64(y) == highlightMin.Line && x < highlightMin.Col {
					highlightMin.Col = x
				}

				inHighlight = true

			} else if inHighlight {
				inHighlight = false
			}

			if isCursor && !ebiten.IsFocused() {
				ebitenutil.DrawRect(tmp, float64(px)+1, float64(py)+1, float64(cellSize.X)-2, float64(cellSize.Y)-2, g.terminal.Theme().DefaultBackground())
			}
		}
		for x := uint16(0); x < buffer.ViewWidth(); x++ {
			cell := buffer.GetCell(x, uint16(y))
			if cell == nil || cell.Rune().Rune == 0 {
				continue
			}

			px := cellSize.X * int(x)
			colour = cell.Fg()
			if g.terminal.GetActiveBuffer().IsCursorVisible() && int(buffer.CursorLine()) == y && buffer.CursorColumn() == x {
				colour = g.terminal.Theme().CursorForeground()
			} else if buffer.InSelection(termutil.Position{
				Line: uint64(y),
				Col:  x,
			}) {
				colour = g.terminal.Theme().SelectionForeground()
			} else if colour == nil {
				colour = defFg
			}

			useFace = regularFace
			if cell.Bold() && cell.Italic() {
				useFace = boldItalicFace
			} else if cell.Bold() {
				useFace = boldFace
			} else if cell.Italic() {
				useFace = italicFace
			}

			if cell.Underline() {
				uly := float64(py + (dotDepth+cellSize.Y)/2)
				ebitenutil.DrawLine(tmp, float64(px), uly, float64(px+cellSize.X), uly, colour)
			}

			text.Draw(tmp, string(cell.Rune().Rune), useFace, px, py+dotDepth, colour)

			if cell.Strikethrough() {
				ebitenutil.DrawLine(tmp, float64(px), float64(py+(cellSize.Y/2)), float64(px+cellSize.X), float64(py+(cellSize.Y/2)), colour)
			}

		}
	}

	for _, sixel := range buffer.GetVisibleSixels() {
		sx := float64(int(sixel.Sixel.X) * cellSize.X)
		sy := float64(sixel.ViewLineOffset * cellSize.Y)

		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(sx, sy)
		tmp.DrawImage(
			ebiten.NewImageFromImage(sixel.Sixel.Image),
			op,
		)
	}

	// draw annotations and overlays
	if highlightRendered {
		if annotation := buffer.GetHighlightAnnotation(); annotation != nil {

			if highlightMin.Col == uint16(g.size.X) {
				highlightMin.Col = 0
			}
			if highlightMin.Line == uint64(g.size.Y) {
				highlightMin.Line = 0
			}

			mx, _ := ebiten.CursorPosition()
			padding := float64(cellSize.X) / 2
			lineX := float64(mx)
			var lineY float64
			var lineHeight float64
			annotationX := mx - cellSize.X*2
			var annotationY float64
			annotationWidth := float64(cellSize.X) * annotation.Width
			var annotationHeight float64

			if annotationX+int(annotationWidth)+int(padding*2) > g.size.X {
				annotationX = g.size.X - (int(annotationWidth) + int(padding*2))
			}
			if annotationX < int(padding) {
				annotationX = int(padding)
			}

			if (highlightMin.Line + (highlightMax.Line-highlightMin.Line)/2) < uint64(buffer.ViewHeight()/2) {
				// annotate underneath max

				pixelsUnderHighlight := float64(g.size.Y) - float64((highlightMax.Line+1)*uint64(cellSize.Y))
				// we need to reserve at least one cell height for the label line
				pixelsAvailableY := pixelsUnderHighlight - float64(cellSize.Y)
				annotationHeight = annotation.Height * float64(cellSize.Y)
				if annotationHeight > pixelsAvailableY {
					annotationHeight = pixelsAvailableY
				}

				lineHeight = pixelsUnderHighlight - padding - annotationHeight
				if lineHeight > annotationHeight {
					if annotationHeight > float64(cellSize.Y)*3 {
						lineHeight = annotationHeight
					} else {
						lineHeight = float64(cellSize.Y) * 3
					}
				}
				annotationY = float64((highlightMax.Line+1)*uint64(cellSize.Y)) + lineHeight + float64(padding)
				lineY = float64((highlightMax.Line + 1) * uint64(cellSize.Y))

			} else {
				//annotate above min

				pixelsAboveHighlight := float64((highlightMin.Line) * uint64(cellSize.Y))
				// we need to reserve at least one cell height for the label line
				pixelsAvailableY := pixelsAboveHighlight - float64(cellSize.Y)
				annotationHeight = annotation.Height * float64(cellSize.Y)
				if annotationHeight > pixelsAvailableY {
					annotationHeight = pixelsAvailableY
				}

				lineHeight = pixelsAboveHighlight - annotationHeight
				if lineHeight > annotationHeight {
					if annotationHeight > float64(cellSize.Y)*3 {
						lineHeight = annotationHeight
					} else {
						lineHeight = float64(cellSize.Y) * 3
					}
				}
				annotationY = float64((highlightMin.Line)*uint64(cellSize.Y)) - lineHeight - float64(padding*2) - annotationHeight
				lineY = annotationY + annotationHeight + +padding
			}

			// draw opaque box below and above highlighted line(s)
			ebitenutil.DrawRect(tmp, 0, float64(highlightMin.Line*uint64(cellSize.Y)), float64(cellSize.X*int(highlightMin.Col)), float64(cellSize.Y), color.RGBA{A: 0x80})
			ebitenutil.DrawRect(tmp, float64((cellSize.X)*int(highlightMax.Col+1)), float64(highlightMax.Line*uint64(cellSize.Y)), float64(g.size.X), float64(cellSize.Y), color.RGBA{A: 0x80})
			ebitenutil.DrawRect(tmp, 0, 0, float64(g.size.X), float64(highlightMin.Line*uint64(cellSize.Y)), color.RGBA{A: 0x80})
			afterLineY := float64((1 + highlightMax.Line) * uint64(cellSize.Y))
			ebitenutil.DrawRect(tmp, 0, afterLineY, float64(g.size.X), float64(g.size.Y)-afterLineY, color.RGBA{A: 0x80})

			// annotation border
			ebitenutil.DrawRect(tmp, float64(annotationX)-padding, annotationY-padding, float64(annotationWidth)+(padding*2), annotationHeight+(padding*2), g.terminal.Theme().SelectionBackground())
			// annotation background
			ebitenutil.DrawRect(tmp, 1+float64(annotationX)-padding, 1+annotationY-padding, float64(annotationWidth)+(padding*2)-2, annotationHeight+(padding*2)-2, g.terminal.Theme().DefaultBackground())

			// vertical line
			ebitenutil.DrawLine(tmp, lineX, float64(lineY), lineX, lineY+lineHeight, g.terminal.Theme().SelectionBackground())

			var tY int
			var tX int

			if annotation.Image != nil {
				tY += annotation.Image.Bounds().Dy() + cellSize.Y/2

				op := &ebiten.DrawImageOptions{}
				op.GeoM.Translate(float64(annotationX), annotationY)
				tmp.DrawImage(
					ebiten.NewImageFromImage(annotation.Image),
					op,
				)
			}

			for _, r := range annotation.Text {
				if r == '\n' {
					tY += cellSize.Y
					tX = 0
					continue
				}
				text.Draw(tmp, string(r), regularFace, annotationX+tX, int(annotationY)+dotDepth+tY, g.terminal.Theme().DefaultForeground())
				tX += cellSize.X
			}

		}
	}

	if len(g.popupMessages) > 0 {
		pad := cellSize.Y / 2 // horizontal and vertical padding
		msgEndY := endY
		for _, msg := range g.popupMessages {

			lines := strings.Split(msg.Text, "\n")

			msgX := pad

			msgY := msgEndY - float64(pad*3) - float64(cellSize.Y*len(lines))

			msgText := msg.Text

			boxWidth := float64(pad*2) + float64(cellSize.X*len(msgText))
			boxHeight := float64(pad*2) + float64(cellSize.Y*len(lines))

			if boxWidth < endX/8 {
				boxWidth = endX / 8
			}

			ebitenutil.DrawRect(tmp, float64(msgX-1), msgY-1, boxWidth+2, boxHeight+2, msg.Foreground)
			ebitenutil.DrawRect(tmp, float64(msgX), msgY, boxWidth, boxHeight, msg.Background)
			for y, line := range lines {
				for x, r := range line {
					text.Draw(tmp, string(r), regularFace, msgX+pad+(x*cellSize.X), pad+(y*cellSize.Y)+int(msgY)+dotDepth, msg.Foreground)
				}
			}
			msgEndY = msgEndY - float64(pad*4) - float64(len(lines)*g.CellSize().Y)
		}
	}

	if g.screenshotRequested {
		g.takeScreenshot(tmp)
	}

	opt := &ebiten.DrawImageOptions{}
	opt.ColorM.Scale(1, 1, 1, g.opacity)
	screen.DrawImage(tmp, opt)
}
