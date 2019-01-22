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

func (fm *FontMap) Free() {
	if fm.defaultFont != nil {
		fm.defaultFont.Free()
		fm.defaultFont = nil
	}

	if fm.defaultBoldFont != nil {
		fm.defaultBoldFont.Free()
		fm.defaultBoldFont = nil
	}
}

func (fm *FontMap) AssignFonts(defaultFont *glfont.Font, defaultBoldFont *glfont.Font) {
	fm.defaultFont.Free()
	fm.defaultBoldFont.Free()

	fm.defaultFont = defaultFont
	fm.defaultBoldFont = defaultBoldFont
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
