package gui

type Option func(g *GUI) error

func WithFontFamily(family string) func(g *GUI) error {
	return func(g *GUI) error {
		return g.fontManager.SetFontByFamilyName(family)
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

func WithStartupFunc(f func(g *GUI)) Option {
	return func(g *GUI) error {
		g.startupFuncs = append(g.startupFuncs, f)
		return nil
	}
}
