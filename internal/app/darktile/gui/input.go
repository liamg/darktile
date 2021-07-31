package gui

import (
	"fmt"

	"github.com/d-tsuji/clipboard"
	"github.com/hajimehoshi/ebiten/v2"
)

var modifiableKeys = map[ebiten.Key]uint8{
	ebiten.KeyA: 'A',
	ebiten.KeyB: 'B',
	ebiten.KeyC: 'C',
	ebiten.KeyD: 'D',
	ebiten.KeyE: 'E',
	ebiten.KeyF: 'F',
	ebiten.KeyG: 'G',
	ebiten.KeyH: 'H',
	ebiten.KeyI: 'I',
	ebiten.KeyJ: 'J',
	ebiten.KeyK: 'K',
	ebiten.KeyL: 'L',
	ebiten.KeyM: 'M',
	ebiten.KeyN: 'N',
	ebiten.KeyO: 'O',
	ebiten.KeyP: 'P',
	ebiten.KeyQ: 'Q',
	ebiten.KeyR: 'R',
	ebiten.KeyS: 'S',
	ebiten.KeyT: 'T',
	ebiten.KeyU: 'U',
	ebiten.KeyV: 'V',
	ebiten.KeyW: 'W',
	ebiten.KeyX: 'X',
	ebiten.KeyY: 'Y',
	ebiten.KeyZ: 'Z',
}

