package terminal

import (
	"context"
	"time"
)

// Wish list here: http://invisible-island.net/xterm/ctlseqs/ctlseqs.html

type TerminalCharSet int

type escapeSequenceHandler func(pty chan rune, terminal *Terminal) error

var escapeSequenceMap = map[rune]escapeSequenceHandler{
	0x1b: ansiHandler,
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

		handler, ok := escapeSequenceMap[b]

		if ok {
			if err := handler(pty, terminal); err != nil {
				terminal.logger.Errorf("Error handling escape sequence: %s", err)
			}
			continue
		}

		terminal.logger.Debugf("Received character 0x%X: %q", b, string(b))

		switch b {
		case 0x0a, 0x0c, 0x0b: // LF, FF, VT
			terminal.ActiveBuffer().NewLine()
		case 0x0d: // CR
			terminal.ActiveBuffer().CarriageReturn()
		case 0x08: // BS
			// backspace
			terminal.ActiveBuffer().Backspace()
		case 0x07: // BEL
			// @todo ring bell - flash red or some shit?
		case 0x05: // ENQ
			terminal.logger.Errorf("Received ENQ!")
		case 0xe, 0xf:
			terminal.logger.Errorf("Received SI/SO")
		case 0x09:
			terminal.logger.Errorf("Received TAB")
		default:
			if b >= 0x20 {
				terminal.ActiveBuffer().Write(b)
			} else {
				terminal.logger.Error("Non-readable rune received: 0x%X", b)
			}
		}

	}
}
