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
			//lineOverflow = false
		}

		switch b {
		case 0x0a:
			terminal.buffer.NewLine()
		case 0x0d:
			terminal.buffer.SetPosition(0, terminal.buffer.CursorLine())
		case 0x08:
			// backspace
			terminal.buffer.MovePosition(-1, 0)
		case 0x07:
			// @todo ring bell
		default:
			// render character at current location
			//		fmt.Printf("%s\n", string([]byte{b}))
			if b >= 0x20 {
				terminal.buffer.Write(b)
			} else {
				terminal.logger.Error("Non-readable rune received: 0x%X", b)
			}
		}

		terminal.triggerOnUpdate()
	}
}
