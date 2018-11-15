package gui

import "github.com/liamg/aminal/glfont"

type FontMap struct {
	defaultFont     *glfont.Font
	defaultBoldFont *glfont.Font
	runeMap         map[rune]*glfont.Font
	ranges          map[runeRange]*glfont.Font
}

type runeRange struct {
	start rune
	end   rune // inclusive
}

func NewFontMap(defaultFont *glfont.Font, defaultBoldFont *glfont.Font) *FontMap {
	return &FontMap{
		defaultFont:     defaultFont,
		defaultBoldFont: defaultBoldFont,
		runeMap:         map[rune]*glfont.Font{},
		ranges:          map[runeRange]*glfont.Font{},
	}
}

func (fm *FontMap) UpdateResolution(w int, h int) {
	fm.defaultFont.UpdateResolution(w, h)
	fm.defaultBoldFont.UpdateResolution(w, h)
	for _, f := range fm.ranges {
		f.UpdateResolution(w, h)
	}
}

func (fm *FontMap) findOverride(r rune) *glfont.Font {

	override, ok := fm.runeMap[r]
	if ok {
		return override
	}

	for rr, f := range fm.ranges {
		if r >= rr.start && r <= rr.end {
			fm.runeMap[r] = f
			return f
		}
	}

	return nil
}

func (fm *FontMap) setOverrideRange(start rune, end rune, font *glfont.Font) {
	fm.ranges[runeRange{start: start, end: end}] = font
}

func (fm *FontMap) GetFont(r rune) *glfont.Font {
	if r <= 0xff {
		return fm.defaultFont
	}

	if f := fm.findOverride(r); f != nil {
		return f
	}

	return fm.defaultFont
}

func (fm *FontMap) GetBoldFont(r rune) *glfont.Font {
	if r <= 0xff {
		return fm.defaultBoldFont
	}

	if f := fm.findOverride(r); f != nil {
		return f
	}

	return fm.defaultBoldFont
}
