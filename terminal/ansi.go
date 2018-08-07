package terminal

// https://www.xfree86.org/4.8.0/ctlseqs.html
// https://vt100.net/docs/vt100-ug/chapter3.html

var ansiSequenceMap = map[rune]escapeSequenceHandler{
	'[':  csiHandler,
	0x5d: oscHandler,
	'7':  saveCursorHandler,
	'8':  restoreCursorHandler,
	'D':  indexHandler,
	'M':  reverseIndexHandler,
}

func indexHandler(buffer chan rune, terminal *Terminal) error {
	// @todo is thus right?
	// "This sequence causes the active position to move downward one line without changing the column position. If the active position is at the bottom margin, a scroll up is performed."
	if terminal.buffer.CursorLine() == terminal.buffer.ViewHeight()-1 {
		terminal.buffer.NewLine()
		return nil
	}
	terminal.buffer.MovePosition(0, 1)
	return nil
}

func reverseIndexHandler(buffer chan rune, terminal *Terminal) error {
	terminal.buffer.MovePosition(0, -1)
	return nil
}

func saveCursorHandler(buffer chan rune, terminal *Terminal) error {
	terminal.buffer.SaveCursor()
	return nil
}

func restoreCursorHandler(buffer chan rune, terminal *Terminal) error {
	terminal.buffer.RestoreCursor()
	return nil
}

func ansiHandler(buffer chan rune, terminal *Terminal) error {
	// if the byte is an escape character, read the next byte to determine which one
	b := <-buffer

	handler, ok := ansiSequenceMap[b]
	if ok {
		return handler(buffer, terminal)
	}

	switch b {

	case 'c':
		terminal.logger.Errorf("RIS not yet supported")
	case '(':
		b = <-buffer
		switch b {
		case 'A': //uk @todo handle these?
			//terminal.charSet = C0
		case 'B': //us
			//terminal.charSet = C0
		}
	case ')':
		b = <-buffer
		switch b {
		case 'A': //uk @todo handle these?
			//terminal.charSet = C1
		case 'B': //us
			//terminal.charSet = C1
		}
	case '*':
		b = <-buffer
		switch b {
		case 'A': //uk @todo handle these?
			//terminal.charSet = C2
		case 'B': //us
			//terminal.charSet = C2
		}
	case '+':
		b = <-buffer
		switch b {
		case 'A': //uk @todo handle these?
			//terminal.charSet = C3
		case 'B': //us
			//terminal.charSet = C3
		}
	case '>':
		// numeric char selection @todo
	case '=':
		//alternate char selection @todo
	case '?':
		pm := ""
		for {
			b = <-buffer
			switch b {
			case 'h':
				switch pm {
				default:
					terminal.logger.Errorf("Unknown private code ESC?%sh", pm)
				}
			case 'l':
				switch pm {
				default:
					terminal.logger.Errorf("Unknown private code ESC?%sl", pm)
				}
			default:
				pm += string(b)
			}
		}
	default:
		terminal.logger.Errorf("Unknown control sequence: 0x%02X [%s]", b, string(b))
	}
	return nil
}
