package terminal

import (
	"fmt"
	"strings"
)

func oscHandler(pty chan rune, terminal *Terminal) error {

	params := []string{}
	param := ""

	for {
		b := <-pty
		if b == 0x07 || b == 0x5c {
			params = append(params, param)
			break
		}
		if b == ';' {
			params = append(params, param)
			param = ""
			continue
		}
		param = fmt.Sprintf("%s%c", param, b)
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
	default:
		return fmt.Errorf("Unknown OSC control sequence: %s", strings.Join(params, ";"))
	}
	return nil
}
