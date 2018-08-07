package terminal

import (
	"fmt"

	"gitlab.com/liamg/raft/buffer"
)

func sgrSequenceHandler(params []string, intermediate string, terminal *Terminal) error {

	for i := range params {
		param := params[i]
		switch param {
		case "00", "0", "":
			attr := terminal.buffer.CursorAttr()
			*attr = buffer.CellAttributes{
				FgColour: terminal.colourScheme.DefaultFg,
				BgColour: terminal.colourScheme.DefaultBg,
			}
		case "1", "01":
			terminal.buffer.CursorAttr().Bold = true
		case "2", "02":
			terminal.buffer.CursorAttr().Dim = true
		case "4", "04":
			terminal.buffer.CursorAttr().Underline = true
		case "5", "05":
			terminal.buffer.CursorAttr().Blink = true
		case "7", "07":
			terminal.buffer.CursorAttr().Reverse = true
		case "8", "08":
			terminal.buffer.CursorAttr().Hidden = true
		case "21":
			terminal.buffer.CursorAttr().Bold = false
		case "22":
			terminal.buffer.CursorAttr().Dim = false
		case "24":
			terminal.buffer.CursorAttr().Underline = false
		case "25":
			terminal.buffer.CursorAttr().Blink = false
		case "27":
			terminal.buffer.CursorAttr().Reverse = false
		case "28":
			terminal.buffer.CursorAttr().Hidden = false
		case "39":
			terminal.buffer.CursorAttr().FgColour = terminal.colourScheme.DefaultFg
		case "30":
			terminal.buffer.CursorAttr().FgColour = terminal.colourScheme.BlackFg
		case "31":
			terminal.buffer.CursorAttr().FgColour = terminal.colourScheme.RedFg
		case "32":
			terminal.buffer.CursorAttr().FgColour = terminal.colourScheme.GreenFg
		case "33":
			terminal.buffer.CursorAttr().FgColour = terminal.colourScheme.YellowFg
		case "34":
			terminal.buffer.CursorAttr().FgColour = terminal.colourScheme.BlueFg
		case "35":
			terminal.buffer.CursorAttr().FgColour = terminal.colourScheme.MagentaFg
		case "36":
			terminal.buffer.CursorAttr().FgColour = terminal.colourScheme.CyanFg
		case "37":
			terminal.buffer.CursorAttr().FgColour = terminal.colourScheme.WhiteFg
		case "90":
			terminal.buffer.CursorAttr().FgColour = terminal.colourScheme.DarkGreyFg
		case "91":
			terminal.buffer.CursorAttr().FgColour = terminal.colourScheme.LightRedFg
		case "92":
			terminal.buffer.CursorAttr().FgColour = terminal.colourScheme.LightGreenFg
		case "93":
			terminal.buffer.CursorAttr().FgColour = terminal.colourScheme.LightYellowFg
		case "94":
			terminal.buffer.CursorAttr().FgColour = terminal.colourScheme.LightBlueFg
		case "95":
			terminal.buffer.CursorAttr().FgColour = terminal.colourScheme.LightMagentaFg
		case "96":
			terminal.buffer.CursorAttr().FgColour = terminal.colourScheme.LightCyanFg
		case "97":
			terminal.buffer.CursorAttr().FgColour = terminal.colourScheme.WhiteFg
		case "49":
			terminal.buffer.CursorAttr().BgColour = terminal.colourScheme.DefaultBg
		case "40":
			terminal.buffer.CursorAttr().BgColour = terminal.colourScheme.BlackBg
		case "41":
			terminal.buffer.CursorAttr().BgColour = terminal.colourScheme.RedBg
		case "42":
			terminal.buffer.CursorAttr().BgColour = terminal.colourScheme.GreenBg
		case "43":
			terminal.buffer.CursorAttr().BgColour = terminal.colourScheme.YellowBg
		case "44":
			terminal.buffer.CursorAttr().BgColour = terminal.colourScheme.BlueBg
		case "45":
			terminal.buffer.CursorAttr().BgColour = terminal.colourScheme.MagentaBg
		case "46":
			terminal.buffer.CursorAttr().BgColour = terminal.colourScheme.CyanBg
		case "47":
			terminal.buffer.CursorAttr().BgColour = terminal.colourScheme.WhiteBg
		case "100":
			terminal.buffer.CursorAttr().BgColour = terminal.colourScheme.DarkGreyBg
		case "101":
			terminal.buffer.CursorAttr().BgColour = terminal.colourScheme.LightRedBg
		case "102":
			terminal.buffer.CursorAttr().BgColour = terminal.colourScheme.LightGreenBg
		case "103":
			terminal.buffer.CursorAttr().BgColour = terminal.colourScheme.LightYellowBg
		case "104":
			terminal.buffer.CursorAttr().BgColour = terminal.colourScheme.LightBlueBg
		case "105":
			terminal.buffer.CursorAttr().BgColour = terminal.colourScheme.LightMagentaBg
		case "106":
			terminal.buffer.CursorAttr().BgColour = terminal.colourScheme.LightCyanBg
		case "107":
			terminal.buffer.CursorAttr().BgColour = terminal.colourScheme.WhiteBg
		default:
			return fmt.Errorf("Unknown SGR control sequence: (ESC[%s%sm)", param, intermediate)
		}

		//terminal.logger.Debugf("SGR control sequence: (ESC[%s%sm)", param, intermediate)
	}
	return nil
}
