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
}

type csiSequenceHandler func(params []string, intermediate string, terminal *Terminal) error

// CSI: Control Sequence Introducer [
func csiHandler(buffer chan rune, terminal *Terminal) error {
	var final rune
	var b rune
	param := ""
	intermediate := ""
CSI:
	for {
		b = <-buffer
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
		return handler(params, intermediate, terminal)
	}

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
		terminal.buffer.MovePosition(0, -int16(distance))
	case 'B':
		distance := 1
		if len(params) > 0 {
			var err error
			distance, err = strconv.Atoi(params[0])
			if err != nil {
				distance = 1
			}
		}

		terminal.buffer.MovePosition(0, int16(distance))
	case 'C':

		distance := 1
		if len(params) > 0 {
			var err error
			distance, err = strconv.Atoi(params[0])
			if err != nil {
				distance = 1
			}
		}

		terminal.buffer.MovePosition(int16(distance), 0)

	case 'D':

		distance := 1
		if len(params) > 0 {
			var err error
			distance, err = strconv.Atoi(params[0])
			if err != nil {
				distance = 1
			}
		}

		terminal.buffer.MovePosition(-int16(distance), 0)

	case 'E':
		distance := 1
		if len(params) > 0 {
			var err error
			distance, err = strconv.Atoi(params[0])
			if err != nil {
				distance = 1
			}
		}

		terminal.buffer.MovePosition(0, int16(distance))
		terminal.buffer.SetPosition(0, terminal.buffer.CursorLine())

	case 'F':

		distance := 1
		if len(params) > 0 {
			var err error
			distance, err = strconv.Atoi(params[0])
			if err != nil {
				distance = 1
			}
		}
		terminal.buffer.MovePosition(0, -int16(distance))
		terminal.buffer.SetPosition(0, terminal.buffer.CursorLine())

	case 'G':

		distance := 1
		if len(params) > 0 {
			var err error
			distance, err = strconv.Atoi(params[0])
			if err != nil || params[0] == "" {
				distance = 1
			}
		}

		terminal.buffer.SetPosition(uint16(distance-1), terminal.buffer.CursorLine())

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

		terminal.buffer.SetPosition(uint16(x-1), uint16(y-1))

	default:
		switch param + intermediate + string(final) {
		case "?25h":
			terminal.buffer.ShowCursor()
		case "?25l":
			terminal.buffer.HideCursor()
		case "?12h":
			terminal.buffer.SetCursorBlink(true)
		case "?12l":
			terminal.buffer.SetCursorBlink(false)
		default:
			return fmt.Errorf("Unknown CSI control sequence: 0x%02X (ESC[%s%s%s)", final, param, intermediate, string(final))
		}

	}
	//terminal.logger.Debugf("Received CSI control sequence: 0x%02X (ESC[%s%s%s)", final, param, intermediate, string(final))
	return nil
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
	_ = n
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
		terminal.buffer.EraseDisplayAfterCursor()
	case "1":
		terminal.buffer.EraseDisplayToCursor()
	case "2":
		terminal.Clear()
	case "3":
		terminal.Clear()

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
		terminal.buffer.EraseLineAfterCursor()
	case "1": // erase to cursor inclusive
		terminal.buffer.EraseLineToCursor()
	case "2": // erase entire
		terminal.buffer.EraseLine()
	default:
		return fmt.Errorf("Unsupported EL: CSI %s K", n)
	}
	return nil
}
