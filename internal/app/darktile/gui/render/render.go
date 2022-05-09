package render

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/liamg/darktile/internal/app/darktile/font"
	"github.com/liamg/darktile/internal/app/darktile/gui/popup"
	"github.com/liamg/darktile/internal/app/darktile/termutil"
	imagefont "golang.org/x/image/font"
)

type Render struct {
	frame           *ebiten.Image
	screen          *ebiten.Image
	terminal        *termutil.Terminal
	buffer          *termutil.Buffer
	theme           *termutil.Theme
	fontManager     *font.Manager
	pixelWidth      int
	pixelHeight     int
	font            Font
	opacity         float64
	popups          []popup.Message
	enableLigatures bool
	cursorImage     *ebiten.Image
}

type Font struct {
	Regular    imagefont.Face
	Bold       imagefont.Face
	Italic     imagefont.Face
	BoldItalic imagefont.Face
	CellSize   image.Point
	DotDepth   int
}

func New(screen *ebiten.Image, terminal *termutil.Terminal, fontManager *font.Manager, popups []popup.Message, opacity float64, enableLigatures bool, cursorImage *ebiten.Image) *Render {
	w, h := screen.Size()
	return &Render{
		screen:      screen,
		frame:       ebiten.NewImage(w, h),
		terminal:    terminal,
		buffer:      terminal.GetActiveBuffer(),
		theme:       terminal.Theme(),
		fontManager: fontManager,
		pixelWidth:  w,
		pixelHeight: h,
		font: Font{
			Regular:    fontManager.RegularFontFace(),
			Bold:       fontManager.BoldFontFace(),
			Italic:     fontManager.ItalicFontFace(),
			BoldItalic: fontManager.BoldItalicFontFace(),
			CellSize:   fontManager.CharSize(),
			DotDepth:   fontManager.DotDepth(),
		},
		opacity:         opacity,
		popups:          popups,
		enableLigatures: enableLigatures,
		cursorImage:     cursorImage,
	}
}

func (r *Render) Draw() {
	r.terminal.Lock()
	defer r.terminal.Unlock()

	// 1. fill frame with default background colour
	r.frame.Fill(r.theme.DefaultBackground())

	// 2. draw content (each row, each cell)
	r.drawContent()

	// 3. draw cursor
	r.drawCursor()

	// // 4. draw sixels
	r.drawSixels()

	// // 5. draw selection
	r.drawSelection()

	// // 6. draw highlight/annotations
	r.drawAnnotation()

	// // 7. draw popups
	r.drawPopups()

	// // 8. apply effects (e.g. transparency)
	r.finalise()

}

func (r *Render) finalise() {
	defer r.frame.Dispose()
	opt := &ebiten.DrawImageOptions{}
	opt.ColorM.Scale(1, 1, 1, r.opacity)
	r.screen.DrawImage(r.frame, opt)
}
