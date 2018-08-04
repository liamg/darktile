package terminal

var ansiSequenceMap = map[rune]escapeSequenceHandler{
	'[': csiHandler,
}

func ansiHandler(buffer chan rune, terminal *Terminal) error {
	// if the byte is an escape character, read the next byte to determine which one
	b := <-buffer

	handler, ok := ansiSequenceMap[b]
	if ok {
		return handler(buffer, terminal)
	}

	switch b {
	case 0x5d: // OSC: Operating System Command
		b = <-buffer
		switch b {
		case rune('0'):
			b = <-buffer
			if b == rune(';') {
				title := []rune{}
				for {
					b = <-buffer
					if b == 0x07 || b == 0x5c { // 0x07 -> BELL, 0x5c -> ST (\)
						break
					}
					title = append(title, b)
				}
				terminal.title = string(title)
			} else {
				terminal.logger.Errorf("Invalid OSC 0 control sequence: 0x%02X", b)
			}
		default:
			terminal.logger.Errorf("Unknown OSC control sequence: 0x%02X", b)
		}
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
