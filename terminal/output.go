package terminal

import (
	"context"
	"time"
)

// Wish list here: http://invisible-island.net/xterm/ctlseqs/ctlseqs.html

type TerminalCharSet int

type escapeSequenceHandler func(pty chan rune, terminal *Terminal) error

var escapeSequenceMap = map[rune]escapeSequenceHandler{
	0x05: enqSequenceHandler,
	0x07: bellSequenceHandler,
	0x08: backspaceSequenceHandler,
	0x09: tabSequenceHandler,
	0x0a: newLineSequenceHandler,
	0x0b: newLineSequenceHandler,
	0x0c: newLineSequenceHandler,
	0x0d: carriageReturnSequenceHandler,
	0x0e: shiftOutSequenceHandler,
	0x0f: shiftInSequenceHandler,
	0x1b: ansiHandler,
}

func newLineSequenceHandler(pty chan rune, terminal *Terminal) error {
	terminal.ActiveBuffer().NewLine()
	return nil
}

func tabSequenceHandler(pty chan rune, terminal *Terminal) error {
	terminal.ActiveBuffer().Tab()
	return nil
}

func carriageReturnSequenceHandler(pty chan rune, terminal *Terminal) error {
	terminal.ActiveBuffer().CarriageReturn()
	return nil
}

func backspaceSequenceHandler(pty chan rune, terminal *Terminal) error {
	terminal.ActiveBuffer().Backspace()
	return nil
}

func bellSequenceHandler(pty chan rune, terminal *Terminal) error {
	// @todo ring bell - flash red or some shit?
	return nil
}

func enqSequenceHandler(pty chan rune, terminal *Terminal) error {
	terminal.logger.Errorf("Received ENQ!")
	return nil
}

func shiftOutSequenceHandler(pty chan rune, terminal *Terminal) error {
	terminal.logger.Errorf("Received shift out")
	return nil
}

func shiftInSequenceHandler(pty chan rune, terminal *Terminal) error {
	terminal.logger.Errorf("Received shift in")
	return nil
}

func (terminal *Terminal) processInput(ctx context.Context, pty chan rune) {

	// https://en.wikipedia.org/wiki/ANSI_escape_code

	for {

		select {
		case <-terminal.pauseChan:
			// @todo alert user when terminal is suspended
			terminal.logger.Debugf("Terminal suspended")
			<-terminal.resumeChan
		case <-ctx.Done():
			break
		default:
		}

		if terminal.config.Slomo {
			time.Sleep(time.Millisecond * 100)
		}

		b := <-pty

		terminal.logger.Debugf("0x%q", string(b))

		handler, ok := escapeSequenceMap[b]

		if ok {
			//terminal.logger.Debugf("Handling escape sequence: 0x%x", b)
			if err := handler(pty, terminal); err != nil {
				terminal.logger.Errorf("Error handling escape sequence: %s", err)
			}
		} else {
			//terminal.logger.Debugf("Received character 0x%X: %q", b, string(b))
			if b >= 0x20 {
				//terminal.logger.Debugf("%c", b)
				terminal.ActiveBuffer().Write(b)
			} else {
				terminal.logger.Error("Non-readable rune received: 0x%X", b)
			}
		}

		terminal.isDirty = true
	}
}
