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
	'M': reverseIndexHandler,
	'c': swallowHandler(0), //RIS
	'(': swallowHandler(1), // character set bullshit
	')': swallowHandler(1), // character set bullshit
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

func indexHandler(pty chan rune, terminal *Terminal) error {
	// @todo is thus right?
	// "This sequence causes the active position to move downward one line without changing the column position. If the active position is at the bottom margin, a scroll up is performed."
	if terminal.ActiveBuffer().CursorLine() == terminal.ActiveBuffer().ViewHeight()-1 {
		terminal.ActiveBuffer().NewLine()
		return nil
	}
	terminal.ActiveBuffer().MovePosition(0, 1)
	return nil
}

func reverseIndexHandler(pty chan rune, terminal *Terminal) error {
	terminal.ActiveBuffer().MovePosition(0, -1)
	return nil
}

func saveCursorHandler(pty chan rune, terminal *Terminal) error {
	terminal.ActiveBuffer().SaveCursor()
	return nil
}

func restoreCursorHandler(pty chan rune, terminal *Terminal) error {
	terminal.ActiveBuffer().RestoreCursor()
	return nil
}

func ansiHandler(pty chan rune, terminal *Terminal) error {
	// if the byte is an escape character, read the next byte to determine which one
	b := <-pty

	handler, ok := ansiSequenceMap[b]
	if ok {
		return handler(pty, terminal)
	}

	return fmt.Errorf("Unknown ANSI control sequence byte: 0x%02X [%v]", b, string(b))
}
