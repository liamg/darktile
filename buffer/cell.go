package buffer

import (
	"image"

	"github.com/go-gl/gl/all-core/gl"
)

type Cell struct {
	r     rune
	attr  CellAttributes
	image *image.RGBA
}

type CellAttributes struct {
	FgColour  [3]float32
	BgColour  [3]float32
	Bold      bool
	Dim       bool
	Underline bool
	Blink     bool
	Reverse   bool
	Hidden    bool
}

func (cell *Cell) Image() *image.RGBA {
	return cell.image
}

func (cell *Cell) SetImage(img *image.RGBA) {

	cell.image = img

}

func (cell *Cell) DrawImage(x, y float32) {

	if cell.image == nil {
		return
	}

	var tex uint32
	gl.Enable(gl.TEXTURE_2D)
	gl.GenTextures(1, &tex)
	gl.BindTexture(gl.TEXTURE_2D, tex)
	gl.TexParameterf(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
	gl.TexParameterf(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
	gl.TexParameterf(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameterf(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)

	gl.TexImage2D(
		gl.TEXTURE_2D,
		0,
		gl.RGBA,
		int32(cell.image.Bounds().Size().X),
		int32(cell.image.Bounds().Size().Y),
		0,
		gl.RGBA,
		gl.UNSIGNED_BYTE,
		gl.Ptr(cell.image.Pix),
	)
	gl.BindTexture(gl.TEXTURE_2D, 0)
	gl.Disable(gl.TEXTURE_2D)

	gl.Disable(gl.BLEND)

	var w float32 = float32(cell.image.Bounds().Size().X)
	var h float32 = float32(cell.image.Bounds().Size().Y)

	var readFboId uint32
	gl.GenFramebuffers(1, &readFboId)
	gl.BindFramebuffer(gl.READ_FRAMEBUFFER, readFboId)

	gl.FramebufferTexture2D(gl.READ_FRAMEBUFFER, gl.COLOR_ATTACHMENT0,
		gl.TEXTURE_2D, tex, 0)
	gl.BlitFramebuffer(0, 0, int32(w), int32(h),
		int32(x), int32(y), int32(x+w), int32(y+h),
		gl.COLOR_BUFFER_BIT, gl.LINEAR)
	gl.BindFramebuffer(gl.READ_FRAMEBUFFER, 0)
	gl.DeleteFramebuffers(1, &readFboId)

}

func (cell *Cell) Attr() CellAttributes {
	return cell.attr
}

func (cell *Cell) Rune() rune {
	return cell.r
}

func (cell *Cell) Fg() [3]float32 {
	return cell.attr.FgColour
}

func (cell *Cell) Bg() [3]float32 {
	return cell.attr.BgColour
}

func (cell *Cell) erase() {
	cell.setRune(0)
}

func (cell *Cell) setRune(r rune) {
	cell.r = r
}

func NewBackgroundCell(colour [3]float32) Cell {
	return Cell{
		attr: CellAttributes{
			BgColour: colour,
		},
	}
}
