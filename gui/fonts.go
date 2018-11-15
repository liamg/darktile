package gui

import (
	"bytes"
	"fmt"

	"github.com/gobuffalo/packr"
	"github.com/liamg/aminal/glfont"
)

func (gui *GUI) getPackedFont(name string) (*glfont.Font, error) {
	box := packr.NewBox("./packed-fonts")
	fontBytes, err := box.MustBytes(name)
	if err != nil {
		return nil, fmt.Errorf("Packaged font '%s' could not be read: %s", name, err)
	}

	font, err := glfont.LoadFont(bytes.NewReader(fontBytes), gui.fontScale, gui.width, gui.height)
	if err != nil {
		return nil, fmt.Errorf("Font '%s' failed to load: %v", name, err)
	}

	return font, nil
}

func (gui *GUI) loadFonts() error {

	defaultFont, err := gui.getPackedFont("Hack-Regular.ttf")
	if err != nil {
		return err
	}

	boldFont, err := gui.getPackedFont("Hack-Bold.ttf")
	if err != nil {
		return err
	}

	gui.fontMap = NewFontMap(defaultFont, boldFont)

	// add special font usage here

	noto, err := gui.getPackedFont("NotoEmoji-Regular.ttf")
	if err != nil {
		return err
	}

	// misc symbols, lightning bolt etc.
	gui.fontMap.setOverrideRange(0x2600, 0x26FF, noto)

	// emoji
	gui.fontMap.setOverrideRange(0x1F600, 0x1F64F, noto)

	return nil
}
