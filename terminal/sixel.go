package terminal

import (
	"fmt"
	"image"
	"image/draw"
	"math"
	"strings"

	"github.com/liamg/aminal/matrix"
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

func filter(src []rune) []rune {
	result := make([]rune, 0, len(src))
	for _, v := range src {
		if v >= 33 {
			result = append(result, v)
		}
	}
	return result
}

func sixelHandler(pty chan rune, terminal *Terminal) error {
	debug := ""

	// data := []rune{}

	// track for Windows formatting workaround
	scrollOffset := uint16(terminal.GetScrollOffset())
	x := terminal.ActiveBuffer().CursorColumn() + 2 // reserve two bytes for Sixel prefix (ESC P)
	y := terminal.ActiveBuffer().CursorLine()
	scrollingLine := terminal.ActiveBuffer().ViewHeight() - 1
	xStart := x
	yStartWithOffset := y + scrollOffset
	matrix := matrix.NewAutoMatrix() // a simplified version of Buffer
	for {
		b := <-pty
		if b == 0x1b {
			t := <-pty
			if t == '[' { // Windows injected a CSI sequence
				final, param, _ := loadCSI(pty)

				if final == 'H' {
					// position cursor
					params := splitParams(param)
					{
						xT, yT := parseCursorPosition(params) // 1 based
						x = uint16(xT - 1)                    // 0 based
						y = uint16(yT - 1)                    // 0 based
					}
				}
				debug += "[CSI " + param + string(final) + "]"
				continue
			}
			if t == ']' { // Windows injected an OSC sequence
				// TODO: pass through as if it came via normal stream
				swallowByFunction(pty, terminal.IsOSCTerminator)
				debug += "[OSC]"
				continue
			}
			// if re-drawing a region beforethe start of sixel sequencce,
			// ignore all possible ESC pairs (including ESC P)
			if y+scrollOffset < yStartWithOffset || (y+scrollOffset == yStartWithOffset && x < xStart) {
				x += 2
				continue
			}
			if t != 0x07 && t != 0x5c {
				return fmt.Errorf("Incorrect terminator in sixel sequence: 0x%02X [%c]", t, t)
			}
			break // terminated by ESC bell or ESC \
		}

		if b == 0x0d {
			// skip
		} else if b == 0x0a {
			terminal.logger.Debugf("Sixel line: %s", debug)
			debug = ""
			if y == scrollingLine {
				scrollOffset++
			} else {
				y++
			}
			x = 0
		} else if y+scrollOffset < yStartWithOffset || (y+scrollOffset == yStartWithOffset && x < xStart) {
			x++
		} else if b < 32 {
			x++ // always?
		} else {
			debug += string(b)
			matrix.SetAt(b, int(x), int(y+scrollOffset-yStartWithOffset))
			x++
		}
		/*
			if b >= 33 {
				data = append(data, b)
			}
		*/
	}

	if debug != "" {
		terminal.logger.Debugf("Sixel last line: %s", debug)
	}

	newData := matrix.ExtractFrom(int(xStart), 0) // , int(x), int(y+scrollOffset))

	terminal.logger.Debugf("Sixel data: %s", string(newData))

	filteredData := filter(newData)

	six, err := sixel.ParseString(string(filteredData))

	if err != nil {
		return fmt.Errorf("Failed to parse sixel data: %s", err)
	}

	drawSixel(six, terminal)

	return nil
}

func drawSixel(six *sixel.Sixel, terminal *Terminal) {
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
}
