package gui

import (
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/liamg/darktile/internal/app/darktile/gui/popup"
)

func (g *GUI) getModifierStr() string {
	switch true {
	case g.keyState.RepeatPressed(ebiten.KeyShift) && g.keyState.RepeatPressed(ebiten.KeyControl) && g.keyState.RepeatPressed(ebiten.KeyAlt):
		return ";8"
	case g.keyState.RepeatPressed(ebiten.KeyAlt) && g.keyState.RepeatPressed(ebiten.KeyControl):
		return ";7"
	case g.keyState.RepeatPressed(ebiten.KeyShift) && g.keyState.RepeatPressed(ebiten.KeyControl):
		return ";6"
	case g.keyState.RepeatPressed(ebiten.KeyControl):
		return ";5"
	case g.keyState.RepeatPressed(ebiten.KeyShift) && g.keyState.RepeatPressed(ebiten.KeyAlt):
		return ";4"
	case g.keyState.RepeatPressed(ebiten.KeyAlt):
		return ";3"
	case g.keyState.RepeatPressed(ebiten.KeyShift):
		return ";2"
	}

	return ""
}

// Update changes the terminal GUI state - all user-initiated modification should happen here.
func (g *GUI) Update() error {

	if err := g.handleInput(); err != nil {
		return err
	}

	g.filterPopupMessages()

	return nil
}

func (g *GUI) filterPopupMessages() {
	var filtered []popup.Message
	for _, msg := range g.popupMessages {
		if time.Since(msg.Expiry) >= 0 {
			continue
		}
		filtered = append(filtered, msg)
	}
	g.popupMessages = filtered
}
