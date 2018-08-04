package terminal

var ansiSequenceMap = map[rune]escapeSequenceHandler{
	'[':  csiHandler,
	0x5d: oscHandler,
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
