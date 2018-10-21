package terminal

import (
	"fmt"

	"github.com/liamg/aminal/buffer"
)

func sgrSequenceHandler(params []string, intermediate string, terminal *Terminal) error {

	for i := range params {
		param := params[i]
		switch param {
		case "00", "0", "":
			attr := terminal.ActiveBuffer().CursorAttr()
			*attr = buffer.CellAttributes{
				FgColour: terminal.config.ColourScheme.Foreground,
				BgColour: terminal.config.ColourScheme.Background,
			}
		case "1", "01":
			terminal.ActiveBuffer().CursorAttr().Bold = true
		case "2", "02":
			terminal.ActiveBuffer().CursorAttr().Dim = true
		case "4", "04":
			terminal.ActiveBuffer().CursorAttr().Underline = true
		case "5", "05":
			terminal.ActiveBuffer().CursorAttr().Blink = true
		case "7", "07":
			terminal.ActiveBuffer().CursorAttr().Reverse = true
		case "8", "08":
			terminal.ActiveBuffer().CursorAttr().Hidden = true
		case "21":
			terminal.ActiveBuffer().CursorAttr().Bold = false
		case "22":
			terminal.ActiveBuffer().CursorAttr().Dim = false
		case "24":
			terminal.ActiveBuffer().CursorAttr().Underline = false
		case "25":
			terminal.ActiveBuffer().CursorAttr().Blink = false
		case "27":
			terminal.ActiveBuffer().CursorAttr().Reverse = false
		case "28":
			terminal.ActiveBuffer().CursorAttr().Hidden = false
		case "39":
			terminal.ActiveBuffer().CursorAttr().FgColour = terminal.config.ColourScheme.Foreground
		case "30":
			terminal.ActiveBuffer().CursorAttr().FgColour = terminal.config.ColourScheme.Black
		case "31":
			terminal.ActiveBuffer().CursorAttr().FgColour = terminal.config.ColourScheme.Red
		case "32":
			terminal.ActiveBuffer().CursorAttr().FgColour = terminal.config.ColourScheme.Green
		case "33":
			terminal.ActiveBuffer().CursorAttr().FgColour = terminal.config.ColourScheme.Yellow
		case "34":
			terminal.ActiveBuffer().CursorAttr().FgColour = terminal.config.ColourScheme.Blue
		case "35":
			terminal.ActiveBuffer().CursorAttr().FgColour = terminal.config.ColourScheme.Magenta
		case "36":
			terminal.ActiveBuffer().CursorAttr().FgColour = terminal.config.ColourScheme.Cyan
		case "37":
			terminal.ActiveBuffer().CursorAttr().FgColour = terminal.config.ColourScheme.White
		case "90":
			terminal.ActiveBuffer().CursorAttr().FgColour = terminal.config.ColourScheme.DarkGrey
		case "91":
			terminal.ActiveBuffer().CursorAttr().FgColour = terminal.config.ColourScheme.LightRed
		case "92":
			terminal.ActiveBuffer().CursorAttr().FgColour = terminal.config.ColourScheme.LightGreen
		case "93":
			terminal.ActiveBuffer().CursorAttr().FgColour = terminal.config.ColourScheme.LightYellow
		case "94":
			terminal.ActiveBuffer().CursorAttr().FgColour = terminal.config.ColourScheme.LightBlue
		case "95":
			terminal.ActiveBuffer().CursorAttr().FgColour = terminal.config.ColourScheme.LightMagenta
		case "96":
			terminal.ActiveBuffer().CursorAttr().FgColour = terminal.config.ColourScheme.LightCyan
		case "97":
			terminal.ActiveBuffer().CursorAttr().FgColour = terminal.config.ColourScheme.White
		case "49":
			terminal.ActiveBuffer().CursorAttr().BgColour = terminal.config.ColourScheme.Background
		case "40":
			terminal.ActiveBuffer().CursorAttr().BgColour = terminal.config.ColourScheme.Black
		case "41":
			terminal.ActiveBuffer().CursorAttr().BgColour = terminal.config.ColourScheme.Red
		case "42":
			terminal.ActiveBuffer().CursorAttr().BgColour = terminal.config.ColourScheme.Green
		case "43":
			terminal.ActiveBuffer().CursorAttr().BgColour = terminal.config.ColourScheme.Yellow
		case "44":
			terminal.ActiveBuffer().CursorAttr().BgColour = terminal.config.ColourScheme.Blue
		case "45":
			terminal.ActiveBuffer().CursorAttr().BgColour = terminal.config.ColourScheme.Magenta
		case "46":
			terminal.ActiveBuffer().CursorAttr().BgColour = terminal.config.ColourScheme.Cyan
		case "47":
			terminal.ActiveBuffer().CursorAttr().BgColour = terminal.config.ColourScheme.White
		case "100":
			terminal.ActiveBuffer().CursorAttr().BgColour = terminal.config.ColourScheme.DarkGrey
		case "101":
			terminal.ActiveBuffer().CursorAttr().BgColour = terminal.config.ColourScheme.LightRed
		case "102":
			terminal.ActiveBuffer().CursorAttr().BgColour = terminal.config.ColourScheme.LightGreen
		case "103":
			terminal.ActiveBuffer().CursorAttr().BgColour = terminal.config.ColourScheme.LightYellow
		case "104":
			terminal.ActiveBuffer().CursorAttr().BgColour = terminal.config.ColourScheme.LightBlue
		case "105":
			terminal.ActiveBuffer().CursorAttr().BgColour = terminal.config.ColourScheme.LightMagenta
		case "106":
			terminal.ActiveBuffer().CursorAttr().BgColour = terminal.config.ColourScheme.LightCyan
		case "107":
			terminal.ActiveBuffer().CursorAttr().BgColour = terminal.config.ColourScheme.White
		default:
			return fmt.Errorf("Unknown SGR control sequence: (ESC[%s%sm)", param, intermediate)
		}

		//terminal.logger.Debugf("SGR control sequence: (ESC[%s%sm)", param, intermediate)
	}
	return nil
}
