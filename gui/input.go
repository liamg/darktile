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

		gui.logger.Debugf("KEY PRESS: key=0x%X scan=0x%X", key, scancode)

		modStr := ""

		switch true {
		case modsPressed(mods, glfw.ModControl, glfw.ModShift, glfw.ModAlt):
			modStr = "8"
		case modsPressed(mods, glfw.ModControl, glfw.ModAlt):
			modStr = "7"
		case modsPressed(mods, glfw.ModControl, glfw.ModShift):
			modStr = "6"
			switch key {
			case glfw.KeyV:
				if s, err := gui.window.GetClipboardString(); err == nil {
					_ = gui.terminal.Paste([]byte(s))
				}
				return
			case glfw.KeySemicolon:
				gui.config.Slomo = !gui.config.Slomo
				return
			}
		case modsPressed(mods, glfw.ModControl):
			modStr = "5"
			switch key {
			case glfw.KeyC: // ctrl^c
				gui.logger.Debugf("Sending CTRL^C")
				gui.terminal.Write([]byte{0x3}) // send EOT
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
			if gui.terminal.IsApplicationCursorKeysModeEnabled() {
				gui.terminal.Write([]byte{
					0x1b,
					'O',
					'M',
				})
			} else {
				gui.terminal.Write([]byte{0x0d})
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
