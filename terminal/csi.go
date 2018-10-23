package terminal

import (
	"fmt"
	"strconv"
	"strings"
)

type csiSequenceHandler func(params []string, intermediate string, terminal *Terminal) error

type csiMapping struct {
	id             rune
	handler        csiSequenceHandler
	description    string
	expectedParams *expectedParams
}

type expectedParams struct {
	min uint8
	max uint8
}

var csiSequences = []csiMapping{
	csiMapping{id: 'c', handler: csiSendDeviceAttributesHandler, description: " Send Device Attributes (Primary/Secondary/Tertiary DA)"},
	csiMapping{id: 'd', handler: csiLinePositionAbsolute, expectedParams: &expectedParams{min: 0, max: 1}, description: "Line Position Absolute  [row] (default = [1,column]) (VPA)"},
	csiMapping{id: 'f', handler: csiCursorPositionHandler, description: "Horizontal and Vertical Position [row;column] (default = [1,1]) (HVP)"},
	csiMapping{id: 'h', handler: csiSetModeHandler, expectedParams: &expectedParams{min: 1, max: 1}, description: "Set Mode (SM)"},
	csiMapping{id: 'l', handler: csiResetModeHandler, expectedParams: &expectedParams{min: 1, max: 1}, description: "Reset Mode (RM)"},
	csiMapping{id: 'm', handler: sgrSequenceHandler, description: "Character Attributes (SGR)"},
	csiMapping{id: 'n', handler: csiDeviceStatusReportHandler, description: "Device Status Report (DSR)"},
	csiMapping{id: 'r', handler: csiSetMarginsHandler, expectedParams: &expectedParams{min: 2, max: 2}, description: "Set Scrolling Region [top;bottom] (default = full size of window) (DECSTBM), VT100"},
	csiMapping{id: 't', handler: csiWindowManipulation, description: "Window manipulation"},
	csiMapping{id: 'A', handler: csiCursorUpHandler, description: "Cursor Up Ps Times (default = 1) (CUU)"},
	csiMapping{id: 'B', handler: csiCursorDownHandler, description: "Cursor Down Ps Times (default = 1) (CUD)"},
	csiMapping{id: 'C', handler: csiCursorForwardHandler, description: "Cursor Forward Ps Times (default = 1) (CUF)"},
	csiMapping{id: 'D', handler: csiCursorBackwardHandler, description: "Cursor Backward Ps Times (default = 1) (CUB)"},
	csiMapping{id: 'E', handler: csiCursorNextLineHandler, description: "Cursor Next Line Ps Times (default = 1) (CNL)"},
	csiMapping{id: 'F', handler: csiCursorPrecedingLineHandler, description: "Cursor Preceding Line Ps Times (default = 1) (CPL)"},
	csiMapping{id: 'G', handler: csiCursorCharacterAbsoluteHandler, description: "Cursor Character Absolute  [column] (default = [row,1]) (CHA)"},
	csiMapping{id: 'H', handler: csiCursorPositionHandler, description: "Cursor Position [row;column] (default = [1,1]) (CUP)"},
	csiMapping{id: 'J', handler: csiEraseInDisplayHandler, description: "Erase in Display (ED), VT100"},
	csiMapping{id: 'K', handler: csiEraseInLineHandler, description: "Erase in Line (EL), VT100"},
	csiMapping{id: 'L', handler: csiInsertLinesHandler, description: "Insert Ps Line(s) (default = 1) (IL)"},
	csiMapping{id: 'P', handler: csiDeleteHandler, description: " Delete Ps Character(s) (default = 1) (DCH)"},
	csiMapping{id: 'S', handler: csiScrollUpHandler, description: "Scroll up Ps lines (default = 1) (SU), VT420, ECMA-48"},
	csiMapping{id: 'T', handler: csiScrollDownHandler, description: "Scroll down Ps lines (default = 1) (SD), VT420"},
	csiMapping{id: 'X', handler: csiEraseCharactersHandler, description: "Erase Ps Character(s) (default = 1) (ECH"},
}

