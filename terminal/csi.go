package terminal

import (
	"fmt"
	"strconv"
	"strings"
)

var csiSequenceMap = map[rune]csiSequenceHandler{
	'm': sgrSequenceHandler,
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
		if terminal.position.Line-distance >= 0 {
			terminal.position.Line -= distance
		} else {
			terminal.position.Line = 0
		}
	case 'B':
		distance := 1
		if len(params) > 0 {
			var err error
			distance, err = strconv.Atoi(params[0])
			if err != nil {
				distance = 1
			}
		}

		_, h := terminal.GetSize()
		if terminal.position.Line+distance >= h {
			terminal.position.Line = h - 1
		} else {
			terminal.position.Line += distance
		}
	case 'C':

		distance := 1
		if len(params) > 0 {
			var err error
			distance, err = strconv.Atoi(params[0])
			if err != nil {
				distance = 1
			}
		}

		terminal.position.Col += distance
		w, _ := terminal.GetSize()
		if terminal.position.Col >= w {
			terminal.position.Col = w - 1
		}

	case 'D':

		distance := 1
		if len(params) > 0 {
			var err error
			distance, err = strconv.Atoi(params[0])
			if err != nil {
				distance = 1
			}
		}

		terminal.position.Col -= distance
		if terminal.position.Col < 0 {
			terminal.position.Col = 0
		}

	case 'E':
		distance := 1
		if len(params) > 0 {
			var err error
			distance, err = strconv.Atoi(params[0])
			if err != nil {
				distance = 1
			}
		}

		terminal.position.Line += distance
		terminal.position.Col = 0

	case 'F':

		distance := 1
		if len(params) > 0 {
			var err error
			distance, err = strconv.Atoi(params[0])
			if err != nil {
				distance = 1
			}
		}
		if terminal.position.Line-distance >= 0 {
			terminal.position.Line -= distance
		} else {
			terminal.position.Line = 0
		}
		terminal.position.Col = 0

	case 'G':

		distance := 1
		if len(params) > 0 {
			var err error
			distance, err = strconv.Atoi(params[0])
			if err != nil || params[0] == "" {
				distance = 1
			}
		}

		terminal.position.Col = distance - 1 // 1 based to 0 based

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
		terminal.position.Col = x - 1
		terminal.position.Line = y - 1

	case 'J':

		n := "0"
		if len(params) > 0 {
			n = params[0]
		}

		switch n {

		case "0", "":
			line := terminal.getBufferedLine(terminal.position.Line)
			if line != nil {
				line.Cells = line.Cells[:terminal.position.Col]
			}
			_, h := terminal.GetSize()
			for i := terminal.position.Line + 1; i < h; i++ {
				line := terminal.getBufferedLine(i)
				if line != nil {
					line.Cells = []Cell{}
				}
			}
		case "1":
			line := terminal.getBufferedLine(terminal.position.Line)
			if line != nil {
				for i := 0; i <= terminal.position.Col; i++ {
					if i < len(line.Cells) {
						line.Cells[i].r = 0
					}
				}
			}
			for i := 0; i < terminal.position.Line; i++ {
				line := terminal.getBufferedLine(i)
				if line != nil {
					line.Cells = []Cell{}
				}
			}

		case "2":
			_, h := terminal.GetSize()
			for i := 0; i < h; i++ {
				line := terminal.getBufferedLine(i)
				if line != nil {
					line.Cells = []Cell{}
				}
			}
		case "3":
			terminal.lines = []Line{}

		default:
			return fmt.Errorf("Unknown CSI ED sequence: %s", n)
		}

	case 'K': // K - EOL - Erase to end of line
		n := "0"
		if len(params) > 0 {
			n = params[0]
		}

		switch n {
		case "0", "":
			line := terminal.getBufferedLine(terminal.position.Line)
			if line != nil {
				line.Cells = line.Cells[:terminal.position.Col]
			}
		case "1":
			line := terminal.getBufferedLine(terminal.position.Line)
			if line != nil {
				for i := 0; i <= terminal.position.Col; i++ {
					if i < len(line.Cells) {
						line.Cells[i].r = 0
					}
				}
			}
		case "2":
			line := terminal.getBufferedLine(terminal.position.Line)
			if line != nil {
				line.Cells = []Cell{}
			}
		default:
			return fmt.Errorf("Unsupported EL: %s", n)
		}

	case 'P': // delete

		n := 1
		if len(params) >= 1 {
			var err error
			n, err = strconv.Atoi(params[0])
			if err != nil {
				n = 1
			}
		}

		_ = terminal.delete(n)

	default:
		switch param + intermediate + string(final) {
		case "?25h":
			terminal.showCursor()
		case "?25l":
			terminal.hideCursor()
		case "?12h":
			// todo enable cursor blink
		case "?12l":
			// todo disable cursor blink
		default:
			return fmt.Errorf("Unknown CSI control sequence: 0x%02X (ESC[%s%s%s)", final, param, intermediate, string(final))
		}

	}
	//terminal.logger.Debugf("Received CSI control sequence: 0x%02X (ESC[%s%s%s)", final, param, intermediate, string(final))
	return nil
}
