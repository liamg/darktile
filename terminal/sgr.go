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
				FgColour: terminal.config.ColourScheme.Foreground,
				BgColour: terminal.config.ColourScheme.Background,
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
			terminal.buffer.CursorAttr().FgColour = terminal.config.ColourScheme.Foreground
		case "30":
			terminal.buffer.CursorAttr().FgColour = terminal.config.ColourScheme.Black
		case "31":
			terminal.buffer.CursorAttr().FgColour = terminal.config.ColourScheme.Red
		case "32":
			terminal.buffer.CursorAttr().FgColour = terminal.config.ColourScheme.Green
		case "33":
			terminal.buffer.CursorAttr().FgColour = terminal.config.ColourScheme.Yellow
		case "34":
			terminal.buffer.CursorAttr().FgColour = terminal.config.ColourScheme.Blue
		case "35":
			terminal.buffer.CursorAttr().FgColour = terminal.config.ColourScheme.Magenta
		case "36":
			terminal.buffer.CursorAttr().FgColour = terminal.config.ColourScheme.Cyan
		case "37":
			terminal.buffer.CursorAttr().FgColour = terminal.config.ColourScheme.White
		case "90":
			terminal.buffer.CursorAttr().FgColour = terminal.config.ColourScheme.DarkGrey
		case "91":
			terminal.buffer.CursorAttr().FgColour = terminal.config.ColourScheme.LightRed
		case "92":
			terminal.buffer.CursorAttr().FgColour = terminal.config.ColourScheme.LightGreen
		case "93":
			terminal.buffer.CursorAttr().FgColour = terminal.config.ColourScheme.LightYellow
		case "94":
			terminal.buffer.CursorAttr().FgColour = terminal.config.ColourScheme.LightBlue
		case "95":
			terminal.buffer.CursorAttr().FgColour = terminal.config.ColourScheme.LightMagenta
		case "96":
			terminal.buffer.CursorAttr().FgColour = terminal.config.ColourScheme.LightCyan
		case "97":
			terminal.buffer.CursorAttr().FgColour = terminal.config.ColourScheme.White
		case "49":
			terminal.buffer.CursorAttr().BgColour = terminal.config.ColourScheme.Background
		case "40":
			terminal.buffer.CursorAttr().BgColour = terminal.config.ColourScheme.Black
		case "41":
			terminal.buffer.CursorAttr().BgColour = terminal.config.ColourScheme.Red
		case "42":
			terminal.buffer.CursorAttr().BgColour = terminal.config.ColourScheme.Green
		case "43":
			terminal.buffer.CursorAttr().BgColour = terminal.config.ColourScheme.Yellow
		case "44":
			terminal.buffer.CursorAttr().BgColour = terminal.config.ColourScheme.Blue
		case "45":
			terminal.buffer.CursorAttr().BgColour = terminal.config.ColourScheme.Magenta
		case "46":
			terminal.buffer.CursorAttr().BgColour = terminal.config.ColourScheme.Cyan
		case "47":
			terminal.buffer.CursorAttr().BgColour = terminal.config.ColourScheme.White
		case "100":
			terminal.buffer.CursorAttr().BgColour = terminal.config.ColourScheme.DarkGrey
		case "101":
			terminal.buffer.CursorAttr().BgColour = terminal.config.ColourScheme.LightRed
		case "102":
			terminal.buffer.CursorAttr().BgColour = terminal.config.ColourScheme.LightGreen
		case "103":
			terminal.buffer.CursorAttr().BgColour = terminal.config.ColourScheme.LightYellow
		case "104":
			terminal.buffer.CursorAttr().BgColour = terminal.config.ColourScheme.LightBlue
		case "105":
			terminal.buffer.CursorAttr().BgColour = terminal.config.ColourScheme.LightMagenta
		case "106":
			terminal.buffer.CursorAttr().BgColour = terminal.config.ColourScheme.LightCyan
		case "107":
			terminal.buffer.CursorAttr().BgColour = terminal.config.ColourScheme.White
		default:
			return fmt.Errorf("Unknown SGR control sequence: (ESC[%s%sm)", param, intermediate)
		}

		//terminal.logger.Debugf("SGR control sequence: (ESC[%s%sm)", param, intermediate)
	}
	return nil
}
