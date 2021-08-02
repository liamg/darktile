package gui

import (
	"fmt"
	"image/color"
	"time"

	"github.com/liamg/darktile/internal/app/darktile/gui/popup"
)

const (
	popupMessageDisplayDuration = time.Second * 5
	popupErrorDisplayDuration   = time.Second * 10
)

func (g *GUI) ShowPopup(msg string, fg color.Color, bg color.Color, duration time.Duration) {
	g.popupMessages = append(g.popupMessages, popup.Message{
		Text:       msg,
		Expiry:     time.Now().Add(duration),
		Foreground: fg,
		Background: bg,
	})
}

func (g *GUI) ShowError(msg string) {
	g.ShowPopup(fmt.Sprintf("Error!\n%s", msg), color.White, color.RGBA{A: 0xff, R: 0xff}, popupErrorDisplayDuration)
}

func (g *GUI) ShowMessage(msg string) {
	g.ShowPopup(msg, color.White, color.RGBA{A: 0xff, G: 0x40, R: 0x40, B: 0xff}, popupMessageDisplayDuration)
}
