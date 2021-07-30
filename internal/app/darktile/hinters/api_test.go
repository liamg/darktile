package hinters

import (
	"image"

	"github.com/liamg/darktile/internal/app/darktile/termutil"
)

type TestAPI struct {
	highlighted string
}

func (a *TestAPI) ShowMessage(_ string) {

}

func (a *TestAPI) Highlight(start termutil.Position, end termutil.Position, label string, img image.Image) {
	a.highlighted = label
}

func (a *TestAPI) ClearHighlight() {
	a.highlighted = ""
}

func (a *TestAPI) CellSize() image.Point {
	return image.Point{}
}

func (a *TestAPI) SetCursorToPointer() {

}

func (a *TestAPI) ResetCursor() {

}
