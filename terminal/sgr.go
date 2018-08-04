package terminal

import "fmt"

func sgrSequenceHandler(params []string, intermediate string, terminal *Terminal) error {

	for i := range params {
		param := params[i]
		switch param {
		case "00", "0", "":
			terminal.cellAttr = terminal.defaultCellAttr
		case "1", "01":
			terminal.cellAttr.Bold = true
		case "2", "02":
			terminal.cellAttr.Dim = true
		case "4", "04":
			terminal.cellAttr.Underline = true
		case "5", "05":
			terminal.cellAttr.Blink = true
		case "7", "07":
			terminal.cellAttr.Reverse = true
		case "8", "08":
			terminal.cellAttr.Hidden = true
		case "21":
			terminal.cellAttr.Bold = false
		case "22":
			terminal.cellAttr.Dim = false
		case "24":
			terminal.cellAttr.Underline = false
		case "25":
			terminal.cellAttr.Blink = false
		case "27":
			terminal.cellAttr.Reverse = false
		case "28":
			terminal.cellAttr.Hidden = false
		case "39":
			terminal.cellAttr.FgColour = terminal.colourScheme.DefaultFg
		case "30":
			terminal.cellAttr.FgColour = terminal.colourScheme.BlackFg
		case "31":
			terminal.cellAttr.FgColour = terminal.colourScheme.RedFg
		case "32":
			terminal.cellAttr.FgColour = terminal.colourScheme.GreenFg
		case "33":
			terminal.cellAttr.FgColour = terminal.colourScheme.YellowFg
		case "34":
			terminal.cellAttr.FgColour = terminal.colourScheme.BlueFg
		case "35":
			terminal.cellAttr.FgColour = terminal.colourScheme.MagentaFg
		case "36":
			terminal.cellAttr.FgColour = terminal.colourScheme.CyanFg
		case "37":
			terminal.cellAttr.FgColour = terminal.colourScheme.WhiteFg
		case "90":
			terminal.cellAttr.FgColour = terminal.colourScheme.DarkGreyFg
		case "91":
			terminal.cellAttr.FgColour = terminal.colourScheme.LightRedFg
		case "92":
			terminal.cellAttr.FgColour = terminal.colourScheme.LightGreenFg
		case "93":
			terminal.cellAttr.FgColour = terminal.colourScheme.LightYellowFg
		case "94":
			terminal.cellAttr.FgColour = terminal.colourScheme.LightBlueFg
		case "95":
			terminal.cellAttr.FgColour = terminal.colourScheme.LightMagentaFg
		case "96":
			terminal.cellAttr.FgColour = terminal.colourScheme.LightCyanFg
		case "97":
			terminal.cellAttr.FgColour = terminal.colourScheme.WhiteFg
		case "49":
			terminal.cellAttr.BgColour = terminal.colourScheme.DefaultBg
		case "40":
			terminal.cellAttr.BgColour = terminal.colourScheme.BlackBg
		case "41":
			terminal.cellAttr.BgColour = terminal.colourScheme.RedBg
		case "42":
			terminal.cellAttr.BgColour = terminal.colourScheme.GreenBg
		case "43":
			terminal.cellAttr.BgColour = terminal.colourScheme.YellowBg
		case "44":
			terminal.cellAttr.BgColour = terminal.colourScheme.BlueBg
		case "45":
			terminal.cellAttr.BgColour = terminal.colourScheme.MagentaBg
		case "46":
			terminal.cellAttr.BgColour = terminal.colourScheme.CyanBg
		case "47":
			terminal.cellAttr.BgColour = terminal.colourScheme.WhiteBg
		case "100":
			terminal.cellAttr.BgColour = terminal.colourScheme.DarkGreyBg
		case "101":
			terminal.cellAttr.BgColour = terminal.colourScheme.LightRedBg
		case "102":
			terminal.cellAttr.BgColour = terminal.colourScheme.LightGreenBg
		case "103":
			terminal.cellAttr.BgColour = terminal.colourScheme.LightYellowBg
		case "104":
			terminal.cellAttr.BgColour = terminal.colourScheme.LightBlueBg
		case "105":
			terminal.cellAttr.BgColour = terminal.colourScheme.LightMagentaBg
		case "106":
			terminal.cellAttr.BgColour = terminal.colourScheme.LightCyanBg
		case "107":
			terminal.cellAttr.BgColour = terminal.colourScheme.WhiteBg
		default:
			return fmt.Errorf("Unknown SGR control sequence: (ESC[%s%sm)", param, intermediate)
		}

		//terminal.logger.Debugf("SGR control sequence: (ESC[%s%sm)", param, intermediate)
	}
	return nil
}
