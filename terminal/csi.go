package terminal

import (
	"fmt"
	"strconv"
	"strings"
)

type csiSequenceHandler func(params []string, terminal *Terminal) error

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
	{id: 'c', handler: csiSendDeviceAttributesHandler, description: " Send Device Attributes (Primary/Secondary/Tertiary DA)"},
	{id: 'd', handler: csiLinePositionAbsolute, expectedParams: &expectedParams{min: 0, max: 1}, description: "Line Position Absolute  [row] (default = [1,column]) (VPA)"},
	{id: 'f', handler: csiCursorPositionHandler, description: "Horizontal and Vertical Position [row;column] (default = [1,1]) (HVP)"},
	{id: 'g', handler: csiTabClearHandler, description: "Tab Clear (TBC)"},
	{id: 'h', handler: csiSetModeHandler, expectedParams: &expectedParams{min: 1, max: ^uint8(0)}, description: "Set Mode (SM)"},
	{id: 'l', handler: csiResetModeHandler, expectedParams: &expectedParams{min: 1, max: ^uint8(0)}, description: "Reset Mode (RM)"},
	{id: 'm', handler: sgrSequenceHandler, description: "Character Attributes (SGR)"},
	{id: 'n', handler: csiDeviceStatusReportHandler, description: "Device Status Report (DSR)"},
	{id: 'r', handler: csiSetMarginsHandler, expectedParams: &expectedParams{min: 0, max: 2}, description: "Set Scrolling Region [top;bottom] (default = full size of window) (DECSTBM), VT100"},
	{id: 't', handler: csiWindowManipulation, description: "Window manipulation"},
	{id: 'A', handler: csiCursorUpHandler, description: "Cursor Up Ps Times (default = 1) (CUU)"},
	{id: 'B', handler: csiCursorDownHandler, description: "Cursor Down Ps Times (default = 1) (CUD)"},
	{id: 'C', handler: csiCursorForwardHandler, description: "Cursor Forward Ps Times (default = 1) (CUF)"},
	{id: 'D', handler: csiCursorBackwardHandler, description: "Cursor Backward Ps Times (default = 1) (CUB)"},
	{id: 'E', handler: csiCursorNextLineHandler, description: "Cursor Next Line Ps Times (default = 1) (CNL)"},
	{id: 'F', handler: csiCursorPrecedingLineHandler, description: "Cursor Preceding Line Ps Times (default = 1) (CPL)"},
	{id: 'G', handler: csiCursorCharacterAbsoluteHandler, description: "Cursor Horizontal Absolute  [column] (default = [row,1]) (CHA)"},
	{id: 'H', handler: csiCursorPositionHandler, description: "Cursor Position [row;column] (default = [1,1]) (CUP)"},
	{id: 'J', handler: csiEraseInDisplayHandler, description: "Erase in Display (ED), VT100"},
	{id: 'K', handler: csiEraseInLineHandler, description: "Erase in Line (EL), VT100"},
	{id: 'L', handler: csiInsertLinesHandler, description: "Insert Ps Line(s) (default = 1) (IL)"},
	{id: 'M', handler: csiDeleteLinesHandler, description: "Delete Ps Line(s) (default = 1) (DL)"},
	{id: 'P', handler: csiDeleteHandler, description: " Delete Ps Character(s) (default = 1) (DCH)"},
	{id: 'S', handler: csiScrollUpHandler, description: "Scroll up Ps lines (default = 1) (SU), VT420, ECMA-48"},
	{id: 'T', handler: csiScrollDownHandler, description: "Scroll down Ps lines (default = 1) (SD), VT420"},
	{id: 'X', handler: csiEraseCharactersHandler, description: "Erase Ps Character(s) (default = 1) (ECH"},
	{id: '@', handler: csiInsertBlankCharactersHandler, description: "Insert Ps (Blank) Character(s) (default = 1) (ICH)"},
}

type runeRange struct {
	min rune
	max rune
}

var csiTerminators = runeRange{0x40, 0x7e}

func loadCSI(pty chan rune) (final rune, param string, intermediate []rune) {
	var b rune
	param = ""
	intermediate = []rune{}
CSI:
	for {
		b = <-pty
		switch true {
		case b >= 0x30 && b <= 0x3F:
			param = param + string(b)
		case b > 0 && b <= 0x2F:
			intermediate = append(intermediate, b)
		case b >= csiTerminators.min && b <= csiTerminators.max:
			final = b
			break CSI
		}
	}

	return final, param, intermediate
}

func splitParams(paramString string) []string {
	params := strings.Split(paramString, ";")
	if paramString == "" {
		params = []string{}
	}
	return params
}

