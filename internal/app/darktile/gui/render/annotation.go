package render

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text"
)

func (r *Render) drawAnnotation() {

	// 1. check if we have anything to highlight/annotate
	highlightStart, highlightEnd, ok := r.buffer.GetViewHighlight()
	if !ok {
		return
	}

	// 2. make everything outside of the highlighted area opaque
	dimColour := color.RGBA{A: 0x80} // 50% alpha black overlay to dim non-highlighted area
	for line := 0; line < int(r.buffer.ViewHeight()); line++ {
		if line < int(highlightStart.Line) || line > int(highlightEnd.Line) {
			ebitenutil.DrawRect(
				r.frame,
				0,
				float64(line*r.font.CellSize.Y),
				float64(r.pixelWidth),
				float64(r.font.CellSize.Y),
				dimColour, // 50% alpha black overlay to dim non-highlighted area
			)
			continue
		}

		if line == int(highlightStart.Line) && highlightStart.Col > 0 {
			// we need to dim some content on this line before the highlight starts
			ebitenutil.DrawRect(
				r.frame,
				0,
				float64(line*r.font.CellSize.Y),
				float64(int(highlightStart.Col)*r.font.CellSize.X),
				float64(r.font.CellSize.Y),
				dimColour,
			)
		}

		if line == int(highlightEnd.Line) && highlightEnd.Col < r.buffer.ViewWidth()-2 {
			// we need to dim some content on this line after the highlight ends
			ebitenutil.DrawRect(
				r.frame,
				float64(int(highlightEnd.Col+1)*r.font.CellSize.X),
				float64(line*r.font.CellSize.Y),
				float64(int(r.buffer.ViewWidth()-(highlightEnd.Col+1))*r.font.CellSize.X),
				float64(r.font.CellSize.Y),
				dimColour,
			)
		}
	}

	// 3. annotate the highlighted area (if there is an annotation)
	annotation := r.buffer.GetHighlightAnnotation()
	if annotation == nil {
		return
	}

	mousePixelX, _ := ebiten.CursorPosition()
	padding := float64(r.font.CellSize.X) / 2

	var lineY float64
	var lineHeight float64
	var annotationY float64
	var annotationHeight float64

	if (highlightStart.Line + (highlightEnd.Line-highlightStart.Line)/2) < uint64(r.buffer.ViewHeight()/2) {
		// annotate underneath max

		pixelsUnderHighlight := float64(r.pixelHeight) - float64((highlightEnd.Line+1)*uint64(r.font.CellSize.Y))
		// we need to reserve at least one cell height for the label line
		pixelsAvailableY := pixelsUnderHighlight - float64(r.font.CellSize.Y)
		annotationHeight = annotation.Height * float64(r.font.CellSize.Y)
		if annotationHeight > pixelsAvailableY {
			annotationHeight = pixelsAvailableY
		}

		lineHeight = pixelsUnderHighlight - padding - annotationHeight
		if lineHeight > annotationHeight {
			if annotationHeight > float64(r.font.CellSize.Y)*3 {
				lineHeight = annotationHeight
			} else {
				lineHeight = float64(r.font.CellSize.Y) * 3
			}
		}
		annotationY = float64((highlightEnd.Line+1)*uint64(r.font.CellSize.Y)) + lineHeight + float64(padding)
		lineY = float64((highlightEnd.Line + 1) * uint64(r.font.CellSize.Y))

	} else {
		//annotate above min

		pixelsAboveHighlight := float64((highlightStart.Line) * uint64(r.font.CellSize.Y))
		// we need to reserve at least one cell height for the label line
		pixelsAvailableY := pixelsAboveHighlight - float64(r.font.CellSize.Y)
		annotationHeight = annotation.Height * float64(r.font.CellSize.Y)
		if annotationHeight > pixelsAvailableY {
			annotationHeight = pixelsAvailableY
		}

		lineHeight = pixelsAboveHighlight - annotationHeight
		if lineHeight > annotationHeight {
			if annotationHeight > float64(r.font.CellSize.Y)*3 {
				lineHeight = annotationHeight
			} else {
				lineHeight = float64(r.font.CellSize.Y) * 3
			}
		}
		annotationY = float64((highlightStart.Line)*uint64(r.font.CellSize.Y)) - lineHeight - float64(padding*2) - annotationHeight
		lineY = annotationY + annotationHeight + +padding
	}

	annotationX := mousePixelX - r.font.CellSize.X*2
	annotationWidth := float64(r.font.CellSize.X) * annotation.Width

	// if the annotation box goes off the right side of the terminal, align it against the right side
	if annotationX+int(annotationWidth)+int(padding*2) > r.pixelWidth {
		annotationX = r.pixelWidth - (int(annotationWidth) + int(padding*2))
	}

	// if the annotation is too far left, align it against the left side
	if annotationX < int(padding) {
		annotationX = int(padding)
	}

	// annotation border
	ebitenutil.DrawRect(r.frame, float64(annotationX)-padding, annotationY-padding, float64(annotationWidth)+(padding*2), annotationHeight+(padding*2), r.theme.SelectionBackground())
	// annotation background
	ebitenutil.DrawRect(r.frame, 1+float64(annotationX)-padding, 1+annotationY-padding, float64(annotationWidth)+(padding*2)-2, annotationHeight+(padding*2)-2, r.theme.DefaultBackground())

	// vertical line
	ebitenutil.DrawLine(r.frame, float64(mousePixelX), float64(lineY), float64(mousePixelX), lineY+lineHeight, r.theme.SelectionBackground())

	var tY int
	var tX int

	if annotation.Image != nil {
		tY += annotation.Image.Bounds().Dy() + r.font.CellSize.Y/2
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(annotationX), annotationY)
		r.frame.DrawImage(
			ebiten.NewImageFromImage(annotation.Image),
			op,
		)
	}

	for _, ch := range annotation.Text {
		if ch == '\n' {
			tY += r.font.CellSize.Y
			tX = 0
			continue
		}
		text.Draw(r.frame, string(ch), r.font.Regular, annotationX+tX, int(annotationY)+r.font.DotDepth+tY, r.theme.DefaultForeground())
		tX += r.font.CellSize.X
	}

}
