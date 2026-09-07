package termutil

import "image/color"

type ThemeFactory struct {
	theme     *Theme
	colourMap map[Colour]color.Color
}

func NewThemeFactory() *ThemeFactory {
	return &ThemeFactory{
		theme: &Theme{
			colourMap: map[Colour]color.Color{},
		},
		colourMap: make(map[Colour]color.Color),
	}
}

func (t *ThemeFactory) Build() *Theme {
	for id, col := range t.colourMap {
		// RGBA() returns 16-bit channels, so scale by 0x101 (not 0xff) to land
		// back on the exact 8-bit value: 0xcccc -> 0xcc, not 0xcd.
		r, g, b, _ := col.RGBA()
		t.theme.colourMap[id] = color.RGBA{
			R: uint8(r >> 8),
			G: uint8(g >> 8),
			B: uint8(b >> 8),
			A: 0xff,
		}
	}
	return t.theme
}

func (t *ThemeFactory) WithColour(key Colour, colour color.Color) *ThemeFactory {
	t.colourMap[key] = colour
	return t
}
