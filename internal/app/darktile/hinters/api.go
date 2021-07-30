package hinters

import (
	"image"

	"github.com/liamg/darktile/internal/app/darktile/termutil"
)

type HintAPI interface {
	ShowMessage(msg string)
	SetCursorToPointer()
	ResetCursor()
	Highlight(start termutil.Position, end termutil.Position, label string, img image.Image)
	ClearHighlight()
	CellSize() image.Point
}
