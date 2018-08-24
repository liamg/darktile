package terminal

import (
	"fmt"
	"strconv"
	"strings"
)

var csiSequenceMap = map[rune]csiSequenceHandler{
	'm': sgrSequenceHandler,
	'P': csiDeleteHandler,
	'J': csiEraseInDisplayHandler,
	'K': csiEraseInLineHandler,
	'h': csiSetModeHandler,
	'l': csiResetModeHandler,
	'd': csiLinePositionAbsolute,
	't': csiWindowManipulation,
	'X': csiEraseCharactersHandler,
	'r': csiSetMarginsHandler,
}

// DECSTBM
func csiSetMarginsHandler(params []string, intermediate string, terminal *Terminal) error {
	top := 1
	bottom := int(terminal.ActiveBuffer().ViewHeight())
	if len(params) > 0 {
		var err error
		top, err = strconv.Atoi(params[0])
		if err != nil {
			top = 1
		}

		if len(params) > 1 {
			var err error
			bottom, err = strconv.Atoi(params[1])
			if err != nil {
				bottom = 1
			}
			if bottom > int(terminal.ActiveBuffer().ViewHeight()) {
				bottom = int(terminal.ActiveBuffer().ViewHeight())
			}
		}
	}
	top -= 1
	bottom -= 1

	terminal.ActiveBuffer().SetPosition(0, 0)

	return nil
}

func csiSetMode(modeStr string, enabled bool, terminal *Terminal) error {

	/*
	   Mouse support

	   		#define SET_X10_MOUSE               9
	        #define SET_VT200_MOUSE             1000
	        #define SET_VT200_HIGHLIGHT_MOUSE   1001
	        #define SET_BTN_EVENT_MOUSE         1002
	        #define SET_ANY_EVENT_MOUSE         1003

	        #define SET_FOCUS_EVENT_MOUSE       1004

	        #define SET_EXT_MODE_MOUSE          1005
	        #define SET_SGR_EXT_MODE_MOUSE      1006
	        #define SET_URXVT_EXT_MODE_MOUSE    1015

	        #define SET_ALTERNATE_SCROLL        1007
	*/

	switch modeStr {
	case "4":
		if enabled { // @todo support replace mode
			terminal.ActiveBuffer().SetInsertMode()
		} else {
			terminal.ActiveBuffer().SetReplaceMode()
		}
	case "?1":
		terminal.modes.ApplicationCursorKeys = enabled
	case "?7":
		// auto-wrap mode
		//DECAWM
		terminal.ActiveBuffer().SetAutoWrap(enabled)
	case "?9":
		if enabled {
			terminal.SetMouseMode(MouseModeX10)
		} else {
			terminal.SetMouseMode(MouseModeNone)
		}
	case "?12", "?13":
		terminal.modes.BlinkingCursor = enabled
	case "?25":
		terminal.modes.ShowCursor = enabled
	case "?47", "?1047":
		if enabled {
			terminal.UseAltBuffer()
		} else {
			terminal.UseMainBuffer()
		}
	case "?1000", "?1006;1000", "?10061000": // ?10061000 seen from htop
		// enable mouse tracking
		if enabled {
			terminal.SetMouseMode(MouseModeVT200)
		} else {
			terminal.SetMouseMode(MouseModeNone)
		}
	case "?1048":
		if enabled {
			terminal.ActiveBuffer().SaveCursor()
		} else {
			terminal.ActiveBuffer().RestoreCursor()
		}
	case "?1049":
		if enabled {
			terminal.UseAltBuffer()
		} else {
			terminal.UseMainBuffer()
		}
	default:
		return fmt.Errorf("Unsupported CSI %sl code", modeStr)
	}

	return nil
}

func csiEraseCharactersHandler(params []string, intermediate string, terminal *Terminal) error {
	count := 1
	if len(params) > 0 {
		var err error
		count, err = strconv.Atoi(params[0])
		if err != nil {
			count = 1
		}
	}

	terminal.ActiveBuffer().EraseCharacters(count)

	return nil
}

func csiResetModeHandler(params []string, intermediate string, terminal *Terminal) error {
	return csiSetMode(strings.Join(params, ""), false, terminal)
}

func csiSetModeHandler(params []string, intermediate string, terminal *Terminal) error {
	return csiSetMode(strings.Join(params, ""), true, terminal)
}

func csiWindowManipulation(params []string, intermediate string, terminal *Terminal) error {
	// @todo this
	return nil
}

func csiLinePositionAbsolute(params []string, intermediate string, terminal *Terminal) error {
	row := 1
	if len(params) > 0 {
		var err error
		row, err = strconv.Atoi(params[0])
		if err != nil {
			row = 1
		}
	}

	terminal.ActiveBuffer().SetPosition(terminal.ActiveBuffer().CursorColumn(), uint16(row-1))

	return nil
}

type csiSequenceHandler func(params []string, intermediate string, terminal *Terminal) error