func csiHandler(pty chan rune, terminal *Terminal) error {
	final, param, intermediate := loadCSI(pty)

	// process intermediate control codes before the CSI
	for _, b := range intermediate {
		terminal.processRune(b)
	}

	params := splitParams(param)

	for _, sequence := range csiSequences {
		if sequence.id == final {
			if sequence.expectedParams != nil && (uint8(len(params)) < sequence.expectedParams.min || uint8(len(params)) > sequence.expectedParams.max) {
				continue
			}
			x, y := terminal.ActiveBuffer().CursorColumn(), terminal.ActiveBuffer().CursorLine()
			err := sequence.handler(params, terminal)
			terminal.logger.Debugf("CSI 0x%02X (ESC[%s%s) %s - %d,%d -> %d,%d", final, param, string(final), sequence.description, x, y, terminal.ActiveBuffer().CursorColumn(), terminal.ActiveBuffer().CursorLine())
			return err
		}
	}

	return fmt.Errorf("Unknown CSI control sequence: 0x%02X (ESC[%s%s)", final, param, string(final))
}

func csiSendDeviceAttributesHandler(params []string, terminal *Terminal) error {

	// we are VT100
	// for DA1 we'll respond ?1;2
	// for DA2 we'll respond >0;0;0

	response := "?1;2"

	if len(params) > 0 && len(params[0]) > 0 && params[0][0] == '>' {
		response = ">0;0;0"
	}

	_ = terminal.Write([]byte("\x1b[" + response + "c"))

	return nil
}

