package gui

import (
	"fmt"

	"github.com/go-gl/glfw/v3.2/glfw"
)

// send typed runes straight through to the pty
func (gui *GUI) char(w *glfw.Window, r rune) {
	gui.terminal.Write([]byte(string(r)))
}

func modsPressed(pressed glfw.ModifierKey, mods ...glfw.ModifierKey) bool {
	for _, mod := range mods {
		if pressed&mod == 0 {
			return false
		}
		pressed ^= mod
	}
	return pressed == 0
}

func (gui *GUI) key(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {

	if action == glfw.Repeat || action == glfw.Press {

		if gui.overlay != nil {
			if key == glfw.KeyEscape {
				gui.setOverlay(nil)
			}
		}

		for userAction, shortcut := range gui.keyboardShortcuts {

			if shortcut.Match(mods, key) {

				f, ok := actionMap[userAction]
				if ok {
					f(gui)
					break
				}

				switch key {
				case glfw.KeyD:

				case glfw.KeyG:

				case glfw.KeyR:
					gui.launchTarget("https://github.com/liamg/aminal/issues/new/choose")
				case glfw.KeySemicolon:
					gui.config.Slomo = !gui.config.Slomo
					return
				}
			}
		}

		modStr := ""
		switch true {
		case modsPressed(mods, glfw.ModControl, glfw.ModShift, glfw.ModAlt):
			modStr = "8"
		case modsPressed(mods, glfw.ModControl, glfw.ModAlt):
			modStr = "7"
		case modsPressed(mods, glfw.ModControl, glfw.ModShift):
			modStr = "6"
		case modsPressed(mods, glfw.ModControl):
			modStr = "5"
			switch key {
			case glfw.KeyA:
				gui.terminal.Write([]byte{0x1})
				return
			case glfw.KeyB:
				gui.terminal.Write([]byte{0x2})
				return
			case glfw.KeyC: // ctrl^c
				gui.terminal.Write([]byte{0x3}) // send EOT
				return
			case glfw.KeyD:
				gui.terminal.Write([]byte{0x4}) // send EOT
				return
			case glfw.KeyE:
				gui.terminal.Write([]byte{0x5})
				return
			case glfw.KeyF:
				gui.terminal.Write([]byte{0x6})
				return
			case glfw.KeyG:
				gui.terminal.Write([]byte{0x7})
				return
			case glfw.KeyH:
				gui.terminal.Write([]byte{0x08})
				return
			case glfw.KeyI:
				gui.terminal.Write([]byte{0x9})
				return
			case glfw.KeyJ:
				gui.terminal.Write([]byte{0x0a})
				return
			case glfw.KeyK:
				gui.terminal.Write([]byte{0x0b})
				return
			case glfw.KeyL:
				gui.terminal.Write([]byte{0x0c})
				return
			case glfw.KeyM:
				gui.terminal.Write([]byte{0x0d})
				return
			case glfw.KeyN:
				gui.terminal.Write([]byte{0x0e})
				return
			case glfw.KeyO:
				gui.terminal.Write([]byte{0x0f})
				return
			case glfw.KeyP:
				gui.terminal.Write([]byte{0x10})
				return
			case glfw.KeyQ:
				gui.terminal.Write([]byte{0x11})
				return
			case glfw.KeyR:
				gui.terminal.Write([]byte{0x12})
				return
			case glfw.KeyS:
				gui.terminal.Write([]byte{0x13})
				return
			case glfw.KeyT:
				gui.terminal.Write([]byte{0x14})
				return
			case glfw.KeyU:
				gui.terminal.Write([]byte{0x15})
				return
			case glfw.KeyV:
				gui.terminal.Write([]byte{0x16})
				return
			case glfw.KeyW:
				gui.terminal.Write([]byte{0x17})
				return
			case glfw.KeyX:
				gui.terminal.Write([]byte{0x18})
				return
			case glfw.KeyY:
				gui.terminal.Write([]byte{0x19})
				return
			case glfw.KeyZ:
				gui.terminal.Write([]byte{0x1a})
				return
			}
		case modsPressed(mods, glfw.ModAlt, glfw.ModShift):
			modStr = "4"
		case modsPressed(mods, glfw.ModAlt):
			modStr = "3"
		case modsPressed(mods, glfw.ModShift):
			modStr = "2"

		}

		switch key {
		case glfw.KeyF1:
			gui.terminal.Write([]byte{
				0x1b,
				'O',
				'P',
			})
		case glfw.KeyF2:
			gui.terminal.Write([]byte{
				0x1b,
				'O',
				'Q',
			})
		case glfw.KeyF3:
			gui.terminal.Write([]byte{
				0x1b,
				'O',
				'R',
			})
		case glfw.KeyF4:
			gui.terminal.Write([]byte{
				0x1b,
				'O',
				'S',
			})
		case glfw.KeyF5:
			gui.terminal.Write([]byte{
				0x1b,
				'[',
				'1', '5', '~',
			})
		case glfw.KeyF6:
			gui.terminal.Write([]byte{
				0x1b,
				'[',
				'1', '7', '~',
			})
		case glfw.KeyF7:
			gui.terminal.Write([]byte{
				0x1b,
				'[',
				'1', '8', '~',
			})
		case glfw.KeyF8:
			gui.terminal.Write([]byte{
				0x1b,
				'[',
				'1', '9', '~',
			})
		case glfw.KeyF9:
			gui.terminal.Write([]byte{
				0x1b,
				'[',
				'2', '0', '~',
			})
		case glfw.KeyF10:
			gui.terminal.Write([]byte{
				0x1b,
				'[',
				'2', '1', '~',
			})
		case glfw.KeyF11:
			gui.terminal.Write([]byte{
				0x1b,
				'[',
				'2', '3', '~',
			})
		case glfw.KeyF12:
			gui.terminal.Write([]byte{
				0x1b,
				'[',
				'2', '4', '~',
			})
		case glfw.KeyInsert:
			gui.terminal.Write([]byte{
				0x1b,
				'[',
				'2', '~',
			})
		case glfw.KeyDelete:
			gui.terminal.Write([]byte{
				0x1b,
				'[',
				'3', '~',
			})
		case glfw.KeyHome:
			if modStr == "" {
				gui.terminal.Write([]byte("\x1b[1~"))
			} else {
				gui.terminal.Write([]byte(fmt.Sprintf("\x1b[1;%s~", modStr)))
			}
		case glfw.KeyEnd:
			if modStr == "" {
				gui.terminal.Write([]byte("\x1b[4~"))
			} else {
				gui.terminal.Write([]byte(fmt.Sprintf("\x1b[4;%s~", modStr)))
			}
		case glfw.KeyPageUp:
			if modStr == "" {
				gui.terminal.Write([]byte("\x1b[5~"))
			} else {
				gui.terminal.Write([]byte(fmt.Sprintf("\x1b[5;%s~", modStr)))
			}
		case glfw.KeyPageDown:
			if modStr == "" {
				gui.terminal.Write([]byte("\x1b[6~"))
			} else {
				gui.terminal.Write([]byte(fmt.Sprintf("\x1b[6;%s~", modStr)))
			}
		case glfw.KeyEscape:
			if gui.terminal.IsApplicationCursorKeysModeEnabled() {
				gui.terminal.Write([]byte{
					0x1b,
				})
			} else {
				gui.terminal.Write([]byte{
					0x1b,
					'[',
					0x1b,
				})
			}
		case glfw.KeyTab:
			if gui.terminal.IsApplicationCursorKeysModeEnabled() {
				gui.terminal.Write([]byte{
					0x1b,
					'O',
					'I',
				})
			} else {
				gui.terminal.Write([]byte{
					0x09,
				})
			}
		case glfw.KeyEnter:
			gui.terminal.Write([]byte{
				0x0d,
			})
		case glfw.KeyKPEnter:
			if gui.terminal.IsApplicationCursorKeysModeEnabled() {
				gui.terminal.Write([]byte{
					0x1b,
					'O',
					'M',
				})
			} else {
				gui.terminal.Write([]byte{
					0x0d,
				})
			}
		case glfw.KeyBackspace:
			gui.terminal.Write([]byte{0x08})
		case glfw.KeyUp:
			if modStr != "" {
				gui.terminal.Write([]byte(fmt.Sprintf("\x1b[1;%sA", modStr)))
			}

			if gui.terminal.IsApplicationCursorKeysModeEnabled() {
				gui.terminal.Write([]byte{
					0x1b,
					'O',
					'A',
				})
			} else {
				gui.terminal.Write([]byte{
					0x1b,
					'[',
					'A',
				})
			}
		case glfw.KeyDown:

			if modStr != "" {
				gui.terminal.Write([]byte(fmt.Sprintf("\x1b[1;%sB", modStr)))
			}

			if gui.terminal.IsApplicationCursorKeysModeEnabled() {
				gui.terminal.Write([]byte{
					0x1b,
					'O',
					'B',
				})
			} else {
				gui.terminal.Write([]byte{
					0x1b,
					'[',
					'B',
				})
			}
		case glfw.KeyLeft:
			if modStr != "" {
				gui.terminal.Write([]byte(fmt.Sprintf("\x1b[1;%sD", modStr)))
			}

			if gui.terminal.IsApplicationCursorKeysModeEnabled() {
				gui.terminal.Write([]byte{
					0x1b,
					'O',
					'D',
				})
			} else {
				gui.terminal.Write([]byte{
					0x1b,
					'[',
					'D',
				})
			}
		case glfw.KeyRight:
			if modStr != "" {
				gui.terminal.Write([]byte(fmt.Sprintf("\x1b[1;%sC", modStr)))
			}

			if gui.terminal.IsApplicationCursorKeysModeEnabled() {
				gui.terminal.Write([]byte{
					0x1b,
					'O',
					'C',
				})
			} else {
				gui.terminal.Write([]byte{
					0x1b,
					'[',
					'C',
				})
			}
		}

		//gui.logger.Debugf("Key pressed: 0x%X %q", key, string([]byte{byte(key)}))
		//gui.terminal.Write([]byte{byte(scancode)})
	}

}
