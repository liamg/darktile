package terminal

import "fmt"

func screenStateHandler(pty chan rune, terminal *Terminal) error {
	b := <-pty
	switch b {
	case '8': // DECALN -- Screen Alignment Pattern
		// hide cursor?
		// reset margins to extreme positions
		buffer := terminal.ActiveBuffer()
		buffer.SetPosition(0, 0)

		// Fill the whole screen with E's
		count := buffer.ViewHeight() * buffer.ViewWidth()
		for count > 0 {
			buffer.Write('E')
			count--
		}
		// restore cursor?
	default:
		return fmt.Errorf("Screen State code not supported: 0x%02X [%v]", b, string(b))
	}
	return nil
}
