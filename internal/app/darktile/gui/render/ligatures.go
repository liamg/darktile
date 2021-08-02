package render

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2/text"
	imagefont "golang.org/x/image/font"
)

var ligatures = map[string]rune{
	":=":  '≔',
	"===": '≡',
	"!=":  '≠',
	"!==": '≢',
	"<=":  '≤',
	">=":  '≥',
	"=>":  '⇒',
	"->":  '→',
	"<-":  '←',
	"<>":  '≷',
}

func (r *Render) handleLigatures(sx uint16, sy uint16, face imagefont.Face, colour color.Color) (length int) {

	var candidate string
	for x := sx; x <= sx+2; x++ {
		cell := r.buffer.GetCell(x, sy)
		if cell == nil || cell.Rune().Rune == 0 {
			break
		}
		candidate += string(cell.Rune().Rune)
	}

	for len(candidate) > 1 {
		if ru, ok := ligatures[candidate]; ok {
			// draw ligature
			ligX := (int(sx) * r.font.CellSize.X) + (((len(candidate) - 1) * r.font.CellSize.X) / 2)
			text.Draw(r.frame, string(ru), face, ligX, (int(sy)*r.font.CellSize.Y)+r.font.DotDepth, colour)
			return len(candidate)
		}
		candidate = candidate[:len(candidate)-1]
	}

	return 0
}
