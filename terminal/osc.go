package terminal

import "fmt"

func oscHandler(buffer chan rune, terminal *Terminal) error {
	b := <-buffer
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
			terminal.SetTitle(string(title))
		} else {
			return fmt.Errorf("Invalid OSC 0 control sequence: 0x%02X", b)
		}
	default:
		return fmt.Errorf("Unknown OSC control sequence: 0x%02X", b)
	}
	return nil
}
