package terminal

import (
	"fmt"
	"strings"

	"github.com/liamg/aminal/buffer"
)

func oscHandler(pty chan rune, terminal *Terminal) error {

	params := []string{}
	param := ""

	{
		isEscaped := false // flag if the previous character was escape

		for {
			b := <-pty
			if terminal.IsOSCTerminator(b, isEscaped) {
				params = append(params, param)
				break
			}
			if isEscaped {
				isEscaped = false
			} else if b == 0x1b {
				isEscaped = true
				continue
			}
			if b == ';' {
				params = append(params, param)
				param = ""
				continue
			}
			param = fmt.Sprintf("%s%c", param, b)
		}
	}

	if len(params) == 0 {
		return fmt.Errorf("OSC with no params")
	}

	pT := params[len(params)-1]
	pS := params[:len(params)-1]

	if len(pS) == 0 {
		pS = []string{pT}
		pT = ""
	}

	switch pS[0] {
	case "0", "2":
		terminal.SetTitle(pT)
	case "8": // hyperlink
		if len(params) > 2 && len(params[2]) > 0 {
			terminal.terminalState.CurrentHyperlink = &buffer.Hyperlink{Uri: params[2]}
		} else {
			terminal.terminalState.CurrentHyperlink = nil
		}
	case "10": // get/set foreground colour
		if len(pS) > 1 {
			if pS[1] == "?" {
				terminal.Write([]byte("\x1b]10;15"))
			}
		}
	case "11": // get/set background colour
		if len(pS) > 1 {
			if pS[1] == "?" {
				terminal.Write([]byte("\x1b]10;0"))
			}
		}
	default:
		return fmt.Errorf("Unknown OSC control sequence: %s", strings.Join(params, ";"))
	}
	return nil
}
