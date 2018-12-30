package terminal

import (
	"fmt"
	"image"
	"image/draw"
	"math"
	"strings"

	"github.com/liamg/aminal/sixel"
)

type boolFormRuneFunc func(rune) bool

func swallowByFunction(pty chan rune, isTerminator boolFormRuneFunc) {
	for {
		b := <-pty
		if isTerminator(b) {
			break
		}
	}
}

func sixelHandler(pty chan rune, terminal *Terminal) error {
	debug := ""

	data := []rune{}

	for {
		b := <-pty
		if b == 0x1b { // terminated by ESC bell or ESC \
			t := <-pty
			if t == '[' { // Windows injected a CSI sequence, let's ignore it
				final, param, _ := loadCSI(pty)
				debug += "[CSI " + param + string(final) + "]"
				continue
			}
			if t == ']' { // Windows injected an OCS sequence, let's ignore it
				swallowByFunction(pty, terminal.IsOSCTerminator)
				debug += "[OSC]"
				continue
			}
			if t != 0x07 && t != 0x5c {
				return fmt.Errorf("Incorrect terminator in sixel sequence: 0x%02X [%c]", t, t)
			}
			break
		}

		if b == 0x0A {
			terminal.logger.Debugf("Sixel line: %s", debug)
			debug = ""
		} else {
			debug += string(b)
		}

		if b >= 33 {
			data = append(data, b)
		}
	}

	if debug != "" {
		terminal.logger.Debugf("Sixel last line: %s", debug)
	}

	six, err := sixel.ParseString(string(data))
	if err != nil {
		return fmt.Errorf("Failed to parse sixel data: %s", err)
	}

	originalImage := six.RGBA()

	w := originalImage.Bounds().Size().X
	h := originalImage.Bounds().Size().Y

	x, y := terminal.ActiveBuffer().CursorColumn(), terminal.ActiveBuffer().CursorLine()

	fromBottom := int(terminal.ActiveBuffer().ViewHeight() - y)
	lines := int(math.Ceil(float64(h) / float64(terminal.charHeight)))
	if fromBottom < lines+2 {
		y -= (uint16(lines+2) - uint16(fromBottom))
	}
	for l := 0; l <= int(lines); l++ {
		terminal.ActiveBuffer().Write([]rune(strings.Repeat(" ", int(terminal.ActiveBuffer().ViewWidth())))...)
		terminal.ActiveBuffer().NewLine()
	}
	cols := int(math.Ceil(float64(w) / float64(terminal.charWidth)))

	for offsetY := 0; offsetY < lines-1; offsetY++ {
		for offsetX := 0; offsetX < cols-1; offsetX++ {

			cell := terminal.ActiveBuffer().GetCell(x+uint16(offsetX), y+uint16((lines-2)-offsetY))
			if cell == nil {
				continue
			}
			img := originalImage.SubImage(image.Rect(
				offsetX*int(terminal.charWidth),
				offsetY*int(terminal.charHeight),
				(offsetX*int(terminal.charWidth))+int(terminal.charWidth),
				(offsetY*int(terminal.charHeight))+int(terminal.charHeight),
			))

			rgba := image.NewRGBA(image.Rect(0, 0, int(terminal.charWidth), int(terminal.charHeight)))
			draw.Draw(rgba, rgba.Bounds(), img, img.Bounds().Min, draw.Src)
			cell.SetImage(rgba)
		}
	}

	return nil
}
