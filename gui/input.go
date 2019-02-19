package gui

import (
	"fmt"
	"strings"

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

func getModStr(mods glfw.ModifierKey) string {

	switch true {
	case modsPressed(mods, glfw.ModControl, glfw.ModShift, glfw.ModAlt):
		return "8"
	case modsPressed(mods, glfw.ModControl, glfw.ModAlt):
		return "7"
	case modsPressed(mods, glfw.ModControl, glfw.ModShift):
		return "6"
	case modsPressed(mods, glfw.ModControl):
		return "5"
	case modsPressed(mods, glfw.ModAlt, glfw.ModShift):
		return "4"
	case modsPressed(mods, glfw.ModAlt):
		return "3"
	case modsPressed(mods, glfw.ModShift):
		return "2"
	}

	return ""
}

func (gui *GUI) key(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {

	if action == glfw.Repeat || action == glfw.Press {

		if gui.overlay != nil {
			if key == glfw.KeyEscape {
				gui.setOverlay(nil)
			}
		}

		// get key name to handle alternative keyboard layouts
		name := glfw.GetKeyName(key, scancode)
		if len(name) == 1 {
			r := rune(strings.ToLower(name)[0])
			for userAction, shortcut := range gui.keyboardShortcuts {
				if shortcut.Match(mods, r) {
					f, ok := actionMap[userAction]
					if ok {
						f(gui)
						break
					}
				}
			}

			// standard ctrl codes e.g. ^C
			if modsPressed(mods, glfw.ModControl) {
				if r >= 97 && r < 123 {
					gui.terminal.Write([]byte{byte(r) - 96})
					return
				} else if r >= 65 && r < 91 {
					gui.terminal.Write([]byte{byte(r) - 64})
					return
				}
			}
		}

		modStr := getModStr(mods)

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
			if gui.terminal.IsApplicationCursorKeysModeEnabled() {
				if modStr == "" {
					gui.terminal.Write([]byte("\x1b[1~"))
				} else {
					gui.terminal.Write([]byte(fmt.Sprintf("\x1b[1;%s~", modStr)))
				}
			} else {
				gui.terminal.Write([]byte("\x1b[H"))
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
			gui.terminal.Write([]byte{
				0x1b,
			})
		case glfw.KeyTab:
			gui.terminal.Write([]byte{
				0x09,
			})
		case glfw.KeyEnter:
			gui.terminal.WriteReturn()
		case glfw.KeyKPEnter:
			if gui.terminal.IsApplicationCursorKeysModeEnabled() {
				gui.terminal.Write([]byte{
					0x1b,
					'O',
					'M',
				})
			} else {
				gui.terminal.WriteReturn()
			}
		case glfw.KeyBackspace:
			if modsPressed(mods, glfw.ModAlt) {
				gui.terminal.Write([]byte{0x17}) // ctrl-w/delete word
			} else {
				gui.terminal.Write([]byte{0x7f}) //0x7f is DEL
			}
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
	}

}
