package render

import (
	"strings"

	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text"
)

func (r *Render) drawPopups() {

	if len(r.popups) == 0 {
		return
	}

	pad := r.font.CellSize.Y / 2 // horizontal and vertical padding
	maxPixelX := float64(r.font.CellSize.X * int(r.buffer.ViewWidth()))
	maxPixelY := float64(r.font.CellSize.Y * int(r.buffer.ViewHeight()))

	for _, msg := range r.popups {

		lines := strings.Split(msg.Text, "\n")
		msgX := pad
		msgY := maxPixelY - float64(pad*3) - float64(r.font.CellSize.Y*len(lines))
		boxWidth := float64(pad*2) + float64(r.font.CellSize.X*len(msg.Text))
		boxHeight := float64(pad*2) + float64(r.font.CellSize.Y*len(lines))

		if boxWidth < maxPixelX/8 {
			boxWidth = maxPixelX / 8
		}

		ebitenutil.DrawRect(r.frame, float64(msgX-1), msgY-1, boxWidth+2, boxHeight+2, msg.Foreground)
		ebitenutil.DrawRect(r.frame, float64(msgX), msgY, boxWidth, boxHeight, msg.Background)
		for y, line := range lines {
			for x, c := range line {
				text.Draw(r.frame, string(c), r.font.Regular, msgX+pad+(x*r.font.CellSize.X), pad+(y*r.font.CellSize.Y)+int(msgY)+r.font.DotDepth, msg.Foreground)
			}
		}
		maxPixelY = maxPixelY - float64(pad*4) - float64(len(lines)*r.font.CellSize.Y)
	}

}
