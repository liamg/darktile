package terminal

import "fmt"

// https://www.xfree86.org/4.8.0/ctlseqs.html
// https://vt100.net/docs/vt100-ug/chapter3.html

var ansiSequenceMap = map[rune]escapeSequenceHandler{
	'[': csiHandler,
	']': oscHandler,
	'7': saveCursorHandler,
	'8': restoreCursorHandler,
	'D': indexHandler,
	'E': nextLineHandler, // NEL
	'H': tabSetHandler,   // HTS
	'M': reverseIndexHandler,
	'P': sixelHandler,
	'c': risHandler, //RIS
	'#': screenStateHandler,
	'^': privacyMessageHandler,
	'(': scs0Handler,       // select character set into G0
	')': scs1Handler,       // select character set into G1
	'*': swallowHandler(1), // character set bullshit
	'+': swallowHandler(1), // character set bullshit
	'>': swallowHandler(0), // numeric char selection  //@todo
	'=': swallowHandler(0), // alt char selection  //@todo
}

func swallowHandler(n int) func(pty chan rune, terminal *Terminal) error {
	return func(pty chan rune, terminal *Terminal) error {
		for i := 0; i < n; i++ {
			<-pty
		}
		return nil
	}
}

func risHandler(pty chan rune, terminal *Terminal) error {
	terminal.Lock()
	defer terminal.Unlock()

	terminal.ActiveBuffer().Clear()
	return nil
}

func indexHandler(pty chan rune, terminal *Terminal) error {
	terminal.Lock()
	defer terminal.Unlock()

	terminal.ActiveBuffer().Index()
	return nil
}

func reverseIndexHandler(pty chan rune, terminal *Terminal) error {
	terminal.Lock()
	defer terminal.Unlock()

	terminal.ActiveBuffer().ReverseIndex()
	return nil
}

func saveCursorHandler(pty chan rune, terminal *Terminal) error {
	// Handler should lock the terminal if there will be write operations to any data read by the renderer
	// terminal.Lock()
	// defer terminal.Unlock()

	terminal.ActiveBuffer().SaveCursor()
	return nil
}

func restoreCursorHandler(pty chan rune, terminal *Terminal) error {
	terminal.Lock()
	defer terminal.Unlock()

	terminal.ActiveBuffer().RestoreCursor()
	return nil
}

func ansiHandler(pty chan rune, terminal *Terminal) error {
	// if the byte is an escape character, read the next byte to determine which one
	b := <-pty

	handler, ok := ansiSequenceMap[b]
	if ok {
		//terminal.logger.Debugf("Handling ansi sequence %c", b)
		return handler(pty, terminal)
	}

	return fmt.Errorf("Unknown ANSI control sequence byte: 0x%02X [%v]", b, string(b))
}

func nextLineHandler(pty chan rune, terminal *Terminal) error {
	terminal.Lock()
	defer terminal.Unlock()

	terminal.ActiveBuffer().NewLineEx(true)
	return nil
}

func tabSetHandler(pty chan rune, terminal *Terminal) error {
	// Handler should lock the terminal if there will be write operations to any data read by the renderer
	// terminal.Lock()
	// defer terminal.Unlock()

	terminal.terminalState.TabSetAtCursor()
	return nil
}

func privacyMessageHandler(pty chan rune, terminal *Terminal) error {
	// Handler should lock the terminal if there will be write operations to any data read by the renderer
	// terminal.Lock()
	// defer terminal.Unlock()

	isEscaped := false
	for {
		b := <-pty
		if b == 0x18 /*CAN*/ || b == 0x1a /*SUB*/ || (b == 0x5c /*backslash*/ && isEscaped) {
			break
		}
		if isEscaped {
			isEscaped = false
		} else if b == 0x1b {
			isEscaped = true
			continue
		}
	}
	return nil
}
