package gui

import (
	"bytes"
	"fmt"

	"github.com/gobuffalo/packr"
	"github.com/liamg/aminal/glfont"
)

func (gui *GUI) getPackedFont(name string, actualWidth int, actualHeight int) (*glfont.Font, error) {
	box := packr.NewBox("./packed-fonts")
	fontBytes, err := box.Find(name)
	if err != nil {
		return nil, fmt.Errorf("packaged font '%s' could not be read: %s", name, err)
	}

	font, err := glfont.LoadFont(bytes.NewReader(fontBytes), gui.fontScale*gui.dpiScale/gui.scale(), actualWidth, actualHeight)
	if err != nil {
		return nil, fmt.Errorf("font '%s' failed to load: %v", name, err)
	}

	return font, nil
}

func (gui *GUI) loadFonts(actualWidth int, actualHeight int) error {

	// from https://github.com/ryanoasis/nerd-fonts/tree/master/patched-fonts/Hack

	defaultFont, err := gui.getPackedFont("Hack Regular Nerd Font Complete.ttf", actualWidth, actualHeight)
	if err != nil {
		return err
	}

	boldFont, err := gui.getPackedFont("Hack Bold Nerd Font Complete.ttf", actualWidth, actualHeight)
	if err != nil {
		return err
	}

	if gui.fontMap == nil {
		gui.fontMap = NewFontMap(defaultFont, boldFont)
	} else {
		gui.fontMap.AssignFonts(defaultFont, boldFont)
	}

	// add special non-ascii fonts here

	return nil
}
