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
}

func csiSetMode(modeStr string, enabled bool, terminal *Terminal) error {
	switch modeStr {
	case "?1":
		terminal.modes.ApplicationCursorKeys = enabled
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
	col := 1
	if len(params) > 0 {
		var err error
		col, err = strconv.Atoi(params[0])
		if err != nil {
			col = 1
		}
	}

	terminal.ActiveBuffer().SetPosition(uint16(col), terminal.ActiveBuffer().CursorLine())

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
