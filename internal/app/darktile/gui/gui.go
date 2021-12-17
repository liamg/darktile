package gui

import (
	"fmt"
	"image"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/liamg/darktile/internal/app/darktile/font"
	"github.com/liamg/darktile/internal/app/darktile/gui/popup"
	"github.com/liamg/darktile/internal/app/darktile/hinters"
	"github.com/liamg/darktile/internal/app/darktile/termutil"

	"github.com/hajimehoshi/ebiten/v2"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type GUI struct {
	mouseStateLeft      MouseState
	mouseStateRight     MouseState
	mouseStateMiddle    MouseState
	mouseDrag           bool
	size                image.Point // pixels
	terminal            *termutil.Terminal
	updateChan          chan struct{}
	lastClick           time.Time
	clickCount          int
	fontManager         *font.Manager
	mousePos            termutil.Position
	hinters             []hinters.Hinter
	activeHinter        int
	popupMessages       []popup.Message
	screenshotRequested bool
	screenshotFilename  string
	startupFuncs        []func(g *GUI)
	keyState            *keyState
	opacity             float64
	enableLigatures     bool
	cursorImage         *ebiten.Image
}

type MouseState uint8

const (
	MouseStateNone MouseState = iota
	MouseStatePressed
)

func New(terminal *termutil.Terminal, options ...Option) (*GUI, error) {

	g := &GUI{
		terminal:        terminal,
		size:            image.Point{80, 30},
		updateChan:      make(chan struct{}),
		fontManager:     font.NewManager(),
		activeHinter:    -1,
		keyState:        newKeyState(),
		enableLigatures: true,
	}

	for _, option := range options {
		if err := option(g); err != nil {
			return nil, err
		}
	}

	terminal.SetWindowManipulator(NewManipulator(g))

	return g, nil
}

func (g *GUI) Run() error {

	go func() {
		if err := g.terminal.Run(g.updateChan, uint16(g.size.X), uint16(g.size.Y)); err != nil {
			fmt.Fprintf(os.Stderr, "Fatal error: %s", err)
			os.Exit(1)
		}
		os.Exit(0)
	}()

	ebiten.SetScreenTransparent(true)
	ebiten.SetScreenClearedEveryFrame(true)
	ebiten.SetWindowResizable(true)
	ebiten.SetRunnableOnUnfocused(true)
	ebiten.SetFPSMode(ebiten.FPSModeVsyncOffMinimum)

	for _, f := range g.startupFuncs {
		go f(g)
	}

	go g.watchForUpdate()

	return ebiten.RunGame(g)
}

func (g *GUI) watchForUpdate() {
	for range g.updateChan {
		ebiten.ScheduleFrame()
		go func() {
			if g.keyState.AnythingPressed() {
				time.Sleep(time.Millisecond * 10)
				ebiten.ScheduleFrame()
			}
		}()
	}
}

func (g *GUI) CellSize() image.Point {
	return g.fontManager.CharSize()
}

func (g *GUI) Highlight(start termutil.Position, end termutil.Position, label string, img image.Image) {

	if label == "" && img == nil {
		g.terminal.GetActiveBuffer().Highlight(start, end, nil)
		return
	}

	annotation := &termutil.Annotation{
		Text:  label,
		Image: img,
	}

	if label != "" {
		lines := strings.Split(label, "\n")
		annotation.Height = float64(len(lines))
		for _, line := range lines {
			if float64(len(line)) > annotation.Width {
				annotation.Width = float64(len(line))
			}
		}
	}

	if img != nil {
		annotation.Height += float64(img.Bounds().Dy() / g.fontManager.CharSize().Y)
		if label != "" {
			annotation.Height += 0.5 // half line spacing between image + text
		}
		imgCellWidth := img.Bounds().Dx() / g.fontManager.CharSize().X
		if float64(imgCellWidth) > annotation.Width {
			annotation.Width = float64(imgCellWidth)
		}
	}

	g.terminal.GetActiveBuffer().Highlight(start, end, annotation)
}

func (g *GUI) ClearHighlight() {
	g.terminal.GetActiveBuffer().ClearHighlight()
}