// CSI: Control Sequence Introducer [
func csiHandler(pty chan rune, terminal *Terminal) error {
	var final rune
	var b rune
	var err error
	param := ""
	intermediate := ""
CSI:
	for {
		b = <-pty
		switch true {
		case b >= 0x30 && b <= 0x3F:
			param = param + string(b)
		case b >= 0x20 && b <= 0x2F:
			//intermediate? useful?
			intermediate += string(b)
		case b >= 0x40 && b <= 0x7e:
			final = b
			break CSI
		}
	}

	params := strings.Split(param, ";")

	handler, ok := csiSequenceMap[final]
	if ok {
		err = handler(params, intermediate, terminal)
	} else {

		switch final {
		case 'A':
			distance := 1
			if len(params) > 0 {
				var err error
				distance, err = strconv.Atoi(params[0])
				if err != nil {
					distance = 1
				}
			}
			terminal.ActiveBuffer().MovePosition(0, -int16(distance))
		case 'B':
			distance := 1
			if len(params) > 0 {
				var err error
				distance, err = strconv.Atoi(params[0])
				if err != nil {
					distance = 1
				}
			}

			terminal.ActiveBuffer().MovePosition(0, int16(distance))
		case 'C':

			distance := 1
			if len(params) > 0 {
				var err error
				distance, err = strconv.Atoi(params[0])
				if err != nil {
					distance = 1
				}
			}

			terminal.ActiveBuffer().MovePosition(int16(distance), 0)

		case 'D':

			distance := 1
			if len(params) > 0 {
				var err error
				distance, err = strconv.Atoi(params[0])
				if err != nil {
					distance = 1
				}
			}

			terminal.ActiveBuffer().MovePosition(-int16(distance), 0)

		case 'E':
			distance := 1
			if len(params) > 0 {
				var err error
				distance, err = strconv.Atoi(params[0])
				if err != nil {
					distance = 1
				}
			}

			terminal.ActiveBuffer().MovePosition(0, int16(distance))
			terminal.ActiveBuffer().SetPosition(0, terminal.ActiveBuffer().CursorLine())

		case 'F':

			distance := 1
			if len(params) > 0 {
				var err error
				distance, err = strconv.Atoi(params[0])
				if err != nil {
					distance = 1
				}
			}
			terminal.ActiveBuffer().MovePosition(0, -int16(distance))
			terminal.ActiveBuffer().SetPosition(0, terminal.ActiveBuffer().CursorLine())

		case 'G':

			distance := 1
			if len(params) > 0 {
				var err error
				distance, err = strconv.Atoi(params[0])
				if err != nil || params[0] == "" {
					distance = 1
				}
			}

			terminal.ActiveBuffer().SetPosition(uint16(distance-1), terminal.ActiveBuffer().CursorLine())

		case 'H', 'f':

			x, y := 1, 1
			if len(params) == 2 {
				var err error
				if params[0] != "" {
					y, err = strconv.Atoi(string(params[0]))
					if err != nil {
						y = 1
					}
				}
				if params[1] != "" {
					x, err = strconv.Atoi(string(params[1]))
					if err != nil {
						x = 1
					}
				}
			}

			terminal.ActiveBuffer().SetPosition(uint16(x-1), uint16(y-1))

		default:
			err = fmt.Errorf("Unknown CSI control sequence: 0x%02X (ESC[%s%s%s)", final, param, intermediate, string(final))
		}
	}
	terminal.logger.Debugf("Received CSI control sequence: 0x%02X (ESC[%s%s%s)", final, param, intermediate, string(final))
	return err
}

func csiDeleteHandler(params []string, intermediate string, terminal *Terminal) error {
	n := 1
	if len(params) >= 1 {
		var err error
		n, err = strconv.Atoi(params[0])
		if err != nil {
			n = 1
		}
	}

	terminal.ActiveBuffer().EraseCharacters(n)

	return nil
}

// CSI Ps J
func csiEraseInDisplayHandler(params []string, intermediate string, terminal *Terminal) error {
	n := "0"
	if len(params) > 0 {
		n = params[0]
	}

	switch n {

	case "0", "":
		terminal.ActiveBuffer().EraseDisplayFromCursor()
	case "1":
		terminal.ActiveBuffer().EraseDisplayToCursor()
	case "2":
		terminal.ActiveBuffer().EraseDisplay()
	default:
		return fmt.Errorf("Unsupported ED: CSI %s J", n)
	}

	return nil
}

// CSI Ps K
func csiEraseInLineHandler(params []string, intermediate string, terminal *Terminal) error {

	n := "0"
	if len(params) > 0 {
		n = params[0]
	}

	switch n {
	case "0", "": //erase adter cursor
		terminal.ActiveBuffer().EraseLineFromCursor()
	case "1": // erase to cursor inclusive
		terminal.ActiveBuffer().EraseLineToCursor()
	case "2": // erase entire
		terminal.ActiveBuffer().EraseLine()
	default:
		return fmt.Errorf("Unsupported EL: CSI %s K", n)
	}
	return nil
}
