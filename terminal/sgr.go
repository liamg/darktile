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
				FgColour: terminal.config.ColourScheme.DefaultFg,
				BgColour: terminal.config.ColourScheme.DefaultBg,
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
			terminal.buffer.CursorAttr().FgColour = terminal.config.ColourScheme.DefaultFg
		case "30":
			terminal.buffer.CursorAttr().FgColour = terminal.config.ColourScheme.BlackFg
		case "31":
			terminal.buffer.CursorAttr().FgColour = terminal.config.ColourScheme.RedFg
		case "32":
			terminal.buffer.CursorAttr().FgColour = terminal.config.ColourScheme.GreenFg
		case "33":
			terminal.buffer.CursorAttr().FgColour = terminal.config.ColourScheme.YellowFg
		case "34":
			terminal.buffer.CursorAttr().FgColour = terminal.config.ColourScheme.BlueFg
		case "35":
			terminal.buffer.CursorAttr().FgColour = terminal.config.ColourScheme.MagentaFg
		case "36":
			terminal.buffer.CursorAttr().FgColour = terminal.config.ColourScheme.CyanFg
		case "37":
			terminal.buffer.CursorAttr().FgColour = terminal.config.ColourScheme.WhiteFg
		case "90":
			terminal.buffer.CursorAttr().FgColour = terminal.config.ColourScheme.DarkGreyFg
		case "91":
			terminal.buffer.CursorAttr().FgColour = terminal.config.ColourScheme.LightRedFg
		case "92":
			terminal.buffer.CursorAttr().FgColour = terminal.config.ColourScheme.LightGreenFg
		case "93":
			terminal.buffer.CursorAttr().FgColour = terminal.config.ColourScheme.LightYellowFg
		case "94":
			terminal.buffer.CursorAttr().FgColour = terminal.config.ColourScheme.LightBlueFg
		case "95":
			terminal.buffer.CursorAttr().FgColour = terminal.config.ColourScheme.LightMagentaFg
		case "96":
			terminal.buffer.CursorAttr().FgColour = terminal.config.ColourScheme.LightCyanFg
		case "97":
			terminal.buffer.CursorAttr().FgColour = terminal.config.ColourScheme.WhiteFg
		case "49":
			terminal.buffer.CursorAttr().BgColour = terminal.config.ColourScheme.DefaultBg
		case "40":
			terminal.buffer.CursorAttr().BgColour = terminal.config.ColourScheme.BlackBg
		case "41":
			terminal.buffer.CursorAttr().BgColour = terminal.config.ColourScheme.RedBg
		case "42":
			terminal.buffer.CursorAttr().BgColour = terminal.config.ColourScheme.GreenBg
		case "43":
			terminal.buffer.CursorAttr().BgColour = terminal.config.ColourScheme.YellowBg
		case "44":
			terminal.buffer.CursorAttr().BgColour = terminal.config.ColourScheme.BlueBg
		case "45":
			terminal.buffer.CursorAttr().BgColour = terminal.config.ColourScheme.MagentaBg
		case "46":
			terminal.buffer.CursorAttr().BgColour = terminal.config.ColourScheme.CyanBg
		case "47":
			terminal.buffer.CursorAttr().BgColour = terminal.config.ColourScheme.WhiteBg
		case "100":
			terminal.buffer.CursorAttr().BgColour = terminal.config.ColourScheme.DarkGreyBg
		case "101":
			terminal.buffer.CursorAttr().BgColour = terminal.config.ColourScheme.LightRedBg
		case "102":
			terminal.buffer.CursorAttr().BgColour = terminal.config.ColourScheme.LightGreenBg
		case "103":
			terminal.buffer.CursorAttr().BgColour = terminal.config.ColourScheme.LightYellowBg
		case "104":
			terminal.buffer.CursorAttr().BgColour = terminal.config.ColourScheme.LightBlueBg
		case "105":
			terminal.buffer.CursorAttr().BgColour = terminal.config.ColourScheme.LightMagentaBg
		case "106":
			terminal.buffer.CursorAttr().BgColour = terminal.config.ColourScheme.LightCyanBg
		case "107":
			terminal.buffer.CursorAttr().BgColour = terminal.config.ColourScheme.WhiteBg
		default:
			return fmt.Errorf("Unknown SGR control sequence: (ESC[%s%sm)", param, intermediate)
		}

		//terminal.logger.Debugf("SGR control sequence: (ESC[%s%sm)", param, intermediate)
	}
	return nil
}