func (g *GUI) handleInput() error {

	if err := g.handleMouse(); err != nil {
		return err
	}

	switch true {

	case ebiten.IsKeyPressed(ebiten.KeyControl) && ebiten.IsKeyPressed(ebiten.KeyShift):

		switch true {
		case g.keyState.RepeatPressed(ebiten.KeyC):
			content, selection := g.terminal.GetActiveBuffer().GetSelection()
			if selection == nil {
				return nil
			}
			return clipboard.Set(content)
		case g.keyState.RepeatPressed(ebiten.KeyV):
			paste, err := clipboard.Get()
			if err != nil {
				return err
			}
			return g.terminal.WriteToPty([]byte(paste))
		case g.keyState.RepeatPressed(ebiten.KeyBracketLeft):
			g.RequestScreenshot("")
		}

	case ebiten.IsKeyPressed(ebiten.KeyControl):

		for key, ch := range modifiableKeys {
			if g.keyState.RepeatPressed(key) {
				if ch >= 97 && ch < 123 {
					return g.terminal.WriteToPty([]byte{ch - 96})
				} else if ch >= 65 && ch < 91 {
					return g.terminal.WriteToPty([]byte{ch - 64})
				}
			}
		}

		switch true {
		case g.keyState.RepeatPressed(ebiten.KeyMinus):
			g.fontManager.DecreaseSize()
			cellSize := g.fontManager.CharSize()
			cols, rows := g.size.X/cellSize.X, g.size.Y/cellSize.Y
			if err := g.terminal.SetSize(uint16(rows), uint16(cols)); err != nil {
				return err
			}
			return nil
		case g.keyState.RepeatPressed(ebiten.KeyEqual):
			g.fontManager.IncreaseSize()
			cellSize := g.fontManager.CharSize()
			cols, rows := g.size.X/cellSize.X, g.size.Y/cellSize.Y
			if err := g.terminal.SetSize(uint16(rows), uint16(cols)); err != nil {
				return err
			}
			return nil
		default:
			return nil
		}

	case ebiten.IsKeyPressed(ebiten.KeyAlt):

		for key, ch := range modifiableKeys {
			if g.keyState.RepeatPressed(key) {
				return g.terminal.WriteToPty([]byte{0x1b, ch})
			}
		}

	case g.keyState.RepeatPressed(ebiten.KeyArrowUp):
		if g.terminal.GetActiveBuffer().IsApplicationCursorKeysModeEnabled() {
			return g.terminal.WriteToPty([]byte{
				0x1b,
				'O',
				'A',
			})
		} else {
			return g.terminal.WriteToPty([]byte(fmt.Sprintf("\x1b[%sA", g.getModifierStr())))
		}
	case g.keyState.RepeatPressed(ebiten.KeyArrowDown):
		if g.terminal.GetActiveBuffer().IsApplicationCursorKeysModeEnabled() {
			return g.terminal.WriteToPty([]byte{
				0x1b,
				'O',
				'B',
			})
		} else {
			return g.terminal.WriteToPty([]byte(fmt.Sprintf("\x1b[%sB", g.getModifierStr())))
		}
	case g.keyState.RepeatPressed(ebiten.KeyArrowRight):
		if g.terminal.GetActiveBuffer().IsApplicationCursorKeysModeEnabled() {
			return g.terminal.WriteToPty([]byte{
				0x1b,
				'O',
				'C',
			})
		} else {
			return g.terminal.WriteToPty([]byte(fmt.Sprintf("\x1b[%sC", g.getModifierStr())))
		}
	case g.keyState.RepeatPressed(ebiten.KeyArrowLeft):
		if g.terminal.GetActiveBuffer().IsApplicationCursorKeysModeEnabled() {
			return g.terminal.WriteToPty([]byte{
				0x1b,
				'O',
				'D',
			})
		} else {
			return g.terminal.WriteToPty([]byte(fmt.Sprintf("\x1b[%sD", g.getModifierStr())))
		}
	case g.keyState.RepeatPressed(ebiten.KeyEnter):
		if g.terminal.GetActiveBuffer().IsNewLineMode() {
			return g.terminal.WriteToPty([]byte{0x0d, 0x0a})
		}
		return g.terminal.WriteToPty([]byte{0x0d})
	case g.keyState.RepeatPressed(ebiten.KeyNumpadEnter):
		if g.terminal.GetActiveBuffer().IsApplicationCursorKeysModeEnabled() {
			g.terminal.WriteToPty([]byte{
				0x1b,
				'O',
				'M',
			})
		} else {
			if g.terminal.GetActiveBuffer().IsNewLineMode() {
				if err := g.terminal.WriteToPty([]byte{0x0d, 0x0a}); err != nil {
					return err
				}
			}
			return g.terminal.WriteToPty([]byte{0x0d})
		}
	case g.keyState.RepeatPressed(ebiten.KeyTab):
		return g.terminal.WriteToPty([]byte{0x09}) // tab

	case g.keyState.RepeatPressed(ebiten.KeyEscape):
		g.terminal.GetActiveBuffer().ClearSelection()
		g.terminal.GetActiveBuffer().ClearHighlight()
		return g.terminal.WriteToPty([]byte{0x1b}) // escape
	case g.keyState.RepeatPressed(ebiten.KeyBackspace):
		if ebiten.IsKeyPressed(ebiten.KeyAlt) {
			return g.terminal.WriteToPty([]byte{0x17}) // ctrl-w/delete word
		} else {
			return g.terminal.WriteToPty([]byte{0x7f}) //0x7f is DEL
		}
	case g.keyState.RepeatPressed(ebiten.KeyF1):
		return g.terminal.WriteToPty([]byte("\x1bOP"))
	case g.keyState.RepeatPressed(ebiten.KeyF2):
		return g.terminal.WriteToPty([]byte("\x1bOQ"))
	case g.keyState.RepeatPressed(ebiten.KeyF3):
		return g.terminal.WriteToPty([]byte("\x1bOR"))
	case g.keyState.RepeatPressed(ebiten.KeyF4):
		return g.terminal.WriteToPty([]byte("\x1bOS"))
	case g.keyState.RepeatPressed(ebiten.KeyF5):
		return g.terminal.WriteToPty([]byte("\x1b[15~"))
	case g.keyState.RepeatPressed(ebiten.KeyF6):
		return g.terminal.WriteToPty([]byte("\x1b[17~"))
	case g.keyState.RepeatPressed(ebiten.KeyF7):
		return g.terminal.WriteToPty([]byte("\x1b[18~"))
	case g.keyState.RepeatPressed(ebiten.KeyF8):
		return g.terminal.WriteToPty([]byte("\x1b[19~"))
	case g.keyState.RepeatPressed(ebiten.KeyF9):
		return g.terminal.WriteToPty([]byte("\x1b[20~"))
	case g.keyState.RepeatPressed(ebiten.KeyF10):
		return g.terminal.WriteToPty([]byte("\x1b[21~"))
	case g.keyState.RepeatPressed(ebiten.KeyF11):
		return g.terminal.WriteToPty([]byte("\x1b[22~"))
	case g.keyState.RepeatPressed(ebiten.KeyF12):
		return g.terminal.WriteToPty([]byte("\x1b[23~"))
	case g.keyState.RepeatPressed(ebiten.KeyInsert):
		return g.terminal.WriteToPty([]byte("\x1b[2~"))
	case g.keyState.RepeatPressed(ebiten.KeyDelete):
		return g.terminal.WriteToPty([]byte("\x1b[3~"))
	case g.keyState.RepeatPressed(ebiten.KeyHome):
		if g.terminal.GetActiveBuffer().IsApplicationCursorKeysModeEnabled() {
			return g.terminal.WriteToPty([]byte(fmt.Sprintf("\x1b[1%s~", g.getModifierStr())))
		} else {
			return g.terminal.WriteToPty([]byte("\x1b[H"))
		}
	case g.keyState.RepeatPressed(ebiten.KeyEnd):
		return g.terminal.WriteToPty([]byte(fmt.Sprintf("\x1b[4%s~", g.getModifierStr())))
	case g.keyState.RepeatPressed(ebiten.KeyPageUp):
		return g.terminal.WriteToPty([]byte(fmt.Sprintf("\x1b[5%s~", g.getModifierStr())))
	case g.keyState.RepeatPressed(ebiten.KeyPageDown):
		return g.terminal.WriteToPty([]byte(fmt.Sprintf("\x1b[6%s~", g.getModifierStr())))
	default:
		input := ebiten.AppendInputChars(nil)
		for _, runePressed := range input {
			if err := g.terminal.WriteToPty([]byte(string(runePressed))); err != nil {
				return err
			}
		}
	}

	return nil
}
