package gui

import (
	"fmt"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/liamg/darktile/internal/app/darktile/hinters"
	"github.com/liamg/darktile/internal/app/darktile/termutil"
)

// time allowed between mouse clicks to chain them into e.g. double-click
const clickChainWindowMS = 500

// max duration of a click before it is counted as a drag
const clickMaxDuration = 100

func (g *GUI) handleMouse() error {

	_, scrollY := ebiten.Wheel()
	if scrollY < 0 {
		g.terminal.GetActiveBuffer().ScrollDown(5)
	} else if scrollY > 0 {
		g.terminal.GetActiveBuffer().ScrollUp(5)
	}

	x, y := ebiten.CursorPosition()
	col := x / g.fontManager.CharSize().X
	line := y / g.fontManager.CharSize().Y
	var moved bool

	if col != int(g.mousePos.Col) || line != int(g.mousePos.Line) {
		if col >= 0 && col < int(g.terminal.GetActiveBuffer().ViewWidth()) && line >= 0 && line < int(g.terminal.GetActiveBuffer().ViewHeight()) {
			// mouse moved!
			moved = true
			g.mousePos = termutil.Position{
				Col:  uint16(col),
				Line: uint64(line),
			}
			if !ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
				if err := g.handleMouseMove(g.mousePos); err != nil {
					return err
				}
			}
		} else if err := g.clearHinters(); err != nil {
			return err
		}
	}

	pressedLeft := ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) && g.mouseStateLeft != MouseStatePressed
	pressedMiddle := ebiten.IsMouseButtonPressed(ebiten.MouseButtonMiddle) && g.mouseStateMiddle != MouseStatePressed
	pressedRight := ebiten.IsMouseButtonPressed(ebiten.MouseButtonRight) && g.mouseStateRight != MouseStatePressed
	released := (!ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) && g.mouseStateLeft == MouseStatePressed) ||
		(!ebiten.IsMouseButtonPressed(ebiten.MouseButtonMiddle) && g.mouseStateMiddle == MouseStatePressed) ||
		(!ebiten.IsMouseButtonPressed(ebiten.MouseButtonRight) && g.mouseStateRight == MouseStatePressed)

	defer func() {
		if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
			g.mouseStateLeft = MouseStatePressed
		} else {
			g.mouseStateLeft = MouseStateNone
		}

		if ebiten.IsMouseButtonPressed(ebiten.MouseButtonMiddle) {
			g.mouseStateMiddle = MouseStatePressed
		} else {
			g.mouseStateMiddle = MouseStateNone
		}

		if ebiten.IsMouseButtonPressed(ebiten.MouseButtonRight) {
			g.mouseStateRight = MouseStatePressed
		} else {
			g.mouseStateRight = MouseStateNone
		}
	}()

	if pressedLeft || pressedMiddle || pressedRight || released {
		if g.handleMouseRemotely(x, y, pressedLeft, pressedMiddle, pressedRight, released, moved) {
			return nil
		}
	}

	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		if g.mouseStateLeft == MouseStatePressed {

			if g.mouseDrag {

				// update selection end
				g.terminal.GetActiveBuffer().SetSelectionEnd(termutil.Position{
					Line: uint64(line),
					Col:  uint16(col),
				})
			} else if time.Since(g.lastClick) > time.Millisecond*clickMaxDuration && !ebiten.IsKeyPressed(ebiten.KeyControl) {
				g.mouseDrag = true
			}
		} else {

			if g.clickCount == 0 || time.Since(g.lastClick) < time.Millisecond*clickChainWindowMS {
				g.clickCount++
			} else {
				g.clickCount = 1
			}

			g.lastClick = time.Now()

			handled, err := g.handleClick(g.clickCount, x, y)
			if err != nil {
				return err
			}
			if handled {
				g.mouseDrag = false
				return nil
			}

			//set selection start
			col := x / g.fontManager.CharSize().X
			line := y / g.fontManager.CharSize().Y

			g.terminal.GetActiveBuffer().SetSelectionStart(termutil.Position{
				Line: uint64(line),
				Col:  uint16(col),
			})
		}

		ebiten.ScheduleFrame()

	} else {
		g.mouseDrag = false
	}

	return nil
}

func (g *GUI) clearHinters() error {
	if g.activeHinter > -1 {
		if err := g.hinters[g.activeHinter].Deactivate(g); err != nil {
			return err
		}
		g.activeHinter = -1
	}
	return nil
}