func csiDeviceStatusReportHandler(params []string, terminal *Terminal) error {

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

func csiCursorUpHandler(params []string, terminal *Terminal) error {
	distance := 1
	if len(params) > 0 {
		var err error
		distance, err = strconv.Atoi(params[0])
		if err != nil || distance < 1 {
			distance = 1
		}
	}
	terminal.ActiveBuffer().MovePosition(0, -int16(distance))
	return nil
}

func csiCursorDownHandler(params []string, terminal *Terminal) error {
	distance := 1
	if len(params) > 0 {
		var err error
		distance, err = strconv.Atoi(params[0])
		if err != nil || distance < 1 {
			distance = 1
		}
	}

	terminal.ActiveBuffer().MovePosition(0, int16(distance))
	return nil
}

func csiCursorForwardHandler(params []string, terminal *Terminal) error {
	distance := 1
	if len(params) > 0 {
		var err error
		distance, err = strconv.Atoi(params[0])
		if err != nil || distance < 1 {
			distance = 1
		}
	}

	terminal.ActiveBuffer().MovePosition(int16(distance), 0)
	return nil
}

func csiCursorBackwardHandler(params []string, terminal *Terminal) error {
	distance := 1
	if len(params) > 0 {
		var err error
		distance, err = strconv.Atoi(params[0])
		if err != nil || distance < 1 {
			distance = 1
		}
	}

	terminal.ActiveBuffer().MovePosition(-int16(distance), 0)
	return nil
}

func csiCursorNextLineHandler(params []string, terminal *Terminal) error {

	distance := 1
	if len(params) > 0 {
		var err error
		distance, err = strconv.Atoi(params[0])
		if err != nil || distance < 1 {
			distance = 1
		}
	}

	terminal.ActiveBuffer().MovePosition(0, int16(distance))
	terminal.ActiveBuffer().SetPosition(0, terminal.ActiveBuffer().CursorLine())
	return nil
}

func csiCursorPrecedingLineHandler(params []string, terminal *Terminal) error {

	distance := 1
	if len(params) > 0 {
		var err error
		distance, err = strconv.Atoi(params[0])
		if err != nil || distance < 1 {
			distance = 1
		}
	}
	terminal.ActiveBuffer().MovePosition(0, -int16(distance))
	terminal.ActiveBuffer().SetPosition(0, terminal.ActiveBuffer().CursorLine())
	return nil
}

func csiCursorCharacterAbsoluteHandler(params []string, terminal *Terminal) error {
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

func parseCursorPosition(params []string) (x, y int) {
	x, y = 1, 1
	if len(params) == 2 {
		var err error
		if params[0] != "" {
			y, err = strconv.Atoi(string(params[0]))
			if err != nil || y < 1 {
				y = 1
			}
		}
		if params[1] != "" {
			x, err = strconv.Atoi(string(params[1]))
			if err != nil || x < 1 {
				x = 1
			}
		}
	}
	return x, y
}

func csiCursorPositionHandler(params []string, terminal *Terminal) error {
	x, y := parseCursorPosition(params)

	terminal.ActiveBuffer().SetPosition(uint16(x-1), uint16(y-1))
	return nil
}

func csiScrollUpHandler(params []string, terminal *Terminal) error {
	distance := 1
	if len(params) > 1 {
		return fmt.Errorf("Not supported")
	}
	if len(params) == 1 {
		var err error
		distance, err = strconv.Atoi(params[0])
		if err != nil || distance < 1 {
			distance = 1
		}
	}
	terminal.logger.Debugf("Scrolling up %d", distance)
	terminal.AreaScrollUp(uint16(distance))
	return nil
}

func csiInsertBlankCharactersHandler(params []string, terminal *Terminal) error {
	count := 1
	if len(params) > 1 {
		return fmt.Errorf("Not supported")
	}
	if len(params) == 1 {
		var err error
		count, err = strconv.Atoi(params[0])
		if err != nil || count < 1 {
			count = 1
		}
	}

	terminal.ActiveBuffer().InsertBlankCharacters(count)

	return nil
}

func csiInsertLinesHandler(params []string, terminal *Terminal) error {
	count := 1
	if len(params) > 1 {
		return fmt.Errorf("Not supported")
	}
	if len(params) == 1 {
		var err error
		count, err = strconv.Atoi(params[0])
		if err != nil || count < 1 {
			count = 1
		}
	}

	terminal.ActiveBuffer().InsertLines(count)

	return nil
}

func csiDeleteLinesHandler(params []string, terminal *Terminal) error {
	count := 1
	if len(params) > 1 {
		return fmt.Errorf("Not supported")
	}
	if len(params) == 1 {
		var err error
		count, err = strconv.Atoi(params[0])
		if err != nil || count < 1 {
			count = 1
		}
	}

	terminal.ActiveBuffer().DeleteLines(count)

	return nil
}

func csiScrollDownHandler(params []string, terminal *Terminal) error {
	distance := 1
	if len(params) > 1 {
		return fmt.Errorf("Not supported")
	}
	if len(params) == 1 {
		var err error
		distance, err = strconv.Atoi(params[0])
		if err != nil || distance < 1 {
			distance = 1
		}
	}
	terminal.logger.Debugf("Scrolling down %d", distance)
	terminal.AreaScrollDown(uint16(distance))
	return nil
}

// DECSTBM
func csiSetMarginsHandler(params []string, terminal *Terminal) error {
	top := 1
	bottom := int(terminal.ActiveBuffer().ViewHeight())

	if len(params) > 2 {
		return fmt.Errorf("Not set margins")
	}

	if len(params) > 0 {
		var err error
		top, err = strconv.Atoi(params[0])
		if err != nil || top < 1 {
			top = 1
		}

		if len(params) > 1 {
			var err error
			bottom, err = strconv.Atoi(params[1])
			if err != nil || bottom > int(terminal.ActiveBuffer().ViewHeight()) || bottom < 1 {
				bottom = int(terminal.ActiveBuffer().ViewHeight())
			}
		}
	}
	top--
	bottom--

	terminal.terminalState.SetVerticalMargins(uint(top), uint(bottom))
	terminal.ActiveBuffer().SetPosition(0, 0)

	return nil
}

func csiEraseCharactersHandler(params []string, terminal *Terminal) error {
	count := 1
	if len(params) > 0 {
		var err error
		count, err = strconv.Atoi(params[0])
		if err != nil || count < 1 {
			count = 1
		}
	}

	terminal.ActiveBuffer().EraseCharacters(count)

	return nil
}

func csiResetModeHandler(params []string, terminal *Terminal) error {
	terminal.ActiveBuffer().ClearSelection()
	return csiSetModes(params, false, terminal)
}

func csiSetModeHandler(params []string, terminal *Terminal) error {
	return csiSetModes(params, true, terminal)
}

func csiWindowManipulation(params []string, terminal *Terminal) error {
	return fmt.Errorf("Window manipulation is not yet supported")
}

func csiLinePositionAbsolute(params []string, terminal *Terminal) error {
	row := 1
	if len(params) > 0 {
		var err error
		row, err = strconv.Atoi(params[0])
		if err != nil || row < 1 {
			row = 1
		}
	}

	terminal.ActiveBuffer().SetPosition(terminal.ActiveBuffer().CursorColumn(), uint16(row-1))

	return nil
}

func csiDeleteHandler(params []string, terminal *Terminal) error {
	n := 1
	if len(params) >= 1 {
		var err error
		n, err = strconv.Atoi(params[0])
		if err != nil || n < 1 {
			n = 1
		}
	}

	terminal.ActiveBuffer().DeleteChars(n)
	return nil
}

func csiTabClearHandler(params []string, terminal *Terminal) error {
	n := "0"
	if len(params) > 0 {
		n = params[0]
	}
	switch n {

	case "0", "":
		terminal.terminalState.TabClearAtCursor()
	case "3":
		terminal.terminalState.TabZonk()
	default:
		return fmt.Errorf("Ignored TBC: CSI %s g", n)
	}

	return nil
}

// CSI Ps J
func csiEraseInDisplayHandler(params []string, terminal *Terminal) error {
	n := "0"
	if len(params) > 0 {
		n = params[0]
	}

	switch n {

	case "0", "":
		terminal.ActiveBuffer().EraseDisplayFromCursor()
	case "1":
		terminal.ActiveBuffer().EraseDisplayToCursor()
	case "2", "3":
		terminal.ActiveBuffer().EraseDisplay()
	default:
		return fmt.Errorf("Unsupported ED: CSI %s J", n)
	}

	return nil
}

// CSI Ps K
func csiEraseInLineHandler(params []string, terminal *Terminal) error {

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
