package terminal

import (
	"fmt"
	"strconv"
	"strings"
)

var csiSequenceMap = map[rune]csiSequenceHandler{
	'd': csiLinePositionAbsolute,
	'h': csiSetModeHandler,
	'l': csiResetModeHandler,
	'm': sgrSequenceHandler,
	'r': csiSetMarginsHandler,
	't': csiWindowManipulation,
	'J': csiEraseInDisplayHandler,
	'K': csiEraseInLineHandler,
	'L': csiInsertLinesHandler,
	'P': csiDeleteHandler,
	'S': csiScrollUpHandler,
	'T': csiScrollDownHandler,
	'X': csiEraseCharactersHandler,
}

func csiScrollUpHandler(params []string, intermediate string, terminal *Terminal) error {
	distance := 1
	if len(params) > 1 {
		return fmt.Errorf("Not supported")
	}
	if len(params) == 1 {
		var err error
		distance, err = strconv.Atoi(params[0])
		if err != nil {
			distance = 1
		}
	}
	terminal.logger.Debugf("Scrolling up %d", distance)
	terminal.ScrollUp(uint16(distance))
	return nil
}

func csiInsertLinesHandler(params []string, intermediate string, terminal *Terminal) error {
	count := 1
	if len(params) > 1 {
		return fmt.Errorf("Not supported")
	}
	if len(params) == 1 {
		var err error
		count, err = strconv.Atoi(params[0])
		if err != nil {
			count = 1
		}
	}
	terminal.logger.Debugf("Inserting %d lines", count)
	return fmt.Errorf("Not supported")
}

func csiScrollDownHandler(params []string, intermediate string, terminal *Terminal) error {
	distance := 1
	if len(params) > 1 {
		return fmt.Errorf("Not supported")
	}
	if len(params) == 1 {
		var err error
		distance, err = strconv.Atoi(params[0])
		if err != nil {
			distance = 1
		}
	}
	terminal.logger.Debugf("Scrolling down %d", distance)
	terminal.ScrollDown(uint16(distance))
	return nil
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
	top--
	bottom--

	terminal.ActiveBuffer().SetVerticalMargins(uint(top), uint(bottom))
	terminal.ActiveBuffer().SetPosition(0, 0)

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
	fmt.Printf("CSI 0x%02X (ESC[%s%s%s)\n", final, param, intermediate, string(final))
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
