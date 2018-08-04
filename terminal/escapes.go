package terminal

import (
	"context"
	"fmt"
)

// Wish list here: http://invisible-island.net/xterm/ctlseqs/ctlseqs.html

type TerminalCharSet int

type escapeSequenceHandler func(buffer chan rune, terminal *Terminal) error

var escapeSequenceMap = map[rune]escapeSequenceHandler{
	0x1b: ansiHandler,
}

func (terminal *Terminal) processInput(ctx context.Context, buffer chan rune) {

	// https://en.wikipedia.org/wiki/ANSI_escape_code

	lineOverflow := false

	for {

		select {
		case <-ctx.Done():
			break
		default:
		}

		b := <-buffer

		handler, ok := escapeSequenceMap[b]

		if ok {
			if err := handler(buffer, terminal); err != nil {
				fmt.Errorf("Error handling escape sequence 0x%X: %s", b, err)
			}
			continue
		}

		if b != 0x0d {
			lineOverflow = false
		}

		switch b {
		case 0x0a:

			_, h := terminal.GetSize()
			if terminal.position.Line+1 >= h {
				terminal.lines = append(terminal.lines, NewLine())
			} else {
				terminal.position.Line++
			}

		case 0x0d:
			if terminal.position.Col == 0 && terminal.position.Line > 0 && lineOverflow {
				terminal.position.Line--
				terminal.logger.Debugf("Swallowing forced new line for CR")
				lineOverflow = false
			}
			terminal.position.Col = 0

		case 0x08:
			// backspace
			terminal.position.Col--
			if terminal.position.Col < 0 {
				terminal.position.Col = 0
			}
		case 0x07:
			// @todo ring bell
		default:
			// render character at current location
			//		fmt.Printf("%s\n", string([]byte{b}))
			if b >= 0x20 {
				terminal.writeRune(b)
				lineOverflow = terminal.position.Col == 0
			} else {
				terminal.logger.Error("Non-readable rune received: 0x%X", b)
			}
		}

		terminal.triggerOnUpdate()
	}
}
