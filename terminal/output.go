package terminal

import (
	"context"
)

// Wish list here: http://invisible-island.net/xterm/ctlseqs/ctlseqs.html

type TerminalCharSet int

type escapeSequenceHandler func(pty chan rune, terminal *Terminal) error

var escapeSequenceMap = map[rune]escapeSequenceHandler{
	0x1b: ansiHandler,
}

func (terminal *Terminal) Suspend() {
	select {
	case terminal.pauseChan <- true:
	default:
	}
}

func (terminal *Terminal) Resume() {
	select {
	case terminal.resumeChan <- true:
	default:
	}
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

		//if terminal.config.slomo
		//time.Sleep(time.Millisecond * 100)

		b := <-pty

		handler, ok := escapeSequenceMap[b]

		if ok {
			if err := handler(pty, terminal); err != nil {
				terminal.logger.Errorf("Error handling escape sequence 0x%X: %s", b, err)
			}
			continue
		}

		terminal.logger.Debugf("Received character 0x%X: %q", b, string(b))

		switch b {
		case 0x0a:
			terminal.ActiveBuffer().NewLine()
		case 0x0d:
			terminal.ActiveBuffer().CarriageReturn()
		case 0x08:
			// backspace
			terminal.ActiveBuffer().Backspace()
		case 0x07:
			// @todo ring bell - flash red or some shit?
		default:
			// render character at current location
			//		fmt.Printf("%s\n", string([]byte{b}))
			if b >= 0x20 {
				terminal.ActiveBuffer().Write(b)
			} else {
				terminal.logger.Error("Non-readable rune received: 0x%X", b)
			}
		}

	}
}