func csiHandler(pty chan rune, terminal *Terminal) error {
	var final rune
	var b rune
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
	if param == "" {
		params = []string{}
	}

	for _, sequence := range csiSequences {
		if sequence.id == final {
			if sequence.expectedParams != nil && (uint8(len(params)) < sequence.expectedParams.min || uint8(len(params)) > sequence.expectedParams.max) {
				continue
			}
			terminal.logger.Debugf("CSI 0x%02X (ESC[%s%s%s) %s", final, param, intermediate, string(final), sequence.description)
			err := sequence.handler(params, intermediate, terminal)
			return err
		}
	}

	return fmt.Errorf("Unknown CSI control sequence: 0x%02X (ESC[%s%s%s)", final, param, intermediate, string(final))
}

func csiSendDeviceAttributesHandler(params []string, intermediate string, terminal *Terminal) error {

	if len(params) > 0 && len(params[0]) > 0 && params[0][0] == '>' { // secondary
		_ = terminal.Write([]byte("\x1b[0;0;0c")) // report VT100
		return nil
	}

	return fmt.Errorf("Unsupported SDA identifier")
}

func csiDeviceStatusReportHandler(params []string, intermediate string, terminal *Terminal) error {

	if len(params) == 0 {
		return fmt.Errorf("Missing Device Status Report identifier")
	}

	switch params[0] {
	case "5":
		_ = terminal.Write([]byte("\x1b[0n")) // everything is cool
	case "6": // report cursor position
		_ = terminal.Write([]byte(fmt.Sprintf(
			"\x1b[%d;%dR",
			terminal.ActiveBuffer().CursorLine()+1,
			terminal.ActiveBuffer().CursorColumn()+1,
		)))
	default:
		return fmt.Errorf("Unknown Device Status Report identifier: %s", params[0])
	}

	return nil
}

func csiCursorUpHandler(params []string, intermediate string, terminal *Terminal) error {
	distance := 1
	if len(params) > 0 {
		var err error
		distance, err = strconv.Atoi(params[0])
		if err != nil {
			distance = 1
		}
	}
	terminal.ActiveBuffer().MovePosition(0, -int16(distance))
	return nil
}

func csiCursorDownHandler(params []string, intermediate string, terminal *Terminal) error {
	distance := 1
	if len(params) > 0 {
		var err error
		distance, err = strconv.Atoi(params[0])
		if err != nil {
			distance = 1
		}
	}

	terminal.ActiveBuffer().MovePosition(0, int16(distance))
	return nil
}

func csiCursorForwardHandler(params []string, intermediate string, terminal *Terminal) error {
	distance := 1
	if len(params) > 0 {
		var err error
		distance, err = strconv.Atoi(params[0])
		if err != nil {
			distance = 1
		}
	}

	terminal.ActiveBuffer().MovePosition(int16(distance), 0)
	return nil
}

func csiCursorBackwardHandler(params []string, intermediate string, terminal *Terminal) error {
	distance := 1
	if len(params) > 0 {
		var err error
		distance, err = strconv.Atoi(params[0])
		if err != nil {
			distance = 1
		}
	}

	terminal.ActiveBuffer().MovePosition(-int16(distance), 0)
	return nil
}

func csiCursorNextLineHandler(params []string, intermediate string, terminal *Terminal) error {

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
	return nil
}

func csiCursorPrecedingLineHandler(params []string, intermediate string, terminal *Terminal) error {

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
	return nil
}

func csiCursorCharacterAbsoluteHandler(params []string, intermediate string, terminal *Terminal) error {
	distance := 1
	if len(params) > 0 {
		var err error
		distance, err = strconv.Atoi(params[0])
		if err != nil || params[0] == "" {
			distance = 1
		}
	}

	terminal.ActiveBuffer().SetPosition(uint16(distance-1), terminal.ActiveBuffer().CursorLine())
	return nil
}

func csiCursorPositionHandler(params []string, intermediate string, terminal *Terminal) error {
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
	return nil
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

	terminal.ActiveBuffer().InsertLines(count)

	return nil
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

	if len(params) > 2 {
		return fmt.Errorf("Not set margins")
	}

	if len(params) > 0 {
		var err error
		top, err = strconv.Atoi(params[0])
		if err != nil {
			top = 1
		}

		if len(params) > 1 {
			var err error
			bottom, err = strconv.Atoi(params[1])
			if err != nil || bottom > int(terminal.ActiveBuffer().ViewHeight()) {
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
	return fmt.Errorf("Window manipulation is not yet supported")
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
