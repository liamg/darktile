package gui

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

type Option func(g *GUI) error

func WithFontFamily(family string) func(g *GUI) error {
	return func(g *GUI) error {
		return g.fontManager.SetFontByFamilyName(family)
	}
}

func WithOpacity(opacity float64) func(g *GUI) error {
	return func(g *GUI) error {
		g.opacity = opacity
		return nil
	}
}

func WithFontSize(size float64) func(g *GUI) error {
	return func(g *GUI) error {
		g.fontManager.SetSize(size)
		return nil
	}
}

func WithFontDPI(dpi float64) func(g *GUI) error {
	return func(g *GUI) error {
		g.fontManager.SetSize(dpi)
		return nil
	}
}

func WithLigatures(enable bool) func(g *GUI) error {
	return func(g *GUI) error {
		g.enableLigatures = enable
		return nil
	}
}

func WithCursorImage(img image.Image) func(g *GUI) error {
	return func(g *GUI) error {
		g.cursorImage = ebiten.NewImageFromImage(img)
		return nil
	}
}

func WithStartupFunc(f func(g *GUI)) Option {
	return func(g *GUI) error {
		g.startupFuncs = append(g.startupFuncs, f)
		return nil
	}
}
