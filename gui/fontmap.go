package gui

import "github.com/liamg/aminal/glfont"

type FontMap struct {
	defaultFont     *glfont.Font
	defaultBoldFont *glfont.Font
}

func NewFontMap(defaultFont *glfont.Font, defaultBoldFont *glfont.Font) *FontMap {
	return &FontMap{
		defaultFont:     defaultFont,
		defaultBoldFont: defaultBoldFont,
	}
}

func (fm *FontMap) UpdateResolution(w int, h int) {
	fm.defaultFont.UpdateResolution(w, h)
	fm.defaultBoldFont.UpdateResolution(w, h)
}

func (fm *FontMap) DefaultFont() *glfont.Font {

	return fm.defaultFont
}

func (fm *FontMap) BoldFont() *glfont.Font {
	return fm.defaultBoldFont
}