// mouse moved to cell (not during click + drag)
func (g *GUI) handleMouseMove(pos termutil.Position) error {

	// start uses raw coords
	start, _, text, index, ok := g.terminal.GetActiveBuffer().GetBoundedTextAtPosition(pos)
	if !ok {
		g.clearHinters()
		return nil
	}

	activeHinter := -1

	for i, hinter := range g.hinters {
		if ok, offset, length := hinter.Match(text, index); ok {
			match := text[offset : offset+length]

			newStartX := int(start.Col) + offset
			newStartY := start.Line
			for newStartX >= int(g.terminal.GetActiveBuffer().ViewWidth()) {
				newStartX -= int(g.terminal.GetActiveBuffer().ViewWidth())
				newStartY++
			}

			newEndX := newStartX + length - 1
			newEndY := newStartY
			for newEndX > int(g.terminal.GetActiveBuffer().ViewWidth()) {
				newEndX -= int(g.terminal.GetActiveBuffer().ViewWidth())
				newEndY++
			}

			matchStart := termutil.Position{
				Col:  uint16(newStartX),
				Line: newStartY,
			}
			matchEnd := termutil.Position{
				Col:  uint16(newEndX),
				Line: newEndY,
			}

			if err := hinter.Activate(g, match, matchStart, matchEnd); err != nil {
				return err
			}

			activeHinter = i
			break
		}
	}

	// hinter was just deactivated
	if g.activeHinter > -1 && activeHinter == -1 {
		if err := g.clearHinters(); err != nil {
			return err
		}
	}

	g.activeHinter = activeHinter

	return nil
}

func WithHinter(h hinters.Hinter) func(g *GUI) error {
	return func(g *GUI) error {
		g.hinters = append(g.hinters, h)
		return nil
	}
}

func (g *GUI) handleClick(clickCount, x, y int) (bool, error) {

	switch clickCount {
	case 1: // single click
		if ebiten.IsKeyPressed(ebiten.KeyControl) { // ctrl + click to run hinters
			if g.activeHinter > -1 {
				g.hinters[g.activeHinter].Click(g)
			}
		} else {
			g.terminal.GetActiveBuffer().ClearSelection()
		}

	case 2: //double click
		col := uint16(x / g.fontManager.CharSize().X)
		line := uint64(y / g.fontManager.CharSize().Y)
		g.terminal.GetActiveBuffer().SelectWordAt(termutil.Position{Col: col, Line: line}, wordMatcher)
		return true, nil
	default: // triple click (or more!)
		g.terminal.GetActiveBuffer().ExtendSelectionToEntireLines()
		return true, nil
	}

	return false, nil
}

func alphaMatcher(r rune) bool {
	if r >= 65 && r <= 90 {
		return true
	}
	if r >= 97 && r <= 122 {
		return true
	}
	return false
}

func numberMatcher(r rune) bool {
	if r >= 48 && r <= 57 {
		return true
	}
	return false
}

func alphaNumericMatcher(r rune) bool {
	return alphaMatcher(r) || numberMatcher(r)
}

func wordMatcher(r rune) bool {

	if alphaNumericMatcher(r) {
		return true
	}

	if r == '_' {
		return true
	}

	return false
}

func (g *GUI) handleMouseRemotely(x, y int, pressedLeft, pressedMiddle, pressedRight, released, moved bool) bool {

	tx, ty := 1+(x/g.fontManager.CharSize().X), 1+(y/g.fontManager.CharSize().Y)

	mode := g.terminal.GetActiveBuffer().GetMouseMode()

	switch mode {
	case termutil.MouseModeNone:
		return false
	case termutil.MouseModeX10:
		var button rune
		switch true {
		case pressedLeft:
			button = 0
		case pressedMiddle:
			button = 1
		case pressedRight:
			button = 2
		default:
			return true
		}
		packet := fmt.Sprintf("\x1b[M%c%c%c", (rune(button + 32)), (rune(tx + 32)), (rune(ty + 32)))
		_ = g.terminal.WriteToPty([]byte(packet))
		return true
	case termutil.MouseModeVT200, termutil.MouseModeButtonEvent:

		var button rune

		extMode := g.terminal.GetActiveBuffer().GetMouseExtMode()

		switch true {
		case pressedLeft:
			button = 0
		case pressedMiddle:
			button = 1
		case pressedRight:
			button = 2
		case released:
			if extMode != termutil.MouseExtSGR {
				button = 3
			}
		default:
			return true
		}

		if moved && mode == termutil.MouseModeButtonEvent {
			button |= 32
		}

		if ebiten.IsKeyPressed(ebiten.KeyShift) {
			button |= 4
		}

		if ebiten.IsKeyPressed(ebiten.KeyMeta) {
			button |= 8
		}

		if ebiten.IsKeyPressed(ebiten.KeyControl) {
			button |= 16
		}

		var packet string

		if extMode == termutil.MouseExtSGR {
			final := 'M'
			if released {
				final = 'm'
			}
			packet = fmt.Sprintf("\x1b[<%d;%d;%d%c", button, tx, ty, final)
		} else {
			packet = fmt.Sprintf("\x1b[M%c%c%c", button+32, tx+32, ty+32)
		}

		g.terminal.WriteToPty([]byte(packet))
		return true

	}

	return false

}

func (g *GUI) SetCursorToPointer() {
	ebiten.SetCursorShape(ebiten.CursorShapePointer)
}

func (g *GUI) ResetCursor() {
	ebiten.SetCursorShape(ebiten.CursorShapeDefault)
}
