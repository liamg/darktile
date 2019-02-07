package terminal

import "fmt"

func screenStateHandler(pty chan rune, terminal *Terminal) error {
	b := <-pty
	switch b {
	case '8': // DECALN -- Screen Alignment Pattern
		// hide cursor?
		buffer := terminal.ActiveBuffer()
		terminal.ResetVerticalMargins()
		terminal.ScrollToEnd()

		// Fill the whole screen with E's
		count := buffer.ViewHeight() * buffer.ViewWidth()
		for count > 0 {
			buffer.Write('E')
			count--
			if count > 0 && !terminal.IsAutoWrap() && count%buffer.ViewWidth() == 0 {
				buffer.Index()
				buffer.CarriageReturn()
			}
		}
		// restore cursor
		buffer.SetPosition(0, 0)
	default:
		return fmt.Errorf("Screen State code not supported: 0x%02X [%v]", b, string(b))
	}
	return nil
}
